package matchers

import "../search"

func init() {
	var matcher defaultMatcher
	search.Register("default", matcher)
}

// 定义一个struct 类型，用于实现Matcher
type defaultMatcher struct {
}

// 实现Matcher 的Search 行为
// 该方法带有接收者defaultMatcher
// 意味着可以使用 defaultMatcher 类型的值或者指向这个类型值的指针来调用 Search 方法
// 无论var dm defaultMatcher 还是dm := new(defaultMatcher)，都可以这样使用：dm.Search(feed, "test")
// 但如果Search 方法是func (dm *defaultMatcher)，则不能使用Matcher 接口的值来调用Search，但可以用Matcher 接口的指针来调用
func (dm defaultMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
	return nil, nil
}
