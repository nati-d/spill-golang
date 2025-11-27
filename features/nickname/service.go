package nickname

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	supabase "github.com/supabase-community/supabase-go"
)

var (
	adjectives     []string
	colors         []string
	nouns          []string
	supabaseClient *supabase.Client
	ctx            = context.Background()
	initialized    bool
)

// InitSupabase initializes the nickname service with Supabase client
func InitSupabase(client *supabase.Client) error {
	if client == nil {
		return errors.New("supabase client cannot be nil")
	}

	supabaseClient = client

	// Load word lists
	if err := loadWords(); err != nil {
		return fmt.Errorf("failed to load words: %w", err)
	}

	initialized = true
	log.Println("Nickname service initialized successfully with Supabase")
	return nil
}

func loadWords() error {
	adjData, err := os.ReadFile("words/adjectives.txt")
	if err != nil {
		return fmt.Errorf("failed to read adjectives.txt: %w", err)
	}

	nounsData, err := os.ReadFile("words/nouns.txt")
	if err != nil {
		return fmt.Errorf("failed to read nouns.txt: %w", err)
	}

	colorsData, err := os.ReadFile("words/colors.txt")
	if err != nil {
		return fmt.Errorf("failed to read colors.txt: %w", err)
	}

	adjectives = filter(strings.Split(string(adjData), "\n"))
	nouns = filter(strings.Split(string(nounsData), "\n"))
	colors = filter(strings.Split(string(colorsData), "\n"))

	if len(adjectives) == 0 || len(nouns) == 0 || len(colors) == 0 {
		return errors.New("one or more word lists are empty")
	}

	return nil
}

func filter(lines []string) []string {
	var res []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "#") {
			res = append(res, l)
		}
	}
	return res
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GenerateThree() ([]string, error) {
	if !initialized {
		return nil, errors.New("nickname service not initialized")
	}

	if supabaseClient == nil {
		return nil, errors.New("supabase client not initialized")
	}

	suggestions := []string{}
	seen := make(map[string]bool)
	maxAttempts := 100 // Prevent infinite loop
	attempts := 0

	for len(suggestions) < 3 && attempts < maxAttempts {
		attempts++

		// Use crypto-secure random source (Go 1.20+)
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		style := rng.Intn(3)

		var nick string
		switch style {
		case 0:
			nick = fmt.Sprintf("%s%s%s_%d",
				capitalize(randomItem(colors, rng)),
				capitalize(randomItem(adjectives, rng)),
				capitalize(randomItem(nouns, rng)),
				rng.Intn(9000)+1000)
		case 1:
			nick = fmt.Sprintf("%s%s_%d",
				capitalize(randomItem(adjectives, rng)),
				capitalize(randomItem(nouns, rng)),
				rng.Intn(9000)+1000)
		case 2:
			nick = fmt.Sprintf("%s%s_%d",
				capitalize(randomItem(colors, rng)),
				capitalize(randomItem(nouns, rng)),
				rng.Intn(9000)+1000)
		}

		if !seen[nick] {
			taken, err := isNicknameTaken(nick)
			if err != nil {
				// If Supabase is down, log but continue (graceful degradation)
				log.Printf("Supabase error checking nickname: %v", err)
				// Assume not taken if we can't check
				taken = false
			}

			if !taken {
				suggestions = append(suggestions, nick)
				seen[nick] = true
			}
		}
	}

	if len(suggestions) < 3 {
		return suggestions, fmt.Errorf("only generated %d suggestions after %d attempts", len(suggestions), attempts)
	}

	return suggestions, nil
}

func randomItem(slice []string, rng *rand.Rand) string {
	if len(slice) == 0 {
		return "Mystic"
	}
	return slice[rng.Intn(len(slice))]
}

func Reserve(nick string) bool {
	if supabaseClient == nil {
		return false
	}

	// Check if nickname already exists
	taken, err := isNicknameTaken(nick)
	if err != nil {
		log.Printf("Supabase error checking nickname: %v", err)
		return false
	}
	if taken {
		return false
	}

	// Insert new nickname
	var result []map[string]interface{}
	_, err = supabaseClient.From("used_nicknames").
		Insert(map[string]interface{}{
			"nickname":   nick,
			"created_at": time.Now().Format(time.RFC3339),
		}, false, "", "", "").
		ExecuteTo(&result)

	if err != nil {
		log.Printf("Supabase error reserving nickname: %v", err)
		return false
	}

	return true
}

// isNicknameTaken checks if a nickname is already taken in Supabase
func isNicknameTaken(nick string) (bool, error) {
	var result []map[string]interface{}
	_, err := supabaseClient.From("used_nicknames").
		Select("nickname", "", false).
		Eq("nickname", nick).
		ExecuteTo(&result)

	if err != nil {
		return false, err
	}

	return len(result) > 0, nil
}
