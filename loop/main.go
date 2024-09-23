package loop

import (
	"einclient/engine"
	"einclient/rgbmatrix"
	"flag"
	"image"
	"time"

	"github.com/fogleman/gg"
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
}

func NewLoop(s *engine.Scene) (*Loop, error) {
	config := &rgbmatrix.DefaultConfig
	m, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		return nil, err
	}
	return &Loop{
		Matrix:    m,
		Animation: NewAnimation(*s),
	}, nil
}

func (l *Loop) Start() error {
	// Loop algorithm
	// 1. Render the scene
	// 2. Display the scene
	// 3. Call the next frame
	// 4. Repeat

	tk := rgbmatrix.NewToolKit(l.Matrix)
	defer tk.Close()

	// gifs, err := loadGIFs(*gifs_repo, *no_resize)
	// if err != nil {
	// 	return err
	// }

	tk.PlayAnimation(l.Animation)
	return nil
}

type Animation struct {
	ctx   *gg.Context
	scene *engine.Scene
}

func NewAnimation(scene engine.Scene) *Animation {
	return &Animation{
		ctx:   gg.NewContext(64, 64),
		scene: &scene,
	}
}

func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
	a.scene.Render(a.ctx)
	return a.ctx.Image(), time.After(time.Millisecond * 50), nil
}

// func (a *Animation) Next() (image.Image, <-chan time.Time, error) {
// 	gif := a.gifs[a.idx]
// 	img := gif.Image[a.frame]
// 	delay := gif.Delay[a.frame]
// 	a.frame = (a.frame + 1) % len(gif.Image)
// 	if a.frame == 0 {
// 		a.idx = (a.idx + 1) % len(a.gifs)
// 	}
// 	return img, time.After(time.Duration(delay) * time.Duration(*gif_delay) * time.Millisecond), nil
// }

// func loadGIFs(dir string, noResize bool) ([]*gif.GIF, error) {
// 	var gifs []*gif.GIF
// 	scaledDir := filepath.Join(dir, ".scaled")
// 	if !noResize {
// 		if err := os.MkdirAll(scaledDir, os.ModePerm); err != nil {
// 			return nil, err
// 		}
// 	}

// 	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if filepath.Ext(path) == ".gif" {
// 			scaledPath := filepath.Join(scaledDir, info.Name())
// 			if !noResize {
// 				if _, err := os.Stat(scaledPath); err == nil {
// 					f, err := os.Open(scaledPath)
// 					if err != nil {
// 						return err
// 					}
// 					defer f.Close()
// 					g, err := gif.DecodeAll(f)
// 					if err != nil {
// 						return err
// 					}
// 					gifs = append(gifs, g)
// 					return nil
// 				}
// 			}

// 			f, err := os.Open(path)
// 			if err != nil {
// 				return err
// 			}
// 			defer f.Close()
// 			g, err := gif.DecodeAll(f)
// 			if err != nil {
// 				return err
// 			}
// 			if !noResize {
// 				for i, img := range g.Image {
// 					scaledImg := image.NewRGBA(image.Rect(0, 0, 64, 64))
// 					draw.CatmullRom.Scale(scaledImg, scaledImg.Rect, img, img.Bounds(), draw.Over, nil)
// 					g.Image[i] = scaledImg
// 				}

// 				sf, err := os.Create(scaledPath)
// 				if err != nil {
// 					return err
// 				}
// 				defer sf.Close()
// 				if err := gif.EncodeAll(sf, g); err != nil {
// 					return err
// 				}
// 			} else {
// 				for _, img := range g.Image {
// 					convertToRGB(img)
// 				}
// 			}

// 			gifs = append(gifs, g)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return gifs, nil
// }

// func convertToRGB(img *image.Paletted) {
// 	for i := 0; i < len(img.Palette); i++ {
// 		c := img.Palette[i]
// 		r, g, b, _ := c.RGBA()
// 		img.Palette[i] = color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
// 	}
// 	for y := 0; y < img.Rect.Dy(); y++ {
// 		for x := 0; x < img.Rect.Dx(); x++ {
// 			if img.ColorIndexAt(x, y) == 0 {
// 				img.SetColorIndex(x, y, 1)
// 			}
// 		}
// 	}
// }
