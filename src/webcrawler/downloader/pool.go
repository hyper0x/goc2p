package downloader

import (
	"errors"
	"fmt"
	"reflect"
	mdw "webcrawler/middleware"
)

// 生成网页下载器的函数类型。
type GenPageDownloader func() PageDownloader

// 网页下载器池的接口类型。
type PageDownloaderPool interface {
	Take() (PageDownloader, error)  // 从池中取出一个网页下载器。
	Return(dl PageDownloader) error // 把一个网页下载器归还给池。
	Total() uint32                  // 获得池的总容量。
	Used() uint32                   // 获得正在被使用的网页下载器的数量。
}

// 创建网页下载器池。
func NewPageDownloaderPool(
	total uint32,
	gen GenPageDownloader) (PageDownloaderPool, error) {
	etype := reflect.TypeOf(gen())
	genEntity := func() mdw.Entity {
		return gen()
	}
	pool, err := mdw.NewPool(total, etype, genEntity)
	if err != nil {
		return nil, err
	}
	dlpool := &myDownloaderPool{pool: pool, etype: etype}
	return dlpool, nil
}

// 网页下载器池的实现类型。
type myDownloaderPool struct {
	pool  mdw.Pool     // 实体池。
	etype reflect.Type // 池内实体的类型。
}

func (dlpool *myDownloaderPool) Take() (PageDownloader, error) {
	entity, err := dlpool.pool.Take()
	if err != nil {
		return nil, err
	}
	dl, ok := entity.(PageDownloader)
	if !ok {
		errMsg := fmt.Sprintf("The type of entity is NOT %s!\n", dlpool.etype)
		panic(errors.New(errMsg))
	}
	return dl, nil
}

func (dlpool *myDownloaderPool) Return(dl PageDownloader) error {
	return dlpool.pool.Return(dl)
}

func (dlpool *myDownloaderPool) Total() uint32 {
	return dlpool.pool.Total()
}
func (dlpool *myDownloaderPool) Used() uint32 {
	return dlpool.pool.Used()
}
