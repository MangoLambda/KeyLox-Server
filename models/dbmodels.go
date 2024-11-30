package models

import "time"

type DBUser struct {
	ID         int       `db:"id"`
	Username   string    `db:"username"`
	Key        string    `db:"key"`
	ClientSalt string    `db:"clientSalt"`
	ServerSalt string    `db:"serverSalt"`
	HashedKey  string    `db:"hashedKey"`
	VaultId    int       `db:"vaultId"`
	LastLogin  time.Time `db:"lastLogin"`
}

type DBVault struct {
	ID        int       `db:"id"`
	UserId    int       `db:"userId"`
	FileName  string    `db:"fileName"`
	CreatedAt time.Time `db:"createdAt"`
}
