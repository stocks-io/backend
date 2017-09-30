package main

import "database/sql"

func setupDB() {
	db, err = sql.Open("mysql", "root@/")
	checkFatalErr(err)
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS stocks")
	checkFatalErr(err)
	db, err = sql.Open("mysql", "root@/stocks")
	checkFatalErr(err)
	cmd := `
	    CREATE TABLE IF NOT EXISTS userinfo
	    (
	      id              	int unsigned NOT NULL auto_increment,
	      first_name		varchar(255) NOT NULL,
	      last_name			varchar(255) NOT NULL,
	      username			varchar(255) NOT NULL UNIQUE,
	      email         	varchar(255) NOT NULL UNIQUE,
	      password         	varchar(255) NOT NULL,
	      added           	varchar(255) NOT NULL,
	      PRIMARY KEY    	(id)
	    );
    `
	stmt, err := db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
	    CREATE TABLE IF NOT EXISTS sessions
	    (
	      id              	int unsigned NOT NULL auto_increment,
	      user_id			int unsigned NOT NULL,
	      token				varchar(255),
	      added           	varchar(255) NOT NULL,
	      expires           varchar(255) NOT NULL,
	      PRIMARY KEY    	string(id)
	    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
	    CREATE TABLE IF NOT EXISTS portfolio
	    (
	      id              	int unsigned NOT NULL auto_increment,
	      user_id          	int unsigned NOT NULL UNIQUE,
	      cash				FLOAT(8),
	      PRIMARY KEY     	(id)
	    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
	    CREATE TABLE IF NOT EXISTS positions
	    (
	      id              	int unsigned NOT NULL auto_increment,
	      user_id          	int unsigned NOT NULL,
	      symbol			varchar(32),
	      units				int unsigned NOT NULL,
	      buy_price			FLOAT(8),
	      PRIMARY KEY     	(id)
	    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
}

func userExists(username string) bool {
	var exists bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM userinfo WHERE username=?", username).Scan(&exists)
	checkErr(err)
	return exists
}
