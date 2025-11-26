package util

// 统一响应结构体（成功/失败都用此格式）
type Response struct {
	Code    int         `json:"code"`           // 业务错误码（0 表示成功）
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 成功时返回的数据（可选）
}

// 自定义错误类型（包含业务错误码和 HTTP 状态码）
type BusinessError struct {
	HttpCode int    // HTTP 响应状态码（如 401、404）
	Code     int    // 业务错误码（如 1001、1002，便于前端区分错误类型）
	Message  string // 错误提示信息
}

// 实现 error 接口（必须，让 BusinessError 满足 error 类型）
func (e *BusinessError) Error() string {
	return e.Message
}

// 快速创建业务错误（工具函数，简化错误创建）
func NewBusinessError(httpCode, code int, message string) *BusinessError {
	return &BusinessError{
		HttpCode: httpCode,
		Code:     code,
		Message:  message,
	}
}

// 预定义常见业务错误（避免硬编码，统一维护）
var (
	// 认证相关
	ErrAuthFailed   = NewBusinessError(401, 1001, "用户认证失败，请重新登录")
	ErrTokenInvalid = NewBusinessError(401, 1002, "Token 无效或已过期")
	// 资源相关
	ErrResourceNotFound = NewBusinessError(404, 2001, "请求的资源不存在")
	ErrArticleNotFound  = NewBusinessError(404, 2002, "文章不存在")
	ErrCommentNotFound  = NewBusinessError(404, 2003, "评论不存在")
	// 参数相关
	ErrInvalidParam = NewBusinessError(400, 3001, "请求参数无效")
	// 系统相关
	ErrDBConnect   = NewBusinessError(500, 5001, "数据库连接失败")
	ErrSystemError = NewBusinessError(500, 5002, "服务器内部错误，请稍后重试")
)

// 成功响应（工具函数）
func Success(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}
