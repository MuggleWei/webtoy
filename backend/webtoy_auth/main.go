package main

import (
	"fmt"
	"net/http"
	"time"

	controller "github.com/MuggleWei/webtoy/backend/webtoy_auth/controller"
	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func initComponents(conf *AuthConfig) {
	// -----------------------------
	// get components
	srClientComponent := base.GetSrClientComponent()
	dbComponent := base.GetDBComponent()

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

	// init db component
	log.Infof("init db component")
	err = dbComponent.InitEngines(conf.dbCfgMaps)
	if err != nil {
		log.Fatalf("failed init mysql: %v", err.Error())
		panic(err)
	}
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()

	authController := controller.GetAuthController()
	router.HandleFunc("/user/auth", authController.UserAuth)
	router.HandleFunc("/user/query", authController.UserQuery)

	router.Use(base.MiddlewareTimeElapsed)

	return router
}

func main() {
	conf, err := InitConfig()
	if err != nil {
		log.Fatalf("failed init config")
		panic(err)
	}

	// init log
	base.InitLog(conf.logLevel, conf.logFile, conf.logEnableConsole)
	log.Infof("webtoy-auth log config")

	PrintConfig(conf)

	log.Info("--------------------")
	log.Infof("webtoy-auth launch")

	// init components
	log.Infof("init components")
	initComponents(conf)

	// init routes
	log.Infof("init routes")
	router := initRoutes()

	// http server
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%v:%v", conf.host, conf.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("run http server")
	log.Fatal(srv.ListenAndServe())
}
