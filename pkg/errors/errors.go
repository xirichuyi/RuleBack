// Package errors 统一错误处理
package errors

import (
	"errors"
	"fmt"
)

// 通用错误码 (10xxx)
const (
	CodeSuccess         = 0
	CodeUnknown         = 10000
	CodeInvalidParams   = 10001
	CodeUnauthorized    = 10002
	CodeForbidden       = 10003
	CodeNotFound        = 10004
	CodeConflict        = 10005
	CodeInternalError   = 10006
	CodeDatabaseError   = 10007
	CodeValidationError = 10008
	CodeRateLimited     = 10009
	CodeTimeout         = 10010
)

// 用户模块错误码 (20xxx)
const (
	CodeUserNotFound      = 20001
	CodeUserExists        = 20002
	CodeUserDisabled      = 20003
	CodePasswordIncorrect = 20004
	CodeTokenExpired      = 20005
	CodeTokenInvalid      = 20006
)

// codeMessages 错误码对应的默认消息
var codeMessages = map[int]string{
	CodeSuccess:           "成功",
	CodeUnknown:           "未知错误",
	CodeInvalidParams:     "参数错误",
	CodeUnauthorized:      "未认证",
	CodeForbidden:         "无权限",
	CodeNotFound:          "资源不存在",
	CodeConflict:          "资源冲突",
	CodeInternalError:     "内部错误",
	CodeDatabaseError:     "数据库错误",
	CodeValidationError:   "验证错误",
	CodeRateLimited:       "请求过于频繁",
	CodeTimeout:           "请求超时",
	CodeUserNotFound:      "用户不存在",
	CodeUserExists:        "用户已存在",
	CodeUserDisabled:      "用户已禁用",
	CodePasswordIncorrect: "密码错误",
	CodeTokenExpired:      "Token已过期",
	CodeTokenInvalid:      "Token无效",
}

// GetMessage 获取错误码对应的默认消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// AppError 应用错误类型
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewWithCode 根据错误码创建错误（使用默认消息）
func NewWithCode(code int) *AppError {
	return &AppError{
		Code:    code,
		Message: GetMessage(code),
	}
}

// Wrap 包装原始错误
func Wrap(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapWithCode 使用默认消息包装原始错误
func WrapWithCode(code int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: GetMessage(code),
		Err:     err,
	}
}

// IsAppError 判断是否为应用错误
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError 获取应用错误
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// GetCode 获取错误码
func GetCode(err error) int {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.Code
	}
	return CodeInternalError
}

// 预定义错误实例
var (
	ErrInvalidParams     = NewWithCode(CodeInvalidParams)
	ErrUnauthorized      = NewWithCode(CodeUnauthorized)
	ErrForbidden         = NewWithCode(CodeForbidden)
	ErrNotFound          = NewWithCode(CodeNotFound)
	ErrConflict          = NewWithCode(CodeConflict)
	ErrInternalError     = NewWithCode(CodeInternalError)
	ErrDatabaseError     = NewWithCode(CodeDatabaseError)
	ErrUserNotFound      = NewWithCode(CodeUserNotFound)
	ErrUserExists        = NewWithCode(CodeUserExists)
	ErrPasswordIncorrect = NewWithCode(CodePasswordIncorrect)
	ErrTokenExpired      = NewWithCode(CodeTokenExpired)
	ErrTokenInvalid      = NewWithCode(CodeTokenInvalid)
)
