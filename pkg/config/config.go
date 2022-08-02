package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Region struct {
	XBounds [2]int `toml:"x_bounds"`
	YBounds [2]int `toml:"y_bounds"`
	ZBounds [2]int `toml:"z_bounds"`
}

type Web struct {
	ListenAddress string `toml:"listen_address"`
	Title         string `toml:"title"`
}

type Renderer struct {
	Workers    int `toml:"workers"`
	ZoomLevels int `toml:"zoom_levels"`
}

type System struct {
	GamePath  string `toml:"game_path"`
	TilesPath string `toml:"tiles_path"`
	WorldPath string `toml:"world_path"`
	WorldDSN  string `toml:"world_dsn"`
}

type Config struct {
	System   System   `toml:"system"`
	Web      Web      `toml:"web"`
	Renderer Renderer `toml:"renderer"`
	Region   Region   `toml:"region"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return config, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return config, err
	}

	_, err = toml.Decode(string(data), &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
