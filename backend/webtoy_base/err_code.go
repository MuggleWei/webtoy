package webtoy_base

const (
	ERROR_OK       = 0
	ERROR_INTERNAL = 501 // 内部错误
	ERROR_AUTH     = 502 // 认证错误(无法获取session, token过期等)
)
