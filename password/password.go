package password

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

const (
	Key = "rAk2PnNL17rd2utEx4zOmnP2ZfI4McSA"
	Iv  = "7ZfLrFhjAx2iyCBh"
)

// VerifyPassWorld ， 验证密码工具函数
// 参数：
//
//	rPass ：输入密码
//	uPass ： 原密码
//
// 返回值：
//
//	bool ：是否验证通过

func VerifyPassWorld(rPass, uPass string) bool {
	infos := strings.Split(uPass, "$")
	if len(infos) != 4 {
		return false
	}
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	if ok := password.Verify(rPass, infos[2], infos[3], options); !ok {
		return false
	}
	return true
}

func Verify2(rPass, uPass string) bool {
	return rPass == uPass
}

// GeneratePassWorld ， 生成密码
// 参数：
//
//	str ： 密码字符串
//
// 返回值：
//
//	string ：密码
func GeneratePassWorld(str string) string {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encryptedPwd := password.Encode(str, options)
	encryptedPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encryptedPwd)
	return encryptedPassword
}

// Padding 对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

// UnPadding 对密文删除填充
func UnPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}

// AesCbcEncrypt AEC加密（CBC模式）
func AesCbcEncrypt(plainText, key, iv []byte) []byte {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = Padding(plainText, block.BlockSize())
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return cipherText
}

/*
	//eg : kv
	//    key: "rAk2PnNL17rd2utEx4zOmnP2ZfI4McSA"
	//    iv: "7ZfLrFhjAx2iyCBh"
*/
// AesCbcDecrypt2 AEC解密（CBC模式）
func AesCbcDecrypt2(cipherText, key, iv []byte) (info []byte, errs error) {
	defer func() {
		if err := recover(); err != nil {
			errs = errors.New("密文解析异常，请确认字段是否加密！")
			return
		}
	}()
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//指定初始化向量IV,和加密的一致
	//iv := []byte("12345678abcdefgh")
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	//删除填充
	plainText = UnPadding(plainText)

	return plainText, errs
}

// Encrypt 加密
func Encrypt(params, k, v string) string {
	encrypt := AesCbcEncrypt([]byte(params), []byte(k), []byte(v))
	encodeToString := base64.StdEncoding.EncodeToString(encrypt)
	return encodeToString
}

// Decrypt 解密
func Decrypt(params, k, v string) (string, error) {
	decodeString, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return "", err
	}
	decrypt2, err := AesCbcDecrypt2(decodeString, []byte(k), []byte(v))
	if err != nil {
		return "加密格式错误，解密失败", err
	}
	return string(decrypt2), nil
}
