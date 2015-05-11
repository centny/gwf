package main

import (
	"os"
	"testing"
)

func TestMexec(t *testing.T) {
	exit = func(code int) {}
	os.Remove("sp.json")
	os.Args = []string{"mexec", "-min=0", "-max=8", "-total=8", "-file=sp.json", "-log", "-id", "/bin/echo", "abc"}
	main()
	os.Args = []string{"mexec", "-min=0", "-max=8", "-total=8", "-log", "-id", "/bin/echo", "abc"}
	main()
	//
	os.Args = []string{"mexec", "-min=0", "-max=8", "-total=8", "-file=/sd/sp.json", "-log", "-id", "/bin/echo", "abc"}
	main()
	os.Args = []string{"mexec", "-min=x", "-max=x", "-total=x", "-file=sp.json", "-log", "-id", "/bin/echo", "abc"}
	main()
}
