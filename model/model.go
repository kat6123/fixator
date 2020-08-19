package model

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	FixationTime  time.Time
	FixationFloat float64

	Fixation struct {
		Datetime FixationTime  `json:"Дата и время фиксации"`
		Car      string        `json:"Номер транспортного средства"`
		Velocity FixationFloat `json:"Скорость движения км/ч,string"`
	}
)

const layout = "02.01.2006 15:04:05"

func (f FixationFloat) String() string {
	parsed := strconv.FormatFloat(float64(f), 'f', 2, 64)
	s := strings.Replace(parsed, ".", ",", 1)
	return s
}

func (t FixationTime) String() string {
	return time.Time(t).Format(layout)
}

func (f Fixation) String() string {
	return fmt.Sprintf("%d %s %.2f", time.Time(f.Datetime).Unix(), f.Car, float64(f.Velocity))
}

func (t *FixationTime) UnmarshalJSON(b []byte) error {
	parsedTime, err := time.Parse(`"`+layout+`"`, string(b))
	if err != nil {
		return fmt.Errorf("parse fixation time: %v", err)
	}
	*t = FixationTime(parsedTime)
	return nil
}

func (t FixationTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
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
	return []byte(`"` + f.String() + `"`), nil
}
