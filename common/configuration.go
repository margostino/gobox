package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Client struct {
	Url         string `yaml:"url"`
	RequestFile string `yaml:"requestFile"`
	CallsNumber int    `yaml:"callsNumber"`
}

type Server struct {
	Port            string `yaml:"port"`
	Host            string `yaml:"host"`
	Path            string `yaml:"path"`
	ResponseFile    string `yaml:"responseFile"`
	HealthcheckPath string `yaml:"healthcheckPath"`
	HealthcheckFile string `yaml:"healthcheckFile"`
}

type Configuration struct {
	Clients []*Client `yaml:"clients"`
	Servers []*Server `yaml:"servers"`
}

func (s *Server) GetAddress() string {
	return s.Host + ":" + s.Port
}

func GetServerConfig(path string) []*Server {
	return GetConfig(path).Servers
}

func GetClientConfig(path string) []*Client {
	return GetConfig(path).Clients
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
