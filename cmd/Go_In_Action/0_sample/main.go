//包和文件夹名可以不同，但最好相同
package main

import (
	//下划线代表导入但不使用其中的标识符，仅执行init
	_ "../0_sample/matchers"
	"../0_sample/search"
	"log"
	"os"
)

// 在main之前执行
func init() {
	// 设置log的输出方式为标准输出
	log.SetOutput(os.Stdout)
}

// main包中有main函数，则会生成可执行文件
func main() {
	//开始搜索，查询包含president的段落
	search.Run("president")
}
