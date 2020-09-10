package cmd

import (
	"io/ioutil"
	"log"
	"registry_benchmark/imggen"

	"gopkg.in/yaml.v3"
)

// Registry is the struct for single registry config
type Registry struct {
	Platform   string
	ImageURL   string `yaml:"image-url,omitempty"`
	URL        string `yaml:"registry-url,omitempty"`
	Username   string
	Password   string
	Repository string
	AccountID  string `yaml:"account-id,omitempty"`
	Region     string
}

// Config is the configuration for the benchmark
type Config struct {
	Registries       []Registry
	ImageGeneration  imggen.ImgGen `yaml:"image-generation,omitempty"`
	ImageName        string        `yaml:"image-name,omitempty"`
	Iterations       int
	StorageURL       string `yaml:"storage-url,omitempty"`
	PullSourceFolder string `yaml:"pull-source-folder,omitempty"`
}

// LoadConfig is the function for loading configuration from yaml file
func loadConfig() (*Config, error) {
	c := Config{}
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &c, nil
}
