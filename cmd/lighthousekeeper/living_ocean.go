package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"time"
)

// Exactly five reusable impact locations. These are visual-only. All stone is
// independently painted and immutable, irrespective of splash timing.
type impact struct {
	x, y, w, h float64
	pose       int
	phase      float64
	behind     bool
}

var rockImpacts = [...]impact{
	{243, 553, 169, 164, 0, .22, true},
	{292, 570, 180, 167, 1, .71, false},
	{167, 515, 132, 139, 0, .52, false},
	{594, 685, 146, 118, 1, .11, false},
	{338, 611, 130, 95, 0, .82, true},
}

func (a *app) initWater() error {
	shader, err := ebiten.NewShader(waterKage)
	if err != nil {
		return fmt.Errorf("water shader: %w", err)
	}
	a.waterShader = shader
	a.waterOptions.Images[0] = a.images["seascape.jpg"]
	a.waterOptions.Images[1] = a.images["water-domain.png"]
	a.waterOptions.Uniforms = map[string]interface{}{"Clock": float32(0)}
	src := a.images["impact-atlas.png"]
	for j := 0; j < 10; j++ {
		col, row := j%5, j/5
		a.impacts[j] = src.SubImage(imageRect(col*128, row*128, 128, 128)).(*ebiten.Image)
	}
	return nil
}

// Minimal rectangle helper permits static atlas slicing with no image allocs.
func imageRect(x, y, w, h int) image.Rectangle { return image.Rect(x, y, x+w, y+h) }

func (a *app) drawLivingWater(dst *ebiten.Image) {
	if a.waterShader == nil {
		dst.DrawImage(a.images["seascape.jpg"], nil)
		return
	}
	a.waterOptions.Uniforms["Clock"] = float32(float64(time.Now().UnixNano()) / 1e9)
	dst.DrawRectShader(1920, 1440, a.waterShader, &a.waterOptions)
}

// One tiny table and ten source sprites. Draw behind or in front of stone by
// scene order, NOT by trying to infer object boundaries from RGB thresholds.
func (a *app) drawImpacts(dst *ebiten.Image, behind bool, sec float64) {
	for _, s := range rockImpacts {
		if s.behind != behind {
			continue
		}
		a.drawOneImpact(dst, s, sec)
	}
}
func (a *app) drawOneImpact(dst *ebiten.Image, s impact, sec float64) {
	period := 2.25 + float64(s.pose)*.16
	cycle := math.Mod(sec/period+s.phase+20, 1)
	if cycle < .23 {
		return
	}
	phase := (cycle - .23) / .77
	f := phase * 5
	at := int(f)
	mix := f - float64(at)
	alpha := .82 * math.Pow(math.Max(0, math.Sin(phase*math.Pi)), .65)
	rise := 22 * math.Sin(phase*math.Pi)
	x := s.x + math.Sin(sec*.51+s.phase*8)*2
	y := s.y - rise
	for q := 0; q < 2; q++ {
		k := min(4, at+q) + 5*s.pose
		im := a.impacts[k]
		o := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		o.GeoM.Translate(-64, -64)
		o.GeoM.Scale(s.w*renderScale/128, s.h*renderScale/128)
		o.GeoM.Translate(x*renderScale, y*renderScale)
		opacity := alpha * (1 - mix)
		if q == 1 {
			opacity = alpha * mix
		}
		o.ColorScale.ScaleAlpha(float32(opacity))
		dst.DrawImage(im, o)
	}
}

// A boat-relative wake does not move simulation state or create an emitter.
func (a *app) drawWake(dst *ebiten.Image, x, y float64, sec float64) {
	s := impact{x: x - 24, y: y + 14, w: 105, h: 43, pose: 0, phase: .18, behind: false}
	a.drawOneImpact(dst, s, sec)
}
