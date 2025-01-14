package main

import "math"

func distance(a []float64, b []float64) float64 {

	d := math.Sqrt(math.Pow((a[0])-float64(b[0]), 2) + math.Pow((a[1]-b[1]), 2))
	return d
}
