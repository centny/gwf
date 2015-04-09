package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"testing"
	"time"
)

func TestCfgs(t *testing.T) {
	os.Args = []string{"cfgs"}
	main()
	os.Args = []string{"cfgs", ":8790", "csd.properties", "./"}
	main()
	os.Args = []string{"cfgs", ":8790", "token.properties", "./"}
	go main()
	time.Sleep(200 * time.Millisecond)
	fmt.Println(util.HGet("http://127.0.0.1:8790/cfg?token=%v", "abc"))
	fmt.Println(util.HGet("http://127.0.0.1:8790/cfg?token=%v", "abd"))
	fmt.Println(util.HGet("http://127.0.0.1:8790/cfg?token=%v", "err"))
	fmt.Println(util.HGet("http://127.0.0.1:8790/cfg?token=%v", "xdsfk"))
	fmt.Println(util.HGet("http://127.0.0.1:8790/cfg?token=%v", ""))
}
