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
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"
	"math/rand"
)

const (
	AESKEY_LEN_16 = 16
	AESKEY_LEN_24 = 24
	AESKEY_LEN_32 = 32
)

func padToBlockSize(payload []byte, blockSize int) (prefix, finalBlock []byte) {
	overrun := len(payload) % blockSize
	paddingLen := blockSize - overrun
	prefix = payload[:len(payload)-overrun]
	finalBlock = make([]byte, blockSize)
	copy(finalBlock, payload[len(payload)-overrun:])
	for i := overrun; i < blockSize; i++ {
		finalBlock[i] = byte(paddingLen)
	}
	return
}

// Encrypter
func AESCBCEncrypter(data, cbcKey []byte) ([]byte, error) {

	key := make([]byte, len([]byte(cbcKey)))
	copy(key, []byte(cbcKey))
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("NewCipher error:", err)
		return nil, err
	}

	// set random IV len : aes.BlockSize(16)
	iv := GetAesKeyRandom(aes.BlockSize)
	encrypter := cipher.NewCBCEncrypter(c, iv)

	//规整数据长度
	prefix, finalBlock  := padToBlockSize(data,encrypter.BlockSize())
	src := append(prefix,finalBlock...)
	length := len(src)

	enData := make([]byte, length)
	encrypter.CryptBlocks(enData, src)
	encBase64Data := []byte(base64.StdEncoding.EncodeToString(enData))

	encBase64Data = append(iv, encBase64Data...)

	return encBase64Data, nil
}

// Decrypter
func AESCBCDecrypter(data, cbcKey []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("data is too short.")
	}
	iv := data[:aes.BlockSize]
	deData, err := base64.StdEncoding.DecodeString(string(data[aes.BlockSize:]))
	if err != nil {
		fmt.Errorf("base64  Decode error:", err)
		return nil, err
	}

	//iv := deDataWithIv[:aes.BlockSize]
	//deData := deDataWithIv[:]
	length := len(deData)

	c, err := aes.NewCipher(cbcKey)
	if err != nil {
		fmt.Errorf("NewCipher[%s] err: %s", cbcKey, err)
		return nil, err
	}

	deCrypter := cipher.NewCBCDecrypter(c, iv)
	// check block size
	if len(deData)%deCrypter.BlockSize() != 0 {
		return nil,fmt.Errorf("crypto/cipher: input not full blocks")
	}
	deCrypter.CryptBlocks(deData, deData)
	outData := deData[:length-int(deData[length-1])]

	return outData, nil
}

// get aes key random
func GetAesKeyRandom(len int) []byte {
	return []byte(GenerateRandomString(len))
}

func GenerateRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}