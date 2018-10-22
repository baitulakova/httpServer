package main

import (
	"net/http"
	"log"

	handlers "github.com/baitulakova/httpServer/handler"
)

var fs=http.FileServer(http.Dir(handlers.createStorage()))

func main() {
	http.HandleFunc("/upload",handlers.uploadFileHandler)
	http.HandleFunc("/download",handlers.downloadHandler)
	http.HandleFunc("/images/",handlers.imagesHandler)
	http.Handle("/", fs)
	log.Println("Server is working on port :8080")
	http.ListenAndServe(":8080", nil)
}
