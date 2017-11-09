package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type mockResponse struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Street   string `json:"street"`
			City     string `json:"city"`
			State    string `json:"state"`
			Postcode int    `json:"postcode"`
		} `json:"location"`
		Email string `json:"email"`
		Login struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Salt     string `json:"salt"`
			Md5      string `json:"md5"`
			Sha1     string `json:"sha1"`
			Sha256   string `json:"sha256"`
		} `json:"login"`
		Dob        string `json:"dob"`
		Registered string `json:"registered"`
		Phone      string `json:"phone"`
		Cell       string `json:"cell"`
		ID         struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"id"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
		Nat string `json:"nat"`
	} `json:"results"`
	Info struct {
		Seed    string `json:"seed"`
		Results int    `json:"results"`
		Page    int    `json:"page"`
		Version string `json:"version"`
	} `json:"info"`
}

func setupDB(name string) *sql.DB {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	if len(dbUsername) == 0 {
		panic("$DB_USERNAME is not set")
	}
	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/", dbUsername, dbPassword))
	checkFatalErr(err)

	_, err = database.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name))
	checkFatalErr(err)
	database, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUsername, dbPassword, name))
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
	stmt, err := database.Prepare(cmd)
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
	stmt, err = database.Prepare(cmd)
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
	stmt, err = database.Prepare(cmd)
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
	stmt, err = database.Prepare(cmd)
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
	stmt, err = database.Prepare(cmd)
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
	stmt, err = database.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	return database
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
		_, err = db.Exec("INSERT positions SET user_id = ?, symbol = ?, units = ?",
			userId, strings.ToUpper(req.Symbol), req.Units)
	} else {
		unitsOwned += req.Units * orderModulator
		_, err = db.Exec("UPDATE positions SET units = ? WHERE user_id = ? AND symbol = ?",
			unitsOwned, userId, req.Symbol)
	}
	checkErr(err)
}

func createOrder(userId string, req orderRequest, price float64, buying int) error {
	_, err := db.Exec("INSERT order_history SET user_id = ?, symbol = ?, units = ?, price = ?, buy = ?, added = ?",
		userId, strings.ToUpper(req.Symbol), req.Units, price, buying, time.Now().Unix())
	return err
}

func mockData(database *sql.DB) {
	log.Printf("mocking DB mock...")
	symbols := []string{
		"TSLA",
		"AMZN",
		"FB",
		"GOOG",
		"MSFT",
		"ANET",
	}
	webClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet,
		"https://randomuser.me/api/?results=50&seed=stocks&nat=us", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "stocks-web-client")
	res, err := webClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp := mockResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().Unix())
	for i, e := range resp.Results {
		t := time.Now()
		_, err = database.Exec("INSERT userinfo SET first_name = ?, last_name = ?, email = ?, password = ?, added = ?",
			e.Name.First, e.Name.Last, e.Email, e.Login.Md5, t.Unix())
		if err != nil {
			log.Fatal(err)
		}
		numSessions := rand.Intn(10)
		for j := 0; j < numSessions; j++ {
			token, err := exec.Command("uuidgen").Output()
			token = token[0 : len(token)-1]
			_, err = database.Exec("INSERT sessions SET user_id = ?, token = ?, added = ?, expires = ?",
				i+1, token, t.Unix(), t.Add(time.Hour*24).Unix())
			if err != nil {
				log.Fatal(err)
			}
		}
		_, err = database.Exec("INSERT portfolio SET user_id = ?, cash = ?",
			i+1, e.Location.Postcode)
		if err != nil {
			log.Fatal(err)
		}
		numPositions := rand.Intn(20)
		for j := 0; j < numPositions; j++ {
			_, err = database.Exec("INSERT positions SET user_id = ?, symbol = ?, units = ?",
				i+1, symbols[rand.Intn(len(symbols))], rand.Intn(30))
			if err != nil {
				log.Fatal(err)
			}
		}
		numOrderHistory := rand.Intn(50)
		for j := 0; j < numOrderHistory; j++ {
			_, err = database.Exec("INSERT order_history SET user_id = ?, symbol = ?, units = ?, price = ?, buy = ?, added = ?",
				i+1, symbols[rand.Intn(len(symbols))], rand.Intn(30), rand.Intn(1000), rand.Intn(2), t.Unix()-int64(86400))
			if err != nil {
				log.Fatal(err)
			}
		}
		numValueHistory := rand.Intn(10)
		for j := 0; j < numValueHistory; j++ {
			_, err = database.Exec("INSERT value_history SET user_id = ?, net_worth = ?, added = ?",
				i+1, e.Location.Postcode*(rand.Intn(5)+1), t.Unix())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Printf("DB mock complete!")
}
