package tutil

import (
	"encoding/xml"
	"fmt"
	"github.com/Centny/gwf/util"
	"io/ioutil"
	"strings"
)

type ReportE struct {
	XMLName xml.Name `xml:"report"`
	Data    DataE    `xml:"data"`
}

type DataE struct {
	All AllE `xml:"all"`
}

type AllE struct {
	Name string      `xml:"name,attr"`
	CS   []CoverageE `xml:"coverage"`
	PS   []PackageE  `xml:"package"`
}

func (a *AllE) AddPkg(pn, class, method, block, line string) error {
	var err error
	var pkg PackageE
	pkg.Name = pn
	_, err = pkg.AddCov("class", class)
	if err != nil {
		return err
	}
	_, err = pkg.AddCov("method", method)
	if err != nil {
		return err
	}
	_, err = pkg.AddCov("block", block)
	if err != nil {
		return err
	}
	_, err = pkg.AddCov("line", line)
	if err != nil {
		return err
	}
	a.PS = append(a.PS, pkg)
	return nil
}
func (a *AllE) CreateCov() {
	mv := map[string][]int{}
	for _, pkg := range a.PS {
		for _, cs := range pkg.CS {
			os := mv[cs.Type]
			if os == nil {
				os = []int{0, 0}
			}
			ns := cov_val_E(cs.Value)
			os[0], os[1] = os[0]+ns[0], os[1]+ns[1]
			mv[cs.Type] = os
		}
	}
	for k, v := range mv {
		a.CS = append(a.CS, create_cov_E(k, v))
	}
}

type PackageE struct {
	Name string      `xml:"name,attr"`
	CS   []CoverageE `xml:"coverage"`
}

func (p *PackageE) AddCov(name, data string) ([]int, error) {
	ds, ce, err := CreateCovE(name, data)
	if err != nil {
		return nil, err
	}
	p.CS = append(p.CS, ce)
	return ds, nil
}
func CreateCovE(name, data string) ([]int, CoverageE, error) {
	ds, err := util.ParseInts2(data, "/")
	if err != nil {
		return nil, CoverageE{}, err
	}
	if len(ds) < 2 {
		return nil, CoverageE{}, util.Err("invalid %v(%v),eg:1/10", name, data)
	}
	return ds, create_cov_E(fmt.Sprintf("%v, %%", name), ds), nil
}
func create_cov_E(typ string, ds []int) CoverageE {
	total := ds[1]
	if total < 1 {
		total = 1
	}
	return CoverageE{
		Type:  typ,
		Value: fmt.Sprintf("%v%% (%v/%v)", int(float64(ds[0])/float64(total)*100), ds[0], ds[1]),
	}
}
func cov_val_E(val string) []int {
	vals := strings.Split(strings.Split(val, ")")[0], "(")
	if len(vals) < 2 {
		return []int{0, 0}
	}
	is, err := util.ParseInts2(vals[1], "/")
	if err == nil && len(is) > 1 {
		return is
	} else {
		return []int{0, 0}
	}
}

type CoverageE struct {
	Type  string `xml:"type"`
	Value string `xml:"value"`
}

func Append(f, pn, class, method, block, line string) error {
	var rep ReportE
	rep.Data.All.Name = "all class"
	if util.Fexists(f) {
		bys, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		err = xml.Unmarshal(bys, &rep)
		if err != nil {
			return err
		}
	}
	err := rep.Data.All.AddPkg(pn, class, method, block, line)
	if err != nil {
		return err
	}
	rep.Data.All.CreateCov()
	bys, _ := xml.Marshal(rep)
	return util.FWrite2(f, bys)
}
