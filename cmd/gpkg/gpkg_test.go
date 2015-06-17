package main

import (
	"fmt"
	"os"
	"testing"
)

func TestLsd(t *testing.T) {
	ef = func(code int) {
	}
	fmt.Println("2")
	os.Args = []string{"gpkg"}
	main()
	fmt.Println("3")
	os.Args = []string{"gpkg", ".."}
	main()
	fmt.Println("4")
	os.Args = []string{"gpkg", "./"}
	main()
	fmt.Println("5")
	os.Args = []string{"gpkg", "-j", ",", ".."}
	main()
	fmt.Println("6")
	os.Args = []string{"gpkg", "-j", ",", "/sds/"}
	main()
	fmt.Println("7")
	os.Args = []string{"gpkg", "-j", ","}
	main()
	fmt.Println("8")
	os.Args = []string{"gpkg", "/sds/"}
	main()
	fmt.Println("9")
	os.Args = []string{"gpkg", "-p", "cmd/", "../../"}
	main()
	fmt.Println("10")
	os.Args = []string{"gpkg", "-h"}
	main()
}
