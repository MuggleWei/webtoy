package service

import "sync"

type UserService struct{}

var (
	singletonUser *UserService
	onceUser      sync.Once
)

func GetUserService() *UserService {
	if singletonUser == nil {
		onceUser.Do(func() {
			singletonUser = &UserService{}
		})
	}
	return singletonUser
}
