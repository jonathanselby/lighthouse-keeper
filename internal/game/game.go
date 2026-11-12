// Package game contains Lighthouse Keeper's deterministic, headless rules.
// It performs no I/O, reads no wall clock and knows nothing about rendering.
package game

import (
	"math"
	"math/rand"
)

const (
	Width          = 960
	Height         = 720
	TicksPerSecond = 60
	LampX          = 480.0
	LampY          = 640.0
	HarbourX       = 820.0
	HarbourY       = 590.0
	BeamRange      = 820.0
	MaxBeamAngle   = 1.50
	MaxBeamSpeed   = 0.025
	MaxSightings   = 6
)

type Screen uint8

const (
	Title Screen = iota
	Predict
	NightScreen
	Harbour
	Book
	Log
	Rest
	BellSetup
	BellScreen
)

type Event uint8

const (
	EvCrossingStart Event = iota
	EvBoatLit
	EvBoatUnlit
	EvSightingAppear
	EvNotice
	EvNewKind
	EvSightingLeft
	EvWhereBoat
	EvCameBack
	EvBell
	EvCheckInAnswered
	EvCrossingEnd
	EvNightEnd
	EvDeskBell
	EvBellEnd
	EvSave
)

type Input struct {
	Left, Right, ParentDown                         bool
	LeftPressed, RightPressed                       bool
	Notice, Yes, No, Start, Book, Log, Bell, Escape bool
	Today                                           string
	DeskElapsedTicks                                int
}

type Point struct{ X, Y float64 }

type Beam struct{ Angle, Omega, HalfWidth float64 }

type Boat struct {
	Path     [4]Point
	T        float64
	Distance float64
	Length   float64
	Arc      [65]float64
	Lit      bool
}

type Night struct {
	Index, Prediction, Crossing                   int
	Boat                                          Boat
	Fog                                           float64
	LitTicks, TotalTicks                          [5]int
	Sightings                                     []Sighting
	NextSighting, OffTicks, LitStreak, ChimeTicks int
	ChimeUsed                                     bool
	CameBacks, Noticed                            int
	BellAt                                        [5]int
	CheckIn                                       bool
	Answers                                       []bool
	Ticks, CrossingTicks, PauseTicks              int
	Result                                        int
	GoodGuess                                     bool
}

type BellSession struct {
	LengthTicks, Elapsed, NextRingAt int
	Answers                          []bool
	Waiting                          bool
	EndedEarly                       bool
}

type Game struct {
	Screen            Screen
	Night             Night
	Bell              BellSession
	Save              Save
	Beam              Beam
	Events            []Event
	RNG               *rand.Rand
	ParentHoldTicks   int
	BookPage, LogPage int
	ReturnScreen      Screen
	BellMinutes       int
	SaveWarning       string // adapter-owned diagnostics can display recovery status
}

func New(seed int64, save Save) *Game {
	if save.Version == 0 {
		save = NewSave()
	}
	g := &Game{Screen: Title, Save: save, RNG: rand.New(rand.NewSource(seed)), Events: make([]Event, 0, 8), BellMinutes: 15}
	return g
}

func (g *Game) emit(e Event) {
	if len(g.Events) < 8 {
		g.Events = append(g.Events, e)
	}
}
func clamp(x, a, b float64) float64 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
func lerp(a, b, t float64) float64 { return a + (b-a)*t }

func (g *Game) Update(in Input) {
	g.Events = g.Events[:0]
	switch g.Screen {
	case Title:
		switch {
		case in.Start:
			g.refreshDate(in.Today)
			if g.restRequired(in.Today) {
				g.Screen = Rest
			} else {
				g.beginPredict()
			}
		case in.Book:
			g.openBook(Title)
		case in.Log:
			g.openLog(Title)
		case in.Bell:
			g.Screen = BellSetup
		}
	case Predict:
		if in.LeftPressed {
			g.Night.Prediction = minInt(5, maxInt(0, g.Night.Prediction-1))
		}
		if in.RightPressed {
			g.Night.Prediction = minInt(5, maxInt(0, g.Night.Prediction+1))
		}
		if in.Start {
			g.startNight()
		}
	case NightScreen:
		g.updateNight(in)
	case Harbour:
		if in.Start {
			g.Screen = Title
		}
		if in.Book {
			g.openBook(Harbour)
		}
		if in.Log {
			g.openLog(Harbour)
		}
	case Book:
		if in.LeftPressed && g.BookPage > 0 {
			g.BookPage--
		}
		if in.RightPressed && g.BookPage < 5 {
			g.BookPage++
		}
		if in.Start || in.Escape {
			g.Screen = g.ReturnScreen
		}
	case Log:
		n := len(g.Save.Log) + len(g.Save.Desk)
		if in.LeftPressed && g.LogPage > 0 {
			g.LogPage--
		}
		if in.RightPressed && g.LogPage < n-1 {
			g.LogPage++
		}
		if in.Start || in.Escape {
			g.Screen = g.ReturnScreen
		}
	case Rest:
		g.refreshDate(in.Today)
		if in.Book {
			g.openBook(Rest)
			return
		}
		if in.Log {
			g.openLog(Rest)
			return
		}
		if in.Start && !g.restRequired(in.Today) {
			g.beginPredict()
			return
		}
		if in.ParentDown && !g.Save.ExtraNightUsed && in.Today != "" {
			g.ParentHoldTicks++
			if g.ParentHoldTicks >= 180 {
				g.Save.ExtraNightUsed = true
				g.emit(EvSave)
				g.ParentHoldTicks = 0
				g.beginPredict()
			}
		} else {
			g.ParentHoldTicks = 0
		}
	case BellSetup:
		if in.LeftPressed {
			g.BellMinutes = 15
		}
		if in.RightPressed {
			g.BellMinutes = 25
		}
		if in.Start {
			g.startDesk()
		}
		if in.Escape {
			g.Screen = Title
		}
	case BellScreen:
		g.updateDesk(in)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (g *Game) openBook(back Screen) { g.ReturnScreen = back; g.BookPage = 0; g.Screen = Book }
func (g *Game) openLog(back Screen)  { g.ReturnScreen = back; g.LogPage = 0; g.Screen = Log }
func (g *Game) restRequired(today string) bool {
	return today != "" && g.Save.LastDate == today && g.Save.NightsToday >= 2
}
func (g *Game) refreshDate(today string) {
	if today == "" {
		return
	}
	if today != g.Save.LastDate {
		g.Save.LastDate = today
		g.Save.NightsToday = 0
		g.Save.ExtraNightUsed = false
		g.emit(EvSave)
	}
}
func (g *Game) beginPredict() {
	prediction := 3
	if len(g.Save.Log) > 0 {
		prediction = g.Save.Log[len(g.Save.Log)-1].Result
	}
	g.Night = Night{Index: minInt(g.Save.NightsCompleted, 10), Prediction: prediction, Sightings: make([]Sighting, 0, 6), Answers: make([]bool, 0, 3)}
	g.Screen = Predict
}

func (g *Game) startNight() {
	g.Screen = NightScreen
	idx := g.Night.Index
	width := lerp(9, 6.5, float64(idx)/10) * math.Pi / 180
	g.Beam = Beam{HalfWidth: width}
	// A night has two or three check-ins, scheduled in guaranteed crossings.
	slots := [3]int{0, 2, 4}
	for i := len(slots) - 1; i > 0; i-- {
		j := g.RNG.Intn(i + 1)
		slots[i], slots[j] = slots[j], slots[i]
	}
	count := 2 + g.RNG.Intn(2)
	for i := 0; i < count; i++ {
		c := slots[i]
		if c == 0 {
			g.Night.BellAt[c] = (19 + g.RNG.Intn(10)) * 60
		} else {
			g.Night.BellAt[c] = (15 + g.RNG.Intn(10)) * 60
		}
	}
	g.Night.NextSighting = g.nextSightingTick()
	g.startCrossing(0)
}

func (g *Game) startCrossing(i int) {
	g.Night.Crossing = i
	g.Night.CrossingTicks = 0
	g.Night.OffTicks = 0
	g.Night.LitStreak = 0
	g.Night.ChimeTicks = 0
	g.Night.ChimeUsed = false
	g.Night.Fog = 0
	g.Night.Boat = g.makeBoat()
	// Starting beam is aimed towards the boat, not given a random input penalty.
	p := g.Night.Boat.Position()
	g.Beam.Angle = Bearing(p)
	g.Beam.Omega = 0
	g.emit(EvCrossingStart)
}

// crossingTargetTicks is the full-light crossing duration. The first sea
// journey remains gentle, then the route becomes more active across the five
// crossings and later nights; an unlit boat still always arrives.
func (g *Game) crossingTargetTicks() int {
	n := g.Night
	ticks := 3300 - minInt(n.Index, 10)*65 - n.Crossing*150
	if ticks < 1950 {
		return 1950
	}
	return ticks
}

// CrossingFullLightTicks exposes duration to the visual interpolation only.
func (g *Game) CrossingFullLightTicks() int { return g.crossingTargetTicks() }

func (g *Game) updateNight(in Input) {
	n := &g.Night
	if n.CheckIn {
		if in.Yes || in.No {
			n.Answers = append(n.Answers, in.Yes)
			n.CheckIn = false
			g.emit(EvCheckInAnswered)
		}
		return
	}
	if n.PauseTicks > 0 {
		n.PauseTicks--
		if n.PauseTicks == 0 && n.CrossingTicks < 0 {
			g.startCrossing(n.Crossing)
		}
		return
	}
	n.Ticks++
	n.CrossingTicks++
	if n.BellAt[n.Crossing] != 0 && n.CrossingTicks == n.BellAt[n.Crossing] {
		n.CheckIn = true
		g.emit(EvBell)
		return
	}
	target := 0.0
	if in.Left && !in.Right {
		target = -MaxBeamSpeed
	} else if in.Right && !in.Left {
		target = MaxBeamSpeed
	}
	g.Beam.Omega += (target - g.Beam.Omega) * 0.2
	g.Beam.Angle = clamp(g.Beam.Angle+g.Beam.Omega, -MaxBeamAngle, MaxBeamAngle)
	boat := &n.Boat
	wasLit := boat.Lit
	boat.Lit = WithinBeam(boat.Position(), g.Beam)
	if boat.Lit {
		n.Fog = clamp(n.Fog-1.0/60, 0, 1)
	} else {
		n.Fog = clamp(n.Fog+1.0/90, 0, 1)
	}
	if !wasLit && boat.Lit {
		g.emit(EvBoatLit)
	}
	if wasLit && !boat.Lit {
		g.emit(EvBoatUnlit)
	}
	n.TotalTicks[n.Crossing]++
	if boat.Lit {
		n.LitTicks[n.Crossing]++
	}
	boat.Advance(lerp(1.0, 0.40, n.Fog) * boat.Length / float64(g.crossingTargetTicks()))
	g.updateReturn()
	g.updateSightings(in.Notice)
	if boat.T >= 1 {
		g.emit(EvCrossingEnd)
		if n.Crossing == 4 {
			g.finishNight(in.Today)
		} else {
			n.PauseTicks = 120
			// Save completion state until pause has elapsed, no extra screen machinery.
			n.Crossing++
			n.CrossingTicks = -120 // see pause rule below: reset on start of next crossing
		}
	}
	// When a pause finishes, generate next crossing with its own path, not in advance.
	if n.PauseTicks == 0 && n.CrossingTicks < 0 {
		g.startCrossing(n.Crossing)
	}
}

func (g *Game) updateReturn() {
	n := &g.Night
	if n.Boat.Lit {
		n.OffTicks = 0
		n.LitStreak++
	} else {
		n.OffTicks++
		n.LitStreak = 0
	}
	if n.LitStreak >= 60 {
		n.ChimeUsed = false
	}
	if n.OffTicks >= 120 && n.CrossingTicks >= 180 && !n.ChimeUsed {
		n.ChimeUsed = true
		n.ChimeTicks = 300
		g.emit(EvWhereBoat)
	}
	if n.ChimeTicks > 0 {
		if n.Boat.Lit {
			n.CameBacks++
			n.ChimeTicks = 0
			g.emit(EvCameBack)
		} else {
			n.ChimeTicks--
		}
	}
}

func (g *Game) finishNight(today string) {
	n := &g.Night
	n.Result = CountLanterns(n.LitTicks, n.TotalTicks)
	n.GoodGuess = absInt(n.Prediction-n.Result) <= 1
	g.refreshDate(today)
	g.Save.NightsCompleted++
	if today != "" {
		g.Save.NightsToday++
	}
	g.Save.Log = append(g.Save.Log, NightRecord{Date: today, Prediction: n.Prediction, Result: n.Result, GoodGuess: n.GoodGuess, CameBacks: n.CameBacks, Noticed: n.Noticed, CheckIns: append([]bool(nil), n.Answers...)})
	if len(g.Save.Log) > 400 {
		g.Save.Log = g.Save.Log[len(g.Save.Log)-400:]
	}
	g.Screen = Harbour
	g.emit(EvNightEnd)
	g.emit(EvSave)
}

func CountLanterns(lit, total [5]int) int {
	count := 0
	for i := 0; i < 5; i++ {
		if total[i] > 0 && float64(lit[i])/float64(total[i]) >= 0.70 {
			count++
		}
	}
	return count
}

func (g *Game) startDesk() {
	g.Bell = BellSession{LengthTicks: g.BellMinutes * 60 * 60, NextRingAt: (180 + g.RNG.Intn(241)) * 60, Answers: make([]bool, 0, 9)}
	g.Screen = BellScreen
}
func (g *Game) updateDesk(in Input) {
	b := &g.Bell
	if in.DeskElapsedTicks > b.Elapsed {
		b.Elapsed = in.DeskElapsedTicks
	}
	if b.Elapsed >= b.LengthTicks || in.Escape {
		g.endDesk(in.Today, in.Escape)
		return
	}
	if b.Waiting {
		if in.Yes || in.No {
			b.Answers = append(b.Answers, in.Yes)
			b.Waiting = false
			b.NextRingAt = b.Elapsed + (180+g.RNG.Intn(241))*60
			g.emit(EvCheckInAnswered)
		}
		return
	}
	if b.Elapsed >= b.NextRingAt {
		b.Waiting = true
		g.emit(EvDeskBell)
	}
}
func (g *Game) endDesk(today string, early bool) {
	b := &g.Bell
	b.EndedEarly = early
	g.Save.Desk = append(g.Save.Desk, DeskRecord{Date: today, LengthMinutes: b.LengthTicks / 3600, Answers: append([]bool(nil), b.Answers...), UnansweredAtEnd: b.Waiting, EndedEarly: early})
	if len(g.Save.Desk) > 400 {
		g.Save.Desk = g.Save.Desk[len(g.Save.Desk)-400:]
	}
	g.Screen = Title
	g.emit(EvBellEnd)
	g.emit(EvSave)
}

// Bearing is clockwise from north in radians. Zero means straight up.
func Bearing(p Point) float64 { return math.Atan2(p.X-LampX, LampY-p.Y) }
func WithinBeam(p Point, b Beam) bool {
	dx, dy := p.X-LampX, p.Y-LampY
	if dx*dx+dy*dy > BeamRange*BeamRange {
		return false
	}
	diff := math.Remainder(Bearing(p)-b.Angle, 2*math.Pi)
	return math.Abs(diff) <= b.HalfWidth
}
func (b Boat) Position() Point { return Bezier(b.Path, b.T) }
func Bezier(path [4]Point, t float64) Point {
	t = clamp(t, 0, 1)
	u := 1 - t
	a := u * u * u
	bb := 3 * u * u * t
	c := 3 * u * t * t
	d := t * t * t
	return Point{a*path[0].X + bb*path[1].X + c*path[2].X + d*path[3].X, a*path[0].Y + bb*path[1].Y + c*path[2].Y + d*path[3].Y}
}
func (b *Boat) Initialise() {
	b.Arc[0] = 0
	prev := Bezier(b.Path, 0)
	for i := 1; i <= 64; i++ {
		p := Bezier(b.Path, float64(i)/64)
		b.Arc[i] = b.Arc[i-1] + math.Hypot(p.X-prev.X, p.Y-prev.Y)
		prev = p
	}
	b.Length = b.Arc[64]
}
func (b *Boat) Advance(d float64) {
	b.Distance = clamp(b.Distance+d, 0, b.Length)
	if b.Distance >= b.Length {
		b.T = 1
		return
	}
	for i := 1; i <= 64; i++ {
		if b.Arc[i] >= b.Distance {
			s := (b.Distance - b.Arc[i-1]) / (b.Arc[i] - b.Arc[i-1])
			b.T = (float64(i-1) + s) / 64
			return
		}
	}
}
func (g *Game) makeBoat() Boat {
	n := g.Night.Index
	edge := g.RNG.Intn(3)
	var start Point
	switch edge {
	case 0:
		start = Point{X: 30, Y: 80 + g.RNG.Float64()*450}
	case 1:
		start = Point{X: 100 + g.RNG.Float64()*730, Y: 40}
	default:
		start = Point{X: 930, Y: 80 + g.RNG.Float64()*430}
	}
	end := Point{X: HarbourX, Y: HarbourY}
	offset := lerp(125, 270, float64(n)/10)
	sign := 1.0
	if g.RNG.Intn(2) == 0 {
		sign = -1
	}
	c1 := Point{X: clamp(lerp(start.X, end.X, 0.34)+sign*offset*.70, 22, 936), Y: clamp(lerp(start.Y, end.Y, 0.34)-offset*.15, 25, 590)}
	c2 := Point{X: clamp(lerp(start.X, end.X, 0.68)-sign*offset*1.15, 22, 936), Y: clamp(lerp(start.Y, end.Y, 0.69)+offset*.11, 25, 590)}
	b := Boat{Path: [4]Point{start, c1, c2, end}}
	b.Initialise()
	return b
}
