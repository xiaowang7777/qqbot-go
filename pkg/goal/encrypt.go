package goal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"qqbot-go/utils"
)

//GenerateRSAKey 生成RSA公钥和私钥文件
func GenerateRSAKey(dir, rsaPrivateFileName, rsaPublicFileName string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	rsaPrivateFilePath := filepath.Join(dir, rsaPrivateFileName)
	rsaPublicFilePath := filepath.Join(dir, rsaPublicFileName)

	err = utils.CreateIfNotExists(rsaPrivateFilePath)
	if err != nil {
		return err
	}
	err = utils.CreateIfNotExists(rsaPublicFilePath)
	if err != nil {
		return err
	}

	privateFile, err := os.OpenFile(rsaPrivateFilePath, os.O_WRONLY|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	publicFile, err := os.OpenFile(rsaPublicFilePath, os.O_WRONLY|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	publicKey := privateKey.PublicKey
	pkcs1PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}

	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: pkcs1PublicKey}
	err = pem.Encode(publicFile, &publicBlock)
	if err != nil {
		return err
	}

	pkcs1PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := &pem.Block{Type: "RSA Private Key", Bytes: pkcs1PrivateKey}
	err = pem.Encode(privateFile, privateBlock)
	if err != nil {
		return err
	}
	return nil
}

//RSAEncrypt 公钥加密
func RSAEncrypt(r io.Reader, size int64, msg []byte) ([]byte, error) {
	byteArr := make([]byte, size)
	if _, err := r.Read(byteArr); err != nil {
		return nil, err
	}
	block, _ := pem.Decode(byteArr)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := key.(*rsa.PublicKey)
	if encryptRes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg); err != nil {
		return nil, err
	} else {
		return encryptRes, nil
	}
}

//RSADecrypt RSA私钥解密
func RSADecrypt(r io.Reader, size int64, msg []byte) ([]byte, error) {
	byteArr := make([]byte, size)
	if _, err := r.Read(byteArr); err != nil {
		return nil, err
	}
	block, _ := pem.Decode(byteArr)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if decryptRes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, msg); err != nil {
		return nil, err
	} else {
		return decryptRes, nil
	}
}
