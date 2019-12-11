
USE db_spot;

CREATE TABLE `g_offset` (
  `group` varchar(255) NOT NULL,
  `partition` bigint(20) NOT NULL DEFAULT '0',
  `log_offset` bigint(20) NOT NULL DEFAULT '0',
  `log_seq` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`group`, `partition`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



