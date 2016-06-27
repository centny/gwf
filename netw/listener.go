package netw

import (
	"errors"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/pool"
	"net"
	"time"
)

//the TCP server listener.
type Listener struct {
	*LConPool              //the connection pool.
	Port      string       //the listen port.
	L         net.Listener //the base listener.
	Running   bool         //whether running accept.
	Wc        chan int     //the wait chan.
	Limit     int64
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
		Limit:    512,
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
	var tempDelay time.Duration
	for l.Running {
		log_d("Pool(%v) waiting tcp connect", l.Id())
		if l.Current() >= l.Limit {
			if tempDelay == 0 {
				tempDelay = 5 * time.Millisecond
			} else {
				tempDelay *= 2
			}
			if max := 1 * time.Second; tempDelay > max {
				tempDelay = max
			}
			log.W("netw: Accept error: opened(%v),limit(%v); retrying in %v", l.Current(), l.Limit, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		con, err := l.L.Accept()
		if err != nil {
			log.E("accept %s error(->%s", l.Port, err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			} else {
				break
			}
		}
		tempDelay = 0
		l.Increase()
		con.(*net.TCPConn).SetNoDelay(true)
		// con.(*net.TCPConn).SetWriteBuffer(5)
		// con.(*net.TCPConn).SetWriteDeadline(t)
		log_d("accepting tcp connect(%v) in pool(%v)", con.RemoteAddr().String(), l.Id())
		go l.do_run(con)
	}
	l.Running = false
	l.Wc <- 0
	log.W("loop accept will exit...")
}

func (l *Listener) do_run(con net.Conn) {
	tcon := l.NewCon(l, l.P, con)
	if l.H.OnConn(tcon) {
		l.RunC_(tcon)
	} else {
		log.W("Pool(%v/%v) rejecting tcp connection from %v", l.Name, l.Id(), con.RemoteAddr().String())
		tcon.Close()
	}
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
