package forms

import "app/model"

type CategoryForm struct {
	Name  string `form:"name" json:"name" binding:"required,min=3,max=20"`
	Desc  string `form:"desc" json:"desc" binding:"required,min=3,max=20"`
	Image string `form:"image" json:"image" `
}

func (c *CategoryForm) GetCategory() *model.Category {
	return &model.Category{
		Name:  c.Name,
		Desc:  c.Desc,
		Image: c.Image,
	}
}
