package model

// 文章: 文章ID 用户ID 分类ID 文章标题 封面 文章内容 创建时间 修改时间 阅读数 评论数
//type Article struct {
//	ID             uint       `gorm:"primary_key" json:"aid"`
//	CreatedAt      time.Time  `json:"-"`
//	UpdatedAt      time.Time  `json:"-"`
//	DeletedAt      *time.Time `sql:"index" json:"-"`
//	UserID         uint       `gorm:"index;not null" json:"-"`
//	CategoryItemID uint       `gorm:"index;not null" json:"cid"`
//	Title          string     `gorm:"type:varchar(100);not null" json:"title"`
//	Content        string     `gorm:"type:text;not null;" json:"content"`
//	Cover          string     `gorm:"not null" json:"cover"`
//	CreateTime     time.Time  `gorm:"not null" json:"create_time"`
//	Views          int        `gorm:"default:0" json:"views"`
//
//	//	Comments int `gorm:"default:0"` //暂时先不做评论
//}

const (
	CreatorAssociation    = "Creator"
	ArticlesAssociation   = "Articles"
	TagsAssociation       = "Tags"
	CategoriesAssociation = "Categories"
	CommentsAssociation   = "Comments"
)

type Article struct {
	BaseModel
	ID         uint       `json:"aid" gorm:"autoIncrement;primaryKey"`
	Title      string     `json:"title" gorm:"type:varchar(100);not null;unique"`
	Content    string     `json:"content" gorm:"type:text;not null"`
	Cover      string     `json:"cover" gorm:"not null"`
	CreatorID  uint       `json:"creatorId"`
	Creator    User       `json:"creator" gorm:"foreignKey:CreatorID"`
	Tags       []Tag      `json:"tags"  gorm:"many2many:tag_articles"`
	Categories []Category `json:"categories" gorm:"many2many:category_articles"`
	//评论
	Comments []Comment `json:"comments"`
	// 文章阅读量
	Views uint `json:"views" gorm:"type:uint"`
	// 点赞数
	Likes     uint `json:"likes" gorm:"-"`
	UserLiked bool `json:"userLiked" gorm:"-"`
	Origin    int  `json:"origin" gorm:"not null"` //是否原创 1原创 0转载
	State     int  `json:"-" gorm:"default:0"`     //0正常发布 2并未发布(草稿箱)
}

func (*Article) TableName() string {
	return "article"
}
