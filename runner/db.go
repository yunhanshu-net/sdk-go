package runner

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var dbLock = new(sync.Mutex)
var dbs = make(map[string]*gorm.DB)

// GetAndInitDB 获取数据库连接，如果不存在则初始化，GetAndInitDB("app.db")
func (c *HttpContext) GetAndInitDB(dbName string) (*gorm.DB, error) {
	dbLock.Lock()
	defer dbLock.Unlock()
	if db, ok := dbs[dbName]; ok {
		return db, nil
	}
	dbName = strings.TrimPrefix(dbName, "../")
	dbName = strings.TrimPrefix(dbName, "./")
	dbName = "./data/" + dbName
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dbs[dbName] = db
	return db, nil
}
