package main

import (
	"encoding/xml"
	"fmt"
	"github.com/Centny/gwf/util"
	"io/ioutil"
	"os"
)

type coverage_ struct {
	XMLName  string     `xml:"coverage"`
	Sources  []string   `xml:"sources>source"`
	Packages []package_ `xml:"packages>package"`
}

func (c *coverage_) Append(t coverage_) {
	for _, v := range t.Sources {
		c.Sources = append(c.Sources, v)
	}
	for _, v := range t.Packages {
		c.Packages = append(c.Packages, v)
	}
}

type package_ struct {
	Name       string   `xml:"name,attr"`
	Linerate   string   `xml:"line-rate,attr"`
	Branchrate string   `xml:"branch-rate,attr"`
	Complexity string   `xml:"complexity,attr"`
	Classes    []class_ `xml:"classes>class"`
}
type class_ struct {
	Name       string    `xml:"name,attr"`
	Filename   string    `xml:"filename,attr"`
	Linerate   string    `xml:"line-rate,attr"`
	Branchrate string    `xml:"branch-rate,attr"`
	Complexity string    `xml:"complexity,attr"`
	Methods    []method_ `xml:"methods>method"`
	Lines      []line_   `xml:"lines>line"`
}
type method_ struct {
	Name       string  `xml:"name,attr"`
	Signature  string  `xml:"signature,attr"`
	Linerate   string  `xml:"line-rate,attr"`
	Branchrate string  `xml:"branch-rate,attr"`
	Lines      []line_ `xml:"lines>line"`
}
type line_ struct {
	Number string `xml:"number,attr"`
	Hits   string `xml:"hits,attr"`
}

func Usage() {
	fmt.Println("Usage:mcobertura -o out.xml in1.xml in2.xml ...")
}
func main() {
	if len(os.Args) < 2 {
		Usage()
		return
	}
	fs := []string{}
	of := ""
	alen := len(os.Args)
	for i := 1; i < alen; i++ {
		switch os.Args[i] {
		case "-o":
			if i < alen-1 {
				of = os.Args[i+1]
				i++
			} else {
				fmt.Println("not file for -o")
				Usage()
				return
			}
		case "-h":
			Usage()
			return
		default:
			fs = append(fs, os.Args[i])
		}
	}
	if len(of) < 1 {
		fmt.Println("not out file")
		Usage()
		return
	}
	var out coverage_
	for _, f := range fs {
		bys, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var in coverage_
		err = xml.Unmarshal(bys, &in)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		out.Append(in)
	}
	bys, _ := xml.Marshal(out)
	err := util.FWrite(of, string(bys))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("create cobertura file to %s by %d file", of, len(fs)))
}
