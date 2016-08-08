//Package dbutil provide some tool for database operation.
//Author:Centny
package dbutil

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Centny/gwf/util"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//convert the sql.Rows to map array.
func DbRow2Map(rows *sql.Rows) []util.Map {
	res := []util.Map{}
	fields, _ := rows.Columns()
	//fmt.Println(fields)
	fieldslen := len(fields)
	for rows.Next() {
		//
		//scan the data to array.
		sary := make([]interface{}, fieldslen) //scan array.
		for i := 0; i < fieldslen; i++ {
			var a interface{}
			sary[i] = &a
		}
		rows.Scan(sary...)
		//
		//convert array to map.
		mm := util.Map{}
		for idx, field := range fields {
			rawValue := reflect.Indirect(reflect.ValueOf(sary[idx]))
			if rawValue.Interface() == nil { //if database data is null.
				mm[field] = nil
				continue
			}
			aa := reflect.TypeOf(rawValue.Interface())
			vv := reflect.ValueOf(rawValue.Interface())
			switch aa.Kind() { //check the value type ant convert.
			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				mm[field] = vv.Int()
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				mm[field] = vv.Uint()
			case reflect.Float32:
				mm[field], _ = strconv.ParseFloat(fmt.Sprintf("%v", vv.Interface()), 64)
			case reflect.Float64:
				mm[field] = vv.Interface().(float64)
			case reflect.Slice:
				mm[field] = string(rawValue.Interface().([]byte))
			case reflect.String:
				mm[field] = vv.String()
			case reflect.Struct:
				mm[field] = rawValue.Interface().(time.Time)
			case reflect.Bool:
				mm[field] = vv.Bool()
			}
		}
		res = append(res, mm)
	}
	return res
}

//query the map result by query string and arguments.
func DbQuery(db *sql.DB, query string, args ...interface{}) ([]util.Map, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return DbRow2Map(rows), nil
}
func DbQuery2(tx *sql.Tx, query string, args ...interface{}) ([]util.Map, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return DbRow2Map(rows), nil
}

//query the struct result by query string and arguments.
func DbQueryS(db *sql.DB, res interface{}, query string, args ...interface{}) error {
	mres, err := DbQuery(db, query, args...)
	if err != nil {
		return err
	}
	util.Ms2Ss(mres, res)
	return nil
}
func DbQueryS2(tx *sql.Tx, res interface{}, query string, args ...interface{}) error {
	mres, err := DbQuery2(tx, query, args...)
	if err != nil {
		return err
	}
	util.Ms2Ss(mres, res)
	return nil
}
func DbQueryI(db *sql.DB, query string, args ...interface{}) (int64, error) {
	ic, err := DbQueryInt(db, query, args...)
	if err != nil {
		return 0, err
	}
	if len(ic) < 1 {
		return 0, util.NOT_FOUND
	} else {
		return ic[0], nil
	}
}
func DbQueryI2(tx *sql.Tx, query string, args ...interface{}) (int64, error) {
	ic, err := DbQueryInt2(tx, query, args...)
	if err != nil {
		return 0, err
	}
	if len(ic) < 1 {
		return 0, util.NOT_FOUND
	} else {
		return ic[0], nil
	}
}
func DbQueryF(db *sql.DB, query string, args ...interface{}) (float64, error) {
	ic, err := DbQueryFloat(db, query, args...)
	if err != nil {
		return 0, err
	}
	if len(ic) < 1 {
		return 0, util.NOT_FOUND
	} else {
		return ic[0], nil
	}
}
func DbQueryF2(tx *sql.Tx, query string, args ...interface{}) (float64, error) {
	ic, err := DbQueryFloat2(tx, query, args...)
	if err != nil {
		return 0, err
	}
	if len(ic) < 1 {
		return 0, util.NOT_FOUND
	} else {
		return ic[0], nil
	}
}
func DbQueryStr(db *sql.DB, query string, args ...interface{}) (string, error) {
	ic, err := DbQueryString(db, query, args...)
	if err != nil {
		return "", err
	}
	if len(ic) < 1 {
		return "", util.NOT_FOUND
	} else {
		return ic[0], nil
	}
}
func DbQueryStr2(tx *sql.Tx, query string, args ...interface{}) (string, error) {
	ic, err := DbQueryString2(tx, query, args...)
	if err != nil {
		return "", err
	}
	if len(ic) < 1 {
		return "", util.NOT_FOUND
	} else {
		return ic[0], nil
	}
}

//
func DbQueryInt(db *sql.DB, query string, args ...interface{}) ([]int64, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []int64{}
	for rows.Next() {
		var iv int64
		rows.Scan(&iv)
		rv = append(rv, iv)
	}
	return rv, nil
}
func DbQueryInt2(tx *sql.Tx, query string, args ...interface{}) ([]int64, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []int64{}
	for rows.Next() {
		var iv int64
		rows.Scan(&iv)
		rv = append(rv, iv)
	}
	return rv, nil
}

func DbQueryFloat(db *sql.DB, query string, args ...interface{}) ([]float64, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []float64{}
	for rows.Next() {
		var fv float64
		rows.Scan(&fv)
		rv = append(rv, fv)
	}
	return rv, nil
}
func DbQueryFloat2(tx *sql.Tx, query string, args ...interface{}) ([]float64, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []float64{}
	for rows.Next() {
		var fv float64
		rows.Scan(&fv)
		rv = append(rv, fv)
	}
	return rv, nil
}

//
func DbQueryString(db *sql.DB, query string, args ...interface{}) ([]string, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []string{}
	for rows.Next() {
		var sv string
		rows.Scan(&sv)
		rv = append(rv, sv)
	}
	return rv, nil
}
func DbQueryString2(tx *sql.Tx, query string, args ...interface{}) ([]string, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []string{}
	for rows.Next() {
		var sv string
		rows.Scan(&sv)
		rv = append(rv, sv)
	}
	return rv, nil
}

func DbQueryInterface(db *sql.DB, query string, args ...interface{}) ([]interface{}, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []interface{}{}
	for rows.Next() {
		var sv interface{}
		rows.Scan(&sv)
		rv = append(rv, sv)
	}
	return rv, nil
}

func DbQueryInterface2(tx *sql.Tx, query string, args ...interface{}) ([]interface{}, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rv := []interface{}{}
	for rows.Next() {
		var sv interface{}
		rows.Scan(&sv)
		rv = append(rv, sv)
	}
	return rv, nil
}

//
func DbInsert(db *sql.DB, query string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, errors.New("db is nil")
	}
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
func DbInsert2(db *sql.Tx, query string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, errors.New("db is nil")
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

//
func DbUpdate(db *sql.DB, query string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, errors.New("db is nil")
	}
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
func DbUpdate2(db *sql.Tx, query string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, errors.New("db is nil")
	}
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

//
func DbExecF(db *sql.DB, file string) error {
	if db == nil {
		return errors.New("db is nil")
	}
	fdata, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return DbExecScript(db, string(fdata))
}
func DbExecScript(db *sql.DB, script string) error {
	script = regexp.MustCompile("(?msU)/\\*.*\\*/\n?").ReplaceAllString(script, "")
	script = regexp.MustCompile("--.*\n?").ReplaceAllString(script, "")
	script = regexp.MustCompile("\n{2,}").ReplaceAllString(script, "\n")
	// fmt.Println(script)
	data := bufio.NewReader(bytes.NewBufferString(script))
	var stackf func(*bytes.Buffer, rune) error
	stackf = func(bw *bytes.Buffer, last rune) error {
		for {
			ru, _, err := data.ReadRune()
			if err != nil {
				return err
			}
			if ru == last {
				bw.WriteRune(ru)
				return nil
			}
			if ru == rune(';') && last == 0 {
				return nil
			}
			bw.WriteRune(ru)
			switch ru {
			case rune('"'):
				stackf(bw, rune('"'))
			case rune('\''):
				stackf(bw, rune('\''))
			}
		}
	}
	var terr error = nil
	for terr == nil {
		buf := bytes.NewBuffer(nil)
		terr = stackf(buf, 0)
		b := strings.Trim(buf.String(), " \t\n\r")
		if len(b) < 1 {
			continue
		}
		_, err := db.Exec(b)
		if err != nil {
			return errors.New(fmt.Sprintf("%v:%v", b, err.Error()))
		}
	}
	return nil
}
func DbExecF2(driver, dsn, file string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	return DbExecF(db, file)
}

func CovInvS(vals []string) string {
	nvals := []string{}
	for _, v := range vals {
		nvals = append(nvals, "'"+v+"'")
	}
	return strings.Join(nvals, ",")
}
func CovInvI(vals []int64) string {
	nvals := []string{}
	for _, v := range vals {
		nvals = append(nvals, fmt.Sprintf("%v", v))
	}
	return strings.Join(nvals, ",")
}

func DbStrJoin(val []string) string {
	var tval = []string{}
	for _, v := range val {
		tval = append(tval, fmt.Sprintf("\"%v\"", v))
	}
	return strings.Join(tval, ",")
}
