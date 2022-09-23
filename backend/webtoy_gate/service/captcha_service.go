package service

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	"github.com/MuggleWei/webtoy/backend/webtoy_gate/utils"
	log "github.com/sirupsen/logrus"
)

type CaptchaService struct {
	transport *http.Transport
}

var (
	singletonCaptcha *CaptchaService
	onceCaptcha      sync.Once
)

func GetCaptchaService() *CaptchaService {
	if singletonCaptcha == nil {
		onceCaptcha.Do(func() {
			singletonCaptcha = &CaptchaService{
				transport: &http.Transport{
					MaxIdleConns:        0,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     30 * time.Second,
				},
			}
		})
	}
	return singletonCaptcha
}

func (this *CaptchaService) Load(captchaSession string, w http.ResponseWriter) ([]byte, error) {
	srClient := base.GetSrClientComponent().Client
	addr, err := srClient.ClientLB.GetService(utils.CaptchaServiceName)
	if err != nil {
		errMsg := fmt.Sprintf("failed get service %v address", utils.CaptchaServiceName)
		log.Errorf("%v", errMsg)
		return nil, errors.New(errMsg)
	}

	url := "http://" + addr + "/captcha/load"
	if captchaSession != "" {
		url = url + "?captcha_session=" + captchaSession
	}

	b, err := base.HttpTransportGet(url, this.transport, w)
	if err != nil {
		errMsg := fmt.Sprintf("failed get transfer to service %v address", utils.CaptchaServiceName)
		log.Errorf("%v", errMsg)
		return nil, errors.New(errMsg)
	}

	return b, nil
}
