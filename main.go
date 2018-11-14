package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"net/http"
)


func main() {

	fs := http.FileServer(http.Dir("public/js"))
	http.Handle("/js/", http.StripPrefix("/js/", fs))
	http.HandleFunc("/", IndexPage)
	http.HandleFunc("/allcomments", ShowComments)
	http.HandleFunc("/addcomments", AddComments)
	http.HandleFunc("/api/addcomments", ApiAddComments)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/clear_comments_a619d974658f3e749b2d88b215baea46", ClearComments)
	http.HandleFunc("/api/signin", ApiSignin)
	http.HandleFunc("/updatefirmware", UpdateFirmware)
	http.HandleFunc("/api/updatefirmware", ApiUpdateFirmware)
	http.HandleFunc("/downloadfirmware", DownloadFirmware)
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
