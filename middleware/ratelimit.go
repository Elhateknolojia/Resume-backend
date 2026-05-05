package middleware

import (
    
    "net/http"
    "sync"
    "time"
)

// TokenBucket represents a per-user limiter
type TokenBucket struct {
    tokens     int
    lastRefill time.Time
}

// Global state
var (
    mu           sync.Mutex
    userBuckets  = make(map[string]*TokenBucket)
    globalTokens = 10000
    globalLast   = time.Now()
)

// refillTokens refills per-user and global buckets every minute
func refillTokens(bucket *TokenBucket, limit int) {
    if time.Since(bucket.lastRefill) >= time.Minute {
        bucket.tokens = limit
        bucket.lastRefill = time.Now()
    }
}

func refillGlobal(limit int) {
    if time.Since(globalLast) >= time.Minute {
        globalTokens = limit
        globalLast = time.Now()
    }
}

// middleware enforcing token bucket limits
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        mu.Lock()
        defer mu.Unlock()

        userID := r.Context().Value("user_id").(string) // replace with session ID or auth token

        // Initialize bucket if new user
        if _, ok := userBuckets[userID]; !ok {
            userBuckets[userID] = &TokenBucket{tokens: 40, lastRefill: time.Now()}
        }

        // Refill user + global buckets
        refillTokens(userBuckets[userID], 40)
        refillGlobal(10000)

        // Check limits
        if userBuckets[userID].tokens <= 0 {
            http.Error(w, "User limit exceeded (40/min)", http.StatusTooManyRequests)
            return
        }
        if globalTokens <= 0 {
            http.Error(w, "System limit exceeded (10000/min)", http.StatusTooManyRequests)
            return
        }

        // Consume tokens
        userBuckets[userID].tokens--
        globalTokens--

        next.ServeHTTP(w, r)
    })
}