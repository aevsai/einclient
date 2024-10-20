package engine

import (
	"testing"

	"github.com/fogleman/gg"
)

func TestLoadSchema(t *testing.T) {
	scenePath := "../scenes/ein.yml"
	ch := make(chan *Scene, 1)
	err := LoadScene(scenePath, ch)
	if err != nil {
		t.Errorf("Failed to load scene: %v", err)
	}
	select {
	case scene := <-ch:
		if scene == nil {
			t.Errorf("Expected a scene but got nil")
		}
	default:
		t.Errorf("No scene was loaded")
	}
}

func TestRender(t *testing.T) {
	scenePath := "../scenes/ein.yml"
	ch := make(chan *Scene, 1)
	err := LoadScene(scenePath, ch)
	if err != nil {
		t.Errorf("Failed to load scene: %v", err)
	}
	select {
	case scene := <-ch:
		if scene == nil {
			t.Errorf("Expected a scene but got nil")
		}
		ctx := gg.NewContext(scene.Frame.Width, scene.Frame.Height)
		scene.Render(ctx)
	default:
		t.Errorf("No scene was loaded")
	}
}
