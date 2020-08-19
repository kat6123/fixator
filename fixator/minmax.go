package fixator

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type minmax struct {
	value          float64
	representation string
}

func find(name, root string, direct []int, comparison func(float64, float64) bool, out chan string) {
	start := time.Now()

	checkNext := true
	ch := make(chan *minmax)
	s, value, valEntry := 0, 0., ""

	for _, i := range direct {
		{
			if !checkNext {
				break
			}
			for j := range Files[i] {
				go findInFile(name, root, Files[i][j], comparison, ch)
			}
			for range Files[i] {
				x := <-ch
				if *x != (minmax{}) && (s == 0 || comparison(x.value, value)) {
					checkNext = false
					value, valEntry = x.value, x.representation
					s++
				}
			}
		}

	}
	elapsed := time.Since(start)
	log.Printf("find %s took %s", name, elapsed)

	out <- valEntry
}

func findInFile(name, root, file string, comparison func(float64, float64) bool, ch chan *minmax) {
	s, val, valEntry := 0, float64(min), ""

	w, err := os.Open(filepath.Join(root, file))
	if os.IsNotExist(err) {
		ch <- &minmax{}
		return
	}

	scanner := bufio.NewScanner(w)
	for scanner.Scan() {
		value := scanner.Text()
		split := strings.Split(value, " ")
		cur, err := strconv.ParseFloat(split[2], 32)
		if err != nil {
			log.Printf("find %s: can't unmarshal %s: %v", name, value, err)
			ch <- &minmax{}
			return
		}

		if s == 0 || comparison(cur, val) {
			val, valEntry = cur, value
			s++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("find %s: scanner error: %v", name, err)
		ch <- &minmax{}
		return
	}
	ch <- &minmax{val, valEntry}
}

func less(a, b float64) bool {
	return a < b
}

func more(a, b float64) bool {
	return a > b
}
