package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

//FileLogRotatePolicy 文件日志循环策略
type FileLogRotatePolicy interface {
	//是否需要循环分割
	CanRotate(fileLogWriter *FileLogWriter) bool
	//循环风格的后置处理
	RotateAfter()
}
type defaultRotatePolicy struct {
}

func (defaultRotatePolicy) CanRotate(fileLogWriter *FileLogWriter) bool {
	return false
}

func (defaultRotatePolicy) RotateAfter() {
}

//SizeRotatePolicy 按文件大小的循环日志的策略
type SizeRotatePolicy struct {
	//文件最大大小
	maxBytes int64
}

//NewSizeRotatePolicy 创建一个按大小循环的策略
func NewSizeRotatePolicy(maxBytes int64) *SizeRotatePolicy {
	return &SizeRotatePolicy{maxBytes: maxBytes}
}

// CanRotate 判断文件是否到最大限制大小
func (policy *SizeRotatePolicy) CanRotate(fileLogWriter *FileLogWriter) bool {
	return fileLogWriter.count > policy.maxBytes
}

// RotateAfter 在循环截断后的处理,未实现
func (*SizeRotatePolicy) RotateAfter() {
}

//TimeRotatePolicy 按时间循环日志的策略
type TimeRotatePolicy struct {
	canRotate bool
}

// CanRotate 判断文件是否能 Rotate
func (policy *TimeRotatePolicy) CanRotate(fileLogWriter *FileLogWriter) bool {
	return policy.canRotate
}

// RotateAfter 在循环截断后的处理
func (policy *TimeRotatePolicy) RotateAfter() {
	policy.canRotate = false
}

//NewTimeRotatePolicy 创建信的按时间循环日志的策略
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

//FileLogWriter 文件日志Writer封装
type FileLogWriter struct {
	err    error
	buf    []byte
	n      int
	wr     *os.File
	config *FileLogConfig
	count  int64
}

//FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path         string //匹配路径
	Timeformat   string //时间格式
	Format       string
	Rotate       int                 //rotate备份个数
	StorePath    string              //日志存放路径,如:/log/test/log.log
	RotatePolicy FileLogRotatePolicy //循环策略
	Level        Level               //日志级别
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
	AddFilter(conf.Level, conf.Path, conf.Timeformat, conf.Format, conf.writer)
}

func addFileFilterForDefualt(level Level, path string, logPath string, timeFormat string, format string) {
	conf := new(FileLogConfig)
	conf.Level = level
	conf.Path = path
	conf.StorePath = logPath
	conf.Timeformat = timeFormat
	conf.Format = format
	addFileFilter(conf)
}

func addFileFilterForTime(level Level, path string, logPath string, times time.Duration, rotate int, timeFormat string, format string) {
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

func addFileFilterForSize(level Level, path string, logPath string, maxBytes int64, rotate int, timeFormat string, format string) {
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

//Rotate 执行循环切割日志操作
func (writer *FileLogWriter) Rotate() {
	if writer.wr == nil {
		err := os.MkdirAll(filepath.Dir(writer.config.StorePath), 0666)
		if err != nil {
			panic(fmt.Sprintf("创建日志目录[%s]失败:%s", filepath.Dir(writer.config.StorePath), err))
		}
		fileInfo, err := os.Stat(writer.config.StorePath)
		if err == nil && fileInfo.Size() != 0 {
			bakName := writer.config.StorePath + "." + strconv.FormatInt(time.Now().Unix(), 10)
			err := os.Rename(writer.config.StorePath, bakName)
			if err != nil {
				panic(fmt.Sprintf("备份老的日志文件失败:%s", err))
			}
		}
		writer.wr, err = os.OpenFile(writer.config.StorePath, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(fmt.Sprintf("打开日志文件失败,%s", err))
		}
	} else {
		writer.wr.Close()
		bakName := writer.config.StorePath + "." + strconv.FormatInt(time.Now().Unix(), 10)
		err := os.Rename(writer.config.StorePath, bakName)
		if err != nil {
			Error(fmt.Sprintf("循环备份老的日志文件失败:%s", err))
		}
		file, err := os.Create(writer.config.StorePath)
		if err != nil {
			//日志创建失败直接让程序挂掉
			println(fmt.Sprintf("创建日志文件失败:%s", err))
			panic(fmt.Sprintf("创建日志文件失败:%s", err))
		}
		writer.wr = file
		go clearLog(writer.config.StorePath, writer.config.Rotate)
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
			log.Fatalf("切割日志出错,%s", r)
			this.count = 0
		}
		this.count = 0
	}()
	this.Flush()
	this.Rotate()
	this.config.RotatePolicy.RotateAfter()
}

//Write 写入日志
func (writer *FileLogWriter) Write(p []byte) (nn int, err error) {
	if writer.config.RotatePolicy.CanRotate(writer) {
		deforeRotateProcess(writer)
	}
	for len(p) > (len(writer.buf)-writer.n) && writer.err == nil {
		var n int
		if writer.n == 0 {
			n, writer.err = writer.wr.Write(p)
		} else {
			n = copy(writer.buf[writer.n:], p)
			writer.n += n
			writer.Flush()
		}
		nn += n
		p = p[n:]
	}
	defer func() {
		writer.count += int64(nn)
	}()
	if writer.err != nil {
		return nn, writer.err
	}
	n := copy(writer.buf[writer.n:], p)
	writer.n += n
	nn += n
	return nn, nil
}

//Flush 落盘操作
func (writer *FileLogWriter) Flush() error {
	if writer.err != nil {
		return writer.err
	}
	if writer.n == 0 {
		return nil
	}
	n, err := writer.wr.Write(writer.buf[0:writer.n])
	if n < writer.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < writer.n {
			copy(writer.buf[0:writer.n-n], writer.buf[n:writer.n])
		}
		writer.n -= n
		writer.err = err
		return err
	}
	writer.n = 0
	return nil
}
