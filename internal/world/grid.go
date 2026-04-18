package world

// File overview:
// grid defines schematic snap and grid spacing in world units for editor and rendering.
// Subsystem: world shared conventions.
// Flow position: referenced by editor snapping and render background grid.

// MinorGridWorld is the fine wire-routing step in world units.
const MinorGridWorld = 1.0

// MajorGridWorld is placement / pin pitch; majors occur every MajorGridWorld/MinorGridWorld minors (4×).
const MajorGridWorld = 4.0
