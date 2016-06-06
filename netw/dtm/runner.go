package dtm

import (
	"bytes"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"os"
	"os/exec"
	"strings"
)

type Runner interface {
	Start() error
	Wait() (util.Map, error)
}

type ResultRunner struct {
	Bash   string                                      `json:"bash"`
	Dir    string                                      `json:"dir"`
	Env    string                                      `json:"env"`
	Cmds   string                                      `json:"cmds"`
	Beg    int64                                       `json:"beg"`
	Buf    *bytes.Buffer                               `json:"-"`
	Runner *exec.Cmd                                   `json:"-"`
	NewCmd func(name string, args ...string) *exec.Cmd `json:"-"`
}

func NewResultRunner(cmds string) *ResultRunner {
	return &ResultRunner{
		Bash:   "bash",
		Dir:    ".",
		Cmds:   cmds,
		NewCmd: exec.Command,
	}
}
func (r *ResultRunner) Start() error {
	r.Beg = util.Now()
	var runner = r.NewCmd(r.Bash, "-c", r.Cmds)
	runner.Dir = r.Dir
	if len(r.Env) > 0 {
		runner.Env = append(os.Environ(), strings.Split(r.Env, ",")...)
	}
	r.Buf = &bytes.Buffer{}
	runner.Stdout = r.Buf
	runner.Stderr = r.Buf
	r.Runner = runner
	return runner.Start()
}

func (r *ResultRunner) Wait() (util.Map, error) {
	args := util.Map{}
	err := r.Runner.Wait()
	used := util.Now() - r.Beg
	res := r.Buf.String()
	if err == nil {
		log.D("ResultRunner run_cmd by running command(\n\t%v\n) success,used(%vms)->\n%v", r.Cmds, used, res)
		args["code"] = r.cmd_do_res(args, r.Cmds, res)
	} else {
		log.E("ResultRunner run_cmd by running command(\n\t%v\n) error(%v)->\n%v", r.Cmds, err, res)
		args["code"] = -1
		args["err"] = err.Error()
	}
	args["used"] = used
	return args, err
}

func (r *ResultRunner) cmd_do_res(args util.Map, cmds, res string) int {
	var res_a = strings.SplitN(res, "----------------result----------------", 2)
	if len(res_a) < 2 {
		return 0
	}
	var mres = util.ParseSectionF("[", "]", res_a[1])
	var jval = mres.StrVal("json")
	if len(jval) < 1 {
		args["data"] = mres
		return 0
	}
	var jval_m, err = util.Json2Map(jval)
	if err == nil {
		args["data"] = jval_m
		return 0
	} else {
		log.E("DTM_C parse json result on command(\n\t%v\n) by data(%v) error->%v", cmds, jval, err)
		args["data"] = mres
		args["err"] = err.Error()
		return -2
	}
}
