# Lighthouse Keeper — Product Requirements Document

**Version:** 0.2
**Status:** Locked product contract for V0.1.0; build and human/hardware acceptance pending
**Product name:** Lighthouse Keeper

## Problem

- Lighthouse Keeper is a calm night-time game for children from 7 years and up who find it hard to stay with one thing.
- She keeps a lighthouse beam on a little boat as it crosses the dark sea to the harbour.
- Lovely things appear at the edges of the sea: whales, shooting stars, a sleepy seal.
- She is allowed, even encouraged, to notice them: one key waves at a sighting and logs it in her Sightings Book.
- Then she brings the beam back to the boat. Chasing a sighting with the beam lets fog roll in and the sighting swims away unlogged.
- The game practises one habit: **notice it, name it, come back.**

> **Noticing is fine. Coming back is the skill.**

**Visual ambition, v0.2.** Present this small game with the finish of a leading contemporary premium iPad, iPhone, Steam or Switch indie title: original painterly art, rich but restrained light, fluid motion, readable interactions and carefully timed transitions. The five visual concepts reviewed in November 2026 are *quality and mood references*, not five selectable themes or ready-to-use runtime sprites. The scene must be rebuilt as layered, independently animatable artwork; the movable beam and moving boat cannot be baked into one static image.

## Product objective

Build a polished, lightweight, calm game in Go that practises noticing distraction, naming it and returning to the task, and carries that habit into real work.

The V0.1 product question is:

> **After a week of short nights, does she use "notice it, name it, come back" on her own, outside the game?**

## Users

- Primary: children from 7 years who are capable but easily distracted, playing on a keyboard.
- Secondary: a parent who reads the Harbour Log with her.

The repository is also a Jonno's Garage engineering example, built with the same method and budget as Astro Answers.

## Product principles

- Never shame. No "wrong", no failure screens, no lost progress.
- Reward noticing, returning and honest answers; never reward performance alone.
- One axis of control plus two single keys.
- A night lasts about 10 minutes and ends on its own.
- The game is a place to practise; real work is where the habit is used.
- World-class presentation is an explicit product requirement, not a reason to expand the gameplay engine.
- Animation should be smooth on current 60 Hz and 120 Hz displays; rendering must not change outcomes or introduce input delay.
- The boat, beam, sighting and question remain identifiable at every visual-quality setting.
- Not a treatment, not a diagnosis, no claims of either.

## Core loop

```text
predict the night -> keep the beam on the boat
                         |
          sighting appears at the edge
           /                        \
   notice it (N), log it        chase it with the beam
   beam stays with the boat      fog rolls in, it swims away
           \                        /
            bell: "were you with the boat?" -> answer honestly
                         |
            harbour log: prediction vs what happened
```

## Product rationale

- **Distraction is normal; fighting it doesn't help.** The useful skill is recognising when attention has wandered and choosing to return. The game rewards both halves and never punishes the wandering.
- **Bright minds wander when the task is dull.** So the distractions are the most charming thing on screen. Resisting something boring teaches nothing; returning from something lovely is the real skill.
- **Self-monitoring has evidence behind it.** Asking yourself "was I paying attention?" at random moments, and recording the answer, reduced off-task behaviour in children in classroom studies. The game builds that check-in into play.
- **Honest self-assessment beats high scores.** Before each night she predicts how well she'll keep the light. Afterwards she compares. Accurate predictions are rewarded, not big numbers.
- **Games train what they train.** Attention games improve the trained task far more than everyday life. So the game carries one habit, the bell and the question, that she can use at her desk, and a mode that rings that bell during homework.
- **Calm, not compulsive.** Nights are short, end naturally, and the game suggests resting rather than pulling her back.

## In scope

- night sea with a lighthouse and one boat crossing at a time;
- beam rotated with Left/Right;
- boat lit when inside the beam; fog thickens and the boat slows when it isn't;
- sightings that appear, linger and leave;
- **N** to notice and log a sighting;
- "came back" moments when the beam returns to the boat after drifting;
- prediction before each night;
- random check-in bell with a yes/no question;
- Sightings Book that persists and grows;
- Harbour Log that persists, with prediction, result, came-backs and check-ins per night;
- gentle progression across nights;
- daily rest suggestion after two nights;
- Keeper's Bell mode for real-world use;
- premium original layered illustration, lighting, animation, legible illustrated interfaces and soft sound;
- smooth presentation at the display refresh rate on suitably capable, tested hardware, including a 120 Hz proof;
- deterministic core gameplay suitable for automated testing;
- Linux desktop V0.1, full-screen 4:3, Linux/ARM64 build proof.

## Out of scope

- high-score tables or competition;
- timers, countdowns or failure states;
- clinical measurement, diagnosis or treatment claims;
- reporting to schools, accounts, cloud sync, telemetry, network;
- notifications outside the game;
- runtime AI;
- music with lyrics or a beat;
- gamepad and touch input in V0.1;
- Windows, macOS and mobile packaging in V0.1.

## V0.2 direction

No mobile, Switch, Windows or macOS implementation is committed in V0.1. The high-end mobile and console games are benchmarks for visual and motion quality, not additional platforms or runtime abstraction requirements.

## Functional requirements

### FR1 — The crossing

A boat crosses from the open sea to the harbour along a gently curving course lasting about 90 seconds. A night is five crossings.

The boat always arrives. Keeping it lit makes the crossing faster and the harbour brighter; losing it never fails the crossing.

### FR2 — The beam

Left and Right rotate the beam around the lighthouse with smooth easing. The beam is a soft cone. When the boat is inside it, the boat's lantern glows and it sails at full speed.

### FR3 — Fog

When the boat is outside the beam, fog gathers around it and it slows to a drift. Fog clears within a second of the beam returning.

### FR4 — Sightings

At irregular moments a sighting appears near the edge of the sea, lingers for several seconds and leaves. Kinds include a whale spout, a shooting star, a seal on a rock, glowing jellyfish, an owl on the cliff and a passing ferry.

- **Notice:** pressing N while a sighting is visible waves the keeper's lantern at it and logs it in the Sightings Book. The beam does not move. A new kind earns a new page.
- **Chase:** if the beam is pointed at a sighting for more than one second while the boat is unlit, the sighting shyly leaves without being logged.

### FR5 — Coming back

When the beam has been off the boat for more than two seconds, a soft chime asks "where's the boat?" Returning the beam to the boat within five seconds of the chime counts as a **came back**, shown with a small sparkle and recorded in the night's log.

Coming back is celebrated every time. Drifting is never mentioned.

### FR6 — Prediction

Before each night she chooses how many of five lanterns she expects to keep lit (a crossing keeps its lantern lit if the boat was in the beam for at least 70% of it). Left/Right choose, Space confirms.

At the end of the night the Harbour Log shows her prediction beside the result. A prediction within one lantern earns a **good guess** star. The game never comments on whether the result was high or low.

### FR7 — Check-in bell

Two or three times a night, at unpredictable moments, a bell rings and the scene softly dims. The question appears with a picture: "Were you with the boat?" She answers Y or N.

Any answer earns a check-in star. The game never contradicts her answer.

### FR8 — Progression

Across nights the boat's course curves more, sightings come a little more often and the beam narrows slightly, to fixed caps. Progression never removes anything she has earned.

### FR9 — Rest

After two nights in a day, the game shows the lighthouse at rest and suggests coming back tomorrow. The Sightings Book and Harbour Log remain open to browse. A parent key allows another night.

### FR10 — Keeper's Bell mode

From the title, K starts Keeper's Bell for real work: choose 15 or 25 minutes. The screen shows a calm, still lighthouse. The bell rings at random intervals of three to seven minutes; each time she answers Y or N. Answers are recorded in the Harbour Log under "at my desk".

### FR11 — Sightings Book

A persistent illustrated book with one page per sighting kind, showing how many times each has been noticed.

### FR12 — Harbour Log

A persistent log with one row per night: date, prediction, lanterns lit, came-backs, check-ins and good-guess stars, plus Keeper's Bell sessions. Readable at a glance, with pictures as well as numbers.

### FR13 — Visual and sound quality

The finished game must feel professionally art-directed and finished: layered hand-painted storybook seascapes, detailed lighthouse and harbour, deep indigo/teal night water, warm amber light, luminous moon and stars, gentle mist, original expressive sightings, elegant illustrated books and restrained UI. The art should match the ambition of the five approved visual concept examples, without reproducing any unrelated game's protected artwork.

- The **lighthouse, boat and sighting** are distinct, legible objects. Keep the playable sea visually open; intricate rocks, flowers, cottages and harbour details frame the action, rather than obscure it.
- The beam moves continuously with its source, with a soft edge and plausible light on the boat/water. Fog gathers and clears smoothly; water and mist move without visible loops or sudden jumps.
- The six sightings have clear identities and a delightful appearance/notice/leave animation; they remain appealing enough to tempt attention without becoming flashing or frantic.
- Lit harbour windows visibly accumulate during the night, matching the five crossing lanterns. The Book, prediction picker, Harbour Log, Rest and Keeper's Bell share the same premium visual language.
- All screen transitions, overlay entrances, page turns and feedback use short, consistent easing; no stutters, hard cuts where a soft transition is appropriate, text jitter or distracting rapid motion.
- Sounds are soft, original and synchronised with events: waves, chime, bell, lantern click, notice and harbour sounds. No alarms and no music with a beat.
- **Presentation only:** animation and decorative particles must not alter boat position, visibility tests, rewards, scheduling, saved data or any other gameplay result.

The game is not required to provide five different graphics modes. A single coherent visual direction must be selected and implemented to this standard.

### FR14 — Display

The gameplay coordinate space remains fixed at **960 × 720 (4:3)**. Default presentation is full-screen, with aspect preserved and neutral bars on wider or taller displays. Never stretch the stage. Draw the final image sharply at the effective display density, using appropriately detailed artwork; logical coordinates do not require a low-resolution final framebuffer. The whole playable area, controls, boat and sightings must remain visible at supported desktop sizes. Provide a windowed mode for Keeper's Bell.

### FR15 — Motion, responsiveness and runtime size

- **60 Hz:** sustained 60 presented frames/second, synchronised with the display, on the designated capable Linux reference computer.
- **High refresh:** sustained **120 presented frames/second at 120 Hz** on a designated capable display and computer, with perceptibly smooth boat, beam and environmental motion; verify with a measured hardware run. If the actual display does not offer 120 Hz, it cannot serve as the high-refresh acceptance device.
- Simulation remains deterministic at 60 logical ticks/second. Higher screen refresh uses presentation-only sub-tick motion, if proven useful, and time-based decorative motion, not additional simulation ticks or delayed input presentation.
- Smooth gameplay takes priority over ornamental post-effects, but degradation may remove only expendable cosmetics, never readability or gameplay feedback.
- No avoidable per-frame asset decoding, texture creation, disk I/O, unbounded rendering allocations or runtime art-generation dependencies. Game contents remain entirely offline.
- Keep final binary, compressed art/audio, decoded texture memory and frame cost within explicit measured budgets; use the smallest assets that pass the visual acceptance test, not an arbitrary quality-reducing size cap.

### FR16 — Local data integrity and privacy

All progress and check-in answers are stored locally only; no telemetry, network, accounts or clinical profiles. An absent save starts a new game. A malformed, oversized or incompatible save is rejected safely, with a retained recovery copy and a child-appropriate explanation; the game remains playable. Writes are atomic. The parent override allows at most one extra night per local calendar day and persists across relaunches.

## Acceptance criteria

V0.1 is accepted when automated and playable evidence shows that:

1. a crossing always completes, faster when lit;
2. Left/Right rotate the beam; the boat is lit only inside the cone;
3. fog gathers when unlit and clears when the beam returns;
4. N logs a visible sighting without moving the beam;
5. chasing a sighting makes it leave unlogged;
6. a return within five seconds of the chime is recorded as a came back;
7. prediction, result and good-guess star are computed as specified;
8. the check-in bell rings two or three times a night and accepts any answer;
9. progression changes course, sighting rate and beam width within caps;
10. the rest screen appears after two nights in a day;
11. Keeper's Bell rings within its interval range and records answers;
12. the Sightings Book and Harbour Log persist between runs;
13. no screen uses failure language;
14. gameplay logic is testable without the renderer and replays deterministically from a seed;
15. full-screen 4:3 presentation with letterboxing;
16. `gofmt`, `go vet`, race tests, Linux/amd64 and Linux/ARM64 builds pass in CI;
17. a representative complete night is visually reviewed against the approved premium concept quality bar, and screenshots and brief capture evidence cover gameplay, title, book, log, rest and desk mode;
18. measured 60 Hz and 120 Hz hardware runs meet the FR15 smoothness and responsiveness targets, with no visible simulation judder, input lag, frame-time spikes or screen-transition hitching;
19. 4:3 letterboxing, high-density rendering, clear boat/sighting readability, correct moving beam and UI legibility pass on different monitor aspect ratios and scale factors;
20. graphics and sound are original or have documented redistribution rights; compressed asset and decoded memory budgets are measured; reducing decorative quality never changes rules or rewards;
21. oversized, malformed or incompatible saves are rejected within bounded work and memory, atomic writes preserve prior data on interruption, and the once-per-day parent extra night persists.

## Definition of done

V0.1 is done when **all 21 acceptance criteria** have executable or observable evidence, the final game matches the approved visual direction on real target hardware, she has played a night without help, the book and log persist, and Keeper's Bell has been used for one real homework session. A concept illustration, CI build or simulated benchmark does not count as visual-performance proof.

The hand-written EDD is then recovered independently from the implementation and compared with the intended design. Save the frame-time evidence, screen captures, asset budget and play-test observations with the release documentation.

**Scope control:** FR1–FR12 and their original numerical rules are unchanged. FR13–FR14 are strengthened, and FR15–FR16 make graphical performance and the existing privacy/persistence expectations explicit. No mobile, Switch or iOS packaging is added to V0.1; those platforms are quality benchmarks, not release targets.

The product test is:

> **She hears a bell at her desk, asks herself the question, and comes back.**
