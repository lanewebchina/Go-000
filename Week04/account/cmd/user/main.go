package main

import (
	"flag"
	"fmt"
	"log"

	"account/pkg/app"
)

func main() {
	config := flag.String("f", "../../configs/config.yaml", "Config file path")
	flag.Parse()
	fmt.Println(config)

	app, err := app.InitializeApp(*config)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
