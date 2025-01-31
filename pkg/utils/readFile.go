package utils

import (
	"io"
	"mime/multipart"
)

// ReadFileContent 将文件转化为比特流
func ReadFileContent(file *multipart.FileHeader) (
	[]byte, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// 读取文件内容
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return content, nil
}
