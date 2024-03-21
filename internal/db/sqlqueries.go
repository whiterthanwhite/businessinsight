package db

const (
	operationType = `
		CREATE TYPE operation_type AS ENUM ('Expense', 'Income', 'Transfer');
	`
	currency = `
		CREATE TABLE currency (
			code char(10) PRIMARY KEY CHECK (code <> ''),
			description varchar(30));
	`
	account = `
		CREATE TABLE account (
			id smallserial PRIMARY KEY,
			name varchar(30) NOT NULL,
			currency_code char(10) REFERENCES currency);
	`
	category = `
		CREATE TABLE category (
			id smallserial PRIMARY KEY,
			type operation_type NOT NULL,
			name varchar(30) NOT NULL,
			description varchar(250));
	`
	operation = `
		CREATE TABLE operation (
			entry_no bigserial PRIMARY KEY,
			date_time timestamp,
			type operation_type NOT NULL,
			amount real,
			source_id smallint REFERENCES account,
			currency_code char(10) REFERENCES currency,
			category_id smallint REFERENCES category,
			transaction_no bigint CHECK ((type = 'Transfer' AND transaction_no <> 0) OR (type = 'Income' AND amount >= 0) OR (type = 'Expense' AND amount <= 0)),
			description varchar(250));	
	`
)
