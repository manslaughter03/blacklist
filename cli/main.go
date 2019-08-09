package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"git.hydra-project.io/banks/blacklist"
)

var (
	verbose *bool
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s [<ip>]:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var err error
	inputFile := flag.String("input-file", "", "Input with the list of ip to check")
	outputFile := flag.String("output-file", "", "Output file")
	limit := flag.Int("limit", 0, "Limit the list to check")
	verbose = flag.Bool("verbose", false, "Verbosity of the program")
	help := flag.Bool("help", false, "Print usage of current program")
	workers := flag.Int("workers", 5, "Workers to run (goroutine)")

	flag.Parse()
	if *help {
		Usage()
	}
	var results = make(map[string]blacklist.Result)

	if *inputFile != "" {
		if *verbose {
			log.Printf("[DEBUG] Try to open file: %s", *inputFile)
		}
		bytes, err := ioutil.ReadFile(*inputFile)
		if err != nil {
			log.Fatal(err)
		}

		ips := strings.Split(string(bytes), "\n")
		// Remove last element empty
		ips = ips[:len(ips)-1]
		if *limit > 0 {
			ips = ips[:*limit]
		}

		for _, ip := range ips {

			res, err := blacklist.Start(ip, *workers)
			if err != nil {
				log.Fatal(err)
			}
			results[ip] = res
		}

	} else if len(os.Args) > 1 {
		ip := os.Args[1]
		results[ip], err = blacklist.Start(ip, *workers)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		Usage()
	}

	fmt.Println(Summary(results))

	if *outputFile != "" {
		bytes, err := json.Marshal(results)
		if err != nil {
			log.Fatalf("can't unmarshall results, got: %v", err)
		}
		if err = ioutil.WriteFile(*outputFile, bytes, 0644); err != nil {
			log.Fatalf("Can't write to file %s, got %v", *outputFile, err)
		}
	}
}

func Summary(results map[string]blacklist.Result) string {
	var resume string
	resume += "----------------------------------------------------\n"
	resume += fmt.Sprintf("%-22s | %-9s | %-9s\n", "IP", "Blacklisted", "Count detect")
	resume += "----------------------------------------------------\n"
	for ip, result := range results {
		resume += fmt.Sprintf("%-22s | %-9v | %-9d\n", ip, !result.Valid, len(result.DomainDetects))
		resume += "----------------------------------------------------\n"
	}

	return resume
}
