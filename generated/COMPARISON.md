# Lighthouse Keeper — Design / Implementation Comparison

**Evidence source commit:** `193a10ee255e43d44c6ad50857f2fd8750779784`  
**Design baseline:** PRD v0.2 and EDD v0.2  
**Generated:** 2026-11-11

| Area | Intended | Recovered implementation | Result |
| --- | --- | --- | --- |
| PRD scope | All required game modes | Five crossings, predictions, six sightings, fog, return, bells, Book, Log, Rest and desk mode | MATCH in code; human UAT open |
| Architecture | One plain-Go deterministic engine | `internal/game.Game`, explicit Update/Events | MATCH |
| Time | Fixed 60 TPS, injected date | Ebitengine 60 TPS; date and real desk elapsed injected | MATCH |
| Crossing | Illuminated boat faster, always arrives | Cubic Bézier and arc-length advance; 35% fog drift speed | MATCH, tested |
| Beam | Reach all boat positions | Cone and distance checks, ±1.50 rad bounds | MATCH, tested |
| Fog | Rise, clear and slow boat | Deterministic fog accumulation and speed scaling | MATCH, tested |
| Sightings | Six kinds, notice vs chase | Bounded spawn, oldest-visible notice, 60-tick chase | MATCH, tested |
| Return | Grace, chime, five-second return | Tick counters and rearm logic | MATCH, tested |
| Prediction | Five lanterns, 70% threshold | Computed result and ±1 good guess | MATCH, tested |
| Check-ins | Two or three pauses, any answer | Crossing schedule with Y/N record | MATCH, tested |
| Rest | Two nights/day, parent override | Injected day, three-second P hold, persisted override | MATCH, tested |
| Keeper's Bell | 15/25 minutes, 3–7-minute rings | Monotonic elapsed ticks and desk records | MATCH in code; homework UAT open |
| Persistence | Local, bounded, recoverable | Strict JSON validation, `.bad` backup, atomic replace | MATCH, tested |
| Privacy | No accounts/network/runtime AI | No corresponding implementation path | MATCH, inspected |
| Graphics | Premium painterly and high-DPI | Original illustrated assets, embedded images, 2× raster composition | IMPLEMENTED; aesthetic UAT open |
| Fluidity | 60 FPS, 120 FPS where available | Render-only sub-tick presentation | UNMEASURED |
| Audio | Soft original effects | Procedural PCM cues and ambient loop | IMPLEMENTED; listening UAT open |
| Build CI | tidy, gofmt, vet, race, amd64, ARM64 | Green run 42962738531 | VERIFIED |
| Production launch | Graphical Linux startup | Green Xvfb smoke run 42962933696 | VERIFIED startup only |
| Release evidence | Independent recovered EDD | This comparison and `generated/EDD.md` | GENERATED; publish verification pending |
| Product UAT | Child night, real homework session | No physical testing evidence | OPEN |
| Constrained hardware | Measured FPS, CPU, RAM, GPU | No physical benchmark | OPEN |

## Material corrections during implementation

1. Binary art was absent from the first text-only GitHub upload. The eight runtime images and original source concept were browser-uploaded and matched exactly by Git blob SHA.
2. `go.mod` originally failed the tidy contract only because of indentation. Formatting was corrected with no dependency change.
3. Missing binary assets initially prevented `//go:embed` compilation; asset publication resolved the error.
4. `render.go` contained `math.Sin(seconds+i)`, adding float64 to int. Corrected to `math.Sin(seconds+float64(i))`, without changing gameplay.
5. Removed the temporary asset-directory `.keep` after upload verification.
6. Enabled release-candidate workflow triggering on its own workflow-file changes to run existing build, graphical-startup and packaging gates.

## Deliberate departures from other Garage games

The Lighthouse Keeper adapter uses more raster painting than Astro Answers and Dandelion Drift, because premium cinematic seascape illustration is an explicit product constraint. It still implements no second game engine, scene graph or general animation framework. Its single local save also holds richer, non-clinical Book/Log/desk records, rather than only high scores.

## Open evidence

The attached five concept artworks are visual benchmarks, not proof of gameplay finish. Neither successful compilation nor six seconds of Xvfb launch demonstrates stable 60/120 FPS, readability in motion, child comprehension, audio subtlety or suitability for real homework. These remain unchecked in `RELEASE.md`.

## Result

**ENGINEERING IMPLEMENTATION MATCH AND VERIFIED LINUX STARTUP; PRODUCT UAT OPEN.** No further architecture is justified by observed CI failures. Visual and physical acceptance must be judged separately.
