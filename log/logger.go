package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// IsDebug 是否输出调试日志
var IsDebug = false

// Debug 输出调试信息
func Debug(format string, a ...interface{}) {
	if !IsDebug {
		return
	}
	color.White("[debug] "+format, a...)
}

func Info(format string, a ...interface{}) {
	color.Cyan(format, a...)
}

func Warn(format string, a ...interface{}) {
	color.Yellow(format, a...)
}

func Error(format string, a ...interface{}) {
	color.Red(format, a...)
}

func Success(format string, a ...interface{}) {
	color.Green(format, a...)
}

func Fatal(err error) {
	color.Red(err.Error())
	os.Exit(1)
}

// PROGRESS_WIDTH 进度条长度
const PROGRESS_WIDTH = 40

// DrawProgressBar 显示进度条
func DrawProgressBar(prefix string, val, max int) {
	proportion := float32(val) / float32(max)
	pos := int(proportion * PROGRESS_WIDTH)
	s := fmt.Sprintf("%s [%s%*s] %6.2f%% \t[%d/%d]",
		prefix, strings.Repeat("■", pos), PROGRESS_WIDTH-pos, "", proportion*100, val, max)
	fmt.Print(color.CyanString("\r" + s))
	if proportion >= 1 {
		fmt.Print("\n")
	}
}
