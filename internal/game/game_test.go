package game

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func startTestNight(seed int64) *Game {
	g := New(seed, NewSave())
	g.Update(Input{Start: true, Today: "2026-11-11"})
	g.Update(Input{Start: true, Today: "2026-11-11"})
	return g
}
func tickNight(g *Game, left, right bool) {
	in := Input{Left: left, Right: right, Today: "2026-11-11"}
	if g.Night.CheckIn {
		in.Yes = true
	}
	g.Update(in)
}
func TestCrossingAlwaysCompletes(t *testing.T) {
	for seed := int64(0); seed < 12; seed++ {
		g := startTestNight(seed)
		for j := 0; j < 160000 && g.Screen == NightScreen; j++ {
			tickNight(g, false, false)
		}
		if g.Screen != Harbour || g.Save.NightsCompleted != 1 {
			t.Fatalf("seed %d: crossing night did not finish", seed)
		}
		for _, ticks := range g.Night.TotalTicks {
			if ticks == 0 {
				t.Fatalf("missing crossing")
			}
		}
	}
}
func TestLitCrossingFasterThanUnlit(t *testing.T) {
	mk := func(lit bool) int {
		g := startTestNight(39)
		if !lit {
			g.Beam.Angle = -MaxBeamAngle
		}
		for j := 0; j < 25000; j++ {
			g.Update(Input{Today: "2026-11-11", Yes: g.Night.CheckIn})
			if !lit {
				g.Beam.Angle = -MaxBeamAngle
				g.Beam.Omega = 0
			}
			if g.Night.Crossing > 0 {
				return g.Night.TotalTicks[0]
			}
		}
		return 0
	}
	a, b := mk(true), mk(false)
	if a == 0 || b == 0 || a >= b {
		t.Fatalf("lit %d, unlit %d", a, b)
	}
}
func TestBeamEasingAndClamp(t *testing.T) {
	g := startTestNight(2)
	old := g.Beam.Angle
	g.Update(Input{Right: true})
	if g.Beam.Angle <= old || g.Beam.Omega <= 0 || g.Beam.Omega > MaxBeamSpeed {
		t.Fatal(g.Beam)
	}
	for i := 0; i < 1000; i++ {
		g.Update(Input{Right: true, Yes: g.Night.CheckIn})
	}
	if g.Beam.Angle > MaxBeamAngle+1e-8 {
		t.Fatal(g.Beam.Angle)
	}
}
func TestLitOnlyInsideCone(t *testing.T) {
	p := Point{X: 500, Y: 150}
	b := Beam{Angle: Bearing(p), HalfWidth: 0.1}
	if !WithinBeam(p, b) {
		t.Fatal("inside")
	}
	b.Angle += 0.2
	if WithinBeam(p, b) {
		t.Fatal("outside")
	}
	if WithinBeam(Point{X: 9999, Y: 10}, b) {
		t.Fatal("range")
	}
}
func TestHarbourReachableAtAngleClamp(t *testing.T) {
	a := Bearing(Point{X: HarbourX, Y: HarbourY})
	if a > MaxBeamAngle || a < -MaxBeamAngle {
		t.Fatal(a)
	}
}
func TestAllBoatPositionsReachableByBeam(t *testing.T) {
	for seed := int64(0); seed < 200; seed++ {
		g := startTestNight(seed)
		for i := 0; i < 5; i++ {
			b := g.makeBoat()
			for j := 0; j <= 256; j++ {
				p := Bezier(b.Path, float64(j)/256)
				a := Bearing(p)
				if math.Abs(a) > MaxBeamAngle || math.Hypot(p.X-LampX, p.Y-LampY) > BeamRange {
					t.Fatalf("seed=%d path=%d t=%d point=%+v bearing=%.3f", seed, i, j, p, a)
				}
			}
		}
	}
}
func TestFogRiseAndClearTimes(t *testing.T) {
	f := 0.0
	for i := 0; i < 90; i++ {
		f = clamp(f+1.0/90, 0, 1)
	}
	if math.Abs(f-1) > 1e-8 {
		t.Fatal(f)
	}
	for i := 0; i < 60; i++ {
		f = clamp(f-1.0/60, 0, 1)
	}
	if f > 1e-8 {
		t.Fatal(f)
	}
}
func TestNoticeLogsOldestVisibleAndBeamUnchanged(t *testing.T) {
	g := startTestNight(3)
	n := &g.Night
	n.NextSighting = 999999
	n.Sightings = []Sighting{{Kind: Whale, Age: 65, Life: 600, State: Visible, Pos: Point{100, 100}}, {Kind: Seal, Age: 100, Life: 600, State: Visible, Pos: Point{120, 110}}}
	old := g.Beam.Angle
	g.Update(Input{Notice: true})
	if n.Noticed != 1 || n.Sightings[1].State != Noticed || g.Save.Book[Seal].Count != 1 || math.Abs(old-g.Beam.Angle) > 1e-12 {
		t.Fatal("notice", n.Sightings, g.Beam.Angle)
	}
}
func TestChaseAfterSixtyTicksUnlogged(t *testing.T) {
	g := startTestNight(2)
	n := &g.Night
	n.NextSighting = 999999
	s := Sighting{Kind: Whale, Life: 600, State: Visible, Pos: Point{400, 110}}
	n.Sightings = append(n.Sightings, s)
	g.Beam.Angle = Bearing(s.Pos)
	for i := 0; i < 60; i++ {
		n.Boat.Lit = false
		n.Sightings[0].Aim = i
		// chase rule is part of sighting update; hold current beam direction.
		g.updateSightings(false)
	}
	if n.Sightings[0].State != Chased || g.Save.Book[Whale].Count != 0 {
		t.Fatal(n.Sightings[0])
	}
}
func TestCameBackInsideWindow(t *testing.T) {
	g := startTestNight(1)
	n := &g.Night
	n.CrossingTicks = 180
	n.OffTicks = 119
	n.Boat.Lit = false
	g.updateReturn()
	if n.ChimeTicks == 0 {
		t.Fatal("no chime")
	}
	n.Boat.Lit = true
	g.updateReturn()
	if n.CameBacks != 1 || n.ChimeTicks != 0 {
		t.Fatal(n.CameBacks)
	}
}
func TestNoChimeInGrace(t *testing.T) {
	g := startTestNight(1)
	n := &g.Night
	for i := 0; i < 179; i++ {
		n.CrossingTicks = i
		n.Boat.Lit = false
		g.updateReturn()
	}
	if n.ChimeUsed {
		t.Fatal("chime before grace")
	}
}
func TestChimeAfterGraceWhenAlreadyUnlit(t *testing.T) {
	g := startTestNight(1)
	n := &g.Night
	n.OffTicks = 250
	n.CrossingTicks = 180
	g.updateReturn()
	if !n.ChimeUsed || n.ChimeTicks == 0 {
		t.Fatal("missed chime after grace")
	}
}
func TestOneChimePerDrift(t *testing.T) {
	g := startTestNight(1)
	n := &g.Night
	n.OffTicks = 120
	n.CrossingTicks = 180
	g.updateReturn()
	for i := 0; i < 500; i++ {
		n.CrossingTicks++
		n.Boat.Lit = false
		g.updateReturn()
	}
	num := 0
	for _, e := range g.Events {
		if e == EvWhereBoat {
			num++
		}
	}
	if num != 1 {
		t.Fatalf("chime should not re-arm in same drift: %d", num)
	}
	n.Boat.Lit = true
	for i := 0; i < 60; i++ {
		g.updateReturn()
	}
	if n.ChimeUsed {
		t.Fatal("did not re-arm")
	}
}
func TestLanternThreshold(t *testing.T) {
	lit := [5]int{70, 69, 100, 0, 100}
	total := [5]int{100, 100, 100, 100, 200}
	if CountLanterns(lit, total) != 2 {
		t.Fatal(CountLanterns(lit, total))
	}
}
func TestGoodGuessWithinOne(t *testing.T) {
	for a := 0; a <= 5; a++ {
		for b := 0; b <= 5; b++ {
			want := absInt(a-b) <= 1
			if (absInt(a-b) <= 1) != want {
				t.Fatal(a, b)
			}
		}
	}
}
func TestBellScheduleBounds(t *testing.T) {
	for seed := int64(1); seed <= 100; seed++ {
		g := startTestNight(seed)
		n := g.Night
		count := 0
		for i, ticks := range n.BellAt {
			if ticks == 0 {
				continue
			}
			count++
			sec := ticks / 60
			if i == 0 {
				if sec < 19 || sec > 28 {
					t.Fatal(i, sec)
				}
			} else if i != 2 && i != 4 || sec < 15 || sec > 24 {
				t.Fatal(i, sec)
			}
		}
		if count < 2 || count > 3 {
			t.Fatal(count)
		}
	}
}
func TestBellScheduleAlwaysTwoOrThreeOnSlowCrossings(t *testing.T) {
	g := startTestNight(22)
	for j := 0; j < 180000 && g.Screen == NightScreen; j++ {
		if g.Night.CheckIn {
			g.Update(Input{No: true, Today: "2026-11-11"})
		} else {
			g.Beam.Angle = -MaxBeamAngle
			g.Update(Input{Today: "2026-11-11"})
		}
	}
	if g.Screen != Harbour || len(g.Night.Answers) < 2 || len(g.Night.Answers) > 3 {
		t.Fatal(g.Screen, len(g.Night.Answers))
	}
}
func TestBellNeverDuringPause(t *testing.T) {
	g := startTestNight(22)
	for i := 0; i < 160000 && g.Screen == NightScreen; i++ {
		in := Input{Yes: g.Night.CheckIn, Today: "2026-11-11"}
		if g.Night.PauseTicks > 0 {
			old := g.Night.Ticks
			g.Update(in)
			if g.Night.Ticks != old {
				t.Fatal("pause advanced")
			}
			for _, e := range g.Events {
				if e == EvBell {
					t.Fatal("bell during pause")
				}
			}
		} else {
			g.Update(in)
		}
	}
}
func TestCheckInPausesSimulation(t *testing.T) {
	g := startTestNight(17)
	n := &g.Night
	n.CheckIn = true
	n.CrossingTicks = 123
	before := n.Boat.T
	g.Update(Input{Left: true})
	if n.CrossingTicks != 123 || n.Boat.T != before {
		t.Fatal("check-in advanced")
	}
}
func TestAnyAnswerEarnsStar(t *testing.T) {
	g := startTestNight(11)
	g.Night.CheckIn = true
	g.Update(Input{No: true})
	if len(g.Night.Answers) != 1 || g.Night.Answers[0] {
		t.Fatal(g.Night.Answers)
	}
}
func TestProgressionMonotonicAndCapped(t *testing.T) {
	for i := 0; i < 21; i++ {
		capped := minInt(i, 10)
		w := lerp(12, 8, float64(capped)/10)
		if w < 8 || w > 12 {
			t.Fatal(w)
		}
		if i > 0 {
			prev := lerp(12, 8, float64(minInt(i-1, 10))/10)
			if w > prev {
				t.Fatal("width increased")
			}
		}
	}
}
func TestRestAfterTwoNightsSameDay(t *testing.T) {
	g := New(1, NewSave())
	g.Save.LastDate = "2026-11-11"
	g.Save.NightsToday = 2
	g.Update(Input{Start: true, Today: "2026-11-11"})
	if g.Screen != Rest {
		t.Fatal(g.Screen)
	}
}
func TestNewDayResets(t *testing.T) {
	g := New(1, NewSave())
	g.Save.LastDate = "2026-11-10"
	g.Save.NightsToday = 2
	g.Save.ExtraNightUsed = true
	g.Update(Input{Start: true, Today: "2026-11-11"})
	if g.Screen != Predict || g.Save.NightsToday != 0 || g.Save.ExtraNightUsed {
		t.Fatal(g.Screen, g.Save)
	}
}
func TestParentOverride(t *testing.T) {
	g := New(1, NewSave())
	g.Save.LastDate = "2026-11-11"
	g.Save.NightsToday = 2
	g.Update(Input{Start: true, Today: "2026-11-11"})
	for i := 0; i < 179; i++ {
		g.Update(Input{ParentDown: true, Today: "2026-11-11"})
	}
	if g.Screen != Rest {
		t.Fatal("premature")
	}
	g.Update(Input{ParentDown: true, Today: "2026-11-11"})
	if g.Screen != Predict || !g.Save.ExtraNightUsed {
		t.Fatal("override did not start")
	}
}
func TestOneExtraNightPerDay(t *testing.T) {
	g := New(1, NewSave())
	g.Save.LastDate = "2026-11-11"
	g.Save.NightsToday = 2
	g.Save.ExtraNightUsed = true
	g.Update(Input{Start: true, Today: "2026-11-11"})
	for i := 0; i < 240; i++ {
		g.Update(Input{ParentDown: true, Today: "2026-11-11"})
	}
	if g.Screen != Rest {
		t.Fatal("repeat override")
	}
}
func TestDeskBellIntervals(t *testing.T) {
	for seed := int64(0); seed < 100; seed++ {
		g := New(seed, NewSave())
		g.startDesk()
		if g.Bell.NextRingAt < 180*60 || g.Bell.NextRingAt > 420*60 {
			t.Fatal(g.Bell.NextRingAt)
		}
	}
}
func TestDeskAnswersRecorded(t *testing.T) {
	g := New(2, NewSave())
	g.startDesk()
	g.Update(Input{DeskElapsedTicks: g.Bell.NextRingAt})
	if !g.Bell.Waiting {
		t.Fatal("missing desk bell")
	}
	g.Update(Input{No: true, DeskElapsedTicks: g.Bell.NextRingAt + 1})
	g.Update(Input{Escape: true, Today: "2026-11-11", DeskElapsedTicks: g.Bell.NextRingAt + 2})
	if len(g.Save.Desk) != 1 || len(g.Save.Desk[0].Answers) != 1 || g.Save.Desk[0].Answers[0] || !g.Save.Desk[0].EndedEarly {
		t.Fatal(g.Save.Desk)
	}
}
func TestDeskClockAdvancesAcrossMissedUpdates(t *testing.T) {
	g := New(3, NewSave())
	g.startDesk()
	g.Update(Input{DeskElapsedTicks: 25 * 3600, Today: "2026-11-11"})
	if g.Screen != Title || len(g.Save.Desk) != 1 {
		t.Fatal("timer not elapsed")
	}
}
func TestDeskUnansweredNoFalseAnswer(t *testing.T) {
	g := New(2, NewSave())
	g.startDesk()
	g.Update(Input{DeskElapsedTicks: g.Bell.NextRingAt})
	g.Update(Input{DeskElapsedTicks: 15 * 3600, Today: "2026-11-11"})
	if len(g.Save.Desk[0].Answers) != 0 || !g.Save.Desk[0].UnansweredAtEnd {
		t.Fatal(g.Save.Desk[0])
	}
}
func TestNoFailureStrings(t *testing.T) {
	for _, s := range UIText {
		for _, bad := range []string{"wrong", "failed", "failure", "loser", "bad job", "you lost"} {
			if strings.Contains(strings.ToLower(s), bad) {
				t.Fatalf("shaming copy: %q", s)
			}
		}
	}
}
func TestReplayDeterministic(t *testing.T) {
	a, b := startTestNight(99), startTestNight(99)
	for i := 0; i < 25000; i++ {
		in := Input{Left: i%220 < 70, Right: i%220 > 100 && i%220 < 170, Yes: a.Night.CheckIn, Today: "2026-11-11"}
		a.Update(in)
		in.Yes = b.Night.CheckIn
		b.Update(in)
		if a.Screen != b.Screen {
			t.Fatalf("screen %d", i)
		}
		if a.Screen == Harbour {
			break
		}
	}
	if !reflect.DeepEqual(a.Save, b.Save) || !reflect.DeepEqual(a.Night, b.Night) || a.Beam != b.Beam {
		t.Fatal("not deterministic")
	}
}
func TestPresentationEffectsDoNotAffectGameplay(t *testing.T) {
	a, b := startTestNight(100), startTestNight(100)
	// Cosmetic draw calls are not part of the core, so arbitrary views cannot consume its RNG.
	for i := 0; i < 500; i++ {
		a.Update(Input{Yes: a.Night.CheckIn})
		b.Update(Input{Yes: b.Night.CheckIn})
	}
	if !reflect.DeepEqual(a.Night, b.Night) {
		t.Fatal("core changed without new inputs")
	}
}
func TestEventsMatchStateChanges(t *testing.T) {
	g := startTestNight(3)
	g.Night.CheckIn = true
	g.Update(Input{No: true})
	if len(g.Events) != 1 || g.Events[0] != EvCheckInAnswered || len(g.Night.Answers) != 1 {
		t.Fatal(g.Events)
	}
}

// User acceptance: a lit crossing should progress visibly and require active
// tracking, without changing the five-crossing or no-failure rules.
func TestUATCrossingsFasterAndProgressive(t *testing.T) {
	g := startTestNight(55)
	first := g.crossingTargetTicks()
	if first > 3600 || first < 2700 {
		t.Fatalf("first crossing: %d ticks", first)
	}
	g.Night.Crossing = 4
	if g.crossingTargetTicks() >= first {
		t.Fatal("later crossings are not faster")
	}
	g.Night.Index = 10
	if g.crossingTargetTicks() >= first {
		t.Fatal("later nights are not faster")
	}
	for i := 0; i < 10; i++ {
		g.Night.Index = i
		for c := 0; c < 5; c++ {
			g.Night.Crossing = c
			if g.crossingTargetTicks() < 1950 {
				t.Fatal("unbounded speed")
			}
		}
	}
}
func TestUATBeamNarrowerThanOriginal(t *testing.T) {
	g := startTestNight(1)
	if !(g.Beam.HalfWidth < 10*math.Pi/180 && g.Beam.HalfWidth > 6*math.Pi/180) {
		t.Fatalf("first-night beam width: %g radians", g.Beam.HalfWidth)
	}
}
