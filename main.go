package main

import (
	"./logger"
	"net/http"
	"io/ioutil"
	"regexp"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"log"
	"fmt"
	"os"
	"strings"
	"./base64kit"
)

const (
	requestHead           string = "http://jandan.net/ooxx"
	requestPageHeadPrefix string = "http://jandan.net/ooxx/page-"
	requestPageHeadSuffix string = "#comments"
	path                  string = "./meizhi"
	concurrent            int    = 3 //并发
)

var (
	maxPageNum        int    = 0
	pageNumSpanPatten string = `<span class="current-comment-page">\[\d*\]</span>`
	pageNumPatten     string = `\d+`

	c chan int
)

func getAllPageNum(html string) int {
	reg := getRegexpPatten(pageNumSpanPatten)
	reg1 := getRegexpPatten(pageNumPatten)
	s := reg.FindString(html)
	page := reg1.FindString(s)
	num, _ := strconv.Atoi(page)
	return num

}

func setAllPicHashFromPage(doc *goquery.Document) []string {
	pagePicHash := make([]string, 100)
	doc.Find(".img-hash").Each(func(i int, selection *goquery.Selection) {
		src := selection.Text()
		pagePicHash[i] = src
	})
	return pagePicHash
}

func getDocumentBody(url string) *goquery.Document {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		logger.Error("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc

}

func getHtml(url string) (error, string, error) {
	response, err := http.Get(url)
	defer response.Body.Close()
	html, err1 := ioutil.ReadAll(response.Body);
	return err, string(html), err1
}

func download(num int) {
	logger.Info("page %d", num)
	var url = requestPageHeadPrefix + strconv.Itoa(num) + requestPageHeadSuffix
	doc := getDocumentBody(url)
	pagePicHash := setAllPicHashFromPage(doc)
	if len(pagePicHash) > 0 {
		for _, hash := range pagePicHash {
			if hash != "" {
				url, _ := base64kit.Base64Decode(hash)
				down(url)
			} else {
				continue
			}
		}
	}

}
func down(s string) {
	logger.Info("http:" + s)
	imgResponse, err := http.Get("http:" + s)
	defer imgResponse.Body.Close()
	if err != nil {
		logger.Error(err.Error())
	}
	path := path + "/"
	imgByte, _ := ioutil.ReadAll(imgResponse.Body)
	pInfo, pErr := os.Stat(path)
	if pErr != nil || pInfo.IsDir() == false {
		errDir := os.Mkdir(path, os.ModePerm)
		if errDir != nil {
			fmt.Println(errDir)
			os.Exit(-1)
		}
	}
	fn := path + s[strings.LastIndex(s, "/")+1:len(s)]
	_, fErr := os.Stat(fn)
	var fh *os.File
	if fErr != nil {
		fh, _ = os.Create(fn)
	} else {
		fh, _ = os.Open(fn)
	}
	defer fh.Close()
	fh.Write(imgByte)

}

func getRegexpPatten(patten string) *regexp.Regexp {
	reg, _ := regexp.Compile(patten);
	return reg
}

func main() {
	c = make(chan int)

	//_, home, _ := getHtml(requestHead)
	//maxPageNum = getAllPageNum(home)

	maxPageNum = 2
	for i := maxPageNum; i > 0; i-- {
		download(i)
	}
}
