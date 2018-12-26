package utils

import (
	"math"
	"math/rand"
)

func RandDecimal() float64 {
	return rand.Float64() / math.MaxFloat64
}

func RandNormDecimal() float64 {
	return rand.NormFloat64() / math.MaxFloat64
}

func RandExpDecimal() float64 {
	return rand.ExpFloat64() / math.MaxFloat64
}
