package gimbal

import (
	"sync"
	"time"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
)

type Gimbal struct {
	cc *dji.CommandController
}

func New(cc *dji.CommandController) *Gimbal {
	wg := sync.WaitGroup{}
	wg.Add(1)
	cc.StartListening(dji.KeyGimbalConnection, func(result *dji.Result) {
		if result.Value().(bool) {
			// Enable chassis and gimbal updates.
			cc.PerformAction(dji.KeyRobomasterOpenChassisSpeedUpdates, nil, nil)
			cc.PerformAction(dji.KeyGimbalOpenAttitudeUpdates, nil, nil)
			wg.Done()
		}
	})

	wg.Wait()

	return &Gimbal{
		cc,
	}
}

func (g *Gimbal) ResetPosition() error {
	return g.cc.PerformAction(dji.KeyGimbalResetPosition, nil, nil)
}

func (g *Gimbal) MoveToAbsolutePosition(yawAngle, pitchAngle int,
	duration time.Duration) error {

	param := absoluteRotationParameter{
		int16(yawAngle * 10),
		int16(pitchAngle * 10),
		int16(duration.Milliseconds()),
	}

	if yawAngle != 0 {
		err := g.cc.PerformAction(dji.KeyGimbalAngleFrontYawRotation,
			param, nil)
		if err != nil {
			return err
		}
	}

	if pitchAngle != 0 {
		err := g.cc.PerformAction(dji.KeyGimbalAngleFrontPitchRotation,
			param, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
