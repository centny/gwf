//Package hrv provide the http reverse connection to require inner netwok http server.
//
//it base netw and netw/impl and netw/rc package.
//
//it map http://<public server address>/<prefix>/<url path> to http://<inner server address>/<url path>
//
package hrv

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const HTML = `
<html>
<body>
<style type="text/css">
table{
	border-collapse:collapse;
	border-spacing:0;
	border-left:1px solid #888;
	border-top:1px solid #888;
	background:#efefef;
}
th,td{
	border-right:1px solid #888;
	border-bottom:1px solid #888;
	padding:5px;
}
th{font-weight:bold;background:#ccc;}
ul{
	list-style-type: none;
	padding:0;
	margin:0;
}
</style>
<table>
<tr><th>Clients</th><th>Args</th><th>Header</th></tr>
<tr>
<td>
<ul>
{{range $ak,$av:=.Cs}}
<li>
<a href="{{$av.Name}}">{{$av.Alias}}</a>
</li>
{{end}}
</ul>
</td>
<td>
<ul>
{{range $ak,$av:=.Args}}
<li>
{{$ak}}:&nbsp;&nbsp;{{$av}}
</li>
{{end}}
</ul>
</td>
<td>
<ul>
{{range $ak,$av:=.Header}}
<li>
{{$ak}}:&nbsp;&nbsp;{{$av}}
</li>
{{end}}
</ul>
</td>
</tr>
</table>
</body>
</html>
`

//http reverse client info struct
type Client struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

//http reverse server handler
type HrvS_H interface {
	netw.CCHandler
	OnLogin(token, name, alias string) error
}

//http reverse server
type HrvS struct {
	*rc.RC_Listener_m //
	//
	ShowLog bool                //show loggin.
	Cs      map[netw.Con]Client //all client
	cs_l    sync.RWMutex        //client lock.
	Pre     string              //doh handler prefix to trim url.
	H       HrvS_H              //handler
	FormMax int64               //the post form max memory
	Headers map[string]bool     //the header list to transfter to reverse client.
	T       *template.Template  //the home html template.
	Pattern []*regexp.Regexp    //the pattern to math reverse api.
	F       http.Handler        //the default file system handler.
	//
	Args  util.Map
	Head  util.Map
	args_ util.Map
	head_ util.Map
}

//new http reverse server by pool/listen address and convert functions.
func NewHrvS(bp *pool.BytePool, addr string, rcm *impl.RCM_S, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) *HrvS {
	sh := &HrvS{
		Cs:      map[netw.Con]Client{},
		FormMax: 102400,
		Headers: map[string]bool{},
		Args:    util.Map{},
		Head:    util.Map{},
		args_:   util.Map{},
		head_:   util.Map{},
	}
	lm := rc.NewRC_Listener_m(bp, addr, sh, rcm, v2b, b2v, na)
	sh.Handle(lm)
	sh.Parse(HTML)
	return sh
}

//new http reverse server by pool/listen address and HRV_V2B/HRV_B2V/Json_NAV
func NewHrvS_j(bp *pool.BytePool, addr string) *HrvS {
	rcm := impl.NewRCM_S_j()
	return NewHrvS(bp, addr, rcm, HRV_V2B, HRV_B2V, impl.Json_NAV)
}

func (h *HrvS) OnCmd(c netw.Cmd) int {
	if h.H == nil {
		return -1
	} else {
		return h.H.OnCmd(c)
	}
}
func (h *HrvS) OnConn(c netw.Con) bool {
	if h.H == nil {
		return true
	} else {
		return h.H.OnConn(c)
	}
}
func (h *HrvS) OnClose(c netw.Con) {
	h.delc(c)
	if h.H != nil {
		h.H.OnClose(c)
	}
}

//show log
func (h *HrvS) slog(f string, args ...interface{}) {
	if h.ShowLog {
		log.D_(1, f, args...)
	}
}

//inital the default file server.
func (h *HrvS) SetWww(www string) {
	h.F = http.FileServer(http.Dir(www))
}

//add pattern to match list.
func (h *HrvS) AddPattern(reg string) {
	h.Pattern = append(h.Pattern, regexp.MustCompile(reg))
}

//parse index html template.
func (h *HrvS) Parse(html string) error {
	t, err := template.New("HrvS").Parse(html)
	if err == nil {
		h.T = t
		return nil
	} else {
		log.E("err->%v", err.Error())
		return err
	}
}

//delete client
func (h *HrvS) delc(c netw.Con) {
	h.cs_l.Lock()
	defer h.cs_l.Unlock()
	delete(h.Cs, c)
}

//adding client
func (h *HrvS) addc(c netw.Con, token, name, alias string) {
	h.cs_l.Lock()
	defer h.cs_l.Unlock()
	h.Cs[c] = Client{
		Token: token,
		Name:  name,
		Alias: alias,
	}
}

//do on login.
func (h *HrvS) onlogin(token, name, alias string) error {
	if h.H == nil {
		return nil
	}
	return h.H.OnLogin(token, name, alias)
}

//login by remote command connection
func (h *HrvS) Login(rc *impl.RCM_Cmd) (interface{}, error) {
	var token, name, alias string
	err := rc.ValidF(`
		token,R|S,L:0;
		name,R|S,L:0;
		alias,O|S,L:0;
		`, &token, &name, &alias)
	if err != nil {
		log.E("HrvS Login valid args error:%v", err.Error())
		return nil, err
	}
	if len(alias) < 1 {
		alias = name
	}
	err = h.onlogin(token, name, alias)
	if err != nil {
		log.E("HrvS OnLogin by name(%v),token(%v),alias(%v) error->%v", name, token, alias, err.Error())
		return util.Map{
			"code": -1,
			"msg":  err.Error(),
		}, nil
	}
	if h.Exist(name) {
		log.E("HrvS OnLogin by name(%v),token(%v),alias(%v) error->%v", name, token, alias, "already login")
		return nil, util.Err("client already login by name(%v)", name)
	}
	h.AddC_rc(name, rc)
	h.addc(rc.BaseCon(), token, name, alias)
	rc.SetWait(true)
	log.D("HrvS->login success by name(%v),token(%v),alias(%v)", name, token, alias)
	return util.Map{
		"code": 0,
	}, nil
}
func (h *HrvS) HB(rc *impl.RCM_Cmd) (interface{}, error) {
	return util.Map{
		"d": "HB-S",
	}, nil
}

//handle http reverse require by normal http handler
func (h *HrvS) Doh(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	hs.R.ParseMultipartForm(h.FormMax)
	uri := hs.R.URL.Path
	h.slog("HrvS doh by prefix(%v),uri(%v)", h.Pre, uri)
	uri = strings.TrimPrefix(uri, h.Pre)
	uri = strings.TrimPrefix(uri, "/")
	qs := strings.SplitN(uri, "/", 2)
	//
	hargs := util.Map{}
	header := util.Map{}
	for k, v := range hs.R.Form {
		if strings.HasPrefix(k, "h:") {
			header[strings.TrimPrefix(k, "h:")] = v[0]
		} else {
			hargs[k] = v[0]
		}
	}
	for k, v := range hs.R.Header {
		tk := strings.ToUpper(k)
		if h.Headers[tk] {
			header[k] = v[0]
		}
	}
	for k, v := range h.Args {
		hargs[k] = v
	}
	for k, v := range h.Head {
		header[k] = v
	}
	if len(qs[0]) < 1 {
		h.args_ = hargs
		h.head_ = header
		h.T.Execute(hs.W, util.Map{
			"Args":   h.args_,
			"Header": h.head_,
			"Cs":     h.Cs,
		})
		return routing.HRES_RETURN
	}
	for k, v := range h.args_ {
		hargs[k] = v
	}
	for k, v := range h.head_ {
		header[k] = v
	}
	//
	//
	rurl := ""
	if len(qs) < 2 {
		rurl = ""
	} else {
		rurl = qs[1]
	}
	args := util.Map{}
	args["U"] = rurl
	args["M"] = hs.R.Method
	args["A"] = hargs
	args["H"] = header
	//
	cmd := h.CmdC(qs[0])
	if cmd == nil {
		log.E("HrvS doh error->client not found by name(%v)", qs[0])
		return hs.MsgResE3(1, "arg-err", "client not found by name("+qs[0]+")")
	}
	h.slog("HrvS doh to remote by url(%v)", rurl)
	for _, reg := range h.Pattern {
		if reg.MatchString(rurl) {
			return h.doh(hs, cmd, args)
		}
	}
	if h.F == nil {
		hs.W.WriteHeader(404)
		hs.W.Write([]byte("404"))
	} else {
		hs.R.URL.Path = rurl
		h.F.ServeHTTP(hs.W, hs.R)
	}
	return routing.HRES_RETURN

}
func (h *HrvS) doh(hs *routing.HTTPSession, cmd *impl.RCM_Con, args util.Map) routing.HResult {
	var res Res
	_, err := cmd.Exec("doh", args, &res)
	if err == nil {
		hs.W.WriteHeader(int(res.GetCode()))
		hs.W.Write(res.GetData())
		return routing.HRES_RETURN
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}

//handler all remote command
func (h *HrvS) Handle(l *rc.RC_Listener_m) {
	h.RC_Listener_m = l
	l.AddHFunc("login", h.Login)
	l.AddHFunc("hb", h.HB)
}

//http reverse client
//if token and name is not empty,it will auto login by token/name/alias
type HrvC struct {
	*rc.RC_Runner_m //client runner.
	//
	ShowLog bool           //show debug log.
	Base    string         //base url.
	H       netw.CCHandler // handler.
	//
	Token string //login token
	Name  string //login name
	Alias string //login alias
}

//new http reverse client by pool/server address/base url and convert functions.
func NewHrvC(bp *pool.BytePool, addr, base string, rcm *impl.RCM_S, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) *HrvC {
	ch := &HrvC{Base: base}
	cr := rc.NewRC_Runner_m(bp, addr, ch, rcm, v2b, b2v, na)
	ch.Handle(cr)
	return ch
}

//new http reverse client by pool/server address/base url/ and HRV_V2B/HRV_B2V/Json_NAV
func NewHrvC_j(bp *pool.BytePool, addr, base string) *HrvC {
	rcm := impl.NewRCM_S_j()
	return NewHrvC(bp, addr, base, rcm, HRV_V2B, HRV_B2V, impl.Json_NAV)
}

func (h *HrvC) OnCmd(c netw.Cmd) int {
	if h.H == nil {
		return -1
	} else {
		return h.H.OnCmd(c)
	}
}
func (h *HrvC) OnConn(c netw.Con) bool {
	h.RC_Runner_m.OnConn(c)
	c.SetWait(true)
	if len(h.Token) > 0 && len(h.Name) > 0 {
		go h.Login(h.Token, h.Name, h.Alias)
	}
	if h.H == nil {
		return true
	} else {
		return h.H.OnConn(c)
	}
}
func (h *HrvC) OnClose(c netw.Con) {
	h.RC_Runner_m.OnClose(c)
	if h.H != nil {
		h.H.OnClose(c)
	}
}

//show log
func (h *HrvC) slog(f string, args ...interface{}) {
	if h.ShowLog {
		log.D_(1, f, args...)
	}
}

//implement the reverse http require on client
func (h *HrvC) Doh(rc *impl.RCM_Cmd) (interface{}, error) {
	header := rc.MapVal("H")
	if header == nil {
		header = util.Map{}
	}
	header_ := map[string]string{}
	for k, v := range header {
		header_[k] = fmt.Sprintf("%v", v)
	}
	hargs := rc.MapVal("A")
	if hargs == nil {
		hargs = util.Map{}
	}

	hargs_ := map[string]string{}
	for k, v := range hargs {
		hargs_[k] = fmt.Sprintf("%v", v)
	}
	url_ := rc.StrVal("U")
	switch rc.StrVal("M") {
	case "GET":
		url_ = fmt.Sprintf("%v/%v?%v", h.Base, url_, util.QueryString(hargs_))
		h.slog("do GET(%v) by header(%v)", url_, header_)
		code, data_, err := util.HTTPClient.HGet_Hv(header_, url_)
		code_ := int32(code)
		return &Res{
			Code: &code_,
			Data: data_,
		}, err
	case "POST":
		h.slog("do POST(%v/%v) by header(%v),args(%v)", h.Base, url_, header_, hargs_)
		code, data_, err := util.HTTPClient.HPostF_Hv(
			fmt.Sprintf("%v/%v", h.Base, url_), hargs_, header_, "", "")
		code_ := int32(code)
		return &Res{
			Code: &code_,
			Data: data_,
		}, err
		return nil, nil
	default:
		log.D("do %v err->not supported", rc.StrVal("M"))
		return nil, util.Err("Method(%v) not supported", rc.StrVal("M"))
	}
}

//login to server by token,name,alias
func (h *HrvC) Login(token, name, alias string) error {
	log.I("do login by token(%v),name(%v),alias(%v)", token, name, alias)
	res, err := h.VExec_m("login", util.Map{
		"token": token,
		"name":  name,
		"alias": alias,
	})
	if err != nil {
		return err
	}
	if res.IntVal("code") == 0 {
		log.I("login by token(%v),name(%v),alias(%v) success", token, name, alias)
		return nil
	} else {
		log.E("login by token(%v),name(%v),alias(%v) error->code(%v),msg(%v)", token, name, alias, res.IntVal("code"), res.StrVal("msg"))
		return util.Err("error code(%v),msg(%v)", res.IntVal("code"), res.StrVal("msg"))
	}
}
func (h *HrvC) HB() error {
	_, err := h.VExec("hb", util.Map{}, &util.Map{})
	if err == nil {
		log.D("doing HB success")
	} else {
		log.W("doing HB err->%v", err.Error())
	}
	return err
}

//handler all client command
func (h *HrvC) Handle(r *rc.RC_Runner_m) {
	h.RC_Runner_m = r
	r.AddHFunc("doh", h.Doh)
}
