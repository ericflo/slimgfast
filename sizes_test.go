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
