package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type game struct {
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	TileSize float64 `json:"tileSize"`
}

type control struct {
	CamSpeed    float64 `json:"camSpeed"`
	CamSpeedMul float64 `json:"camSpeedMul"`
}

type video struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type config struct {
	Game    game    `json:"game"`
	Control control `json:"control"`
	Video   video   `json:"video"`
}

var cfg config

var Game game
var Control control
var Video video

func LoadConfig() error {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &cfg)

	Game = cfg.Game
	Control = cfg.Control
	Video = cfg.Video

	fmt.Println("Successfully loaded config.json")

	return nil
}
