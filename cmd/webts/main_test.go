package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	os.Args = []string{"webtest", ":8090", "."}
	go main()
	time.Sleep(time.Second)
	os.Args = []string{"webtest", "-h"}
	go main()
	time.Sleep(time.Second)
}
