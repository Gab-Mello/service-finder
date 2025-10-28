package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

type SessionManager struct {
	ttl  time.Duration
	byID map[string]session
}
type session struct {
	UserID string
	ExpAt  time.Time
}

func NewSessionManager(ttl time.Duration) *SessionManager {
	return &SessionManager{
		ttl:  ttl,
		byID: make(map[string]session),
	}
}

func (m *SessionManager) New(userID string) string {
	var b [32]byte
	_, _ = rand.Read(b[:])
	sid := hex.EncodeToString(b[:])
	m.byID[sid] = session{UserID: userID, ExpAt: time.Now().Add(m.ttl)}
	return sid
}

func (m *SessionManager) Get(sid string) (string, bool) {
	s, ok := m.byID[sid]
	if !ok || time.Now().After(s.ExpAt) {
		return "", false
	}
	return s.UserID, true
}

func (m *SessionManager) Delete(sid string) {
	delete(m.byID, sid)
}

func (m *SessionManager) SetCookie(w http.ResponseWriter, sid string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,

		MaxAge: int(m.ttl.Seconds()),
	})
}

func (m *SessionManager) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
