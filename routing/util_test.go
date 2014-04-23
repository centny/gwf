package routing

import (
	"errors"
	"fmt"
	"github.com/Centny/Cny4go/util"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type Tw struct {
	*os.File
}

func (t *Tw) Stat() (os.FileInfo, error) {
	return nil, errors.New("test error")
}

func TestSendFErr(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendF(w, r, "kkk", "/tmp", "", false)
		sendf(w, nil, "", "", false)
		fmt.Println("SendF")
	}))
	util.HGet("%v?", ts.URL)
}
