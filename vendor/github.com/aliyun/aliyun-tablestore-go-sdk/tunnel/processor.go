package tunnel

import (
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	DefaultChannelSize        = 10
	DefaultCheckpointInterval = 10 * time.Second
)

type ChannelProcessorFactory interface {
	NewProcessor(tunnelId, clientId, channelId string, checkpointer Checkpointer) (ChannelProcessor, error)
}

type ChannelProcessor interface {
	Process(records []*Record, binaryRecords []byte, recordCount int, nextToken, traceId string, streamChannel bool, manager ParallelReleaseManager) error
	Shutdown()
	Error() bool
	Finished() bool
	CommitToken(token string) error
}

type SimpleProcessFactory struct {
	CustomValue interface{}

	CpInterval time.Duration

	ProcessFunc     func(channelCtx *ChannelContext, records []*Record) error
	ChannelOpenFunc func(channelCtx *ChannelContext) error
	ShutdownFunc    func(channelCtx *ChannelContext)

	Logger               *zap.Logger
	NeedBinaryRecord     bool
	SyncCloseResource    bool // sync close tunnel resource before daemon exit
	AlwaysCallBack       bool // when AlwaysCallBack is true, processFunc will be called even if the record is empty
	RecordPipePerChannel *int
}

func (s *SimpleProcessFactory) NewProcessor(tunnelId, clientId, channelId string, checkpointer Checkpointer) (ChannelProcessor, error) {
	var lg *zap.Logger
	if s.Logger == nil {
		lg, _ = DefaultLogConfig.Build(ReplaceLogCore(DefaultSyncer, DefaultLogConfig))
	} else {
		lg = s.Logger
	}
	interval := DefaultCheckpointInterval
	if s.CpInterval > 0 {
		interval = s.CpInterval
	}
	if s.RecordPipePerChannel == nil || *s.RecordPipePerChannel < 0 {
		s.RecordPipePerChannel = &PipeChannelSize
	}
	p := &defaultProcessor{
		ctx:                  NewChannelContext(tunnelId, clientId, channelId, s.CustomValue),
		checkpointer:         checkpointer,
		processFunc:          s.ProcessFunc,
		channelOpenFunc:      s.ChannelOpenFunc,
		shutdownFunc:         s.ShutdownFunc,
		checkpointCh:         make(chan string, DefaultChannelSize),
		closeCh:              make(chan struct{}),
		ticker:               time.NewTicker(interval),
		wg:                   new(sync.WaitGroup),
		lg:                   lg,
		finished:             atomic.NewBool(false),
		needBinaryRecords:    s.NeedBinaryRecord,
		alwaysCallBack:       s.AlwaysCallBack,
		recordPipePerChannel: *s.RecordPipePerChannel,
	}
	p.wg.Add(1)
	go p.cpLoop()
	lg.Info("new sync processor success, ", zap.String("tunnelId", tunnelId), zap.String("clientId", clientId),
		zap.String("channelId", channelId), zap.Bool("needBinaryRecords", s.NeedBinaryRecord), zap.Bool("syncCloseResource", s.SyncCloseResource))
	return p, nil
}

// AsyncProcessFactory support async commit token
type AsyncProcessFactory struct {
	CustomValue interface{}

	CpInterval time.Duration

	ProcessFunc     func(channelCtx *ChannelContext, records []*Record) error
	ChannelOpenFunc func(channelCtx *ChannelContext) error
	ShutdownFunc    func(channelCtx *ChannelContext)

	Logger               *zap.Logger
	NeedBinaryRecord     bool
	SyncCloseResource    bool // sync close tunnel resource before daemon exit
	AlwaysCallBack       bool // when AlwaysCallBack is true, processFunc will be called even if the record is empty
	RecordPipePerChannel *int
}

func (s *AsyncProcessFactory) NewProcessor(tunnelId, clientId, channelId string, checkpointer Checkpointer) (ChannelProcessor, error) {
	var lg *zap.Logger
	if s.Logger == nil {
		lg, _ = DefaultLogConfig.Build(ReplaceLogCore(DefaultSyncer, DefaultLogConfig))
	} else {
		lg = s.Logger
	}
	interval := DefaultCheckpointInterval
	if s.CpInterval > 0 {
		interval = s.CpInterval
	}
	if s.RecordPipePerChannel == nil || *s.RecordPipePerChannel < 0 {
		s.RecordPipePerChannel = &PipeChannelSize
	}
	p := &defaultProcessor{
		ctx:                  NewChannelContext(tunnelId, clientId, channelId, s.CustomValue),
		checkpointer:         checkpointer,
		processFunc:          s.ProcessFunc,
		channelOpenFunc:      s.ChannelOpenFunc,
		shutdownFunc:         s.ShutdownFunc,
		checkpointCh:         make(chan string, DefaultChannelSize),
		closeCh:              make(chan struct{}),
		ticker:               time.NewTicker(interval),
		wg:                   new(sync.WaitGroup),
		lg:                   lg,
		finished:             atomic.NewBool(false),
		asyncProcessFlag:     true,
		needBinaryRecords:    s.NeedBinaryRecord,
		alwaysCallBack:       s.AlwaysCallBack,
		recordPipePerChannel: *s.RecordPipePerChannel,
	}
	lg.Info("new async processor success, ", zap.String("tunnelId", tunnelId), zap.String("clientId", clientId),
		zap.String("channelId", channelId), zap.Bool("asyncProcessFlag", true),
		zap.Bool("needBinaryRecords", s.NeedBinaryRecord), zap.Bool("SyncCloseResource", s.SyncCloseResource))
	return p, nil
}

type defaultProcessor struct {
	ctx *ChannelContext

	checkpointer    Checkpointer
	processFunc     func(channelCtx *ChannelContext, records []*Record) error
	channelOpenFunc func(channelCtx *ChannelContext) error
	shutdownFunc    func(channelCtx *ChannelContext)

	checkpointCh chan string
	closeCh      chan struct{}
	closeOnce    sync.Once
	ticker       *time.Ticker
	wg           *sync.WaitGroup

	finished *atomic.Bool

	lg *zap.Logger

	needBinaryRecords    bool
	asyncProcessFlag     bool
	alwaysCallBack       bool
	recordPipePerChannel int
}

func (p *defaultProcessor) Process(records []*Record, binaryRecords []byte, recordCount int, nextToken, traceId string, isStreamChannel bool, manager ParallelReleaseManager) error {
	if recordCount != 0 || p.alwaysCallBack {
		ctx := &ChannelContext{
			TunnelId:               p.ctx.TunnelId,
			ClientId:               p.ctx.ClientId,
			ChannelId:              p.ctx.ChannelId,
			IsStreamChannel:        isStreamChannel,
			TraceId:                traceId,
			RecordCount:            recordCount,
			NextToken:              nextToken,
			BinaryRecords:          binaryRecords,
			Processor:              p,
			CustomValue:            p.ctx.CustomValue,
			ParallelReleaseManager: manager,
		}
		err := p.processFunc(ctx, records)
		if err != nil {
			return err
		}
	}
	//if p.async is true, then skip checkpoint
	if !p.asyncProcessFlag {
		select {
		case p.checkpointCh <- nextToken:
		case <-p.closeCh:
		}
	}
	if nextToken == FinishTag {
		p.finished.Store(true)
		p.Shutdown()
	}
	return nil
}

func (p *defaultProcessor) Error() bool {
	return false
}

func (p *defaultProcessor) Shutdown() {
	p.closeOnce.Do(func() {
		close(p.closeCh)
		p.ticker.Stop()
		p.wg.Wait()
		if p.shutdownFunc != nil {
			ctx := &ChannelContext{
				TunnelId:    p.ctx.TunnelId,
				ClientId:    p.ctx.ClientId,
				ChannelId:   p.ctx.ChannelId,
				CustomValue: p.ctx.CustomValue,
			}
			p.shutdownFunc(ctx)
		}
	})
}

func (p *defaultProcessor) Finished() bool {
	return p.finished.Load()
}

func (p *defaultProcessor) cpLoop() {
	defer p.wg.Done()
	newCp := ""
	cpFlush := make(chan string)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			cp, ok := <-cpFlush
			if !ok {
				return
			}
			err := p.checkpointer.Checkpoint(cp)
			if err != nil {
				p.lg.Error("make checkpoint failed", zap.String("checkpoint", cp), zap.Error(err))
			} else {
				p.lg.Info("checkpoint progress", zap.String("context", p.ctx.String()), zap.String("checkpoint", cp), zap.Bool("asyncCheckpoint", p.asyncProcessFlag))
			}
		}
	}()
	for {
		select {
		case cp := <-p.checkpointCh:
			newCp = cp
		case <-p.ticker.C:
			if newCp != "" {
				select {
				case cpFlush <- newCp:
					newCp = ""
				default:
				}
			}
		case <-p.closeCh:
			if newCp != "" {
				cpFlush <- newCp
			}
			close(cpFlush)
			return
		}
	}
}

func (p *defaultProcessor) CommitToken(token string) error {
	err := p.checkpointer.Checkpoint(token)
	if err != nil {
		p.lg.Error("async commit checkpoint failed", zap.String("checkpoint", token), zap.Error(err))
		return err
	} else {
		p.lg.Info("async commit checkpoint progress", zap.String("context", p.ctx.String()), zap.String("checkpoint", token))
	}
	return nil
}
