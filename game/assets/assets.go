package assets

import (
	"autocell/game/config"
	"autocell/game/util"
	"encoding/json"
	"image"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type Assets struct {
	Sprites  map[string]*ebiten.Image
	Textures map[string]*image.Image

	fonts     map[string]map[float64]*font.Face
	_rawFonts map[string]*sfnt.Font
}

func (a *Assets) GetFont(name string, fontSize float64) *font.Face {
	if a.fonts[name][fontSize] == nil {
		face, err := generateFontFace(a._rawFonts[name], fontSize)

		if err != nil {
			panic(err)
		}

		if a.fonts[name] == nil {
			a.fonts[name] = make(map[float64]*font.Face)
		}

		a.fonts[name][fontSize] = face
	}

	return a.fonts[name][fontSize]
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

	sprites, textures, err := loadSprites(&assetsJson)
	if err != nil {
		return err
	}

	rawFonts, err := loadFonts(&assetsJson)
	if err != nil {
		return err
	}

	AssetsInstance = &Assets{
		Sprites:  sprites,
		Textures: textures,

		fonts:     make(map[string]map[float64]*font.Face),
		_rawFonts: rawFonts,
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
	Fonts []string `json:"fonts"`
}

func loadSprites(assets *assetsJson) (map[string]*ebiten.Image, map[string]*image.Image, error) {
	sprites := make(map[string]*ebiten.Image)
	textures := make(map[string]*image.Image)

	for _, img := range assets.Image {
		f, err := os.Open(img.Path)
		if err != nil {
			println("Warning: " + img.Path + " failed to load")
			continue
		}
		defer f.Close()

		image, _, err := image.Decode(f)
		if err != nil {
			return nil, nil, err
		}

		if img.Type == "tile" {
			image = resize.Resize(uint(config.Game.TileSize), uint(config.Game.TileSize), image, resize.NearestNeighbor)
		} else if img.Size.Height > 0 && img.Size.Width > 0 {
			image = resize.Resize(uint(img.Size.Width), uint(img.Size.Height), image, resize.NearestNeighbor)
		}

		textures[img.Name] = &image
		sprites[img.Name] = ebiten.NewImageFromImage(image)
	}

	return sprites, textures, nil
}

func loadFonts(assets *assetsJson) (map[string]*sfnt.Font, error) {
	rawFonts := make(map[string]*sfnt.Font)

	for _, path := range assets.Fonts {
		name := util.FilePathToName(path)

		fontBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		font, err := opentype.Parse(fontBytes)
		if err != nil {
			return nil, err
		}

		rawFonts[name] = font
	}

	return rawFonts, nil
}

func generateFontFace(rawFont *sfnt.Font, fontSize float64) (*font.Face, error) {
	face, err := opentype.NewFace(rawFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return &face, err
}
