package token

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractJwtClaims(t *testing.T) {
	token := "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjJmOGI1NTdjMWNkMWUxZWM2ODBjZTkyYWFmY2U0NTIxMWUxZTRiNDEiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vZm94Zm94LWxlYXJuIiwiYXVkIjoiZm94Zm94LWxlYXJuIiwiYXV0aF90aW1lIjoxNjA1MDMyMzI5LCJ1c2VyX2lkIjoiMU9TM1FBZ1lRc2d0bjNwTHpIOXJLQ0dqY1A5MyIsInN1YiI6IjFPUzNRQWdZUXNndG4zcEx6SDlyS0NHamNQOTMiLCJpYXQiOjE2MDUwMzIzMjksImV4cCI6MTYwNTAzNTkyOSwiZW1haWwiOiJzc2RmQGFzZGZhZGYuc2RmZC5jb20iLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZW1haWwiOlsic3NkZkBhc2RmYWRmLnNkZmQuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoicGFzc3dvcmQifX0.c7SZROTypgKcrEQEozAMQvJxkcD4ZnGD6XSuMYN4X6UxwSFhCaC7Gs1ns-ur1XMokM-7eO2BQK9-1R0VV5enOfIrzJV1ltqNMbj3sqRvKyC93dHTJ8lVEbSHzcI4Zo08R40I8tmqJWJDoWLdHUA1hg_rqOFIy8nvDcRvz6vzy66JJi_JdQyVFDECv08P5UNe7LhCn4cOELdoaR2uSUtilJDrX4ykEZHISQWotcY3adHlvCmA7KBXOFLnRRcSLHpo2a3aUS3Y5o7nk_BITyby2shqXdSWyKe2iyiUCArzlvlBrg8coryFRyV6qXHjtyOEAnpDqMfbeIwHN9kudgAyVw"

	claim, err := ExtractJwtClaims(token)

	assert.NoError(t, err)

	fmt.Println(utils.StringifyIndent(claim))
	fmt.Println(claim.IssuedAt())
	fmt.Println(claim.ExpiresAt())
}

/*
{
  "iss": "https://securetoken.google.com/foxfox-learn",
  "aud": "foxfox-learn",
  "auth_time": 1605032329,
  "user_id": "1OS3QAgYQsgtn3pLzH9rKCGjcP93",
  "sub": "1OS3QAgYQsgtn3pLzH9rKCGjcP93",
  "iat": 1605032329,
  "exp": 1605035929,
  "email": "ssdf@asdfadf.sdfd.com",
  "email_verified": false,
  "firebase": {
    "identities": {
      "email": [
        "ssdf@asdfadf.sdfd.com"
      ]
    },
    "sign_in_provider": "password"
  }
}
*/
