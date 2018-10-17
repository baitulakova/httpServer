package main

import (
	"net/http"
	"os/exec"
	"os"
	"fmt"
	"strings"
	"io"
	"log"
)

func createStorage() (path string){
	var UN = make([]string,0)
	out, _ := exec.Command("whoami").Output()
	username:=string(out)
	for i:=0;i<len(username);i++{
		if username[i]==10{
			continue
		}else {UN=append(UN,string(username[i]))}
	}
	fileStorage := "/home/"+strings.Join(UN,"")+"/httpServerStorage"
	err:=os.MkdirAll(fileStorage,os.ModePerm)
	if err!=nil{
		fmt.Println("error",err)
	}
	return fileStorage
}

func uploadFileHandler(w http.ResponseWriter,r *http.Request){
	if r.Method=="POST" {
		file, h, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("Can't get file from request"))
		}
		fileStorage := createStorage()
		src, oserror := os.Create(fileStorage + h.Filename)
		if oserror != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't create file"))
		}
		defer file.Close()
		f := io.Reader(file)
		defer src.Close()
		io.Copy(src, f)
	}
}

func main() {
	http.HandleFunc("/upload",uploadFileHandler)
	log.Println("Server is working on port :8080")
	http.ListenAndServe(":8080", nil)
}
