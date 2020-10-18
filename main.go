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

func main() {
	path := os.Args[1]
	fi, err := os.Stat(path)
	if err != nil {
		log.Print(err)
	}
	if fi.Mode().IsDir() {
		walkDir(path)
	} else {
		processFile(path)
	}
}

func walkDir(rootPath string) {
	path := rootPath + "test/"
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if info.IsDir() == false {
			processFile(path)
		}
		return nil
	})
}

func processFile(path string) {

	data, _ := ioutil.ReadFile(path)
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
		var includeFinal []byte
		assert, err := ioutil.ReadFile(os.Args[2] + "assert.js")
		if err != nil {
			log.Print(err)
		}
		sta, err := ioutil.ReadFile(os.Args[2] + "sta.js")
		if err != nil {
			log.Print(err)
		}
		for _, include := range includes {
			data, _ := ioutil.ReadFile(os.Args[2] + include.(string))
			includeFinal = append(includeFinal, data...)
		}
		includeFinal = append(includeFinal, assert...)
		includeFinal = append(includeFinal, sta...)
		var finalFile []byte
		finalFile = append(finalFile, includeFinal...)
		finalFile = append(finalFile, data...)
		currDir, _ := os.Getwd()
		err = ioutil.WriteFile(currDir+"/tmp.js", finalFile, 0777)
		if err != nil {
			log.Print(err)
		}
		var outbuf, errbuf bytes.Buffer
		var exitCode int
		cmd := exec.Command("node", "tmp.js")
		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf

		err = cmd.Run()
		stdout := outbuf.String()
		stderr := errbuf.String()

		if err != nil {
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
		}
		log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	}
}
