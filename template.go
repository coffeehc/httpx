// template
package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coffeehc/logger"
)

func initTemplate(tempDir string) *template.Template {
	logger.Debug("初始化模版系统,模版文件夹:%s", tempDir)
	baseTemp := template.New("Base")
	baseTemp.Parse("{{coffeeWeb.}}")
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if !info.IsDir() && ext == ".tmpl" {
			name, _ := filepath.Rel(tempDir, path)
			logger.Info("加载模版文件:%s", name)
			loadTemplate(tempDir, name, baseTemp)
		}
		return nil
	})
	return baseTemp
}

func loadTemplate(tempDir, name string, baseTemp *template.Template) *template.Template {
	t := baseTemp.New(name)
	b, err := ioutil.ReadFile(filepath.Join(tempDir, name))
	if err != nil {
		panic(fmt.Sprintf("读取模版文件失败:%s", err))
	}
	_, err = t.Parse(string(b))
	if err != nil {
		panic(fmt.Sprintf("解析模版文件失败:%s", err))
	}
	return t
}
