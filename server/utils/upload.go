package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// 允许的图片类型
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// 允许的文件扩展名
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// ValidateImageFile 验证上传的文件是否为图片
func ValidateImageFile(file *multipart.FileHeader, maxSize int64) error {
	// 检查文件大小
	if file.Size > maxSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %d MB", maxSize/(1024*1024))
	}

	// 检查文件类型
	contentType := file.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return fmt.Errorf("不支持的文件类型，仅支持 jpg, png, gif, webp")
	}

	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("不支持的文件扩展名")
	}

	return nil
}

// SaveUploadedFile 保存上传的文件到指定目录
func SaveUploadedFile(file *multipart.FileHeader, destPath string) error {
	// 确保目标目录存在
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}
