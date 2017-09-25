package rcmd

import (
	"testing"
	"time"
)

func TestRCmd(t *testing.T) {
	err := StartMaster(":2832", map[string]int{
		"Ctrl-abc":  1,
		"Slave-abc": 1,
	})
	if err != nil {
		t.Errorf("start master fail with %v", err)
		return
	}
	err = StartSlave("slave1", "127.0.0.1:2832", "Slave-abc")
	if err != nil {
		t.Errorf("start slave fail with %v", err)
		return
	}
	err = StartControl("control", "127.0.0.1:2832", "Ctrl-abc")
	if err != nil {
		t.Errorf("start contorl fail with %v", err)
		return
	}
	time.Sleep(time.Second) //wait for salve connected
	//test start and done
	res, err := SharedControl.StartCmd("sleep 1 && echo abc", 0, "")
	if err != nil {
		t.Errorf("start cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "#1" {
		t.Errorf("start cmd fial with not result found %v", res)
		return
	}
	res, err = SharedControl.List()
	if err != nil {
		t.Errorf("list cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "#1" {
		t.Errorf("list cmd fial with not result found %v", res)
		return
	}
	time.Sleep(2 * time.Second)
	res, err = SharedControl.List()
	if err != nil {
		t.Errorf("list cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "" {
		t.Errorf("list cmd fial with result found %v", res)
		return
	}
	//test start and stop all
	res, err = SharedControl.StartCmd("sleep 10", 0, "")
	if err != nil {
		t.Errorf("start cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "#2" {
		t.Errorf("start cmd fial with not result found %v", res)
		return
	}
	res, err = SharedControl.StopCmd("", "#2")
	if err != nil {
		t.Errorf("stop cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "ok" {
		t.Errorf("stop cmd fial with not result found %v", res)
		return
	}
	res, err = SharedControl.List()
	if err != nil {
		t.Errorf("list cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "" {
		t.Errorf("list cmd fial with result found %v", res)
		return
	}
	//
	//test start and stop one
	res, err = SharedControl.StartCmd("sleep 10", 0, "")
	if err != nil {
		t.Errorf("start cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "#3" {
		t.Errorf("start cmd fial with not result found %v", res)
		return
	}
	res, err = SharedControl.StopCmd("slave1", "#3")
	if err != nil {
		t.Errorf("stop cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "ok" {
		t.Errorf("stop cmd fial with not result found %v", res)
		return
	}
	res, err = SharedControl.List()
	if err != nil {
		t.Errorf("list cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "" {
		t.Errorf("list cmd fial with result found %v", res)
		return
	}
	//test start shell and done
	res, err = SharedControl.StartCmd("sleep 1 && echo abc", 1, "")
	if err != nil {
		t.Errorf("start cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "#4" {
		t.Errorf("start cmd fial with not result found %v", res)
		return
	}
	res, err = SharedControl.List()
	if err != nil {
		t.Errorf("list cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "#4" {
		t.Errorf("list cmd fial with not result found %v", res)
		return
	}
	time.Sleep(2 * time.Second)
	res, err = SharedControl.List()
	if err != nil {
		t.Errorf("list cmd fail with %v", err)
		return
	}
	if len(res) < 1 || res.StrVal("slave1") != "" {
		t.Errorf("list cmd fial with result found %v", res)
		return
	}
}
