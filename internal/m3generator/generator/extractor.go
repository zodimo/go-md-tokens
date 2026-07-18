package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bep/godartsass/v2"
)

type Schema struct {
	Components    []ComponentSchema            `json:"components"`
	SystemModules map[string]map[string]string `json:"systemModules"`
}

type ComponentSchema struct {
	Name              string            `json:"name"`
	SupportedTokens   []string          `json:"supportedTokens"`
	UnsupportedTokens []string          `json:"unsupportedTokens"`
	RenamedTokens     map[string]string `json:"renamedTokens"`
	Dependencies      []string          `json:"dependencies"`
	StatePrefixes     []string          `json:"statePrefixes"`
}

var renamedRegex = regexp.MustCompile(`\$renamed-tokens:\s*\(([^)]+)\)`)
var keyValueRegex = regexp.MustCompile(`'([^']+)'\s*:\s*(?:'([^']+)'|[^,\n]+)`)
var cssVarRegex = regexp.MustCompile(`--([^:]+):\s*([^;]+);`)

func BuildSchema(transpiler *godartsass.Transpiler, tokensDir string) (*Schema, error) {
	schema := &Schema{
		Components:    []ComponentSchema{},
		SystemModules: make(map[string]map[string]string),
	}

	files, err := os.ReadDir(tokensDir)
	if err != nil {
		return nil, err
	}

	// 1. Parse System Modules
	sysModules := []string{"md-sys-color", "md-sys-elevation", "md-sys-shape", "md-sys-state", "md-sys-typescale"}
	for _, sys := range sysModules {
		sysMap, err := extractSystemModule(transpiler, tokensDir, sys)
		if err != nil {
			return nil, err
		}
		schema.SystemModules[sys] = sysMap
	}

	// 2. Parse Components
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "_md-comp-") && strings.HasSuffix(file.Name(), ".scss") {
			compSchema, err := extractComponentSchema(transpiler, filepath.Join(tokensDir, file.Name()))
			if err != nil {
				return nil, err
			}
			if compSchema != nil {
				schema.Components = append(schema.Components, *compSchema)
			}
		}
	}

	return schema, nil
}

func extractSystemModule(transpiler *godartsass.Transpiler, dir, sysName string) (map[string]string, error) {
	valFunc := "values()"
	if sysName == "md-sys-color" {
		valFunc = "values-light()"
	}
	bridgeCode := fmt.Sprintf(`
@use '_%s' as sys;
.dump {
  @each $k, $v in sys.%s {
    --sys-#{$k}: #{$v};
  }
}
`, sysName, valFunc)

	res, err := transpiler.Execute(godartsass.Args{
		Source:       bridgeCode,
		IncludePaths: []string{dir},
	})
	if err != nil {
		return nil, fmt.Errorf("sass error for %s: %v", sysName, err)
	}

	sysMap := make(map[string]string)
	matches := cssVarRegex.FindAllStringSubmatch(res.CSS, -1)
	for _, m := range matches {
		key := strings.TrimSpace(m[1])
		val := strings.TrimSpace(m[2])
		if strings.HasPrefix(key, "sys-") {
			sysMap[strings.TrimPrefix(key, "sys-")] = val
		}
	}
	return sysMap, nil
}

func extractComponentSchema(transpiler *godartsass.Transpiler, path string) (*ComponentSchema, error) {
	filename := filepath.Base(path)
	name := strings.TrimPrefix(filename, "_md-comp-")
	name = strings.TrimSuffix(name, ".scss")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)

	// Check if $supported-tokens exists
	if !strings.Contains(content, "$supported-tokens:") {
		return nil, nil // Some files are helpers, skip if no supported tokens
	}

	bridgeCode := fmt.Sprintf(`
@use 'sass:list';
@use 'sass:map';
@use 'sass:meta';
@use 'md-comp-%s' as comp;

.dump {
  @each $t in comp.$supported-tokens {
    --supported-#{$t}: true;
  }
  @if map.has-key(meta.module-variables("comp"), "unsupported-tokens") {
    @each $t in comp.$unsupported-tokens {
      --unsupported-#{$t}: true;
    }
  }
}
`, name)

	res, err := transpiler.Execute(godartsass.Args{
		Source:       bridgeCode,
		IncludePaths: []string{filepath.Dir(path)},
	})
	if err != nil {
		return nil, fmt.Errorf("sass error for %s: %v", name, err)
	}

	compSchema := &ComponentSchema{
		Name:              name,
		SupportedTokens:   []string{},
		UnsupportedTokens: []string{},
		RenamedTokens:     make(map[string]string),
		Dependencies:      []string{},
		StatePrefixes:     []string{},
	}

	matches := cssVarRegex.FindAllStringSubmatch(res.CSS, -1)
	for _, m := range matches {
		key := strings.TrimSpace(m[1])
		if strings.HasPrefix(key, "supported-") {
			compSchema.SupportedTokens = append(compSchema.SupportedTokens, strings.TrimPrefix(key, "supported-"))
		} else if strings.HasPrefix(key, "unsupported-") {
			compSchema.UnsupportedTokens = append(compSchema.UnsupportedTokens, strings.TrimPrefix(key, "unsupported-"))
		}
	}

	// Regex for dependencies (from $_default to the next ;)
	if idx := strings.Index(content, "$_default:"); idx != -1 {
		endIdx := strings.Index(content[idx:], ";")
		if endIdx != -1 {
			block := content[idx : idx+endIdx]
			var stringRegex = regexp.MustCompile(`'([^']+)'`)
			for _, m := range stringRegex.FindAllStringSubmatch(block, -1) {
				compSchema.Dependencies = append(compSchema.Dependencies, m[1])
			}
		}
	}

	// Regex for renamed-tokens
	if m := renamedRegex.FindStringSubmatch(content); len(m) > 1 {
		for _, kv := range keyValueRegex.FindAllStringSubmatch(m[1], -1) {
			if len(kv) >= 3 {
				compSchema.RenamedTokens[kv[1]] = kv[2]
			}
		}
	}

	// Identify state prefixes (compound greedily)
	stateVocab := []string{"hover-", "focus-", "pressed-", "disabled-", "dragged-", "selected-", "error-", "active-", "toggle-"}
	prefixMap := make(map[string]bool)
	for _, t := range compSchema.SupportedTokens {
		prefix := ""
		for {
			matched := false
			for _, state := range stateVocab {
				if strings.HasPrefix(t, state) {
					prefix += state
					t = strings.TrimPrefix(t, state)
					matched = true
					break
				}
			}
			if !matched {
				break
			}
		}
		if prefix != "" {
			prefixMap[prefix] = true
		}
	}
	for prefix := range prefixMap {
		compSchema.StatePrefixes = append(compSchema.StatePrefixes, prefix)
	}

	sort.Strings(compSchema.SupportedTokens)
	sort.Strings(compSchema.UnsupportedTokens)
	sort.Strings(compSchema.Dependencies)
	sort.Strings(compSchema.StatePrefixes)

	return compSchema, nil
}
