package util

import "strconv"

type FloatParser struct {
	Err error
}

func (p *FloatParser) ParseFloat(str string, bitSize int) float64 {
	if p.Err != nil {
		return 0
	}
	res, err := strconv.ParseFloat(str, bitSize)
	p.Err = err
	return res
}
