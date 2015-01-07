package pool

// import (
// 	"fmt"
// 	"runtime"
// 	"syscall"
// 	"testing"
// 	"time"
// )

// func runtime_procPin() int {
// 	return 0
// }
// func runtime_procUnpin() {

// }

// func TestList(t *testing.T) {
// 	runtime.GOMAXPROCS(10)
// 	for i := 0; i < 1000; i++ {
// 		go func() {
// 			fmt.Println(syscall.RawSyscall(syscall.SYS_GETTID, 0, 0, 0))
// 		}()
// 	}
// 	time.Sleep(time.Second)
// }
