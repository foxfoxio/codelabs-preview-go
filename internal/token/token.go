package token

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"golang.org/x/oauth2"
	"strings"
)

type JwtClaims struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Hd            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Nonce         string `json:"nonce"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
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
