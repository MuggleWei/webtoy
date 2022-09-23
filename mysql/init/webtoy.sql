CREATE DATABASE IF NOT EXISTS webtoy; 
USE webtoy;

-- ----------------
-- user table
-- ----------------
DROP TABLE IF EXISTS t_user;
CREATE TABLE t_user (
	id INT AUTO_INCREMENT NOT NULL COMMENT 'id',
	name VARCHAR(64) NOT NULL COMMENT 'user name',
	passwd VARCHAR(128) NOT NULL COMMENT 'user password bcrypt hash',
	show_name VARCHAR(32) DEFAULT NULL COMMENT 'user show name',
	email VARCHAR(32) DEFAULT NULL COMMENT 'email',
	phone VARCHAR(32) DEFAULT NULL COMMENT 'phone number',
	PRIMARY KEY (`id`),
	UNIQUE INDEX unique_name(`name`),
	UNIQUE INDEX unique_email(`email`),
	UNIQUE INDEX unique_phone(`phone`)
) ENGINE=INNODB DEFAULT CHARSET=UTF8 COLLATE utf8_bin;
