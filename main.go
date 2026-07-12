package main

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/bep/godartsass/v2"
)

func main() {
	transpiler, err := godartsass.Start(godartsass.Options{})
	if err != nil {
		log.Fatal(err)
	}

	tokensDir := filepath.Join("m3web", "v2.4.1", "tokens")
	absDir, err := filepath.Abs(tokensDir)
	if err != nil {
		log.Fatal(err)
	}

	source := `
@use 'md-comp-elevated-button';

:root {
  @each $token, $value in md-comp-elevated-button.values() {
    --md-elevated-button-#{$token}: #{$value};
  }
}
`

	res, err := transpiler.Execute(godartsass.Args{
		Source:       source,
		URL:          (&url.URL{Scheme: "file", Path: filepath.Join(absDir, "gen.css.scss")}).String(),
		IncludePaths: []string{absDir},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.CSS)
}
