package game

// State is the full Qwixx scorecard state sent to/from the client.
type State struct {
	Crosses     [4][12]bool `json:"crosses"`
	MissedRolls [4]bool     `json:"missedRolls"`
}

// ScoresByXCount maps number of crosses in a row to points (0..12).
var ScoresByXCount = [13]int{0, 1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 66, 78}

// CountCrossesInRow returns how many cells are crossed in a row.
func CountCrossesInRow(row [12]bool) int {
	n := 0
	for _, v := range row {
		if v {
			n++
		}
	}
	return n
}

// Score computes total score from crosses and missed rolls.
func (s *State) Score() int {
	score := 0
	for i := range 4 {
		score += ScoresByXCount[CountCrossesInRow(s.Crosses[i])]
	}
	missed := 0
	for _, v := range s.MissedRolls {
		if v {
			missed++
		}
	}
	return score - 5*missed
}

// UnlockedLock returns true if row has at least 5 crosses (lock is unlocked).
func (s *State) UnlockedLock(row int) bool {
	return CountCrossesInRow(s.Crosses[row]) >= 5
}

// DisabledIndex returns the rightmost crossed index in the row; cells at or left of it are disabled.
func (s *State) DisabledIndex(row int) int {
	idx := -1
	for i, v := range s.Crosses[row] {
		if v {
			idx = i
		}
	}
	return idx
}

// SetCross sets the cross at (row, col). No validation (client can send invalid; we just apply).
func (s *State) SetCross(row, col int) {
	if row >= 0 && row < 4 && col >= 0 && col < 12 {
		s.Crosses[row][col] = true
	}
}

// UnsetCross clears the cross at (row, col).
func (s *State) UnsetCross(row, col int) {
	if row >= 0 && row < 4 && col >= 0 && col < 12 {
		s.Crosses[row][col] = false
	}
}

// SetMissed sets the missed roll for color index i.
func (s *State) SetMissed(i int) {
	if i >= 0 && i < 4 {
		s.MissedRolls[i] = true
	}
}

// UnsetMissed clears the missed roll for color index i.
func (s *State) UnsetMissed(i int) {
	if i >= 0 && i < 4 {
		s.MissedRolls[i] = false
	}
}

// ClearAll resets the entire scorecard.
func (s *State) ClearAll() {
	for i := range 4 {
		for j := range 12 {
			s.Crosses[i][j] = false
		}
		s.MissedRolls[i] = false
	}
}
