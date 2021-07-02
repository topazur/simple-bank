package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	// 密码将使用 bcrypt 散列函数进行散列以产生散列值。
	// cost：此参数将决定算法的密钥扩展轮数或迭代次数。
	// 字符串切片生成一个22个字符的字符串Salt（随机的，即便是相同密码），采用base54编码
	// ALG：是哈希算法标识符。
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, hashedPassword string) error {
	// password 登录时传递的参数；hashedPassword数据库中保存的hash密码
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
