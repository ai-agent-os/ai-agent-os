package model

import (
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	// 先迁移外键父表，再迁移子表，避免外键约束错误
	err := db.AutoMigrate(
		&Nats{},
		&Host{},
		&User{},
		&App{},
		&ServiceTree{},
		&Function{},
		&Package{},
		// 新增的认证相关表
		&EmailVerification{},
		&EmailCode{},
		&UserSession{},
	)
	if err != nil {
		return err
	}

	// 创建默认的NATS和Host记录
	return initDefaultData(db)
}

// initDefaultData 初始化默认数据
func initDefaultData(db *gorm.DB) error {
	// 检查是否已有NATS记录
	var natsCount int64
	if err := db.Model(&Nats{}).Count(&natsCount).Error; err != nil {
		return err
	}

	// 如果没有NATS记录，创建默认记录
	if natsCount == 0 {
		defaultNats := &Nats{
			Host: "localhost",
			Port: 4222,
		}
		if err := db.Create(defaultNats).Error; err != nil {
			return err
		}
	}

	// 检查是否已有Host记录
	var hostCount int64
	if err := db.Model(&Host{}).Count(&hostCount).Error; err != nil {
		return err
	}

	// 如果没有Host记录，创建默认记录
	if hostCount == 0 {
		// 获取第一个NATS记录
		var nats Nats
		if err := db.First(&nats).Error; err != nil {
			return err
		}

		defaultHost := &Host{
			Domain:   "localhost",
			NatsID:   nats.ID,
			Status:   "active",
			AppCount: 0,
		}
		if err := db.Create(defaultHost).Error; err != nil {
			return err
		}
	}

	return nil
}
