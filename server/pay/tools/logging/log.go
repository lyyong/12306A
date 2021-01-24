// @Author liuYong
// @Created at 2020-01-05
// @Modified at 2020-01-05
package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	logger        *log.Logger
	levelFlags    = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	logPrefix     = ""
	DefaultPrefix = ""

	DefaultCallerDepth = 2
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	logger = log.New(os.Stdout, DefaultPrefix, log.LstdFlags)

}

func setPrefix(level Level) {
	// 返回当前程序栈中的文件和执行的行数
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}
