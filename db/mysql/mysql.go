package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Options defines optsions for mysql database.
type Options struct {
	Host                  string // MySQL 数据库主机地址
	User                  string // MySQL 数据库用户名
	Password              string // MySQL 数据库密码
	DB                    string // MySQL 数据库名称
	Port                  int    // MySQL 数据库端口号
	Timeout               string // 连接超时时间
	ReadTimeout           string // 读取超时时间
	WriteTimeout          string // 写入超时时间
	Loc                   string // 时区
	Charset               string // 字符集
	ParseTime             bool   // 是否解析时间
	MaxOpenConns          int    // 最大连接数
	MaxIdleConns          int    // 最大空闲连接数
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}

var (
	dbConn *gorm.DB
)

// New create a new gorm db instance with the given options.
func New(opts *Options) error {
	var dsn string = mysqlDsn(opts)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键生成
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: opts.Logger,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConns)

	dbConn = db
	return nil
}

func mysqlDsn(opts *Options) string {
	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s",
		opts.User, opts.Password, opts.Host, opts.Port, opts.DB, opts.Charset,
		opts.ParseTime, opts.Loc, opts.Timeout, opts.ReadTimeout, opts.WriteTimeout)
	return dsn
}

func GetDB() *gorm.DB {
	return dbConn.Debug()
}
