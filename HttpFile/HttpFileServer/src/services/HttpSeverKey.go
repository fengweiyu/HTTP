package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
// kUserTokenUrl = rsurl + "/login/va1/%s/%s.rs"
)

type usertokenResponse struct {
	Msg  string   `json:"msg"`
	Code int      `json:"code"`
	Data userdata `json:"data"`
}
type userdata struct {
	AccessToken string `json:"accessToken"`
}

// generate key
func KeyFilter(timeMillis, appSecret string) string {
	var key string
	if len(timeMillis)/2 > 0 {
		key = timeMillis[len(timeMillis)/2:]
	}
	key += appSecret
	length := len(key)
	if length > 16 {
		key = key[:16]
	} else {
		for i := 0; i < 16-length; i++ {
			key += strconv.Itoa(i)
		}
	}

	return key
}

func AesEncryption(src, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	if src == "" {
		return "", fmt.Errorf("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	encodeData := hex.EncodeToString(crypted)
	encodeData = strings.ToUpper(encodeData)
	return encodeData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
