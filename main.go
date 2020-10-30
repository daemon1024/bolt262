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

var count, pass, fail int64 = 0, 0, 0
var includePath string

func main() {
	path := os.Args[1]
	includePath = os.Args[2]
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
	// includePath = "../test262/harness/" //hardcoded for bench
	assert, err := ioutil.ReadFile(includePath + "assert.js")
	if err != nil {
		log.Print(err)
	}
	sta, err := ioutil.ReadFile(includePath + "sta.js")
	if err != nil {
		log.Print(err)
	}
	var wg sync.WaitGroup
	// filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		log.Fatalf(err.Error())
	// 	}
	// 	if info.IsDir() == false {
	// 		wg.Add(1)
	// 		go func(path string) {
	// 			defer wg.Done()
	// 			processFile(path, assert, sta)
	// 		}(path)
	// 		fmt.Printf("\033[36mTotal files %d\n\033[0m\033[31mFailed Tests %d\n\033[0m\033[34mPassed Tests %d\033[0m", count, fail, pass)

	// 	}
	// 	return nil
	// })
	err = godirwalk.Walk(path, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {

			if de.IsDir() == false {
				wg.Add(1)
				go func(path string) {
					//note to self : check file directive limit
					defer wg.Done()
					processFile(osPathname, assert, sta)
				}(path)
				fmt.Printf("\033[36mTotal files %d\n\033[0m\033[31mFailed Tests %d\n\033[0m\033[34mPassed Tests %d\033[0m", count, fail, pass)
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
	// filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		log.Fatalf(err.Error())
	// 	}
	// 	if info.IsDir() == false {
	// 		processFile(path, assert, sta)
	// 	}
	// 	return nil
	// })
	wg.Wait()
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
			data, _ := ioutil.ReadFile(includePath + include.(string))
			includeFinal = append(includeFinal, data...)
		}
	}
	var finalFile []byte
	finalFile = append(finalFile, includeFinal...)
	finalFile = append(finalFile, data...)
	currDir, _ := os.Getwd()
	// err = ioutil.WriteFile(currDir+"/tmp.js", finalFile, 0777)
	// if err != nil {
	// 	log.Print(err)
	// }
	tmpFile, err := ioutil.TempFile(currDir+"/tmp/", "tmptest-")
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
