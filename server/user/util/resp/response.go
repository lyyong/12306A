/**
 * @Author fzh
 * @Date 2020/2/1
 */
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

func (r *Response) SetCode(code string) *Response {
	r.Code = code
	return r
}

func (r *Response) SetMsg(msg string) *Response {
	r.Msg = msg
	return r
}
