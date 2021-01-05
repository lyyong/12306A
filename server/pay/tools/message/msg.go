package message

var msg = map[int]string{
	ERROR: "错误",
	PAYOK: "支付完成",
}

func GetMsg(code int) string {
	m, ok := msg[code]
	if ok {
		return m
	}
	return msg[ERROR]
}
