package main

import (
	"flag"
	"fmt"
)

func main() {
	config := flag.String("f", "../../configs/config.yaml", "Config file path")
	flag.Parse()
	fmt.Println(config)
}
