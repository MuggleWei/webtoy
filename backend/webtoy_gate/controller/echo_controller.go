package controller

import (
	"io"
	"net/http"
	"sync"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
)

type EchoController struct{}

var (
	singletonEcho *EchoController
	onceEcho      sync.Once
)

func GetEchoController() *EchoController {
	if singletonEcho == nil {
		onceEcho.Do(func() {
			singletonEcho = &EchoController{}
		})
	}
	return singletonEcho
}

func (this *EchoController) Echo(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		base.HttpResponse(w, &base.MessageRsp{
			Code:   -1,
			ErrMsg: err.Error(),
		})
	}

	base.HttpResponse(w, &base.MessageRsp{
		Data: string(b),
	})
}
