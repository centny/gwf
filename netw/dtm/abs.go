package dtm

import (
	"bytes"
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type AbsCreator interface {
	Create(sec string, fcfg *util.Fcfg) (Abs, error)
}

type FuncCreator func(sec string, fcfg *util.Fcfg) (Abs, error)

func (f FuncCreator) Create(sec string, fcfg *util.Fcfg) (Abs, error) {
	return f(sec, fcfg)
}

type Abs interface {
	Match(dtcm *DTCM_S, id, info interface{}, args ...interface{}) bool
	Build(dtcm *DTCM_S, id, info interface{}, args ...interface{}) (interface{}, interface{}, []interface{}, error)
}

var creators = map[string]AbsCreator{
	"N":   FuncCreator(NewAbsN),
	"CMD": FuncCreator(NewAbsC),
}

func AddCreator(name string, creator AbsCreator) {
	creators[name] = creator
}

func CreateAbs(sec string, fcfg *util.Fcfg) (Abs, error) {
	var atype = fcfg.Val2(sec+"/type", "")
	if len(atype) < 1 {
		return nil, util.Err("the AbsBuilder type is empty on section(%v)", sec)
	}
	if creator, ok := creators[atype]; ok {
		return creator.Create(sec, fcfg)
	} else {
		return nil, util.Err("the AbsBuilder by type(%v) on section(%v) is not exist", atype, sec)
	}
}

type AbsInfo struct {
	Alias string      `bson:"alias" json:"alias"`
	Info  interface{} `bson:"info" json:"info"`
}

type AbsN struct {
	Sec   string
	Cfg   *util.Fcfg
	Regs  []*regexp.Regexp
	Args  []string
	Envs  string
	WDir  string
	Alias string
}

func NewAbsN(sec string, cfg *util.Fcfg) (Abs, error) {
	regs_, err := ParseRegs(sec, cfg)
	if err != nil {
		return nil, err
	}
	args := cfg.Val2(sec+"/args", "")
	args_ := []string{}
	if len(args) > 0 {
		args_ = util.ParseArgs(args)
	}
	wdir := cfg.Val2(sec+"/wdir", ".")
	return &AbsN{
		Sec:   sec,
		Cfg:   cfg,
		Regs:  regs_,
		Args:  args_,
		Envs:  cfg.Val2(sec+"/envs", ""),
		WDir:  wdir,
		Alias: cfg.Val2(sec+"/alias", sec),
	}, err
}

func (a *AbsN) Match(dtcm *DTCM_S, id, info interface{}, args ...interface{}) bool {
	var arg = fmt.Sprintf("%v", args[0])
	for _, reg := range a.Regs {
		if reg.MatchString(arg) {
			return true
		}
	}
	return false
}
func (a *AbsN) Build(dtcm *DTCM_S, id, info interface{}, args ...interface{}) (interface{}, interface{}, []interface{}, error) {
	return id, &AbsInfo{Alias: a.Alias, Info: info}, args, nil
}

type AbsC struct {
	*AbsN
	Cmds string
	Seq  string
}

func NewAbsC(sec string, cfg *util.Fcfg) (Abs, error) {
	var n, err = NewAbsN(sec, cfg)
	if err != nil {
		return nil, err
	}
	var cmds = cfg.Val2(sec+"/cmds", "")
	if len(cmds) < 1 {
		return nil, util.Err("create AbsC fail with %v/cmds is empty", sec)
	}
	var seq = cfg.Val2(sec+"/seq", ",")
	return &AbsC{
		AbsN: n.(*AbsN),
		Cmds: cmds,
		Seq:  seq,
	}, nil
}
func (a *AbsC) Build(dtcm *DTCM_S, id, info interface{}, args ...interface{}) (interface{}, interface{}, []interface{}, error) {
	var cfg = util.NewFcfg3()
	for idx, arg := range args {
		cfg.SetVal(fmt.Sprintf("v%v", idx), fmt.Sprintf("%v", arg))
	}
	for idx, arg := range a.Args {
		cfg.SetVal(fmt.Sprintf("arg%v", idx), fmt.Sprintf("%v", arg))
	}
	cfg.SetVal("v_id", fmt.Sprintf("%v", id))
	var cmds = cfg.EnvReplaceV(a.Cmds, false)
	var runner = exec.Command(a.Cfg.Val2("bash_c", "bash"), "-c", cmds)
	runner.Dir = a.WDir
	if len(a.Envs) > 0 {
		runner.Env = append(os.Environ(), strings.Split(a.Envs, ",")...)
	}
	buf := &bytes.Buffer{}
	runner.Stdout = buf
	runner.Stderr = buf
	var err = runner.Run()
	if err != nil {
		return nil, nil, nil, util.Err("AbsC run cmds(%v) error->%v", cmds, err)
	}
	var res = strings.Trim(buf.String(), " \t\n")
	if len(res) < 1 {
		return nil, nil, nil, util.Err("AbsC run cmds(%v) result is epmty", cmds)
	}
	var bargs = []interface{}{}
	for _, arg := range strings.Split(res, a.Seq) {
		bargs = append(bargs, arg)
	}
	return id, &AbsInfo{Alias: a.Alias, Info: info}, bargs, nil
}
