package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/smallfish/simpleyaml"
)

func main() {
	rootPath := os.Args[1]
	iterate(rootPath)
}

func iterate(rootPath string) {
	path := rootPath + "test/"
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if info.IsDir() == false {
			data, _ := ioutil.ReadFile(path)
			re, _ := regexp.Compile(`\/\*---([\s\S]*?)---\*\/`)
			y := re.FindSubmatch([]byte(data))
			if y != nil {
				yaml, err := simpleyaml.NewYaml(y[1])
				if err != nil {
					log.Print(path)
					return nil
				}
				// bar := yaml
				// fmt.Printf("Value: %#v\n", bar)
				includes, _ := yaml.Get("includes").Array()
				fmt.Print(includes)
				for _, include := range includes {
					data, _ := ioutil.ReadFile(rootPath + "/harness/" + include.(string))
					fmt.Printf("%s\n", data)

				}

			}
		}
		return nil
	})
}
