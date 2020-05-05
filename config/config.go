package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/Yesterday17/pug-backend/models"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/json-iterator/go"
)

type Config struct {
	// Auto loaded keys
	KeyPrivate *rsa.PrivateKey `json:"-"`
	KeyPublic  *rsa.PublicKey  `json:"-"`

	models.ModelSettings

	PublicKeyPath   string `json:"public_key_path"`
	PrivateKeyPath  string `json:"private_key_path"`
	PublicKeyString string `json:"-"`

	CrossOrigin bool   `json:"cross_origin"`
	Listen      string `json:"listen"`
}

func LoadConfig() *Config {
	config := Config{}

	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to read config.json", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Failed to unmarshal json file", err)
	}

	if config.PublicKeyPath != "" && config.PrivateKeyPath != "" {
		// Public key
		key, err := ioutil.ReadFile(config.PublicKeyPath)
		if err != nil {
			log.Fatal("Failed to load public key", err)
		}
		config.PublicKeyString = string(key)
		config.KeyPublic, err = jwt.ParseRSAPublicKeyFromPEM(key)
		if err != nil {
			log.Fatal("Failed to parse public key")
		}

		// Private key
		key, err = ioutil.ReadFile(config.PrivateKeyPath)
		if err != nil {
			log.Fatal("Failed to load private key", err)
		}
		config.KeyPrivate, err = jwt.ParseRSAPrivateKeyFromPEM(key)
		if err != nil {
			log.Fatal("Failed to parse private key")
		}
	} else {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
		bytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		p := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytes,
		})
		_ = ioutil.WriteFile("key/pub.pem", p, 0644)

		bytes = x509.MarshalPKCS1PrivateKey(privateKey)
		p = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: bytes,
		})
		_ = ioutil.WriteFile("key/key.pem", p, 0644)

		config.PrivateKeyPath = "key/key.pem"
		config.PublicKeyPath = "key/pub.pem"
		SaveConfig(&config)
	}

	return &config
}

func SaveConfig(config *Config) {
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal("Failed to marshal config struct", err)
	}

	err = ioutil.WriteFile("config.json", data, 0777)
	if err != nil {
		log.Fatal("Failed to save config file", err)
	}
}
