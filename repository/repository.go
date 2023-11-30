package repository

import (
	"app/database"
	"context"
	"gorm.io/gorm"
)

type repository struct {
	user     UserRepository
	article  ArticleRepository
	category CategoryRepository
	comment  CommentRepository
	tag      TagRepository
	like     LikeRepository
	db       *gorm.DB
	rdb      *database.RedisDB
	migrants []Migranter
}

func NewRepository(db *gorm.DB, rdb *database.RedisDB) Repository {
	r := &repository{
		user:     NewUserRepository(db, rdb),
		article:  NewArticleRepository(db),
		tag:      NewTagRepository(db),
		comment:  NewCommentRepository(db),
		like:     NewLikeRepository(db),
		category: NewCategoryRepository(db),
		db:       db,
		rdb:      rdb,
	}
	r.migrants = getMigrants(
		r.user,
		r.article,
		r.article,
		r.like,
		r.category,
		r.tag)
	return r
}

func getMigrants(objs ...interface{}) []Migranter {
	var migrants []Migranter
	for _, obj := range objs {
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

func (r *repository) Migrate() error {
	return r.Migrant()
}

func (r *repository) User() UserRepository {
	return r.user
}

func (r *repository) Article() ArticleRepository {
	return r.article
}

func (r *repository) Comment() CommentRepository {
	return r.comment
}

func (r *repository) Category() CategoryRepository {
	return r.category
}

func (r *repository) Like() LikeRepository {
	return r.like
}

func (r *repository) Tag() TagRepository {
	return r.tag
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
