package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	MongoAddress string `json:"mongoAddress"`
	MQAddress    string `json:"mqAddress"`
	ResourceRoot string `json:"resourceRoot"`
}

var once sync.Once
var config Config

func Get() Config {
	once.Do(func() {
		initConfig()
	})
	return config
}

func initConfig() {
	open, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(open).Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	validateConfig(config)
}

func validateConfig(c Config) {
	if len(c.MongoAddress) == 0 {
		log.Fatal("MongoDB address is empty")
	}
	if len(c.ResourceRoot) == 0 {
		log.Fatal("Resource root directory is empty")
	}
	if len(c.MQAddress) == 0 {
		log.Fatal("Message Queue address is empty")
	}
}
