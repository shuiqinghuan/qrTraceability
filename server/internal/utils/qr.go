package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode 生成包含产品ID的二维码
func GenerateQRCode(productID int) (string, error) {
	// 构建产品详情页URL
	url := fmt.Sprintf("http://localhost:3000/product/%d", productID)

	// 生成二维码
	var buf bytes.Buffer
	err := qrcode.WriteFile(url, qrcode.Medium, 256, "./qrcodes/product_"+fmt.Sprintf("%d", productID)+".png")
	if err != nil {
		return "", err
	}

	// 读取生成的二维码文件并转换为base64
	// 这里简化处理，实际项目中可能需要更复杂的存储方案
	return "./qrcodes/product_"+fmt.Sprintf("%d", productID)+".png", nil
}