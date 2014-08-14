package main

import (
	"database/sql"
	"fmt"
	"github.com/Centny/Cny4go/dbutil"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/routing"
	"github.com/Centny/Cny4go/test"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

type TSt struct {
	Tid    int64     `m2s:"TID"`
	Tname  string    `m2s:"TNAME"`
	Titem  string    `m2s:"TITEM"`
	Tval   string    `m2s:"TVAL"`
	Status string    `m2s:"STATUS"`
	Time   time.Time `m2s:"TIME"`
	T      int64     `m2s:"TIME" it:"Y"`
	Fval   float64   `m2s:"FVAL"`
	Uival  int64     `m2s:"UIVAL"`
	Add1   string    `m2s:"ADD1" json:"-"`
	Add2   string    `m2s:"Add2" json:"-"`
}

func main() {
	db, err := sql.Open("mysql", test.TDbCon)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	sb := routing.NewSrvSessionBuilder("", "/", "example", 60*60*1000, 10000)
	mux := routing.NewSessionMux("/example", sb)
	mux.HFilterFunc("^.*$", func(hs *routing.HTTPSession) routing.HResult {
		log.D("filt 001")
		return routing.HRES_CONTINUE
	})
	//http://localhost:8080/example/ok
	mux.HFunc("^/ok(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.MsgRes("OK")
		return routing.HRES_RETURN
	})
	//http://localhost:8080/example/data
	mux.HFunc("^/data(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		var tid int64
		var name string
		err := hs.ValidRVal(`//valid the argument
			tid,R|I,R:0,tid is error;
			name,R|S,L:0,name is empty;
			`, &tid, &name)
		if err != nil {
			return hs.MsgResE(1, err.Error())
		}
		name = "%" + name + "%"
		var ts []TSt
		err = dbutil.DbQueryS(db, &ts, "select * from ttable where tid>? and tname like ?", tid, name)
		if err != nil {
			return hs.MsgResE(1, err.Error())
		}
		return hs.MsgRes(ts)
	})
	mux.HFunc("^/mdata(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.W.Write([]byte("some data\n"))
		return routing.HRES_RETURN
	})
	s := http.Server{Addr: ":8080", Handler: mux}
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
