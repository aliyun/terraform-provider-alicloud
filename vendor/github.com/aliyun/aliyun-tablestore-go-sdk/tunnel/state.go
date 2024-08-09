package tunnel

import (
	"errors"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tunnel/protocol"
	"go.uber.org/zap"
	"sync"
)

type BatchGetStatusReq struct {
	callbackCh chan *sync.Map
}

func NewBatchGetStatusReq() *BatchGetStatusReq {
	return &BatchGetStatusReq{
		callbackCh: make(chan *sync.Map, 1)}
}

type TunnelStateMachine struct {
	batchUpdateStatusCh chan []*protocol.Channel
	batchGetStatusCh    chan *BatchGetStatusReq
	updateStatusCh      chan *protocol.Channel
	closeCh             chan struct{}
	wg                  sync.WaitGroup
	closeOnce           sync.Once

	dialer   ChannelDialer
	pFactory ChannelProcessorFactory
	tunnelId string
	clientId string

	api *TunnelApi

	channelConn     map[string]ChannelConn
	currentChannels *sync.Map

	lg *zap.Logger
}

func NewTunnelStateMachine(tunnelId, clientId string, dialer ChannelDialer, factory ChannelProcessorFactory, api *TunnelApi, lg *zap.Logger) *TunnelStateMachine {
	stateMachine := &TunnelStateMachine{
		batchUpdateStatusCh: make(chan []*protocol.Channel),
		batchGetStatusCh:    make(chan *BatchGetStatusReq),
		updateStatusCh:      make(chan *protocol.Channel),
		closeCh:             make(chan struct{}),
		dialer:              dialer,
		pFactory:            factory,
		tunnelId:            tunnelId,
		clientId:            clientId,
		api:                 api,
		channelConn:         make(map[string]ChannelConn),
		currentChannels:     new(sync.Map),
		lg:                  lg,
	}
	stateMachine.wg.Add(1)
	go stateMachine.bgLoop()
	return stateMachine
}

func (s *TunnelStateMachine) BatchUpdateStatus(batchChannels []*protocol.Channel) {
	select {
	case s.batchUpdateStatusCh <- batchChannels:
	case <-s.closeCh:
	}
}

func (s *TunnelStateMachine) UpdateStatus(channel *ChannelStatus) {
	pbChannel := channel.ToPbChannel()
	select {
	case s.updateStatusCh <- pbChannel:
	case <-s.closeCh:
	}
}

func (s *TunnelStateMachine) BatchGetStatus(req *BatchGetStatusReq) ([]*protocol.Channel, error) {
	select {
	case s.batchGetStatusCh <- req:
		ret := <-req.callbackCh
		length := getSyncMapLength(ret)
		channels := make([]*protocol.Channel, 0, length)
		ret.Range(func(cid, pbChannel interface{}) bool {
			channels = append(channels, pbChannel.(*protocol.Channel))
			return true
		})
		return channels, nil
	case <-s.closeCh:
		return nil, errors.New("state machine is closed")
	}
}

func (s *TunnelStateMachine) Close() {
	s.closeOnce.Do(func() {
		close(s.closeCh)
		s.wg.Wait()
		var syncCloseResource bool
		if factory, ok := s.pFactory.(*SimpleProcessFactory); ok {
			syncCloseResource = factory.SyncCloseResource
		} else if factory, ok := s.pFactory.(*AsyncProcessFactory); ok {
			syncCloseResource = factory.SyncCloseResource
		}
		for _, conn := range s.channelConn {
			conn := conn
			if syncCloseResource {
				conn.Close()
			} else {
				go conn.Close()
			}
		}
		s.lg.Info("state machine is closed")
	})
}

func (s *TunnelStateMachine) bgLoop() {
	defer s.wg.Done()
	for {
		select {
		case channels := <-s.batchUpdateStatusCh:
			s.doBatchUpdateStatus(channels)
		case channel := <-s.updateStatusCh:
			s.doUpdateStatus(channel)
		case req := <-s.batchGetStatusCh:
			req.callbackCh <- s.currentChannels
		case <-s.closeCh:
			s.lg.Info("state machine background loop is going to quite...")
			return
		}
	}
}

func (s *TunnelStateMachine) doUpdateStatus(channel *protocol.Channel) {
	cid := channel.GetChannelId()
	loadChannel, ok := s.currentChannels.Load(cid)
	if !ok {
		s.lg.Info("redundant channel", zap.String("channelId", cid),
			zap.String("status", protocol.ChannelStatus_name[int32(channel.GetStatus())]))
		return
	}
	curChannel := loadChannel.(*protocol.Channel)
	if curChannel.GetVersion() >= channel.GetVersion() {
		s.lg.Info("expired channel version", zap.String("channelId", cid),
			zap.Int64("current version", curChannel.GetVersion()), zap.Int64("old version", channel.GetVersion()))
		return
	}
	s.currentChannels.Store(cid, channel)
	s.lg.Info("update channel", zap.String("channelId", cid),
		zap.String("status", protocol.ChannelStatus_name[int32(channel.GetStatus())]))
	if conn, ok := s.channelConn[cid]; ok {
		s.lg.Info("check channel", zap.String("channelId", cid), zap.Bool("closed", conn.Closed()))
		if conn.Closed() {
			delete(s.channelConn, cid)
		}
	}
}

func (s *TunnelStateMachine) doBatchUpdateStatus(batchChannels []*protocol.Channel) {
	s.currentChannels = validateChannels(batchChannels, s.currentChannels)
	//update
	s.currentChannels.Range(func(channelId, pbChannel interface{}) bool {
		cid := channelId.(string)
		channel := pbChannel.(*protocol.Channel)
		conn, ok := s.channelConn[cid]
		if !ok {
			token, sequenceNumber, err := s.api.GetCheckpoint(s.tunnelId, s.clientId, cid)
			if err != nil {
				s.lg.Error("get channel checkpoint failed", zap.String("tunnelId", s.tunnelId), zap.String("clientId", s.clientId),
					zap.String("channelId", cid), zap.Error(err))
				//failConn will turn channel to closed if necessary
				conn = &failConn{state: s}
			} else {
				processor, err := s.pFactory.NewProcessor(s.tunnelId, s.clientId, cid,
					newCheckpointer(s.api, s.tunnelId, s.clientId, cid, sequenceNumber+1))
				if err != nil {
					s.lg.Error("new processor failed", zap.String("tunnelId", s.tunnelId), zap.String("clientId", s.clientId),
						zap.String("channelId", cid), zap.Error(err))
					//failConn will turn channel to closed if necessary
					conn = &failConn{state: s}
				} else {
					conn = s.dialer.ChannelDial(s.tunnelId, s.clientId, cid, token, processor, s)
				}
			}
			s.channelConn[cid] = conn
		}
		go conn.NotifyStatus(ToChannelStatus(channel))
		return true
	})
	//clean
	for cid, conn := range s.channelConn {
		cid, conn := cid, conn
		if _, ok := s.currentChannels.Load(cid); !ok {
			s.lg.Info("redundant channel conn", zap.String("channelId", cid))
			if !conn.Closed() {
				go conn.Close()
			}
			delete(s.channelConn, cid)
		}
	}

}

func validateChannels(newChans []*protocol.Channel, currentChans *sync.Map) *sync.Map {
	updateChannels := new(sync.Map)
	for _, newChannel := range newChans {
		id := newChannel.GetChannelId()
		if oldChannel, ok := currentChans.Load(id); ok {
			if newChannel.GetVersion() >= oldChannel.(*protocol.Channel).GetVersion() {
				updateChannels.Store(id, newChannel)
			} else {
				updateChannels.Store(id, oldChannel)
			}
		} else {
			updateChannels.Store(id, newChannel)
		}
	}
	return updateChannels
}

func getSyncMapLength(m *sync.Map) int {
	length := 0
	m.Range(func(cid, pbChannel interface{}) bool {
		length++
		return true
	})
	return length
}
