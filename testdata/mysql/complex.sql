CREATE DATABASE `client`;

USE `client`;

DROP TABLE IF EXISTS `site`;

CREATE TABLE `site` (
  `site_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `lang_id` int(10) unsigned NOT NULL DEFAULT '0',
  `zone_code` smallint(5) unsigned NOT NULL DEFAULT '0',
  `shard_id` tinyint(1) unsigned DEFAULT NULL,
  `site_name` varchar(64) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `site_url` varchar(128) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `has_logo` smallint(1) DEFAULT '0',
  `site_active` smallint(1) unsigned NOT NULL DEFAULT '1',
  `site_status` smallint(5) unsigned NOT NULL DEFAULT '0',
  `creation_ts` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`site_id`),
  KEY `lang_id` (`lang_id`),
  KEY `site_status` (`site_status`),
  KEY `idx_shardid` (`shard_id`),
  KEY `idx_siteurl` (`site_url`),
  KEY `idx_gz_shardid` (`zone_code`,`shard_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8751457 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `action`;

CREATE TABLE `action` (
  `action_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `org_id` int(11) NOT NULL,
  `product_name` enum('P1','P2','P3') COLLATE utf8_unicode_ci NOT NULL,
  `user_id` int(11) NOT NULL,
  `type` enum('T1','T2') COLLATE utf8_unicode_ci NOT NULL,
  PRIMARY KEY (`action_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_product_name` (`product_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;