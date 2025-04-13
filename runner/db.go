package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

var dbLock = new(sync.Mutex)
var dbs = make(map[string]*gorm.DB)

func (c *Context) MustGetOrInitDB(dbName string) *gorm.DB {
	dbLock.Lock()
	defer dbLock.Unlock()
	if db, ok := dbs[dbName]; ok {
		return db
	}
	dbName = strings.TrimPrefix(dbName, "../")
	dbName = strings.TrimPrefix(dbName, "./")
	dbName = "./data/" + dbName
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		// 添加 PRAGMA 配置
		DisableForeignKeyConstraintWhenMigrating: true,                                           // 可选：禁用外键约束（按需）
		NowFunc:                                  func() time.Time { return time.Now().Local() }, // 可选：本地时间

	})
	if err != nil {
		logrus.Errorf("gorm Open db %s err:%s", dbName, err.Error())
		//不必紧张，框架层面有做recover 处理，也可以自己recover来捕获错误
		panic(fmt.Errorf("gorm Open db %s err:%s", dbName, err.Error()))
	}
	dbs[dbName] = db
	return db
}
