package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decodeBase64(data []byte) []byte {
	b, _ := base64.StdEncoding.DecodeString(string(data))
	return b
}

func getLocalKey() string {
	env := os.Getenv("ENIGMA_KEY")
	if env == "" {
		env = "thebrownfoxjumpedoverthefence"
	}

	return env
}

func getSizedKey(key string) string {
	// Get the correct key length
	l := len(key)
	if l < 16 {
		for i := 0; i < 16-l; i++ {
			key += "."
		}
	} else if l < 24 {
		for i := 0; i < 24-l; i++ {
			key += "."
		}
	} else if l < 32 {
		for i := 0; i < 32-l; i++ {
			key += "."
		}
	} else {
		key = key[:32]
	}

	return key
}

// Encrypt the password book
func encrypt(text, passphrase string) ([]byte, error) {
	key := []byte(getSizedKey(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipherText := make([]byte, aes.BlockSize+len(text))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(cipherText[aes.BlockSize:], []byte(text))
	return cipherText, nil
}

// Decrypt the password book
func decrypt(text []byte, passphrase string) ([]byte, error) {
	key := []byte(getSizedKey(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("Cipher text too short")
	}
	iv := text[:aes.BlockSize]
	data := text[aes.BlockSize:]
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(data, data)
	return data, nil
}
