package fixator

import (
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fixator/model"
)

const (
	FSRoot  = "./db"
	LogFile = "/home/katya/go/src/fixator/fixator/test.log"
)

func initBench(b *testing.B) (*Fixator, io.Writer) {
	// Set file for log output from server.
	w, err := os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		b.Fatal(err)
	}
	log.SetOutput(w)

	// Set random seed for fixation generations.
	rand.Seed(time.Now().UnixNano())

	// Define a config for fixator.
	conf := &Config{
		MinVelocity: min,
		MaxVelocity: max,
		FSRoot:      FSRoot,
	}

	return New(*conf), w
}

func initDB(b *testing.B, N int, f *Fixator) {
	start := time.Now()

	ch := make(chan struct{})
	for i := 0; i < N; i++ {
		fixation := &model.Fixation{
			Datetime: model.FixationTime(ranDate()),
			Car:      randStringRunes(9),
			Velocity: model.FixationFloat(f.minVelocity + rand.Float64()*(f.maxVelocity-f.minVelocity)),
		}
		go func(r *model.Fixation) {
			if err := f.Fix(r); err != nil {
				b.Logf("err while fix the fixation: %v", err)
			}
			ch <- struct{}{}
		}(fixation)
	}

	for i := 0; i < N; i++ {
		<-ch
	}

	elapsed := time.Since(start)
	b.Logf("init of %d entries took %s", N, elapsed)
}

func cleanBench(b *testing.B, f *Fixator, w io.Writer) {
	if err := removeContents(f.root); err != nil {
		b.Logf("error when clean out: %v", err)
	}
	w.(*os.File).Close()
}

var result interface{}

func BenchmarkFixator(b *testing.B) {
	fixator, w := initBench(b)
	initDB(b, 5000, fixator)
	// XXX: wait for some time until all entries are really written in the filesystem.
	time.Sleep(5 * time.Second)

	date, err := time.Parse(dayLayout, "01.01.2020")
	if err != nil {
		b.Fatal(err)
	}

	b.Run("MinMax", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var r interface{}
			for pb.Next() {
				r, err = fixator.SelectRange(model.FixationTime(date))
				if err != nil {
					b.Fatal(err)
				}
			}
			result = r
		})
	})

	b.Run("Select-Mid", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var r interface{}
			for pb.Next() {
				r, err = fixator.Select(model.FixationTime(date), (max-min)/2)
				if err != nil {
					b.Fatal(err)
				}
			}
			result = r
		})
	})

	b.Run("Select-Max", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var r interface{}
			for pb.Next() {
				r, err = fixator.Select(model.FixationTime(date), max)
				if err != nil {
					b.Fatal(err)
				}
			}
			result = r
		})
	})

	b.Run("Select-Min", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var r interface{}
			for pb.Next() {
				r, err = fixator.Select(model.FixationTime(date), min)
				if err != nil {
					b.Fatal(err)
				}
			}
			result = r
		})
	})

	cleanBench(b, fixator, w)
}

func ranDate() time.Time {
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2020, 1, 1, 24, 60, 60, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min

	return time.Unix(sec, 0)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		d, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		defer d.Close()
		names1, err := d.Readdirnames(-1)
		if err != nil {
			return err
		}
		for _, name1 := range names1 {
			err = os.RemoveAll(filepath.Join(dir, name, name1))
			if err != nil {
				return err
			}
		}
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	err = os.RemoveAll(filepath.Join(dir))
	if err != nil {
		return err
	}
	return nil
}
