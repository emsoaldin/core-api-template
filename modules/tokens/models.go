package tokens

import "time"

// Token model
type Token struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"userId" binding:"required"`
	Token       string    `json:"token" binding:"required"`
	TokenTypeID uint64    `json:"tokenTypeId" binding:"required"`
	Meta        string    `json:"meta"`
	ExpiresAt   time.Time `json:"expiresAt" binding:"required"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
