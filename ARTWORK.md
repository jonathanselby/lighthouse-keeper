# Lighthouse Keeper — Artwork, provenance and export

**Visual revision:** painterly presentation pass, following the first release-candidate graphical UAT rejection. **Status:** implementation candidate; premium quality and fluidity require on-device approval.

## Originals and rights

Three source studies are original ChatGPT-generated project illustrations created for Lighthouse Keeper on 11 November 2026, rather than third-party game art or web images:

- `assets-source/concept-approved.png` — owner's approved lighthouse, moonlight, sea and harbour direction; retained as the title-image source.
- `assets-source/approved-ocean-study.png` — original painterly coastal study used as texture and character-isolation source for the interactive sea, lantern and boat.
- `assets-source/original-ui-study.png` — original illustrated scene-and-book study; the six subject paintings supply the Sightings Book and drifting sightings. Any wording visible in the source study is not gameplay policy.

`assets-source/export_premium.py` assembles and masks these **offline** into the twelve embedded runtime images. The painterly sea is newly composed from multiple illustration regions. The moving boat, spotlight, sightings and lantern are **not** baked into that ocean image: they are separate sprites controlled by real game state.

The runtime sprites are ordinary optimised JPEG and PNG files, inspectable in GitHub. The game does **not** call image-generation services or execute Python at runtime. Fonts use Go's standard redistributable font modules. Audio is generated in Go from original mathematical cues.

## Reproducible build

From the repository root, with Python 3, Pillow and NumPy available on the **authoring** machine:

```text
python assets-source/export_premium.py
```

The script is seeded; source originals, coordinates, quality and output sizes are fixed. Production is otherwise self-contained. Do not commit obsolete prototypes, base64 image copies or transient mockups.

Shipped files in `cmd/lighthousekeeper/assets/`:

- `seascape.jpg`, `title_art.jpg` — opaque 1920 × 1440 backdrops;
- `lighthouse.png`, `boat.png`, `sightings.png`, `lantern.png` — painted interactive objects / miniature status lights;
- `beam.png`, `fog.png`, `glow.png`, `shimmer.png` — reusable bounded soft lighting/water overlays;
- `parchment.png`, `ui_rail.png` — reusable illustrated book and brass/navy UI treatments.

Ebitengine handles texture drawing and scaling. One reusable nine-slice panel is prepared at startup, never sliced again during steady-state drawing. The visual-only water glints are sixteen bounded translated draws, not a fluid simulation. Sighting illustrations reuse one six-frame atlas.

## Budget and acceptance

The game keeps the PRD's full 960 × 720 logical space, while raster artwork renders to a 1920 × 1440 target. Painterly detail is prioritised over an arbitrary byte limit. Record packed bytes, decoded texture memory and native 60/120 Hz frame-time evidence after graphical CI and hardware UAT; these cannot be deduced from the illustration alone.

Before public release, capture actual running footage of all screens and compare against the five approved painterly references. Confirm the beam's visually illuminated boat matches `WithinBeam` in all movement directions, fog does not obscure the player, the harbour remains readable, and every sighting is identifiable and collectible. Only real hardware can prove high-refresh fluidity.

## v0.3 — Arcade-inspired source-owned living ocean

`assets-source/bake_living_ocean.py` is now the canonical authoring tool, using the original `assets-source/approved-ocean-study.png` already tracked in GitHub. It performs one-time donor retouch of the baked lighthouse beam and immovable foam, then exports a water-only painting, a geometric water-domain, transparent fixed shore and boulders and a ten-pose spray atlas. The former code-generated shimmer is no longer a runtime asset.

**Paint ordering:** clean water shader → rear spray → immutable rock and shoreline PNG cut-outs → foreground spray → the original moving spotlight and boat. Static stone never enters the animated texture. All significant spray on rocks must be separately animated, not flattened under it. The whitewater atlas uses original project art and is a visual overlay, not a particle engine.

This authoring pipeline requires only build-time Python, Pillow, NumPy and OpenCV. It is run by the art-bake GitHub workflow and validated by `assets-source/test_living_ocean.py`. Native end users only receive one Linux executable with embedded assets, and need no Python or internet connection. If the rebuilt illustrations are visually unsatisfactory, improve the one-time painting, never add a runtime classifier to compensate.
