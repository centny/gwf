package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"io"
	"io/ioutil"
	"strings"
)

func GetArgs(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	var args = util.Map{}
	for key, _ := range hs.R.Form {
		args.SetVal(key, hs.RVal(key))
	}
	return hs.JRes(args)
}

func PostArgs(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	var args = util.Map{}
	for key, _ := range hs.R.PostForm {
		args.SetVal(key, hs.RVal(key))
	}
	return hs.JRes(args)
}

func MultipartArgs(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseMultipartForm(102400)
	var args = util.Map{}
	for key, _ := range hs.R.MultipartForm.Value {
		args.SetVal(key, hs.R.MultipartForm.Value[key][0])
	}
	return hs.JRes(args)
}

func SetSs(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	var args = util.Map{}
	for key, _ := range hs.R.Form {
		hs.SetVal(key, hs.RVal(key))
		args.SetVal(key, hs.RVal(key))
	}
	return hs.JRes(args)
}

func GetSs(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	var args = util.Map{}
	for _, key := range strings.Split(hs.RVal("keys"), ",") {
		args.SetVal(key, hs.Val(key))
	}
	return hs.JRes(args)
}

func Upload(hs *routing.HTTPSession) routing.HResult {
	var fn, md5, sha, size, err = hs.RecFv2("file", WWW)
	return hs.JRes(util.Map{
		"fn":   fn,
		"md5":  md5,
		"sha":  sha,
		"size": size,
		"err":  err,
	})
}

func Body(hs *routing.HTTPSession) routing.HResult {
	io.Copy(hs.W, hs.R.Body)
	return routing.HRES_RETURN
}

func ReqCType(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	var args = util.Map{}
	for key, _ := range hs.R.Header {
		args.SetVal(key, hs.R.Header.Get(key))
	}
	return hs.JRes(args)
}

func ResCType(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	var args = util.Map{}
	var h = hs.W.Header()
	for key, _ := range hs.R.Form {
		h.Set(key, hs.RVal(key))
		args.SetVal(key, hs.RVal(key))
	}
	return hs.JRes(args)
}

func Echo(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	fmt.Println(">Headers>>")
	for key, _ := range hs.R.Header {
		fmt.Println(key, ":", hs.R.Header.Get(key))
	}
	fmt.Println("\n")
	fmt.Println(">Get>>")
	for key, _ := range hs.R.Form {
		fmt.Println("  ", key, ":", hs.R.FormValue(key))
	}
	fmt.Println("\n")
	fmt.Println(">Post>>")
	for key, _ := range hs.R.PostForm {
		fmt.Println("  ", key, ":", hs.R.PostFormValue(key))
	}
	fmt.Println("\n")
	fmt.Println(">Body>>")
	bys, _ := ioutil.ReadAll(hs.R.Body)
	fmt.Println(string(bys))
	fmt.Println("\n\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Fprintf(hs.W, "OK")
	return routing.HRES_RETURN
}
