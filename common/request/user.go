package request

type Login struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type Register struct {
	Name       string `json:"name" form:"name" binding:"required"`
	Password   string `json:"password" form:"password" binding:"required"`
	RePassword string `json:"repassword" form:"password" binding:"required"`
	Email      string `json:"email" form:"email" binding:"required"`
	Mobile     string `json:"mobile" form:"mobile" binding:"required,mobile"` //手机号码格式规范， 自定义validator
}
