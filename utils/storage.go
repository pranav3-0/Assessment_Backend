package utils

import (
	"os"
	"path/filepath"
)

func SaveFileToDisk(folder string, fileName string, data []byte) (string, error) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(folder, fileName)

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", err
	}

	urlPath := "/" + filepath.ToSlash(filePath)

	return urlPath, nil
}