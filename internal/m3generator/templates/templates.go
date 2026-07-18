package templates

import (
	"embed"
	"text/template"
)

//go:embed *.tmpl
var templateFiles embed.FS

// ParsedTemplates contains all the loaded templates ready for execution.
var ParsedTemplates *template.Template

func init() {
	ParsedTemplates = template.Must(template.ParseFS(templateFiles, "*.tmpl"))
}
