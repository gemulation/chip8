package chip8

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	DisplayWidth       = 64
	DisplayHeight      = 32
	DisplayScaleFactor = 20
)

type Display struct {
	config pixelgl.WindowConfig
	window *pixelgl.Window
	memory [DisplayWidth * DisplayHeight]byte
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
	d.window.Clear(colornames.Lightgreen)
}

func (d *Display) Update() {
	for x := 0; x < DisplayWidth; x++ {
		for y := 0; y < DisplayHeight; y++ {
			if d.memory[y*DisplayWidth+x] == 1 {
				x := float64(x * DisplayScaleFactor)
				y := float64(y * DisplayScaleFactor)

				cube := imdraw.New(nil)
				cube.Color = colornames.Black
				cube.Push(pixel.V(x, y))
				cube.Push(pixel.V(x, y+DisplayScaleFactor))
				cube.Push(pixel.V(x+DisplayScaleFactor, y+DisplayScaleFactor))
				cube.Push(pixel.V(x+DisplayScaleFactor, y))
				cube.Polygon(0)
				cube.Draw(d.window)
			}
		}
	}
	d.window.Update()
}
