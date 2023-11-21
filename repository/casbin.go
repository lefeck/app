package repository

import (
	"app/global"
	"app/model"
	"errors"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
	gormadapter "github.com/casbin/gorm-adapter"
	"strings"
)

type CasbinRepository struct {
}

// 要实现的接口的方法
type ICasbinRepository interface {
	Create(c *model.Casbin) (*model.Casbin, error)
	Delete(c *model.Casbin) error
	Update(c *model.Casbin, value interface{}) error
	List(c *model.Casbin) [][]string
	//CasbinCreate(roleId string, path, method string) error
	CasbinDelete(roleId string, path, method string) error
	CasbinList(roleID string) [][]string
	Migrate() error
}

// 实例化
func NewCasbinRepository() ICasbinRepository {
	return &CasbinRepository{}
}

func (c *CasbinRepository) Create(cs *model.Casbin) (*model.Casbin, error) {
	e := Casbin()
	if ok := e.AddPolicy(cs.RoleId, cs.Path, cs.Method); ok == false {
		return nil, errors.New("存在相同的API，添加失败")
	}
	return cs, nil
}

func (c *CasbinRepository) Delete(cs *model.Casbin) error {
	e := Casbin()
	if ok := e.DeletePermission(cs.RoleId, cs.Path, cs.Method); ok == false {
		return errors.New("不存在相同的API，删除失败")
	}
	return nil
}

func (c *CasbinRepository) Update(cs *model.Casbin, value interface{}) error {
	if err := global.DB.Model(c).Where("v1 = ? AND v2 = ?", cs.Path, cs.Method).Update(value).Error; err != nil {
		return err
	}
	return nil
}

func (c *CasbinRepository) List(cs *model.Casbin) [][]string {
	e := Casbin()
	policy := e.GetFilteredPolicy(0, cs.RoleId)
	return policy
}

//var DB *gorm.DB

// @function: Casbin
// @description: 持久化到数据库  引入自定义规则
// @return: *casbin.Enforcer
func Casbin() *casbin.Enforcer {
	adapter := gormadapter.NewAdapterByDB(global.DB)
	enforcer := casbin.NewEnforcer("./config/rbac.conf", adapter)
	enforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
	return enforcer
}

// @function: ParamsMatch
// @description: 自定义规则函数
// @param: fullNameKey1 string, key2 string
// @return: bool
func ParamsMatch(fullNameKey, key string) bool {
	key1 := strings.Split(fullNameKey, "?")[0]
	fmt.Println(key1)
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key)
}

// @function: ParamsMatchFunc
// @description: 自定义规则函数
// @param: args ...interface{}
// @return: interface{}, error
func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return ParamsMatch(name1, name2), nil
}

//func (c *CasbinRepository) CasbinCreate(roleId string, path, method string) error {
//	cm := model.Casbin{
//		PType:  "p",
//		RoleId: roleId,
//		Path:   path,
//		Method: method,
//	}
//	return c.Create(&cm)
//}

func (c *CasbinRepository) CasbinDelete(roleId string, path, method string) error {
	cm := model.Casbin{
		PType:  "p",
		RoleId: roleId,
		Path:   path,
		Method: method,
	}
	return c.Delete(&cm)
}

func (c *CasbinRepository) CasbinList(roleID string) [][]string {
	cm := model.Casbin{RoleId: roleID}
	return c.List(&cm)
}

func (u *CasbinRepository) Migrate() error {
	global.DB.AutoMigrate(&model.Casbin{})
	return nil
}
