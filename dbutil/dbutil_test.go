package dbutil

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Centny/Cny4go/test"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

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
}
