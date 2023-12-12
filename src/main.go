package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	configFileName    = "config.json"
	defaultKeyMessage = "put a valid API key from Openrouter here"
)

func main() {

	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	op := ebiten.RunGameOptions{}
	op.ScreenTransparent = true
	op.SkipTaskbar = true
	game := NewGame(config)
	if err := ebiten.RunGameWithOptions(game, &op); err != nil {
		log.Fatal(err)
		return
	}
}

func loadConfig() (*Config, error) {
	file, err := os.Open(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Config file does not exist. Creating a new one.")
			return createConfig()
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func createConfig() (*Config, error) {
	config := Config{API_Key: defaultKeyMessage}

	err := saveConfig(config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config Config) error {
	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
