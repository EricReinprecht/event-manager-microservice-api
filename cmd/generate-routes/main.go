package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func main() {

	// Load env
	_ = godotenv.Load(".env")

	if _, err := os.Stat(".env.local"); err == nil {
		_ = godotenv.Overload(".env.local")
	}

	output := os.Getenv(
		"ROUTES_OUTPUT_PATH",
	)

	if output == "" {
		log.Fatal(
			"ROUTES_OUTPUT_PATH is missing",
		)
	}

	// Allow CLI override
	if len(os.Args) > 1 {
		output = os.Args[1]
	}

	data, err := json.MarshalIndent(
		constants.Routes,
		"",
		"  ",
	)

	if err != nil {
		log.Fatal(err)
	}

	// Create missing folder structure
	dir := filepath.Dir(output)

	err = os.MkdirAll(
		dir,
		0755,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Create or overwrite file
	err = os.WriteFile(
		output,
		data,
		0644,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"routes generated:",
		output,
	)
}
