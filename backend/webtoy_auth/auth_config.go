package main

import (
	"fmt"
	"os"
	"time"

	base "github.com/MuggleWei/webtoy/backend/webtoy_base"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type AuthConfig struct {
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

	// mysql
	dbCfgMaps map[string]base.DBConfig
}

func InitConfig() (*AuthConfig, error) {
	// pflag
	pflag.String("host", "0.0.0.0", "bind host")
	pflag.Uint("port", 8080, "listen port")

	pflag.String("log.level", "info", "log level")
	pflag.String("log.file", "./log/webtoy-auth.log", "log file path")
	pflag.Bool("log.console", false, "enable/disable log console output")

	pflag.String("sr.host", "127.0.0.1", "sr host")
	pflag.Uint("sr.port", 8500, "sr port")
	pflag.String("sr.service.host", "0.0.0.0", "service bind host")
	pflag.Uint("sr.service.port", 8080, "service listen port")
	pflag.String("sr.service.name", "webtoy-auth", "service name")
	pflag.String("sr.service.id", "webtoy-auth-0", "service id")
	pflag.String("sr.service.tag", "", "sr service tags")
	pflag.String("sr.service.ttl", "3s", "sr ttl")

	pflag.String("mysql.main.addr", "127.0.0.1:3306", "mysql main address")
	pflag.String("mysql.main.params", "charset=utf8", "mysql main params")
	pflag.String("mysql.main.db", "webtoy", "mysql main db")
	pflag.String("mysql.main.user", "muggle", "mysql main user")
	pflag.String("mysql.main.passwd", "wsz123", "mysql main passwd")

	pflag.Parse()

	// config
	viper.SetConfigName("auth")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/webtoy_auth")
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

	return &AuthConfig{
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

		dbCfgMaps: map[string]base.DBConfig{
			"main": {
				Driver: viper.GetString("db.main.driver"),
				Net:    viper.GetString("db.main.net"),
				Addr:   viper.GetString("db.main.addr"),
				Params: viper.GetString("db.main.params"),
				Db:     viper.GetString("db.main.db"),
				User:   viper.GetString("db.main.user"),
				Passwd: viper.GetString("db.main.passwd"),
			},
		},
	}, nil
}

func PrintConfig(conf *AuthConfig) {
	log.Info("--------------------")
	log.Info("auth config:")
	log.Infof("host=%v, port=%v", conf.host, conf.port)
	log.Infof("log.level=%v, log.file=%v, log.console=%v",
		conf.logLevel, conf.logFile, conf.logEnableConsole)
	log.Infof("sr.host=%v, sr.port=%v, sr.service.host=%v, sr.service.port=%v, sr.service.name=%v, sr.service.id=%v, sr.service.tag=%v, sr.service.ttl=%v",
		conf.srHost, conf.srPort,
		conf.srServiceHost, conf.srServicePort,
		conf.srServiceName, conf.srServiceID, conf.srServiceTag, conf.srServiceTTL)
	for name, dbCfg := range conf.dbCfgMaps {
		log.Infof("db.%v: driver=%v, net=%v, addr=%v, params=%v, db=%v, user=%v, passwd=******",
			name, dbCfg.Driver, dbCfg.Net, dbCfg.Addr, dbCfg.Params, dbCfg.Db, dbCfg.User)
	}
}
