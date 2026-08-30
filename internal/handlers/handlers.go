package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func MainHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

func UploadHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("[TRACE] UploadHandler, method=%s", req.Method)

	if req.Method != http.MethodPost {
		log.Printf("[WARN] Wrong method: %s", req.Method)
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := req.FormFile("myFile")
	if err != nil {
		log.Printf("[WARN] FormFile error (client issue): %v", err)
		http.Error(res, "No file uploaded: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	log.Printf("[INFO] File received: %q, size=%d bytes", header.Filename, header.Size)

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[ERROR] ReadAll error: %v", err)
		http.Error(res, "Cannot read file content", http.StatusInternalServerError)
		return
	}
	inputString := string(data)

	preview := inputString
	if len(preview) > 80 {
		preview = preview[:80] + "..."
	}
	log.Printf("[DEBUG] Input (preview): %q", preview)

	log.Printf("[TRACE] Calling service.Convert...")
	result, err := service.Convert(inputString)
	if err != nil {
		log.Printf("[ERROR] Convert error: %v", err)
		http.Error(res, "Conversion failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] Conversion successful")

	name := time.Now().UTC().Format("2006-01-02_15-04-05") + filepath.Ext(header.Filename)
	out, err := os.Create(name)
	if err != nil {
		log.Printf("[ERROR] Create file error: %v", err)
		http.Error(res, "Cannot create result file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = out.Write([]byte(result))
	if err != nil {
		log.Printf("[ERROR] Write to file error: %v", err)
		http.Error(res, "Cannot write result file", http.StatusInternalServerError)
		return
	}
	log.Printf("[OK] Result saved to: %s", name)

	_, err = res.Write([]byte(result))
	if err != nil {
		log.Printf("[ERROR] Response write error: %v", err)
		return
	}
}
