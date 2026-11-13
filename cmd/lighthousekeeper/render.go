package main

import (
	"fmt"
	"image/color"
	"math"

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
	src := a.images["ui_rail.png"]
	if src == nil {
		return
	}
	sw, sh := src.Bounds().Dx(), src.Bounds().Dy()
	xm := [4]int{0, 52, sw - 52, sw}
	ym := [4]int{0, 32, sh - 32, sh}
	dw, dh := w*renderScale, h*renderScale
	xd := [4]float64{0, 52, dw - 52, dw}
	yd := [4]float64{0, 32, dh - 32, dh}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			part := a.rail[j*3+i]
			op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
			op.GeoM.Scale((xd[i+1]-xd[i])/float64(xm[i+1]-xm[i]),
				(yd[j+1]-yd[j])/float64(ym[j+1]-ym[j]))
			op.GeoM.Translate(x*renderScale+xd[i], y*renderScale+yd[j])
			dst.DrawImage(part, op)
		}
	}
}

// This affine presentation transform keeps every simulated boat in the water.
// It is applied equally to all visible gameplay objects and the light cone;
// gameplay geometry and deterministic collision decisions remain unchanged.
func sceneY(y float64) float64 { return 338 + y*.42 }

func (a *app) drawLantern(dst *ebiten.Image, x, y, scale float64, lit bool) {
	img := a.images["lantern.png"]
	if img == nil {
		return
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x*renderScale, y*renderScale)
	if !lit {
		op.ColorScale.Scale(.32, .41, .52, .74)
	}
	dst.DrawImage(img, op)
	if lit {
		a.imageAt(dst, a.images["glow.png"], x, y, scale*.48)
	}
}

func (a *app) lamp(dst *ebiten.Image) {
	x, y := game.LampX, sceneY(game.LampY)
	// The authored lighthead is large enough to read against a detailed sea.
	a.imageAt(dst, a.images["glow.png"], x, y-8, 1.8)
	a.imageAt(dst, a.images["lighthouse.png"], x, y-18, .78)
}

func (a *app) spotlight(dst *ebiten.Image, angle, halfWidth float64) {
	img := a.images["beam.png"]
	if img == nil {
		return
	}
	bw, bh := img.Bounds().Dx(), img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	// Authored cone points upwards and its narrow end is the lamp position.
	op.GeoM.Translate(-float64(bw)/2, -float64(bh))
	op.GeoM.Scale(1.18*halfWidth/(12*math.Pi/180), 1.17)
	op.GeoM.Rotate(angle)
	// Same transform as dynamic boat positions: visual and logical aim agree.
	op.GeoM.Scale(1, .42)
	op.GeoM.Translate(game.LampX*renderScale, sceneY(game.LampY)*renderScale)
	dst.DrawImage(img, op)
}

func (a *app) drawSea(dst *ebiten.Image) {
	// CPS/Amiga priority planes: water samples CLEAN paint, then explicit
	// spray is drawn behind and in front of separate, immutable stone art.
	seconds := float64(a.presentationClock())
	a.drawLivingWater(dst)
	a.drawImpacts(dst, true, seconds)
	dst.DrawImage(a.images["shore-fixed.png"], &ebiten.DrawImageOptions{GeoM: moved(0, 760)})
	dst.DrawImage(a.images["rocks-fixed.png"], &ebiten.DrawImageOptions{GeoM: moved(314, 1040)})
	dst.DrawImage(a.images["rocks-bottom.png"], &ebiten.DrawImageOptions{GeoM: moved(900, 1290)})
	a.drawImpacts(dst, false, seconds)
}

// Translate already-authored transparent foreground layers without mutation.
func moved(x, y float64) ebiten.GeoM { var m ebiten.GeoM; m.Translate(x, y); return m }

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
		_ = col
		a.drawLantern(dst, x, 380, .94, lit)
	}
	a.subText(dst, fmt.Sprintf("%d of 5", a.game.Night.Prediction), 480, 441)
	a.smallText(dst, "← / → choose    •    SPACE begin", 480, 512)
}
func (a *app) drawNight(dst *ebiten.Image) {
	a.drawSea(dst)
	g := a.game
	n := &g.Night
	dt := a.subTick()
	seconds := float64(a.presentationClock())
	angle := g.Beam.Angle
	boat := n.Boat
	if !n.CheckIn && n.PauseTicks == 0 {
		angle = math.Max(-game.MaxBeamAngle, math.Min(game.MaxBeamAngle, angle+g.Beam.Omega*dt))
		boat.Advance((1 - .60*n.Fog) * boat.Length / float64(g.CrossingFullLightTicks()) * dt)
	}
	a.spotlight(dst, angle, g.Beam.HalfWidth)
	for i, s := range n.Sightings {
		a.drawLiveSighting(dst, s, i, seconds)
	}
	p := boat.Position()
	px, py := p.X, sceneY(p.Y)
	if n.Fog > .01 {
		im := a.images["fog.png"]
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		w, h := im.Bounds().Dx(), im.Bounds().Dy()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Scale(.67+n.Fog*.35, .66+n.Fog*.30)
		op.GeoM.Translate(px*renderScale, py*renderScale)
		op.ColorScale.ScaleAlpha(float32(n.Fog * .85))
		dst.DrawImage(im, op)
	}
	if boat.Lit {
		a.imageAt(dst, a.images["glow.png"], px, py, .76)
	}
	a.drawWake(dst, px, py, seconds)
	a.drawBoat(dst, px, py+math.Sin(seconds*2.5)*1.3, math.Sin(seconds*1.3)*.035)
	a.lamp(dst)
	for i := 0; i < 5; i++ {
		lit := n.TotalTicks[i] > 0 && float64(n.LitTicks[i])/float64(n.TotalTicks[i]) >= .7 &&
			(i < n.Crossing || n.Crossing == 4 && boat.T == 1)
		a.drawLantern(dst, 740+float64(i)*38, 598, .40, lit)
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
		alpha := float64(p.life-p.age) / float64(p.life)
		col := color.RGBA{255, 218, 142, uint8(220 * alpha)}
		vector.DrawFilledCircle(dst, float32(p.x*renderScale), float32(sceneY(p.y)*renderScale),
			float32(3+int(p.hue)), col, true)
	}
}

func (a *app) drawBoat(dst *ebiten.Image, x, y, rock float64) {
	im := a.images["boat.png"]
	if im == nil {
		return
	}
	w, h := im.Bounds().Dx(), im.Bounds().Dy()
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(.80, .80)
	op.GeoM.Rotate(rock)
	op.GeoM.Translate(x*renderScale, y*renderScale)
	dst.DrawImage(im, op)
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
		_ = c
		a.drawLantern(dst, x, 326, 1.00, n.TotalTicks[i] > 0 && float64(n.LitTicks[i])/float64(n.TotalTicks[i]) >= .7)
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
	op.GeoM.Translate(-120, -120)
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
