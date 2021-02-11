/**
 * @Author fzh
 * @Date 2021/2/8
 */
package user

import (
	"common/tools/logging"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claim struct {
	Userid   uint   `json:"userid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var (
	signingKey = []byte("SIGNING123KEY")

	ErrParseFailure = errors.New("Token解析失败")
	ErrInvalidToken = errors.New("Token无效")
)

func Generate(userid uint, username string) (string, error) {
	c := &Claim{
		Userid:   userid,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,
			ExpiresAt: time.Now().Unix() + 30*60,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		logging.Error("生成Token失败", err)
		return "", err
	}
	return ss, nil
}

func Parse(tokenString string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		logging.Error("Token解析失败", err)
		return nil, ErrParseFailure
	}
	if claims, ok := token.Claims.(*Claim); ok && token.Valid {
		return claims, nil
	} else {
		return nil, ErrInvalidToken
	}
}
