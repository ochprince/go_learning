package matchers

import (
	"../search"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// 批量定义rss 数据结构对应的结构体
// 用于存储rss 数据
type (
	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}

	channel struct {
		XMLName       xml.Name
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		Description   string `xml:"description"`
		Language      string `xml:"language"`
		Copyright     string `xml:"copyright"`
		Generator     string `xml:"generator"`
		LastBuildDate string `xml:"lastBuildDate"`
		Image         image  `xml:"image"`
		Item          []item `xml:"item"`
	}

	image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	item struct {
		XMLName        xml.Name `xml:"item"`
		Title          string   `xml:"title"`
		Description    string   `xml:"description"`
		PubDate        string   `xml:"pubDate"`
		Link           string   `xml:"link"`
		GUID           string   `xml:"guid"`
		ContentEncoded string   `xml:"content:encoded"`
		DcCreator      string   `xml:"dc:creator"`
	}
)

func init() {
	var matcher rssMatcher
	search.Register("rss", matcher)
}

type rssMatcher struct {
}

func (rm rssMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
	var results []*search.Result

	log.Printf("Search Feed Type[%s] Site[%s] For Uri[%s] \n", feed.Type, feed.Name, feed.URI)

	document, err := rm.retrieve(feed)
	if err != nil {
		return nil, err
	}

	for _, channelItem := range document.Channel.Item {
		rm.matchString("Title", channelItem.Title, searchTerm, &results)
		rm.matchString("Description", channelItem.Description, searchTerm, &results)
	}
	return results, nil
}

// 这里是在原著基础上改动的部分
// 最早results 参数的类型写成了 []*search.Result
// 执行完后发现上一级results 返回的nil
// 原因是Golang 参数只有值传递，这一级对results 切片的任何处理都不会影响上一级的results
func (rm rssMatcher) matchString(typeOfItem string, contentOfItem string, searchTerm string, results *[]*search.Result) {
	// 使用正则匹配查询item 中的标题中是否包含搜索项
	matched, err := regexp.MatchString(searchTerm, contentOfItem)
	// 如果找到匹配的项，将其作为结果保存
	if err == nil && matched {
		// 这里使用了append 方法，将元素添加到切片中，返回了新切片
		*results = append(*results, &search.Result{
			Field:   typeOfItem,
			Content: contentOfItem,
		})
	}
}

// retrieve 发送Http Get 请求获取feed 对应rss 数据源并解码
func (rm rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
	if feed.URI == "" {
		// 这里使用了errors.New 方法创建了一个错误对象
		return nil, errors.New("No rss feed URI provided")
	}

	// 通过http Get 获取远程资源
	resp, err := http.Get(feed.URI)
	if err != nil {
		return nil, err
	}
	// 函数返回时关闭响应连接
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// 这里使用fmt.Errorf 方法构建了一个错误对象
		return nil, fmt.Errorf("Http response error with code %d\n", resp.StatusCode)
	}

	// 和json.NewDecoder 类似，将响应内容解码到rssDocument 中
	var document rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&document)
	return &document, err
}
