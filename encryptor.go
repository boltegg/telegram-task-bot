package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

type Encryptor struct {
	key          []byte
	errorHandler func(error)
}

func NewEncryptor(passphrase string, errorHandler func(error)) *Encryptor {

	key := sha256.Sum256([]byte(passphrase))
	return &Encryptor{key: []byte(key[:]), errorHandler: errorHandler}
}

func (e *Encryptor) EncryptString(src string) string {
	enc, err := e.encrypt([]byte(src))
	if err != nil {
		e.errorHandler(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(enc)
}

func (e *Encryptor) DecryptString(src string) string {
	b, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		e.errorHandler(err)
		return ""
	}

	dec, err := e.decrypt(b)
	if err != nil {
		e.errorHandler(err)
		return ""
	}

	return string(dec)
}

func (e *Encryptor) encrypt(data []byte) ([]byte, error) {
	block, _ := aes.NewCipher(e.key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *Encryptor) decrypt(data []byte) ([]byte, error) {

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
