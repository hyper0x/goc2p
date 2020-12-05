package lib

import (
	"time"
)

// 调用器的接口。
type Caller interface {
	// 构建请求。
	BuildReq() RawReq
	// 调用。
	Call(req []byte, timeoutNs time.Duration) ([]byte, error)
	// 检查响应。
	CheckResp(rawReq RawReq, rawResp RawResp) *CallResult
}
