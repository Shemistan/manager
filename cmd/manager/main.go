package main

import (
	"log"
	"os"

	"github.com/Shemistan/manager/internal/app/manager"
)

func main() {
	if err := manager.Run(); err != nil {
		log.Printf("application error: %v", err)
		os.Exit(1)
	}
}
