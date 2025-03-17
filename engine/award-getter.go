package engine

import (
	"fmt"
)

type AwardGetter struct {
	wager int64
}

func NewAwardGetter(wager int64) (AwardGetter, error) {
	if wager < 0 {
		return AwardGetter{}, fmt.Errorf("negative wager: %v", wager)
	}

	return AwardGetter{wager: wager}, nil
}

func (a AwardGetter) GetAward(symbol, size int) int64 {
	return multipliers[symbol][size] * a.wager / multiplicationDivider
}
