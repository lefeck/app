package repository

import (
	"app/database"
	"context"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, rdb *database.RedisDB) Repository {
	r := &repository{
		user: NewUserRepository(db, rdb),
		post: NewPostRepository(db, rdb),
		db:   db,
		rdb:  rdb,
	}
	r.migrants = getMigrants(r.user, r.post)
	return r
}

func getMigrants(obj ...interface{}) []Migranter {
	var migrants []Migranter
	for _, obj := range migrants {
		if m, ok := obj.(Migranter); ok {
			migrants = append(migrants, m)
		}
	}
	return migrants
}

// 迁移接口
type Migranter interface {
	Migrant() error
}

// 创建数据库的表结构
func (r *repository) Migrant() error {
	for _, m := range r.migrants {
		if err := m.Migrant(); err != nil {
			return err
		}
	}
	return nil
}

type repository struct {
	user UserRepository
	//group    GroupRepository
	post PostRepository
	//rbac     RBACRepository
	db       *gorm.DB
	rdb      *database.RedisDB
	migrants []Migranter
}

func (r *repository) Migrate() error {
	return r.Migrant()
}

func (r *repository) User() UserRepository {
	return r.user
}

func (r *repository) Post() PostRepository {
	return r.post
}

func (r *repository) Close() error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}
	if r.rdb != nil {
		if err := r.rdb.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Ping(ctx context.Context) error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.PingContext(ctx); err != nil {
			return err
		}
	}
	if r.rdb != nil {
		if _, err := r.rdb.Ping(ctx).Result(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Init() error {
	//TODO implement me
	panic("implement me")
}
