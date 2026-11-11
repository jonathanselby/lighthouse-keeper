package main

import (
	"encoding/binary"
	"math"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/jonathanselby/lighthouse-keeper/internal/game"
)

const sampleRate = 44100

type soundBank struct {
	context *audio.Context
	players map[string]*audio.Player
	ambient *audio.Player
}
type cue struct{ freq, dur, vol float64 }

var cues = map[string]cue{
	"lit": {560, .055, .16}, "appear": {670, .11, .10}, "notice": {810, .18, .14},
	"newkind": {930, .24, .13}, "chime": {450, .42, .17}, "back": {720, .19, .18},
	"bell": {520, .42, .16}, "checkin": {800, .13, .13}, "harbour": {430, .31, .16},
	"desk_end": {740, .46, .14},
}

func procedural(spec cue) []byte {
	n := int(spec.dur * sampleRate)
	b := make([]byte, n*4)
	for i := 0; i < n; i++ {
		sec := float64(i) / sampleRate
		phase := float64(i) / float64(n)
		wave := math.Sin(2*math.Pi*spec.freq*sec) + .17*math.Sin(4*math.Pi*spec.freq*sec)
		envelope := math.Sin(math.Pi*phase) * math.Sin(math.Pi*phase)
		v := int16(math.Max(-1, math.Min(1, wave*spec.vol*envelope)) * 32767)
		binary.LittleEndian.PutUint16(b[i*4:], uint16(v))
		binary.LittleEndian.PutUint16(b[i*4+2:], uint16(v))
	}
	return b
}
func ambientPCM() []byte {
	n := sampleRate * 8
	b := make([]byte, n*4)
	for i := 0; i < n; i++ {
		t := float64(i) / sampleRate
		// Four slow, phase-locked wave layers form a seamless repeating eight-second bed.
		wave := .26*math.Sin(2*math.Pi*float64(4)*t/8) + .17*math.Sin(2*math.Pi*float64(7)*t/8)
		wave += .08 * math.Sin(2*math.Pi*float64(12)*t/8)
		v := int16(wave * 3300)
		binary.LittleEndian.PutUint16(b[i*4:], uint16(v))
		binary.LittleEndian.PutUint16(b[i*4+2:], uint16(v))
	}
	return b
}
func audioLikelyAvailable() bool {
	if runtime.GOOS != "linux" {
		return true
	}
	if os.Getenv("ALSOFT_DRIVERS") == "null" {
		return true
	}
	_, err := os.Stat("/dev/snd")
	return err == nil
}
func loadSoundBank() *soundBank {
	b := &soundBank{players: make(map[string]*audio.Player)}
	if !audioLikelyAvailable() {
		return b
	}
	b.context = audio.NewContext(sampleRate)
	for name, s := range cues {
		b.players[name] = b.context.NewPlayerFromBytes(procedural(s))
	}
	b.ambient = b.context.NewPlayerFromBytes(ambientPCM())
	return b
}
func (b *soundBank) ambientTick() {
	if b == nil || b.ambient == nil {
		return
	}
	if !b.ambient.IsPlaying() {
		_ = b.ambient.Rewind()
		b.ambient.Play()
	}
}
func (b *soundBank) play(name string) {
	if b == nil || name == "" {
		return
	}
	p := b.players[name]
	if p == nil {
		return
	}
	if err := p.Rewind(); err == nil {
		p.Play()
	}
}
func (b *soundBank) playEvent(ev game.Event) {
	var name string
	switch ev {
	case game.EvBoatLit:
		name = "lit"
	case game.EvSightingAppear:
		name = "appear"
	case game.EvNotice:
		name = "notice"
	case game.EvNewKind:
		name = "newkind"
	case game.EvWhereBoat:
		name = "chime"
	case game.EvCameBack:
		name = "back"
	case game.EvBell, game.EvDeskBell:
		name = "bell"
	case game.EvCheckInAnswered:
		name = "checkin"
	case game.EvCrossingEnd:
		name = "harbour"
	case game.EvBellEnd:
		name = "desk_end"
	}
	b.play(name)
}
