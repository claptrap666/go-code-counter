package core

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var CurrentDB *gorm.DB = nil

func DSN() string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=Asia/Shanghai",
		CurrentConfig.DB.Host,
		CurrentConfig.DB.Port,
		CurrentConfig.DB.User,
		CurrentConfig.DB.Password,
		CurrentConfig.DB.Database)
}

func ConnectDB() {
	dsn := DSN()
	var err error
	CurrentDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.WithError(err).Errorln()
	}
	// CurrentDB.AutoMigrate(&models.Plan{}, &models.Subscribe{})
	sqlDB, err := CurrentDB.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		log.WithError(err).Errorln()
	}
}
