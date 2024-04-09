package result

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func SuccessWithData(data interface{}) *Result {
	return &Result{
		Code: 200,
		Data: data,
	}
}

func Success() *Result {
	return &Result{
		Code: 200,
	}
}

func CreateSuccess() *Result {
	return &Result{
		Code: 201,
	}
}

func DeleteSuccess() *Result {
	return &Result{
		Code: 204,
	}
}

func Fail(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
	}
}
