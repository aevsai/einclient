package engine

import (
	"einclient/engine/objects"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/expr-lang/expr"
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
	Version    string                 `yaml:"version"`
	Env        map[string]interface{} `yaml:"env"`
	Frame      Frame                  `yaml:"frame"`
	Objects    []ObjectWrapper        `yaml:"objects"`
	Animations []AnimationWrapper     `yaml:"animations"`
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

type AnimationWrapper struct {
	Name      string            `yaml:"name"`
	Duration  string            `yaml:"duration"`
	Repeat    string            `yaml:"repeat"`
	Delay     string            `yaml:"delay"`
	Keyframes []KeyframeWrapper `yaml:"keyframes"`
	PlayedAt  int64
}

type KeyframeWrapper struct {
	Time       float64                `yaml:"time"`
	Properties map[string]interface{} `yaml:"properties"`
}

func (a *AnimationWrapper) getKeyframe(t int64) (int, *KeyframeWrapper) {
	for idx, keyframe := range a.Keyframes {
		if t-a.PlayedAt <= int64(keyframe.Time*float64(time.Second)) {
			return idx, &keyframe
		}
	}
	return -1, nil
}
func (scene *Scene) ComputeAnimations() error {
	var t int64 = time.Now().UnixNano()
	for _, animation := range scene.Animations {
		duration, err := EvaluateExpression(animation.Duration, scene.Env)
		if err != nil {
			continue
		}
		repeat, err := EvaluateExpression(animation.Repeat, scene.Env)
		if err != nil {
			continue
		}
		delay, err := EvaluateExpression(animation.Delay, scene.Env)
		if err != nil {
			continue
		}
		startAt := animation.PlayedAt + int64(duration.(float64)*float64(time.Second)) + int64(delay.(float64)*float64(time.Second))
		if animation.PlayedAt == 0 || (repeat == true && startAt < t) {
			animation.PlayedAt = t
		}
		idx, kf := animation.getKeyframe(t)
		if kf == nil {
			continue
		}
		prevKeyframe := animation.Keyframes[idx-1]
		for key, targetValue := range kf.Properties {
			var value interface{}
			if prevKeyframe.Properties[key] != nil {
				prevValue := prevKeyframe.Properties[key]
				if prevValue != nil {
					if targetValue == nil {
						targetValue = prevValue
					}
					prevDurations := time.Duration(0)
					for _, pkf := range animation.Keyframes[:idx] {
						prevDurations += time.Duration(pkf.Time * float64(time.Second))
					}
					switch prevValueTyped := prevValue.(type) {
					case float64:
						kfDuration := time.Duration((kf.Time - prevKeyframe.Time) * float64(time.Second))
						elapsed := time.Duration(t - animation.PlayedAt - prevDurations.Nanoseconds())
						progress := float64(elapsed) / float64(kfDuration)
						value = prevValueTyped + progress*(targetValue.(float64)-prevValueTyped)
					case int:
						kfDuration := time.Duration((kf.Time - prevKeyframe.Time) * float64(time.Second))
						elapsed := time.Duration(t - animation.PlayedAt - prevDurations.Nanoseconds())
						progress := float64(elapsed) / float64(kfDuration)
						prevValueFloat := float64(prevValueTyped)
						targetValueFloat := float64(targetValue.(int))
						value = int(prevValueFloat + progress*(targetValueFloat-prevValueFloat))
					case string:
						// String interpolation isn't straightforward; leaving this as-is or handling it differently
						value = targetValue.(string)
					}
				}
			} else {
				value = targetValue
			}
			scene.Env[key] = value
		}
	}
	return nil
}

func (wrapper *ObjectWrapper) Render(ctx *gg.Context, env map[string]interface{}) error {

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
	computedProperties, err := Process(wrapper.Properties, env)
	data, err := json.Marshal(computedProperties)
	if err != nil {
		return err
	}

	constructor, ok := typeConstructorMap[wrapper.Type]
	if !ok {
		return fmt.Errorf("unknown object type: %s", wrapper.Type)
	}

	if err := unmarshalObject(data, &wrapper.Object, constructor()); err != nil {
		return err
	}
	wrapper.Object.Render(ctx)
	return nil
}

func unmarshalObject(data []byte, target *objects.Renderable, obj objects.Renderable) error {
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}
	*target = obj
	return nil
}

func EvaluateExpression(expression string, variables map[string]interface{}) (interface{}, error) {
	program, err := expr.Compile(expression)
	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, variables)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Process(obj map[string]interface{}, env map[string]interface{}) (map[string]interface{}, error) {
	var err error
	var res interface{}
	computed := make(map[string]interface{})
	for key, value := range obj {
		switch key {
		case "points":
			var output []interface{}
			for _, item := range value.([]interface{}) {
				pItem, err := Process(item.(map[string]interface{}), env)
				if err != nil {
					fmt.Printf("Error processing points %v\n", item)
					return nil, err
				}
				output = append(output, pItem)
			}
			res = output
			break
		case "startPoint", "endPoint":
			res, err = Process(value.(map[string]interface{}), env)
			if err != nil {
				return nil, err
			}
			break
		default:
			res, err = EvaluateExpression(value.(string), env)
			if err != nil {
				fmt.Printf("Error processing expression %v\n", err.Error())
				return nil, err
			}
			break
		}
		computed[key] = res
	}
	// fmt.Printf("Processed object %v\n", obj)
	return computed, nil
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
	s.ComputeAnimations()
	for _, wrapper := range s.Objects {
		wrapper.Render(ctx, s.Env)
	}
	return nil
}
