/**
 * @Author fzh
 * @Date 2020/2/1
 */
package service

import (
	"common/middleware/token/usertoken"
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
	"user/model"
)

// 用户注册
func Register(username, password string) error {
	salt := generateSalt()
	hashedPassword := hash2(password, salt)

	u := &model.User{
		Model:             gorm.Model{},
		CreatedBy:         "",
		UpdatedBy:         "",
		DeletedBy:         "",
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
	// TODO: 具体错误类型判断
	if err := model.InsertUser(DB, u); err != nil {
		return errortype.ErrUserNameHasExist
	}
	return nil
}

// 用户登录 返回token
func Login(username, password string) (string, error) {
	// 根据用户名获取用户信息
	u, err := model.GetUserByUsername(DB, username)
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
	token, _ := usertoken.Generate(u.ID, u.Username)
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

func GetUser(id uint) (*model.User, error) {
	u, err := model.GetUserById(DB, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Debug("[用户信息查询] 用户ID不存在:", id)
			return nil, errortype.ErrUserNotExist
		} else {
			logging.Error(err)
			return nil, errortype.ErrUnknown
		}
	}
	return u, nil
}
