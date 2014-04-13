//Package igtest provide the integration testing framework.
package igtest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/util"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

//task for testing
type Task struct {
	Rely []string    //rely data key
	F    interface{} //run function
}

//the integration tesing builder
type Builder struct {
	In    string   //the rely json data file path for read
	Out   string   //the json data file path for write
	Dset  util.Map //the data set read from json file or store to json file.
	Kvs   util.Map //the extern key value to run the task
	tasks []Task   //all testing task.
}

//add data to data set by key,
//if val is nil,it will delete the key value.
func (b *Builder) Set(key string, val interface{}) *Builder {
	if val == nil {
		delete(b.Kvs, key)
	} else {
		b.Kvs[key] = val
	}
	log.D("adding the cache data:%v", key)
	return b
}

//store data to data set by key.
func (b *Builder) StoreData(key string, val interface{}) *Builder {
	b.Dset[key] = val
	log.D("adding the test data to set:%v", key)
	return b
}

//delete the data from data set by key
func (b *Builder) DeleteData(key string) *Builder {
	delete(b.Dset, key)
	log.D("remove the test data from set:%v", key)
	return b
}

//append task to the testing queue by task object.
func (b *Builder) Append(ts ...Task) *Builder {
	for _, t := range ts {
		b.tasks = append(b.tasks, t)
	}
	return b
}

//append task to the testing queue by rely key and function
func (b *Builder) Append2(rely string, f interface{}) *Builder {
	b.Append(Task{
		Rely: strings.Split(rely, ","),
		F:    f,
	})
	return b
}

//append not rely task to the testing queue by function.
func (b *Builder) Append3(f interface{}) *Builder {
	b.Append(Task{
		Rely: []string{},
		F:    f,
	})
	return b
}

//clear all task.
func (b *Builder) cls_task() {
	b.tasks = []Task{}
}

//clear all data.
func (b *Builder) cls_dset() {
	b.Dset = util.Map{}
}

//execute on testing task.
func (b *Builder) exec_f(task Task) error {
	fval := reflect.ValueOf(task.F)
	if fval.Kind() != reflect.Func { //check function kind
		err := errors.New(fmt.Sprintf("not func in task{Rely:%v,F:%v}", task.Rely, task.F))
		return err
	}
	ftype := fval.Type()
	if len(task.Rely) != ftype.NumIn() { //check the function in count
		err := errors.New(fmt.Sprintf(
			"func(%v) argument count(%v) is not equal rely count(%v)",
			ftype.String(), ftype.NumIn(), len(task.Rely),
		))
		return err
	}
	//covert all in value.
	vals := []reflect.Value{}
	for i, r := range task.Rely {
		atype := ftype.In(i)
		var data interface{} = nil
		if td, ok := b.Dset[r]; ok {
			data = td
		}
		if td, ok := b.Kvs[r]; ok {
			data = td
		}
		if data == nil {
			err := errors.New(fmt.Sprintf(
				"rely data(%v) not found for func(%v)",
				r, ftype.String(),
			))
			return err
		}
		dtype := reflect.TypeOf(data)
		if atype.Kind() == dtype.Kind() {
			vals = append(vals, reflect.ValueOf(data))
			continue
		}
		if dtype.Kind() == reflect.Float64 {
			switch atype.Kind() {
			case reflect.Int64:
				vals = append(vals, reflect.ValueOf(b.Dset.IntVal(r)))
				continue
			case reflect.Uint64:
				vals = append(vals, reflect.ValueOf(b.Dset.UintVal(r)))
				continue
			case reflect.String:
				vals = append(vals, reflect.ValueOf(b.Dset.StrVal(r)))
				continue
			}
		}
		if dtype.Kind().String() == "map" && atype.Kind() == reflect.Struct {
			tval := reflect.New(atype)
			util.M2S(util.Map(data.(map[string]interface{})), tval.Interface())
			vals = append(vals, reflect.Indirect(tval))
			continue
		}
		if atype.Kind() == reflect.Ptr {
			ptype := atype.Elem()
			tval := reflect.New(ptype)
			util.M2S(util.Map(data.(map[string]interface{})), tval.Interface())
			vals = append(vals, tval)
			continue
		}
		return errors.New(fmt.Sprintf(
			"rely data(%v) type(%v) not match func(%v) argument type(%v)",
			r, dtype.Kind(), ftype.String(), atype.Kind()))
	}
	fval.Call(vals)
	return nil
}

//execute all task in queue.
func (b *Builder) Exec() error {
	b.Set("@b", b)
	if len(b.In) > 0 {
		idata, err := ioutil.ReadFile(b.In)
		if err != nil {
			return err
		}
		err = json.Unmarshal(idata, &b.Dset)
		if err != nil {
			return err
		}
	}
	for _, ts := range b.tasks {
		b.Set("@t", ts)
		if err := b.exec_f(ts); err != nil {
			b.Set("@t", nil)
			return err
		}
	}
	b.Set("@t", nil)
	if len(b.Out) > 0 {
		odata, err := json.Marshal(b.Dset)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(b.Out, odata, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

//new build by in and out file path.
func NewBuilder(in string, out string) *Builder {
	b := &Builder{}
	b.Dset = util.Map{}
	b.Kvs = util.Map{}
	b.tasks = []Task{}
	b.In, b.Out = in, out
	for _, a := range os.Args {
		if strings.HasPrefix(a, "@ig-i=") {
			b.In = strings.TrimPrefix(a, "@ig-i=")
		} else if strings.HasPrefix(a, "@ig-o=") {
			b.Out = strings.TrimPrefix(a, "@ig-o=")
		}
	}
	return b
}

//new the not in and out file path builder.
func NewBuilder2() *Builder {
	return NewBuilder("", "")
}
