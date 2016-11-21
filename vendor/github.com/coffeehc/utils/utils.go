// utils project utils.go
package utils

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}

/*
	获取App执行文件目录
*/
func GetAppDir() string {
	return filepath.Dir(GetAppPath())
}

/**
读取文件夹的一级文件
*/
func DirList(path string, hasSubDir bool) []string {
	files := make([]string, 0)
	ls, err := ioutil.ReadDir(path)
	if err != nil {
		return files
	}
	for _, fileInfo := range ls {
		if fileInfo.IsDir() {
			if hasSubDir {
				subList := DirList(filepath.Join(path, fileInfo.Name()), hasSubDir)
				for _, sub := range subList {
					files = append(files, sub)
				}

			}
			continue
		}
		files = append(files, filepath.Join(path, fileInfo.Name()))
	}
	return files
}

/*
	获取文件清单
*/
func FileList(path string) []string {
	list := make([]string, 0)
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		list = append(list, path)
		return nil
	})
	return list
}

/**
读取行,此处需要优化成iobuf.sanc
*/
func ReadLine(reader *bufio.Reader) ([]byte, error) {
	base, isPrefix, err := reader.ReadLine()
	if isPrefix {
		return readline(reader, base)
	} else {
		return base, err
	}
}

func readline(reader *bufio.Reader, base []byte) ([]byte, error) {
	//TODO 限制最大大小
	bytes, isPrefix, err := reader.ReadLine()
	if isPrefix {
		return readline(reader, append(base, bytes...))
	} else {
		return append(base, bytes...), err
	}
}
