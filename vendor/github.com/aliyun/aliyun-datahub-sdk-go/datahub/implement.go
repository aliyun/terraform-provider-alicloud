package datahub

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/aliyun/aliyun-datahub-sdk-go/datahub/util"
    "time"
)

type DataHub struct {
    Client *RestClient
}

// ListProjects list all projects
func (datahub *DataHub) ListProject() (*ListProjectResult, error) {
    path := projectsPath
    responseBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    return NewListProjectResult(responseBody)
}

// CreateProject create new project
func (datahub *DataHub) CreateProject(projectName, comment string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckComment(comment) {
        return NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    requestBody := &CreateProjectRequest{
        Comment: comment,
    }
    if _, err := datahub.Client.Post(path, requestBody); err != nil {
        return err
    }
    return nil
}

// UpdateProject update project
func (datahub *DataHub) UpdateProject(projectName, comment string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckComment(comment) {
        return NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    requestBody := &UpdateProjectRequest{
        Comment: comment,
    }
    if _, err := datahub.Client.Put(path, requestBody); err != nil {
        return err
    }
    return nil
}

// DeleteProject delete project
func (datahub *DataHub) DeleteProject(projectName string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    if _, err := datahub.Client.Delete(path); err != nil {
        return err
    }
    return nil
}

// GetProject get a project deatil named the given name
func (datahub *DataHub) GetProject(projectName string) (*GetProjectResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(projectPath, projectName)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }

    result, err := NewGetProjectResult(respBody)
    if err != nil {
        return nil, err
    }

    result.ProjectName = projectName
    return result, nil

}

func (datahub *DataHub) WaitAllShardsReady(projectName, topicName string) bool {
    return datahub.WaitAllShardsReadyWithTime(projectName, topicName, minWaitingTimeInMs/1000)
}

func (datahub *DataHub) WaitAllShardsReadyWithTime(projectName, topicName string, timeout int64) bool {
    ready := make(chan bool)
    if timeout > 0 {
        go func(timeout int64) {
            time.Sleep(time.Duration(timeout) * time.Second)
            ready <- false
        }(timeout)
    }
    go func(datahub DataHubApi) {
        for {
            ls, err := datahub.ListShard(projectName, topicName)
            shards := ls.Shards
            if err != nil {
                time.Sleep(1 * time.Microsecond)
                continue
            }
            ok := true
            for _, shard := range shards {
                switch shard.State {
                case ACTIVE, CLOSED:
                    continue
                default:
                    ok = false
                    break
                }
            }
            if ok {
                break
            }
        }
        ready <- true
    }(datahub)
    return <-ready
}

func (datahub *DataHub) CreateBlobTopic(projectName, topicName, comment string, shardCount, lifeCycle int) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckComment(comment) {
        return NewInvalidParameterErrorWithMessage(commentInvalid)
    }
    path := fmt.Sprintf(topicPath, projectName, topicName)
    ctr := &CreateTopicRequest{
        Action:      "create",
        ProjectName: projectName,
        TopicName:   topicName,
        ShardCount:  shardCount,
        Lifecycle:   lifeCycle,
        RecordType:  BLOB,
        Comment:     comment,
    }
    if _, err := datahub.Client.Post(path, ctr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) CreateTupleTopic(projectName, topicName, comment string, shardCount, lifeCycle int, recordSchema *RecordSchema) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckComment(comment) {
        return NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    ctr := &CreateTopicRequest{
        Action:       "create",
        ProjectName:  projectName,
        TopicName:    topicName,
        ShardCount:   shardCount,
        Lifecycle:    lifeCycle,
        RecordType:   TUPLE,
        RecordSchema: recordSchema,
        Comment:      comment,
    }

    if _, err := datahub.Client.Post(path, ctr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) UpdateTopic(projectName, topicName, comment string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckComment(comment) {
        return NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    ut := &UpdateTopicRequest{
        Comment: comment,
    }

    if _, err := datahub.Client.Put(path, ut); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) GetTopic(projectName, topicName string) (*GetTopicResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    result, err := NewGetTopicResult(respBody)

    if err != nil {
        return nil, err
    }
    result.ProjectName = projectName
    result.TopicName = topicName
    return result, nil
}

func (datahub *DataHub) DeleteTopic(projectName, topicName string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    path := fmt.Sprintf(topicPath, projectName, topicName)
    if _, err := datahub.Client.Delete(path); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) ListTopic(projectName string) (*ListTopicResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }

    path := fmt.Sprintf(topicsPath, projectName)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    return NewListTopicResult(respBody)
}

func (datahub *DataHub) ListShard(projectName, topicName string) (*ListShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    return NewListShardResult(respBody)
}

func (datahub *DataHub) SplitShard(projectName, topicName, shardId string) (*SplitShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    splitKey, err := generateSpliteKey(projectName, topicName, shardId, datahub)
    if err != nil {
        return nil, err
    }
    path := fmt.Sprintf(shardsPath, projectName, topicName)
    ssr := &SplitShardRequest{
        Action:   "split",
        ShardId:  shardId,
        SplitKey: splitKey,
    }

    respBody, err := datahub.Client.Post(path, ssr)
    if err != nil {
        return nil, err
    }
    return NewSplitShardResult(respBody)

}

func (datahub *DataHub) SplitShardBySplitKey(projectName, topicName, shardId, splitKey string) (*SplitShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    ssr := &SplitShardRequest{
        Action:   "split",
        ShardId:  shardId,
        SplitKey: splitKey,
    }

    respBody, err := datahub.Client.Post(path, ssr)
    if err != nil {
        return nil, err
    }
    return NewSplitShardResult(respBody)
}

func (datahub *DataHub) MergeShard(projectName, topicName, shardId, adjacentShardId string) (*MergeShardResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    if !util.CheckShardId(shardId) || !util.CheckShardId(adjacentShardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    mss := &MergeShardRequest{
        Action:          "merge",
        ShardId:         shardId,
        AdjacentShardId: adjacentShardId,
    }

    respBody, err := datahub.Client.Post(path, mss)
    if err != nil {
        return nil, err
    }
    return NewMergeShardResult(respBody)

}

func (datahub *DataHub) GetCursor(projectName, topicName, shardId string, ctype CursorType, param ...int64) (*GetCursorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }
    if len(param) > 1 {
        return nil, NewInvalidParameterErrorWithMessage(parameterNumInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    gcr := &GetCursorRequest{
        Action:     "cursor",
        CursorType: ctype,
    }

    switch ctype {
    case OLDEST, LATEST:
        if len(param) != 0 {
            return nil, NewInvalidParameterErrorWithMessage(parameterNumInvalid)
        }
    case SYSTEM_TIME:
        if len(param) != 1 {
            return nil, NewInvalidParameterErrorWithMessage(parameterNumInvalid)
        }
        gcr.SystemTime = param[0]
    case SEQUENCE:
        if len(param) != 1 {
            return nil, NewInvalidParameterErrorWithMessage(parameterNumInvalid)
        }
        gcr.Sequence = param[0]
    }

    respBody, err := datahub.Client.Post(path, gcr)
    if err != nil {
        return nil, err
    }
    return NewGetCursorResult(respBody)
}
func (datahub *DataHub) PutRecords(projectName, topicName string, records []IRecord) (*PutRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    prr := &PutRecordsRequest{
        Action:  "pub",
        Records: records,
    }
    respBody, err := datahub.Client.Post(path, prr)
    if err != nil {
        return nil, err
    }
    return NewPutRecordsResult(respBody)
}

func (datahub *DataHub) PutRecordsByShard(projectName, topicName, shardId string, records []IRecord) error {
    return errors.New("not support this method")
}

func (datahub *DataHub) GetTupleRecords(projectName, topicName, shardId, cursor string, limit int, recordSchema *RecordSchema) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    grr := &GetRecordRequest{
        Action: "sub",
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, err := datahub.Client.Post(path, grr)
    if err != nil {
        return nil, err
    }
    return NewGetRecordsResult(respBody, recordSchema)
}

func (datahub *DataHub) GetBlobRecords(projectName, topicName, shardId, cursor string, limit int) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    grr := &GetRecordRequest{
        Action: "sub",
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, err := datahub.Client.Post(path, grr)
    if err != nil {
        return nil, err
    }
    return NewGetRecordsResult(respBody, nil)
}

func (datahub *DataHub) AppendField(projectName, topicName string, field Field) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(topicPath, projectName, topicName)
    afr := &AppendFieldRequest{
        Action:    "AppendField",
        FieldName: field.Name,
        FieldType: field.Type,
    }
    _, err := datahub.Client.Post(path, afr)
    if err != nil {
        return err
    }
    return nil

}

func (datahub *DataHub) GetMeterInfo(projectName, topicName, shardId string) (*GetMeterInfoResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    gmir := &GetMeterInfoRequest{
        Action: "meter",
    }
    respBody, err := datahub.Client.Post(path, gmir)
    if err != nil {
        return nil, err
    }
    return NewGetMeterInfoResult(respBody)
}

func (datahub *DataHub) CreateConnector(projectName, topicName string, cType ConnectorType, columnFields []string, config interface{}) (*CreateConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !validateConnectorType(cType) {
        return nil, NewInvalidParameterErrorWithMessage(parameterTypeInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, cType.String())
    ccr := &CreateConnectorRequest{
        Action:        "create",
        Type:          cType,
        SinkStartTime: -1,
        ColumnFields:  columnFields,
        Config:        config,
    }
    respBody, err := datahub.Client.Post(path, ccr)
    if err != nil {
        return nil, err
    }
    return NewCreateConnectorResult(respBody)
}

func (datahub *DataHub) CreateConnectorWithStartTime(projectName, topicName string, cType ConnectorType,
    columnFields []string, sinkStartTime int64, config interface{}) (*CreateConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !validateConnectorType(cType) {
        return nil, NewInvalidParameterErrorWithMessage(parameterTypeInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, cType.String())
    ccr := &CreateConnectorRequest{
        Action:        "create",
        Type:          cType,
        SinkStartTime: sinkStartTime,
        ColumnFields:  columnFields,
        Config:        config,
    }
    respBody, err := datahub.Client.Post(path, ccr)
    if err != nil {
        return nil, err
    }
    return NewCreateConnectorResult(respBody)
}

func (datahub *DataHub) GetConnector(projectName, topicName, connectorId string) (*GetConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    return NewGetConnectorResult(respBody)
}

func (datahub *DataHub) UpdateConnector(projectName, topicName, connectorId string, config interface{}) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)

    ucr := &UpdateConnectorRequest{
        Action: "updateconfig",
        Config: config,
    }
    if _, err := datahub.Client.Post(path, ucr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) ListConnector(projectName, topicName string) (*ListConnectorResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorsPath, projectName, topicName)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    return NewListConnectorResult(respBody)
}

func (datahub *DataHub) DeleteConnector(projectName, topicName, connectorId string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    if _, err := datahub.Client.Delete(path); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) GetConnectorDoneTime(projectName, topicName, connectorId string) (*GetConnectorDoneTimeResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorDoneTimePath, projectName, topicName, connectorId)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    return NewGetConnectorDoneTimeResult(respBody)
}

func (datahub *DataHub) ReloadConnector(projectName, topicName, connectorId string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    rcr := &ReloadConnectorRequest{
        Action: "Reload",
    }
    if _, err := datahub.Client.Post(path, rcr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) ReloadConnectorByShard(projectName, topicName, connectorId, shardId string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    rcr := &ReloadConnectorRequest{
        Action:  "Reload",
        ShardId: shardId,
    }
    if _, err := datahub.Client.Post(path, rcr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) UpdateConnectorState(projectName, topicName, connectorId string, state ConnectorState) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !validateConnectorState(state) {
        return NewInvalidParameterErrorWithMessage(parameterTypeInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    ucsr := &UpdateConnectorStateRequest{
        Action: "updatestate",
        State:  state,
    }
    if _, err := datahub.Client.Post(path, ucsr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) UpdateConnectorOffset(projectName, topicName, connectorId, shardId string, offset ConnectorOffset) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    ucor := &UpdateConnectorOffsetRequest{
        Action:    "updateshardcontext",
        ShardId:   shardId,
        Timestamp: offset.Timestamp,
        Sequence:  offset.Sequence,
    }

    if _, err := datahub.Client.Post(path, ucor); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) GetConnectorShardStatus(projectName, topicName, connectorId string) (*GetConnectorShardStatusResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    gcss := &GetConnectorShardStatusRequest{
        Action: "Status",
    }
    respBody, err := datahub.Client.Post(path, gcss)
    if err != nil {
        return nil, err
    }
    return NewGetConnectorShardStatusResult(respBody)
}

func (datahub *DataHub) GetConnectorShardStatusByShard(projectName, topicName, connectorId, shardId string) (*ConnectorShardStatusEntry, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    gcss := &GetConnectorShardStatusRequest{
        Action:  "Status",
        ShardId: shardId,
    }
    respBody, err := datahub.Client.Post(path, gcss)
    if err != nil {
        return nil, err
    }
    csse := &ConnectorShardStatusEntry{}
    if err := json.Unmarshal(respBody, csse); err != nil {
        return nil, err
    }
    return csse, nil
}

func (datahub *DataHub) AppendConnectorField(projectName, topicName, connectorId, fieldName string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(connectorPath, projectName, topicName, connectorId)
    acfr := &AppendConnectorFieldRequest{
        Action:    "appendfield",
        FieldName: fieldName,
    }
    if _, err := datahub.Client.Post(path, acfr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) CreateSubscription(projectName, topicName, comment string) (*CreateSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckComment(comment) {
        return nil, NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(subscriptionsPath, projectName, topicName)
    csr := &CreateSubscriptionRequest{
        Action:  "create",
        Comment: comment,
    }
    respBody, err := datahub.Client.Post(path, csr)
    if err != nil {
        return nil, err
    }
    return NewCreateSubscriptionResult(respBody)
}

func (datahub *DataHub) GetSubscription(projectName, topicName, subId string) (*GetSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    respBody, err := datahub.Client.Get(path)
    if err != nil {
        return nil, err
    }
    return NewGetSubscriptionResult(respBody)
}

func (datahub *DataHub) DeleteSubscription(projectName, topicName, subId string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    if _, err := datahub.Client.Delete(path); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) ListSubscription(projectName, topicName string, pageIndex, pageSize int) (*ListSubscriptionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionsPath, projectName, topicName)
    lsr := &ListSubscriptionRequest{
        Action:    "list",
        PageIndex: pageIndex,
        PageSize:  pageSize,
    }
    respBody, err := datahub.Client.Post(path, lsr)
    if err != nil {
        return nil, err
    }
    return NewListSubscriptionResult(respBody)
}

func (datahub *DataHub) UpdateSubscription(projectName, topicName, subId, comment string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckComment(comment) {
        return NewInvalidParameterErrorWithMessage(commentInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    usr := &UpdateSubscriptionRequest{
        Comment: comment,
    }
    if _, err := datahub.Client.Put(path, usr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) UpdateSubscriptionState(projectName, topicName, subId string, state SubscriptionState) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(subscriptionPath, projectName, topicName, subId)
    usr := &UpdateSubscriptionStateRequest{
        State: state,
    }
    if _, err := datahub.Client.Put(path, usr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) OpenSubscriptionSession(projectName, topicName, subId string, shardIds []string) (*OpenSubscriptionSessionResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    for _, id := range shardIds {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    ossr := &OpenSubscriptionSessionRequest{
        Action:   "open",
        ShardIds: shardIds,
    }
    respBody, err := datahub.Client.Post(path, ossr)
    if err != nil {
        return nil, err
    }
    return NewOpenSubscriptionSessionResult(respBody)
}

func (datahub *DataHub) GetSubscriptionOffset(projectName, topicName, subId string, shardIds []string) (*GetSubscriptionOffsetResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    for _, id := range shardIds {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    gsor := &GetSubscriptionOffsetRequest{
        Action:   "get",
        ShardIds: shardIds,
    }
    respBody, err := datahub.Client.Post(path, gsor)
    if err != nil {
        return nil, err
    }
    return NewGetSubscriptionOffsetResult(respBody)
}

func (datahub *DataHub) CommitSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    req := &CommitSubscriptionOffsetRequest{
        Action:  "commit",
        Offsets: offsets,
    }
    if _, err := datahub.Client.Put(path, req); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) ResetSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(offsetsPath, projectName, topicName, subId)
    req := &ResetSubscriptionOffsetRequest{
        Action:  "reset",
        Offsets: offsets,
    }
    if _, err := datahub.Client.Put(path, req); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHub) Heartbeat(projectName, topicName, consumerGroup, consumerId string, versionId int64, holdShardList, readEndShardList []string) (*HeartbeatResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    for _, id := range holdShardList {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }
    for _, id := range readEndShardList {
        if !util.CheckShardId(id) {
            return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    hr := &HeartbeatRequest{
        Action:           "heartbeat",
        ConsumerId:       consumerId,
        VersionId:        versionId,
        HoldShardList:    holdShardList,
        ReadEndShardList: readEndShardList,
    }
    respBody, err := datahub.Client.Post(path, hr)
    if err != nil {
        return nil, err
    }
    return NewHeartbeatResult(respBody)
}

func (datahub *DataHub) JoinGroup(projectName, topicName, consumerGroup string, sessionTimeout int64) (*JoinGroupResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    jgr := &JoinGroupRequest{
        Action:         "joinGroup",
        SessionTimeout: sessionTimeout,
    }
    respBody, err := datahub.Client.Post(path, jgr)
    if err != nil {
        return nil, err
    }
    return NewJoinGroupResult(respBody)

}
func (datahub *DataHub) SyncGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64, releaseShardList, readEndShardList []string) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if len(releaseShardList) == 0 || len(readEndShardList) == 0 {
        return NewInvalidParameterErrorWithMessage(shardListInvalid)
    }
    for _, id := range releaseShardList {
        if !util.CheckShardId(id) {
            return NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }
    for _, id := range readEndShardList {
        if !util.CheckShardId(id) {
            return NewInvalidParameterErrorWithMessage(shardIdInvalid)
        }
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    sgr := &SyncGroupRequest{
        Action:           "syncGroup",
        ConsumerId:       consumerId,
        VersionId:        versionId,
        ReleaseShardList: releaseShardList,
        ReadEndShardList: readEndShardList,
    }
    if _, err := datahub.Client.Post(path, sgr); err != nil {
        return err
    }
    return nil
}
func (datahub *DataHub) LeaveGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(consumerGroupPath, projectName, topicName, consumerGroup)
    lgr := &LeaveGroupRequest{
        Action:     "leaveGroup",
        ConsumerId: consumerId,
        VersionId:  versionId,
    }
    if _, err := datahub.Client.Post(path, lgr); err != nil {
        return err
    }
    return nil

}

type DataHubPB struct {
    DataHub
}

func (datahub *DataHubPB) PutRecords(projectName, topicName string, records []IRecord) (*PutRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }

    path := fmt.Sprintf(shardsPath, projectName, topicName)
    prr := &PutPBRecordsRequest{
        Records: records,
    }
    respBody, err := datahub.Client.Post(path, prr)
    if err != nil {
        return nil, err
    }
    return NewPutPBRecordsResult(respBody)
}

func (datahub *DataHubPB) PutRecordsByShard(projectName, topicName, shardId string, records []IRecord) error {
    if !util.CheckProjectName(projectName) {
        return NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    prr := &PutPBRecordsRequest{
        Records: records,
    }
    if _, err := datahub.Client.Post(path, prr); err != nil {
        return err
    }
    return nil
}

func (datahub *DataHubPB) GetTupleRecords(projectName, topicName, shardId, cursor string, limit int, recordSchema *RecordSchema) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    grr := &GetPBRecordRequest{
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, err := datahub.Client.Post(path, grr)
    if err != nil {
        return nil, err
    }
    return NewGetPBRecordsResult(respBody, recordSchema)
}

func (datahub *DataHubPB) GetBlobRecords(projectName, topicName, shardId, cursor string, limit int) (*GetRecordsResult, error) {
    if !util.CheckProjectName(projectName) {
        return nil, NewInvalidParameterErrorWithMessage(projectNameInvalid)
    }
    if !util.CheckTopicName(topicName) {
        return nil, NewInvalidParameterErrorWithMessage(topicNameInvalid)
    }
    if !util.CheckShardId(shardId) {
        return nil, NewInvalidParameterErrorWithMessage(shardIdInvalid)
    }

    path := fmt.Sprintf(shardPath, projectName, topicName, shardId)
    grr := &GetPBRecordRequest{
        Cursor: cursor,
        Limit:  limit,
    }
    respBody, err := datahub.Client.Post(path, grr)
    if err != nil {
        return nil, err
    }
    return NewGetPBRecordsResult(respBody, nil)
}
