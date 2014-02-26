package test

const TDbCon string = "cny:123@tcp(127.0.0.1:3306)/cny?charset=utf8"

// func HTTPGetSrv(addr string) {
// 	var lck sync.WaitGroup
// 	lck.Add(1)
// 	http.HandleFunc("/w", func(w http.ResponseWriter, r *http.Request) {
// 		for {
// 			fmt.Println("writing...")
// 			w.Write([]byte("testing\n"))
// 			time.Sleep(time.Second)
// 		}
// 	})
// 	http.HandleFunc("/e", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("exit")
// 		lck.Done()
// 	})
// 	fmt.Println("listen server:", addr)
// 	go http.ListenAndServe(addr, nil)
// 	fmt.Println("waiting ...")
// 	lck.Wait()
// 	fmt.Println("end ....")
// }
