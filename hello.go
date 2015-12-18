package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	PORT         = ":80"
	UPLOAD_DIR   = "./uploads"
	TPL_DIR      = "./views"
	DOWNLOAD_DIR = "./downloads"
	TimeoutLimit = 10
)

func main() {
	time.Now().Minute()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/up", uploadHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/spider", spiderHandler)
	http.HandleFunc("/img", imgSpiderHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err.Error())
		return
	}

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}

//get输出上传页面 post上传文件
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//使用Go模板库输出
		t, err := template.ParseFiles(TPL_DIR + "/upload.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)

		//以字符串输出L
		//renderHtml(w, TPL_DIR+"/upload.html")
		return
	}

	if r.Method == "POST" {
		f, h, err := r.FormFile("img")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileNmae := h.Filename
		defer f.Close()
		t, err := os.Create(UPLOAD_DIR + "/" + fileNmae)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer t.Close()

		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//跳转到 显示图片
		http.Redirect(w, r, "/view?id="+fileNmae, http.StatusFound)
	}
}

//输出上传之后的文件
func viewHandler(w http.ResponseWriter, r *http.Request) {
	fileNmae := r.FormValue("id")
	filePath := UPLOAD_DIR + "/" + fileNmae

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, filePath)
}

//请求一个链接 并输出内容
func spiderHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://www.yinwang.org/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(body))
}

//解析HTML
func imgSpiderHandler(w http.ResponseWriter, r *http.Request) {
	url := "http://mt.locojoy.com/chengka"
	index := r.FormValue("index")

	var i int
	if strings.EqualFold(index, "") {
		i = 0
	} else {
		i, _ = strconv.Atoi(index)
	}

	fmt.Println(strconv.Itoa(i))

	if i > 1 {
		url = url + "/" + strconv.Itoa(i) + ".html"
	}

	fmt.Println(url)

	doc, err := goquery.NewDocument(url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	doc.Find(".js_ulWrap_SH").Find("li").Each(func(i int, li *goquery.Selection) {
		url, exists := li.Find("img").First().Attr("src")
		title := li.Find("a").Eq(1).Text()
		if !exists {
			fmt.Println("no exists")
		} else {
			fmt.Println(url + "---" + title)

			go download(url) //使用并发下载

		}
	})

	x, _ := doc.Html()
	io.WriteString(w, x)
}

//下载
func download(url string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	defer res.Body.Close()

	fileName := filepath.Base(url)

	temp, err := os.Create(DOWNLOAD_DIR + "/" + fileName)

	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}

	io.Copy(temp, res.Body)

	defer temp.Close()
}

//输出HTML
func renderHtml(w http.ResponseWriter, htmlPath string) {
	s, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(s))
}
