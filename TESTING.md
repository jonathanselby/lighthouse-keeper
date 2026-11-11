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
