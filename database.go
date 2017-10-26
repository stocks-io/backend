package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

func setupDB() {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	if len(dbUsername) == 0 {
		panic("$DB_USERNAME is not set")
	}
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/", dbUsername, dbPassword))
	checkFatalErr(err)
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS stocks")
	checkFatalErr(err)
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/stocks", dbUsername, dbPassword))
	checkFatalErr(err)
	cmd := `
	    CREATE TABLE IF NOT EXISTS userinfo
	    (
	      id              	int unsigned NOT NULL auto_increment,
	      first_name		varchar(255) NOT NULL,
	      last_name			varchar(255) NOT NULL,
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
	      token				varchar(255) NOT NULL,
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
	      cash				FLOAT(8) NOT NULL,
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
	      symbol			varchar(32) NOT NULL,
	      units				int unsigned NOT NULL,
	      PRIMARY KEY     	(id)
	    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
	    CREATE TABLE IF NOT EXISTS order_history
	    (
	      id              	int unsigned NOT NULL auto_increment,
	      user_id          	int unsigned NOT NULL,
	      symbol			varchar(32) NOT NULL,
	      units				int unsigned NOT NULL,
	      price				float(8) NOT NULL,
	      buy				bit NOT NULL,
	      added				varchar(255) NOT NULL,
	      PRIMARY KEY     	(id)
	    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	cmd = `
	    CREATE TABLE IF NOT EXISTS value_history
	    (
	      id              	int unsigned NOT NULL auto_increment,
	      user_id          	int unsigned NOT NULL,
	      net_worth		float(8) NOT NULL,
	      added				varchar(255) NOT NULL,
	      PRIMARY KEY     	(id)
	    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
}

func getUserIdFromToken(token string) (int, error) {
	rows, err := db.Query("SELECT user_id FROM sessions WHERE token=?", token)
	checkErr(err)
	defer rows.Close()
	var id int
	rows.Next()
	err = rows.Scan(&id)
	return id, err
}

func getUserId(email string) (int, error) {
	rows, err := db.Query("SELECT id FROM userinfo WHERE email=?", email)
	checkErr(err)
	defer rows.Close()
	var id int
	rows.Next()
	err = rows.Scan(&id)
	return id, err
}

func emailExists(email string) bool {
	var exists bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM userinfo WHERE email=?", email).Scan(&exists)
	checkErr(err)
	return exists
}

func userExists(email string) bool {
	var exists bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM userinfo WHERE email=?", email).Scan(&exists)
	checkErr(err)
	return exists
}

func tokenToUserId(token string) string {
	var userId string
	err := db.QueryRow("SELECT user_id FROM sessions WHERE token=?", token).Scan(&userId)
	if err == sql.ErrNoRows {
		return ""
	}
	checkErr(err)
	return userId
}

func getCash(userId string) float64 {
	var cash float64
	err := db.QueryRow("SELECT cash FROM portfolio WHERE user_id=?", userId).Scan(&cash)
	checkErr(err)
	return cash
}

func setCash(userId string, cash float64) error {
	_, err := db.Exec("UPDATE portfolio SET cash = ? WHERE user_id = ?", cash, userId)
	return err
}

func getUnitsOwned(userId string, symbol string) int {
	var unitsOwned int
	err := db.QueryRow("SELECT units FROM positions WHERE symbol = ? AND user_id = ?", symbol, userId).Scan(&unitsOwned)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		return -1
	}
	return unitsOwned
}

func updateUnitsOwned(userId string, req orderRequest, buying bool) {
	unitsOwned := getUnitsOwned(userId, req.Symbol)
	var err error
	orderModulator := 1
	if !buying {
		orderModulator *= -1
	}
	if unitsOwned == -1 {
		panic("Could not get units owned")
	}
	if unitsOwned == 0 && !buying {
		panic("Cannot sell with no inventory")
	}
	if unitsOwned == 0 {
		_, err = db.Exec("INSERT positions SET user_id = ?, symbol = ?, units = ?", userId, strings.ToUpper(req.Symbol), req.Units)
	} else {
		unitsOwned += req.Units * orderModulator
		_, err = db.Exec("UPDATE positions SET units = ? WHERE user_id = ? AND symbol = ?", unitsOwned, userId, req.Symbol)
	}
	checkErr(err)
}

func createOrder(userId string, req orderRequest, price float64, buying int) error {
	_, err := db.Exec("INSERT order_history SET user_id = ?, symbol = ?, units = ?, price = ?, buy = ?, added = ?", userId, strings.ToUpper(req.Symbol), req.Units, price, buying, time.Now().Unix())
	return err
}
