package test

import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http"
)

//the x1 api
//
//@url,the request url
//	~/xx/api/x1?a=1&b=xx&c=n	GET
//@arg,the normal http query arguments:
//	a	required	the required arguments, using for ....
//	b	optional	the optional arguments, using for ....
//	c	R	the required arguments, using for ....
//	d	O	the optional arguments, using for ....
//
//	a=1&b=2&c=3&d=4
//@ret,the json return:
//	xx	S	the xx data field(string)
//	xa	I	the xa data field(int)
//	xb	F	the xb data field(float)
//	xc	A	the xb data field(array)
//	xd	O	the xb data field(object)
/*
	{
		"code": 0,
		"data": {
			"xa": "xx",
			"xb": "aa",
			"xx": 3
		}
	}
*/
//@tag,test,x
func X1(hs *routing.HTTPSession) routing.HResult {
	var a, b int
	var c, d string
	err := hs.ValidCheckVal(`
		a,R|I,R:0;
		b,O|I,R:0;
		c,R|S,L:0;
		d,O|S,L:0;
		`, &a, &b, &c, &d)
	if err == nil {
		return hs.MsgRes(util.Map{
			"xx": a + b,
			"xa": c,
			"xb": d,
		})
	} else {
		return hs.MsgResErr2(1, "arg-err", err)
	}
}

//
func X2(hs *routing.HTTPSession) (x routing.HResult) {
	return routing.HRES_RETURN
}

//command is empty
func X3(hs *routing.HTTPSession) (x routing.HResult) {
	return routing.HRES_RETURN
}

//command is empty
//@url,xxxx
func X4(hs *routing.HTTPSession) (x routing.HResult) {
	return routing.HRES_RETURN
}

//command is empty
//@url,xxxx
//	~/xx/api/x1?a=1&b=xx&c=n	GET	application/json
//@arg,xxx
func X5(hs *routing.HTTPSession) (x routing.HResult) {
	return routing.HRES_RETURN
}

func V1(w http.ResponseWriter, r *http.Request) {

}

type M struct {
}

func (m *M) X1(hs *routing.HTTPSession) routing.HResult {
	return routing.HRES_RETURN
}
func (m *M) V2(w http.ResponseWriter, r *http.Request) {
}

func A1(hs *routing.HTTPSession) {

}
func A2(hs *routing.HTTPSession) int {
	return 0
}

func A3() routing.HResult {
	return routing.HRES_RETURN
}

func A4(v int) routing.HResult {
	return routing.HRES_RETURN
}

func A5(hs *routing.HTTPSession) routing.HTTPSession {
	return *hs
}

func A6(hs *routing.HTTPSession) *routing.HTTPSession {
	return nil
}

func A7(hs *routing.HTTPSession) (x, a routing.HResult) {
	return routing.HRES_RETURN, routing.HRES_RETURN
}

func B1(w http.ResponseWriter, r http.Request) {
}

func B2(w *http.ResponseWriter, r *http.Request) {
}

func B3(w routing.HTTPSession, r *http.Request) {
}

func c1(hs *routing.HTTPSession) routing.HResult {
	return routing.HRES_RETURN
}
