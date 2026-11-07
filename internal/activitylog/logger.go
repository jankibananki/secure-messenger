package activitylog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var mu sync.Mutex

func Record(event string) error {
	mu.Lock()
	defer mu.Unlock()

	// Pronalazi direktorijum za cuvanje podataka aplikacije
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	logDir := filepath.Join(
		configDir, "secure-messenger", "logs",
	)

	// Kreira direktorijum ako ne postoji
	if err := os.MkdirAll(logDir, 0700); err != nil {
		return err
	}

	// Otvara log fajl u append modu
	file, err := os.OpenFile(
		filepath.Join(logDir, "activity.log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	// Belezi dogadjaj zajedno sa vremenom
	_, err = fmt.Fprintf(
		file,
		"[%s] %s\n",
		time.Now().Format(time.RFC3339),
		event,
	)

	return err
}
