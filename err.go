package std

import (
	"errors"
	"fmt"
	"strings"
)

type Err struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误信息
	Data interface{} `json:"data,omitempty"`
}

func (this *Err) Error() string {
	return this.Msg
}

func NewError(code int, msg string, data interface{}) *Err {
	return &Err{
		code,
		msg,
		data,
	}
}

//错误码 = -1
func ErrorMsg(msg string) *Err {
	return NewError(-1, msg, nil)
}

func ErrorWithErr(msg string, err error) *Err {
	return ErrorMsg(fmt.Sprintf("%s :causedBy %s", msg, err.Error()))
}

func NewSuccess() *Err {
	return NewSuccessWithData(nil)
}

func NewSuccessWithData(data interface{}) *Err {
	return NewError(0, "success", data)
}

func AssertError(err error, msg string) {
	if err == nil {
		return
	}
	output := fmt.Sprintf("%s :causedBy %s", msg, err.Error())
	//defer func() {
	//	Assert(false, output)
	//}()
	panic(output)
}

func Assert(cond bool, msg string) {
	if cond {
		return
	}
	//defer func() {
	//	const abortCode = 6
	//	os.Exit(abortCode)
	//}()
	e := errors.New(fmt.Sprintf("assertFailed!,%s", msg))
	panic(e)

}

type CombinedErrors []error

func (this CombinedErrors) Error() string {
	builder := strings.Builder{}

	for idx, it := range this {
		builder.WriteString(fmt.Sprintf("err %d:%s", idx+1, it.Error()))
	}
	return builder.String()
}
