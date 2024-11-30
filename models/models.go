package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Key        string    `json:"key"`
	ClientSalt string    `json:"clientSalt"`
	ServerSalt string    `json:"-"`
	VaultId    int       `json:"-"`
	LastLogin  time.Time `json:"-"`
}

// Rework the structs. The filename and created at should not be allowed to be sent from the client.
type Vault struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userId"`
	FileName  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
}
