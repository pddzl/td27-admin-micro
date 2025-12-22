package initialization

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志
//func (w *writer) Printf(message string, data ...interface{}) {
//	if global.TD27_CONFIG.Mysql.LogZap {
//		global.TD27_LOG.Info(fmt.Sprintf(message+"\n", data...))
//	} else {
//		w.Writer.Printf(message, data...)
//	}
//}

func Gorm(
	username string,
	password string,
	host string,
	port uint,
	dbName string,
	config string,
	maxOpenConn int,
	maxIdleConn int,
	maxLifetimeSec int64,
) (*gorm.DB, error) {

	dsn := username + ":" + password + "@tcp(" + fmt.Sprintf("%s:%d", host, port) + ")/" + dbName + "?" + config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig()); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(maxIdleConn)
		sqlDB.SetMaxOpenConns(maxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeSec) * time.Second)
		return db, nil
	}
}

func gormConfig() *gorm.Config {
	newLogger := logger.New(
		NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.Warn,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                  // 禁用彩色打印
		},
	)
	config := &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	return config
}
