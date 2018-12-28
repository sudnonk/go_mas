package utils

import (
	"math/rand"
)

//-4～4σを0～1の小数にする
func RandNormDecimal(ra *rand.Rand) float64 {
	//NormFloat64はμ=0、σ=1
	t := ra.NormFloat64()/4 + 0.5
	if t < 0 || t > 1 {
		return RandNormDecimal(ra)
	} else {
		return t
	}
}

func RandIntSlice(max int64, num int64, e int64, ra *rand.Rand) []int64 {
	m := ra.Int63n(num)
	r := make(map[int64]struct{}, m)
	r[e] = struct{}{}

	for i := int64(0); i < m; i++ {
		r[ra.Int63n(max)] = struct{}{}
	}

	keys := make([]int64, len(r))
	i := 0
	for k := range r {
		keys[i] = k
		i++
	}

	return keys
}
