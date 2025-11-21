# Boss Implementation Plan

This document outlines the plan for porting bosses from the Java codebase to the Go microservices architecture.

## Overview

We need to implement specific logic for various boss groups. Each group has unique behaviors, spawn conditions, and skill sets.

## Boss Groups

### 1. Androids (Tiểu Đội Sát Thủ) ✅

**Reference**: `src/boss/Android`
**Status**: Implemented in `boss_android.go`

- **Bosses**: Android 13, 14, 15, 19, 20, Poc, Pic, King Kong.
- **Behavior**:
  - Often appear in groups.
  - Can fuse (e.g., Android 13 + 14 + 15 → Super Android 13).
  - Self-destruct mechanics (Android 20).
  - Energy absorption (Android 20).

### 2. Black Goku & Zamasu ✅

**Reference**: `src/boss/Black_Goku`
**Status**: Implemented in `boss_black_goku.go`

- **Bosses**: Black Goku, Zamasu.
- **Behavior**:
  - Transformation (Super Saiyan Rose at 50% HP).
  - Healing/Immortality (Zamasu regenerates 5% HP every 5s).
  - Fusion (Zamasu + Black Goku → Fused Zamasu at 30% HP).

### 3. Broly ✅

**Reference**: `src/boss/Broly`
**Status**: Implemented in `boss_broly.go`

- **Bosses**: Broly, Super Broly.
- **Behavior**:
  - Increasing power as HP drops.
  - Rage mode at 50% HP (+50% damage).
  - Berserk mode at 25% HP (2x damage).

### 4. Cell (Xên Bọ Hung) ✅

**Reference**: `src/boss/Cell`
**Status**: Implemented in `boss_cell.go`

- **Bosses**: Cell Saga (Xên cấp 1, 2, Hoàn thiện, Xên con).
- **Behavior**:
  - Evolution through 3 forms (Imperfect → Semi-Perfect → Perfect).
  - Regeneration (25% HP once per minute).
  - Self-destruct ability (Perfect form only).

### 5. Frieza (Fide) ✅

**Reference**: `src/boss/Frieza`, `src/boss/Golden_fireza`
**Status**: Implemented in `boss_frieza.go`

- [x] **Skill System**: Support boss-specific skills.
- [x] **AI Scripting**: Create `BossAI` interface to handle unique behaviors.

### Phase 2: Boss Data Migration

- `Template *BossTemplate`: Static boss data
- State tracking fields for chat, combat, and lifecycle

## Next Steps

1. **Boss Data Migration**: Create JSON templates for all bosses
2. **Integration**: Connect boss AI to `BossManager` update loop
3. **Combat Integration**: Link boss attacks to `CombatService`
4. **Zone Integration**: Spawn bosses in appropriate zones
5. **Network Packets**: Send boss spawn/death/chat messages to clients
6. **Remaining Bosses**: Implement Nappa, Cold, and minor bosses
