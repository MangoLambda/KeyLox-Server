package models

import "time"

type DBUser struct {
	ID         int       `db:"id"`
	Username   string    `db:"username"`
	ClientSalt string    `db:"client_salt"`
	ServerSalt string    `db:"server_salt"`
	HashedKey  string    `db:"hashed_key"`
	LastLogin  time.Time `db:"last_login"`
	VaultId    int       `db:"vault_id"`
}

type DBVault struct {
	ID        int       `db:"id"`
	FileName  string    `db:"filename"`
	CreatedAt time.Time `db:"created_at"`
	UserId    int       `db:"user_id"`
}
