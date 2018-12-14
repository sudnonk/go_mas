package utils

import "math"

func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func Round(n float64) int64 {
	return int64(math.Round(n))
}
