package config

import (
	"io/ioutil"
	"log"

	"src/gopkg.in/yaml.v2"
)

type Config struct {
	Hits int64 `yaml:"hits"`
	Time int64 `yaml:"time"`
}

func New(filePath string) *Config {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	c := &Config{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
