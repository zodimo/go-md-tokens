package main

import (
	"bytes"
	"fmt"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

func main() {
	// true because this is the content of an inline style attribute
	p := css.NewParser(parse.NewInput(bytes.NewBufferString("color: red;")), true)
	out := ""
	for {
		gt, _, data := p.Next()
		if gt == css.ErrorGrammar {
			break
		} else if gt == css.AtRuleGrammar || gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar || gt == css.DeclarationGrammar {
			out += string(data)
			if gt == css.DeclarationGrammar {
				out += ":"
			}
			for _, val := range p.Values() {
				out += string(val.Data)
			}
			if gt == css.BeginAtRuleGrammar || gt == css.BeginRulesetGrammar {
				out += "{"
			} else if gt == css.AtRuleGrammar || gt == css.DeclarationGrammar {
				out += ";"
			}
		} else {
			out += string(data)
		}
	}
	fmt.Println(out)
}
