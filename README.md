# Blacklist checker

Check if an address ip (or a list) is blacklist by the common lists.

Simple example of implementation:
```golang
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
```
