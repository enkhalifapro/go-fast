package utilities

import (
	"crypto/rand"
	"math/big"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type ICryptUtil interface {
	RandomString(n int) string
	Encrypt(str string) string
	NewEncryptedToken() string
	Bcrypt(str string) string
	CompareHashAndPassword(hashedPassword string, password string) bool
}

type CryptUtil struct {

}

func NewCryptUtil() ICryptUtil {
	util := CryptUtil{}
	return &util
}

func (util CryptUtil) RandomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols := big.NewInt(int64(len(alphanum)))
	states := big.NewInt(0)
	states.Exp(symbols, big.NewInt(int64(n)), nil)
	r, err := rand.Int(rand.Reader, states)
	if err != nil {
		panic(err)
	}
	var bytes = make([]byte, n)
	r2 := big.NewInt(0)
	symbol := big.NewInt(0)
	for i := range bytes {
		r2.DivMod(r, symbols, symbol)
		r, r2 = r2, r
		bytes[i] = alphanum[symbol.Int64()]
	}
	return string(bytes)
}

func (util CryptUtil) Encrypt(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	passwordSha256Hash := hex.EncodeToString(h.Sum(nil))
	return passwordSha256Hash
}

func (util CryptUtil) Bcrypt(str string) string {
	byteStr := []byte(str)
	hashedPassword, err := bcrypt.GenerateFromPassword(byteStr, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func (util CryptUtil) CompareHashAndPassword(hashedPassword string, password string) bool {
	hashedPasswordByte := []byte(hashedPassword)
	passwordByte := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedPasswordByte, passwordByte)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

func (util CryptUtil) NewEncryptedToken() string {
	randomStr := util.RandomString(100)
	token := util.Encrypt(randomStr)
	return token
}