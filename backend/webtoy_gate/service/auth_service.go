package service

import (
	"net/http"
	"sync"
	"time"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	"github.com/MuggleWei/webtoy/backend/webtoy_gate/utils"
	msgAuth "github.com/MuggleWei/webtoy/backend/webtoy_msg_auth"
)

type AuthService struct {
	transport *http.Transport
}

var (
	singletonAuth *AuthService
	onceAuth      sync.Once
)

func GetAuthService() *AuthService {
	if singletonAuth == nil {
		onceAuth.Do(func() {
			singletonAuth = &AuthService{
				transport: &http.Transport{
					MaxIdleConns:        0,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     30 * time.Second,
				},
			}
		})
	}
	return singletonAuth
}

func (this *AuthService) UserAuth(req *msgAuth.MsgAuthUserReq) (*base.MessageRsp, error) {
	urlPath := "/user/auth"
	var rspData msgAuth.MsgAuthUserRsp
	return base.HttpSRPost(utils.AuthServiceName, urlPath, this.transport, req, &rspData)
}

func (this *AuthService) UserQuery(req *msgAuth.MsgQueryUserReq) (*base.MessageRsp, error) {
	urlPath := "/user/query"
	var rspData msgAuth.MsgQueryUserRsp
	return base.HttpSRPost(utils.AuthServiceName, urlPath, this.transport, req, &rspData)
}
