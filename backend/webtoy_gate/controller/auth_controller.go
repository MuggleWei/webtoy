package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	"github.com/MuggleWei/webtoy/backend/webtoy_gate/service"
	msgAuth "github.com/MuggleWei/webtoy/backend/webtoy_msg_auth"
	log "github.com/sirupsen/logrus"
)

type AuthController struct {
	authService    *service.AuthService
	captchaService *service.CaptchaService
}

var (
	singletonUser *AuthController
	onceUser      sync.Once
)

func GetAuthController() *AuthController {
	if singletonUser == nil {
		onceUser.Do(func() {
			singletonUser = &AuthController{
				authService:    service.GetAuthService(),
				captchaService: service.GetCaptchaService(),
			}
		})
	}
	return singletonUser
}

// user check
func (this *AuthController) UserCheck(w http.ResponseWriter, r *http.Request) {
	base.HttpResponse(w, &base.MessageRsp{})
}

// user login
func (this *AuthController) UserLogin(w http.ResponseWriter, r *http.Request) {
	// parse request
	var req msgAuth.MsgAuthUserReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errMsg := fmt.Sprintf("failed parse body: %v", err.Error())
		log.Errorf("%v", errMsg)
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: errMsg,
		})
		return
	}

	log.Debugf("UserController Login, name=%v, email=%v, phone=%v, passwd=******, captcha_session=%v, captcha_value=%v",
		req.Name, req.Email, req.Phone, req.CaptchaSession, req.CaptchaValue)

	// verify captcher
	ret, err := this.captchaService.Verify(req.CaptchaSession, req.CaptchaValue)
	if err != nil {
		log.Errorf("failed captcha verify: %v", err.Error())
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}

	if !ret {
		log.Debugf("captcha verify failed")
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_CAPTCHA,
			ErrMsg: "captcha verify error",
		})
		return
	}
	log.Debugf("captcha verify pass")

	// user auth
	rsp, err := this.authService.UserAuth(&req)
	if err != nil {
		log.Errorf("auth exception")
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}

	if rsp.Code != 0 {
		log.Debugf("auth failed")
		base.HttpResponse(w, rsp)
		return
	}

	log.Debugf("user auth pass")

	rspData, ok := rsp.Data.(*msgAuth.MsgAuthUserRsp)
	if !ok {
		errMsg := fmt.Sprint("failed get reponse data")
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: errMsg,
		})
		return
	}

	// save session
	sessionHandler := base.GetSessionComponent().Handler
	userSession, err := sessionHandler.GenSession(rspData.UserID)
	if err != nil {
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}
	rspData.Session = userSession.Session
	rspData.Token = userSession.Token

	// response
	base.HttpResponse(w, rsp)
}

// user register
func (this *AuthController) UserRegister(w http.ResponseWriter, r *http.Request) {
}

// get user profile
func (this *AuthController) UserProfile(w http.ResponseWriter, r *http.Request) {
	// already auth by middleware
	userID := r.Header.Get(base.SESSION_USER)

	rsp, err := this.authService.UserQuery(&msgAuth.MsgQueryUserReq{
		UserID: userID,
	})
	if err != nil {
		log.Error("failed query user, %v", err.Error())
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}

	base.HttpResponse(w, rsp)
}
