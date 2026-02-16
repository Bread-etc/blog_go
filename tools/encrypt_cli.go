package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

// 使用 /user/public_key 所复制的内容
const ServerPublicKey = `-----BEGIN PUBLIC KEY-----xxx-----END PUBLIC KEY-----\n`

func main() {
	// 处理公钥格式
	pemStr := strings.ReplaceAll(ServerPublicKey, `\n`, "\n")
	pemStr = strings.TrimSpace(pemStr)

	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		fmt.Println("❌ 解析公钥失败：请检查 PEM 格式是否正确")
		return
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Printf("❌ 解析 PKIX 公钥失败: %v\n", err)
		return
	}
	pubKey := pubInterface.(*rsa.PublicKey)

	// 交互式输入密码
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("🔒 请输入要加密的密码：")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if password == "" {
		fmt.Println("❌ 密码不能为空")
		return
	}

	// 加密 (PKCS1v1.5)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(password))
	if err != nil {
		fmt.Printf("❌ 加密失败: %v\n", err)
		return
	}

	// Base64 编码
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	fmt.Println("\n✅ 加密成功! (RSA/PKCS1v1.5 + Base64)")
	fmt.Println("--------------------------------------------------")
	fmt.Println(encoded)
	fmt.Println("--------------------------------------------------")
	fmt.Println("📋 已输出到控制台，请手动复制")
}
