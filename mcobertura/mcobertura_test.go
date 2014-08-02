package main

import (
	"os"
	"testing"
)

func TestR(t *testing.T) {
	os.Args = []string{"abc"}
	main()
	os.Args = []string{"abc", "j"}
	main()
	os.Args = []string{"abc", "j", "-o"}
	main()
	os.Args = []string{"abc", "c1.xml", "c3.xml", "-o", "c.xml"}
	main()
	os.Args = []string{"abc", "c1.xml", "ce.xml", "-o", "c.xml"}
	main()
	os.Args = []string{"abc", "c1.xml", "c2.xml", "-o", "/kdfd/c.xml"}
	main()
	os.Args = []string{"abc", "c1.xml", "c2.xml", "-o", "c.xml"}
	main()
}
