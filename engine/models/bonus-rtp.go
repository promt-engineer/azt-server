package models

type BonusRtp struct {
	Low  float64 `json:"low"`
	High float64 `json:"high"`
}

func (b BonusRtp) CalcLowProb(realAvgFreeWin float64) float64 {
	return (b.High - realAvgFreeWin) / (b.High - b.Low)
}
