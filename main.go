package main

import (
	"./logger"
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"strconv"
)

const (
	requestHead           = "http://jandan.net/ooxx"
	requestPageHeadPrefix = "http://jandan.net/ooxx/page-"
	requestPageHeadSuffix = "#comments"
	path                  = "./meizhi"
	concurrent            = 3 //并发
)

var (
	maxPageNum int = 0
)

func getPageNum(html string) int {
	reg, _ := regexp.Compile(`<span class="current-comment-page">\[\d*\]</span>`)
	reg1, _ := regexp.Compile(`\d+`)
	s := reg.FindString(html)
	page := reg1.FindString(s)
	num, _ := strconv.Atoi(page)
	return num

}

func getPicUrl() string {
	//TODO
	return ""
}

func getHtml(url string) (error, string, error) {
	response, err := http.Get(url)
	defer response.Body.Close()
	html, err1 := ioutil.ReadAll(response.Body);
	return err, string(html), err1
}

func download() {

}

func main() {
	_, html, _ := getHtml(requestHead)
	maxPageNum = getPageNum(html)
	download()
	fmt.Print("hello")
	logger.Info("hello")
}
