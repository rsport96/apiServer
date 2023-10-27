package main

func getMostProbableCountry(prbs []countryJson) string {
	var max float64
	for _, prb := range prbs {
		if prb.Probability > max {
			max = prb.Probability
		}
	}
	for _, prb := range prbs {
		if prb.Probability == max {
			return prb.Country
		}
	}
	return ""
}

func rmZeroes(bs []byte) []byte {
	var ans []byte
	for _, el := range bs {
		if el != 0 {
			ans = append(ans, el)
		}
	}
	return ans
}
