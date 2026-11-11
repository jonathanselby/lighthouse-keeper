# Lighthouse Keeper

**Notice it. Name it. Come back.**

Lighthouse Keeper is a gentle, illustrated nighttime game for children aged seven and above. Guide a soft lighthouse beam to keep a small boat company as it crosses toward the harbour. Whales, seals, jellyfish and shooting stars appear at the edges: notice them, collect sightings and then return your attention to the boat.

There are no scores, failure screens, online accounts or clinical claims. Each night has five crossings, a prediction, quiet self-check bells, a Sightings Book and a Harbour Log. A separate Keeper's Bell session brings the same small habit to fifteen or twenty-five minutes of real-world work.

## Play

```text
SPACE          start / confirm / return
LEFT / RIGHT   rotate the beam, choose or turn pages
N              notice a sighting; answer "no" during a bell
Y              answer "yes" during a bell
B              Sightings Book
L              Harbour Log
K              Keeper's Bell
P (hold 3 s)   parent extra night on Rest screen (once daily)
ESC            exit; in Keeper's Bell, record an early finish
```

## V0.1 scope

Linux desktop, 4:3 aspect-preserved full-screen play, with a windowed homework bell. The logical gameplay area is 960×720; the presentation renders onto a 1920×1440 image to retain detail. Core game rules update at 60 ticks/sec; drawing is independent of the fixed simulation rate. Premium-quality imagery and 60/120 Hz visual smoothness are acceptance targets requiring measured hardware and human evidence, not claims established by source code alone.

All artwork is embedded into the executable and generated offline from an approved original painterly concept plus authored layer assets. No network, telemetry, accounts or runtime AI. Progress lives in the user's local configuration directory (`$XDG_CONFIG_HOME/lighthouse-keeper/save.json`, or the platform user configuration directory) rather than the launch directory.

## Build

Requires Go 1.25+ and Ebitengine v2.10.6. On Linux with a graphical desktop:

```sh
go mod download
go run ./cmd/lighthousekeeper
```

For the standalone executable:

```sh
CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o Lighthouse-Keeper ./cmd/lighthousekeeper
./Lighthouse-Keeper
```

The included `.github/workflows/test.yml` is the authoritative compile/test gate; Go dependencies are obtained in CI. The `assets-source/export.py` script is *offline production tooling* using Python/Pillow, not used by the game at runtime.

## Engineering gate

```sh
go mod tidy
gofmt -l ./cmd ./internal
go vet ./...
go test -race -count=1 ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/lighthousekeeper
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./cmd/lighthousekeeper
```

The gameplay rules live in `internal/game` and are independent of Ebitengine. They can be tested headless using fixed seeds and injected dates and elapsed homework time. The thin adapter in `cmd/lighthousekeeper` owns input, local persistence, graphics and sound.

## Engineering and evidence

- [`PRD.md`](PRD.md) — locked V0.2 product requirements and 21 acceptance criteria.
- [`EDD.md`](EDD.md) — V0.2 engineering design and resource budgets.
- [`DESIGN-NOTES.md`](DESIGN-NOTES.md) — material design decisions, simplifications and rejected alternatives.
- [`TESTING.md`](TESTING.md) — observed test results, unresolved graphical gates, and eventual UAT evidence.
- [`RELEASE.md`](RELEASE.md) — checklist with evidence status.
- [`architecture.txt`](architecture.txt) — small runtime architecture.
- [`ARTWORK.md`](ARTWORK.md) — source, rights, asset format and offline export instructions.
- `generated/EDD.md` and `generated/COMPARISON.md` — to be created **only after** final source commits are validated, independently from implementation.

## Status

This is the initial implementation candidate, not a certified polished release. Follow the evidence and unchecked release gates in [`RELEASE.md`](RELEASE.md); CI success does not prove frame pacing, final graphical quality or child/hardware acceptance.

## Licence

MIT. See [`LICENSE`](LICENSE).

## Painterly visual implementation

The game uses original 4:3 painterly illustration composites and separately animated boat, lighthouse, sightings, beam and fog. Rendering remains a thin Ebitengine adapter over the unchanged deterministic core. New artwork is made offline via `assets-source/export_premium.py`; see [`ARTWORK.md`](ARTWORK.md) for source provenance, asset inventory and quality gates.

The rendered world is not the game state. A simple presentation-only vertical transform keeps travelling boats over the sea, while the authoritative beam test, timing, rewards, Sightings Book and Harbour Log use exactly the same rules as before. Brass/navy interface panels share one nine-slice image; there is no additional runtime UI or graphics framework.
