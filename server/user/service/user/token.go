/**
 * @Author fzh
 * @Date 2021/2/1
 */
package user

import (
	"common/tools/logging"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func generateToken(username, password string) (string, error) {
	signingKey := []byte("SIGNING123KEY")

	c := &Claim{
		Username: username,
		Password: password,
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
