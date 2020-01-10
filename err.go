package std

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Error interface {
	SetCode(code int) Error
	SetMsg(msg string) Error
	SetData(data interface{}) Error
	GetData() interface{}
	GetMsg() string
	GetCode() int
	Is(err error) bool
}

type _stdError struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误信息
	Data interface{} `json:"data,omitempty"`
}

func (this *_stdError) GetCode() int {
	return this.Code
}

func (this *_stdError) GetMsg() string {
	return this.Msg
}

func (this *_stdError) GetData() interface{} {
	return this.Data
}

func (this *_stdError) Error() string {
	return this.Msg
}

func (this *_stdError) Is(err error) bool {
	if err == nil {
		return false
	}
	if xe, ok := err.(*_stdError); ok {
		return xe.Code == this.Code
	}
	return false
}

func NewError(msg string) Error {
	return &_stdError{
		0, msg, nil,
	}
}

func Errorf(format string, args ...interface{}) Error {
	out := NewError(fmt.Sprintf(format, args...))
	return out
}

func ErrorWrap(err error, msg string) Error {
	out := NewError(msg)
	out.SetMsg(fmt.Sprintf("%s causedBy:%v", out.GetMsg(), err))
	if err != nil {
		if xe, ok := err.(*_stdError); ok {
			out.SetCode(xe.Code)
			out.SetData(xe.Data)
		}
	}
	return out
}

func ErrorWrapf(err error, format string, args ...interface{}) Error {
	return ErrorWrap(err, fmt.Sprintf(format, args...))
}

func (this *_stdError) SetMsg(msg string) Error {
	this.Msg = msg
	return this
}

func (this *_stdError) SetCode(code int) Error {
	this.Code = code
	return this
}

func (this *_stdError) SetData(data interface{}) Error {
	this.Data = data
	return this
}

func IsStdError(e error) (err Error, ok bool) {
	if e == nil {
		return nil, false
	}
	err, ok = e.(*_stdError)
	return
}

func AssertError(err error, msg string) {
	if err == nil {
		return
	}
	//defer func() {
	//	Assert(false, output)
	//}()
	panic(errors.Wrapf(err, "%s :causedBy", msg))
}

func Assert(cond bool, msg string) {
	if cond {
		return
	}
	//defer func() {
	//	const abortCode = 6
	//	os.Exit(abortCode)
	//}()
	e := errors.Errorf("assertFailed!,%s", msg)
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
