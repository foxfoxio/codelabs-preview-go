package token

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"golang.org/x/oauth2"
	"strings"
	"time"
)

type JwtClaims struct {
	Iss           string `json:"iss,omitempty"`
	Azp           string `json:"azp,omitempty"`
	Aud           string `json:"aud,omitempty"`
	Sub           string `json:"sub,omitempty"`
	Hd            string `json:"hd,omitempty"`
	UserId        string `json:"user_id,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	AtHash        string `json:"at_hash,omitempty"`
	Nonce         string `json:"nonce,omitempty"`
	Iat           int    `json:"iat,omitempty"`
	Exp           int    `json:"exp,omitempty"`
	Name          string `json:"name,omitempty"`
	Picture       string `json:"picture,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	Locale        string `json:"locale,omitempty"`
}

func (j *JwtClaims) IssuedAt() time.Time {
	return time.Unix(int64(j.Iat), 0)
}

func (j *JwtClaims) ExpiresAt() time.Time {
	return time.Unix(int64(j.Exp), 0)
}

func (j *JwtClaims) Valid() bool {
	return j.Email != "" && j.UserId != "" && time.Now().Before(j.ExpiresAt())
}

func ExtractJwtClaims(token string) (*JwtClaims, error) {
	tokenStruct := &JwtClaims{}
	jwtParts := strings.Split(token, ".")
	out, _ := base64.RawURLEncoding.DecodeString(jwtParts[1])
	err := json.Unmarshal(out, &tokenStruct)
	if err != nil {
		return nil, err
	}

	return tokenStruct, nil
}

func EncodeBase64(token *oauth2.Token) (string, error) {
	var buffer bytes.Buffer
	if e := gob.NewEncoder(&buffer).Encode(token); e != nil {
		return "", e
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func DecodeBase64(encodedToken string) (*oauth2.Token, error) {
	var token oauth2.Token
	var tokenBytes []byte
	if b, e := base64.StdEncoding.DecodeString(encodedToken); e != nil {
		return nil, e
	} else {
		tokenBytes = b
	}
	reader := bytes.NewReader(tokenBytes)
	if ee := gob.NewDecoder(reader).Decode(&token); ee != nil {
		return nil, ee
	}

	return &token, nil
}
