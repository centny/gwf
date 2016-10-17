package plugin

import (
	"testing"

	"time"

	"fmt"

	"github.com/Centny/gwf/netw/rc/rctest"
	"github.com/Centny/gwf/util"
)

func TestNotify(t *testing.T) {
	ShowLog = 1
	var err error
	var rct *rctest.RCTest
	var mdb *NotifyMemDb
	var srv *NotifySrv
	var client *NotifyClient
	var received = 0

	//
	//test normal post
	{
		received = 0
		rct = rctest.NewRCTest_j2(":9332")
		mdb = NewNotifyMemDb()
		srv = NewNotifySrv(mdb)
		srv.Hand(rct.L)
		srv.Start()
		//
		client = NewNotifyClient(NotifyHandlerF(func(n *NotifyClient, m *Message) error {
			if m.Type == "testing" || m.Type == "xx" {
				var err = n.Mark(m.ID, "testing")
				if err != nil {
					t.Error(err)
					return err
				}
				received++
				return err
			}
			return nil
		}))
		client.SetRunner(rct.R, rct.Rmh)
		//
		err = client.Monitor("testing,xx", 10)
		if err != nil {
			t.Error(err)
			return
		}
		for i := 0; i < 100; i++ {
			err = srv.PostMessage(&Message{
				Type: "testing",
			})
			if err != nil {
				t.Error(err)
				return
			}
			err = srv.PostMessage(&Message{
				Type: "xx",
			})
			if err != nil {
				t.Error(err)
				return
			}
		}
		for x := 0; x < 10 && received < 200; x++ {
			time.Sleep(100 * time.Millisecond)
		}
		if received != 200 {
			t.Error("not received")
			return
		}
		if len(mdb.MS) > 0 {
			t.Error("error")
			return
		}
		srv.Stop()
		rct.R.Stop()
		rct.L.Close()
		rct.L.Wait()
		fmt.Printf("\n\n\n")
	}
	//
	//test offline
	{
		received = 0
		rct = rctest.NewRCTest_j2(":9332")
		mdb = NewNotifyMemDb()
		srv = NewNotifySrv(mdb)
		srv.Hand(rct.L)
		srv.Start()
		for i := 0; i < 100; i++ {
			err = srv.PostMessage(&Message{
				Type: "testing",
			})
			if err != nil {
				t.Error(err)
				return
			}
			err = srv.PostMessage(&Message{
				Type: "xx",
			})
			if err != nil {
				t.Error(err)
				return
			}
		}
		if len(mdb.MS) != 200 {
			t.Error("error")
			return
		}
		client = NewNotifyClient(NotifyHandlerF(func(n *NotifyClient, m *Message) error {
			if m.Type == "testing" {
				var err = n.Mark(m.ID, "testing")
				if err != nil {
					t.Error(err)
					return err
				}
				received++
				return err
			}
			return nil
		}))
		client.SetRunner(rct.R, rct.Rmh)
		err = client.Monitor("testing", 10)
		if err != nil {
			t.Error(err)
			return
		}
		for x := 0; x < 10 && received < 100; x++ {
			time.Sleep(100 * time.Millisecond)
		}
		if received != 100 {
			t.Error("not received")
			return
		}
		if len(mdb.MS) != 100 {
			t.Error("error")
			return
		}
		srv.Stop()
		rct.R.Stop()
		rct.L.Close()
		rct.L.Wait()
		fmt.Printf("\n\n\n")
	}
	//
	//test mark count
	{
		received = 0
		rct = rctest.NewRCTest_j2(":9332")
		mdb = NewNotifyMemDb()
		srv = NewNotifySrv(mdb)
		srv.Hand(rct.L)
		srv.Start()
		mdb.Count["testing"] = 2
		mdb.Count["xx"] = 3
		client = NewNotifyClient(NotifyHandlerF(func(n *NotifyClient, m *Message) error {
			if m.Type == "testing" || m.Type == "xx" {
				// fmt.Printf("---->%v\n", m.ID)
				err = n.Mark(m.ID, "testing")
				if err != nil {
					t.Error(err)
					return err
				}
				err = n.Mark(m.ID, "testing")
				if err != nil {
					t.Error(err)
					return err
				}
				received++
				return err
			}
			return nil
		}))
		client.SetRunner(rct.R, rct.Rmh)
		err = client.Monitor("testing,xx", 10)
		if err != nil {
			t.Error(err)
			return
		}
		for i := 0; i < 1; i++ {
			err = srv.PostMessage(&Message{
				Type: "testing",
			})
			if err != nil {
				t.Error(err)
				return
			}
			err = srv.PostMessage(&Message{
				Type: "xx",
			})
			if err != nil {
				t.Error(err)
				return
			}
		}
		if len(mdb.MS) != 2 {
			t.Error("error")
			return
		}
		for x := 0; x < 10 && received < 200; x++ {
			time.Sleep(100 * time.Millisecond)
		}
		if received != 2 {
			t.Error("not received")
			return
		}
		if len(mdb.MS) != 1 {
			t.Error("error")
			return
		}
		srv.Stop()
		rct.R.Stop()
		rct.L.Close()
		rct.L.Wait()
		fmt.Printf("\n\n\n")
	}
	//
	//test error
	{
		received = 0
		rct = rctest.NewRCTest_j2(":9332")
		mdb = NewNotifyMemDb()
		srv = NewNotifySrv(mdb)
		srv.Hand(rct.L)
		srv.Start()
		mdb.Count["testing"] = 2
		mdb.Count["xx"] = 3
		client = NewNotifyClient(NotifyHandlerF(func(n *NotifyClient, m *Message) error {
			if m.Type == "testing" || m.Type == "xx" {
				received++
				return err
			} else if m.Type == "" {
				return util.Err("error")
			}
			return nil
		}))
		client.SetRunner(rct.R, rct.Rmh)
		err = client.Monitor("testing,xx", 10)
		if err != nil {
			t.Error(err)
			return
		}
		//
		var cid string
		for x := range srv.clients["testing"] {
			cid = x
		}
		//parsing error
		srv.L.MsgC(cid).Writev2([]byte{NotifyMessageMark}, util.Map{
			"attrs": "xxxx",
			"count": "xxxx",
		})
		//type error
		srv.L.MsgC(cid).Writev2([]byte{NotifyMessageMark}, util.Map{})
		//argument error
		err = client.Monitor("", 10)
		if err == nil {
			t.Error("error")
			return
		}
		//argument error
		err = client.Mark("", "")
		if err == nil {
			t.Error("error")
			return
		}
		//done not found
		err = client.Mark("xxxdd", "xxxx")
		if err == nil {
			t.Error("error")
			return
		}
		//client not found
		srv.clients["kjd"] = map[string]int{"kkk": 10}
		srv.notifyMessage(&Message{Type: "kjd"})
		srv.notifyClient("xxx", "xx")
		//
		srv.Start()
		//
		time.Sleep(300 * time.Millisecond)
		srv.Stop()
		rct.R.Stop()
		rct.L.Close()
		rct.L.Wait()
		fmt.Printf("\n\n\n")
	}
}
