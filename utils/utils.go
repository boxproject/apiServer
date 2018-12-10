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
	"io/ioutil"
	"net/http"
	"strconv"
	logger "github.com/alecthomas/log4go"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"crypto/tls"
	"math/big"
	"reflect"
	"sort"
	"os"
	"path/filepath"
	"sync"
	"github.com/boxproject/apiServer/common"
)

// 发起https请求
func HttpRequest(method, urls string, jsonStr string) (r []byte, err error) {
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	payload := strings.NewReader(jsonStr)
	req, _ := http.NewRequest(method, urls, payload)

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	}
	c := &http.Client{Transport: tr}
	res, err := c.Do(req)

	if err != nil {
		logger.Error("请求代理服务器", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// 生成sha256哈希字符串
func GenHashStr(msg string) string {
	h := sha256.New()
	h.Write([]byte(msg))
	bs := h.Sum(nil)
	transHash := hex.EncodeToString(bs)
	return transHash
}

// 密码加盐
func PwdSalt(pwd string, salt string) string {
	h := sha256.New()
	h.Write([]byte(pwd))
	h.Write([]byte(salt))
	bs := h.Sum(nil)
	transHash := hex.EncodeToString(bs)
	return transHash
}

// 单位换算 amount/(10^factor)
func UnitConversion(amount float64, factor int64, base int64) (float64, error) {
	decimal := new(big.Int).Exp(big.NewInt(10), big.NewInt(factor), big.NewInt(0)).String()
	decimal_f, err := strconv.ParseFloat(decimal, 64)

	if err != nil {
		return 0, err
	}

	result := DivFloat64(amount, decimal_f)

	return result, nil
}

// 单位反换算 amount*(10^factor)
func UnitReConversion(amount float64, factor int64, base int64) (float64, error) {
	decimal := new(big.Int).Exp(big.NewInt(10), big.NewInt(factor), big.NewInt(0)).String()
	decimal_f, err := strconv.ParseFloat(decimal, 64)

	if err != nil {
		return 0, err
	}

	result := MulFloat64(amount, decimal_f)

	return result, nil
}

func GetTypeValues(data interface{}) (string, error) {
	var typeStr string
	var t reflect.Value
	t = reflect.ValueOf(data)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			msg := ToStr(t.Field(i).Interface())
			typeStr = typeStr + msg
		}
	} else if t.Kind() == reflect.Map {
		keys := t.MapKeys()
		for _, k := range keys {
			msg := ToStr(t.MapIndex(k))
			typeStr = typeStr + msg
		}
	} else if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		for i := 0; i < t.Len(); i++ {
			msg := ToStr(t.Index(i))
			typeStr = typeStr + msg
		}
	} else {
		typeStr = ToStr(t.Interface())
	}

	return typeStr, nil
}

func SortByValue(info interface{}) []byte {
	typeValues, err := GetTypeValues(info)
	if err != nil {
		logger.Error("GetTypeValues error:%v", err)
		return nil
	}
	return SortStrByRune(typeValues)
}

func SortStrByRune(str string) []byte {
	msg_rune := []rune(str)
	msgs := make([]string, len(msg_rune))
	for i := 0; i < len(msg_rune); i++ {
		msgs[i] = string(msg_rune[i])
	}
	sort.Strings(msgs)
	result := strings.Join(msgs, "")
	return []byte(result)
}

// 获取当前路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// 获取父级路径
func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

type MyMap struct {
	IsUsing map[string]int
	sync.RWMutex
}

func (this *MyMap) Set(key string, value int) {
	this.Lock()
	defer this.Unlock()
	this.IsUsing[key] = value
}
func (this *MyMap) Get(key string) int {
	this.RLock()
	defer this.RUnlock()
	return this.IsUsing[key]
}

// 转换语言
func ConvertLangTag(lang string) string {
	switch lang {
	case common.HeaderZhLang:
		return common.LangZhType
	case common.HeaderEnLang:
		return common.LangEnType
	default:
		return common.LangZhType
	}
}
