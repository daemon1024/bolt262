package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"syscall"

	"github.com/karrick/godirwalk"
	"github.com/smallfish/simpleyaml"
)

func runTests(testPath string, includePath string) {
	count, pass, fail = 0, 0, 0
	fmt.Println(testPath, includePath)
	path := testPath
	// includePath = "../test262/harness/" //hardcoded for bench
	assert, err := ioutil.ReadFile(includePath + "/assert.js")
	if err != nil {
		log.Print(err)
	}
	sta, err := ioutil.ReadFile(includePath + "/sta.js")
	if err != nil {
		log.Print(err)
	}
	var mustIncludes []byte
	mustIncludes = append(mustIncludes, assert...)
	mustIncludes = append(mustIncludes, sta...)
	var wg sync.WaitGroup
	err = godirwalk.Walk(path, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if !de.IsDir() {
				wg.Add(1)
				go func(path string) {
					//note to self : check file directive limit
					defer wg.Done()
					processFile(osPathname, mustIncludes, includePath)
				}(path)
			}

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {

			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	wg.Wait()
	fmt.Printf("\033[36mTotal files %d\n\033[0m\033[31mFailed Tests %d\n\033[0m\033[34mPassed Tests %d\033[0m\n", count, fail, pass)
}

func processFile(path string, mustIncludes []byte, includePath string) {

	data, _ := ioutil.ReadFile(path)

	var includeFinal []byte
	includeFinal = append(includeFinal, mustIncludes...)

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
		for _, include := range includes {
			data, _ := ioutil.ReadFile(includePath + include.(string))
			// fmt.Printf("%s Value: %#v\n", includePath+include.(string), data)
			includeFinal = append(includeFinal, data...)
		}
	}

	var finalFile []byte
	finalFile = append(finalFile, includeFinal...)
	finalFile = append(finalFile, data...)
	currDir, _ := os.Getwd()

	tmpFile, err := ioutil.TempFile(currDir+"/tmp/", "tmptest-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// Remember to clean up the file afterwards
	defer os.Remove(tmpFile.Name())
	count++

	// Example writing to the file
	if _, err = tmpFile.Write(finalFile); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	// run the test
	runFile(tmpFile.Name(), path)

	// Close the file
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

}

func runFile(tmpName string, filePath string) {
	var outbuf, errbuf bytes.Buffer
	var exitCode int
	cmd := exec.Command("node", tmpName)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
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
	if exitCode == 0 {
		fmt.Printf("\033[34mPASS\033[0m %s \n", filePath)
	} else {
		fmt.Printf("\033[31mFAIL\033[0m %s \n", filePath)
		log.Println("\033[31m", stderr, "\033[0m")
	}
}
