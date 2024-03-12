package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// create Token creates a new token for a specific username and  duration of  time
	CreateToken(username string, duration time.Duration) (string, error)
	//verify token
	VerifyToken(token string) (*Payload, error)
}
