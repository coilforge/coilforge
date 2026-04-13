package sim

type unionFind struct {
	parent []int
}

func newUnionFind(size int) *unionFind {
	parent := make([]int, size)
	for i := range parent {
		parent[i] = i
	}
	return &unionFind{parent: parent}
}

func (uf *unionFind) Find(a int) int {
	if a < 0 || a >= len(uf.parent) {
		return a
	}
	if uf.parent[a] != a {
		uf.parent[a] = uf.Find(uf.parent[a])
	}
	return uf.parent[a]
}

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
