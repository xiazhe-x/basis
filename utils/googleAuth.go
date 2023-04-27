package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (t *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (t *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (t *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (t *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (t *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (t *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (t *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := t.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := t.toUint32(hashParts)
	return number % 1000000
}

// 获取秘钥
func (t *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, t.un())
	return strings.ToUpper(t.base32encode(t.hmacSha1(buf.Bytes(), nil)))
}

// 获取动态码
func (t *GoogleAuth) GetCode(secret string,offset int64) string {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := t.base32decode(secretUpper)
	if err != nil {
		return ""
	}
	number := t.oneTimePassword(secretKey, t.toBytes((time.Now().Unix()+ offset)/30))
	return fmt.Sprintf("%06d", number)
}

// 获取动态码二维码内容
func (t *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

// 获取动态码二维码图片地址,这里是第三方二维码api
func (t *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := t.GetQrcode(user, secret)
	width := "200"
	height := "200"
	data := url.Values{}
	data.Set("data", qrcode)
	return "https://api.qrserver.com/v1/create-qr-code/?" + data.Encode() + "&size=" + width + "x" + height + "&ecc=M";
}

// 验证动态码
func (t *GoogleAuth) VerifyCode(secret, code string) bool {
	//_code, err := t.GetCode(secret)
	if t.GetCode(secret, 0) == code {
		return true
	}
	if t.GetCode(secret, -30) == code {
		return true
	}
	if t.GetCode(secret, 30) == code {
		return true
	}
	//fmt.Println(_code, code)
	return false
}