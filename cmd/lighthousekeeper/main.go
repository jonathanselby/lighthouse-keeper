package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jonathanselby/lighthouse-keeper/internal/game"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/*
var assetFS embed.FS

const renderScale = 2.0 // 1920x1440 internal render for sharp native-HD artwork and type

type app struct {
	game                 *game.Game
	images               map[string]*ebiten.Image
	sightseeing          [6]*ebiten.Image
	rail                 [9]*ebiten.Image
	heading, font, small font.Face
	audio                *soundBank
	savePath             string
	saveMessage          string
	lastUpdate           time.Time
	deskStart            time.Time
	deskOffset           int
	lastScreen           game.Screen
	idleTime             time.Time
	now                  time.Time
	particles            []particle
	frame                uint64
}

type particle struct {
	x, y, vx, vy float64
	age, life    int
	hue          uint8
}

func main() {
	path, err := game.SavePath()
	if err != nil {
		log.Fatal(err)
	}
	saved, loadErr := game.LoadSave(path)
	a := &app{game: game.New(time.Now().UnixNano(), saved), images: make(map[string]*ebiten.Image), savePath: path, lastUpdate: time.Now(), particles: make([]particle, 0, 256)}
	if loadErr != nil {
		log.Printf("journal recovery: %v", loadErr)
		a.saveMessage = game.UIText[25]
	}
	for _, name := range []string{"seascape.jpg", "title_art.jpg", "lighthouse.png", "boat.png", "sightings.png", "parchment.png", "beam.png", "fog.png", "shimmer.png", "ui_rail.png", "glow.png", "lantern.png"} {
		data, e := assetFS.ReadFile("assets/" + name)
		if e != nil {
			log.Fatalf("required art %s: %v", name, e)
		}
		im, _, e := image.Decode(bytes.NewReader(data))
		if e != nil {
			log.Fatalf("required art %s: %v", name, e)
		}
		a.images[name] = ebiten.NewImageFromImage(im)
	}
	for i := 0; i < 6; i++ {
		a.sightseeing[i] = a.images["sightings.png"].SubImage(image.Rect(240*i, 0, 240*(i+1), 240)).(*ebiten.Image)
	}
	// Prepare brass-panel nine-slice regions once, not on every draw.
	rail := a.images["ui_rail.png"]
	rw, rh := rail.Bounds().Dx(), rail.Bounds().Dy()
	xp, yp := [4]int{0, 52, rw - 52, rw}, [4]int{0, 32, rh - 32, rh}
	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			a.rail[j*3+i] = rail.SubImage(image.Rect(xp[i], yp[j], xp[i+1], yp[j+1])).(*ebiten.Image)
		}
	}
	a.heading = loadFont(gobold.TTF, 56)
	a.font = loadFont(goregular.TTF, 36)
	a.small = loadFont(goregular.TTF, 25)
	a.audio = loadSoundBank()
	a.lastScreen = a.game.Screen
	ebiten.SetTPS(60)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowTitle("Lighthouse Keeper — Notice it. Name it. Come back.")
	ebiten.SetWindowSize(960, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(true)
	if err = ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}
func loadFont(data []byte, size float64) font.Face {
	tt, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}
	return face
}
func key(k ebiten.Key) bool { return inpututil.IsKeyJustPressed(k) }
func (a *app) Update() error {
	a.now = time.Now()
	in := game.Input{Left: ebiten.IsKeyPressed(ebiten.KeyArrowLeft), Right: ebiten.IsKeyPressed(ebiten.KeyArrowRight),
		LeftPressed: key(ebiten.KeyArrowLeft), RightPressed: key(ebiten.KeyArrowRight), Notice: key(ebiten.KeyN),
		Yes: key(ebiten.KeyY), No: key(ebiten.KeyN), Start: key(ebiten.KeySpace) || key(ebiten.KeyEnter),
		Book: key(ebiten.KeyB), Log: key(ebiten.KeyL), Bell: key(ebiten.KeyK), Escape: key(ebiten.KeyEscape),
		ParentDown: ebiten.IsKeyPressed(ebiten.KeyP), Today: a.now.Format("2006-01-02")}
	if a.game.Screen == game.NightScreen && !a.game.Night.CheckIn {
		in.No = false
	}
	if a.game.Screen == game.BellScreen {
		elapsed := int(time.Since(a.deskStart).Seconds()*60) + a.deskOffset
		if elapsed < 0 {
			elapsed = 0
		}
		in.DeskElapsedTicks = elapsed
	}
	before := a.game.Screen
	if in.Escape && before != game.BellScreen && before != game.Book && before != game.Log && before != game.BellSetup {
		return ebiten.Termination
	}
	a.game.Update(in)
	// Book and log mutations that are not a finished night persist only on completed nights.
	for _, ev := range a.game.Events {
		if ev == game.EvSave {
			if err := game.WriteSave(a.savePath, &a.game.Save); err != nil {
				log.Printf("journal save: %v", err)
				a.saveMessage = game.UIText[26]
			} else {
				a.saveMessage = ""
			}
		}
		a.audio.playEvent(ev)
		a.spawnEffects(ev)
	}
	if before != a.game.Screen {
		if a.game.Screen == game.BellScreen {
			a.deskStart = a.now
			ebiten.SetFullscreen(false)
			ebiten.SetWindowSize(960, 720)
		}
		if before == game.BellScreen && a.game.Screen != game.BellScreen {
			ebiten.SetFullscreen(true)
		}
	}
	a.audio.ambientTick()
	a.updateParticles()
	a.lastUpdate = a.now
	a.frame++
	return nil
}
func (a *app) Draw(dst *ebiten.Image) { a.drawFrame(dst) }
func (a *app) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(game.Width * renderScale), int(game.Height * renderScale)
}
func (a *app) subTick() float64 {
	return math.Max(0, math.Min(1, time.Since(a.lastUpdate).Seconds()*60))
}
func (a *app) spawnEffects(ev game.Event) {
	if ev != game.EvNotice && ev != game.EvCameBack && ev != game.EvNewKind && ev != game.EvCrossingEnd {
		return
	}
	var x, y float64
	if ev == game.EvNotice || ev == game.EvNewKind {
		x = game.LampX
		y = game.LampY - 10
		for _, s := range a.game.Night.Sightings {
			if s.State == game.Noticed {
				x = s.Pos.X
				y = s.Pos.Y
				break
			}
		}
	} else {
		p := a.game.Night.Boat.Position()
		x = p.X
		y = p.Y
	}
	for i := 0; i < 24 && len(a.particles) < 256; i++ {
		alpha := float64(i) * math.Pi * 2 / 24
		a.particles = append(a.particles, particle{x: x, y: y, vx: math.Cos(alpha) * (0.5 + float64(i%5)*0.16), vy: math.Sin(alpha) * (0.5 + float64(i%4)*.13), life: 30 + i%25, hue: uint8(i % 3)})
	}
}
func (a *app) updateParticles() {
	live := a.particles[:0]
	for _, p := range a.particles {
		p.x += p.vx
		p.y += p.vy
		p.vy += .006
		p.age++
		if p.age < p.life {
			live = append(live, p)
		}
	}
	a.particles = live
}
func (a *app) readmeStatus() string {
	names := []string{"Title", "Predict", "Night", "Harbour", "Book", "Log", "Rest", "BellSetup", "Bell"}
	return fmt.Sprintf("%s — %d sightings", names[a.game.Screen], len(a.game.Night.Sightings))
}
