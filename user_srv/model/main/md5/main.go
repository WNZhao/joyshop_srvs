package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

// genMd5 生成字符串的 MD5 哈希值（返回十六进制字符串）
func genMd5(code string) string {
	hash := md5.New()
	hash.Write([]byte(code))
	return hex.EncodeToString(hash.Sum(nil))
}

func main() {
	// Using the default options
	//salt, encodedPwd := password.Encode("generic password", nil)
	//fmt.Println(encodedPwd) // 8c7a9f3b
	//fmt.Println(salt)       // 8c7a9f3b
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check) // true

	// Using custom options
	options := &password.Options{10, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	//把盐值和加密后的密码存储到数据库中 盐值整合到加密后的密码中
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPassword)
	fmt.Println(salt) // 8c7a9f3b
	passwordInfo := strings.Split(newPassword, "$")
	//fmt.Println(passwordInfo)
	check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check) // true
}
