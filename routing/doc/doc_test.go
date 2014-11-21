package doc

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"net/http"
	"testing"
)

type Abcd struct {
}

func (a *Abcd) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	return routing.HRES_RETURN
}

func (a *Abcd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

type Abcd2 struct {
}

func (a *Abcd2) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	return routing.HRES_RETURN
}

func (a *Abcd2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
func (a *Abcd2) Doc() *Desc {
	return &Desc{}
}

var _ = Desc{
	Detail: "jjsdfsdf",
	SeeV:   []interface{}{TTN, &Abcd2{}},
}.Api(TT)

func TT(hs *routing.HTTPSession) routing.HResult {
	return routing.HRES_RETURN
}
func TT2(hs *routing.HTTPSession) routing.HResult {
	return routing.HRES_RETURN
}

var _ = Desc{
	Detail: "jjsdfsdf",
	SeeV:   []interface{}{TT},
}.Api(TTN)

func TTN(w http.ResponseWriter, r *http.Request) {

}
func TTN2(w http.ResponseWriter, r *http.Request) {

}
func TestTt(t *testing.T) {
	ts := httptest.NewMuxServer()
	ts.Mux.H("/abc.*", NewDocViewer())
	dv := NewDocViewer()
	dv.HTML = "{{sds}}"
	ts.Mux.H("/bgt_d.*", dv)
	ts.Mux.H("/abd.*", NewDocViewerInc(".*abccc01.*"))
	ts.Mux.H("/abe.*", NewDocViewerExc(".*abccc01.*"))
	//
	ts.Mux.H("/abccc01.*", &Abcd{})
	ts.Mux.Handler("/abccc02.*", &Abcd{})
	ts.Mux.H("/abccc03.*", &Abcd2{})
	ts.Mux.Handler("/abccc04.*", &Abcd2{})
	//
	ts.Mux.HFilter("/abccc052.*", &Abcd{})
	ts.Mux.HFilterFunc("/abccc062.*", TT)
	//
	ts.Mux.HFunc("/abccc05.*", TT)
	ts.Mux.HandleFunc("/abccc06.*", TTN)
	ts.Mux.HFunc("/abccc07.*", TT2)
	ts.Mux.HandleFunc("/abccc08.*", TTN2)
	//
	fmt.Println("-------->\n")
	fmt.Println(ts.G("/bgt_d"))
	fmt.Println(ts.G("/abc?fmt=json"))
	fmt.Println("-------->\n")
	fmt.Println(ts.G("/abd?fmt=json"))
	fmt.Println("-------->\n")
	fmt.Println(ts.G("/abe?fmt=json"))
	fmt.Println(ts.G("/abe?fmt=tttt"))
}

func TestErr(t *testing.T) {
	func() {
		defer func() {
			err := recover()
			fmt.Println(err)
		}()
		Desc{
			Detail: "jjsdfsdf",
		}.Api(TTN)
	}()
}
