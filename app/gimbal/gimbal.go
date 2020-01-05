package gimbal

import (
	"fmt"
	"time"

	"git.bug-br.org.br/bga/robomasters1/app/internal"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
)

type Gimbal struct {
	cc *internal.CommandController
}

func New(cc *internal.CommandController) *Gimbal {
	cc.StartListening(dji.KeyGimbalConnection, func(result *dji.Result) {
		if result.Value().(bool) {
			// Enable chassis and gimbal updates.
			fmt.Println("Gimbal connecvtion stablished.")
			cc.PerformAction(dji.KeyRobomasterOpenChassisSpeedUpdates, nil, nil)
			cc.PerformAction(dji.KeyGimbalOpenAttitudeUpdates, nil, nil)
		}
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
