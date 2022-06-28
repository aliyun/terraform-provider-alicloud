package tunnel

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/common"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tunnel/protocol"
	"github.com/cenkalti/backoff"
	"github.com/golang/protobuf/proto"
	"io"
	"time"
)

var (
	randomizationFactor = 0.33
	backOffMultiplier   = 2.0
)

const (
	//for serialize binary record
	TAG_VERSION       = 0x1
	TAG_RECORD_COUNT  = 0x2
	TAG_ACTION_TYPE   = 0x3
	TAG_RECORD_LENGTH = 0x4
	TAG_RECORD        = 0x5
)

func StreamRecordSequenceLess(a, b *SequenceInfo) bool {
	if a.Epoch < b.Epoch {
		return true
	}
	if a.Epoch == b.Epoch {
		if a.Timestamp < b.Timestamp {
			return true
		}
		if a.Timestamp == b.Timestamp {
			return a.RowIndex < b.RowIndex
		}
	}
	return false
}

func ParseActionType(pbType *protocol.ActionType) (ActionType, error) {
	switch *pbType {
	case protocol.ActionType_PUT_ROW:
		return AT_Put, nil
	case protocol.ActionType_UPDATE_ROW:
		return AT_Update, nil
	case protocol.ActionType_DELETE_ROW:
		return AT_Delete, nil
	default:
		return ActionType(-1), &TunnelError{Code: ErrCodeClientError, Message: fmt.Sprintf("Unexpected action type %s", pbType.String())}
	}
}

func DeserializeRecordFromRawBytes(data []byte, originData []byte, actionType ActionType) (*Record, error) {
	rows, err := protocol.ReadRowsWithHeader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	record := new(Record)
	record.PrimaryKey = &PrimaryKey{}
	record.Type = actionType

	for _, pk := range rows[0].PrimaryKey {
		pkColumn := &PrimaryKeyColumn{ColumnName: string(pk.CellName), Value: pk.CellValue.Value}
		record.PrimaryKey.PrimaryKeys = append(record.PrimaryKey.PrimaryKeys, pkColumn)
	}

	if rows[0].Extension != nil {
		record.Timestamp = rows[0].Extension.Timestamp
		record.SequenceInfo = &SequenceInfo{
			Timestamp: rows[0].Extension.Timestamp,
			RowIndex:  rows[0].Extension.RowIndex,
		}
	}

	for _, cell := range rows[0].Cells {
		cellName := (string)(cell.CellName)
		dataColumn := &RecordColumn{Name: &cellName, Timestamp: &cell.CellTimestamp}
		if cell.CellValue != nil {
			dataColumn.Value = cell.CellValue.Value
		}
		switch cell.CellType {
		case protocol.DELETE_ONE_VERSION:
			dataColumn.Type = RCT_DeleteOneVersion
		case protocol.DELETE_ALL_VERSION:
			dataColumn.Type = RCT_DeleteAllVersions
		default:
			dataColumn.Type = RCT_Put
		}
		record.Columns = append(record.Columns, dataColumn)
	}

	if originData != nil {
		originRows, err := protocol.ReadRowsWithHeader(bytes.NewReader(originData))
		if err != nil {
			return nil, err
		}
		for _, originCell := range originRows[0].Cells {
			cellName := (string)(originCell.CellName)
			dataColumn := &RecordColumn{Name: &cellName, Timestamp: &originCell.CellTimestamp}
			if originCell.CellValue != nil {
				dataColumn.Value = originCell.CellValue.Value
			}
			switch originCell.CellType {
			case protocol.DELETE_ONE_VERSION:
				dataColumn.Type = RCT_DeleteOneVersion
			case protocol.DELETE_ALL_VERSION:
				dataColumn.Type = RCT_DeleteAllVersions
			default:
				dataColumn.Type = RCT_Put
			}
			record.OriginColumns = append(record.OriginColumns, dataColumn)
		}
	}
	return record, nil
}

func ExponentialBackoff(interval, maxInterval, maxElapsed time.Duration, multiplier, jitter float64) *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.RandomizationFactor = jitter
	b.Multiplier = multiplier
	b.InitialInterval = interval
	b.MaxInterval = maxInterval
	b.MaxElapsedTime = maxElapsed
	b.Reset()
	return b
}

func getBackoffConfForDiffUri(uri string, duration time.Duration) (time.Duration, time.Duration, time.Duration) {
	switch uri {
	case readRecordsUri:
		return initRetryIntervalForDataApi, maxRetryIntervalForDataApi, duration
	default:
		return initRetryInterValForMetaApi, maxRetryIntervalForMetaApi, retryMaxElapsedTimeForMetaApi
	}
}

func ParseRequestToken(token string) (*protocol.TokenContentV2, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	tokenPb := new(protocol.Token)
	err = proto.Unmarshal(decoded, tokenPb)
	if err != nil {
		return nil, err
	}

	if tokenPb.Version == nil {
		return nil, errors.New("Token miss must filed: version.")
	}

	innerMessage := tokenPb.Content
	if *tokenPb.Version == 1 {
		innerTokenPb := new(protocol.TokenContent)
		err = proto.Unmarshal(innerMessage, innerTokenPb)

		if err != nil {
			return nil, err
		} else {
			initCount := int64(0)
			return &protocol.TokenContentV2{
				PrimaryKey: innerTokenPb.PrimaryKey,
				Timestamp:  innerTokenPb.Timestamp,
				Iterator:   innerTokenPb.Iterator,
				TotalCount: &initCount,
			}, nil
		}
	} else if *tokenPb.Version == 2 {
		innerTokenPbV2 := new(protocol.TokenContentV2)
		err = proto.Unmarshal(innerMessage, innerTokenPbV2)

		if err != nil {
			return nil, err
		} else {
			return innerTokenPbV2, nil
		}
	} else {
		return nil, fmt.Errorf("not support")
	}
}

func streamToken(token string) (bool, error) {
	tok, err := ParseRequestToken(token)
	if err != nil {
		return false, err
	}
	if tok.GetIterator() == "" {
		return false, nil
	}
	return true, nil
}

func parseTunnelStreamConfig(config *StreamTunnelConfig) *protocol.StreamTunnelConfig {
	if config == nil {
		return nil
	}
	pbConfig := new(protocol.StreamTunnelConfig)
	if config.StartOffset != 0 {
		pbConfig.StartOffset = proto.Uint64(config.StartOffset)
	} else {
		pbConfig.Flag = config.Flag.Enum()
	}
	if config.EndOffset != 0 {
		pbConfig.EndOffset = proto.Uint64(config.EndOffset)
	}
	return pbConfig
}

func parseProtoTunnelStreamConfig(pbConfig *protocol.StreamTunnelConfig) *StreamTunnelConfig {
	if pbConfig == nil {
		return nil
	}
	return &StreamTunnelConfig{
		Flag:        pbConfig.GetFlag(),
		StartOffset: pbConfig.GetStartOffset(),
		EndOffset:   pbConfig.GetEndOffset(),
	}
}

func UnSerializeBatchBinaryRecordFromBytes(data []byte) (record []*Record, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			if _, ok := err2.(error); ok {
				err = err2.(error)
			}
			return
		}
	}()
	records := make([]*Record, 0)
	r := bytes.NewReader(data)
	recordCount := readHeaderInfo(r)
	for ; recordCount > 0; recordCount-- {
		rec, err := UnSerializeBinaryRecordFromBytes(r)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func UnSerializeBinaryRecordFromBytes(r *bytes.Reader) (*Record, error) {
	tag := readTag(r)
	if tag != TAG_ACTION_TYPE {
		panic(ErrUnExpectBinaryRecordTag)
	}
	typ := readRawDataInt32(r)
	tag = readTag(r)
	if tag != TAG_RECORD_LENGTH {
		panic(ErrUnExpectBinaryRecordTag)
	}
	length := readRawDataInt32(r)
	tag = readTag(r)
	if tag != TAG_RECORD {
		panic(ErrUnExpectBinaryRecordTag)
	}
	b := readBytes(r, length)
	pbActionTyp := protocol.ActionType(typ)
	actionTyp, err := ParseActionType(&pbActionTyp)
	if err != nil {
		return nil, err
	}
	rec, err := DeserializeRecordFromRawBytes(b, nil, actionTyp)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func serializeBinaryRecord(b *bytes.Buffer, record *protocol.Record) []byte {
	writeTag(b, TAG_ACTION_TYPE)
	writeRawDataInt32(b, int32(*record.ActionType))
	writeTag(b, TAG_RECORD_LENGTH)
	writeRawDataInt32(b, int32(len(record.Record)))
	writeTag(b, TAG_RECORD)
	b.Write(record.GetRecord())
	return b.Bytes()
}

func writeHeaderInfo(b *bytes.Buffer, count int32) {
	writeTag(b, TAG_VERSION)
	writeRawDataInt32(b, 0)
	writeTag(b, TAG_RECORD_COUNT)
	writeRawDataInt32(b, count)
}

func writeRawDataInt32(w io.Writer, value int32) {
	w.Write([]byte{byte((value) & 0xFF)})
	w.Write([]byte{byte((value >> 8) & 0xFF)})
	w.Write([]byte{byte((value >> 16) & 0xFF)})
	w.Write([]byte{byte((value >> 24) & 0xFF)})
}

func writeTag(w io.Writer, tag byte) {
	writeRawByte(w, tag)
}

func writeRawByte(w io.Writer, value byte) {
	w.Write([]byte{value})
}

func readHeaderInfo(r *bytes.Reader) int32 {
	if readTag(r) != TAG_VERSION {
		panic(ErrUnExpectBinaryRecordTag)
	} else {
		version := readRawDataInt32(r)
		if version != 0 {
			panic(ErrUnSupportRecordVersion)
		}
	}

	if readTag(r) != TAG_RECORD_COUNT {
		panic(ErrUnExpectBinaryRecordTag)
	}
	return readRawDataInt32(r)
}

func readRawDataInt32(r *bytes.Reader) int32 {
	if r.Len() < 4 {
		panic("read raw data panic")
	}
	var v int32
	binary.Read(r, binary.LittleEndian, &v)
	return v
}

func readTag(r *bytes.Reader) int {
	return int(readRawByte(r))
}

func readRawByte(r *bytes.Reader) byte {
	if r.Len() == 0 {
		panic("read tag panic")
	}
	b, _ := r.ReadByte()
	return b
}

func readBytes(r *bytes.Reader, size int32) []byte {
	if int32(r.Len()) < size {
		panic("size not enough")
	}
	v := make([]byte, size)
	r.Read(v)
	return v
}

// SetCredentialsProvider sets funciton for get the user's ak
func SetCredentialsProvider(provider common.CredentialsProvider) ClientOption {
	return func(client *TunnelApi) {
		client.credentialsProvider = provider
	}
}
