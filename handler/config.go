package handler

import (
	"fmt"
	"time"
)

const periodLayout = "15:04"

type (
	// TimePeriod is a type to store specific moment of the day as offset from midnight in minutes.
	TimePeriod int

	Config struct {
		SelectStartHour TimePeriod `yaml:"select_start"`
		SelectEndHour   TimePeriod `yaml:"select_end"`
	}
)

// Implements the Unmarshaler interface of the yaml pkg.
func (t *TimePeriod) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var b string
	err := unmarshal(&b)
	if err != nil {
		return err
	}

	var parsedTime time.Time
	parsedTime, err = time.Parse(periodLayout, b)
	if err != nil {
		return fmt.Errorf("parse period time: %v", err)
	}
	*t = TimePeriod(parsedTime.Hour()*60 + parsedTime.Minute())

	return nil
}
