package filter

import (
	"fmt"
	"net"
	"strings"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
)

// 是否ip地址
func IsIPAddress(input string) bool {
	ip := net.ParseIP(input)
	return ip != nil
}

// 是否ip段
func IsIPRange(input string) bool {
	_, _, err := net.ParseCIDR(input)
	return err == nil
}

type WhitelistFilter struct {
	whitelist []*net.IPNet
}

// input: ip或ip段
// eg. NewWhitelistFilter("192.168.0.1") or NewWhitelistFilter("192.168.0.1/16")
func NewWhitelistFilter(input ...string) (*WhitelistFilter, error) {
	f := &WhitelistFilter{
		whitelist: []*net.IPNet{},
	}
	for _, in := range input {
		if IsIPAddress(in) {
			err := f.AddIP(in)
			if err != nil {
				return nil, err
			}
			continue
		}
		if IsIPRange(in) {
			err := f.AddIPRange(in)
			if err != nil {
				return nil, err
			}
			continue
		}
		return nil, fmt.Errorf("input(%v) is not a valid ip or ip range", in)
	}

	return f, nil
}

func (f *WhitelistFilter) AddIP(ipStr string) error {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	f.whitelist = append(f.whitelist, &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(32, 32),
	})

	return nil
}

func (f *WhitelistFilter) AddIPRange(ipRangeStr string) error {
	_, ipNet, err := net.ParseCIDR(ipRangeStr)
	if err != nil {
		return fmt.Errorf("invalid IP range: %s", ipRangeStr)
	}

	f.whitelist = append(f.whitelist, ipNet)
	return nil
}

func (f *WhitelistFilter) IsAllowed(ipStr string) bool {
	if ipStr == "" {
		return false
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, subnet := range f.whitelist {
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}

func (f *WhitelistFilter) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var addr = strings.Split(hs.R.Header.Get("X-Real-IP"), ":")[0]
	if len(addr) < 1 {
		addr = strings.Split(hs.R.RemoteAddr, ":")[0]
	}
	if f.IsAllowed(addr) {
		log.D("WhiteListFilter found ip(%v) in white lists, from request path(%v), xip(%v), remote(%v)",
			addr, hs.R.URL.Path, hs.R.Header.Get("X-Real-IP"), hs.R.RemoteAddr)
		return routing.HRES_CONTINUE
	}

	log.W("WhiteListFilter found ip(%v) which not in white lists, from request path(%v), xip(%v), remote(%v)",
		addr, hs.R.URL.Path, hs.R.Header.Get("X-Real-IP"), hs.R.RemoteAddr)
	return hs.MsgResErr(401, "access err", util.Err("access error"))
}
