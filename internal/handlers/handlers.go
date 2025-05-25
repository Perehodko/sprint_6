package handlers

import (
    "fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	
    "github.com/go-chi/chi/v5"
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func BackHTML(res http.ResponseWriter, req *http.Request) {
	htmlContent, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(res, "Could not read HTML file", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	res.Write(htmlContent)
}

func HandleUpload(res http.ResponseWriter, req *http.Request) {
    err := req.ParseMultipartForm(10 << 20) 
    if err != nil {
        http.Error(res, "Failed to parse form", http.StatusBadRequest)
        return
    }

    file, header, err := req.FormFile("myFile")
    if err != nil {
        http.Error(res, "File not found in form", http.StatusBadRequest)
        return
    }
    defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(res, "Error reading file content", http.StatusInternalServerError)
		return
	}

	convertedData, err := service.AutoDetectAndConvert(string(fileBytes))
	if err != nil {
		http.Error(res, fmt.Sprintf("Conversion error: %v", err), http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf("result_%s%s", 
		time.Now().UTC().Format("20060102_150405"),
		filepath.Ext(header.Filename))

	err = os.WriteFile(fileName, []byte(convertedData), 0644)
	if err != nil {
		http.Error(res, "Error saving result file", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(convertedData))
}

func main() {
    r := chi.NewRouter()

    r.Get("/", BackHTML)
	r.Get("/upload", HandleUpload)

    err := http.ListenAndServe(":8080", r) 
    if err != nil {
        panic(err)
    }
}
