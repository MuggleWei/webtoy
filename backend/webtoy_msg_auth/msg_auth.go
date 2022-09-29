package webtoymsgauth

type MsgAuthUserReq struct {
	Name           string `json:"name,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Passwd         string `json:"passwd"`
	CaptchaSession string `json:"captcha_session,omitempty"`
	CaptchaValue   string `json:"captcha_value,omitempty"`
}

type MsgAuthUserRsp struct {
	UserID  string `json:"uid"`
	Session string `json:"session"`
	Token   string `json:"token"`
}

type MsgAddUserReq struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Passwd string `json:"passwd"`
}

type MsgAddUserRsp struct {
	UserID string `json:"uid"`
}

type MsgQueryUserReq struct {
	UserID string `json:"uid,omitempty"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`
}

type MsgQueryUserRsp struct {
	UserID   string `json:"uid"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	ShowName string `json:"show_name,omitempty"`
}
