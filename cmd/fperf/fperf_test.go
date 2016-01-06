package main

import (
	"github.com/Centny/gwf/util"
	"os"
	"testing"
)

func TestPerf(t *testing.T) {
	ef = func(int) {
	}
	os.Args = []string{"sss", "-h"}
	main()
	//rw
	os.Args = []string{"sss"}
	main()
	util.Exec("rm", "-f", "test_*")
	//rw
	os.Args = []string{"sss", "-M", "W"}
	main()
	os.Args = []string{"sss", "-M", "R"}
	main()
	util.Exec("rm", "-f", "test_*")
	//rw
	os.Args = []string{"sss", "/tmp"}
	main()
	util.Exec("rm", "-f", "/tmp/test_*")
	//rw
	os.Args = []string{"sss", "-M", "RW", "-p", "xxs_", "-B", "102400", "-c", "11", "-t", "111", "-m", "8", "-b", "0", "-e", "10"}
	main()
	util.Exec("rm", "-f", "xxs_*")
	//
	//error
	//rw
	os.Args = []string{"sss", "-M", "RW", "/xsds/"}
	main()
	util.Exec("rm", "-f", "/tmp/xxs_*")
	os.Args = []string{"sss", "-M", "RW", "-B", "/xsds/"}
	main()
}
