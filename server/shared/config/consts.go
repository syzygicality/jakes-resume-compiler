package config

type contextKey string

const RequestIDKey contextKey = "requestID"

const RequestIDHeader = "X-Request-ID"

const APIKeyHeader = "X-API-Key"
