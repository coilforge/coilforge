# CoilForge — Gamification Ideas

This document collects all gamification concepts discussed for CoilForge  
It serves as a living reference for future development after the core architecture is stable.

**Core Philosophy Reminder**  
Keep the foundation clean and understandable. All gamification features must reuse existing Part interface, world/, editor/, sim/, and render/ as much as possible. New features live in gamify/ or higher layers.

## 1. Overall Gamification Modes

### NAND2Tetris-style Linear Campaign
- Progressive tutorial from basic relay logic to complex systems.
- ~25–30 levels, unlocked sequentially.
- Focus: learning concepts, correctness, and gentle efficiency.

**Example early progression:**
1. The Relay
2. NOT Gate
3. AND/OR Gates
4. SR Latch
5. Clock
6. D Flip-Flop
7. 4-Bit Counter
8. 4-Bit Adder
9. 8-Bit Register
10. Simple ALU
11. Basic CPU cycle
12. Minimal Pong

Later levels introduce peripherals and advanced relay techniques.

### Design Challenges (Open-Ended)
Creative problem-solving with loose constraints.  
Examples:
- 60-second timer with ≤ 20 relays
- 4-way traffic light controller
- Binary-to-7-segment decoder
- Small Turing Machine
- Relay-based RNG using jitter

### Golfing Challenges (Optimization)
Same problem, scored on efficiency.  
Metrics (cheat-resistant):
- Part count
- Total wire length
- Simulation settle time (ticks until stable)
- Hybrid score formula: (parts × w1) + (wire_length × w2) + (settle_ticks × w3)

**Important Decision**:  
Do NOT score real-world wall-clock solving time. It discourages learning, enables easy cheating on replays, and adds pressure. Use only in-simulation metrics.

## 2. Win Conditions

Win conditions are defined per level and checked by the gamify layer after each simulation step or on stable state.

### Supported Win Condition Types
- netStableHigh / netStableLow: A specific pin/net must stay high/low for N ticks
- partCountUnder: Total relays/wires ≤ max
- settleTimeUnder: Circuit stabilizes within X simulation ticks
- opponentHealthZero: Arena mode — reduce enemy health to 0
- opponentPushedOut: Arena mode — push enemy out of the ring
- surviveLongest: Arena mode — outlast opponent until time expires
- correctOutputSequence: Output matches expected pattern (e.g. counter values)
- displayShowsValue: 7-segment or dot-matrix shows correct number/text
- melodyPlayed: Buzzer plays correct sequence of tones (for music challenges)

Multiple win conditions can be combined (e.g. correctness + efficiency bonus).

**Implementation note**:  
Create a small gamify/winchecker.go that evaluates conditions against current world state and simulation results. This keeps the core sim/ package clean.

## 3. Input / Output Peripherals (to make it fun)

Essential for moving beyond switch + bulb:

| Priority | Peripheral                     | Why It Matters                              | Typical Use Cases                  | Recommended Sim Speed |
|----------|--------------------------------|---------------------------------------------|------------------------------------|-----------------------|
| 1        | 7-Segment Display (1–4 digits) | Instant visual feedback                     | Counters, scores, stopwatches      | 10–50×               |
| 2        | Piezo Buzzer / Speaker         | Audio payoff — tones, melodies, beeps       | Music box, alarms, Pong sounds     | 50–100×              |
| 3        | 8×8 or 5×7 Dot-Matrix          | Graphics and animations                     | Pong, Snake, scrolling text        | 20–100×              |
| 4        | Hex / Matrix Keyboard          | Proper multi-button input                   | Calculators, opcode entry          | Real-time            |
| 5        | Stepper Motor (visual)         | Mechanical movement                         | Plotter, clock hands               | 10–50×               |
| 6        | Simple VGA-style Output        | Memory-mapped screen                        | Full Pong, text terminal           | 50–100×              |
| 7        | Bell / Solenoid                | Satisfying "ding" reward                    | End-of-level, quiz buzzer          | Real-time            |

**Implementation note**: All are normal Part implementations. A global simulation speed multiplier (in world/) allows fun peripherals without forcing unrealistic 5–10 ms relay timing.

## 4. Arena / Competitive Modes (Relay Battles)

**Concept**: "Advanced BEAM Wars" — players build small hardcoded relay state machines that fight each other. No CPU or software required.

### Primary Idea: Relay Duel (Sumo / Tag style)
- Arena: 48×48 grid shown on a large dot-matrix display.
- Each player builds one "drone" circuit.
- Arena provides 4–6 sensor input pins (wall detection, enemy proximity, health) and 5 actuator output pins (move, turn, ram, shield).
- ArenaController part bridges drones to the battlefield and runs at 50–100× speed.
- Win conditions: push opponent out, reduce health to 0, or survive longest.

**Why it works well**:
- Simple to build (12–60 relays per drone).
- Highly visual and audible (clicking relays + dot-matrix action + buzzer).
- Naturally leads to golfing: "Win the arena with ≤ 25 relays".
- Excellent for sharing recorded battles.

**Simpler alternatives** (if full arena is too ambitious initially):
- Tug-of-War (shared flag net)
- Light-Cycle / Tron trails on dot-matrix
- Pure Sumo push-out

**Architecture fit**:
- New part/catalog/arena/ with ArenaController part.
- Temporary combined world.Parts list for the battle.
- Fully reuses existing sim loop and net resolution.

## 5. Persistence & Social Features

- Personal Save/Load: Desktop files + WASM IndexedDB
- Export/Import: JSON circuits for easy sharing
- Shared Repository: Optional upload for community circuits and leaderboards
- Screen Recording:
  - Desktop: GIF export
  - WASM: Canvas captureStream() + MediaRecorder (WebM)
  - Perfect for sharing arena battles and golf solutions

These enable leaderboards per challenge and easy sharing of winning designs.

## 6. Level Format & Level Builder

**Recommended**: Build a lightweight level builder early (reuses the existing editor heavily).

The level format should include:
- Level ID and title
- Mode (nand2tetris, design, golf, or arena)
- Description
- Starting schematic
- Allowed tools
- Win conditions
- Golf scoring rules (if applicable)
- Recommended simulation speed
- Next level reference

**Level JSON skeleton**:
```json
{
  "id": "nand2tetris-04-latch",
  "title": "Build a Set-Reset Latch",
  "mode": "nand2tetris",        // or "design", "golf", "arena"
  "description": "...",
  "startingParts": [ ... ],
  "allowedTools": ["relay", "wire", ...],
  "winConditions": [
    { "type": "netStableHigh", "pin": 42, "forTicks": 500 },
    { "type": "partCountUnder", "max": 12 }
  ],
  "golfScoring": { ... },
  "simSpeed": 50,
  "nextLevel": "..."
}```

**Plan**:
1. Ship first 8–10 levels hardcoded as JSON.
2. Add level builder (initially hidden/debug, then full UI).
3. Enable user-created and shared levels.

## 7. Implementation Priorities (Post Core v2)

1. Core peripherals: 7-Segment + Buzzer + Dot-Matrix
2. Personal save/load + WASM IndexedDB
3. Screen recording (especially in browser)
4. Level system + level builder
5. First arena mode (simple Relay Duel)
6. Shared circuits + basic leaderboards

---

**Key Principle**  
All gamification must stay true to the v2 design philosophy:  
- Concrete structs and methods  
- Parts own their state and behavior  
- Clear package structure  
- No magic, no heavy abstractions  

Gamification lives on top of the clean core — never inside it.

---

*Last updated: April 2026*  
*Living document — update as new ideas emerge.*