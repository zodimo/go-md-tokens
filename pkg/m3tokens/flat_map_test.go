package m3tokens

import (
	"testing"

	m3theme "github.com/zodimo/go-md-tokens/theme"
)

func TestFlatMap(t *testing.T) {
	th, err := NewTheme(m3theme.GetM3LightTokens())
	if err != nil {
		t.Fatalf("Failed to create theme: %v", err)
	}

	val1 := "#abcdef"
	val2 := "#123456"
	val3 := "#789012"

	th.FilledButtonTokens.Container.Color = val1

	th.FilledButtonTokens.Hover = &FilledButtonHoverOverlay{
		LabelTextColor: &val2,
	}

	th.RegisterCustomComponent("my-widget", map[string]string{
		"custom-token": val3,
	})

	m := th.FlatMap()

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
