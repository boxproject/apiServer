// Copyright 2018. bolaxy.org authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package middleware

import (
	"errors"
	"fmt"
	"time"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	err "github.com/boxproject/apiServer/errors"
	reterror "github.com/boxproject/apiServer/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

const TOKEN_EXP = 24 // token有效周期

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		//非法token定义
		msgtoken := err.ErrModel{Code: reterror.ErrToken, Err: reterror.MSG_ErrToken}
		if strings.ToUpper(c.Request.Method) == "GET" {
			token = c.Query("token")
		}
		if strings.ToUpper(c.Request.Method) == "POST" {
			token = c.PostForm("token")
		}
		if token == "" {
			log.Error("JWTAuth")
			c.Abort()
			msgtoken.RetErr(c)
			return
		}

		j := NewJWT()
		claims, err := j.ParseToken(token)
		if claims == nil || claims.AppID == "" {
			log.Error("JWTAuth")
			c.Abort()
			msgtoken.RetErr(c)
			return
		}

		if err != nil {
			c.Abort()
			log.Error("token", err.Error())
			msgtoken.RetErr(c)
			return
		}

		appid := claims.AppID
		// 校验账号
		acc, msg := validateUser(appid)

		if msg.Code != 0 {
			log.Error("JWTAuth")
			c.Abort()
			msgtoken.RetErr(c)
			return
		}
		fmt.Println("jwt-------------", acc.UserType, claims.UserType)
		if acc.UserType != claims.UserType {
			log.Error("员工角色改变，重新登陆")
			c.Abort()
			msgtoken.RetErr(c)
			return
		}
		claims.PubKey = acc.PubKey
		claims.ID = acc.ID
		claims.Level = acc.Level
		claims.UserType = acc.UserType
		c.Set("claims", claims)

	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "box.la"
)

type CustomClaims struct {
	ID       int
	AppID    string `json:"appid"`
	Account  string `json:"account"`
	UserType int    `json:"userType"`
	PubKey   string
	Level    int
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}
func GetSignKey() string {
	return SignKey
}
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil || token.Valid == false {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Add(TOKEN_EXP * time.Hour).Unix()
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}

// 校验账号
func validateUser(appid string) (*db.Account, reterror.ErrModel) {
	as := db.AccDBService{}
	accInfo, err := as.FindAccountByAppId(appid)
	if err != nil {
		log.Debug("find account error", appid)
		return nil, reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if accInfo == nil {
		log.Debug("account not found", appid)
		return nil, reterror.ErrModel{Code: reterror.User_3004, Err: reterror.MSG_3004}
	}

	if accInfo.IsDeleted == 1 {
		return nil, reterror.ErrModel{Code: reterror.User_3005, Err: reterror.MSG_3005}
	}
	if accInfo.Frozen == 1 {
		nowTime := time.Now().Unix()
		timeStep := accInfo.FrozenTo.Unix() - nowTime
		if timeStep <= 0 {
			// 重置账户状态
			as.ResetAccount(appid)
			accInfo.Frozen = 0
			return accInfo, reterror.ErrModel{Err: nil, Code: reterror.Success}
		}

		return nil, reterror.ErrModel{Code: reterror.User_3006, Err: reterror.MSG_3006, Data: map[string]int64{"frozenTo": accInfo.FrozenTo.Unix()}}
	}
	return accInfo, reterror.ErrModel{Err: nil, Code: reterror.Success}
}
