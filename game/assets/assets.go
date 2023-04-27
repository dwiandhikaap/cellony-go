package assets

import (
	"cellony/game/config"
	"encoding/json"
	"image"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

type Assets struct {
	Sprites map[string]*ebiten.Image
}

var AssetsInstance *Assets

func InitializeAssets() error {
	var assetsJson assetsJson

	jsonFile, err := os.Open("assets.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &assetsJson)

	sprites, err := loadSprites(&assetsJson)
	if err != nil {
		return err
	}

	AssetsInstance = &Assets{
		Sprites: sprites,
	}

	return err
}

type assetsJson struct {
	Image []struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Size struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"size"`
		Type string `json:"type,omitempty"`
	} `json:"image"`
	Audio []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"audio"`
}

func loadSprites(assets *assetsJson) (map[string]*ebiten.Image, error) {
	sprites := make(map[string]*ebiten.Image)

	for _, img := range assets.Image {
		f, err := os.Open(img.Path)
		if err != nil {
			println("Warning: " + img.Path + " failed to load")
			continue
		}
		defer f.Close()

		image, _, err := image.Decode(f)
		if err != nil {
			return nil, err
		}

		if img.Type == "tile" {
			image = resize.Resize(uint(config.Game.TileSize), uint(config.Game.TileSize), image, resize.NearestNeighbor)
		} else if img.Size.Height > 0 && img.Size.Width > 0 {
			image = resize.Resize(uint(img.Size.Width), uint(img.Size.Height), image, resize.NearestNeighbor)
		}

		sprites[img.Name] = ebiten.NewImageFromImage(image)
	}

	return sprites, nil
}
