package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type MyMux struct{}

// Emb は入力内容をuser.htmlに埋め込むためのもの
type Emb struct {
	Username string
	Age      int
	Gender   string
	Fruit    string
	Interest []string
}

// Msg は入力内容に誤りがある場合にlogin.htmlに埋め込むもの
type Msg struct {
	Message string
}

var emb Emb = Emb{}
var msg Msg = Msg{""}

func main() {
	port := "8080"
	mux := &MyMux{}

	log.Printf("server is running on http://localhost:%s", port)
	log.Print(http.ListenAndServe("localhost:"+port, mux))
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.Redirect(w, r, "/login", 301)
	case "/login":
		handleLogin(w, r)
	case "/user":
		handleUser(w, r)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 入力フォームを返す
		t, _ := template.ParseFiles("public/login.html")
		t.Execute(w, msg)
	}
	if r.Method == "POST" {
		msg.Message = ""
		r.ParseForm()
		// 名前入力チェック ＊usernameはhtmlのnameから来たもの
		if len(r.Form.Get("username")) == 0 {
			msg.Message = msg.Message + "名前が入力されていません。"
		}
		// 年齢入力チェック
		integerized_age, interr := strconv.Atoi(r.Form.Get("age")) // 文字列の変換
		if interr != nil || integerized_age < 0 || integerized_age > 100 {
			msg.Message = msg.Message + "年齢を入力してください。"
		}

		// 性別入力チェック
		if r.Form.Get("gender") == "" {
			msg.Message = msg.Message + "性別を入力して下さい。"
		}

		//fruitとinterestはblankで良い
		fruit := r.Form.Get("fruit")
		if fruit == "" {
			fruit = "未登録"
		}
		var interest []string = r.Form["interest"]
		fmt.Println(interest)
		if len(interest) == 0 {
			interest = append(interest, "未登録")
		}

		// 入力内容によってページ遷移
		if msg.Message == "" {
			emb = Emb{
				Username: r.Form.Get("username"),
				Age:      integerized_age,
				Gender:   r.Form.Get("gender"),
				Fruit:    fruit,
				Interest: interest,
			}
			http.Redirect(w, r, "/user", 301)
		} else {
			http.Redirect(w, r, "/login", 301)
		}
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	// 直接 /user へ行った場合にはリダイレクト
	if emb.Username == "" {
		http.Redirect(w, r, "/login", 301)
	}

	t, _ := template.ParseFiles("public/user.html")
	t.Execute(w, emb)
}
