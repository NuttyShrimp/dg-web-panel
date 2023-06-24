package business

import (
	"degrens/panel/internal/api"
	"degrens/panel/lib/cache"
	"time"

	"github.com/sirupsen/logrus"
)

var permissionCache = cache.InitCache[string, uint](1 * time.Hour)

func refetchBusinessPerms() bool {
	perms := make(map[string]uint)
	ei, err := api.CfxApi.DoRequest("GET", "/business/permissions", nil, &perms)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch business permissions")
		return false
	}
	if ei.Message != "" {
		logrus.WithField("msg", ei.Message).Error("Failed to fetch business permissions")
		return false
	}
	for perm, mask := range perms {
		permissionCache.AddEntry(mask, perm)
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
