package main

import (
	"database/sql"
	"fmt"
	"github-com/gorilla/mux"
	"html/template"
	"net/http"

	//$GOROOT = C:\Program Files\Go\src\github.com\go-sql-driver\mysql
	//$GOPATH = C\src\github.com\go-sql-driver\mysql
	//$GOPATH = \Users\ivsiver\go\src\github.com\go-sql-driver\mysql
	//При использовании модулей Go переменная GOPATH (по умолчанию в $HOME/go для Unix и %USERPROFILE%\go для Windows)
	_ "github-com/go-sql-driver/mysql"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:ghjnjrjk19@tcp(localhost:3306)/gomysqldb")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Выборка данных
	res, err := db.Query("SELECT * FROM `article`")
	if err != nil {
		panic(err)
	}

	//Обнуление списка статей
	posts = []Article{}
	//Формирование списка статей
	for res.Next() {
		var post Article //post = каждая отдельная статья
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
		//fmt.Println(fmt.Sprintf("Статья: %s with id: %d", post.Title, post.Id))
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все поля заполнены")
	} else {
		db, err := sql.Open("mysql", "root:ghjnjrjk19@tcp(localhost:3306)/gomysqldb")

		if err != nil {
			panic(err)
		}

		defer db.Close()

		//Установка данных
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `article` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')",
			title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:ghjnjrjk19@tcp(localhost:3306)/gomysqldb")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Выборка данных
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `article` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	//Обнуление списка статей
	showPost = Article{}
	//Формирование списка статей
	for res.Next() {
		var post Article //post = каждая отдельная статья
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)
}

func handleFunc() {
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/save_article", save_article).Methods("POST")
	router.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
