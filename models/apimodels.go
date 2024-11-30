package models

import "time"

type RegisterRequest struct {
	Username   string `json:"username"`
	Key        string `json:"key"`
	ClientSalt string `json:"clientSalt"`
}

type UserResponse struct {
	Salt string `json:"clientSalt"`
}

type VaultResponse struct {
	ModifiedAt time.Time `json:"modifiedAt"`
}
