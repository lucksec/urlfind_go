package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	//获取参数
	var url string
	flag.StringVar(&url, "u", "", "use:main -u url")
	flag.Parse()
	if len(url) <= 0 {
		fmt.Println("urlfind \n\t by:luckone\n建议配合httpx使用   \nuse:\n\tmain -u url ｜ httpx")
	}
	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// 解析 HTML 响应并提取所有的链接
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // 忽略无效的 URL
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)

	// 打印所有提取到的链接
	for i := 0; i < len(links); i++ {
		fmt.Println(links[i])
	}

}

func resHead(url string) {
	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// 获取响应头内容
	contentType := resp.Header.Get("server")
	fmt.Println("server:" + contentType)
}

// forEachNode 遍历 n 及其子节点
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
