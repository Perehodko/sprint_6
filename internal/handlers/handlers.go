package handlers

import (
    "fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	
	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func BackHTML(res http.ResponseWriter, req *http.Request) {
	htmlContent, err := os.ReadFile("index.html")
	if err != nil {
		http.Error(res, "Could not read HTML file", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	_, err = res.Write(htmlContent)
	if err != nil {
		http.Error(res, "Could not write HTML content: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandlerUpload(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(res, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := req.FormFile("myFile")
	if err != nil {
		http.Error(res, "Failed to get file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(res, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	content := string(fileBytes)
	result, err := service.AutoDetectAndConvert(content)
	if err != nil {
		http.Error(res, "Conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	originalExt := filepath.Ext(fileHeader.Filename)
	timestamp := time.Now().UTC().Format("20060102_150405")
	newFilename := fmt.Sprintf("converted_%s%s", timestamp, originalExt)
	
	err = os.WriteFile(newFilename, []byte(result), 0644)
	if err != nil {
		http.Error(res, "Failed to save result file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(result))
}

