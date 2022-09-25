package webtoy_base

const (
	ERROR_OK       = 0
	ERROR_COMMON   = -1  // 通用错误
	ERROR_INTERNAL = 501 // 内部错误
	ERROR_AUTH     = 502 // 认证错误(无法获取session, token过期等)
	ERROR_BAD_REQ  = 503 // 客户端错误的请求
	ERROR_CAPTCHA  = 504 // 验证码错误
)
