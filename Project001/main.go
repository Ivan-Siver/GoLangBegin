package main

import (
	"database/sql"
	"fmt"

	//$GOROOT = C:\Program Files\Go\src\github.com\go-sql-driver\mysql
	//$GOPATH = C\src\github.com\go-sql-driver\mysql
	//$GOPATH = \Users\ivsiver\go\src\github.com\go-sql-driver\mysql
	//При использовании модулей Go переменная GOPATH (по умолчанию в $HOME/go для Unix и %USERPROFILE%\go для Windows)
	_ "github-com/go-sql-driver/mysql"
)

type User struct {
	Name string `json:"name"`
	Age  uint16 `json:"age"`
}

func main() {
	//127.0.0.1
	db, err := sql.Open("mysql", "root:ghjnjrjk19@tcp(localhost:3306)/gomysqldb")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Установка данных
	//insert, err := db.Query("INSERT INTO `user` (`name`, `age`) VALUES('IVS_2015', 23)")
	//if err != nil {
	//	panic(err)
	//}
	//defer insert.Close()
	fmt.Println("Connected to MySQL-Database")

	//Выборка данных
	res, err := db.Query("SELECT `name`, `age` FROM `user`")
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var user User
		err = res.Scan(&user.Name, &user.Age)
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("User: %s with age %d", user.Name, user.Age))
	}
}
