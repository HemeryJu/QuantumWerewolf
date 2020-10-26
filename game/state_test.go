package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func getTestStates() *State {
	// State 1 :
	//    Evil : 3, 4
	//    Good : 0, 1, 2
	//    All roles possible
	//    Everyone alive
	// State 2 :
	//    Evil : 2, 3
	//    Good : 0, 1, 4
	//    Only 1 Seer, 3 dead
	// State 3 :
	//    Evil : 1, 4
	//    Good : 0, 2, 3
	//    0, 2 Seer, 2 dead
	// State 4 :
	//    Evil : 2, 4
	//    Good : 0, 1, 3
	//    0, 1 Seer, 3 dead
	return &State{
		states: []*RoleState{
			{
				camp:        []filter{24, 7},
				roleStates:  [][]filter{{8, 1}, {8, 2}, {8, 4}, {16, 1}, {16, 2}, {16, 4}},
				deathStates: []filter{0, 0, 0, 0, 0, 0},
				total:       6,
			}, {
				camp:        []filter{12, 19},
				roleStates:  [][]filter{{4, 2}, {16, 2}},
				deathStates: []filter{8, 8},
				total:       2,
			}, {
				camp:        []filter{18, 13},
				roleStates:  [][]filter{{2, 1}, {2, 4}, {16, 1}, {16, 4}},
				deathStates: []filter{4, 4, 4, 4},
				total:       4,
			}, {
				camp:        []filter{20, 11},
				roleStates:  [][]filter{{4, 1}, {4, 2}, {16, 1}, {16, 2}},
				deathStates: []filter{8, 8, 8, 8},
				total:       4,
			},
		},
		stateCount: 16,
		campIndex: map[Camp]int{
			evilCamp: 0,
			goodCamp: 1,
		},
		roleIndex: map[Role]int{
			dominantWolve: 0,
			seer:          1,
		},
	}
}

func TestState_GetPlayerCampProb(t *testing.T) {

	stateTest := getTestStates()

	type expectedCase struct {
		good float64
		evil float64
	}

	type testCase struct {
		player       int
		playerStates []PlayerState
		expectCases  expectedCase
	}

	cases := []testCase{
		{
			player:       0,
			playerStates: nil,
			expectCases: expectedCase{
				good: 1.0,
				evil: 0.0,
			},
		}, {
			player:       1,
			playerStates: nil,
			expectCases: expectedCase{
				good: 0.75,
				evil: 0.25,
			},
		}, {
			player:       2,
			playerStates: nil,
			expectCases: expectedCase{
				good: 0.625,
				evil: 0.375,
			},
		}, {
			player:       3,
			playerStates: nil,
			expectCases: expectedCase{
				good: 0.5,
				evil: 0.5,
			},
		}, {
			player:       4,
			playerStates: nil,
			expectCases: expectedCase{
				good: 0.125,
				evil: 0.875,
			},
		}, {
			player: 3,
			playerStates: []PlayerState{
				{
					id:     0,
					role:   seer,
					status: Alive,
				}, {
					id:     3,
					status: Alive,
				},
			},
			expectCases: expectedCase{
				good: 0.5,
				evil: 0.5,
			},
		}, {
			player: 3,
			playerStates: []PlayerState{
				{
					id:     1,
					role:   seer,
					status: Alive,
				}, {
					id:     3,
					status: Alive,
				},
			},
			expectCases: expectedCase{
				good: 0.0,
				evil: 1.0,
			},
		}, {
			player: 3,
			playerStates: []PlayerState{
				{
					id:     2,
					role:   seer,
					status: Alive,
				}, {
					id:     3,
					status: Alive,
				},
			},
			expectCases: expectedCase{
				good: 0.0,
				evil: 1.0,
			},
		},
	}

	for i, tc := range cases {
		probs := stateTest.GetPlayerCampProb(tc.player, tc.playerStates...)
		require.InDelta(t, tc.expectCases.good, probs[goodCamp], 0.01, "error while testing good, test %d", i)
		require.InDelta(t, tc.expectCases.evil, probs[evilCamp], 0.01, "error while testing evil, test %d", i)
	}
}

func TestState_RemoveIfMatch(t *testing.T) {

	type testCase struct {
		playerStates []PlayerState
		expectCount  int
	}

	cases := []testCase{
		{
			playerStates: []PlayerState{
				{
					id:     0,
					role:   seer,
					status: Alive,
				}, {
					id: 3,
				},
			},
			expectCount: 10,
		}, {
			playerStates: []PlayerState{
				{
					id:     1,
					role:   seer,
					status: Alive,
				}, {
					id:     3,
					status: Dead,
				},
			},
			expectCount: 12,
		}, {
			playerStates: []PlayerState{
				{
					id:     2,
					role:   seer,
					status: Alive,
				}, {
					id:     3,
					status: Alive,
				},
			},
			expectCount: 14,
		},
	}

	for i, tc := range cases {
		stateTest := getTestStates()
		stateTest.RemoveIfMatch(tc.playerStates...)
		require.Equal(t, tc.expectCount, stateTest.stateCount, "error while testing, test %d", i)
	}
}

func TestState_KillIfMatch(t *testing.T) {

	type testCase struct {
		player        int
		playerStates  []PlayerState
		expectedDeath float64
	}

	testCases := []testCase{
		{
			player:        3,
			playerStates:  nil,
			expectedDeath: 1.0,
		},
	}

	for i, tc := range testCases {
		stateTest := getTestStates()
		stateTest.KillIfMatch(tc.player, tc.playerStates...)
		prob := stateTest.GetPlayerLiveProb(tc.player)
		for _, x := range stateTest.states {
			fmt.Println(x.deathStates)
		}
		require.InDelta(t, tc.expectedDeath, prob[Dead], 0.01, "error while testing, test %d", i)
	}
}
