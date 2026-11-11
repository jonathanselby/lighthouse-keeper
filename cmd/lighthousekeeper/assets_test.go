package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

func TestRequiredArtDecodes(t *testing.T) {
	required := []string{"title_art.jpg", "seascape.jpg", "beam.png", "boat.png", "fog.png", "lighthouse.png", "sightings.png", "parchment.png"}
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
