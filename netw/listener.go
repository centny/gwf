package netw

import (
	"errors"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/pool"
	"net"
)

//the TCP server listener.
type Listener struct {
	*LConPool              //the connection pool.
	Port      string       //the listen port.
	L         net.Listener //the base listener.
	Running   bool         //whether running accept.
	Wc        chan int     //the wait chan.
}

//new one listener.
func NewListener(p *pool.BytePool, port string, h CmdHandler) *Listener {
	return &Listener{
		Port:     port,
		LConPool: NewLConPool(p, h),
		Wc:       make(chan int),
	}
}

//listen on the special port.
func (l *Listener) Listen() error {
	if len(l.Port) < 1 {
		return errors.New("port is empty")
	}
	ln, err := net.Listen("tcp", l.Port)
	if err != nil {
		return err
	}
	l.L = ln
	log.I("listen tcp on port:%s", l.Port)
	return nil
}

//run all async.
func (l *Listener) Run() error {
	err := l.Listen()
	if err != nil {
		log.E("run listener error:%s", err.Error())
		return err
	}
	go l.LoopAccept()
	go l.LoopTimeout()
	return nil
}

//looping the accept
func (l *Listener) LoopAccept() {
	l.Running = true
	for l.Running {
		con, err := l.L.Accept()
		if err != nil {
			log_d("accept %s error:%s", l.Port, err.Error())
			break
		}
		log_d("accept tcp connect from %s", con.RemoteAddr().String())
		l.RunC(NewCon(l.P, con))
	}
	l.Running = false
	l.Wc <- 0
}

//close the listener.
func (l *Listener) Close() {
	l.LConPool.Close()
	l.L.Close()
}

//wait the listener close.
func (l *Listener) Wait() {
	<-l.LConPool.Wc
	<-l.Wc
}
