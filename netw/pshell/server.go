package pshell

import (
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"

	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"golang.org/x/crypto/ssh"
)

var SharedServer *Server

func StartServer(rcaddr string, ts map[string]int, hosts ...*Host) (err error) {
	SharedServer = NewServer(hosts...)
	err = SharedServer.Run(rcaddr, ts)
	return
}

const (
	SessionStatusNormal = "normal"
	SessionStatusError  = "error"
)

type Host struct {
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	*Host
	Status         string
	session        *ssh.Session
	client         *ssh.Client
	Srv            *Server
	out            io.WriteCloser
	stdout, stderr io.Reader
}

func (s *Session) Start() (err error) {
	log.D("Session start dail to %v by username(%v),password(%v)", s.Addr, s.Username, s.Password)
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	s.client, err = ssh.Dial("tcp", s.Addr, config)
	if err == nil {
		s.session, err = s.client.NewSession()
		if err == nil {
			s.out, _ = s.session.StdinPipe()
			s.session.Stderr = s
			s.session.Stdout = s
			err = s.session.Start(s.Srv.Shell)
		}
	}
	if err == nil {
		s.Status = SessionStatusNormal
		log.D("Session dail to %v success", s.Addr)
	} else {
		log.D("Session dail to %v fail with %v", s.Addr, err)
	}
	return
}

func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
	if s.client != nil {
		s.client.Close()
	}
}

func (s *Session) Runf(format string, args ...interface{}) (err error) {
	return s.Run(fmt.Sprintf(format, args...))
}

func (s *Session) Run(cmds string) (err error) {
	_, err = s.out.Write([]byte(cmds + "\n"))
	// err = s.session.Run(cmds)
	if err == nil {
		return
	}
	etype := reflect.TypeOf(err)
	if etype.String() != "*ExitMissingError" || etype.String() != "*ExitError" {
		log.D("Session(%v) run command fail with %v, will set session status to error", s.Name, err)
		s.Status = SessionStatusError
		go s.Start()
	}
	return
}

func (s *Session) Write(p []byte) (n int, err error) {
	s.notify(string(p))
	n = len(p)
	return
}

func (s *Session) notify(msg string) {
	msgcs := s.Srv.L.MsgCs()
	for _, msgc := range msgcs {
		msgc.Writev(util.Map{
			"n": s.Name,
			"m": msg,
		})
	}
}

type Server struct {
	L        *rc.RC_Listener_m
	Shell    string
	sessions map[string]*Session
	slck     sync.RWMutex
}

func NewServer(hosts ...*Host) *Server {
	srv := &Server{
		Shell:    "bash",
		sessions: map[string]*Session{},
	}
	for _, host := range hosts {
		srv.sessions[host.Name] = &Session{
			Host: host,
			Srv:  srv,
		}
	}
	return srv
}

func (s *Server) Run(rcaddr string, ts map[string]int) (err error) {
	s.L = rc.NewRC_Listener_m_j(pool.BP, rcaddr, s)
	s.L.Name = "Server"
	s.L.LCH = s
	s.L.AddToken(ts)
	s.L.AddHFunc("add_session", s.AddSessionH)
	s.L.AddHFunc("exec", s.ExecH)
	s.L.AddHFunc("list", s.ListH)
	err = s.L.Run()
	if err == nil {
		s.connectAllSession()
	}
	return
}

func (s *Server) connectAllSession() {
	for _, ss := range s.sessions {
		if ss.Status == SessionStatusNormal {
			continue
		}
		ss.Start()
	}
}

func (s *Server) OnLogin(rc *impl.RCM_Cmd, token string) (cid string, err error) {
	cid, err = s.L.RCH.OnLogin(rc, token)
	if err != nil {
		return
	}
	sid := rc.StrVal("sid")
	rc.Kvs().SetVal("sid", sid)
	return
}

func (s *Server) ListH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var status = util.Map{}
	s.slck.RLock()
	for name, ss := range s.sessions {
		status[name] = ss.Status
	}
	s.slck.RUnlock()
	return status, nil
}

func (s *Server) AddSessionH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var host Host
	err = rc.ValidF(`
		name,R|S,L:0;
		addr,R|S,L:0;
		username,R|S,L:0;
		password,O|S,L:0;
		`, &host.Name, &host.Addr, &host.Username, &host.Password)
	if err != nil {
		return
	}

	s.slck.RLock()
	if _, ok := s.sessions[host.Name]; ok {
		err = fmt.Errorf("the session is exists by name(%v)", host.Name)
		s.slck.RUnlock()
		return
	}
	s.slck.RUnlock()
	log.D("Server start add session by %v", util.S2Json(host))
	ss := &Session{
		Host: &host,
		Srv:  s,
	}
	err = ss.Start()
	if err == nil {
		s.slck.Lock()
		s.sessions[ss.Name] = ss
		s.slck.Unlock()
	}
	res = "OK"
	return
}

func (s *Server) ExecH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var shell, cmds string
	var cids []string
	err = rc.ValidF(`
		shell,O|S,L:0;
		cmds,O|S,L:0;
		cids,O|S,L:0;
		`, &shell, &cmds, &cids)
	if err != nil {
		return
	}
	cidsMap := map[string]int{}
	for _, cid := range cids {
		cidsMap[cid] = 1
	}
	reply := util.Map{}
	wg := sync.WaitGroup{}
	s.slck.RLock()
	for _, ss := range s.sessions {
		if ss.Status != SessionStatusNormal {
			continue
		}
		if len(cidsMap) > 0 && cidsMap[ss.Name] < 1 {
			continue
		}
		wg.Add(1)
		go func(session *Session) {
			defer wg.Done()
			var xerr error
			if len(shell) > 0 {
				log.D("Server execte cmds(%v) on host(%v,%v) by shell:\n%v", cmds, session.Name, session.Addr, shell)
				tmpf := fmt.Sprintf("/tmp/%v.sh", util.UUID())
				xerr = session.Runf("echo '%v' > %v", shell, tmpf)
				if xerr == nil {
					xerr = session.Runf("%v -e %v %v", s.Shell, tmpf, cmds)
					session.Runf("rm -f %v", tmpf)
				}
			} else {
				log.D("Server execte cmds(%v) on host(%v,%v)", cmds, session.Name, session.Addr)
				xerr = session.Run(cmds)
			}
			if xerr == nil {
				reply[session.Name] = "ok"
			} else {
				reply[session.Name] = xerr
			}
		}(ss)
	}
	s.slck.RUnlock()
	wg.Wait()
	res = reply
	return
}

//OnConn see ConHandler for detail
func (s *Server) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}

//OnClose see ConHandler for detail
func (s *Server) OnClose(c netw.Con) {
}

//OnCmd see ConHandler for detail
func (s *Server) OnCmd(c netw.Cmd) int {
	return 0
}

func (s *Server) Wait() {
	s.L.Wait()
}
