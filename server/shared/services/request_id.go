package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"jakes-resume-compiler/server/shared/config"
)

func NewRequestID() (reqID string) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		reqID = "buy-a-lottery-ticket-bro"
	} else {
		reqID = hex.EncodeToString(bytes)
	}
	return reqID
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(config.RequestIDKey).(string); ok {
		return id
	}
	return ""
}
