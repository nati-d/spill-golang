package auth

type Profile struct {
  ID         string `json:"id"`
  TelegramID int64  `json:"telegram_id"`
  Nickname   string `json:"nickname"`
  AvatarURL  string `json:"avatar_url,omitempty"`
  CreatedAt  string `json:"created_at"`
}