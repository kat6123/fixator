package fixator

import (
	"fixator/model"
	"fmt"
	"strconv"
	"time"
)

const (
	rangeStart, rangeEnd, step, min, max = 40, 150, 10, 0, 250
	dayLayout                            = "02.01.2006"
)

var (
	// velocityRanges should be SORTED
	velocityRanges []string
	Files          [][]string
)

func init() {
	velocityRanges = []string{strconv.Itoa(rangeStart)}
	for i := rangeStart + step; i <= rangeEnd; i += step {
		velocityRanges = append(velocityRanges, strconv.Itoa(i))
	}
	velocityRanges = append(velocityRanges, strconv.Itoa(max))

	Files = make([][]string, len(velocityRanges))
	for i, v := range velocityRanges {
		Files[i] = make([]string, 24)
		for j := 1; j <= 24; j++ {
			Files[i][j-1] = fmt.Sprintf("%s-%d", v, j)
		}
	}
}

func getFixationInfo(fixation *model.Fixation) (dir string, channel string, path string, value string) {
	dir = time.Time(fixation.Datetime).Format(dayLayout)
	channel = rangeByFloat(fixation.Velocity)
	path = fmt.Sprintf("%s-%d", channel, time.Time(fixation.Datetime).Hour())
	value = fixation.String()

	return
}

func rangeByFloat(v model.FixationFloat) string {
	vel, start, end := int(v)/step, rangeStart/step, rangeEnd/step

	velocity := ""
	if vel < start {
		velocity = strconv.Itoa(rangeStart)
	} else if vel >= end {
		velocity = strconv.Itoa(max)
	} else {
		velocity = strconv.Itoa((vel + 1) * step)
	}
	return velocity
}
