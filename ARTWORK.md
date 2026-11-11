# Lighthouse Keeper — Artwork and rights register

## Provenance

- `assets-source/concept-approved.png`: one of the five original cinematic concepts generated for Lighthouse Keeper in the ChatGPT working session on 11 November 2026 and approved by the project owner. It is not stock imagery or an unrelated commercial game's art.
- `cmd/lighthousekeeper/assets/title_art.jpg`: colour-graded, 4:3-cropped export of that original concept; used only as a static title backdrop (where the baked illustrative beam is not gameplay).
- `cmd/lighthousekeeper/assets/seascape.jpg`: independently authored layered sky and ocean, with a moon sampled from the original concept; no baked-in moving gameplay object.
- `lighthouse.png`, `boat.png`, `sightings.png`, `parchment.png`, `beam.png`, `fog.png`: original procedural illustration exports authored specifically for the game with offline Pillow code.
- Typography: `goregular` and `gobold` from the Go font sources in `golang.org/x/image`. No proprietary font files are bundled.
- Audio: original mathematically generated PCM cues and an ambient wave bed; no sampled third-party music or audio files.

## Offline reproducibility

With Python 3 and Pillow installed:

```sh
python assets-source/export.py
```

The export script reads `assets-source/concept-approved.png` and writes exactly the eight runtime files into `cmd/lighthousekeeper/assets/`. At runtime, the game needs no Python, Pillow, image-generation service, asset-server or external download.

The source artwork and shipped assets are bundled in Git for auditability; a user releases the **compiled binary** rather than asking players to run the export tool.

## Asset/resource budget

The initial embedded art total is less than 2 MiB compressed (plus approximately 2.7 MiB of offline source concept). The decoded GPU textures, actual executable, frame times and runtime RAM must be measured in an integrated build. Do not substitute these file sizes for that evidence. The core never opens an art file per frame, and no authoring image is decoded during play.

## Quality gate

The six sighting illustrations are separate collectible identities. All movement, beam illumination and fog are render-only transforms over genuine gameplay state. The complete art pass needs visual comparison, multi-display screenshots and child observation before being called premium or final.
