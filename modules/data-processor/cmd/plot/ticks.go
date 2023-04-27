package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
)

type MyTicker struct {
	width      float64
	labelwidth float64
	labelfmt   string
}

func (m MyTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	start := float64(int(min/10) * 10)
	mulval := float64(1000)
	labelmul := m.labelwidth * mulval
	for i := start; i < max+m.width; i += m.width {
		var label string
		muli := math.Round(i * mulval)
		showLabel := int(muli)%int(labelmul) == 0
		if showLabel {
			label = fmt.Sprintf(m.labelfmt, i)
		} else {
			label = ""
		}

		ticks = append(ticks, plot.Tick{Value: float64(i), Label: label})
	}
	return ticks
}
