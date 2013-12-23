package web

import (
	"errors"
	"github.com/coffeehc/logger"
	"io"
	"sync"
)

type pageEntity struct {
	io.WriterTo
	key    string
	isMem  bool
	data   []byte
	rwLock *sync.RWMutex
}

var webCache map[string]*pageEntity = make(map[string]*pageEntity)

func (this *pageEntity) WriteTo(w io.Writer) (n int64, err error) {
	if this.isMem {
		nn, e := w.Write(this.data)
		n = int64(nn)
		err = e
		return
	} else {
		logger.Warn("暂时没有实现磁盘缓存")
		return 0, errors.New("暂时没有实现磁盘缓存")
	}

}

func GetCacheEntity(key string) *pageEntity {
	return webCache[key]
}
