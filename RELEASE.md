# Lighthouse Keeper — V0.1.0 release checklist

## Contract

- [x] PRD v0.2 retained in `PRD.md` with all 16 original and five added criteria.
- [x] EDD v0.2 retained in `EDD.md` as the current design baseline.
- [x] Linux desktop, full-screen 4:3, offline, keyboard-first scope.
- [x] No network, telemetry, accounts, advertising or runtime AI.
- [ ] Source and documents reconciled to the final accepted gameplay implementation.

## Engineering gate

- [x] Headless game-core tests passed locally under Go 1.23.2 (module-off fallback; *not* full CI).
- [x] Go 1.25 / pinned dependency module contract green in GitHub CI (run `42974831234`).
- [x] `gofmt` clean in CI.
- [x] `go vet ./...` green in CI.
- [x] `go test -race -count=1 ./...` green in CI.
- [x] Linux/amd64 and Linux/arm64 builds green in CI.
- [x] Embedded art and procedural sound decode tests green in CI.
- [x] Actual Linux GUI launch smoke test green (run `42974831280`).

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

- [ ] Source, tests and assets final; visual rework outstanding. Candidate evidence commit `f5029e2999223a2410848931fbaf8ef7f0c8c342` recorded.
- [x] Candidate `generated/EDD.md` reconciled against `f5029e2999223a2410848931fbaf8ef7f0c8c342`; regenerate after final graphics changes.
- [x] Candidate `generated/COMPARISON.md` records material graphics deviation; final sign-off pending.
- [ ] `TESTING.md` includes UAT and hardening failures, not only successes.

## Publication

- [x] GitHub repository source and binary asset blobs verified.
- [x] CI and release-build workflow green at source commit `f5029e2999223a2410848931fbaf8ef7f0c8c342`.
- [x] Linux executable built, smoke-tested and packaged as an **unapproved release candidate**.
- [ ] Release tag `v0.1.0` created only after code/evidence gates.
- [ ] Release binary attached to GitHub Release, not committed to source history.
- [ ] GitHub About and topics completed.

## Rule

Unobserved hardware, visual and child tests remain unchecked. The repository may contain implementation candidates before release, but no published v0.1.0 release is called complete until accepted evidence exists.

## Corrective graphics revision

- [x] Previous graphical launch inspected and premium-art failure recorded.
- [x] Painterly assets and new renderer prepared with unchanged rules.
- [x] Authoring/source asset provenance and current exporter documented.
- [x] Changed source and binary artwork uploaded and byte-verified in GitHub.
- [x] Changed-source CI and production graphical launch gate green.
- [x] Title and active-night screenshots captured and reviewed; **night scene rejected for visible compositing seams**. Remaining screens require inspection.
- [ ] Native 60/120 Hz fluidity and input responsiveness measured.
- [ ] Independent code-derived EDD/comparison regenerated from the **new** final source commit.
- [ ] Owner/child UAT and release sign-off.

## Actual premium candidate inspection — 11 November 2026

- CI: successful [run 42974831234](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/42974831234).
- GitHub release-candidate build and Xvfb launch: successful [run 42974831280](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/42974831280).
- The downloaded GitHub-built executable was launched under a 1920×1440 Xvfb display and its title and active-night frames captured; these are actual executable screenshots, not promotional composites.
- **FAIL — art direction:** gameplay background shows conspicuous horizontal collage seams, especially between sky, middle waves and foreground. The lighthouse sits awkwardly at the bottom edge and the boat appears near the shore. This is not yet the cohesive premium game artwork required by the PRD.
- This proves startup and first-screen rendering only. It does not prove normal 60/120 Hz hardware frame pacing or all-screen UAT.
- No public release tag is authorised until the repaired running scene, hardware performance and required human UAT are accepted.
