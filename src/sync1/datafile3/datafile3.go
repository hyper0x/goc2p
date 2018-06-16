package datafile3

import (
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

// 数据的类型
type Data []byte

// 数据文件的接口类型。
type DataFile interface {
	// 读取一个数据块。
	Read() (rsn int64, d Data, err error)
	// 写入一个数据块。
	Write(d Data) (wsn int64, err error)
	// 获取最后读取的数据块的序列号。
	Rsn() int64
	// 获取最后写入的数据块的序列号。
	Wsn() int64
	// 获取数据块的长度
	DataLen() uint32
}

// 数据文件的实现类型。
type myDataFile struct {
	f       *os.File     // 文件。
	fmutex  sync.RWMutex // 被用于文件的读写锁。
	rcond   *sync.Cond   // 读操作需要用到的条件变量
	woffset int64        // 写操作需要用到的偏移量。
	roffset int64        // 读操作需要用到的偏移量。
	dataLen uint32       // 数据块长度。
}

func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}
	df := &myDataFile{f: f, dataLen: dataLen}
	df.rcond = sync.NewCond(df.fmutex.RLocker())
	return df, nil
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新读偏移量
	var offset int64
	for {
		offset = atomic.LoadInt64(&df.roffset)
		if atomic.CompareAndSwapInt64(&df.roffset, offset, (offset + int64(df.dataLen))) {
			break
		}
	}

	//读取一个数据块
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	df.fmutex.RLock()
	defer df.fmutex.RUnlock()
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				df.rcond.Wait()
				continue
			}
			return
		}
		d = bytes
		return
	}
}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	// 读取并更新写偏移量
	var offset int64
	for {
		offset = atomic.LoadInt64(&df.woffset)
		if atomic.CompareAndSwapInt64(&df.woffset, offset, (offset + int64(df.dataLen))) {
			break
		}
	}

	//写入一个数据块
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	_, err = df.f.Write(bytes)
	df.rcond.Signal()
	return
}

func (df *myDataFile) Rsn() int64 {
	offset := atomic.LoadInt64(&df.roffset)
	return offset / int64(df.dataLen)
}

func (df *myDataFile) Wsn() int64 {
	offset := atomic.LoadInt64(&df.woffset)
	return offset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}
