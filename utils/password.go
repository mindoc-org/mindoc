package utils


import (
	"crypto/rand"
	mt "math/rand"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strconv"
	"strings"
)

const (
	saltSize            = 16
	delmiter            = "$"
	stretching_password = 500
	salt_local_secret   = "ahfw*&TGdsfnbi*^Wt"
)
//加密密码
func PasswordHash(pass string) (string, error) {

	salt_secret, err := salt_secret()
	if err != nil {
		return "", err
	}

	salt, err := salt(salt_local_secret + salt_secret)
	if err != nil {
		return "", err
	}

	interation := randInt(1, 20)

	hash, err := hash(pass, salt_secret, salt, int64(interation))
	if err != nil {
		return "", err
	}
	interation_string := strconv.Itoa(interation)
	password := salt_secret + delmiter + interation_string + delmiter + hash + delmiter + salt

	return password, nil

}
//校验密码是否有效
func  PasswordVerify(hashing string, pass string) (bool, error) {
	data := trim_salt_hash(hashing)

	interation, _ := strconv.ParseInt(data["interation_string"], 10, 64)

	has, err := hash(pass, data["salt_secret"], data["salt"], int64(interation))
	if err != nil {
		return false, err
	}

	if (data["salt_secret"]+delmiter+data["interation_string"]+delmiter+has+delmiter+data["salt"]) == hashing {
		return true, nil
	} else {
		return false, nil
	}

}

func hash(pass string, salt_secret string, salt string, interation int64) (string, error) {
	var pass_salt string = salt_secret + pass + salt + salt_secret + pass + salt + pass + pass + salt
	var i int

	hash_pass := salt_local_secret
	hash_start := sha512.New()
	hash_center := sha256.New()
	hash_output := sha256.New224()

	i = 0
	for i <= stretching_password {
		i = i + 1
		hash_start.Write([]byte(pass_salt + hash_pass))
		hash_pass = hex.EncodeToString(hash_start.Sum(nil))
	}

	i = 0
	for int64(i) <= interation {
		i = i + 1
		hash_pass = hash_pass + hash_pass
	}

	i = 0
	for i <= stretching_password {
		i = i + 1
		hash_center.Write([]byte(hash_pass + salt_secret))
		hash_pass = hex.EncodeToString(hash_center.Sum(nil))
	}
	hash_output.Write([]byte(hash_pass + salt_local_secret))
	hash_pass = hex.EncodeToString(hash_output.Sum(nil))

	return hash_pass, nil
}

func trim_salt_hash(hash string) map[string]string {
	str := strings.Split(hash, delmiter)

	return map[string]string{
		"salt_secret":       str[0],
		"interation_string": str[1],
		"hash":              str[2],
		"salt":              str[3],
	}
}
func salt(secret string) (string, error) {

	buf := make([]byte, saltSize, saltSize + md5.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(buf)
	hash.Write([]byte(secret))
	return hex.EncodeToString(hash.Sum(buf)), nil
}

func salt_secret() (string, error) {
	rb := make([]byte, randInt(10, 100))
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

func randInt(min int, max int) int {
	return min + mt.Intn(max-min)
}