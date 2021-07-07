package datahub

func New(accessId, accessKey, endpoint string) DataHubApi {
    config := NewDefaultConfig()
    return &DataHubPB{
        DataHub: DataHub{
            Client: NewRestClient(endpoint, config.UserAgent, config.HttpClient,
                NewAliyunAccount(accessId, accessKey), config.CompressorType),
        },
    }
}

func NewClientWithConfig(endpoint string, config *Config, account Account) DataHubApi {
    if config.UserAgent == "" {
        config.UserAgent = DefaultUserAgent()
    }
    if config.HttpClient == nil {
        config.HttpClient = DefaultHttpClient()
    }
    if !validateCompressorType(config.CompressorType) {
        config.CompressorType = NOCOMPRESS
    }

    dh := &DataHub{
        Client: NewRestClient(endpoint, config.UserAgent, config.HttpClient,
            account, config.CompressorType),
    }
    if !config.EnableBinary {
        return dh
    }

    //return &DataHubJson{}
    return &DataHubPB{
        DataHub: *dh,
    }
}

// Datahub provides restful apis for visiting examples service.
type DataHubApi interface {
    // Get the information of the specified project.
    GetProject(projectName string) (*GetProjectResult, error)

    // List all projects the user owns.
    ListProject() (*ListProjectResult, error)

    // Create a examples project.
    CreateProject(projectName, comment string) error

    // Update project information. Only support comment
    UpdateProject(projectName, comment string) error

    // Delete the specified project. If any topics exist in the project, the delete operation will fail.
    DeleteProject(projectName string) error

    // Wait for all shards' status of this topic is ACTIVE. Default timeout is 60s.
    WaitAllShardsReady(projectName, topicName string) bool

    // Wait for all shards' status of this topic is ACTIVE.
    // The unit is seconds.
    // If timeout < 0, it will block util all shards ready
    WaitAllShardsReadyWithTime(projectName, topicName string, timeout int64) bool

    // Create a examples topic with type: BLOB
    CreateBlobTopic(projectName, topicName, comment string, shardCount, lifeCycle int) error

    // Create a examples topic with type: TUPLE
    CreateTupleTopic(projectName, topicName, comment string, shardCount, lifeCycle int, recordSchema *RecordSchema) error

    // Update topic meta information. Now only support modify comment info.
    UpdateTopic(projectName, topicName, comment string) error

    // Get the information of the specified topic.
    GetTopic(projectName, topicName string) (*GetTopicResult, error)

    // Delete a specified topic.
    DeleteTopic(projectName, topicName string) error

    // List all topics in the project.
    ListTopic(projectName string) (*ListTopicResult, error)

    // List shard information {ShardEntry} of a topic.
    ListShard(projectName, topicName string) (*ListShardResult, error)

    // Split a shard. In function, sdk will automatically compute the split key which is used to split shard.
    SplitShard(projectName, topicName, shardId string) (*SplitShardResult, error)

    // Split a shard by the specified splitKey.
    SplitShardBySplitKey(projectName, topicName, shardId, splitKey string) (*SplitShardResult, error)

    // Merge the specified shard and its adjacent shard. Only adjacent shards can be merged.
    MergeShard(projectName, topicName, shardId, adjacentShardId string) (*MergeShardResult, error)

    // Get the data cursor of a shard. This function support OLDEST, LATEST, SYSTEM_TIME and SEQUENCE.
    // If choose OLDEST or LATEST, the last parameter will not be needed.
    // if choose SYSTEM_TIME or SEQUENCE. it needs to a parameter as sequence num or timestamp.
    GetCursor(projectName, topicName, shardId string, ctype CursorType, param ...int64) (*GetCursorResult, error)

    // Write data records into a DataHub topic.
    // The PutRecordsResult includes unsuccessfully processed records.
    // Datahub attempts to process all records in each record.
    // A single record failure does not stop the processing of subsequent records.
    PutRecords(projectName, topicName string, records []IRecord) (*PutRecordsResult, error)

    PutRecordsByShard(projectName, topicName, shardId string, records []IRecord) error

    // Get the TUPLE records of a shard.
    GetTupleRecords(projectName, topicName, shardId, cursor string, limit int, recordSchema *RecordSchema) (*GetRecordsResult, error)

    // Get the BLOB records of a shard.
    GetBlobRecords(projectName, topicName, shardId, cursor string, limit int) (*GetRecordsResult, error)

    // Append a field to a TUPLE topic.
    // Field AllowNull should be true.
    AppendField(projectName, topicName string, field Field) error

    // Get metering info of the specified shard
    GetMeterInfo(projectName, topicName, shardId string) (*GetMeterInfoResult, error)

    // Create data connectors.
    CreateConnector(projectName, topicName string, cType ConnectorType, columnFields []string, config interface{}) (*CreateConnectorResult, error)

    // Create connector with start time(unit:ms)
    CreateConnectorWithStartTime(projectName, topicName string, cType ConnectorType,
        columnFields []string, sinkStartTime int64, config interface{}) (*CreateConnectorResult, error)

    // Get information of the specified data connector.
    GetConnector(projectName, topicName, connectorId string) (*GetConnectorResult, error)

    // Update connector config of the specified data connector.
    // Config should be SinkOdpsConfig, SinkOssConfig ...
    UpdateConnector(projectName, topicName, connectorId string, config interface{}) error

    // List name of connectors.
    ListConnector(projectName, topicName string) (*ListConnectorResult, error)

    // Delete a data connector.
    DeleteConnector(projectName, topicName, connectorId string) error

    // Get the done time of a data connector. This method mainly used to get MaxCompute synchronize point.
    GetConnectorDoneTime(projectName, topicName, connectorId string) (*GetConnectorDoneTimeResult, error)

    // Reload a data connector.
    ReloadConnector(projectName, topicName, connectorId string) error

    // Reload the specified shard of the data connector.
    ReloadConnectorByShard(projectName, topicName, connectorId, shardId string) error

    // Update the state of the data connector
    UpdateConnectorState(projectName, topicName, connectorId string, state ConnectorState) error

    // Update connector sink offset. The operation must be operated after connector stopped.
    UpdateConnectorOffset(projectName, topicName, connectorId, shardId string, offset ConnectorOffset) error

    // Get the detail information of the shard task which belongs to the specified data connector.
    GetConnectorShardStatus(projectName, topicName, connectorId string) (*GetConnectorShardStatusResult, error)

    // Get the detail information of the shard task which belongs to the specified data connector.
    GetConnectorShardStatusByShard(projectName, topicName, connectorId, shardId string) (*ConnectorShardStatusEntry, error)

    // Append data connector field.
    // Before run this method, you should ensure that this field is in both the topic and the connector.
    AppendConnectorField(projectName, topicName, connectorId, fieldName string) error

    // Create a subscription, and then you should commit offsets with this subscription.
    CreateSubscription(projectName, topicName, comment string) (*CreateSubscriptionResult, error)

    // Get the detail information of a subscription.
    GetSubscription(projectName, topicName, subId string) (*GetSubscriptionResult, error)

    // Delete a subscription.
    DeleteSubscription(projectName, topicName, subId string) error

    // List subscriptions in the topic.
    ListSubscription(projectName, topicName string, pageIndex, pageSize int) (*ListSubscriptionResult, error)

    // Update a subscription. Now only support update comment information.
    UpdateSubscription(projectName, topicName, subId, comment string) error

    // Update a subscription' state. You can change the state of a subscription to SUB_ONLINE or SUB_OFFLINE.
    // When offline, you can not commit offsets of the subscription.
    UpdateSubscriptionState(projectName, topicName, subId string, state SubscriptionState) error

    // Init and get a subscription session, and returns offset if any offset stored before.
    // Subscription should be initialized before use. This operation makes sure that only one client use this subscription.
    // If this function be called in elsewhere, the seesion will be invalid and can not commit offsets of the subscription.
    OpenSubscriptionSession(projectName, topicName, subId string, shardIds []string) (*OpenSubscriptionSessionResult, error)

    // Get offsets of a subscription.This method dost not return sessionId in SubscriptionOffset.
    // Only the SubscriptionOffset containing sessionId can commit offset.
    GetSubscriptionOffset(projectName, topicName, subId string, shardIds []string) (*GetSubscriptionOffsetResult, error)

    // Update offsets of shards to server. This operation allows you store offsets on the server side.
    CommitSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) error

    // Reset offsets of shards to server. This operation allows you reset offsets on the server side.
    ResetSubscriptionOffset(projectName, topicName, subId string, offsets map[string]SubscriptionOffset) error

    // Heartbeat request to let server know consumer status.
    Heartbeat(projectName, topicName, consumerGroup, consumerId string, versionId int64, holdShardList, readEndShardList []string) (*HeartbeatResult, error)

    // Join a consumer group.
    JoinGroup(projectName, topicName, consumerGroup string, sessionTimeout int64) (*JoinGroupResult, error)

    // Sync consumer group info.
    SyncGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64, releaseShardList, readEndShardList []string) error

    // Leave consumer group info.
    LeaveGroup(projectName, topicName, consumerGroup, consumerId string, versionId int64) error
}
