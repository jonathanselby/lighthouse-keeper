# Lighthouse Keeper — Design Notes

**Build status:** implementation candidate, pending integration CI and real visual/hardware acceptance.

## One mechanic, not a learning platform

The boat is the task, sightings are the interesting world, and the beam is attention. The game intentionally keeps noticing, returning and honest self-report together. There are no penalties, scores, clinical judgements or engagement mechanics unrelated to this habit.

## Changes from EDD v0.1, retained in v0.2

The original beam geometry could not illuminate the harbour mouth. Widening the permitted angle to ±1.50 radians and range to 820 pixels preserves all of the sea and the complete crossing. The original chime's exact OffTicks equality could miss the trigger after the initial grace period; `OffTicks>=120 && CrossingTicks>=180` corrects this.

Fixed `[5]int` arrays replace dynamic crossing score slices. Cosmetic particles and audio belong to the application adapter, not the authoritative gameplay engine. The two or three bell triggers are selected from guaranteed crossings rather than an optimistic global timeline. Keeper's Bell measures injected monotonic elapsed real time rather than counting Update calls.

## Art is authored once, played many times

The approved cinematic title concept is used for the static attract view. The playable sea is an **independently rendered original layer** with no baked-in boat, lighthouse beam or sightings. The boat, lighthouse, beam, mist and six sighting kinds are individually movable/visible assets. Offline art export uses Pillow; the runtime only decodes embedded images and draws transforms. The presentation is built at two times logical resolution to keep text and hero art crisp on ordinary HD monitors. This is a starting strategy, not measured proof of 4K support or 120 Hz performance.

## Smooth motion without a second simulation

Gameplay is a deterministic 60-tick update. The renderer can project the latest boat arc distance and beam angular velocity a fraction of a tick forward using elapsed draw time, avoiding a deliberately delayed state buffer. Environmental movement is based on wall time in the adapter only. This can improve higher-refresh motion but must be tested for overshoot and actual frame pacing; the gameplay state is unchanged.

## State boundaries

One `Game` owns the night, beam, six-kind book, saved counts and bell. `Save` is serialised by plain functions into a bounded, validated local JSON document. No registry, persistence interface, event queue, scene graph, ECS or game-wide animation framework exists. Small render files are a legibility aid rather than new application architecture.

## Release lessons from the two previous Garage games

Astro Answers used AppImage successfully. Dandelion Drift's target Ubuntu hardware exposed AppImage/FUSE friction, and a native executable inside tar.gz proved more reliable. Lighthouse Keeper initially follows the simpler native archive path but will preserve human hardware evidence as the deciding factor; an AppImage would be justified only by successful target testing.

## Gaps intentionally left open

The first runtime art pass is a **candidate**, not an assertion that the premium benchmark has been reached. Real graphical launch, measured frame times at 60 and 120 Hz, child/UAT, actual sound device behaviour, and independent code-derived design evidence remain mandatory. The implementation must be corrected from evidence rather than claiming success from unit tests alone.

> **Does every piece of machinery remove more complexity than it introduces?**

## Corrective painterly pass — visual quality is a product requirement

The first CI-green candidate proved the underlying mechanics but not the approved graphical standard. Its sea used thousands of mathematically drawn arc strokes, its lighthouse was incorrectly cropped at the bottom of the screen, and flat animal silhouettes did not belong to the title painting. A graphical launch check could not catch this artistic mismatch.

The corrective design keeps `internal/game` byte-for-byte unchanged. Original approved painterly studies are composited and isolated offline into a true moving boat, six independent sightings, a foreground lantern, warm gradient beam, fog and one detailed 4:3 sea. The lantern is now visibly aligned with the light origin, and the actual moving gameplay components are never baked into the background.

A monotonic `sceneY` transform repositions the visual boat within the painted sea and is applied to the beam, fog, sightings and particles. Logical angles, distances, rewards and persistence remain those of the headless game. Visual-only rocking, painted shimmer and soft lighting use frame time; they never consume gameplay randomness.

The UI shares a single brass/navy frame, sliced once at startup and drawn with nine transforms for crisp corners. It replaces thick debug-style boxes without introducing a GUI framework. The Book still uses a single reusable parchment illustration, and all six paintings share one atlas.

We deliberately rejected a real-time fluid renderer, generic shader/post-processing stack, scene graph, and frame-buffer-based lighting. None removes more complexity than it adds for this bounded single-scene game. The art-source pipeline is offline and not part of the shipped application.

This is **not** a declaration of premium approval. The source painting, preview and compiled interactive render are different kinds of evidence. The CI, visual UAT and refresh-rate tests remain distinct gates.
