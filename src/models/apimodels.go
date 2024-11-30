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

type InvalidInputError struct {
	Message string
}

func (i *InvalidInputError) Error() string {
	return i.Message
}
