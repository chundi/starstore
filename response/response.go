package response

type Response struct {
	Code    int         `json:"code"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
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

type DataListObject struct {
	Count int         `json:"count"`
	List  interface{} `json:"list"`
}
