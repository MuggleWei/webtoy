package controller

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	"github.com/MuggleWei/webtoy/backend/webtoy_gate/service"
	log "github.com/sirupsen/logrus"
)

type CaptchaController struct {
	captchaService *service.CaptchaService
}

var (
	singletonCaptcha *CaptchaController
	onceCaptcha      sync.Once
)

func GetCaptchaController() *CaptchaController {
	if singletonCaptcha == nil {
		onceCaptcha.Do(func() {
			singletonCaptcha = &CaptchaController{
				captchaService: service.GetCaptchaService(),
			}
		})
	}
	return singletonCaptcha
}

func (this *CaptchaController) Load(w http.ResponseWriter, r *http.Request) {
	captchaSession := ""
	cookie, err := r.Cookie("captcha_session")
	if err == nil {
		captchaSession = cookie.Value
	}

	b, err := this.captchaService.Load(captchaSession, w)
	if err != nil {
		log.Errorf("failed captcha service load")
		base.HttpResponse(w, &base.MessageRsp{
			Code:   base.ERROR_INTERNAL,
			ErrMsg: err.Error(),
		})
		return
	}

	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(b))
}
