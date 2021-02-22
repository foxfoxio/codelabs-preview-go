package token

import "encoding/base64"

func Swizzle(token string) string {
	buf := []byte(token)

	for i, b := range buf {
		buf[i] = byte(int(b) - 2*(i%2) + 1)
	}

	b64Str := base64.StdEncoding.EncodeToString(buf)
	b64 := []byte(b64Str)

	return string(reverse(b64))
}

func Unswizzle(token string) string {
	revToken := reverse([]byte(token))
	buf, err := base64.StdEncoding.DecodeString(string(revToken))

	if err != nil {
		return ""
	}

	for i, b := range buf {
		buf[i] = byte(int(b) + 2*(i%2) - 1)
	}

	return string(buf)
}

func reverse(bytes []byte) []byte {
	for i := 0; i < len(bytes)/2; i++ {
		j := len(bytes) - i - 1
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

/*
export const swizzle = (idToken: string) => {
  const shifted_idToken = idToken
    .split('')
    .map((e, i) => String.fromCharCode(e.charCodeAt(0) - 2 * (i % 2) + 1))
    .join('')
  const idToken_base64 = encodeBase64(decodeUTF8(shifted_idToken))
  const reversed_idToken_base64 = idToken_base64.split('').reverse().join('')
  return reversed_idToken_base64
}
*/
