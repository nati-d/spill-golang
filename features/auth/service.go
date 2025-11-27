package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/nati-d/spill-backend/features/nickname"
	"github.com/nati-d/spill-backend/internal/telegram"
	supabase "github.com/supabase-community/supabase-go"
)

type authService struct {
	client *supabase.Client
}

var AuthService *authService

// InitAuthService initializes the auth service with Supabase client
func InitAuthService(client *supabase.Client) {
	AuthService = &authService{client: client}
}

func (s *authService) LoginOrRegister(ctx context.Context, tgUser telegram.TelegramUser) (Profile, []string, error) {
	var profile Profile

	data := map[string]interface{}{
		"telegram_id":       tgUser.ID,
		"telegram_username": tgUser.Username,
		"telegram_data":     tgUser,
	}

	_, err := s.client.From("profiles").
		Upsert(data, "", "", "").
		ExecuteTo(&profile)
	if err != nil {
		return profile, nil, fmt.Errorf("failed to upsert profile: %w", err)
	}

	if profile.Nickname == "" {
		suggestions, err := nickname.GenerateThree()
		if err != nil {
			return profile, nil, fmt.Errorf("failed to generate suggestions: %w", err)
		}
		return profile, suggestions, nil
	}

	return profile, nil, nil
}

func (s *authService) SetNickname(ctx context.Context, telegramID int64, nick string) error {
	if !nickname.Reserve(nick) {
		return errors.New("nickname taken")
	}

	_, _, err := s.client.From("profiles").
		Update(map[string]interface{}{"nickname": nick}, "", "").
		Eq("telegram_id", fmt.Sprintf("%d", telegramID)).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to update nickname: %w", err)
	}

	return nil
}
