package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/MuggleWei/webtoy/backend/webtoy_gate/model"
	"github.com/MuggleWei/webtoy/backend/webtoy_gate/service"
	log "github.com/sirupsen/logrus"
)

type UserController struct {
	userService *service.UserService
}

var (
	singletonUser *UserController
	onceUser      sync.Once
)

func GetUserController() *UserController {
	if singletonUser == nil {
		onceUser.Do(func() {
			singletonUser = &UserController{
				userService: service.GetUserService(),
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
		log.Warningf("failed parse body: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: check captcher

	// TODO: check user password
}

// user register
func (this *UserController) Register(w http.ResponseWriter, r *http.Request) {
}

// get user profile
func (this *UserController) Profile(w http.ResponseWriter, r *http.Request) {
}
