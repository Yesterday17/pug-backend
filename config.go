package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Yesterday17/pug-backend/models"
	_ "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
)

type Config struct {
	// Auto generated per startup
	// It means user login would break after server restart
	Secret string

	// Auto loaded keys
	PublicKey  []byte
	PrivateKey []byte

	models.ModelSettings

	PublicKeyPath  string `json:"public_key_path"`
	PrivateKeyPath string `json:"private_key_path"`
}

func LoadConfig() *Config {
	config := Config{
		Secret: uuid.NewV4().String(),
	}

	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to read config.json", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Failed to unmarshal json file", err)
	}
	return &config
}
