package main

import "log"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
