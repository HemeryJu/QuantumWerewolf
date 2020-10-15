package game

type filter uint64

// newFilter creates a new filter for n player
func newFilter() filter {
	return 0
}

func (f filter) addPlayer(i int) filter {
	return f | (1 << i)
}

func (f filter) matchAnd(g filter) bool {
	return f&g == g
}

func (f filter) matchOr(g filter) bool {
	return f&g != 0
}

func (f filter) countMatch(g filter) int {
	return (f & g).countBits()
}

func (f filter) isNil() bool {
	return f == 0
}

func (f filter) countBits() int {
	res := 0
	g := f
	for g != 0 {
		if g&1 == 1 {
			res++
		}
		g = g >> 1
	}
	return res
}

type simpleMatch func(filter) bool

var simpleMatchAll simpleMatch = func(filter) bool { return true }

type multiMatch func([]filter) bool

var multiMatchAll multiMatch = func([]filter) bool { return true }
