package main

import (
	"fmt"
	"log"
	"main/config"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("must provide input file path")
	}
	fmt.Println("Parsing request...")
	cfg, err := config.Parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	if err := cfg.Exec(); err != nil {
		log.Fatalln(err)
	}
}
