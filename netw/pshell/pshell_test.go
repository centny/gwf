package pshell

import (
	"testing"
	"time"
)

func TestPshell(t *testing.T) {
	err := StartServer(":2832", map[string]int{
		"Ctrl-abc":  1,
		"Slave-abc": 1,
	}, &Host{
		Name:     "loc.m",
		Addr:     "loc.m:22",
		Username: "root",
		Password: "sco",
	})
	if err != nil {
		t.Errorf("start master fail with %v", err)
		return
	}
	err = StartControl("control", "127.0.0.1:2832", "Ctrl-abc")
	if err != nil {
		t.Errorf("start contorl fail with %v", err)
		return
	}
	time.Sleep(time.Second) //wait for salve connected
	//
	err = SharedControl.AddSession("loc.m2", "loc.m:22", "root", "sco")
	if err != nil {
		t.Errorf("add session fail with %v", err)
		return
	}
	//
	res, err := SharedControl.Exec("loc.m", "", "uptime")
	if err != nil || len(res) < 1 {
		t.Errorf("exec fail with %v", err)
		return
	}
	res, err = SharedControl.Exec("", "", "uptime")
	if err != nil || len(res) < 1 {
		t.Errorf("exec fail with %v", err)
		return
	}
	time.Sleep(time.Second)
	res, err = SharedControl.Exec("", `
	cd /tmp
	pwd
	echo $1
		`, "abc")
	if err != nil || len(res) < 1 {
		t.Errorf("exec fail with %v", err)
		return
	}
	time.Sleep(time.Second)
}
