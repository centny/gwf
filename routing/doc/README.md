#Doc
>building handler api document in routing.SessionMux

###1.first add document description to handler

* for func
	
	```
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
```
* for handler

	```
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
```

###2.adding document viewer to routing.SessionMux

* view hander

	```
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
```