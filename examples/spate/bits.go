package main

import (
	"errors"
	"strings"

	"github.com/alecthomas/units"
)

var (
	bytesUnitMap       = units.MakeUnitMap("iB", "B", 1024)
	metricBytesUnitMap = units.MakeUnitMap("B", "B", 1000)
	bitsUnitMap        = units.MakeUnitMap("ib", "b", 1024)
	metricBitsUnitMap  = units.MakeUnitMap("b", "b", 1000)
)

func ParseBitsPerSecond(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "/s") || strings.HasSuffix(s, "ps") {
		s = s[:len(s)-2]
	} else {
		return 0, errors.New("Invalid denominator in unit: must be '/s' or 'ps' like in 'KB/s' or 'KBps'")
	}

	n, err := units.ParseUnit(s, bytesUnitMap)
	if err == nil {
		n *= 8
	} else {
		n, err = units.ParseUnit(s, metricBytesUnitMap)

		if err == nil {
			n *= 8
		} else {
			n, err = units.ParseUnit(s, bitsUnitMap)

			if err != nil {
				n, err = units.ParseUnit(s, metricBitsUnitMap)
			}
		}
	}

	return n, err
}
