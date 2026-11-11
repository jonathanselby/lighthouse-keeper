# Lighthouse Keeper — V0.1.0 release checklist

## Contract

- [x] PRD v0.2 retained in `PRD.md` with all 16 original and five added criteria.
- [x] EDD v0.2 retained in `EDD.md` as the current design baseline.
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
