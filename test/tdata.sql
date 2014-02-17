drop table if exists TTABLE;

create table TTABLE
(
   TID                  integer not null,
   TNAME                varchar(255),
   TITEM                varchar(255),
   TVAL                 varchar(255),
   STATUS               varchar(255),
   TIME                 timestamp,
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (TID)
);
INSERT INTO TTABLE (TID,TNAME,TITEM,TVAL,STATUS,TIME) VALUES(1,"测试数据1","测试数据1","abc_1","N",NOW());
INSERT INTO TTABLE (TID,TNAME,TITEM,TVAL,STATUS,TIME) VALUES(2,"测试数据2","测试数据2","abc_2","N",NOW());
INSERT INTO TTABLE (TID,TNAME,TITEM,TVAL,STATUS,TIME) VALUES(3,"测试数据3","测试数据3","abc_3","N",NOW());