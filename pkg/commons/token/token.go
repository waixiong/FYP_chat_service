package token

import (
	"strconv"
	"strings"
)

type Get_Token struct {
	User    string `json:"u"`
	Exp     int64  `json:"e"`
	Refresh string `json:"r"`
}

func (t *Get_Token) toString() string {
	return t.User + "." + strconv.FormatInt(t.Exp, 36) + "." + t.Refresh
}

func get_token_FromString(token string) (*Get_Token, error) {
	dot1 := strings.Index(token, ".")
	dot2 := strings.Index(token[dot1+1:], ".") + dot1 + 1
	n, err := strconv.ParseInt(token[dot1+1:dot2], 36, 64)
	return &Get_Token{
		User:    token[:dot1],
		Exp:     n,
		Refresh: token[dot2+1:],
	}, err
}
