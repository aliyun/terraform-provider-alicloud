package fc

// TriggerConfig defines the event bridge trigger config which maps to FunctionComputeConfiguration.
type EventBridgeTriggerConfig struct {
	TriggerEnable          *bool              `json:"triggerEnable"`
	AsyncInvocationType    *bool              `json:"asyncInvocationType"`
	EventSourceConfig      *EventSourceConfig `json:"eventSourceConfig"`
	EventRuleFilterPattern *string            `json:"eventRuleFilterPattern"`
}


// NewEventBridgeTriggerConfig ...
func NewEventBridgeTriggerConfig() *EventBridgeTriggerConfig {
	return &EventBridgeTriggerConfig{}
}

func (ebtc *EventBridgeTriggerConfig) WithTriggerEnable(triggerEnable bool) *EventBridgeTriggerConfig {
	ebtc.TriggerEnable = &triggerEnable
	return ebtc
}

func (ebtc *EventBridgeTriggerConfig) WithAsyncInvocationType(asyncInvocationType bool) *EventBridgeTriggerConfig {
	ebtc.AsyncInvocationType = &asyncInvocationType
	return ebtc
}

func (ebtc *EventBridgeTriggerConfig) WithEventSourceConfig(eventSourceConfig *EventSourceConfig) *EventBridgeTriggerConfig {
	ebtc.EventSourceConfig = eventSourceConfig
	return ebtc
}

func (ebtc *EventBridgeTriggerConfig) WithEventRuleFilterPattern(eventRuleFilterPattern string) *EventBridgeTriggerConfig {
	ebtc.EventRuleFilterPattern = &eventRuleFilterPattern
	return ebtc
}

// EventSourceConfig ...
type EventSourceConfig struct {
	EventSourceType       *string                `json:"eventSourceType"`
	EventSourceParameters *EventSourceParameters `json:"eventSourceParameters"`
}

// NewEventSourceConfig ...
func NewEventSourceConfig() *EventSourceConfig {
	return &EventSourceConfig{}
}

// WithEventSourceType ...
func (esc *EventSourceConfig) WithEventSourceType(eventSourceType string) *EventSourceConfig {
	esc.EventSourceType = &eventSourceType
	return esc
}

// WithEventSourceParameters ...
func (esc *EventSourceConfig) WithEventSourceParameters(eventSourceParameters *EventSourceParameters) *EventSourceConfig {
	esc.EventSourceParameters = eventSourceParameters
	return esc
}

// EventSourceParameters ...
type EventSourceParameters struct {
	SourceMNSParameters      *SourceMNSParameters      `json:"sourceMNSParameters"`
	SourceRocketMQParameters *SourceRocketMQParameters `json:"sourceRocketMQParameters"`

	SourceRabbitMQParameters *SourceRabbitMQParameters `json:"sourceRabbitMQParameters"`
}

// NewEventSourceParameters ...
func NewEventSourceParameters() *EventSourceParameters {
	return &EventSourceParameters{}
}

// WithSourceMNSParameters
func (esp *EventSourceParameters) WithSourceMNSParameters(sourceMNSParameters *SourceMNSParameters) *EventSourceParameters {
	esp.SourceMNSParameters = sourceMNSParameters
	return esp
}

// WithSourceRocketMQParameters
func (esp *EventSourceParameters) WithSourceRocketMQParameters(sourceRocketMQParameters *SourceRocketMQParameters) *EventSourceParameters {
	esp.SourceRocketMQParameters = sourceRocketMQParameters
	return esp
}

// WithSourceRabbitMQParameters
func (esp *EventSourceParameters) WithSourceRabbitMQParameters(sourceRabbitMQParameters *SourceRabbitMQParameters) *EventSourceParameters {
	esp.SourceRabbitMQParameters = sourceRabbitMQParameters
	return esp
}

// SourceMNSParameters refers to github.com/alibabacloud-go/eventbridge-sdk v1.2.10
type SourceMNSParameters struct {
	RegionId       *string `json:"RegionId,omitempty"`
	QueueName      *string `json:"QueueName,omitempty"`
	IsBase64Decode *bool   `json:"IsBase64Decode,omitempty"`
}

// NewSourceMNSParameters ...
func NewSourceMNSParameters() *SourceMNSParameters {
	return &SourceMNSParameters{}
}

// WithRegionId ...
func (mns *SourceMNSParameters) WithRegionId(regionId string) *SourceMNSParameters {
	mns.RegionId = &regionId
	return mns
}

// WithQueueName ...
func (mns *SourceMNSParameters) WithQueueName(queueName string) *SourceMNSParameters {
	mns.QueueName = &queueName
	return mns
}

// WithIsBase64Decode ...
func (mns *SourceMNSParameters) WithIsBase64Decode(isBase64Decode bool) *SourceMNSParameters {
	mns.IsBase64Decode = &isBase64Decode
	return mns
}

// SourceRocketMQParameters refers to github.com/alibabacloud-go/eventbridge-sdk v1.2.10
type SourceRocketMQParameters struct {
	RegionId   *string `json:"RegionId,omitempty"`
	InstanceId *string `json:"InstanceId,omitempty"`
	Topic      *string `json:"Topic,omitempty"`
	Tag        *string `json:"Tag,omitempty"`
	Offset     *string `json:"Offset,omitempty"`
	GroupID    *string `json:"GroupID,omitempty"`
	Timestamp  *int    `json:"Timestamp,omitempty"`
}

// NewSourceRocketMQParameters ...
func NewSourceRocketMQParameters() *SourceRocketMQParameters {
	return &SourceRocketMQParameters{}
}

// WithRegionId ...
func (rocket *SourceRocketMQParameters) WithRegionId(regionId string) *SourceRocketMQParameters {
	rocket.RegionId = &regionId
	return rocket
}

// WithInstanceId ...
func (rocket *SourceRocketMQParameters) WithInstanceId(instanceId string) *SourceRocketMQParameters {
	rocket.InstanceId = &instanceId
	return rocket
}

// WithTopic ...
func (rocket *SourceRocketMQParameters) WithTopic(topic string) *SourceRocketMQParameters {
	rocket.Topic = &topic
	return rocket
}

// WithTag ...
func (rocket *SourceRocketMQParameters) WithTag(tag string) *SourceRocketMQParameters {
	rocket.Tag = &tag
	return rocket
}

// WithOffset ...
func (rocket *SourceRocketMQParameters) WithOffset(offset string) *SourceRocketMQParameters {
	rocket.Offset = &offset
	return rocket
}

// WithGroupID ...
func (rocket *SourceRocketMQParameters) WithGroupID(groupID string) *SourceRocketMQParameters {
	rocket.GroupID = &groupID
	return rocket
}

// WithTimestamp ...
func (rocket *SourceRocketMQParameters) WithTimestamp(timestamp int) *SourceRocketMQParameters {
	rocket.Timestamp = &timestamp
	return rocket
}

// SourceRocketMQParameters refers to ithub.com/alibabacloud-go/eventbridge-sdk v1.2.10
type SourceRabbitMQParameters struct {
	RegionId        *string `json:"RegionId,omitempty"`
	InstanceId      *string `json:"InstanceId,omitempty"`
	VirtualHostName *string `json:"VirtualHostName,omitempty"`
	QueueName       *string `json:"QueueName,omitempty"`
}

// SourceRabbitMQParameters ...
func NewSourceRabbitMQParameters() *SourceRabbitMQParameters {
	return &SourceRabbitMQParameters{}
}

// WithRegionId ...
func (rabbit *SourceRabbitMQParameters) WithRegionId(regionId string) *SourceRabbitMQParameters {
	rabbit.RegionId = &regionId
	return rabbit
}

// WithInstanceId ...
func (rabbit *SourceRabbitMQParameters) WithInstanceId(instanceId string) *SourceRabbitMQParameters {
	rabbit.InstanceId = &instanceId
	return rabbit
}

// WithVirtualHostName ...
func (rabbit *SourceRabbitMQParameters) WithVirtualHostName(virtualHostName string) *SourceRabbitMQParameters {
	rabbit.VirtualHostName = &virtualHostName
	return rabbit
}

// WithQueueName ...
func (rabbit *SourceRabbitMQParameters) WithQueueName(queueName string) *SourceRabbitMQParameters {
	rabbit.QueueName = &queueName
	return rabbit
}
