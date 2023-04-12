package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
)

const (
	SecretKey = "c787ead050b555ca99c5f570a7d"
)

func ParseToken(s string) (*jwt.Token, error) {
	fn := func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	}
	return jwt.Parse(s, fn)
}

func CreateToken(Claims jwt.MapClaims) string {
	//jwt.SigningMethodRS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	if tokenString, err := token.SignedString([]byte(SecretKey)); err == nil {
		return tokenString
	} else {
		return ""
	}
}

func GetFromClaims(claims jwt.Claims) map[string]string {
	result := make(map[string]string)
	v := reflect.ValueOf(claims)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			value := v.MapIndex(k)
			result[fmt.Sprintf("%s", k.Interface())] = fmt.Sprintf("%v", value.Interface())
		}
	}
	return result
}
