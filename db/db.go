package db

import (
	"github.com/zenpk/safe-csv-db"
)

var TableUser *scd.Table
var TableRefreshToken *scd.Table
var TableClient *scd.Table

func Init(prepared, done chan struct{}) error {
	var err error

	TableUser, err = scd.OpenTable("./db/user.csv", User{})
	if err != nil {
		return err
	}
	defer TableUser.Close()
	go func() {
		if err := TableUser.ListenChange(); err != nil {
			panic(err)
		}
	}()

	TableClient, err = scd.OpenTable("./db/client.csv", Client{})
	if err != nil {
		return err
	}
	defer TableClient.Close()
	go func() {
		if err := TableClient.ListenChange(); err != nil {
			panic(err)
		}
	}()

	TableRefreshToken, err = scd.OpenTable("./db/refresh_token.csv", RefreshToken{})
	if err != nil {
		return err
	}
	defer TableRefreshToken.Close()
	go func() {
		if err := TableRefreshToken.ListenChange(); err != nil {
			panic(err)
		}
	}()

	prepared <- struct{}{}
	<-done
	return nil
}
