package dbutil

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/Centny/TDb"
	"github.com/Centny/gwf/test"
	"github.com/Centny/gwf/util"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

type TSt struct {
	Tid    int64     `m2s:"TID"`
	Tname  string    `m2s:"TNAME"`
	Titem  string    `m2s:"TITEM"`
	Tval   string    `m2s:"TVAL"`
	Status string    `m2s:"STATUS"`
	Time   time.Time `m2s:"TIME"`
	T      int64     `m2s:"TIME" it:"Y"`
	Fval   float64   `m2s:"FVAL"`
	Uival  int64     `m2s:"UIVAL"`
	Add1   string    `m2s:"ADD1"`
	Add2   string    `m2s:"Add2"`
	LTime  int64     `m2s:"V_TIME" it:"Y"`
}

func TestDbUtilX(t *testing.T) {
	db, _ := sql.Open("mysql", test.TDbCon)
	defer db.Close()
	var tss []TSt
	DbQueryS(db, &tss, `SELECT V_TIME,TID FROM AGS_ORDER where tid<10`)
	fmt.Println(tss)
}
func TestDbUtil(t *testing.T) {
	db, _ := sql.Open("mysql", test.TDbCon)
	defer db.Close()
	err := DbExecF(db, "ttable.sql")
	if err != nil {
		t.Error(err.Error())
	}
	res, err := DbQuery(db, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(res) < 1 {
		t.Error("not data")
		return
	}
	if len(res[0]) < 1 {
		t.Error("data is empty")
		return
	}
	bys, err := json.Marshal(res)
	fmt.Println(string(bys))
	fmt.Println("T-->00")
	//
	var mres []TSt
	err = DbQueryS(db, &mres, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(mres) < 1 {
		t.Error("not data")
		return
	}
	fmt.Println("...", mres[0].T, util.Timestamp(mres[0].Time), util.Timestamp(time.Now()))
	fmt.Println(mres, mres[0].Add1)
	var mres2 []*TSt
	err = DbQueryS(db, &mres2, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(mres2) < 1 {
		t.Error("not data")
		return
	}
	//
	tx, _ := db.Begin()
	err = DbQueryS2(tx, &mres, "select * from ttable where tid>?", 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	tx.Commit()
	fmt.Println("T-->01")
	//
	ivs, err := DbQueryInt(db, "select * from ttable where tid")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(ivs) < 1 {
		t.Error("not data")
		return
	}
	_, err = DbQueryI(db, "select count(*) from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = DbQueryF(db, "select count(*) from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = DbQueryI(db, "select tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryI(db, "selects tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryF(db, "select tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryF(db, "selects tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryStr(db, "select count(*) from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = DbQueryStr(db, "select tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryStr(db, "selectsfs tidfrom ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	tx, _ = db.Begin()
	_, err = DbQueryI2(tx, "select count(*) from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = DbQueryF2(tx, "select count(*) from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = DbQueryI2(tx, "select tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryF2(tx, "select tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryI2(tx, "selects tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryF2(tx, "selects tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryFloat2(tx, "select tid from ttable where tid<?")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryInt2(tx, "select tid from ttable where tid<?")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryStr2(tx, "select count(*) from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = DbQueryStr2(tx, "select tid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryStr2(tx, "selectsfsd fsftid from ttable where tid<1")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryString2(tx, "select2 tid from ttable where tid<?")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryString2(tx, "select tid from ttable where tid<?")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQuery2(tx, "select2 tid from ttable where tid<?")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQuery2(tx, "select tid from ttable where tid<?")
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Commit()
	fmt.Println("T-->02")
	//
	svs, err := DbQueryString(db, "select tname from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(svs) < 1 {
		t.Error("not data")
		return
	}
	fmt.Println("T-->03")
	//
	tx, _ = db.Begin()
	svs, err = DbQueryString2(tx, "select tname from ttable")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(svs) < 1 {
		t.Error("not data")
		return
	}
	tx.Commit()
	//
	iid, err := DbInsert(db, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(iid)
	fmt.Println("T-->04")
	//
	tx, _ = db.Begin()
	iid2, err := DbInsert2(tx, "insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now())")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(iid2)
	tx.Commit()
	fmt.Println("T-->05")
	//
	erow, err := DbUpdate(db, "delete from ttable where tid=?", iid)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(erow)
	fmt.Println("T-->06")
	//
	tx, _ = db.Begin()
	erow, err = DbUpdate2(tx, "delete from ttable where tid=?", iid2)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(erow)
	tx.Commit()
	//
	_, err = DbQuery(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryInt(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryFloat(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryString(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbInsert(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx, _ = db.Begin()
	_, err = DbInsert2(tx, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = DbUpdate(db, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = DbUpdate2(tx, "selectt * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = DbQuery(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryInt(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryFloat(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbQueryString(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = DbInsert(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = DbInsert2(tx, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	_, err = DbUpdate(db, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	tx, _ = db.Begin()
	_, err = DbUpdate2(tx, "select * from ttable where tid>?", 1, 2)
	if err == nil {
		t.Error("not error")
		return
	}
	tx.Rollback()
	//
	err = DbQueryS(nil, nil, "select * from ttable where tid>?", 1)
	if err == nil {
		t.Error("not error")
		return
	}
	DbQueryInt(nil, "select * from ttable where tid>?", 1, 2)
	DbQueryFloat(nil, "select * from ttable where tid>?", 1, 2)
	DbQueryString(nil, "select * from ttable where tid>?", 1, 2)
	DbInsert(nil, "select * from ttable where tid>?", 1, 2)
	DbUpdate(nil, "select * from ttable where tid>?", 1, 2)
	DbInsert2(nil, "select * from ttable where tid>?", 1, 2)
	DbUpdate2(nil, "select * from ttable where tid>?", 1, 2)
	DbQuery2(nil, "select * from ttable where tid>?", 1, 2)
	DbQueryInt2(nil, "select * from ttable where tid>?", 1, 2)
	DbQueryFloat2(nil, "select * from ttable where tid>?", 1, 2)
	DbQueryString2(nil, "select * from ttable where tid>?", 1, 2)
	DbQueryS2(nil, nil, "query", 11)
	//
	fmt.Println("------->")
}

func Map2Val2(columns []string, row map[string]interface{}, dest []driver.Value) {
	for i, c := range columns {
		if v, ok := row[c]; ok {
			switch c {
			case "INT":
				dest[i] = int(v.(float64))
			case "UINT":
				dest[i] = uint32(v.(float64))
			case "FLOAT":
				dest[i] = float32(v.(float64))
			case "SLICE":
				dest[i] = []byte(v.(string))
			case "STRING":
				dest[i] = v.(string)
			case "STRUCT":
				dest[i] = time.Now()
			case "BOOL":
				dest[i] = true
			}
		} else {
			dest[i] = nil
		}
	}
}
func TestDbUtil2(t *testing.T) {
	TDb.Map2Val = Map2Val2
	db, _ := sql.Open("TDb", "td@tdata.json")
	defer db.Close()
	res, err := DbQuery(db, "SELECT * FROM TESTING WHERE INT=? AND STRING=?", 1, "cny")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(res)
}
func TestDbUtilErr(t *testing.T) {
	db2, _ := sql.Open("mysql", test.TDbCon)
	db2.Close()
	DbQuery(db2, "select * from ttable where tid>?", 1, 2)
	DbInsert(db2, "select * from ttable where tid>?", 1, 2)
	DbUpdate(db2, "select * from ttable where tid>?", 1, 2)
	DbQueryString(db2, "select * from ttable where tid>?", 1, 2)
	DbQueryInt(db2, "select * from ttable where tid>?", 1, 2)
	DbQueryFloat(db2, "select * from ttable where tid>?", 1, 2)
}
func TestDbExecF(t *testing.T) {
	db, _ := sql.Open("mysql", test.TDbCon)
	defer db.Close()
	err := DbExecF(db, "ttable.sql")
	if err != nil {
		t.Error(err.Error())
	}
	DbExecF(nil, "ttable.sql")
	DbExecF(db, "ttables.sql")
	db.Close()
	DbExecF(db, "ttable.sql")
}
func TestDbExecF2(t *testing.T) {
	err := DbExecF2("mysql", test.TDbCon, "ttable.sql")
	if err != nil {
		t.Error(err.Error())
	}
	err = DbExecF2("myl", test.TDbCon, "ttable.sql")
	if err == nil {
		t.Error("not error")
	}
}

func TestCov(t *testing.T) {
	if "'ab','dd'" != CovInvS([]string{"ab", "dd"}) {
		t.Error("error")
	}
	if "1,2" != CovInvI([]int64{1, 2}) {
		t.Error("error")
	}
}

var fsrv_sql string = `/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     4/11/2015 6:39:35 PM                         */
/*==============================================================*/


drop table if exists RCP_ACTIVITY;

drop table if exists RCP_AUDIT_COURSE;

drop table if exists RCP_AUDIT_STRATEGY;

drop table if exists RCP_CATEGORY;

drop table if exists RCP_COURSE;

drop table if exists RCP_COURSE_CATEGORY;

drop table if exists RCP_COURSE_REF;

drop table if exists RCP_LIVE;

drop table if exists RCP_QUALIFICATION;

drop table if exists RCP_SCORE;

drop table if exists RCP_SECTION;

drop table if exists RCP_STUDY_STAT;

drop table if exists RCP_TEACH;

drop table if exists RCP_TRACE_RECORD;

drop table if exists RCP_U_C_AUTH;

/*==============================================================*/
/* Table: RCP_ACTIVITY                                          */
/*==============================================================*/
create table RCP_ACTIVITY
(
   TID                  integer not null auto_increment,
   CID                  integer,
   NAME                 varchar(255),
   START_TIME           timestamp default '0000-00-00 00:00:00',
   END_TIME             timestamp default '0000-00-00 00:00:00',
   CONTENT              longtext,
   STATUS               varchar(255),
   TIME                 timestamp default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   ADD1                 varchar(255),
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_AUDIT_COURSE                                      */
/*==============================================================*/
create table RCP_AUDIT_COURSE
(
   TID                  int not null auto_increment,
   COURSE_ID            int comment '课程id',
   AUDITOR              int comment '审核人',
   ORG_ID               int comment '机构ID/所属ID',
   STATUS               int comment '待审核10、审核未通过20、审核通过30、',
   REASON               varchar(255) comment '通过或未通过的理由',
   AUDIT_TIME           timestamp,
   AUTO_CHECK           int comment '是否自动审核通过0 否 1是',
   CREATED_TIME         timestamp,
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_AUDIT_STRATEGY                                    */
/*==============================================================*/
create table RCP_AUDIT_STRATEGY
(
   TID                  int not null auto_increment,
   AUTO_CHECK           int,
   ORG_ID               int,
   UID                  int comment '老师',
   TIME                 timestamp,
   OPERATOR             int comment '操作者',
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_CATEGORY                                          */
/*==============================================================*/
create table RCP_CATEGORY
(
   TID                  int not null auto_increment,
   PID                  int,
   UID                  int,
   NAME                 varchar(255),
   TYPE                 int,
   SYS_FLAG             int,
   TIME                 timestamp,
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_COURSE                                            */
/*==============================================================*/
create table RCP_COURSE
(
   TID                  INTEGER not null auto_increment,
   NAME                 VARCHAR(255),
   IMGS                 VARCHAR(1024),
   START_TIME           timestamp default '0000-00-00 00:00:00',
   BURDEN_TYPE          integer comment '1:每周；2:每天',
   BURDEN               float comment '课程负载',
   GUIDE                varchar(255) comment '导学老师，id逗号分割',
   CREDIT               float comment '学分',
   DESCRIPTION          longtext,
   TOTAL_PRICE          float,
   SECTION_PRICE        float,
   TEACH_PRICE          float,
   ANSWER_PRICE         float,
   ACTIVITY_PRICE       float,
   TEST_PRICE           float,
   USER_NAME            varchar(255),
   USER                 integer,
   COURSE_TYPE          INTEGER comment '10:课程;20：题库;30:活动',
   BANK_ID              INTEGER comment '评测系统题库id',
   EXT                  varchar(1024) comment 'json格式备用字段',
   TAGS                 varchar(2048),
   STATUS               varchar(255),
   TIME                 timestamp default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   ADD1                 varchar(255),
   key AK_KEY_1 (TID)
);

/*==============================================================*/
/* Table: RCP_COURSE_CATEGORY                                   */
/*==============================================================*/
create table RCP_COURSE_CATEGORY
(
   TID                  int not null auto_increment,
   COURSE_ID            int,
   CATEGORY_ID          int,
   TIME                 timestamp,
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_COURSE_REF                                        */
/*==============================================================*/
create table RCP_COURSE_REF
(
   TID                  integer not null auto_increment,
   L_ID                 integer,
   R_ID                 integer comment '被指向的课程',
   TYPE                 integer comment '20：相关课程，左id为本课程，右id为相关课程；
            10：前置课程，左id为本课程，右id为指向课程；
            30：后置课程，左id为本课程，右id为指向课程；',
   CID                  integer comment '添加此记录的课程id',
   STATUS               varchar(255),
   TIME                 timestamp default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   ADD1                 varchar(255),
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_LIVE                                              */
/*==============================================================*/
create table RCP_LIVE
(
   TID                  int not null auto_increment,
   CONF_ID              varchar(255) comment '会议id',
   TOPIC                varchar(255) comment '主题名称',
   CHAIRMAN_PASS        varchar(255) comment '主播密码',
   ACTIVE_PASS          varchar(255) comment '参与者密码',
   BEGINTIME            varchar(255) comment '开始时间',
   ENDTIME              varchar(255) comment '结束时间',
   ATTEND_NUM           int comment '参与人数',
   COURSE_ID            int comment '课程id',
   UID                  int comment '课程者uid',
   TYPE                 integer comment '10:教学；20:答疑；30:评测',
   STATUS               varchar(255),
   TIME                 timestamp default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   ADD1                 varchar(255),
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_QUALIFICATION                                     */
/*==============================================================*/
create table RCP_QUALIFICATION
(
   TID                  int not null auto_increment,
   UID                  int,
   TITLE                varchar(255),
   PICS                 varchar(2048),
   DESCRIPTION          text,
   TIME                 timestamp,
   key AK_KEY_1 (TID)
);

/*==============================================================*/
/* Table: RCP_SCORE                                             */
/*==============================================================*/
create table RCP_SCORE
(
   TID                  integer not null auto_increment,
   UID                  integer,
   STUDENT_ID           integer,
   STUDENT_NO           varchar(255),
   STUDENT_NAME         varchar(255),
   OWNER                integer comment '成绩类型 ：10 课程 20 题库',
   OID                  integer,
   USUALLY              float comment '平时成绩',
   USUALLY_SOURCE       varchar(255) comment '平时成绩来源试卷',
   USUALLY_PERCENT      float comment '平时成绩比例',
   PERIOD               float comment '期中成绩',
   PERIOD_SOURCE        varchar(255) comment '期中成绩来源试卷',
   PERIOD_PERCENT       float comment '期中成绩比例',
   ENDING               float comment '期末成绩',
   ENDING_SOURCE        varchar(255) comment '期末成绩来源试卷',
   ENDING_PERCENT       float comment '期末成绩比例',
   TOTAL                float comment '总分',
   TOTAL_TYPE           integer comment '总分类型  10 百分制 20 五分制',
   CREDIT               float comment '学分',
   REMARK               varchar(255) comment '备注',
   STATUS               varchar(255),
   TIME                 timestamp default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   ADD1                 varchar(255),
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_SECTION                                           */
/*==============================================================*/
create table RCP_SECTION
(
   TID                  integer not null auto_increment,
   NAME                 varchar(255),
   CONTENT              longtext,
   CONTENT_TYPE         varchar(255),
   EXTRA                longtext comment '用来保存额外信息 格式自定',
   PRE_TAG              varchar(255),
   NEXT_TAG             varchar(255),
   REF_TAG              varchar(255) comment '相关知识点标签',
   CID                  integer,
   PID                  integer,
   SEQ                  integer comment '序号',
   PRICE                float,
   STATUS               varchar(255),
   TIME                 timestamp default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   ADD1                 varchar(255),
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_STUDY_STAT                                        */
/*==============================================================*/
create table RCP_STUDY_STAT
(
   TID                  int not null auto_increment,
   UID                  int,
   CID                  int,
   SECTION_ID           int,
   SECTION_PID          int,
   STUDY_TIME           bigint,
   USE_COUNT            int,
   PROGRESS             float,
   LAST_STUDY_TIME      timestamp,
   TIME                 timestamp,
   primary key (TID)
);

/*==============================================================*/
/* Index: CID                                                   */
/*==============================================================*/
create index CID on RCP_STUDY_STAT
(
   CID
);

/*==============================================================*/
/* Table: RCP_TEACH                                             */
/*==============================================================*/
create table RCP_TEACH
(
   TID                  integer not null auto_increment,
   NAME                 varchar(255),
   CID                  integer,
   PRICE                float,
   START_DATE           date,
   END_DATE             date,
   START_TIME           time,
   END_TIME             time,
   TIME_GROUP           varchar(255),
   DESCRIPTION          longtext comment '描述',
   TEACHER              varchar(255) comment '老师ID',
   TEACHER_NAME         varchar(255) comment '老师名称',
   TYPE                 integer comment '10:教学；20:答疑；30:评测',
   EXTRA                longtext comment '保存额外信息',
   USR_LIMIT            integer comment '在线答疑每人限制次数，-1为无限次数',
   STATUS               varchar(255),
   TIME                 timestamp default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   ADD1                 varchar(255),
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_TRACE_RECORD                                      */
/*==============================================================*/
create table RCP_TRACE_RECORD
(
   TID                  int not null auto_increment,
   OPERATION_ID         varchar(24),
   START_TIME           timestamp,
   END_TIME             timestamp,
   TIME                 timestamp,
   primary key (TID)
);

/*==============================================================*/
/* Table: RCP_U_C_AUTH                                          */
/*==============================================================*/
create table RCP_U_C_AUTH
(
   TID                  integer not null auto_increment,
   CID                  integer,
   UID                  integer,
   OID                  integer,
   OWNER                integer comment '10:整门课程；""20:教学、评测；30:活动；40:章节;',
   STATUS               varchar(255),
   TIME                 timestamp,
   ADD1                 varchar(255),
   primary key (TID)
);

/*==============================================================*/
/* Index: INDEX_1                                               */
/*==============================================================*/
create index INDEX_1 on RCP_U_C_AUTH
(
   CID
)
`

func TestDbs(t *testing.T) {
	db, _ := sql.Open("mysql", test.TDbCon)
	defer db.Close()
	err := DbExecScript(db, fsrv_sql)
	if err != nil {
		t.Error(err.Error())
	}
}
