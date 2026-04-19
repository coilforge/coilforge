package flatten

import (
	"coilforge/internal/core"
	"coilforge/internal/part"
	"coilforge/internal/world"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "coilforge/internal/part/catalog/gnd"
	_ "coilforge/internal/part/catalog/indicator"
	_ "coilforge/internal/part/catalog/vcc"
)

type fileFormat struct {
	NextPartID int           `json:"nextPartID"`
	NextPinID  core.PinID    `json:"nextPinID"`
	Parts      []part.Record `json:"parts"`
}

func loadWorldFromJSON(path string) error {
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

func anchorSummary() (keys []string, pinsPerKey map[string][]core.PinID) {
	var anchors []core.PinAnchor
	for _, p := range world.Parts {
		anchors = append(anchors, p.Anchors()...)
	}
	pinsPerKey = make(map[string][]core.PinID)
	for _, a := range anchors {
		k := pointKey(a.Pt)
		pinsPerKey[k] = append(pinsPerKey[k], a.PinID)
	}
	for k := range pinsPerKey {
		keys = append(keys, k)
	}
	return keys, pinsPerKey
}

func TestPinKeysCoilforgeVsOK(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Caller")
	}
	root := filepath.Join(filepath.Dir(thisFile), "..", "..")

	for _, name := range []string{"coilforge-ok.json", "coilforge.json"} {
		path := filepath.Join(root, name)
		world.Reset()
		if err := loadWorldFromJSON(path); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		BuildNets()

		keys, byKey := anchorSummary()
		t.Logf("%s: %d nets, PinNet entries=%d anchor keys=%d", name, len(world.Nets), len(world.PinNet), len(keys))
		for k, ids := range byKey {
			t.Logf("  key=%s pins=%v", k, ids)
		}
		if len(byKey) != 2 {
			t.Errorf("%s: expected 2 merged pin positions (GND net + VCC net), got %d", name, len(byKey))
		}
	}
}
