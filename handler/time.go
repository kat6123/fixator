package handler

import (
	"fmt"
	"time"
)

// TimePeriod is a type to store specific moment of the day as offset from midnight in minutes.
type TimePeriod int

const periodLayout = "15:04"

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
