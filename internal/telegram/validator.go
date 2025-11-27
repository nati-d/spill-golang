package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

type TelegramUser struct {
	ID int64 `json:"id"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	PhotoURL string `json:"photo_url"`
}

func ValidateInitData(initData string) (TelegramUser, error) {
	var user TelegramUser

	params, err := url.ParseQuery(initData)
	if err != nil {
		return user, err
	}

	hash := params.Get("hash")
	params.Del("hash")

	var dataCheckString []string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		dataCheckString = append(dataCheckString, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}
	

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(os.Getenv("TELEGRAM_BOT_TOKEN")))
	hmacKey := secretKey.Sum(nil)

	h := hmac.New(sha256.New, hmacKey)
	h.Write([]byte(strings.Join(dataCheckString, "\n")))

	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash != hash {
		return user, errors.New("invalid hash")
	}

	userJson := params.Get("user")
	if err := json.Unmarshal([]byte(userJson), &user); err != nil {
		return user, err
	}

	return user, nil
}

func mustInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}