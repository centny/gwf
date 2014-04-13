package dbutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/Centny/Cny4go/test"
	"github.com/Centny/Cny4go/util"
	"github.com/Centny/TDb"
	_ "github.com/go-sql-driver/mysql"
	"testing"
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
	Add1   string    `m2s:"ADD1"`
	Add2   string    `m2s:"Add2"`
}

func TestDbUtil(t *testing.T) {
	db, _ := sql.Open("mysql", test.TDbCon)
	defer db.Close()
	err := DbExecF(db, "ttable.sql")
	if err != nil {
		t.Error(err.Error())
	}
	res, err := DbQuery(db, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(res) < 1 {
		t.Error("not data")
		return
	}
	if len(res[0]) < 1 {
		t.Error("data is empty")
		return
	}
	bys, err := json.Marshal(res)
	fmt.Println(string(bys))
	//
	var mres []TSt
	err = DbQueryS(db, &mres, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(mres) < 1 {
		t.Error("not data")
		return
	}
	fmt.Println("...", mres[0].T, util.Timestamp(mres[0].Time), util.Timestamp(time.Now()))
	fmt.Println(mres, mres[0].Add1)
	//
	ivs, err := DbQueryInt(db, "select * from ttable where tid")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(ivs) < 1 {
		t.Error("not data")
		return
	}
	//
	svs, err := DbQueryString(db, "select tname from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(svs) < 1 {
		t.Error("not data")
		return
	}
	//
	iid, err := DbInsert(db, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(iid)
	//
	tx, _ := db.Begin()
	iid2, err := DbInsert2(tx, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(iid2)
	tx.Commit()
	//
	erow, err := DbUpdate(db, "delete from ttable where tid=?", iid)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(erow)
	//
	tx, _ = db.Begin()
	erow, err = DbUpdate2(tx, "delete from ttable where tid=?", iid2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(erow)
	tx.Commit()
	//
	_, err = DbQuery(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryInt(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryString(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbInsert(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx, _ = db.Begin()
	_, err = DbInsert2(tx, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = DbUpdate(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = DbUpdate2(tx, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = DbQuery(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryInt(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryString(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbInsert(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = DbInsert2(tx, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = DbUpdate(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = DbUpdate2(tx, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	err = DbQueryS(nil, nil, "select * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	DbQueryInt(nil, "select * from ttable where tid>?", 1, 2)
	DbQueryString(nil, "select * from ttable where tid>?", 1, 2)
	DbInsert(nil, "select * from ttable where tid>?", 1, 2)
	DbUpdate(nil, "select * from ttable where tid>?", 1, 2)
	DbInsert2(nil, "select * from ttable where tid>?", 1, 2)
	DbUpdate2(nil, "select * from ttable where tid>?", 1, 2)
}
func Map2Val2(columns []string, row map[string]interface{}, dest []driver.Value) {
	for i, c := range columns {
		if v, ok := row[c]; ok {
			switch c {
			case "INT":
				dest[i] = int(v.(float64))
			case "UINT":
				dest[i] = uint32(v.(float64))
			case "FLOAT":
				dest[i] = float32(v.(float64))
			case "SLICE":
				dest[i] = []byte(v.(string))
			case "STRING":
				dest[i] = v.(string)
			case "STRUCT":
				dest[i] = time.Now()
			case "BOOL":
				dest[i] = true
			}
		} else {
			dest[i] = nil
		}
	}
}
func TestDbUtil2(t *testing.T) {
	TDb.Map2Val = Map2Val2
	db, _ := sql.Open("TDb", "td@tdata.json")
	defer db.Close()
	res, err := DbQuery(db, "SELECT * FROM TESTING WHERE INT=? AND STRING=?", 1, "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(res)
}

func TestDbExecF(t *testing.T) {
	db, _ := sql.Open("mysql", test.TDbCon)
	defer db.Close()
	err := DbExecF(db, "ttable.sql")
	if err != nil {
		t.Error(err.Error())
	}
	DbExecF(nil, "ttable.sql")
	DbExecF(db, "ttables.sql")
	db.Close()
	DbExecF(db, "ttable.sql")
}
