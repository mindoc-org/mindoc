package cryptil

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
	"io"
	"crypto/rand"
)

//对称加密与解密之加密【从Beego中提取出来的】
//@param            value           需要加密的字符串
//@param            secret          加密密钥
//@return           encrypt			返回的加密后的字符串
func Encrypt(value, secret string) (encrypt string) {
	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := hmac.New(sha1.New, []byte(secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", h.Sum(nil))
	return strings.Join([]string{vs, timestamp, sig}, ".")
}

//对称加密与解密之解密【从Beego中提取出来的】
//@param            value           需要解密的字符串
//@param            secret          密钥
//@return           decrypt			返回解密后的字符串
func Decrypt(value, secret string) (decrypt string) {
	parts := strings.SplitN(value, ".", 3)
	if len(parts) != 3 {
		return ""
	}
	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]
	h := hmac.New(sha1.New, []byte(secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)
	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return ""
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res)
}

//MD5加密
//@param			str			需要加密的字符串
//@param			salt		盐值
//@return			CryptStr	加密后返回的字符串
func Md5Crypt(str string, salt ...interface{}) (CryptStr string) {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

//SHA1加密
//@param			str			需要加密的字符串
//@param			salt		盐值
//@return			CryptStr	加密后返回的字符串
func Sha1Crypt(str string, salt ...interface{}) (CryptStr string) {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}


//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5Crypt(base64.URLEncoding.EncodeToString(b))
}