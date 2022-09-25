package mapper

import (
	"sync"

	"github.com/MuggleWei/webtoy/backend/webtoy_auth/dao"
	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	msgAuth "github.com/MuggleWei/webtoy/backend/webtoy_msg_auth"
	log "github.com/sirupsen/logrus"
)

type UserMapper struct{}

var (
	singletonUserMapper *UserMapper
	onceUserMapper      sync.Once
)

func GetUserMapper() *UserMapper {
	if singletonUserMapper == nil {
		onceUserMapper.Do(func() {
			singletonUserMapper = &UserMapper{}
		})
	}
	return singletonUserMapper
}

func (this *UserMapper) Query(req *msgAuth.MsgQueryUserReq) ([]dao.DaoUser, error) {
	sourceName := "main"
	engine, err := base.GetDBComponent().GetEngine(sourceName)
	if err != nil {
		log.Errorf("failed get db instance '%v'", sourceName)
		return nil, err
	}

	session := engine.Where("")
	if req.UserID != "" {
		session = session.And("id=?", req.UserID)
	}
	if req.Name != "" {
		session = session.And("name=?", req.Name)
	}
	if req.Email != "" {
		session = session.And("email=?", req.Email)
	}
	if req.Phone != "" {
		session = session.And("phone=?", req.Phone)
	}

	var daos []dao.DaoUser
	err = session.Find(&daos)
	if err != nil {
		log.Errorf("failed query user, req=%+v, err=%v", *req, err.Error())
		return nil, err
	}

	return daos, nil
}

func (this *UserMapper) Insert(daoUser *dao.DaoUser) (int64, error) {
	sourceName := "main"
	engine, err := base.GetDBComponent().GetEngine(sourceName)
	if err != nil {
		log.Errorf("failed get db instance '%v'", sourceName)
		return 0, err
	}

	affect, err := engine.Insert(daoUser)
	if err != nil {
		passwd := daoUser.Passwd
		daoUser.Passwd = "******"
		log.Errorf("failed get insert new user: %+v", *daoUser)
		daoUser.Passwd = passwd
		return 0, err
	}

	passwd := daoUser.Passwd
	daoUser.Passwd = "******"
	log.Infof("success insert new user: %+v", *daoUser)
	daoUser.Passwd = passwd

	return affect, nil
}
