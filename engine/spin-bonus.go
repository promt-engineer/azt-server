package engine

import "aztec-pyramids/engine/models"

type Bonus struct {
	Win int64 `json:"award"`

	Spins []SpinBonus `json:"spins"`
}

type SpinBonus struct {
	Stops               []int `json:"stops"`
	Award               int64 `json:"award"`
	BonusSpinsLeft      int   `json:"bonus_spins_left"`
	BonusSpinsTriggered int   `json:"bonus_spins_triggered"`
	Multiplier          []int `json:"multiplier"`

	Avalanches []models.Avalanche `json:"avalanches"`
}

func (s *SpinBonus) deepCopy() SpinBonus {
	spin := SpinBonus{
		Award:               s.Award,
		BonusSpinsLeft:      s.BonusSpinsLeft,
		BonusSpinsTriggered: s.BonusSpinsTriggered,

		Stops:      make([]int, len(s.Stops)),
		Multiplier: make([]int, len(s.Multiplier)),
		Avalanches: make([]models.Avalanche, len(s.Avalanches)),
	}

	copy(spin.Stops, s.Stops)
	copy(spin.Multiplier, s.Multiplier)
	copy(spin.Avalanches, s.Avalanches)

	return spin
}

func (b *Bonus) Award() int64 {
	if b == nil {
		return 0
	}

	return b.Win
}

func (b *Bonus) deepCopy() *Bonus {
	bonus := &Bonus{
		Win:   b.Win,
		Spins: make([]SpinBonus, len(b.Spins)),
	}

	for i, spin := range b.Spins {
		bonus.Spins[i] = spin.deepCopy()
	}

	return bonus
}
