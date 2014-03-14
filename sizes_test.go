package slimgfast

import (
	"os"
	"testing"
)

func GetTmpSizesJson() string {
	return "/tmp/sizes.json"
}

func ClearTmpSizesJson() error {
	filename := "/tmp/sizes.json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filename)
}

func TestSizeFromKey(t *testing.T) {
	var key string = "15x15"
	var expected uint = 15
	size, err := SizeFromKey(key)
	if err != nil {
		t.Error(err.Error())
	}
	if size.Width != expected {
		t.Error("Expected width:", expected, "Actual:", size.Width)
	}
	if size.Height != expected {
		t.Error("Expected height:", expected, "Actual:", size.Height)
	}
}

func TestSizeFromBadKey(t *testing.T) {
	badKeys := []string{"Zx15", "-12x-12", "-12x12", "asdf", "x", "nullxnull"}
	for _, key := range badKeys {
		size, err := SizeFromKey(key)
		if err == nil {
			t.Error("Expected an error parsing:", key, "Got:", size)
		}
	}
}

func TestNewSizeCounter(t *testing.T) {
	if err := ClearTmpSizesJson(); err != nil {
		t.Error("Could not clear temporary sizes.json file")
	}
	if counter, err := NewSizeCounter(GetTmpSizesJson()); err != nil {
		t.Error("Could not create a new size counter")
	} else {
		counter.Close()
	}
}
