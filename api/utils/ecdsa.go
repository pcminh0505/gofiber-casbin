package utils

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func LoadEcdsaPrivateKeyKey() *ecdsa.PrivateKey {
	privateKeyFile, err := os.Open("private.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	privateKeyFile.Close()

	// data, _ := pem.Decode([]byte(pembytes))

	privateKeyImported, err := jwt.ParseECPrivateKeyFromPEM(pembytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return privateKeyImported
}
