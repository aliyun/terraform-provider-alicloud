package tunnel

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tunnel/protocol"
	"github.com/cenkalti/backoff"
	"go.uber.org/zap"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// PipeChannelSize is channelDialer's data pipeline channel size
	PipeChannelSize = 1

	rpoBar     = 500
	rpoSizeBar = 900 * 1024 //900K bytes
)

type ChannelStatus struct {
	ChannelId string
	Version   int64
	Status    protocol.ChannelStatus
}

func ToChannelStatus(c *protocol.Channel) *ChannelStatus {
	return &ChannelStatus{
		ChannelId: c.GetChannelId(),
		Version:   c.GetVersion(),
		Status:    c.GetStatus(),
	}
}

func (c *ChannelStatus) ToPbChannel() *protocol.Channel {
	clone := c
	return &protocol.Channel{
		ChannelId: &clone.ChannelId,
		Version:   &clone.Version,
		Status:    clone.Status.Enum(),
	}
}

type ChannelConn interface {
	NotifyStatus(channel *ChannelStatus)
	Closed() bool
	Close()
}

// failConn do nothing but turn channel state to close
type failConn struct {
	mu           sync.Mutex
	state        *TunnelStateMachine
	currentState *ChannelStatus
}

func (c *failConn) NotifyStatus(channel *ChannelStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currentState != nil && c.currentState.Version > channel.Version {
		return
	}
	c.currentState = channel
	switch c.currentState.Status {
	case protocol.ChannelStatus_CLOSE:
	case protocol.ChannelStatus_CLOSING:
		c.currentState.Version += 1
		c.currentState.Status = protocol.ChannelStatus_CLOSE
		c.state.UpdateStatus(c.currentState)
	case protocol.ChannelStatus_OPEN:
		c.currentState.Version += 1
		c.currentState.Status = protocol.ChannelStatus_CLOSE
		c.state.UpdateStatus(c.currentState)
	case protocol.ChannelStatus_TERMINATED:
	default:
	}
}

func (c *failConn) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currentState == nil {
		return false
	}
	return c.currentState.Status == protocol.ChannelStatus_CLOSE ||
		c.currentState.Status == protocol.ChannelStatus_TERMINATED
}

func (c *failConn) Close() {}

type ChannelDialer interface {
	ChannelDial(tunnelId, clientId, channelId, token string, p ChannelProcessor, state *TunnelStateMachine) ChannelConn
}

type channelDialer struct {
	api                 tunnelDataApi
	lg                  *zap.Logger
	bc                  *ChannelBackoffConfig
	channelParallelChan chan bool
	needManualRelease   bool
	syncReadRecords     bool
}

func (d *channelDialer) ChannelDial(tunnelId, clientId, channelId, token string, p ChannelProcessor, state *TunnelStateMachine) ChannelConn {
	isStream, err := streamToken(token)
	if err != nil {
		/*
			There are two situations where errors will occur, one is base channel finished, the other is parent stream channel finished.
			we judge the channel to be Base depend on base channelId always begin with  "BaseData_"
		*/
		if strings.HasPrefix(channelId, "BaseData_") {
			isStream = false
		}
		isStream = true //treat as stream token with flow control
	}
	conn := &channelConn{
		tunnelId:      tunnelId,
		clientId:      clientId,
		channelId:     channelId,
		token:         token,
		api:           d.api,
		p:             p,
		state:         state,
		lg:            d.lg,
		bc:            d.bc,
		streamChannel: isStream,
	}
	if d.channelParallelChan != nil {
		conn.parallelReleaseManager = &defaultParallelReleaseManager{
			mu:                  new(sync.Mutex),
			quotaAcquired:       false,
			channelParallelChan: d.channelParallelChan,
		}
		conn.needManualRelease = d.needManualRelease
	}
	if d.syncReadRecords {
		conn.syncReadRecordChan = make(chan bool)
	}
	return conn
}

var (
	waitStatus    = int32(0)
	runningStatus = int32(1)
	closingStatus = int32(2)
	closedStatus  = int32(3)
)

type tunnelDataApi interface {
	ReadRecords(req *ReadRecordRequest) (*ReadRecordResponse, error)
}

type channelConn struct {
	mu sync.Mutex

	tunnelId  string
	clientId  string
	channelId string
	token     string

	api tunnelDataApi
	p   ChannelProcessor

	currentState *ChannelStatus
	state        *TunnelStateMachine

	lg *zap.Logger
	bc *ChannelBackoffConfig

	status int32

	ticker *backoff.Ticker

	streamChannel          bool
	parallelReleaseManager *defaultParallelReleaseManager
	needManualRelease      bool
	syncReadRecordChan     chan bool
}

func (c *channelConn) NotifyStatus(channel *ChannelStatus) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currentState != nil && c.currentState.Version > channel.Version {
		return
	}
	c.currentState = channel
	switch channel.Status {
	case protocol.ChannelStatus_CLOSE:
		c.lg.Info("closed channel status", zap.String("tunnelId", c.tunnelId),
			zap.String("clientId", c.clientId), zap.String("channelId", c.channelId), zap.Int64("version", channel.Version))
		c.close()
	case protocol.ChannelStatus_CLOSING: //draw closing action and check closed/finish status
		if atomic.LoadInt32(&c.status) == waitStatus {
			atomic.StoreInt32(&c.status, closedStatus)
		} else {
			atomic.CompareAndSwapInt32(&c.status, runningStatus, closingStatus)
		}
		c.checkUpdateStatus()
	case protocol.ChannelStatus_OPEN: //draw processing action or check closed/finish status
		if atomic.CompareAndSwapInt32(&c.status, waitStatus, runningStatus) {
			if d, ok := c.p.(*defaultProcessor); ok && d.asyncProcessFlag && d.channelOpenFunc != nil {
				c.lg.Info("open channel", zap.String("tunnelId", c.tunnelId),
					zap.String("clientId", c.clientId), zap.String("channelId", c.channelId), zap.Int64("version", channel.Version))
				d.channelOpenFunc(d.ctx)
			}
			go c.workLoop()
		} else {
			c.checkUpdateStatus()
		}
	case protocol.ChannelStatus_TERMINATED:
		c.lg.Info("terminated channel status", zap.String("tunnelId", c.tunnelId),
			zap.String("clientId", c.clientId), zap.String("channelId", c.channelId), zap.Int64("version", channel.Version))
		c.close()
	default:
		c.lg.Error("Unexpected channel status", zap.String("tunnelId", c.tunnelId),
			zap.String("clientId", c.clientId), zap.String("channelId", c.channelId),
			zap.Int32("Status", int32(c.currentState.Status)), zap.Int64("version", channel.Version))
	}
}

func (c *channelConn) Closed() bool {
	return atomic.LoadInt32(&c.status) == closedStatus
}

func (c *channelConn) Close() {
	c.close()
}

func (c *channelConn) close() {
	c.p.Shutdown()
	atomic.StoreInt32(&c.status, closedStatus)
}

func (c *channelConn) checkUpdateStatus() {
	if c.p.Finished() {
		c.closeAndUpdate(protocol.ChannelStatus_TERMINATED)
		return
	}
	if c.p.Error() {
		c.closeAndUpdate(protocol.ChannelStatus_CLOSE)
		return
	}
	if atomic.LoadInt32(&c.status) == closedStatus {
		if c.currentState.Status != protocol.ChannelStatus_TERMINATED {
			if c.p.Finished() {
				c.closeAndUpdate(protocol.ChannelStatus_TERMINATED)
			} else {
				c.closeAndUpdate(protocol.ChannelStatus_CLOSE)
			}
		}
	}
}

func (c *channelConn) closeAndUpdate(status protocol.ChannelStatus) {
	c.close()
	c.currentState.Version += 1
	c.currentState.Status = status
	c.state.UpdateStatus(c.currentState)
}

type pipeResult struct {
	finished      bool
	records       []*Record
	binaryRecords []byte
	traceId       string
	nextToken     string
	recordCount   int
	error         error
}

func (c *channelConn) workLoop() {
	var pipeSize int
	if p, ok := c.p.(*defaultProcessor); ok {
		pipeSize = p.recordPipePerChannel
	} else {
		//for data delivery
		pipeSize = PipeChannelSize
	}
	pipeCh := make(chan *pipeResult, pipeSize)
	c.lg.Info("run channel workLoop", zap.Int("PipeChannelSize", pipeSize),
		zap.String("cid", c.channelId))
	closeCh := make(chan struct{})
	//clean channel channelParallelChan when processRecords exit or processRecords is not executed
	if c.parallelReleaseManager != nil {
		defer c.parallelReleaseManager.Release()
	}
	defer close(closeCh)
	go c.readRecordsPipe(pipeCh, closeCh)

	for atomic.LoadInt32(&c.status) == runningStatus {
		drained, err := c.processRecords(pipeCh)
		if err != nil {
			c.lg.Info("channel shutdown with error", zap.String("cid", c.channelId), zap.Error(err))
			break
		}
		if drained {
			c.lg.Info("channel read finished", zap.String("cid", c.channelId))
			return
		}
	}
	c.lg.Info("channel shutdown", zap.String("cid", c.channelId))
	c.close()
}

func (c *channelConn) readRecordsPipe(outCh chan *pipeResult, closeCh chan struct{}) {
	//channelParallelChan must be released when readRecords exist
	if c.parallelReleaseManager != nil {
		defer c.parallelReleaseManager.Release()
	}
	var bkoff *backoff.ExponentialBackOff
	if c.streamChannel {
		bkoff = ExponentialBackoff(c.bc.baseDelay, c.bc.MaxDelay, 0, c.bc.factor, c.bc.jitter)
	}
	var needBinaryRecords bool
	if processor, ok := c.p.(*defaultProcessor); ok {
		needBinaryRecords = processor.needBinaryRecords
	}
	for {
		ret := new(pipeResult)
		select {
		case <-closeCh:
			return
		default:
		}

		if c.token == FinishTag {
			ret.finished = true
		} else {
			channelWaitChan := time.Now()
			//whether there is a limit on the concurrent number of channel
			if c.parallelReleaseManager != nil {
				if !c.parallelReleaseManager.checkQuotaAcquired() {
					//there can only be one goroutine acquire token at a time
					c.parallelReleaseManager.mu.Lock()
					select {
					case c.parallelReleaseManager.channelParallelChan <- true:
						c.parallelReleaseManager.quotaAcquired = true
						c.parallelReleaseManager.mu.Unlock()
						if c.needManualRelease {
							c.lg.Info("channel get lock", zap.String("channelId: ", c.channelId), zap.String("tunnelId", c.tunnelId), zap.String("clientId", c.clientId))
						}
					case <-closeCh:
						c.parallelReleaseManager.mu.Unlock()
						return
					}
				}
			}
			s := time.Now()
			request := &ReadRecordRequest{
				TunnelId:         c.tunnelId,
				ChannelId:        c.channelId,
				ClientId:         c.clientId,
				Token:            c.token,
				NeedBinaryRecord: needBinaryRecords,
			}
			resp, err := c.api.ReadRecords(request)
			if err != nil {
				ret.error = err
			} else {
				if needBinaryRecords {
					ret.binaryRecords = resp.BinaryRecords
				} else {
					ret.records = resp.Records
				}
				ret.recordCount = resp.RecordCount
				ret.nextToken = resp.NextToken
				ret.traceId = resp.RequestId
				if bkoff != nil {
					if streamFullData(resp.RecordCount, resp.Size) {
						bkoff.Reset()
					}
				}
				c.token = resp.NextToken
			}
			if ret.recordCount != 0 {
				c.lg.Info("Metric info", zap.String("tunnelId", c.tunnelId), zap.String("clientId", c.clientId),
					zap.String("channelId", c.channelId), zap.String("token", c.token), zap.String("channelWaitChanLatency", s.Sub(channelWaitChan).String()),
					zap.String("GetTunnelRecordLatency", time.Now().Sub(s).String()), zap.Int("recordCount", ret.recordCount), zap.Bool("needBinaryRecords", needBinaryRecords))
			}
		}

		select {
		case outCh <- ret:
			if c.syncReadRecordChan != nil {
				select {
				case c.syncReadRecordChan <- true:
				case <-closeCh:
					return
				}
			}
		case <-closeCh:
			return
		}
		if bkoff != nil {
			time.Sleep(bkoff.NextBackOff())
		}
	}
}

func (c *channelConn) processRecords(inCh chan *pipeResult) (bool, error) {
	var manager ParallelReleaseManager
	if c.parallelReleaseManager != nil {
		manager = c.parallelReleaseManager
		//channel switch control by sdk
		if !c.needManualRelease {
			manager = nil
		}
	}
	if c.syncReadRecordChan != nil {
		defer func() {
			<-c.syncReadRecordChan
		}()
	}
	ret := <-inCh
	if ret.error != nil {
		c.lg.Error("Channel read records failed",
			zap.String("tunnelId", c.tunnelId), zap.String("clientId", c.clientId),
			zap.String("channelId", c.channelId), zap.Error(ret.error))
		return false, ret.error
	}
	if ret.finished {
		return true, nil
	}
	s := time.Now()
	if err := c.p.Process(ret.records, ret.binaryRecords, ret.recordCount, ret.nextToken, ret.traceId, c.streamChannel, manager); err != nil {
		c.lg.Error("Processor process records failed",
			zap.String("tunnelId", c.tunnelId), zap.String("clientId", c.clientId),
			zap.String("channelId", c.channelId), zap.Error(err))
		return false, err
	}
	if ret.recordCount != 0 {
		c.lg.Info("Metric info", zap.String("tunnelId", c.tunnelId), zap.String("clientId", c.clientId),
			zap.String("channelId", c.channelId), zap.String("token", ret.nextToken),
			zap.String("ClientProcessLatency", time.Now().Sub(s).String()))
	}
	//channel switch control by sdk
	if c.parallelReleaseManager != nil && !c.needManualRelease {
		c.parallelReleaseManager.Release()
	}
	return false, nil
}

func streamFullData(numRec int, size int) bool {
	return numRec > rpoBar || size > rpoSizeBar
}

type ParallelReleaseManager interface {
	Release() bool
}

type defaultParallelReleaseManager struct {
	mu                  *sync.Mutex
	quotaAcquired       bool
	channelParallelChan chan bool
}

/*
	Manually release the quota occupied by this channel, after release, Other channels can preempt this channel.
	Only the channel that occupies quota can read and process data.
*/
func (c *defaultParallelReleaseManager) Release() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.quotaAcquired {
		<-c.channelParallelChan
		c.quotaAcquired = false
		return true
	} else {
		return false
	}
}

func (c *defaultParallelReleaseManager) checkQuotaAcquired() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.quotaAcquired
}
