package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strconv"
	"time"
)

/*
将Excel导入Mysql，内容：
	链式操作
	反射
	类型转换
*/

var DB *gorm.DB

const file = "/Users/jinhuaiwang/go/src/app/xlsx/test.xlsx"

type ExcelData interface {
	// 把excel中每行数据转换成map
	CreateMap(arr []string) map[string]interface{}
	ChangeTime(source string) time.Time
}

type ExcelStrcut struct {
	// 二维数组
	temp  [][]string
	Model interface{}
	Info  []map[string]string
}

// 确保定义的结构体字段和表中的字段一致，与定义的顺序无关。
type Temp struct {
	Uuid         uint64
	GoodName     string
	GoodMainImg  string
	GoodDescLink string
	CategoryName string
	TaobaokeLink string
	GoodPrice    float64
	TicketId     string
	TicketCount  uint64
	TicketLast   uint64
	StartTime    time.Time
	EndTime      time.Time
}

func (t *Temp) TableName() string {
	return "user"
}

// 读取Excel中数据 转换成二维数组
func (excel *ExcelStrcut) ReadExcel(file string) *ExcelStrcut {
	f, err := excelize.OpenFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rows)
	excel.temp = rows
	return excel
}

// 将二维数组中的每行转成对应的map
func (excel *ExcelStrcut) CreateMap() *ExcelStrcut {
	//利用反射得到字段名
	for _, v := range excel.temp {
		//将二维数组的每行转成对应切片类型的map
		var info = make(map[string]string)
		for i := 0; i < reflect.ValueOf(excel.Model).NumField(); i++ {
			obj := reflect.TypeOf(excel.Model).Field(i)
			//fmt.Printf("key:%s--val:%s\n", obj.Name, v[i])
			info[obj.Name] = v[i]
		}
		excel.Info = append(excel.Info, info)
	}
	return excel
}

// 时间做格式化
func (excel *ExcelStrcut) ChangeTime(source string) time.Time {
	times, err := time.Parse("2006-01-02", source)
	if err != nil {
		log.Fatalf("转换时间错误:%s", err)
	}
	return times
}

// 处理数据，写入mysql中
func (excel *ExcelStrcut) SaveDB(temp *Temp) *ExcelStrcut {
	//忽略标题行
	for i := 1; i < len(excel.Info); i++ {
		t := reflect.ValueOf(temp).Elem()
		for k, v := range excel.Info[i] {
			// 从map中读取出字段的值和类型
			//fmt.Println("key:%v---val:%v", t.FieldByName(k), t.FieldByName(k).Kind())
			switch t.FieldByName(k).Kind() {
			case reflect.String:
				//把map中的value写入到struct
				t.FieldByName(k).Set(reflect.ValueOf(v))
			case reflect.Float64:
				strToFloat64, err := strconv.ParseFloat(v, 64)
				if err != nil {
					log.Printf("string to float64 err：%v", err)
				}
				t.FieldByName(k).Set(reflect.ValueOf(strToFloat64))
			case reflect.Uint64:
				strToUint64, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					log.Printf("string to uint64 err：%v", err)
				}
				t.FieldByName(k).Set(reflect.ValueOf(strToUint64))
			case reflect.Struct:
				times, err := time.Parse("2006-01-02", v)
				if err != nil {
					log.Printf("string to time err：%v", err)
				}
				t.FieldByName(k).Set(reflect.ValueOf(times))
			default:
				fmt.Println("type err")
			}
		}
		//fmt.Println(temp)
		if err := DB.Create(&temp).Error; err != nil {
			log.Fatalf("save temp table err:%v", err)
		}
		fmt.Printf("导入临时表成功")
	}
	return excel
}

func init() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:123456@(192.168.10.168:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"))
	// DB.LogMode(true)//打印sql语句
	if err != nil {
		log.Fatalf("database connect is err:%s", err.Error())
	} else {
		log.Print("connect database is success")
	}
	DB.AutoMigrate(&Temp{})
}

func main() {
	e := ExcelStrcut{}
	temp := Temp{}
	e.Model = temp
	e.ReadExcel(file).CreateMap().SaveDB(&temp)
}
