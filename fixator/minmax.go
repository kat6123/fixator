package fixator

import (
	"bufio"
	"fixator/model"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type minmax struct {
	value          float64
	representation string
}

func findMIN(root string, out chan string) {
	start := time.Now()

	checkNext := true
	minCh := make(chan *minmax)
	s, min, minVal := 0, 0., ""

	var cur model.FixationFloat
	for i := 0; i < len(Files); i++ {
		{
			if !checkNext {
				break
			}
			for j := range Files[i] {
				go func(file string) {
					min := float64(max)
					minVal := ""

					w, err := os.Open(filepath.Join(root, file))
					if os.IsNotExist(err) {
						minCh <- &minmax{}
						return
					}
					// atomic??????
					checkNext = false

					scanner := bufio.NewScanner(w)
					for scanner.Scan() {
						value := scanner.Text()
						split := strings.Split(value, " ")
						if err := cur.UnmarshalJSON([]byte(split[3])); err != nil {
							log.Printf("find min: can't unmarshal %s: %v", value, err)
							minCh <- &minmax{}
							return
						}

						if float64(cur) < min {
							min, minVal = float64(cur), value
						}
					}

					if err := scanner.Err(); err != nil {
						log.Printf("find min: scanner error: %v", err)
						minCh <- &minmax{}
						return
					}
					minCh <- &minmax{min, minVal}
				}(Files[i][j])
			}
			for range Files[i] {
				x := <-minCh
				if *x != (minmax{}) && (s == 0 || x.value < min) {
					min, minVal = x.value, x.representation
					s++
				}
			}
		}
	}

	elapsed := time.Since(start)
	log.Printf("find min took %s", elapsed)

	out <- minVal
}

func findMAX(root string, out chan string) {
	start := time.Now()

	checkNext := true
	maxCh := make(chan *minmax)
	s, max, maxVal := 0, 0., ""

	var cur model.FixationFloat
	for i := len(Files) - 1; i >= 0; i-- {
		{
			if !checkNext {
				break
			}
			for j := range Files[i] {
				go func(file string) {
					max := float64(min)
					maxVal := ""

					w, err := os.Open(filepath.Join(root, file))
					if os.IsNotExist(err) {
						maxCh <- &minmax{}
						return
					}
					// atomic??????
					checkNext = false

					scanner := bufio.NewScanner(w)
					for scanner.Scan() {
						value := scanner.Text()
						split := strings.Split(value, " ")
						if err := cur.UnmarshalJSON([]byte(split[3])); err != nil {
							log.Printf("find max: can't unmarshal %s: %v", value, err)
							maxCh <- &minmax{}
							return
						}

						if float64(cur) > max {
							max, maxVal = float64(cur), value
						}
					}

					if err := scanner.Err(); err != nil {
						log.Printf("find max: scanner error: %v", err)
						maxCh <- &minmax{}
						return
					}
					maxCh <- &minmax{max, maxVal}
				}(Files[i][j])
			}
			for range Files[i] {
				x := <-maxCh
				if *x != (minmax{}) && (s == 0 || x.value > max) {
					max, maxVal = x.value, x.representation
					s++
				}
			}
		}

	}
	elapsed := time.Since(start)
	log.Printf("find max took %s", elapsed)

	out <- maxVal
}