package m3tokens

import (
	"testing"

	m3theme "github.com/zodimo/go-md-tokens/theme"
)

func TestStateFallbackScenarios(t *testing.T) {
	th, err := NewTheme(m3theme.GetM3LightTokens())
	if err != nil {
		t.Fatalf("Failed to create theme: %v", err)
	}

	baseColor := "#111111"
	hoverColor := "#222222"

	th.FilledButtonTokens.Container.Color = baseColor

	th.FilledButtonTokens.Hover = &FilledButtonHoverOverlay{
		LabelTextColor: &hoverColor,
	}

	if got := th.Resolver.FilledButton().ContainerColor(StateHover); got != baseColor {
		t.Errorf("Expected fallback to base ContainerColor %q, got %q", baseColor, got)
	}

	if got := th.Resolver.FilledButton().ContainerColor(StateFocus); got != baseColor {
		t.Errorf("Expected fallback to base ContainerColor %q, got %q", baseColor, got)
	}

	baseLabelColor := "#000000"
	th.FilledButtonTokens.LabelText.Color = baseLabelColor
	if got := th.Resolver.FilledButton().LabelTextColor(StateHover); got != hoverColor {
		t.Errorf("Expected overlay LabelTextColor %q, got %q", hoverColor, got)
	}
}

func TestCompoundStateResolution(t *testing.T) {
	th, err := NewTheme(m3theme.GetM3LightTokens())
	if err != nil {
		t.Fatalf("Failed to create theme: %v", err)
	}

	baseColor := "base"
	selectedHoverColor := "selected-hover"
	hoverColor := "hover"

	th.CheckboxTokens.StateLayer.Color = baseColor

	th.CheckboxTokens.SelectedHover = &CheckboxSelectedHoverOverlay{
		StateLayerColor: &selectedHoverColor,
	}

	th.CheckboxTokens.Hover = &CheckboxHoverOverlay{
		StateLayerColor: &hoverColor,
	}

	compoundState := StateSelected | StateHover
	if got := th.Resolver.Checkbox().StateLayerColor(compoundState); got != selectedHoverColor {
		t.Errorf("Expected compound state value %q, got %q", selectedHoverColor, got)
	}

	th.CheckboxTokens.SelectedHover = &CheckboxSelectedHoverOverlay{
		// StateLayerColor is nil
	}

	if got := th.Resolver.Checkbox().StateLayerColor(compoundState); got != hoverColor {
		t.Errorf("Expected fallback to single state value %q, got %q", hoverColor, got)
	}
}

func TestCustomComponentResolution(t *testing.T) {
	th, err := NewTheme(m3theme.GetM3LightTokens())
	if err != nil {
		t.Fatalf("Failed to create theme: %v", err)
	}

	customTokens := map[string]string{
		"container-color":                "custom-base",
		"hover-container-color":          "custom-hover",
		"selected-hover-container-color": "custom-selected-hover",
	}

	th.RegisterCustomComponent("my-custom-widget", customTokens)

	// Test base token
	if got := th.Resolver.Custom("my-custom-widget", "container-color", StateDefault); got != "custom-base" {
		t.Errorf("Expected custom component base value %q, got %q", "custom-base", got)
	}

	// Test hover token
	if got := th.Resolver.Custom("my-custom-widget", "container-color", StateHover); got != "custom-hover" {
		t.Errorf("Expected custom component hover value %q, got %q", "custom-hover", got)
	}

	// Test compound token
	if got := th.Resolver.Custom("my-custom-widget", "container-color", StateSelected|StateHover); got != "custom-selected-hover" {
		t.Errorf("Expected custom component selected-hover value %q, got %q", "custom-selected-hover", got)
	}

	// Test fallback for missing state (focus)
	if got := th.Resolver.Custom("my-custom-widget", "container-color", StateFocus); got != "custom-base" {
		t.Errorf("Expected custom component fallback to base value %q, got %q", "custom-base", got)
	}

	// Test missing component
	if got := th.Resolver.Custom("non-existent", "container-color", StateDefault); got != "" {
		t.Errorf("Expected empty string for non-existent component, got %q", got)
	}
}
