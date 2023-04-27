package generator

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

func GenerateTextureShades(path string, shades int) error {
	f, err := os.Open(path)
	if err != nil {
		println("Error: " + path + " failed to load")
		return nil
	}
	defer f.Close()

	filename := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	filename = strings.Split(filename, ".")[0]

	filepath := path[:len(path)-len(filename)-len(".png")]

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	img = imaging.Resize(img, 10, 10, imaging.NearestNeighbor)

	shadeIndex := shades - 1
	newFile, err := os.Create(fmt.Sprintf("%s/%s-%d.png", filepath, filename, shadeIndex))

	if err != nil {
		return err
	}

	png.Encode(newFile, img)

	for i := shades - 1; i > 0; i-- {
		newImg := imaging.AdjustBrightness(img, -float64(100/(shades+1)*(shades-i)))
		shadeIndex--
		newFile, err := os.Create(fmt.Sprintf("%s/%s-%d.png", filepath, filename, shadeIndex))

		if err != nil {
			return err
		}

		png.Encode(newFile, newImg)
	}

	return nil
}
