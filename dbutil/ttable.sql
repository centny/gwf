DROP TABLE IF EXISTS `ttable`;
CREATE TABLE `ttable` (
  `TID` int(11) NOT NULL AUTO_INCREMENT,
  `TNAME` varchar(255) DEFAULT NULL,
  `TITEM` varchar(255) DEFAULT NULL,
  `TVAL` varchar(255) DEFAULT NULL,
  `STATUS` varchar(255) DEFAULT NULL,
  `TIME` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `FVAL` float DEFAULT NULL,
  `EVAL` int(10) unsigned DEFAULT NULL,
  `ADD1` varchar(255) DEFAULT NULL,
  `ADD2` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`TID`)
);
insert into ttable(tname,titem,tval,status,time) values('name','item','val','N',now());