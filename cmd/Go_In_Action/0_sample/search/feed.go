package search

import (
	"encoding/json"
	"os"
)

//Golang 中的相对路径，均为相对GoPath的路径
const dataFile = "../go_learning/cmd/Go_In_Action/0_sample/data/data.json"

// 定义一个struct 类型
type Feed struct {
	// 每个字段声明后的``被称为标记tag
	// 用于将结构体的字段与JSON文档中的字段对应起来
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

// 读取数据源文件data.json 并解析文件，存入到Feed 切片中
func RetrieveFeeds() ([]*Feed, error) {
	// 打开文件，将返回指向File 的指针
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	// defer 关键字标识的方法会在函数返回时再执行
	// 类似java 的finally
	defer file.Close()

	// 定义一个切片，并将json 解码到切片中
	var feeds []*Feed
	// NewDecoder方法接受一个io.Reader 的参数
	// File 类型实现了该interface 的Read 方法，因此可以传入
	err = json.NewDecoder(file).Decode(&feeds)

	return feeds, err
}
