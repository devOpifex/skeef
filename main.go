package main

import (
	"fmt"
	"log"

	"github.com/devOpifex/opiflex/config"
)

func main() {
	config, err := config.Read()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}
