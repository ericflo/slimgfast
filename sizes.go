package slimgfast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Size is a Width and a Height
type Size struct {
	Width  uint
	Height uint
}

// SizeFile is a struct that is used to serialize the aggregated Size counts
// to a JSON file.
type SizeFile struct {
	Counts map[string]uint
}

// SizeCounter keeps track of how many times a certain size was requested, and
// persists this information to disk on a periodic basis.
type SizeCounter struct {
	filename string
	counts   map[string]uint
	done     chan struct{}
	mut      *sync.RWMutex
}

// SizeFromKey takes a string of the form WIDTHxHEIGHT and parses it into a
// Size struct.
func SizeFromKey(key string) (Size, error) {
	var s Size
	splitKey := strings.Split(key, "x")
	if width, err := strconv.Atoi(splitKey[0]); err == nil {
		if width < 1 {
			return s, fmt.Errorf("Got a key with an invalid width: %s", key)
		}
		s.Width = uint(width)
	} else {
		return s, err
	}
	if height, err := strconv.Atoi(splitKey[1]); err == nil {
		if height < 1 {
			return s, fmt.Errorf("Got a key with an invalid height: %s", key)
		}
		s.Height = uint(height)
	} else {
		return s, err
	}
	return s, nil
}

// Key takes the current Size object's Width and Height and generates a string
// of the form WIDTHxHEIGHT
func (size Size) Key() string {
	return fmt.Sprintf("%dx%d", size.Width, size.Height)
}

// NewSizeCounter initializes a *SizeCounter struct, loads in, and parses the
// persisted sizes.
func NewSizeCounter(filename string) (*SizeCounter, error) {
	counts, err := getCountsFromFilename(filename)
	if err != nil {
		return nil, err
	}
	mut := &sync.RWMutex{}
	return &SizeCounter{
		filename: filename,
		done:     make(chan struct{}),
		counts:   counts,
		mut:      mut,
	}, nil
}

// getCountsFromFilename loads in and parses the persisted sizes stored at the
// specified filename.
func getCountsFromFilename(filename string) (map[string]uint, error) {
	sizes := make(map[string]uint)

	// If there's no file, then there's nothing to read
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return sizes, nil
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Print("Could not open file: " + err.Error())
		return sizes, err
	}
	defer file.Close()

	// Decode the json
	decoder := json.NewDecoder(file)
	var sizeFile SizeFile
	if err := decoder.Decode(&sizeFile); err != nil {
		log.Fatal(err)
	}

	return sizeFile.Counts, nil
}

// saveFile serializes and persists the aggregated size stats to the filesystem.
func saveFile(counter *SizeCounter) error {
	counter.mut.RLock()
	defer counter.mut.RUnlock()
	bytes, err := json.Marshal(SizeFile{Counts: counter.counts})
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(counter.filename, bytes, 0666); err != nil {
		log.Print("Error writing out sizes file: " + err.Error())
		return err
	}
	return nil
}

// Start starts the counter persisting its aggregated size stats to disk
// periodically.
func (counter *SizeCounter) Start(every time.Duration) {
	go func() {
		ticks := time.Tick(every)
		for {
			select {
			case <-ticks:
				saveFile(counter)
			case <-counter.done:
				saveFile(counter)
				return
			}
		}
	}()
	return
}

// CountSize notes the size of one request.
func (counter *SizeCounter) CountSize(size *Size) {
	counter.mut.Lock()
	defer counter.mut.Unlock()
	key := size.Key()
	if val, ok := counter.counts[key]; ok {
		counter.counts[key] = val + 1
	} else {
		counter.counts[key] = 1
	}
	return
}

// GetAllSizes gets a list of all the sizes that we've seen.
func (counter *SizeCounter) GetAllSizes() ([]Size, error) {
	counter.mut.RLock()
	defer counter.mut.RUnlock()
	sizes := []Size{}
	for key := range counter.counts {
		if size, err := SizeFromKey(key); err == nil {
			sizes = append(sizes, size)
		} else {
			return nil, err
		}
	}
	return sizes, nil
}

// GetTopSizesByCount gets a list of all the sizes who have been requested at
// least `count` times.
func (counter *SizeCounter) GetTopSizesByCount(count uint) ([]Size, error) {
	// TODO: Sort first or something, this is not correct right now, it'll pick random values
	counter.mut.RLock()
	defer counter.mut.RUnlock()
	sizes := make([]Size, 0)
	i := 0
	for key := range counter.counts {
		if uint(i) == count {
			return sizes, nil
		}
		if size, err := SizeFromKey(key); err == nil {
			sizes = append(sizes, size)
		} else {
			return nil, err
		}
		i += 1
	}
	return sizes, nil
}

// TODO: GetTopSizesByPercentage?

// Close stops the SizeCounter from doing any more persistence.
func (counter *SizeCounter) Close() {
	close(counter.done)
	return
}
