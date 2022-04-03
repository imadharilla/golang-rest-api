package main

import (
	"backend-homecase/redisdb"
	"backend-homecase/utills"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/upload", uploadFileHandler)
	mux.HandleFunc("/serve/", serveFileHandler)

	// Prepare upload directory
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating the uploads directory:", err.Error())
		return
	}

	// Start Webserver
	if err := http.ListenAndServe(":8090", mux); err != nil {
		fmt.Println("Error starting the webserver:", err.Error())
		return
	}
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Pong")
}

func uploadFileHandler(w http.ResponseWriter, req *http.Request) {
	file, fileHandler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, "Error when parsing the uploaded file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if !utills.IsFileFormatSupported(fileHandler)  {
		http.Error(w, "File Format Not supported! ", http.StatusBadRequest)
		return
	}
	
	// Create a new file in the uploads directory
	destFile, err := os.Create(fmt.Sprintf("./uploads/%s", fileHandler.Filename))
	if err != nil {
		http.Error(w, "Error when trying to write the uploaded file to the upload directory:"+err.Error(), http.StatusInternalServerError)
		return
	}

	defer destFile.Close()

	// Write uploaded file to filesystem
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Error when trying to write the uploaded file to the upload directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File successfully uploaded!")
}

func serveFileHandler(w http.ResponseWriter, req *http.Request) {
	imageWidth, ok := req.URL.Query()["width"]
	// Case of no width query provided
	if !ok {
		serveOriginalImageHandler(w, req)
		return
	}

	width , _ := strconv.Atoi(imageWidth[0])
	requestedFile := strings.TrimPrefix(req.URL.Path, "/serve/")

	// Check if data in Redis 
	imageKey := requestedFile + "?width=" + imageWidth[0]
	val, err := redisdb.GetValue(imageKey)
	if 	err == nil  {
		image := []byte(val)
		w.Write(image)
		return
	}

	f, err := os.Open("uploads/"+requestedFile)
	if err != nil {
		http.Error(w, "The requested file doesn't exist", http.StatusBadRequest)
		return
	}
	resizedImg := utills.ResizeImage(f, uint(width))
	
	// store in redis 
	h := redisdb.SetValue(imageKey, string(resizedImg))
	if h != nil {
		println(err)
	}

	
	w.Write(resizedImg)
}

func serveOriginalImageHandler(w http.ResponseWriter, req *http.Request) {
	requestedFile := strings.TrimPrefix(req.URL.Path, "/serve/")
	http.ServeFile(w, req, "uploads/"+requestedFile)
}
