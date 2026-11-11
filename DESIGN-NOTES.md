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
