package main

import (
	"log"

	"github.com/reinp/event-platform/backend/internal/app"
)

func main() {

	application, err :=
		app.New()

	if err != nil {

		log.Fatal(err)

	}

	if err :=
		application.Start(); err != nil {

		log.Fatal(err)

	}
}
