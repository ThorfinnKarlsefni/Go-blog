package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "<h1>Hello ,Cheung!<h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求页面未找到:(</h1></br>"+
			"<p>如有疑惑,请联系我们.</p>")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;chartext=utf-8")
	fmt.Fprintf(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:402832626@qq.com\">Nick@qq.com </a>")
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", aboutHandler)

	// article details
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 3)[0]
		fmt.Fprint(w, "文章 ID:"+id)
	})

	// lists or created
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w, "访问文章列表")
		case "POST":
			fmt.Fprint(w, "创建新的文章")
		}
	})
	http.ListenAndServe(":3000", router)

}
