package response

type Response struct {
	Code int 	`json:"code"`
	Message string	`json:"message"`
	Data interface{}	`json:"data"`
}

func (res *Response) SetCode(code int) {
	res.Code = code
}

func (res *Response) SetMessage(msg string) {
	res.Message = msg
}

func (res *Response) SetData(data interface{}) {
	res.Data = data
}

