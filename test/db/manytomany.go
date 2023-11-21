package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*

多对多关系

多对多关系，需要用第三张表存储两张表的关系

*/

// 多对多就必须得加上many2many的tag。article_tags是用来指定第三张表
type Tag struct {
	ID       uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	Name     string    `json:"name" gorm:"size:256;not null;unique"`
	Articles []Article `gorm:"many2many:articles_tags;"`
}

type Article struct {
	ID    uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Title string `json:"title" gorm:"size:256;not null;unique"`
	Tags  []Tag  `gorm:"many2many:articles_tags;"`
}

func main() {
	dsn := "root:123456@tcp(192.168.10.168:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Article{}, &Tag{})

	//多对多添加
	//添加文章,并创建标签
	db.Create(&Article{
		Title: "python base",
		Tags: []Tag{
			{Name: "python"},
			{Name: "backend"},
		},
	})

	db.Create(&Article{
		Title: "golang study",
		Tags: []Tag{
			{Name: "golang"},
		},
	})

	// 添加文章选择标签
	var tags []Tag
	db.Find(&tags, "name in ?", []string{"python"})
	fmt.Println(tags)

	db.AutoMigrate(&Article{}, &Tag{})
	db.Create(&Article{
		Title: "python study",
		Tags:  tags,
	})

	//多对对查询

	// 查询文章, 显示文章的标签列表
	var article Article
	db.Preload("Tags").Take(&article, 1)
	fmt.Println(article)

	// 查询标签, 显示文章的列表
	var tag Tag
	db.Preload("Articles").Take(&tag, 2)
	fmt.Println(tag)
}

// check databases
/*
MySQL [(none)]> use test;
MySQL [test]> select * from articles;
+----+-------------+
| id | title       |
+----+-------------+
|  1 | python base |
+----+-------------+
1 row in set (0.01 sec)

MySQL [test]> select * from tags;
+----+---------+
| id | name    |
+----+---------+
|  2 | backend |
|  1 | python  |
+----+---------+
2 rows in set (0.00 sec)

MySQL [test]> select * from articles_tags;
+--------+------------+
| tag_id | article_id |
+--------+------------+
|      1 |          1 |
|      2 |          1 |
+--------+------------+
2 rows in set (0.00 sec)

*/
