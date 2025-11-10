package requests

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8services"
	"github.com/saichler/l8utils/go/utils/strings"
)

var (
	ErrDuplicateRequest = errors.New("duplicate request")
)

type Requests struct {
	pending sync.Map // map[string]*Request - eliminates mutex bottleneck
}

type Request struct {
	msgSource      string
	msgNum         uint32
	timeout        time.Duration
	timeoutReached bool
	response       ifs.IElements
	log            ifs.ILogger
	responseChan   chan ifs.IElements
	cancel         context.CancelFunc
	once           sync.Once
	mtx            sync.RWMutex
	requests       *Requests // reference to parent for auto-cleanup
}

func NewRequests() *Requests {
	return &Requests{}
}

func (this *Requests) NewRequest(msgNum uint32, msgSource string, timeoutInSeconds int, log ifs.ILogger) (*Request, error) {
	key := requestKey(msgSource, msgNum)

	request := &Request{
		msgNum:       msgNum,
		msgSource:    msgSource,
		timeout:      time.Duration(timeoutInSeconds) * time.Second,
		log:          log,
		responseChan: make(chan ifs.IElements, 1),
		requests:     this,
	}

	// Check for duplicate and store atomically
	if _, loaded := this.pending.LoadOrStore(key, request); loaded {
		return nil, ErrDuplicateRequest
	}

	return request, nil
}

func (this *Requests) GetRequest(msgNum uint32, msgSource string) *Request {
	key := requestKey(msgSource, msgNum)
	val, ok := this.pending.Load(key)
	if !ok {
		return nil
	}
	return val.(*Request)
}

func (this *Requests) DelRequest(msgNum uint32, msgSource string) {
	key := requestKey(msgSource, msgNum)
	this.pending.Delete(key)
}

func (this *Request) MsgNum() uint32 {
	return this.msgNum
}

func (this *Request) MsgSource() string {
	return this.msgSource
}

func (this *Request) Response() ifs.IElements {
	this.mtx.RLock()
	defer this.mtx.RUnlock()

	if this.timeoutReached {
		return object.NewError("Timeout Reached!")
	}
	return this.response
}

func (this *Request) Wait() ifs.IElements {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	this.cancel = cancel
	defer cancel()
	defer this.cleanup()

	select {
	case resp := <-this.responseChan:
		this.mtx.Lock()
		this.response = resp
		this.mtx.Unlock()
		return resp
	case <-ctx.Done():
		this.mtx.Lock()
		this.timeoutReached = true
		this.mtx.Unlock()
		return object.NewError("Timeout Reached!")
	}
}

func (this *Request) SetResponse(resp ifs.IElements) {
	if this == nil || resp == nil {
		return
	}

	tr, ok := resp.Element().(*l8services.L8Transaction)

	// If transaction is not ended, don't complete the request yet
	if ok && tr.End == 0 {
		return
	}

	// Use once to ensure we only send the final response once
	this.once.Do(func() {
		// Send response to channel (non-blocking)
		select {
		case this.responseChan <- resp:
			// Response sent successfully
		default:
			// Channel already has a response or is closed
		}
	})
}

func (this *Request) cleanup() {
	// Auto-cleanup: remove from pending map
	if this.requests != nil {
		this.requests.DelRequest(this.msgNum, this.msgSource)
	}

	// Close the response channel to free resources
	close(this.responseChan)

	// Cancel context if not already canceled
	if this.cancel != nil {
		this.cancel()
	}
}

func requestKey(msgSource string, msgNum uint32) string {
	return strings.New(msgSource, int(msgNum)).String()
}
