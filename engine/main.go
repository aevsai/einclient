package engine

import (
	"einclient/engine/objects"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fogleman/gg"
	"gopkg.in/yaml.v2"
)

const (
	ObjectTypeCircle    = "circle"
	ObjectTypeRectangle = "rectangle"
	ObjectTypeArc       = "arc"
	ObjectTypeLine      = "line"
	ObjectTypePolygon   = "polygon"
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

func LoadScene(filePath string) (*Scene, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
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

func (s *Scene) Render(ctx *gg.Context) error {
	for _, wrapper := range s.Objects {
		fmt.Printf("Rendering %v\n", wrapper.Object)
		wrapper.Object.Render(ctx)
	}
	return nil
}
