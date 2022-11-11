package main

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/logger"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "123456",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 连接 mysql
	db, err = sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接 失败会报错
	err = db.Ping()
	logger.LogError(err)

}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

type Article struct {
	Title, Body string
	ID          int64
}

func (a Article) Link() string {
	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
	if err != nil {
		logger.LogError(err)
		return ""
	}
	return showURL.String()
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"

	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	return article, err
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {

	id := getRouteVariable("id", r)
	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		updateURL, _ := router.Get("articles.update").URL("id", id)

		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")

		logger.LogError(err)

		err = tmpl.Execute(w, data)

		logger.LogError(err)
	}
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {

	id := getRouteVariable("id", r)

	_, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if title == "" {
			errors["title"] = "标题不能为空"
		} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 30 {
			errors["title"] = "标题长度介于 3 - 40"
		}

		if body == "" {
			errors["body"] = "内容不能为空"
		} else if utf8.RuneCountInString(body) < 10 {
			errors["body"] = "内容长度不能小于10"
		}

		if len(errors) == 0 {
			query := "update articles set title = ?, body = ? where id = ?"
			rs, err := db.Exec(query, title, body, id)
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改!")
			}
		} else {
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}

			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)

			err = tmpl.Execute(w, data)
			logger.LogError(err)
		}
	}

}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * from articles")
	logger.LogError(err)
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		logger.LogError(err)

		articles = append(articles, article)
	}

	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")

	logger.LogError(err)

	err = tmpl.Execute(w, articles)

	logger.LogError(err)
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3 - 40"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于10个字节"
	}

	return errors
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")

	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需要大于或等于10个字节"
	}

	// 检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功 ID 为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {

		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func saveArticleToDB(title string, body string) (int64, error) {
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1. 获取一个 prepare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles(title,body) VALUES(?,?)")
	// 例行的错误检测
	if err != nil {
		return 0, err
	}

	// 2.在此函数运行结束关闭此语句 防止占用SQL链接
	defer stmt.Close()

	// 3.执行请求 传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	// 4.插入成功的话 会返回自增 ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, err
}

// 使用中间件

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(w, r)
	})
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
 <html lang="en">
 <head>
     <title>创建文章 —— 我的技术博客</title>
 </head>
 <body>
     <form action="%s?test=data" method="post">
         <p><input type="text" name="title"></p>
         <p><textarea name="body" cols="30" rows="10"></textarea></p>
         <p><button type="submit">提交</button></p>
     </form>
 </body>
 </html>
	`
	storeURL, _ := router.Get("articles.store").URL()

	fmt.Fprintf(w, html, storeURL)
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)

	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		rowsAffected, err := article.Delete()

		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			if rowsAffected > 0 {
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("delete from articles where id =" + strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {

	bootstrap.SetupDB()

	router = bootstrap.SetupRoute()

	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")

	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.created")

	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)
	// 通过命名路由 URL 事例

	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "1")
	fmt.Println("articles: ", articleURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))

}
