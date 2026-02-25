package token

import "encoding/base64"

func generateToken() string {
	return base64.StdEncoding.EncodeToString(make([]byte, 32))
}