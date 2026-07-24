package helpers

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func loadTestEnv() error {

	dir, err := os.Getwd()

	if err != nil {
		return err
	}

	for {

		envPath := filepath.Join(
			dir,
			".env.test",
		)

		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			break
		}

		dir = parent
	}

	return os.ErrNotExist
}
