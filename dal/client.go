package dal

import (
	"database/sql"
	"errors"
	"time"
)

type IClient interface {
	Init() error
	Insert(client *Client) error
	SelectById(id int64) (*Client, error)
	SelectByClientId(clientId string) (*Client, error)
	DeleteById(id int64) error
}

type Client struct {
	db              *sql.DB
	Id              int64
	ClientId        string
	ClientSecret    string
	Redirects       string
	AccessTokenAge  time.Duration
	RefreshTokenAge time.Duration
}

func (c Client) Init() error {
	if _, err := c.db.Exec(`
	CREATE TABLE IF NOT EXISTS clients (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    client_id TEXT NOT NULL UNIQUE,
	    client_secret TEXT NOT NULL,
	    redirects TEXT NOT NULL,
	    access_token_age INTEGER NOT NULL,
	    refresh_token_age INTEGER NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0
	);`); err != nil {
		return err
	}
	rows, err := c.db.Query(`SELECT * FROM sqlite_master WHERE type = "index" AND tbl_name = "clients" AND name = "idx_client_id";`)
	if err != nil {
		return err
	}
	if !rows.Next() {
		if _, err = c.db.Exec(`CREATE UNIQUE INDEX idx_client_id ON clients(client_id);`); err != nil {
			return err
		}
	}
	return rows.Close()
}

func (c Client) Insert(client *Client) error {
	_, err := c.db.Exec("INSERT INTO clients (client_id, client_secret, redirects, access_token_age, refresh_token_age) VALUES (?, ?, ?, ?, ?);",
		client.ClientId, client.ClientSecret, client.Redirects, client.AccessTokenAge, client.RefreshTokenAge)
	return err
}

func (c Client) SelectById(id int64) (client *Client, err error) {
	rows, err := c.db.Query("SELECT * FROM clients WHERE (id = ? AND deleted = 0);", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()
	client = new(Client)
	if rows.Next() {
		if err := rows.Scan(&client.Id, &client.ClientId, &client.ClientSecret, &client.Redirects, &client.AccessTokenAge, &client.RefreshTokenAge); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return client, err
}

func (c Client) SelectByClientId(clientId string) (client *Client, err error) {
	rows, err := c.db.Query("SELECT * FROM clients WHERE (client_id = ? AND deleted = 0);", clientId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()
	client = new(Client)
	if rows.Next() {
		if err := rows.Scan(&client.Id, &client.ClientId, &client.ClientSecret, &client.Redirects, &client.AccessTokenAge, &client.RefreshTokenAge); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return client, err
}

func (c Client) DeleteById(id int64) error {
	_, err := c.db.Exec("UPDATE clients SET deleted = 1 WHERE id = ?;", id)
	return err
}
