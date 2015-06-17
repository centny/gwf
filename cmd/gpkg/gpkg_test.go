package main

import (
	"os"
	"testing"
)

func TestLsd(t *testing.T) {
	ef = func(code int) {
	}
	os.Args = []string{"gpkg"}
	main()
	os.Args = []string{"gpkg", ".."}
	main()
	os.Args = []string{"gpkg", "-j", ",", ".."}
	main()
}
