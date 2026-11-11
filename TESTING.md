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

- [x] Go 1.25+/Ebitengine 2.10.6 module resolution and `go mod tidy` contract (CI `42974831234`).
- [x] Full `go vet ./...`.
- [x] Full `go test -race -count=1 ./...`.
- [x] Linux/amd64 and Linux/arm64 production builds.
- [x] Embedded asset decode and procedural audio tests under the build toolchain.
- [x] Production executable Xvfb graphical-launch smoke test (`42974831280`).
- [ ] Full-screen title and one active-night frame observed; remaining transitions, Book, Harbour Log, Rest, Desk mode and save recovery not yet manually confirmed.
- [ ] Actual 60 Hz hardware frame-time evidence and input responsiveness.
- [ ] Actual 120 Hz hardware frame-time evidence on a verified 120 Hz display.
- [ ] Hi-DPI and multiple aspect ratio checks.
- [ ] A child plays a night without help.
- [ ] One real 15/25-minute homework session with Keeper's Bell.
- [ ] Constrained-hardware Garage benchmark (CPU, RAM, GPU memory, frame rate, executable size, scene entity caps).
- [ ] Original layered artwork approved against five cinematic quality references — **not yet**: actual gameplay frame shows visual seam failures.

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

## Premium runtime inspection — 11 November 2026

**Provenance:** source commit `f5029e2999223a2410848931fbaf8ef7f0c8c342`; successful Go CI `42974831234`, successful packaging/Xvfb startup `42974831280`. The GitHub Actions release artifact was downloaded and its actual native executable launched under a software-rendered virtual Linux display (1920×1440).

**Observed:** the title uses the approved painterly source artwork and legible panel treatment. The first active crossing displays the painted water, harbour, boat, illuminated lighthouse and newly framed HUD, confirming asset decode and rendering. The active scene has obvious horizontal edges from multiple incompatible water-image crops (near the sky-water join and between the successive ocean panels). The result is not a continuous painted environment. The positioned lighthouse and boat still need composition review. **Artistic acceptance fails.**

**Corrective direction:** replace the composite `seascape.jpg` with one coherently composed full-scene painting, rather than stacking unrelated ocean crops; retain independently animated, transparent boat, beam, fog and sightings, and preserve deterministic engine geometry. Validate against a newly captured *running* screenshot and live motion before any release approval.

**Not measured:** a virtual display with software rasterisation is not a 60/120 Hz hardware validation. No real GPU frame-time, input response or mobile/desktop high-refresh numbers are claimed. No child or homework UAT has occurred.
