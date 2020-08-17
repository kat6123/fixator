package model

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

type (
	FixationTime  time.Time
	FixationFloat float64
)

type Fixation struct {
	Datetime FixationTime  `json:"Дата и время фиксации"`
	Car      string        `json:"Номер транспортного средства"`
	Velocity FixationFloat `json:"Скорость движения км/ч,string"`
}

const layout = "\"02.01.2006 15:04:05\""

func (t *FixationTime) UnmarshalJSON(b []byte) error {
	parsedTime, err := time.Parse(layout, string(b))
	if err != nil {
		return fmt.Errorf("parse fixation time: %v", err)
	}
	*t = FixationTime(parsedTime)
	return nil
}

func (t FixationTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(layout)), nil
}

func (f *FixationFloat) UnmarshalJSON(b []byte) error {
	b = bytes.Replace(b, []byte{','}, []byte{'.'}, 1)
	parsed, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return fmt.Errorf("parse fixation velocity: %v", err)
	}
	*f = FixationFloat(parsed)
	return nil
}

func (f FixationFloat) MarshalJSON() ([]byte, error) {
	parsed := []byte(
		`"` + strconv.FormatFloat(float64(f), 'g', -1, 64) + `"`)
	b := bytes.Replace(parsed, []byte{'.'}, []byte{','}, 1)
	return b, nil
}
