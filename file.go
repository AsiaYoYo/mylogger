package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件里输出日志相关代码

// FileLogger 结构体
type FileLogger struct {
	Level       LogLevel
	filePath    string
	fileName    string
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64
}

// NewFileLogger 构造函数
func NewFileLogger(levelStr, filePath, fileName string, maxSize int64) *FileLogger {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	// 初始化FileLogger结构体
	fl := &FileLogger{
		Level:       level,
		filePath:    filePath,
		fileName:    filePath,
		maxFileSize: maxSize,
	}
	fl.initFile(filePath, fileName)
	return fl
}

// initFile 函数
func (f *FileLogger) initFile(fp, fn string) {
	pathName := path.Join(fp, fn)
	fileObj, err := os.OpenFile(pathName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开日志文件失败，err:%s", err)
	}
	errFileObj, err := os.OpenFile(pathName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开错误日志文件失败，err:%s", err)
	}
	f.fileObj = fileObj
	f.errFileObj = errFileObj
}

func (f *FileLogger) enable(logLevel LogLevel) bool {
	return logLevel >= f.Level
}

func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		fmt.Fprintf(f.fileObj, "%s [%s] %s:%s:%d %s\n", now.Format("2006-01-02 15:04:05"), parseLogLevelStr(lv), fileName, funcName, lineNo, msg)
		if lv >= ERROR {
			fmt.Fprintf(f.errFileObj, "%s [%s] %s:%s:%d %s\n", now.Format("2006-01-02 15:04:05"), parseLogLevelStr(lv), fileName, funcName, lineNo, msg)
		}
	}
}

// Debug ...
func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}

// Info ...
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}

// Warning ...
func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}

// Error ...
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}
