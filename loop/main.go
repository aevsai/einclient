package loop

import (
	"einclient/engine"
	"einclient/rgbmatrix"
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	"time"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

const (
	Height = 64
	Width  = 64
)

var (
	gifs_repo = flag.String("gifs", "./gifs", "directory containing GIFs to play")
	gif_delay = flag.Int("gif-delay", 10, "delay between GIFs in milliseconds")
	no_resize = flag.Bool("no-resize", false, "play GIFs without resizing")
)

func init() {
	flag.Parse()
}

type Loop struct {
	Matrix    rgbmatrix.Matrix
	Animation *Animation
	Toolkit   *rgbmatrix.ToolKit
	Chan      chan *engine.Scene
}

func NewLoop(ch chan *engine.Scene) (*Loop, error) {
	config := &rgbmatrix.DefaultConfig
	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		return nil, err
	}
	return &Loop{
		Matrix:    m,
		Animation: NewAnimation(*<-ch),
		Toolkit:   rgbmatrix.NewToolKit(m),
		Chan:      ch,
	}, nil
}

func (l *Loop) Start() error {
	var err error
	var i image.Image
	var n <-chan time.Time
	fmt.Println("Starting loop")
	for {
		select {
		case scene := <-l.Chan:
			l.Animation = NewAnimation(*scene)
		default:
		}

		i, n, err = l.Animation.Next()
		if err != nil {
			break
		}

		if err := l.Toolkit.PlayImageUntil(i, n); err != nil {
			return err
		}
	}

	if err == io.EOF {
		return nil
	}

	return err
}

func (l *Loop) Stop() error {
	return l.Toolkit.Close()
}

type Animation struct {
	ctx   *gg.Context
	scene *engine.Scene
}

func NewAnimation(scene engine.Scene) *Animation {
	return &Animation{
		ctx:   gg.NewContext(scene.Frame.Width, scene.Frame.Height),
		scene: &scene,
	}
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	a.ctx.SetColor(color.Black)
	a.ctx.Clear()
	a.scene.Render(a.ctx)
	img := resize.Resize(64, 64, a.ctx.Image(), resize.Lanczos2)
	return img, time.After(time.Millisecond * 50), nil
}
