# Lighthouse Keeper — Testing and Debugging Record

**Status:** implementation candidate. Do not interpret a design target as a measured result.

## Tests actually run locally

Environment: Go 1.23.2, Linux container, without external Go module downloads.

```text
GO111MODULE=off go test ./internal/game -count=1
PASS
```

Headless tests run without a renderer, native window or hardware audio. They exercise complete crossing and deterministic replay, speed difference under illumination, beam geometry, fog, sighting notice/chase, chime/window rearm, lantern thresholds, scheduled check-ins, rest and extra-night rules, desk timing and persistence safety. These tests do not prove that the Ebitengine application compiles on Go 1.25, nor that a Linux desktop can actually launch it.

## GitHub engineering and launch gates observed — 11 November 2026

- **CI success:** [run 42962738531](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/42962738531), source commit `46e44f0a884e74305b5ad578c4ff2398e4e7beb6`: module contract, gofmt, vet, race-enabled tests, Linux/amd64 and Linux/ARM64 builds, and named core acceptance all passed.
- **CI reconfirmation:** [run 42962933707](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/42962933707), source/workflow commit `193a10ee255e43d44c6ad50857f2fd8750779784`, passed.
- **Production release candidate:** [run 42962933696](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/42962933696), source/workflow commit `193a10ee255e43d44c6ad50857f2fd8750779784`, passed full validation, native amd64 build, Xvfb GUI smoke launch, permission-preserving tar archive verification and Actions artifact upload.
- **Binary asset integrity:** all eight embedded art assets and the approved source PNG in GitHub were compared byte-for-byte using Git blob SHA against the local approved files.
- **Initial failed CI:** first module run failed because go.mod indentation was noncanonical; next run passed module/format but failed before compilation because the binary assets were missing. After asset upload the first complete graphical compile found a `float64 + int` error in the water line expression; corrected to `float64(i)`. No gameplay mechanic changed.
- **Actual GUI gameplay and UAT:** NOT observed by this automated gate. Six seconds of successful Xvfb startup is not a completed gameplay or visual test.

## Integration and hardening gates — checklist with evidence

- [x] Go 1.25+/Ebitengine 2.10.6 module resolution and `go mod tidy` contract (green GitHub CI).
- [x] Full `go vet ./...` (green GitHub CI).
- [x] Full `go test -race -count=1 ./...` (green GitHub CI).
- [x] Linux/amd64 and Linux/arm64 CI production builds.
- [x] Embedded asset decode and procedural audio tests under the Go 1.25 CI toolchain.
- [x] Production executable Xvfb graphical-launch smoke test (run 42962933696).
- [ ] Full-screen gameplay, visual transitions, Book, Harbour Log, Rest, Desk mode and save recovery manually observed.
- [ ] Actual 60 Hz hardware frame-time evidence and input responsiveness.
- [ ] Actual 120 Hz hardware frame-time evidence on a verified 120 Hz display.
- [ ] Hi-DPI and multiple aspect ratio checks.
- [ ] A child plays a night without help.
- [ ] One real 15/25-minute homework session with Keeper's Bell.
- [ ] Constrained-hardware Garage benchmark (CPU, RAM, GPU memory, frame rate, executable size, scene entity caps).
- [ ] Original layered artwork approved against five cinematic quality references.

## Graphical visual inspection — Linux Xvfb, 2026-11-11

The production Linux x86-64 executable from release-candidate workflow 42962933696 was unpacked, executed in a 1920×1080 virtual display using software OpenGL, and photographed at Title and at two points during Crossing 1. It stayed running and showed a functioning screen transition from Title → Predict → Night via keyboard input.

**Observed success:** the title painting displays at correct 4:3 aspect ratio with black pillarboxing; keyboard input advances screens; the night scene draws the beam, sea, lamp, crossing HUD, moving boat and sightings.

**Observed visual-quality failure:** the current night scene is substantially flatter and more schematic than the five approved premium painterly references, especially the water, illustrated sightings, lighthouse and boat sprites. The title menu is functional but uses generic rectangular panels. These visuals are implementation placeholders and do **not** satisfy the locked premium graphics acceptance criterion.

**Additional usability risk:** the first part of a crossing may position the boat near the HUD/sky region, where readability is compromised. The photographic snapshots are not sufficient to certify all seeded paths or whether noticing and return prompts remain readable at every boat position.

**Next gate:** replace or materially improve the active scene's hero art and layout while preserving the core engine, then repeat captured-screen inspection and hardware frame pacing measurements. No high-refresh performance claim is made from these software-rendered screenshots.

## Test evidence required at release

Name the hardware, display refresh/resolution, OS, actual measured FPS/TPS and frame times rather than assuming them. Capture all screens and a representative night. When an issue is found, record the initial failure, smallest corrective change and retest result here. Do not check a release gate merely because a script exists.

## Known engineering risks to confirm

- Frame projection at 120 Hz must neither display obvious overshoot nor introduce a buffered-input frame.
- Soft beam alpha must match its authoritative angular hit test visually.
- Saving at a transition must not cause a noticeable hitch during animated scenes.
- Local audio access may vary across Linux desktops; silent mode must remain playable.
- Large 1920×1440 render target must sustain chosen refresh rates on real GPU hardware.
- Asset visual quality must meet human premium-art approval; no automated test can certify taste.
