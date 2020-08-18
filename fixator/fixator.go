package fixator

import (
	"fixator/model"
	"fmt"
	"log"
	"os"
)

type (
	Config struct {
		MinVelocity float64 `yaml:"min_velocity"`
		MaxVelocity float64 `yaml:"max_velocity"`
		FSRoot      string  `yaml:"root"`
	}

	entity struct {
		Path  string
		Value string
	}

	Fixator struct {
		minVelocity float64
		maxVelocity float64
		root        string
		dirs        map[string]bool
		channel     chan *model.Fixation
		workers     map[string]chan *entity
	}
)

const buffer = 50

func New(config Config) *Fixator {
	fixator := &Fixator{
		minVelocity: config.MinVelocity,
		maxVelocity: config.MaxVelocity,
		root:        config.FSRoot,
		channel:     make(chan *model.Fixation, buffer),
		dirs:        make(map[string]bool),
		workers:     initWorkers(),
	}

	if _, err := os.Stat(fixator.root); os.IsNotExist(err) {
		if err := os.MkdirAll(fixator.root, os.ModePerm); err != nil {
			log.Fatalf("failed to create a file system root: %e", err)
		}
		log.Printf("created FS root at %s", fixator.root)
	}

	for k, ch := range fixator.workers {
		go fixator.save(ch)
		log.Printf("start saver for %s range\n", k)
	}
	go fixator.work()
	log.Printf("start worker with buffer size=%d\n", buffer)

	return fixator
}

func initWorkers() map[string]chan *entity {
	workers := make(map[string]chan *entity)
	velRange := getVelocityRanges()
	workerBuffer := buffer / len(velRange)
	for _, vel := range velRange {
		workers[vel] = make(chan *entity, buffer/workerBuffer)
	}
	return workers
}

func (f *Fixator) Fix(fixation *model.Fixation) error {
	f.channel <- fixation
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
