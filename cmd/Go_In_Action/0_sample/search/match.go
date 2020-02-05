package search

import (
	"fmt"
	"log"
)

// 保存match 到的结果
type Result struct {
	Field   string
	Content string
}

// 定义一个interface 类型
// 规范matcher 的行为
type Matcher interface {
	// 接收一个指向Feed的指针和一个查询项
	// 返回指向结果的指针的切片和一个可能的错误
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

// 使用Matcher 执行搜索，并将结果写入channel
// chan<- 代表send-only，意味着只能向其中写数据，不能从中读数据
func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	// 这里的matcher 即为Matcher 接口的值
	searchResults, err := matcher.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}

	for _, result := range searchResults {
		// 结果写入通道
		results <- result
	}
}

// 从每个goroutine 的结果通道中获取值并输出到终端
func Display(results chan *Result) {
	// 通道会一直被阻塞，知道有结果写入
	// 一旦通道被关闭，for 循环会终止
	for result := range results {
		fmt.Printf("%s:\n%s\n\n", result.Field, result.Content)
	}
}
