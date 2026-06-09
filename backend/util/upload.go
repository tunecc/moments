package fs_util

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

// 允许上传的文件扩展名白名单(图片与常见视频)
var allowedUploadExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".svg":  false, // SVG 可内嵌脚本,明确禁止
	".mp4":  true,
	".webm": true,
	".mov":  true,
	".ogg":  true,
}

// 允许的 MIME 前缀
var allowedMimePrefixes = []string{"image/", "video/"}

// MaxUploadSize 单个上传文件的最大体积(50MB)
const MaxUploadSize = 50 << 20

// ValidateUploadExt 校验上传文件扩展名是否在白名单内。
func ValidateUploadExt(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return fmt.Errorf("文件缺少扩展名")
	}
	allowed, ok := allowedUploadExts[ext]
	if !ok || !allowed {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}
	return nil
}

// ValidateUploadContentType 基于文件内容嗅探的 MIME 校验。
// content 为文件头部字节(至少 512 字节)。
func ValidateUploadContentType(content []byte) error {
	mime := http.DetectContentType(content)
	for _, prefix := range allowedMimePrefixes {
		if strings.HasPrefix(mime, prefix) {
			return nil
		}
	}
	// http.DetectContentType 对部分视频(如 mp4)可能返回 application/octet-stream,
	// 这里放行 octet-stream,交由扩展名白名单兜底。
	if mime == "application/octet-stream" {
		return nil
	}
	return fmt.Errorf("不支持的文件内容类型: %s", mime)
}
