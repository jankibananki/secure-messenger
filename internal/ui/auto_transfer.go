package ui

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/jankibananki/secure-messenger/internal/activitylog"
	"github.com/jankibananki/secure-messenger/internal/watcher"
)

func newAutoTransferPage(window fyne.Window) fyne.CanvasObject {
	var selectedDir string
	var stopMonitoring func()
	var generation uint64

	folderPath := widget.NewLabel("No folder selected")
	folderPath.Wrapping = fyne.TextWrapWord

	monitoringStatus := widget.NewLabel("INACTIVE")

	// Belezenje dogadjaja iskljucivo u log fajl
	addLog := func(message string) {
		if err := activitylog.Record(message); err != nil {
			log.Printf("Failed to write activity log: %v", err)
		}
	}

	var startButton *widget.Button
	var stopButton *widget.Button
	var chooseButton *widget.Button

	// Pokretanje lokalnog monitoringa
	startButton = widget.NewButton("Start monitoring", func() {
		if selectedDir == "" || stopMonitoring != nil {
			return
		}

		generation++
		id := generation

		stop, err := watcher.Start(
			selectedDir,

			// Detekcija novog fajla
			func(path string) {
				fyne.Do(func() {
					if generation != id || stopMonitoring == nil {
						return
					}

					addLog("New file: " + filepath.Base(path))
				})
			},

			// Obrada gresaka watchera
			func(err error) {
				fyne.Do(func() {
					if generation != id || stopMonitoring == nil {
						return
					}

					addLog("Watcher error: " + err.Error())
				})
			},
		)

		if err != nil {
			addLog("Failed to start monitoring: " + err.Error())
			dialog.ShowError(err, window)
			return
		}

		stopMonitoring = stop

		monitoringStatus.SetText("ACTIVE - Monitoring directory")

		startButton.Disable()
		stopButton.Enable()
		chooseButton.Disable()

		addLog("Monitoring started: " + selectedDir)
	})

	startButton.Importance = widget.HighImportance
	startButton.Disable()

	// Zaustavljanje lokalnog monitoringa
	stopButton = widget.NewButton("Stop monitoring", func() {
		if stopMonitoring == nil {
			return
		}

		stopMonitoring()
		stopMonitoring = nil
		generation++

		monitoringStatus.SetText("INACTIVE")

		stopButton.Disable()
		startButton.Enable()
		chooseButton.Enable()

		addLog("Monitoring stopped")
	})

	stopButton.Disable()

	// Izbor direktorijuma
	chooseButton = widget.NewButtonWithIcon(
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

					path := folder.Path()

					info, err := os.Stat(path)
					if err != nil || !info.IsDir() {
						dialog.ShowInformation(
							"Invalid directory",
							"The selected directory is not accessible.",
							window,
						)
						return
					}

					selectedDir = path
					folderPath.SetText(path)
					startButton.Enable()
				},
				window,
			)
		},
	)

	// Panel za izbor direktorijuma
	folderPanel := widget.NewCard(
		"Target directory",
		"Choose the directory to monitor",
		container.NewVBox(
			folderPath,
			chooseButton,
		),
	)

	// Panel za kontrolu monitoringa
	monitoringPanel := widget.NewCard(
		"Monitoring",
		"Local test mode - no files are transferred",
		container.NewVBox(
			monitoringStatus,
			container.NewGridWithColumns(
				2,
				startButton,
				stopButton,
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
	)
}
