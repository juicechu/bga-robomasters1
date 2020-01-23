package pid

type PIDController struct {
	pController Controller
	iController Controller
	dController Controller
}

func NewPIDController(kp, ki, kd float64) Controller {
	return &PIDController{
		NewPController(kp),
		NewIController(ki),
		NewDController(kd),
	}
}

func (p *PIDController) Output(currentError float64) float64 {
	return p.pController.Output(currentError) +
		p.iController.Output(currentError) +
		p.dController.Output(currentError)
}
