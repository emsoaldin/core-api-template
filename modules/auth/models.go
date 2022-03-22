package auth

import "time"

// AuthProvider model
type AuthProvider struct {
	UserID    uint64
	Provider  string
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
