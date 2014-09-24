package main

import (
	"bufio"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"io"
	"os"
)

func main() {
	//create and open log file
	util.FTouch("/tmp/tt.tmp")
	f, err := os.OpenFile("/tmp/tt.tmp", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
		return
	}
	//create file buffer writer.
	bo := bufio.NewWriter(f)
	defer f.Close()
	defer bo.Flush()
	//set the log writer to file and stdout.
	//or log.SetWriter(bo) to only file.
	log.SetWriter(io.MultiWriter(bo, os.Stdout))
	//set the log level
	log.SetLevel(log.DEBUG)
	//use
	log.I("tesing")
	log.D("tesing")
}
