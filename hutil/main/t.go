package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/simlab/", handler)
	err := http.ListenAndServe(":8800", nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "http://lab.szjhjhb.com/simlab/labUserLogin.do?method=login&loginName=dyuser&passwd=dy123456", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	for _, c := range resp.Cookies() {
		w.Header().Add("Set-Cookie", c.Raw)
	}
	w.WriteHeader(resp.StatusCode)
	result, _ := ioutil.ReadAll(resp.Body)
	w.Write(result)
}

// func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.URL.Path)
// 		log.Println(r.Host)
// 		r.Host = "lab.szjhjhb.com"
// 		p.ServeHTTP(w, r)
// 	}
// }
