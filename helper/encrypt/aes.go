package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// DefaultPepper .
var DefaultPepper = Pepper("This is a default pepper, you should change it in deployment environment.")
var pepper = DefaultPepper

// SetPepper set the default pepper to the give key.
func SetPepper(key string) { pepper = Pepper(key) }

// Pepper would padding the give key to multiple of 8, minimum is 16 and maximum is 32.
// (k < 16 => k = 16, 16 < k < 24 => k = 24, 24 < k < 32 => k = 32, 32 < k => k = 32)
func Pepper(key string) []byte {
	l := len(key)
	var pepper []byte
	if l < 16 {
		pepper = pkcs7Padding([]byte(key), 16)
	} else if l > 32 {
		pepper = []byte(key)[:32]
	} else {
		pepper = pkcs7Padding([]byte(key), 8)
	}
	return pepper
}

// AesEncrry .
func AesEncrry(plain, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()

	plain = pkcs7Padding(plain, bs)
	secret := make([]byte, len(plain))
	bm := cipher.NewCBCEncrypter(block, key[:bs])
	bm.CryptBlocks(secret, plain)
	return secret, nil
}

// AesDecrypt .
func AesDecrypt(secret, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()

	bm := cipher.NewCBCDecrypter(block, key[:bs])
	plain := make([]byte, len(secret))
	bm.CryptBlocks(plain, secret)
	plain = pkcs7UnPadding(plain, bs)
	return plain, nil
}

func pkcs7Padding(plain []byte, blockSize int) []byte {
	padLen := blockSize - len(plain)%blockSize
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(plain, padText...)
}

func pkcs7UnPadding(padded []byte, blockSize int) []byte {
	length := len(padded)
	if length < 1 {
		return nil
	}
	padLen := int(padded[length-1])
	if padLen > length {
		// fmt.Println(padded, padLen, length)
		return nil
	}
	return padded[:(length - padLen)]
}
