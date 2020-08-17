package fixator

import (
	"fixator/model"
	"fmt"
)

type Fixator struct {
}

func (f *Fixator) Fix(fixation *model.Fixation) error {
	fmt.Printf("%+v\n", fixation)
	return nil
}

func (f *Fixator) Select(date model.FixationTime, velocity model.FixationFloat) ([]*model.Fixation, error) {
	fmt.Printf("%v, %v\n", date, velocity)
	velocityRange := make([]*model.Fixation, 2)

	velocityRange[0] = &model.Fixation{
		Datetime: model.FixationTime{},
		Car:      "sAS",
		Velocity: 65,
	}
	velocityRange[1] = &model.Fixation{
		Datetime: model.FixationTime{},
		Car:      "sAS",
		Velocity: 65,
	}
	return velocityRange, nil
}

func (f *Fixator) SelectRange(date model.FixationTime) ([2]*model.Fixation, error) {
	var velocityRange [2]*model.Fixation

	velocityRange[0] = &model.Fixation{
		Datetime: model.FixationTime{},
		Car:      "sAS",
		Velocity: 65,
	}
	velocityRange[1] = &model.Fixation{
		Datetime: model.FixationTime{},
		Car:      "sAS",
		Velocity: 65,
	}

	return velocityRange, nil
}
