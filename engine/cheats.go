package engine

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"encoding/json"
)

type Cheats struct {
	Stops                  []int `json:"stops"`
	AdditionalTriggerCount int   `json:"additional_trigger_count"`
}

func (s *Cheats) Eval(payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, s)
}

func getCheatsFromCtx(ctx engine.Context) (*Cheats, error) {
	var cheats *Cheats

	if ctx.Cheats != nil {
		cheats = new(Cheats)

		if err := cheats.Eval(ctx.Cheats); err != nil {
			return nil, err
		}
	}

	return cheats, nil
}
