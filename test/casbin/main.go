package main

import (
	"fmt"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
)

var DB *gorm.DB

func init() {
	dsn := "root:123456@tcp(192.168.10.168:3306)/testuser?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("connect DB error")
		panic(err)
	}
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second,   // 慢 SQL 阈值
	//		LogLevel:      logger.Silent, // Log level
	//		Colorful:      true,          // 禁用彩色打印
	//	},
	//)
	//
	//var err error
	//db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true,
	//	},
	//	Logger: newLogger,
	//})
	DB = db
}

func EnforcerTool() *casbin.Enforcer {
	adapter := gormadapter.NewAdapterByDB(DB)
	enforcer := casbin.NewEnforcer("/Users/jinhuaiwang/go/src/app/config/rbac.conf", adapter)
	enforcer.LoadPolicy()
	return enforcer
}

func Interceptor(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Request.URL.Path
		act := c.Request.Method
		sub := "admin"

		if ok := enforcer.AddPolicy(sub, obj, act); !ok {
			fmt.Println("没有通过权限")
			c.Abort()
		} else {
			fmt.Println("通过权限")
			c.Next()
		}
	}
}

func GetUser(c *gin.Context) {
	message := "成功"
	code := 200
	data := "data"
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
		"result":  "true",
	})
}

func main() {
	r := gin.New()
	r.Use(Interceptor(EnforcerTool()))
	r.GET("/api/v1/user", GetUser)
	r.Run()
}
