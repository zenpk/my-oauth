package dal

import (
	"database/sql"
	"errors"
	"time"
)

type IRefreshToken interface {
	Init() error
	Insert(token *RefreshToken) error
	SelectByToken(token string) (*RefreshToken, error)
	CleanExpired() error
	DeleteById(id int64) error
}

type RefreshToken struct {
	db         *sql.DB
	Id         int64
	Token      string
	ClientId   int64
	UserId     int64
	ExpireTime *time.Time
	Deleted    bool
}

func (r RefreshToken) Init() error {
	if _, err := r.db.Exec(`
	CREATE TABLE IF NOT EXISTS refresh_tokens (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    token TEXT NOT NULL,
	    client_id INTEGER NOT NULL,
	    user_id INTEGER NOT NULL,
		expire_time INTEGER NOT NULL,
		deleted INTEGER NOT NULL DEFAULT 0
	);`); err != nil {
		return err
	}
	// create index if not exists
	rows, err := r.db.Query(`SELECT * FROM sqlite_master WHERE type = "index" AND tbl_name = "refresh_tokens" AND name = "idx_token";`)
	if err != nil {
		return err
	}
	if !rows.Next() {
		if _, err = r.db.Exec(`CREATE INDEX idx_token ON refresh_tokens(token);`); err != nil {
			return err
		}
	}
	if err := rows.Close(); err != nil {
		return err
	}
	rows, err = r.db.Query(`SELECT * FROM sqlite_master WHERE type = "index" AND tbl_name = "refresh_tokens" AND name = "idx_expire_time";`)
	if err != nil {
		return err
	}
	if !rows.Next() {
		if _, err = r.db.Exec(`CREATE INDEX idx_expire_time ON refresh_tokens(expire_time);`); err != nil {
			return err
		}
	}
	return rows.Close()
}

func (r RefreshToken) Insert(token *RefreshToken) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (token, client_id, user_id, expire_time) VALUES (?, ?, ?, ?);", token.Token, token.ClientId, token.UserId, token.ExpireTime.Unix())
	return err
}

func (r RefreshToken) SelectByToken(token string) (refreshToken *RefreshToken, err error) {
	rows, err := r.db.Query("SELECT * FROM refresh_tokens WHERE (token = ? AND deleted = 0);", token)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()
	refreshToken = new(RefreshToken)
	if rows.Next() {
		var unixTime int64
		if err := rows.Scan(&refreshToken.Id, &refreshToken.Token, &refreshToken.ClientId, &refreshToken.UserId, &unixTime, &refreshToken.Deleted); err != nil {
			return nil, err
		}
		expireTime := time.Unix(unixTime, 0)
		refreshToken.ExpireTime = &expireTime
	} else {
		return nil, nil
	}
	return refreshToken, err
}

func (r RefreshToken) CleanExpired() error {
	timeNow := time.Now().Unix()
	_, err := r.db.Exec("UPDATE refresh_tokens SET deleted = 1 WHERE expire_time < ?;", timeNow)
	return err
}

func (r RefreshToken) DeleteById(id int64) error {
	_, err := r.db.Exec("UPDATE refresh_tokens SET deleted = 1 WHERE id = ?;", id)
	return err
}
