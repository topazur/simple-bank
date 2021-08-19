package util

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

/// init() 函数将在第一次使用包时自动调用。
/*
	通过调用 rand.Seed() 为随机生成器设置种子值。
	通常，种子值通常设置为当前时间。
	由于 rand.Seed() 期望 int64 作为输入，所以转换为UnixNano。
*/
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 返回一个介于 min 和 max 之间的随机 int64 数
func RandomInt(min, max int64) int64 {
	// rand.Int63n() 函数返回 0 和 n-1 之间的随机整数
	return min + rand.Int63n(max-min+1)
}

// RandomString 返回一个指定长度的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner 返回username
func RandomOwner() string {
	return RandomString(6)
}

// RandomAmountMoney 返回随机金额
func RandomAmountMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency 生成随机货币代码
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "RMB"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail 生成随机电子邮件
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// Mkdir 定义一个创建文件目录的方法
func MkdirFolder(basePath string) string {
	//	1.获取当前时间,并且格式化时间
	folderName := time.Now().Format("2006/01/02")
	folderPath := filepath.Join(basePath, folderName)
	//使用mkdirall会创建多层级目录
	os.MkdirAll(folderPath, os.ModePerm)
	return folderPath
}
