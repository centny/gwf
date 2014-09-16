package routing

import (
	"encoding/json"
	"github.com/Centny/Cny4go/log"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type JsonINT struct {
	Path    string
	Default string
	Local   string
	Kvs     map[string]map[string]string
	lock    sync.RWMutex
}

func (j *JsonINT) SetLocal(hs *HTTPSession, local string) {
	j.Local = local
}
func (j *JsonINT) LocalVal(hs *HTTPSession, key string) string {
	if len(j.Local) > 0 {
		val := j.LangVal(j.Local, key)
		if len(val) > 0 {
			return val
		}
	}
	als := hs.AcceptLanguages()
	if len(als) < 1 {
		return j.LangVal(j.Default, key)
	} else {
		return j.LangQesVal(als, key)
	}
}
func (j *JsonINT) LangQesVal(als LangQes, key string) string {
	for _, al := range als {
		val := j.LangVal(al.Lang, key)
		if len(val) > 0 {
			return val
		}
	}
	return j.LangVal(j.Default, key)
}
func (j *JsonINT) LangVal(lang string, key string) string {
	if _, ok := j.Kvs[lang]; !ok {
		j.lock.Lock()
		err := j.LoadJson(lang)
		if err != nil { //load error or not found,marked loaded.
			j.Kvs[lang] = map[string]string{}
		}
		j.lock.Unlock()
	}
	return j.Kvs[lang][key]
}
func (j *JsonINT) LoadJson(lang string) error {
	fpath := filepath.Join(j.Path, lang+".json")
	bys, err := ioutil.ReadFile(fpath)
	if err != nil {
		// fmt.Println(err.Error())
		return err
	}
	mv := map[string]string{}
	err = json.Unmarshal(bys, &mv)
	if err != nil {
		log.D("load the lang(%s) json file(%s) error:%s", lang, fpath, err.Error())
		return err
	}
	j.Kvs[lang] = mv
	return nil
}
func NewJsonINT(path string) (*JsonINT, error) {
	ji := &JsonINT{}
	ji.Kvs = map[string]map[string]string{}
	ji.lock = sync.RWMutex{}
	return ji, nil
}
