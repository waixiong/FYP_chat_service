package token

import (
	"time"
)

type Status int

const (
	Valid    Status = 1
	Refresh  Status = 2
	NotValid Status = 3
)

func ValidateAccessToken(accessToken string) (string, Status) {
	token, err := RawTokenFromAccessToken(accessToken)
	if err != nil {
		return "", NotValid
	}
	if token.Exp <= time.Now().Unix() {
		return "", Refresh
	}
	return token.User, Valid
}

func ValidateRefreshToken(accessToken string, refreshToken string) (string, Status) {
	token, err := RawTokenFromAccessToken(accessToken)
	if err != nil {
		return "", NotValid
	}
	if token.Refresh != refreshToken[8:] {
		return "", NotValid
	}
	if token.Exp > time.Now().Unix() {
		return token.User, Valid // no refresh
	}
	if token.Exp < time.Now().Add(-71*time.Hour).Unix() {
		return token.User, NotValid // expired
	}
	return token.User, Refresh
}
