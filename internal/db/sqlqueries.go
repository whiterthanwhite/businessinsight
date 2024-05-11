package db

const (
	QUERY_CREATE_OPERATION_TYPE = `
		CREATE TYPE operation_type AS ENUM ('Expense', 'Income', 'Transfer');
	`
	QUERY_CREATE_TABLE_CURRENCY = `
		CREATE TABLE currency (
			code varchar(10) PRIMARY KEY CHECK (code <> ''),
			description varchar(30));
	`
	QUERY_CREATE_TABLE_ACCOUNT = `
		CREATE TABLE account (
			id smallserial PRIMARY KEY,
			name varchar(30) NOT NULL,
			currency_code varchar(10) REFERENCES currency);
	`
	QUERY_CREATE_TABLE_CATEGORY = `
		CREATE TABLE category (
			id smallserial PRIMARY KEY,
			type operation_type NOT NULL,
			name varchar(30) NOT NULL,
			description varchar(250));
	`
	QUERY_CREATE_TABLE_OPERATION = `
		CREATE TABLE operation (
			entry_no bigserial PRIMARY KEY,
			date_time timestamp,
			creation_date date,
			creation_time time,
			type operation_type NOT NULL,
			amount DECIMAL(20, 10),
			source_id smallint REFERENCES account,
			currency_code varchar(10) REFERENCES currency,
			category_id smallint REFERENCES category,
			transaction_no bigint CHECK ((type = 'Transfer' AND transaction_no <> 0) OR (type = 'Income' AND amount >= 0) OR (type = 'Expense' AND amount <= 0)),
			description varchar(250));	
	`
)
