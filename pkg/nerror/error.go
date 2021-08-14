package nerror

import "errors"

var (
	SendSMSError      = errors.New("sms send error")
	TypeNotFoundError = errors.New("type not found")
	UnknownError      = errors.New("unknown error")
)
