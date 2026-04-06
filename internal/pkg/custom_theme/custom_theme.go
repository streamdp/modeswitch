package custom_theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct {
	fyne.Theme
}

func (t *Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameDisabled {
		if variant == theme.VariantLight {
			return color.Gray16{Y: 0x5b5b} // Bright Gray
		}

		return color.White
	}

	return theme.DefaultTheme().Color(name, variant)
}
