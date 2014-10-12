package hsr

import (
	"database/sql"
	"github.com/Centny/gwf/dbutil"
	"github.com/Centny/gwf/log"
)

const chk_hsr string = `SELECT * FROM HSR_DEMO WHERE TID<?`

func CheckHSR(db *sql.DB) error {
	if _, err := dbutil.DbQuery(db, chk_hsr, 1); err != nil {
		log.D("HSR table not found, auto creating...")
		err = dbutil.DbExecScript(db, HSR)
		if err != nil {
			return err
		}
	}
	return nil
}

const HSR string = `
/*==============================================================*/
/* Table: HSR_ARG                                               */
/*==============================================================*/
CREATE TABLE HSR_ARG
(
   TID                  INTEGER NOT NULL AUTO_INCREMENT,
   RID                  INTEGER NOT NULL,
   NAME                 VARCHAR(255) NOT NULL,
   VAL                  VARCHAR(512) NOT NULL,
   TYPE                 VARCHAR(255) NOT NULL,
   PRIMARY KEY (TID)
);

/*==============================================================*/
/* Table: HSR_DEMO                                              */
/*==============================================================*/
CREATE TABLE HSR_DEMO
(
   TID                  INTEGER NOT NULL AUTO_INCREMENT,
   DEMO                 VARCHAR(255) NOT NULL,
   CPU                  NUMERIC(8,4) NOT NULL,
   MEM                  NUMERIC(8,4) NOT NULL,
   HD                   NUMERIC(8,4) NOT NULL,
   NET                  NUMERIC(8,4) NOT NULL,
   P_CPU                NUMERIC(8,4) NOT NULL,
   P_MEM                NUMERIC(8,4) NOT NULL,
   P_HD                 NUMERIC(8,4) NOT NULL,
   P_NET                NUMERIC(8,4) NOT NULL,
   P_PROC               INTEGER NOT NULL,
   P_THR                VARCHAR(255) NOT NULL,
   P_GOS                INTEGER NOT NULL,
   TIME                 TIMESTAMP NOT NULL,
   STATUS               VARCHAR(255) NOT NULL,
   PRIMARY KEY (TID)
);

/*==============================================================*/
/* Table: HSR_H                                                 */
/*==============================================================*/
CREATE TABLE HSR_H
(
   TID                  INTEGER NOT NULL AUTO_INCREMENT,
   RID                  INTEGER NOT NULL,
   NAME                 VARCHAR(255) NOT NULL,
   LOST                 INTEGER NOT NULL,
   TYPE                 VARCHAR(255) NOT NULL,
   PRIMARY KEY (TID)
);

/*==============================================================*/
/* Table: HSR_R                                                 */
/*==============================================================*/
CREATE TABLE HSR_R
(
   TID                  INTEGER NOT NULL AUTO_INCREMENT,
   DEMO                 VARCHAR(255) NOT NULL COMMENT 'the demo name.',
   U                    VARCHAR(2048) NOT NULL COMMENT 'the URL.',
   P                    VARCHAR(255) NOT NULL COMMENT 'the URL pattern.',
   M                    VARCHAR(255) NOT NULL COMMENT 'the method.',
   T                    VARCHAR(255) NOT NULL COMMENT 'the type.',
   L                    INTEGER NOT NULL COMMENT 'the lost time.',
   TIME                 TIMESTAMP NOT NULL,
   STATUS               VARCHAR(255) NOT NULL,
   PRIMARY KEY (TID)
);
`
