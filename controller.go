package autoscaler

import (
	"math"
)

type KaburayaController struct {
	Rho      float64
	Mu       float64
	Lambda   float64
	cnt      uint
	actuator Component
}

func NewKaburayaController(rho float64) *KaburayaController {
	return &KaburayaController{
		Rho:      rho,
		actuator: &RoundAndMinimum{Minimum: 1.0},
	}
}

func (c *KaburayaController) Calculate(lambda_, mu_, ts_ float64) float64 {
	if math.IsNaN(mu_) || math.IsNaN(ts_) || mu_ == 0.0 || ts_ == 0.0 {
		return c.actuator.Work(0.0)
	}

	c.Lambda = lambda_
	c.Mu = onlineAvgFloat(math.Max(mu_, 1.0/ts_), c.cnt, c.Mu)
	c.cnt++
	s := c.Lambda / (c.Rho * c.Mu)
	return c.actuator.Work(s)
}

func onlineAvgFloat(x float64, n uint, avg float64) float64 {
	return (float64(n)*avg + x) / float64(n+1)
}
