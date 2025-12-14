package password

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	saltPassword    = "qkhPAGA13HocW3GAEWwb"
	defaultPassword = "123456"
)

func GeneratePassword(str string) (password string) {
	// md5
	// 对原始密码进行MD5哈希计算，生成128位（16 字节）的哈希值
	m := md5.New()
	m.Write([]byte(str))
	mByte := m.Sum(nil)

	// hmac
	// 使用 HMAC（密钥哈希消息认证码）对MD5 结果进行二次加密
	// 使用固定的盐值
	// 生成 256 位（32 字节）的 HMAC 哈希值
	// 转换为十六进制字符串输出
	h := hmac.New(sha256.New, []byte(saltPassword))
	h.Write(mByte)

	password = hex.EncodeToString(h.Sum(nil))

	return
}

func ResetPassword() (password string) {
	m := md5.New()
	m.Write([]byte(defaultPassword))
	mStr := hex.EncodeToString(m.Sum(nil))

	password = GeneratePassword(mStr)

	return
}

func GenerateLoginToken(id int32) (token string) {
	// 对用户ID和固定盐值进行MD5哈希计算
	// 生成 128 位（16 字节）的哈希值
	m := md5.New()
	m.Write([]byte(fmt.Sprintf("%d%s", id, saltPassword)))

	// 转换为十六进制字符串输出，共32字符
	token = hex.EncodeToString(m.Sum(nil))

	return
}
