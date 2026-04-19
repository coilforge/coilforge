package uidebug

// File overview:
// uidebug writes optional per-frame snapshots for diagnosing UI vs simulation timing.
// Enable with COILFORGE_UI_DEBUG=1 (writes coilforge-ui-debug.log in cwd) or set to a file path.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"coilforge/internal/part/catalog/clock"
	"coilforge/internal/part/catalog/indicator"
	"coilforge/internal/world"
)

var (
	logMu     sync.Mutex
	logFile   *os.File
	logOpened bool

	prevWall  time.Time
	prevSimUs uint64
	prevHas   bool
)

// logPath returns the destination path when logging is enabled, or "" when disabled.
func logPath() string {
	v := strings.TrimSpace(os.Getenv("COILFORGE_UI_DEBUG"))
	if v == "" {
		return ""
	}
	if v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes") {
		return "coilforge-ui-debug.log"
	}
	return v
}

// LogUpdateFrame appends one line per Ebiten Update while in run mode (when COILFORGE_UI_DEBUG is set).
// It holds SimMu briefly to sample sim time and catalog runtime fields consistently with the sim loop.
func LogUpdateFrame() {
	path := logPath()
	if path == "" {
		return
	}
	if !world.RunMode {
		logMu.Lock()
		prevHas = false
		logMu.Unlock()
		return
	}

	now := time.Now()

	world.SimMu.RLock()
	simUs := world.SimTimeMicros
	phase := clock.PhaseLabel(simUs)

	var clkParts []string
	var indIDs []int
	indLit := map[int]bool{}
	indLabel := map[int]string{}

	for _, p := range world.Parts {
		if c, ok := p.(*clock.Clock); ok {
			b := c.Base()
			clkParts = append(clkParts, fmt.Sprintf("id=%d label=%q", b.ID, b.Label))
		}
		if in, ok := p.(*indicator.Indicator); ok {
			b := in.Base()
			indIDs = append(indIDs, b.ID)
			indLit[b.ID] = in.Lit
			indLabel[b.ID] = b.Label
		}
	}
	world.SimMu.RUnlock()

	half := clock.HalfPeriodMicros()
	var halfIdx uint64
	phaseMod2 := 0
	if half > 0 {
		halfIdx = simUs / half
		phaseMod2 = int(halfIdx % 2)
	}

	clkJoined := strings.Join(clkParts, "; ")
	if clkJoined == "" {
		clkJoined = "(none)"
	}
	sort.Ints(indIDs)

	var indParts []string
	for _, id := range indIDs {
		indParts = append(indParts, fmt.Sprintf("id=%d label=%q lit=%t", id, indLabel[id], indLit[id]))
	}
	indJoined := strings.Join(indParts, "; ")
	if indJoined == "" {
		indJoined = "(none)"
	}

	logMu.Lock()
	defer logMu.Unlock()

	dtStr := "NA"
	dSimStr := "NA"
	if prevHas {
		dtStr = fmt.Sprintf("%.3f", now.Sub(prevWall).Seconds()*1000)
		dSimStr = fmt.Sprintf("%d", int64(simUs)-int64(prevSimUs))
	}

	line := fmt.Sprintf(
		"%s dt_ms=%s d_sim_us=%s sim_us=%d half_us=%d half_idx=%d phase_mod2=%d clock_phase=%s clocks %s indicators %s\n",
		now.Format(time.RFC3339Nano),
		dtStr,
		dSimStr,
		simUs,
		half,
		halfIdx,
		phaseMod2,
		phase,
		clkJoined,
		indJoined,
	)

	if err := appendLogLine(path, line); err != nil {
		return
	}

	prevWall = now
	prevSimUs = simUs
	prevHas = true
}

func appendLogLine(path, line string) error {
	if !logOpened {
		path = filepath.Clean(path)
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		logFile = f
		logOpened = true
	}
	_, err := logFile.WriteString(line)
	return err
}
