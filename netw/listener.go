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
func NewListener(p *pool.BytePool, port string, n string, h CCHandler) *Listener {
	return NewListenerN(p, port, n, h, NewCon)
}
func NewListener2(p *pool.BytePool, port string, h CCHandler) *Listener {
	return NewListener(p, port, "S-", h)
}
func NewListenerN(p *pool.BytePool, port string, n string, h CCHandler, ncf NewConF) *Listener {
	ls := &Listener{
		Port:     port,
		LConPool: NewLConPoolV(p, h, n, ncf),
		Wc:       make(chan int),
	}
	return ls
}
func NewListenerN2(p *pool.BytePool, port string, h CCHandler, ncf NewConF) *Listener {
	return NewListenerN(p, port, "S-", h, ncf)
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
	log.I("Server(%v) listen tcp on port:%s", l.Name, l.Port)
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
		log_d("Pool(%v) waiting tcp connect", l.Id())
		con, err := l.L.Accept()
		if err != nil {
			log.I("accept %s error->%s", l.Port, err.Error())
			break
		}
		l.Increase()
		con.(*net.TCPConn).SetNoDelay(true)
		// con.(*net.TCPConn).SetWriteBuffer(5)
		// con.(*net.TCPConn).SetWriteDeadline(t)
		log_d("accepting tcp connect(%v) in pool(%v)", con.RemoteAddr().String(), l.Id())
		tcon := l.NewCon(l, l.P, con)
		if l.H.OnConn(tcon) {
			l.RunC(tcon)
		} else {
			log.W("Pool(%v/%v) rejecting tcp connection from %v", l.Name, l.Id(), con.RemoteAddr().String())
			tcon.Close()
		}
	}
	l.Running = false
	l.Wc <- 0
	log.I("loop accept will exit...")
}

//close the listener.
func (l *Listener) Close() {
	l.Running = false
	l.LConPool.Close()
	l.L.Close()
}

//wait the listener close.
func (l *Listener) Wait() {
	<-l.Wc
}
