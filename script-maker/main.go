package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg *sync.WaitGroup

func main() {
	outputFile := "/some/resources.txt"
	folder := "/some"
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	count := 0
	baseScript := "./bash-script action base-folder"
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if strings.HasSuffix(info.Name(), ".yml") {
			config := GetConfig(fmt.Sprintf("%s/%s", folder, info.Name()))
			for environment := range config.Application.Environments {
				partition := strings.Split(environment, "-")[0]
				stage := strings.Split(environment, "-")[1]
				line := fmt.Sprintf("%s%s %s %s\n", baseScript, info.Name(), partition, stage)
				_, err2 := f.WriteString(line)

				if err2 != nil {
					log.Fatal(err2)
				}
			}
			count += 1
		}
		return nil
	})
	fmt.Printf("Total: %d", count)
}

type Application struct {
	Environments map[string]interface{} `yaml:"environments"`
}

type Configuration struct {
	Application Application `yaml:"application"`
}

func GetConfig(file string) *Configuration {
	var configuration Configuration
	ymlFile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &configuration
}
