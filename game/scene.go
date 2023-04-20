package game

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)

	Open()
	Close()
}

type SceneManager struct {
	scenes            []Scene
	currentSceneIndex int
}

func (s *SceneManager) Update() error {
	if len(s.scenes) == 0 {
		return nil
	}

	return s.scenes[s.currentSceneIndex].Update()
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	if len(s.scenes) == 0 {
		return
	}

	s.scenes[s.currentSceneIndex].Draw(screen)
}

func (s *SceneManager) SelectScene(index int) {
	s.scenes[s.currentSceneIndex].Close()
	s.currentSceneIndex = index
	s.scenes[s.currentSceneIndex].Open()
}
