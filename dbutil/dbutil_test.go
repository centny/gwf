package dbutil

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Centny/Cny4go/test"
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
	Add1   string    `m2s:"ADD1"`
	Add2   string    `m2s:"Add2"`
}

func TestDbUtil(t *testing.T) {
	db, _ := sql.Open("mysql", test.TDbCon)
	defer db.Close()
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
	fmt.Println(mres, mres[0].Add1)
	//

	_, err = DbQuery(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQuery(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	err = DbQueryS(nil, nil, "select * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
}
