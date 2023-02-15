package discord

import (
	"crypto/rand"
	"degrens/panel/internal/config"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	dgerrors "degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type DiscordIdentity struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	// To get actual url: https://discord.com/developers/docs/reference#image-formatting
	AvatarHash string   `json:"avatar"`
	Roles      []uint64 `json:"string,omitempty"`
}

type discordGuildMember struct {
	Roles []string `json:"roles"`
}

type discordInfo struct {
	Conf    *oauth2.Config
	GuildId string
	Roles   []config.ConfigRole
}

type DiscordToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

var info *discordInfo
var logger log.Logger

func InitDiscordConf(conf *config.Config, pLogger log.Logger) {
	info = &discordInfo{

		Conf: &oauth2.Config{
			RedirectURL: conf.Discord.RedirectURL,
			// This next 2 lines must be edited before running this.
			ClientID:     conf.Discord.ClientID,
			ClientSecret: conf.Discord.ClientSecret,
			Scopes: []string{
				"identify",
				"guilds.members.read",
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://discord.com/api/oauth2/authorize",
				TokenURL:  "https://discord.com/api/oauth2/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
		GuildId: conf.Discord.Guild,
		Roles:   conf.Discord.Roles,
	}
	logger = pLogger.With("module", "discord")
}

func GetOAuthConf() *oauth2.Config {
	return info.Conf
}

func GenerateOAuthState(c *gin.Context) (string, error) {
	// Generate a unique state string
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", n), nil
}

func fetchUserIdentity(token *oauth2.Token) (*DiscordIdentity, error) {
	var user DiscordIdentity
	err := fetchFromDiscordAPI(token, "users/@me", &user)

	if err != nil {
		logger.Error("Error while reading user info", "error", err)
		return nil, errors.New("error while getting user info")
	}

	return &user, nil
}

func fetchGuildInfo(token *oauth2.Token) (*discordGuildMember, error) {
	var member discordGuildMember
	err := fetchFromDiscordAPI(token, fmt.Sprintf("users/@me/guilds/%s/member", info.GuildId), &member)

	if err != nil {
		logger.Error("Error while reading member info", "error", err.Error())
		return nil, errors.New("error while getting member info")
	}

	return &member, nil
}

func GetUserInfoViaToken(token *oauth2.Token) (*DiscordIdentity, error) {
	user, err := fetchUserIdentity(token)
	if err != nil {
		return nil, err
	}
	guildMember, err := fetchGuildInfo(token)
	if err != nil {
		return nil, err
	}

	for _, role := range guildMember.Roles {
		number, _ := strconv.ParseUint(role, 10, 64)
		user.Roles = append(user.Roles, number)
	}

	return user, nil
}

func RevokeAuthToken(token string) error {
	params := url.Values{}
	params.Add("token", token)
	resp, err := http.PostForm("https://discord.com/api/oauth2/token/revoke", params)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()
	return err
}

func RefreshToken(userId uint) {
	var DBToken panel_models.DiscordTokens
	db.MariaDB.Client.First(&DBToken, "user_id = ?", userId)
	params := url.Values{}
	params.Set("client_id", info.Conf.ClientID)
	params.Set("client_secret", info.Conf.ClientSecret)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", DBToken.RefreshToken)
	res, err := http.PostForm(info.Conf.Endpoint.TokenURL, params)
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to refresh token for user %d", userId), "error", err.Error())
		return
	}
	DCToken := DiscordToken{}
	if err := json.NewDecoder(res.Body).Decode(&DCToken); err != nil {
		dgerrors.HandleJsonError(err, logger)
		return
	}
	if DCToken.AccessToken == "" {
		return
	}
	Token := oauth2.Token{}
	Token.AccessToken = DCToken.AccessToken
	Token.RefreshToken = DCToken.RefreshToken
	Token.Expiry = time.Now().Add(time.Second * time.Duration(DCToken.ExpiresIn))
	AssignTokenToUser(userId, &Token)
}
