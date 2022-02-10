package main

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type WorldConfig struct {
	Backend    string
	Connection string
}

type GameConfig struct {
	Path string
	Desc string
}

type Config struct {
	World WorldConfig
	Game  GameConfig
}

func LoadConfig(path string) Config {
	var config Config

	file, err := os.Open(path)
	if err != nil {
		file.Close()
		panic(err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(string(data), &config)
	if err != nil {
		panic(err)
	}

	return config
}
