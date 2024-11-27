package token

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"example/config"
	"io"
)

const (
	aesKey = "1234567887654321"
)

type Token struct {
	UserID int
	IP     string
}

func NewToken(userID int, ip string) config.IToken {
	return &Token{UserID: userID, IP: ip}
}

func (o Token) GetUserID() int {
	return o.UserID
}

func (o Token) EnToken() (string, error) {
	bs, err := json.Marshal(&o)
	if err != nil {
		return "", err
	}
	result, err := encrypt(bs, []byte(aesKey))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(result), nil
}

func (o *Token) DeToken(token string) error {
	if token == "" || token == "undefined" {
		return errors.New("need token")
	}
	tmp, err := hex.DecodeString(token)
	if err != nil {
		return errors.New("token error")
	}
	de, err := decrypt(tmp, []byte(aesKey))
	if err != nil {
		return errors.New("token error")
	}
	if err = json.Unmarshal(de, o); err != nil {
		return errors.New("token error")
	}
	if o.UserID == 0 {
		return errors.New("please login")
	}
	return nil
}

func encrypt(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	return ciphertext, nil
}

func decrypt(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}
