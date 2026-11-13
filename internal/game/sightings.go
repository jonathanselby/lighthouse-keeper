package game

import "math"

type Kind uint8

const (
	Whale Kind = iota
	ShootingStar
	Seal
	Jellyfish
	Owl
	Ferry
	KindCount
)

var KindNames = [6]string{"Whale", "Shooting Star", "Seal", "Jellyfish", "Owl", "Ferry"}

type SightState uint8

const (
	Appearing SightState = iota
	Visible
	Noticed
	Chased
	Leaving
)

type Sighting struct {
	Kind      Kind
	Pos       Point
	Age, Life int
	State     SightState
	Aim       int
}

func (g *Game) nextSightingTick() int {
	progress := float64(g.Night.Index) / 10
	low := int(lerp(12, 8, progress) * 60)
	high := int(lerp(25, 16, progress) * 60)
	return g.Night.Ticks + low + g.RNG.Intn(high-low+1)
}
func (g *Game) updateSightings(notice bool) {
	n := &g.Night
	for i := range n.Sightings {
		s := &n.Sightings[i]
		s.Age++
		switch s.State {
		case Appearing:
			if s.Age >= 48 {
				s.State = Visible
			}
		case Visible:
			if s.Age >= s.Life {
				s.State = Leaving
				s.Age = 0
			}
			if s.State == Visible && !n.Boat.Lit && WithinBeam(s.Pos, g.Beam) {
				s.Aim++
				if s.Aim >= 60 {
					s.State = Chased
					s.Age = 0
					g.emit(EvSightingLeft)
				}
			} else {
				s.Aim = 0
			}
		case Noticed:
			if s.Age >= 60 {
				s.State = Leaving
				s.Age = 0
			}
		}
	}
	if notice {
		oldest := -1
		for i := range n.Sightings {
			if n.Sightings[i].State == Visible && (oldest < 0 || n.Sightings[i].Age > n.Sightings[oldest].Age) {
				oldest = i
			}
		}
		if oldest >= 0 {
			s := &n.Sightings[oldest]
			s.State = Noticed
			s.Age = 0
			n.Noticed++
			g.Save.Book[s.Kind].Count++
			g.emit(EvNotice)
			if g.Save.Book[s.Kind].Count == 1 {
				g.emit(EvNewKind)
			}
		}
	}
	// Stable compaction; no iterator or entity service.
	alive := n.Sightings[:0]
	for _, s := range n.Sightings {
		if s.State == Leaving && s.Age >= 30 || s.State == Chased && s.Age >= 30 {
			continue
		}
		alive = append(alive, s)
	}
	n.Sightings = alive
	if n.Ticks >= n.NextSighting {
		visible := 0
		for _, s := range n.Sightings {
			if s.State == Appearing || s.State == Visible {
				visible++
			}
		}
		if visible < 3 && len(n.Sightings) < MaxSightings {
			g.spawnSighting()
		}
		n.NextSighting = g.nextSightingTick()
	}
}
// Habitat-specific positions authored against the approved sea/shore painting.
// Sea animals cannot be placed on land. Seal remains on nearshore rocks (FR4).
func (g *Game) sightingPosition(kind Kind) Point {
 switch kind {
 case Whale,Jellyfish: return Point{X:415+g.RNG.Float64()*340,Y:205+g.RNG.Float64()*245}
 case Seal: return Point{X:240+g.RNG.Float64()*80,Y:510+g.RNG.Float64()*55}
 case Ferry: return Point{X:540+g.RNG.Float64()*240,Y:215+g.RNG.Float64()*210}
 case Owl: return Point{X:65+g.RNG.Float64()*75,Y:160+g.RNG.Float64()*140}
 default: return Point{X:330+g.RNG.Float64()*390,Y:20+g.RNG.Float64()*70}
 }
}
func (g *Game) spawnSighting() {
 n:=&g.Night
 weights:=[6]int{5,1,5,5,2,2}
 draw:=g.RNG.Intn(20)
 var kind Kind
 for i,w:=range weights {
   if draw<w {kind=Kind(i);break}
   draw-=w
 }
 for attempt:=0;attempt<24;attempt++ {
   p:=g.sightingPosition(kind)
   if !clearOfPath(p,n.Boat.Path,115) {continue}
   n.Sightings=append(n.Sightings,Sighting{Kind:kind,Pos:p,Life:(6+g.RNG.Intn(4))*60})
   g.emit(EvSightingAppear)
   return
 }
}
func clearOfPath(p Point, path [4]Point, minDistance float64) bool {
	prev := Bezier(path, 0)
	for i := 1; i <= 64; i++ {
		next := Bezier(path, float64(i)/64)
		dx, dy := next.X-prev.X, next.Y-prev.Y
		l2 := dx*dx + dy*dy
		t := 0.0
		if l2 > 0 {
			t = clamp(((p.X-prev.X)*dx+(p.Y-prev.Y)*dy)/l2, 0, 1)
		}
		if math.Hypot(p.X-(prev.X+t*dx), p.Y-(prev.Y+t*dy)) < minDistance {
			return false
		}
		prev = next
	}
	return true
}
