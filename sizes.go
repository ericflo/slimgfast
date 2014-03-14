package slimgfast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Size struct {
	Width  uint
	Height uint
}

type SizeFile struct {
	Sizes  []Size
	Counts map[string]uint
}

type SizeCounter struct {
	filename string
	counts   map[string]uint
	done     chan struct{}
}

/* Size */

func SizeFromKey(key string) (Size, error) {
	var s Size
	splitKey := strings.Split(key, "x")
	if width, err := strconv.Atoi(splitKey[0]); err == nil {
		s.Width = uint(width)
	} else {
		return s, err
	}
	if height, err := strconv.Atoi(splitKey[1]); err == nil {
		s.Height = uint(height)
	} else {
		return s, err
	}
	return s, nil
}

func (size Size) Key() string {
	return fmt.Sprintf("%dx%d", size.Width, size.Height)
}

/* SizeCounter */

func NewSizeCounter(filename string) (*SizeCounter, error) {
	counts, err := getCountsFromFilename(filename)
	if err != nil {
		return nil, err
	}
	return &SizeCounter{
		filename: filename,
		done:     make(chan struct{}),
		counts:   counts,
	}, nil
}

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

	// Build the sizes map from the decoded json
	for _, s := range sizeFile.Sizes {
		key := s.Key()
		if sizeFile.Counts == nil {
			sizes[key] = 1
		} else {
			sizes[key] = sizeFile.Counts[key]
		}
	}

	return sizes, nil
}

func saveFile(counter *SizeCounter) error {
	sizes, err := counter.GetAllSizes()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(SizeFile{
		Sizes:  sizes,
		Counts: counter.counts,
	})
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(counter.filename, bytes, 0666); err != nil {
		log.Print("Error writing out sizes file: " + err.Error())
		return err
	}
	return nil
}

func (counter *SizeCounter) Start(every time.Duration) {
	go func() {
		ticks := time.Tick(every)
		select {
		case <-ticks:
			saveFile(counter)
		case <-counter.done:
			saveFile(counter)
			return
		}
	}()
	return
}

func (counter *SizeCounter) CountSize(size Size) {
	key := size.Key()
	if val, ok := counter.counts[key]; ok {
		counter.counts[key] = val + 1
	} else {
		counter.counts[key] = 1
	}
	return
}

func (counter *SizeCounter) GetAllSizes() ([]Size, error) {
	sizes := make([]Size, 0)
	for key := range counter.counts {
		if size, err := SizeFromKey(key); err == nil {
			sizes = append(sizes, size)
		} else {
			return nil, err
		}
	}
	return sizes, nil
}

func (counter *SizeCounter) GetTopSizesByCount(count uint) ([]Size, error) {
	sizes := make([]Size, 0)
	var i uint = 0
	for key := range counter.counts {
		if i == count {
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

func (counter *SizeCounter) Close() {
	close(counter.done)
	return
}
