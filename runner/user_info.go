package runner

const (
	LoginUserHeaderKey = "X-Login-User"
)

type LoginUserInfo struct {
	IsLoggedIn bool   //登陆用户的请求
	Username   string //登陆用户名称
}

// GetLoginUserInfo 获取请求用户信息
func (c *Context) GetLoginUserInfo() LoginUserInfo {
	user, ok := c.Request.Headers[LoginUserHeaderKey]
	if !ok {
		return LoginUserInfo{IsLoggedIn: false, Username: ""}
	}
	return LoginUserInfo{IsLoggedIn: true, Username: user}
}
