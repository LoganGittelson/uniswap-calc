package main

import (
	"log"

	"github.com/LoganGittelson/uniswap-calc/pkg/unicalc"
)

func main() {
	err := unicalc.RunCalcs()
	if err != nil {
		log.Print("Apllication returned error")
		log.Fatal(err)
	}
}
