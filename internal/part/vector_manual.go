package part

import "coilforge/internal/core"

func drawManualVectorAsset(name string, ctx DrawContext, base core.BasePart) bool {
	_ = ctx
	_ = base
	// Manual vector names (wire-*, template) removed with catalog cleanup; regenerate if re-added.
	_ = name
	return false
}
