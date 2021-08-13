package message

type MSG map[string]interface{}

func (m MSG) IsSuccess() bool {
	return m["status"] == "success"
}
