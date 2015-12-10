package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"os"
	"regexp"
	"strings"
	"time"
)

type TimeoutKiller struct {
	Timeout int64
	PS      []string
	PT      map[string]int64
}

func (t *TimeoutKiller) Do() error {
	res, err := util.Exec(t.PS...)
	if err != nil {
		return err
	}
	reg := regexp.MustCompile("\\s+")
	lines := strings.Split(res, "\n")
	now := util.Now()
	lsed := map[string]int{}
	for _, line := range lines {
		line = strings.Trim(line, " \t")
		items := reg.Split(line, -1)
		if len(items) < 2 {
			continue
		}
		log.D("checking line->%v", line)
		pid := strings.Trim(items[1], " \t")
		lsed[pid] = 1
		ot, ok := t.PT[pid]
		if ok {
			if now-ot > t.Timeout {
				res, err = util.Exec("kill", pid)
				log.D("kill %v->%v,%v", pid, res, err)
			}
		} else {
			t.PT[pid] = now
		}
	}
	for k, _ := range t.PT {
		if _, ok := lsed[k]; ok {
			continue
		}
		log.D("remove pid->%v", k)
		delete(t.PT, k)
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: kzp script timeout")
		return
	}
	timeout, err := util.ParseInt64(os.Args[2])
	if err != nil {
		panic(err.Error())
	}
	tk := TimeoutKiller{
		Timeout: timeout * 1000,
		PS:      []string{os.Args[1]},
		PT:      map[string]int64{},
	}
	for {
		log.D("do Timeout Killer...")
		err := tk.Do()
		if err != nil {
			log.E("%v", err)
		}
		time.Sleep(10 * time.Second)
	}
}
