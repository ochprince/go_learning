package search

import (
	"log"
	"sync"
)

// 创建包级变量且使用make将其初始化为
// key为string类型 value为Matcher类型的map
// 注意该变量是小写字母开头，意味着包外不可访问
var matchers = make(map[string]Matcher)

func Register(feedType string, matcher Matcher) {
	if _, exist := matchers[feedType]; exist {
		log.Fatalln(feedType, "Matcher already registered")
	}

	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}

// 该方法公开，接受一个string类型的参数，没有返回值
func Run(searchTerm string) {
	// 获取数据源feeds，是一个切片类型
	// 这里使用了简化变量声明运算符 :=
	// 一般来说，声明初始值为零值的变量使用var，非零值初始化变量使用 :=
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个无缓冲的通道，接受匹配后结果
	// 关键字chan 和map 一样，只能由make 进行初始化
	// 该段代码意味着通道中仅可传递指向Result 对象的指针
	results := make(chan *Result)

	// 创建一个计数信号量waitGroup用于统计和等待goroutine
	// 防止程序在全部执行完成前终止
	var waitGroup sync.WaitGroup

	// 设置waitGroup的大小
	waitGroup.Add(len(feeds))

	// 遍历切片feeds
	// range 可以用于遍历数组、字符串、切片、映射和通道
	// range 会返回index和元素的副本
	for _, feed := range feeds {
		// 获取一个合适的匹配器用于查找
		// map 返回一个值时，若不存在则返回零值
		// map 返回两个值时，第二个为bool
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		// 启动一个goroutine
		// 类似java的lambda表达式
		// 这里的第二个参数为指向Feed 类型的指针
		// Golang中，所有方法的参数都是值传递
		go func(matcher Matcher, feed *Feed) {
			// 下面的searchTerm 和results，以及waitGroup 都是应用了闭包，直接访问外层函数的变量本身
			// matcher 和feed 没有直接访问外层函数的变量，这是因为其值会变化，导致每个goroutine 使用的变量都会改变
			Match(matcher, feed, searchTerm, results)
			//完成工作后，waitGroup递减
			waitGroup.Done()
		}(matcher, feed)
	}

	// 启动一个goroutine 来监控所有工作是否完成
	go func() {
		// 等候所有任务完成（waitGroup递减到0时）
		waitGroup.Wait()

		// 关闭通道
		close(results)
	}()

	// 实时显示结果，并在results 被关闭之后返回
	Display(results)
}
