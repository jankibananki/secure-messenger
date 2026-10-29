package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	colorBackground = color.NRGBA{14, 10, 19, 255}
	colorSurface    = color.NRGBA{28, 21, 31, 255}
	colorControl    = color.NRGBA{41, 32, 48, 255}
	colorPurple     = color.NRGBA{122, 31, 162, 255}
	colorPink       = color.NRGBA{255, 0, 85, 255}
	colorNeon       = color.NRGBA{187, 255, 31, 255}
	colorText       = color.NRGBA{239, 248, 255, 255}
	colorMuted      = color.NRGBA{164, 154, 175, 255}
)

type messengerTheme struct{}

func (m *messengerTheme) Color(
	name fyne.ThemeColorName,
	variant fyne.ThemeVariant,
) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colorBackground

	case theme.ColorNamePrimary:
		return colorPink

	case theme.ColorNameForegroundOnPrimary:
		return color.White

	case theme.ColorNameButton:
		return colorControl

	case theme.ColorNameInputBackground:
		return colorSurface

	case theme.ColorNameInputBorder:
		return colorControl

	case theme.ColorNameForeground:
		return colorText

	case theme.ColorNamePlaceHolder:
		return colorMuted

	case theme.ColorNameSuccess:
		return colorNeon

	case theme.ColorNameError:
		return color.NRGBA{255, 111, 97, 255}

	case theme.ColorNameFocus:
		return colorPink

	case theme.ColorNameSeparator:
		return colorControl

	case theme.ColorNameHover:
		return color.NRGBA{122, 31, 162, 60}
	}

	return theme.DarkTheme().Color(
		name, theme.VariantDark,
	)
}

func (m *messengerTheme) Font(
	style fyne.TextStyle,
) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

func (m *messengerTheme) Icon(
	name fyne.ThemeIconName,
) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (m *messengerTheme) Size(
	name fyne.ThemeSizeName,
) float32 {
	return theme.DarkTheme().Size(name)
}
