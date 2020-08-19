package fixator

import (
	"bufio"
	"container/heap"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fixator/model"
)

func sortHour(root string, hour int, boundary model.FixationFloat, out chan []*model.Fixation) {
	start := rangeByFloat(boundary)
	toSelect := make([]string, 0)
	for i := range velocityRanges {
		if velocityRanges[i] == start {
			toSelect = velocityRanges[i:]
		}
	}

	channels := make([]chan *Item, len(toSelect))
	for i := range channels {
		channels[i] = make(chan *Item)
		go sort(
			filepath.Join(root, toSelect[i]+"-"+strconv.Itoa(hour)),
			float64(boundary), i, channels[i])
	}

	MERGED := make(PriorityQueue, 0)
	for i := 0; i < len(toSelect); i++ {
		// XXX- add check for close if the channel is empty
		it, ok := <-channels[i]
		if !ok {
			continue
		}
		MERGED = append(MERGED, it)
	}
	heap.Init(&MERGED)

	result := make([]*model.Fixation, 0)
	chToLook := 0
	for MERGED.Len() > 0 {
		item := MERGED[0]
		chToLook = item.FromChannel
		result = append(result, item.Value)

		//took from chanel
		it, ok := <-channels[chToLook]
		if !ok {
			heap.Remove(&MERGED, 0)
			continue
		}

		MERGED[0] = it
		heap.Fix(&MERGED, 0)
	}
	out <- result
	//close(out)
}

func sort(path string, bound float64, chnum int, out chan *Item) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("file %s not exists", path)
			close(out)
			return
		}
		log.Fatal(err)
	}
	defer file.Close()

	i := 0
	pq := make(PriorityQueue, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := scanner.Text()
		split := strings.Split(value, " ")
		unix, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		vel, err := strconv.ParseFloat(split[2], 32)
		if err != nil {
			log.Fatal(err)
		}

		if vel < bound {
			continue
		}

		pq = append(pq, &Item{
			Priority: unix,
			Value: &model.Fixation{
				Datetime: model.FixationTime(time.Unix(unix, 0)),
				Car:      split[1],
				Velocity: model.FixationFloat(vel),
			},
			Index:       i,
			FromChannel: chnum,
		})
		i++
	}
	heap.Init(&pq)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		out <- item
	}
	close(out)
}
