package model

type ModelUserLoginReq struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Passwd         string `json:"passwd"`
	CaptchaSession string `json:"captcha_session"`
	CaptchaValue   string `json:"captcha_value"`
}

type ModelUserLoginRsp struct {
	Id      int64  `json:"uid"`
	Session string `json:"session"`
	Token   string `json:"token"`
}
