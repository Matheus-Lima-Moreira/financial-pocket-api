package auth

import (
	"strings"
	"sync"
	"time"
)

type AuthRateLimitAction string

const (
	AuthRateLimitRegister                AuthRateLimitAction = "register"
	AuthRateLimitLogin                   AuthRateLimitAction = "login"
	AuthRateLimitSendResetPassword       AuthRateLimitAction = "send_reset_password_email"
	AuthRateLimitResendVerificationEmail AuthRateLimitAction = "resend_verification_email"
)

type authRateLimitPolicy struct {
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	ResetAfter time.Duration
}

type authAttemptState struct {
	Attempts     int
	BlockedUntil time.Time
	LastRequest  time.Time
}

type AuthRateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*authAttemptState
	policies map[AuthRateLimitAction]authRateLimitPolicy
}

func NewAuthRateLimiter() *AuthRateLimiter {
	return &AuthRateLimiter{
		attempts: make(map[string]*authAttemptState),
		policies: map[AuthRateLimitAction]authRateLimitPolicy{
			AuthRateLimitRegister: {
				BaseDelay:  time.Minute,
				MaxDelay:   10 * time.Minute,
				ResetAfter: 24 * time.Hour,
			},
			AuthRateLimitLogin: {
				BaseDelay:  30 * time.Second,
				MaxDelay:   10 * time.Minute,
				ResetAfter: 12 * time.Hour,
			},
			AuthRateLimitSendResetPassword: {
				BaseDelay:  time.Minute,
				MaxDelay:   15 * time.Minute,
				ResetAfter: 24 * time.Hour,
			},
			AuthRateLimitResendVerificationEmail: {
				BaseDelay:  time.Minute,
				MaxDelay:   15 * time.Minute,
				ResetAfter: 24 * time.Hour,
			},
		},
	}
}

func (r *AuthRateLimiter) Allow(action AuthRateLimitAction, identity, ip string, now time.Time) (bool, time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	policy, exists := r.policies[action]
	if !exists {
		return true, 0
	}

	key := buildAuthRateLimitKey(action, identity, ip)

	state, exists := r.attempts[key]
	if !exists {
		state = &authAttemptState{}
		r.attempts[key] = state
	}

	if state.LastRequest.IsZero() || now.Sub(state.LastRequest) > policy.ResetAfter {
		state.Attempts = 0
		state.BlockedUntil = time.Time{}
	}

	if now.Before(state.BlockedUntil) {
		return false, state.BlockedUntil.Sub(now)
	}

	state.Attempts++
	state.LastRequest = now

	exponent := state.Attempts - 1
	if exponent > 10 {
		exponent = 10
	}

	delay := policy.BaseDelay * time.Duration(1<<exponent)
	if delay > policy.MaxDelay {
		delay = policy.MaxDelay
	}
	state.BlockedUntil = now.Add(delay)

	return true, 0
}

func buildAuthRateLimitKey(action AuthRateLimitAction, identity, ip string) string {
	normalizedIdentity := strings.ToLower(strings.TrimSpace(identity))
	if normalizedIdentity == "" {
		normalizedIdentity = "anonymous"
	}
	normalizedIP := strings.TrimSpace(ip)
	return string(action) + "|" + normalizedIdentity + "|" + normalizedIP
}
