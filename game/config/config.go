package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type game struct {
	Resolution []int `json:"resolution"`
	FieldSize  []int `json:"fieldSize"`
}

type control struct {
	CamSpeed    float64 `json:"camSpeed"`
	CamSpeedMul float64 `json:"camSpeedMul"`
}

type config struct {
	Game    game    `json:"game"`
	Control control `json:"control"`
}

var cfg config

var Game game
var Control control

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

	fmt.Println("Successfully loaded config.json")

	return nil
}
