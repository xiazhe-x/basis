package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

/*json  操作*/
var jsonApi jsoniter.API
var oJson sync.Once

func getJsonIter() jsoniter.API {
	oJson.Do(func() {
		jsonApi = jsoniter.ConfigCompatibleWithStandardLibrary
	})
	return jsonApi
}

func DataByJsonStr(params interface{}) string {
	by, err := getJsonIter().Marshal(params)
	if err != nil {
		logrus.Error("DataByJsonStr err ", err)
	}
	return string(by)
}

func StrJsonByData(str string, data interface{}) {
	err := getJsonIter().Unmarshal([]byte(str), data)
	if err != nil {
		logrus.Error("StrJsonByData err ", err)
	}
}

func ByteJsonByData(by []byte, data interface{}) error {
	err := getJsonIter().Unmarshal(by, data)
	if err != nil {
		return err
	}
	return nil
}

func RandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func RandForInt(max int, n int) int {
	rand.Seed(time.Now().UnixNano() + int64(n))
	return rand.Intn(max)
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

const authCode = "2b0f949a0d1"

//密码MD5方法
func PwdMD5(str string) string {
	_md5 := md5.New()
	_md5.Write([]byte(authCode + str))
	return hex.EncodeToString(_md5.Sum([]byte("")))
}
func HashFileMd5(file multipart.File) (md5Str string) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}
	hashInBytes := hash.Sum(nil)[:16]
	md5Str = hex.EncodeToString(hashInBytes)
	md5Str = fmt.Sprintf("%s_%d", md5Str, time.Now().Unix())
	return
}

// hmac 加密
func ComputeHmacSha256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

var CHARS = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

/*RandString  生成随机字符串(a~zA~Z])
  lenNum 长度
*/
func RandString(lenNum int) string {
	str := strings.Builder{}
	length := 52
	for i := 0; i < lenNum; i++ {
		str.WriteString(CHARS[rand.Intn(length)])
	}
	return str.String()
}

/*RandNumString  生成随机数字字符串([0~9])
  lenNum 长度
*/
func RandNumString(lenNum int) string {
	str := strings.Builder{}
	length := 10
	for i := 0; i < lenNum; i++ {
		str.WriteString(CHARS[52+rand.Intn(length)])
	}
	return str.String()
}

func RandAllString(lenNum int) string {
	str := strings.Builder{}
	length := len(CHARS)
	for i := 0; i < lenNum; i++ {
		l := CHARS[rand.Intn(length)]
		str.WriteString(l)
	}
	return str.String()
}

//计算当天还剩多少秒
func SameDaySurplusSecond() int64 {
	now := time.Now()
	endAt, _ := time.ParseInLocation("2006-01-02 15:04:05", now.Format("2006-01-02")+" 23:59:59", time.Local)
	return endAt.Unix() - now.Unix()
}

//------------------------------------------------		其他
func SliIndex(data []int, value int) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}
	return false
}

func DataByRedisMap(params interface{}) map[string]interface{} {
	//	结构体序列化反序列化map 你也可以用其他包进行转化为map,但是注意转化后的键名大小写问题
	m := make(map[string]interface{})
	buf, _ := getJsonIter().Marshal(params)
	_ = getJsonIter().Unmarshal(buf, &m)
	return m
}
