package app

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dao struct {
	app *App
	*gorm.DB
}

type DaoConfig struct {
	UserName     string
	Password     string
	Host         string
	Port         int
	DatabaseName string
}

func newDao(cfg *DaoConfig, app *App) *Dao {
	return &Dao{
		DB:  initDb(cfg),
		app: app,
	}
}

func (d *Dao) UsePlugin(plugins []gorm.Plugin) error {

	for _, plugin := range plugins {

		if err := d.Use(plugin); err != nil {
			return err
		}
	}

	return nil
}

func initDb(cfg *DaoConfig) *gorm.DB {
	db, err := gorm.Open(mysqlConfig(cfg), gormConfig(cfg))
	if err != nil {
		panic(err)
	}
	return db
}

func mysqlConfig(config *DaoConfig) mysql.Dialector {

	return mysql.Dialector{Config: &mysql.Config{
		DSN:                       mysqlDsn(config), // DSN data source name
		DefaultStringSize:         191,              // string 类型字段的默认长度
		DisableDatetimePrecision:  true,             // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,             // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,             // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,            // 根据版本自动配置
	}}
}

func mysqlDsn(c *DaoConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&&loc=Local", c.UserName, c.Password, c.Host, c.Port, c.DatabaseName)
}

func gormConfig(config *DaoConfig) *gorm.Config {

	return &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   logger.Default.LogMode(logger.Info),
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}
}
