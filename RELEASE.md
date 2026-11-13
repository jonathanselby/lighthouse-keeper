# Lighthouse Keeper — V0.1.0 release checklist

## Contract

- [x] PRD v0.3 retained in `PRD.md` with all 16 original and five added criteria.
- [x] EDD v0.3 retained in `EDD.md` as the current design baseline.
- [x] Linux desktop, full-screen 4:3, offline, keyboard-first scope.
- [x] No network, telemetry, accounts, advertising or runtime AI.
- [ ] Source and documents reconciled to the final accepted gameplay implementation.

## Engineering gate

- [x] Headless game-core tests passed locally under Go 1.23.2 (module-off fallback; *not* full CI).
- [ ] Go 1.25 / pinned dependency module contract green in GitHub CI.
- [ ] `gofmt` clean in CI.
- [ ] `go vet ./...` green in CI.
- [ ] `go test -race -count=1 ./...` green in CI.
- [ ] Linux/amd64 and Linux/arm64 builds green in CI.
- [ ] Embedded art and procedural sound decode tests green in CI.
- [ ] Actual Linux GUI launch smoke test green.

## Product acceptance

- [ ] Complete interactive night meets all gameplay requirements under human UAT.
- [ ] All eight screens and transitions approved.
- [ ] Sightings Book and Harbour Log persist across relaunch.
- [ ] Painterly art meets the five approved quality benchmarks.
- [ ] Native 60 Hz frametime, fluidity and responsiveness proven.
- [ ] Native 120 Hz frametime, fluidity and responsiveness proven.
- [ ] Hi-DPI / aspect-preserving letterboxing verified.
- [ ] Save recovery, atomicity, and one-extra-night limit verified in the application.
- [ ] Child completes an entire night independently.
- [ ] Keeper's Bell used during one real homework session.
- [ ] Garage constrained-hardware benchmark recorded.

## Independent evidence

- [ ] Source, tests and assets final; release source-changing commit recorded.
- [ ] `generated/EDD.md` independently recovered from that commit.
- [ ] `generated/COMPARISON.md` reconciles design and observed source.
- [ ] `TESTING.md` includes UAT and hardening failures, not only successes.

## Publication

- [ ] GitHub repository source and binary asset blobs verified.
- [ ] CI and release-build workflow green.
- [ ] Linux executable built, smoke-tested and packaged.
- [ ] Release tag `v0.1.0` created only after code/evidence gates.
- [ ] Release binary attached to GitHub Release, not committed to source history.
- [ ] GitHub About and topics completed.

## Rule

Unobserved hardware, visual and child tests remain unchecked. The repository may contain implementation candidates before release, but no published v0.1.0 release is called complete until accepted evidence exists.

## Corrective graphics revision

- [x] Previous graphical launch inspected and premium-art failure recorded.
- [x] Painterly assets and new renderer prepared with unchanged rules.
- [x] Authoring/source asset provenance and current exporter documented.
- [ ] Changed source and binary artwork uploaded and byte-verified in GitHub.
- [ ] Changed-source CI and production graphical launch gate green.
- [ ] Running-scene and transition screenshots visually reviewed.
- [ ] Native 60/120 Hz fluidity and input responsiveness measured.
- [ ] Independent code-derived EDD/comparison regenerated from the **new** final source commit.
- [ ] Owner/child UAT and release sign-off.

## v0.3 UAT fix status

- [x] All prior PRD features retained in the revised v0.3 PRD/EDD.
- [x] Independent fixed-rock and moving-water art pipeline; no colour-based runtime rock classification.
- [x] Prominent baked spotlight retouched in the *gameplay* image; separately moving spotlight remains.
- [x] Reduced and progressive crossing durations with strengthened beam-tracking challenge; headless tests pass locally.
- [x] Changed bell schedule retains 2–3 check-ins during short crossings; headless tests pass locally.
- [x] Independent spray sprites cut from original project art, fixed shore and boulders; offline motion/immutability tests pass.
- [ ] CI full toolchain and graphical shader compilation with newly baked binaries green.
- [ ] Linux gameplay and full-night interactive UAT, including screenshots/motion capture.
- [ ] GPU 60/120 Hz hardware evidence and artistic approval.
- [ ] Independently generated engineering design reconciled to final source commit.
- [ ] Human child and Keeper's Bell homework acceptance.

**The public release tag remains withheld until outstanding gates are evidenced.**

## v0.3 integrated UAT candidate — 13 November 2026

- [x] Implement layered moving sea, fixed foreground rocks, independent spray, no baked static beam, faster progressive crossings and retuned beam tolerance.
- [x] CI with race tests and cross-compilation: [43327777321](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43327777321).
- [x] Reproducible artwork bake and committed binaries: [43327777298](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43327777298).
- [x] Native Linux executable and Xvfb launch: [43327777251](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/43327777251).
- [ ] Full-scene real-time visual inspection of all rock/spray areas and dynamic beam alignment.
- [ ] 60/120 Hz native hardware pacing, gameplay feel, child and homework UAT.
- [ ] Final release approval. **This is a candidate, not a signed-off release.**

## v0.4 UAT graphics remediation — 13 November 2026

The prior build's shader passed absolute Unix epoch seconds (about 1.8 billion) into a float32 uniform, which quantises time in roughly 128-second increments and freezes waves. The GPU clock now uses relative 60 TPS frame count plus fractional tick for high-refresh presentation; spray and water share the clock. Deterministic simulation is untouched.

The moving boat previously included a static painted wake, and the six live sighting sprites were circular illustrated scenes with fixed land/water/spray. Offline reproducible asset authoring now extracts the boat and six transparent subjects. The Sightings Book retains its original illustrated atlas. The water/spray/wake remain separate animations.

Spawn selection is now kind-first, with explicit approved-art habitats: whale/jellyfish/ferry on open water; seal on nearshore rock; owl on cliff; shooting star in sky. Aquatic sightings rise across 48 ticks, with waterline cropping and separately animated spray; departing sightings submerge. Notice, chase, progress, saved book/log, crossing speeds and beam difficulty remain unchanged.

New acceptance gates: 600-seed habitat regression; 48-tick emergence regression; sub-frame shader clock precision; offline transparency and water-domain assertions; unchanged layered-rock checks. Real hardware 60/120 Hz pacing, whole-night visual UAT, child UAT and final artistic sign-off remain open.
