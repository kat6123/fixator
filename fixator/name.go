package fixator

import (
	"fixator/model"
	"fmt"
	"strconv"
	"time"
)

const rangeStart, rangeEnd, step, max = 40, 150, 10, 250

func getVelocityRanges() []string {
	v := []string{strconv.Itoa(rangeStart), strconv.Itoa(max)}
	for i := rangeStart + step; i <= rangeEnd; i += step {
		v = append(v, strconv.Itoa(i))
	}

	return v
}

func getFixationInfo(fixation *model.Fixation) (dir string, channel string, path string, value string) {
	dir = time.Time(fixation.Datetime).Format("02.01.2006")
	channel, path = getFixationPath(fixation)
	value = fixation.String()

	return
}

func getFixationPath(f *model.Fixation) (string, string) {
	vel, start, end := int(f.Velocity)/step, rangeStart/step, rangeEnd/step
	hour := time.Time(f.Datetime).Hour()
	velocity := ""

	if vel < start {
		velocity = strconv.Itoa(rangeStart)
	} else if vel >= end {
		velocity = strconv.Itoa(max)
	} else {
		velocity = strconv.Itoa((vel + 1) * step)
	}

	return velocity, fmt.Sprintf("%s-%d", velocity, hour)
}
