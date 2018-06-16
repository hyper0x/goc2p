package loadgen

import (
	"bytes"
	"errors"
	"fmt"
	lib "loadgen/lib"
	"math"
	"time"
)

// 载荷发生器的实现。
type myGenerator struct {
	caller      lib.Caller           // 调用器。
	timeoutNs   time.Duration        // 处理超时时间，单位：纳秒。
	lps         uint32               // 每秒载荷量。
	durationNs  time.Duration        // 负载持续时间，单位：纳秒。
	concurrency uint32               // 并发量。
	tickets     lib.GoTickets        // Goroutine票池。
	stopSign    chan byte            // 停止信号的传递通道。
	cancelSign  byte                 // 取消发送后续结果的信号。
	endSign     chan uint64          // 完结信号的传递通道，同时被用于传递调用执行计数。
	callCount   uint64               // 调用执行计数。
	status      lib.GenStatus        // 状态。
	resultCh    chan *lib.CallResult // 调用结果通道。
}

func NewGenerator(
	caller lib.Caller,
	timeoutNs time.Duration,
	lps uint32,
	durationNs time.Duration,
	resultCh chan *lib.CallResult) (lib.Generator, error) {
	logger.Infoln("New a load generator...")
	logger.Infoln("Checking the parameters...")
	var errMsg string
	if caller == nil {
		errMsg = fmt.Sprintln("Invalid caller!")
	}
	if timeoutNs == 0 {
		errMsg = fmt.Sprintln("Invalid timeoutNs!")
	}
	if lps == 0 {
		errMsg = fmt.Sprintln("Invalid lps(load per second)!")
	}
	if durationNs == 0 {
		errMsg = fmt.Sprintln("Invalid durationNs!")
	}
	if resultCh == nil {
		errMsg = fmt.Sprintln("Invalid result channel!")
	}
	if errMsg != "" {
		return nil, errors.New(errMsg)
	}
	gen := &myGenerator{
		caller:     caller,
		timeoutNs:  timeoutNs,
		lps:        lps,
		durationNs: durationNs,
		stopSign:   make(chan byte, 1),
		cancelSign: 0,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   resultCh,
	}
	logger.Infof("Passed. (timeoutNs=%v, lps=%d, durationNs=%v)",
		timeoutNs, lps, durationNs)
	err := gen.init()
	if err != nil {
		return nil, err
	}
	return gen, nil
}

func (gen *myGenerator) init() error {
	logger.Infoln("Initializing the load generator...")
	// 载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
	var total64 int64 = int64(gen.timeoutNs)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	gen.concurrency = uint32(total64)
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets
	logger.Infof("Initialized. (concurrency=%d)", gen.concurrency)
	return nil
}

func (gen *myGenerator) interact(rawReq *lib.RawReq) *lib.RawResp {
	if rawReq == nil {
		return &lib.RawResp{Id: -1, Err: errors.New("Invalid raw request.")}
	}
	start := time.Now().Nanosecond()
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNs)
	end := time.Now().Nanosecond()
	elapsedTime := time.Duration(end - start)
	var rawResp lib.RawResp
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResp{
			Id:     rawReq.Id,
			Err:    errors.New(errMsg),
			Elapse: elapsedTime}
	} else {
		rawResp = lib.RawResp{
			Id:     rawReq.Id,
			Resp:   resp,
			Elapse: elapsedTime}
	}
	return &rawResp
}

func (gen *myGenerator) asyncCall() {
	gen.tickets.Take()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				err, ok := interface{}(p).(error)
				var buff bytes.Buffer
				buff.WriteString("Async Call Panic! (")
				if ok {
					buff.WriteString("error: ")
					buff.WriteString(err.Error())
				} else {
					buff.WriteString("clue: ")
					buff.WriteString(fmt.Sprintf("%v", p))
				}
				buff.WriteString(")")
				errMsg := buff.String()
				logger.Fatalln(errMsg)
				result := &lib.CallResult{
					Id:   -1,
					Code: lib.RESULT_CODE_FATAL_CALL,
					Msg:  errMsg}
				gen.sendResult(result)
			}
		}()
		rawReq := gen.caller.BuildReq()
		var timeout bool
		timer := time.AfterFunc(gen.timeoutNs, func() {
			timeout = true
			result := &lib.CallResult{
				Id:   rawReq.Id,
				Req:  rawReq,
				Code: lib.RESULT_CODE_WARNING_CALL_TIMEOUT,
				Msg:  fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNs)}
			gen.sendResult(result)
		})
		rawResp := gen.interact(&rawReq)
		if !timeout {
			timer.Stop()
			var result *lib.CallResult
			if rawResp.Err != nil {
				result = &lib.CallResult{
					Id:     rawResp.Id,
					Req:    rawReq,
					Code:   lib.RESULT_CODE_ERROR_CALL,
					Msg:    rawResp.Err.Error(),
					Elapse: rawResp.Elapse}
			} else {
				result = gen.caller.CheckResp(rawReq, *rawResp)
				result.Elapse = rawResp.Elapse
			}
			gen.sendResult(result)
		}
		gen.tickets.Return()
	}()
}

func (gen *myGenerator) sendResult(result *lib.CallResult) bool {
	if gen.status == lib.STATUS_STARTED && gen.cancelSign == 0 {
		gen.resultCh <- result
		return true
	}
	logger.Warnf("Ignore result: %s.\n",
		fmt.Sprintf(
			"Id=%d, Code=%d, Msg=%s, Elapse=%v",
			result.Id, result.Code, result.Msg, result.Elapse))
	return false
}

func (gen *myGenerator) handleStopSign(callCount uint64) {
	gen.cancelSign = 1
	logger.Infof("Closing result channel...")
	close(gen.resultCh)
	gen.endSign <- callCount
	gen.endSign <- callCount
}

func (gen *myGenerator) genLoad(throttle <-chan time.Time) {
	callCount := uint64(0)
Loop:
	for ; ; callCount++ {
		select {
		case <-gen.stopSign:
			gen.handleStopSign(callCount)
			break Loop
		default:
		}
		gen.asyncCall()
		if gen.lps > 0 {
			select {
			case <-throttle:
			case <-gen.stopSign:
				gen.handleStopSign(callCount)
				break Loop
			}
		}
	}
}

func (gen *myGenerator) Start() {
	logger.Infoln("Starting load generator...")

	// 设定节流阀
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		logger.Infof("Setting throttle (%v)...", interval)
		throttle = time.Tick(interval)
	}

	// 初始化停止信号
	go func() {
		time.AfterFunc(gen.durationNs, func() {
			logger.Infof("Stopping load generator...")
			gen.stopSign <- 0
		})
	}()

	// 初始化完结信号通道
	gen.endSign = make(chan uint64, 2)

	// 初始化调用执行计数
	gen.callCount = 0

	// 设置已启动状态
	gen.status = lib.STATUS_STARTED

	go func() {
		// 生成载荷
		logger.Infoln("Generating loads...")
		gen.genLoad(throttle)

		// 接收调用执行计数
		callCount := <-gen.endSign
		gen.status = lib.STATUS_STOPPED
		logger.Infof("Stopped. (callCount=%d)\n", callCount)
	}()
}

func (gen *myGenerator) Stop() (uint64, bool) {
	if gen.stopSign == nil {
		return 0, false
	}
	if gen.status != lib.STATUS_STARTED {
		return 0, false
	}
	gen.status = lib.STATUS_STOPPED
	gen.stopSign <- 1
	callCount := <-gen.endSign
	return callCount, true
}

func (gen *myGenerator) Status() lib.GenStatus {
	return gen.status
}
