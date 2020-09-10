package imggen

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// ImgGen is a struct containing relevant config for Image Generation tool
type ImgGen struct {
	ImgSizeMb      int  `yaml:"img-size-mb,omitempty"`
	LayerNumber    int  `yaml:"layer-number,omitempty"`
	GenerateRandom bool `yaml:"generate-random,omitempty"`
}

// Config is the config for image generation
type Config struct {
	ImageGeneration ImgGen `yaml:"image-generation,omitempty"`
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

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

// Generate a docker image out of yaml config file
func Generate() string {
	log.Printf("Loading config file")
	config, _ := loadConfig()
	layerSize := int64((config.ImageGeneration.ImgSizeMb / config.ImageGeneration.LayerNumber) * 1024 * 1024)
	randhex, _ := randomHex(3)
	filepath := strconv.Itoa(config.ImageGeneration.ImgSizeMb) + "-" + strconv.Itoa(config.ImageGeneration.LayerNumber) + "-" + randhex + "/"
	for i := 0; i < config.ImageGeneration.LayerNumber; i++ {
		hexval, _ := randomHex(32)
		fd, err := create(filepath + hexval)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		if config.ImageGeneration.GenerateRandom {
			_, err = fd.Seek(layerSize-9, 0)
			randbytes := make([]byte, 8)
			rand.Read(randbytes)
			_, err = fd.Write(randbytes)
			_, err = fd.Write([]byte{0})
			err = fd.Close()
			if err != nil {
				log.Fatal("Failed to close file")
			}
		}

		digest, err := sha256Digest(filepath + hexval)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Rename(filepath+hexval, filepath+digest)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Docker layer generated")
	}
	return filepath
}
