package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/MuggleWei/webtoy/backend/webtoy_auth/service"
	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	msgAuth "github.com/MuggleWei/webtoy/backend/webtoy_msg_auth"
	log "github.com/sirupsen/logrus"
)

type AuthController struct {
	authService *service.AuthService
}

var (
	singletonUser *AuthController
	onceUser      sync.Once
)

func GetAuthController() *AuthController {
	if singletonUser == nil {
		onceUser.Do(func() {
			singletonUser = &AuthController{
				authService: service.GetAuthService(),
			}
		})
	}
	return singletonUser
}

// user check
func (this *AuthController) UserAuth(w http.ResponseWriter, r *http.Request) {
	var req msgAuth.MsgAuthUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("%v", err.Error())
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}

	log.Debugf("user auth, name=%v, email=%v, phone=%v, passwd=******",
		req.Name, req.Email, req.Phone)

	rsp, err := this.authService.UserAuth(&req)
	if err != nil {
		log.Debugf("user auth failed, %v", err.Error())
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_AUTH,
			ErrMsg: err.Error(),
		})
	} else {
		log.Debugf("user auth success")
		base.HttpResponse(w, &base.MessageRsp{
			Data: rsp,
		})
	}
}

// user query
func (this *AuthController) UserQuery(w http.ResponseWriter, r *http.Request) {
	var req msgAuth.MsgQueryUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("%v", err.Error())
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}

	log.Debugf("user query, req=%+v", req)

	rsp, err := this.authService.UserQuery(&req)
	if err != nil {
		log.Debugf("user query failed, %v", err.Error())
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_COMMON,
			ErrMsg: err.Error(),
		})
	} else {
		log.Debugf("user query success, %+v", *rsp)
		base.HttpResponse(w, &base.MessageRsp{
			Data: rsp,
		})
	}
}

// user add
func (this *AuthController) UserAdd(w http.ResponseWriter, r *http.Request) {
	var req msgAuth.MsgAddUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("%v", err.Error())
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}

	log.Debugf("user add, name=%v, email=%v, phone=%v",
		req.Name, req.Email, req.Phone)

	rsp, err := this.authService.UserAdd(&req)
	if err != nil {
		log.Debugf("user add failed")
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_COMMON,
			ErrMsg: err.Error(),
		})
	} else {
		log.Debugf("user add success")
		base.HttpResponse(w, &base.MessageRsp{
			Data: rsp,
		})
	}
}
