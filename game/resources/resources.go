package resources

type uiResources struct {
	Fonts *fonts
}

func CreateUIResource() (*uiResources, error) {
	fonts, err := loadFonts()
	if err != nil {
		return nil, err
	}

	return &uiResources{
		Fonts: fonts,
	}, nil
}

func (u *uiResources) close() {
	u.Fonts.close()
}
