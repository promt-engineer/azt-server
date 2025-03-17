package main

import (
	"flag"
	"fmt"

	"aztec-pyramids/engine"

	"bitbucket.org/play-workspace/base-slot-server/pkg/app"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/config"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	baseEngine "bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
)

func main() {
	tmpConfig := flag.String("config", "config.yml", "Path to the configuration file")
	flag.Parse()

	application, err := app.New(*tmpConfig, engine.Bootstrap)
	if err != nil {
		panic(err)
	}

	cfg := application.Ctn().Get(constants.ConfigName).(*config.Config)

	rand := application.Ctn().Get(constants.RNGName).(rng.Client)
	if cfg.EngineConfig.MockRNG {
		rand = application.Ctn().Get(constants.RNGMockName).(rng.Client)
	}

	if err := application.RunOrSimulate(keepGenerateWrap(rand)); err != nil {
		panic(err)
	}

	for k, v := range engine.WinsMap.Map() {
		fmt.Printf("Scatter count: %d, count: %d\n", k, v)
	}
}

func keepGenerateWrap(rand rng.Client) services.KeepGenerateWrapper {
	return func(ctx baseEngine.Context, spin baseEngine.Spin, factory baseEngine.SpinFactory) (baseEngine.Spin, error) {
		sp, ok := spin.(*engine.SpinBase)
		if !ok {
			return nil, fmt.Errorf("unexpected spin type")
		}

		if len(sp.BonusChoice) > 0 {
			ctx.LastSpin = spin

			index, err := rand.Rand(4)
			if err != nil {
				return nil, err
			}

			newSpin, ok, err := factory.KeepGenerate(ctx, sp.BonusChoice[index])
			if err != nil {
				return nil, err
			}

			if ok {
				spin = newSpin
			}
		}

		return spin, nil
	}
}
