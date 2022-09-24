package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	"github.com/MuggleWei/webtoy/backend/webtoy_gate/model"
	"github.com/MuggleWei/webtoy/backend/webtoy_gate/service"
	log "github.com/sirupsen/logrus"
)

type UserController struct {
	userService    *service.UserService
	captchaService *service.CaptchaService
}

var (
	singletonUser *UserController
	onceUser      sync.Once
)

func GetUserController() *UserController {
	if singletonUser == nil {
		onceUser.Do(func() {
			singletonUser = &UserController{
				userService:    service.GetUserService(),
				captchaService: service.GetCaptchaService(),
			}
		})
	}
	return singletonUser
}

// user login
func (this *UserController) Login(w http.ResponseWriter, r *http.Request) {
	// parse request
	var req model.ModelUserLoginReq
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

	log.Debugf("UserController Login, name=%v, email=%v, phone=%v, captcha_session=%v, captcha_value=%v",
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
			Code:   base.ERROR_AUTH,
			ErrMsg: "captcha verify error",
		})
		return
	}
	log.Debugf("captcha verify pass")

	// TODO: check user password
	userID := int64(10000)

	// gen session & token
	sessionHandler := base.GetSessionComponent().Handler
	sessionID, session, err := sessionHandler.GenSession(userID)
	if err != nil {
		log.Errorf("failed generate session for user: %v", userID)
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: "failed generate session",
		})
		return
	}

	base.HttpResponse(w, &base.MessageRsp{
		Code: 0,
		Data: model.ModelUserLoginRsp{
			Id:      userID,
			Session: sessionID,
			Token:   session.Token,
		},
	})
}

// user register
func (this *UserController) Register(w http.ResponseWriter, r *http.Request) {
}

// get user profile
func (this *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	// TODO:
	base.HttpResponse(w, &base.MessageRsp{})
}
