package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

func TestRequiredArtDecodes(t *testing.T) {
	required := []string{"title_art.jpg", "seascape.jpg", "beam.png", "boat.png", "fog.png", "lighthouse.png", "sightings.png", "parchment.png", "ui_rail.png", "glow.png", "lantern.png"}
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
	if sceneY(40) < 130 || sceneY(640) > 650 {
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

func TestSeascapeIsApprovedSinglePainting(t *testing.T) {
	data, err := assetFS.ReadFile("assets/seascape.jpg")
	if err != nil {
		t.Fatal(err)
	}
	// Prevent an old stitched seascape from silently passing CI merely
	// because it decodes and has the correct pixel dimensions.
	const expected = "d1ac25cc66411dd87627b26b32b00aa90c527cfa15ab429179220f742818de6a"
	got := fmt.Sprintf("%x", sha256.Sum256(data))
	if got != expected {
		t.Fatalf("seascape source mismatch: %s (expected approved master %s)", got, expected)
	}
}
