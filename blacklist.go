package blacklist

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

const (
	blacklistFile = "./blacklist.txt"
)

// Result struct
type Result struct {
	Valid         bool
	DomainDetects []string
}

// String function
func (result Result) String() string {
	return fmt.Sprintf("Valid: %v, DomainDetects (%d): %v", result.Valid, len(result.DomainDetects), result.DomainDetects)
}

// Work struct
type Work struct {
	Blacklists []string
	Result     Result
}

// String function
func (work Work) String() string {
	return fmt.Sprintf("Blacklists: %v, Result: %v", work.Blacklists, work.Result)
}

func (work *Work) worker(ip string, blacklistChan chan string, resultChan chan bool) {
	for blacklist := range blacklistChan {
		ok := CheckBlackList(ip, blacklist)
		if ok {
			resultChan <- false
			work.Result.DomainDetects = append(work.Result.DomainDetects, blacklist)
		} else {
			resultChan <- true
		}
	}
}

// Start blacklist on ip with n workers
func Start(ip string, workers int) (Result, error) {
	var work = new(Work)
	blacklists, err := LoadBlacklistDomainTest(blacklistFile)
	if err != nil {
		return Result{}, err
	}
	work.Blacklists = blacklists

	resultChan := make(chan bool)
	blacklistChan := make(chan string)

	for i := 0; i < workers-1; i++ {
		go work.worker(ip, blacklistChan, resultChan)
	}

	go func() {
		for _, blacklist := range work.Blacklists {
			blacklistChan <- blacklist
		}
	}()

	work.Result.Valid = true
	for range blacklists {
		res := <-resultChan
		if !res {
			work.Result.Valid = false
		}

	}

	return work.Result, nil

}

// LoadBlacklistDomainTest Load blacklist domains test from a file in parameter
func LoadBlacklistDomainTest(filename string) ([]string, error) {
	var blacklists []string
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return blacklists, err
	}
	return strings.Split(string(bytes), "\n"), nil
}

// CheckBlackList Check if ip givent is blacklist on dnsTest
func CheckBlackList(ip, dnsTest string) bool {
	reverseIP := reverse(strings.Split(ip, "."))
	dnsNameTest := fmt.Sprintf("%s.%s", strings.Join(reverseIP, "."), dnsTest)
	ips, err := net.LookupAddr(dnsNameTest)
	if err == nil && len(ips) > 0 {
		return true
	}
	ips, err = net.LookupTXT(dnsNameTest)
	if err == nil && len(ips) > 0 {
		return true
	}

	return false
}

func reverse(numbers []string) []string {
	newNumbers := make([]string, len(numbers))
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		newNumbers[i], newNumbers[j] = numbers[j], numbers[i]
	}
	return newNumbers
}
