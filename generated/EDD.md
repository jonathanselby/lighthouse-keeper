# Lighthouse Keeper — Code-Derived Engineering Design Document

**Evidence source commit:** `f5029e2999223a2410848931fbaf8ef7f0c8c342`  
**Method:** independent source inspection, using the Garage codebase-to-EDD evidence discipline  
**Date:** 2026-11-11  
**Status:** compiled and smoke-tested implementation; human UAT and physical performance measurements outstanding

This document records the implementation recovered from code, not claims inferred from the hand-written PRD or EDD.

## 1. Executive description

**KNOWN:** Lighthouse Keeper is a Linux desktop game written in Go 1.25 with Ebitengine 2.10.6. The deterministic gameplay engine is plain Go in `internal/game`; `cmd/lighthousekeeper` handles input, drawing, audio, persistence and desktop display.

## 2. Evidence inspected

- `internal/game/game.go`: authoritative state; nine screens; boat, beam, fog, chimes, five crossings, bells, prediction, daily rest and real-work desk session.
- `internal/game/sightings.go`: bounded, weighted sighting generation; notice/chase mechanics.
- `internal/game/save.go`: save schema, validation, recovery, atomic persistence.
- `internal/game/strings.go`: central child-facing copy.
- `internal/game/game_test.go`, `save_test.go`: headless rules, long runs and persistence.
- `cmd/lighthousekeeper/main.go`, `render.go`, `audio.go`, `assets_test.go`: app input and display, 2D painting and effects, generated audio, asset validation.
- `cmd/lighthousekeeper/assets/`: two illustrated JPG images and six PNG overlays/sprites, all embedded.
- `assets-source/concept-approved.png` and `export.py`: art source and production export path.
- `.github/workflows/test.yml`, `release-build.yml`, `go.mod` and `go.sum`: actual CI and dependency contracts.

## 3. Runtime architecture

```text
keyboard + date + real-time desk elapsed ticks
                  |
                  v
       Ebitengine adapter
       /       |        \
      v        v         v
 input     renderer   sound / JSON I/O
      |
      v
  Game.Update(Input)
      |
      +-- beam / path / fog
      +-- sightings / return
      +-- prediction / night bells
      +-- daily rest / Keeper's Bell
      |
      v
authoritative Game + explicit Events
```

There is no general game engine, ECS, event bus, physics system, UI framework, network path, telemetry or runtime AI.

## 4. State and time

**KNOWN:** One `Game` owns `Screen`, `Night`, `Bell`, `Save`, `Beam`, a seeded `RNG`, small navigation values and at most eight Events per update. Cosmetic images, particles, font faces, audio and save error messages remain in the adapter.

**KNOWN:** Ebitengine invokes fixed 60-TPS updates. Game logic uses integer ticks; the date and monotonic desk elapsed time are injected via `Input`. The renderer separately estimates a sub-tick for fluid cosmetic movement. This does not constitute a measured high-refresh result.

## 5. Gameplay rules

**KNOWN:** Each night contains five cubic-Bézier boat crossings with a 65-entry cumulative arc-length table, inter-crossing pauses and a positive fogged drift speed. The beam uses eased angular input and explicit bearing/range tests around the lighthouse. Fog slows the boat when it is not lit.

**KNOWN:** Six sighting kinds (Whale, Shooting Star, Seal, Jellyfish, Owl, Ferry) appear within a bounded sea-edge band and are kept away from the boat path. N notices and logs the oldest visible sighting without moving the beam. Keeping the beam on a sighting for 60 ticks while the boat is unlit causes that sighting to leave unlogged.

**KNOWN:** A two-second drift after the grace interval can open a five-second chime return window. Prediction scores use the number of crossings lit for 70% or more of their ticks, with a good-guess star awarded within one lantern.

**KNOWN:** Each night schedules two or three Y/N check-ins across selected crossings, pausing gameplay; honest answers are recorded without grading. Two nights per date trigger Rest, with one parent override requiring P held for three seconds. Keeper's Bell uses 15/25-minute sessions and 3–7-minute randomized ring intervals based on externally injected elapsed time.

## 6. Persistence and integrity

**KNOWN:** A single local `save.json` stores the six-kind Book, night and desk logs and daily progression. Schema version 2, strict unknown-field rejection, a 256-KiB input cap, structural validation, 400-record history limits and a 200-KiB output target are implemented. Corrupt input is preserved as a numbered `.bad` file where possible; writes use a permissions-restricted temporary file, flush, rename and directory sync.

## 7. Presentation and audio

**KNOWN:** The adapter embeds and decodes eight runtime graphics, renders a 2× internal 4:3 image and composes the title, prediction, night, harbour, book, log, rest, bell setup and desk screens. The sea illustration, lighthouse, beam, fog, boat and sightings are raster elements with restrained overlay effects. The app contains ten generated PCM sound cues and an eight-second quiet ambient loop, with silent Linux fallback.

**UNVERIFIED:** The design's premium visual standard, stable 60/120 FPS, frame cadence, GPU memory budget and physical input latency have not been measured or approved in user UAT.

## 8. Testing and build evidence

**OBSERVED:** GitHub CI run [42962738531](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/42962738531) passed module tidiness, format, vet, race-enabled unit tests, Linux amd64 and ARM64 builds, and named core acceptance tests.

**OBSERVED:** Release build run [42962933696](https://github.com/jonathanselby/lighthouse-keeper/actions/runs/42962933696) passed source validation, native amd64 compilation, software-rendered Xvfb graphical startup and tar.gz packaging, uploading a binary workflow artifact.

**OPEN:** Child completes a night unaided; actual 15/25-minute homework mode use; screenshot and visual UAT; hardware performance and constrained-hardware benchmark; persistence and navigation checked manually across application relaunches.

## 9. Release boundary

`VERSION` is `0.1.0`. A green workflow is engineering evidence, not permission to claim human UAT complete. No release tag or published release is justified until acceptance is reviewed. The artifact currently belongs to Actions rather than a public GitHub Release.

## 10. Conclusion

The recovered application is a small deterministic Go game with a separate Ebitengine presentation layer and bounded persistence. Its correctness/build/startup gates have evidence; its premium graphics and high-refresh targets remain physical acceptance criteria.

## Premium artwork implementation reconciliation — candidate f5029e29

The visual adapter now loads **12 runtime raster images** (two JPG scenic images, ten PNG sprites/overlays), including lantern, glow, shimmer and a nine-slice UI rail. It pre-slices the rail and six sighting tiles at startup; `drawSea` places a single 1920×1440 composited ocean painting and sixteen bounded glint draws. Dynamic scene objects, beam, fog, book, log, and other screens remain Ebitengine-owned presentation decisions. Offline art preparation is in `assets-source/export_premium.py`, with three original project source illustrations. Gameplay rules in `internal/game` remain the authoritative core.

**Observed running implementation:** CI `42974831234` and release smoke `42974831280` succeeded. An independently captured native title and night screenshot exposed visible horizontal seams in the artwork compositing, so the artistically approved build is **not** yet achieved. Only visual composition requires further correction; no gameplay scope or security change is justified. Human UAT and native 60/120 Hz performance measurement remain open.
