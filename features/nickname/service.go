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

	"github.com/redis/go-redis/v9"
)

var (
	adjectives  []string
	colors      []string
	nouns       []string
	rdb         *redis.Client
	ctx         = context.Background()
	initialized bool
)

func InitRedis(client *redis.Client) error {
	if client == nil {
		return errors.New("redis client cannot be nil")
	}

	rdb = client

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Load word lists
	if err := loadWords(); err != nil {
		return fmt.Errorf("failed to load words: %w", err)
	}

	initialized = true
	log.Println("Nickname service initialized successfully")
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

	if rdb == nil {
		return nil, errors.New("redis client not initialized")
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
			taken, err := rdb.SIsMember(ctx, "used_nicknames", nick).Result()
			if err != nil {
				// If Redis is down, log but continue (graceful degradation)
				log.Printf("Redis error checking nickname: %v", err)
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
	if rdb == nil {
		return false
	}

	added, err := rdb.SAdd(ctx, "used_nicknames", nick).Result()
	if err != nil {
		log.Printf("Redis error reserving nickname: %v", err)
		return false
	}
	return added == 1
}
