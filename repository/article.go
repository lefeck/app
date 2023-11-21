package repository

import (
	"app/database"
	"app/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type articleRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func NewArticleRepository(db *gorm.DB, rdb *database.RedisDB) ArticleRepository {
	return &articleRepository{
		db:  db,
		rdb: rdb,
	}
}

func (p *articleRepository) List() ([]model.Article, error) {
	articles := make([]model.Article, 0)
	if err := p.db.Omit("content").Preload(model.CreatorAssociation).Preload(model.TagAssociation).Preload(model.CategoriesAssociation).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&articles).Error; err != nil {
		return nil, err
	}

	ids := make([]uint, len(articles))
	for i, a := range articles {
		ids[i] = a.ID
	}
	type result struct {
		ID    uint
		Likes uint
	}

	results := []result{}
	if err := p.db.Model(&model.Like{}).Select("post_id as id, count(likes.post_id) as likes").Where("post_id in ?", ids).Group("post_id").Scan(&results).Error; err != nil {
		return nil, err
	}

	resMap := make(map[uint]uint, len(results))
	for _, r := range results {
		resMap[r.ID] = r.Likes
	}

	for i := range articles {
		articles[i].Likes = resMap[articles[i].ID]
	}
	return articles, nil
}

func (p *articleRepository) Create(user *model.User, article *model.Article) (*model.Article, error) {
	article.CreatorID = user.ID
	article.Creator = *user
	err := p.db.Create(article).Error
	return article, err
}

func (p *articleRepository) GetArticleByID(id uint) (*model.Article, error) {
	post := new(model.Article)
	if err := p.db.Preload(model.CreatorAssociation).Preload(model.TagAssociation).Preload(model.CategoriesAssociation).Preload("Comments.User").Preload(model.CommentsAssociation).
		First(post, id).Error; err != nil {
		return nil, err
	}
	//获取帖子的点赞数
	likes, err := p.CountLike(id)
	if err != nil {
		return nil, err
	}
	post.Likes = uint(likes)
	return post, nil
}

func (p *articleRepository) CountLike(id uint) (int64, error) {
	var count int64
	if err := p.db.Model(&model.Like{}).Where("pos_id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 获取用户创建的文章
func (p *articleRepository) GetArticleByName(name string) (*model.Article, error) {
	article := new(model.Article)
	err := p.db.Where("name = ?", name).Find(article).Error
	return article, err
}

func (p *articleRepository) Update(article *model.Article) (*model.Article, error) {
	err := p.db.Model(article).Omit("views", "creator_id").Updates(article).Error
	return article, err
}

func (p *articleRepository) Delete(id uint) error {
	return p.db.Delete(&model.Article{}, id).Error
}

// 文章阅读量增加方法
func (p *articleRepository) IncView(id uint) error {
	article := model.Article{ID: id}
	return p.db.Model(article).UpdateColumn("views", gorm.Expr("view + 1")).Error
}

func (p *articleRepository) Migrate() error {
	return p.db.AutoMigrate(&model.Article{}, &model.Like{}, &model.Category{}, &model.Tag{}, &model.Comment{})
}
