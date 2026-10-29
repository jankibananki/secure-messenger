package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func newChatBubble(
	message string,
	algorithm string,
	outgoing bool,
) fyne.CanvasObject {

	backgroundColor := colorControl
	sender := "CONTACT"

	if outgoing {
		backgroundColor = colorPurple
		sender = "YOU"
	}

	background := canvas.NewRectangle(backgroundColor)
	background.CornerRadius = 12
	metadata := widget.NewLabel(
		sender + "  |  " + algorithm + "  |  " +
			time.Now().Format("15:04"),
	)
	metadata.SizeName = theme.SizeNameCaptionText
	metadata.Alignment = fyne.TextAlignLeading
	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord
	messageLabel.Alignment = fyne.TextAlignLeading
	content := container.NewVBox(
		metadata,
		messageLabel,
	)
	paddedContent := container.New(
		layout.NewCustomPaddedLayout(6, 6, 10, 10),
		content,
	)

	bubble := container.NewStack(
		background,
		paddedContent,
	)

	empty := layout.NewSpacer()

	if outgoing {
		return container.NewGridWithColumns(
			2, empty, bubble,
		)
	}

	return container.NewGridWithColumns(
		2, bubble, empty,
	)
}

func newChatPage() fyne.CanvasObject {
	ipEntry := widget.NewEntry()
	ipEntry.SetText("127.0.0.1")

	portEntry := widget.NewEntry()
	portEntry.SetText("5050")

	connectButton := widget.NewButton("Connect", nil)
	connectButton.Disable()

	connectionStatus := canvas.NewText(
		"OFFLINE", colorMuted,
	)
	connectionStatus.TextSize = 12

	connectionBar := container.NewBorder(
		nil,
		nil,
		connectionStatus,
		connectButton,
		container.NewGridWithColumns(
			2, ipEntry, portEntry,
		),
	)
	messages := container.NewVBox()

	messageScroll := container.NewVScroll(messages)
	messageInput := widget.NewMultiLineEntry()
	messageInput.SetPlaceHolder("Write your message...")
	messageInput.SetMinRowsVisible(2)

	algorithmSelect := widget.NewSelect(
		[]string{"TEA", "LEA-CTR"},
		nil,
	)
	algorithmSelect.SetSelected("LEA-CTR")

	previewButton := widget.NewButton(
		"Preview message",
		func() {
			text := messageInput.Text

			if text == "" {
				return
			}

			messages.Add(
				newChatBubble(
					text,
					algorithmSelect.Selected,
					true,
				),
			)

			messageScroll.ScrollToBottom()
			messageInput.SetText("")
		},
	)

	previewButton.Importance = widget.HighImportance

	composer := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			widget.NewLabel("Encryption:"),
			algorithmSelect,
		),
		container.NewBorder(
			nil,
			nil,
			nil,
			previewButton,
			messageInput,
		),
	)

	header := container.NewVBox(
		widget.NewLabelWithStyle(
			"ENCRYPTED CHAT",
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		widget.NewLabel(
			"Peer-to-peer communication workspace",
		),
		widget.NewSeparator(),
		connectionBar,
		widget.NewSeparator(),
	)

	return container.NewBorder(
		header,
		composer,
		nil,
		nil,
		messageScroll,
	)
}
