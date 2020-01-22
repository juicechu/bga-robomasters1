package gimbal

import (
	"fmt"
	"sync"
	"time"

	"git.bug-br.org.br/bga/robomasters1/app/internal"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
)

type Gimbal struct {
	cc *internal.CommandController
}

func New(cc *internal.CommandController) *Gimbal {
	cc.StartListening(dji.KeyGimbalConnection,
		func(result *dji.Result, wg *sync.WaitGroup) {
			if result.Value().(bool) {
				// Enable chassis and gimbal updates.
				fmt.Println("Gimbal connection established.")
				cc.PerformAction(
					dji.KeyRobomasterOpenChassisSpeedUpdates, nil,
					nil)
				cc.PerformAction(dji.KeyGimbalOpenAttitudeUpdates, nil,
					nil)
			}

			wg.Done()
		})

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
		Time: int16(duration.Milliseconds()),
	}

	if yawAngle != 0 {
		param.Pitch = 0
		param.Yaw = int16(yawAngle * 10)
		err := g.cc.PerformAction(dji.KeyGimbalAngleFrontYawRotation,
			param, nil)
		if err != nil {
			return err
		}
	}

	if pitchAngle != 0 {
		param.Pitch = int16(pitchAngle * 10)
		param.Yaw = 0
		err := g.cc.PerformAction(dji.KeyGimbalAngleFrontPitchRotation,
			param, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
