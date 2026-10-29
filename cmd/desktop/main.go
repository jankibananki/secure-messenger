
package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
)

func main() {
    application := app.NewWithID("com.jankibananki.securemessenger")
    application.Settings().SetTheme(theme.DarkTheme())

    window := application.NewWindow("Secure Messenger")
    window.Resize(fyne.NewSize(1100, 700))
    window.CenterOnScreen()

    // Gornja traka
    title := widget.NewLabelWithStyle(
        "SECURE MESSENGER",
        fyne.TextAlignLeading,
        fyne.TextStyle{Bold: true},
    )

    status := widget.NewLabel("Disconnected")

    header := container.NewHBox(
        title,
        layout.NewSpacer(),
        status,
    )

    // Centralni sadrzaj
    pageTitle := widget.NewLabelWithStyle(
        "Chat",
        fyne.TextAlignLeading,
        fyne.TextStyle{Bold: true},
    )

    pageDescription := widget.NewLabel(
        "Encrypted messaging will be available here.",
    )

    showPage := func(name, description string) {
        pageTitle.SetText(name)
        pageDescription.SetText(description)
    }

    // Leva navigacija
    sidebar := container.NewVBox(
        widget.NewButton("Chat", func() {
            showPage("Chat", "Encrypted messaging")
        }),
        widget.NewButton("Files", func() {
            showPage("Files", "Secure file transfers")
        }),
        widget.NewButton("Auto Transfer", func() {
            showPage("Auto Transfer", "File System Watcher")
        }),
        widget.NewButton("History", func() {
            showPage("History", "Encrypted conversation history")
        }),
        widget.NewButton("Activity Log", func() {
            showPage("Activity Log", "Application events")
        }),
        widget.NewButton("Settings", func() {
            showPage("Settings", "Application configuration")
        }),
    )

    content := container.NewVBox(
        pageTitle,
        widget.NewSeparator(),
        pageDescription,
    )

    window.SetContent(
        container.NewPadded(
            container.NewBorder(
                header,
                nil,
                sidebar,
                nil,
                container.NewPadded(content),
            ),
        ),
    )

    window.ShowAndRun()
}
