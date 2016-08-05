package routing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Centny/gwf/hooks"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type HResult int

const (
	HRES_CONTINUE HResult = iota
	HRES_RETURN
)
const (
	//the hook name
	HK_ROUTING = "ROUTING"

	//
	//filter begin,
	//the hook parameter
	// val:nil
	// args[]:
	//  *HTTTPSession	the HTTTPSession
	HK_R_BEG = "R_BEG"

	//
	//filter end,
	//the hook parameter
	// val: the HTTTPSession.V, it will be converted by SessionMux.FIND_V returned function.
	// args[]:
	//  *HTTTPSession	the HTTTPSession
	//  matched(bool)	if having filter matched.
	HK_R_END = "R_END"

	//
	//filter begin,
	//the hook parameter
	// val:nil
	// args[]:
	//  *HTTTPSession	the HTTTPSession
	HK_F_BEG = "F_BEG"

	//
	//filter end,
	//the hook parameter
	// val:nil
	// args[]:
	//  *HTTTPSession	the HTTTPSession
	//  matched(bool)	if having filter matched.
	//  HResult			the execute result.
	HK_F_END = "F_END" //filter end

	//
	//handler begin,
	//the hook parameter
	// val:nil
	// args[]:
	//  *HTTTPSession	the HTTTPSession
	HK_H_BEG = "H_BEG"

	//
	//handler end,
	//the hook parameter
	// val:nil
	// args[]:
	//  *HTTTPSession	the HTTTPSession
	//  matched(bool)	if having filter matched.
	//  HResult			the execute result.
	HK_H_END = "H_END"
)

func (h HResult) String() string {
	if h == HRES_CONTINUE {
		return "CONTINUE"
	} else {
		return "RETURN"
	}
}

type SessionEvHFunc func(string, Session)

func (h SessionEvHFunc) OnCreate(s Session) {
	h("CREATE", s)
}
func (h SessionEvHFunc) OnTimeout(s Session) {
	h("TIMEOUT", s)
}

type SessionEvHandler interface {
	OnCreate(s Session)
	OnTimeout(s Session)
}

type SessionBuilder interface {
	FindSession(w http.ResponseWriter, r *http.Request) Session
	SetEvH(h SessionEvHandler)
}
type Session interface {
	Val(key string) interface{}
	Set(key string, val interface{})
	Flush() error
}
type HandleFunc func(*HTTPSession) HResult
type Handler interface {
	SrvHTTP(*HTTPSession) HResult
}
type HandleFuncv func(v util.Validable) (interface{}, error)

func (h HandleFuncv) SrvHTTP(hs *HTTPSession) HResult {
	v, err := h(hs)
	if err == nil {
		return hs.MsgRes(v)
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}

type International interface {
	SetLocal(hs *HTTPSession, local string)
	LocalVal(hs *HTTPSession, key string) string
}

type HTTPSession struct {
	W   http.ResponseWriter
	R   *http.Request
	S   Session
	Mux *SessionMux
	Kvs map[string]interface{}
	INT International
	V   interface{} //response value.
}

func (h *HTTPSession) SetCookie(key string, val string) {
	cookie := &http.Cookie{}
	cookie.Name = key
	cookie.Domain = h.Mux.Domain
	cookie.Path = h.Mux.Path
	cookie.Value = val
	cookie.MaxAge = 0
	if len(val) < 1 {
		cookie.Expires = util.Time(0)
	}
	http.SetCookie(h.W, cookie)
}

func (h *HTTPSession) Cookie(key string) string {
	c, err := h.R.Cookie(key)
	if c == nil || err != nil {
		return ""
	}
	return c.Value
}

/* Redirect */
func (h *HTTPSession) Redirect(url string) {
	http.Redirect(h.W, h.R, url, http.StatusTemporaryRedirect)
}

func (h *HTTPSession) SetVal(key string, val interface{}) {
	h.S.Set(key, val)
}

func (h *HTTPSession) Flush() error {
	return h.S.Flush()
}

func (h *HTTPSession) Val(key string) interface{} {
	return h.S.Val(key)
}
func (h *HTTPSession) UintVal(key string) uint64 {
	v := h.S.Val(key)
	if v == nil {
		return 0
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.Uint:
		return uint64(v.(uint))
	case reflect.Uint8:
		return uint64(v.(uint8))
	case reflect.Uint16:
		return uint64(v.(uint16))
	case reflect.Uint32:
		return uint64(v.(uint32))
	case reflect.Uint64:
		return v.(uint64)
	default:
		return uint64(h.IntVal(key))
	}
}
func (h *HTTPSession) IntVal(key string) int64 {
	v := h.S.Val(key)
	if v == nil {
		return 0
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Int64:
		return v.(int64)
	default:
		return int64(h.FloatVal(key))
	}
}
func (h *HTTPSession) FloatVal(key string) float64 {
	v := h.S.Val(key)
	if v == nil {
		return 0
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.Float32:
		return float64(v.(float32))
	case reflect.Float64:
		return v.(float64)
	default:
		return 0
	}
}
func (h *HTTPSession) StrVal(key string) string {
	v := h.S.Val(key)
	if v == nil {
		return ""
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.String:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}
func (h *HTTPSession) JsonVal(key string) (util.Map, error) {
	json := h.CheckVal(key)
	if len(json) < 1 {
		return nil, util.NewNotFound("json valus not found by key(%v)", key)
	}
	json = strings.Trim(json, " \t")
	return util.Json2Map(json)
}

//converting json string to struct by util.J2S
//struct tags is m2s,not json
func (h *HTTPSession) JsonObjVal(key string, v interface{}) error {
	jdata := h.CheckVal(key)
	if len(jdata) < 1 {
		return util.NewNotFound("json valus not found by key(%v)", key)
	}
	return util.J2S(jdata, v)
}

//converting json string to struct by json.
func (h *HTTPSession) JsonObjVal2(key string, v interface{}) error {
	jdata := h.CheckVal(key)
	if len(jdata) < 1 {
		return util.NewNotFound("json valus not found by key(%v)", key)
	}
	return json.Unmarshal([]byte(jdata), v)
}

//
func (h *HTTPSession) CheckVal(key string) string {
	v := h.RVal(key)
	if len(v) > 0 {
		return v
	}
	return h.StrVal(key)
}

//check all value order by request,session,cookie.
func (h *HTTPSession) CheckValA(key string) string {
	v := h.CheckVal(key)
	if len(v) > 0 {
		return v
	}
	return h.Cookie(key)
}
func (h *HTTPSession) RVal(key string) string {
	v := h.R.FormValue(key)
	if len(v) > 0 {
		return v
	}
	v = h.R.PostFormValue(key)
	return v
}
func (h *HTTPSession) FormFInfo(name string) (int64, string, error) {
	src, fh, err := h.R.FormFile(name)
	if err != nil {
		return 0, "", err
	}
	err = errors.New("file size error")
	fsize := util.FormFSzie(src)
	if fsize > 0 {
		err = nil
	}
	return fsize, fh.Filename, err
}
func (h *HTTPSession) SendF(fname, tfile, ctype string, attach bool) {
	SendF(h.W, h.R, fname, tfile, ctype, attach)
}
func (h *HTTPSession) SendF2(tfile string) error {
	tf, err := os.Open(tfile)
	if err != nil {
		return err
	}
	defer tf.Close()
	_, err = io.Copy(h.W, tf)
	return err
}

//sending string by target context type.
func (h *HTTPSession) SendT(data string, ctype string) {
	header := h.W.Header()
	header.Set("Content-Type", ctype)
	header.Set("Content-Length", fmt.Sprintf("%v", len(data)))
	// header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	h.W.Write([]byte(data))
}
func (h *HTTPSession) SendT2(data string) {
	h.SendT(data, "text/plain")
}

//valid require value by format,limit require.
func (h *HTTPSession) ValidRVal(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.RVal, true, args...)
}

//valid require value by format,not limit require.
func (h *HTTPSession) ValidRValN(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.RVal, false, args...)
}

//valid all value by format,limit require.
func (h *HTTPSession) ValidCheckVal(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.CheckVal, true, args...)
}

//valid all value by format,not limit require.
func (h *HTTPSession) ValidCheckValN(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.CheckVal, false, args...)
}

//the same as ValidCheckVal,for impl util.Validable
func (h *HTTPSession) ValidF(f string, args ...interface{}) error {
	return h.ValidCheckVal(f, args...)
}

func (h *HTTPSession) AllRVal() util.Map {
	h.R.ParseForm()
	kvs := util.Map{}
	for k, v := range h.R.Form {
		kvs[k] = v
	}
	for k, v := range h.R.PostForm {
		kvs[k] = v
	}
	return kvs
}
func (h *HTTPSession) ParseQuery() error {
	vals, err := url.ParseQuery(h.R.URL.RawQuery)
	if err == nil {
		h.R.Form = vals
		h.R.PostForm = vals
	}
	return err
}
func http_res(code int, data interface{}, msg string, dmsg string) util.Map {
	res := make(util.Map)
	res["code"] = code
	if len(msg) > 0 {
		res["msg"] = msg
	}
	if data != nil {
		res["data"] = data
	}
	if len(dmsg) > 0 {
		res["dmsg"] = dmsg
	}
	return res
}
func http_res_ext(code int, data interface{}, msg string, dmsg string, ext interface{}, pa interface{}) util.Map {
	res := make(util.Map)
	res["code"] = code
	if len(msg) > 0 {
		res["msg"] = msg
	}
	if data != nil {
		res["data"] = data
	}
	if len(dmsg) > 0 {
		res["dmsg"] = dmsg
	}
	if ext != nil {
		res["ext"] = ext
	}
	if pa != nil {
		res["pa"] = pa
	}
	return res
}

// func json_res(code int, data interface{}, msg string, dmsg string) []byte {
// 	res := http_res(code, data, msg, dmsg)
// 	dbys, _ := json.Marshal(res)
// 	return dbys
// }
func (h *HTTPSession) JRes(data interface{}) HResult {
	err := h.JsonRes(data)
	if err != nil {
		log.E("JRes convert value(%v) to json err:%v", data, err)
	}
	return HRES_RETURN
}
func (h *HTTPSession) JsonRes(data interface{}) error {
	h.V = data
	h.W.Header().Set("Content-Type", "application/json;charset=utf-8")
	dbys, err := json.Marshal(data)
	if err != nil {
		return err
	}
	h.W.Write(dbys)
	return nil
}
func (h *HTTPSession) MsgResF(code int, data interface{}) HResult {
	h.V = data
	h.W.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprintf(h.W, `{"code":%v,"data":%v}`, code, data)
	return HRES_RETURN
}
func (h *HTTPSession) MsgRes(data interface{}) HResult {
	return h.JRes(http_res(0, data, "", ""))
}
func (h *HTTPSession) MsgResP(data interface{}, pn, ps, total int64) HResult {
	return h.JRes(http_res_ext(0, data, "", "", nil, map[string]int64{
		"pn":    pn,
		"ps":    ps,
		"total": total,
	}))
}
func (h *HTTPSession) MsgResExt(data interface{}, ext interface{}) HResult {
	return h.JRes(http_res_ext(0, data, "", "", ext, nil))
}
func (h *HTTPSession) MsgRes2(code int, data interface{}) HResult {
	return h.JRes(http_res(code, data, "", ""))
}

func (h *HTTPSession) MsgResE(code int, msg string) HResult {
	return h.JRes(http_res(code, nil, msg, ""))
}
func (h *HTTPSession) MsgResE2(code int, msg string, dmsg string) HResult {
	return h.JRes(http_res(code, nil, msg, dmsg))
}
func (h *HTTPSession) MsgResE3(code int, key string, dmsg string) HResult {
	return h.JRes(http_res(code, nil, h.LocalVal(key), dmsg))
}
func (h *HTTPSession) MsgResErr(code int, msg string, err error) HResult {
	return h.JRes(http_res(code, nil, msg, err.Error()))
}

//using the local value by key for error message.
func (h *HTTPSession) MsgResErr2(code int, key string, err error) HResult {
	return h.JRes(http_res(code, nil, h.LocalVal(key), err.Error()))
}

func (h *HTTPSession) Printf(format string, args ...interface{}) HResult {
	fmt.Fprintf(h.W, format, args...)
	return HRES_RETURN
}

/* International */
func (h *HTTPSession) SetLocal(local string) {
	if h.INT != nil {
		h.INT.SetLocal(h, local)
	}
}
func (h *HTTPSession) LocalVal(key string) string {
	if h.INT != nil {
		return h.INT.LocalVal(h, key)
	} else {
		return ""
	}
}

/* Host */
func (h *HTTPSession) Host() string {
	return h.R.Host
}

/* --------------- Access-Language --------------- */
type LangQ struct {
	Lang string
	Q    float64
}
type LangQes []LangQ

func (l LangQes) Len() int           { return len(l) }
func (l LangQes) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l LangQes) Less(i, j int) bool { return l[i].Q > l[j].Q }

func (h *HTTPSession) AcceptLanguages() LangQes {
	if len(h.R.Header["Accept-Language"]) < 1 {
		return LangQes{}
	}
	lstr := h.R.Header["Accept-Language"][0]
	var als LangQes = LangQes{} //all access languages.
	regexp.MustCompile("[^;]*;q?[^,]*").ReplaceAllStringFunc(lstr, func(src string) string {
		src = strings.Trim(src, "\t \n,")
		lq := strings.Split(src, ";")
		qua, err := strconv.ParseFloat(strings.Replace(lq[1], "q=", "", -1), 64)
		if err != nil {
			log.D("invalid Accept-Language q:%s", src)
			return src
		}
		for _, lan := range strings.Split(lq[0], ",") {
			als = append(als, LangQ{
				Lang: lan,
				Q:    qua,
			})
		}
		return src
	})
	sort.Sort(als)
	return als
}

/* --------------- Access-Language --------------- */

type SessionMux struct {
	Pre    string
	Domain string
	Path   string
	//
	Sb SessionBuilder
	//
	Filters      map[*regexp.Regexp]Handler
	FilterFunc   map[*regexp.Regexp]HandleFunc
	Handlers     map[*regexp.Regexp]Handler
	HandlerFunc  map[*regexp.Regexp]HandleFunc
	NHandlers    map[*regexp.Regexp]http.Handler
	NHandlerFunc map[*regexp.Regexp]http.HandlerFunc
	regex_f_ary  []*regexp.Regexp
	regex_f      map[*regexp.Regexp]int
	regex_h_ary  []*regexp.Regexp
	regex_h      map[*regexp.Regexp]int
	regex_m      map[*regexp.Regexp]string
	rs_m         map[*http.Request]*HTTPSession //request to session
	rs_l         sync.RWMutex
	Kvs          map[string]interface{}
	FilterEnable bool
	HandleEnable bool
	ShowLog      bool
	INT          International
	ShowSlow     int64
	M            *tutil.Monitor
	//provide the convert function to convert HTTPSesion.V as the hook HK_R_END value argument.
	FIND_V func(hs *HTTPSession) func(v interface{}) interface{}
}

func NewSessionMux2(pre string) *SessionMux {
	return NewSessionMux(pre, NewDefaultSessionBuilder())
}

func NewSessionMux(pre string, sb SessionBuilder) *SessionMux {
	if sb == nil {
		fmt.Println("session builder is nil")
		return nil
	}
	mux := SessionMux{}
	mux.Pre = pre
	mux.Domain = ""
	mux.Path = "/"
	mux.Sb = sb
	mux.Filters = map[*regexp.Regexp]Handler{}
	mux.Handlers = map[*regexp.Regexp]Handler{}
	mux.NHandlers = map[*regexp.Regexp]http.Handler{}
	mux.FilterFunc = map[*regexp.Regexp]HandleFunc{}
	mux.HandlerFunc = map[*regexp.Regexp]HandleFunc{}
	mux.NHandlerFunc = map[*regexp.Regexp]http.HandlerFunc{}
	mux.regex_f = map[*regexp.Regexp]int{}
	mux.regex_f_ary = []*regexp.Regexp{}
	mux.regex_h = map[*regexp.Regexp]int{}
	mux.regex_h_ary = []*regexp.Regexp{}
	mux.regex_m = map[*regexp.Regexp]string{}
	mux.rs_m = map[*http.Request]*HTTPSession{}
	mux.Kvs = map[string]interface{}{}
	mux.FilterEnable = true
	mux.HandleEnable = true
	mux.ShowLog = false
	mux.INT = nil
	mux.M = nil
	mux.ShowSlow = 10000
	return &mux
}

func (s *SessionMux) RSession(r *http.Request) *HTTPSession {
	s.rs_l.RLock()
	defer s.rs_l.RUnlock()
	if v, ok := s.rs_m[r]; ok {
		return v
	} else {
		return nil
	}
}
func (s *SessionMux) HFilter(pattern string, h Handler) {
	s.HFilterM(pattern, h, "*")
}
func (s *SessionMux) HFilterM(pattern string, h Handler, m string) {
	reg := regexp.MustCompile(pattern)
	s.Filters[reg] = h
	s.regex_f[reg] = 1
	s.regex_f_ary = append(s.regex_f_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) HFilterFunc(pattern string, h HandleFunc) {
	s.HFilterFuncM(pattern, h, "*")
}
func (s *SessionMux) HFilterFuncM(pattern string, h HandleFunc, m string) {
	reg := regexp.MustCompile(pattern)
	s.FilterFunc[reg] = h
	s.regex_f[reg] = 2
	s.regex_f_ary = append(s.regex_f_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) H(pattern string, h Handler) {
	s.HM(pattern, h, "*")
}
func (s *SessionMux) HM(pattern string, h Handler, m string) {
	reg := regexp.MustCompile(pattern)
	s.Handlers[reg] = h
	s.regex_h[reg] = 1
	s.regex_h_ary = append(s.regex_h_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) HFunc(pattern string, h HandleFunc) {
	s.HFuncM(pattern, h, "*")
}
func (s *SessionMux) HFuncv(pattern string, h HandleFuncv) {
	s.H(pattern, Handler(h))
}
func (s *SessionMux) HFuncM(pattern string, h HandleFunc, m string) {
	reg := regexp.MustCompile(pattern)
	s.HandlerFunc[reg] = h
	s.regex_h[reg] = 2
	s.regex_h_ary = append(s.regex_h_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) Handler(pattern string, h http.Handler) {
	s.HandlerM(pattern, h, "*", true)
}
func (s *SessionMux) HandlerM(pattern string, h http.Handler, m string, ret bool) {
	reg := regexp.MustCompile(pattern)
	s.NHandlers[reg] = h
	s.regex_h[reg] = 3
	s.regex_h_ary = append(s.regex_h_ary, reg)
	if ret {
		m = fmt.Sprintf("%s,:RETURN", m)
	} else {
		m = fmt.Sprintf("%s,:CONTINUE", m)
	}
	s.regex_m[reg] = m
}
func (s *SessionMux) HandleFunc(pattern string, h http.HandlerFunc) {
	s.HandleFuncM(pattern, h, "*", true)
}
func (s *SessionMux) HandleFuncM(pattern string, h http.HandlerFunc, m string, ret bool) {
	reg := regexp.MustCompile(pattern)
	s.NHandlerFunc[reg] = h
	s.regex_h[reg] = 4
	s.regex_h_ary = append(s.regex_h_ary, reg)
	if ret {
		m = fmt.Sprintf("%s,:RETURN", m)
	} else {
		m = fmt.Sprintf("%s,:CONTINUE", m)
	}
	s.regex_m[reg] = m
}

func (s *SessionMux) slog(fmt string, args ...interface{}) {
	if s.ShowLog {
		log.D(fmt, args...)
	}
}
func (s *SessionMux) check_m(reg *regexp.Regexp, m string) bool {
	if tm, ok := s.regex_m[reg]; ok {
		return strings.Contains(tm, "*") || strings.Contains(tm, m)
	}
	return false
}

func (s *SessionMux) check_continue(reg *regexp.Regexp) bool {
	if tm, ok := s.regex_m[reg]; ok {
		return strings.Contains(tm, ":CONTINUE")
	}
	return false
}

func (s *SessionMux) exec_f(hs *HTTPSession) (bool, HResult) {
	url := hs.R.URL.Path
	var matched bool = false
	for _, k := range s.regex_f_ary {
		if !k.MatchString(url) {
			continue
		}
		if !s.check_m(k, hs.R.Method) {
			s.slog("not mathced method %v to %v", hs.R.Method, s.regex_m[k])
			continue
		}
		var mid = ""
		matched = true
		switch s.regex_f[k] {
		case 1:
			if s.M != nil {
				mid = s.M.Start(fmt.Sprintf("F_%v", k.String()))
			}
			rv := s.Filters[k]
			res := rv.SrvHTTP(hs)
			if s.M != nil {
				s.M.Done(mid)
			}
			s.slog("mathced filter %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		case 2:
			if s.M != nil {
				mid = s.M.Start(fmt.Sprintf("F_%v", k.String()))
			}
			rv := s.FilterFunc[k]
			res := rv(hs)
			if s.M != nil {
				s.M.Done(mid)
			}
			s.slog("mathced filter func %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		}
	}
	return matched, HRES_CONTINUE
}
func (s *SessionMux) exec_h(hs *HTTPSession) (bool, HResult) {
	url := hs.R.URL.Path
	var matched bool = false
	for _, k := range s.regex_h_ary {
		if !k.MatchString(url) {
			continue
		}
		if !s.check_m(k, hs.R.Method) {
			s.slog("not mathced method %v to %v", hs.R.Method, s.regex_m[k])
			continue
		}
		var mid = ""
		matched = true
		switch s.regex_h[k] {
		case 1:
			if s.M != nil {
				mid = s.M.Start(fmt.Sprintf("H_%v", k.String()))
			}
			rv := s.Handlers[k]
			res := rv.SrvHTTP(hs)
			if s.M != nil {
				s.M.Done(mid)
			}
			s.slog("mathced handler %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		case 2:
			if s.M != nil {
				mid = s.M.Start(fmt.Sprintf("H_%v", k.String()))
			}
			rv := s.HandlerFunc[k]
			res := rv(hs)
			if s.M != nil {
				s.M.Done(mid)
			}
			s.slog("mathced handler func %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		case 3:
			if s.M != nil {
				mid = s.M.Start(fmt.Sprintf("H_%v", k.String()))
			}
			rv := s.NHandlers[k]
			rv.ServeHTTP(hs.W, hs.R)
			if s.M != nil {
				s.M.Done(mid)
			}
			if s.check_continue(k) {
				s.slog("mathced normal handler %v to %v (%v)",
					k, hs.R.URL.Path, HRES_CONTINUE.String())
				continue
			} else {
				s.slog("mathced normal handler %v to %v (%v)",
					k, hs.R.URL.Path, HRES_RETURN.String())
				return matched, HRES_RETURN
			}
		case 4:
			if s.M != nil {
				mid = s.M.Start(fmt.Sprintf("H_%v", k.String()))
			}
			rv := s.NHandlerFunc[k]
			rv(hs.W, hs.R)
			if s.M != nil {
				s.M.Done(mid)
			}
			if s.check_continue(k) {
				s.slog("mathced normal handler func %v to %v (%v)",
					k, hs.R.URL.Path, HRES_CONTINUE.String())
				continue
			} else {
				s.slog("mathced normal handler func %v to %v (%v)", k,
					hs.R.URL.Path, HRES_RETURN.String())
				return matched, HRES_RETURN
			}
		}
	}
	return matched, HRES_CONTINUE
}

//
func (s *SessionMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	beg := util.Now()
	r.URL.Path = strings.TrimPrefix(r.URL.Path, s.Pre)
	session := s.Sb.FindSession(w, r)
	hs := &HTTPSession{
		W:   w,
		R:   r,
		S:   session,
		Mux: s,
		Kvs: map[string]interface{}{},
		INT: s.INT,
	}
	s.rs_l.Lock()
	s.rs_m[r] = hs
	s.rs_l.Unlock()
	defer func() {
		s.rs_l.Lock()
		delete(s.rs_m, r) //remove the http session object.
		s.rs_l.Unlock()
		used := util.Now() - beg
		if s.ShowSlow > 0 && used > s.ShowSlow {
			log.W("SessionMux slow request found->%v", r.URL.String())
		}
	}()
	//
	var matched bool = false
	//
	defer func() {
		if !matched { //if not matched
			s.slog("not matchd any filter:%s", r.URL.Path)
			http.NotFound(w, r)
		}
		var tv interface{} = hs.V
		if s.FIND_V != nil {
			if fv := s.FIND_V(hs); fv != nil {
				tv = fv(hs.V)
			}
		}
		hooks.Call(HK_ROUTING, HK_R_END, tv, hs, matched)
	}()
	hooks.Call(HK_ROUTING, HK_R_BEG, nil, hs)
	//match filter.
	if s.FilterEnable {
		hooks.Call(HK_ROUTING, HK_F_BEG, nil, hs)
		mrv, res := s.exec_f(hs)
		matched = mrv
		hooks.Call(HK_ROUTING, HK_F_END, nil, hs, mrv, res)
		if res == HRES_RETURN {
			return
		}
	}
	//match handle
	if s.HandleEnable {
		hooks.Call(HK_ROUTING, HK_H_BEG, nil, hs)
		mrv, res := s.exec_h(hs)
		matched = matched || mrv
		hooks.Call(HK_ROUTING, HK_H_END, nil, hs, mrv, res)
	}
}

func (s *SessionMux) Print() {
	if len(s.Filters) > 0 {
		fmt.Println(" >Filters---->")
		for reg, h := range s.Filters {
			fmt.Printf("\t%v->%p\n", reg.String(), h)
		}
	}
	if len(s.FilterFunc) > 0 {
		fmt.Println(" >FilterFunc---->")
		for reg, h := range s.FilterFunc {
			fmt.Printf("\t%v->%p\n", reg.String(), h)
		}
	}
	if len(s.Handlers) > 0 {
		fmt.Println(" >Handlers---->")
		for reg, h := range s.Handlers {
			fmt.Printf("\t%v->%p\n", reg.String(), h)
		}
	}
	if len(s.HandlerFunc) > 0 {
		fmt.Println(" >HandlerFunc---->")
		for reg, h := range s.HandlerFunc {
			fmt.Printf("\t%v->%p\n", reg.String(), h)
		}
	}
	if len(s.NHandlers) > 0 {
		fmt.Println(" >NHandlers---->")
		for reg, h := range s.NHandlers {
			fmt.Printf("\t%v->%p\n", reg.String(), h)
		}
	}
	if len(s.NHandlerFunc) > 0 {
		fmt.Println(" >NHandlerFunc---->")
		for reg, h := range s.NHandlerFunc {
			fmt.Printf("\t%v->%p\n", reg.String(), h)
		}
	}
}

//set the show slow log
func (s *SessionMux) SetShowSlow(v int64) {
	s.ShowSlow = v
}

//start monitor
func (s *SessionMux) StartMonitor() {
	s.M = tutil.NewMonitor()
}

func (s *SessionMux) State() (interface{}, error) {
	if s.M == nil {
		return nil, nil
	} else {
		return s.M.State()
	}
}
