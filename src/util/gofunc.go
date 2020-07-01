package util

import (
	"runtime"

	"github.com/szxby/tools/log"
)

// Gofunc create a goroutine to invoke f
func Gofunc(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Error("goroutine err:%v", err)
				log.Error("%s", buf)
			}
		}()
		f()
	}()
}
