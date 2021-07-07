package datahub

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/golang/protobuf/proto"
    "net/http"

    pbmodel "github.com/aliyun/aliyun-datahub-sdk-go/datahub/pbmodel"
    "github.com/aliyun/aliyun-datahub-sdk-go/datahub/util"
)

// for the common response and detect error
type commonResponseResult struct {
    // StatusCode http return code
    StatusCode int

    // RequestId examples request id return by server
    RequestId string
}

func newCommonResponseResult(code int, header *http.Header, body []byte) (*commonResponseResult, error) {
    result := &commonResponseResult{
        StatusCode: code,
        RequestId:  header.Get(httpHeaderRequestId),
    }
    var err error
    switch {
    case code >= 400:
        var datahubErr DatahubClientError
        if err = json.Unmarshal(body, &datahubErr); err != nil {
            return nil, err
        }
        err = errorHandler(code, result.RequestId, datahubErr.Code, datahubErr.Message)
    default:
        err = nil
    }
    return result, err
}

//  the result of ListProject
type ListProjectResult struct {
    ProjectNames []string `json:"ProjectNames"`
}

// convert the response body to ListProjectResult
func NewListProjectResult(data []byte) (*ListProjectResult, error) {
    lpr := &ListProjectResult{}
    lpr.ProjectNames = make([]string, 0, 0)
    if err := json.Unmarshal(data, lpr); err != nil {
        return nil, err
    }
    return lpr, nil
}

// the result of GetProject
type GetProjectResult struct {
    ProjectName    string
    CreateTime     int64  `json:"CreateTime"`
    LastModifyTime int64  `json:"LastModifyTime"`
    Comment        string `json"Comment"`
}

// convert the response body to GetProjectResult
func NewGetProjectResult(data []byte) (*GetProjectResult, error) {
    gpr := &GetProjectResult{}
    if err := json.Unmarshal(data, gpr); err != nil {
        return nil, err
    }
    return gpr, nil
}

type ListTopicResult struct {
    TopicNames [] string `json:"TopicNames"`
}

func NewListTopicResult(data []byte) (*ListTopicResult, error) {
    lt := &ListTopicResult{}
    if err := json.Unmarshal(data, lt); err != nil {
        return nil, err
    }
    return lt, nil
}

type GetTopicResult struct {
    ProjectName    string
    TopicName      string
    ShardCount     int           `json:"ShardCount"`
    LifeCycle      int           `json:"LifeCycle"`
    RecordType     RecordType    `json:"RecordType"`
    RecordSchema   *RecordSchema `json:"RecordSchema"`
    Comment        string        `json:"Comment"`
    CreateTime     int64         `json:"CreateTime"`
    LastModifyTime int64         `json:"LastModifyTime"`
}

// for deserialize the RecordSchema
func (gtr *GetTopicResult) UnmarshalJSON(data []byte) error {
    msg := &struct {
        ShardCount     int        `json:"ShardCount"`
        LifeCycle      int        `json:"LifeCycle"`
        RecordType     RecordType `json:"RecordType"`
        RecordSchema   string     `json:"RecordSchema"`
        Comment        string     `json:"Comment"`
        CreateTime     int64      `json:"CreateTime"`
        LastModifyTime int64      `json:"LastModifyTime"`
    }{}
    if err := json.Unmarshal(data, msg); err != nil {
        return err
    }

    gtr.ShardCount = msg.ShardCount
    gtr.LifeCycle = msg.LifeCycle
    gtr.RecordType = msg.RecordType
    gtr.Comment = msg.Comment
    gtr.CreateTime = msg.CreateTime
    gtr.LastModifyTime = msg.LastModifyTime
    if msg.RecordType == TUPLE {
        rs := &RecordSchema{}
        if err := json.Unmarshal([]byte(msg.RecordSchema), rs); err != nil {
            return err
        }
        for idx := range rs.Fields {
            rs.Fields[idx].AllowNull = !rs.Fields[idx].AllowNull
        }
        gtr.RecordSchema = rs
    }
    return nil
}

func NewGetTopicResult(data []byte) (*GetTopicResult, error) {
    gr := &GetTopicResult{}
    if err := json.Unmarshal(data, gr); err != nil {
        return nil, err
    }
    return gr, nil
}

type ListShardResult struct {
    Shards []ShardEntry `json:"Shards"`
}

func NewListShardResult(data []byte) (*ListShardResult, error) {
    lsr := &ListShardResult{}
    if err := json.Unmarshal(data, lsr); err != nil {
        return nil, err
    }
    return lsr, nil
}

type SplitShardResult struct {
    NewShards []ShardEntry `json:"NewShards"`
}

func NewSplitShardResult(data []byte) (*SplitShardResult, error) {
    ssr := &SplitShardResult{}
    if err := json.Unmarshal(data, ssr); err != nil {
        return nil, err
    }
    return ssr, nil
}

type MergeShardResult struct {
    ShardId      string `json:"ShardId"`
    BeginHashKey string `json:"BeginHashKey"`
    EndHashKey   string `json:"EndHashKey"`
}

func NewMergeShardResult(data []byte) (*MergeShardResult, error) {
    ssr := &MergeShardResult{}
    if err := json.Unmarshal(data, ssr); err != nil {
        return nil, err
    }
    return ssr, nil
}

type GetCursorResult struct {
    Cursor     string `json:"Cursor"`
    RecordTime int64  `json:"RecordTime"`
    Sequence   int64  `json:"Sequence"`
}

func NewGetCursorResult(data []byte) (*GetCursorResult, error) {
    gcr := &GetCursorResult{}
    if err := json.Unmarshal(data, gcr); err != nil {
        return nil, err
    }
    return gcr, nil
}

type PutRecordsResult struct {
    FailedRecordCount int            `json:"FailedRecordCount"`
    FailedRecords     []FailedRecord `json:"FailedRecords"`
}

func NewPutRecordsResult(data []byte) (*PutRecordsResult, error) {
    prr := &PutRecordsResult{}
    if err := json.Unmarshal(data, prr); err != nil {
        return nil, err
    }
    return prr, nil
}

func NewPutPBRecordsResult(data []byte) (*PutRecordsResult, error) {
    pr := &PutRecordsResult{}
    data, err := util.UnwrapMessage(data)
    if err != nil {
        return nil, err
    }
    prr := &pbmodel.PutRecordsResponse{}
    if err := proto.Unmarshal(data, prr); err != nil {
        return nil, err
    }

    pr.FailedRecordCount = int(*prr.FailedCount)
    if pr.FailedRecordCount > 0 {
        records := make([]FailedRecord, pr.FailedRecordCount)
        for idx, v := range prr.FailedRecords {
            records[idx].ErrorCode = *v.ErrorCode
            records[idx].ErrorMessage = *v.ErrorMessage
            records[idx].Index = int(*v.Index)
        }
        pr.FailedRecords = records
    }
    return pr, nil
}

type GetRecordsResult struct {
    NextCursor    string        `json:"NextCursor"`
    RecordCount   int           `json:"RecordCount"`
    StartSequence int64         `json:"StartSeq"`
    Records       []IRecord     `json:"Records"`
    RecordSchema  *RecordSchema `json:"-"`
}

func (grr *GetRecordsResult) UnmarshalJSON(data []byte) error {
    msg := &struct {
        NextCursor  string `json:"NextCursor"`
        RecordCount int    `json:"RecordCount"`
        Records     []*struct {
            SystemTime    int64                  `json:"SystemTime"`
            NextCursor    string                 `json:"NextCursor"`
            CurrentCursor string                 `json:"Cursor"`
            Sequence      int64                  `json:"Sequence"`
            Attributes    map[string]interface{} `json:"Attributes"`
            Data          interface{}            `json:"Data"`
        } `json:"Records"`
    }{}
    err := json.Unmarshal(data, msg)
    if err != nil {
        fmt.Printf("%v\n", err)
        return err
    }
    grr.NextCursor = msg.NextCursor
    grr.RecordCount = msg.RecordCount
    grr.Records = make([]IRecord, len(msg.Records))
    for idx, record := range msg.Records {
        switch dt := record.Data.(type) {
        case []interface{}, []string:
            if grr.RecordSchema == nil {
                return errors.New("tuple record type must set record schema")
            }
            grr.Records[idx] = NewTupleRecord(grr.RecordSchema, record.SystemTime)
        case string:
            grr.Records[idx] = NewBlobRecord([]byte(dt), record.SystemTime)
        default:
            return errors.New(fmt.Sprintf("illegal record data type[%T]", dt))
        }
        if err := grr.Records[idx].FillData(record.Data); err != nil {
            return err
        }
        for key, val := range record.Attributes {
            grr.Records[idx].SetAttribute(key, val)
        }
        br := BaseRecord{
            SystemTime: msg.Records[idx].SystemTime,
            NextCursor: msg.Records[idx].NextCursor,
            Cursor:     msg.Records[idx].CurrentCursor,
            Sequence:   msg.Records[idx].Sequence,
            Attributes: msg.Records[idx].Attributes,
        }
        grr.Records[idx].SetBaseRecord(br)
    }
    return nil
}

func NewGetRecordsResult(data []byte, schema *RecordSchema) (*GetRecordsResult, error) {
    grr := &GetRecordsResult{
        RecordSchema: schema,
    }
    if err := json.Unmarshal(data, grr); err != nil {
        return nil, err
    }
    return grr, nil
}

func NewGetPBRecordsResult(data []byte, schema *RecordSchema) (*GetRecordsResult, error) {
    data, err := util.UnwrapMessage(data)
    if err != nil {
        return nil, err
    }
    grr := &pbmodel.GetRecordsResponse{}
    if err := proto.Unmarshal(data, grr); err != nil {
        return nil, err
    }

    result := &GetRecordsResult{
        RecordSchema: schema,
    }
    if grr.NextCursor != nil {
        result.NextCursor = *(grr.NextCursor)
    }
    if grr.StartSequence != nil {
        result.StartSequence = *grr.StartSequence
    }
    if grr.RecordCount != nil {
        result.RecordCount = int(*grr.RecordCount)
        if result.RecordCount > 0 {
            result.Records = make([]IRecord, result.RecordCount)
            for idx, record := range grr.Records {
                //Tuple topic

                if result.RecordSchema != nil {
                    tr := NewTupleRecord(result.RecordSchema, *record.SystemTime)
                    if err := fillTupleData(tr, record); err != nil {
                        return nil, err
                    }
                    result.Records[idx] = tr
                } else {
                    br := NewBlobRecord(record.Data.Data[0].Value, *record.SystemTime)
                    if err := fillBlobData(br, record); err != nil {
                        return nil, err
                    }
                    result.Records[idx] = br
                }
            }
        }
    }
    return result, nil
}
func fillTupleData(tr *TupleRecord, recordEntry *pbmodel.RecordEntry) error {
    if recordEntry.ShardId != nil {
        tr.ShardId = *recordEntry.ShardId
    }
    if recordEntry.HashKey != nil {
        tr.HashKey = *recordEntry.HashKey
    }
    if recordEntry.PartitionKey != nil {
        tr.Sequence = *recordEntry.Sequence
    }
    if recordEntry.Cursor != nil {
        tr.Cursor = *recordEntry.Cursor
    }
    if recordEntry.NextCursor != nil {
        tr.NextCursor = *recordEntry.NextCursor
    }
    if recordEntry.Sequence != nil {
        tr.Sequence = *recordEntry.Sequence
    }
    if recordEntry.SystemTime != nil {
        tr.SystemTime = *recordEntry.SystemTime
    }
    if recordEntry.Attributes != nil {
        for _, pair := range recordEntry.Attributes.Attributes {
            tr.Attributes[*pair.Key] = pair.Value
        }
    }
    data := recordEntry.Data.Data

    for idx, v := range data {
        if v.Value != nil {
            tv, err := castValueFromString(string(v.Value), tr.RecordSchema.Fields[idx].Type)
            if err != nil {
                return err
            }
            tr.Values[idx] = tv
        }
    }
    return nil

}

func fillBlobData(br *BlobRecord, recordEntry *pbmodel.RecordEntry) error {
    if recordEntry.ShardId != nil {
        br.ShardId = *recordEntry.ShardId
    }
    if recordEntry.HashKey != nil {
        br.HashKey = *recordEntry.HashKey
    }
    if recordEntry.PartitionKey != nil {
        br.Sequence = *recordEntry.Sequence
    }
    if recordEntry.Cursor != nil {
        br.Cursor = *recordEntry.Cursor
    }
    if recordEntry.NextCursor != nil {
        br.NextCursor = *recordEntry.NextCursor
    }
    if recordEntry.Sequence != nil {
        br.Sequence = *recordEntry.Sequence
    }
    if recordEntry.SystemTime != nil {
        br.SystemTime = *recordEntry.SystemTime
    }
    if recordEntry.Attributes != nil {
        for _, pair := range recordEntry.Attributes.Attributes {
            br.Attributes[*pair.Key] = pair.Value
        }
    }
    br.RawData = recordEntry.Data.Data[0].Value
    br.StoreData = string(br.RawData)
    return nil
}

type GetMeterInfoResult struct {
    ActiveTime int64 `json:"ActiveTime"`
    Storage    int64 `json:"Storage"`
}

func NewGetMeterInfoResult(data []byte) (*GetMeterInfoResult, error) {
    gmir := &GetMeterInfoResult{}
    if err := json.Unmarshal(data, gmir); err != nil {
        return nil, err
    }
    return gmir, nil
}

type CreateConnectorResult struct {
    ConnectorId string `json:"ConnectorId"`
}

func NewCreateConnectorResult(data []byte) (*CreateConnectorResult, error) {
    ccr := &CreateConnectorResult{}
    if err := json.Unmarshal(data, ccr); err != nil {
        return nil, err
    }
    return ccr, nil
}

type GetConnectorResult struct {
    CreateTime     int64             `json:"CreateTime"`
    LastModifyTime int64             `json:"LastModifyTime"`
    ConnectorId    string            `json:"ConnectorId"`
    ClusterAddress string            `json:"ClusterAddress"`
    Type           ConnectorType     `json:"Type"`
    State          ConnectorState    `json:"State"`
    ColumnFields   []string          `json:"ColumnFields"`
    ExtraConfig    map[string]string `json:"ExtraInfo"`
    Creator        string            `json:"Creator"`
    Owner          string            `json:"Owner"`
    Config         interface{}       `json:"Config"`
}

func NewGetConnectorResult(data []byte) (*GetConnectorResult, error) {
    gcr := &GetConnectorResult{}
    cType := &struct {
        Type ConnectorType `json:"Type"`
    }{}
    if err := json.Unmarshal(data, cType); err != nil {
        return nil, err
    }
    switch cType.Type {
    case SinkOdps:
        if err := unmarshalGetOdpsConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    case SinkOss:
        if err := unmarshalGetOssConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    case SinkEs:
        if err := unmarshalGetEsConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    case SinkAds:
        if err := unmarshalGetAdsConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    case SinkMysql:
        if err := unmarshalGetMysqlConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    case SinkFc:
        if err := unmarshalGetFcConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    case SinkOts:
        if err := unmarshalGetOtsConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    case SinkDatahub:
        if err := unmarshalGetDatahubConnector(gcr, data); err != nil {
            return nil, err
        }
        return gcr, nil
    default:
        return nil, errors.New(fmt.Sprintf("not support connector type %s", cType.Type.String()))

    }
}

type ListConnectorResult struct {
    ConnectorIds []string `json:"Connectors"`
}

func NewListConnectorResult(data []byte) (*ListConnectorResult, error) {
    lcr := &ListConnectorResult{}
    if err := json.Unmarshal(data, lcr); err != nil {
        return nil, err
    }
    return lcr, nil
}

type GetConnectorDoneTimeResult struct {
    DoneTime int64  `json:"DoneTime"`
    TimeZone string `json:"TimeZone"`
}

func NewGetConnectorDoneTimeResult(data []byte) (*GetConnectorDoneTimeResult, error) {
    gcdt := &GetConnectorDoneTimeResult{}
    if err := json.Unmarshal(data, gcdt); err != nil {
        return nil, err
    }
    return gcdt, nil
}

type GetConnectorShardStatusResult struct {
    ShardStatus map[string]ConnectorShardStatusEntry `json:"ShardStatusInfos"`
}

func NewGetConnectorShardStatusResult(data []byte) (*GetConnectorShardStatusResult, error) {
    gcss := &GetConnectorShardStatusResult{}
    if err := json.Unmarshal(data, gcss); err != nil {
        return nil, err
    }
    return gcss, nil
}

type CreateSubscriptionResult struct {
    SubId string `json:"SubId"`
}

func NewCreateSubscriptionResult(data []byte) (*CreateSubscriptionResult, error) {
    csr := &CreateSubscriptionResult{}
    if err := json.Unmarshal(data, csr); err != nil {
        return nil, err
    }
    return csr, nil
}

type GetSubscriptionResult struct {
    SubscriptionEntry
}

func NewGetSubscriptionResult(data []byte) (*GetSubscriptionResult, error) {
    gsr := &GetSubscriptionResult{}
    if err := json.Unmarshal(data, gsr); err != nil {
        return nil, err
    }
    return gsr, nil
}

type ListSubscriptionResult struct {
    TotalCount    int64               `json:"TotalCount"`
    Subscriptions []SubscriptionEntry `json:"Subscriptions"`
}

func NewListSubscriptionResult(data []byte) (*ListSubscriptionResult, error) {
    lsr := &ListSubscriptionResult{}
    if err := json.Unmarshal(data, lsr); err != nil {
        return nil, err
    }
    return lsr, nil
}

type OpenSubscriptionSessionResult struct {
    Offsets map[string]SubscriptionOffset `json:"Offsets"`
}

func NewOpenSubscriptionSessionResult(data []byte) (*OpenSubscriptionSessionResult, error) {
    ossr := &OpenSubscriptionSessionResult{}
    if err := json.Unmarshal(data, ossr); err != nil {
        return nil, err
    }
    return ossr, nil
}

type GetSubscriptionOffsetResult struct {
    Offsets map[string]SubscriptionOffset `json:"Offsets"`
}

func NewGetSubscriptionOffsetResult(data []byte) (*GetSubscriptionOffsetResult, error) {
    gsor := &GetSubscriptionOffsetResult{}
    if err := json.Unmarshal(data, gsor); err != nil {
        return nil, err
    }
    return gsor, nil
}

type HeartbeatResult struct {
    PlanVersion int64    `json:"PlanVersion"`
    ShardList   []string `json:"ShardList"`
    TotalPlan   string   `json:"TotalPlan"`
}

func NewHeartbeatResult(data []byte) (*HeartbeatResult, error) {
    hr := &HeartbeatResult{}
    if err := json.Unmarshal(data, hr); err != nil {
        return nil, err
    }
    return hr, nil
}

type JoinGroupResult struct {
    ConsumerId     string `json:"ConsumerId"`
    VersionId      int64  `json:"VersionId"`
    SessionTimeout int64  `json:"SessionTimeout"`
}

func NewJoinGroupResult(data []byte) (*JoinGroupResult, error) {
    jgr := &JoinGroupResult{}
    if err := json.Unmarshal(data, jgr); err != nil {
        return nil, err
    }
    return jgr, nil
}
