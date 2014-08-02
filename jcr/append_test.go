package jcr

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWalkC(t *testing.T) {
	filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		return walk_c("dir", "out", path, fi, errors.New("kkkkk"), func(path string) error {
			return nil
		})
	})
	filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		fmt.Println(path)
		return walk_c("dir", "out", path, fi, nil, func(path string) error {
			return nil
		})
	})
}
