package utils

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode 生成二维码
// url: 要编码的URL
// size: 二维码大小（像素）
// 返回二维码的字节数组（PNG格式）
func GenerateQRCode(url string, size int) ([]byte, error) {
	// 直接生成到字节数组
	png, err := qrcode.Encode(url, qrcode.Medium, size)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code: %w", err)
	}

	return png, nil
}
