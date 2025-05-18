package runner

//
//// Context 上下文结构体，封装请求上下文信息
//type Context struct {
//	context.Context
//	userInfo     *UserInfo
//	sessionData  map[string]interface{}
//	sessionMutex sync.RWMutex
//	startTime    time.Time
//	requestID    string
//}
//
//func (c *Context) getTraceId() string {
//	value := c.Context.Value(constants.TraceID)
//	if value == nil {
//		return ""
//	}
//	v, ok := value.(string)
//	if ok {
//		return v
//	}
//	return ""
//}
//
//// UserInfo 用户信息结构
//type UserInfo struct {
//	Username string
//	UserID   string
//	Email    string
//	Roles    []string
//}
//
//// NewContext 创建新的上下文
//func NewContext(ctx context.Context) *Context {
//	return &Context{
//		Context:     ctx,
//		sessionData: make(map[string]interface{}),
//		startTime:   time.Now(),
//	}
//}
//
//// WithRequestID 设置请求ID
//func (c *Context) WithRequestID(requestID string) *Context {
//	c.requestID = requestID
//	return c
//}
//
//// RequestID 获取请求ID
//func (c *Context) RequestID() string {
//	return c.requestID
//}
//
//// ElapsedTime 获取自上下文创建以来的经过时间
//func (c *Context) ElapsedTime() time.Duration {
//	return time.Since(c.startTime)
//}
//
//// GetUsername 获取当前用户名
//func (c *Context) GetUsername() string {
//	if c.userInfo == nil {
//		return ""
//	}
//	return c.userInfo.Username
//}
//
//// GetUserID 获取当前用户ID
//func (c *Context) GetUserID() string {
//	if c.userInfo == nil {
//		return ""
//	}
//	return c.userInfo.UserID
//}
//
//// SetUserInfo 设置用户信息
//func (c *Context) SetUserInfo(info *UserInfo) {
//	c.userInfo = info
//	logrus.Debugf("用户信息已设置: %s", info.Username)
//}
//
//// HasRole 检查用户是否拥有指定角色
//func (c *Context) HasRole(role string) bool {
//	if c.userInfo == nil || len(c.userInfo.Roles) == 0 {
//		return false
//	}
//
//	for _, r := range c.userInfo.Roles {
//		if r == role {
//			return true
//		}
//	}
//
//	return false
//}
//
//// Set 存储会话数据
//func (c *Context) Set(key string, value interface{}) {
//	c.sessionMutex.Lock()
//	defer c.sessionMutex.Unlock()
//	c.sessionData[key] = value
//}
//
//// Get 获取会话数据
//func (c *Context) Get(key string) (interface{}, bool) {
//	c.sessionMutex.RLock()
//	defer c.sessionMutex.RUnlock()
//	value, exists := c.sessionData[key]
//	return value, exists
//}
//
//// GetString 获取字符串类型的会话数据
//func (c *Context) GetString(key string) (string, bool) {
//	value, exists := c.Get(key)
//	if !exists {
//		return "", false
//	}
//
//	str, ok := value.(string)
//	return str, ok
//}
//
//// Clear 清除会话数据
//func (c *Context) Clear() {
//	c.sessionMutex.Lock()
//	defer c.sessionMutex.Unlock()
//	c.sessionData = make(map[string]interface{})
//}
