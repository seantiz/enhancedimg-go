package main

import (
	"log"
	"os"

	"github.com/seantiz/enhancedimg-go/enhancedimg"
)

func main() {
	if err := os.MkdirAll("static/processed", 0755); err != nil {
		log.Fatal(err)
	}

	if err := enhancedimg.FindAllImageElements("routes"); err != nil {
		log.Fatal(err)
	}
}
