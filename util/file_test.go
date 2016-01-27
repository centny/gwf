package util

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	fmt.Println(FilterDir("/Users/vty/vgo/src/w.gdy.io/dyf", nil, []string{".*\\.git"}))
}
