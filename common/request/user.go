package request

type Login struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	//Email    string `json:"email" form:"email" binding:"required"`
}
