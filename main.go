package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Println("Crawling...")
	Crawl("output")
	fmt.Println("Finish...")
}

// Crawl all the url content and output the file
func Crawl(outputFile string) {
	var text []string
	for _, url := range GetURL("url") {
		// for i, s := range AllText(url) {
		// 	fmt.Println(i, " : ", s)
		// }
		text = append(text, AllText(url)...)
	}
	text = RemoveC(text)
	text = ReplaceC(text)
	Write(outputFile, RemoveEmpty(RemoveC(text)))
}
func RemoveC(text []string) []string {
	var ret []string

	for _, v := range text {
		if v == "-" || v == ">" || v == "｜" || v == "(" || v == ")" || v == "。" || v == "|" || v == "︽" || v == "/" || v == "[" || v == "]" {
			continue
		}
		v = strings.TrimSpace(v)

		ret = append(ret, v)
	}
	return ret
}

// replace content by replacer
func ReplaceC(text []string) []string {
	var ret []string
	var r = []string{",", "\n", "。", "\n", "，", "\n", "∘", "\n", "&nbsp;", ""}
	replacer := strings.NewReplacer(r...)
	for _, v := range text {
		v = strings.TrimSpace(v)
		ret = append(ret, replacer.Replace(v))
	}
	return ret
}

//GetURL  get the urls from file
func GetURL(fileName string) []string {
	content := Read(fileName)
	var urls []string
	for _, s := range RemoveEmpty(content) {
		urls = append(urls, "https"+strings.Split(s, "https")[1])

	}
	return urls
}

// AllText crawl the web content by tags
func AllText(url string) []string {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var textList []string
	// textList = append(textList, url)
	for _, tag := range AllTag() {
		if tag == "script" || tag == "noscript" || tag == "style" || tag == "body" || tag == "head" || tag == "html" || tag == "div" || tag == "template" || tag == "section" || tag == "main" || tag == "nav" || tag == "header" {
			continue
		}
		// debug
		textList = append(textList, tag)
		doc.Find(tag).Each(func(i int, s *goquery.Selection) {
			var contents []string
			// debug
			contents = append(contents, tag)
			s.Text()
			contents = append(contents, s.Text())
			textList = append(textList, RemoveEmpty(contents)...)
		})
	}
	return RemoveEmpty(textList)
}

// AllTag get all the html tags
func AllTag() []string {
	// Request the HTML page.
	res, err := http.Get("https://www.w3schools.com/tags/default.asp")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var tag string
	// Find the review items
	doc.Find(".ws-table-all").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		tag = s.Find("a").Text()
	})

	var tagList []string
	tags := strings.SplitN(tag, ">", -1)
	for _, v := range tags {
		v = strings.Replace(v, "<", "", -1)
		v = strings.TrimSpace(v)
		tagList = append(tagList, v)
	}

	tagList = tagList[:len(tagList)-1]
	Write("tag", tagList)
	return tagList
}

// CheckEmpty will cehck the list of string contain ""
func CheckEmpty(list []string) bool {
	for _, s := range list {
		if len(s) == 0 {
			return true
		}
	}
	return false
}

// RemoveEmpty remove len ==0 string
func RemoveEmpty(list []string) []string {
	var noEmptyList []string
	for _, s := range list {
		if len(s) != 0 {
			noEmptyList = append(noEmptyList, s)
		}
	}
	removed := []string{string(byte(9)), string(byte(10))}
	ret := Split(noEmptyList, removed)
	for i := 0; i < len(ret); i++ {
		ret[i] = strings.TrimSpace(ret[i])
	}
	if CheckEmpty(list) {
		return RemoveEmpty(ret)
	}
	return ret
}

// Split can input all the seperated characters
func Split(strs []string, sep []string) []string {
	if len(sep) == 0 {
		return strs
	}
	var ret = []string{}
	for _, s := range strs {
		ret = append(ret, strings.SplitN(s, sep[0], -1)...)
	}
	return Split(ret, sep[1:])
}

// Write can wrtte string list to file by line
func Write(fileName string, content []string) {
	f, err := os.Create(fileName + ".txt")
	if err != nil {
		fmt.Println("os Create error: ", err)
		return
	}
	defer f.Close()
	bw := bufio.NewWriter(f)
	for _, c := range content {
		bw.WriteString(c)
		bw.WriteString("\n")
	}
	bw.Flush()
}

// Read can read the file by line to string list
func Read(fileName string) []string {
	f, err := os.Open(fileName + ".txt")
	if err != nil {
		fmt.Println("os Open error: ", err)
		return nil
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var ret []string
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("br ReadLine error: ", err)
			return nil
		}
		ret = append(ret, string(line))

	}
	return ret
}
