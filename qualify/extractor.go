//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//

// [ Coming together is a beginning; keeping together is progress; working ]
// [ together is success.                                   - Henry Ford - ]

package qualify

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ufuchs/itplus/base/fcc"
)

type Measurement struct {
	Host           string  `json:"host"`
	Num            int     `json:"num"`
	Alias          string  `json:"alias"`
	PhenomenonTime int64   `json:"phenomenontime"`
	Lon            float64 `json:"lon"`
	Lat            float64 `json:"lat"`
	Alt            float64 `json:"alt"`
	Temp           float64 `json:"temp"`
	Hum            int     `json:"pressure"`
	Press          int     `json:"humidity"`
	LowBatt        bool    `json:"lowbattery"`
}

type extractfn func([]string, int64, int) Measurement

type Extractor struct{}

func (e *Extractor) calcTemp(tempH int, tempL int) int {
	return ((tempH*256 + tempL) - 1000)
}

func (e *Extractor) calcPress(pressH int, pressL int) int {
	return pressH*256 + pressL
}

//
//
//
func (e *Extractor) spacesToUnderline(s string) string {
	s1 := strings.Split(s, " ")
	return strings.Join(s1, "_")
}

//
//
//
func (e *Extractor) split_TX25TP_Alias(alias string) (string, string) {

	var (
		sep  = ":"
		part = make([]string, 2)
	)

	copy(part, strings.Split(alias, sep))

	if len(part[1]) == 0 {
		part[1] = "outer temp"
	}

	return part[0], part[1]

}

//
// extractData_9
//
func (e *Extractor) extractData_9(data []string, when int64, num int) Measurement {

	hum, _ := strconv.Atoi(data[6])
	tempH, _ := strconv.Atoi(data[4])
	tempL, _ := strconv.Atoi(data[5])

	return Measurement{
		Num:            num,
		PhenomenonTime: when,
		Temp:           float64(e.calcTemp(tempH, tempL)) / 10.0,
		Press:          -999,
		Hum:            hum,
		LowBatt:        (hum & 0x80) == 0x80,
	}

}

//
// extractData_0
//
func (e *Extractor) extractData_0(data []string, when int64, num int) Measurement {

	tempH, _ := strconv.Atoi(data[4])
	tempL, _ := strconv.Atoi(data[5])

	lowBatt, _ := strconv.Atoi(data[15])

	pressH, _ := strconv.Atoi(data[16])
	pressL, _ := strconv.Atoi(data[17])

	return Measurement{
		Num:            num,
		PhenomenonTime: when,
		Temp:           float64(e.calcTemp(tempH, tempL)),
		Press:          e.calcPress(pressH, pressL),
		Hum:            -999,
		LowBatt:        (lowBatt & 4) == 4,
	}

}

//
//
//
func (e *Extractor) ExtractData(ts fcc.Timestamped, device fcc.Device) (Measurement, error) {

	var extract extractfn

	data := strings.Fields(string((ts.Data)))

	switch data[1] {
	case "WS":
		extract = e.extractData_0
		break
	case "9":
		extract = e.extractData_9
		break
	default:
		// 2017/04/01 08:26:44 ==>New sensor type detected --> OK LS 6 0 12 88 3 200 55
		s := "==> New sensor type detected --> " + strings.Join(data, " ")
		fmt.Println(s)
		return Measurement{}, errors.New(s)

	}

	m := extract(data, ts.When, device.Num)

	alias := device.Alias
	if device.Type == "TX25TP-IT" {

		inner, outer := e.split_TX25TP_Alias(alias)

		switch m.Hum {
		case 106:
			alias = inner
		case 125:
			alias = outer
		}

		m.Hum = -999

	}

	alias = e.spacesToUnderline(alias)

	m.Alias = alias
	m.Lon = device.Lon
	m.Lat = device.Lat
	m.Alt = device.Alt

	return m, nil

}
