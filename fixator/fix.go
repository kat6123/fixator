package fixator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (f *Fixator) save(in <-chan *entity) {
	for fixation := range in {
		w, err := os.OpenFile(fixation.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Printf("failed to open file %s: value %s: %v", fixation.Path, fixation.Value, err)
			continue
		}

		if _, err := fmt.Fprintln(w, fixation.Value); err != nil {
			log.Printf("failed to write value %s: %v", fixation.Value, err)
			continue
		}

		if err := w.Close(); err != nil {
			log.Printf("failed to close file %s: value %s: %v", fixation.Path, fixation.Value, err)
			continue
		}
		log.Printf("saved at %s: value %s", fixation.Path, fixation.Value)
	}
}

func (f *Fixator) work() {
	for fixation := range f.channel {
		dir, channel, path, value := getFixationInfo(fixation)

		dirPath := filepath.Join(f.root, dir)
		if _, ok := f.dirs[dir]; !ok {
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
					log.Printf("failed to create folder for the day %s: value %s: %v", dir, value, err)
					continue
				}
				log.Printf("created folder for the day %s", dir)
			}
			f.dirs[dir] = true
		}

		f.workers[channel] <- &entity{filepath.Join(dirPath, path), value}
	}
}
