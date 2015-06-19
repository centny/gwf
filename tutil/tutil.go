package tutil

import (
	"bufio"
	"errors"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const LINE_LEN_LIMIT int = 1024000

type MSG_H func(tsk *TSK_C, msg string) error
type CON_H func(tsk *TSK_C) error

type TSK struct {
	L       net.Listener
	Running bool
	Port    string
	H       CON_H
	Wg      sync.WaitGroup
	Cm      map[net.Conn]*TSK_C
}

func (t *TSK) Listen() error {
	if len(t.Port) < 1 {
		return errors.New("port is empty")
	}
	ln, err := net.Listen("tcp", t.Port)
	if err != nil {
		return err
	}
	t.L = ln
	log.I("listen tcp on port:%s", t.Port)
	return nil
}

func (t *TSK) LoopAccept() {
	t.Running = true
	for t.Running {
		con, err := t.L.Accept()
		if err != nil {
			log.D("accept %s error:%s", t.Port, err.Error())
			break
		}
		log.D("accept tcp connect from %s", con.RemoteAddr().String())
		go func() {
			tc := &TSK_C{
				Tsk: t,
				C:   con,
				W:   bufio.NewWriter(con),
				R:   bufio.NewReader(con),
			}
			t.Cm[con] = tc
			defer func() {
				delete(t.Cm, con)
				tc.Close()
			}()
			t.H(tc)
		}()
	}
}
func (t *TSK) Run() {
	err := t.Listen()
	if err != nil {
		log.E("run listener error:%s", err.Error())
		return
	}
	t.LoopAccept()
}
func (t *TSK) Stop() {
	t.Running = false
	if t.L != nil {
		t.L.Close()
	}
}

func (t *TSK) NewC() (*TSK_C, error) {
	return NewTSk_C("127.0.0.1" + t.Port)
}

//async do
func (t *TSK) Conn(h MSG_H) (*TSK_C, error) {
	tc, err := t.NewC()
	if err != nil {
		return tc, err
	}
	return tc.Do(h), nil
}

//new
func NewTSK(port string, h CON_H) *TSK {
	return &TSK{
		Port: port,
		H:    h,
		Wg:   sync.WaitGroup{},
		Cm:   map[net.Conn]*TSK_C{},
	}
}

func NewTSK2(port string, h MSG_H) *TSK {
	return NewTSK(port, func(t *TSK_C) error {
		return t.Do_(h)
	})
}

type TSK_C struct {
	Tsk     *TSK
	C       net.Conn
	W       *bufio.Writer
	R       *bufio.Reader
	Running bool
}

func (t *TSK_C) Close() {
	t.Running = false
	if t.C != nil {
		t.C.Close()
	}
}

//sync do
func (t *TSK_C) Do_(h MSG_H) error {
	t.Running = true
	var terr error = nil
	for t.Running {
		bys, err := util.ReadLine(t.R, LINE_LEN_LIMIT, false)
		if err != nil {
			terr = err
			break
		}
		err = h(t, string(bys))
		if err != nil {
			terr = err
			break
		}
	}
	return terr
}

//async do
func (t *TSK_C) Do(h MSG_H) *TSK_C {
	go t.Do_(h)
	time.Sleep(100 * time.Millisecond)
	return t
}
func (t *TSK_C) Write(m string) (int, error) {
	defer t.W.Flush()
	return t.W.WriteString(m)
}
func NewTSk_C(con string) (*TSK_C, error) {
	cc, err := net.Dial("tcp", con)
	if err != nil {
		return nil, err
	}
	return &TSK_C{
		Tsk: nil,
		C:   cc,
		W:   bufio.NewWriter(cc),
		R:   bufio.NewReader(cc),
	}, nil
}

func IgMain(f func()) {
	log.I("IgMain start...")
	nargs := []string{}
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test") {
			continue
		} else {
			nargs = append(nargs, arg)
		}
	}
	os.Args = nargs
	go func() {
		log.I("IgMain main start...")
		f()
		log.I("IgMain main done...")
	}()
	for !util.Fexists(filepath.Join(os.TempDir(), "/.gwf.ig.exit")) {
		time.Sleep(time.Second)
	}
	log.I("IgMain done...")
}
