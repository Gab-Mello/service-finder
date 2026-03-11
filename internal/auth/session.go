package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type SessionManager struct {
	mu     sync.RWMutex
	ttl    time.Duration
	byID   map[string]session
	secure bool
	done   chan struct{}
}

type session struct {
	UserID string
	ExpAt  time.Time
}

func NewSessionManager(ttl time.Duration) *SessionManager {
	m := &SessionManager{
		ttl:    ttl,
		byID:   make(map[string]session),
		secure: false,
		done:   make(chan struct{}),
	}
	go m.cleanupLoop()
	return m
}

func (m *SessionManager) SetSecure(secure bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.secure = secure
}

func (m *SessionManager) Close() {
	close(m.done)
}

func (m *SessionManager) cleanupLoop() {
	ticker := time.NewTicker(m.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.done:
			return
		}
	}
}

func (m *SessionManager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for sid, s := range m.byID {
		if now.After(s.ExpAt) {
			delete(m.byID, sid)
		}
	}
}

func (m *SessionManager) New(userID string) (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	sid := hex.EncodeToString(b[:])

	m.mu.Lock()
	defer m.mu.Unlock()
	m.byID[sid] = session{UserID: userID, ExpAt: time.Now().Add(m.ttl)}
	return sid, nil
}

func (m *SessionManager) Get(sid string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.byID[sid]
	if !ok || time.Now().After(s.ExpAt) {
		return "", false
	}
	return s.UserID, true
}

func (m *SessionManager) Delete(sid string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.byID, sid)
}

func (m *SessionManager) SetCookie(w http.ResponseWriter, sid string) {
	m.mu.RLock()
	secure := m.secure
	m.mu.RUnlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(m.ttl.Seconds()),
	})
}

func (m *SessionManager) ClearCookie(w http.ResponseWriter) {
	m.mu.RLock()
	secure := m.secure
	m.mu.RUnlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		MaxAge:   -1,
	})
}
