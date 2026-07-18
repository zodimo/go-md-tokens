package m3tokens

import (
	"testing"
)

func TestFlatMap(t *testing.T) {
	theme := NewTheme("light")

	val1 := "#abcdef"
	val2 := "#123456"
	val3 := "#789012"

	theme.FilledButtonTokens.Container.Color = val1

	theme.FilledButtonTokens.Hover = &FilledButtonHoverOverlay{
		LabelTextColor: &val2,
	}

	theme.RegisterCustomComponent("my-widget", map[string]string{
		"custom-token": val3,
	})

	m := theme.FlatMap()

	if m["md.comp.filled.button.container.color"] != val1 {
		t.Errorf("Expected FlatMap to contain base token 'md.comp.filled.button.container.color' with value %q, got %q", val1, m["md.comp.filled.button.container.color"])
	}

	if m["md.comp.filled.button.hover.label.text.color"] != val2 {
		t.Errorf("Expected FlatMap to contain overlay token 'md.comp.filled.button.hover.label.text.color' with value %q, got %q", val2, m["md.comp.filled.button.hover.label.text.color"])
	}

	if m["md.comp.my.widget.custom.token"] != val3 {
		t.Errorf("Expected FlatMap to contain custom token 'md.comp.my.widget.custom.token' with value %q, got %q", val3, m["md.comp.my.widget.custom.token"])
	}
}
