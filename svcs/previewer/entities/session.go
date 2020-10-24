package entities

import (
	"encoding/json"
	tokenUtil "github.com/foxfoxio/codelabs-preview-go/internal/token"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

const (
	SessionKeyUserSession = "user-session"
)

type UserSession struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	UserId      string    `json:"userId"`
	State       string    `json:"state"`
	Token       string    `json:"token"`
	CreatedAt   time.Time `json:"createdAt"`
	RedirectUrl string    `json:"redirectUrl"`
	session     *sessions.Session
}

func NewUserSession(session *sessions.Session) *UserSession {
	userSession := &UserSession{
		Id:        utils.NewID(),
		CreatedAt: time.Now(),
		session:   session,
	}

	if session == nil {
		return userSession
	}

	if s, ok := session.Values[SessionKeyUserSession].(string); ok && s != "" {
		userSession.Unmarshal(s)
	}

	return userSession
}

func (s *UserSession) Empty() bool {
	return s == nil || s.Id == ""
}

func (s *UserSession) Oauth2Token() *oauth2.Token {
	if s.Token == "" {
		return nil
	}
	if t, err := tokenUtil.DecodeBase64(s.Token); err == nil {
		return t
	}

	return nil
}

func (s *UserSession) IsValid() bool {
	if s == nil {
		return false
	}

	t := s.Oauth2Token()
	return t != nil && t.Valid()
}

func (s *UserSession) String() string {
	if s == nil {
		return ""
	}

	return utils.Stringify(*s)
}

func (s *UserSession) Invalidate(r *http.Request, w http.ResponseWriter) error {
	s.session.Values[SessionKeyUserSession] = ""
	s.session.Options.MaxAge = -1 // immediately expires the cookies
	return s.Save(r, w)
}

func (s *UserSession) Save(r *http.Request, w http.ResponseWriter) error {
	s.session.Values[SessionKeyUserSession] = s.Marshal()
	return s.session.Save(r, w)
}

func (s UserSession) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *UserSession) Unmarshal(value string) {
	_ = json.Unmarshal([]byte(value), s)
}

func UnmarshalSession(value string) *UserSession {
	session := &UserSession{}
	_ = json.Unmarshal([]byte(value), session)

	return session
}
