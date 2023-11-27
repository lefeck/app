package service

import (
	"app/common/request"
	"app/model"
	"app/repository"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (u *userService) List(pageSize int, pageNum int) (int, []interface{}) {
	return u.userRepository.List(pageSize, pageNum)
}

func (u *userService) Create(user *model.User) (*model.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(password)

	//if user == nil {
	//	return nil, errors.New("user is empty")
	//}
	//if user.Name == "" {
	//	return nil, errors.New("user name is empty")
	//}
	//if len(user.Password) < MinPasswordLength {
	//	return nil, fmt.Errorf("password length must great than %d", MinPasswordLength)
	//}
	//if user.Email == "" {
	//	user.Email = fmt.Sprintf("%s@163.com", user.Name)
	//}
	return u.userRepository.Create(user)
}

func (u *userService) Get(id string) (*model.User, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return u.userRepository.GetUserByID(uint(uid))
}

const (
	MinPasswordLength = 6
)

func (u *userService) Validate(user *model.User) error {
	if user == nil {
		return errors.New("user is empty")
	}
	if user.Name == "" {
		return errors.New("user name is empty")
	}
	if len(user.Password) < MinPasswordLength {
		return fmt.Errorf("password length must great than %d", MinPasswordLength)
	}
	return nil
}

func (u *userService) Update(user *model.User) (*model.User, error) {
	if len(user.Password) > 0 {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(password)
	}
	return u.userRepository.Update(user)
}

func (u *userService) Delete(id string) error {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return u.userRepository.Delete(uid)
}

func (u *userService) FindAll(userlist []model.User) ([]model.User, error) {
	return u.userRepository.FindAll(userlist)
}

func (u *userService) Login(param request.Login) (*model.User, error) {
	if u == nil || param.Name == "" || param.Password == "" {
		return nil, fmt.Errorf("name or password is empty")
	}

	user, err := u.userRepository.GetUserByName(param.Name)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password)); err != nil {
		return nil, err
	}

	return user, nil
}

// Register 注册
func (u *userService) Register(params request.Register) (user *model.User, err error) {

	if params.Password != params.RePassword {
		return nil, fmt.Errorf("Password do not match")
	}

	users := &model.User{Name: params.Name, Mobile: params.Mobile, Password: params.Password, RePassword: params.RePassword, Email: params.Email}
	user, err = u.userRepository.Create(users)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) Export(data *[]model.User, headerName []string, filename string, c *gin.Context) error {
	sheetName := "Sheet1"
	sheetWords := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U",
		"V", "W", "X", "Y", "Z",
	}
	file := excelize.NewFile()
	//设置表头
	for k, v := range headerName {
		file.SetCellValue(sheetName, sheetWords[k]+"1", v)
	}
	for k, v := range *data {
		k = k + 2
		// 向单元格中设置值
		file.SetCellValue(sheetName, "A"+strconv.Itoa(k), v.ID)
		file.SetCellValue(sheetName, "B"+strconv.Itoa(k), v.Name)
		file.SetCellValue(sheetName, "C"+strconv.Itoa(k), v.Password)
		file.SetCellValue(sheetName, "D"+strconv.Itoa(k), v.Email)
		file.SetCellValue(sheetName, "E"+strconv.Itoa(k), v.RePassword)
		file.SetCellValue(sheetName, "F"+strconv.Itoa(k), v.Avatar)
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	if err := file.Write(c.Writer); err != nil {
		return err
	}
	//// 将电子表格进行保存
	if err := file.SaveAs(filename); err != nil {
		return err
	}
	return nil
}
