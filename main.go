package main

import (
	"fmt"
	"log"

	"github.com/utam0k/ws-dbg/crt"
)

func main() {
	fmt.Println("Hello, World")
	cc, err := crt.NewClient()
	if err != nil {
		log.Fatalf("cannot connet to cllient: %v", err)
	}
	cc.FetchContainers()
}
