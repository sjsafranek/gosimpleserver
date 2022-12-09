package utils

import (
	"os"
	"path/filepath"
)

func GetSubDirectories(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if nil != info && info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}