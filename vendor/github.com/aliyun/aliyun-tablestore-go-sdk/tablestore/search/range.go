package search

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"
	"math"
)

var Inf = math.Inf(1)
var NegInf = math.Inf(-1)

type Range struct {
	from float64	//math.Inf(-1) means negative infinity
	to float64		//math.Inf(1) means positive infinity
}

func BuildPBForRanges(ranges []Range) ([]*otsprotocol.Range, error) {
	if len(ranges) == 0 {
		return nil, errors.New("no range given")
	}
	var pbRanges []*otsprotocol.Range
	for i := 0; i < len(ranges); i++ {
		pbR := new(otsprotocol.Range)
		if !math.IsInf(ranges[i].from, -1) {
			pbR.From = &ranges[i].from
		}
		if !math.IsInf(ranges[i].to, 1) {
			pbR.To = &ranges[i].to
		}
		pbRanges = append(pbRanges, pbR)
	}
	return pbRanges, nil
}
