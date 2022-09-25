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

-- ----------------
-- insert user table
-- passwd: 123456
-- ----------------
INSERT INTO t_user
(name, passwd)
VALUES
('foo', '$2a$10$mU.6kLEoCjOfq6RbLNPEIuOSYl5R6t5z.2BvS8Gs7mIcq40QL4f66'),
('bar', '$2a$10$3ZbjN2wxw82r2WjuOHGn4ei.MbRs10.1k9Y3Tbqg2Omw8HVC3vmWe'),
('baz', '$2a$10$yEl8tbZWQlGJBKjvS2WTr.Gtjha.EOjZcqj/NklSb/ydXzE4KdYJq');
