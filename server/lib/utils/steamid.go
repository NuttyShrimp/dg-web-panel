package utils

import "strings"

func ValidateSteamId(steamid string) bool {
	return strings.HasPrefix(steamid, "steam:")
}
