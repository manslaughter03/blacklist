package main

import (
	"log"

	"git.hydra-project.io/banks/blacklist"
)

func main() {
	ip := "8.8.4.4"
	workers := 25
	res, err := blacklist.Start(ip, workers)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("res: %v", res)
}
