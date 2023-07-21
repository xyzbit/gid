package loader

import (
	"time"

	"github.com/pkg/errors"

	"github.com/xyzbit/gid/core"
	"github.com/xyzbit/gid/core/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StatusLoader struct {
	db *gorm.DB
}

func NewStatusLoader(cfg *conf.DBConfig) core.StatusLoader {
	return &StatusLoader{
		db: NewDB(cfg),
	}
}

// NewDB new db instance.
func NewDB(cfg *conf.DBConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.Source))
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	if cfg.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.OpenConns != 0 {
		sqlDB.SetMaxOpenConns(cfg.OpenConns)
	}
	if cfg.MaxLifetimeSec != 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeSec) * time.Second)
	}

	return db
}

type Status struct {
	ID     int32  `gorm:"primaryKey"`
	Name   string `gorm:"column:name"`
	LastID int64  `gorm:"column:last_id"`
}

func (s *Status) TableName() string {
	return "schedule"
}

func (l *StatusLoader) AllStatus() ([]*core.Status, error) {
	var ss []*Status
	if err := l.db.Find(&ss).Error; err != nil {
		return nil, errors.Wrap(err, "status load error")
	}
	ret := make([]*core.Status, len(ss))
	for i, s := range ss {
		ret[i] = &core.Status{
			ModuleID:   s.ID,
			ModuleName: s.Name,
			LastID:     s.LastID,
		}
	}
	return ret, nil
}

func (l *StatusLoader) Status(name string) (*core.Status, error) {
	s := &Status{}
	if err := l.db.Where("name=?", name).Take(s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || s.ID == 0 {
			return &core.Status{
				ModuleName: name,
				LastID:     0,
			}, nil
		}
		return nil, errors.Wrap(err, "status load error")
	}
	return &core.Status{
		ModuleID:   s.ID,
		ModuleName: s.Name,
		LastID:     s.LastID,
	}, nil
}

// 注销 配置和状态都删除
// 重置 配置不变，状态删除（会自动重建）；重置到某一个id

func (l *StatusLoader) SaveStatus(name string, s *core.Status) error {
	status := &Status{
		ID:     s.ModuleID,
		Name:   s.ModuleName,
		LastID: s.LastID,
	}

	var err error
	if s.ModuleID == 0 {
		err = l.db.Create(status).Error
		s.ModuleID = status.ID
	} else {
		err = l.db.Save(status).Error
	}
	if err != nil {
		return errors.Wrap(err, "status save error")
	}
	return nil
}
