package main

import (
	"net/http"
	"os"
	"fmt"
	"io"
	"log"
	"strings"
	"path"
	"strconv"
)

var fs=http.FileServer(http.Dir(createStorage()))

var uploadForm=[]byte(
	`<html>
		<body>
		<form action="/upload" method="post" enctype=multipart/form-data>
			Выберите файл: <input type="file" name=file">
			<br>
			<input type="submit" value="Upload">
		</form>
		</body>
</html>
`)

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
	w.Write(uploadForm)
	if r.Method=="POST" {
		r.ParseMultipartForm(32 << 20)
		file, h, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		defer file.Close()
		fileStorage := createStorage()
		src, oserror := os.Create(fileStorage +h.Filename)
		defer src.Close()
		if oserror != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Can't create file"))
		}
		f := io.Reader(file)
		io.Copy(src, f)
		log.Println("Uploaded ",h.Filename," file")
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request){
	file:=r.URL.Query().Get("filename")
	if file==""{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}
	f,err:=os.Open(createStorage()+file)
	defer f.Close()
	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}
	fileInfo,_:=f.Stat()
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(),10))
	log.Println(r.Header)
	io.Copy(w,f)
}

func imagesHandler(w http.ResponseWriter,r *http.Request){
	a:=r.URL.EscapedPath()
	s:=strings.Split(a,"/")
	last:=s[len(s)-1]
	imgPath:=path.Join(createStorage(),"images",last)
	log.Println(imgPath)
	http.ServeFile(w,r,imgPath)
}

func main() {
	http.HandleFunc("/upload",uploadFileHandler)
	http.HandleFunc("/download",downloadHandler)
	http.HandleFunc("/images/",imagesHandler)
	http.Handle("/", fs)
	log.Println("Server is working on port :8080")
	http.ListenAndServe(":8080", nil)
}
