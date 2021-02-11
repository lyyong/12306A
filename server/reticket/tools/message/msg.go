// @Author LiuYong
// @Created at 2021-02-02
// @Modified at 2021-02-02
package message

var msg = map[int]string{
	ERROR:        "错误",
	PARAMS_ERROR: "参数错误",
	OK:           "成功",
}

func GetMsg(code int) string {
	m, ok := msg[code]
	if ok {
		return m
	}
	return msg[ERROR]
}
