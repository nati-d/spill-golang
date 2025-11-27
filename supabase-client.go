// supabase_client.go
package main

import (
	"os"

	supabase "github.com/supabase-community/supabase-go"
)

var supabaseClient *supabase.Client

// Call once at startup (in main())
func InitSupabase() {
	url := os.Getenv("SUPABASE_URL")
	anonKey := os.Getenv("SUPABASE_KEY") 

	if url == "" || anonKey == "" {
		panic("SUPABASE_URL and SUPABASE_KEY must be set in .env")
	}

	// NewClient returns only ONE value in the current version
	client, err := supabase.NewClient(url, anonKey, nil)
	if err != nil {
		panic("Failed to create Supabase client: " + err.Error())
	}

	supabaseClient = client
}

// Safe getter — call anywhere in your code
func Supabase() *supabase.Client {
	if supabaseClient == nil {
		panic("Supabase client not initialized — call InitSupabase() first")
	}
	return supabaseClient
}