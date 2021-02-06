/**
 * @Author fzh
 * @Date 2020/2/1
 */
package user

import (
	"common/tools/logging"
	"crypto/md5"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"strconv"
	"time"
	. "user/global/database"
	"user/global/errortype"
	"user/model/user"
)

// 用户注册
func Register(username, password string) error {
	salt := generateSalt()
	hashedPassword := hash2(password, salt)

	u := &user.User{
		Model:             gorm.Model{},
		Username:          username,
		Password:          hashedPassword,
		State:             0,
		Salt:              salt,
		CertificateType:   0,
		Name:              "",
		CertificateNumber: "",
		PhoneNumber:       "",
		Email:             "",
		PassengerType:     0,
	}

	logging.Debug("[用户注册] 用户名:", username)
	if err := user.InsertUser(DB, u); err != nil {
		return errortype.ErrUserNameHasExist
	}
	return nil
}

// 用户登录 返回token
func Login(username, password string) (string, error) {
	// 根据用户名获取用户信息
	u, err := user.GetUserByUsername(DB, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Debug("[用户登录] 用户名不存在:", username)
			return "", errortype.ErrUserNotExist
		} else {
			logging.Error(err)
			return "", errortype.ErrUnknown
		}
	}

	// 验证密码
	hashedPassword := hash2(password, u.Salt)
	if hashedPassword != u.Password {
		return "", errortype.ErrWrongPassword
	}

	logging.Debug(username, "登录成功")
	token, _ := generateToken(username, password)
	return token, nil
}

func generateSalt() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// 拼接两个字符串，计算Hash
func hash2(s1 string, s2 string) string {
	h := md5.New()
	if _, err := io.WriteString(h, s1); err != nil {
		logging.Error(err)
	}
	if _, err := io.WriteString(h, s2); err != nil {
		logging.Error(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}