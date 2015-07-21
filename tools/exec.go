package tools

import (
	"bytes"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"html/template"
	"io"
	"os/exec"
	"sort"
	"sync"
	"sync/atomic"
)

type ExeH interface {
	NewCmd(e *Exec, idx string) *exec.Cmd
	OnDone(e *Exec, idx string, res *util.Map, wg *util.WaitGroup)
}
type Exec struct {
	H       ExeH
	ShowLog bool
	//
	//
	Res   map[string]*util.Map
	cmds  map[string]*exec.Cmd
	Bin   string
	Args  []string
	res_l sync.RWMutex
	exe_w util.WaitGroup
	idxc  int32
}

func NewExec(bin string, args ...string) *Exec {
	return &Exec{
		H:    nil,
		Res:  map[string]*util.Map{},
		cmds: map[string]*exec.Cmd{},
		Bin:  bin,
		Args: args,
		idxc: 0,
	}
}
func (e *Exec) log(f string, args ...interface{}) {
	if e.ShowLog {
		log.D(f, args...)
	}
}
func (e *Exec) exec() {
	idx := fmt.Sprintf("%v", atomic.AddInt32(&e.idxc, 1)-1)
	res := &util.Map{}
	//
	defer func() {
		e.res_l.Lock()
		delete(e.cmds, idx)
		e.res_l.Unlock()
		if e.H == nil {
			e.exe_w.Done()
		} else {
			e.H.OnDone(e, idx, res, &e.exe_w)
		}
	}()
	e.log("start run %v by %v", e.Bin, e.Args)
	//
	res.SetVal("beg", util.Now())
	var cmd *exec.Cmd
	if e.H == nil {
		cmd = exec.Command(e.Bin, e.Args...)
	} else {
		cmd = e.H.NewCmd(e, idx)
	}
	e.res_l.Lock()
	e.Res[idx] = res
	e.cmds[idx] = cmd
	e.res_l.Unlock()

	outb, errb := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd.Stdout, cmd.Stderr = outb, errb
	err := cmd.Start()
	if err != nil {
		res.SetVal("err", err)
		res.SetVal("end", util.Now())
		e.log("run %v error:%v", e.Bin, err)
		return
	}
	res.SetVal("err", cmd.Wait())
	res.SetVal("o_out", outb.String())
	res.SetVal("o_err", errb.String())
	res.SetVal("end", util.Now())
	e.log("exec end...")
}

func (e *Exec) Run(tc int) {
	e.exe_w.Add(tc)
	for i := 0; i < tc; i++ {
		go e.exec()
	}
}

func (e *Exec) Wait() {
	e.exe_w.Wait()
}
func (e *Exec) Data() (string, int, int) {
	var total, suc int = 0, 0
	for _, kv := range e.Res {
		total++
		if !kv.Exist("err") {
			suc++
		}
	}
	return util.S2Json(e.Res), suc, total
}
func (e *Exec) emma(suc, total int) string {
	return fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<report>
  <data>
    <all name="all classes">
      <coverage type="class, %%" value="100%%  (1/1)"/>
      <coverage type="method, %%" value="100%%  (1/1)"/>
      <coverage type="block, %%" value="100%%  (1/1)"/>
      <coverage type="line, %%" value="%v%%  (%v/%v)"/>
    </all>
  </data>
</report>`, float64(suc)/float64(total)*100, suc, total)
}
func (e *Exec) SaveP(fp, emma string) error {
	data, suc, total := e.Data()
	if len(emma) > 0 {
		err := util.FWrite(emma, e.emma(suc, total))
		if err != nil {
			return err
		}
	}
	if len(fp) > 0 {
		return util.FWrite(fp, data)
	}
	return nil
}
func (e *Exec) Save(w io.Writer) (int, error) {
	return w.Write([]byte(util.S2Json(e.Res)))
}
func (e *Exec) Execing() int {
	return e.exe_w.Size()
}
func (e *Exec) List(hs *routing.HTTPSession) routing.HResult {
	t, _ := template.New("Exec").Parse(`
		<html>
		<body>
		<ul>
		{{range $idx,$val:=.Items}}
			<li><span style="width:30px;display:inline-block;">{{$val}}</span><a href="logs?id={{$val}}&key=o_out">out.log</a>&nbsp;&nbsp;&nbsp;&nbsp;<a href="logs?id={{$val}}&key=o_err">err.log</a></li>
		{{end}}
		</ul>
		</body>
		</html>
		`)
	keys := []int{}
	for key, _ := range e.Res {
		iv, _ := util.ParseInt(key)
		keys = append(keys, iv)
	}
	sort.Sort(sort.IntSlice(keys))
	t.Execute(hs.W, map[string]interface{}{
		"Items": keys,
	})
	return routing.HRES_RETURN
}
func (e *Exec) Logs(hs *routing.HTTPSession) routing.HResult {
	var id string
	var key string = "o_out"
	err := hs.ValidRVal(`
		id,R|S,L:0;
		key,O|S,L:0;
		`, &id, &key)
	if err != nil {
		hs.SendT("id must not empty", routing.CT_TEXT)
	}
	if res, ok := e.Res[id]; ok {
		hs.SendT(res.StrVal(key), routing.CT_TEXT)
	} else {
		hs.SendT("id not found", routing.CT_TEXT)
	}
	return routing.HRES_RETURN
}

type ExeK struct {
	*Exec
	CmdF  func(exe *Exec, exk *ExeK, idx string) *exec.Cmd
	MT    int64
	last  int64
	Min   int
	Max   int
	Total int
	done  int
	d_lck sync.RWMutex
	//
	Done sync.WaitGroup
}

func NewExeK(min, max, total int, bin string, args ...string) *ExeK {
	eh := &ExeK{
		MT:    0,
		Min:   min,
		Max:   max,
		Total: total,
	}
	eh.Exec = NewExec(bin, args...)
	eh.Exec.H = eh
	return eh
}
func (e *ExeK) Start() {
	e.Done.Add(1)
	e.Run(e.Max)
	e.last = util.Now()
}
func (e *ExeK) NewCmd(exe *Exec, idx string) *exec.Cmd {
	if e.CmdF == nil {
		return exec.Command(exe.Bin, exe.Args...)
	} else {
		return e.CmdF(exe, e, idx)
	}
}
func (e *ExeK) OnDone(exe *Exec, idx string, res *util.Map, wg *util.WaitGroup) {
	e.d_lck.Lock()
	defer e.d_lck.Unlock()
	wg.Done()
	e.done++
	r := wg.Size()
	tc := e.Total - e.done - r
	//
	if tc == 0 {
		if r == 0 {
			e.Done.Done()
		}
		return
	}
	if tc > 0 && (util.Now()-e.last < e.MT) {
		exe.Run(1)
		return
	}
	if r > e.Min {
		return
	}
	if tc > e.Max {
		exe.Run(e.Max - e.Min)
		e.last = util.Now()
	} else {
		exe.Run(tc)
	}
}
func (e *ExeK) Wait() {
	e.Done.Wait()
}
func (e *ExeK) DoneSize() int {
	return e.done
}
