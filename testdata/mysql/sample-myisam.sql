CREATE DATABASE country;

CREATE TABLE IF NOT EXISTS city (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` CHAR(35) NOT NULL DEFAULT '',
  `country_code` CHAR(3) NOT NULL DEFAULT '',
  `district` CHAR(20) NOT NULL DEFAULT '',
  `population` INT(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `country_code` (`country_code`)
) engine=MyISAM default charset=latin1 collate=latin1_general_ci;