package conchapp

// tables design

// funds tables
/*
	address | amount | create_time | update_time
	1XXXXX| 100000| 2018-08-27 16:39:12| 2018-08-27 16:39:15

	create table funds (
		address CHAR(50) PRIMARY KEY NOT NULL default '',
		amount UNSIGNED BIGINT not null default '0',
		create_time INTEGER,
		update_time INTEGER
	);
*/

// transaction_records tables
/*
	id| sender | receiver | amount | input | expired | time_stamp | nonce | ref_block_num | block_num |sign

	create table transaction_records (
		id VARCHAR(64)  NOT NULL default '',
		sender CHAR(64) not null default '',
		receiver CHAR(64) not null default '',
		amount UNSIGNED BIGINT not null default '0',
		input TEXT not null default '',
		expired unsigned INTEGER not null default '0',
		time_stamp INTEGER not null default '0',
		nonce CHAR(64) not null default '',
		ref_block_num unsigned INTEGER not null default '0',
		block_num unsigned INTEGER not null default '0',
		sign VARCHAR(255)  NOT NULL default ''
	);
	CREATE UNIQUE INDEX hash_unique on transaction_records (id);
*/

// transaction_receipts tables
/*
	id | status | fee | block_num | tx_hash | log

		create table transaction_receipts (
		id VARCHAR(64)  NOT NULL default '',
		status INTEGER not null default '0',
		fee UNSIGNED BIGINT not null default '0',
		block_num unsigned INTEGER not null default '0',
		tx_hash VARCHAR(64) not null default '',
		log TEXT  NOT NULL default ''
	);
	CREATE UNIQUE INDEX id_unique on transaction_receipts (id);
	CREATE UNIQUE INDEX txhash_unique on transaction_receipts (tx_hash);
*/

// block_records tables
/*
	apphash | block_hash | block_num | tx_root | receipt_root | time_stamp

		create table block_records (
		apphash VARCHAR(64)  NOT NULL default '',
		block_hash VARCHAR(64)  NOT NULL default '',
		block_num unsigned INTEGER not null default '0',
		tx_root VARCHAR(64)  NOT NULL default '',
		receipt_root VARCHAR(64)  NOT NULL default '',
		time_stamp unsigned INTEGER not null default '0'
	);
	CREATE UNIQUE INDEX apphash_unique on block_records (apphash);
	CREATE UNIQUE INDEX block_unique on block_records (block_hash);

*/

// state table
/*
	id |state|
	  json

	create table state (
		id INTEGER not null default '1',
		json TEXT not null default '{}'
	);
*/
