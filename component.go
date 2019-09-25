package autoscaler

import (
	"container/list"
	"math"
)

type Component interface {
	Work(float64) float64
}

type RoundAndMinimum struct {
	Minimum float64
}

func (c *RoundAndMinimum) Work(x float64) float64 {
	return math.Max(math.Round(x), c.Minimum)
}

type Delay struct {
	Gamma int
	queue *list.List
}

func (c *Delay) Work(x float64) float64 {
	if c.queue == nil {
		c.queue = list.New()
	}
	c.queue.PushBack(x)
	if c.queue.Len() <= c.Gamma {
		return 0.0
	}
	elem := c.queue.Front()
	defer c.queue.Remove(elem)

	y, _ := elem.Value.(float64)
	return y
}
