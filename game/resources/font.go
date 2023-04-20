package resources

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	fontFaceRegular = "assets/fonts/NotoSans-Regular.ttf"
	fontFaceBold    = "assets/fonts/NotoSans-Bold.ttf"
)

type fonts struct {
	Face         font.Face
	TitleFace    font.Face
	BigTitleFace font.Face
	ToolTipFace  font.Face
}

func loadFonts() (*fonts, error) {
	fontFace, err := loadFont(fontFaceRegular, 20)
	if err != nil {
		return nil, err
	}

	titleFontFace, err := loadFont(fontFaceBold, 24)
	if err != nil {
		return nil, err
	}

	bigTitleFontFace, err := loadFont(fontFaceBold, 28)
	if err != nil {
		return nil, err
	}

	toolTipFace, err := loadFont(fontFaceRegular, 15)
	if err != nil {
		return nil, err
	}

	return &fonts{
		Face:         fontFace,
		TitleFace:    titleFontFace,
		BigTitleFace: bigTitleFontFace,
		ToolTipFace:  toolTipFace,
	}, nil
}

func (f *fonts) close() {
	if f.Face != nil {
		_ = f.Face.Close()
	}

	if f.TitleFace != nil {
		_ = f.TitleFace.Close()
	}

	if f.BigTitleFace != nil {
		_ = f.BigTitleFace.Close()
	}
}

func loadFont(path string, size float64) (font.Face, error) {
	fontData, err := embeddedAssets.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
