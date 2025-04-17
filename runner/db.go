package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var dbLock = new(sync.Mutex)
var dbs = make(map[string]*gorm.DB)

func MustGetOrInitDB(dbName string) *gorm.DB {
	dbLock.Lock()
	defer dbLock.Unlock()
	dbName = strings.TrimPrefix(dbName, "../")
	dbName = strings.TrimPrefix(dbName, "./")
	dbName = "./data/" + dbName
	if db, ok := dbs[dbName]; ok {
		return db
	}
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		logrus.Errorf("gorm Open db %s err:%s", dbName, err.Error())
		//不必慌张，sdk层有做recover 处理，也可以自己recover来提前捕获错误
		panic(fmt.Errorf("gorm Open db %s err:%s", dbName, err.Error()))
	}
	db.Exec("PRAGMA journal_mode=WAL;PRAGMA temp_store=MEMORY;")
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(1)
	dbs[dbName] = db
	return db
}

func (c *Context) MustGetOrInitDB(dbName string) *gorm.DB {
	return MustGetOrInitDB(dbName)
}
