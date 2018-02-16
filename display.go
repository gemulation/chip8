package chip8

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	DisplayWidth       = 64
	DisplayHeight      = 32
	DisplayScaleFactor = 10
)

type Display struct {
	config pixelgl.WindowConfig
	window *pixelgl.Window
}

func NewDisplay() *Display {
	config := pixelgl.WindowConfig{
		Title: "Chip 8",
		Bounds: pixel.R(
			0, 0,
			DisplayWidth*DisplayScaleFactor,
			DisplayHeight*DisplayScaleFactor,
		),
		VSync: true,
	}
	return &Display{config: config}
}

func (d *Display) Init() {
	window, err := pixelgl.NewWindow(d.config)
	if err != nil {
		panic(err)
	}
	d.window = window
}

func (d *Display) Clear() {
	d.window.Clear(colornames.Black)
}

func (d *Display) Update() {
	d.window.Update()
}
