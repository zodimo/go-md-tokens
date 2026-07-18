package m3tokens

import (
	"testing"
)

func TestStateFallbackScenarios(t *testing.T) {
	theme := NewTheme("light")
	
	baseColor := "#111111"
	hoverColor := "#222222"
	
	theme.FilledButtonTokens.Container.Color = baseColor
	
	theme.FilledButtonTokens.Hover = &FilledButtonHoverOverlay{
		LabelTextColor: &hoverColor,
	}
	
	if got := theme.Resolver.FilledButton().ContainerColor(StateHover); got != baseColor {
		t.Errorf("Expected fallback to base ContainerColor %q, got %q", baseColor, got)
	}
	
	if got := theme.Resolver.FilledButton().ContainerColor(StateFocus); got != baseColor {
		t.Errorf("Expected fallback to base ContainerColor %q, got %q", baseColor, got)
	}
	
	baseLabelColor := "#000000"
	theme.FilledButtonTokens.LabelText.Color = baseLabelColor
	if got := theme.Resolver.FilledButton().LabelTextColor(StateHover); got != hoverColor {
		t.Errorf("Expected overlay LabelTextColor %q, got %q", hoverColor, got)
	}
}

func TestCompoundStateResolution(t *testing.T) {
	theme := NewTheme("light")
	
	baseColor := "base"
	selectedHoverColor := "selected-hover"
	hoverColor := "hover"
	
	theme.CheckboxTokens.StateLayer.Color = baseColor
	
	theme.CheckboxTokens.SelectedHover = &CheckboxSelectedHoverOverlay{
		StateLayerColor: &selectedHoverColor,
	}
	
	theme.CheckboxTokens.Hover = &CheckboxHoverOverlay{
		StateLayerColor: &hoverColor,
	}
	
	compoundState := StateSelected | StateHover
	if got := theme.Resolver.Checkbox().StateLayerColor(compoundState); got != selectedHoverColor {
		t.Errorf("Expected compound state value %q, got %q", selectedHoverColor, got)
	}
	
	theme.CheckboxTokens.SelectedHover = &CheckboxSelectedHoverOverlay{
		// StateLayerColor is nil
	}
	
	if got := theme.Resolver.Checkbox().StateLayerColor(compoundState); got != hoverColor {
		t.Errorf("Expected fallback to single state value %q, got %q", hoverColor, got)
	}
}
