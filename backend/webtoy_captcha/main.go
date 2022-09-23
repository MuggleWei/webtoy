package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func initComponents(conf *CaptchaConfig) {
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

var expireSecond int = 60

func main() {
	conf, err := InitConfig()
	if err != nil {
		log.Fatalf("failed init config")
		panic(err)
	}

	// init log
	base.InitLog(conf.logLevel, conf.logFile, conf.logEnableConsole)
	log.Infof("webtoy-captcha log config")

	PrintConfig(conf)

	expireSecond = conf.sessionExpireSecond

	log.Info("--------------------")
	log.Infof("webtoy-captcha launch")

	// init components
	log.Infof("init components")
	initComponents(conf)

	// http server
	log.Infof("init http server")
	router := mux.NewRouter()
	router.HandleFunc("/captcha/load", Load)
	router.HandleFunc("/captcha/verify", Verify)

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

func Load(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("captcha_session")
	if sessionID == "" {
		var err error
		sessionHandler := base.GetSessionComponent().Handler
		sessionID, err = sessionHandler.GenSessionID()
		if err != nil {
			log.Errorf("failed gen sessionID, %v", err.Error())
			panic(err)
		}
		sessionID = "captcha_" + sessionID

		cookie := http.Cookie{Name: "captcha_session", Value: sessionID, MaxAge: expireSecond}
		http.SetCookie(w, &cookie)
	}

	content, val, err := GenCaptcha()
	if err != nil {
		log.Errorf("failed gen captcha, %v", err.Error())
		panic(err)
	}

	err = SaveCaptchaVal(val, sessionID)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(content.Bytes()))
}

func Verify(w http.ResponseWriter, r *http.Request) {
	// parse request
	var req ModelCaptchaVerify
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warningf("failed parse body: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	redisClient := base.GetRedisComponent().Client
	statusCmd := redisClient.Get(req.CaptchaSessionID)
	rsp := &base.MessageRsp{}
	if statusCmd.Err() == nil && statusCmd.Val() == req.CaptchaValue {
		rsp.Code = 0
	} else {
		rsp.Code = -1
		rsp.ErrMsg = "failed verify captcha"
	}

	redisClient.Del(req.CaptchaSessionID)

	base.HttpResponse(w, rsp)
}

func GenCaptcha() (*bytes.Buffer, []byte, error) {
	d := captcha.RandomDigits(captcha.DefaultLen)

	var content bytes.Buffer
	img := captcha.NewImage("", d, captcha.StdWidth, captcha.StdHeight)
	if img == nil {
		return nil, nil, errors.New("failed to generate captcha")
	}
	_, err := img.WriteTo(&content)
	if err != nil {
		return nil, nil, err
	}

	return &content, d, nil
}

func SaveCaptchaVal(val []byte, uuid string) error {
	s := ""
	for _, v := range val {
		s = s + strconv.Itoa(int(v))
	}

	statusCmd := base.GetRedisComponent().Client.Set(uuid, s, time.Second*time.Duration(expireSecond))
	return statusCmd.Err()
}
