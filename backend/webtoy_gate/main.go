package main

import (
	"fmt"
	"net/http"
	"time"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	controller "github.com/MuggleWei/webtoy/backend/webtoy_gate/controller"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func initComponents(conf *GateConfig) {
	// -----------------------------
	// get components
	srClientComponent := base.GetSrClientComponent()
	redisComponent := base.GetRedisComponent()
	sessionComponent := base.GetSessionComponent()

	// -----------------------------
	// init srclient component
	log.Infof("init service registry component")
	srClientArgs := &base.SRClientArgs{
		SrAddr:        fmt.Sprintf("%v:%v", conf.srHost, conf.srPort),
		SrServiceID:   conf.srServiceID,
		SrServiceName: conf.srServiceName,
		SrServiceHost: conf.srServiceHost,
		SrServicePort: conf.srServicePort,
		SrServiceTag:  conf.srServiceTag,
		SrServiceTTL:  conf.srServiceTTL,
	}
	err := srClientComponent.Init(srClientArgs)
	if err != nil {
		log.Fatalf("failed init srclient: %v", err.Error())
		panic(err)
	}

	// init redis component
	log.Infof("init redis component")
	err = redisComponent.Init(conf.redisHost, conf.redisPort, conf.redisPasswd, conf.redisDb)
	if err != nil {
		log.Fatalf("failed init redis: %v", err.Error())
		panic(err)
	}

	// init session component
	log.Infof("init session component")
	sessionComponent.Handler.SessionExpireSecond = conf.sessionExpireSecond

	// -----------------------------
	// Dependency Injection
	log.Infof("Dependency Injection")
	sessionComponent.Handler.RedisClient = redisComponent.Client
}

func main() {
	conf, err := InitConfig()
	if err != nil {
		log.Fatalf("failed init config")
		panic(err)
	}

	// init log
	base.InitLog(conf.logLevel, conf.logFile, conf.logEnableConsole)
	log.Infof("webtoy-gate log config")

	PrintConfig(conf)

	log.Info("--------------------")
	log.Infof("webtoy-gate launch")

	// init components
	log.Infof("init components")
	initComponents(conf)

	// http server
	log.Infof("init http server")
	router := mux.NewRouter()

	echoController := controller.GetEchoController()
	router.HandleFunc("/api/v1/echo", echoController.Echo)

	userController := controller.GetUserController()
	userRouter := router.PathPrefix("/api/v1/user").Subrouter()
	userRouter.HandleFunc("/login", userController.Login)
	userRouter.HandleFunc("/register", userController.Register)
	userRouter.HandleFunc("/profile", userController.Profile)

	captchaController := controller.GetCaptchaController()
	router.HandleFunc("/api/v1/captcha/load", captchaController.Load)

	router.Use(base.MiddlewareTimeElapsed)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%v:%v", conf.host, conf.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("run http server")
	log.Fatal(srv.ListenAndServe())
}