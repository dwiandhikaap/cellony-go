package input

import (
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ScrollDown input.Action = iota
	ScrollUp
	ZoomIn
)

var Keymap = input.Keymap{
	ScrollDown: {input.KeyWheelDown},
	ScrollUp:   {input.KeyWheelUp},
	ZoomIn:     {input.KeyW},
}

var InputSystem = input.System{}
var Handler = InputSystem.NewHandler(0, Keymap)

func init() {
	InputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
}

func Update() {
	InputSystem.Update()
}
