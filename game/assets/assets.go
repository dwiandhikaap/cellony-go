package assets

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

type Assets struct {
	Sprites map[string]*ebiten.Image
}

var AssetsInstance *Assets

func InitializeAssets() error {
	sprites, err := loadSprites()
	if err != nil {
		return err
	}

	AssetsInstance = &Assets{
		Sprites: sprites,
	}

	return err
}

func loadSprites() (map[string]*ebiten.Image, error) {
	sprites := make(map[string]*ebiten.Image)

	// create string to string dict for asset name -> asset path
	// cell -> assets/image/cell.png

	sprite_dicts := map[string]string{
		"cell":     "assets/image/cell.png",
		"circle64": "assets/image/circle-64.png",
		"wall0":    "assets/image/wall-0.png",
		"wall1":    "assets/image/wall-1.png",
		"wall2":    "assets/image/wall-2.png",
		"wall3":    "assets/image/wall-3.png",
		"wall4":    "assets/image/wall-4.png",
		"wall5":    "assets/image/wall-5.png",
		"wall6":    "assets/image/wall-6.png",
		"wall7":    "assets/image/wall-7.png",
		"wall8":    "assets/image/wall-8.png",
		"wall9":    "assets/image/wall-9.png",
	}

	for name, path := range sprite_dicts {
		f, err := os.Open(path)
		if err != nil {
			println("Warning: " + path + " failed to load")
			continue
		}
		defer f.Close()

		image, _, err := image.Decode(f)
		if err != nil {
			return nil, err
		}

		// scale image to 8x8
		if name == "circle64" {
			image = resize.Resize(8, 8, image, resize.NearestNeighbor)
		}

		sprites[name] = ebiten.NewImageFromImage(image)
	}

	return sprites, nil
}
