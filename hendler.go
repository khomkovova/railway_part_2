package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)
type AllCommands struct {
	Comments string `json:"comments"`
}
var PrivateKey, err2 = rsa.GenerateKey(rand.Reader, 2048)
var PublicKey = &PrivateKey.PublicKey

type CredentialsSignin struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
type Comments struct {
	Comments string `json:"comments"`
}

func IndexPage(w http.ResponseWriter, r *http.Request)  {
	data, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(data))
}

func ShowComments(w http.ResponseWriter, r *http.Request) {
	data:= generateShowCooments()
	w.Write([]byte(data))

}

func AddComments(w http.ResponseWriter, r *http.Request)  {
	data, err := ioutil.ReadFile("public/addcomments.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(data))
}
func ApiAddComments(w http.ResponseWriter, r *http.Request)  {
	var comments Comments
	err := json.NewDecoder(r.Body).Decode(&comments)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	f, err := os.OpenFile("public/comments.db", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	com := commentsCheck(comments.Comments)
	com += "<br><h3>---</h3>"
	_, err = f.Write([]byte(com))
	if err != nil {
		return
	}
	f.Close()
}


func Signin(w http.ResponseWriter, r *http.Request)  {
	data, err := ioutil.ReadFile("public/signin.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(data))
}
func ApiSignin(w http.ResponseWriter, r *http.Request) {
	var creds CredentialsSignin

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if creds.Username != "admin" || creds.Password != "pass"{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cookie, err3 := encodeCookie(creds.Username)
	if err3 == false{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &cookie)
	//http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ApiUpdateFirmware(w http.ResponseWriter, r *http.Request)  {
	user := decodeCookie(r)
	if user == ""{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if user != "admin"{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("firmware")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Upload file name = ", handler.Filename);
	err = os.Remove("mycheck.cpp")
	if err != nil {
		fmt.Println("This file is deleted")
	}
	f, err := os.OpenFile("./"+"mycheck.cpp", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	io.Copy(f, file)
	file.Close()
	f.Close()
	if !compileFirmware(){
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("File uploaded success"))

}
func UpdateFirmware(w http.ResponseWriter, r *http.Request)  {
	user := decodeCookie(r)
	if user == ""{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if user != "admin"{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	data, err := ioutil.ReadFile("public/updatefirmware.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(data))


}
func DownloadFirmware(w http.ResponseWriter, r *http.Request)  {
	var token string
	for _, cookie := range r.Cookies() {
		token = cookie.Value
	}
	fmt.Println("token =" + token)
	if token != "server1234567890"{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b, err := ioutil.ReadFile("mycheck")
	if err != nil{
		return
	}
	w.Write(b)
	return
}

func compileFirmware() bool {
	//gcc mycheck.cpp -o mycheck
	cmd := exec.Command("gcc", "mycheck.cpp", "-o", "mycheck" )
	//cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func decodeCookie(r *http.Request) string {
	var token string
	for _, cookie := range r.Cookies() {
		token = cookie.Value
	}
	if token == ""{
		return ""
	}
	sDec, _ := b64.StdEncoding.DecodeString(token)
	label := []byte("")
	hash := sha256.New()
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, PrivateKey, sDec, label)
	if err != nil{
		return ""
	}
	user := string(plainText)
	fmt.Println("token=",token)
	return user
}

func encodeCookie(user string) (http.Cookie, bool)  {
	message := []byte(user)
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, PublicKey, message, label)
	if err != nil{
		return http.Cookie{}, false
	}
	sEnc := b64.StdEncoding.EncodeToString(ciphertext)
	cookie := http.Cookie{Name:"token", Value:sEnc}
	cookie.Path = "/"
	return cookie, true

}

func generateShowCooments() string {
	comments, err := ioutil.ReadFile("public/comments.db")
	if err != nil{

		return ""
	}
	data :=`<!DOCTYPE html>
		<html lang="en" xmlns:v-on="http://www.w3.org/1999/xhtml">
	<head>
	<meta charset="UTF-8">
		<title>Comments</title>
 <link rel="stylesheet" href="js/main.css">
		</head>
		<body>
		<footer>
		<a class="link" href="signin">Signin</a>
		<a class="link"  href="/">Main</a>
		<a  class="link" href="allcomments">All Comments</a>
		<a  class="link" href="addcomments">Add Comments</a>
        </footer>
<div class="header">
    <h1>All comments</h1>
</div>
		<dir>
		<p>` + string(comments) + `
		</p>
		</dir>
		
</body>
</html>`
return data
}

func commentsCheck(data string) string {
	data = strings.ToLower(data)
	var blacklist []string
	blacklist = []string{"script", "http", ".", "//", "</", "img", "src", "body", "style", "br", "bgsoung", "link", "meta", "div", "iframe", "object", "data", "href", "alert", "document", "cookie", "0x" }
	//a := blacklist
	for _, word := range blacklist {
		data = strings.Replace(data, word, "",-1)
	}
	return data
}