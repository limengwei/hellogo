package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	PORT       = ":80"
	UPLOAD_DIR = "./uploads"
	TPL_DIR    = "./views"
)

func main() {

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/up", uploadHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/spider", spiderHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err.Error())
		return
	}

}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "hello")
}

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
		http.Redirect(w, r, "/view?id="+fileNmae, http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fileNmae := r.FormValue("id")
	filePath := UPLOAD_DIR + "/" + fileNmae

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, filePath)
}

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

func renderHtml(w http.ResponseWriter, htmlPath string) {
	s, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(s))
}
