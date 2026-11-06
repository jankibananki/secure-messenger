package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func newFilesPage(window fyne.Window) fyne.CanvasObject {
	fileName := widget.NewLabel("No file selected")
	filePath := widget.NewLabel("")
	fileInfo := widget.NewLabel("Size: - | Extension: -")

	filePath.Wrapping = fyne.TextWrapWord
	fileName.Wrapping = fyne.TextWrapWord

	algorithmSelect := widget.NewSelect(
		[]string{"TEA", "LEA-CTR"},
		nil,
	)
	algorithmSelect.SetSelected("LEA-CTR")

	sendButton := widget.NewButton("Send file", nil)
	sendButton.Importance = widget.HighImportance
	sendButton.Disable()

	chooseButton := widget.NewButtonWithIcon(
		"Choose file",
		theme.FolderOpenIcon(),
		func() {
			dialog.ShowFileOpen(
				func(file fyne.URIReadCloser, err error) {
					if err != nil {
						dialog.ShowError(err, window)
						return
					}
					if file == nil {
						return
					}
					defer file.Close()

					uri := file.URI()
					name := uri.Name()

					fileName.SetText(name)
					filePath.SetText(uri.Path())

					extension := strings.ToUpper(
						strings.TrimPrefix(filepath.Ext(name), "."),
					)
					if extension == "" {
						extension = "None"
					}

					sizeText := "Unknown"
					if uri.Scheme() == "file" {
						if info, err := os.Stat(uri.Path()); err == nil {
							sizeText = formatFileSize(info.Size())
						}
					}

					fileInfo.SetText(fmt.Sprintf(
						"Size: %s | Extension: %s",
						sizeText, extension,
					))
				},
				window,
			)
		},
	)

	filePanel := widget.NewCard(
		"Selected file",
		"File information",
		container.NewVBox(
			chooseButton,
			widget.NewSeparator(),
			fileName,
			fileInfo,
			filePath,
		),
	)

	encryptionPanel := widget.NewCard(
		"Encryption",
		"Select encryption algorithm",
		algorithmSelect,
	)

	transferPanel := widget.NewCard(
		"Transfer",
		"TCP transfer is not implemented yet",
		sendButton,
	)

	return container.NewVBox(
		widget.NewLabelWithStyle(
			"SECURE FILE TRANSFER",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		widget.NewSeparator(),
		filePanel,
		encryptionPanel,
		transferPanel,
	)
}

func formatFileSize(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2f MB",
			float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.2f GB",
			float64(size)/(1024*1024*1024))
	}
}
