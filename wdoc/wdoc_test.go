package wdoc

import (
	"fmt"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	os.Setenv("PATH", os.Getenv("PATH")+":/usr/local/bin")
	pp := NewParser("", "/doc", "/usr/local/bin/pandoc %v -s --highlight-style tango")
	var wait = make(chan int)
	go func() {
		pp.LoopParse(os.Getenv("GOPATH")+"/src/github.com/Centny/gwf/wdoc", nil, nil, 1000)
		wait <- 0
	}()
	time.Sleep(2 * time.Second)
	pp.Running = false
	<-wait
	var res = pp.ToMv("x1", "test,x")
	if len(res.Pkgs) < 1 {
		t.Error("error")
		return
	}
	if len(pp.Web.HS) < 1 {
		t.Error("error")
		return
	}
	res = pp.ToMv("xxxx1", "test,x")
	if len(res.Pkgs) > 0 {
		t.Error("error")
		return
	}
	res = pp.ToM()
	res.RateV()
	bys, err := res.Marshal()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(string(bys))
	//
	fmt.Println("------>")
	//
	fmt.Println(util.S2Json(res))
	//
	ts := httptest.NewMuxServer()
	ts.Mux.H("^/doc.*$", pp)
	// ts.Mux.HFunc("^/html/.*$", pp.LoadHtml)
	_, err = ts.G2("/doc/all")
	if err != nil {
		t.Error("error")
		return
	}
	fp, _ := filepath.Abs("README_cn.md")
	wid := util.Crc32([]byte(fp))

	hres, err := ts.G("/doc/html/" + wid)
	if err != nil {
		t.Error(err)
		return
	}
	if len(hres) < 50 {
		fmt.Println(hres)
		t.Error("errror")
		return
	}
	fmt.Println("xxxx->")
	fmt.Println(hres)
	fmt.Println(ts.G("/doc/html/" + wid + "/t.sh"))
	fmt.Println(ts.G("/doc/html/" + wid + "/xxx.go"))
	//
	//test case
	mres, err := ts.G2("/doc/case/list")
	if err != nil {
		t.Error(err)
		return
	}
	if mres.IntVal("wdoc") != 1 {
		t.Error("error")
		fmt.Println(util.S2Json(mres))
		return
	}
	mres, err = ts.G2("/doc/case/data?key=%v", "wdoc")
	if err != nil {
		t.Error(err)
		return
	}
	if len(mres.AryMapVal("text")) != 1 {
		t.Error("error")
		fmt.Println(util.S2Json(mres))
		return
	}
	ts.G2("/doc/case/xx")
	ts.G2("/doc/case/data?key=%v", "abc")
	//
	//
	//test error
	NewParser("", "", "").Parse("/sdfk/sds")
	NewParser("", "", "").ParseDir("/sdfk/sds", nil, nil)
	go pp.LoopParse("/dsfsfd", nil, nil, 1000)
	time.Sleep(2 * time.Second)
	pp.Running = false
	pkgs_l([]*Pkg{&Pkg{}, &Pkg{}}).Swap(0, 1)
}

func TestReg(t *testing.T) {
	var ta = "xx	R	sss"
	var tb = "x1	O	sss"
	var tc = "x2	optional	sss"
	var td = "x3	required	sss"
	if !ARG_REG.MatchString(ta) {
		t.Error("error->a")
	}
	if !ARG_REG.MatchString(tb) {
		t.Error("error->b")
	}
	if !ARG_REG.MatchString(tc) {
		t.Error("error->c")
	}
	if !ARG_REG.MatchString(td) {
		t.Error("error->d")
	}
	//
	var texts = []string{
		"xx	S	sss",
		"x1	I	sss",
		"x2	F	sss",
		"x3	A	sss",
		"x3	O	sss",
	}
	for _, text := range texts {
		if !RET_REG.MatchString(text) {
			t.Error("error->" + text)
		}
	}
	fmt.Println("done...")
}

func TestCmd(t *testing.T) {
	var res = `
登录
通过用户名密码登录
@url,不需求登录，GET请求
	~/sso/api/login		POST	application/x-www-form-urlencoded
@arg,POST参数
	usr		R	用户名
	pwd		R	用户密码
	usr=abc&pwd=123

@ret,当不使用url参数时，返回通用code/data
	code	I	0：登录成功，1：参数错误，2：用户名或密码不能，3：登录失败
	token	S	成功登录的token
	usr		O	登录成功的用户对象, 详细查看~/sso/api/uinfo
	样例
	{
		"code": 0,
		"data": {
			"token": "56ACF73479C0DE596804D024",
			"usr": {
				"id": "u14",
				"usr": ["1454173753676"],
				"status": 10,
				"last": 1454173753678,
				"time": "email": "xxx@xx.com"
			}
		}
	}

@tag,用户,登录
@author,wensh,2016-01-31
@html,readme_cn,README_cn.md,the read me
	`
	var xx = `
{
		"code": 0,
		"data": {
			"token": "56ACF73479C0DE596804D024",
			"usr": {
				"id": "u14",
				"usr": ["1454173753676"],
				"status": 10,
				"last": 1454173753678,
				"time": "email": "xxx@xx.com"
			}
		}
}
	`
	var ress = cmd_m.FindAllString(res, -1)
	for _, r := range ress {
		fmt.Println("->\n", r)
	}
	fmt.Println(json_m.MatchString(strings.Trim(xx, "\n\t ")))
}

// type coverage_ struct {
// 	XMLName  string     `xml:"coverage"`
// 	Packages []package_ `xml:"packages>package"`
// }

// type package_ struct {
// 	Name     string   `xml:"name,attr"`
// 	Linerate string   `xml:"line-rate,attr"`
// 	Classes  []class_ `xml:"classes>class"`
// }
// type class_ struct {
// 	Name     string    `xml:"name,attr"`
// 	Linerate string    `xml:"line-rate,attr"`
// 	Methods  []method_ `xml:"methods>method"`
// }
// type method_ struct {
// 	Name     string `xml:"name,attr"`
// 	Linerate string `xml:"line-rate,attr"`
// }

// func TestXx(t *testing.T) {
// 	bys, _ := xml.Marshal(&coverage_{
// 		Packages: []package_{
// 			package_{
// 				Classes: []class_{
// 					class_{
// 						Methods: []method_{
// 							method_{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	fmt.Println(string(bys))
// }

// func TestXx(t *testing.T) {
// 	wdoc := &Wdoc{
// 		Pkgs: []Pkg{
// 			Pkg{
// 				Funcs: []Func{
// 					Func{
// 						Methods: []Method{
// 							Method{},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	bys, _ := wdoc.Marshal()
// 	fmt.Println(string(bys))
// }

func TestMarkDown(t *testing.T) {
	cmds := exec.Command("bash", "-c", "./t.sh")
	// cmds.Env = []string{"PATH=/usr/local"}
	bys, err := cmds.Output()
	fmt.Println(string(bys), err)
}

func TestMatch(t *testing.T) {
	fmt.Println(regexp.MustCompile(".*\\.[(md)(MD)]+$").MatchString("xxxdd.md"))
}

func TestAbs(t *testing.T) {
	fmt.Println(filepath.Abs("../../"))
}
