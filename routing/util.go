package routing

import (
	"fmt"
	"github.com/Centny/Cny4go/log"
	"io"
	"net/http"
	"os"
)

func SendF(w http.ResponseWriter, r *http.Request, fname, tfile, ctype string, attach bool) {
	src, err := os.OpenFile(tfile, os.O_RDONLY, 0)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fi, err := src.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fsize := fi.Size()
	header := w.Header()
	if len(ctype) < 1 {
		header.Set("Content-Type", "application/octet-stream")
	} else {
		header.Set("Content-Type", ctype)
	}
	if attach {
		header.Set("Content-Disposition", fmt.Sprintf("attachment; filename='%s'", fname))
	}
	header.Set("Content-Length", fmt.Sprintf("%v", fsize))
	header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	_, err = io.Copy(w, src)
	src.Close()
	if err != nil {
		log.E("sending file(%v) error:%s", tfile, err.Error())
	}
}
