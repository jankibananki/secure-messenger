package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jankibananki/secure-messenger/internal/watcher"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func newAutoTransferPage(window fyne.Window) fyne.CanvasObject {
	var selectedDir string
	var stopMonitoring func()
	var generation uint64

	folderPath := widget.NewLabel("No folder selected")
	folderPath.Wrapping = fyne.TextWrapWord

	monitoringStatus := widget.NewLabel("INACTIVE")

	eventLog := widget.NewMultiLineEntry()
	eventLog.SetMinRowsVisible(7)
	eventLog.SetText("No events yet.")
	eventLog.Disable()

	var logLines []string

	addLog := func(message string) {
		line := fmt.Sprintf(
			"[%s] %s",
			time.Now().Format("15:04:05"),
			message,
		)

		logLines = append(logLines, line)

		if len(logLines) > 100 {
			logLines = logLines[len(logLines)-100:]
		}

		eventLog.SetText(strings.Join(logLines, "\n"))
	}

	var startButton *widget.Button
	var stopButton *widget.Button
	var chooseButton *widget.Button

	// Start local monitoring
	startButton = widget.NewButton("Start monitoring", func() {
		if selectedDir == "" || stopMonitoring != nil {
			return
		}

		generation++
		id := generation

		stop, err := watcher.Start(
			selectedDir,

			func(path string) {
				// Watcher callbacks run in a goroutine.
				fyne.Do(func() {
					if generation != id || stopMonitoring == nil {
						return
					}

					addLog("New file: " + filepath.Base(path))
				})
			},

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

	// Stop local monitoring
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

	// Choose local directory
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

	folderPanel := widget.NewCard(
		"Target directory",
		"Choose the directory to monitor",
		container.NewVBox(
			folderPath,
			chooseButton,
		),
	)

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

	eventsPanel := widget.NewCard(
		"Detected files",
		"File system events",
		eventLog,
	)

	return container.NewVScroll(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"AUTO TRANSFER",
				fyne.TextAlignLeading,
				fyne.TextStyle{Bold: true},
			),
			widget.NewSeparator(),
			folderPanel,
			monitoringPanel,
			eventsPanel,
		),
	)
}
