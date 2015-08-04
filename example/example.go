package main

import (
	"database/sql"
	"fmt"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/filter"
	"github.com/Centny/gwf/test"
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
	/*
		json arguments:
		{
		  "a":{
		  	"a1":"val",
		  	"a2":100
		  },
		  "b":{
			"b1":"val2",
			"b2":200
		  }
		}
	*/
	mux.HFunc("^/jsonv(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		mv, err := hs.JsonVal("json")
		if err != nil {
			return hs.MsgResErr2(1, "arg-err", err)
		}
		var a1, b1 string
		var a2, b2 int64
		err = mv.ValidF(`
			a/a1,R|S,L:0;
			a/a2,R|I,R:0;
			b/b1,R|S,L:0;
			b/b2,R|I,R:0;
			`, &a1, &a2, &b1, &b2)
		if err != nil {
			return hs.MsgResErr2(1, "arg-err", err)
		}
		fmt.Println(a1, a2, b1, b2)
		return hs.MsgRes("ok")
	})
	var cache1, cache2 int = 0, 0
	mux.HFunc("^/cache1(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		cache1++
		hs.W.Header().Set("Last-Modified", "Tue, 01 Jan 1980 1:00:00 GMT")
		return hs.MsgRes(fmt.Sprintf("%v", cache1))
	})
	mux.HFilterFunc("^/cache2(\\?.*)?$", filter.NoCacheFilter)
	mux.HFunc("^/cache2(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		cache2++
		return hs.MsgRes(fmt.Sprintf("%v", cache2))
	})
	mux.Handler("^/.*", http.FileServer(http.Dir("www")))
	s := http.Server{Addr: ":8080", Handler: mux}
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
