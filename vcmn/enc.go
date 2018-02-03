package vcmn

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/varunamachi/vaali/vlog"
)

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}
	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, errors.New("unpad error. " +
			"This could happen when incorrect encryption key is used")
	}
	return src[:(length - unpadding)], nil
}

//Encrypt - encrypts input text with given key using AES algo
func Encrypt(key []byte, text string) (encrypted string, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err == nil {
		msg := pad([]byte(text))
		ciphertext := make([]byte, aes.BlockSize+len(msg))
		iv := ciphertext[:aes.BlockSize]
		if _, err = io.ReadFull(rand.Reader, iv); err == nil {
			cfb := cipher.NewCFBEncrypter(block, iv)
			cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
			encrypted = removeBase64Padding(
				base64.URLEncoding.EncodeToString(ciphertext))
		}
	}
	return encrypted, err
}

//Decrypt - decrypts input text with given key using AES algo
func Decrypt(key []byte, text string) (decrypted string, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err == nil {
		var decodedMsg []byte
		decodedMsg, err = base64.URLEncoding.DecodeString(
			addBase64Padding(text))
		if err == nil {
			if (len(decodedMsg) % aes.BlockSize) == 0 {
				iv := decodedMsg[:aes.BlockSize]
				msg := decodedMsg[aes.BlockSize:]
				cfb := cipher.NewCFBDecrypter(block, iv)
				cfb.XORKeyStream(msg, msg)
				var unpadMsg []byte
				unpadMsg, err = unpad(msg)
				if err == nil {
					decrypted = string(unpadMsg)
				}
			} else {
				err = errors.New(
					"Blocksize must be multipe of decoded message length")
			}
		}
	}
	return decrypted, vlog.LogError("Cmn:Enc", err)
}

//EncryptStr - AES encrypts input text with given key string
func EncryptStr(key string, text string) (encrypted string, err error) {
	encrypted, err = Encrypt([]byte(key), text)
	return encrypted, err
}

//DecryptStr - AES decrypts input text with given key string
func DecryptStr(key string, text string) (decrypted string, err error) {
	decrypted, err = Decrypt([]byte(key), text)
	return decrypted, err
}

//Hash - creates a SHA1 hash of input string
func Hash(in string) (out string) {
	hasher := sha1.New()
	hasher.Write([]byte(in))
	out = fmt.Sprintf("%x", hasher.Sum(nil))
	return out
}
