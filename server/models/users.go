package models

type UserInfo struct {
	Username  string   `json:"username"`
	AvatarUrl string   `json:"avatarUrl"`
	Roles     []string `json:"roles"`
}
