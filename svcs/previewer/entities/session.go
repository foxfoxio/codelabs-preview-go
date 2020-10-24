package entities

import (
	"encoding/json"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"time"
)

const (
	SessionKeyUserSession = "user-session"
)

type UserSession struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewSession(userId string) *UserSession {
	return &UserSession{
		Id:        utils.NewID(),
		UserId:    userId,
		CreatedAt: time.Now(),
	}
}

func (s *UserSession) Empty() bool {
	return s == nil || s.Id == ""
}

func (s *UserSession) IsValid() bool {
	return s != nil && s.UserId != ""
}

func (s *UserSession) String() string {
	if s == nil {
		return ""
	}

	return utils.Stringify(*s)
}

func (s UserSession) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func UnmarshalSession(value string) *UserSession {
	session := &UserSession{}
	_ = json.Unmarshal([]byte(value), session)

	return session
}
