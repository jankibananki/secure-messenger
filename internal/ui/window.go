package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	application := app.NewWithID(
		"com.jankibananki.securemessenger",
	)

	application.Settings().SetTheme(&messengerTheme{})

	window := application.NewWindow("Secure Messenger")
	window.Resize(fyne.NewSize(1150, 720))
	window.CenterOnScreen()

	pages := map[string]fyne.CanvasObject{
		"Chat": newChatPage(),
		"Files": newSectionPage(
			"FILES",
		),
		"Auto Transfer": newSectionPage(
			"AUTO TRANSFER",
		),
		"History": newSectionPage(
			"HISTORY",
		),
		"Activity Log": newSectionPage(
			"ACTIVITY LOG",
		),
		"Settings": newSectionPage(
			"SETTINGS",
		),
	}

	mainArea := container.NewMax(pages["Chat"])
	title := canvas.NewText(
		"SECURE MESSENGER", colorText,
	)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}

	status := canvas.NewText("● OFFLINE", colorMuted)
	status.TextSize = 13

	header := container.NewVBox(
		container.NewHBox(
			title,
			layout.NewSpacer(),
			status,
		),
		widget.NewSeparator(),
	)
	navigation := []struct {
		name string
		icon fyne.Resource
	}{
		{"Chat", theme.MailComposeIcon()},
		{"Files", theme.FolderIcon()},
		{"Auto Transfer", theme.FolderOpenIcon()},
		{"History", theme.HistoryIcon()},
		{"Activity Log", theme.InfoIcon()},
		{"Settings", theme.SettingsIcon()},
	}

	backgrounds := make(map[string]*canvas.Rectangle)

	showPage := func(name string) {
		mainArea.Objects = []fyne.CanvasObject{pages[name]}
		mainArea.Refresh()

		for key, bg := range backgrounds {
			if key == name {
				bg.FillColor = colorPurple
			} else {
				bg.FillColor = colorSurface
			}
			bg.Refresh()
		}
	}

	navRows := make([]fyne.CanvasObject, 0, len(navigation))

	for _, item := range navigation {
		pageName := item.name

		bg := canvas.NewRectangle(colorSurface)
		bg.CornerRadius = 8

		if pageName == "Chat" {
			bg.FillColor = colorPurple
		}

		backgrounds[pageName] = bg

		button := widget.NewButtonWithIcon(
			pageName,
			item.icon,
			func() {
				showPage(pageName)
			},
		)

		button.Importance = widget.LowImportance
		button.Alignment = widget.ButtonAlignLeading

		navRows = append(
			navRows,
			container.NewStack(bg, button),
		)
	}

	sidebar := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("WORKSPACE"),
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewGridWrap(
			fyne.NewSize(195, 44),
			navRows...,
		),
	)

	content := container.NewBorder(
		header,
		nil,
		container.NewPadded(sidebar),
		nil,
		container.NewPadded(mainArea),
	)

	window.SetContent(container.NewPadded(content))
	window.ShowAndRun()
}

func newSectionPage(
	title string,
) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle(
			title,
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		widget.NewSeparator(),
	)
}
