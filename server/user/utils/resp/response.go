package resp

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func R(data interface{}) *Response {
	r := &Response{
		Data: data,
	}
	return r
}
