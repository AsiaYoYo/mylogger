package mylogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件里输出日志相关代码

var maxChanSize = 50000

// FileLogger 结构体
type FileLogger struct {
	Level       LogLevel
	filePath    string
	fileName    string
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64
	logChan     chan *logMsg
}

type logMsg struct {
	level     LogLevel
	timeStamp string
	fileName  string
	funcName  string
	msg       string
	lineNo    int
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
		logChan:     make(chan *logMsg, maxChanSize),
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
	// 开启一个后台goroutine用来异步写日志
	go f.writeLogBackground()
}

// enable 判断需要打印哪些级别的日志
func (f *FileLogger) enable(logLevel LogLevel) bool {
	return logLevel >= f.Level
}

// writeLogBackground 从通道中取日志，然后后台往文件里写日志
func (f *FileLogger) writeLogBackground() {
	for {
		select {
		// 从通道中去log
		case logTmp := <-f.logChan:
			logInfo := fmt.Sprintf("%s [%s] %s:%s:%d %s\n", logTmp.timeStamp,
				parseLogLevelStr(logTmp.level), logTmp.fileName, logTmp.funcName, logTmp.lineNo, logTmp.msg)
			fmt.Fprintf(f.fileObj, logInfo)
			if logTmp.level >= ERROR {
				fmt.Fprintf(f.errFileObj, logInfo)
			}
		// 取不到就跳过，先休息500毫秒
		default:
			time.Sleep(time.Millisecond * 500)
		}

	}
}

// log 处理日志，将日志发送到通道中
func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		// 先把日志发送到通道中
		// 造一个logMsg对象
		logTmp := &logMsg{
			level:     lv,
			timeStamp: now.Format("2006-01-02 15:04:05"),
			fileName:  fileName,
			funcName:  funcName,
			msg:       msg,
			lineNo:    lineNo,
		}
		select {
		// 将log放入通道中
		case f.logChan <- logTmp:
		// 把日志丢掉保证不出现阻塞
		default:
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
