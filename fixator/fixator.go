package fixator

import (
	"fixator/model"
	"log"
	"os"
	"path/filepath"
	"time"
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
	workerBuffer := buffer / len(velocityRanges)
	for _, vel := range velocityRanges {
		workers[vel] = make(chan *entity, buffer/workerBuffer)
	}
	return workers
}

func (f *Fixator) Fix(fixation *model.Fixation) error {
	f.channel <- fixation
	return nil
}

func (f *Fixator) Select(date model.FixationTime, velocity model.FixationFloat) ([]*model.Fixation, error) {
	start := time.Now()

	root := filepath.Join(f.root, time.Time(date).Format(dayLayout))

	channels := make([]chan []*model.Fixation, 24)
	for i := range channels {
		channels[i] = make(chan []*model.Fixation)
	}

	for i := 0; i < 24; i++ {
		go sortHour(root, i, velocity, channels[i])
	}

	result := make([]*model.Fixation, 0)
	for i := range channels {
		hourSorted := <-channels[i]
		result = append(result, hourSorted...)
	}

	elapsed := time.Since(start)
	log.Printf("select of %d entries took %s", len(result), elapsed)
	return result, nil
}

func (f *Fixator) SelectRange(date model.FixationTime) ([2]string, error) {
	start := time.Now()

	var minmax [2]string
	min, max := make(chan string), make(chan string)

	root := filepath.Join(f.root, time.Time(date).Format(dayLayout))
	go findMIN(root, min)
	go findMAX(root, max)

	for i := 0; i < 2; i++ {
		select {
		case x := <-min:
			minmax[0] = x
		case x := <-max:
			minmax[1] = x
		}
	}
	elapsed := time.Since(start)
	log.Printf("minmax took %s", elapsed)

	return minmax, nil
}
