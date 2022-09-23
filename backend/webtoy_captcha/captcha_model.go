package main

type ModelCaptchaVerify struct {
	CaptchaSessionID string `json:"k"`
	CaptchaValue     string `json:"v"`
}
