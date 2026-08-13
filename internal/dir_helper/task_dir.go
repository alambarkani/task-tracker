package dir_helper

import (
	"log"
	"os"
	"path/filepath"
)

func TaskDir() string {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	appDir := filepath.Join(baseDir, "task_tracker")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Fatal(err)
	}

	return filepath.Join(appDir, "tasks.json")
}
