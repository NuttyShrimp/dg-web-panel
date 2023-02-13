package apikeys

import (
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/lib/graylogger"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func generateKey() string {
	for {
		id := uuid.NewString()
		var count int64
		db.MariaDB.Client.Table("api_keys").Where("api_key = ?", id).Count(&count)
		if count == 0 {
			return id
		}
	}
}

func filterKeys(keys []*panel_models.APIKey) []*panel_models.APIKey {
	filteredKeys := []*panel_models.APIKey{}
	for _, key := range keys {
		if key.Expired() {
			DeleteAPIKey(0, key.ApiKey)
		} else {
			filteredKeys = append(filteredKeys, key)
		}
	}
	return filteredKeys
}

func GetAPIKey(key string) *panel_models.APIKey {
	var entry panel_models.APIKey
	db.MariaDB.Client.Model(&panel_models.APIKey{}).Preload("User").Where("api_key = ?", key).First(&entry)
	if entry.Expired() {
		DeleteAPIKey(0, key)
		return nil
	}
	return &entry
}

func GetAPIKeyForUser(key string, userId uint) *panel_models.APIKey {
	var entry panel_models.APIKey
	db.MariaDB.Client.Preload("User").Where("api_key = ? AND user_id = ?", key, userId).First(&entry)
	if entry.Expired() {
		DeleteAPIKey(0, key)
		return nil
	}
	return &entry
}

func GetAPIKeys(userId uint) []*panel_models.APIKey {
	keys := []*panel_models.APIKey{}
	db.MariaDB.Client.Preload("User").Where("user_id = ?", userId).Find(&keys)
	return filterKeys(keys)
}

// TODO: Test if gorm can read into slice of pointers
func GetAllAPIKeys() []*panel_models.APIKey {
	keys := []*panel_models.APIKey{}
	db.MariaDB.Client.Preload("User").Find(&keys)
	return filterKeys(keys)
}

func CreateAPIKey(userId uint, comment string, duration time.Duration) (string, error) {
	id := generateKey()
	key := panel_models.APIKey{
		ApiKey:  id,
		UserID:  userId,
		Comment: comment,
		Expiry:  time.Now().Add(duration),
	}
	err := db.MariaDB.Client.Create(&key).Error
	graylogger.Log("apikeys:create", fmt.Sprintf("created a new API key for user %d", userId), "key", key, "comment", comment)
	return id, err
}

func DeleteAPIKey(userId uint, apiKey string) {
	graylogger.Log("apikeys:delete", fmt.Sprintf("user %d deleted the an API key", userId), "key", apiKey)
	db.MariaDB.Client.Where("api_key = ?", apiKey).Delete(&panel_models.APIKey{})
}

func DeleteAPIKeys(userId uint, apiKeys []string) {
	graylogger.Log("apikeys:batch_delete", fmt.Sprintf("user %d deleted multiple API keys", userId), "keys", apiKeys)
	db.MariaDB.Client.Where("api_key IN ?", apiKeys).Delete(&panel_models.APIKey{})
}
