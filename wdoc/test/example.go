package test

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http"
)

//接口标题(第一行)
//接口详细描述（可以多行）
//
//
//@url,接口url描述
//	~/xx/api/example1?a=1&b=xx&c=n	POST	application/json
//@arg,接口参数的详细描述
//	opta	required	必选参数，此处参数的值域要求，来源等
//	optb	optional	可选参数
//	optc	R	必选参数
//	optd	O	可选参数
/*
	{//样例数据
		"opta": "xx",
		"optb": "aa",
		"optc": 3,
		"optd": 3
	}
*/
//@ret,接口返回数据描述
//	code	I	通用code/data返回格式，0表示成功，其他表示失败
//	xx	S	string类型返回参数描述
//	xa	I	int类型返回参数描述
//	xb	F	float类型返回参数描述
//	xc	A	array类型返回参数描述
//	xd	O	object类型返回参数描述
/*
	{
		"code": 0,
		"data": {
			"xx": "xx",
			"xa": 1,
			"xb": 3.2,
			"xc":["arrr","ssss"],
			"xd":{"v1":1,"v2":2}
		}
	}
*/
//@tag,接口分类tags
//@author,作者,创建时间,完成
func XXV(hs *routing.HTTPSession) routing.HResult { //也可以是golang默认的http handler
	return routing.HRES_CONTINUE
}

type User struct {
	Id     string `json:"id"`
	Usr    string `json:"usr"`
	Pwd    string `json:"pwd"`
	Alias  string `json:"alias"`
	Gender int    `json:"usr"`
}

//添加用户(Json)
//通过用户、密码、别名、性别创建用户
//
//
//@url,公开接口，不需要登录
//	~/usr/api/createUser	POST	application/json
//@arg,json对象参数
//	usr		R	用户输入的用户名，不少于3个字符
//	pwd		R	用户输入的密码，不少于6个字符
//	alias	O	用户昵称
//	gender	R	性别，1表示男，2表示女
/*	样例
	{
		"usr": "abc",
		"pwd": "123456",
		"alias": "测试用户",
		"gender": 1
	}
*/
//@ret,以下返回值在code为0时才有返回
//	code	I	通用code/data返回格式，0表示成功，其他表示失败
//	id		S	用户id
//	usr		S	与请求参数意义相同
//	alias	S	与请求参数意义相同
//	gender	I	与请求参数意义相同
/*	样例
	{
		"code": 0,
		"data": {
			"id": "u_001",
			"usr": "abc",
			"alias":"测试用户",
			"gender":1
		}
	}
*/
//@tag,创建用户,用户
//@author,Centny,2016-01-28
func AddUser_j(hs *routing.HTTPSession) routing.HResult {
	var user User
	err := hs.UnmarshalJ(&user)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", util.Err("unmarshal json body error->%v", err))
	}
	if len(user.Usr) < 1 || len(user.Pwd) < 1 {
		return hs.MsgResErr2(2, "arg-err", util.Err("user or pwd is empty"))
	}
	if user.Gender != 1 && user.Gender != 2 {
		return hs.MsgResErr2(3, "arg-err", util.Err("the gender must be 1 or 2, but %v found", user.Gender))
	}
	user.Id = "new id"
	//do add user
	return hs.MsgRes(&user)
}

//添加用户(Query)
//通过用户、密码、别名、性别创建用户
//
//
//@url,公开接口，不需要登录
//	~/usr/api/createUser	GET
//@arg,url参数
//	usr		R	用户输入的用户名，不少于3个字符
//	pwd		R	用户输入的密码，不少于6个字符
//	alias	O	用户昵称
//	gender	R	性别，1表示男，2表示女
/*	样例
	~/usr/api/createUser?usr=abc&pwd=123456&alias=测试用户&gender=1
*/
//@ret,以下返回值在code为0时才有返回
//	code	I	通用code/data返回格式，0表示成功，其他表示失败
//	id		S	用户id
//	usr		S	与请求参数意义相同
//	alias	S	与请求参数意义相同
//	gender	I	与请求参数意义相同
/*	样例
	{
		"code": 0,
		"data": {
			"id": "u_001",
			"usr": "abc",
			"alias":"测试用户",
			"gender":1
		}
	}
*/
//@tag,创建用户,用户
//@author,Centny,2016-01-28
func AddUser_f(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user User
	user.Usr = r.FormValue("usr")
	user.Pwd = r.FormValue("pwd")
	user.Alias = r.FormValue("alias")
	gender_ := r.FormValue("gender")
	if len(user.Usr) < 3 || len(user.Pwd) < 6 || len(gender_) < 1 {
		fmt.Fprintf(w, "%v", util.S2Json(util.Map{
			"code": 1,
			"msg":  "user or pwd or gender is empty",
		}))
		return
	}
	if gender_ != "1" && gender_ != "2" {
		fmt.Fprintf(w, "%v", util.S2Json(util.Map{
			"code": 2,
			"msg":  fmt.Sprintf("the gender must be 1 or 2, but %v found", user.Gender),
		}))
		return
	}
	//do add user
	fmt.Fprintf(w, "%v", util.S2Json(util.Map{
		"code": 0,
		"data": user,
	}))
}
