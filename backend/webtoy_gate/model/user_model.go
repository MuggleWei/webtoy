package model

type ModelUserLoginReq struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Passwd string `json:"passwd"`
}

type ModelUserLoginRsp struct {
	Id      int64  `json:"user_id"`
	Session string `json:"session"`
	Token   string `json:"token"`
}
