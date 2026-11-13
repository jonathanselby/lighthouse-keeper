package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	"github.com/jonathanselby/lighthouse-keeper/internal/game"
)

func TestRequiredArtDecodes(t *testing.T) {
	required := []string{"title_art.jpg", "seascape.jpg", "beam.png", "boat.png", "sightings-live.png", "fog.png", "lighthouse.png", "sightings.png", "parchment.png", "water-domain.png", "shore-fixed.png", "rocks-fixed.png", "rocks-bottom.png", "impact-atlas.png", "ui_rail.png", "glow.png", "lantern.png"}
	for _, name := range required {
		data, err := assetFS.ReadFile("assets/" + name)
		if err != nil {
			t.Fatal(name, err)
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(name, err)
		}
		if img.Bounds().Dx() < 1 || img.Bounds().Dy() < 1 {
			t.Fatal(name, "empty")
		}
	}
}
func TestPainterlyAtlasDimensions(t *testing.T) {
	for _, item := range []struct {
		name string
		w, h int
	}{
		{"seascape.jpg", 1920, 1440},
		{"title_art.jpg", 1920, 1440},
		{"sightings.png", 1440, 240},
		{"sightings-live.png", 1440, 240},
		{"ui_rail.png", 960, 220},
	} {
		data, err := assetFS.ReadFile("assets/" + item.name)
		if err != nil {
			t.Fatal(err)
		}
		cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			t.Fatal(item.name, err)
		}
		if cfg.Width != item.w || cfg.Height != item.h {
			t.Fatalf("%s: got %dx%d, want %dx%d", item.name, cfg.Width, cfg.Height, item.w, item.h)
		}
	}
}

func TestPresentationMappingPreservesOrder(t *testing.T) {
	last := sceneY(-1)
	for y := 0.; y <= 720; y += 10 {
		got := sceneY(y)
		if got <= last {
			t.Fatalf("visual coordinate mapping not monotonic: %g", y)
		}
		last = got
	}
	if sceneY(40) < 340 || sceneY(640) > 650 {
		t.Fatal("artwork projection outside scene bounds")
	}
}

func TestSoundCuesProducePCM(t *testing.T) {
	for name, s := range cues {
		b := procedural(s)
		if len(b) < 100 || len(b)%4 != 0 {
			t.Fatalf("%s produced invalid PCM", name)
		}
	}
}
func TestViewportPreservesAspect(t *testing.T) {
	a := &app{}
	for _, wh := range [][2]int{{960, 720}, {1920, 1080}, {2560, 1440}, {800, 1200}} {
		w, h := a.Layout(wh[0], wh[1])
		if w*3 != h*4 {
			t.Fatalf("invalid 4:3 internal layout at %v", wh)
		}
	}
}

func TestLivingOceanArtLayerDimensions(t *testing.T) {
	want := map[string][2]int{
		"water-domain.png": {1920, 1440},
		"impact-atlas.png": {640, 256},
		"shore-fixed.png":  {1090, 680},
		"rocks-fixed.png":  {416, 242},
		"rocks-bottom.png": {650, 150},
	}
	for name, size := range want {
		b, e := assetFS.ReadFile("assets/" + name)
		if e != nil {
			t.Fatal(e)
		}
		c, _, e := image.DecodeConfig(bytes.NewReader(b))
		if e != nil {
			t.Fatal(e)
		}
		if c.Width != size[0] || c.Height != size[1] {
			t.Fatalf("%s: %dx%d", name, c.Width, c.Height)
		}
	}
}

func TestWaterShaderClockHasSubFramePrecision(t *testing.T) {
 for _,ticks:=range []uint64{0,60,3600,216000} {
  a,b,c:=presentationClock(ticks,0),presentationClock(ticks,.5),presentationClock(ticks+1,0)
  if !(a<b&&b<c) {t.Fatalf("frozen shader at %d: %f %f %f",ticks,a,b,c)}
 }
}
func TestWaterSightingEmergenceIsContinuous(t *testing.T) {
 previous:=-1.0
 for age:=0;age<=48;age++ {
  p:=emergenceProgress(game.Sighting{State:game.Appearing,Age:age})
  if p<previous||p<0||p>1 {t.Fatalf("jump %d %f",age,p)}
  previous=p
 }
 if previous!=1 {t.Fatal("did not finish")}
}
