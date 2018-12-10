// Copyright 2018. bolaxy.org authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func GenerateRandomRsaKeyPair() (prvBytes, pubBytes []byte) {
	return ExportKey(GenRsaKey())
}

func GenRsaKey() *rsa.PrivateKey {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	return priv
}

func ExportKey(key *rsa.PrivateKey) (prvBytes, pybBytes []byte) {
	if key == nil {
		return nil, nil
	}

	priv_byte := x509.MarshalPKCS1PrivateKey(key)
	priv_bs64 := base64.StdEncoding.EncodeToString(priv_byte)

	// 公钥
	pub_byte := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	pub_bs64 := base64.StdEncoding.EncodeToString(pub_byte)

	return []byte(priv_bs64), []byte(pub_bs64)
}

func RsaEncrypt(pubkey, msg []byte) ([]byte, error) {
	bKey, err := base64.StdEncoding.DecodeString(string(pubkey))
	if err != nil {
		return nil, err
	}
	pub, err := x509.ParsePKCS1PublicKey(bKey)
	if err != nil {
		return nil, err
	}
	data, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	if err != nil {
		return nil, err
	}
	base64Data := []byte(base64.StdEncoding.EncodeToString(data))
	return base64Data, nil
}

func RsaDecrypt(priv, msg []byte) ([]byte, error) {
	bKey, err := base64.StdEncoding.DecodeString(string(priv))
	if err != nil {
		return nil, err
	}
	priv_key, err := x509.ParsePKCS1PrivateKey(bKey)
	if err != nil {
		return nil, err
	}

	msg_byte, err := base64.StdEncoding.DecodeString(string(msg))
	if err != nil {
		return nil, err
	}
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv_key, msg_byte)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func RsaSign(priv, msg []byte) ([]byte, error) {
	bKey, err := base64.StdEncoding.DecodeString(string(priv))
	if err != nil {
		return nil, err
	}
	priv_key, err := x509.ParsePKCS1PrivateKey(bKey)
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256(msg)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv_key, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	signature_bs64 := []byte(base64.StdEncoding.EncodeToString(signature))
	return signature_bs64, nil
}

func RsaVerify(pubKeyBytes, msg, signature []byte) error {
	bSignData, err := base64.StdEncoding.DecodeString(string(signature))
	if err != nil {
		return fmt.Errorf("DecodeString sign error:%v", err)
	}
	hashed := sha256.Sum256(msg)
	bKey, err := base64.StdEncoding.DecodeString(string(pubKeyBytes))
	if err != nil {
		return fmt.Errorf("DecodeString pubKeyBytes  error:%v", err)
	}
	pub, err := x509.ParsePKCS1PublicKey(bKey)
	if err != nil {
		return fmt.Errorf("ParsePKCS1PublicKey error:%v", err)
	}
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], bSignData)
	if err != nil {
		return fmt.Errorf("VerifyPKCS1v15 error:%v", err)
	}
	return nil
}

//func RsaVerify(pubKeyBytes, msg, signature []byte) error {
//	bSignData, err := base64.StdEncoding.DecodeString(string(signature))
//	if err != nil {
//		return fmt.Errorf("DecodeString sign error:%v", err)
//	}
//	hashed := sha256.Sum256(msg)
//	bKey, err := base64.StdEncoding.DecodeString(string(pubKeyBytes))
//	if err != nil {
//		return fmt.Errorf("DecodeString pubKeyBytes  error:%v", err)
//	}
//	pub, err := x509.ParsePKCS1PublicKey(bKey)
//	if err != nil {
//		return fmt.Errorf("ParsePKCS1PublicKey error:%v", err)
//	}
//	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], bSignData)
//	if err != nil {
//		return fmt.Errorf("VerifyPKCS1v15 error:%v", err)
//	}
//	return nil
//}
