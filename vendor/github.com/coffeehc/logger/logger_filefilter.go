// logger_filefilter
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

//文件日志循环策略
type FileLogRotatePolicy interface {
	//是否需要循环分割
	CanRotate(fileLogWriter *FileLogWriter) bool
	//循环风格的后置处理
	RotateAfter()
}
type defaultRotatePolicy struct {
}

func (this *defaultRotatePolicy) CanRotate(fileLogWriter *FileLogWriter) bool {
	return false
}
func (this *defaultRotatePolicy) RotateAfter() {
}

//按文件大小的循环日志的策略
type SizeRotatePolicy struct {
	//文件最大大小
	maxBytes int64
}

//创建一个按大小循环的策略
func NewSizeRotatePolicy(maxBytes int64) *SizeRotatePolicy {
	return &SizeRotatePolicy{maxBytes: maxBytes}
}

func (this *SizeRotatePolicy) CanRotate(fileLogWriter *FileLogWriter) bool {
	return fileLogWriter.count > this.maxBytes
}
func (this *SizeRotatePolicy) RotateAfter() {
}

//按时间循环日志的策略
type TimeRotatePolicy struct {
	canRotate bool
}

func (this *TimeRotatePolicy) CanRotate(fileLogWriter *FileLogWriter) bool {
	return this.canRotate
}
func (this *TimeRotatePolicy) RotateAfter() {
	this.canRotate = false
}

//创建信的按时间循环日志的策略
func NewTimeRotatePolicy(delay time.Duration) *TimeRotatePolicy {
	this := new(TimeRotatePolicy)
	this.canRotate = false
	now := time.Now()
	//对齐
	nowTime := now.Truncate(delay)
	nowTime = nowTime.Add(delay)
	go func() {
		for {
			select {
			case <-time.After(nowTime.Sub(now)):
				this.canRotate = true
				break
			}
		}
	}()
	return this
}

//文件日志Writer封装
type FileLogWriter struct {
	err    error
	buf    []byte
	n      int
	wr     *os.File
	config *FileLogConfig
	count  int64
}

//文件日志配置
type FileLogConfig struct {
	Path         string //匹配路径
	Timeformat   string //时间格式
	Format       string
	Rotate       int                 //rotate备份个数
	StorePath    string              //日志存放路径,如:/log/test/log.log
	RotatePolicy FileLogRotatePolicy //循环策略
	Level        byte                //日志级别
	writer       *FileLogWriter
}

func checkConfig(conf *FileLogConfig) {
	if conf.Path == "" {
		panic("没有指定需要记录的日志路径")
	}
	if conf.StorePath == "" {
		panic("没有指定日志的存储路径")
	}
	conf.StorePath = strings.Replace(conf.StorePath, "\\", "/", 100)
	if conf.Rotate == 0 {
		conf.Rotate = 3
	}
	if conf.RotatePolicy == nil {
		conf.RotatePolicy = new(defaultRotatePolicy)
	}
}
func addFileFilter(conf *FileLogConfig) {
	checkConfig(conf)
	conf.writer = new(FileLogWriter)
	conf.writer.config = conf
	conf.writer.count = 0
	conf.writer.Rotate()
	AddFileter(conf.Level, conf.Path, conf.Timeformat, conf.Format, conf.writer)
}

func addFileFilterForDefualt(level byte, path string, logPath string, timeFormat string, format string) {
	conf := new(FileLogConfig)
	conf.Level = level
	conf.Path = path
	conf.StorePath = logPath
	conf.Timeformat = timeFormat
	conf.Format = format
	addFileFilter(conf)
}

func addFileFilterForTime(level byte, path string, logPath string, times time.Duration, rotate int, timeFormat string, format string) {
	conf := new(FileLogConfig)
	conf.Level = level
	conf.Path = path
	conf.StorePath = logPath
	conf.Rotate = rotate
	conf.RotatePolicy = NewTimeRotatePolicy(times)
	conf.Timeformat = timeFormat
	conf.Format = format
	addFileFilter(conf)
}

func addFileFilterForSize(level byte, path string, logPath string, maxBytes int64, rotate int, timeFormat string, format string) {
	conf := new(FileLogConfig)
	conf.Level = level
	conf.Path = path
	conf.StorePath = logPath
	conf.Rotate = rotate
	conf.RotatePolicy = NewSizeRotatePolicy(maxBytes)
	conf.Timeformat = timeFormat
	conf.Format = format
	addFileFilter(conf)
}

//执行循环切割日志操作
func (this *FileLogWriter) Rotate() {
	if this.wr == nil {
		err := os.MkdirAll(filepath.Dir(this.config.StorePath), 0666)
		if err != nil {
			panic(fmt.Sprintf("创建日志目录[%s]失败:%s", filepath.Dir(this.config.StorePath), err))
		}
		fileInfo, err := os.Stat(this.config.StorePath)
		if err == nil && fileInfo.Size() != 0 {
			bakName := this.config.StorePath + "." + strconv.FormatInt(time.Now().Unix(), 10)
			err := os.Rename(this.config.StorePath, bakName)
			if err != nil {
				panic(fmt.Sprintf("备份老的日志文件失败:%s", err))
			}
		}
		this.wr, err = os.OpenFile(this.config.StorePath, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(fmt.Sprintf("打开日志文件失败,%s", err))
		}
	} else {
		this.wr.Close()
		bakName := this.config.StorePath + "." + strconv.FormatInt(time.Now().Unix(), 10)
		err := os.Rename(this.config.StorePath, bakName)
		if err != nil {
			Error(fmt.Sprintf("循环备份老的日志文件失败:%s", err))
		}
		file, err := os.Create(this.config.StorePath)
		if err != nil {
			//日志创建失败直接让程序挂掉
			println(fmt.Sprintf("创建日志文件失败:%s", err))
			panic(fmt.Sprintf("创建日志文件失败:%s", err))
		}
		this.wr = file
		go clearLog(this.config.StorePath, this.config.Rotate)
	}
}

func clearLog(logPath string, rotateSize int) {
	dirIndex := strings.LastIndex(logPath, "/")
	files := make([]string, 0)
	logPath = logPath[:dirIndex]
	filepath.Walk(logPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		path = strings.Replace(path, "\\", "/", 100)
		if strings.LastIndex(path, "/") == dirIndex {
			files = append(files, path)
		}
		return nil
	})
	deleteFileSize := len(files) - rotateSize
	if deleteFileSize > 1 {
		sort.Sort(sort.StringSlice(files))
		for i := 1; i < deleteFileSize; i++ {
			os.Remove(files[i])
		}
	}
}

func deforeRotateProcess(this *FileLogWriter) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("切割日志出错,%s", r)
			this.count = 0
		}
		this.count = 0
	}()
	this.Flush()
	this.Rotate()
	this.config.RotatePolicy.RotateAfter()
}

//写入日志
func (this *FileLogWriter) Write(p []byte) (nn int, err error) {
	if this.config.RotatePolicy.CanRotate(this) {
		deforeRotateProcess(this)
	}
	for len(p) > (len(this.buf)-this.n) && this.err == nil {
		var n int
		if this.n == 0 {
			n, this.err = this.wr.Write(p)
		} else {
			n = copy(this.buf[this.n:], p)
			this.n += n
			this.Flush()
		}
		nn += n
		p = p[n:]
	}
	defer func() {
		this.count += int64(nn)
	}()
	if this.err != nil {
		return nn, this.err
	}
	n := copy(this.buf[this.n:], p)
	this.n += n
	nn += n
	return nn, nil
}

//落盘操作
func (this *FileLogWriter) Flush() error {
	if this.err != nil {
		return this.err
	}
	if this.n == 0 {
		return nil
	}
	n, err := this.wr.Write(this.buf[0:this.n])
	if n < this.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < this.n {
			copy(this.buf[0:this.n-n], this.buf[n:this.n])
		}
		this.n -= n
		this.err = err
		return err
	}
	this.n = 0
	return nil
}
