package chip8

import (
	"os"

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
		Bounds: pixel.R(
			0, 0,
			DisplayWidth*DisplayScaleFactor,
			DisplayHeight*DisplayScaleFactor,
		),
		VSync: true,
	}
	return &Display{config: config}
}

func (display *Display) Init() {
	window, err := pixelgl.NewWindow(display.config)
	if err != nil {
		panic(err)
	}
	display.window = window

	go func() {
		for !display.window.Closed() {
		}
		os.Exit(0)
	}()

	go func() {
		for {
			display.Update()
		}
	}()
}

func (display *Display) Clear() {
	for i := 0; i < DisplayWidth*DisplayHeight; i++ {
		display.memory[i] = 0
	}
}

func (display *Display) Update() {
	display.window.Clear(colornames.Greenyellow)
	for x := 0; x < DisplayWidth; x++ {
		for y := 0; y < DisplayHeight; y++ {
			if display.memory[y*DisplayWidth+x] == 1 {
				x := float64(x * DisplayScaleFactor)
				y := float64((DisplayHeight - y) * DisplayScaleFactor)

				cube := imdraw.New(nil)
				cube.Color = colornames.Black
				cube.Push(pixel.V(x, y))
				cube.Push(pixel.V(x, y-DisplayScaleFactor))
				cube.Push(pixel.V(x+DisplayScaleFactor, y-DisplayScaleFactor))
				cube.Push(pixel.V(x+DisplayScaleFactor, y))
				cube.Polygon(0)
				cube.Draw(display.window)
			}
		}
	}
	display.window.Update()
}
