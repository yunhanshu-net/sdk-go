package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	dbLock = new(sync.Mutex)
	dbs    = make(map[string]*gorm.DB)
)

// MustGetOrInitDB 获取或初始化数据库连接
// 如果数据库不存在，会自动创建
func MustGetOrInitDB(dbName string) *gorm.DB {
	dbLock.Lock()
	defer dbLock.Unlock()

	// 安全处理数据库名称，防止目录穿越攻击
	dbName = sanitizeDBName(dbName)

	// 检查缓存是否已存在连接
	if db, ok := dbs[dbName]; ok {
		return db
	}

	// 确保数据目录存在
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logrus.Errorf("创建数据目录失败: %v", err)
		panic(fmt.Errorf("创建数据目录失败: %v", err))
	}

	// 设置GORM日志配置
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// 创建数据库连接
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		logrus.Errorf("打开数据库失败 %s: %v", dbName, err)
		panic(fmt.Errorf("打开数据库失败 %s: %v", dbName, err))
	}

	// 设置SQLite优化参数
	db.Exec("PRAGMA journal_mode=WAL;PRAGMA temp_store=MEMORY;PRAGMA synchronous=NORMAL;")

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("获取原生数据库连接失败: %v", err)
		panic(fmt.Errorf("获取原生数据库连接失败: %v", err))
	}

	sqlDB.SetMaxOpenConns(5)            // 增加连接数，支持多协程
	sqlDB.SetMaxIdleConns(2)            // 保持一些空闲连接
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最长生命周期

	// 缓存连接
	dbs[dbName] = db
	logrus.Infof("数据库连接已创建: %s", dbName)

	return db
}

// sanitizeDBName 安全处理数据库名称，防止目录穿越
func sanitizeDBName(dbName string) string {
	// 移除路径前缀
	dbName = strings.TrimPrefix(dbName, "../")
	dbName = strings.TrimPrefix(dbName, "./")

	// 确保只取基本文件名，防止目录穿越
	dbName = filepath.Base(dbName)

	// 确保有.db后缀
	if !strings.HasSuffix(dbName, ".db") {
		dbName = dbName + ".db"
	}

	return "./data/" + dbName
}

// Context的MustGetOrInitDB方法
func (c *Context) MustGetOrInitDB(dbName string) *gorm.DB {
	return MustGetOrInitDB(dbName)
}

// CloseAllDBs 关闭所有数据库连接，用于程序退出时清理
func CloseAllDBs() {
	dbLock.Lock()
	defer dbLock.Unlock()

	for name, db := range dbs {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logrus.Warnf("关闭数据库连接失败 %s: %v", name, err)
			} else {
				logrus.Infof("数据库连接已关闭: %s", name)
			}
		}
	}

	// 清空连接池
	dbs = make(map[string]*gorm.DB)
}

// GetDB 获取或初始化数据库连接
func GetDB(dbName string) (*gorm.DB, error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	// 安全处理数据库名称，防止目录穿越攻击
	dbName = sanitizeDBName(dbName)

	// 检查缓存是否已存在连接
	if db, ok := dbs[dbName]; ok && db != nil {
		return db, nil
	}

	// 确保数据目录存在
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logrus.Errorf("创建数据目录失败: %v", err)
		panic(fmt.Errorf("创建数据目录失败: %v", err))
	}

	// 设置GORM日志配置
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// 创建数据库连接
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		logrus.Errorf("打开数据库失败 %s: %v", dbName, err)
		return nil, err
	}

	// 设置SQLite优化参数
	db.Exec("PRAGMA journal_mode=WAL;PRAGMA temp_store=MEMORY;PRAGMA synchronous=NORMAL;")

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("获取原生数据库连接失败: %v", err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(5)            // 增加连接数，支持多协程
	sqlDB.SetMaxIdleConns(2)            // 保持一些空闲连接
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最长生命周期

	// 缓存连接
	dbs[dbName] = db
	logrus.Infof("数据库连接已创建: %s", dbName)

	return db, nil
}
