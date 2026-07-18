package m3tokens

// TokenState represents the state of a UI component token.
// It is a bitmask to allow for compound states.
type TokenState uint32

// TokenCategory represents a category of tokens (e.g. Container, Icon, LabelText)
type TokenCategory string

// TokenProperty represents a specific property of a token (e.g. Color, Opacity, Elevation)
type TokenProperty string
