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
		image = resize.Resize(8, 8, image, resize.NearestNeighbor)

		sprites[name] = ebiten.NewImageFromImage(image)
	}

	return sprites, nil
}
