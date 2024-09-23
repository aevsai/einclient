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

	for _, wrapper := range scene.Objects {
		data, err := json.Marshal(wrapper.Properties)
		if err != nil {
			return nil, err
		}
		switch wrapper.Type {
		case ObjectTypeCircle:
			var object objects.Circle
			if err := json.Unmarshal(data, &object); err != nil {
				return nil, err
			}
		case ObjectTypeRectangle:
			var object objects.Rectangle
			if err := json.Unmarshal(data, &object); err != nil {
				return nil, err
			}
		case ObjectTypeArc:
			var object objects.Arc
			if err := json.Unmarshal(data, &object); err != nil {
				return nil, err
			}
		case ObjectTypeLine:
			var object objects.Line
			if err := json.Unmarshal(data, &object); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown object type: %s", wrapper.Type)
		}

	}
	return &scene, nil
}

func (s *Scene) Render(ctx *gg.Context) error {
	for _, wrapper := range s.Objects {
		wrapper.Object.Render(ctx)
	}
	return nil
}
