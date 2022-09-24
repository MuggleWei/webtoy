package webtoymsgcaptcha

type MsgCaptchaVerifyReq struct {
	CaptchaSessionID string `json:"captcha_session"`
	CaptchaValue     string `json:"captcha_value"`
}
