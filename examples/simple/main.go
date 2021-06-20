//+build ignore

package main

import (
	"log"
	"os"

	env "github.com/kaatinga/settings"
)

type settings struct {
	Port       string `env:"PORT" validate:"numeric"`
	Database   string `env:"DATABASE"`
	CacheSize  byte   `env:"CACHE_SIZE"`
	LaunchMode string `env:"LAUNCH_MODE"`
}

func main() {
	err := os.Setenv("PORT", "8080") // for validate pass
	if err != nil {
		log.Fatalf("set env error, %s", err)
	}
	var sett settings
	err = env.LoadSettings(&sett)
	if err != nil {
		log.Fatalf("load error happened, %s", err)
	}
	log.Printf("settings: %+v\n", sett)
}
