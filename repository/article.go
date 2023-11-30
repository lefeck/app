package repository

import (
	"app/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

// 通过id获取文章
func (a *articleRepository) GetArticleByID(id uint) (*model.Article, error) {
	article := new(model.Article)
	//文章关联了,creator, tag, category, comments, comment.user(评论的用户)
	if err := a.db.Preload(model.CreatorAssociation).Preload(model.TagsAssociation).Preload(model.CategoriesAssociation).Preload("Comments.User").Preload(model.CommentsAssociation).Find(article).Error; err != nil {
		return nil, err
	}
	like, err := a.CountLike(id)
	if err != nil {
		return nil, err
	}
	article.Likes = uint(like)
	return article, nil
}

// 统计文章的点赞数
func (a *articleRepository) CountLike(id uint) (int64, error) {
	var count int64
	if err := a.db.Model(&model.Like{}).Where("article_id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 通过名字获取文章
func (a *articleRepository) GetArticleByName(name string) (*model.Article, error) {
	article := new(model.Article)
	if err := a.db.Where("name = ?", name).First(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

// 获取所有的文章
func (a *articleRepository) List() ([]model.Article, error) {
	articles := make([]model.Article, 0)
	//db.Omit 在create，update和query 忽略哪些字段
	/*
		order 从数据库检索记录时指定顺序
		db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
		// 表示按照创建时间顺序来获取输出
	*/

	if err := a.db.Omit("content").Preload(model.CreatorAssociation).Preload(model.TagsAssociation).Preload(model.CategoriesAssociation).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&articles).Error; err != nil {
		return nil, err
	}

	ids := make([]uint, len(articles))
	for i, article := range articles {
		ids[i] = article.ID
	}

	type result struct {
		ID    uint
		Likes uint
	}

	results := []result{}
	if err := a.db.Model(&model.Like{}).Select("article_id as id, count(likes.article_id) as likes").Where("article_id in ?", ids).Group("article_id").Scan(&results).Error; err != nil {
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

// 创建文章
func (a *articleRepository) Create(user *model.User, article *model.Article) (*model.Article, error) {
	article.CreatorID = user.ID
	article.Creator = *user
	err := a.db.Create(article).Error
	return article, err
}

// 更新文章
func (a *articleRepository) Update(article *model.Article) (*model.Article, error) {
	err := a.db.Model(article).Omit("views", "creator_id").Updates(article).Error
	return article, err
}

// 删除文章
func (a *articleRepository) Delete(id uint) error {
	article := &model.Article{}
	return a.db.Delete(article, id).Error
}

// 文章阅读数
func (a articleRepository) IncView(id uint) error {
	article := &model.Article{ID: id}
	return a.db.Model(article).UpdateColumn("views", gorm.Expr("views + 1")).Error
}

// 自动创建表结构到db
func (a *articleRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Article{})
}
