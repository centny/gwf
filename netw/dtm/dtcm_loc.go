package dtm

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
)

func (d *DTCM_S) RunLocTaskV(id, info interface{}, args ...interface{}) (util.Map, error) {
	var err error
	id, info, args, err = d.BuildLocArgs(id, info, args...)
	if err != nil {
		log.E("%v", err)
		return nil, err
	}
	var res = util.Map{}
	for _, cmd := range d.LocCmds {
		if !cmd.Match(args...) {
			continue
		}
		cmds := cmd.ParseCmd(args...)
		cmds = d.Local.EnvReplaceV(cmds, false)
		runner := NewResultRunner(cmds)
		runner.Dir, runner.Env, runner.Bash = d.Cfg.Val2("proc_ws", "."),
			d.Cfg.Val2("proc_env", ""), d.Cfg.Val2("bash_c", "bash")
		err = runner.Start()
		if err != nil {
			err = util.Err("RunLocTaskV start runner(%v) error(%v)", util.S2Json(runner), err)
			log.E("%v", err)
			return nil, err
		}
		res[cmd.Name], err = runner.Wait()
		if err != nil {
			err = util.Err("RunLocTaskV wail runner(%v) done error(%v)", util.S2Json(runner), err)
			log.E("%v", err)
			return nil, err
		}
	}
	return res, nil
}

func (d *DTCM_S) BuildLocArgs(id, info interface{}, args ...interface{}) (interface{}, interface{}, []interface{}, error) {
	for _, abs := range d.LocAbsL {
		if abs.Match(d, id, info, args...) {
			return abs.Build(d, id, info, args...)
		}
	}
	return nil, nil, nil, NewNotMatchedErr("DTCM_S not local abs matched by id(%v),info(%v),args(%v)", id, util.S2Json(info), util.S2Json(args))
}

func (d *DTCM_S) MatchLocArgsV(id, info interface{}, args ...interface{}) bool {
	for _, abs := range d.LocAbsL {
		if abs.Match(d, id, info, args...) {
			return true
		}
	}
	return false
}
