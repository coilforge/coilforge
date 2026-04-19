package flatten

// Disjoint-set for grouping quantized schematic nodes into electrical nets.

type uf struct {
	parent []int
}

func newUnionFind(n int) *uf {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &uf{parent: p}
}

func (u *uf) find(a int) int {
	if a < 0 || a >= len(u.parent) {
		return a
	}
	if u.parent[a] != a {
		u.parent[a] = u.find(u.parent[a])
	}
	return u.parent[a]
}

func (u *uf) union(a, b int) {
	if a < 0 || b < 0 || a >= len(u.parent) || b >= len(u.parent) {
		return
	}
	ra := u.find(a)
	rb := u.find(b)
	if ra != rb {
		u.parent[rb] = ra
	}
}
