package config

import (
	"time"
)

const (
	JWTSecret    = "your-secret-key" // In production, use environment variables
	JWTExpiresIn = time.Hour * 24    // 24 hours
)
