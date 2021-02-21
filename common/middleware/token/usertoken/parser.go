/**
 * @Author fzh
 * @Date 2021/2/8
 */
package usertoken

import (
	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserId   uint
	UserName string
}

func TokenParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从header中获取token
		tokenList, exist := c.Request.Header["Token"]
		if !exist || len(tokenList) == 0 {
			return
		}
		token := tokenList[0]

		// 解析token
		claims, err := Parse(token)
		if err != nil {
			return
		}
		info := &UserInfo{
			UserId:   claims.Userid,
			UserName: claims.Username,
		}
		c.Set("UserInfo", info)
	}
}

func GetUserInfo(c *gin.Context) (info *UserInfo, exist bool) {
	infoInf, exist := c.Get("UserInfo")
	if exist {
		info = infoInf.(*UserInfo)
	} else {
		info = nil
	}
	return
}
