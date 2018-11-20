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

var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

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

//生成指定长度的字符串.
func NewRandChars(length int) string {
	if length == 0 {
		return ""
	}
	clen := len(stdChars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = stdChars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
