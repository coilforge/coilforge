Icon 1
  name:button_open.svg
  width:512 height:256
  endpoint#1   x:0 y:128
  endpoint#2   x:512 y:128
  Electrical toggle switch symbol in the open position. A horizontal wire runs from endpoint#1 to the left pivot at x:128, and another from the right contact area at x:343 out to endpoint#2, all at y:128. The left pivot is a filled black circle (r=12) at (128,128). At x:343 there is a short upward stem from y:128 to y:88 topped with a small chevron arrowhead (±10px wings at y:100). The switch arm is a line from the pivot (128,128) angled steeply upward to (350,46), clearly not touching the arrowhead — indicating open.

Icon 2
  name:button_closed.svg
  width:512 height:256
  endpoint#1   x:0 y:128
  endpoint#2   x:512 y:128
  Identical to button_open.svg (same wires, pivot dot, stem, and arrowhead) except the arm. The closed arm goes from the pivot (128,128) to (382,76), angling upward so its end rests on the arrowhead tip — indicating closed.

Icon 3
  name:indicator_off.svg
  width:512 height:256
  endpoint#1   x:0 y:128
  endpoint#2   x:512 y:128
  This should be a the electical symbol for an indicator light -  a circle with a x-style cross in it. The circle itelf should be 30% of the total width and having "wires/stems" going from the switch out to the respective endpoints. Drawn horizontally between the two endpoints. It should be in the unlit mode.

Icon 4
  name:indicator_on.svg
  width:512 height:256
  endpoint#1   x:0 y:128
  endpoint#2   x:512 y:128
  This should be identical to the indicator_off.svg, but the background porton inside of the circle should be orangey-yellow to indicate that the bulb is lit.

Icon 5
  name:power_gnd.svg
  width:512 height:128
  endpoint#1  x:256 y:0
  Electrical ground/earth symbol. A vertical stem runs from the endpoint at the top center down to ~30% of the height, where it meets the largest horizontal bar (~30% of the width). Two more bars follow below at even spacing, each progressively shorter (the middle bar ~20% of the width, the smallest ~9%). The bars are centered horizontally. The lowest bar should sit roughly 25% above the bottom edge.

Icon 6
  name:power_vcc.svg
  width:512 height:128
  endpoint#1  x:256 y:128
  VCC/positive power symbol. An upward-pointing equilateral triangle (outline only, no fill) centered horizontally, with its apex near the top (~15% from top) and base at roughly 70% of the height. The triangle width should be ~24% of the image width. A vertical stem runs from the base center of the triangle straight down to the endpoint at the bottom center. Same general proportions and spacing as power_gnd.svg but flipped — triangle at top, stem at bottom.

Icon 7
  name:rch.svg
  width:512 height:256
  endpoint#1  x:256 y:256
  RC delay / hold timer symbol. A vertical stem rises from the endpoint at the bottom center up ~20% of the height, ending at a horizontal capacitor bar (~27% of the image width, centered). Above the bar, with a visible gap (~12% of the height), sits an outline-only square of the same width and height as the bar. Inside the square is an upper-case "H", drawn with stroke-width 10 (same as all other strokes) and generous padding (~20% inset from the square edges) so it remains legible when scaled down.

Icon 7b
  name:rch-charged.svg
  width:512 height:256
  endpoint#1  x:256 y:256
  Identical to rch.svg but the square has a solid orangey-yellow (#FFBF30) fill, indicating the module is fully charged.

Icon 7c
  name:rch-draining.svg
  width:512 height:256
  endpoint#1  x:256 y:256
  Identical to rch.svg but only the bottom half of the square is filled with orangey-yellow (#FFBF30), using a clipPath at the vertical midpoint of the square. The top half remains transparent. This indicates the module is being used/drained.

Icon 10
  name:clock_high.svg
  width:512 height:256
  endpoint#1  x:512 y:128
  Clock / free-running pulse generator symbol showing output-high state. A rounded outline-only rectangle (~220×210, r=20, drawn as a path with arc corners matching the relay style) centered horizontally. A horizontal stem runs from the box's right edge to the endpoint at the right middle. Inside the box is an all-black square wave polyline: short-low, long-high, short-low — pulse height is ~40% of the box height, vertically centered, with round stroke-linejoin. A band of orangey-yellow (#FFBF30) fills the top of the box from the outline down, stopping ~22px above the high rail of the pulse, clipped to the rounded box shape.

Icon 10b
  name:clock_low.svg
  width:512 height:256
  endpoint#1  x:512 y:128
  Same box and stem as clock_high.svg. The waveform is inverted: short-high, long-low, short-high, with round stroke-linejoin. A band of orangey-yellow (#FFBF30) fills the bottom of the box from the outline up, stopping ~22px below the low rail of the pulse, clipped to the rounded box shape.

Icon 8
  name:diode.svg
  width:512 height:256
  endpoint#1   x:0 y:128
  endpoint#2   x:512 y:128
  Standard diode symbol drawn horizontally between the two endpoints. A right-pointing triangle (outline only, no fill) centered in the image, ~25% of the total width (slightly narrower than the switch). The triangle's left edge is vertical, apex points right. A vertical cathode bar the same height as the triangle sits flush against the apex. Horizontal wires connect from each endpoint to the triangle's left edge and from the cathode bar to the right endpoint, all at y=128.

Icon 9a
  name:relay_bottom.svg
  width:512 height:128
  endpoint#1   x:128 y:128  (25% of width)
  endpoint#2   x:384 y:128  (75% of width)
  Bottom part of a 3-piece stackable relay. The outline (stroke-width 5 — halved from the standard 10 because the relay renders at 2× scale in the app) covers the left, right, and bottom edges with rounded corners (r≈20) — the top edge is open for stacking. The outline is inset ~1% from the viewBox edges so strokes don't clip. Two horizontal wire segments at the top edge (y:0) provide full stroke thickness at the stacking seam: one from x:0 to x:128 (COM) and one from x:512 to x:350 (NC), both stroke-width 5. Short vertical stems rise from each endpoint at the bottom up to ~62% of the height, where a relay coil is drawn horizontally between them. The coil is four connected upward-bulging semicircle arcs (half-circles, tall — nearly as tall as the distance from the coil baseline to the top of the part). De-energised variant: coil is black stroke only, no fill.

Icon 9b
  name:relay_bottom_on.svg
  width:512 height:128
  endpoint#1   x:128 y:128  (25% of width)
  endpoint#2   x:384 y:128  (75% of width)
  Identical to relay_bottom.svg but the coil stroke is colored orangey-yellow (#FFBF30) instead of black, indicating the relay is energised. No fill on the coil — only the stroke changes color.

Icon 9c
  name:relay_middle_nc.svg
  width:512 height:256
  endpoint#1  x:0 y:256 (common, 100%)
  endpoint#2  x:512 y:256 (NC, 100%)
  endpoint#3  x:512 y:128 (NO, 50%)
  Middle/armature part of the stackable relay, non-energised (NC) position. Side outlines (stroke-width 5) on left and right edges, top and bottom open. All relay SVG stroke widths are halved (wires sw 5, arrowheads sw 3.75, pivot r=6) because the relay renders at 2× the scale of other parts in the app. Common wire runs horizontally from the left endpoint at y:256 to x:128, then a 30px vertical stem up to a filled pivot dot (r=6) at (128,226). NC contact: horizontal wire from the right endpoint at y:256 to x:350, then a 30px vertical stem up to y:226, with an upward-pointing chevron arrowhead (±10px wings, stroke-width 3.75). NO contact: horizontal wire from the right endpoint at y:128 to x:350, then a 30px vertical stem down to y:158, with a downward-pointing chevron arrowhead (±10px wings, stroke-width 3.75). The pivot L-shape and NC L-shape are mirrors of each other. Armature line from the pivot (128,226) horizontal to (382,226), resting firmly on the NC arrowhead tip — NC (de-energised) position.

Icon 9c-float
  name:relay_middle_float.svg
  width:512 height:256
  endpoint#1  x:0 y:256 (common, 100%)
  endpoint#2  x:512 y:256 (NC, 100%)
  endpoint#3  x:512 y:128 (NO, 50%)
  Identical to relay_middle_nc.svg except the armature tilts moderately upward from the pivot (128,226) to (382,192), floating between both contacts.

Icon 9c-no
  name:relay_middle_no.svg
  width:512 height:256
  endpoint#1  x:0 y:256 (common, 100%)
  endpoint#2  x:512 y:256 (NC, 100%)
  endpoint#3  x:512 y:128 (NO, 50%)
  Identical to relay_middle_nc.svg except the armature tilts steeply upward from the pivot (128,226) to (382,150), sinking slightly past the NO arrowhead tip at y:158 for a firm contact appearance — NO (energised) position.

Icon 9d
  name:relay_top.svg
  width:512 height:128
  Top cap of the stackable relay. The outline (stroke-width 5, halved for 2× render scale) is a horizontal bar with rounded corners (r≈20) sitting at the bottom of the viewBox (y:108→128 via arcs). No vertical side segments — the path starts and ends at the bottom edge (y:128) with arcs curving up to the bar. This minimises empty space above the relay when the armature sits lower in the middle part.
