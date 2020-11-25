package game

import (
	"math/rand"
	"sort"
)

// TODO : Init state

// RoleState represents a possible campIndex repartition
// It contains all the possible roleIndex repartition within each campIndex
type RoleState struct {
	camp []filter

	roleStates  [][]filter
	deathStates []filter

	total int
}

func (s *RoleState) FindCamp(match simpleMatch) int {
	for i := range s.camp {
		if match(s.camp[i]) {
			return i
		}
	}
	return -1
}

func (s *RoleState) FindRole(match simpleMatch, index int) int {
	for i, f := range s.roleStates[index] {
		if match(f) {
			return i
		}
	}
	return -1
}

func (s *RoleState) MatchCamp(match multiMatch) bool {
	return match(s.camp)
}

func (s *RoleState) MatchAndCount(roleMatch multiMatch, deathMatch simpleMatch) int {
	count := 0
	for i := range s.roleStates {
		if roleMatch(s.roleStates[i]) && deathMatch(s.deathStates[i]) {
			count++
			continue
		}
	}
	return count
}

func (s *RoleState) MatchAndCountDeath(roleMatch multiMatch, deathMatch simpleMatch, playerFilter simpleMatch) (int, int) {
	alive := 0
	death := 0
	for i := range s.roleStates {
		if !roleMatch(s.roleStates[i]) || !deathMatch(s.deathStates[i]) {
			continue
		}
		if playerFilter(s.deathStates[i]) {
			death++
			continue
		}
		alive++
	}
	return alive, death
}

func (s *RoleState) MatchAndRemoveRole(roleMatch multiMatch, deathMatch simpleMatch) {
	i := 0
	r := 0
	n := len(s.roleStates)
	for i < n-r {
		if roleMatch(s.roleStates[i]) && deathMatch(s.deathStates[i]) {
			s.roleStates[i], s.roleStates[n-1-r] = s.roleStates[n-1-r], s.roleStates[i]
			s.deathStates[i], s.deathStates[n-1-r] = s.deathStates[n-1-r], s.deathStates[i]
			s.total--
			r++
			continue
		}
		i++
	}
	s.roleStates = s.roleStates[:n-r]
	s.deathStates = s.deathStates[:n-r]
}

func (s *RoleState) MatchAndKill(roleMatch multiMatch, deathMatch simpleMatch, player int) {
	for i := range s.roleStates {
		if !roleMatch(s.roleStates[i]) || !deathMatch(s.deathStates[i]) {
			continue
		}
		s.deathStates[i] = s.deathStates[i].addPlayer(player)
	}
}

func (s *RoleState) VoteAndKill(voteFilter map[int]filter, voteCount map[int]int, victims []int) {
	for i := range s.deathStates {
		max := 0
		victim := 0
		draw := false
		for j := 0; j < len(victims) && max <= voteCount[victims[j]]; j++ {
			v := victims[j]
			realVoteCount := voteCount[v] - s.deathStates[i].countMatch(voteFilter[v])
			if realVoteCount < max {
				continue
			}
			draw = realVoteCount == max
			max = realVoteCount
			victim = v
		}
		if draw {
			continue
		}
		s.deathStates[i] |= newFilter().addPlayer(victim)
	}
}

// State represents all the states for the current game
type State struct {
	states      []*RoleState
	globalDeath filter
	stateCount  int

	camps []Camp
	roles []Role

	campIndex map[Camp]int
	roleIndex map[Role]int

	initialPlayers int
	deadPlayers    map[int]struct{}

	randModule *rand.Rand
}

type PlayerState struct {
	id     int
	camp   Camp
	role   Role
	status Status
}

func (s *State) GetPlayerCampProb(player int, playerStates ...PlayerState) map[Camp]float64 {
	matchCamp := s.getMatchCamp(playerStates)
	matchPlayer := s.getMatchPlayer(player)
	matchRole := s.getMatchRole(playerStates)
	matchDeath := s.getMatchDeath(playerStates)

	campCount := make(map[int]int)
	total := 0
	for _, roleState := range s.states {
		if !roleState.MatchCamp(matchCamp) {
			continue
		}
		count := roleState.MatchAndCount(matchRole, matchDeath)
		if count == 0 {
			continue
		}
		camp := roleState.FindCamp(matchPlayer)
		total += count
		campCount[camp] += count
	}

	prob := make(map[Camp]float64)
	for camp, index := range s.campIndex {
		prob[camp] = float64(campCount[index]) / float64(total)
	}
	return prob
}

func (s *State) GetPlayerLiveProb(player int, playerStates ...PlayerState) map[Status]float64 {
	matchCamp := s.getMatchCamp(playerStates)
	matchPlayer := s.getMatchPlayer(player)
	matchRole := s.getMatchRole(playerStates)
	matchDeath := s.getMatchDeath(playerStates)

	aliveCount := 0
	deathCount := 0
	for _, roleState := range s.states {
		if !roleState.MatchCamp(matchCamp) {
			continue
		}
		alive, death := roleState.MatchAndCountDeath(matchRole, matchDeath, matchPlayer)
		aliveCount += alive
		deathCount += death
	}

	prob := make(map[Status]float64, 2)
	prob[Alive] = float64(aliveCount) / float64(aliveCount+deathCount)
	prob[Dead] = float64(deathCount) / float64(aliveCount+deathCount)

	return prob
}

func (s *State) RemoveIfMatch(playerStates ...PlayerState) {
	if playerStates == nil {
		return
	}

	matchCamp := s.getMatchCamp(playerStates)
	matchRole := s.getMatchRole(playerStates)
	matchDeath := s.getMatchDeath(playerStates)

	match := func(rs *RoleState) bool {
		if !rs.MatchCamp(matchCamp) {
			return false
		}
		rs.MatchAndRemoveRole(matchRole, matchDeath)
		return rs.total == 0
	}

	s.removeIfMatch(match)
}

func (s *State) KillIfMatch(player int, playerStates ...PlayerState) {
	matchCamp := s.getMatchCamp(playerStates)
	matchRole := s.getMatchRole(playerStates)
	matchDeath := s.getMatchDeath(playerStates)

	for _, rs := range s.states {
		if !rs.MatchCamp(matchCamp) {
			continue
		}
		rs.MatchAndKill(matchRole, matchDeath, player)
	}
}

func (s *State) CheckDeath() []int {
	var newDeaths []int
	for i := 0; i < s.initialPlayers; i++ {
		if _, ok := s.deadPlayers[i]; ok {
			continue
		}
		if s.globalDeath.matchAnd(newFilter().addPlayer(i)) {
			newDeaths = append(newDeaths, i)
			s.deadPlayers[i] = struct{}{}
		}
	}
	return newDeaths
}

func (s *State) CollapsePlayers(players ...int) []PlayerState {
	index := s.randModule.Intn(s.stateCount)
	roleState, baseIndex := s.findRoleStateByIndex(index)

	playerStates := make([]PlayerState, len(players))
	for i, player := range players {
		matchPlayer := s.getMatchPlayer(player)
		playerStates[i] = PlayerState{
			id:   player,
			camp: s.camps[roleState.FindCamp(matchPlayer)],
		}
		roleIndex := roleState.FindRole(matchPlayer, index-baseIndex)
		if roleIndex >= 0 {
			playerStates[i].role = s.roles[roleIndex]
		}
	}

	matchCamp := s.getMatchCamp(playerStates)
	matchRole := s.getMatchRole(playerStates)

	antiMatchRole := func(f []filter) bool { return !matchRole(f) }

	match := func(rs *RoleState) bool {
		if !rs.MatchCamp(matchCamp) {
			return true
		}
		rs.MatchAndRemoveRole(antiMatchRole, simpleMatchAll)
		return rs.total == 0
	}
	s.removeIfMatch(match)

	return playerStates
}

func (s *State) VoteAndKillPlayer(votes map[int]int) {
	voteFilter := make(map[int]filter)
	voteCount := make(map[int]int)

	for id, vote := range votes {
		voteFilter[vote] = voteFilter[vote].addPlayer(id)
		voteCount[vote]++
	}

	victims := make([]int, 0, len(voteCount))
	for i := range voteCount {
		victims = append(victims, i)
	}
	sort.Slice(victims, func(i, j int) bool {
		return voteCount[victims[i]] > voteCount[victims[j]]
	})

	for _, rs := range s.states {
		rs.VoteAndKill(voteFilter, voteCount, victims)
	}
}

func (s *State) getMatchCamp(states []PlayerState) multiMatch {
	if states == nil {
		return multiMatchAll
	}

	masks := make(map[int]filter)
	for _, playerState := range states {
		if playerState.camp.IsNil() {
			continue
		}
		i := s.campIndex[playerState.camp]
		mask, ok := masks[i]
		if !ok {
			mask = newFilter()
		}
		mask = mask.addPlayer(playerState.id)
		masks[i] = mask
	}
	if len(masks) == 0 {
		return multiMatchAll
	}

	return func(filters []filter) bool {
		for i, mask := range masks {
			if !filters[i].matchAnd(mask) {
				return false
			}
		}
		return true
	}
}

func (s *State) getMatchPlayer(player int) simpleMatch {
	mask := newFilter().addPlayer(player)
	return func(f filter) bool {
		return f.matchAnd(mask)
	}
}

func (s *State) getMatchRole(states []PlayerState) multiMatch {
	if states == nil {
		return multiMatchAll
	}

	masks := make(map[int]filter)
	for _, playerState := range states {
		if playerState.role.IsNil() {
			continue
		}
		i := s.roleIndex[playerState.role]
		mask, ok := masks[i]
		if !ok {
			mask = newFilter()
		}
		mask = mask.addPlayer(playerState.id)
		masks[i] = mask
	}

	if len(masks) == 0 {
		return multiMatchAll
	}

	return func(filters []filter) bool {
		for i, mask := range masks {
			if !filters[i].matchAnd(mask) {
				return false
			}
		}
		return true
	}
}

func (s *State) getMatchDeath(states []PlayerState) simpleMatch {
	if states == nil {
		return simpleMatchAll
	}

	aliveMask := newFilter()
	deathMask := newFilter()
	for _, playerState := range states {
		if playerState.status.IsNil() {
			continue
		}
		if playerState.status == Alive {
			aliveMask = aliveMask.addPlayer(playerState.id)
		}
		if playerState.status == Dead {
			deathMask = deathMask.addPlayer(playerState.id)
		}
	}

	if aliveMask.isNil() && deathMask.isNil() {
		return simpleMatchAll
	}

	return func(f filter) bool {
		return f.matchAnd(deathMask) && !f.matchOr(aliveMask)
	}
}

func (s *State) findRoleStateByIndex(index int) (*RoleState, int) {
	j := 0
	for _, rs := range s.states {
		if j+rs.total > index {
			return rs, j
		}
		j += rs.total
	}
	return nil, s.stateCount
}

func (s *State) removeIfMatch(match func(rs *RoleState) bool) {
	i := 0
	r := 0
	n := len(s.states)
	total := 0
	for i < n-r {
		m := match(s.states[i])
		total += s.states[i].total
		if m {
			s.states[i], s.states[n-1-r] = s.states[n-1-r], s.states[i]
			r++
			continue
		}
		i++
	}
	s.states = s.states[:n-r]
	s.stateCount = total
}
