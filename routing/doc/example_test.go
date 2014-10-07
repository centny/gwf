package doc

import (
	"github.com/Centny/gwf/routing"
)

func Abc(hs *routing.HTTPSession) routing.HResult {
	return routing.HRES_RETURN
}
func ExampleFunc() {
	//adding before routing handler func,
	//like func Abc(hs *routing.HTTPSession) routing.HResult
	var _ = Desc{
		Title: "List Some",
		ArgsR: map[string]interface{}{
			"type": "the type",
		},
		ArgsO: map[string]interface{}{
			"limit": "the list limit",
		},
		Option: map[string]map[string]interface{}{
			"type": map[string]interface{}{
				"A": "A type",
				"B": "B type",
			},
		},
		ResV: []map[string]interface{}{
			map[string]interface{}{
				"Type": "A",
				"Name": "AAA",
			},
		},
	}.Api(Abc)
}

func ExampleHandler() {
	//implementation Docable in handler,
	//like: func (e *Exm) Doc() *Desc
	_ = func() *Desc { // func (e *Exm) Doc() *Desc {
		return &Desc{
			Title: "List Some",
			ArgsR: map[string]interface{}{
				"type": "the type",
			},
			ArgsO: map[string]interface{}{
				"limit": "the list limit",
			},
			Option: map[string]map[string]interface{}{
				"type": map[string]interface{}{
					"A": "A type",
					"B": "B type",
				},
			},
			ResV: []map[string]interface{}{
				map[string]interface{}{
					"Type": "A",
					"Name": "AAA",
				},
			},
		}
	}
}

func ExampleListDoc() {
	mux := routing.NewSessionMux2("")
	//
	//adding document viewer handler to Mux.
	//http://localhost/abc for doc
	mux.H("/abc.*", NewDocViewer())
	//
	//adding document viewer handler(include special handler) to Mux.
	//http://localhost/abd for doc
	mux.H("/abd.*", NewDocViewerInc(".*abccc01.*"))
	//
	//adding document viewer handler(exclude some handler) to Mux.
	//http://localhost/abe for doc
	mux.H("/abe.*", NewDocViewerExc(".*abccc01.*"))
}
