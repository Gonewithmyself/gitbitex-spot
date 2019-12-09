CREATE TABLE `g_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `g_account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `user_id` bigint(20) NOT NULL,
  `currency` varchar(255) NOT NULL,
  `hold` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `available` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_currency` (`user_id`,`currency`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `g_product` (
  `id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `base_currency` varchar(255) NOT NULL,
  `quote_currency` varchar(255) NOT NULL,
  `base_min_size` decimal(32,16) NOT NULL,
  `base_max_size` decimal(32,16) NOT NULL,
  `base_scale` int(11) NOT NULL,
  `quote_scale` int(11) NOT NULL,
  `quote_increment` double NOT NULL,
  `quote_min_size` decimal(32,16) NOT NULL,
  `quote_max_size` decimal(32,16) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `g_bill` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `user_id` bigint(20) NOT NULL,
  `currency` varchar(255) NOT NULL,
  `available` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `hold` decimal(32,16) NOT NULL DEFAULT '0.0000000000000000',
  `type` varchar(255) NOT NULL,
  `settled` tinyint(1) NOT NULL DEFAULT '0',
  `notes` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_gsoci` (`user_id`,`currency`,`settled`,`id`),
  KEY `idx_s` (`settled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `g_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `key` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into `g_product`(`id`,`created_at`,`updated_at`,`base_currency`,`quote_currency`,`base_min_size`,`base_max_size`,`base_scale`,`quote_scale`,`quote_increment`,`quote_min_size`,`quote_max_size`) values
('BCH-USDT',null,null,'BCH','USDT',0.0000100000000000,10000.0000000000000000,4,2,0.01,0E-16,0E-16),
('BTC-USDT',null,null,'BTC','USDT',0.0000100000000000,10000000.0000000000000000,6,2,0.01,0E-16,0E-16),
('EOS-USDT',null,null,'EOS','USDT',0.0001000000000000,1000.0000000000000000,4,3,0,0E-16,0E-16),
('ETH-USDT',null,null,'ETH','USDT',0.0001000000000000,10000.0000000000000000,4,2,0.01,0E-16,0E-16),
('LTC-USDT',null,null,'LTC','USDT',0.0010000000000000,1000.0000000000000000,4,2,0.01,0E-16,0E-16);
INSERT INTO `g_user` (`created_at`, `updated_at`, `user_id`, `email`, `password_hash`) VALUES
	('2019-11-28 14:10:56', '2019-11-28 14:10:56', 0, '124369976@qq.com', '16ede86aa3a32052c9b218c72063d968');
INSERT INTO `g_account` (`created_at`, `updated_at`, `user_id`, `currency`, `hold`, `available`) VALUES
	(NULL, NULL, 41, 'usdt', 0.0000000000000000, 1000000.0000000000000000),
	(NULL, NULL, 41, 'btc', 0.0000000000000000, 100000.0000000000000000),
	(NULL, NULL, 41, 'eth', 0.0000000000000000, 1000000.0000000000000000);