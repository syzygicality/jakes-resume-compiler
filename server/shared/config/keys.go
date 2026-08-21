package config

import (
	"time"
)

type contextKey string

const PreambleFormat = "preamble"

const DocumentMarker = `\begin{document}`

const CompileTimeout = 5 * time.Second

const APIKeyHeader = "X-API-Key"

const RequestIDHeader = "X-Request-ID"

const RequestIDKey contextKey = "requestID"
