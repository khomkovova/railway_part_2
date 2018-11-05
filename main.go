package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"log"
	"net"
	"net/http"
)

var cache redis.Conn

func main() {

	//InitDataBase()
	//TestDb()
	// "Signin" and "Signup" are handler that we will implement
	fs := http.FileServer(http.Dir("public/js"))
	http.Handle("/js/", http.StripPrefix("/js/", fs))
	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/allcomments", ShowComments)
	http.HandleFunc("/addcomments", AddComments)
	http.HandleFunc("/api/addcomments", ApiAddComments)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/api/signin", ApiSignin)
	http.HandleFunc("/updatefirmware", UpdateFirmware)
	http.HandleFunc("/api/updatefirmware", ApiUpdateFirmware)
	//http.HandleFunc("/api/getcomments",GetComments)
	http.HandleFunc("/downloadfirmware", DownloadFirmware)
	//http.HandleFunc("/", Info)
	//http.HandleFunc("/upload", Upload)
	//http.HandleFunc("/refresh", Refresh)
	// start the server on port 8000
	l, err := net.Listen("tcp4", ":12345")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(l, nil))
}

func TestDb(){
	db, err := sql.Open("mysql", "root:Remidolov12345@@/test?charset=utf8")
	if err != nil{
		return
	}
	rows, err := db.Query("SELECT * FROM test")
	if err != nil{
		return
	}
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		if err != nil{
			return
		}
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}
}
