package file

import (
	"io"
	"os"
)

func ReadBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}
