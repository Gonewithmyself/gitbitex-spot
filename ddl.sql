
USE db_spot;

CREATE TABLE `g_offset` (
  `group` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `log_offset` bigint(20) NOT NULL DEFAULT '0',
  `log_seq` bigint(20) NOT NULL DEFAULT '0',
  `log_partition` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`group`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



