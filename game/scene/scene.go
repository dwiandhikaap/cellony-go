package scene

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)

	Open()
	Close()
}

type SceneManager struct {
	Scenes            []Scene
	CurrentSceneIndex int
}

func (s *SceneManager) Update() error {
	if len(s.Scenes) == 0 {
		return nil
	}

	return s.Scenes[s.CurrentSceneIndex].Update()
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	if len(s.Scenes) == 0 {
		return
	}

	s.Scenes[s.CurrentSceneIndex].Draw(screen)
}

func (s *SceneManager) SelectScene(index int) {
	s.Scenes[s.CurrentSceneIndex].Close()
	s.CurrentSceneIndex = index
	s.Scenes[s.CurrentSceneIndex].Open()
}
