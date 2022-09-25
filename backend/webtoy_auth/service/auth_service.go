package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/MuggleWei/webtoy/backend/webtoy_auth/mapper"
	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	msgAuth "github.com/MuggleWei/webtoy/backend/webtoy_msg_auth"
	log "github.com/sirupsen/logrus"
)

type AuthService struct {
	userMapper *mapper.UserMapper
}

var (
	singletonAuth *AuthService
	onceAuth      sync.Once
)

func GetAuthService() *AuthService {
	if singletonAuth == nil {
		onceAuth.Do(func() {
			singletonAuth = &AuthService{
				userMapper: mapper.GetUserMapper(),
			}
		})
	}
	return singletonAuth
}

func (this *AuthService) UserAuth(req *msgAuth.MsgAuthUserReq) (*msgAuth.MsgAuthUserRsp, error) {
	qry := &msgAuth.MsgQueryUserReq{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
	daos, err := this.userMapper.Query(qry)
	if err != nil {
		log.Errorf("failed query user, query=%v, err=%v", *qry, err.Error())
		return nil, err
	}

	if len(daos) != 1 {
		err = errors.New("result not equal 1")
		log.Errorf("failed auth user, query=%v, err=%v", *qry, err.Error())
		return nil, err
	}

	daoUser := &daos[0]
	if base.BCryptMatchPasswd(daoUser.Passwd, req.Passwd) != true {
		err = errors.New("incorrect password")
		log.Errorf("failed auth user, query=%v, err=%v", qry, err.Error())
		return nil, err
	}

	rsp := &msgAuth.MsgAuthUserRsp{
		UserID: fmt.Sprint(daoUser.Id),
	}

	return rsp, nil
}

func (this *AuthService) UserAdd(req *msgAuth.MsgAddUserReq) (*msgAuth.MsgAddUserRsp, error) {
	// TODO
	return nil, nil
}
