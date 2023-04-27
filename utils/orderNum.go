package utils

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

var num int64

//生成24位订单号
//前面17位代表时间精确到毫秒，中间3位代表进程id，最后4位代表序号
func Generate(t time.Time) string {
	// s := t.Format(timeformat.Continuity)
	s := t.Format("02150405")
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}

//对长度不足n的数字前面补0
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

//------------------------------------------------		签名
func SignStr(m map[string]interface{}, signKey string) string {
	keys := make([]string, 0, len(m))
	for k, v := range m {
		if v == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	v := url.Values{}
	for _, k := range keys {
		va := fmt.Sprintf("%v", m[k])
		if va == "" || va == "0" {
			continue
		}
		v.Add(k, va)
	}
	//转成QueryUnescape
	body := v.Encode()
	de, _ := url.QueryUnescape(body)

	str := fmt.Sprintf("%s&key=%s", de, signKey)
	//再将得到的字符串所有字符转换为大写，得到sign值signValue
	//fmt.Println("签名字符串>", str)
	sign := strings.ToUpper(MD5([]byte(str)))
	//fmt.Println("签名结果>", sign)
	return sign
}

//换算金额
func ChByAmount(currency string, amount int) (i int, err error) {
	switch currency {
	case "CNY":
		if amount < 10 {
			err = errors.New("支付金额过小")
		}
		i = amount * 100
	case "VND":
		if amount < 350 {
			err = errors.New("支付金额过小")
		}
		i = amount / 100
	case "USD":
		if amount < 2 {
			err = errors.New("支付金额过小")
		}
		i = amount * 100
	}
	return
}

//商户汇率时币种转换
func MchByAmount(currency string, amount uint) (i float64) {
	switch currency {
	case "CNY":
		i = float64(amount) / 100
	case "VND":
		i = float64(amount)
	case "USD":
		i = float64(amount) / 100
	}
	return
}
