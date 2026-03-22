# Celebration System

## Context

Mavis needs positive feedback when users hit milestones. Currently, drill completion shows a plain text message regardless of performance. A tiered celebration system provides visual reinforcement scaled to achievement level.

## Tiers

| Tier | Trigger | Scope | Visual | Duration |
|------|---------|-------|--------|----------|
| Pass | Accuracy >= 95% | Per-phase (after 10th drill) | Green checkmark + green-styled message | Instant (no animation) |
| Perfect | Accuracy = 100% | Per-drill (any single drill) | 12 sparkle particles + bold golden "PERFECT!" | ~1.5s |
| Level-up | All phases for a level complete | Per-level | 20 sparkle particles (wider spread) + "LEVEL N COMPLETE!" with color cycling | ~1.5s |

Pass tier has no animation. Perfect and Level-up tiers play a particle animation before showing the "Press Enter" prompt.

**Precedence:** Higher tiers supersede lower ones. If the final drill in a phase is 100% and the phase passes, only the phase/level celebration fires (not both Perfect and Pass). Level-up supersedes all. Specifically:
1. Check level-up first (all phases complete)
2. Then check phase pass (>= 95% on 10th drill)
3. Then check perfect drill (100% on any individual drill)
4. Only one celebration fires per completion event

## Particle System

Each particle has:
- Position: x, y offset from message center
- Character: randomly chosen from `*`, `+`, `.`
- Color: hex color, cycles through a palette each tick
- Lifetime: remaining ticks before removal

On trigger, N particles spawn at random positions within a bounding box around the message area. Each tick (every 100ms), lifetime decrements and color shifts. Dead particles are removed. When all particles expire, the celebration ends.

The palette cycles through the finger-zone colors (pink, purple, blue, green, yellow) for visual consistency with the rest of the UI.

## State Machine Integration

New state `stateCelebrating` sits between completion detection and the existing completion states.

Flow:
1. Drill completes with 100% accuracy (or level completes)
2. State transitions to `stateCelebrating`
3. Particles spawn, tick command starts firing every 100ms
4. View renders particles overlaid on the celebration message
5. When all particles expire, state transitions to the appropriate completion state (`stateDrillComplete`, `statePhaseComplete`, or `stateLevelComplete`)
6. Normal "Press Enter" flow resumes

The 95% pass tier skips `stateCelebrating` entirely. It just styles the existing message green with a checkmark prefix.

User can press Enter during the celebration to skip it and go straight to the completion state.

## Implementation

**New file: `internal/app/celebrate.go`**

Contains:
- `particle` struct (x, y, char, color index, lifetime)
- `celebration` struct (particles slice, message string, message style, ticks elapsed, next state)
- `spawnCelebration(tier, nextState)` creates particles and sets the message
- `tickCelebration()` ages particles, returns whether celebration is still active
- `renderCelebration(width)` renders particles and message as a string

**Tick mechanism:**
- `tickMsg` struct (empty, just a signal)
- `tickCmd` function returns a `tea.Tick(100ms, func() tea.Msg { return tickMsg{} })`
- Only fires during `stateCelebrating`

**Changes to `internal/app/app.go`:**
- Add `stateCelebrating` to the state enum
- Add `celebration *celebration` field to Model
- In `onDrillComplete`: check accuracy, spawn celebration if 100% or level-up
- In `Update`: handle `tickMsg` when in `stateCelebrating`
- In `Update`: allow Enter during `stateCelebrating` to skip
- In `View`: render celebration when in `stateCelebrating`
- Style the 95% pass message green with checkmark (no state change needed)

## Scope Boundaries

**In scope:**
- Three-tier celebration (pass, perfect, level-up)
- Particle animation for perfect and level-up
- Green styled message for pass
- Enter to skip animation
- ~1.5 second animation duration

**Out of scope:**
- Sound effects
- Persistent streak tracking
- Configurable animation duration
- Custom particle characters
