package main

import (
	"log"
	"os"

	"github.com/Shemistan/manager/internal/app/migrator"
)

func main() {
	if err := migrator.Run(); err != nil {
		log.Printf("migrator error: %v", err)
		os.Exit(1)
	}
}
