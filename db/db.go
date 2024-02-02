package db

import (
	"github.com/zenpk/safe-csv-db"
)

type Db struct {
	TableUser         *scd.Table
	TableRefreshToken *scd.Table
	TableClient       *scd.Table
}

func (d *Db) Init(preparing chan<- struct{}, stop <-chan struct{}) error {
	var err error

	d.TableUser, err = scd.OpenTable("./db/user.csv", User{})
	if err != nil {
		return err
	}
	defer d.TableUser.Close()
	go func() {
		if err := d.TableUser.ListenChange(); err != nil {
			panic(err)
		}
	}()

	d.TableClient, err = scd.OpenTable("./db/client.csv", Client{})
	if err != nil {
		return err
	}
	defer d.TableClient.Close()
	go func() {
		if err := d.TableClient.ListenChange(); err != nil {
			panic(err)
		}
	}()

	d.TableRefreshToken, err = scd.OpenTable("./db/refresh_token.csv", RefreshToken{})
	if err != nil {
		return err
	}
	defer d.TableRefreshToken.Close()
	go func() {
		if err := d.TableRefreshToken.ListenChange(); err != nil {
			panic(err)
		}
	}()

	close(preparing)
	<-stop
	return nil
}
