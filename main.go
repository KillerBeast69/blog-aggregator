package main

import (
	"fmt"
	"log"

	"github.com/KillerBeast69/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	fmt.Printf("initial config: %+v\n", cfg)

	err = cfg.SetUser("om")
	if err != nil {
		log.Fatalf("failed to set user: %v", err)
	}
	fmt.Println("successfully updated user")

	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read updated file: %v", err)
	}
	fmt.Printf("updated config: %+v\n", updatedCfg)

}
