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
)

var PrivateKey, err2 = rsa.GenerateKey(rand.Reader, 2048)
var PublicKey = &PrivateKey.PublicKey

type CredentialsSignin struct {
	Password string `json:"password"`
	Username string `json:"username"`
}


func Signin(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UpdateFirmware(w http.ResponseWriter, r *http.Request)  {
	user := decodeCookie(r)
	if user == ""{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if user != "admin"{
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("firmware")
	if err != nil {
		fmt.Println("123")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("./"+"mycheck.cpp", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("124")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.Copy(f, file)
	file.Close()
	f.Close()
	if !compileFirmware(){
		fmt.Println("125")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("Ok upload file"))

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
	fmt.Printf("Check return %s", out.String())
	return true
}

func decodeCookie(r *http.Request) string {
	var token string
	for _, cookie := range r.Cookies() {
		token = cookie.Value
	}
	if token == ""{
		fmt.Println("error")
		return ""
	}
	sDec, _ := b64.StdEncoding.DecodeString(token)
	label := []byte("")
	hash := sha256.New()
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, PrivateKey, sDec, label)
	if err != nil{
		fmt.Println("error")
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
	return cookie, true

}
