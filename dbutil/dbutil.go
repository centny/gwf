//Package dbutil provide some tool for database operation.
//Author:Centny
package dbutil

import (
	"database/sql"
	"errors"
	"github.com/Centny/Cny4go/util"
	"reflect"
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
			case reflect.Float32, reflect.Float64:
				mm[field] = vv.Float()
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
	stmt, err := db.Prepare(query)
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
