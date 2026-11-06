package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func newAutoTransferPage(window fyne.Window) fyne.CanvasObject {

	folderPath := widget.NewLabel("No folder selected")
	folderPath.Wrapping = fyne.TextWrapWord

	monitoringStatus := widget.NewLabel(
		"Inactive - no directory is being monitored",
	)

	chooseButton := widget.NewButtonWithIcon(
		"Choose folder",
		theme.FolderOpenIcon(),
		func() {
			dialog.ShowFolderOpen(
				func(folder fyne.ListableURI, err error) {
					if err != nil {
						dialog.ShowError(err, window)
						return
					}

					if folder == nil {
						return
					}

					if folder.Scheme() != "file" {
						dialog.ShowInformation(
							"Invalid directory",
							"Please select a local folder.",
							window,
						)
						return
					}

					folderPath.SetText(folder.Path())
				},
				window,
			)
		},
	)

	folderPanel := widget.NewCard(
		"Target directory",
		"Select the folder to monitor",
		container.NewVBox(
			folderPath,
			chooseButton,
		),
	)

	startButton := widget.NewButton("Start monitoring", nil)
	startButton.Importance = widget.HighImportance
	startButton.Disable()

	stopButton := widget.NewButton("Stop monitoring", nil)
	stopButton.Disable()

	monitoringPanel := widget.NewCard(
		"Monitoring",
		"Automatic file transfer",
		container.NewVBox(
			monitoringStatus,
			container.NewGridWithColumns(
				2,
				startButton,
				stopButton,
			),
			widget.NewLabel(
				"Remote activation is not implemented yet.",
			),
		),
	)

	return container.NewVBox(
		widget.NewLabelWithStyle(
			"AUTO TRANSFER",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		widget.NewSeparator(),
		folderPanel,
		monitoringPanel,
		infoPanel,
	)
}
