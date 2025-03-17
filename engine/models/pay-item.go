package models

type PayItem struct {
	Symbol              int     `json:"symbol"`
	Indexes             [][]int `json:"indexes"`
	Award               int64   `json:"award"`
	Multiplier          int     `json:"multiplier"`
	AwardWithMultiplier int64   `json:"award_with_multiplier"`
}
