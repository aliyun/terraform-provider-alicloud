package tablestore

import "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/otsprotocol"

type TaskType int32

const (
	BaseTask TaskType = iota
	IncTask
	BaseIncTask
)

type TimeFormat int32

const (
	EventColumnRFC822  TimeFormat = 1
	EventColumnRFC850  TimeFormat = 2
	EventColumnRFC1123 TimeFormat = 3
	EventColumnRFC3339 TimeFormat = 4
	EventColumnUnix    TimeFormat = 5
)

type EventColumn struct {
	ColumnName string
	TimeFormat TimeFormat
}

type Format int32

const ParquetFormat Format = 0

type ParquetDataType int32

const (
	ParquetBool ParquetDataType = iota
	ParquetInt64
	ParquetUtf8
	ParquetDouble
	ParquetDate
	ParquetDecimal
	ParquetTimeMills
	ParquetTimeMicros
	ParquetTimestampMills
	ParquetTimestampMicros
)

type Encoding int32

const (
	EncodingPlain Encoding = iota
	EncodingPlainDictionary
	EncodingDeltaBinaryPacked
	EncodingDeltaByteArray
	EncodingDeltaLengthByteArray
)

type ErrorType int32

const (
	ErrorTypeUnauthorized     ErrorType = 1
	ErrorTypeInvalidOssBucket ErrorType = 2
)

type TaskSchema struct {
	ColumnName    string
	OssColumnName string
	Type          ParquetDataType
	Encode        Encoding
	TypeExtend    string
}

type TaskSyncStat struct {
	TaskSyncPhase        TaskSyncPhase
	CurrentSyncTimestamp int64
	ErrorCode            ErrorType
	Detail               string
}

type TaskSyncPhase int32

const (
	TaskInitStat TaskSyncPhase = iota
	TaskBaseStat
	TaskIncStat
)

type OSSTaskConfig struct {
	OssPrefix       string
	OssBucket       string
	OssEndpoint     string
	OssRoleName     string
	EventTimeColumn *EventColumn
	Format          Format
	Schema          []*TaskSchema
}

type CreateDeliveryTaskRequest struct {
	TableName  string
	TaskName   string
	TaskType   TaskType
	TaskConfig *OSSTaskConfig
}

type CreateDeliveryTaskResponse struct {
	ResponseInfo
}

func toTaskPbConfig(config *OSSTaskConfig) *otsprotocol.OSSTaskConfig {
	pbSchema := make([]*otsprotocol.ParquetSchema, len(config.Schema))
	for i, schema := range config.Schema {
		pbSchema[i] = &otsprotocol.ParquetSchema{
			ColumnName:    &schema.ColumnName,
			OssColumnName: &schema.OssColumnName,
			Type:          otsprotocol.ParquetSchema_DataType(schema.Type).Enum(),
			Encode:        otsprotocol.Encoding(schema.Encode).Enum(),
		}
	}
	var eventColumn *otsprotocol.EventColumn
	if config.EventTimeColumn != nil {
		eventColumn = &otsprotocol.EventColumn{
			ColumnName: &config.EventTimeColumn.ColumnName,
			TimeFormat: otsprotocol.EventColumnEventTimeFormat(config.EventTimeColumn.TimeFormat).Enum(),
		}
	}
	return &otsprotocol.OSSTaskConfig{
		OssPrefix:       &config.OssPrefix,
		OssBucket:       &config.OssBucket,
		OssEndpoint:     &config.OssEndpoint,
		OssStsRole:      &config.OssRoleName,
		EventTimeColumn: eventColumn,
		Format:          otsprotocol.Format(config.Format).Enum(),
		Schema:          pbSchema,
	}
}

type DeleteDeliveryTaskRequest struct {
	TableName string
	TaskName  string
}

type DeleteDeliveryTaskResponse struct {
	ResponseInfo
}

type ListDeliveryTaskRequest struct {
	TableName string
}

type DeliveryTaskInfo struct {
	TableName string
	TaskName  string
	TaskType  TaskType
}

type ListDeliveryTaskResponse struct {
	Tasks []*DeliveryTaskInfo
	ResponseInfo
}

type DescribeDeliveryTaskRequest struct {
	TableName string
	TaskName  string
}

type DescribeDeliveryTaskResponse struct {
	TaskConfig   *OSSTaskConfig
	TaskSyncStat *TaskSyncStat
	TaskType     TaskType
	ResponseInfo
}

func toOSSTaskConfig(pbConf *otsprotocol.OSSTaskConfig) *OSSTaskConfig {
	var (
		eventColumn *EventColumn
		schemas     []*TaskSchema = make([]*TaskSchema, len(pbConf.Schema))
	)
	if pbConf.EventTimeColumn != nil {
		eventColumn = &EventColumn{
			ColumnName: pbConf.EventTimeColumn.GetColumnName(),
			TimeFormat: TimeFormat(pbConf.EventTimeColumn.GetTimeFormat()),
		}
	}
	for i, schema := range pbConf.Schema {
		schemas[i] = &TaskSchema{
			ColumnName:    schema.GetColumnName(),
			OssColumnName: schema.GetOssColumnName(),
			Type:          ParquetDataType(schema.GetType()),
			Encode:        Encoding(schema.GetEncode()),
			TypeExtend:    schema.GetTypeExtend(),
		}
	}
	return &OSSTaskConfig{
		OssPrefix:       pbConf.GetOssPrefix(),
		OssBucket:       pbConf.GetOssBucket(),
		OssEndpoint:     pbConf.GetOssEndpoint(),
		OssRoleName:     pbConf.GetOssStsRole(),
		EventTimeColumn: eventColumn,
		Format:          Format(pbConf.GetFormat()),
		Schema:          schemas,
	}
}

func toTaskSyncStat(stat *otsprotocol.TaskSyncStat) *TaskSyncStat {
	if stat == nil {
		return new(TaskSyncStat)
	}
	syncStat := &TaskSyncStat{
		CurrentSyncTimestamp: stat.GetCurrentSyncTimestamp(),
		Detail:               stat.GetDetail(),
	}
	if stat.ErrorCode != nil {
		syncStat.ErrorCode = ErrorType(stat.GetErrorCode())
	}
	if stat.TaskSyncPhase != nil {
		syncStat.TaskSyncPhase = TaskSyncPhase(stat.GetTaskSyncPhase())
	}
	return syncStat
}
