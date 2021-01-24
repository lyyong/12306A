// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package message

var msg = map[int]string{
	ERROR:        "错误",
	PAYOK:        "支付完成",
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
