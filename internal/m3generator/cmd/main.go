package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bep/godartsass/v2"
	"github.com/xeipuuv/gojsonschema"
	"github.com/zodimo/go-md-tokens/internal/m3generator/generator"
	"github.com/zodimo/go-md-tokens/internal/m3generator/schemadef"
	"github.com/zodimo/go-md-tokens/internal/m3generator/templates"
	"gopkg.in/yaml.v3"
)

type ThemeFileContext struct {
	Modes []ThemeModeContext
}

type ThemeModeContext struct {
	FuncName string
	ModeName string
	Tokens   map[string]string // Used for deterministic iteration during Go code rendering
}

type SassBridgeContext struct {
	FileURL          string
	CustomColorsSass string
	Component        string
}

type Config struct {
	Pipeline struct {
		ThemeSource   string `yaml:"theme_source"`
		SassTokensDir string `yaml:"sass_tokens_dir"`
	} `yaml:"pipeline"`
	Generator struct {
		SchemaOutputDir  string `yaml:"schema_output_dir"`
		SchemaOutputFile string `yaml:"schema_output_file"`
	} `yaml:"generator"`
	ThemeOutput struct {
		OutputDir        string            `yaml:"output_dir"`
		OutputFile       string            `yaml:"output_file"`
		Components       []string          `yaml:"components"`
		Targets          []string          `yaml:"targets"`
		StructuralTokens map[string]string `yaml:"structural_tokens"`
	} `yaml:"theme_output"`
}

var (
	cssPropRegex     = regexp.MustCompile(`([^:\s]+)\s*:\s*([^;\n]+);`)
	aliasRegex       = regexp.MustCompile(`\{([^}]+)\}`)
	varFallbackRegex = regexp.MustCompile(`var\(--md-[^,)]+,\s*([^)]+)\)`)
)

func main() {
	configPath := flag.String("config", "m3gen.yaml", "Path to configuration file")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Could not parse config file: %v", err)
	}

	if _, err := os.Stat(cfg.Pipeline.ThemeSource); os.IsNotExist(err) {
		log.Fatalf("Could not find theme file at: %s", cfg.Pipeline.ThemeSource)
	}

	themeDataRaw, err := os.ReadFile(cfg.Pipeline.ThemeSource)
	if err != nil {
		log.Fatal(err)
	}

	var themeData struct {
		Schemes map[string]map[string]string `json:"schemes"`
	}
	if err := json.Unmarshal(themeDataRaw, &themeData); err != nil {
		log.Fatal(err)
	}

	transpiler, err := godartsass.Start(godartsass.Options{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Extracting token schema from SCSS...")
	schema, err := generator.BuildSchema(transpiler, cfg.Pipeline.SassTokensDir)
	if err != nil {
		log.Fatalf("Failed to extract schema: %v", err)
	}

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize schema: %v", err)
	}

	if err := os.MkdirAll(cfg.Generator.SchemaOutputDir, 0755); err != nil {
		log.Fatal(err)
	}
	schemaPath := filepath.Join(cfg.Generator.SchemaOutputDir, cfg.Generator.SchemaOutputFile)
	if err := os.WriteFile(schemaPath, schemaJSON, 0644); err != nil {
		log.Fatalf("Failed to write schema: %v", err)
	}

	absSchemaDefPath := schemadef.Path("def.json")
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + absSchemaDefPath)
	documentLoader := gojsonschema.NewStringLoader(string(schemaJSON))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Fatalf("Schema validation error: %v", err)
	}
	if !result.Valid() {
		for _, desc := range result.Errors() {
			log.Printf("- %s\n", desc)
		}
		log.Fatal("Extracted schema does not match schema-def.json")
	}
	log.Println("Schema extracted and validated successfully!")

	log.Println("Generating typed token system...")
	if err := generator.Generate(schema, cfg.Generator.SchemaOutputDir); err != nil {
		log.Fatalf("Code generation failed: %v", err)
	}
	log.Println("Typed token system generated successfully!")

	absDir, err := filepath.Abs(cfg.Pipeline.SassTokensDir)
	if err != nil {
		log.Fatal(err)
	}

	var fileContext ThemeFileContext

	for _, targetMode := range cfg.ThemeOutput.Targets {
		log.Printf("Processing theme mode: [%s]", targetMode)

		allTokens := make(map[string]string, len(cfg.ThemeOutput.StructuralTokens))
		for k, v := range cfg.ThemeOutput.StructuralTokens {
			allTokens[k] = v
		}

		schemeColors, ok := themeData.Schemes[targetMode]
		if !ok {
			continue
		}

		for colorKey, hexVal := range schemeColors {
			kebabKey := camelToKebab(colorKey)
			allTokens[fmt.Sprintf("md.sys.color.%s", kebabKey)] = hexVal
		}

		// Pre-populate all supported tokens with empty strings
		// This ensures that even if SASS drops a token (e.g., null value), it exists in the map
		for _, compSchema := range schema.Components {
			for _, tok := range compSchema.SupportedTokens {
				tokenKey := normalize(fmt.Sprintf("md.comp.%s.%s", compSchema.Name, tok))
				allTokens[tokenKey] = ""
			}
		}

		customColorsSass := buildSassColorMap(schemeColors)

		for _, comp := range cfg.ThemeOutput.Components {
			filename := fmt.Sprintf("_md-comp-%s.scss", comp)
			localFilePath := filepath.Join(cfg.Pipeline.SassTokensDir, filename)

			if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
				continue
			}

			absFilePath, err := filepath.Abs(localFilePath)
			if err != nil {
				log.Printf("Error resolving path for %s: %v", filename, err)
				continue
			}

			fileUrl := (&url.URL{Scheme: "file", Path: absFilePath}).String()

			sassCtx := SassBridgeContext{
				FileURL:          fileUrl,
				CustomColorsSass: customColorsSass,
				Component:        comp,
			}

			var sassBridgeBuf bytes.Buffer
			if err := templates.ParsedTemplates.ExecuteTemplate(&sassBridgeBuf, "sass_bridge.tmpl", sassCtx); err != nil {
				log.Printf("Failed to execute sass bridge template for %s: %v", filename, err)
				continue
			}

			res, err := transpiler.Execute(godartsass.Args{
				Source:       sassBridgeBuf.String(),
				URL:          (&url.URL{Scheme: "file", Path: filepath.Join(absDir, "gen.css.scss")}).String(),
				IncludePaths: []string{absDir},
			})
			if err != nil {
				log.Printf("Sass runtime exception for %s: %v", filename, err)
				continue
			}

			matches := cssPropRegex.FindAllStringSubmatch(res.CSS, -1)

			for _, match := range matches {
				property := strings.TrimPrefix(strings.TrimSpace(match[1]), "--")
				value := strings.TrimSpace(match[2])
				value = varFallbackRegex.ReplaceAllString(value, "$1")

				if strings.Contains(value, "$") || strings.Contains(value, "values-light") {
					continue
				}

				tokenKey := normalize(fmt.Sprintf("md.comp.%s.%s", comp, property))
				allTokens[tokenKey] = normalize(value)
			}
		}

		resolvedTokens := resolveAliases(allTokens)

		modeCtx := ThemeModeContext{
			FuncName: toPascalCase(targetMode),
			ModeName: targetMode,
			Tokens:   resolvedTokens, // Unordered map is fine here, we will sort keys in the template if needed, or we sort here to pass an ordered slice.
			// Wait, text/template 'range' over maps sorts by key natively!
			// "The default behavior is to iterate over map keys in sorted order." (since Go 1.12)
			// So we can just pass the map directly.
		}

		fileContext.Modes = append(fileContext.Modes, modeCtx)
	}

	if err := os.MkdirAll(cfg.ThemeOutput.OutputDir, 0755); err != nil {
		log.Fatal(err)
	}

	var goCodeBuf bytes.Buffer
	if err := templates.ParsedTemplates.ExecuteTemplate(&goCodeBuf, "go_theme.tmpl", fileContext); err != nil {
		log.Fatalf("Failed to execute go theme template: %v", err)
	}

	if err := os.WriteFile(cfg.ThemeOutput.OutputFile, goCodeBuf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("Go maps utilizing official M3 token naming conventions written to ./theme/tokens.go!")

	log.Println("Running gofmt on generated files...")
	if err := exec.Command("gofmt", "-w", cfg.Generator.SchemaOutputDir).Run(); err != nil {
		log.Printf("Warning: gofmt failed on %s: %v", cfg.Generator.SchemaOutputDir, err)
	}
	if err := exec.Command("gofmt", "-w", cfg.ThemeOutput.OutputDir).Run(); err != nil {
		log.Printf("Warning: gofmt failed on %s: %v", cfg.ThemeOutput.OutputDir, err)
	}
}

func camelToKebab(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('-')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "-", ".")
	s = strings.ReplaceAll(s, "_", ".")
	return s
}

func buildSassColorMap(schemeColors map[string]string) string {
	var b strings.Builder
	keys := make([]string, 0, len(schemeColors))
	for k := range schemeColors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		v := schemeColors[k]
		kebabKey := camelToKebab(k)
		if i > 0 {
			b.WriteString(",\n")
		}
		b.WriteString(fmt.Sprintf("  '%s': %s", kebabKey, v))
	}
	return b.String()
}

func resolveAliases(allTokens map[string]string) map[string]string {
	resolved := make(map[string]string, len(allTokens))

	for key, value := range allTokens {
		resolvedValue := value
		if aliasRegex.MatchString(value) {
			resolvedValue = aliasRegex.ReplaceAllStringFunc(value, func(match string) string {
				tokenPath := match[1 : len(match)-1]
				normalizedPath := normalize(tokenPath)
				if v, ok := allTokens[normalizedPath]; ok {
					return v
				}
				if v, ok := allTokens["md."+normalizedPath]; ok {
					return v
				}
				return match
			})
		}
		resolvedValue = strings.ReplaceAll(strings.ReplaceAll(resolvedValue, "{", ""), "}", "")
		officialM3Key := normalize(key)
		resolved[officialM3Key] = resolvedValue
	}

	return resolved
}
