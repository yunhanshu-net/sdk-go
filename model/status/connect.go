package status

import "fmt"

type ErrorStatus struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	ErrInfo string `json:"err_info"`
}

func (s *ErrorStatus) WithError(err error) *ErrorStatus {
	s.ErrInfo = err.Error()
	return s
}
func (s *ErrorStatus) Error() string {
	return fmt.Sprintf("message:%s,status:%d,err_info:%s", s.Message, s.Status, s.ErrInfo)
}

var (
	ConnectError = &ErrorStatus{Message: "建立连接失败", Status: -1}
)
