package middleware

import (
	"sync"

	"github.com/didip/tollbooth/limiter"
)

var (
	once              sync.Once
	allowedOriginsMap map[string]struct{}
	limit             *limiter.Limiter
)
