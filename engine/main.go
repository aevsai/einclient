package engine

import (
	"einclient/engine/objects"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fogleman/gg"
	"gopkg.in/yaml.v3"
)

const (
	ObjectTypeCircle        = "circle"
	ObjectTypeRectangle     = "rectangle"
	ObjectTypeArc           = "arc"
	ObjectTypeLine          = "line"
	ObjectTypeSimplePolygon = "simple"
	ObjectTypePolygon       = "polygon"
)

type Scene struct {
	Version string          `yaml:"version"`
	Frame   Frame           `yaml:"frame"`
	Objects []ObjectWrapper `yaml:"objects"`
}

type Frame struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type ObjectWrapper struct {
	Name       string                 `yaml:"name"`
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties"`
	Object     objects.Renderable
}

func unmarshalObject(data []byte, target *objects.Renderable, obj objects.Renderable) error {
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}
	*target = obj
	return nil
}

func LoadScene(filePath string, reloadChan chan *Scene) error {
	load := func() (*Scene, error) {
		var data []byte
		var err error
		for len(data) == 0 {
			data, err = os.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
		}

		if len(data) == 0 {
			return nil, fmt.Errorf("file is empty")
		}
		var scene Scene
		err = yaml.Unmarshal(data, &scene)
		if err != nil {
			return nil, err
		}

		typeConstructorMap := map[string]func() objects.Renderable{
			ObjectTypeCircle:    func() objects.Renderable { return new(objects.Circle) },
			ObjectTypeRectangle: func() objects.Renderable { return new(objects.Rectangle) },
			ObjectTypeArc:       func() objects.Renderable { return new(objects.Arc) },
			ObjectTypeLine:      func() objects.Renderable { return new(objects.Line) },
			ObjectTypePolygon:   func() objects.Renderable { return new(objects.Polygon) },
			ObjectTypeSimplePolygon: func() objects.Renderable {
				return new(objects.SimplePolygon)
			},
		}

		for i := 0; i < len(scene.Objects); i++ {
			wrapper := &scene.Objects[i]
			data, err := json.Marshal(wrapper.Properties)
			if err != nil {
				return nil, err
			}

			constructor, ok := typeConstructorMap[wrapper.Type]
			if !ok {
				return nil, fmt.Errorf("unknown object type: %s", wrapper.Type)
			}

			if err := unmarshalObject(data, &wrapper.Object, constructor()); err != nil {
				return nil, err
			}
		}
		return &scene, nil
	}

	scene, err := load()
	if err != nil {
		return err
	}
	reloadChan <- scene

	go func() {
		watcher, err := os.Stat(filePath)
		if err != nil {
			fmt.Println("Error watching file:", err)
			return
		}
		modTime := watcher.ModTime()

		for {
			watcher, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("Error watching file:", err)
				return
			}

			if watcher.ModTime() != modTime {
				scene, err := load()
				if err != nil {
					fmt.Println("Error loading scene:", err)
				} else {
					fmt.Printf("Reloading scene %s\n", filePath)
					fmt.Printf("Loaded %d objects\n", len(scene.Objects))
					reloadChan <- scene
					modTime = watcher.ModTime()
				}
			}
		}
	}()

	return nil
}

func (s *Scene) Render(ctx *gg.Context) error {
	for _, wrapper := range s.Objects {
		wrapper.Object.Render(ctx)
	}
	return nil
}
