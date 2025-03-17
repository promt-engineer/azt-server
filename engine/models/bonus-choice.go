package models

type BonusChoice struct {
	Spins      int  `json:"spins"`
	Multiplier int  `json:"multiplier"`
	Random     bool `json:"random"`
}
