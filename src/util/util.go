package util

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/szxby/tools/log"
)

const (
	key = "7inrmpd5DSQTfDxnAnOH"
)

// CalculateHash calculate hash
func CalculateHash(data string) string {
	h := sha256.New()
	h.Write([]byte(key + data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// RandomString return a random string with length len
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// CheckDir 检查目录是否存在，不存在创建
func CheckDir(dir string) {
	local, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error("get local fail %v", err)
		return
	}
	_, err = os.Stat(local + dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(local+dir, os.ModePerm); err != nil {
			log.Error("make dir fail %v", err)
		}
	}
}

// Decimal 保留2位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, v1 := range maps {
		for k2, v2 := range v1 {
			ret[k2] = v2
		}
	}
	return ret
}

// FormatFloat 取小数点后n位非零小数
func FormatFloat(num float64, decimal int) string {
	if math.Trunc(num) == num || decimal == 0 {
		return fmt.Sprintf("%.f", math.Trunc(num))
	}
	format := "%." + strconv.Itoa(decimal) + "f"
	return fmt.Sprintf(format, num)
}

// GetFirstDateOfMonth 获取本月第一天零点
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// GetLastDateOfMonth 获取本月最后一天零点
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
