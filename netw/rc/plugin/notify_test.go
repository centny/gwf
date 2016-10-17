package plugin

import (
	"testing"

	"time"

	"github.com/Centny/gwf/netw/rc/rctest"
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
	}
}
