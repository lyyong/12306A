/**
 * @Author fzh
 * @Date 2020/2/1
 */
package errortype

import "errors"

var (
	ErrUserNotExist     = errors.New("用户不存在")
	ErrWrongPassword    = errors.New("密码不正确")
	ErrUserNameHasExist = errors.New("用户名已存在")

	ErrUnknown = errors.New("未知错误")
)
