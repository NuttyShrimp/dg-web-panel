package models

type CfxPlayer struct {
	Name      string   `json:"name"`
	SteamID   string   `json:"steamId"`
	DiscordID string   `json:"discordId"`
	Roles     []string `json:"roles"`
}

type CfxCharacter struct {
	SteamID   string `json:"steamId"`
	Cid       uint   `json:"cid"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type CfxBankPermissions struct {
	Deposit      bool `json:"deposit"`
	Withdraw     bool `json:"withdraw"`
	Transfer     bool `json:"transfer"`
	Transactions bool `json:"transactions"`
}

type CfxReportAnnouncement struct {
	ID    uint     `json:"id"`
	Recvs []string `json:"receivers"`
}
