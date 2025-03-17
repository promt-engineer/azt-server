package engine

import (
	"encoding/json"
	"fmt"
)

type Features struct {
	AnteBet  bool `json:"ante_bet"`
	BuyBonus bool `json:"buy_bonus"`
}

func UnmarshalTo[T any](payload interface{}) (T, error) {
	var b T

	bytes, err := json.Marshal(payload)
	if err != nil {
		return b, err
	}

	if err := json.Unmarshal(bytes, &b); err != nil {
		return b, err
	}

	return b, nil
}

func getFeaturesFromParams(params interface{}) (isAnteBet, isBuyFreeSpin bool, err error) {
	data, err := UnmarshalTo[Features](params)
	if err != nil {
		return
	}

	if data.AnteBet {
		isAnteBet = true
	}

	if data.BuyBonus {
		isBuyFreeSpin = true
	}

	if isBuyFreeSpin && isAnteBet {
		return false, false, fmt.Errorf("can not buy free spin and use ante bet at the same time")
	}

	return
}
