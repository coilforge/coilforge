package sim

// File overview:
// union provides disjoint-set helpers used when grouping electrically connected pins.
// Subsystem: simulation support.
// It is used by net derivation logic in sim/flatten without depending on higher layers.
// Flow position: internal algorithm utility beneath simulation orchestration.

type unionFind struct {
	parent []int // parent value.
}

// newUnionFind handles new union find.
func newUnionFind(size int) *unionFind {
	parent := make([]int, size)
	for i := range parent {
		parent[i] = i
	}
	return &unionFind{parent: parent}
}

// Find handles find.
func (uf *unionFind) Find(a int) int {
	if a < 0 || a >= len(uf.parent) {
		return a
	}
	if uf.parent[a] != a {
		uf.parent[a] = uf.Find(uf.parent[a])
	}
	return uf.parent[a]
}

// Union handles union.
func (uf *unionFind) Union(a, b int) {
	if a < 0 || b < 0 || a >= len(uf.parent) || b >= len(uf.parent) {
		return
	}
	ra := uf.Find(a)
	rb := uf.Find(b)
	if ra != rb {
		uf.parent[rb] = ra
	}
}
