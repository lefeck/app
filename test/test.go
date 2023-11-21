package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	//key1 := "/api/v1/users?pagesize=10&pagenum=1"
	//key := ""
	//ParamsMatch(key1, key)
	//key3 := "/swagger/*any"
	//KeyMatch(key3, key)
	//test()

	Get(&contexts{
		[]param{{Key: "name", Name: "tom"}},
	})
}

// 获取key对应的值
type param struct {
	Key  string
	Name string
}

type params []param

type contexts struct {
	Params params
}

func (c *contexts) Parames(key string) string {
	return c.Params.ByName(key)
}

func (p params) ByName(name string) (value string) {
	value, _ = p.Get(name)
	return
}

func (p params) Get(name string) (string, bool) {
	for _, entry := range p {
		if entry.Key == name {
			return entry.Name, true
		}
	}
	return "", false
}

func Get(c *contexts) {
	name := c.Parames("name")
	fmt.Println(name)
}

func ParamsMatch(fullNameKey, key string) bool {
	key1 := strings.Split(fullNameKey, "?")[0]
	fmt.Println(key1)
	// 剥离路径后再使用casbin的keyMatch2
	return KeyMatch(key1, key)
}

func KeyMatch(key1 string, key2 string) bool {
	i := strings.Index(key2, "*")
	fmt.Println(i)
	if i == -1 {
		return key1 == key2
	}
	fmt.Println(key1)

	if len(key1) > i {
		return key1[:i] == key2[:i]
	}
	return key1 == key2[:i]
}

func test() {
	str := regexp.MustCompile(`^\/(.*)\/(.*)$`)
	params := str.FindStringSubmatch("/common/vs_test1")
	fmt.Println(params)
	value := params[len(params)-1]
	fmt.Println(value)
	//for _, value := range params {
	//	fmt.Println(value)
	//}
}

func DeleteExtraSpace(s string) string {
	s1 := strings.Replace(s, "  ", " ", -1)
	regstr := "\\s{2,}"
	reg, _ := regexp.Compile(regstr)
	s2 := make([]byte, len(s1))
	copy(s2, s1)
	spec_index := reg.FindStringIndex(string((s2)))
	//fmt.Println(spec_index)
	for len(spec_index) > 0 {
		fmt.Println(s2)
		s2 = append(s2[:spec_index[0]+1], s2[spec_index[1]:]...)
		fmt.Println(s2)
		spec_index = reg.FindStringIndex(string(s2))
		fmt.Println(spec_index)
	}
	return string(s2)
}
