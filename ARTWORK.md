# Lighthouse Keeper — Artwork, provenance and export

**Visual revision:** painterly presentation pass, following the first release-candidate graphical UAT rejection. **Status:** implementation candidate; premium quality and fluidity require on-device approval.

## Originals and rights

Three source studies are original ChatGPT-generated project illustrations created for Lighthouse Keeper on 11 November 2026, rather than third-party game art or web images:

- `assets-source/concept-approved.png` — owner's approved lighthouse, moonlight, sea and harbour direction; retained as the title-image source.
- `assets-source/approved-ocean-study.png` — original painterly coastal study used as texture and character-isolation source for the interactive sea, lantern and boat.
- `assets-source/original-ui-study.png` — original illustrated scene-and-book study; the six subject paintings supply the Sightings Book and drifting sightings. Any wording visible in the source study is not gameplay policy.

`assets-source/export_master_scene.py` produces the **single contiguous 1920 × 1440 night scene** directly from the approved original coastal illustration, retouching the non-interactive boat offline. `assets-source/export_premium.py` exports the other eleven original artwork sprites and UI surfaces. All moving boats, beam, fog and sightings remain separate from the night scene and are controlled only by game state. The painterly scene has no combined ocean strips, horizon joins, or runtime image-generation calls.

The runtime sprites are ordinary optimised JPEG and PNG files, inspectable in GitHub. The game does **not** call image-generation services or execute Python at runtime. Fonts use Go's standard redistributable font modules. Audio is generated in Go from original mathematical cues.

## Reproducible build

From the repository root, on an **art authoring** computer with Python 3, Pillow, NumPy, and OpenCV (authoring-only dependencies):

```text
python assets-source/export_premium.py
python assets-source/export_master_scene.py
```

The second command deliberately runs last so the one coherent night scene is never overwritten by a stitched export.

The script is seeded; source originals, coordinates, quality and output sizes are fixed. Production is otherwise self-contained. Do not commit obsolete prototypes, base64 image copies or transient mockups.

Shipped files in `cmd/lighthousekeeper/assets/`:

- `seascape.jpg`, `title_art.jpg` — opaque 1920 × 1440 backdrops;
- `lighthouse.png`, `boat.png`, `sightings.png`, `lantern.png` — painted interactive objects / miniature status lights;
- `beam.png`, `fog.png`, `glow.png` — reusable soft lighting/fog masks;
- `parchment.png`, `ui_rail.png` — reusable illustrated book and brass/navy UI treatments.

Ebitengine handles texture drawing and scaling. One reusable nine-slice panel is prepared at startup, never sliced again during steady-state drawing. Twelve vertical-adjacent horizontal water strips sampled from the **same** background are prepared once and drawn at subtle alpha and sinusoidal subpixel offsets (historical line scrolling). No glint texture, extra ocean painting, runtime physics or Kage shader is required. Sighting illustrations reuse one six-frame atlas.

## Budget and acceptance

The game keeps the PRD's full 960 × 720 logical space, while raster artwork renders to a 1920 × 1440 target. Painterly detail is prioritised over an arbitrary byte limit. Record packed bytes, decoded texture memory and native 60/120 Hz frame-time evidence after graphical CI and hardware UAT; these cannot be deduced from the illustration alone.

Before public release, capture actual running footage of all screens and compare against the five approved painterly references. Confirm the beam's visually illuminated boat matches `WithinBeam` in all movement directions, fog does not obscure the player, the harbour remains readable, and every sighting is identifiable and collectible. Only real hardware can prove high-refresh fluidity.
