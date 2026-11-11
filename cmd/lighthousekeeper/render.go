package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jonathanselby/lighthouse-keeper/internal/game"
	"golang.org/x/image/font"
)

var (
	ivory = color.RGBA{249, 238, 204, 255}
	amber = color.RGBA{255, 206, 122, 255}
	muted = color.RGBA{182, 207, 217, 255}
	dark  = color.RGBA{10, 27, 45, 239}
	ink   = color.RGBA{57, 49, 46, 255}
)

func (a *app) sprite(dst *ebiten.Image, name string, x, y, scale float64) {
	a.imageAt(dst, a.images[name], x, y, scale)
}
func (a *app) imageAt(dst *ebiten.Image, img *ebiten.Image, x, y, scale float64) {
	if img == nil {
		return
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x*renderScale, y*renderScale)
	dst.DrawImage(img, op)
}
func (a *app) drawText(dst *ebiten.Image, s string, x, y float64, face font.Face, c color.Color, center bool) {
	sx := int(x * renderScale)
	if center {
		sx -= font.MeasureString(face, s).Round() / 2
	}
	text.Draw(dst, s, face, sx, int(y*renderScale), c)
}
func (a *app) headingText(dst *ebiten.Image, s string, x, y float64) {
	a.drawText(dst, s, x+1.5, y+2, a.heading, dark, true)
	a.drawText(dst, s, x, y, a.heading, ivory, true)
}
func (a *app) subText(dst *ebiten.Image, s string, x, y float64) {
	a.drawText(dst, s, x, y, a.font, ivory, true)
}
func (a *app) smallText(dst *ebiten.Image, s string, x, y float64) {
	a.drawText(dst, s, x, y, a.small, ivory, true)
}
func (a *app) panel(dst *ebiten.Image, x, y, w, h float64) {
	ebitenutil.DrawRect(dst, x*2, y*2, w*2, h*2, color.RGBA{9, 27, 43, 220})
	vector.StrokeRect(dst, float32(x*2+3), float32(y*2+3), float32(w*2-6), float32(h*2-6), 2, color.RGBA{233, 185, 113, 195}, true)
}
func (a *app) lamp(dst *ebiten.Image) {
	img := a.images["lighthouse.png"]
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-205, -95)
	op.GeoM.Translate(game.LampX*2, game.LampY*2)
	dst.DrawImage(img, op)
	// A thin halo over the known lamp centre, not gameplay illumination.
	rad := float32(16 + 3*math.Sin(float64(time.Now().UnixNano())/1e9*1.2))
	vector.DrawFilledCircle(dst, float32(game.LampX*2), float32(game.LampY*2), rad, color.RGBA{255, 226, 151, 95}, true)
}
func (a *app) spotlight(dst *ebiten.Image, angle, halfWidth float64) {
	img := a.images["beam.png"]
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-245, -850)
	op.GeoM.Scale(1.44*(halfWidth/(12*math.Pi/180)), 1.92)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(game.LampX*2, game.LampY*2)
	dst.DrawImage(img, op)
}
func (a *app) drawSea(dst *ebiten.Image) {
	dst.DrawImage(a.images["seascape.jpg"], nil)
	// Low-cost, very subtle water motion layered over the static painted sea.
	seconds := float64(time.Now().UnixNano()) / 1e9
	for i := 0; i < 15; i++ {
		y := 445 + float64(i*13)
		x := math.Mod(float64(i*149)+seconds*8, 1110) - 90
		alpha := uint8(23 + i%4*8)
		ebitenutil.DrawLine(dst, x*2, y*2, (x+65)*2, (y+math.Sin(seconds+float64(i))*1.2)*2, color.RGBA{190, 226, 226, alpha})
	}
}
func (a *app) drawFrame(dst *ebiten.Image) {
	switch a.game.Screen {
	case game.Title:
		dst.DrawImage(a.images["title_art.jpg"], nil)
		a.drawTitle(dst)
	case game.NightScreen:
		a.drawNight(dst)
	case game.Predict:
		a.drawSea(dst)
		a.lamp(dst)
		a.drawPredict(dst)
	case game.Harbour:
		a.drawSea(dst)
		a.lamp(dst)
		a.drawHarbour(dst)
	case game.Book:
		a.drawSea(dst)
		a.drawBook(dst)
	case game.Log:
		a.drawSea(dst)
		a.drawLog(dst)
	case game.Rest:
		a.drawSea(dst)
		a.lamp(dst)
		a.drawRest(dst)
	case game.BellSetup:
		a.drawSea(dst)
		a.lamp(dst)
		a.drawBellSetup(dst)
	case game.BellScreen:
		a.drawSea(dst)
		a.lamp(dst)
		a.drawBell(dst)
	}
	if a.saveMessage != "" {
		a.panel(dst, 36, 660, 888, 46)
		a.drawText(dst, a.saveMessage, 480, 689, a.small, amber, true)
	}
}
func (a *app) drawTitle(dst *ebiten.Image) {
	a.panel(dst, 125, 70, 710, 190)
	a.headingText(dst, "LIGHTHOUSE", 480, 145)
	a.headingText(dst, "KEEPER", 480, 205)
	a.drawText(dst, "NOTICE IT  •  NAME IT  •  COME BACK", 480, 245, a.small, amber, true)
	a.panel(dst, 232, 305, 496, 226)
	a.subText(dst, "SPACE    Start a night", 480, 356)
	a.subText(dst, "B    Sightings Book", 480, 404)
	a.subText(dst, "L    Harbour Log", 480, 452)
	a.subText(dst, "K    Keeper's Bell", 480, 501)
	a.smallText(dst, "A small, gentle game about finding your way back", 480, 610)
}
func (a *app) drawPredict(dst *ebiten.Image) {
	a.panel(dst, 123, 135, 714, 425)
	a.headingText(dst, "A NIGHT AT SEA", 480, 211)
	a.subText(dst, "How many lanterns will glow?", 480, 276)
	for i := 0; i < 5; i++ {
		x := 264 + float64(i)*105
		lit := i < a.game.Night.Prediction
		col := color.RGBA{95, 123, 142, 190}
		if lit {
			col = amber
		}
		vector.DrawFilledCircle(dst, float32(x*2), float32(380*2), float32(24*2), col, true)
		vector.StrokeCircle(dst, float32(x*2), float32(380*2), float32(26*2), 3, ivory, true)
	}
	a.subText(dst, fmt.Sprintf("%d of 5", a.game.Night.Prediction), 480, 441)
	a.smallText(dst, "← / → choose    •    SPACE begin", 480, 512)
}
func (a *app) drawNight(dst *ebiten.Image) {
	a.drawSea(dst)
	g := a.game
	n := &g.Night
	dt := a.subTick()
	angle := g.Beam.Angle
	boat := n.Boat
	if !n.CheckIn && n.PauseTicks == 0 {
		angle = math.Max(-game.MaxBeamAngle, math.Min(game.MaxBeamAngle, angle+g.Beam.Omega*dt))
		boat.Advance((1 - 0.65*n.Fog) * boat.Length / 5400 * dt)
	}
	a.spotlight(dst, angle, g.Beam.HalfWidth)
	for i, s := range n.Sightings {
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		im := a.sightseeing[s.Kind]
		wobble := math.Sin(float64(time.Now().UnixNano())/1e9*0.96+float64(i)) * 3
		op.GeoM.Translate(-75, -165)
		op.GeoM.Scale(0.78, 0.78)
		op.GeoM.Translate(s.Pos.X*2, s.Pos.Y*2+wobble)
		opacity := float32(1)
		if s.State == game.Appearing {
			opacity = float32(math.Min(1, float64(s.Age)/24))
		}
		if s.State == game.Leaving || s.State == game.Chased {
			opacity = float32(math.Max(0, 1-float64(s.Age)/30))
		}
		if s.State == game.Noticed {
			opacity = float32(.80 + .2*math.Sin(float64(time.Now().UnixNano())/1e9*7.2))
		}
		op.ColorScale.ScaleAlpha(opacity)
		dst.DrawImage(im, op)
		if s.State == game.Noticed {
			vector.StrokeCircle(dst, float32(s.Pos.X*2), float32(s.Pos.Y*2), float32(38+math.Sin(float64(time.Now().UnixNano())/1e9*6)*6), 2, amber, true)
		}
	}
	p := boat.Position()
	if n.Fog > 0.01 {
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		im := a.images["fog.png"]
		op.GeoM.Translate(-180, -115)
		op.GeoM.Scale(.9+n.Fog*.5, .7+n.Fog*.45)
		op.GeoM.Translate(p.X*2, p.Y*2)
		op.ColorScale.ScaleAlpha(float32(n.Fog * .85))
		dst.DrawImage(im, op)
	}
	a.imageAt(dst, a.images["boat.png"], p.X, p.Y, 0.73)
	if boat.Lit {
		vector.DrawFilledCircle(dst, float32(p.X*2+15), float32(p.Y*2-16), float32(6+math.Sin(float64(time.Now().UnixNano())/1e9*7.2)*2), amber, true)
	}
	a.lamp(dst)
	// Progress is conveyed by five accumulating amber harbour lanterns.
	for i := 0; i < 5; i++ {
		c := color.RGBA{101, 130, 155, 210}
		if n.TotalTicks[i] > 0 && float64(n.LitTicks[i])/float64(n.TotalTicks[i]) >= .7 && (i < n.Crossing || n.Crossing == 4 && boat.T == 1) {
			c = amber
		}
		vector.DrawFilledCircle(dst, float32((730+30*i)*2), float32(596*2), 12, c, true)
	}
	a.panel(dst, 26, 18, 908, 53)
	a.drawText(dst, fmt.Sprintf("Crossing %d of 5", n.Crossing+1), 165, 54, a.small, ivory, true)
	a.drawText(dst, "←  →  guide the beam", 470, 54, a.small, muted, true)
	a.drawText(dst, "N  notice", 780, 54, a.small, amber, true)
	if n.ChimeTicks > 0 {
		a.panel(dst, 310, 92, 340, 56)
		a.drawText(dst, "Where's the boat?", 480, 130, a.font, ivory, true)
	}
	if n.PauseTicks > 0 {
		a.panel(dst, 260, 280, 440, 95)
		a.subText(dst, "A lantern for the harbour", 480, 335)
	}
	if n.CheckIn {
		ebitenutil.DrawRect(dst, 0, 0, 1920, 1440, color.RGBA{5, 20, 39, 158})
		a.panel(dst, 180, 217, 600, 261)
		a.headingText(dst, "A LITTLE CHECK-IN", 480, 290)
		a.subText(dst, "Were you with the boat?", 480, 360)
		a.smallText(dst, "Y  yes     •     N  no", 480, 424)
	}
	for _, p := range a.particles {
		op := float64(p.life-p.age) / float64(p.life)
		col := color.RGBA{255, 218, 142, uint8(220 * op)}
		vector.DrawFilledCircle(dst, float32(p.x*2), float32(p.y*2), float32(3+int(p.hue)), col, true)
	}
}
func (a *app) drawHarbour(dst *ebiten.Image) {
	n := a.game.Night
	a.panel(dst, 128, 110, 704, 495)
	a.headingText(dst, "THE HARBOUR LOG", 480, 190)
	a.subText(dst, "Five little journeys home", 480, 245)
	for i := 0; i < 5; i++ {
		x := 277 + float64(i)*101
		c := color.RGBA{90, 120, 138, 190}
		if n.TotalTicks[i] > 0 && float64(n.LitTicks[i])/float64(n.TotalTicks[i]) >= .7 {
			c = amber
		}
		vector.DrawFilledCircle(dst, float32(x*2), float32(326*2), 40, c, true)
	}
	a.subText(dst, fmt.Sprintf("Your guess: %d    •    Lanterns: %d", n.Prediction, n.Result), 480, 411)
	a.smallText(dst, fmt.Sprintf("Noticed: %d    Came back: %d    Check-ins: %d", n.Noticed, n.CameBacks, len(n.Answers)), 480, 463)
	if n.GoodGuess {
		a.smallText(dst, "★  You made a good guess  ★", 480, 511)
	}
	a.smallText(dst, "SPACE  continue     •     B  book     •     L  log", 480, 570)
}
func (a *app) drawParchment(dst *ebiten.Image) { a.sprite(dst, "parchment.png", 480, 358, 1) }
func (a *app) drawBook(dst *ebiten.Image) {
	a.drawParchment(dst)
	g := a.game
	idx := g.BookPage
	a.drawText(dst, "THE SIGHTINGS BOOK", 480, 155, a.heading, ink, true)
	a.drawText(dst, game.KindNames[idx], 480, 228, a.font, ink, true)
	im := a.sightseeing[idx]
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-75, -170)
	op.GeoM.Scale(1.85, 1.85)
	op.GeoM.Translate(480*2, 310*2)
	dst.DrawImage(im, op)
	a.drawText(dst, fmt.Sprintf("Noticed %d times", g.Save.Book[idx].Count), 480, 455, a.font, ink, true)
	a.drawText(dst, fmt.Sprintf("Page %d of 6    •    ← / → turn    •    SPACE return", idx+1), 480, 572, a.small, ink, true)
}
func (a *app) drawLog(dst *ebiten.Image) {
	a.drawParchment(dst)
	g := a.game
	a.drawText(dst, "THE HARBOUR LOG", 480, 151, a.heading, ink, true)
	nNight := len(g.Save.Log)
	nTotal := nNight + len(g.Save.Desk)
	if nTotal == 0 {
		a.drawText(dst, "Your journeys will appear here", 480, 348, a.font, ink, true)
	} else {
		index := nTotal - 1 - g.LogPage
		if index < nNight {
			r := g.Save.Log[index]
			a.drawText(dst, "Night at sea  •  "+r.Date, 480, 233, a.font, ink, true)
			a.drawText(dst, fmt.Sprintf("Guess %d    |    Lanterns %d", r.Prediction, r.Result), 480, 321, a.font, ink, true)
			a.drawText(dst, fmt.Sprintf("Noticed %d    |    Came back %d", r.Noticed, r.CameBacks), 480, 388, a.font, ink, true)
			a.drawText(dst, fmt.Sprintf("Check-ins %d    |    Good guess: %v", len(r.CheckIns), r.GoodGuess), 480, 454, a.small, ink, true)
		} else {
			r := g.Save.Desk[index-nNight]
			a.drawText(dst, "At my desk  •  "+r.Date, 480, 233, a.font, ink, true)
			a.drawText(dst, fmt.Sprintf("%d minutes    |    Check-ins %d", r.LengthMinutes, len(r.Answers)), 480, 330, a.font, ink, true)
			a.drawText(dst, "Every answer was welcome", 480, 405, a.small, ink, true)
		}
	}
	a.drawText(dst, fmt.Sprintf("Page %d of %d    •    ← / → turn    •    SPACE return", g.LogPage+1, maxIntRender(1, nTotal)), 480, 570, a.small, ink, true)
}
func maxIntRender(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (a *app) drawRest(dst *ebiten.Image) {
	a.panel(dst, 130, 160, 700, 407)
	a.headingText(dst, "THE LIGHTHOUSE RESTS", 480, 249)
	a.subText(dst, "Two nights of lovely noticing", 480, 319)
	a.subText(dst, "Come back tomorrow", 480, 383)
	a.smallText(dst, "B  book    •    L  log", 480, 454)
	if !a.game.Save.ExtraNightUsed {
		a.smallText(dst, "A grown-up may hold P for one extra night", 480, 517)
	}
}
func (a *app) drawBellSetup(dst *ebiten.Image) {
	a.panel(dst, 166, 164, 628, 390)
	a.headingText(dst, "KEEPER'S BELL", 480, 249)
	a.subText(dst, "A little lighthouse for your desk", 480, 326)
	a.subText(dst, fmt.Sprintf("%d minutes", a.game.BellMinutes), 480, 414)
	a.smallText(dst, "← 15 minutes    → 25 minutes    •    SPACE begin", 480, 505)
}
func (a *app) drawBell(dst *ebiten.Image) {
	a.panel(dst, 145, 100, 670, 182)
	a.headingText(dst, "KEEPER'S BELL", 480, 177)
	a.subText(dst, "Your quiet companion", 480, 238)
	if a.game.Bell.Waiting {
		ebitenutil.DrawRect(dst, 0, 0, 1920, 1440, color.RGBA{6, 20, 34, 145})
		a.panel(dst, 170, 278, 620, 188)
		a.subText(dst, "Were you with the boat?", 480, 355)
		a.smallText(dst, "Y  yes     •     N  no", 480, 428)
	} else {
		a.panel(dst, 275, 520, 410, 90)
		a.smallText(dst, "Your light is keeping you company", 480, 556)
		a.smallText(dst, "ESC  finish early", 480, 593)
	}
}
