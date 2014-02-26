package test

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const TDbCon string = "cny:123@tcp(127.0.0.1:3306)/cny?charset=utf8"

func HTTPGetSrv(addr string) {
	var lck sync.WaitGroup
	lck.Add(1)
	http.HandleFunc("/w", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("writing...")
		for {
			w.Write([]byte("testing\n"))
			time.Sleep(time.Second)
		}
	})
	http.HandleFunc("/e", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("exit")
		lck.Done()
	})
	fmt.Println("listen server:", addr)
	go http.ListenAndServe(addr, nil)
	lck.Wait()
}
