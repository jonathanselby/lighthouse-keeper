# Lighthouse Keeper — V0.1.0 release checklist

## Contract

- [x] PRD v0.2 retained in `PRD.md` with all 16 original and five added criteria.
- [x] EDD v0.2 retained in `EDD.md` as the current design baseline.
- [x] Linux desktop, full-screen 4:3, offline, keyboard-first scope.
- [x] No network, telemetry, accounts, advertising or runtime AI.
- [x] Code-derived EDD and comparison published against source commit `193a10ee255e43d44c6ad50857f2fd8750779784`; physical product acceptance remains open.

## Engineering gate

- [x] Headless game-core tests passed locally under Go 1.23.2 (module-off fallback; *not* full CI).
- [x] Go 1.25 / pinned dependency module contract green in GitHub CI (run 42962933707).
- [x] `gofmt` clean in CI.
- [x] `go vet ./...` green in CI.
- [x] `go test -race -count=1 ./...` green in CI.
- [x] Linux/amd64 and Linux/arm64 builds green in CI.
- [x] Embedded art and procedural sound decode tests green in CI.
- [x] Linux GUI launch smoke test green under Xvfb (release candidate run 42962933696; not full UAT).

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

- [x] Current source/tests/assets snapshot recorded at commit `193a10ee255e43d44c6ad50857f2fd8750779784`; subsequent changes to this implementation require regenerated evidence.
- [x] `generated/EDD.md` written from direct source inspection and pinned to that commit.
- [x] `generated/COMPARISON.md` records implementation matches, changes and open physical evidence.
- [ ] `TESTING.md` includes UAT and hardening failures, not only successes.

## Publication

- [x] GitHub source present and all nine image blob hashes match local originals.
- [x] CI and release-build workflow green (runs 42962933707 and 42962933696).
- [x] Linux executable built, smoke-tested and packaged as GitHub Actions artifact; not yet attached to a public Release.
- [ ] Release tag `v0.1.0` created only after code/evidence gates.
- [ ] Release binary attached to GitHub Release, not committed to source history.
- [ ] GitHub About and topics completed.

## Rule

Unobserved hardware, visual and child tests remain unchecked. The repository may contain implementation candidates before release, but no published v0.1.0 release is called complete until accepted evidence exists.
