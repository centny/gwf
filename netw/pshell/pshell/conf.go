package main

import (
	"strings"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/pshell"

	"github.com/Centny/gwf/util"
)

var Conf = util.NewFcfg3()

func Listen() string {
	return Conf.Val("listen")
}

func ReadHosts() (hosts []*pshell.Host) {
	for key := range Conf.Map {
		if !strings.HasPrefix(key, "hosts/") {
			continue
		}
		name := strings.TrimPrefix(key, "hosts/")
		uri := Conf.Val(key)
		parts := strings.SplitN(uri, "@", 2)
		if len(parts) < 2 {
			log.W("pars config line %v=%v fail", key, uri)
			continue
		}
		user := strings.Split(parts[0], ":")
		host := &pshell.Host{
			Name:     name,
			Addr:     parts[1],
			Username: user[0],
		}
		if len(user) > 1 {
			host.Password = user[1]
		}
		hosts = append(hosts, host)
	}
	return
}

func ReadTokens() (tokens map[string]int) {
	tokens = map[string]int{}
	for key := range Conf.Map {
		if !strings.HasPrefix(key, "token/") {
			continue
		}
		token := strings.TrimPrefix(key, "token/")
		tokens[token] = Conf.IntVal(key)
	}
	return
}
