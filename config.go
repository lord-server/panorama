package main

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type RegionConfig struct {
	XBounds    [2]int `toml:"x_bounds"`
	YBounds    [2]int `toml:"y_bounds"`
	UpperLimit int    `toml:"upper_limit"`
	LowerLimit int    `toml:"lower_limit"`
}

type Config struct {
	ListenAddress string       `toml:"listen_address"`
	GamePath      string       `toml:"game_path"`
	WorldPath     string       `toml:"world_path"`
	WorldDSN      string       `toml:"world_dsn"`
	RegionConfig  RegionConfig `toml:"region"`
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
