package main

import "fmt"
import "net/http"
import "log"
import "io/ioutil"
import "regexp"
import "os"

import iconv "iconv-go-master"

type Spider struct {
	Page int
}

func test_write(mystr string, filename string) {
	//fout, err := os.Create(filename)
	fout, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fout.Close()
	fout.WriteString(mystr)
	fout.WriteString("\n")
}

func (this *Spider) test_http(content string) {

	filename := "./test.txt"
	http_reg := regexp.MustCompile(`<h4> <a href="(.*?)"`)
	http_result := http_reg.FindAllStringSubmatch(content, -1)
	for _, mystr_slice := range http_result {
		mystr_url := fmt.Sprintf("%s%s", "https://www.neihanba.com", mystr_slice[1])
		fmt.Printf("============ %s ================\n", mystr_url)
		content, status := this.HttpGet(mystr_url)
		if status != 200 {
			continue
		}
		/*
			http_reg = regexp.MustCompile(`<h1>(.*?)</h1>`)
			http_result = http_reg.FindAllStringSubmatch(content, -1)
			for _, mystr_text := range http_result {
				fmt.Println(mystr_text[1])
			}
		*/
		http_reg = regexp.MustCompile(`<td><p>(?s:(.*?))</p></td>`)
		http_result = http_reg.FindAllStringSubmatch(content, -1)
		for _, mystr_text := range http_result {
			fmt.Println(mystr_text[1])
			test_write(mystr_text[1], filename)
		}
	}
}

func (this *Spider) Spider_one_page() {
	fmt.Println("正在爬取 ", this.Page, " 页")
	url := ""
	if this.Page == 1 {
		url = "https://www.neihanba.com/dz/index.html"
	} else {
		url = fmt.Sprintf("%s%d%s", "https://www.neihanba.com/dz/list_", this.Page, ".html")
	}
	content, status := this.HttpGet(url)
	if status != 200 {
		fmt.Println("http Get error, status is ", status)
		return
	}
	this.test_http(content)

}

func (this *Spider) HttpGet(url string) (content string, statusCode int) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		content = ""
		statusCode = -100
		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		content = ""
		statusCode = resp.StatusCode
		return
	}

	//content = string(data)
	//statusCode = resp.StatusCode
	out := make([]byte, len(data))
	out = out[:]

	iconv.Convert(data, out, "gb2312", "utf-8")
	content = string(out)
	statusCode = resp.StatusCode

	return
}

func (this *Spider) DoWork() {
	fmt.Println("Spider begin to work")
	this.Page = 1
	var cmd string

	for {
		fmt.Println("请输入任意键爬取下一页，输入exit退出")
		fmt.Scanf("%s", &cmd)

		if cmd == "exit" {
			fmt.Println("exit")
			break
		}
		this.Spider_one_page()
		this.Page++
	}
}

func main() {
	sp := new(Spider)
	sp.DoWork()

}
