# Lighthouse Keeper — Engineering Design Document

**Version:** 0.2 (Garage simplification review + premium visual design)
**Status:** Locked engineering baseline for V0.1.0; initial implementation candidate, integration/performance/UAT evidence pending
**Implementation:** Go 1.25+, Ebitengine 2.10.6 (pinned)
**Product source:** Lighthouse Keeper PRD v0.2 (FR1–FR16, AC1–AC21)
**Review source:** https://github.com/jonathanselby/engineering-under-constraints — `METHODOLOGY.md`, `principles/modern-software.md`, `principles/constraint-review.md`, `c64/lessons.md`

**Version choice:** 2.10.6 includes a fix for uneven update timing near the TPS or its multiples (Ebitengine's 7 November 2026 release). It is preferred to 2.10.5 specifically because frame pacing is a product requirement.

## 1. Purpose

Lighthouse Keeper is one small calm game. A beam follows a boat; sightings tempt the beam away; noticing, coming back and answering honestly are what the game rewards.

> **Keep the game as small as the game actually is.**

The behaviour that matters comes from five interacting rules: the beam lights the boat, an unlit boat fogs and slows, sightings appear at the edges, noticing logs them while chasing loses them, and a drift triggers a chance to come back. Prediction, check-ins and Keeper's Bell are thin layers of recording around that loop.

## 2. Design constraints

1. **Mission before architecture.** Full PRD scope and premium appearance are hard requirements; only redundant implementation complexity may be removed.
2. **One authoritative `Game` state.** Rendering, music, particles and viewport caches are disposable projections, never business/game truth.
3. **One explicit, single-threaded 60 TPS gameplay loop.** Input, seeded RNG and injected local date / elapsed time determine the outcome. No goroutines, runtime AI or network.
4. **Separate TPS from FPS.** Render as often as the display supports (60/120 Hz where proven), interpolating presentation-only values; no extra gameplay ticks.
5. **Finite budgets, known behaviour at limits.** 5 crossings, 6 kinds, 3 simultaneously visible sightings, finite log length, bounded effects and file reads.
6. **No abstraction before its need.** A separate render *file* earns its place; a scene graph, renderer framework, ECS, physics, tween or particle engine does not.
7. **Interfaces first.** Game -> read-only state/events -> presentation. Render never feeds values back into gameplay.
8. **Degrade gracefully.** Optional glints, spray and ornament may be reduced before essential visual clarity or smoothness.
9. **Test the real environment.** Headless rule tests plus actual 60/120 Hz, Hi-DPI, full-night, accessibility/readability and child play-tests.
10. **Production evidence.** PRD, EDD, `TESTING.md`, `RELEASE.md` with genuinely verified ticks, independently recovered EDD and comparison pinned to the release commit.

Implementation restriction: Go standard library plus Ebitengine, and its ordinary transitive module dependencies, only. Native Linux graphics/audio libraries required by Ebitengine are not counted as added services. No network, telemetry, accounts, engine plugins or runtime art generation.

**Review outcome:** keep the five core rules; simplify duplicated state, bounded collections, bell scheduling and cosmetic ownership; fix unreachable beam geometry and grace-window chimes; test all performance claims on hardware. Details and retained risks are in `Lighthouse_Keeper_Garage_Simplification_Review_v0.2.md`. If a new component is proposed, apply the Garage removal test first.

## 3. Runtime dependency

```text
               OFFLINE / DEVELOPMENT ONLY
    author layered source art + animation + sound
                        |
           optimise, crop, pack, verify
                        v
              shipped local assets
                        |
                        v
          +--------------------------------+
          | cmd/lighthousekeeper           |
          | window / input / time adapter   |
          | rendering / sound / save I/O    |
          +----------------+---------------+
                           | Input / events / readonly values
                           v
                 internal/game
         pure rules; no images, I/O or wall clock
```

The gameplay package never imports Ebitengine. The adapter creates all textures and decodes all essential assets during loading, before gameplay; it receives deterministic gameplay state/events and presents them. A rendered screenshot or visual tween is not an input to any rule. One source image may be reused in many placements via Ebitengine transforms; no custom rendering framework is introduced.

## 4. Units and geometry

| Quantity | Unit |
| --- | --- |
| time | ticks (60 per second) |
| position | logical screen pixels, 960×720, origin top-left |
| angle | radians, 0 = straight up from the lighthouse, positive clockwise |

```text
lighthouse lamp  (480, 640)
sea              0 <= y < 600
harbour mouth    (820, 590)
beam range       820 px
beam angle clamp [-1.50, +1.50] rad
```

## 5. Authoritative state

```go
type Game struct {
    Screen Screen // Title, Predict, Night, Harbour, Book, Rest, BellSetup, Bell
    Night  Night
    Bell   BellSession
    Save   Save
    Beam   Beam
    Events []Event // reused backing buffer, cleared at each Update
    RNG    *rand.Rand // gameplay only
    ParentHoldTicks int // active on Rest screen only
}

type Beam struct {
    Angle, Omega float64
    HalfWidth    float64
}

type Night struct {
    Index, Prediction, Crossing int  // Crossing: 0..4
    Boat       Boat
    Fog        float64
    LitTicks   [5]int // known fixed five crossings
    TotalTicks [5]int
    Sightings  []Sighting // capacity bounded; max 3 visible
    NextSighting, OffTicks, LitStreak, ChimeTicks int
    ChimeUsed bool
    CameBacks, Noticed int
    BellAt        [5]int // crossing-local active ticks, 0 = no bell
    CheckIn       bool
    Answers       []bool // at most 3
    Ticks         int   // active play ticks, excludes crossing pauses/check-ins
    CrossingTicks int
    PauseTicks    int   // 120 between crossings
}
```

Removed from the core: duplicated `Game.Today` (use `Input.Today` at state transitions), `CosmeticRNG`, and `Particles`. The adapter alone owns transient decorative particle/sound state, which may be discarded without changing outcomes. Five-element arrays are simpler and allocate less than slices for permanent, known-size crossing records. Keep the six-entry persisted Book as a slice so its decoded length can be validated exactly (a fixed Go array could silently discard extra JSON array elements). An array of five crossing-local bell triggers replaces a mutable global-timeline schedule. The current screen and `Night.CheckIn` remain explicit, not inferred from animation.

Presentation consumes the latest read-only rule state immediately. If higher-refresh motion needs it, the adapter can cache a previous pose or use bounded sub-tick forward prediction of visual positions; no buffered one-tick-late display is permitted. These are disposable view values, not a second game state; all decisions use current core state.

## 6. The boat and crossing

```go
type Boat struct {
    Path [4]Point // cubic Bézier from sea entry to harbour mouth
    T    float64  // progress 0..1
    Lit  bool
}
```

A crossing generates a cubic Bézier with entry on the left, top or right sea edge and two control points whose lateral offset grows with progression (section 12). Constrain entry, control points and harbour endpoint to the sea rectangle; the convex-hull property then keeps the entire curve within beam range and rotation. This is a bounded construction check, not a runtime path-finding system. Keep the existing 64-entry cumulative arc-length lookup: it cheaply approximates constant speed and preserves the softly curved path. Advance by distance along the lookup, not by uniformly adding to Bézier parameter `T` (which would vary speed). The boat and harbour positions are gameplay coordinates; layered artwork must be composed around them.

**Geometry guard:** v0.1's harbour `(820,590)` is at ~81.6° from the lamp, beyond its ±77.35° clamp. The extreme sea corners are ~800 px from the lamp, beyond the former 760 px range. Widen only the permitted pointing envelope to ±1.50 rad (~85.9°) and range to 820 px. That covers the complete bounded sea rectangle, including the existing harbour, without changing the path, progression or skill mechanic. Test every generated point and full curve (including extrema), at first and maximum progression. The beam visual must use these same gameplay geometry parameters, never a separate artist-authored hitbox.

```text
full speed  = path length / (90 s * 60)     per tick
drift speed = 0.35 * full speed
speed       = lerp(full, drift, Fog)
```

The crossing ends when `T >= 1`. It always ends: drift speed is positive.

### Lit test

```text
bearing = angle from lamp to boat
Lit = |wrap(bearing - Beam.Angle)| <= Beam.HalfWidth  and  distance <= 820
```

## 7. Beam

```text
input: Left held -> target = -OmegaMax; Right -> +OmegaMax; neither or both -> 0
Omega += (target - Omega) * 0.2
Angle += Omega; clamp to [-1.50, +1.50] rad (all reachable sea and harbour)
```

`OmegaMax = 0.025` rad/tick (about 86°/s). Easing makes the beam feel heavy and calm rather than twitchy.

## 8. Fog

```text
unlit: Fog = min(1, Fog + 1/90)   (fully foggy after 1.5 s)
lit:   Fog = max(0, Fog - 1/60)   (clear within 1 s)
```

Fog is gameplay (it sets speed) and is also drawn.

## 9. Sightings

```go
type Kind uint8 // Whale, ShootingStar, Seal, Jellyfish, Owl, Ferry

type Sighting struct {
    Kind     Kind
    Pos      Point
    Age, Life int
    State    SightState // Appearing, Visible, Noticed, Chased, Leaving
    Aim      int        // consecutive ticks the beam has covered it
}
```

### Spawning

The next spawn tick is drawn from the gameplay RNG: 12–25 s on the first night, narrowing to 8–16 s at the progression cap. Position is a point in the sea-edge band (within 140 px of the left, right or top edge), at least 160 px from the boat's path. Life is 6–9 s. At most three sightings are *visible* at once; exiting animation records are also bounded (maximum six total). Kind weights: Whale, Seal, Jellyfish common; Owl, Ferry less common; ShootingStar rare.

Use a fixed small number of seeded position attempts; if none meets the clearance rule, skip that spawn and draw the next interval. Do not retry indefinitely or force an invalid location. Skipped spawns do not change the six collectible kinds or any rewards. The artist supplies original, distinct frames for each kind; animation does not change its coordinates or spawn timing.

### Notice

`N` just pressed, not in check-in, with at least one sighting `Visible`: the oldest `Visible` sighting becomes `Noticed`, its kind's count in the book increments, emit `Notice` (and `NewKind` if the count became 1). The beam does not move. A noticed sighting plays a happy animation for 60 ticks and leaves.

### Chase

Each tick a `Visible` sighting is inside the beam cone while the boat is unlit, `Aim++`; otherwise `Aim = 0`. At `Aim >= 60` it becomes `Chased`, leaves shyly, is not logged, emit `SightingLeft`.

## 10. Coming back

```text
grace: no chime until CrossingTicks >= 180
unlit: OffTicks++, LitStreak = 0
lit:   OffTicks = 0, LitStreak++
if OffTicks >= 120 and CrossingTicks >= 180 and !ChimeUsed:
    ChimeTicks = 300; ChimeUsed = true; emit WhereBoat
while ChimeTicks > 0: if lit -> CameBacks++, emit CameBack, ChimeTicks = 0
ChimeTicks counts down on active crossing ticks; reaching 0 closes silently
if LitStreak >= 60: ChimeUsed = false (a new drift episode may earn another)
```

Nothing is emitted when a window closes without a return. **Correction:** `OffTicks == 120` in v0.1 could be passed during the first 180 grace ticks; the chime then never occurred for that drift. `>= 120` together with the grace condition guarantees the promised opportunity. The 2-second threshold, 5-second return window and 1-second re-arm remain unchanged.

## 11. Prediction, lanterns and good guess

The Predict screen precedes each night: Left/Right choose 0–5, Space confirms. Default is the previous night's result, or 3.

```text
crossing lantern lit = LitTicks[i] / TotalTicks[i] >= 0.70
result               = count of lit lanterns (0..5)
goodGuess            = |Prediction - result| <= 1
```

## 12. Progression

`Night.Index = min(Save.NightsCompleted, 10)`, interpolated linearly from first-night values to caps:

| Parameter | Night 0 | Cap (night 10) |
| --- | --- | --- |
| path lateral offset | 60 px | 220 px |
| sighting interval | 12–25 s | 8–16 s |
| beam half-width | 12° | 8° |

Nothing earned is ever removed.

## 13. Check-in bell

At night start choose **two or three distinct bell slots** from crossing indices `{0,2,4}` using the gameplay RNG. If two, randomly select which two of the three. For crossing 0, the bell falls at a random active crossing tick from 60–75 s; for crossing 2 or 4, at 20–50 s. Store tick or 0 in `Night.BellAt[5]`. A full-speed crossing lasts 90 s; every slot is therefore reached even when a later crossing slows. The shortest gap across two candidate slots exceeds 90 s; no slot overlaps the 2-second between-crossing pause. This replaces the need to forecast night end time or reschedule bells when fog changes boat speed.

When `CrossingTicks == BellAt[Crossing]`, `CheckIn = true`, emit `Bell`. The simulation pauses (boat, beam, fog, sightseeing clocks and crossing time do not advance). `Y` or `N` records the answer, emits `CheckInAnswered`, adds a star, then resumes. Only `Y`/`N` is live in this overlay. The overlay can continue a gentle render-only glow.

The player cannot infer the exact time or count; no two playthroughs need use the same chosen slots/offsets unless their seed is the same. The core records only her answer; it does not calculate or store whether she "was" with the boat, or contradict that answer.

## 14. Rest and daily counts

`Input.Today` is the adapter-supplied local calendar date `YYYY-MM-DD` (empty only when the clock is unavailable). At a night-end / night-start / Rest transition, if `Save.LastDate != Today`, reset `NightsToday = 0` and `ExtraNightUsed = false`; update `LastDate`. On successful night end, increment `NightsToday`. With two completed nights, Title routes Start to Rest. On Rest, `Input.ParentDown` held for 180 active ticks **immediately starts** the third night at Predict and atomically marks `ExtraNightUsed = true` at the transition; do not persist a separate unused grant. Repeated overrides that day are ignored. The book and log remain browseable even during Rest.

Date missing: disable only the daily restriction for that process, not saved history or other gameplay. The adapter does not expose local file paths or date editing through in-game input. On relaunch the persisted `ExtraNightUsed` and count preserve the one-extra-night limit. A relaunch before the third night starts consumes no grant; a third night started and abandoned still consumes that day's one extra opportunity.

## 15. Keeper's Bell mode

```go
type BellSession struct {
    LengthTicks int // 15 or 25 minutes, expressed in 60 Hz units
    Elapsed     int // injected real elapsed time, not number of Update calls
    NextRingAt  int // one pending deadline, never an unbounded queue
    Answers     []bool // bounded by practical session length
    Waiting     bool
    EndedEarly  bool
}
```

At session start select `NextRingAt` in 180–420 s (seeded). For each answered ring, select the next one in 180–420 s after that answer; at no point keep more than one pending question. The screen is a still lighthouse with gently animated lamp glow. A ring emits `DeskBell` and shows `Were you with the boat?` with `Y` / `N`; only her response is recorded. An unanswered question stays visible until answered or session end, while elapsed real time continues; no backlog of bells fires on resumption. At 15/25 minutes emit `BellEnd` and append the desk log, including any pending unanswered state and `EndedEarly` if Esc was used. Never invent an answer.

**Time injection:** the adapter supplies nondecreasing `Input.DeskElapsedTicks` derived from *real elapsed time* since the session started; `Game.Update` never calls the system clock. Do **not** calculate the homework timer by adding one on each Update: lost-focus throttling or OS suspension would make a 25-minute session run long. `SetRunnableOnUnfocused(true)` permits normal background updates; on restoration after suspension, compare the injected elapsed time to the deadline and process **at most one** overdue bell or the session end (end wins). Large backwards clock adjustments are clamped and observed in diagnostics, not turned into negative elapsed time. Timer accuracy across actual suspend/resume and clock changes is a release test, not assumed from this design alone. Provide a windowed mode for desk use.

## 16. Persistence

One local `save.json`:

```go
type Save struct {
    Version         int
    Book            []BookEntry // exactly six in fixed Kind order; validate length
    NightsCompleted int
    LastDate        string
    NightsToday     int
    ExtraNightUsed  bool
    Log             []NightRecord // newest last, capped at 400
    Desk            []DeskRecord  // newest last, capped at 400
}

type NightRecord struct {
    Date       string
    Prediction, Result int
    GoodGuess  bool
    CameBacks, Noticed int
    CheckIns   []bool
}

type DeskRecord struct {
    Date            string
    LengthMinutes   int // 15 or 25
    Answers         []bool
    UnansweredAtEnd bool
    EndedEarly      bool
}
```

- Store in the user-local application configuration directory as `save.json` (not an arbitrary current working directory), with directory permissions restricted to the user where supported. Save via a same-directory temporary file, flush and atomic rename. Write permissions restrict other local users; do not introduce network storage.
- Missing file: start empty. Malformed, oversized, unknown-version, unknown-field or structurally invalid: display a gentle recovery message, retain a `.bad` copy without overwriting an earlier recovery copy, and start a new in-memory save; do not crash or trust unchecked counts. The parent can inspect the recovery copy.
- **Bound before decoding:** read at most 256 KiB plus a byte to detect overlength, then `json.Decoder.DisallowUnknownFields`, reject trailing data, validate `Version` (2 for this unreleased design; no earlier released-save migration required) and type/length/range invariants (six known kinds in order; nonnegative finite counts; nights and results 0..5; at most 3 check-ins per night; valid date shape; 400-record caps; bounded desk-answer count). Never allocate based on unchecked values or untrusted JSON counts.
- Cap slices in memory before writing and compare output size to the 200 KiB save-file target. If it would exceed the bound, trim oldest bounded history entries (never the book or current record) before writing. No JSON maps or unordered iteration in persistence.
- Cache decoded presentation assets at startup. **Never perform a save write, image decode or texture upload from the `Draw` path.** Save at the specified low-frequency transition boundaries, where a short synchronous atomic write is simpler and avoids introducing a worker/queue; measure the worst transition hitch.

## 17. Screens and input

```text
TITLE      Space: night   B: book   L: log   K: Keeper's Bell   Esc: quit
PREDICT    Left/Right choose, Space confirm
NIGHT      Left/Right beam, N notice; check-in overlay: Y/N
HARBOUR    end-of-night log page; Space continues
BOOK       Left/Right turn pages
REST       lighthouse at rest; P held 3 s: parent override
BELLSETUP  Left/Right 15 or 25 min, Space start
BELL       Y/N when asked; Esc ends early (recorded as ended early)
```

```go
type Input struct {
    Left, Right, ParentDown bool // held
    LeftPressed, RightPressed bool // key edges for picker / page turning
    Notice, Yes, No, Start, Book, Log, Bell, Escape bool // key edges
    Today            string // adapter-injected local YYYY-MM-DD or empty
    DeskElapsedTicks int    // adapter-injected real elapsed time; Bell only
}
```

Esc quits from any screen except BELL, where it ends the session first.

## 18. Update loop (Night)

```text
reuse Events buffer, set length 0
if CheckIn: accept Y or N only; emit answer; return
if PauseTicks > 0: PauseTicks--; return
Ticks++; CrossingTicks++
if CrossingTicks == BellAt[Crossing]: CheckIn = true; emit Bell; return
beam easing + angle; boat bearing/lit test; fog and boat advance
per-crossing TotalTicks and LitTicks; coming-back threshold/window
sightings: bounded spawn, age, aim/chase/leave; N notices oldest visible
crossing end -> record lantern, then start next with PauseTicks=120
if fifth crossing ends -> calculate prediction/result, append Save record and emit EvNightEnd
adapter writes Save at the transition, never inside the pure game package
emit bounded events for adapter; never generate particles in Game
```

## 19. Events and sound

```go
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
)
```

Sound is soft throughout: an ambient wave loop, a lantern click on a **state transition** to lit (not every illuminated tick), a gentle chime for "where's the boat?", a sparkle for came back, a wave-and-giggle for notice, a page turn for a new kind, a ship's bell for check-ins, and harbour bells at crossing end. No alarm sounds. Use a small number of original, volume-matched short clips; short WAV is acceptable, while a compact Ogg/Vorbis ambient loop may be streamed with Ebitengine's included audio/vorbis package if it saves material size. Load/decode or prepare playback outside `Draw`, and fail silently if audio is unavailable. Sound consumes events and never gates gameplay.

## 20. Premium rendering — small mechanisms, high artistic finish

### 20.1 Art direction and composition

The five November 2026 sample images are **approved ambition and mood references**, not final level artwork. Select one coherent finished visual treatment using their shared strengths: layered painterly detail, navy/teal ocean, indigo-violet sky, amber lamplight, pearl moon, soft mist, rock/flower framing and a warmly welcoming harbour. Do not add five skin modes. Produce original, rights-cleared layers rather than shipping a full-screen AI concept image as a game background.

The stage remains **4:3 with a fixed 960 × 720 logical coordinate system**. Place the painterly lighthouse lantern exactly over the core lamp `(480,640)`, keep the sea and boat path clear, and put detailed architecture/flowers on the margins. A composition originally painted at 16:9 must be recomposed for 4:3, not squashed or centre-cropped until controls are obscured. Boat, beam, six sightings, five lanterns and check-in question have priority over decorative detail.

### Implementation clarification — painterly corrective pass

The initial vector-like graphical pass was rejected after actual capture. The corrective implementation replaces it with three original source paintings composed **offline** by `assets-source/export_premium.py` into the sea, lantern, boat, six sighting sprites and reused UI/light assets. This preserves the runtime's sole rendering dependency, Ebitengine; Python, Pillow and NumPy are **authoring-only** tools. The final result is still subject to owner/hardware acceptance.

The authoritative scene remains 960 × 720. To display all boat positions against a predominantly painted ocean, the renderer applies the presentation-only map `visualY = 124 + 0.79 × logicalY`. Beam texture coordinates receive the identical affine projection *after* rotation, as do sightings, fog and particles. In particular, the lighthouse's visible lighthead is centred at `visualY(LampY)` instead of assuming a source crop centre is the gameplay beam origin. This changes no rule, crossing, bearing, reward or replay value.

The twelve runtime art files are listed in `ARTWORK.md`. The six sightseeing paintings are atlas regions of 240 × 240 pixels. A single UI rail is sliced at **loading**, not per frame, and reused as nine image draws per panel. The initial water-curl renderer is deleted; sixteen faint translated shimmer draws are bounded visual work and can be omitted on slower hardware. The boat rocking is a bounded render-only angle, not a physics update. No Kage shader or new runtime framework is introduced.

This is an implementation-level clarification of §20's requirement for painterly animation and no added gameplay latency. The remaining mandatory evidence is CI success on the changed build, graphical screenshots and real 60/120 Hz frame measurements. A static preview does not close those gates.

### 20.2 Layer and asset plan

| Visual layer | Shipped form | Runtime mechanism |
| --- | --- | --- |
| Sky, distant cliffs, moon, town (unlit) | 1–3 compressed opaque background images | Draw static once per frame, optional gentle subpixel drift |
| Water surface | 2–3 painted repeatable strips/masks | Slow translated/scaled layers and colour modulation; no fluid simulation |
| Lighthouse and foreground coast | Pre-cut painted transparent assets | Static transforms, glow/overlay at lamp |
| Harbour windows and five lanterns | Small lit/unlit sprites | Draw according to actual crossing results |
| Boat and lantern | Small illustrated frames | Translate/rotate/rock from core pose; lit frame from `Boat.Lit` |
| Six sightings | Compact authored atlas, several frames where visible | Enter/notice/chase/leave using state + age; no gameplay changes |
| Beam, water glints and fog | Small reusable soft-gradient textures | Transform/blend; one optional shader only if visually and measurably superior |
| Book, prediction, Harbour, Rest and Bell UI | Painted panels, icons, high-quality type | Shared typography and panel transforms; no UI framework |
| Decorative particles | Reused tiny textures | Renderer-only pool, <=256 with overflow dropped |

**Shipped formats:** JPEG for genuinely opaque artwork if it passes quality review; PNG for alpha/transparency; small original short WAV or compact Ogg for ambient audio. Source art stays at editable/high resolution in an offline asset source folder or referenced production archive; release includes only optimised exported assets, licences and the instructions necessary to regenerate them. Do not add an image processing SDK to the runtime.

**Animation:** water uses slow phase-offset movement, the boat rocks by a small sine function, lanterns breathe subtly, sightings use authored frame changes plus gentle movement, fog curls using a few layers, and UI transitions use simple cubic easing. No skeletal framework or tween library. Each motion parameter is drawn from state/elapsed visual time and is bounded.

### 20.3 Rendering pipeline and frame pacing

```text
fixed 60 TPS pure game Update
     |
     +--> authoritative boat / beam / sighting / check-in state
     |
     v
adapter reads current visual pose (read-only)
     |
     v
Draw at display VSync (60, 120 or measured capability)
     |
     +--> bounded sub-tick visual motion and looping decoration
     +--> draw ordered image layers / soft beam / fog / sprites / UI
     v
high-density framebuffer -> centred 4:3 viewport + bars -> display
```

Use Ebitengine's normal vsync; `Update` continues at 60 TPS and `Draw` may run at the monitor refresh rate. At >60 Hz, draw the *latest* authoritative beam, boat and input response immediately: do **not** add a one-tick render buffer that delays player feedback. If tests show 60 TPS poses appear to step at 120 Hz, use a bounded sub-tick visual projection forward from the latest boat velocity/beam angular velocity, with a time fraction in `[0,1]`. This projection only changes on-screen transforms; never alter the *lit decision*, rewards, check-ins, save, collision or event generation. On direction reversal, crossing changes and screen transitions, snap cleanly to the newest authoritative pose rather than extrapolating stale momentum. Decorative twinkles and waves may use visual elapsed time in `Draw` without changing gameplay RNG. `Draw` is side-effect-free with respect to the game.

The simplest no-buffer renderer is the baseline. Authorise sub-tick visual projection only if an actual 120 Hz video/frame-time and response comparison shows visibly smoother motion **without extra input latency**. Since Ebitengine already schedules `Update` separately from `Draw`, do not build custom timing threads or run 120 game ticks just to make the display look fluid.

### 20.4 Hi-DPI and 4:3 display mapping

`960 × 720` defines gameplay *coordinates*, not the final number of display pixels. Derive a render target that is sharp at the effective display scale (using Ebitengine's current monitor device-scale factor / layout mechanism). For a full screen/window `(W,H)` in actual target pixels:

```text
scale = min(W/960, H/720)
viewW = 960*scale; viewH = 720*scale
left = (W-viewW)/2; top = (H-viewH)/2
screenXY = (left,top) + logicalXY*scale
```

Letterbox with dark neutral bars; never distort or clip the active sea. The implementation may use a `Layout`-sized physical framebuffer and a common transform, or a final-screen presentation hook; choose **one** after validating Hi-DPI behaviour, not both. Test 4:3 and 16:9, 100%/150%/200% desktop scale, and movement between monitors. The single highest-density asset set must remain sharp when scaled; cap render resolution *only* by measured GPU limits and re-test legibility. Use Ebitengine's Hi-DPI example and FAQ as the implementation reference.

### 20.5 Batching, overdraw and asset lifetime

- Decode/create textures once while loading (or once on a screen transition for nonessential UI). Do not create textures, parse files, render to CPU pixels or generate full-screen noise every frame.
- Order opaque/background passes first, then water, lighthouse/coast, beam and targeted water illumination, boat/sightings, fog, particles, UI. Use a small number of reusable soft-gradient textures; avoid stacking multiple full-screen translucent passes.
- Ebitengine already maintains automatic texture atlases and batches compatible successive draw calls. Prefer its built-in handling; introduce manual atlas packing only if profiling proves a benefit. Group compatible blend/filter calls; avoid heavy CPU pixel reads, `ReplacePixels` each frame or screen-as-render-source feedback.
- Start with **zero custom Kage shaders**. Add at most one or two only if art quality/performance is measurably better than reused textures; cache shader compilation at load. No software post-process pipeline or additional GPU framework.
- Keep animations and optional particles finite and pooled (<=256). Missing optional spray/glints degrades gracefully; missing essential boat, lamp or sighting art fails visibly at loading, not halfway through a night.

### 20.6 Screens and UI finish

Title, Predict, Night, Harbour, Book, Rest, BellSetup and Bell all reuse the same palette, fonts, panel corners, transition durations and icon scale. Build the Book as one illustrated open-book frame with six content pages, *not* six separate large screen backgrounds. The Log is readable as a parchment page with both icons and numbers; no information only in colour. Bell overlay softly dims the screen without pausing visual compositing; gameplay itself remains paused. Rest feels restful, not like a warning. Provide no failure or diagnostic language in child-facing screens.

The visual reference must be judged **in a real running scene**: a static concept painting alone cannot prove interactive beams, fog, sighting visibility, animations, screens or 120 Hz fluidity.

### 20.7 Performance contract and controlled degradation

| Metric | Required evidence / provisional engineering target |
| --- | --- |
| Rule update | 60 TPS, seeded replay; never depends on render FPS |
| 60 Hz experience | Measured 60 FPS with normal vsync, correct input and no visible hitches on designated Linux reference hardware |
| 120 Hz experience | Measured 120 FPS on a capable 120 Hz display/GPU with boat, beam and decorations visibly smooth |
| Total present cadence | Target <=1% missed refresh opportunities during a five-minute representative test; report p50/p95/p99 frame intervals |
| Per-frame app work | Initial target p99 <=12 ms at 60 Hz and <=6 ms at 120 Hz (CPU-side Draw+Update measured separately from present intervals); tune using real profiles |
| Visual asset download size | Aspirational <=24 MiB *compressed*; not permission to lower quality |
| Peak graphics memory | Provisional 256 MiB on reference 1080p profile, report decoded RGBA + GPU/atlas overhead; higher resolution evaluated separately |
| Per-frame allocations | Profile to zero *avoidable* allocations in steady-state rendering; do not add cleverness solely for an absolute-zero metric |
| Optional effect overflow | Drop excess particles/mist decorations, never gameplay events or required feedback |

If hardware falls behind: first turn off only ornamental particles/glints, then lower optional water overlay count, then choose a smaller *render target* if still sharp and legible. Preserve the boat, beam, fog gameplay indication, sighting states, text and all sounds. Do not silently promise 120 FPS on a 60 Hz monitor; a real 120 Hz acceptance machine is mandatory. No assertions of measured performance are made before that run.

### 20.8 Visual production and acceptance gate

**Before implementing all eight screens, build one vertical-slice gameplay scene** containing the final lighthouse, boat, beam, water, fog, one animated sighting, a check-in overlay and a harbour window change. This is a release-quality slice, not a placeholder. Validate it at 60 Hz and 120 Hz, on a Hi-DPI screen and with a child following the boat. Only then reuse these mechanisms across the remaining screens.

Record: actual screenshot and short moving capture, frame-time distribution, FPS/TPS, input response observations, asset manifest and licence evidence, compressed pack size, decoded texture memory, and a screenshot comparison against the approved concept-quality bar. `TESTING.md` must show rejections and corrections, not only successful tests.

## 21. Testing and acceptance evidence

All tests run headless in `internal/game`, with seeded RNG and injected `Today`.

| PRD acceptance criterion | Evidence |
| --- | --- |
| 1 crossing completes, faster when lit | `TestCrossingAlwaysCompletes`, `TestLitCrossingFasterThanUnlit` |
| 2 beam rotates; lit only in cone | `TestBeamEasingAndClamp`, `TestLitOnlyInsideCone` |
| 3 fog gathers and clears | `TestFogRiseAndClearTimes` |
| 4 notice logs without moving beam | `TestNoticeLogsOldestVisibleAndBeamUnchanged` |
| 5 chasing loses the sighting | `TestChaseAfterSixtyTicksUnlogged` |
| 6 came back within window | `TestCameBackInsideWindow`, `TestNoChimeInGrace`, `TestOneChimePerDrift` |
| 7 prediction and good guess | `TestLanternThreshold`, `TestGoodGuessWithinOne` |
| 8 bells and any-answer star | `TestBellScheduleBounds`, `TestCheckInPausesSimulation`, `TestAnyAnswerEarnsStar` |
| 9 progression within caps | `TestProgressionMonotonicAndCapped` |
| 10 rest after two nights | `TestRestAfterTwoNightsSameDay`, `TestNewDayResets`, `TestParentOverride` |
| 11 Keeper's Bell | `TestDeskBellIntervals`, `TestDeskAnswersRecorded` |
| 12 persistence | `TestSaveRoundTrip`, `TestMissingSaveStartsEmpty`, `TestMalformedSaveReportedAndKept`, `TestUnknownFieldsRejected`, `TestLogCaps` |
| 13 no failure language | `TestNoFailureStrings` over all UI strings in one table |
| 14 headless determinism | `TestReplayDeterministic`, `TestPresentationEffectsDoNotAffectGameplay`, `TestEventsMatchStateChanges` |
| 15 full-screen 4:3 | play-test evidence in TESTING.md |
| 16 CI gates | `.github/workflows/test.yml` |
| 17 premium visual quality | visual sign-off in `TESTING.md`, all screens and the vertical-slice capture |
| 18 smooth 60/120 Hz | measured hardware frame-time report, `ActualFPS`/`ActualTPS`, transitions and input observations |
| 19 aspect / Hi-DPI / readability | `TestViewportPreservesAspect`, multi-scale screenshots and live monitor checks |
| 20 art licensing and graphics budgets | asset licence/source register, bytes/texture audit, `TestCosmeticStateDoesNotAffectGameplay` |
| 21 data resilience and daily override | `TestOversizedSaveRejectedBeforeDecode`, `TestInvalidSaveFieldsRejected`, `TestAtomicSaveInterruptionKeepsOldFile`, `TestOneExtraNightPerDay` |

Additional regression tests required for revised rules: `TestAllBoatPositionsReachableByBeam`, `TestHarbourReachableAtAngleClamp`, `TestChimeAfterGraceWhenAlreadyUnlit`, `TestBellScheduleAlwaysTwoOrThreeOnSlowCrossings`, `TestBellNeverDuringPause`, `TestDeskClockAdvancesAcrossMissedUpdates`, `TestDeskUnansweredNoFalseAnswer`, `TestNoRepeatedSoundOnStableLit`, `TestPresentationMotionNoRuleMutation`, `TestRendererShowsLatestInputState`. All UI strings live in one table in the game package so test 13 can scan them.

The first 16 AC remain in force. AC17–AC21 require human/hardware evidence as well as tests; CI alone cannot certify the artistic bar, native frame cadence, legibility or homework-session success. Do not tick these gates without actual evidence.

## 22. Resource budgets

```text
sightings visible          <= 3 (total active/exiting <= 6)
particles                  <= 256 (renderer-only, overflow dropped)
bells per night            2..3, one per selected crossing slot
check-ins per night        <= 3
crossings                  exactly 5
night log entries         <= 400; desk entries <= 400
save file output          target < 200 KiB; bounded read <= 256 KiB
events per gameplay tick   <= 8; reuse backing memory
art pack (compressed)      aspirational <= 24 MiB, subject to visual gate
GPU/texture memory         provisional <= 256 MiB at 1080p; measure
steady 60 Hz present       <= 16.67 ms budget; 120 Hz <= 8.33 ms
optional effects           bounded, first to degrade
```

## 23. Repository shape

```text
cmd/
    lighthousekeeper/
        main.go        adapter: input, injected clocks, save I/O, game loop
        render.go      one painterly layer renderer, viewport and interpolation
        audio.go       event-to-sound playback; omit file if small enough for main
        assets/        optimised sky/water/lighthouse/boat/sighting/UI images and audio
internal/
    game/
        game.go        state, screens, update loop, beam, boat, fog, bell rules
        sightings.go   spawn, notice, chase
        save.go        Save type, validation, atomic persistence
        strings.go     every child-facing UI string
        *_test.go
TESTING.md             deterministic, visual, performance, child test evidence
RELEASE.md             honest unticked acceptance gates until verified
assets-source/         optional offline editable source art and export recipe; no runtime dependency
```

`render.go` earns its separate file because the art/render concern is materially larger than the rule engine; it is **not** a new scene graph or rendering package. `sightings.go` earns its file because spawn/aim/book are one concern; `strings.go` makes the no-failure-language test possible. Use Ebitengine's automatic sprite atlas unless measured otherwise. No new microcomponent/repository is needed.

## 24. Failure behaviour

| Failure | Behaviour |
| --- | --- |
| save absent | start empty |
| save malformed, over 256 KiB, invalid version/fields | report gently, preserve recovery copy, start empty; avoid unbounded allocations |
| save write fails | report on Harbour/appropriate transition screen, preserve previous save, keep playing |
| sound missing | continue silently |
| essential sprite missing | fail loading visibly with a repairable error; do not play with invisible boat/beam |
| optional particle/glint missing | suppress cosmetic effect; preserve core visuals |
| display unable to sustain effect load | drop optional decoration first; report measured profile; preserve 60 TPS logic |
| 120 Hz display unavailable | use measured native refresh; do not claim 120 Hz acceptance |
| date unavailable | adapter supplies empty date; daily rest rule disabled for this session only |
| long desk-mode suspension | resume from injected elapsed real time, no burst of overdue bells; end if time expired |
| programming error | fail visibly (diagnostics separate from child-facing game strings) |

## 25. Privacy

No network, telemetry, accounts or runtime AI. The save file holds only the book, nightly counts and yes/no answers on the local machine; no clinical measure, diagnosis, user identity, profile or clinical claim. Validate local JSON **before** allocating or accepting state, restrict save-file permissions where supported, and reject unknown fields and unexpected versions. Only user-triggered local config data is read; no plugins, remote content, executable mods or untrusted asset downloads. Original images, fonts and audio must have redistributable licences recorded in the asset manifest. None of these controls is relaxed to obtain visual speed or a smaller binary.

## 26. CI and release

Copy Astro Answers' workflow: module tidy contract, `gofmt`, `go vet`, `go test -race -count=1`, Linux/amd64 and Linux/ARM64 builds, named acceptance tests from section 21, AppImage. Pin Go/Ebitengine versions, record module provenance and original asset rights. Make the staged vertical-slice screenshot/capture and actual 60/120 Hz benchmark required non-CI evidence. `RELEASE.md` leaves high-refresh proof, constrained hardware, independent child night and real homework session **unticked** until completed. Play-test fixes are committed with a reason and recorded in `TESTING.md`.

**No invented proof:** this EDD is the designed solution. The geometry calculation and bounded schedule can be checked analytically, but named Go test results, final binary size, art quality and frame-time measurements do not exist until built and tested.

## 27. Deliberately absent

```text
scores, leaderboards       clinical metrics          parent dashboard
scene graph / ECS           UI framework               state-machine library
physics / fluid simulation particle engine            notification system
network / telemetry / AI   gamepad / touch            music with a beat
custom frame scheduler     animation/tween framework   multiple art skin engines
```

## 28. Open questions for play-testing

1. Is N the right notice key, or is Space easier to reach without looking?
2. Is 120 ticks unlit the right moment for the chime, or does it feel like nagging?
3. Do six sighting kinds give enough to collect, or does the book need seasonal rarities?
4. Is 90 seconds per crossing calm or slow?
5. Does the 4:3 recomposition keep the boat easy to find within the detailed, approved painterly direction?
6. Does any optional sub-tick visual motion improve 120 Hz fluidity without adding input latency or overshoot?
7. What is the smallest visually indistinguishable compressed/decoded art set on the reference Linux hardware?

## 29. Design conclusion

The beam is the attention, the boat is the task, the sightings are the world being interesting, and the chime is the moment of choice. Everything the game rewards (noticing, coming back, predicting honestly, answering the bell) is a behaviour she can take to her desk.

**Engineering conclusion:** a small rules engine can support a lavish-looking game when most artistic complexity is authored once, outside runtime; the runtime only moves, blends and displays a few carefully made layers. The success condition is *premium visual quality and measured 60/120 Hz fluidity with the entire PRD intact*, not an impressive number of abstractions removed.

> **Does every piece of machinery remove more complexity than it introduces?**

## 30. Historical renderer simplification — pending visual acceptance

The final scene is one coherent authored 1920 × 1440 JPEG from an original painterly source (not three stitched ocean paintings). A short offline script removes the source's decorative boat using same-horizon water; dynamic boat, spotlight, fog, sightings, lantern and HUD remain distinct. The old `shimmer.png` texture and sixteen glint drawings have been deleted. Motion instead uses **twelve prepared image subregions** drawn with bounded horizontal sine offsets and very low alpha over their original location. This directly reuses Ebitengine's existing `SubImage` and `DrawImage`; no shaders, 3D engine, physics or runtime asset pipeline. The painted scene is 479,359 bytes compared with the previous 695,231-byte stitched image (~31% reduction for that asset). The gameplay core is unchanged.

**Acceptance:** full CI, real executable screenshots and continuous gameplay inspection, 60 Hz / 120 Hz hardware frame pacing, child play and homework UAT remain separate gates. A successful artist-side preview is not release proof.
