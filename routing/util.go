package routing

import (
	"fmt"
	"github.com/Centny/Cny4go/log"
	"io"
	"net/http"
	"os"
)

func SendF(w http.ResponseWriter, r *http.Request, fname, tfile, ctype string, attach bool) {
	src, err := os.Open(tfile)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = sendf(w, src, fname, ctype, attach)
	if err != nil {
		log.E("sending file(%v) error:%s", tfile, err.Error())
		http.NotFound(w, r)
		return
	}
}
func sendf(w http.ResponseWriter, file *os.File, fname, ctype string, attach bool) error {
	fi, err := file.Stat()
	if err != nil {
		return err
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
	_, err = io.Copy(w, file)
	return err
}
