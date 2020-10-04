package game

type Camp int

const (
	evilCamp Camp = iota + 1
	goodCamp
)

var (
	// AllCamp defines all the possible campIndex
	AllCamp = []Camp{
		evilCamp,
		goodCamp,
	}
)

func (c Camp) IsNil() bool {
	return c == 0
}

type Role int

const (
	dominantWolve Role = iota + 1
	seer
)

var (
	RoleToCamp = map[Role]Camp{
		dominantWolve: evilCamp,
		seer:          goodCamp,
	}
)

func (r Role) IsNil() bool {
	return r == 0
}

type Status int

const (
	Alive Status = iota + 1
	Dead
)

func (s Status) IsNil() bool {
	return s == 0
}
