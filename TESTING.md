# Lighthouse Keeper — Testing and Debugging Record

**Status:** implementation candidate. Do not interpret a design target as a measured result.

## Tests actually run locally

Environment: Go 1.23.2, Linux container, without external Go module downloads.

```text
GO111MODULE=off go test ./internal/game -count=1
PASS
```

Headless tests run without a renderer, native window or hardware audio. They exercise complete crossing and deterministic replay, speed difference under illumination, beam geometry, fog, sighting notice/chase, chime/window rearm, lantern thresholds, scheduled check-ins, rest and extra-night rules, desk timing and persistence safety. These tests do not prove that the Ebitengine application compiles on Go 1.25, nor that a Linux desktop can actually launch it.

## Integration and hardening gates — outstanding until verified

- [ ] Go 1.25+/Ebitengine 2.10.6 module resolution and `go mod tidy` contract.
- [ ] Full `go vet ./...`.
- [ ] Full `go test -race -count=1 ./...`.
- [ ] Linux/amd64 and Linux/arm64 production builds.
- [ ] Embedded asset decode and procedural audio tests under the build toolchain.
- [ ] Production executable Xvfb graphical-launch smoke test.
- [ ] Full-screen gameplay, visual transitions, Book, Harbour Log, Rest, Desk mode and save recovery manually observed.
- [ ] Actual 60 Hz hardware frame-time evidence and input responsiveness.
- [ ] Actual 120 Hz hardware frame-time evidence on a verified 120 Hz display.
- [ ] Hi-DPI and multiple aspect ratio checks.
- [ ] A child plays a night without help.
- [ ] One real 15/25-minute homework session with Keeper's Bell.
- [ ] Constrained-hardware Garage benchmark (CPU, RAM, GPU memory, frame rate, executable size, scene entity caps).
- [ ] Original layered artwork approved against five cinematic quality references.

## Test evidence required at release

Name the hardware, display refresh/resolution, OS, actual measured FPS/TPS and frame times rather than assuming them. Capture all screens and a representative night. When an issue is found, record the initial failure, smallest corrective change and retest result here. Do not check a release gate merely because a script exists.

## Known engineering risks to confirm

- Frame projection at 120 Hz must neither display obvious overshoot nor introduce a buffered-input frame.
- Soft beam alpha must match its authoritative angular hit test visually.
- Saving at a transition must not cause a noticeable hitch during animated scenes.
- Local audio access may vary across Linux desktops; silent mode must remain playable.
- Large 1920×1440 render target must sustain chosen refresh rates on real GPU hardware.
- Asset visual quality must meet human premium-art approval; no automated test can certify taste.

## Corrective painterly build — visual implementation review

**Baseline evidence:** prior implementation was CI-green at GitHub Actions run `42962933707` and the production release candidate launched under a virtual Linux display in run `42962933696`. Its actual running night screenshot failed the owner's premium-art standard, despite passing the technical build gates. That is the reason for this focused renderer-only replacement.

**Local work completed for the corrective candidate:**

- Kept the entire `internal/game` implementation unchanged; the headless Go game package still passes.
- Replaced generated water arcs with a composited original painterly sea and bounded moving water-light texture.
- Moved the visible boat and light cone into the same monotonic painterly scene projection without affecting logical hit tests.
- Replaced the primitive boat, lighthouse and sightings assets with masked original illustrations.
- Replaced plain geometric panels with one nine-slice brass/navy panel, pre-sliced at startup.
- Added separate lanterns, glow, fog, and art atlas dimension/decoding tests.
- Generated an offline composite review still: `lighthouse-premium-night-preview.jpg` (this is **not** a captured game frame).

**Gates still open on the changed source:** Go 1.25/Ebitengine CI, graphical launch of the revised binary, actual live motion and high-DPI screenshots, native 60/120 Hz frame-time and responsiveness measurements, child UAT, homework UAT and final owner approval of artistic quality. The previous successful CI and release build must not be represented as covering new unmerged source.

## v0.3 UAT remediation gates — 13 November 2026

**Local checks completed:** Go `GO111MODULE=off go test ./internal/game -count=1` passed with new UAT tests for bounded progressive crossing times and reduced beam width, alongside all existing deterministic rule, bell, fog, persistence and no-failure tests. The self-contained offline art bake ran successfully. `assets-source/test_living_ocean.py` passed: 4,470 fully opaque hero-rock pixels unchanged across sampled water frames; 30,191 sea pixels changed in a measured open-water ROI; the ten atlas poses differ; incremental baked graphics measured at 727,014 bytes. These measurements are offline source-art checks, **not** actual GPU FPS.

**Open release evidence:** GitHub CI against Go 1.25.0/Ebitengine 2.10.6; native Kage shader compilation and Linux executable startup; at least one complete five-crossing game with two or three bells; inspection of all screens; beam/boat visual agreement at path extremes; 10–15 second actual smooth ocean capture with rocks and whitewater; 60 and 120 Hz real hardware frame pacing and input-response tests; child and homework UAT. Do not mark final graphics quality as accepted based on still previews or offline simulations.

## v0.3 CI and native smoke evidence — 13 November 2026

- Source commit `6f4dcc32e08990201909a735c04472c374a3c323`: [CI run 43327777321](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43327777321) **passed**, including Python art bake/immutability checks, gofmt, vet, race/unit tests, Linux amd64 and arm64 builds and named core acceptance.
- [Art bake run 43327777298](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43327777298) **passed** and committed regenerated binary artwork to main.
- [Native release candidate run 43327777251](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43327777251) **passed**, including Kage shader startup under Xvfb and published a Linux x86-64 artifact. Xvfb is a virtual-display startup check, not sustained 60/120 Hz hardware evidence.
- This does **not** close visual acceptance: still need full moving-scene capture from native executable, beam/boat tracking correspondence, art review across all rocks and foam, hardware FPS, child UAT and full session acceptance. Keep the build labelled candidate.

## v0.4 UAT graphics remediation — 13 November 2026

The prior build's shader passed absolute Unix epoch seconds (about 1.8 billion) into a float32 uniform, which quantises time in roughly 128-second increments and freezes waves. The GPU clock now uses relative 60 TPS frame count plus fractional tick for high-refresh presentation; spray and water share the clock. Deterministic simulation is untouched.

The moving boat previously included a static painted wake, and the six live sighting sprites were circular illustrated scenes with fixed land/water/spray. Offline reproducible asset authoring now extracts the boat and six transparent subjects. The Sightings Book retains its original illustrated atlas. The water/spray/wake remain separate animations.

Spawn selection is now kind-first, with explicit approved-art habitats: whale/jellyfish/ferry on open water; seal on nearshore rock; owl on cliff; shooting star in sky. Aquatic sightings rise across 48 ticks, with waterline cropping and separately animated spray; departing sightings submerge. Notice, chase, progress, saved book/log, crossing speeds and beam difficulty remain unchanged.

New acceptance gates: 600-seed habitat regression; 48-tick emergence regression; sub-frame shader clock precision; offline transparency and water-domain assertions; unchanged layered-rock checks. Real hardware 60/120 Hz pacing, whole-night visual UAT, child UAT and final artistic sign-off remain open.

## v0.4 bake and formatting verification

Offline cut-outs and Go formatting were committed by the reproducible GitHub art-bake workflow. The final commit must be independently validated by CI and the native Linux release workflow. Do not substitute previous failed runs for current acceptance.

## v0.4 verified CI and native build — 13 November 2026

Commit `d177ef99307c4c057b26bf48261dbe9aad4718a8` passed [CI 43339549074](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43339549074), including Go race/unit tests, new water-clock precision and habitat tests, offline moving-subject art checks, fixed rock regression, Linux amd64 and arm64 compilation. The [native Linux release workflow 43339549011](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43339549011) passed including graphical startup and published artifact `Lighthouse-Keeper-v0.4.0-uat-linux-x86_64`. The separate art bake [43339438643](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43339438643) passed and committed reproducible binary cut-outs.

These technical checks do not replace human visual UAT. Actual GPU ocean motion over a sustained crossing, live animal surfacing/occlusion, all sighting habitats, 60/120 Hz frame pacing and child acceptance remain open.
