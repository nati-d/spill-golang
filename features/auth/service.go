package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

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

	// Convert telegram_data to JSON string for storage
	telegramDataJSON, err := json.Marshal(tgUser)
	if err != nil {
		return profile, nil, fmt.Errorf("failed to marshal telegram data: %w", err)
	}

	data := map[string]interface{}{
		"telegram_id":       tgUser.ID,
		"telegram_username": tgUser.Username,
		"telegram_data":     string(telegramDataJSON),
	}

	log.Printf("Attempting to upsert profile for telegram_id: %d", tgUser.ID)

	// Try to find existing profile first
	var existingProfiles []Profile
	_, err = s.client.From("profiles").
		Select("*", "", false).
		Eq("telegram_id", fmt.Sprintf("%d", tgUser.ID)).
		ExecuteTo(&existingProfiles)

	if err != nil {
		log.Printf("Error checking existing profile: %v", err)
		// Continue to try upsert anyway
	}

	if len(existingProfiles) > 0 {
		// Update existing profile
		profile = existingProfiles[0]
		_, _, err = s.client.From("profiles").
			Update(data, "", "").
			Eq("telegram_id", fmt.Sprintf("%d", tgUser.ID)).
			Execute()
		if err != nil {
			log.Printf("Error updating profile: %v", err)
			return profile, nil, fmt.Errorf("failed to update profile: %w", err)
		}
		// Re-fetch to get updated data
		_, err = s.client.From("profiles").
			Select("*", "", false).
			Eq("telegram_id", fmt.Sprintf("%d", tgUser.ID)).
			ExecuteTo(&existingProfiles)
		if err == nil && len(existingProfiles) > 0 {
			profile = existingProfiles[0]
		}
	} else {
		// Insert new profile
		var newProfiles []Profile
		_, err = s.client.From("profiles").
			Insert(data, false, "", "", "").
			ExecuteTo(&newProfiles)
		if err != nil {
			log.Printf("Error inserting profile: %v", err)
			return profile, nil, fmt.Errorf("failed to insert profile: %w", err)
		}
		if len(newProfiles) > 0 {
			profile = newProfiles[0]
		}
	}

	if profile.Nickname == "" {
		suggestions, err := nickname.GenerateThree()
		if err != nil {
			log.Printf("Error generating suggestions: %v", err)
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
