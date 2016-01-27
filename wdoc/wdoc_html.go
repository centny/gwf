package wdoc

var HTML = `
<html>
<head>
	<style type="text/css">
		.nav-pkg h2{
			margin: 0px;
			padding: 0px;
		}
		.nav-pkg p{
			margin: 0 0 0 16px;
		}
	</style>
</head>
<body>
	<div id="nav-header">
		{{range $pkey,$pval:=.pkgs}}
		<div id="{{$pval.name}}" class="nav-pkg">
			<h2>{{$pval.name}}</h2>
			{{range $ikey,$ival:=$pval.items}}
			<p id="{{$pval.name}}/{{$ival.name}}">
				<a href="#{{$pval.name}}_{{$ival.name}}">{{$ival.title}}</a>
				<span>{{$ival.desc}}</span>
			</p>
			{{end}}
		</div>
		{{end}}
	</div>
	<div id="nav-list">
	{{range $pkey,$pval:=.pkgs}}
		{{range $ikey,$ival:=$pval.items}}
		<div class="nav-list">
			<h3>{{$pval.name}}/{{$ival.name}}&nbsp;&nbsp;{{$ival.title}}</h3>
			<p class="nav-list-url">
				{{$ival.url.path}}&nbsp;&nbsp;{{$ival.url.method}}&nbsp;&nbsp;{{$ival.url.ctype}}&nbsp;&nbsp;{{$ival.url.desc}}
			</p>
		</div>
		{{end}}
	{{end}}
	</div>
</body>
</html>
`
