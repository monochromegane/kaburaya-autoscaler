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
	iDelay   *Delay
	iWeights []float64
	oDelay   *Delay
	oWeights []float64
}

func NewKaburayaController(rho, iDelay, oDelay float64) *KaburayaController {
	id, iw := newDelayWithWeight(iDelay)
	od, ow := newDelayWithWeight(oDelay)
	return &KaburayaController{
		Rho:      rho,
		actuator: &RoundAndMinimum{Minimum: 1.0},
		iDelay:   id,
		iWeights: iw,
		oDelay:   od,
		oWeights: ow,
	}
}

func newDelayWithWeight(delay float64) (*Delay, []float64) {
	gamma := int(math.Ceil(delay))
	weights := make([]float64, gamma)
	for i := 0; i < gamma; i++ {
		weights[i] = 1.0
	}
	if f := math.Floor(delay); gamma != int(f) {
		weights[0] = delay - f
	}
	return NewDelay(gamma), weights
}

func (c *KaburayaController) Calculate(lambda_, mu_, ts_ float64) float64 {
	if math.IsNaN(mu_) || math.IsNaN(ts_) || mu_ == 0.0 || ts_ == 0.0 {
		return c.actuator.Work(0.0)
	}

	c.Lambda = lambda_
	c.Mu = onlineAvgFloat(math.Max(mu_, 1.0/ts_), c.cnt, c.Mu)
	c.cnt++
	s := (c.Lambda +
		c.predictDelayedLambda(c.iDelay.list(), c.iWeights, c.Lambda, c.Mu) +
		c.predictDelayedLambda(c.oDelay.list(), c.oWeights, c.Lambda, c.Mu)) / (c.Rho * c.Mu)
	s = c.actuator.Work(s)
	c.iDelay.Work(s)
	c.oDelay.Work(s)
	return s
}

func (c *KaburayaController) predictDelayedLambda(servers, weights []float64, lambda_, mu_ float64) float64 {
	delayedLambda := 0.0
	for i, s := range servers {
		delayedLambda += math.Max((lambda_-(mu_*s))*weights[i], 0.0)
	}
	return delayedLambda
}

func onlineAvgFloat(x float64, n uint, avg float64) float64 {
	return (float64(n)*avg + x) / float64(n+1)
}
