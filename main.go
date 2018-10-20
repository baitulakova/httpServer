package main

import (
	"net/http"
	"os"
	"fmt"
	"io"
	"log"
)

func createStorage() (path string){
	userHome:=os.Getenv("HOME")
	fileStorage := userHome+"/httpServerStorage/"
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
		src, oserror := os.Create(fileStorage +h.Filename)
		if oserror != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't create file"))
		}
		defer file.Close()
		f := io.Reader(file)
		defer src.Close()
		io.Copy(src, f)
		log.Println("Uploaded ",h.Filename," file")
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request){
	file:=r.URL.Query().Get("filename")
	f,err:=os.Open(createStorage()+file)
	defer f.Close()
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.ServeFile(w,r,createStorage()+file)
	io.Copy(w,f)
}

func main() {
	http.HandleFunc("/upload",uploadFileHandler)
	http.HandleFunc("/download",downloadHandler)
	log.Println("Server is working on port :8080")
	http.ListenAndServe(":8080", nil)
}
