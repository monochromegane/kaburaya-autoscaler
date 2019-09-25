package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	autoscaler "github.com/monochromegane/kaburaya-autoscaler"
	simulator "github.com/monochromegane/queuing-theory-simulator"
)

var (
	seed   int
	step   int
	DT     float64
	rho    float64
	lambda string
	mu     string
	delay  float64

	outDir     string
	params     string
	simulation string
)

func init() {
	flag.IntVar(&seed, "seed", 1, "Seed")
	flag.IntVar(&step, "step", 1, "Number of step")
	flag.Float64Var(&rho, "rho", 0.7, "Utilization of the server (< 1.0)")
	flag.Float64Var(&DT, "DT", 1.0, "DT")
	flag.StringVar(&lambda, "lambda", "1.0", "Lambda (Number of arrival in UnitTime). Allow CSV style like 1.0,5,2.0,10,3.0")
	flag.StringVar(&mu, "mu", "1.0", "Mu (Number of service in UnitTime). Allow CSV style like 1.0,5,2.0,10,3.0")
	flag.Float64Var(&delay, "delay", 0.0, "The delay step until changing the number of servers")

	flag.StringVar(&outDir, "dir", "out", "Output directory")
	flag.StringVar(&params, "params", "params.csv", "File name for parameters")
	flag.StringVar(&simulation, "simulation", "simulation.csv", "File name for simulation result")
}

func main() {
	flag.Parse()

	fParams, fSimulation, err := setup()
	if err != nil {
		panic(err)
	}
	defer func() {
		fParams.Close()
		fSimulation.Close()
	}()

	lambdas, err := toParams(lambda)
	if err != nil {
		panic(err)
	}

	mus, err := toParams(mu)
	if err != nil {
		panic(err)
	}

	controller := autoscaler.NewKaburayaController(rho)
	plant := NewPlant(int64(seed), DT, lambdas, mus, delay)

	lambda_, mu_, ts_, waiting, server := 0.0, 0.0, 0.0, 0, 0.0
	for i := 0; i < step; i++ {
		s := controller.Calculate(lambda_, mu_, ts_)
		lambda_, mu_, ts_, waiting, server = plant.Run(int(s))
		fmt.Printf("Server: %.1f [%.1f], Lambda: %.1f, Mu: %.1f, Ts: %.5f, Waiting: %d\n", s, server, lambda_, mu_, ts_, waiting)
		fmt.Fprintf(fSimulation, "%f,%f,%d,%f,%f,%f\n", s, server, waiting, ts_/DT, lambda_, mu_)
	}
}

type Plant struct {
	DT    float64
	model simulator.Model
	delay autoscaler.Component
}

func NewPlant(seed int64, DT float64, lambda, mu func(int) float64, delay float64) *Plant {
	return &Plant{
		DT:    DT,
		model: simulator.NewMMSModel(seed, ToHighResolution(DT, lambda), ToHighResolution(DT, mu)),
		delay: &autoscaler.Delay{Gamma: int(delay / DT)},
	}
}

func (p *Plant) Run(s int) (float64, float64, float64, int, float64) {
	step := int(1.0 / p.DT)
	responseTimes := []int{}
	lambda := 0
	mu := 0
	waiting := 0
	avgServer := 0.0
	for i := 0; i < step; i++ {
		server := int(p.delay.Work(float64(s)))
		arrival, _, waiting_, ts := p.model.Progress(server)
		lambda += arrival
		mu += len(ts)
		waiting = waiting_
		responseTimes = append(responseTimes, ts...)
		avgServer = onlineAvg(server, i, avgServer)
	}
	return float64(lambda), float64(mu) / avgServer, average(responseTimes) * p.DT, waiting, avgServer
}

func average(xs []int) float64 {
	avg := 0
	for _, x := range xs {
		avg += x
	}
	return float64(avg) / float64(len(xs))
}

func ToHighResolution(DT float64, params func(int) float64) func(int) float64 {
	fn := func(i int) float64 {
		di := int(float64(i) * DT)
		return params(di) * DT
	}
	return fn
}

func toParams(param string) (func(int) float64, error) {
	params := append([]string{"0"}, strings.Split(param, ",")...)
	fParams := make([]float64, len(params))
	for i, p := range params {
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return func(t int) float64 { return 0.0 }, err
		}
		fParams[i] = v
	}
	return func(t int) float64 {
		for i := 0; i < len(fParams)/2; i++ {
			min := int(fParams[i*2])
			if len(fParams)-1 > i*2+2 {
				max := int(fParams[i*2+2])
				if min <= t && t < max {
					return fParams[i*2+1]
				}
			} else {
				if min <= t {
					return fParams[i*2+1]
				}
			}
		}
		return 0.0
	}, nil
}

func setup() (*os.File, *os.File, error) {
	if outDir != "" {
		err := os.MkdirAll(outDir, 0755)
		if err != nil {
			return nil, nil, err
		}
	}

	fParams, err := os.OpenFile(filepath.Join(outDir, params), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, nil, err
	}
	fmt.Fprintf(fParams, "seed,step,rho,DT,delay\n")
	fmt.Fprintf(fParams, "%d,%d,%f,%f,%f\n", seed, step, rho, DT, delay)

	fSimulation, err := os.OpenFile(filepath.Join(outDir, simulation), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, nil, err
	}
	fmt.Fprintf(fSimulation, "servers,delayedServers,waiting,averageResponseTime,lambda,mu\n")

	return fParams, fSimulation, nil
}

func onlineAvg(x, n int, avg float64) float64 {
	return (float64(n)*avg + float64(x)) / float64(n+1)
}
