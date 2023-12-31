package repository

import (
	"app/database"
	"app/model"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type userRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

// 实例化
func NewUserRepository(db *gorm.DB, rdb *database.RedisDB) UserRepository {
	return &userRepository{
		db:  db,
		rdb: rdb,
	}
}

func (u *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user *model.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	if err := u.SetCache(user); err != nil {
		logrus.Errorf("failed to set user: %v", err)
	}

	return user, nil
}

func (u *userRepository) GetUserByAuthID(authType, authID string) (*model.User, error) {
	authInfo := new(model.UserInfo)
	if err := u.db.Where("auth_type = ? and auth_id = ?", authType, authID).First(&authInfo).Error; err != nil {
		return nil, err
	}
	return u.GetUserByID(authInfo.UserID)
}

func (u *userRepository) GetUserByName(name string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetUserByMobile(number string) (*model.User, error) {
	var user model.User

	var result = u.db.Where("mobile = ?", number).Select("id").First(user)
	if result.RowsAffected != 0 {
		err := errors.New("手机号已存在")
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) List(pageSize int, pageNum int) (int, []interface{}) {
	var users []model.User
	userList := make([]interface{}, 0, len(users))

	if err := u.db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error; err != nil {
		return 0, nil
	}
	total := len(users)
	for _, user := range users {
		userItemMap := map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"password": user.Password,
			"email   ": user.Email,
			"avatar  ": user.Avatar,
		}
		userList = append(userList, userItemMap)
	}
	return total, userList

}

func (u *userRepository) Create(user *model.User) (*model.User, error) {
	userCreateField := []string{"name", "email", "password", "repassword", "mobile", "avatar"}
	if err := u.db.Select(userCreateField).Create(user).Error; err != nil {
		return nil, err
	}
	u.SetCache(user)
	return user, nil
}

func (u *userRepository) Update(user *model.User) (*model.User, error) {
	if err := u.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(&user).Error; err != nil {
		return nil, err
	}
	u.rdb.HDel(user.CacheKey(), strconv.Itoa(int(user.ID)))
	return user, nil
}

func (u *userRepository) Delete(id int) error {
	var user model.User
	if err := u.db.Delete(&user, id).Error; err != nil {
		return err
	}

	u.rdb.HDel(user.CacheKey(), strconv.Itoa(int(user.ID)))
	return nil
}

func (u *userRepository) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *userRepository) FindAll(userlist []model.User) ([]model.User, error) {
	if err := u.db.Find(&userlist).Error; err != nil {
		return nil, err
	}
	return userlist, nil
}

func (u *userRepository) SetCache(user *model.User) error {
	if user == nil {
		return nil
	}
	return u.rdb.HSet(user.CacheKey(), strconv.Itoa(int(user.ID)), user)
}
