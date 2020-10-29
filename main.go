package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/smallfish/simpleyaml"
)

var count, pass, fail int64 = 0, 0, 0

func main() {
	path := os.Args[1]
	fi, err := os.Stat(path)
	if err != nil {
		log.Print(err)
	}
	if fi.Mode().IsDir() {
		walkDir(path)
	} else {
		// processFile(path)
		log.Print("commented process")
	}
}

func walkDir(rootPath string) {
	path := rootPath
	assert, err := ioutil.ReadFile(os.Args[2] + "assert.js")
	if err != nil {
		log.Print(err)
	}
	sta, err := ioutil.ReadFile(os.Args[2] + "sta.js")
	if err != nil {
		log.Print(err)
	}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if info.IsDir() == false {
			processFile(path, assert, sta)
		}
		return nil
	})
	fmt.Printf("\033[36mTotal files %d\n\033[0m\033[31mFailed Tests %d\n\033[0m\033[34mPassed Tests %d\033[0m", count, fail, pass)
}

func processFile(path string, assert []byte, sta []byte) {
	data, _ := ioutil.ReadFile(path)
	fmt.Println("\033[36m", path)
	var includeFinal []byte
	includeFinal = append(includeFinal, assert...)
	includeFinal = append(includeFinal, sta...)
	re, _ := regexp.Compile(`\/\*---([\s\S]*?)---\*\/`)
	y := re.FindSubmatch([]byte(data))
	if y != nil {
		yaml, err := simpleyaml.NewYaml(y[1])
		if err != nil {
			log.Print(path)
		}
		// bar := yaml
		// fmt.Printf("Value: %#v\n", bar)
		includes, _ := yaml.Get("includes").Array()
		fmt.Print(includes)
		for _, include := range includes {
			data, _ := ioutil.ReadFile(os.Args[2] + include.(string))
			includeFinal = append(includeFinal, data...)
		}
	}
	var finalFile []byte
	finalFile = append(finalFile, includeFinal...)
	finalFile = append(finalFile, data...)
	// currDir, _ := os.Getwd()
	// err = ioutil.WriteFile(currDir+"/tmp.js", finalFile, 0777)
	// if err != nil {
	// 	log.Print(err)
	// }
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmptest-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// Remember to clean up the file afterwards
	defer os.Remove(tmpFile.Name())
	count++
	fmt.Println("Created File: " + tmpFile.Name())

	// Example writing to the file
	if _, err = tmpFile.Write(finalFile); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	var outbuf, errbuf bytes.Buffer
	var exitCode int
	cmd := exec.Command("node", tmpFile.Name())
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err = cmd.Run()
	stderr := errbuf.String()

	if err != nil {
		fail++
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			log.Print("Could not get exit code for failed program")
			exitCode = 1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
		pass++
	}
	fmt.Println("\033[31m", stderr)
	fmt.Println("\033[36m", exitCode)
	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
}
