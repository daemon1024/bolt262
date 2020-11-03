package main

import (
	"log"
	"os"

	"github.com/nakabonne/gosivy/agent"
)

var count, pass, fail int64 = 0, 0, 0
var includePath string

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	defer agent.Close()
	path := os.Args[1]
	includePath = os.Args[2]
	fi, err := os.Stat(path)
	if err != nil {
		log.Print(err)
	}
	if fi.Mode().IsDir() {
		runTests(path, includePath)
	} else {
		// processFile(path)
		log.Print("commented process")
	}
}
