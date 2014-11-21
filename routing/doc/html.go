package doc

const HTML = `
<html>
<body>
{{range $hkey,$hval:=.}}
<h1 id="{{$hkey}}" style="background:#AABBCC;font-size:25px;">{{$hkey}}</h1>
{{range $hval}}
<div style="margin-left:10px;margin-right:10px;">
	<h2 id="{{.Abs}}" style="padding:10px;background:#E0EBF5;font-size:20px;">
		{{if .Marked}}
			{{.Doc.Title}}
		{{else}}
			{{.Name}}
		{{end}}
	</h2>
	<div style="margin-left:10px;margin-right:10px;">
			<div style="margin:0px;padding-left:5px;padding-bottom:10px;background:#EEEBF5;">
				<div style="padding-top:10px;">Path: <a href="#{{.Abs}}">{{.Pkg}}/{{.Name}}</a></div>
				<div style="margin-top:10px;">Pattern: {{.Pattern}}</div>
				{{if .Marked}}
				{{if .Doc.Url}}
				<div style="margin-top:10px;">Example: 
					{{range .Doc.Url}}
					&nbsp;<a href="{{.}}">{{.}}</a>
					{{end}}
				</div>
				{{end}}
				{{end}}
			</div>
		{{if .Marked}}
			{{if .Doc.Detail}}
			<div style="margin-top:5px;padding-left:5px;padding-top:10px;padding-bottom:10px;background:#F0F0F0;">{{.Doc.Detail}}</div>
			{{end}}
			<br/>

			<b>Parameters(Required)</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range $key,$val:=.Doc.ArgsR}}
				<li><b style="font-size:16px;">{{$key}}</b> {{$val}}</li>
				{{end}}
			</ul>
			<br/>

			<b>Parameters(Optioned)</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range $key,$val:=.Doc.ArgsO}}
				<li><b style="font-size:16px;">{{$key}}</b> {{$val}}</li>
				{{end}}
			</ul>
			<br/>

			<b>Parameter Value Option</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range $key,$val:=.Doc.Option}}
				<li>
					<b style="font-size:16px;">{{$key}}</b>
					<ul style="margin:0px;padding-left:10px;">
						{{range $key1,$val1:=$val}}
						<li><b style="font-size:14px;">{{$key1}}</b> {{$val1}}</li>
						{{end}}
					</ul>
				</li>
				{{end}}
			</ul>
			<br/>

			<b>Return</b>
			<div style="margin:0px;padding:10px;background:#EEE;">
			{{.Doc.ResHTML}}
			</div>

			{{if .Doc.See}}
			<b>See</b>
			<ul style="margin:0px;padding-left:30px;background:#EEE;">
				{{range .Doc.See}}
				<li>
					<a href="#{{.Abs}}">{{.Pkg}}/{{.Name}}</a>
				</li>
				{{end}}
			</ul>
			{{end}}
		{{else}}
			<div style="margin-top:5px;padding-left:5px;padding-top:10px;padding-bottom:10px;background:#F0F0F0;">Not Marked</div>
		{{end}}
	</div>
</div>
{{end}}
{{end}}
</body>
</html>
`
