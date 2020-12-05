package downloader

import (
	"logging"
	"net/http"
	base "webcrawler/base"
	mdw "webcrawler/middleware"
)

// 日志记录器。
var logger logging.Logger = base.NewLogger()

// ID生成器。
var downloaderIdGenerator mdw.IdGenerator = mdw.NewIdGenerator()

// 生成并返回ID。
func genDownloaderId() uint32 {
	return downloaderIdGenerator.GetUint32()
}

// 网页下载器的接口类型。
type PageDownloader interface {
	Id() uint32                                        // 获得ID。
	Download(req base.Request) (*base.Response, error) // 根据请求下载网页并返回响应。
}

// 创建网页下载器。
func NewPageDownloader(client *http.Client) PageDownloader {
	id := genDownloaderId()
	if client == nil {
		client = &http.Client{}
	}
	return &myPageDownloader{
		id:         id,
		httpClient: *client,
	}
}

// 网页下载器的实现类型。
type myPageDownloader struct {
	id         uint32      // ID。
	httpClient http.Client // HTTP客户端。
}

func (dl *myPageDownloader) Id() uint32 {
	return dl.id
}

func (dl *myPageDownloader) Download(req base.Request) (*base.Response, error) {
	httpReq := req.HttpReq()
	logger.Infof("Do the request (url=%s)... \n", httpReq.URL)
	httpResp, err := dl.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	return base.NewResponse(httpResp, req.Depth()), nil
}
