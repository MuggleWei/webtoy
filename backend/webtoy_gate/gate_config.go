package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type GateConfig struct {
	// service
	host string
	port uint

	// log
	logLevel         string
	logFile          string
	logEnableConsole bool

	// service registry
	srHost string
	srPort uint

	srServiceHost string
	srServicePort uint
	srServiceName string
	srServiceID   string
	srServiceTag  string
	srServiceTTL  time.Duration

	// redis
	redisHost   string
	redisPort   uint
	redisPasswd string
	redisDb     int

	// session
	sessionExpireSecond int
}

func InitConfig() (*GateConfig, error) {
	// pflag
	pflag.String("host", "0.0.0.0", "bind host")
	pflag.Uint("port", 8080, "listen port")

	pflag.String("log.level", "info", "log level")
	pflag.String("log.file", "./log/webtoy-gate.log", "log file path")
	pflag.Bool("log.console", false, "enable/disable log console output")

	pflag.String("sr.host", "127.0.0.1", "sr host")
	pflag.Uint("sr.port", 8500, "sr port")
	pflag.String("sr.service.host", "0.0.0.0", "service bind host")
	pflag.Uint("sr.service.port", 8080, "service listen port")
	pflag.String("sr.service.name", "webtoy-gate", "service name")
	pflag.String("sr.service.id", "webtoy-gate-0", "service id")
	pflag.String("sr.service.tag", "", "sr service tags")
	pflag.String("sr.service.ttl", "3s", "sr ttl")

	pflag.String("redis.host", "127.0.0.1", "redis host")
	pflag.Uint("redis.port", 6379, "redis port")
	pflag.Int("redis.db", 0, "redis db")
	pflag.String("redis.passwd", "", "redis.passwd")

	pflag.String("session.expired", "1d", "session expired time duration")

	pflag.Parse()

	// config
	viper.SetConfigName("gate")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/webtoy_gate")
	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			fmt.Fprintf(os.Stdout, "config file not found\n")
		} else {
			panic(fmt.Errorf("error config file: %v", err))
		}
	}

	// viper bind command line
	viper.BindPFlags(pflag.CommandLine)

	ttl, err := time.ParseDuration(viper.GetString("sr.service.ttl"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed parse sr ttl value: %v", err.Error())
		return nil, err
	}

	sessionExpired, err := time.ParseDuration(viper.GetString("session.expired"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed parse sesseion expired value: %v", err.Error())
		return nil, err
	}

	return &GateConfig{
		host: viper.GetString("host"),
		port: viper.GetUint("port"),

		logLevel:         viper.GetString("log.level"),
		logFile:          viper.GetString("log.file"),
		logEnableConsole: viper.GetBool("log.console"),

		srHost:        viper.GetString("sr.host"),
		srPort:        viper.GetUint("sr.port"),
		srServiceHost: viper.GetString("sr.service.host"),
		srServicePort: viper.GetUint("sr.service.port"),
		srServiceName: viper.GetString("sr.service.name"),
		srServiceID:   viper.GetString("sr.service.id"),
		srServiceTag:  viper.GetString("sr.service.tag"),
		srServiceTTL:  ttl,

		redisHost:   viper.GetString("redis.host"),
		redisPort:   viper.GetUint("redis.port"),
		redisDb:     viper.GetInt("redis.db"),
		redisPasswd: viper.GetString("redis.passwd"),

		sessionExpireSecond: int(sessionExpired.Seconds()),
	}, nil
}

func PrintConfig(conf *GateConfig) {
	log.Info("--------------------")
	log.Info("gate config:")
	log.Infof("host=%v, port=%v", conf.host, conf.port)
	log.Infof("log.level=%v, log.file=%v, log.console=%v",
		conf.logLevel, conf.logFile, conf.logEnableConsole)
	log.Infof("sr.host=%v, sr.port=%v, sr.service.host=%v, sr.service.port=%v, sr.service.name=%v, sr.service.id=%v, sr.service.tag=%v, sr.service.ttl=%v",
		conf.srHost, conf.srPort,
		conf.srServiceHost, conf.srServicePort,
		conf.srServiceName, conf.srServiceID, conf.srServiceTag, conf.srServiceTTL)
	log.Infof("redis.host=%v, redis.port=%v, redis.passwd=******, redis.db=%v",
		conf.redisHost, conf.redisPort, conf.redisDb, conf.redisPasswd)
	log.Infof("session.expired=%vs", conf.sessionExpireSecond)
}
