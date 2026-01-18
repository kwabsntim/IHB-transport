package utils

import (
	"fmt"
	"sync"
	"time"
)
type RateLimiter struct{
	tokens chan struct{}
	refillTicker *time.Ticker
}
// EmailRateLimiter manages email sending rate limits
type EmailRateLimiter struct {
	mu                sync.Mutex
	dailyLimit        int
	sentToday         int
	lastResetDate     string
	clientEmailCounts map[string]int // Track per-client email counts
	clientDailyLimit  int            // Limit per client per day
}

var (
	globalLimiter *EmailRateLimiter
	once          sync.Once
)

// GetEmailRateLimiter returns the singleton instance
func GetEmailRateLimiter() *EmailRateLimiter {
	once.Do(func() {
		globalLimiter = &EmailRateLimiter{
			dailyLimit:        300, // Gmail allows 500/day, we use 300 to be safe
			sentToday:         0,
			lastResetDate:     time.Now().Format("2006-01-02"),
			clientEmailCounts: make(map[string]int),
			clientDailyLimit:  10, // Max 10 emails per client per day (prevents abuse)
		}
	})
	return globalLimiter
}

// CanSendEmail checks if we can send an email within rate limits
func (rl *EmailRateLimiter) CanSendEmail(clientEmail string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Reset counters if it's a new day
	today := time.Now().Format("2006-01-02")
	if today != rl.lastResetDate {
		rl.sentToday = 0
		rl.clientEmailCounts = make(map[string]int)
		rl.lastResetDate = today
		fmt.Printf("📊 Rate limiter reset for new day: %s\n", today)
	}

	// Check global daily limit
	if rl.sentToday >= rl.dailyLimit {
		return fmt.Errorf("daily email limit reached (%d/%d)", rl.sentToday, rl.dailyLimit)
	}

	// Check per-client limit
	clientCount := rl.clientEmailCounts[clientEmail]
	if clientCount >= rl.clientDailyLimit {
		return fmt.Errorf("email limit reached for %s (%d/%d emails today)",
			clientEmail, clientCount, rl.clientDailyLimit)
	}

	return nil
}

// RecordEmailSent records that an email was sent successfully
func (rl *EmailRateLimiter) RecordEmailSent(clientEmail string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.sentToday++
	rl.clientEmailCounts[clientEmail]++

	fmt.Printf("📧 Email sent (%d/%d today) to %s (%d/%d for this client)\n",
		rl.sentToday, rl.dailyLimit,
		clientEmail, rl.clientEmailCounts[clientEmail], rl.clientDailyLimit)
}

// GetStats returns current rate limiter statistics
func (rl *EmailRateLimiter) GetStats() (int, int, string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	return rl.sentToday, rl.dailyLimit, rl.lastResetDate
}

// SetDailyLimit allows changing the daily limit (useful for scaling)
func (rl *EmailRateLimiter) SetDailyLimit(newLimit int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.dailyLimit = newLimit
	fmt.Printf("📊 Daily email limit updated to: %d\n", newLimit)
}

// SetClientDailyLimit allows changing per-client limit
func (rl *EmailRateLimiter) SetClientDailyLimit(newLimit int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.clientDailyLimit = newLimit
	fmt.Printf("📊 Per-client daily limit updated to: %d\n", newLimit)
}

