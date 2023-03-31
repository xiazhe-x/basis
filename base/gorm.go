package base

import (
	"github.com/sirupsen/logrus"
	"github.com/xiazhe-x/basis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func OrmDatabase() *gorm.DB {
	Check(db)
	return db
}

type DatabaseStarter struct {
	basis.BaseStarter
}

func (s *DatabaseStarter) Setup(ctx basis.StarterContext) {
	conf := ctx.Props()
	var err error
	//数据库配置
	dbStr, err := conf.Get("mysql.DSN")
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dbStr, // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		CreateBatchSize: 1000, //分批创建时，指定每批的数量
		//QueryFields:     true, //QueryFields 模式会根据当前 model 的所有字段名称进行 select。
		AllowGlobalUpdate:      false, //在没有任何条件的情况下执行批量更新，默认情况下 GORM 不会执行该操作
		SkipDefaultTransaction: true,  //禁用默认事务
		PrepareStmt:            true,  //执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
	})
	if err != nil {
		logrus.Panicf("数据库连接错误：", err)
		panic(err)
	}
	// 获取通用数据库对象 sql.DB，然后使用其提供的功能
	sqlDB, err := db.DB()
	// Ping
	err = sqlDB.Ping()
	if err != nil {
		logrus.Panicf("数据库Ping错误：", err)
		panic(err)
	}
}
