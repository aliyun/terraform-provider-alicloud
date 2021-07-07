package datahub

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "reflect"
    "strconv"
)

type AuthMode string

const (
    AK  AuthMode = "ak"
    STS AuthMode = "sts"
)

type ConnectorType string

const (
    SinkOdps    ConnectorType = "sink_odps"
    SinkOss     ConnectorType = "sink_oss"
    SinkEs      ConnectorType = "sink_es"
    SinkAds     ConnectorType = "sink_ads"
    SinkMysql   ConnectorType = "sink_mysql"
    SinkFc      ConnectorType = "sink_fc"
    SinkOts     ConnectorType = "sink_ots"
    SinkDatahub ConnectorType = "sink_datahub"
)

func (ct *ConnectorType) String() string {
    return string(*ct)
}

func validateConnectorType(ct ConnectorType) bool {
    switch ct {
    case SinkOdps, SinkOss, SinkEs, SinkAds, SinkMysql, SinkFc, SinkOts, SinkDatahub:
        return true
    default:
        return false
    }
}

type ConnectorState string

const (
    ConnectorStopped ConnectorState = "CONNECTOR_PAUSED"
    ConnectorRunning ConnectorState = "CONNECTOR_RUNNING"
)

func validateConnectorState(ct ConnectorState) bool {
    switch ct {
    case ConnectorStopped, ConnectorRunning:
        return true
    default:
        return false
    }
}

type PartitionMode string

const (
    UserDefineMode PartitionMode = "USER_DEFINE"
    SystemTimeMode PartitionMode = "SYSTEM_TIME"
    EventTimeMode  PartitionMode = "EVENT_TIME"
)

func (pm *PartitionMode) String() string {
    return string(*pm)
}

func NewPartitionConfig() *PartitionConfig {
    pc := &PartitionConfig{
        ConfigMap: make([]map[string]string, 0, 0),
    }
    return pc
}

type PartitionConfig struct {
    ConfigMap []map[string]string
}

func (pc *PartitionConfig) AddConfig(key, value string) {
    m := map[string]string{
        key: value,
    }
    pc.ConfigMap = append(pc.ConfigMap, m)
}

func (pc *PartitionConfig) MarshalJSON() ([]byte, error) {
    if pc == nil || len(pc.ConfigMap) == 0 {
        return nil, nil
    }
    buf := &bytes.Buffer{}
    buf.Write([]byte{'{'})

    length := len(pc.ConfigMap)
    for i, m := range pc.ConfigMap {
        for k, v := range m {
            if _, err := fmt.Fprintf(buf, "\"%s\":\"%s\"", k, v); err != nil {
                return nil, errors.New(fmt.Sprintf("partition config is invalid"))
            }
        }
        if i < length-1 {
            buf.WriteByte(',')
        }
    }
    buf.WriteByte('}')

    return buf.Bytes(), nil
}

func (pc *PartitionConfig) UnmarshalJSON(data []byte) error {
    //the data is "xxxxxx",should convert to xxxx, remove the ""
    var str *string = new(string)
    if err := json.Unmarshal(data, str); err != nil {
        return err
    }

    confParser := make([]map[string]string, 0)
    if err := json.Unmarshal([]byte(*str), &confParser); err != nil {
        return err
    }
    confMap := make([]map[string]string, len(confParser))

    //convert {"key":"ds","value":"%Y%m%d",...} to {"ds":"%Y%m%d",...}
    for i, m := range confParser {
        confMap[i] = map[string]string{
            m["key"]: m["value"],
        }
    }
    pc.ConfigMap = confMap
    return nil
}

/** ODPS CONFIG **/
type SinkOdpsConfig struct {
    Endpoint        string          `json:"OdpsEndpoint"`
    Project         string          `json:"Project"`
    Table           string          `json:"Table"`
    AccessId        string          `json:"AccessId"`
    AccessKey       string          `json:"AccessKey"`
    TimeRange       int             `json:"TimeRange"`
    TimeZone        string          `json:"TimeZone"`
    PartitionMode   PartitionMode   `json:"PartitionMode"`
    PartitionConfig PartitionConfig `json:"PartitionConfig"`
    TunnelEndpoint  string          `json:"TunnelEndpoint,omitempty"`
}

func marshalCreateOdpsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkOdpsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOdpsConfig", reflect.TypeOf(ccr.Config)))
    }
    ct := &struct {
        Action       string         `json:"Action"`
        Type         string         `json:"Type"`
        ColumnFields []string       `json:"ColumnFields"`
        Config       SinkOdpsConfig `json:"Config"`
    }{
        Action:       ccr.Action,
        Type:         ccr.Type.String(),
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetOdpsConnector(gcr *GetConnectorResult, data []byte) error {

    //the api return TimeRange is string, so need to convert to int64
    type SinkOdpsConfigHelper struct {
        SinkOdpsConfig
        TimeRange string `json:"TimeRange"`
    }
    ct := &struct {
        GetConnectorResult
        Config SinkOdpsConfigHelper `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }
    //if gcr == nil {
    //    *gcr = GetConnectorResult{}
    //}
    *gcr = ct.GetConnectorResult

    soConf := ct.Config.SinkOdpsConfig
    t, err := strconv.Atoi(ct.Config.TimeRange)
    if err != nil {
        return err
    }
    soConf.TimeRange = t
    gcr.Config = soConf
    return nil
}

func marshalUpdateOdpsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkOdpsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOdpsConfig", reflect.TypeOf(ucr.Config)))
    }
    ct := &struct {
        Action string         `json:"Action"`
        Config SinkOdpsConfig `json:"Config"`
    }{
        Action: ucr.Action,
        Config: soConf,
    }
    return json.Marshal(ct)
}

/*  Oss Config */
type SinkOssConfig struct {
    Endpoint   string   `json:"Endpoint"`
    Bucket     string   `json:"Bucket"`
    Prefix     string   `json:"Prefix"`
    TimeFormat string   `json:"TimeFormat"`
    TimeRange  int      `json:"TimeRange"`
    AuthMode   AuthMode `json:"AuthMode"`
    AccessId   string   `json:"AccessId"`
    AccessKey  string   `json:"AccessKey"`
}

func marshalCreateOssConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkOssConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOssConfig", reflect.TypeOf(ccr.Config)))
    }

    ct := &struct {
        Action       string        `json:"Action"`
        Type         ConnectorType `json:"Type"`
        ColumnFields []string      `json:"ColumnFields"`
        Config       SinkOssConfig `json:"Config"`
    }{
        Action:       "create",
        Type:         ccr.Type,
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetOssConnector(gcr *GetConnectorResult, data []byte) error {
    type SinkOssConfigHelper struct {
        SinkOssConfig
        TimeRange string `json:"TimeRange"`
    }
    ct := &struct {
        GetConnectorResult
        Config SinkOssConfigHelper `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }

    *gcr = ct.GetConnectorResult
    soConf := ct.Config.SinkOssConfig
    t, err := strconv.Atoi(ct.Config.TimeRange)
    if err != nil {
        return err
    }
    soConf.TimeRange = t
    gcr.Config = soConf
    return nil
}

func marshalUpdateOssConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkOssConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOssConfig", reflect.TypeOf(ucr.Config)))
    }

    ct := &struct {
        Action string        `json:"Action"`
        Config SinkOssConfig `json:"Config"`
    }{
        Action: "create",
        Config: soConf,
    }
    return json.Marshal(ct)
}

/*  mysql Config */
type SinkMysqlConfig struct {
    Host     string     `json:"Host"`
    Port     string     `json:"Port"`
    Database string     `json:"Database"`
    Table    string     `json:"Table"`
    User     string     `json:"User"`
    Password string     `json:"Password"`
    Ignore   InsertMode `json:"Ignore"`
}

type InsertMode string

const (
    IGNORE    InsertMode = "true"
    OVERWRITE InsertMode = "false"
)

func marshalCreateMysqlConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkMysqlConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkMysqlConfig", reflect.TypeOf(ccr.Config)))
    }

    ct := &struct {
        Action       string          `json:"Action"`
        Type         ConnectorType   `json:"Type"`
        ColumnFields []string        `json:"ColumnFields"`
        Config       SinkMysqlConfig `json:"Config"`
    }{
        Action:       "create",
        Type:         ccr.Type,
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetMysqlConnector(gcr *GetConnectorResult, data []byte) error {
    ct := &struct {
        GetConnectorResult
        Config SinkMysqlConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }

    *gcr = ct.GetConnectorResult
    gcr.Config = ct.Config
    return nil
}

func marshalUpdateMysqlConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkMysqlConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkMysqlConfig", reflect.TypeOf(ucr.Config)))
    }

    ct := &struct {
        Action string          `json:"Action"`
        Config SinkMysqlConfig `json:"Config"`
    }{
        Action: "create",
        Config: soConf,
    }
    return json.Marshal(ct)
}

/*  Ads Config */
type SinkAdsConfig struct {
    SinkMysqlConfig
}

func marshalCreateAdsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkAdsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkAdsConfig", reflect.TypeOf(ccr.Config)))
    }

    ct := &struct {
        Action       string        `json:"Action"`
        Type         ConnectorType `json:"Type"`
        ColumnFields []string      `json:"ColumnFields"`
        Config       SinkAdsConfig `json:"Config"`
    }{
        Action:       "create",
        Type:         ccr.Type,
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetAdsConnector(gcr *GetConnectorResult, data []byte) error {
    ct := &struct {
        GetConnectorResult
        Config SinkMysqlConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }

    *gcr = ct.GetConnectorResult
    gcr.Config = ct.Config
    return nil
}

func marshalUpdateAdsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkMysqlConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkAdsConfig", reflect.TypeOf(ucr.Config)))
    }

    ct := &struct {
        Action string          `json:"Action"`
        Config SinkMysqlConfig `json:"Config"`
    }{
        Action: "create",
        Config: soConf,
    }
    return json.Marshal(ct)
}

/*  datahub Config */
type SinkDatahubConfig struct {
    Endpoint  string   `json:"Endpoint"`
    Project   string   `json:"Project"`
    Topic     string   `json:"Topic"`
    AuthMode  AuthMode `json:"AuthMode"`
    AccessId  string   `json:"AccessId"`
    AccessKey string   `json:"AccessKey"`
}

func marshalCreateDatahubConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkDatahubConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkDatahubConfig", reflect.TypeOf(ccr.Config)))
    }

    ct := &struct {
        Action       string            `json:"Action"`
        Type         ConnectorType     `json:"Type"`
        ColumnFields []string          `json:"ColumnFields"`
        Config       SinkDatahubConfig `json:"Config"`
    }{
        Action:       "create",
        Type:         ccr.Type,
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetDatahubConnector(gcr *GetConnectorResult, data []byte) error {
    ct := &struct {
        GetConnectorResult
        Config SinkDatahubConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }

    *gcr = ct.GetConnectorResult
    gcr.Config = ct.Config
    return nil
}

func marshalUpdateDatahubConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkDatahubConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkDatahubConfig", reflect.TypeOf(ucr.Config)))
    }

    ct := &struct {
        Action string            `json:"Action"`
        Config SinkDatahubConfig `json:"Config"`
    }{
        Action: "create",
        Config: soConf,
    }
    return json.Marshal(ct)
}

/*  ES Config */
type SinkEsConfig struct {
    Index      string   `json:"Index"`
    Endpoint   string   `json:"Endpoint"`
    User       string   `json:"User"`
    Password   string   `json:"Password"`
    IDFields   []string `json:"IDFields"`
    TypeFields []string `json:"TypeFields"`
    ProxyMode  bool     `json:"ProxyMode"`
}

func marshalCreateEsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkEsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkEsConfig", reflect.TypeOf(ccr.Config)))
    }

    ct := &struct {
        Action       string        `json:"Action"`
        Type         ConnectorType `json:"Type"`
        ColumnFields []string      `json:"ColumnFields"`
        Config       SinkEsConfig  `json:"Config"`
    }{
        Action:       "create",
        Type:         ccr.Type,
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetEsConnector(gcr *GetConnectorResult, data []byte) error {
    ct := &struct {
        GetConnectorResult
        Config SinkEsConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }

    *gcr = ct.GetConnectorResult
    gcr.Config = ct.Config
    return nil
}

func marshalUpdateEsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkEsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkEsConfig", reflect.TypeOf(ucr.Config)))
    }

    ct := &struct {
        Action string       `json:"Action"`
        Config SinkEsConfig `json:"Config"`
    }{
        Action: "create",
        Config: soConf,
    }
    return json.Marshal(ct)
}

/*  FC Config */
type SinkFcConfig struct {
    Endpoint  string   `json:"Endpoint"`
    Service   string   `json:"Service"`
    Function  string   `json:"Function"`
    AuthMode  AuthMode `json:"AuthMode"`
    AccessId  string   `json:"AccessId"`
    AccessKey string   `json:"AccessKey"`
}

func marshalCreateFcConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkFcConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkFcConfig", reflect.TypeOf(ccr.Config)))
    }

    ct := &struct {
        Action       string        `json:"Action"`
        Type         ConnectorType `json:"Type"`
        ColumnFields []string      `json:"ColumnFields"`
        Config       SinkFcConfig  `json:"Config"`
    }{
        Action:       "create",
        Type:         ccr.Type,
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetFcConnector(gcr *GetConnectorResult, data []byte) error {
    ct := &struct {
        GetConnectorResult
        Config SinkFcConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }

    *gcr = ct.GetConnectorResult
    gcr.Config = ct.Config
    return nil
}

func marshalUpdateFcConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkFcConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkFcConfig", reflect.TypeOf(ucr.Config)))
    }

    ct := &struct {
        Action string       `json:"Action"`
        Config SinkFcConfig `json:"Config"`
    }{
        Action: "create",
        Config: soConf,
    }
    return json.Marshal(ct)
}

/*  Ots Config */
type SinkOtsConfig struct {
    Endpoint     string   `json:"Endpoint"`
    InstanceName string   `json:"InstanceName"`
    TableName    string   `json:"TableName"`
    AuthMode     AuthMode `json:"AuthMode"`
    AccessId     string   `json:"AccessId"`
    AccessKey    string   `json:"AccessKey"`
}

func marshalCreateOtsConnector(ccr *CreateConnectorRequest) ([]byte, error) {
    soConf, ok := ccr.Config.(SinkOtsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkOtsConfig", reflect.TypeOf(ccr.Config)))
    }

    ct := &struct {
        Action       string        `json:"Action"`
        Type         ConnectorType `json:"Type"`
        ColumnFields []string      `json:"ColumnFields"`
        Config       SinkOtsConfig `json:"Config"`
    }{
        Action:       "create",
        Type:         ccr.Type,
        ColumnFields: ccr.ColumnFields,
        Config:       soConf,
    }
    return json.Marshal(ct)
}

func unmarshalGetOtsConnector(gcr *GetConnectorResult, data []byte) error {
    ct := &struct {
        GetConnectorResult
        Config SinkOtsConfig `json:"Config"`
    }{}

    if err := json.Unmarshal(data, ct); err != nil {
        return err
    }

    *gcr = ct.GetConnectorResult
    gcr.Config = ct.Config
    return nil
}

func marshalUpdateOtsConnector(ucr *UpdateConnectorRequest) ([]byte, error) {
    soConf, ok := ucr.Config.(SinkOtsConfig)
    if !ok {
        return nil, NewInvalidParameterErrorWithMessage(fmt.Sprintf("config type error,your input config type is %s,should be SinkMysqlConfig", reflect.TypeOf(ucr.Config)))
    }

    ct := &struct {
        Action string        `json:"Action"`
        Config SinkOtsConfig `json:"Config"`
    }{
        Action: "create",
        Config: soConf,
    }
    return json.Marshal(ct)
}

type ConnectorOffset struct {
    Timestamp int64 `json:"Timestamp"`
    Sequence  int64 `json:"Sequence"`
}

type ConnectorShardState string

const (
    Created   ConnectorShardState = "CONTEXT_PLANNED"
    Eexcuting ConnectorShardState = "CONTEXT_EXECUTING"
    Stopped   ConnectorShardState = "CONTEXT_PAUSED"
    Finished  ConnectorShardState = "CONTEXT_FINISHED"
)

type ConnectorShardStatusEntry struct {
    StartSequence    int64               `json:"StartSequence"`
    EndSequence      int64               `json:"EndSequence"`
    CurrentSequence  int64               `json:"CurrentSequence"`
    CurrentTimestamp int64               `json:"CurrentTimestamp"`
    UpdateTime       int64               `json:"UpdateTime"`
    State            ConnectorShardState `json:"State"`
    LastErrorMessage string              `json:"LastErrorMessage"`
    DiscardCount     int64               `json:"DiscardCount"`
    DoneTime         int64               `json:"DoneTime"`
    WorkerAddress    string              `json:"WorkerAddress"`
}
