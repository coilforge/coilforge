package sim

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/world"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"coilforge/internal/part/catalog/indicator"

	_ "coilforge/internal/part/catalog/gnd"
	_ "coilforge/internal/part/catalog/indicator"
	_ "coilforge/internal/part/catalog/vcc"
)

type fileFormat struct {
	NextPartID int           `json:"nextPartID"`
	NextPinID  core.PinID    `json:"nextPinID"`
	Parts      []part.Record `json:"parts"`
}

func loadWorld(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var file fileFormat
	if err := json.Unmarshal(data, &file); err != nil {
		return err
	}
	parts := make([]part.Part, 0, len(file.Parts))
	for _, record := range file.Parts {
		p, err := part.DecodeRecord(record)
		if err != nil {
			return err
		}
		parts = append(parts, p)
	}
	world.Parts = parts
	world.NextPartID = file.NextPartID
	world.NextPinID = file.NextPinID
	return nil
}

func TestNetStatesCoilforgeOKAnd45(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Caller")
	}
	root := filepath.Join(filepath.Dir(thisFile), "..", "..")

	for _, name := range []string{"coilforge-ok.json", "coilforge.json"} {
		path := filepath.Join(root, name)
		world.Reset()
		if err := loadWorld(path); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		Start()

		for nid, st := range world.NetStates {
			t.Logf("%s NetStates[%d]=%d", name, nid, st)
		}
		ind := world.Parts[0].(*indicator.Indicator)
		t.Logf("%s indicator Lit=%v", name, ind.Lit)

		if world.NetStates == nil {
			t.Fatalf("%s: NetStates nil", name)
		}
		haveH := false
		haveL := false
		for _, st := range world.NetStates {
			if st == core.NetHigh {
				haveH = true
			}
			if st == core.NetLow {
				haveL = true
			}
		}
		if !haveH || !haveL {
			t.Errorf("%s: expected both NetHigh and NetLow among nets", name)
		}
		if !ind.Lit {
			t.Errorf("%s: expected indicator lit when straddling VCC and GND nets", name)
		}
		Stop()
	}
}
