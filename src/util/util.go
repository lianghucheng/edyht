package util

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

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
	_, err = os.Stat(local + MatchIconDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(local+MatchIconDir, os.ModePerm); err != nil {
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
