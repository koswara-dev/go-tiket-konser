package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StorageProvider interface {
	UploadFile(file *multipart.FileHeader, c *gin.Context) (string, error)
}

type localStorageProvider struct {
	UploadDir string
	BaseURL   string
}

// pastikan folder penyimpanan local terbentuk
func NewLocalStorageProvider(uploadDir string, baseURL string) StorageProvider {
	_ = os.MkdirAll(uploadDir, os.ModePerm)
	return &localStorageProvider{
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}
}

// implementasi method UploadFile
func (s *localStorageProvider) UploadFile(file *multipart.FileHeader, c *gin.Context) (string, error) {
	// 1. validasi batas ukuran file (maksimal 2MB)
	var maxFileSize int64 = 2 * 1024 * 1024 // 2MB
	if file.Size > maxFileSize {
		return "", errors.New("file size exceeds 2MB limit")
	}

	// 2. validasi ekstensi gambar & pdf
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".pdf" {
		return "", errors.New("invalid file extension")
	}

	// 3. generating unique filename
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.UploadDir, fileName)

	// 4. save file ke local storage
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// 5. return URL file disimpan
	return fmt.Sprintf("%s/uploads/%s", s.BaseURL, fileName), nil
}
