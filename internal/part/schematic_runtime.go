package part

// File overview:
// schematic_runtime resets simulation-derived fields that must not persist in saves
// and should start from defaults after load or when entering run mode.

// SchematicRuntimeClear is implemented by parts that carry runtime-only schematic state
// (e.g. indicator lit driven by nets). ClearSchematicRuntime restores edit-mode defaults.
type SchematicRuntimeClear interface {
	ClearSchematicRuntime()
}

// ClearAllSchematicRuntime calls [SchematicRuntimeClear.ClearSchematicRuntime] on each part when implemented.
func ClearAllSchematicRuntime(parts []Part) {
	for _, p := range parts {
		if c, ok := p.(SchematicRuntimeClear); ok {
			c.ClearSchematicRuntime()
		}
	}
}
