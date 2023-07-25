package main

import (
	"github.com/zenpk/my-oauth/db"
	"log"
)

func main() {
	done := make(chan struct{})
	prepared := make(chan struct{})
	if err := db.Init(prepared, done); err != nil {
		log.Fatalln(err)
	}
	<-prepared
	
}
