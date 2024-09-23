package objects

import "github.com/fogleman/gg"

type Renderable interface {
	Render(ctx *gg.Context)
}

type Circle struct {
	X      float64
	Y      float64
	Radius float64
}

func (c Circle) Render(ctx *gg.Context) {
	ctx.DrawCircle(c.X, c.Y, c.Radius)
}

type Rectangle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func (r Rectangle) Render(ctx *gg.Context) {
	ctx.DrawRectangle(r.X, r.Y, r.Width, r.Height)
}

type Line struct {
	X1 float64
	Y1 float64
	X2 float64
	Y2 float64
}

func (l Line) Render(ctx *gg.Context) {
	ctx.DrawLine(l.X1, l.Y1, l.X2, l.Y2)
}

type Arc struct {
	X      float64
	Y      float64
	Radius float64
	Start  float64
	End    float64
}

func (a Arc) Render(ctx *gg.Context) {
	ctx.DrawArc(a.X, a.Y, a.Radius, a.Start, a.End)
}
