package std

import (
	"os"
	"testing"
)

func TestUnZip(t *testing.T) {
	err := os.Mkdir("zip_output", 0755)
	AssertError(err, "mkdir zip_output")
	defer func() {
		_ = os.RemoveAll("zip_output")
	}()
	err = UnZip("sample.zip", "zip_output")
	AssertError(err, "unzip")
}
