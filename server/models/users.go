package models

type UserInfo struct {
	ID        uint     `json:"-"`
	Username  string   `json:"username"`
	AvatarUrl string   `json:"avatarUrl"`
	Roles     []string `json:"roles"`
}
