package business

import (
	"degrens/panel/internal/api"
	"degrens/panel/lib/cache"
	"time"
)

var permissionCache = cache.InitCache[string, uint](1 * time.Hour)

func refetchBusinessPerms() bool {
	perms := make(map[string]uint)
	ei, err := api.CfxApi.DoRequest("GET", "/business/permissions", nil, &perms)
	if err != nil {
		br.Logger.Error("Failed to fetch business permissions", "error", err)
		return false
	}
	if ei.Message != "" {
		br.Logger.Error("Failed to fetch business permissions", "msg", ei.Message)
		return false
	}
	for perm, mask := range perms {
		permissionCache.AddEntry(uint(mask), perm)
	}
	return true
}

func bitToPermList(bit uint) []string {
	list := []string{}
	mask := uint(0)
	for bit > (1 << mask) {
		rMask := uint(1 << mask)
		if bit&rMask == rMask {
			perm, exists := permissionCache.GetEntry(rMask)
			if !exists || perm == nil {
				if succ := refetchBusinessPerms(); !succ {
					break
				}
				perm, _ = permissionCache.GetEntry(rMask)
			}
			if perm != nil {
				list = append(list, *perm)
			}
		}
		mask++
	}
	return list
}
