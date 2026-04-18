# Button Icons

Style rules for toolbar button icons:
- 192x192 pixels (downsampled to 48x48 for iPad toolbar buttons)
- ViewBox="0 0 192 192"
- Pure black line drawings on transparent background
- Stroke width: 8-12 (slightly thinner than part icons to keep clarity at 48px)
- Bold, simplified versions of the part symbols — optimised for small-size legibility
- Fill the full 192x192 area; external margins are handled by the toolbar

Button 1
  name: btn_button.svg
  Toggle switch symbol (open position). Left wire from x=8 to pivot dot at x=52, right wire from x=148 to x=184, all at y=112. Pivot is a filled circle (r=10). Right contact has an upward stem with a small chevron arrowhead. Arm angles steeply upward from pivot to ~(142,48), clearly not touching the arrowhead.

Button 2
  name: btn_clock.svg
  Clock/pulse generator. A rounded rectangle (path with r=16 arc corners) nearly filling the viewBox with ~10px margin. Inside: a bold square wave polyline — short-low, long-high, short-low — vertically spanning most of the box height.

Button 3
  name: btn_diode.svg
  Standard diode symbol centred in the frame. Right-pointing triangle (outline, no fill) from x=48 to x=144, spanning y=48 to y=144. Vertical cathode bar at x=144. Short horizontal wires connect to left and right edges at y=96.

Button 4
  name: btn_indicator.svg
  Indicator light symbol. A circle (r=70) centred at (96,96) with an X-cross inside. No wires — the circle fills most of the frame.

Button 5
  name: btn_gnd.svg
  Ground/earth symbol centred in frame. Vertical stem from top (y=16) down to the largest bar at y=64. Three horizontal bars below, evenly spaced (~36px apart), progressively shorter: ~112px, ~80px, ~40px. Bars use stroke-width=12 for boldness at small sizes.

Button 6
  name: btn_vcc.svg
  VCC/positive power symbol. Upward-pointing triangle (outline, no fill) with apex at (96,32) and base at y=128, ~112px wide. Vertical stem from base centre down to y=176.

Button 7
  name: btn_rch.svg
  RC hold/timer symbol. Vertical stem from bottom (y=184) up to a horizontal capacitor bar at y=148 (~100px wide, centred). Above with a gap: a square (100x100, rx=4) from y=16. Inside the square: a bold uppercase "H" with generous padding.

Button 8
  name: btn_relay.svg
  Relay symbol. Two vertical stems rising from the bottom (y=184) at x=32 and x=160 up to a coil at y=108. The coil is four upward-bulging semicircle arcs connecting the stems. Above the coil: a horizontal armature bar at y=48 spanning nearly the full width, with two short upward contact posts at x=48 and x=144.
