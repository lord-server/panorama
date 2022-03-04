package main

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type MapConfig struct {
	CenterX int `toml:"center_x"`
	CenterZ int `toml:"center_z"`
	XSize   int `toml:"x_size"`
	ZSize   int `toml:"z_size"`
}

type Config struct {
	GamePath  string `toml:"game_path"`
	WorldPath string `toml:"world_path"`
	MapConfig MapConfig
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
