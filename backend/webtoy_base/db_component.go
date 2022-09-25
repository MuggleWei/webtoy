package webtoy_base

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

type DBComponent struct {
	instanceMap map[string]*xorm.Engine
}

var (
	singletonDB *DBComponent
	onceDB      sync.Once
)

func GetDBComponent() *DBComponent {
	if singletonDB == nil {
		onceDB.Do(func() {
			singletonDB = &DBComponent{
				instanceMap: make(map[string]*xorm.Engine),
			}
		})
	}
	return singletonDB
}

func (this *DBComponent) GetEngine(name string) (*xorm.Engine, error) {
	engine, ok := this.instanceMap[name]
	if !ok {
		errMsg := fmt.Sprintf("Db GetInstance not exists: %v", name)
		err := errors.New(errMsg)
		log.Error(err.Error())
		return nil, err
	} else {
		return engine, nil
	}
}

func (this *DBComponent) InitEngines(dbMap map[string]DBConfig) error {
	this.instanceMap = make(map[string]*xorm.Engine)

	for name, cfg := range dbMap {
		_, ok := this.instanceMap[name]
		if ok {
			errMsg := fmt.Sprintf("Repeated add data source: %v", name)
			log.Errorf(errMsg)
			return errors.New(errMsg)
		}

		switch cfg.Driver {
		case "mysql":
			this.AddDBMysql(name, cfg)
		default:
			log.Errorf("unsupport driver: %v", cfg.Driver)
		}
	}

	return nil
}

func (this *DBComponent) AddDBMysql(name string, cfg DBConfig) error {
	log.Infof("insert db mysql: name=%v, net=%v, addr=%v, params=%v, db=%v, user=%v, passwd=******",
		cfg.Net, cfg.Net, cfg.Addr, cfg.Params, cfg.User)

	params := make(map[string]string)
	paramList := strings.Split(cfg.Params, "&")
	for _, param := range paramList {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			errMsg := fmt.Sprintf("invalid db params: %v", param)
			log.Errorf(errMsg)
			return errors.New(errMsg)
		}
		params[kv[0]] = kv[1]
	}

	mysqlCfg := mysql.Config{
		DBName:               cfg.Db,
		User:                 cfg.User,
		Passwd:               cfg.Passwd,
		Net:                  cfg.Net,
		Addr:                 cfg.Addr,
		Params:               params,
		AllowNativePasswords: true,
	}
	engine, err := xorm.NewEngine("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return err
	}

	// set max connection number
	engine.SetMaxOpenConns(64)

	// set max idle connection number, MaxIdleConns need <= MaxOpenConns
	engine.SetMaxIdleConns(8)

	// set max life time
	engine.SetConnMaxLifetime(time.Hour)

	this.instanceMap[name] = engine

	return nil
}
