// This file is auto-generated, don't edit it. Thanks.
/**
 *
 */
package client

import (
	fcutil "github.com/alibabacloud-go/alibabacloud-gateway-fc-util/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	endpointutil "github.com/alibabacloud-go/endpoint-util/service"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	
	"net/http"
)

type AccelerationInfo struct {
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s AccelerationInfo) String() string {
	return tea.Prettify(s)
}

func (s AccelerationInfo) GoString() string {
	return s.String()
}

func (s *AccelerationInfo) SetStatus(v string) *AccelerationInfo {
	s.Status = &v
	return s
}

type AsyncConfigMeta struct {
	FunctionName *string `json:"functionName,omitempty" xml:"functionName,omitempty"`
	Qualifier    *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	ServiceName  *string `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
}

func (s AsyncConfigMeta) String() string {
	return tea.Prettify(s)
}

func (s AsyncConfigMeta) GoString() string {
	return s.String()
}

func (s *AsyncConfigMeta) SetFunctionName(v string) *AsyncConfigMeta {
	s.FunctionName = &v
	return s
}

func (s *AsyncConfigMeta) SetQualifier(v string) *AsyncConfigMeta {
	s.Qualifier = &v
	return s
}

func (s *AsyncConfigMeta) SetServiceName(v string) *AsyncConfigMeta {
	s.ServiceName = &v
	return s
}

type AvailableAZ struct {
	AvailableAZs *string `json:"availableAZs,omitempty" xml:"availableAZs,omitempty"`
}

func (s AvailableAZ) String() string {
	return tea.Prettify(s)
}

func (s AvailableAZ) GoString() string {
	return s.String()
}

func (s *AvailableAZ) SetAvailableAZs(v string) *AvailableAZ {
	s.AvailableAZs = &v
	return s
}

type BatchWindow struct {
	CountBasedWindow *int64 `json:"CountBasedWindow,omitempty" xml:"CountBasedWindow,omitempty"`
	TimeBasedWindow  *int64 `json:"TimeBasedWindow,omitempty" xml:"TimeBasedWindow,omitempty"`
}

func (s BatchWindow) String() string {
	return tea.Prettify(s)
}

func (s BatchWindow) GoString() string {
	return s.String()
}

func (s *BatchWindow) SetCountBasedWindow(v int64) *BatchWindow {
	s.CountBasedWindow = &v
	return s
}

func (s *BatchWindow) SetTimeBasedWindow(v int64) *BatchWindow {
	s.TimeBasedWindow = &v
	return s
}

type CDNEventsTriggerConfig struct {
	EventName    *string              `json:"eventName,omitempty" xml:"eventName,omitempty"`
	EventVersion *string              `json:"eventVersion,omitempty" xml:"eventVersion,omitempty"`
	Filter       map[string][]*string `json:"filter,omitempty" xml:"filter,omitempty"`
	Notes        *string              `json:"notes,omitempty" xml:"notes,omitempty"`
}

func (s CDNEventsTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s CDNEventsTriggerConfig) GoString() string {
	return s.String()
}

func (s *CDNEventsTriggerConfig) SetEventName(v string) *CDNEventsTriggerConfig {
	s.EventName = &v
	return s
}

func (s *CDNEventsTriggerConfig) SetEventVersion(v string) *CDNEventsTriggerConfig {
	s.EventVersion = &v
	return s
}

func (s *CDNEventsTriggerConfig) SetFilter(v map[string][]*string) *CDNEventsTriggerConfig {
	s.Filter = v
	return s
}

func (s *CDNEventsTriggerConfig) SetNotes(v string) *CDNEventsTriggerConfig {
	s.Notes = &v
	return s
}

type CertConfig struct {
	CertName    *string `json:"certName,omitempty" xml:"certName,omitempty"`
	Certificate *string `json:"certificate,omitempty" xml:"certificate,omitempty"`
	PrivateKey  *string `json:"privateKey,omitempty" xml:"privateKey,omitempty"`
}

func (s CertConfig) String() string {
	return tea.Prettify(s)
}

func (s CertConfig) GoString() string {
	return s.String()
}

func (s *CertConfig) SetCertName(v string) *CertConfig {
	s.CertName = &v
	return s
}

func (s *CertConfig) SetCertificate(v string) *CertConfig {
	s.Certificate = &v
	return s
}

func (s *CertConfig) SetPrivateKey(v string) *CertConfig {
	s.PrivateKey = &v
	return s
}

type Code struct {
	OssBucketName *string `json:"ossBucketName,omitempty" xml:"ossBucketName,omitempty"`
	OssObjectName *string `json:"ossObjectName,omitempty" xml:"ossObjectName,omitempty"`
	ZipFile       *string `json:"zipFile,omitempty" xml:"zipFile,omitempty"`
}

func (s Code) String() string {
	return tea.Prettify(s)
}

func (s Code) GoString() string {
	return s.String()
}

func (s *Code) SetOssBucketName(v string) *Code {
	s.OssBucketName = &v
	return s
}

func (s *Code) SetOssObjectName(v string) *Code {
	s.OssObjectName = &v
	return s
}

func (s *Code) SetZipFile(v string) *Code {
	s.ZipFile = &v
	return s
}

type CustomContainerConfig struct {
	AccelerationType *string `json:"accelerationType,omitempty" xml:"accelerationType,omitempty"`
	Args             *string `json:"args,omitempty" xml:"args,omitempty"`
	Command          *string `json:"command,omitempty" xml:"command,omitempty"`
	Image            *string `json:"image,omitempty" xml:"image,omitempty"`
	InstanceID       *string `json:"instanceID,omitempty" xml:"instanceID,omitempty"`
	WebServerMode    *bool   `json:"webServerMode,omitempty" xml:"webServerMode,omitempty"`
}

func (s CustomContainerConfig) String() string {
	return tea.Prettify(s)
}

func (s CustomContainerConfig) GoString() string {
	return s.String()
}

func (s *CustomContainerConfig) SetAccelerationType(v string) *CustomContainerConfig {
	s.AccelerationType = &v
	return s
}

func (s *CustomContainerConfig) SetArgs(v string) *CustomContainerConfig {
	s.Args = &v
	return s
}

func (s *CustomContainerConfig) SetCommand(v string) *CustomContainerConfig {
	s.Command = &v
	return s
}

func (s *CustomContainerConfig) SetImage(v string) *CustomContainerConfig {
	s.Image = &v
	return s
}

func (s *CustomContainerConfig) SetInstanceID(v string) *CustomContainerConfig {
	s.InstanceID = &v
	return s
}

func (s *CustomContainerConfig) SetWebServerMode(v bool) *CustomContainerConfig {
	s.WebServerMode = &v
	return s
}

type CustomContainerConfigInfo struct {
	AccelerationInfo *AccelerationInfo `json:"accelerationInfo,omitempty" xml:"accelerationInfo,omitempty"`
	AccelerationType *string           `json:"accelerationType,omitempty" xml:"accelerationType,omitempty"`
	Args             *string           `json:"args,omitempty" xml:"args,omitempty"`
	Command          *string           `json:"command,omitempty" xml:"command,omitempty"`
	Image            *string           `json:"image,omitempty" xml:"image,omitempty"`
	InstanceID       *string           `json:"instanceID,omitempty" xml:"instanceID,omitempty"`
	WebServerMode    *bool             `json:"webServerMode,omitempty" xml:"webServerMode,omitempty"`
}

func (s CustomContainerConfigInfo) String() string {
	return tea.Prettify(s)
}

func (s CustomContainerConfigInfo) GoString() string {
	return s.String()
}

func (s *CustomContainerConfigInfo) SetAccelerationInfo(v *AccelerationInfo) *CustomContainerConfigInfo {
	s.AccelerationInfo = v
	return s
}

func (s *CustomContainerConfigInfo) SetAccelerationType(v string) *CustomContainerConfigInfo {
	s.AccelerationType = &v
	return s
}

func (s *CustomContainerConfigInfo) SetArgs(v string) *CustomContainerConfigInfo {
	s.Args = &v
	return s
}

func (s *CustomContainerConfigInfo) SetCommand(v string) *CustomContainerConfigInfo {
	s.Command = &v
	return s
}

func (s *CustomContainerConfigInfo) SetImage(v string) *CustomContainerConfigInfo {
	s.Image = &v
	return s
}

func (s *CustomContainerConfigInfo) SetInstanceID(v string) *CustomContainerConfigInfo {
	s.InstanceID = &v
	return s
}

func (s *CustomContainerConfigInfo) SetWebServerMode(v bool) *CustomContainerConfigInfo {
	s.WebServerMode = &v
	return s
}

type CustomDNS struct {
	DnsOptions  []*DNSOption `json:"dnsOptions,omitempty" xml:"dnsOptions,omitempty" type:"Repeated"`
	NameServers []*string    `json:"nameServers,omitempty" xml:"nameServers,omitempty" type:"Repeated"`
	Searches    []*string    `json:"searches,omitempty" xml:"searches,omitempty" type:"Repeated"`
}

func (s CustomDNS) String() string {
	return tea.Prettify(s)
}

func (s CustomDNS) GoString() string {
	return s.String()
}

func (s *CustomDNS) SetDnsOptions(v []*DNSOption) *CustomDNS {
	s.DnsOptions = v
	return s
}

func (s *CustomDNS) SetNameServers(v []*string) *CustomDNS {
	s.NameServers = v
	return s
}

func (s *CustomDNS) SetSearches(v []*string) *CustomDNS {
	s.Searches = v
	return s
}

type CustomHealthCheckConfig struct {
	FailureThreshold    *int32  `json:"failureThreshold,omitempty" xml:"failureThreshold,omitempty"`
	HttpGetUrl          *string `json:"httpGetUrl,omitempty" xml:"httpGetUrl,omitempty"`
	InitialDelaySeconds *int32  `json:"initialDelaySeconds,omitempty" xml:"initialDelaySeconds,omitempty"`
	PeriodSeconds       *int32  `json:"periodSeconds,omitempty" xml:"periodSeconds,omitempty"`
	SuccessThreshold    *int32  `json:"successThreshold,omitempty" xml:"successThreshold,omitempty"`
	TimeoutSeconds      *int32  `json:"timeoutSeconds,omitempty" xml:"timeoutSeconds,omitempty"`
}

func (s CustomHealthCheckConfig) String() string {
	return tea.Prettify(s)
}

func (s CustomHealthCheckConfig) GoString() string {
	return s.String()
}

func (s *CustomHealthCheckConfig) SetFailureThreshold(v int32) *CustomHealthCheckConfig {
	s.FailureThreshold = &v
	return s
}

func (s *CustomHealthCheckConfig) SetHttpGetUrl(v string) *CustomHealthCheckConfig {
	s.HttpGetUrl = &v
	return s
}

func (s *CustomHealthCheckConfig) SetInitialDelaySeconds(v int32) *CustomHealthCheckConfig {
	s.InitialDelaySeconds = &v
	return s
}

func (s *CustomHealthCheckConfig) SetPeriodSeconds(v int32) *CustomHealthCheckConfig {
	s.PeriodSeconds = &v
	return s
}

func (s *CustomHealthCheckConfig) SetSuccessThreshold(v int32) *CustomHealthCheckConfig {
	s.SuccessThreshold = &v
	return s
}

func (s *CustomHealthCheckConfig) SetTimeoutSeconds(v int32) *CustomHealthCheckConfig {
	s.TimeoutSeconds = &v
	return s
}

type CustomRuntimeConfig struct {
	Args    []*string `json:"args,omitempty" xml:"args,omitempty" type:"Repeated"`
	Command []*string `json:"command,omitempty" xml:"command,omitempty" type:"Repeated"`
}

func (s CustomRuntimeConfig) String() string {
	return tea.Prettify(s)
}

func (s CustomRuntimeConfig) GoString() string {
	return s.String()
}

func (s *CustomRuntimeConfig) SetArgs(v []*string) *CustomRuntimeConfig {
	s.Args = v
	return s
}

func (s *CustomRuntimeConfig) SetCommand(v []*string) *CustomRuntimeConfig {
	s.Command = v
	return s
}

type DNSOption struct {
	Name  *string `json:"name,omitempty" xml:"name,omitempty"`
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s DNSOption) String() string {
	return tea.Prettify(s)
}

func (s DNSOption) GoString() string {
	return s.String()
}

func (s *DNSOption) SetName(v string) *DNSOption {
	s.Name = &v
	return s
}

func (s *DNSOption) SetValue(v string) *DNSOption {
	s.Value = &v
	return s
}

type DeadLetterQueue struct {
	Arn *string `json:"Arn,omitempty" xml:"Arn,omitempty"`
}

func (s DeadLetterQueue) String() string {
	return tea.Prettify(s)
}

func (s DeadLetterQueue) GoString() string {
	return s.String()
}

func (s *DeadLetterQueue) SetArn(v string) *DeadLetterQueue {
	s.Arn = &v
	return s
}

type DeliveryOption struct {
	EventSchema *string `json:"eventSchema,omitempty" xml:"eventSchema,omitempty"`
	Mode        *string `json:"mode,omitempty" xml:"mode,omitempty"`
}

func (s DeliveryOption) String() string {
	return tea.Prettify(s)
}

func (s DeliveryOption) GoString() string {
	return s.String()
}

func (s *DeliveryOption) SetEventSchema(v string) *DeliveryOption {
	s.EventSchema = &v
	return s
}

func (s *DeliveryOption) SetMode(v string) *DeliveryOption {
	s.Mode = &v
	return s
}

type Destination struct {
	Destination *string `json:"destination,omitempty" xml:"destination,omitempty"`
}

func (s Destination) String() string {
	return tea.Prettify(s)
}

func (s Destination) GoString() string {
	return s.String()
}

func (s *Destination) SetDestination(v string) *Destination {
	s.Destination = &v
	return s
}

type DestinationConfig struct {
	OnFailure *Destination `json:"onFailure,omitempty" xml:"onFailure,omitempty"`
	OnSuccess *Destination `json:"onSuccess,omitempty" xml:"onSuccess,omitempty"`
}

func (s DestinationConfig) String() string {
	return tea.Prettify(s)
}

func (s DestinationConfig) GoString() string {
	return s.String()
}

func (s *DestinationConfig) SetOnFailure(v *Destination) *DestinationConfig {
	s.OnFailure = v
	return s
}

func (s *DestinationConfig) SetOnSuccess(v *Destination) *DestinationConfig {
	s.OnSuccess = v
	return s
}

type Error struct {
	ErrorCode    *string `json:"errorCode,omitempty" xml:"errorCode,omitempty"`
	ErrorMessage *string `json:"errorMessage,omitempty" xml:"errorMessage,omitempty"`
}

func (s Error) String() string {
	return tea.Prettify(s)
}

func (s Error) GoString() string {
	return s.String()
}

func (s *Error) SetErrorCode(v string) *Error {
	s.ErrorCode = &v
	return s
}

func (s *Error) SetErrorMessage(v string) *Error {
	s.ErrorMessage = &v
	return s
}

type ErrorInfo struct {
	ErrorMessage *string `json:"errorMessage,omitempty" xml:"errorMessage,omitempty"`
	StackTrace   *string `json:"stackTrace,omitempty" xml:"stackTrace,omitempty"`
}

func (s ErrorInfo) String() string {
	return tea.Prettify(s)
}

func (s ErrorInfo) GoString() string {
	return s.String()
}

func (s *ErrorInfo) SetErrorMessage(v string) *ErrorInfo {
	s.ErrorMessage = &v
	return s
}

func (s *ErrorInfo) SetStackTrace(v string) *ErrorInfo {
	s.StackTrace = &v
	return s
}

type EventBridgeTriggerConfig struct {
	AsyncInvocationType    *bool              `json:"asyncInvocationType,omitempty" xml:"asyncInvocationType,omitempty"`
	EventRuleFilterPattern *string            `json:"eventRuleFilterPattern,omitempty" xml:"eventRuleFilterPattern,omitempty"`
	EventSinkConfig        *EventSinkConfig   `json:"eventSinkConfig,omitempty" xml:"eventSinkConfig,omitempty"`
	EventSourceConfig      *EventSourceConfig `json:"eventSourceConfig,omitempty" xml:"eventSourceConfig,omitempty"`
	RunOptions             *RunOptions        `json:"runOptions,omitempty" xml:"runOptions,omitempty"`
	TriggerEnable          *bool              `json:"triggerEnable,omitempty" xml:"triggerEnable,omitempty"`
}

func (s EventBridgeTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s EventBridgeTriggerConfig) GoString() string {
	return s.String()
}

func (s *EventBridgeTriggerConfig) SetAsyncInvocationType(v bool) *EventBridgeTriggerConfig {
	s.AsyncInvocationType = &v
	return s
}

func (s *EventBridgeTriggerConfig) SetEventRuleFilterPattern(v string) *EventBridgeTriggerConfig {
	s.EventRuleFilterPattern = &v
	return s
}

func (s *EventBridgeTriggerConfig) SetEventSinkConfig(v *EventSinkConfig) *EventBridgeTriggerConfig {
	s.EventSinkConfig = v
	return s
}

func (s *EventBridgeTriggerConfig) SetEventSourceConfig(v *EventSourceConfig) *EventBridgeTriggerConfig {
	s.EventSourceConfig = v
	return s
}

func (s *EventBridgeTriggerConfig) SetRunOptions(v *RunOptions) *EventBridgeTriggerConfig {
	s.RunOptions = v
	return s
}

func (s *EventBridgeTriggerConfig) SetTriggerEnable(v bool) *EventBridgeTriggerConfig {
	s.TriggerEnable = &v
	return s
}

type EventSinkConfig struct {
	DeliveryOption *DeliveryOption `json:"deliveryOption,omitempty" xml:"deliveryOption,omitempty"`
}

func (s EventSinkConfig) String() string {
	return tea.Prettify(s)
}

func (s EventSinkConfig) GoString() string {
	return s.String()
}

func (s *EventSinkConfig) SetDeliveryOption(v *DeliveryOption) *EventSinkConfig {
	s.DeliveryOption = v
	return s
}

type EventSourceConfig struct {
	EventSourceParameters *EventSourceParameters `json:"eventSourceParameters,omitempty" xml:"eventSourceParameters,omitempty"`
	EventSourceType       *string                `json:"eventSourceType,omitempty" xml:"eventSourceType,omitempty"`
}

func (s EventSourceConfig) String() string {
	return tea.Prettify(s)
}

func (s EventSourceConfig) GoString() string {
	return s.String()
}

func (s *EventSourceConfig) SetEventSourceParameters(v *EventSourceParameters) *EventSourceConfig {
	s.EventSourceParameters = v
	return s
}

func (s *EventSourceConfig) SetEventSourceType(v string) *EventSourceConfig {
	s.EventSourceType = &v
	return s
}

type EventSourceParameters struct {
	SourceKafkaParameters    *SourceKafkaParameters    `json:"sourceKafkaParameters,omitempty" xml:"sourceKafkaParameters,omitempty"`
	SourceMNSParameters      *SourceMNSParameters      `json:"sourceMNSParameters,omitempty" xml:"sourceMNSParameters,omitempty"`
	SourceRabbitMQParameters *SourceRabbitMQParameters `json:"sourceRabbitMQParameters,omitempty" xml:"sourceRabbitMQParameters,omitempty"`
	SourceRocketMQParameters *SourceRocketMQParameters `json:"sourceRocketMQParameters,omitempty" xml:"sourceRocketMQParameters,omitempty"`
}

func (s EventSourceParameters) String() string {
	return tea.Prettify(s)
}

func (s EventSourceParameters) GoString() string {
	return s.String()
}

func (s *EventSourceParameters) SetSourceKafkaParameters(v *SourceKafkaParameters) *EventSourceParameters {
	s.SourceKafkaParameters = v
	return s
}

func (s *EventSourceParameters) SetSourceMNSParameters(v *SourceMNSParameters) *EventSourceParameters {
	s.SourceMNSParameters = v
	return s
}

func (s *EventSourceParameters) SetSourceRabbitMQParameters(v *SourceRabbitMQParameters) *EventSourceParameters {
	s.SourceRabbitMQParameters = v
	return s
}

func (s *EventSourceParameters) SetSourceRocketMQParameters(v *SourceRocketMQParameters) *EventSourceParameters {
	s.SourceRocketMQParameters = v
	return s
}

type HTTPTriggerConfig struct {
	AuthConfig *string `json:"authConfig,omitempty" xml:"authConfig,omitempty"`
	AuthType   *string `json:"authType,omitempty" xml:"authType,omitempty"`
	// 禁用默认公网域名访问的开关，设置为true 时，访问函数默认提供的公网URL地址会返回403错误。设置为 false 则不会有任何影响。
	DisableURLInternet *bool     `json:"disableURLInternet,omitempty" xml:"disableURLInternet,omitempty"`
	Methods            []*string `json:"methods,omitempty" xml:"methods,omitempty" type:"Repeated"`
}

func (s HTTPTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s HTTPTriggerConfig) GoString() string {
	return s.String()
}

func (s *HTTPTriggerConfig) SetAuthConfig(v string) *HTTPTriggerConfig {
	s.AuthConfig = &v
	return s
}

func (s *HTTPTriggerConfig) SetAuthType(v string) *HTTPTriggerConfig {
	s.AuthType = &v
	return s
}

func (s *HTTPTriggerConfig) SetDisableURLInternet(v bool) *HTTPTriggerConfig {
	s.DisableURLInternet = &v
	return s
}

func (s *HTTPTriggerConfig) SetMethods(v []*string) *HTTPTriggerConfig {
	s.Methods = v
	return s
}

type InstanceLifecycleConfig struct {
	PreFreeze *LifecycleHook `json:"preFreeze,omitempty" xml:"preFreeze,omitempty"`
	PreStop   *LifecycleHook `json:"preStop,omitempty" xml:"preStop,omitempty"`
}

func (s InstanceLifecycleConfig) String() string {
	return tea.Prettify(s)
}

func (s InstanceLifecycleConfig) GoString() string {
	return s.String()
}

func (s *InstanceLifecycleConfig) SetPreFreeze(v *LifecycleHook) *InstanceLifecycleConfig {
	s.PreFreeze = v
	return s
}

func (s *InstanceLifecycleConfig) SetPreStop(v *LifecycleHook) *InstanceLifecycleConfig {
	s.PreStop = v
	return s
}

type JWTAuthConfig struct {
	BlackList   *string   `json:"blackList,omitempty" xml:"blackList,omitempty"`
	ClaimPassBy []*string `json:"claimPassBy,omitempty" xml:"claimPassBy,omitempty" type:"Repeated"`
	Jwks        *string   `json:"jwks,omitempty" xml:"jwks,omitempty"`
	TokenLookup []*string `json:"tokenLookup,omitempty" xml:"tokenLookup,omitempty" type:"Repeated"`
	WhiteList   []*string `json:"whiteList,omitempty" xml:"whiteList,omitempty" type:"Repeated"`
}

func (s JWTAuthConfig) String() string {
	return tea.Prettify(s)
}

func (s JWTAuthConfig) GoString() string {
	return s.String()
}

func (s *JWTAuthConfig) SetBlackList(v string) *JWTAuthConfig {
	s.BlackList = &v
	return s
}

func (s *JWTAuthConfig) SetClaimPassBy(v []*string) *JWTAuthConfig {
	s.ClaimPassBy = v
	return s
}

func (s *JWTAuthConfig) SetJwks(v string) *JWTAuthConfig {
	s.Jwks = &v
	return s
}

func (s *JWTAuthConfig) SetTokenLookup(v []*string) *JWTAuthConfig {
	s.TokenLookup = v
	return s
}

func (s *JWTAuthConfig) SetWhiteList(v []*string) *JWTAuthConfig {
	s.WhiteList = v
	return s
}

type JaegerConfig struct {
	Endpoint *string `json:"endpoint,omitempty" xml:"endpoint,omitempty"`
}

func (s JaegerConfig) String() string {
	return tea.Prettify(s)
}

func (s JaegerConfig) GoString() string {
	return s.String()
}

func (s *JaegerConfig) SetEndpoint(v string) *JaegerConfig {
	s.Endpoint = &v
	return s
}

type JobConfig struct {
	MaxRetryTime    *int64 `json:"maxRetryTime,omitempty" xml:"maxRetryTime,omitempty"`
	TriggerInterval *int64 `json:"triggerInterval,omitempty" xml:"triggerInterval,omitempty"`
}

func (s JobConfig) String() string {
	return tea.Prettify(s)
}

func (s JobConfig) GoString() string {
	return s.String()
}

func (s *JobConfig) SetMaxRetryTime(v int64) *JobConfig {
	s.MaxRetryTime = &v
	return s
}

func (s *JobConfig) SetTriggerInterval(v int64) *JobConfig {
	s.TriggerInterval = &v
	return s
}

type JobLogConfig struct {
	Logstore *string `json:"logstore,omitempty" xml:"logstore,omitempty"`
	Project  *string `json:"project,omitempty" xml:"project,omitempty"`
}

func (s JobLogConfig) String() string {
	return tea.Prettify(s)
}

func (s JobLogConfig) GoString() string {
	return s.String()
}

func (s *JobLogConfig) SetLogstore(v string) *JobLogConfig {
	s.Logstore = &v
	return s
}

func (s *JobLogConfig) SetProject(v string) *JobLogConfig {
	s.Project = &v
	return s
}

type Layer struct {
	Acl               *int32     `json:"acl,omitempty" xml:"acl,omitempty"`
	Arn               *string    `json:"arn,omitempty" xml:"arn,omitempty"`
	ArnV2             *string    `json:"arnV2,omitempty" xml:"arnV2,omitempty"`
	Code              *LayerCode `json:"code,omitempty" xml:"code,omitempty"`
	CodeChecksum      *string    `json:"codeChecksum,omitempty" xml:"codeChecksum,omitempty"`
	CodeSize          *int64     `json:"codeSize,omitempty" xml:"codeSize,omitempty"`
	CompatibleRuntime []*string  `json:"compatibleRuntime,omitempty" xml:"compatibleRuntime,omitempty" type:"Repeated"`
	CreateTime        *string    `json:"createTime,omitempty" xml:"createTime,omitempty"`
	Description       *string    `json:"description,omitempty" xml:"description,omitempty"`
	LayerName         *string    `json:"layerName,omitempty" xml:"layerName,omitempty"`
	License           *string    `json:"license,omitempty" xml:"license,omitempty"`
	Version           *int32     `json:"version,omitempty" xml:"version,omitempty"`
}

func (s Layer) String() string {
	return tea.Prettify(s)
}

func (s Layer) GoString() string {
	return s.String()
}

func (s *Layer) SetAcl(v int32) *Layer {
	s.Acl = &v
	return s
}

func (s *Layer) SetArn(v string) *Layer {
	s.Arn = &v
	return s
}

func (s *Layer) SetArnV2(v string) *Layer {
	s.ArnV2 = &v
	return s
}

func (s *Layer) SetCode(v *LayerCode) *Layer {
	s.Code = v
	return s
}

func (s *Layer) SetCodeChecksum(v string) *Layer {
	s.CodeChecksum = &v
	return s
}

func (s *Layer) SetCodeSize(v int64) *Layer {
	s.CodeSize = &v
	return s
}

func (s *Layer) SetCompatibleRuntime(v []*string) *Layer {
	s.CompatibleRuntime = v
	return s
}

func (s *Layer) SetCreateTime(v string) *Layer {
	s.CreateTime = &v
	return s
}

func (s *Layer) SetDescription(v string) *Layer {
	s.Description = &v
	return s
}

func (s *Layer) SetLayerName(v string) *Layer {
	s.LayerName = &v
	return s
}

func (s *Layer) SetLicense(v string) *Layer {
	s.License = &v
	return s
}

func (s *Layer) SetVersion(v int32) *Layer {
	s.Version = &v
	return s
}

type LayerCode struct {
	Location       *string `json:"location,omitempty" xml:"location,omitempty"`
	RepositoryType *string `json:"repositoryType,omitempty" xml:"repositoryType,omitempty"`
}

func (s LayerCode) String() string {
	return tea.Prettify(s)
}

func (s LayerCode) GoString() string {
	return s.String()
}

func (s *LayerCode) SetLocation(v string) *LayerCode {
	s.Location = &v
	return s
}

func (s *LayerCode) SetRepositoryType(v string) *LayerCode {
	s.RepositoryType = &v
	return s
}

type LifecycleHook struct {
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	Timeout *int32  `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s LifecycleHook) String() string {
	return tea.Prettify(s)
}

func (s LifecycleHook) GoString() string {
	return s.String()
}

func (s *LifecycleHook) SetHandler(v string) *LifecycleHook {
	s.Handler = &v
	return s
}

func (s *LifecycleHook) SetTimeout(v int32) *LifecycleHook {
	s.Timeout = &v
	return s
}

type LogConfig struct {
	EnableInstanceMetrics *bool   `json:"enableInstanceMetrics,omitempty" xml:"enableInstanceMetrics,omitempty"`
	EnableRequestMetrics  *bool   `json:"enableRequestMetrics,omitempty" xml:"enableRequestMetrics,omitempty"`
	LogBeginRule          *string `json:"logBeginRule,omitempty" xml:"logBeginRule,omitempty"`
	Logstore              *string `json:"logstore,omitempty" xml:"logstore,omitempty"`
	Project               *string `json:"project,omitempty" xml:"project,omitempty"`
}

func (s LogConfig) String() string {
	return tea.Prettify(s)
}

func (s LogConfig) GoString() string {
	return s.String()
}

func (s *LogConfig) SetEnableInstanceMetrics(v bool) *LogConfig {
	s.EnableInstanceMetrics = &v
	return s
}

func (s *LogConfig) SetEnableRequestMetrics(v bool) *LogConfig {
	s.EnableRequestMetrics = &v
	return s
}

func (s *LogConfig) SetLogBeginRule(v string) *LogConfig {
	s.LogBeginRule = &v
	return s
}

func (s *LogConfig) SetLogstore(v string) *LogConfig {
	s.Logstore = &v
	return s
}

func (s *LogConfig) SetProject(v string) *LogConfig {
	s.Project = &v
	return s
}

type LogTriggerConfig struct {
	Enable            *bool              `json:"enable,omitempty" xml:"enable,omitempty"`
	FunctionParameter map[string]*string `json:"functionParameter,omitempty" xml:"functionParameter,omitempty"`
	JobConfig         *JobConfig         `json:"jobConfig,omitempty" xml:"jobConfig,omitempty"`
	LogConfig         *JobLogConfig      `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	SourceConfig      *SourceConfig      `json:"sourceConfig,omitempty" xml:"sourceConfig,omitempty"`
}

func (s LogTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s LogTriggerConfig) GoString() string {
	return s.String()
}

func (s *LogTriggerConfig) SetEnable(v bool) *LogTriggerConfig {
	s.Enable = &v
	return s
}

func (s *LogTriggerConfig) SetFunctionParameter(v map[string]*string) *LogTriggerConfig {
	s.FunctionParameter = v
	return s
}

func (s *LogTriggerConfig) SetJobConfig(v *JobConfig) *LogTriggerConfig {
	s.JobConfig = v
	return s
}

func (s *LogTriggerConfig) SetLogConfig(v *JobLogConfig) *LogTriggerConfig {
	s.LogConfig = v
	return s
}

func (s *LogTriggerConfig) SetSourceConfig(v *SourceConfig) *LogTriggerConfig {
	s.SourceConfig = v
	return s
}

type MeteringConfig struct {
	LogConfig *LogConfig `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	PayerId   *string    `json:"payerId,omitempty" xml:"payerId,omitempty"`
	Role      *string    `json:"role,omitempty" xml:"role,omitempty"`
}

func (s MeteringConfig) String() string {
	return tea.Prettify(s)
}

func (s MeteringConfig) GoString() string {
	return s.String()
}

func (s *MeteringConfig) SetLogConfig(v *LogConfig) *MeteringConfig {
	s.LogConfig = v
	return s
}

func (s *MeteringConfig) SetPayerId(v string) *MeteringConfig {
	s.PayerId = &v
	return s
}

func (s *MeteringConfig) SetRole(v string) *MeteringConfig {
	s.Role = &v
	return s
}

type MnsTopicTriggerConfig struct {
	FilterTag           *string `json:"filterTag,omitempty" xml:"filterTag,omitempty"`
	NotifyContentFormat *string `json:"notifyContentFormat,omitempty" xml:"notifyContentFormat,omitempty"`
	NotifyStrategy      *string `json:"notifyStrategy,omitempty" xml:"notifyStrategy,omitempty"`
}

func (s MnsTopicTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s MnsTopicTriggerConfig) GoString() string {
	return s.String()
}

func (s *MnsTopicTriggerConfig) SetFilterTag(v string) *MnsTopicTriggerConfig {
	s.FilterTag = &v
	return s
}

func (s *MnsTopicTriggerConfig) SetNotifyContentFormat(v string) *MnsTopicTriggerConfig {
	s.NotifyContentFormat = &v
	return s
}

func (s *MnsTopicTriggerConfig) SetNotifyStrategy(v string) *MnsTopicTriggerConfig {
	s.NotifyStrategy = &v
	return s
}

type NASConfig struct {
	GroupId     *int32                  `json:"groupId,omitempty" xml:"groupId,omitempty"`
	MountPoints []*NASConfigMountPoints `json:"mountPoints,omitempty" xml:"mountPoints,omitempty" type:"Repeated"`
	UserId      *int32                  `json:"userId,omitempty" xml:"userId,omitempty"`
}

func (s NASConfig) String() string {
	return tea.Prettify(s)
}

func (s NASConfig) GoString() string {
	return s.String()
}

func (s *NASConfig) SetGroupId(v int32) *NASConfig {
	s.GroupId = &v
	return s
}

func (s *NASConfig) SetMountPoints(v []*NASConfigMountPoints) *NASConfig {
	s.MountPoints = v
	return s
}

func (s *NASConfig) SetUserId(v int32) *NASConfig {
	s.UserId = &v
	return s
}

type NASConfigMountPoints struct {
	EnableTLS  *bool   `json:"enableTLS,omitempty" xml:"enableTLS,omitempty"`
	MountDir   *string `json:"mountDir,omitempty" xml:"mountDir,omitempty"`
	ServerAddr *string `json:"serverAddr,omitempty" xml:"serverAddr,omitempty"`
}

func (s NASConfigMountPoints) String() string {
	return tea.Prettify(s)
}

func (s NASConfigMountPoints) GoString() string {
	return s.String()
}

func (s *NASConfigMountPoints) SetEnableTLS(v bool) *NASConfigMountPoints {
	s.EnableTLS = &v
	return s
}

func (s *NASConfigMountPoints) SetMountDir(v string) *NASConfigMountPoints {
	s.MountDir = &v
	return s
}

func (s *NASConfigMountPoints) SetServerAddr(v string) *NASConfigMountPoints {
	s.ServerAddr = &v
	return s
}

type OSSMountConfig struct {
	MountPoints []*OSSMountConfigMountPoints `json:"mountPoints,omitempty" xml:"mountPoints,omitempty" type:"Repeated"`
}

func (s OSSMountConfig) String() string {
	return tea.Prettify(s)
}

func (s OSSMountConfig) GoString() string {
	return s.String()
}

func (s *OSSMountConfig) SetMountPoints(v []*OSSMountConfigMountPoints) *OSSMountConfig {
	s.MountPoints = v
	return s
}

type OSSMountConfigMountPoints struct {
	BucketName *string `json:"bucketName,omitempty" xml:"bucketName,omitempty"`
	BucketPath *string `json:"bucketPath,omitempty" xml:"bucketPath,omitempty"`
	Endpoint   *string `json:"endpoint,omitempty" xml:"endpoint,omitempty"`
	MountDir   *string `json:"mountDir,omitempty" xml:"mountDir,omitempty"`
	ReadOnly   *bool   `json:"readOnly,omitempty" xml:"readOnly,omitempty"`
}

func (s OSSMountConfigMountPoints) String() string {
	return tea.Prettify(s)
}

func (s OSSMountConfigMountPoints) GoString() string {
	return s.String()
}

func (s *OSSMountConfigMountPoints) SetBucketName(v string) *OSSMountConfigMountPoints {
	s.BucketName = &v
	return s
}

func (s *OSSMountConfigMountPoints) SetBucketPath(v string) *OSSMountConfigMountPoints {
	s.BucketPath = &v
	return s
}

func (s *OSSMountConfigMountPoints) SetEndpoint(v string) *OSSMountConfigMountPoints {
	s.Endpoint = &v
	return s
}

func (s *OSSMountConfigMountPoints) SetMountDir(v string) *OSSMountConfigMountPoints {
	s.MountDir = &v
	return s
}

func (s *OSSMountConfigMountPoints) SetReadOnly(v bool) *OSSMountConfigMountPoints {
	s.ReadOnly = &v
	return s
}

type OSSTriggerConfig struct {
	Events []*string         `json:"events,omitempty" xml:"events,omitempty" type:"Repeated"`
	Filter *OSSTriggerFilter `json:"filter,omitempty" xml:"filter,omitempty"`
}

func (s OSSTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s OSSTriggerConfig) GoString() string {
	return s.String()
}

func (s *OSSTriggerConfig) SetEvents(v []*string) *OSSTriggerConfig {
	s.Events = v
	return s
}

func (s *OSSTriggerConfig) SetFilter(v *OSSTriggerFilter) *OSSTriggerConfig {
	s.Filter = v
	return s
}

type OSSTriggerFilter struct {
	Key *OSSTriggerKey `json:"key,omitempty" xml:"key,omitempty"`
}

func (s OSSTriggerFilter) String() string {
	return tea.Prettify(s)
}

func (s OSSTriggerFilter) GoString() string {
	return s.String()
}

func (s *OSSTriggerFilter) SetKey(v *OSSTriggerKey) *OSSTriggerFilter {
	s.Key = v
	return s
}

type OSSTriggerKey struct {
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	Suffix *string `json:"suffix,omitempty" xml:"suffix,omitempty"`
}

func (s OSSTriggerKey) String() string {
	return tea.Prettify(s)
}

func (s OSSTriggerKey) GoString() string {
	return s.String()
}

func (s *OSSTriggerKey) SetPrefix(v string) *OSSTriggerKey {
	s.Prefix = &v
	return s
}

func (s *OSSTriggerKey) SetSuffix(v string) *OSSTriggerKey {
	s.Suffix = &v
	return s
}

type OnDemandConfig struct {
	MaximumInstanceCount *int64  `json:"maximumInstanceCount,omitempty" xml:"maximumInstanceCount,omitempty"`
	Resource             *string `json:"resource,omitempty" xml:"resource,omitempty"`
}

func (s OnDemandConfig) String() string {
	return tea.Prettify(s)
}

func (s OnDemandConfig) GoString() string {
	return s.String()
}

func (s *OnDemandConfig) SetMaximumInstanceCount(v int64) *OnDemandConfig {
	s.MaximumInstanceCount = &v
	return s
}

func (s *OnDemandConfig) SetResource(v string) *OnDemandConfig {
	s.Resource = &v
	return s
}

type OpenReservedCapacity struct {
	CreatedTime      *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	Cu               *int64  `json:"cu,omitempty" xml:"cu,omitempty"`
	Deadline         *string `json:"deadline,omitempty" xml:"deadline,omitempty"`
	InstanceId       *string `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	IsRefunded       *string `json:"isRefunded,omitempty" xml:"isRefunded,omitempty"`
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
}

func (s OpenReservedCapacity) String() string {
	return tea.Prettify(s)
}

func (s OpenReservedCapacity) GoString() string {
	return s.String()
}

func (s *OpenReservedCapacity) SetCreatedTime(v string) *OpenReservedCapacity {
	s.CreatedTime = &v
	return s
}

func (s *OpenReservedCapacity) SetCu(v int64) *OpenReservedCapacity {
	s.Cu = &v
	return s
}

func (s *OpenReservedCapacity) SetDeadline(v string) *OpenReservedCapacity {
	s.Deadline = &v
	return s
}

func (s *OpenReservedCapacity) SetInstanceId(v string) *OpenReservedCapacity {
	s.InstanceId = &v
	return s
}

func (s *OpenReservedCapacity) SetIsRefunded(v string) *OpenReservedCapacity {
	s.IsRefunded = &v
	return s
}

func (s *OpenReservedCapacity) SetLastModifiedTime(v string) *OpenReservedCapacity {
	s.LastModifiedTime = &v
	return s
}

type OutputCodeLocation struct {
	Location       *string `json:"location,omitempty" xml:"location,omitempty"`
	RepositoryType *string `json:"repositoryType,omitempty" xml:"repositoryType,omitempty"`
}

func (s OutputCodeLocation) String() string {
	return tea.Prettify(s)
}

func (s OutputCodeLocation) GoString() string {
	return s.String()
}

func (s *OutputCodeLocation) SetLocation(v string) *OutputCodeLocation {
	s.Location = &v
	return s
}

func (s *OutputCodeLocation) SetRepositoryType(v string) *OutputCodeLocation {
	s.RepositoryType = &v
	return s
}

type PathConfig struct {
	FunctionName  *string        `json:"functionName,omitempty" xml:"functionName,omitempty"`
	Methods       []*string      `json:"methods,omitempty" xml:"methods,omitempty" type:"Repeated"`
	Path          *string        `json:"path,omitempty" xml:"path,omitempty"`
	Qualifier     *string        `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	RewriteConfig *RewriteConfig `json:"rewriteConfig,omitempty" xml:"rewriteConfig,omitempty"`
	ServiceName   *string        `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
}

func (s PathConfig) String() string {
	return tea.Prettify(s)
}

func (s PathConfig) GoString() string {
	return s.String()
}

func (s *PathConfig) SetFunctionName(v string) *PathConfig {
	s.FunctionName = &v
	return s
}

func (s *PathConfig) SetMethods(v []*string) *PathConfig {
	s.Methods = v
	return s
}

func (s *PathConfig) SetPath(v string) *PathConfig {
	s.Path = &v
	return s
}

func (s *PathConfig) SetQualifier(v string) *PathConfig {
	s.Qualifier = &v
	return s
}

func (s *PathConfig) SetRewriteConfig(v *RewriteConfig) *PathConfig {
	s.RewriteConfig = v
	return s
}

func (s *PathConfig) SetServiceName(v string) *PathConfig {
	s.ServiceName = &v
	return s
}

type PolicyItem struct {
	Key      *string `json:"key,omitempty" xml:"key,omitempty"`
	Operator *string `json:"operator,omitempty" xml:"operator,omitempty"`
	Type     *string `json:"type,omitempty" xml:"type,omitempty"`
	Value    *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s PolicyItem) String() string {
	return tea.Prettify(s)
}

func (s PolicyItem) GoString() string {
	return s.String()
}

func (s *PolicyItem) SetKey(v string) *PolicyItem {
	s.Key = &v
	return s
}

func (s *PolicyItem) SetOperator(v string) *PolicyItem {
	s.Operator = &v
	return s
}

func (s *PolicyItem) SetType(v string) *PolicyItem {
	s.Type = &v
	return s
}

func (s *PolicyItem) SetValue(v string) *PolicyItem {
	s.Value = &v
	return s
}

type PreFreeze struct {
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	Timeout *int32  `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s PreFreeze) String() string {
	return tea.Prettify(s)
}

func (s PreFreeze) GoString() string {
	return s.String()
}

func (s *PreFreeze) SetHandler(v string) *PreFreeze {
	s.Handler = &v
	return s
}

func (s *PreFreeze) SetTimeout(v int32) *PreFreeze {
	s.Timeout = &v
	return s
}

type PreStop struct {
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	Timeout *int32  `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s PreStop) String() string {
	return tea.Prettify(s)
}

func (s PreStop) GoString() string {
	return s.String()
}

func (s *PreStop) SetHandler(v string) *PreStop {
	s.Handler = &v
	return s
}

func (s *PreStop) SetTimeout(v int32) *PreStop {
	s.Timeout = &v
	return s
}

type RdsTriggerConfig struct {
	Concurrency         *int64    `json:"concurrency,omitempty" xml:"concurrency,omitempty"`
	EventFormat         *string   `json:"eventFormat,omitempty" xml:"eventFormat,omitempty"`
	Retry               *int64    `json:"retry,omitempty" xml:"retry,omitempty"`
	SubscriptionObjects []*string `json:"subscriptionObjects,omitempty" xml:"subscriptionObjects,omitempty" type:"Repeated"`
}

func (s RdsTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s RdsTriggerConfig) GoString() string {
	return s.String()
}

func (s *RdsTriggerConfig) SetConcurrency(v int64) *RdsTriggerConfig {
	s.Concurrency = &v
	return s
}

func (s *RdsTriggerConfig) SetEventFormat(v string) *RdsTriggerConfig {
	s.EventFormat = &v
	return s
}

func (s *RdsTriggerConfig) SetRetry(v int64) *RdsTriggerConfig {
	s.Retry = &v
	return s
}

func (s *RdsTriggerConfig) SetSubscriptionObjects(v []*string) *RdsTriggerConfig {
	s.SubscriptionObjects = v
	return s
}

type Resource struct {
	ResourceArn *string            `json:"resourceArn,omitempty" xml:"resourceArn,omitempty"`
	Tags        map[string]*string `json:"tags,omitempty" xml:"tags,omitempty"`
}

func (s Resource) String() string {
	return tea.Prettify(s)
}

func (s Resource) GoString() string {
	return s.String()
}

func (s *Resource) SetResourceArn(v string) *Resource {
	s.ResourceArn = &v
	return s
}

func (s *Resource) SetTags(v map[string]*string) *Resource {
	s.Tags = v
	return s
}

type RetryStrategy struct {
	MaximumEventAgeInSeconds *int64  `json:"MaximumEventAgeInSeconds,omitempty" xml:"MaximumEventAgeInSeconds,omitempty"`
	MaximumRetryAttempts     *int64  `json:"MaximumRetryAttempts,omitempty" xml:"MaximumRetryAttempts,omitempty"`
	PushRetryStrategy        *string `json:"PushRetryStrategy,omitempty" xml:"PushRetryStrategy,omitempty"`
}

func (s RetryStrategy) String() string {
	return tea.Prettify(s)
}

func (s RetryStrategy) GoString() string {
	return s.String()
}

func (s *RetryStrategy) SetMaximumEventAgeInSeconds(v int64) *RetryStrategy {
	s.MaximumEventAgeInSeconds = &v
	return s
}

func (s *RetryStrategy) SetMaximumRetryAttempts(v int64) *RetryStrategy {
	s.MaximumRetryAttempts = &v
	return s
}

func (s *RetryStrategy) SetPushRetryStrategy(v string) *RetryStrategy {
	s.PushRetryStrategy = &v
	return s
}

type RewriteConfig struct {
	EqualRules    []*RewriteConfigEqualRules    `json:"equalRules,omitempty" xml:"equalRules,omitempty" type:"Repeated"`
	RegexRules    []*RewriteConfigRegexRules    `json:"regexRules,omitempty" xml:"regexRules,omitempty" type:"Repeated"`
	WildcardRules []*RewriteConfigWildcardRules `json:"wildcardRules,omitempty" xml:"wildcardRules,omitempty" type:"Repeated"`
}

func (s RewriteConfig) String() string {
	return tea.Prettify(s)
}

func (s RewriteConfig) GoString() string {
	return s.String()
}

func (s *RewriteConfig) SetEqualRules(v []*RewriteConfigEqualRules) *RewriteConfig {
	s.EqualRules = v
	return s
}

func (s *RewriteConfig) SetRegexRules(v []*RewriteConfigRegexRules) *RewriteConfig {
	s.RegexRules = v
	return s
}

func (s *RewriteConfig) SetWildcardRules(v []*RewriteConfigWildcardRules) *RewriteConfig {
	s.WildcardRules = v
	return s
}

type RewriteConfigEqualRules struct {
	Match       *string `json:"match,omitempty" xml:"match,omitempty"`
	Replacement *string `json:"replacement,omitempty" xml:"replacement,omitempty"`
}

func (s RewriteConfigEqualRules) String() string {
	return tea.Prettify(s)
}

func (s RewriteConfigEqualRules) GoString() string {
	return s.String()
}

func (s *RewriteConfigEqualRules) SetMatch(v string) *RewriteConfigEqualRules {
	s.Match = &v
	return s
}

func (s *RewriteConfigEqualRules) SetReplacement(v string) *RewriteConfigEqualRules {
	s.Replacement = &v
	return s
}

type RewriteConfigRegexRules struct {
	Match       *string `json:"match,omitempty" xml:"match,omitempty"`
	Replacement *string `json:"replacement,omitempty" xml:"replacement,omitempty"`
}

func (s RewriteConfigRegexRules) String() string {
	return tea.Prettify(s)
}

func (s RewriteConfigRegexRules) GoString() string {
	return s.String()
}

func (s *RewriteConfigRegexRules) SetMatch(v string) *RewriteConfigRegexRules {
	s.Match = &v
	return s
}

func (s *RewriteConfigRegexRules) SetReplacement(v string) *RewriteConfigRegexRules {
	s.Replacement = &v
	return s
}

type RewriteConfigWildcardRules struct {
	Match       *string `json:"match,omitempty" xml:"match,omitempty"`
	Replacement *string `json:"replacement,omitempty" xml:"replacement,omitempty"`
}

func (s RewriteConfigWildcardRules) String() string {
	return tea.Prettify(s)
}

func (s RewriteConfigWildcardRules) GoString() string {
	return s.String()
}

func (s *RewriteConfigWildcardRules) SetMatch(v string) *RewriteConfigWildcardRules {
	s.Match = &v
	return s
}

func (s *RewriteConfigWildcardRules) SetReplacement(v string) *RewriteConfigWildcardRules {
	s.Replacement = &v
	return s
}

type RouteConfig struct {
	Routes []*PathConfig `json:"routes,omitempty" xml:"routes,omitempty" type:"Repeated"`
}

func (s RouteConfig) String() string {
	return tea.Prettify(s)
}

func (s RouteConfig) GoString() string {
	return s.String()
}

func (s *RouteConfig) SetRoutes(v []*PathConfig) *RouteConfig {
	s.Routes = v
	return s
}

type RoutePolicy struct {
	Condition   *string       `json:"condition,omitempty" xml:"condition,omitempty"`
	PolicyItems []*PolicyItem `json:"policyItems,omitempty" xml:"policyItems,omitempty" type:"Repeated"`
}

func (s RoutePolicy) String() string {
	return tea.Prettify(s)
}

func (s RoutePolicy) GoString() string {
	return s.String()
}

func (s *RoutePolicy) SetCondition(v string) *RoutePolicy {
	s.Condition = &v
	return s
}

func (s *RoutePolicy) SetPolicyItems(v []*PolicyItem) *RoutePolicy {
	s.PolicyItems = v
	return s
}

type RunOptions struct {
	BatchWindow     *BatchWindow     `json:"batchWindow,omitempty" xml:"batchWindow,omitempty"`
	DeadLetterQueue *DeadLetterQueue `json:"deadLetterQueue,omitempty" xml:"deadLetterQueue,omitempty"`
	ErrorsTolerance *string          `json:"errorsTolerance,omitempty" xml:"errorsTolerance,omitempty"`
	MaximumTasks    *int64           `json:"maximumTasks,omitempty" xml:"maximumTasks,omitempty"`
	Mode            *string          `json:"mode,omitempty" xml:"mode,omitempty"`
	RetryStrategy   *RetryStrategy   `json:"retryStrategy,omitempty" xml:"retryStrategy,omitempty"`
}

func (s RunOptions) String() string {
	return tea.Prettify(s)
}

func (s RunOptions) GoString() string {
	return s.String()
}

func (s *RunOptions) SetBatchWindow(v *BatchWindow) *RunOptions {
	s.BatchWindow = v
	return s
}

func (s *RunOptions) SetDeadLetterQueue(v *DeadLetterQueue) *RunOptions {
	s.DeadLetterQueue = v
	return s
}

func (s *RunOptions) SetErrorsTolerance(v string) *RunOptions {
	s.ErrorsTolerance = &v
	return s
}

func (s *RunOptions) SetMaximumTasks(v int64) *RunOptions {
	s.MaximumTasks = &v
	return s
}

func (s *RunOptions) SetMode(v string) *RunOptions {
	s.Mode = &v
	return s
}

func (s *RunOptions) SetRetryStrategy(v *RetryStrategy) *RunOptions {
	s.RetryStrategy = v
	return s
}

type ScheduledActions struct {
	EndTime            *string `json:"endTime,omitempty" xml:"endTime,omitempty"`
	Name               *string `json:"name,omitempty" xml:"name,omitempty"`
	ScheduleExpression *string `json:"scheduleExpression,omitempty" xml:"scheduleExpression,omitempty"`
	StartTime          *string `json:"startTime,omitempty" xml:"startTime,omitempty"`
	Target             *int64  `json:"target,omitempty" xml:"target,omitempty"`
}

func (s ScheduledActions) String() string {
	return tea.Prettify(s)
}

func (s ScheduledActions) GoString() string {
	return s.String()
}

func (s *ScheduledActions) SetEndTime(v string) *ScheduledActions {
	s.EndTime = &v
	return s
}

func (s *ScheduledActions) SetName(v string) *ScheduledActions {
	s.Name = &v
	return s
}

func (s *ScheduledActions) SetScheduleExpression(v string) *ScheduledActions {
	s.ScheduleExpression = &v
	return s
}

func (s *ScheduledActions) SetStartTime(v string) *ScheduledActions {
	s.StartTime = &v
	return s
}

func (s *ScheduledActions) SetTarget(v int64) *ScheduledActions {
	s.Target = &v
	return s
}

type SourceConfig struct {
	Logstore *string `json:"logstore,omitempty" xml:"logstore,omitempty"`
}

func (s SourceConfig) String() string {
	return tea.Prettify(s)
}

func (s SourceConfig) GoString() string {
	return s.String()
}

func (s *SourceConfig) SetLogstore(v string) *SourceConfig {
	s.Logstore = &v
	return s
}

type SourceKafkaParameters struct {
	ConsumerGroup   *string `json:"ConsumerGroup,omitempty" xml:"ConsumerGroup,omitempty"`
	InstanceId      *string `json:"InstanceId,omitempty" xml:"InstanceId,omitempty"`
	Network         *string `json:"Network,omitempty" xml:"Network,omitempty"`
	OffsetReset     *string `json:"OffsetReset,omitempty" xml:"OffsetReset,omitempty"`
	RegionId        *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	SecurityGroupId *string `json:"SecurityGroupId,omitempty" xml:"SecurityGroupId,omitempty"`
	Topic           *string `json:"Topic,omitempty" xml:"Topic,omitempty"`
	VSwitchIds      *string `json:"VSwitchIds,omitempty" xml:"VSwitchIds,omitempty"`
	VpcId           *string `json:"VpcId,omitempty" xml:"VpcId,omitempty"`
}

func (s SourceKafkaParameters) String() string {
	return tea.Prettify(s)
}

func (s SourceKafkaParameters) GoString() string {
	return s.String()
}

func (s *SourceKafkaParameters) SetConsumerGroup(v string) *SourceKafkaParameters {
	s.ConsumerGroup = &v
	return s
}

func (s *SourceKafkaParameters) SetInstanceId(v string) *SourceKafkaParameters {
	s.InstanceId = &v
	return s
}

func (s *SourceKafkaParameters) SetNetwork(v string) *SourceKafkaParameters {
	s.Network = &v
	return s
}

func (s *SourceKafkaParameters) SetOffsetReset(v string) *SourceKafkaParameters {
	s.OffsetReset = &v
	return s
}

func (s *SourceKafkaParameters) SetRegionId(v string) *SourceKafkaParameters {
	s.RegionId = &v
	return s
}

func (s *SourceKafkaParameters) SetSecurityGroupId(v string) *SourceKafkaParameters {
	s.SecurityGroupId = &v
	return s
}

func (s *SourceKafkaParameters) SetTopic(v string) *SourceKafkaParameters {
	s.Topic = &v
	return s
}

func (s *SourceKafkaParameters) SetVSwitchIds(v string) *SourceKafkaParameters {
	s.VSwitchIds = &v
	return s
}

func (s *SourceKafkaParameters) SetVpcId(v string) *SourceKafkaParameters {
	s.VpcId = &v
	return s
}

type SourceMNSParameters struct {
	IsBase64Decode *bool   `json:"IsBase64Decode,omitempty" xml:"IsBase64Decode,omitempty"`
	QueueName      *string `json:"QueueName,omitempty" xml:"QueueName,omitempty"`
	RegionId       *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
}

func (s SourceMNSParameters) String() string {
	return tea.Prettify(s)
}

func (s SourceMNSParameters) GoString() string {
	return s.String()
}

func (s *SourceMNSParameters) SetIsBase64Decode(v bool) *SourceMNSParameters {
	s.IsBase64Decode = &v
	return s
}

func (s *SourceMNSParameters) SetQueueName(v string) *SourceMNSParameters {
	s.QueueName = &v
	return s
}

func (s *SourceMNSParameters) SetRegionId(v string) *SourceMNSParameters {
	s.RegionId = &v
	return s
}

type SourceRabbitMQParameters struct {
	InstanceId      *string `json:"InstanceId,omitempty" xml:"InstanceId,omitempty"`
	QueueName       *string `json:"QueueName,omitempty" xml:"QueueName,omitempty"`
	RegionId        *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	VirtualHostName *string `json:"VirtualHostName,omitempty" xml:"VirtualHostName,omitempty"`
}

func (s SourceRabbitMQParameters) String() string {
	return tea.Prettify(s)
}

func (s SourceRabbitMQParameters) GoString() string {
	return s.String()
}

func (s *SourceRabbitMQParameters) SetInstanceId(v string) *SourceRabbitMQParameters {
	s.InstanceId = &v
	return s
}

func (s *SourceRabbitMQParameters) SetQueueName(v string) *SourceRabbitMQParameters {
	s.QueueName = &v
	return s
}

func (s *SourceRabbitMQParameters) SetRegionId(v string) *SourceRabbitMQParameters {
	s.RegionId = &v
	return s
}

func (s *SourceRabbitMQParameters) SetVirtualHostName(v string) *SourceRabbitMQParameters {
	s.VirtualHostName = &v
	return s
}

type SourceRocketMQParameters struct {
	GroupID    *string `json:"GroupID,omitempty" xml:"GroupID,omitempty"`
	InstanceId *string `json:"InstanceId,omitempty" xml:"InstanceId,omitempty"`
	Offset     *string `json:"Offset,omitempty" xml:"Offset,omitempty"`
	RegionId   *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	Tag        *string `json:"Tag,omitempty" xml:"Tag,omitempty"`
	Timestamp  *int64  `json:"Timestamp,omitempty" xml:"Timestamp,omitempty"`
	Topic      *string `json:"Topic,omitempty" xml:"Topic,omitempty"`
}

func (s SourceRocketMQParameters) String() string {
	return tea.Prettify(s)
}

func (s SourceRocketMQParameters) GoString() string {
	return s.String()
}

func (s *SourceRocketMQParameters) SetGroupID(v string) *SourceRocketMQParameters {
	s.GroupID = &v
	return s
}

func (s *SourceRocketMQParameters) SetInstanceId(v string) *SourceRocketMQParameters {
	s.InstanceId = &v
	return s
}

func (s *SourceRocketMQParameters) SetOffset(v string) *SourceRocketMQParameters {
	s.Offset = &v
	return s
}

func (s *SourceRocketMQParameters) SetRegionId(v string) *SourceRocketMQParameters {
	s.RegionId = &v
	return s
}

func (s *SourceRocketMQParameters) SetTag(v string) *SourceRocketMQParameters {
	s.Tag = &v
	return s
}

func (s *SourceRocketMQParameters) SetTimestamp(v int64) *SourceRocketMQParameters {
	s.Timestamp = &v
	return s
}

func (s *SourceRocketMQParameters) SetTopic(v string) *SourceRocketMQParameters {
	s.Topic = &v
	return s
}

type StatefulAsyncInvocation struct {
	AlreadyRetriedTimes    *int64                          `json:"alreadyRetriedTimes,omitempty" xml:"alreadyRetriedTimes,omitempty"`
	DestinationStatus      *string                         `json:"destinationStatus,omitempty" xml:"destinationStatus,omitempty"`
	EndTime                *int64                          `json:"endTime,omitempty" xml:"endTime,omitempty"`
	Events                 []*StatefulAsyncInvocationEvent `json:"events,omitempty" xml:"events,omitempty" type:"Repeated"`
	FunctionName           *string                         `json:"functionName,omitempty" xml:"functionName,omitempty"`
	InstanceId             *string                         `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	InvocationErrorMessage *string                         `json:"invocationErrorMessage,omitempty" xml:"invocationErrorMessage,omitempty"`
	InvocationId           *string                         `json:"invocationId,omitempty" xml:"invocationId,omitempty"`
	InvocationPayload      *string                         `json:"invocationPayload,omitempty" xml:"invocationPayload,omitempty"`
	Qualifier              *string                         `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	RequestId              *string                         `json:"requestId,omitempty" xml:"requestId,omitempty"`
	ServiceName            *string                         `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
	StartedTime            *int64                          `json:"startedTime,omitempty" xml:"startedTime,omitempty"`
	Status                 *string                         `json:"status,omitempty" xml:"status,omitempty"`
}

func (s StatefulAsyncInvocation) String() string {
	return tea.Prettify(s)
}

func (s StatefulAsyncInvocation) GoString() string {
	return s.String()
}

func (s *StatefulAsyncInvocation) SetAlreadyRetriedTimes(v int64) *StatefulAsyncInvocation {
	s.AlreadyRetriedTimes = &v
	return s
}

func (s *StatefulAsyncInvocation) SetDestinationStatus(v string) *StatefulAsyncInvocation {
	s.DestinationStatus = &v
	return s
}

func (s *StatefulAsyncInvocation) SetEndTime(v int64) *StatefulAsyncInvocation {
	s.EndTime = &v
	return s
}

func (s *StatefulAsyncInvocation) SetEvents(v []*StatefulAsyncInvocationEvent) *StatefulAsyncInvocation {
	s.Events = v
	return s
}

func (s *StatefulAsyncInvocation) SetFunctionName(v string) *StatefulAsyncInvocation {
	s.FunctionName = &v
	return s
}

func (s *StatefulAsyncInvocation) SetInstanceId(v string) *StatefulAsyncInvocation {
	s.InstanceId = &v
	return s
}

func (s *StatefulAsyncInvocation) SetInvocationErrorMessage(v string) *StatefulAsyncInvocation {
	s.InvocationErrorMessage = &v
	return s
}

func (s *StatefulAsyncInvocation) SetInvocationId(v string) *StatefulAsyncInvocation {
	s.InvocationId = &v
	return s
}

func (s *StatefulAsyncInvocation) SetInvocationPayload(v string) *StatefulAsyncInvocation {
	s.InvocationPayload = &v
	return s
}

func (s *StatefulAsyncInvocation) SetQualifier(v string) *StatefulAsyncInvocation {
	s.Qualifier = &v
	return s
}

func (s *StatefulAsyncInvocation) SetRequestId(v string) *StatefulAsyncInvocation {
	s.RequestId = &v
	return s
}

func (s *StatefulAsyncInvocation) SetServiceName(v string) *StatefulAsyncInvocation {
	s.ServiceName = &v
	return s
}

func (s *StatefulAsyncInvocation) SetStartedTime(v int64) *StatefulAsyncInvocation {
	s.StartedTime = &v
	return s
}

func (s *StatefulAsyncInvocation) SetStatus(v string) *StatefulAsyncInvocation {
	s.Status = &v
	return s
}

type StatefulAsyncInvocationEvent struct {
	EventDetail *string `json:"eventDetail,omitempty" xml:"eventDetail,omitempty"`
	EventId     *int64  `json:"eventId,omitempty" xml:"eventId,omitempty"`
	Status      *string `json:"status,omitempty" xml:"status,omitempty"`
	Timestamp   *int64  `json:"timestamp,omitempty" xml:"timestamp,omitempty"`
}

func (s StatefulAsyncInvocationEvent) String() string {
	return tea.Prettify(s)
}

func (s StatefulAsyncInvocationEvent) GoString() string {
	return s.String()
}

func (s *StatefulAsyncInvocationEvent) SetEventDetail(v string) *StatefulAsyncInvocationEvent {
	s.EventDetail = &v
	return s
}

func (s *StatefulAsyncInvocationEvent) SetEventId(v int64) *StatefulAsyncInvocationEvent {
	s.EventId = &v
	return s
}

func (s *StatefulAsyncInvocationEvent) SetStatus(v string) *StatefulAsyncInvocationEvent {
	s.Status = &v
	return s
}

func (s *StatefulAsyncInvocationEvent) SetTimestamp(v int64) *StatefulAsyncInvocationEvent {
	s.Timestamp = &v
	return s
}

type TLSConfig struct {
	CipherSuites []*string `json:"cipherSuites,omitempty" xml:"cipherSuites,omitempty" type:"Repeated"`
	MaxVersion   *string   `json:"maxVersion,omitempty" xml:"maxVersion,omitempty"`
	MinVersion   *string   `json:"minVersion,omitempty" xml:"minVersion,omitempty"`
}

func (s TLSConfig) String() string {
	return tea.Prettify(s)
}

func (s TLSConfig) GoString() string {
	return s.String()
}

func (s *TLSConfig) SetCipherSuites(v []*string) *TLSConfig {
	s.CipherSuites = v
	return s
}

func (s *TLSConfig) SetMaxVersion(v string) *TLSConfig {
	s.MaxVersion = &v
	return s
}

func (s *TLSConfig) SetMinVersion(v string) *TLSConfig {
	s.MinVersion = &v
	return s
}

type TargetTrackingPolicies struct {
	EndTime      *string  `json:"endTime,omitempty" xml:"endTime,omitempty"`
	MaxCapacity  *int64   `json:"maxCapacity,omitempty" xml:"maxCapacity,omitempty"`
	MetricTarget *float64 `json:"metricTarget,omitempty" xml:"metricTarget,omitempty"`
	MetricType   *string  `json:"metricType,omitempty" xml:"metricType,omitempty"`
	MinCapacity  *int64   `json:"minCapacity,omitempty" xml:"minCapacity,omitempty"`
	Name         *string  `json:"name,omitempty" xml:"name,omitempty"`
	StartTime    *string  `json:"startTime,omitempty" xml:"startTime,omitempty"`
}

func (s TargetTrackingPolicies) String() string {
	return tea.Prettify(s)
}

func (s TargetTrackingPolicies) GoString() string {
	return s.String()
}

func (s *TargetTrackingPolicies) SetEndTime(v string) *TargetTrackingPolicies {
	s.EndTime = &v
	return s
}

func (s *TargetTrackingPolicies) SetMaxCapacity(v int64) *TargetTrackingPolicies {
	s.MaxCapacity = &v
	return s
}

func (s *TargetTrackingPolicies) SetMetricTarget(v float64) *TargetTrackingPolicies {
	s.MetricTarget = &v
	return s
}

func (s *TargetTrackingPolicies) SetMetricType(v string) *TargetTrackingPolicies {
	s.MetricType = &v
	return s
}

func (s *TargetTrackingPolicies) SetMinCapacity(v int64) *TargetTrackingPolicies {
	s.MinCapacity = &v
	return s
}

func (s *TargetTrackingPolicies) SetName(v string) *TargetTrackingPolicies {
	s.Name = &v
	return s
}

func (s *TargetTrackingPolicies) SetStartTime(v string) *TargetTrackingPolicies {
	s.StartTime = &v
	return s
}

type TimeTriggerConfig struct {
	CronExpression *string `json:"cronExpression,omitempty" xml:"cronExpression,omitempty"`
	Enable         *bool   `json:"enable,omitempty" xml:"enable,omitempty"`
	Payload        *string `json:"payload,omitempty" xml:"payload,omitempty"`
}

func (s TimeTriggerConfig) String() string {
	return tea.Prettify(s)
}

func (s TimeTriggerConfig) GoString() string {
	return s.String()
}

func (s *TimeTriggerConfig) SetCronExpression(v string) *TimeTriggerConfig {
	s.CronExpression = &v
	return s
}

func (s *TimeTriggerConfig) SetEnable(v bool) *TimeTriggerConfig {
	s.Enable = &v
	return s
}

func (s *TimeTriggerConfig) SetPayload(v string) *TimeTriggerConfig {
	s.Payload = &v
	return s
}

type TracingConfig struct {
	Params map[string]*string `json:"params,omitempty" xml:"params,omitempty"`
	Type   *string            `json:"type,omitempty" xml:"type,omitempty"`
}

func (s TracingConfig) String() string {
	return tea.Prettify(s)
}

func (s TracingConfig) GoString() string {
	return s.String()
}

func (s *TracingConfig) SetParams(v map[string]*string) *TracingConfig {
	s.Params = v
	return s
}

func (s *TracingConfig) SetType(v string) *TracingConfig {
	s.Type = &v
	return s
}

type Trigger struct {
	CreatedTime      *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	DomainName       *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	InvocationRole   *string `json:"invocationRole,omitempty" xml:"invocationRole,omitempty"`
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	Qualifier        *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	SourceArn        *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
	Status           *string `json:"status,omitempty" xml:"status,omitempty"`
	TargetArn        *string `json:"targetArn,omitempty" xml:"targetArn,omitempty"`
	TriggerConfig    *string `json:"triggerConfig,omitempty" xml:"triggerConfig,omitempty"`
	TriggerId        *string `json:"triggerId,omitempty" xml:"triggerId,omitempty"`
	TriggerName      *string `json:"triggerName,omitempty" xml:"triggerName,omitempty"`
	TriggerType      *string `json:"triggerType,omitempty" xml:"triggerType,omitempty"`
	UrlInternet      *string `json:"urlInternet,omitempty" xml:"urlInternet,omitempty"`
	UrlIntranet      *string `json:"urlIntranet,omitempty" xml:"urlIntranet,omitempty"`
}

func (s Trigger) String() string {
	return tea.Prettify(s)
}

func (s Trigger) GoString() string {
	return s.String()
}

func (s *Trigger) SetCreatedTime(v string) *Trigger {
	s.CreatedTime = &v
	return s
}

func (s *Trigger) SetDomainName(v string) *Trigger {
	s.DomainName = &v
	return s
}

func (s *Trigger) SetInvocationRole(v string) *Trigger {
	s.InvocationRole = &v
	return s
}

func (s *Trigger) SetLastModifiedTime(v string) *Trigger {
	s.LastModifiedTime = &v
	return s
}

func (s *Trigger) SetQualifier(v string) *Trigger {
	s.Qualifier = &v
	return s
}

func (s *Trigger) SetSourceArn(v string) *Trigger {
	s.SourceArn = &v
	return s
}

func (s *Trigger) SetStatus(v string) *Trigger {
	s.Status = &v
	return s
}

func (s *Trigger) SetTargetArn(v string) *Trigger {
	s.TargetArn = &v
	return s
}

func (s *Trigger) SetTriggerConfig(v string) *Trigger {
	s.TriggerConfig = &v
	return s
}

func (s *Trigger) SetTriggerId(v string) *Trigger {
	s.TriggerId = &v
	return s
}

func (s *Trigger) SetTriggerName(v string) *Trigger {
	s.TriggerName = &v
	return s
}

func (s *Trigger) SetTriggerType(v string) *Trigger {
	s.TriggerType = &v
	return s
}

func (s *Trigger) SetUrlInternet(v string) *Trigger {
	s.UrlInternet = &v
	return s
}

func (s *Trigger) SetUrlIntranet(v string) *Trigger {
	s.UrlIntranet = &v
	return s
}

type VPCConfig struct {
	Role            *string   `json:"role,omitempty" xml:"role,omitempty"`
	SecurityGroupId *string   `json:"securityGroupId,omitempty" xml:"securityGroupId,omitempty"`
	VSwitchIds      []*string `json:"vSwitchIds,omitempty" xml:"vSwitchIds,omitempty" type:"Repeated"`
	VpcId           *string   `json:"vpcId,omitempty" xml:"vpcId,omitempty"`
}

func (s VPCConfig) String() string {
	return tea.Prettify(s)
}

func (s VPCConfig) GoString() string {
	return s.String()
}

func (s *VPCConfig) SetRole(v string) *VPCConfig {
	s.Role = &v
	return s
}

func (s *VPCConfig) SetSecurityGroupId(v string) *VPCConfig {
	s.SecurityGroupId = &v
	return s
}

func (s *VPCConfig) SetVSwitchIds(v []*string) *VPCConfig {
	s.VSwitchIds = v
	return s
}

func (s *VPCConfig) SetVpcId(v string) *VPCConfig {
	s.VpcId = &v
	return s
}

type VendorConfig struct {
	MeteringConfig *MeteringConfig `json:"meteringConfig,omitempty" xml:"meteringConfig,omitempty"`
}

func (s VendorConfig) String() string {
	return tea.Prettify(s)
}

func (s VendorConfig) GoString() string {
	return s.String()
}

func (s *VendorConfig) SetMeteringConfig(v *MeteringConfig) *VendorConfig {
	s.MeteringConfig = v
	return s
}

type WAFConfig struct {
	EnableWAF *bool `json:"enableWAF,omitempty" xml:"enableWAF,omitempty"`
}

func (s WAFConfig) String() string {
	return tea.Prettify(s)
}

func (s WAFConfig) GoString() string {
	return s.String()
}

func (s *WAFConfig) SetEnableWAF(v bool) *WAFConfig {
	s.EnableWAF = &v
	return s
}

type ClaimGPUInstanceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ClaimGPUInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s ClaimGPUInstanceHeaders) GoString() string {
	return s.String()
}

func (s *ClaimGPUInstanceHeaders) SetCommonHeaders(v map[string]*string) *ClaimGPUInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ClaimGPUInstanceHeaders) SetXFcAccountId(v string) *ClaimGPUInstanceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ClaimGPUInstanceHeaders) SetXFcDate(v string) *ClaimGPUInstanceHeaders {
	s.XFcDate = &v
	return s
}

func (s *ClaimGPUInstanceHeaders) SetXFcTraceId(v string) *ClaimGPUInstanceHeaders {
	s.XFcTraceId = &v
	return s
}

type ClaimGPUInstanceRequest struct {
	// The disk performance level of the GPU rendering instance.
	DiskPerformanceLevel *string `json:"diskPerformanceLevel,omitempty" xml:"diskPerformanceLevel,omitempty"`
	// The system disk space of the GPU rendering instance.
	DiskSizeGigabytes *string `json:"diskSizeGigabytes,omitempty" xml:"diskSizeGigabytes,omitempty"`
	// The image ID of the GPU-rendered instance.
	ImageId *string `json:"imageId,omitempty" xml:"imageId,omitempty"`
	// The specifications of the GPU rendering instance.
	InstanceType *string `json:"instanceType,omitempty" xml:"instanceType,omitempty"`
	// The outbound Internet bandwidth of the GPU rendering instance.
	InternetBandwidthOut *string `json:"internetBandwidthOut,omitempty" xml:"internetBandwidthOut,omitempty"`
	// The password of the GPU rendering instance.
	Password *string `json:"password,omitempty" xml:"password,omitempty"`
	// The role of the user.
	Role *string `json:"role,omitempty" xml:"role,omitempty"`
	// The ID of the security group.
	SgId *string `json:"sgId,omitempty" xml:"sgId,omitempty"`
	// The source IPv4 CIDR block of the GPU rendering instance.
	SourceCidrIp *string `json:"sourceCidrIp,omitempty" xml:"sourceCidrIp,omitempty"`
	// The range of TCP ports that are open to the security group of the GPU rendering instance.
	TcpPortRange []*string `json:"tcpPortRange,omitempty" xml:"tcpPortRange,omitempty" type:"Repeated"`
	// The range of UDP ports that are open to the security group of the GPU rendering instance.
	UdpPortRange []*string `json:"udpPortRange,omitempty" xml:"udpPortRange,omitempty" type:"Repeated"`
	// The ID of the VPC in which the instance resides.
	VpcId *string `json:"vpcId,omitempty" xml:"vpcId,omitempty"`
	// The ID of the vSwitch.
	VswId *string `json:"vswId,omitempty" xml:"vswId,omitempty"`
}

func (s ClaimGPUInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s ClaimGPUInstanceRequest) GoString() string {
	return s.String()
}

func (s *ClaimGPUInstanceRequest) SetDiskPerformanceLevel(v string) *ClaimGPUInstanceRequest {
	s.DiskPerformanceLevel = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetDiskSizeGigabytes(v string) *ClaimGPUInstanceRequest {
	s.DiskSizeGigabytes = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetImageId(v string) *ClaimGPUInstanceRequest {
	s.ImageId = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetInstanceType(v string) *ClaimGPUInstanceRequest {
	s.InstanceType = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetInternetBandwidthOut(v string) *ClaimGPUInstanceRequest {
	s.InternetBandwidthOut = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetPassword(v string) *ClaimGPUInstanceRequest {
	s.Password = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetRole(v string) *ClaimGPUInstanceRequest {
	s.Role = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetSgId(v string) *ClaimGPUInstanceRequest {
	s.SgId = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetSourceCidrIp(v string) *ClaimGPUInstanceRequest {
	s.SourceCidrIp = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetTcpPortRange(v []*string) *ClaimGPUInstanceRequest {
	s.TcpPortRange = v
	return s
}

func (s *ClaimGPUInstanceRequest) SetUdpPortRange(v []*string) *ClaimGPUInstanceRequest {
	s.UdpPortRange = v
	return s
}

func (s *ClaimGPUInstanceRequest) SetVpcId(v string) *ClaimGPUInstanceRequest {
	s.VpcId = &v
	return s
}

func (s *ClaimGPUInstanceRequest) SetVswId(v string) *ClaimGPUInstanceRequest {
	s.VswId = &v
	return s
}

type ClaimGPUInstanceResponseBody struct {
	// The time when the product instance is created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The ID of the instance that you query.
	InstanceId *string `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	// The public IP address of the server.
	PublicIp *string `json:"publicIp,omitempty" xml:"publicIp,omitempty"`
}

func (s ClaimGPUInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ClaimGPUInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *ClaimGPUInstanceResponseBody) SetCreatedTime(v string) *ClaimGPUInstanceResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *ClaimGPUInstanceResponseBody) SetInstanceId(v string) *ClaimGPUInstanceResponseBody {
	s.InstanceId = &v
	return s
}

func (s *ClaimGPUInstanceResponseBody) SetPublicIp(v string) *ClaimGPUInstanceResponseBody {
	s.PublicIp = &v
	return s
}

type ClaimGPUInstanceResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ClaimGPUInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ClaimGPUInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s ClaimGPUInstanceResponse) GoString() string {
	return s.String()
}

func (s *ClaimGPUInstanceResponse) SetHeaders(v map[string]*string) *ClaimGPUInstanceResponse {
	s.Headers = v
	return s
}

func (s *ClaimGPUInstanceResponse) SetStatusCode(v int32) *ClaimGPUInstanceResponse {
	s.StatusCode = &v
	return s
}

func (s *ClaimGPUInstanceResponse) SetBody(v *ClaimGPUInstanceResponseBody) *ClaimGPUInstanceResponse {
	s.Body = v
	return s
}

type CreateAliasHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s CreateAliasHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateAliasHeaders) GoString() string {
	return s.String()
}

func (s *CreateAliasHeaders) SetCommonHeaders(v map[string]*string) *CreateAliasHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateAliasHeaders) SetXFcAccountId(v string) *CreateAliasHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *CreateAliasHeaders) SetXFcDate(v string) *CreateAliasHeaders {
	s.XFcDate = &v
	return s
}

func (s *CreateAliasHeaders) SetXFcTraceId(v string) *CreateAliasHeaders {
	s.XFcTraceId = &v
	return s
}

type CreateAliasRequest struct {
	// The canary release version to which the alias points and the weight of the canary release version.
	//
	// *   The canary release version takes effect only when the function is invoked.
	// *   The value consists of a version number and a specific weight. For example, 2:0.05 indicates that when a function is invoked, Version 2 is the canary release version, 5% of the traffic is distributed to the canary release version, and 95% of the traffic is distributed to the major version.
	AdditionalVersionWeight map[string]*float32 `json:"additionalVersionWeight,omitempty" xml:"additionalVersionWeight,omitempty"`
	// The name of the alias.  The name contains only letters, digits, hyphens (-), and underscores (\_). The name must be 1 to 128 characters in length and cannot start with a digit or hyphen (-).  The name cannot be **LATEST**.
	AliasName *string `json:"aliasName,omitempty" xml:"aliasName,omitempty"`
	// The description of the alias.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The canary release mode. Valid values:
	//
	// *   **Random**: random canary release. This is the default value.
	// *   **Content**: rule-based canary release.
	ResolvePolicy *string `json:"resolvePolicy,omitempty" xml:"resolvePolicy,omitempty"`
	// The canary release rule. Traffic that meets the canary release rule is routed to the canary release instance.
	RoutePolicy *RoutePolicy `json:"routePolicy,omitempty" xml:"routePolicy,omitempty"`
	// The ID of the version to which the alias points.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s CreateAliasRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAliasRequest) GoString() string {
	return s.String()
}

func (s *CreateAliasRequest) SetAdditionalVersionWeight(v map[string]*float32) *CreateAliasRequest {
	s.AdditionalVersionWeight = v
	return s
}

func (s *CreateAliasRequest) SetAliasName(v string) *CreateAliasRequest {
	s.AliasName = &v
	return s
}

func (s *CreateAliasRequest) SetDescription(v string) *CreateAliasRequest {
	s.Description = &v
	return s
}

func (s *CreateAliasRequest) SetResolvePolicy(v string) *CreateAliasRequest {
	s.ResolvePolicy = &v
	return s
}

func (s *CreateAliasRequest) SetRoutePolicy(v *RoutePolicy) *CreateAliasRequest {
	s.RoutePolicy = v
	return s
}

func (s *CreateAliasRequest) SetVersionId(v string) *CreateAliasRequest {
	s.VersionId = &v
	return s
}

type CreateAliasResponseBody struct {
	// The canary release version to which the alias points and the weight of the canary release version.
	//
	// *   The canary release version takes effect only when the function is invoked.
	// *   The value consists of a version number and a specific weight. For example, 2:0.05 indicates that when a function is invoked, Version 2 is the canary release version, 5% of the traffic is distributed to the canary release version, and 95% of the traffic is distributed to the major version.
	AdditionalVersionWeight map[string]*float32 `json:"additionalVersionWeight,omitempty" xml:"additionalVersionWeight,omitempty"`
	// The name of the alias.
	AliasName *string `json:"aliasName,omitempty" xml:"aliasName,omitempty"`
	// The time when the alias was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the alias.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The time when the alias was last modified.
	LastModifiedTime *string      `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	ResolvePolicy    *string      `json:"resolvePolicy,omitempty" xml:"resolvePolicy,omitempty"`
	RoutePolicy      *RoutePolicy `json:"routePolicy,omitempty" xml:"routePolicy,omitempty"`
	// The ID of the version to which the alias points.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s CreateAliasResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateAliasResponseBody) GoString() string {
	return s.String()
}

func (s *CreateAliasResponseBody) SetAdditionalVersionWeight(v map[string]*float32) *CreateAliasResponseBody {
	s.AdditionalVersionWeight = v
	return s
}

func (s *CreateAliasResponseBody) SetAliasName(v string) *CreateAliasResponseBody {
	s.AliasName = &v
	return s
}

func (s *CreateAliasResponseBody) SetCreatedTime(v string) *CreateAliasResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *CreateAliasResponseBody) SetDescription(v string) *CreateAliasResponseBody {
	s.Description = &v
	return s
}

func (s *CreateAliasResponseBody) SetLastModifiedTime(v string) *CreateAliasResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *CreateAliasResponseBody) SetResolvePolicy(v string) *CreateAliasResponseBody {
	s.ResolvePolicy = &v
	return s
}

func (s *CreateAliasResponseBody) SetRoutePolicy(v *RoutePolicy) *CreateAliasResponseBody {
	s.RoutePolicy = v
	return s
}

func (s *CreateAliasResponseBody) SetVersionId(v string) *CreateAliasResponseBody {
	s.VersionId = &v
	return s
}

type CreateAliasResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateAliasResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateAliasResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAliasResponse) GoString() string {
	return s.String()
}

func (s *CreateAliasResponse) SetHeaders(v map[string]*string) *CreateAliasResponse {
	s.Headers = v
	return s
}

func (s *CreateAliasResponse) SetStatusCode(v int32) *CreateAliasResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateAliasResponse) SetBody(v *CreateAliasResponseBody) *CreateAliasResponse {
	s.Body = v
	return s
}

type CreateCustomDomainHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the operation is called. The format is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s CreateCustomDomainHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateCustomDomainHeaders) GoString() string {
	return s.String()
}

func (s *CreateCustomDomainHeaders) SetCommonHeaders(v map[string]*string) *CreateCustomDomainHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateCustomDomainHeaders) SetXFcAccountId(v string) *CreateCustomDomainHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *CreateCustomDomainHeaders) SetXFcDate(v string) *CreateCustomDomainHeaders {
	s.XFcDate = &v
	return s
}

func (s *CreateCustomDomainHeaders) SetXFcTraceId(v string) *CreateCustomDomainHeaders {
	s.XFcTraceId = &v
	return s
}

type CreateCustomDomainRequest struct {
	// The configurations of the HTTPS certificate.
	CertConfig *CertConfig `json:"certConfig,omitempty" xml:"certConfig,omitempty"`
	// The domain name. Enter a custom domain name that has obtained an ICP filing in the Alibaba Cloud ICP Filing system, or a custom domain name whose ICP filing information includes Alibaba Cloud as a service provider.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The protocol types supported by the domain name. Valid values:
	//
	// *   **HTTP**: Only HTTP is supported.
	// *   **HTTPS**: Only HTTPS is supported.
	// *   **HTTP,HTTPS**: HTTP and HTTPS are supported.
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// The route table that maps the paths to functions when the functions are invoked by using the custom domain name.
	RouteConfig *RouteConfig `json:"routeConfig,omitempty" xml:"routeConfig,omitempty"`
	// The Transport Layer Security (TLS) configuration.
	TlsConfig *TLSConfig `json:"tlsConfig,omitempty" xml:"tlsConfig,omitempty"`
	// The Web Application Firewall (WAF) configuration.
	WafConfig *WAFConfig `json:"wafConfig,omitempty" xml:"wafConfig,omitempty"`
}

func (s CreateCustomDomainRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateCustomDomainRequest) GoString() string {
	return s.String()
}

func (s *CreateCustomDomainRequest) SetCertConfig(v *CertConfig) *CreateCustomDomainRequest {
	s.CertConfig = v
	return s
}

func (s *CreateCustomDomainRequest) SetDomainName(v string) *CreateCustomDomainRequest {
	s.DomainName = &v
	return s
}

func (s *CreateCustomDomainRequest) SetProtocol(v string) *CreateCustomDomainRequest {
	s.Protocol = &v
	return s
}

func (s *CreateCustomDomainRequest) SetRouteConfig(v *RouteConfig) *CreateCustomDomainRequest {
	s.RouteConfig = v
	return s
}

func (s *CreateCustomDomainRequest) SetTlsConfig(v *TLSConfig) *CreateCustomDomainRequest {
	s.TlsConfig = v
	return s
}

func (s *CreateCustomDomainRequest) SetWafConfig(v *WAFConfig) *CreateCustomDomainRequest {
	s.WafConfig = v
	return s
}

type CreateCustomDomainResponseBody struct {
	// The ID of your Alibaba Cloud account.
	AccountId *string `json:"accountId,omitempty" xml:"accountId,omitempty"`
	// The version of the API.
	ApiVersion *string `json:"apiVersion,omitempty" xml:"apiVersion,omitempty"`
	// The configurations of the HTTPS certificate.
	CertConfig *CertConfig `json:"certConfig,omitempty" xml:"certConfig,omitempty"`
	// The time when the domain name was added.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The domain name.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The time when the domain name was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The protocol types supported by the domain name. Valid values:
	//
	// *   **HTTP**: Only HTTP is supported.
	// *   **HTTPS**: Only HTTPS is supported.
	// *   **HTTP,HTTPS**: HTTP and HTTPS are supported.
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// The route table that maps the paths to functions when the functions are invoked by using the custom domain name.
	RouteConfig *RouteConfig `json:"routeConfig,omitempty" xml:"routeConfig,omitempty"`
	// The Transport Layer Security (TLS) configuration.
	TlsConfig *TLSConfig `json:"tlsConfig,omitempty" xml:"tlsConfig,omitempty"`
	// The Web Application Firewall (WAF) configuration.
	WafConfig *WAFConfig `json:"wafConfig,omitempty" xml:"wafConfig,omitempty"`
}

func (s CreateCustomDomainResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateCustomDomainResponseBody) GoString() string {
	return s.String()
}

func (s *CreateCustomDomainResponseBody) SetAccountId(v string) *CreateCustomDomainResponseBody {
	s.AccountId = &v
	return s
}

func (s *CreateCustomDomainResponseBody) SetApiVersion(v string) *CreateCustomDomainResponseBody {
	s.ApiVersion = &v
	return s
}

func (s *CreateCustomDomainResponseBody) SetCertConfig(v *CertConfig) *CreateCustomDomainResponseBody {
	s.CertConfig = v
	return s
}

func (s *CreateCustomDomainResponseBody) SetCreatedTime(v string) *CreateCustomDomainResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *CreateCustomDomainResponseBody) SetDomainName(v string) *CreateCustomDomainResponseBody {
	s.DomainName = &v
	return s
}

func (s *CreateCustomDomainResponseBody) SetLastModifiedTime(v string) *CreateCustomDomainResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *CreateCustomDomainResponseBody) SetProtocol(v string) *CreateCustomDomainResponseBody {
	s.Protocol = &v
	return s
}

func (s *CreateCustomDomainResponseBody) SetRouteConfig(v *RouteConfig) *CreateCustomDomainResponseBody {
	s.RouteConfig = v
	return s
}

func (s *CreateCustomDomainResponseBody) SetTlsConfig(v *TLSConfig) *CreateCustomDomainResponseBody {
	s.TlsConfig = v
	return s
}

func (s *CreateCustomDomainResponseBody) SetWafConfig(v *WAFConfig) *CreateCustomDomainResponseBody {
	s.WafConfig = v
	return s
}

type CreateCustomDomainResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateCustomDomainResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateCustomDomainResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateCustomDomainResponse) GoString() string {
	return s.String()
}

func (s *CreateCustomDomainResponse) SetHeaders(v map[string]*string) *CreateCustomDomainResponse {
	s.Headers = v
	return s
}

func (s *CreateCustomDomainResponse) SetStatusCode(v int32) *CreateCustomDomainResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateCustomDomainResponse) SetBody(v *CreateCustomDomainResponseBody) *CreateCustomDomainResponse {
	s.Body = v
	return s
}

type CreateFunctionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The CRC-64 value of the function code package.
	XFcCodeChecksum *string `json:"X-Fc-Code-Checksum,omitempty" xml:"X-Fc-Code-Checksum,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request. The value is the same as that of the requestId parameter in the response.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s CreateFunctionHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateFunctionHeaders) GoString() string {
	return s.String()
}

func (s *CreateFunctionHeaders) SetCommonHeaders(v map[string]*string) *CreateFunctionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateFunctionHeaders) SetXFcAccountId(v string) *CreateFunctionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *CreateFunctionHeaders) SetXFcCodeChecksum(v string) *CreateFunctionHeaders {
	s.XFcCodeChecksum = &v
	return s
}

func (s *CreateFunctionHeaders) SetXFcDate(v string) *CreateFunctionHeaders {
	s.XFcDate = &v
	return s
}

func (s *CreateFunctionHeaders) SetXFcTraceId(v string) *CreateFunctionHeaders {
	s.XFcTraceId = &v
	return s
}

type CreateFunctionRequest struct {
	// The port on which the HTTP server listens for the custom runtime or custom container runtime.
	CaPort *int32 `json:"caPort,omitempty" xml:"caPort,omitempty"`
	// The code of the function. The code must be packaged into a ZIP file. Choose **code** or **customContainerConfig** for the function.
	Code *Code `json:"code,omitempty" xml:"code,omitempty"`
	// The number of vCPUs of the function. The value must be a multiple of 0.05.
	Cpu *float32 `json:"cpu,omitempty" xml:"cpu,omitempty"`
	// The configurations of the custom container runtime. After you configure the custom container runtime, Function Compute can execute the function in a container created from a custom image. Choose **code** or **customContainerConfig** for the function.
	CustomContainerConfig *CustomContainerConfig `json:"customContainerConfig,omitempty" xml:"customContainerConfig,omitempty"`
	// The custom Domain Name System (DNS) configurations of the function.
	CustomDNS *CustomDNS `json:"customDNS,omitempty" xml:"customDNS,omitempty"`
	// The custom health check configurations of the function. This parameter is applicable to only custom runtimes and custom containers.
	CustomHealthCheckConfig *CustomHealthCheckConfig `json:"customHealthCheckConfig,omitempty" xml:"customHealthCheckConfig,omitempty"`
	// The configurations of the custom runtime.
	CustomRuntimeConfig *CustomRuntimeConfig `json:"customRuntimeConfig,omitempty" xml:"customRuntimeConfig,omitempty"`
	// The description of the function.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The disk size of the function. Unit: MB. Valid values: 512 and 10240.
	DiskSize *int32 `json:"diskSize,omitempty" xml:"diskSize,omitempty"`
	// The environment variables that you configured for the function. You can obtain the values of the environment variables from the function. For more information, see [Overview](~~69777~~).
	EnvironmentVariables map[string]*string `json:"environmentVariables,omitempty" xml:"environmentVariables,omitempty"`
	// The name of the function. The name can contain letters, digits, underscores (\_), and hyphens (-) only. The name cannot start with a digit or a hyphen (-). The name must be 1 to 64 characters in length.
	FunctionName *string `json:"functionName,omitempty" xml:"functionName,omitempty"`
	// The GPU memory capacity for the function. Unit: MB. The value must be a multiple of 1,024.
	GpuMemorySize *int32 `json:"gpuMemorySize,omitempty" xml:"gpuMemorySize,omitempty"`
	// The handler of the function. The format varies based on the programming language. For more information, see [Function handlers](~~157704~~).
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	// The timeout period for the execution of the Initializer hook. Unit: seconds. Default value: 3. Valid values: 1 to 300. When this period expires, the execution of the Initializer hook is terminated.
	InitializationTimeout *int32 `json:"initializationTimeout,omitempty" xml:"initializationTimeout,omitempty"`
	// The handler of the Initializer hook. For more information, see [Initializer hook](~~157704~~).
	Initializer *string `json:"initializer,omitempty" xml:"initializer,omitempty"`
	// The number of requests that can be concurrently processed by a single instance.
	InstanceConcurrency *int32 `json:"instanceConcurrency,omitempty" xml:"instanceConcurrency,omitempty"`
	// The lifecycle configurations of the instance.
	InstanceLifecycleConfig *InstanceLifecycleConfig `json:"instanceLifecycleConfig,omitempty" xml:"instanceLifecycleConfig,omitempty"`
	// The soft concurrency of the instance. You can use this parameter to implement graceful scale-up of instances. If the number of concurrent requests on an instance is greater than the value of soft concurrency, an instance scale-up is triggered. For example, if your instance requires a long time to start, you can specify a suitable soft concurrency to start the instance in advance.
	//
	// The value must be less than or equal to that of the **instanceConcurrency** parameter.
	InstanceSoftConcurrency *int32 `json:"instanceSoftConcurrency,omitempty" xml:"instanceSoftConcurrency,omitempty"`
	// The instance type of the function. Valid values:
	//
	// *   **e1**: elastic instance
	// *   **c1**: performance instance
	// *   **fc.gpu.tesla.1**: GPU-accelerated instance (Tesla T4)
	// *   **fc.gpu.ampere.1**: GPU-accelerated instance (Ampere A10)
	// *   **g1**: same as **fc.gpu.tesla.1**
	InstanceType *string `json:"instanceType,omitempty" xml:"instanceType,omitempty"`
	// The information about layers.
	//
	// > Multiple layers are merged based on the order of array subscripts. The content of a layer with a smaller subscript overwrites the file with the same name as a layer with a larger subscript.
	Layers []*string `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The memory size for the function. Unit: MB. The memory size must be a multiple of 64 MB. The memory size varies based on the function instance type. For more information, see [Instance types](~~179379~~).
	MemorySize *int32 `json:"memorySize,omitempty" xml:"memorySize,omitempty"`
	// The runtime environment of the function. Valid values: **nodejs16**, **nodejs14**, **nodejs12**, **nodejs10**, **nodejs8**, **nodejs6**, **nodejs4.4**, **python3.9**, **python3**, **python2.7**, **java11**, **java8**, **go1**, **php7.2**, **dotnetcore3.1**, **dotnetcore2.1**, **custom** and **custom-container**. For more information, see [Supported function runtime environments](~~73338~~).
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The timeout period for the execution of the function. Unit: seconds. Default value: 3. Minimum value: 1. When the period ends, the execution of the function is terminated.
	Timeout *int32 `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s CreateFunctionRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateFunctionRequest) GoString() string {
	return s.String()
}

func (s *CreateFunctionRequest) SetCaPort(v int32) *CreateFunctionRequest {
	s.CaPort = &v
	return s
}

func (s *CreateFunctionRequest) SetCode(v *Code) *CreateFunctionRequest {
	s.Code = v
	return s
}

func (s *CreateFunctionRequest) SetCpu(v float32) *CreateFunctionRequest {
	s.Cpu = &v
	return s
}

func (s *CreateFunctionRequest) SetCustomContainerConfig(v *CustomContainerConfig) *CreateFunctionRequest {
	s.CustomContainerConfig = v
	return s
}

func (s *CreateFunctionRequest) SetCustomDNS(v *CustomDNS) *CreateFunctionRequest {
	s.CustomDNS = v
	return s
}

func (s *CreateFunctionRequest) SetCustomHealthCheckConfig(v *CustomHealthCheckConfig) *CreateFunctionRequest {
	s.CustomHealthCheckConfig = v
	return s
}

func (s *CreateFunctionRequest) SetCustomRuntimeConfig(v *CustomRuntimeConfig) *CreateFunctionRequest {
	s.CustomRuntimeConfig = v
	return s
}

func (s *CreateFunctionRequest) SetDescription(v string) *CreateFunctionRequest {
	s.Description = &v
	return s
}

func (s *CreateFunctionRequest) SetDiskSize(v int32) *CreateFunctionRequest {
	s.DiskSize = &v
	return s
}

func (s *CreateFunctionRequest) SetEnvironmentVariables(v map[string]*string) *CreateFunctionRequest {
	s.EnvironmentVariables = v
	return s
}

func (s *CreateFunctionRequest) SetFunctionName(v string) *CreateFunctionRequest {
	s.FunctionName = &v
	return s
}

func (s *CreateFunctionRequest) SetGpuMemorySize(v int32) *CreateFunctionRequest {
	s.GpuMemorySize = &v
	return s
}

func (s *CreateFunctionRequest) SetHandler(v string) *CreateFunctionRequest {
	s.Handler = &v
	return s
}

func (s *CreateFunctionRequest) SetInitializationTimeout(v int32) *CreateFunctionRequest {
	s.InitializationTimeout = &v
	return s
}

func (s *CreateFunctionRequest) SetInitializer(v string) *CreateFunctionRequest {
	s.Initializer = &v
	return s
}

func (s *CreateFunctionRequest) SetInstanceConcurrency(v int32) *CreateFunctionRequest {
	s.InstanceConcurrency = &v
	return s
}

func (s *CreateFunctionRequest) SetInstanceLifecycleConfig(v *InstanceLifecycleConfig) *CreateFunctionRequest {
	s.InstanceLifecycleConfig = v
	return s
}

func (s *CreateFunctionRequest) SetInstanceSoftConcurrency(v int32) *CreateFunctionRequest {
	s.InstanceSoftConcurrency = &v
	return s
}

func (s *CreateFunctionRequest) SetInstanceType(v string) *CreateFunctionRequest {
	s.InstanceType = &v
	return s
}

func (s *CreateFunctionRequest) SetLayers(v []*string) *CreateFunctionRequest {
	s.Layers = v
	return s
}

func (s *CreateFunctionRequest) SetMemorySize(v int32) *CreateFunctionRequest {
	s.MemorySize = &v
	return s
}

func (s *CreateFunctionRequest) SetRuntime(v string) *CreateFunctionRequest {
	s.Runtime = &v
	return s
}

func (s *CreateFunctionRequest) SetTimeout(v int32) *CreateFunctionRequest {
	s.Timeout = &v
	return s
}

type CreateFunctionResponseBody struct {
	// The port on which the HTTP server listens for the custom runtime or custom container runtime.
	CaPort *int32 `json:"caPort,omitempty" xml:"caPort,omitempty"`
	// The CRC-64 value of the function code package.
	CodeChecksum *string `json:"codeChecksum,omitempty" xml:"codeChecksum,omitempty"`
	// The size of the function code package that is returned by the system. Unit: bytes.
	CodeSize *int64 `json:"codeSize,omitempty" xml:"codeSize,omitempty"`
	// The number of vCPUs of the function. The value is a multiple of 0.05.
	Cpu *float32 `json:"cpu,omitempty" xml:"cpu,omitempty"`
	// The time when the function was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The configurations of the custom container runtime. After you configure the custom container runtime, Function Compute can execute the function in a container created from a custom image.
	CustomContainerConfig *CustomContainerConfig `json:"customContainerConfig,omitempty" xml:"customContainerConfig,omitempty"`
	// The custom DNS configurations of the function.
	CustomDNS *CustomDNS `json:"customDNS,omitempty" xml:"customDNS,omitempty"`
	// The custom health check configuration of the function. This parameter is applicable only to custom runtimes and custom containers.
	CustomHealthCheckConfig *CustomHealthCheckConfig `json:"customHealthCheckConfig,omitempty" xml:"customHealthCheckConfig,omitempty"`
	// The configurations of the custom runtime.
	CustomRuntimeConfig *CustomRuntimeConfig `json:"customRuntimeConfig,omitempty" xml:"customRuntimeConfig,omitempty"`
	// The description of the function.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The disk size of the function. Unit: MB. Valid values: 512 and 10240.
	DiskSize *int32 `json:"diskSize,omitempty" xml:"diskSize,omitempty"`
	// The environment variables that are configured for the function. You can obtain the values of the environment variables from the function. For more information, see [Environment variables](~~69777~~).
	EnvironmentVariables map[string]*string `json:"environmentVariables,omitempty" xml:"environmentVariables,omitempty"`
	// The unique ID that is generated by the system for the function.
	FunctionId *string `json:"functionId,omitempty" xml:"functionId,omitempty"`
	// The name of the function.
	FunctionName *string `json:"functionName,omitempty" xml:"functionName,omitempty"`
	// The GPU memory capacity for the function. Unit: MB. The value is a multiple of 1,024.
	GpuMemorySize *int32 `json:"gpuMemorySize,omitempty" xml:"gpuMemorySize,omitempty"`
	// The handler of the function.
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	// The timeout period for the execution of the Initializer hook. Unit: seconds. Default value: 3. Minimum value: 1. When the period ends, the execution of the Initializer hook is terminated.
	InitializationTimeout *int32 `json:"initializationTimeout,omitempty" xml:"initializationTimeout,omitempty"`
	// The handler of the Initializer hook. The format is determined by the programming language.
	Initializer *string `json:"initializer,omitempty" xml:"initializer,omitempty"`
	// The number of requests that can be concurrently processed by a single instance.
	InstanceConcurrency *int32 `json:"instanceConcurrency,omitempty" xml:"instanceConcurrency,omitempty"`
	// The lifecycle configurations of the instance.
	InstanceLifecycleConfig *InstanceLifecycleConfig `json:"instanceLifecycleConfig,omitempty" xml:"instanceLifecycleConfig,omitempty"`
	// The soft concurrency of the instance. You can use this parameter to implement graceful scale-up of instances. If the number of concurrent requests on an instance is greater than the value of soft concurrency, an instance scale-up is triggered. For example, if your instance requires a long time to start, you can specify a suitable soft concurrency to start the instance in advance.
	//
	// The value must be less than or equal to that of the **instanceConcurrency** parameter.
	InstanceSoftConcurrency *int32 `json:"instanceSoftConcurrency,omitempty" xml:"instanceSoftConcurrency,omitempty"`
	// The instance type of the function. Valid values:
	//
	// *   **e1**: elastic instance
	// *   **c1**: performance instance
	// *   **fc.gpu.tesla.1**: GPU-accelerated instance (Tesla T4)
	// *   **fc.gpu.ampere.1**: GPU-accelerated instance (Ampere A10)
	// *   **g1**: same as **fc.gpu.tesla.1**
	InstanceType *string `json:"instanceType,omitempty" xml:"instanceType,omitempty"`
	// The time when the function was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// An array that consists of the information of layers.
	//
	// > Multiple layers are merged based on the order of array subscripts. The content of a layer with a smaller subscript overwrites the file with the same name as a layer with a larger subscript.
	Layers []*string `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The memory size that is configured for the function. Unit: MB.
	MemorySize *int32 `json:"memorySize,omitempty" xml:"memorySize,omitempty"`
	// The runtime environment of the function. Valid values: **nodejs16**, **nodejs14**, **nodejs12**, **nodejs10**, **nodejs8**, **nodejs6**, **nodejs4.4**, **python3.9**, **python3**, **python2.7**, **java11**, **java8**, **go1**, **php7.2**, **dotnetcore3.1**, **dotnetcore2.1**, **custom** and **custom-container**. For more information, see [Supported function runtime environments](~~73338~~).
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The timeout period for the execution of the function. Unit: seconds. Default value: 60. Valid values: 1 to 600. When this period expires, the execution of the function is terminated.
	Timeout *int32 `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s CreateFunctionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateFunctionResponseBody) GoString() string {
	return s.String()
}

func (s *CreateFunctionResponseBody) SetCaPort(v int32) *CreateFunctionResponseBody {
	s.CaPort = &v
	return s
}

func (s *CreateFunctionResponseBody) SetCodeChecksum(v string) *CreateFunctionResponseBody {
	s.CodeChecksum = &v
	return s
}

func (s *CreateFunctionResponseBody) SetCodeSize(v int64) *CreateFunctionResponseBody {
	s.CodeSize = &v
	return s
}

func (s *CreateFunctionResponseBody) SetCpu(v float32) *CreateFunctionResponseBody {
	s.Cpu = &v
	return s
}

func (s *CreateFunctionResponseBody) SetCreatedTime(v string) *CreateFunctionResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *CreateFunctionResponseBody) SetCustomContainerConfig(v *CustomContainerConfig) *CreateFunctionResponseBody {
	s.CustomContainerConfig = v
	return s
}

func (s *CreateFunctionResponseBody) SetCustomDNS(v *CustomDNS) *CreateFunctionResponseBody {
	s.CustomDNS = v
	return s
}

func (s *CreateFunctionResponseBody) SetCustomHealthCheckConfig(v *CustomHealthCheckConfig) *CreateFunctionResponseBody {
	s.CustomHealthCheckConfig = v
	return s
}

func (s *CreateFunctionResponseBody) SetCustomRuntimeConfig(v *CustomRuntimeConfig) *CreateFunctionResponseBody {
	s.CustomRuntimeConfig = v
	return s
}

func (s *CreateFunctionResponseBody) SetDescription(v string) *CreateFunctionResponseBody {
	s.Description = &v
	return s
}

func (s *CreateFunctionResponseBody) SetDiskSize(v int32) *CreateFunctionResponseBody {
	s.DiskSize = &v
	return s
}

func (s *CreateFunctionResponseBody) SetEnvironmentVariables(v map[string]*string) *CreateFunctionResponseBody {
	s.EnvironmentVariables = v
	return s
}

func (s *CreateFunctionResponseBody) SetFunctionId(v string) *CreateFunctionResponseBody {
	s.FunctionId = &v
	return s
}

func (s *CreateFunctionResponseBody) SetFunctionName(v string) *CreateFunctionResponseBody {
	s.FunctionName = &v
	return s
}

func (s *CreateFunctionResponseBody) SetGpuMemorySize(v int32) *CreateFunctionResponseBody {
	s.GpuMemorySize = &v
	return s
}

func (s *CreateFunctionResponseBody) SetHandler(v string) *CreateFunctionResponseBody {
	s.Handler = &v
	return s
}

func (s *CreateFunctionResponseBody) SetInitializationTimeout(v int32) *CreateFunctionResponseBody {
	s.InitializationTimeout = &v
	return s
}

func (s *CreateFunctionResponseBody) SetInitializer(v string) *CreateFunctionResponseBody {
	s.Initializer = &v
	return s
}

func (s *CreateFunctionResponseBody) SetInstanceConcurrency(v int32) *CreateFunctionResponseBody {
	s.InstanceConcurrency = &v
	return s
}

func (s *CreateFunctionResponseBody) SetInstanceLifecycleConfig(v *InstanceLifecycleConfig) *CreateFunctionResponseBody {
	s.InstanceLifecycleConfig = v
	return s
}

func (s *CreateFunctionResponseBody) SetInstanceSoftConcurrency(v int32) *CreateFunctionResponseBody {
	s.InstanceSoftConcurrency = &v
	return s
}

func (s *CreateFunctionResponseBody) SetInstanceType(v string) *CreateFunctionResponseBody {
	s.InstanceType = &v
	return s
}

func (s *CreateFunctionResponseBody) SetLastModifiedTime(v string) *CreateFunctionResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *CreateFunctionResponseBody) SetLayers(v []*string) *CreateFunctionResponseBody {
	s.Layers = v
	return s
}

func (s *CreateFunctionResponseBody) SetMemorySize(v int32) *CreateFunctionResponseBody {
	s.MemorySize = &v
	return s
}

func (s *CreateFunctionResponseBody) SetRuntime(v string) *CreateFunctionResponseBody {
	s.Runtime = &v
	return s
}

func (s *CreateFunctionResponseBody) SetTimeout(v int32) *CreateFunctionResponseBody {
	s.Timeout = &v
	return s
}

type CreateFunctionResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateFunctionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateFunctionResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateFunctionResponse) GoString() string {
	return s.String()
}

func (s *CreateFunctionResponse) SetHeaders(v map[string]*string) *CreateFunctionResponse {
	s.Headers = v
	return s
}

func (s *CreateFunctionResponse) SetStatusCode(v int32) *CreateFunctionResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateFunctionResponse) SetBody(v *CreateFunctionResponseBody) *CreateFunctionResponse {
	s.Body = v
	return s
}

type CreateLayerVersionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s CreateLayerVersionHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateLayerVersionHeaders) GoString() string {
	return s.String()
}

func (s *CreateLayerVersionHeaders) SetCommonHeaders(v map[string]*string) *CreateLayerVersionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateLayerVersionHeaders) SetXFcAccountId(v string) *CreateLayerVersionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *CreateLayerVersionHeaders) SetXFcDate(v string) *CreateLayerVersionHeaders {
	s.XFcDate = &v
	return s
}

func (s *CreateLayerVersionHeaders) SetXFcTraceId(v string) *CreateLayerVersionHeaders {
	s.XFcTraceId = &v
	return s
}

type CreateLayerVersionRequest struct {
	// The code of the layer.
	Code *Code `json:"Code,omitempty" xml:"Code,omitempty"`
	// The list of runtime environments that are supported by the layer.
	CompatibleRuntime []*string `json:"compatibleRuntime,omitempty" xml:"compatibleRuntime,omitempty" type:"Repeated"`
	// The description of the layer.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
}

func (s CreateLayerVersionRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateLayerVersionRequest) GoString() string {
	return s.String()
}

func (s *CreateLayerVersionRequest) SetCode(v *Code) *CreateLayerVersionRequest {
	s.Code = v
	return s
}

func (s *CreateLayerVersionRequest) SetCompatibleRuntime(v []*string) *CreateLayerVersionRequest {
	s.CompatibleRuntime = v
	return s
}

func (s *CreateLayerVersionRequest) SetDescription(v string) *CreateLayerVersionRequest {
	s.Description = &v
	return s
}

type CreateLayerVersionResponseBody struct {
	// The access mode of the layer.
	Acl *int32 `json:"acl,omitempty" xml:"acl,omitempty"`
	// The name of the layer.
	Arn *string `json:"arn,omitempty" xml:"arn,omitempty"`
	// The information about the layer code package.
	Code *OutputCodeLocation `json:"code,omitempty" xml:"code,omitempty"`
	// The checksum of the layer code package.
	CodeChecksum *string `json:"codeChecksum,omitempty" xml:"codeChecksum,omitempty"`
	// The size of the layer code package. Unit: Byte.
	Codesize *int64 `json:"codesize,omitempty" xml:"codesize,omitempty"`
	// The list of runtime environments that are supported by the layer.
	CompatibleRuntime []*string `json:"compatibleRuntime,omitempty" xml:"compatibleRuntime,omitempty" type:"Repeated"`
	// The time when the layer version was created. The time follows the **yyyy-MM-ddTHH:mm:ssZ** format.
	CreateTime *string `json:"createTime,omitempty" xml:"createTime,omitempty"`
	// The description of the layer version.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The name of the layer.
	LayerName *string `json:"layerName,omitempty" xml:"layerName,omitempty"`
	// The version of the layer.
	Version *int32 `json:"version,omitempty" xml:"version,omitempty"`
}

func (s CreateLayerVersionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateLayerVersionResponseBody) GoString() string {
	return s.String()
}

func (s *CreateLayerVersionResponseBody) SetAcl(v int32) *CreateLayerVersionResponseBody {
	s.Acl = &v
	return s
}

func (s *CreateLayerVersionResponseBody) SetArn(v string) *CreateLayerVersionResponseBody {
	s.Arn = &v
	return s
}

func (s *CreateLayerVersionResponseBody) SetCode(v *OutputCodeLocation) *CreateLayerVersionResponseBody {
	s.Code = v
	return s
}

func (s *CreateLayerVersionResponseBody) SetCodeChecksum(v string) *CreateLayerVersionResponseBody {
	s.CodeChecksum = &v
	return s
}

func (s *CreateLayerVersionResponseBody) SetCodesize(v int64) *CreateLayerVersionResponseBody {
	s.Codesize = &v
	return s
}

func (s *CreateLayerVersionResponseBody) SetCompatibleRuntime(v []*string) *CreateLayerVersionResponseBody {
	s.CompatibleRuntime = v
	return s
}

func (s *CreateLayerVersionResponseBody) SetCreateTime(v string) *CreateLayerVersionResponseBody {
	s.CreateTime = &v
	return s
}

func (s *CreateLayerVersionResponseBody) SetDescription(v string) *CreateLayerVersionResponseBody {
	s.Description = &v
	return s
}

func (s *CreateLayerVersionResponseBody) SetLayerName(v string) *CreateLayerVersionResponseBody {
	s.LayerName = &v
	return s
}

func (s *CreateLayerVersionResponseBody) SetVersion(v int32) *CreateLayerVersionResponseBody {
	s.Version = &v
	return s
}

type CreateLayerVersionResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateLayerVersionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateLayerVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateLayerVersionResponse) GoString() string {
	return s.String()
}

func (s *CreateLayerVersionResponse) SetHeaders(v map[string]*string) *CreateLayerVersionResponse {
	s.Headers = v
	return s
}

func (s *CreateLayerVersionResponse) SetStatusCode(v int32) *CreateLayerVersionResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateLayerVersionResponse) SetBody(v *CreateLayerVersionResponseBody) *CreateLayerVersionResponse {
	s.Body = v
	return s
}

type CreateServiceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s CreateServiceHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateServiceHeaders) GoString() string {
	return s.String()
}

func (s *CreateServiceHeaders) SetCommonHeaders(v map[string]*string) *CreateServiceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateServiceHeaders) SetXFcAccountId(v string) *CreateServiceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *CreateServiceHeaders) SetXFcDate(v string) *CreateServiceHeaders {
	s.XFcDate = &v
	return s
}

func (s *CreateServiceHeaders) SetXFcTraceId(v string) *CreateServiceHeaders {
	s.XFcTraceId = &v
	return s
}

type CreateServiceRequest struct {
	// The description of the service.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Specifies whether to allow functions to access the Internet. Valid values:
	//
	// *   **true**: allows functions to access the Internet. This is the default value.
	// *   **false**: does not allow functions to access the Internet.
	InternetAccess *bool `json:"internetAccess,omitempty" xml:"internetAccess,omitempty"`
	// The log configuration. Function Compute writes function execution logs to the specified Logstore.
	LogConfig *LogConfig `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	// The configuration of the Apsara File Storage NAS (NAS) file system. The configurations allow functions in the specified service to access the NAS file system.
	NasConfig *NASConfig `json:"nasConfig,omitempty" xml:"nasConfig,omitempty"`
	// The OSS mount configurations.
	OssMountConfig *OSSMountConfig `json:"ossMountConfig,omitempty" xml:"ossMountConfig,omitempty"`
	// The RAM role that is used to grant required permissions to Function Compute. The RAM role is used in the following scenarios:
	//
	// *   Sends function execution logs to your Logstore.
	// *   Generates a token for a function to access other cloud resources during function execution.
	Role *string `json:"role,omitempty" xml:"role,omitempty"`
	// The name of the service. The name can contain only letters, digits, hyphens (-), and underscores (\_). It cannot start with a digit or hyphen (-). It must be 1 to 128 characters in length.
	ServiceName *string `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
	// The configuration of Tracing Analysis. After Function Compute is integrated with Tracing Analysis, you can record the duration of a request in Function Compute, view the cold start time of a function, and record the execution duration of a function. For more information, see [Tracing Analysis](~~189804~~).
	TracingConfig *TracingConfig `json:"tracingConfig,omitempty" xml:"tracingConfig,omitempty"`
	// The VPC configurations. The configurations allow functions in the specified service to access the specified VPC.
	VpcConfig *VPCConfig `json:"vpcConfig,omitempty" xml:"vpcConfig,omitempty"`
}

func (s CreateServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateServiceRequest) GoString() string {
	return s.String()
}

func (s *CreateServiceRequest) SetDescription(v string) *CreateServiceRequest {
	s.Description = &v
	return s
}

func (s *CreateServiceRequest) SetInternetAccess(v bool) *CreateServiceRequest {
	s.InternetAccess = &v
	return s
}

func (s *CreateServiceRequest) SetLogConfig(v *LogConfig) *CreateServiceRequest {
	s.LogConfig = v
	return s
}

func (s *CreateServiceRequest) SetNasConfig(v *NASConfig) *CreateServiceRequest {
	s.NasConfig = v
	return s
}

func (s *CreateServiceRequest) SetOssMountConfig(v *OSSMountConfig) *CreateServiceRequest {
	s.OssMountConfig = v
	return s
}

func (s *CreateServiceRequest) SetRole(v string) *CreateServiceRequest {
	s.Role = &v
	return s
}

func (s *CreateServiceRequest) SetServiceName(v string) *CreateServiceRequest {
	s.ServiceName = &v
	return s
}

func (s *CreateServiceRequest) SetTracingConfig(v *TracingConfig) *CreateServiceRequest {
	s.TracingConfig = v
	return s
}

func (s *CreateServiceRequest) SetVpcConfig(v *VPCConfig) *CreateServiceRequest {
	s.VpcConfig = v
	return s
}

type CreateServiceResponseBody struct {
	// The time when the service was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the service.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Specifies whether to allow functions to access the Internet. Valid values:
	//
	// *   **true**: allows functions in the specified service to access the Internet.
	// *   **false**: does not allow functions in the specified service to access the Internet.
	InternetAccess *bool `json:"internetAccess,omitempty" xml:"internetAccess,omitempty"`
	// The time when the service was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The log configuration. Function Compute writes function execution logs to the specified Logstore.
	LogConfig *LogConfig `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	// The configuration of the NAS file system. The configurations allow functions in the specified service to access the NAS file system.
	NasConfig *NASConfig `json:"nasConfig,omitempty" xml:"nasConfig,omitempty"`
	// The OSS mount configurations.
	OssMountConfig *OSSMountConfig `json:"ossMountConfig,omitempty" xml:"ossMountConfig,omitempty"`
	// The RAM role that is used to grant required permissions to Function Compute. The RAM role is used in the following scenarios:
	//
	// *   Sends function execution logs to your Logstore.
	// *   Generates a token for a function to access other cloud resources during function execution.
	Role *string `json:"role,omitempty" xml:"role,omitempty"`
	// The unique ID generated by the system for the service.
	ServiceId *string `json:"serviceId,omitempty" xml:"serviceId,omitempty"`
	// The name of the service.
	ServiceName *string `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
	// The configuration of Tracing Analysis. After Function Compute is integrated with Tracing Analysis, you can record the duration of a request in Function Compute, view the cold start time of a function, and record the execution duration of a function. For more information, see [Tracing Analysis](~~189804~~).
	TracingConfig *TracingConfig `json:"tracingConfig,omitempty" xml:"tracingConfig,omitempty"`
	// The VPC configurations. The configurations allow functions in the specified service to access the specified VPC.
	VpcConfig *VPCConfig `json:"vpcConfig,omitempty" xml:"vpcConfig,omitempty"`
}

func (s CreateServiceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateServiceResponseBody) GoString() string {
	return s.String()
}

func (s *CreateServiceResponseBody) SetCreatedTime(v string) *CreateServiceResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *CreateServiceResponseBody) SetDescription(v string) *CreateServiceResponseBody {
	s.Description = &v
	return s
}

func (s *CreateServiceResponseBody) SetInternetAccess(v bool) *CreateServiceResponseBody {
	s.InternetAccess = &v
	return s
}

func (s *CreateServiceResponseBody) SetLastModifiedTime(v string) *CreateServiceResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *CreateServiceResponseBody) SetLogConfig(v *LogConfig) *CreateServiceResponseBody {
	s.LogConfig = v
	return s
}

func (s *CreateServiceResponseBody) SetNasConfig(v *NASConfig) *CreateServiceResponseBody {
	s.NasConfig = v
	return s
}

func (s *CreateServiceResponseBody) SetOssMountConfig(v *OSSMountConfig) *CreateServiceResponseBody {
	s.OssMountConfig = v
	return s
}

func (s *CreateServiceResponseBody) SetRole(v string) *CreateServiceResponseBody {
	s.Role = &v
	return s
}

func (s *CreateServiceResponseBody) SetServiceId(v string) *CreateServiceResponseBody {
	s.ServiceId = &v
	return s
}

func (s *CreateServiceResponseBody) SetServiceName(v string) *CreateServiceResponseBody {
	s.ServiceName = &v
	return s
}

func (s *CreateServiceResponseBody) SetTracingConfig(v *TracingConfig) *CreateServiceResponseBody {
	s.TracingConfig = v
	return s
}

func (s *CreateServiceResponseBody) SetVpcConfig(v *VPCConfig) *CreateServiceResponseBody {
	s.VpcConfig = v
	return s
}

type CreateServiceResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateServiceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateServiceResponse) GoString() string {
	return s.String()
}

func (s *CreateServiceResponse) SetHeaders(v map[string]*string) *CreateServiceResponse {
	s.Headers = v
	return s
}

func (s *CreateServiceResponse) SetStatusCode(v int32) *CreateServiceResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateServiceResponse) SetBody(v *CreateServiceResponseBody) *CreateServiceResponse {
	s.Body = v
	return s
}

type CreateTriggerHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the request is initiated on the client. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s CreateTriggerHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateTriggerHeaders) GoString() string {
	return s.String()
}

func (s *CreateTriggerHeaders) SetCommonHeaders(v map[string]*string) *CreateTriggerHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateTriggerHeaders) SetXFcAccountId(v string) *CreateTriggerHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *CreateTriggerHeaders) SetXFcDate(v string) *CreateTriggerHeaders {
	s.XFcDate = &v
	return s
}

func (s *CreateTriggerHeaders) SetXFcTraceId(v string) *CreateTriggerHeaders {
	s.XFcTraceId = &v
	return s
}

type CreateTriggerRequest struct {
	// The description of the trigger.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The role that is used by the event source such as OSS to invoke the function. For more information, see [Overview](~~53102~~).
	InvocationRole *string `json:"invocationRole,omitempty" xml:"invocationRole,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The Alibaba Cloud Resource Name (ARN) of the event source for the trigger.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
	// The configurations of the trigger. The configurations vary based on the trigger type. For more information about the format, see the following topics:
	//
	// * OSS trigger: [OSSTriggerConfig](~~struct:OSSTriggerConfig~~).
	// * Log Service trigger: [LogTriggerConfig](~~struct:LogTriggerConfig~~).
	// * Time trigger: [TimeTriggerConfig](~~struct:LogTriggerConfig~~).
	// * HTTP trigger: [HTTPTriggerConfig](~~struct:HTTPTriggerConfig~~).
	// * Tablestore trigger: Specify the **SourceArn** parameter and leave this parameter empty.
	// * Alibaba Cloud CDN event trigger: [CDNEventsTriggerConfig](~~struct:CDNEventsTriggerConfig~~).
	// * MNS topic trigger: [MnsTopicTriggerConfig](~~struct:MnsTopicTriggerConfig~~).
	TriggerConfig *string `json:"triggerConfig,omitempty" xml:"triggerConfig,omitempty"`
	// The name of the trigger. The name contains only letters, digits, hyphens (-), and underscores (\_). The name must be 1 to 128 characters in length and cannot start with a digit or hyphen (-).
	TriggerName *string `json:"triggerName,omitempty" xml:"triggerName,omitempty"`
	// The type of the trigger. Valid values:
	//
	// *   **oss**: OSS event trigger. For more information, see [Overview](~~62922~~).
	// *   **log**: Log Service trigger. For more information, see [Overview](~~84386~~).
	// *   **timer**: time trigger. For more information, see [Overview](~~68172~~).
	// *   **http**: HTTP trigger. For more information, see [Overview](~~71229~~).
	// *   **tablestore**: Tablestore trigger. For more information, see [Overview](~~100092~~).
	// *   **cdn_events**: CDN event trigger. For more information, see [Overview](~~73333~~).
	// *   **mns_topic**: MNS topic trigger. For more information, see [Overview](~~97032~~).
	TriggerType *string `json:"triggerType,omitempty" xml:"triggerType,omitempty"`
}

func (s CreateTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateTriggerRequest) GoString() string {
	return s.String()
}

func (s *CreateTriggerRequest) SetDescription(v string) *CreateTriggerRequest {
	s.Description = &v
	return s
}

func (s *CreateTriggerRequest) SetInvocationRole(v string) *CreateTriggerRequest {
	s.InvocationRole = &v
	return s
}

func (s *CreateTriggerRequest) SetQualifier(v string) *CreateTriggerRequest {
	s.Qualifier = &v
	return s
}

func (s *CreateTriggerRequest) SetSourceArn(v string) *CreateTriggerRequest {
	s.SourceArn = &v
	return s
}

func (s *CreateTriggerRequest) SetTriggerConfig(v string) *CreateTriggerRequest {
	s.TriggerConfig = &v
	return s
}

func (s *CreateTriggerRequest) SetTriggerName(v string) *CreateTriggerRequest {
	s.TriggerName = &v
	return s
}

func (s *CreateTriggerRequest) SetTriggerType(v string) *CreateTriggerRequest {
	s.TriggerType = &v
	return s
}

type CreateTriggerResponseBody struct {
	// The time when the trigger was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the trigger.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The domain name used to invoke the function by using HTTP. You can add this domain name as the prefix to the endpoint of Function Compute. This way, you can invoke the function that corresponds to the trigger by using HTTP. For example, `{domainName}.cn-shanghai.fc.aliyuncs.com`.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The ARN of the RAM role that is used by the event source to invoke the function.
	InvocationRole *string `json:"invocationRole,omitempty" xml:"invocationRole,omitempty"`
	// The time when the trigger was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The version of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The ARN of the event source.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
	// The configurations of the trigger. The configurations vary based on the trigger type.
	TriggerConfig *string `json:"triggerConfig,omitempty" xml:"triggerConfig,omitempty"`
	// The unique ID of the trigger.
	TriggerId *string `json:"triggerId,omitempty" xml:"triggerId,omitempty"`
	// The name of the trigger. The name contains only letters, digits, hyphens (-), and underscores (\_). The name must be 1 to 128 characters in length and cannot start with a digit or hyphen (-).
	TriggerName *string `json:"triggerName,omitempty" xml:"triggerName,omitempty"`
	// The trigger type, such as **oss**, **log**, **tablestore**, **timer**, **http**, **cdn_events**, and **mns_topic**.
	TriggerType *string `json:"triggerType,omitempty" xml:"triggerType,omitempty"`
	// The public domain address. You can access HTTP triggers over the Internet by using HTTP or HTTPS.
	UrlInternet *string `json:"urlInternet,omitempty" xml:"urlInternet,omitempty"`
	// The private endpoint. In a VPC, you can access HTTP triggers by using HTTP or HTTPS.
	UrlIntranet *string `json:"urlIntranet,omitempty" xml:"urlIntranet,omitempty"`
}

func (s CreateTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *CreateTriggerResponseBody) SetCreatedTime(v string) *CreateTriggerResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *CreateTriggerResponseBody) SetDescription(v string) *CreateTriggerResponseBody {
	s.Description = &v
	return s
}

func (s *CreateTriggerResponseBody) SetDomainName(v string) *CreateTriggerResponseBody {
	s.DomainName = &v
	return s
}

func (s *CreateTriggerResponseBody) SetInvocationRole(v string) *CreateTriggerResponseBody {
	s.InvocationRole = &v
	return s
}

func (s *CreateTriggerResponseBody) SetLastModifiedTime(v string) *CreateTriggerResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *CreateTriggerResponseBody) SetQualifier(v string) *CreateTriggerResponseBody {
	s.Qualifier = &v
	return s
}

func (s *CreateTriggerResponseBody) SetSourceArn(v string) *CreateTriggerResponseBody {
	s.SourceArn = &v
	return s
}

func (s *CreateTriggerResponseBody) SetTriggerConfig(v string) *CreateTriggerResponseBody {
	s.TriggerConfig = &v
	return s
}

func (s *CreateTriggerResponseBody) SetTriggerId(v string) *CreateTriggerResponseBody {
	s.TriggerId = &v
	return s
}

func (s *CreateTriggerResponseBody) SetTriggerName(v string) *CreateTriggerResponseBody {
	s.TriggerName = &v
	return s
}

func (s *CreateTriggerResponseBody) SetTriggerType(v string) *CreateTriggerResponseBody {
	s.TriggerType = &v
	return s
}

func (s *CreateTriggerResponseBody) SetUrlInternet(v string) *CreateTriggerResponseBody {
	s.UrlInternet = &v
	return s
}

func (s *CreateTriggerResponseBody) SetUrlIntranet(v string) *CreateTriggerResponseBody {
	s.UrlIntranet = &v
	return s
}

type CreateTriggerResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateTriggerResponse) GoString() string {
	return s.String()
}

func (s *CreateTriggerResponse) SetHeaders(v map[string]*string) *CreateTriggerResponse {
	s.Headers = v
	return s
}

func (s *CreateTriggerResponse) SetStatusCode(v int32) *CreateTriggerResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateTriggerResponse) SetBody(v *CreateTriggerResponseBody) *CreateTriggerResponse {
	s.Body = v
	return s
}

type CreateVpcBindingHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s CreateVpcBindingHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateVpcBindingHeaders) GoString() string {
	return s.String()
}

func (s *CreateVpcBindingHeaders) SetCommonHeaders(v map[string]*string) *CreateVpcBindingHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateVpcBindingHeaders) SetXFcAccountId(v string) *CreateVpcBindingHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *CreateVpcBindingHeaders) SetXFcDate(v string) *CreateVpcBindingHeaders {
	s.XFcDate = &v
	return s
}

func (s *CreateVpcBindingHeaders) SetXFcTraceId(v string) *CreateVpcBindingHeaders {
	s.XFcTraceId = &v
	return s
}

type CreateVpcBindingRequest struct {
	// The ID of the VPC to be bound.
	VpcId *string `json:"vpcId,omitempty" xml:"vpcId,omitempty"`
}

func (s CreateVpcBindingRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateVpcBindingRequest) GoString() string {
	return s.String()
}

func (s *CreateVpcBindingRequest) SetVpcId(v string) *CreateVpcBindingRequest {
	s.VpcId = &v
	return s
}

type CreateVpcBindingResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s CreateVpcBindingResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateVpcBindingResponse) GoString() string {
	return s.String()
}

func (s *CreateVpcBindingResponse) SetHeaders(v map[string]*string) *CreateVpcBindingResponse {
	s.Headers = v
	return s
}

func (s *CreateVpcBindingResponse) SetStatusCode(v int32) *CreateVpcBindingResponse {
	s.StatusCode = &v
	return s
}

type DeleteAliasHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// If the ETag specified in the request matches the ETag value of the object, OSS transmits the object and returns 200 OK. If the ETag specified in the request does not match the ETag value of the object, OSS returns 412 Precondition Failed.
	// The ETag value of a resource is used to check whether the resource has changed. You can check data integrity by using the ETag value.
	// Default value: null
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteAliasHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteAliasHeaders) GoString() string {
	return s.String()
}

func (s *DeleteAliasHeaders) SetCommonHeaders(v map[string]*string) *DeleteAliasHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteAliasHeaders) SetIfMatch(v string) *DeleteAliasHeaders {
	s.IfMatch = &v
	return s
}

func (s *DeleteAliasHeaders) SetXFcAccountId(v string) *DeleteAliasHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteAliasHeaders) SetXFcDate(v string) *DeleteAliasHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteAliasHeaders) SetXFcTraceId(v string) *DeleteAliasHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteAliasResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteAliasResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAliasResponse) GoString() string {
	return s.String()
}

func (s *DeleteAliasResponse) SetHeaders(v map[string]*string) *DeleteAliasResponse {
	s.Headers = v
	return s
}

func (s *DeleteAliasResponse) SetStatusCode(v int32) *DeleteAliasResponse {
	s.StatusCode = &v
	return s
}

type DeleteCustomDomainHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteCustomDomainHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteCustomDomainHeaders) GoString() string {
	return s.String()
}

func (s *DeleteCustomDomainHeaders) SetCommonHeaders(v map[string]*string) *DeleteCustomDomainHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteCustomDomainHeaders) SetXFcAccountId(v string) *DeleteCustomDomainHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteCustomDomainHeaders) SetXFcDate(v string) *DeleteCustomDomainHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteCustomDomainHeaders) SetXFcTraceId(v string) *DeleteCustomDomainHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteCustomDomainResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteCustomDomainResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteCustomDomainResponse) GoString() string {
	return s.String()
}

func (s *DeleteCustomDomainResponse) SetHeaders(v map[string]*string) *DeleteCustomDomainResponse {
	s.Headers = v
	return s
}

func (s *DeleteCustomDomainResponse) SetStatusCode(v int32) *DeleteCustomDomainResponse {
	s.StatusCode = &v
	return s
}

type DeleteFunctionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ETag value of the resource. This value is used to ensure that the modified resource is consistent with the resource to be modified. The ETag value is returned in the responses of the CREATE, GET, and UPDATE operations.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API. The value is the same as that of the requestId parameter in the response.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteFunctionHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionHeaders) GoString() string {
	return s.String()
}

func (s *DeleteFunctionHeaders) SetCommonHeaders(v map[string]*string) *DeleteFunctionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteFunctionHeaders) SetIfMatch(v string) *DeleteFunctionHeaders {
	s.IfMatch = &v
	return s
}

func (s *DeleteFunctionHeaders) SetXFcAccountId(v string) *DeleteFunctionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteFunctionHeaders) SetXFcDate(v string) *DeleteFunctionHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteFunctionHeaders) SetXFcTraceId(v string) *DeleteFunctionHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteFunctionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteFunctionResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionResponse) GoString() string {
	return s.String()
}

func (s *DeleteFunctionResponse) SetHeaders(v map[string]*string) *DeleteFunctionResponse {
	s.Headers = v
	return s
}

func (s *DeleteFunctionResponse) SetStatusCode(v int32) *DeleteFunctionResponse {
	s.StatusCode = &v
	return s
}

type DeleteFunctionAsyncInvokeConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteFunctionAsyncInvokeConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionAsyncInvokeConfigHeaders) GoString() string {
	return s.String()
}

func (s *DeleteFunctionAsyncInvokeConfigHeaders) SetCommonHeaders(v map[string]*string) *DeleteFunctionAsyncInvokeConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteFunctionAsyncInvokeConfigHeaders) SetXFcAccountId(v string) *DeleteFunctionAsyncInvokeConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteFunctionAsyncInvokeConfigHeaders) SetXFcDate(v string) *DeleteFunctionAsyncInvokeConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteFunctionAsyncInvokeConfigHeaders) SetXFcTraceId(v string) *DeleteFunctionAsyncInvokeConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteFunctionAsyncInvokeConfigRequest struct {
	// The qualifier.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s DeleteFunctionAsyncInvokeConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionAsyncInvokeConfigRequest) GoString() string {
	return s.String()
}

func (s *DeleteFunctionAsyncInvokeConfigRequest) SetQualifier(v string) *DeleteFunctionAsyncInvokeConfigRequest {
	s.Qualifier = &v
	return s
}

type DeleteFunctionAsyncInvokeConfigResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteFunctionAsyncInvokeConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionAsyncInvokeConfigResponse) GoString() string {
	return s.String()
}

func (s *DeleteFunctionAsyncInvokeConfigResponse) SetHeaders(v map[string]*string) *DeleteFunctionAsyncInvokeConfigResponse {
	s.Headers = v
	return s
}

func (s *DeleteFunctionAsyncInvokeConfigResponse) SetStatusCode(v int32) *DeleteFunctionAsyncInvokeConfigResponse {
	s.StatusCode = &v
	return s
}

type DeleteFunctionOnDemandConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// If the ETag specified in the request matches the ETag value of the OndemandConfig, FC returns 200 OK. If the ETag specified in the request does not match the ETag value of the object, FC returns 412 Precondition Failed.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The start time when the function is invoked. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API, which is also the unique ID of the request.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteFunctionOnDemandConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionOnDemandConfigHeaders) GoString() string {
	return s.String()
}

func (s *DeleteFunctionOnDemandConfigHeaders) SetCommonHeaders(v map[string]*string) *DeleteFunctionOnDemandConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteFunctionOnDemandConfigHeaders) SetIfMatch(v string) *DeleteFunctionOnDemandConfigHeaders {
	s.IfMatch = &v
	return s
}

func (s *DeleteFunctionOnDemandConfigHeaders) SetXFcAccountId(v string) *DeleteFunctionOnDemandConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteFunctionOnDemandConfigHeaders) SetXFcDate(v string) *DeleteFunctionOnDemandConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteFunctionOnDemandConfigHeaders) SetXFcTraceId(v string) *DeleteFunctionOnDemandConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteFunctionOnDemandConfigRequest struct {
	// The alias of the service or LATEST.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s DeleteFunctionOnDemandConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionOnDemandConfigRequest) GoString() string {
	return s.String()
}

func (s *DeleteFunctionOnDemandConfigRequest) SetQualifier(v string) *DeleteFunctionOnDemandConfigRequest {
	s.Qualifier = &v
	return s
}

type DeleteFunctionOnDemandConfigResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteFunctionOnDemandConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteFunctionOnDemandConfigResponse) GoString() string {
	return s.String()
}

func (s *DeleteFunctionOnDemandConfigResponse) SetHeaders(v map[string]*string) *DeleteFunctionOnDemandConfigResponse {
	s.Headers = v
	return s
}

func (s *DeleteFunctionOnDemandConfigResponse) SetStatusCode(v int32) *DeleteFunctionOnDemandConfigResponse {
	s.StatusCode = &v
	return s
}

type DeleteLayerVersionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteLayerVersionHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteLayerVersionHeaders) GoString() string {
	return s.String()
}

func (s *DeleteLayerVersionHeaders) SetCommonHeaders(v map[string]*string) *DeleteLayerVersionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteLayerVersionHeaders) SetXFcAccountId(v string) *DeleteLayerVersionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteLayerVersionHeaders) SetXFcDate(v string) *DeleteLayerVersionHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteLayerVersionHeaders) SetXFcTraceId(v string) *DeleteLayerVersionHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteLayerVersionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteLayerVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteLayerVersionResponse) GoString() string {
	return s.String()
}

func (s *DeleteLayerVersionResponse) SetHeaders(v map[string]*string) *DeleteLayerVersionResponse {
	s.Headers = v
	return s
}

func (s *DeleteLayerVersionResponse) SetStatusCode(v int32) *DeleteLayerVersionResponse {
	s.StatusCode = &v
	return s
}

type DeleteProvisionConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XFcAccountId  *string            `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	XFcDate       *string            `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	XFcTraceId    *string            `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteProvisionConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteProvisionConfigHeaders) GoString() string {
	return s.String()
}

func (s *DeleteProvisionConfigHeaders) SetCommonHeaders(v map[string]*string) *DeleteProvisionConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteProvisionConfigHeaders) SetXFcAccountId(v string) *DeleteProvisionConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteProvisionConfigHeaders) SetXFcDate(v string) *DeleteProvisionConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteProvisionConfigHeaders) SetXFcTraceId(v string) *DeleteProvisionConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteProvisionConfigRequest struct {
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s DeleteProvisionConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteProvisionConfigRequest) GoString() string {
	return s.String()
}

func (s *DeleteProvisionConfigRequest) SetQualifier(v string) *DeleteProvisionConfigRequest {
	s.Qualifier = &v
	return s
}

type DeleteProvisionConfigResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteProvisionConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteProvisionConfigResponse) GoString() string {
	return s.String()
}

func (s *DeleteProvisionConfigResponse) SetHeaders(v map[string]*string) *DeleteProvisionConfigResponse {
	s.Headers = v
	return s
}

func (s *DeleteProvisionConfigResponse) SetStatusCode(v int32) *DeleteProvisionConfigResponse {
	s.StatusCode = &v
	return s
}

type DeleteServiceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ETag value of the service. This value is used to ensure that the modified service is consistent with the service to be modified. The ETag value is returned in the responses of the [CreateService](~~175256~~), [UpdateService](~~188167~~), and [GetService](~~189225~~) operations.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteServiceHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteServiceHeaders) GoString() string {
	return s.String()
}

func (s *DeleteServiceHeaders) SetCommonHeaders(v map[string]*string) *DeleteServiceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteServiceHeaders) SetIfMatch(v string) *DeleteServiceHeaders {
	s.IfMatch = &v
	return s
}

func (s *DeleteServiceHeaders) SetXFcAccountId(v string) *DeleteServiceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteServiceHeaders) SetXFcDate(v string) *DeleteServiceHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteServiceHeaders) SetXFcTraceId(v string) *DeleteServiceHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteServiceResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteServiceResponse) GoString() string {
	return s.String()
}

func (s *DeleteServiceResponse) SetHeaders(v map[string]*string) *DeleteServiceResponse {
	s.Headers = v
	return s
}

func (s *DeleteServiceResponse) SetStatusCode(v int32) *DeleteServiceResponse {
	s.StatusCode = &v
	return s
}

type DeleteServiceVersionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteServiceVersionHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteServiceVersionHeaders) GoString() string {
	return s.String()
}

func (s *DeleteServiceVersionHeaders) SetCommonHeaders(v map[string]*string) *DeleteServiceVersionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteServiceVersionHeaders) SetXFcAccountId(v string) *DeleteServiceVersionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteServiceVersionHeaders) SetXFcDate(v string) *DeleteServiceVersionHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteServiceVersionHeaders) SetXFcTraceId(v string) *DeleteServiceVersionHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteServiceVersionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteServiceVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteServiceVersionResponse) GoString() string {
	return s.String()
}

func (s *DeleteServiceVersionResponse) SetHeaders(v map[string]*string) *DeleteServiceVersionResponse {
	s.Headers = v
	return s
}

func (s *DeleteServiceVersionResponse) SetStatusCode(v int32) *DeleteServiceVersionResponse {
	s.StatusCode = &v
	return s
}

type DeleteTriggerHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// This parameter is used to ensure that the modified resource is consistent with the resource to be modified. You can obtain the parameter value from the responses of [CreateTrigger](~~415729~~), [GetTrigger](~~415732~~), and [UpdateTrigger](~~415731~~) operations.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the request is initiated on the client. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteTriggerHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteTriggerHeaders) GoString() string {
	return s.String()
}

func (s *DeleteTriggerHeaders) SetCommonHeaders(v map[string]*string) *DeleteTriggerHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteTriggerHeaders) SetIfMatch(v string) *DeleteTriggerHeaders {
	s.IfMatch = &v
	return s
}

func (s *DeleteTriggerHeaders) SetXFcAccountId(v string) *DeleteTriggerHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteTriggerHeaders) SetXFcDate(v string) *DeleteTriggerHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteTriggerHeaders) SetXFcTraceId(v string) *DeleteTriggerHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteTriggerResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteTriggerResponse) GoString() string {
	return s.String()
}

func (s *DeleteTriggerResponse) SetHeaders(v map[string]*string) *DeleteTriggerResponse {
	s.Headers = v
	return s
}

func (s *DeleteTriggerResponse) SetStatusCode(v int32) *DeleteTriggerResponse {
	s.StatusCode = &v
	return s
}

type DeleteVpcBindingHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeleteVpcBindingHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteVpcBindingHeaders) GoString() string {
	return s.String()
}

func (s *DeleteVpcBindingHeaders) SetCommonHeaders(v map[string]*string) *DeleteVpcBindingHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteVpcBindingHeaders) SetXFcAccountId(v string) *DeleteVpcBindingHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeleteVpcBindingHeaders) SetXFcDate(v string) *DeleteVpcBindingHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeleteVpcBindingHeaders) SetXFcTraceId(v string) *DeleteVpcBindingHeaders {
	s.XFcTraceId = &v
	return s
}

type DeleteVpcBindingResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteVpcBindingResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteVpcBindingResponse) GoString() string {
	return s.String()
}

func (s *DeleteVpcBindingResponse) SetHeaders(v map[string]*string) *DeleteVpcBindingResponse {
	s.Headers = v
	return s
}

func (s *DeleteVpcBindingResponse) SetStatusCode(v int32) *DeleteVpcBindingResponse {
	s.StatusCode = &v
	return s
}

type DeregisterEventSourceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s DeregisterEventSourceHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeregisterEventSourceHeaders) GoString() string {
	return s.String()
}

func (s *DeregisterEventSourceHeaders) SetCommonHeaders(v map[string]*string) *DeregisterEventSourceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeregisterEventSourceHeaders) SetXFcAccountId(v string) *DeregisterEventSourceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *DeregisterEventSourceHeaders) SetXFcDate(v string) *DeregisterEventSourceHeaders {
	s.XFcDate = &v
	return s
}

func (s *DeregisterEventSourceHeaders) SetXFcTraceId(v string) *DeregisterEventSourceHeaders {
	s.XFcTraceId = &v
	return s
}

type DeregisterEventSourceRequest struct {
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s DeregisterEventSourceRequest) String() string {
	return tea.Prettify(s)
}

func (s DeregisterEventSourceRequest) GoString() string {
	return s.String()
}

func (s *DeregisterEventSourceRequest) SetQualifier(v string) *DeregisterEventSourceRequest {
	s.Qualifier = &v
	return s
}

type DeregisterEventSourceResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeregisterEventSourceResponse) String() string {
	return tea.Prettify(s)
}

func (s DeregisterEventSourceResponse) GoString() string {
	return s.String()
}

func (s *DeregisterEventSourceResponse) SetHeaders(v map[string]*string) *DeregisterEventSourceResponse {
	s.Headers = v
	return s
}

func (s *DeregisterEventSourceResponse) SetStatusCode(v int32) *DeregisterEventSourceResponse {
	s.StatusCode = &v
	return s
}

type GetAccountSettingsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetAccountSettingsHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAccountSettingsHeaders) GoString() string {
	return s.String()
}

func (s *GetAccountSettingsHeaders) SetCommonHeaders(v map[string]*string) *GetAccountSettingsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAccountSettingsHeaders) SetXFcAccountId(v string) *GetAccountSettingsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetAccountSettingsHeaders) SetXFcDate(v string) *GetAccountSettingsHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetAccountSettingsHeaders) SetXFcTraceId(v string) *GetAccountSettingsHeaders {
	s.XFcTraceId = &v
	return s
}

type GetAccountSettingsResponseBody struct {
	// The list of zones.
	AvailableAZs []*string `json:"availableAZs,omitempty" xml:"availableAZs,omitempty" type:"Repeated"`
	// The default RAM role.
	DefaultRole *string `json:"defaultRole,omitempty" xml:"defaultRole,omitempty"`
}

func (s GetAccountSettingsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAccountSettingsResponseBody) GoString() string {
	return s.String()
}

func (s *GetAccountSettingsResponseBody) SetAvailableAZs(v []*string) *GetAccountSettingsResponseBody {
	s.AvailableAZs = v
	return s
}

func (s *GetAccountSettingsResponseBody) SetDefaultRole(v string) *GetAccountSettingsResponseBody {
	s.DefaultRole = &v
	return s
}

type GetAccountSettingsResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetAccountSettingsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetAccountSettingsResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccountSettingsResponse) GoString() string {
	return s.String()
}

func (s *GetAccountSettingsResponse) SetHeaders(v map[string]*string) *GetAccountSettingsResponse {
	s.Headers = v
	return s
}

func (s *GetAccountSettingsResponse) SetStatusCode(v int32) *GetAccountSettingsResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAccountSettingsResponse) SetBody(v *GetAccountSettingsResponseBody) *GetAccountSettingsResponse {
	s.Body = v
	return s
}

type GetAliasHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The start time when the function is invoked. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetAliasHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAliasHeaders) GoString() string {
	return s.String()
}

func (s *GetAliasHeaders) SetCommonHeaders(v map[string]*string) *GetAliasHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAliasHeaders) SetXFcAccountId(v string) *GetAliasHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetAliasHeaders) SetXFcDate(v string) *GetAliasHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetAliasHeaders) SetXFcTraceId(v string) *GetAliasHeaders {
	s.XFcTraceId = &v
	return s
}

type GetAliasResponseBody struct {
	// The canary release version to which the alias points and the weight of the canary release version.
	//
	// - The canary release version takes effect only when the function is invoked.
	// - The value consists of a version number and the corresponding weight. For example, 2:0.05 indicates that when a function is invoked, Version 2 is the canary release version, 5% of the traffic is distributed to the canary release version, and 95% of the traffic is distributed to the major version.
	AdditionalVersionWeight map[string]*float32 `json:"additionalVersionWeight,omitempty" xml:"additionalVersionWeight,omitempty"`
	// The name of the alias.
	AliasName *string `json:"aliasName,omitempty" xml:"aliasName,omitempty"`
	// The time when the alias was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the alias.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The time when the alias was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The canary release mode. Valid values:
	//
	// - **Random**: random canary release. This is the default value.
	// - **Content**: rule-based canary release.
	ResolvePolicy *string `json:"resolvePolicy,omitempty" xml:"resolvePolicy,omitempty"`
	// Canary release rule. The traffic that meets the conditions of the canary release rule is diverted to the canary release instances.
	RoutePolicy *RoutePolicy `json:"routePolicy,omitempty" xml:"routePolicy,omitempty"`
	// The version to which the alias points.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s GetAliasResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAliasResponseBody) GoString() string {
	return s.String()
}

func (s *GetAliasResponseBody) SetAdditionalVersionWeight(v map[string]*float32) *GetAliasResponseBody {
	s.AdditionalVersionWeight = v
	return s
}

func (s *GetAliasResponseBody) SetAliasName(v string) *GetAliasResponseBody {
	s.AliasName = &v
	return s
}

func (s *GetAliasResponseBody) SetCreatedTime(v string) *GetAliasResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *GetAliasResponseBody) SetDescription(v string) *GetAliasResponseBody {
	s.Description = &v
	return s
}

func (s *GetAliasResponseBody) SetLastModifiedTime(v string) *GetAliasResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *GetAliasResponseBody) SetResolvePolicy(v string) *GetAliasResponseBody {
	s.ResolvePolicy = &v
	return s
}

func (s *GetAliasResponseBody) SetRoutePolicy(v *RoutePolicy) *GetAliasResponseBody {
	s.RoutePolicy = v
	return s
}

func (s *GetAliasResponseBody) SetVersionId(v string) *GetAliasResponseBody {
	s.VersionId = &v
	return s
}

type GetAliasResponse struct {
	Headers    map[string]*string    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetAliasResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetAliasResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAliasResponse) GoString() string {
	return s.String()
}

func (s *GetAliasResponse) SetHeaders(v map[string]*string) *GetAliasResponse {
	s.Headers = v
	return s
}

func (s *GetAliasResponse) SetStatusCode(v int32) *GetAliasResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAliasResponse) SetBody(v *GetAliasResponseBody) *GetAliasResponse {
	s.Body = v
	return s
}

type GetCustomDomainHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the operation is called. The format is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetCustomDomainHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetCustomDomainHeaders) GoString() string {
	return s.String()
}

func (s *GetCustomDomainHeaders) SetCommonHeaders(v map[string]*string) *GetCustomDomainHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetCustomDomainHeaders) SetXFcAccountId(v string) *GetCustomDomainHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetCustomDomainHeaders) SetXFcDate(v string) *GetCustomDomainHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetCustomDomainHeaders) SetXFcTraceId(v string) *GetCustomDomainHeaders {
	s.XFcTraceId = &v
	return s
}

type GetCustomDomainResponseBody struct {
	// The ID of your Alibaba Cloud account.
	AccountId *string `json:"accountId,omitempty" xml:"accountId,omitempty"`
	// The version of the API.
	ApiVersion *string `json:"apiVersion,omitempty" xml:"apiVersion,omitempty"`
	// The configurations of the HTTPS certificate.
	CertConfig *CertConfig `json:"certConfig,omitempty" xml:"certConfig,omitempty"`
	// The time when the custom domain name was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The domain name.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The time when the domain name was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The protocol types supported by the domain name. Valid values:
	//
	// *   **HTTP**: Only HTTP is supported.
	// *   **HTTPS**: Only HTTPS is supported.
	// *   **HTTP,HTTPS**: HTTP and HTTPS are supported.
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// The route table that maps the paths to functions when the functions are invoked by using the custom domain name.
	RouteConfig *RouteConfig `json:"routeConfig,omitempty" xml:"routeConfig,omitempty"`
	// The Transport Layer Security (TLS) configuration.
	TlsConfig *TLSConfig `json:"tlsConfig,omitempty" xml:"tlsConfig,omitempty"`
	// The Web Application Firewall (WAF) configuration.
	WafConfig *WAFConfig `json:"wafConfig,omitempty" xml:"wafConfig,omitempty"`
}

func (s GetCustomDomainResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetCustomDomainResponseBody) GoString() string {
	return s.String()
}

func (s *GetCustomDomainResponseBody) SetAccountId(v string) *GetCustomDomainResponseBody {
	s.AccountId = &v
	return s
}

func (s *GetCustomDomainResponseBody) SetApiVersion(v string) *GetCustomDomainResponseBody {
	s.ApiVersion = &v
	return s
}

func (s *GetCustomDomainResponseBody) SetCertConfig(v *CertConfig) *GetCustomDomainResponseBody {
	s.CertConfig = v
	return s
}

func (s *GetCustomDomainResponseBody) SetCreatedTime(v string) *GetCustomDomainResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *GetCustomDomainResponseBody) SetDomainName(v string) *GetCustomDomainResponseBody {
	s.DomainName = &v
	return s
}

func (s *GetCustomDomainResponseBody) SetLastModifiedTime(v string) *GetCustomDomainResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *GetCustomDomainResponseBody) SetProtocol(v string) *GetCustomDomainResponseBody {
	s.Protocol = &v
	return s
}

func (s *GetCustomDomainResponseBody) SetRouteConfig(v *RouteConfig) *GetCustomDomainResponseBody {
	s.RouteConfig = v
	return s
}

func (s *GetCustomDomainResponseBody) SetTlsConfig(v *TLSConfig) *GetCustomDomainResponseBody {
	s.TlsConfig = v
	return s
}

func (s *GetCustomDomainResponseBody) SetWafConfig(v *WAFConfig) *GetCustomDomainResponseBody {
	s.WafConfig = v
	return s
}

type GetCustomDomainResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetCustomDomainResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetCustomDomainResponse) String() string {
	return tea.Prettify(s)
}

func (s GetCustomDomainResponse) GoString() string {
	return s.String()
}

func (s *GetCustomDomainResponse) SetHeaders(v map[string]*string) *GetCustomDomainResponse {
	s.Headers = v
	return s
}

func (s *GetCustomDomainResponse) SetStatusCode(v int32) *GetCustomDomainResponse {
	s.StatusCode = &v
	return s
}

func (s *GetCustomDomainResponse) SetBody(v *GetCustomDomainResponseBody) *GetCustomDomainResponse {
	s.Body = v
	return s
}

type GetFunctionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetFunctionHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionHeaders) GoString() string {
	return s.String()
}

func (s *GetFunctionHeaders) SetCommonHeaders(v map[string]*string) *GetFunctionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetFunctionHeaders) SetXFcAccountId(v string) *GetFunctionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetFunctionHeaders) SetXFcDate(v string) *GetFunctionHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetFunctionHeaders) SetXFcTraceId(v string) *GetFunctionHeaders {
	s.XFcTraceId = &v
	return s
}

type GetFunctionRequest struct {
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s GetFunctionRequest) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionRequest) GoString() string {
	return s.String()
}

func (s *GetFunctionRequest) SetQualifier(v string) *GetFunctionRequest {
	s.Qualifier = &v
	return s
}

type GetFunctionResponseBody struct {
	// The port on which the HTTP server listens for the custom runtime or custom container runtime.
	CaPort *int32 `json:"caPort,omitempty" xml:"caPort,omitempty"`
	// The CRC-64 value of the function code package.
	CodeChecksum *string `json:"codeChecksum,omitempty" xml:"codeChecksum,omitempty"`
	// The size of the function code package. Unit: byte.
	CodeSize *int64 `json:"codeSize,omitempty" xml:"codeSize,omitempty"`
	// The number of vCPUs of the function. The value must be a multiple of 0.05.
	Cpu *float32 `json:"cpu,omitempty" xml:"cpu,omitempty"`
	// The time when the function was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The configurations of the custom container runtime. After you configure the custom container runtime, Function Compute can execute the function in a container created from a custom image.
	CustomContainerConfig *CustomContainerConfigInfo `json:"customContainerConfig,omitempty" xml:"customContainerConfig,omitempty"`
	// The custom DNS configurations of the function.
	CustomDNS *CustomDNS `json:"customDNS,omitempty" xml:"customDNS,omitempty"`
	// The custom health check configuration of the function. This parameter is applicable only to custom runtimes and custom containers.
	CustomHealthCheckConfig *CustomHealthCheckConfig `json:"customHealthCheckConfig,omitempty" xml:"customHealthCheckConfig,omitempty"`
	// The configurations of the custom runtime.
	CustomRuntimeConfig *CustomRuntimeConfig `json:"customRuntimeConfig,omitempty" xml:"customRuntimeConfig,omitempty"`
	// The description of the function.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The disk size of the function. Unit: MB. Valid values: 512 and 10240.
	DiskSize *int32 `json:"diskSize,omitempty" xml:"diskSize,omitempty"`
	// The environment variables that are configured for the function. You can obtain the values of the environment variables from the function. For more information, see [Environment variables](~~69777~~).
	EnvironmentVariables map[string]*string `json:"environmentVariables,omitempty" xml:"environmentVariables,omitempty"`
	// The ID that is generated by the system for the function. Each function ID is unique in Function Compute.
	FunctionId *string `json:"functionId,omitempty" xml:"functionId,omitempty"`
	// The name of the function.
	FunctionName *string `json:"functionName,omitempty" xml:"functionName,omitempty"`
	// The GPU memory capacity for the function. Unit: MB. The memory capacity must be a multiple of 1024 MB.
	GpuMemorySize *int32 `json:"gpuMemorySize,omitempty" xml:"gpuMemorySize,omitempty"`
	// The handler of the function. For more information, see [Function handler](~~157704~~).
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	// The timeout period for the execution of the initializer function. Unit: seconds. Default value: 3. Valid values: 1 to 300. When this period ends, the execution of the initializer function is terminated.
	InitializationTimeout *int32 `json:"initializationTimeout,omitempty" xml:"initializationTimeout,omitempty"`
	// The handler of the initializer function. The format of the value is determined by the programming language that you use. For more information, see [Initializer function](~~157704~~).
	Initializer *string `json:"initializer,omitempty" xml:"initializer,omitempty"`
	// The number of requests that can be concurrently processed by a single instance.
	InstanceConcurrency *int32 `json:"instanceConcurrency,omitempty" xml:"instanceConcurrency,omitempty"`
	// The lifecycle configurations of the instance.
	InstanceLifecycleConfig *InstanceLifecycleConfig `json:"instanceLifecycleConfig,omitempty" xml:"instanceLifecycleConfig,omitempty"`
	// The soft concurrency of the instance. You can use this parameter to implement graceful scale-up of instances. If the number of concurrent requests on an instance is greater than the number of the soft concurrency, the instance scale-up is triggered. For example, if your instance requires a long time to start, you can specify a suitable soft concurrency to start the instance in advance.
	//
	// The value must be less than or equal to that of the **instanceConcurrency** parameter.
	InstanceSoftConcurrency *int32 `json:"instanceSoftConcurrency,omitempty" xml:"instanceSoftConcurrency,omitempty"`
	// The instance type of the function. Valid values:
	//
	// *   **e1**: elastic instance
	// *   **c1**: performance instance
	// *   **fc.gpu.tesla.1**: GPU-accelerated instances (Tesla T4)
	// *   **fc.gpu.ampere.1**: GPU-accelerated instances (Ampere A10)
	// *   **g1**: same fc.gpu.tesla.1
	InstanceType *string `json:"instanceType,omitempty" xml:"instanceType,omitempty"`
	// The time when the function was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The list of layers (ARN V1 version).
	//
	// > If multiple layers exist, the layers are merged based on the order of array subscripts. The content of a layer with a smaller subscript overwrites the file with the same name in the layer with a larger subscript. >
	//
	// **
	//
	// **Warning:** This parameter is to be deprecated. Use layersArnV2.
	Layers []*string `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The list of layers (ARN V2 version).
	//
	// > If multiple layers exist, the layers are merged based on the order of array subscripts. The content of a layer with a smaller subscript overwrites the file that has the same name and a larger subscript in the layer.
	LayersArnV2 []*string `json:"layersArnV2,omitempty" xml:"layersArnV2,omitempty" type:"Repeated"`
	// The memory size for the function. Unit: MB. The memory size must be a multiple of 64 MB. The memory size varies based on the function instance type. For more information, see [Instance types](~~179379~~).
	MemorySize *int32 `json:"memorySize,omitempty" xml:"memorySize,omitempty"`
	// The runtime environment of the function. Valid values: **nodejs16**, **nodejs14**, **nodejs12**, **nodejs10**, **nodejs8**, **nodejs6**, **nodejs4.4**, **python3.9**, **python3**, **python2.7**, **java11**, **java8**, **go1**, **php7.2**, **dotnetcore2.1**, **custom**, and **custom-container**.
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The timeout period for the execution of the function. Unit: seconds. Default value: 60. Valid values: 1 to 600. When this period expires, the execution of the function is terminated.
	Timeout *int32 `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s GetFunctionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionResponseBody) GoString() string {
	return s.String()
}

func (s *GetFunctionResponseBody) SetCaPort(v int32) *GetFunctionResponseBody {
	s.CaPort = &v
	return s
}

func (s *GetFunctionResponseBody) SetCodeChecksum(v string) *GetFunctionResponseBody {
	s.CodeChecksum = &v
	return s
}

func (s *GetFunctionResponseBody) SetCodeSize(v int64) *GetFunctionResponseBody {
	s.CodeSize = &v
	return s
}

func (s *GetFunctionResponseBody) SetCpu(v float32) *GetFunctionResponseBody {
	s.Cpu = &v
	return s
}

func (s *GetFunctionResponseBody) SetCreatedTime(v string) *GetFunctionResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *GetFunctionResponseBody) SetCustomContainerConfig(v *CustomContainerConfigInfo) *GetFunctionResponseBody {
	s.CustomContainerConfig = v
	return s
}

func (s *GetFunctionResponseBody) SetCustomDNS(v *CustomDNS) *GetFunctionResponseBody {
	s.CustomDNS = v
	return s
}

func (s *GetFunctionResponseBody) SetCustomHealthCheckConfig(v *CustomHealthCheckConfig) *GetFunctionResponseBody {
	s.CustomHealthCheckConfig = v
	return s
}

func (s *GetFunctionResponseBody) SetCustomRuntimeConfig(v *CustomRuntimeConfig) *GetFunctionResponseBody {
	s.CustomRuntimeConfig = v
	return s
}

func (s *GetFunctionResponseBody) SetDescription(v string) *GetFunctionResponseBody {
	s.Description = &v
	return s
}

func (s *GetFunctionResponseBody) SetDiskSize(v int32) *GetFunctionResponseBody {
	s.DiskSize = &v
	return s
}

func (s *GetFunctionResponseBody) SetEnvironmentVariables(v map[string]*string) *GetFunctionResponseBody {
	s.EnvironmentVariables = v
	return s
}

func (s *GetFunctionResponseBody) SetFunctionId(v string) *GetFunctionResponseBody {
	s.FunctionId = &v
	return s
}

func (s *GetFunctionResponseBody) SetFunctionName(v string) *GetFunctionResponseBody {
	s.FunctionName = &v
	return s
}

func (s *GetFunctionResponseBody) SetGpuMemorySize(v int32) *GetFunctionResponseBody {
	s.GpuMemorySize = &v
	return s
}

func (s *GetFunctionResponseBody) SetHandler(v string) *GetFunctionResponseBody {
	s.Handler = &v
	return s
}

func (s *GetFunctionResponseBody) SetInitializationTimeout(v int32) *GetFunctionResponseBody {
	s.InitializationTimeout = &v
	return s
}

func (s *GetFunctionResponseBody) SetInitializer(v string) *GetFunctionResponseBody {
	s.Initializer = &v
	return s
}

func (s *GetFunctionResponseBody) SetInstanceConcurrency(v int32) *GetFunctionResponseBody {
	s.InstanceConcurrency = &v
	return s
}

func (s *GetFunctionResponseBody) SetInstanceLifecycleConfig(v *InstanceLifecycleConfig) *GetFunctionResponseBody {
	s.InstanceLifecycleConfig = v
	return s
}

func (s *GetFunctionResponseBody) SetInstanceSoftConcurrency(v int32) *GetFunctionResponseBody {
	s.InstanceSoftConcurrency = &v
	return s
}

func (s *GetFunctionResponseBody) SetInstanceType(v string) *GetFunctionResponseBody {
	s.InstanceType = &v
	return s
}

func (s *GetFunctionResponseBody) SetLastModifiedTime(v string) *GetFunctionResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *GetFunctionResponseBody) SetLayers(v []*string) *GetFunctionResponseBody {
	s.Layers = v
	return s
}

func (s *GetFunctionResponseBody) SetLayersArnV2(v []*string) *GetFunctionResponseBody {
	s.LayersArnV2 = v
	return s
}

func (s *GetFunctionResponseBody) SetMemorySize(v int32) *GetFunctionResponseBody {
	s.MemorySize = &v
	return s
}

func (s *GetFunctionResponseBody) SetRuntime(v string) *GetFunctionResponseBody {
	s.Runtime = &v
	return s
}

func (s *GetFunctionResponseBody) SetTimeout(v int32) *GetFunctionResponseBody {
	s.Timeout = &v
	return s
}

type GetFunctionResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetFunctionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetFunctionResponse) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionResponse) GoString() string {
	return s.String()
}

func (s *GetFunctionResponse) SetHeaders(v map[string]*string) *GetFunctionResponse {
	s.Headers = v
	return s
}

func (s *GetFunctionResponse) SetStatusCode(v int32) *GetFunctionResponse {
	s.StatusCode = &v
	return s
}

func (s *GetFunctionResponse) SetBody(v *GetFunctionResponseBody) *GetFunctionResponse {
	s.Body = v
	return s
}

type GetFunctionAsyncInvokeConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the Function Compute is called. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetFunctionAsyncInvokeConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionAsyncInvokeConfigHeaders) GoString() string {
	return s.String()
}

func (s *GetFunctionAsyncInvokeConfigHeaders) SetCommonHeaders(v map[string]*string) *GetFunctionAsyncInvokeConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetFunctionAsyncInvokeConfigHeaders) SetXFcAccountId(v string) *GetFunctionAsyncInvokeConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigHeaders) SetXFcDate(v string) *GetFunctionAsyncInvokeConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigHeaders) SetXFcTraceId(v string) *GetFunctionAsyncInvokeConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type GetFunctionAsyncInvokeConfigRequest struct {
	// The qualifier.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s GetFunctionAsyncInvokeConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionAsyncInvokeConfigRequest) GoString() string {
	return s.String()
}

func (s *GetFunctionAsyncInvokeConfigRequest) SetQualifier(v string) *GetFunctionAsyncInvokeConfigRequest {
	s.Qualifier = &v
	return s
}

type GetFunctionAsyncInvokeConfigResponseBody struct {
	// The time when the desktop group was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The configuration structure of the destination for asynchronous invocations.
	DestinationConfig *DestinationConfig `json:"destinationConfig,omitempty" xml:"destinationConfig,omitempty"`
	// The name of the function.
	Function *string `json:"function,omitempty" xml:"function,omitempty"`
	// The time when the configuration was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The maximum validity period of a message.
	MaxAsyncEventAgeInSeconds *int64 `json:"maxAsyncEventAgeInSeconds,omitempty" xml:"maxAsyncEventAgeInSeconds,omitempty"`
	// The maximum number of retries allowed after an asynchronous invocation fails.
	MaxAsyncRetryAttempts *int64 `json:"maxAsyncRetryAttempts,omitempty" xml:"maxAsyncRetryAttempts,omitempty"`
	// The version or alias of the service to which the function belongs.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The name of the service.
	Service *string `json:"service,omitempty" xml:"service,omitempty"`
	// Indicates whether the asynchronous task feature is enabled.
	//
	// *   **true**: The asynchronous task feature is enabled.
	// *   **false**: The asynchronous task feature is disabled.
	StatefulInvocation *bool `json:"statefulInvocation,omitempty" xml:"statefulInvocation,omitempty"`
}

func (s GetFunctionAsyncInvokeConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionAsyncInvokeConfigResponseBody) GoString() string {
	return s.String()
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetCreatedTime(v string) *GetFunctionAsyncInvokeConfigResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetDestinationConfig(v *DestinationConfig) *GetFunctionAsyncInvokeConfigResponseBody {
	s.DestinationConfig = v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetFunction(v string) *GetFunctionAsyncInvokeConfigResponseBody {
	s.Function = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetLastModifiedTime(v string) *GetFunctionAsyncInvokeConfigResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetMaxAsyncEventAgeInSeconds(v int64) *GetFunctionAsyncInvokeConfigResponseBody {
	s.MaxAsyncEventAgeInSeconds = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetMaxAsyncRetryAttempts(v int64) *GetFunctionAsyncInvokeConfigResponseBody {
	s.MaxAsyncRetryAttempts = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetQualifier(v string) *GetFunctionAsyncInvokeConfigResponseBody {
	s.Qualifier = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetService(v string) *GetFunctionAsyncInvokeConfigResponseBody {
	s.Service = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponseBody) SetStatefulInvocation(v bool) *GetFunctionAsyncInvokeConfigResponseBody {
	s.StatefulInvocation = &v
	return s
}

type GetFunctionAsyncInvokeConfigResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetFunctionAsyncInvokeConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetFunctionAsyncInvokeConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionAsyncInvokeConfigResponse) GoString() string {
	return s.String()
}

func (s *GetFunctionAsyncInvokeConfigResponse) SetHeaders(v map[string]*string) *GetFunctionAsyncInvokeConfigResponse {
	s.Headers = v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponse) SetStatusCode(v int32) *GetFunctionAsyncInvokeConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *GetFunctionAsyncInvokeConfigResponse) SetBody(v *GetFunctionAsyncInvokeConfigResponseBody) *GetFunctionAsyncInvokeConfigResponse {
	s.Body = v
	return s
}

type GetFunctionCodeHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetFunctionCodeHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionCodeHeaders) GoString() string {
	return s.String()
}

func (s *GetFunctionCodeHeaders) SetCommonHeaders(v map[string]*string) *GetFunctionCodeHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetFunctionCodeHeaders) SetXFcAccountId(v string) *GetFunctionCodeHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetFunctionCodeHeaders) SetXFcDate(v string) *GetFunctionCodeHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetFunctionCodeHeaders) SetXFcTraceId(v string) *GetFunctionCodeHeaders {
	s.XFcTraceId = &v
	return s
}

type GetFunctionCodeRequest struct {
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s GetFunctionCodeRequest) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionCodeRequest) GoString() string {
	return s.String()
}

func (s *GetFunctionCodeRequest) SetQualifier(v string) *GetFunctionCodeRequest {
	s.Qualifier = &v
	return s
}

type GetFunctionCodeResponseBody struct {
	// The CRC-64 value of the function code package.
	Checksum *string `json:"checksum,omitempty" xml:"checksum,omitempty"`
	// The URL of the function code package.
	Url *string `json:"url,omitempty" xml:"url,omitempty"`
}

func (s GetFunctionCodeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionCodeResponseBody) GoString() string {
	return s.String()
}

func (s *GetFunctionCodeResponseBody) SetChecksum(v string) *GetFunctionCodeResponseBody {
	s.Checksum = &v
	return s
}

func (s *GetFunctionCodeResponseBody) SetUrl(v string) *GetFunctionCodeResponseBody {
	s.Url = &v
	return s
}

type GetFunctionCodeResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetFunctionCodeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetFunctionCodeResponse) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionCodeResponse) GoString() string {
	return s.String()
}

func (s *GetFunctionCodeResponse) SetHeaders(v map[string]*string) *GetFunctionCodeResponse {
	s.Headers = v
	return s
}

func (s *GetFunctionCodeResponse) SetStatusCode(v int32) *GetFunctionCodeResponse {
	s.StatusCode = &v
	return s
}

func (s *GetFunctionCodeResponse) SetBody(v *GetFunctionCodeResponseBody) *GetFunctionCodeResponse {
	s.Body = v
	return s
}

type GetFunctionOnDemandConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The unique ID of the trace.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetFunctionOnDemandConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionOnDemandConfigHeaders) GoString() string {
	return s.String()
}

func (s *GetFunctionOnDemandConfigHeaders) SetCommonHeaders(v map[string]*string) *GetFunctionOnDemandConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetFunctionOnDemandConfigHeaders) SetXFcAccountId(v string) *GetFunctionOnDemandConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetFunctionOnDemandConfigHeaders) SetXFcDate(v string) *GetFunctionOnDemandConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetFunctionOnDemandConfigHeaders) SetXFcTraceId(v string) *GetFunctionOnDemandConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type GetFunctionOnDemandConfigRequest struct {
	// Service alias or LATEST. Other versions are not supported.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s GetFunctionOnDemandConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionOnDemandConfigRequest) GoString() string {
	return s.String()
}

func (s *GetFunctionOnDemandConfigRequest) SetQualifier(v string) *GetFunctionOnDemandConfigRequest {
	s.Qualifier = &v
	return s
}

type GetFunctionOnDemandConfigResponseBody struct {
	// The maximum number of instances.
	MaximumInstanceCount *int64 `json:"maximumInstanceCount,omitempty" xml:"maximumInstanceCount,omitempty"`
	// The description of the resource.
	Resource *string `json:"resource,omitempty" xml:"resource,omitempty"`
}

func (s GetFunctionOnDemandConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionOnDemandConfigResponseBody) GoString() string {
	return s.String()
}

func (s *GetFunctionOnDemandConfigResponseBody) SetMaximumInstanceCount(v int64) *GetFunctionOnDemandConfigResponseBody {
	s.MaximumInstanceCount = &v
	return s
}

func (s *GetFunctionOnDemandConfigResponseBody) SetResource(v string) *GetFunctionOnDemandConfigResponseBody {
	s.Resource = &v
	return s
}

type GetFunctionOnDemandConfigResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetFunctionOnDemandConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetFunctionOnDemandConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s GetFunctionOnDemandConfigResponse) GoString() string {
	return s.String()
}

func (s *GetFunctionOnDemandConfigResponse) SetHeaders(v map[string]*string) *GetFunctionOnDemandConfigResponse {
	s.Headers = v
	return s
}

func (s *GetFunctionOnDemandConfigResponse) SetStatusCode(v int32) *GetFunctionOnDemandConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *GetFunctionOnDemandConfigResponse) SetBody(v *GetFunctionOnDemandConfigResponseBody) *GetFunctionOnDemandConfigResponse {
	s.Body = v
	return s
}

type GetLayerVersionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetLayerVersionHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetLayerVersionHeaders) GoString() string {
	return s.String()
}

func (s *GetLayerVersionHeaders) SetCommonHeaders(v map[string]*string) *GetLayerVersionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetLayerVersionHeaders) SetXFcAccountId(v string) *GetLayerVersionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetLayerVersionHeaders) SetXFcDate(v string) *GetLayerVersionHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetLayerVersionHeaders) SetXFcTraceId(v string) *GetLayerVersionHeaders {
	s.XFcTraceId = &v
	return s
}

type GetLayerVersionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *Layer             `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetLayerVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s GetLayerVersionResponse) GoString() string {
	return s.String()
}

func (s *GetLayerVersionResponse) SetHeaders(v map[string]*string) *GetLayerVersionResponse {
	s.Headers = v
	return s
}

func (s *GetLayerVersionResponse) SetStatusCode(v int32) *GetLayerVersionResponse {
	s.StatusCode = &v
	return s
}

func (s *GetLayerVersionResponse) SetBody(v *Layer) *GetLayerVersionResponse {
	s.Body = v
	return s
}

type GetProvisionConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The start time when the function is invoked. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetProvisionConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetProvisionConfigHeaders) GoString() string {
	return s.String()
}

func (s *GetProvisionConfigHeaders) SetCommonHeaders(v map[string]*string) *GetProvisionConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetProvisionConfigHeaders) SetXFcAccountId(v string) *GetProvisionConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetProvisionConfigHeaders) SetXFcDate(v string) *GetProvisionConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetProvisionConfigHeaders) SetXFcTraceId(v string) *GetProvisionConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type GetProvisionConfigRequest struct {
	// The name of the alias.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s GetProvisionConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s GetProvisionConfigRequest) GoString() string {
	return s.String()
}

func (s *GetProvisionConfigRequest) SetQualifier(v string) *GetProvisionConfigRequest {
	s.Qualifier = &v
	return s
}

type GetProvisionConfigResponseBody struct {
	// Specifies whether to always allocate CPU to a function instance.
	AlwaysAllocateCPU *bool `json:"alwaysAllocateCPU,omitempty" xml:"alwaysAllocateCPU,omitempty"`
	// The actual number of provisioned instances.
	Current *int64 `json:"current,omitempty" xml:"current,omitempty"`
	// The error message returned if a provisioned instance fails to be created.
	CurrentError *string `json:"currentError,omitempty" xml:"currentError,omitempty"`
	// The description of the resource.
	Resource *string `json:"resource,omitempty" xml:"resource,omitempty"`
	// The configurations of scheduled auto scaling.
	ScheduledActions []*ScheduledActions `json:"scheduledActions,omitempty" xml:"scheduledActions,omitempty" type:"Repeated"`
	// The expected number of provisioned instances.
	Target *int64 `json:"target,omitempty" xml:"target,omitempty"`
	// The configurations of metric-based auto scaling.
	TargetTrackingPolicies []*TargetTrackingPolicies `json:"targetTrackingPolicies,omitempty" xml:"targetTrackingPolicies,omitempty" type:"Repeated"`
}

func (s GetProvisionConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetProvisionConfigResponseBody) GoString() string {
	return s.String()
}

func (s *GetProvisionConfigResponseBody) SetAlwaysAllocateCPU(v bool) *GetProvisionConfigResponseBody {
	s.AlwaysAllocateCPU = &v
	return s
}

func (s *GetProvisionConfigResponseBody) SetCurrent(v int64) *GetProvisionConfigResponseBody {
	s.Current = &v
	return s
}

func (s *GetProvisionConfigResponseBody) SetCurrentError(v string) *GetProvisionConfigResponseBody {
	s.CurrentError = &v
	return s
}

func (s *GetProvisionConfigResponseBody) SetResource(v string) *GetProvisionConfigResponseBody {
	s.Resource = &v
	return s
}

func (s *GetProvisionConfigResponseBody) SetScheduledActions(v []*ScheduledActions) *GetProvisionConfigResponseBody {
	s.ScheduledActions = v
	return s
}

func (s *GetProvisionConfigResponseBody) SetTarget(v int64) *GetProvisionConfigResponseBody {
	s.Target = &v
	return s
}

func (s *GetProvisionConfigResponseBody) SetTargetTrackingPolicies(v []*TargetTrackingPolicies) *GetProvisionConfigResponseBody {
	s.TargetTrackingPolicies = v
	return s
}

type GetProvisionConfigResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetProvisionConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetProvisionConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s GetProvisionConfigResponse) GoString() string {
	return s.String()
}

func (s *GetProvisionConfigResponse) SetHeaders(v map[string]*string) *GetProvisionConfigResponse {
	s.Headers = v
	return s
}

func (s *GetProvisionConfigResponse) SetStatusCode(v int32) *GetProvisionConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *GetProvisionConfigResponse) SetBody(v *GetProvisionConfigResponseBody) *GetProvisionConfigResponse {
	s.Body = v
	return s
}

type GetResourceTagsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetResourceTagsHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetResourceTagsHeaders) GoString() string {
	return s.String()
}

func (s *GetResourceTagsHeaders) SetCommonHeaders(v map[string]*string) *GetResourceTagsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetResourceTagsHeaders) SetXFcAccountId(v string) *GetResourceTagsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetResourceTagsHeaders) SetXFcDate(v string) *GetResourceTagsHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetResourceTagsHeaders) SetXFcTraceId(v string) *GetResourceTagsHeaders {
	s.XFcTraceId = &v
	return s
}

type GetResourceTagsRequest struct {
	// The Alibaba Cloud Resource Name (ARN) of the resource.
	//
	// > You can use the value of this parameter to query the information about the resource, such as the account, service, and region information of the resource. You can manage tags only for services for top level resources.
	ResourceArn *string `json:"resourceArn,omitempty" xml:"resourceArn,omitempty"`
}

func (s GetResourceTagsRequest) String() string {
	return tea.Prettify(s)
}

func (s GetResourceTagsRequest) GoString() string {
	return s.String()
}

func (s *GetResourceTagsRequest) SetResourceArn(v string) *GetResourceTagsRequest {
	s.ResourceArn = &v
	return s
}

type GetResourceTagsResponseBody struct {
	// The ARN of the resource.
	//
	// > You can use the value of this parameter to query the information about the resource, such as the account, service, and region information of the resource.
	ResourceArn *string `json:"resourceArn,omitempty" xml:"resourceArn,omitempty"`
	// The tag dictionary.
	Tags map[string]*string `json:"tags,omitempty" xml:"tags,omitempty"`
}

func (s GetResourceTagsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetResourceTagsResponseBody) GoString() string {
	return s.String()
}

func (s *GetResourceTagsResponseBody) SetResourceArn(v string) *GetResourceTagsResponseBody {
	s.ResourceArn = &v
	return s
}

func (s *GetResourceTagsResponseBody) SetTags(v map[string]*string) *GetResourceTagsResponseBody {
	s.Tags = v
	return s
}

type GetResourceTagsResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetResourceTagsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetResourceTagsResponse) String() string {
	return tea.Prettify(s)
}

func (s GetResourceTagsResponse) GoString() string {
	return s.String()
}

func (s *GetResourceTagsResponse) SetHeaders(v map[string]*string) *GetResourceTagsResponse {
	s.Headers = v
	return s
}

func (s *GetResourceTagsResponse) SetStatusCode(v int32) *GetResourceTagsResponse {
	s.StatusCode = &v
	return s
}

func (s *GetResourceTagsResponse) SetBody(v *GetResourceTagsResponseBody) *GetResourceTagsResponse {
	s.Body = v
	return s
}

type GetServiceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the function is invoked. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetServiceHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetServiceHeaders) GoString() string {
	return s.String()
}

func (s *GetServiceHeaders) SetCommonHeaders(v map[string]*string) *GetServiceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetServiceHeaders) SetXFcAccountId(v string) *GetServiceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetServiceHeaders) SetXFcDate(v string) *GetServiceHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetServiceHeaders) SetXFcTraceId(v string) *GetServiceHeaders {
	s.XFcTraceId = &v
	return s
}

type GetServiceRequest struct {
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s GetServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s GetServiceRequest) GoString() string {
	return s.String()
}

func (s *GetServiceRequest) SetQualifier(v string) *GetServiceRequest {
	s.Qualifier = &v
	return s
}

type GetServiceResponseBody struct {
	// The time when the service was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the service.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Specifies whether to allow functions to access the Internet. Valid values:
	//
	// *   **true**: allows functions in the specified service to access the Internet.
	// *   **false**: does not allow functions in the specified service to access the Internet.
	InternetAccess *bool `json:"internetAccess,omitempty" xml:"internetAccess,omitempty"`
	// The time when the service was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The log configuration, which specifies a Logstore to store function execution logs.
	LogConfig *LogConfig `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	// The configuration of the NAS file system. The configuration allows functions in the specified service in Function Compute to access the NAS file system.
	NasConfig *NASConfig `json:"nasConfig,omitempty" xml:"nasConfig,omitempty"`
	// The OSS mount configurations.
	OssMountConfig *OSSMountConfig `json:"ossMountConfig,omitempty" xml:"ossMountConfig,omitempty"`
	// The RAM role that is used to grant required permissions to Function Compute. Scenarios:
	//
	// *   Sends function execution logs to your Logstore.
	// *   Generates a token for a function to access other cloud resources during function execution.
	Role *string `json:"role,omitempty" xml:"role,omitempty"`
	// The unique ID generated by the system for the service.
	ServiceId *string `json:"serviceId,omitempty" xml:"serviceId,omitempty"`
	// The name of the service.
	ServiceName *string `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
	// The configurations of Tracing Analysis. After you configure Tracing Analysis for a service in Function Compute, you can record the execution duration of a request, view the amount of cold start time for a function, and record the execution duration of a function. For more information, see [Overview](~~189804~~).
	TracingConfig *TracingConfig `json:"tracingConfig,omitempty" xml:"tracingConfig,omitempty"`
	// The VPC configuration. The configuration allows a function to access the specified VPC.
	VpcConfig *VPCConfig `json:"vpcConfig,omitempty" xml:"vpcConfig,omitempty"`
}

func (s GetServiceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetServiceResponseBody) GoString() string {
	return s.String()
}

func (s *GetServiceResponseBody) SetCreatedTime(v string) *GetServiceResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *GetServiceResponseBody) SetDescription(v string) *GetServiceResponseBody {
	s.Description = &v
	return s
}

func (s *GetServiceResponseBody) SetInternetAccess(v bool) *GetServiceResponseBody {
	s.InternetAccess = &v
	return s
}

func (s *GetServiceResponseBody) SetLastModifiedTime(v string) *GetServiceResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *GetServiceResponseBody) SetLogConfig(v *LogConfig) *GetServiceResponseBody {
	s.LogConfig = v
	return s
}

func (s *GetServiceResponseBody) SetNasConfig(v *NASConfig) *GetServiceResponseBody {
	s.NasConfig = v
	return s
}

func (s *GetServiceResponseBody) SetOssMountConfig(v *OSSMountConfig) *GetServiceResponseBody {
	s.OssMountConfig = v
	return s
}

func (s *GetServiceResponseBody) SetRole(v string) *GetServiceResponseBody {
	s.Role = &v
	return s
}

func (s *GetServiceResponseBody) SetServiceId(v string) *GetServiceResponseBody {
	s.ServiceId = &v
	return s
}

func (s *GetServiceResponseBody) SetServiceName(v string) *GetServiceResponseBody {
	s.ServiceName = &v
	return s
}

func (s *GetServiceResponseBody) SetTracingConfig(v *TracingConfig) *GetServiceResponseBody {
	s.TracingConfig = v
	return s
}

func (s *GetServiceResponseBody) SetVpcConfig(v *VPCConfig) *GetServiceResponseBody {
	s.VpcConfig = v
	return s
}

type GetServiceResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetServiceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s GetServiceResponse) GoString() string {
	return s.String()
}

func (s *GetServiceResponse) SetHeaders(v map[string]*string) *GetServiceResponse {
	s.Headers = v
	return s
}

func (s *GetServiceResponse) SetStatusCode(v int32) *GetServiceResponse {
	s.StatusCode = &v
	return s
}

func (s *GetServiceResponse) SetBody(v *GetServiceResponseBody) *GetServiceResponse {
	s.Body = v
	return s
}

type GetStatefulAsyncInvocationHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The CRC-64 value of the function code package. This value is used to check data integrity. The value is automatically calculated by the tool.
	XFcCodeChecksum *string `json:"X-Fc-Code-Checksum,omitempty" xml:"X-Fc-Code-Checksum,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The invocation method.
	//
	// - **Sync**: synchronous invocation
	// - **Async**: asynchronous invocation
	XFcInvocationType *string `json:"X-Fc-Invocation-Type,omitempty" xml:"X-Fc-Invocation-Type,omitempty"`
	// The method used to return logs. Valid values:
	//
	// - **Tail**: returns the last 4 KB of logs that are generated for the current request.
	// - **None**: does not return logs for the current request. This is the default value.
	XFcLogType *string `json:"X-Fc-Log-Type,omitempty" xml:"X-Fc-Log-Type,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetStatefulAsyncInvocationHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetStatefulAsyncInvocationHeaders) GoString() string {
	return s.String()
}

func (s *GetStatefulAsyncInvocationHeaders) SetCommonHeaders(v map[string]*string) *GetStatefulAsyncInvocationHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetStatefulAsyncInvocationHeaders) SetXFcAccountId(v string) *GetStatefulAsyncInvocationHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetStatefulAsyncInvocationHeaders) SetXFcCodeChecksum(v string) *GetStatefulAsyncInvocationHeaders {
	s.XFcCodeChecksum = &v
	return s
}

func (s *GetStatefulAsyncInvocationHeaders) SetXFcDate(v string) *GetStatefulAsyncInvocationHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetStatefulAsyncInvocationHeaders) SetXFcInvocationType(v string) *GetStatefulAsyncInvocationHeaders {
	s.XFcInvocationType = &v
	return s
}

func (s *GetStatefulAsyncInvocationHeaders) SetXFcLogType(v string) *GetStatefulAsyncInvocationHeaders {
	s.XFcLogType = &v
	return s
}

func (s *GetStatefulAsyncInvocationHeaders) SetXFcTraceId(v string) *GetStatefulAsyncInvocationHeaders {
	s.XFcTraceId = &v
	return s
}

type GetStatefulAsyncInvocationRequest struct {
	// The version or alias of the service to which the asynchronous task belongs.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s GetStatefulAsyncInvocationRequest) String() string {
	return tea.Prettify(s)
}

func (s GetStatefulAsyncInvocationRequest) GoString() string {
	return s.String()
}

func (s *GetStatefulAsyncInvocationRequest) SetQualifier(v string) *GetStatefulAsyncInvocationRequest {
	s.Qualifier = &v
	return s
}

type GetStatefulAsyncInvocationResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *StatefulAsyncInvocation `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetStatefulAsyncInvocationResponse) String() string {
	return tea.Prettify(s)
}

func (s GetStatefulAsyncInvocationResponse) GoString() string {
	return s.String()
}

func (s *GetStatefulAsyncInvocationResponse) SetHeaders(v map[string]*string) *GetStatefulAsyncInvocationResponse {
	s.Headers = v
	return s
}

func (s *GetStatefulAsyncInvocationResponse) SetStatusCode(v int32) *GetStatefulAsyncInvocationResponse {
	s.StatusCode = &v
	return s
}

func (s *GetStatefulAsyncInvocationResponse) SetBody(v *StatefulAsyncInvocation) *GetStatefulAsyncInvocationResponse {
	s.Body = v
	return s
}

type GetTriggerHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the request is initiated on the client. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s GetTriggerHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetTriggerHeaders) GoString() string {
	return s.String()
}

func (s *GetTriggerHeaders) SetCommonHeaders(v map[string]*string) *GetTriggerHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetTriggerHeaders) SetXFcAccountId(v string) *GetTriggerHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *GetTriggerHeaders) SetXFcDate(v string) *GetTriggerHeaders {
	s.XFcDate = &v
	return s
}

func (s *GetTriggerHeaders) SetXFcTraceId(v string) *GetTriggerHeaders {
	s.XFcTraceId = &v
	return s
}

type GetTriggerResponseBody struct {
	// The time when the trigger was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the trigger.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The domain name used to invoke the function by using HTTP. You can add this domain name as the prefix to the endpoint of Function Compute. This way, you can invoke the function that corresponds to the trigger by using HTTP. For example, `{domainName}.cn-shanghai.fc.aliyuncs.com`.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The ARN of the RAM role that is used by the event source to invoke the function.
	InvocationRole *string `json:"invocationRole,omitempty" xml:"invocationRole,omitempty"`
	// The time when the trigger was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The ARN of the event source.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
	// The configurations of the trigger. The configurations vary based on the trigger type. For more information about the format, see the following topics:
	//
	// *   OSS trigger: [OSSTriggerConfig](~~struct:OSSTriggerConfig~~).
	// *   Log Service trigger: [LogTriggerConfig](~~struct:LogTriggerConfig~~).
	// *   Time trigger: [TimeTriggerConfig](~~struct:TimeTriggerConfig~~).
	// *   HTTP trigger: [HTTPTriggerConfig](~~struct:HTTPTriggerConfig~~).
	// *   Tablestore trigger: Specify the **SourceArn** parameter and leave this parameter empty.
	// *   Alibaba Cloud CDN event trigger: [CDNEventsTriggerConfig](~~struct:CDNEventsTriggerConfig~~).
	// *   MNS topic trigger: [MnsTopicTriggerConfig](~~struct:MnsTopicTriggerConfig~~).
	TriggerConfig *string `json:"triggerConfig,omitempty" xml:"triggerConfig,omitempty"`
	// The unique ID of the trigger.
	TriggerId *string `json:"triggerId,omitempty" xml:"triggerId,omitempty"`
	// The name of the trigger.
	TriggerName *string `json:"triggerName,omitempty" xml:"triggerName,omitempty"`
	// The trigger type, such as **oss**, **log**, **tablestore**, **timer**, **http**, **cdn_events**, and **mns_topic**.
	TriggerType *string `json:"triggerType,omitempty" xml:"triggerType,omitempty"`
	// The public domain address. You can access HTTP triggers over the Internet by using HTTP or HTTPS.
	UrlInternet *string `json:"urlInternet,omitempty" xml:"urlInternet,omitempty"`
	// The private endpoint. In a VPC, you can access HTTP triggers by using HTTP or HTTPS.
	UrlIntranet *string `json:"urlIntranet,omitempty" xml:"urlIntranet,omitempty"`
}

func (s GetTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *GetTriggerResponseBody) SetCreatedTime(v string) *GetTriggerResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *GetTriggerResponseBody) SetDescription(v string) *GetTriggerResponseBody {
	s.Description = &v
	return s
}

func (s *GetTriggerResponseBody) SetDomainName(v string) *GetTriggerResponseBody {
	s.DomainName = &v
	return s
}

func (s *GetTriggerResponseBody) SetInvocationRole(v string) *GetTriggerResponseBody {
	s.InvocationRole = &v
	return s
}

func (s *GetTriggerResponseBody) SetLastModifiedTime(v string) *GetTriggerResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *GetTriggerResponseBody) SetQualifier(v string) *GetTriggerResponseBody {
	s.Qualifier = &v
	return s
}

func (s *GetTriggerResponseBody) SetSourceArn(v string) *GetTriggerResponseBody {
	s.SourceArn = &v
	return s
}

func (s *GetTriggerResponseBody) SetTriggerConfig(v string) *GetTriggerResponseBody {
	s.TriggerConfig = &v
	return s
}

func (s *GetTriggerResponseBody) SetTriggerId(v string) *GetTriggerResponseBody {
	s.TriggerId = &v
	return s
}

func (s *GetTriggerResponseBody) SetTriggerName(v string) *GetTriggerResponseBody {
	s.TriggerName = &v
	return s
}

func (s *GetTriggerResponseBody) SetTriggerType(v string) *GetTriggerResponseBody {
	s.TriggerType = &v
	return s
}

func (s *GetTriggerResponseBody) SetUrlInternet(v string) *GetTriggerResponseBody {
	s.UrlInternet = &v
	return s
}

func (s *GetTriggerResponseBody) SetUrlIntranet(v string) *GetTriggerResponseBody {
	s.UrlIntranet = &v
	return s
}

type GetTriggerResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s GetTriggerResponse) GoString() string {
	return s.String()
}

func (s *GetTriggerResponse) SetHeaders(v map[string]*string) *GetTriggerResponse {
	s.Headers = v
	return s
}

func (s *GetTriggerResponse) SetStatusCode(v int32) *GetTriggerResponse {
	s.StatusCode = &v
	return s
}

func (s *GetTriggerResponse) SetBody(v *GetTriggerResponseBody) *GetTriggerResponse {
	s.Body = v
	return s
}

type InvokeFunctionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the function is invoked. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The method used to invoke the function. Valid values:
	//
	// *   **Sync**: synchronous
	// *   **Async**: asynchronous
	XFcInvocationType *string `json:"X-Fc-Invocation-Type,omitempty" xml:"X-Fc-Invocation-Type,omitempty"`
	// The method used to return logs. Valid values:
	//
	// *   **Tail**: returns the last 4 KB of logs that are generated for the current request.
	// *   **None**: No logs are returned for the current request. Default value: None.
	XFcLogType *string `json:"X-Fc-Log-Type,omitempty" xml:"X-Fc-Log-Type,omitempty"`
	// The ID of the asynchronous task. You must enable the asynchronous task feature in advance.
	//
	// > When you use an SDK to invoke a function, we recommend that you specify a business-related ID to facilitate subsequent operations. For example, you can use the video name as the invocation ID for a video-processing function. This way, you can use the ID to check whether the video is processed or terminate the processing of the video. The ID must start with a letter or an underscore (\_) and can contain letters, digits, underscores (\_), and hyphens (-). The ID can be up to 128 characters in length. If you do not specify the ID of the asynchronous invocation, Function Compute automatically generates an ID.
	XFcStatefulAsyncInvocationId *string `json:"X-Fc-Stateful-Async-Invocation-Id,omitempty" xml:"X-Fc-Stateful-Async-Invocation-Id,omitempty"`
	// The trace ID of the request for Function Compute API. The value is the same as that of the **requestId** parameter in the response.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s InvokeFunctionHeaders) String() string {
	return tea.Prettify(s)
}

func (s InvokeFunctionHeaders) GoString() string {
	return s.String()
}

func (s *InvokeFunctionHeaders) SetCommonHeaders(v map[string]*string) *InvokeFunctionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *InvokeFunctionHeaders) SetXFcAccountId(v string) *InvokeFunctionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *InvokeFunctionHeaders) SetXFcDate(v string) *InvokeFunctionHeaders {
	s.XFcDate = &v
	return s
}

func (s *InvokeFunctionHeaders) SetXFcInvocationType(v string) *InvokeFunctionHeaders {
	s.XFcInvocationType = &v
	return s
}

func (s *InvokeFunctionHeaders) SetXFcLogType(v string) *InvokeFunctionHeaders {
	s.XFcLogType = &v
	return s
}

func (s *InvokeFunctionHeaders) SetXFcStatefulAsyncInvocationId(v string) *InvokeFunctionHeaders {
	s.XFcStatefulAsyncInvocationId = &v
	return s
}

func (s *InvokeFunctionHeaders) SetXFcTraceId(v string) *InvokeFunctionHeaders {
	s.XFcTraceId = &v
	return s
}

type InvokeFunctionRequest struct {
	// The event to be processed by the function. Set this parameter to a binary string. Function Compute passes the event to the function for processing.
	Body []byte `json:"body,omitempty" xml:"body,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s InvokeFunctionRequest) String() string {
	return tea.Prettify(s)
}

func (s InvokeFunctionRequest) GoString() string {
	return s.String()
}

func (s *InvokeFunctionRequest) SetBody(v []byte) *InvokeFunctionRequest {
	s.Body = v
	return s
}

func (s *InvokeFunctionRequest) SetQualifier(v string) *InvokeFunctionRequest {
	s.Qualifier = &v
	return s
}

type InvokeFunctionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []byte             `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s InvokeFunctionResponse) String() string {
	return tea.Prettify(s)
}

func (s InvokeFunctionResponse) GoString() string {
	return s.String()
}

func (s *InvokeFunctionResponse) SetHeaders(v map[string]*string) *InvokeFunctionResponse {
	s.Headers = v
	return s
}

func (s *InvokeFunctionResponse) SetStatusCode(v int32) *InvokeFunctionResponse {
	s.StatusCode = &v
	return s
}

func (s *InvokeFunctionResponse) SetBody(v []byte) *InvokeFunctionResponse {
	s.Body = v
	return s
}

type ListAliasesHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The start time when the function is invoked. Specify the time in the yyyy-mm-ddhh:mm:ss format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListAliasesHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListAliasesHeaders) GoString() string {
	return s.String()
}

func (s *ListAliasesHeaders) SetCommonHeaders(v map[string]*string) *ListAliasesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListAliasesHeaders) SetXFcAccountId(v string) *ListAliasesHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListAliasesHeaders) SetXFcDate(v string) *ListAliasesHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListAliasesHeaders) SetXFcTraceId(v string) *ListAliasesHeaders {
	s.XFcTraceId = &v
	return s
}

type ListAliasesRequest struct {
	// The maximum number of resources to return.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token used to obtain more results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The prefix.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The starting position of the result list.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListAliasesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListAliasesRequest) GoString() string {
	return s.String()
}

func (s *ListAliasesRequest) SetLimit(v int32) *ListAliasesRequest {
	s.Limit = &v
	return s
}

func (s *ListAliasesRequest) SetNextToken(v string) *ListAliasesRequest {
	s.NextToken = &v
	return s
}

func (s *ListAliasesRequest) SetPrefix(v string) *ListAliasesRequest {
	s.Prefix = &v
	return s
}

func (s *ListAliasesRequest) SetStartKey(v string) *ListAliasesRequest {
	s.StartKey = &v
	return s
}

type ListAliasesResponseBody struct {
	// The list of aliases.
	Aliases []*ListAliasesResponseBodyAliases `json:"aliases,omitempty" xml:"aliases,omitempty" type:"Repeated"`
	// The token used to obtain more results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListAliasesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListAliasesResponseBody) GoString() string {
	return s.String()
}

func (s *ListAliasesResponseBody) SetAliases(v []*ListAliasesResponseBodyAliases) *ListAliasesResponseBody {
	s.Aliases = v
	return s
}

func (s *ListAliasesResponseBody) SetNextToken(v string) *ListAliasesResponseBody {
	s.NextToken = &v
	return s
}

type ListAliasesResponseBodyAliases struct {
	// The weight of the canary release version.
	AdditionalVersionWeight map[string]*float32 `json:"additionalVersionWeight,omitempty" xml:"additionalVersionWeight,omitempty"`
	// The name of the alias.
	AliasName *string `json:"aliasName,omitempty" xml:"aliasName,omitempty"`
	// The creation time.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the alias.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The last update time.
	LastModifiedTime *string      `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	ResolvePolicy    *string      `json:"resolvePolicy,omitempty" xml:"resolvePolicy,omitempty"`
	RoutePolicy      *RoutePolicy `json:"routePolicy,omitempty" xml:"routePolicy,omitempty"`
	// The ID of the version.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s ListAliasesResponseBodyAliases) String() string {
	return tea.Prettify(s)
}

func (s ListAliasesResponseBodyAliases) GoString() string {
	return s.String()
}

func (s *ListAliasesResponseBodyAliases) SetAdditionalVersionWeight(v map[string]*float32) *ListAliasesResponseBodyAliases {
	s.AdditionalVersionWeight = v
	return s
}

func (s *ListAliasesResponseBodyAliases) SetAliasName(v string) *ListAliasesResponseBodyAliases {
	s.AliasName = &v
	return s
}

func (s *ListAliasesResponseBodyAliases) SetCreatedTime(v string) *ListAliasesResponseBodyAliases {
	s.CreatedTime = &v
	return s
}

func (s *ListAliasesResponseBodyAliases) SetDescription(v string) *ListAliasesResponseBodyAliases {
	s.Description = &v
	return s
}

func (s *ListAliasesResponseBodyAliases) SetLastModifiedTime(v string) *ListAliasesResponseBodyAliases {
	s.LastModifiedTime = &v
	return s
}

func (s *ListAliasesResponseBodyAliases) SetResolvePolicy(v string) *ListAliasesResponseBodyAliases {
	s.ResolvePolicy = &v
	return s
}

func (s *ListAliasesResponseBodyAliases) SetRoutePolicy(v *RoutePolicy) *ListAliasesResponseBodyAliases {
	s.RoutePolicy = v
	return s
}

func (s *ListAliasesResponseBodyAliases) SetVersionId(v string) *ListAliasesResponseBodyAliases {
	s.VersionId = &v
	return s
}

type ListAliasesResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListAliasesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListAliasesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListAliasesResponse) GoString() string {
	return s.String()
}

func (s *ListAliasesResponse) SetHeaders(v map[string]*string) *ListAliasesResponse {
	s.Headers = v
	return s
}

func (s *ListAliasesResponse) SetStatusCode(v int32) *ListAliasesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListAliasesResponse) SetBody(v *ListAliasesResponseBody) *ListAliasesResponse {
	s.Body = v
	return s
}

type ListCustomDomainsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the operation is called. The format is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListCustomDomainsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListCustomDomainsHeaders) GoString() string {
	return s.String()
}

func (s *ListCustomDomainsHeaders) SetCommonHeaders(v map[string]*string) *ListCustomDomainsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListCustomDomainsHeaders) SetXFcAccountId(v string) *ListCustomDomainsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListCustomDomainsHeaders) SetXFcDate(v string) *ListCustomDomainsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListCustomDomainsHeaders) SetXFcTraceId(v string) *ListCustomDomainsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListCustomDomainsRequest struct {
	// The maximum number of resources to return. Valid values: \[0,100]. Default value: 20. The number of returned results is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The pagination token to use to request the next page of results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The prefix that the returned domain names must contain.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The returned resources are sorted in alphabetical order, and the resources that include and follow the resource specified by the startKey parameter are returned.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListCustomDomainsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListCustomDomainsRequest) GoString() string {
	return s.String()
}

func (s *ListCustomDomainsRequest) SetLimit(v int32) *ListCustomDomainsRequest {
	s.Limit = &v
	return s
}

func (s *ListCustomDomainsRequest) SetNextToken(v string) *ListCustomDomainsRequest {
	s.NextToken = &v
	return s
}

func (s *ListCustomDomainsRequest) SetPrefix(v string) *ListCustomDomainsRequest {
	s.Prefix = &v
	return s
}

func (s *ListCustomDomainsRequest) SetStartKey(v string) *ListCustomDomainsRequest {
	s.StartKey = &v
	return s
}

type ListCustomDomainsResponseBody struct {
	// The information about custom domain names.
	CustomDomains []*ListCustomDomainsResponseBodyCustomDomains `json:"customDomains,omitempty" xml:"customDomains,omitempty" type:"Repeated"`
	// The pagination token to use to request the next page of results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListCustomDomainsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListCustomDomainsResponseBody) GoString() string {
	return s.String()
}

func (s *ListCustomDomainsResponseBody) SetCustomDomains(v []*ListCustomDomainsResponseBodyCustomDomains) *ListCustomDomainsResponseBody {
	s.CustomDomains = v
	return s
}

func (s *ListCustomDomainsResponseBody) SetNextToken(v string) *ListCustomDomainsResponseBody {
	s.NextToken = &v
	return s
}

type ListCustomDomainsResponseBodyCustomDomains struct {
	// The ID of your Alibaba Cloud account.
	AccountId *string `json:"accountId,omitempty" xml:"accountId,omitempty"`
	// The version of the API.
	ApiVersion *string `json:"apiVersion,omitempty" xml:"apiVersion,omitempty"`
	// The configurations of the HTTPS certificate.
	CertConfig *CertConfig `json:"certConfig,omitempty" xml:"certConfig,omitempty"`
	// The time when the custom domain name was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The domain name.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The time when the domain name was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The protocol type that is supported by the custom domain name.
	//
	// *   **HTTP**: Only HTTP is supported.
	// *   **HTTPS**: Only HTTPS is supported.
	// *   **HTTP,HTTPS**: HTTP and HTTPS are supported.
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// The route table that maps the paths to functions when the functions are invoked by using the custom domain name.
	RouteConfig *RouteConfig `json:"routeConfig,omitempty" xml:"routeConfig,omitempty"`
	// The Transport Layer Security (TLS) configuration.
	TlsConfig *TLSConfig `json:"tlsConfig,omitempty" xml:"tlsConfig,omitempty"`
	// The Web Application Firewall (WAF) configuration.
	WafConfig *WAFConfig `json:"wafConfig,omitempty" xml:"wafConfig,omitempty"`
}

func (s ListCustomDomainsResponseBodyCustomDomains) String() string {
	return tea.Prettify(s)
}

func (s ListCustomDomainsResponseBodyCustomDomains) GoString() string {
	return s.String()
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetAccountId(v string) *ListCustomDomainsResponseBodyCustomDomains {
	s.AccountId = &v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetApiVersion(v string) *ListCustomDomainsResponseBodyCustomDomains {
	s.ApiVersion = &v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetCertConfig(v *CertConfig) *ListCustomDomainsResponseBodyCustomDomains {
	s.CertConfig = v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetCreatedTime(v string) *ListCustomDomainsResponseBodyCustomDomains {
	s.CreatedTime = &v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetDomainName(v string) *ListCustomDomainsResponseBodyCustomDomains {
	s.DomainName = &v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetLastModifiedTime(v string) *ListCustomDomainsResponseBodyCustomDomains {
	s.LastModifiedTime = &v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetProtocol(v string) *ListCustomDomainsResponseBodyCustomDomains {
	s.Protocol = &v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetRouteConfig(v *RouteConfig) *ListCustomDomainsResponseBodyCustomDomains {
	s.RouteConfig = v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetTlsConfig(v *TLSConfig) *ListCustomDomainsResponseBodyCustomDomains {
	s.TlsConfig = v
	return s
}

func (s *ListCustomDomainsResponseBodyCustomDomains) SetWafConfig(v *WAFConfig) *ListCustomDomainsResponseBodyCustomDomains {
	s.WafConfig = v
	return s
}

type ListCustomDomainsResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListCustomDomainsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListCustomDomainsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListCustomDomainsResponse) GoString() string {
	return s.String()
}

func (s *ListCustomDomainsResponse) SetHeaders(v map[string]*string) *ListCustomDomainsResponse {
	s.Headers = v
	return s
}

func (s *ListCustomDomainsResponse) SetStatusCode(v int32) *ListCustomDomainsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListCustomDomainsResponse) SetBody(v *ListCustomDomainsResponseBody) *ListCustomDomainsResponse {
	s.Body = v
	return s
}

type ListEventSourcesHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListEventSourcesHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListEventSourcesHeaders) GoString() string {
	return s.String()
}

func (s *ListEventSourcesHeaders) SetCommonHeaders(v map[string]*string) *ListEventSourcesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListEventSourcesHeaders) SetXFcAccountId(v string) *ListEventSourcesHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListEventSourcesHeaders) SetXFcDate(v string) *ListEventSourcesHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListEventSourcesHeaders) SetXFcTraceId(v string) *ListEventSourcesHeaders {
	s.XFcTraceId = &v
	return s
}

type ListEventSourcesRequest struct {
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s ListEventSourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListEventSourcesRequest) GoString() string {
	return s.String()
}

func (s *ListEventSourcesRequest) SetQualifier(v string) *ListEventSourcesRequest {
	s.Qualifier = &v
	return s
}

type ListEventSourcesResponseBody struct {
	// The information about event sources.
	EventSources []*ListEventSourcesResponseBodyEventSources `json:"eventSources,omitempty" xml:"eventSources,omitempty" type:"Repeated"`
}

func (s ListEventSourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListEventSourcesResponseBody) GoString() string {
	return s.String()
}

func (s *ListEventSourcesResponseBody) SetEventSources(v []*ListEventSourcesResponseBodyEventSources) *ListEventSourcesResponseBody {
	s.EventSources = v
	return s
}

type ListEventSourcesResponseBodyEventSources struct {
	// The time when the event source was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The ARN of the event source.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
}

func (s ListEventSourcesResponseBodyEventSources) String() string {
	return tea.Prettify(s)
}

func (s ListEventSourcesResponseBodyEventSources) GoString() string {
	return s.String()
}

func (s *ListEventSourcesResponseBodyEventSources) SetCreatedTime(v string) *ListEventSourcesResponseBodyEventSources {
	s.CreatedTime = &v
	return s
}

func (s *ListEventSourcesResponseBodyEventSources) SetSourceArn(v string) *ListEventSourcesResponseBodyEventSources {
	s.SourceArn = &v
	return s
}

type ListEventSourcesResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListEventSourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListEventSourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListEventSourcesResponse) GoString() string {
	return s.String()
}

func (s *ListEventSourcesResponse) SetHeaders(v map[string]*string) *ListEventSourcesResponse {
	s.Headers = v
	return s
}

func (s *ListEventSourcesResponse) SetStatusCode(v int32) *ListEventSourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListEventSourcesResponse) SetBody(v *ListEventSourcesResponseBody) *ListEventSourcesResponse {
	s.Body = v
	return s
}

type ListFunctionAsyncInvokeConfigsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The CRC-64 value of the function code package. This value is used to check data integrity. The value is automatically calculated by the tool.
	XFcCodeChecksum *string `json:"X-Fc-Code-Checksum,omitempty" xml:"X-Fc-Code-Checksum,omitempty"`
	// The time when the Function Compute is called. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The invocation method.
	//
	// *   **Sync**: synchronous
	// *   **Async**: asynchronous
	XFcInvocationType *string `json:"X-Fc-Invocation-Type,omitempty" xml:"X-Fc-Invocation-Type,omitempty"`
	// The method used to return logs. Valid values:
	//
	// *   **Tail**: returns the last 4 KB of logs that are generated for the current request.
	// *   **None**: No logs are returned for the current request. Default value: None.
	XFcLogType *string `json:"X-Fc-Log-Type,omitempty" xml:"X-Fc-Log-Type,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListFunctionAsyncInvokeConfigsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionAsyncInvokeConfigsHeaders) GoString() string {
	return s.String()
}

func (s *ListFunctionAsyncInvokeConfigsHeaders) SetCommonHeaders(v map[string]*string) *ListFunctionAsyncInvokeConfigsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsHeaders) SetXFcAccountId(v string) *ListFunctionAsyncInvokeConfigsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsHeaders) SetXFcCodeChecksum(v string) *ListFunctionAsyncInvokeConfigsHeaders {
	s.XFcCodeChecksum = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsHeaders) SetXFcDate(v string) *ListFunctionAsyncInvokeConfigsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsHeaders) SetXFcInvocationType(v string) *ListFunctionAsyncInvokeConfigsHeaders {
	s.XFcInvocationType = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsHeaders) SetXFcLogType(v string) *ListFunctionAsyncInvokeConfigsHeaders {
	s.XFcLogType = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsHeaders) SetXFcTraceId(v string) *ListFunctionAsyncInvokeConfigsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListFunctionAsyncInvokeConfigsRequest struct {
	// The maximum number of resources to return.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token required to obtain more results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListFunctionAsyncInvokeConfigsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionAsyncInvokeConfigsRequest) GoString() string {
	return s.String()
}

func (s *ListFunctionAsyncInvokeConfigsRequest) SetLimit(v int32) *ListFunctionAsyncInvokeConfigsRequest {
	s.Limit = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsRequest) SetNextToken(v string) *ListFunctionAsyncInvokeConfigsRequest {
	s.NextToken = &v
	return s
}

type ListFunctionAsyncInvokeConfigsResponseBody struct {
	// The list of asynchronous invocation configurations.
	Configs []*ListFunctionAsyncInvokeConfigsResponseBodyConfigs `json:"configs,omitempty" xml:"configs,omitempty" type:"Repeated"`
	// The token used to obtain more results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListFunctionAsyncInvokeConfigsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionAsyncInvokeConfigsResponseBody) GoString() string {
	return s.String()
}

func (s *ListFunctionAsyncInvokeConfigsResponseBody) SetConfigs(v []*ListFunctionAsyncInvokeConfigsResponseBodyConfigs) *ListFunctionAsyncInvokeConfigsResponseBody {
	s.Configs = v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBody) SetNextToken(v string) *ListFunctionAsyncInvokeConfigsResponseBody {
	s.NextToken = &v
	return s
}

type ListFunctionAsyncInvokeConfigsResponseBodyConfigs struct {
	// The time when the desktop group was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The configuration structure of the destination for asynchronous invocations. If you have not configured this parameter, this parameter is null.
	DestinationConfig *DestinationConfig `json:"destinationConfig,omitempty" xml:"destinationConfig,omitempty"`
	// The name of the function.
	Function *string `json:"function,omitempty" xml:"function,omitempty"`
	// The time when the configuration was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The maximum validity period of a message. If you have not configured this parameter, this parameter is null.
	MaxAsyncEventAgeInSeconds *int64 `json:"maxAsyncEventAgeInSeconds,omitempty" xml:"maxAsyncEventAgeInSeconds,omitempty"`
	// The maximum number of retries allowed after an asynchronous invocation fails. If you have not configured this parameter, this parameter is null.
	MaxAsyncRetryAttempts *int64 `json:"maxAsyncRetryAttempts,omitempty" xml:"maxAsyncRetryAttempts,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The name of the service.
	Service *string `json:"service,omitempty" xml:"service,omitempty"`
	// Indicates whether the asynchronous task feature is enabled.
	//
	// *   **true**: The asynchronous task feature is enabled.
	// *   **false**: The asynchronous task feature is disabled.
	//
	// If you have not configured this parameter, this parameter is null.
	StatefulInvocation *bool `json:"statefulInvocation,omitempty" xml:"statefulInvocation,omitempty"`
}

func (s ListFunctionAsyncInvokeConfigsResponseBodyConfigs) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionAsyncInvokeConfigsResponseBodyConfigs) GoString() string {
	return s.String()
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetCreatedTime(v string) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.CreatedTime = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetDestinationConfig(v *DestinationConfig) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.DestinationConfig = v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetFunction(v string) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.Function = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetLastModifiedTime(v string) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.LastModifiedTime = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetMaxAsyncEventAgeInSeconds(v int64) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.MaxAsyncEventAgeInSeconds = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetMaxAsyncRetryAttempts(v int64) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.MaxAsyncRetryAttempts = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetQualifier(v string) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.Qualifier = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetService(v string) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.Service = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponseBodyConfigs) SetStatefulInvocation(v bool) *ListFunctionAsyncInvokeConfigsResponseBodyConfigs {
	s.StatefulInvocation = &v
	return s
}

type ListFunctionAsyncInvokeConfigsResponse struct {
	Headers    map[string]*string                          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListFunctionAsyncInvokeConfigsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListFunctionAsyncInvokeConfigsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionAsyncInvokeConfigsResponse) GoString() string {
	return s.String()
}

func (s *ListFunctionAsyncInvokeConfigsResponse) SetHeaders(v map[string]*string) *ListFunctionAsyncInvokeConfigsResponse {
	s.Headers = v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponse) SetStatusCode(v int32) *ListFunctionAsyncInvokeConfigsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListFunctionAsyncInvokeConfigsResponse) SetBody(v *ListFunctionAsyncInvokeConfigsResponseBody) *ListFunctionAsyncInvokeConfigsResponse {
	s.Body = v
	return s
}

type ListFunctionsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the function is invoked. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListFunctionsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionsHeaders) GoString() string {
	return s.String()
}

func (s *ListFunctionsHeaders) SetCommonHeaders(v map[string]*string) *ListFunctionsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListFunctionsHeaders) SetXFcAccountId(v string) *ListFunctionsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListFunctionsHeaders) SetXFcDate(v string) *ListFunctionsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListFunctionsHeaders) SetXFcTraceId(v string) *ListFunctionsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListFunctionsRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token required to obtain more results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The prefix that the names of returned resources must contain.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The returned resources are sorted in alphabetical order, and the resources that include and follow the resource specified by the startKey parameter are returned.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListFunctionsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionsRequest) GoString() string {
	return s.String()
}

func (s *ListFunctionsRequest) SetLimit(v int32) *ListFunctionsRequest {
	s.Limit = &v
	return s
}

func (s *ListFunctionsRequest) SetNextToken(v string) *ListFunctionsRequest {
	s.NextToken = &v
	return s
}

func (s *ListFunctionsRequest) SetPrefix(v string) *ListFunctionsRequest {
	s.Prefix = &v
	return s
}

func (s *ListFunctionsRequest) SetQualifier(v string) *ListFunctionsRequest {
	s.Qualifier = &v
	return s
}

func (s *ListFunctionsRequest) SetStartKey(v string) *ListFunctionsRequest {
	s.StartKey = &v
	return s
}

type ListFunctionsResponseBody struct {
	// The information about functions.
	Functions []*ListFunctionsResponseBodyFunctions `json:"functions,omitempty" xml:"functions,omitempty" type:"Repeated"`
	// The token used to obtain more results. If this parameter is left empty, all the results are returned.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListFunctionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionsResponseBody) GoString() string {
	return s.String()
}

func (s *ListFunctionsResponseBody) SetFunctions(v []*ListFunctionsResponseBodyFunctions) *ListFunctionsResponseBody {
	s.Functions = v
	return s
}

func (s *ListFunctionsResponseBody) SetNextToken(v string) *ListFunctionsResponseBody {
	s.NextToken = &v
	return s
}

type ListFunctionsResponseBodyFunctions struct {
	// The port on which the HTTP server listens for the custom runtime or custom container runtime.
	CaPort *int32 `json:"caPort,omitempty" xml:"caPort,omitempty"`
	// The CRC-64 value of the function code package.
	CodeChecksum *string `json:"codeChecksum,omitempty" xml:"codeChecksum,omitempty"`
	// The size of the function code package that is returned by the system. Unit: byte.
	CodeSize *int64 `json:"codeSize,omitempty" xml:"codeSize,omitempty"`
	// The number of vCPUs of the function. The value must be a multiple of 0.05.
	Cpu *float32 `json:"cpu,omitempty" xml:"cpu,omitempty"`
	// The time when the function is created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The configurations of the custom container runtime.
	CustomContainerConfig *CustomContainerConfig `json:"customContainerConfig,omitempty" xml:"customContainerConfig,omitempty"`
	// The custom health check configuration of the function. This parameter is applicable only to custom runtimes and custom containers.
	CustomHealthCheckConfig *CustomHealthCheckConfig `json:"customHealthCheckConfig,omitempty" xml:"customHealthCheckConfig,omitempty"`
	// The description of the function.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The disk size of the function. Unit: MB. Valid values: 512 and 10240.
	DiskSize *int32 `json:"diskSize,omitempty" xml:"diskSize,omitempty"`
	// The environment variables that you configured for the function. You can obtain the values of the environment variables from the function.
	EnvironmentVariables map[string]*string `json:"environmentVariables,omitempty" xml:"environmentVariables,omitempty"`
	// The unique ID that is generated by the system for the function.
	FunctionId *string `json:"functionId,omitempty" xml:"functionId,omitempty"`
	// The name of the function.
	FunctionName *string `json:"functionName,omitempty" xml:"functionName,omitempty"`
	// The GPU memory capacity for the function. Unit: MB. The memory capacity must be a multiple of 1024 MB.
	GpuMemorySize *int32 `json:"gpuMemorySize,omitempty" xml:"gpuMemorySize,omitempty"`
	// The handler of the function.
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	// The timeout period for the execution of the initializer function. Unit: seconds. Default value: 3. Valid values: 1 to 300. When this period ends, the execution of the initializer function is terminated.
	InitializationTimeout *int32 `json:"initializationTimeout,omitempty" xml:"initializationTimeout,omitempty"`
	// The handler of the initializer function. The format of the value is determined by the programming language that you use. For more information, see [Initializer function](~~157704~~).
	Initializer *string `json:"initializer,omitempty" xml:"initializer,omitempty"`
	// The number of requests that can be concurrently processed by a single instance.
	InstanceConcurrency *int32 `json:"instanceConcurrency,omitempty" xml:"instanceConcurrency,omitempty"`
	// The lifecycle configurations of the instance.
	InstanceLifecycleConfig *InstanceLifecycleConfig `json:"instanceLifecycleConfig,omitempty" xml:"instanceLifecycleConfig,omitempty"`
	// The soft concurrency of the instance. You can use this parameter to implement graceful scale-up of instances. If the number of concurrent requests on an instance is greater than the number of the soft concurrency, the instance scale-up is triggered. For example, if your instance requires a long time to start, you can specify a suitable soft concurrency to start the instance in advance.
	//
	// The value must be less than or equal to that of the **instanceConcurrency** parameter.
	InstanceSoftConcurrency *int32 `json:"instanceSoftConcurrency,omitempty" xml:"instanceSoftConcurrency,omitempty"`
	// The instance type of the function. Valid values:
	//
	// *   **e1**: elastic instance
	// *   **c1**: performance instance
	// *   **fc.gpu.tesla.1**: GPU-accelerated instances (Tesla T4)
	// *   **fc.gpu.ampere.1**: GPU-accelerated instances (Ampere A10)
	// *   **g1**: same fc.gpu.tesla.1
	InstanceType *string `json:"instanceType,omitempty" xml:"instanceType,omitempty"`
	// The time when the function was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// An array that consists of the information of layers.
	//
	// > If multiple layers exist, the layers are merged based on the order of array subscripts. The content of a layer with a smaller subscript overwrites the file that has the same name and a larger subscript in the layer.
	Layers []*string `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The memory size that is configured for the function. Unit: MB.
	MemorySize *int32 `json:"memorySize,omitempty" xml:"memorySize,omitempty"`
	// The runtime environment of the function. Valid values: **nodejs16**, **nodejs14**, **nodejs12**, **nodejs10**, **nodejs8**, **nodejs6**, **nodejs4.4**, **python3.9**, **python3**, **python2.7**, **java11**, **java8**, **go1**, **php7.2**, **dotnetcore3.1**, **dotnetcore2.1**, **custom** and **custom-container**. For more information, see [Supported function runtime environments](~~73338~~).
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The timeout period for the execution of the function. Unit: seconds. Default value: 60. Valid values: 1 to 600. When this period expires, the execution of the function is terminated.
	Timeout *int32 `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s ListFunctionsResponseBodyFunctions) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionsResponseBodyFunctions) GoString() string {
	return s.String()
}

func (s *ListFunctionsResponseBodyFunctions) SetCaPort(v int32) *ListFunctionsResponseBodyFunctions {
	s.CaPort = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetCodeChecksum(v string) *ListFunctionsResponseBodyFunctions {
	s.CodeChecksum = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetCodeSize(v int64) *ListFunctionsResponseBodyFunctions {
	s.CodeSize = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetCpu(v float32) *ListFunctionsResponseBodyFunctions {
	s.Cpu = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetCreatedTime(v string) *ListFunctionsResponseBodyFunctions {
	s.CreatedTime = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetCustomContainerConfig(v *CustomContainerConfig) *ListFunctionsResponseBodyFunctions {
	s.CustomContainerConfig = v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetCustomHealthCheckConfig(v *CustomHealthCheckConfig) *ListFunctionsResponseBodyFunctions {
	s.CustomHealthCheckConfig = v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetDescription(v string) *ListFunctionsResponseBodyFunctions {
	s.Description = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetDiskSize(v int32) *ListFunctionsResponseBodyFunctions {
	s.DiskSize = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetEnvironmentVariables(v map[string]*string) *ListFunctionsResponseBodyFunctions {
	s.EnvironmentVariables = v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetFunctionId(v string) *ListFunctionsResponseBodyFunctions {
	s.FunctionId = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetFunctionName(v string) *ListFunctionsResponseBodyFunctions {
	s.FunctionName = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetGpuMemorySize(v int32) *ListFunctionsResponseBodyFunctions {
	s.GpuMemorySize = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetHandler(v string) *ListFunctionsResponseBodyFunctions {
	s.Handler = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetInitializationTimeout(v int32) *ListFunctionsResponseBodyFunctions {
	s.InitializationTimeout = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetInitializer(v string) *ListFunctionsResponseBodyFunctions {
	s.Initializer = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetInstanceConcurrency(v int32) *ListFunctionsResponseBodyFunctions {
	s.InstanceConcurrency = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetInstanceLifecycleConfig(v *InstanceLifecycleConfig) *ListFunctionsResponseBodyFunctions {
	s.InstanceLifecycleConfig = v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetInstanceSoftConcurrency(v int32) *ListFunctionsResponseBodyFunctions {
	s.InstanceSoftConcurrency = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetInstanceType(v string) *ListFunctionsResponseBodyFunctions {
	s.InstanceType = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetLastModifiedTime(v string) *ListFunctionsResponseBodyFunctions {
	s.LastModifiedTime = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetLayers(v []*string) *ListFunctionsResponseBodyFunctions {
	s.Layers = v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetMemorySize(v int32) *ListFunctionsResponseBodyFunctions {
	s.MemorySize = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetRuntime(v string) *ListFunctionsResponseBodyFunctions {
	s.Runtime = &v
	return s
}

func (s *ListFunctionsResponseBodyFunctions) SetTimeout(v int32) *ListFunctionsResponseBodyFunctions {
	s.Timeout = &v
	return s
}

type ListFunctionsResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListFunctionsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListFunctionsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListFunctionsResponse) GoString() string {
	return s.String()
}

func (s *ListFunctionsResponse) SetHeaders(v map[string]*string) *ListFunctionsResponse {
	s.Headers = v
	return s
}

func (s *ListFunctionsResponse) SetStatusCode(v int32) *ListFunctionsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListFunctionsResponse) SetBody(v *ListFunctionsResponseBody) *ListFunctionsResponse {
	s.Body = v
	return s
}

type ListInstancesHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	XFcDate      *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	XFcTraceId   *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListInstancesHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListInstancesHeaders) GoString() string {
	return s.String()
}

func (s *ListInstancesHeaders) SetCommonHeaders(v map[string]*string) *ListInstancesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListInstancesHeaders) SetXFcAccountId(v string) *ListInstancesHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListInstancesHeaders) SetXFcDate(v string) *ListInstancesHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListInstancesHeaders) SetXFcTraceId(v string) *ListInstancesHeaders {
	s.XFcTraceId = &v
	return s
}

type ListInstancesRequest struct {
	// The IDs of the instance.
	InstanceIds []*string `json:"instanceIds,omitempty" xml:"instanceIds,omitempty" type:"Repeated"`
	// The maximum number of resources to return. Valid values: \[0,1000].
	//
	// The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The version or alias.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s ListInstancesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListInstancesRequest) GoString() string {
	return s.String()
}

func (s *ListInstancesRequest) SetInstanceIds(v []*string) *ListInstancesRequest {
	s.InstanceIds = v
	return s
}

func (s *ListInstancesRequest) SetLimit(v int32) *ListInstancesRequest {
	s.Limit = &v
	return s
}

func (s *ListInstancesRequest) SetQualifier(v string) *ListInstancesRequest {
	s.Qualifier = &v
	return s
}

type ListInstancesResponseBody struct {
	// The information about instances.
	Instances []*ListInstancesResponseBodyInstances `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
}

func (s ListInstancesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListInstancesResponseBody) GoString() string {
	return s.String()
}

func (s *ListInstancesResponseBody) SetInstances(v []*ListInstancesResponseBodyInstances) *ListInstancesResponseBody {
	s.Instances = v
	return s
}

type ListInstancesResponseBodyInstances struct {
	// The ID of the instance.
	InstanceId *string `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	// The version of the service to which the instance belongs. If the instance belongs to the LATEST alias, 0 is returned as the version.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s ListInstancesResponseBodyInstances) String() string {
	return tea.Prettify(s)
}

func (s ListInstancesResponseBodyInstances) GoString() string {
	return s.String()
}

func (s *ListInstancesResponseBodyInstances) SetInstanceId(v string) *ListInstancesResponseBodyInstances {
	s.InstanceId = &v
	return s
}

func (s *ListInstancesResponseBodyInstances) SetVersionId(v string) *ListInstancesResponseBodyInstances {
	s.VersionId = &v
	return s
}

type ListInstancesResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListInstancesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListInstancesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListInstancesResponse) GoString() string {
	return s.String()
}

func (s *ListInstancesResponse) SetHeaders(v map[string]*string) *ListInstancesResponse {
	s.Headers = v
	return s
}

func (s *ListInstancesResponse) SetStatusCode(v int32) *ListInstancesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListInstancesResponse) SetBody(v *ListInstancesResponseBody) *ListInstancesResponse {
	s.Body = v
	return s
}

type ListLayerVersionsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListLayerVersionsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListLayerVersionsHeaders) GoString() string {
	return s.String()
}

func (s *ListLayerVersionsHeaders) SetCommonHeaders(v map[string]*string) *ListLayerVersionsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListLayerVersionsHeaders) SetXFcAccountId(v string) *ListLayerVersionsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListLayerVersionsHeaders) SetXFcDate(v string) *ListLayerVersionsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListLayerVersionsHeaders) SetXFcTraceId(v string) *ListLayerVersionsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListLayerVersionsRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The initial version of the layer.
	StartVersion *int32 `json:"startVersion,omitempty" xml:"startVersion,omitempty"`
}

func (s ListLayerVersionsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListLayerVersionsRequest) GoString() string {
	return s.String()
}

func (s *ListLayerVersionsRequest) SetLimit(v int32) *ListLayerVersionsRequest {
	s.Limit = &v
	return s
}

func (s *ListLayerVersionsRequest) SetStartVersion(v int32) *ListLayerVersionsRequest {
	s.StartVersion = &v
	return s
}

type ListLayerVersionsResponseBody struct {
	// The information about layer versions.
	Layers []*Layer `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The initial version of the layer for the next query.
	NextVersion *int32 `json:"nextVersion,omitempty" xml:"nextVersion,omitempty"`
}

func (s ListLayerVersionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListLayerVersionsResponseBody) GoString() string {
	return s.String()
}

func (s *ListLayerVersionsResponseBody) SetLayers(v []*Layer) *ListLayerVersionsResponseBody {
	s.Layers = v
	return s
}

func (s *ListLayerVersionsResponseBody) SetNextVersion(v int32) *ListLayerVersionsResponseBody {
	s.NextVersion = &v
	return s
}

type ListLayerVersionsResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListLayerVersionsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListLayerVersionsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListLayerVersionsResponse) GoString() string {
	return s.String()
}

func (s *ListLayerVersionsResponse) SetHeaders(v map[string]*string) *ListLayerVersionsResponse {
	s.Headers = v
	return s
}

func (s *ListLayerVersionsResponse) SetStatusCode(v int32) *ListLayerVersionsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListLayerVersionsResponse) SetBody(v *ListLayerVersionsResponseBody) *ListLayerVersionsResponse {
	s.Body = v
	return s
}

type ListLayersHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the function is invoked. The format is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListLayersHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListLayersHeaders) GoString() string {
	return s.String()
}

func (s *ListLayersHeaders) SetCommonHeaders(v map[string]*string) *ListLayersHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListLayersHeaders) SetXFcAccountId(v string) *ListLayersHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListLayersHeaders) SetXFcDate(v string) *ListLayersHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListLayersHeaders) SetXFcTraceId(v string) *ListLayersHeaders {
	s.XFcTraceId = &v
	return s
}

type ListLayersRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned configurations is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token required to obtain more results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// Specifies whether to obtain the official public layer. When the official parameter is set to true, the public field does not take effect. The default value is false.
	Official *bool `json:"official,omitempty" xml:"official,omitempty"`
	// The name prefix of the layer. The names of returned resources must contain the prefix. If the name prefix is a, the names of returned resources must start with a.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// Specifies whether to obtain only the common layer. Default value: false.
	Public *bool `json:"public,omitempty" xml:"public,omitempty"`
	// The name of the start layer. The returned layers are sorted in alphabetical order, and the layers that include and follow the layer specified by the startKey parameter are returned.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListLayersRequest) String() string {
	return tea.Prettify(s)
}

func (s ListLayersRequest) GoString() string {
	return s.String()
}

func (s *ListLayersRequest) SetLimit(v int32) *ListLayersRequest {
	s.Limit = &v
	return s
}

func (s *ListLayersRequest) SetNextToken(v string) *ListLayersRequest {
	s.NextToken = &v
	return s
}

func (s *ListLayersRequest) SetOfficial(v bool) *ListLayersRequest {
	s.Official = &v
	return s
}

func (s *ListLayersRequest) SetPrefix(v string) *ListLayersRequest {
	s.Prefix = &v
	return s
}

func (s *ListLayersRequest) SetPublic(v bool) *ListLayersRequest {
	s.Public = &v
	return s
}

func (s *ListLayersRequest) SetStartKey(v string) *ListLayersRequest {
	s.StartKey = &v
	return s
}

type ListLayersResponseBody struct {
	// The information about layers.
	Layers []*Layer `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The name of the start layer for the next query, which is also the token used to obtain more results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListLayersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListLayersResponseBody) GoString() string {
	return s.String()
}

func (s *ListLayersResponseBody) SetLayers(v []*Layer) *ListLayersResponseBody {
	s.Layers = v
	return s
}

func (s *ListLayersResponseBody) SetNextToken(v string) *ListLayersResponseBody {
	s.NextToken = &v
	return s
}

type ListLayersResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListLayersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListLayersResponse) String() string {
	return tea.Prettify(s)
}

func (s ListLayersResponse) GoString() string {
	return s.String()
}

func (s *ListLayersResponse) SetHeaders(v map[string]*string) *ListLayersResponse {
	s.Headers = v
	return s
}

func (s *ListLayersResponse) SetStatusCode(v int32) *ListLayersResponse {
	s.StatusCode = &v
	return s
}

func (s *ListLayersResponse) SetBody(v *ListLayersResponseBody) *ListLayersResponse {
	s.Body = v
	return s
}

type ListOnDemandConfigsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListOnDemandConfigsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListOnDemandConfigsHeaders) GoString() string {
	return s.String()
}

func (s *ListOnDemandConfigsHeaders) SetCommonHeaders(v map[string]*string) *ListOnDemandConfigsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListOnDemandConfigsHeaders) SetXFcAccountId(v string) *ListOnDemandConfigsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListOnDemandConfigsHeaders) SetXFcDate(v string) *ListOnDemandConfigsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListOnDemandConfigsHeaders) SetXFcTraceId(v string) *ListOnDemandConfigsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListOnDemandConfigsRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token used to obtain more results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The prefix that the names of returned resources must contain. If the name prefix is a, the names of returned resources must start with a.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The returned resources are sorted in alphabetical order, and the resources that include and follow the resource specified by the startKey parameter are returned.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListOnDemandConfigsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListOnDemandConfigsRequest) GoString() string {
	return s.String()
}

func (s *ListOnDemandConfigsRequest) SetLimit(v int32) *ListOnDemandConfigsRequest {
	s.Limit = &v
	return s
}

func (s *ListOnDemandConfigsRequest) SetNextToken(v string) *ListOnDemandConfigsRequest {
	s.NextToken = &v
	return s
}

func (s *ListOnDemandConfigsRequest) SetPrefix(v string) *ListOnDemandConfigsRequest {
	s.Prefix = &v
	return s
}

func (s *ListOnDemandConfigsRequest) SetStartKey(v string) *ListOnDemandConfigsRequest {
	s.StartKey = &v
	return s
}

type ListOnDemandConfigsResponseBody struct {
	// The information about the provisioned configuration.
	Configs []*OnDemandConfig `json:"configs,omitempty" xml:"configs,omitempty" type:"Repeated"`
	// The token used to obtain more results. If this parameter is left empty, all the results are returned.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListOnDemandConfigsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListOnDemandConfigsResponseBody) GoString() string {
	return s.String()
}

func (s *ListOnDemandConfigsResponseBody) SetConfigs(v []*OnDemandConfig) *ListOnDemandConfigsResponseBody {
	s.Configs = v
	return s
}

func (s *ListOnDemandConfigsResponseBody) SetNextToken(v string) *ListOnDemandConfigsResponseBody {
	s.NextToken = &v
	return s
}

type ListOnDemandConfigsResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListOnDemandConfigsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListOnDemandConfigsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListOnDemandConfigsResponse) GoString() string {
	return s.String()
}

func (s *ListOnDemandConfigsResponse) SetHeaders(v map[string]*string) *ListOnDemandConfigsResponse {
	s.Headers = v
	return s
}

func (s *ListOnDemandConfigsResponse) SetStatusCode(v int32) *ListOnDemandConfigsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListOnDemandConfigsResponse) SetBody(v *ListOnDemandConfigsResponseBody) *ListOnDemandConfigsResponse {
	s.Body = v
	return s
}

type ListProvisionConfigsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListProvisionConfigsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListProvisionConfigsHeaders) GoString() string {
	return s.String()
}

func (s *ListProvisionConfigsHeaders) SetCommonHeaders(v map[string]*string) *ListProvisionConfigsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListProvisionConfigsHeaders) SetXFcAccountId(v string) *ListProvisionConfigsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListProvisionConfigsHeaders) SetXFcDate(v string) *ListProvisionConfigsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListProvisionConfigsHeaders) SetXFcTraceId(v string) *ListProvisionConfigsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListProvisionConfigsRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int64 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token used to obtain more results. You do not need to provide this parameter in the first call. The tokens for subsequent queries are obtained from the returned results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The qualifier of the service to which resources belong. The qualifier must be aliasName and used together with the serviceName parameter.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The name of the service to which resources belong.
	ServiceName *string `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
}

func (s ListProvisionConfigsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListProvisionConfigsRequest) GoString() string {
	return s.String()
}

func (s *ListProvisionConfigsRequest) SetLimit(v int64) *ListProvisionConfigsRequest {
	s.Limit = &v
	return s
}

func (s *ListProvisionConfigsRequest) SetNextToken(v string) *ListProvisionConfigsRequest {
	s.NextToken = &v
	return s
}

func (s *ListProvisionConfigsRequest) SetQualifier(v string) *ListProvisionConfigsRequest {
	s.Qualifier = &v
	return s
}

func (s *ListProvisionConfigsRequest) SetServiceName(v string) *ListProvisionConfigsRequest {
	s.ServiceName = &v
	return s
}

type ListProvisionConfigsResponseBody struct {
	// The token used to obtain more results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The information about provisioned instances.
	ProvisionConfigs []*ListProvisionConfigsResponseBodyProvisionConfigs `json:"provisionConfigs,omitempty" xml:"provisionConfigs,omitempty" type:"Repeated"`
}

func (s ListProvisionConfigsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListProvisionConfigsResponseBody) GoString() string {
	return s.String()
}

func (s *ListProvisionConfigsResponseBody) SetNextToken(v string) *ListProvisionConfigsResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListProvisionConfigsResponseBody) SetProvisionConfigs(v []*ListProvisionConfigsResponseBodyProvisionConfigs) *ListProvisionConfigsResponseBody {
	s.ProvisionConfigs = v
	return s
}

type ListProvisionConfigsResponseBodyProvisionConfigs struct {
	// Specifies whether to always allocate CPU to a function instance.
	AlwaysAllocateCPU *bool `json:"alwaysAllocateCPU,omitempty" xml:"alwaysAllocateCPU,omitempty"`
	// The actual number of provisioned instances.
	Current *int64 `json:"current,omitempty" xml:"current,omitempty"`
	// The error message returned if a provisioned instance fails to be created.
	CurrentError *string `json:"currentError,omitempty" xml:"currentError,omitempty"`
	// The description of the resource.
	Resource *string `json:"resource,omitempty" xml:"resource,omitempty"`
	// The configurations of scheduled auto scaling.
	ScheduledActions []*ScheduledActions `json:"scheduledActions,omitempty" xml:"scheduledActions,omitempty" type:"Repeated"`
	// The expected number of provisioned instances.
	Target *int64 `json:"target,omitempty" xml:"target,omitempty"`
	// The configurations of metric-based auto scaling.
	TargetTrackingPolicies []*TargetTrackingPolicies `json:"targetTrackingPolicies,omitempty" xml:"targetTrackingPolicies,omitempty" type:"Repeated"`
}

func (s ListProvisionConfigsResponseBodyProvisionConfigs) String() string {
	return tea.Prettify(s)
}

func (s ListProvisionConfigsResponseBodyProvisionConfigs) GoString() string {
	return s.String()
}

func (s *ListProvisionConfigsResponseBodyProvisionConfigs) SetAlwaysAllocateCPU(v bool) *ListProvisionConfigsResponseBodyProvisionConfigs {
	s.AlwaysAllocateCPU = &v
	return s
}

func (s *ListProvisionConfigsResponseBodyProvisionConfigs) SetCurrent(v int64) *ListProvisionConfigsResponseBodyProvisionConfigs {
	s.Current = &v
	return s
}

func (s *ListProvisionConfigsResponseBodyProvisionConfigs) SetCurrentError(v string) *ListProvisionConfigsResponseBodyProvisionConfigs {
	s.CurrentError = &v
	return s
}

func (s *ListProvisionConfigsResponseBodyProvisionConfigs) SetResource(v string) *ListProvisionConfigsResponseBodyProvisionConfigs {
	s.Resource = &v
	return s
}

func (s *ListProvisionConfigsResponseBodyProvisionConfigs) SetScheduledActions(v []*ScheduledActions) *ListProvisionConfigsResponseBodyProvisionConfigs {
	s.ScheduledActions = v
	return s
}

func (s *ListProvisionConfigsResponseBodyProvisionConfigs) SetTarget(v int64) *ListProvisionConfigsResponseBodyProvisionConfigs {
	s.Target = &v
	return s
}

func (s *ListProvisionConfigsResponseBodyProvisionConfigs) SetTargetTrackingPolicies(v []*TargetTrackingPolicies) *ListProvisionConfigsResponseBodyProvisionConfigs {
	s.TargetTrackingPolicies = v
	return s
}

type ListProvisionConfigsResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListProvisionConfigsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListProvisionConfigsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListProvisionConfigsResponse) GoString() string {
	return s.String()
}

func (s *ListProvisionConfigsResponse) SetHeaders(v map[string]*string) *ListProvisionConfigsResponse {
	s.Headers = v
	return s
}

func (s *ListProvisionConfigsResponse) SetStatusCode(v int32) *ListProvisionConfigsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListProvisionConfigsResponse) SetBody(v *ListProvisionConfigsResponseBody) *ListProvisionConfigsResponse {
	s.Body = v
	return s
}

type ListReservedCapacitiesHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the Function Compute API is called. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListReservedCapacitiesHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListReservedCapacitiesHeaders) GoString() string {
	return s.String()
}

func (s *ListReservedCapacitiesHeaders) SetCommonHeaders(v map[string]*string) *ListReservedCapacitiesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListReservedCapacitiesHeaders) SetXFcAccountId(v string) *ListReservedCapacitiesHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListReservedCapacitiesHeaders) SetXFcDate(v string) *ListReservedCapacitiesHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListReservedCapacitiesHeaders) SetXFcTraceId(v string) *ListReservedCapacitiesHeaders {
	s.XFcTraceId = &v
	return s
}

type ListReservedCapacitiesRequest struct {
	// The maximum number of resources to return. Valid values: \[1, 100].
	Limit *string `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token that determines the start point of the query.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListReservedCapacitiesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListReservedCapacitiesRequest) GoString() string {
	return s.String()
}

func (s *ListReservedCapacitiesRequest) SetLimit(v string) *ListReservedCapacitiesRequest {
	s.Limit = &v
	return s
}

func (s *ListReservedCapacitiesRequest) SetNextToken(v string) *ListReservedCapacitiesRequest {
	s.NextToken = &v
	return s
}

type ListReservedCapacitiesResponseBody struct {
	// The token used to obtain more results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The information about subscription instances.
	ReservedCapacities []*OpenReservedCapacity `json:"reservedCapacities,omitempty" xml:"reservedCapacities,omitempty" type:"Repeated"`
}

func (s ListReservedCapacitiesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListReservedCapacitiesResponseBody) GoString() string {
	return s.String()
}

func (s *ListReservedCapacitiesResponseBody) SetNextToken(v string) *ListReservedCapacitiesResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListReservedCapacitiesResponseBody) SetReservedCapacities(v []*OpenReservedCapacity) *ListReservedCapacitiesResponseBody {
	s.ReservedCapacities = v
	return s
}

type ListReservedCapacitiesResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListReservedCapacitiesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListReservedCapacitiesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListReservedCapacitiesResponse) GoString() string {
	return s.String()
}

func (s *ListReservedCapacitiesResponse) SetHeaders(v map[string]*string) *ListReservedCapacitiesResponse {
	s.Headers = v
	return s
}

func (s *ListReservedCapacitiesResponse) SetStatusCode(v int32) *ListReservedCapacitiesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListReservedCapacitiesResponse) SetBody(v *ListReservedCapacitiesResponseBody) *ListReservedCapacitiesResponse {
	s.Body = v
	return s
}

type ListServiceVersionsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListServiceVersionsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListServiceVersionsHeaders) GoString() string {
	return s.String()
}

func (s *ListServiceVersionsHeaders) SetCommonHeaders(v map[string]*string) *ListServiceVersionsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListServiceVersionsHeaders) SetXFcAccountId(v string) *ListServiceVersionsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListServiceVersionsHeaders) SetXFcDate(v string) *ListServiceVersionsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListServiceVersionsHeaders) SetXFcTraceId(v string) *ListServiceVersionsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListServiceVersionsRequest struct {
	// The order in which the returned versions are sorted. Valid values:
	//   - **FORWARD**: in ascending order.
	//   - **BACKWARD**: in descending order. This is the default value.
	Direction *string `json:"direction,omitempty" xml:"direction,omitempty"`
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token used to obtain more results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The starting position of the result list. The returned resources are sorted based on the version number, and the resources that include and follow the resource specified by the startKey parameter are returned.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListServiceVersionsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListServiceVersionsRequest) GoString() string {
	return s.String()
}

func (s *ListServiceVersionsRequest) SetDirection(v string) *ListServiceVersionsRequest {
	s.Direction = &v
	return s
}

func (s *ListServiceVersionsRequest) SetLimit(v int32) *ListServiceVersionsRequest {
	s.Limit = &v
	return s
}

func (s *ListServiceVersionsRequest) SetNextToken(v string) *ListServiceVersionsRequest {
	s.NextToken = &v
	return s
}

func (s *ListServiceVersionsRequest) SetStartKey(v string) *ListServiceVersionsRequest {
	s.StartKey = &v
	return s
}

type ListServiceVersionsResponseBody struct {
	// The order in which the returned versions are sorted. Valid values:
	//   - **FORWARD**: in ascending order.
	//   - **BACKWARD**: in descending order. This is the default value.
	Direction *string `json:"direction,omitempty" xml:"direction,omitempty"`
	// The token used to obtain more results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The list of versions.
	Versions []*ListServiceVersionsResponseBodyVersions `json:"versions,omitempty" xml:"versions,omitempty" type:"Repeated"`
}

func (s ListServiceVersionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListServiceVersionsResponseBody) GoString() string {
	return s.String()
}

func (s *ListServiceVersionsResponseBody) SetDirection(v string) *ListServiceVersionsResponseBody {
	s.Direction = &v
	return s
}

func (s *ListServiceVersionsResponseBody) SetNextToken(v string) *ListServiceVersionsResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListServiceVersionsResponseBody) SetVersions(v []*ListServiceVersionsResponseBodyVersions) *ListServiceVersionsResponseBody {
	s.Versions = v
	return s
}

type ListServiceVersionsResponseBodyVersions struct {
	// The time when the service version was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the service version.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The time when the service version was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The version of the service.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s ListServiceVersionsResponseBodyVersions) String() string {
	return tea.Prettify(s)
}

func (s ListServiceVersionsResponseBodyVersions) GoString() string {
	return s.String()
}

func (s *ListServiceVersionsResponseBodyVersions) SetCreatedTime(v string) *ListServiceVersionsResponseBodyVersions {
	s.CreatedTime = &v
	return s
}

func (s *ListServiceVersionsResponseBodyVersions) SetDescription(v string) *ListServiceVersionsResponseBodyVersions {
	s.Description = &v
	return s
}

func (s *ListServiceVersionsResponseBodyVersions) SetLastModifiedTime(v string) *ListServiceVersionsResponseBodyVersions {
	s.LastModifiedTime = &v
	return s
}

func (s *ListServiceVersionsResponseBodyVersions) SetVersionId(v string) *ListServiceVersionsResponseBodyVersions {
	s.VersionId = &v
	return s
}

type ListServiceVersionsResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListServiceVersionsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListServiceVersionsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListServiceVersionsResponse) GoString() string {
	return s.String()
}

func (s *ListServiceVersionsResponse) SetHeaders(v map[string]*string) *ListServiceVersionsResponse {
	s.Headers = v
	return s
}

func (s *ListServiceVersionsResponse) SetStatusCode(v int32) *ListServiceVersionsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListServiceVersionsResponse) SetBody(v *ListServiceVersionsResponseBody) *ListServiceVersionsResponse {
	s.Body = v
	return s
}

type ListServicesHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListServicesHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListServicesHeaders) GoString() string {
	return s.String()
}

func (s *ListServicesHeaders) SetCommonHeaders(v map[string]*string) *ListServicesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListServicesHeaders) SetXFcAccountId(v string) *ListServicesHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListServicesHeaders) SetXFcDate(v string) *ListServicesHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListServicesHeaders) SetXFcTraceId(v string) *ListServicesHeaders {
	s.XFcTraceId = &v
	return s
}

type ListServicesRequest struct {
	// The maximum number of resources to return. Default value: 20. The value cannot exceed 100. The number of returned configurations is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The starting position of the query. If this parameter is left empty, the query starts from the beginning. You do not need to specify this parameter in the first query. If the number of asynchronous tasks exceeds the limit, the nextToken parameter is returned, the value of which can be used in subsequent calls to obtain more results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The prefix that the names of returned resources must contain. If the name prefix is a, the names of returned resources must start with a.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The returned resources are sorted in alphabetical order, and the resources that include and follow the resource specified by the startKey parameter are returned.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListServicesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListServicesRequest) GoString() string {
	return s.String()
}

func (s *ListServicesRequest) SetLimit(v int32) *ListServicesRequest {
	s.Limit = &v
	return s
}

func (s *ListServicesRequest) SetNextToken(v string) *ListServicesRequest {
	s.NextToken = &v
	return s
}

func (s *ListServicesRequest) SetPrefix(v string) *ListServicesRequest {
	s.Prefix = &v
	return s
}

func (s *ListServicesRequest) SetStartKey(v string) *ListServicesRequest {
	s.StartKey = &v
	return s
}

type ListServicesResponseBody struct {
	// The token used to obtain more results. If this parameter is left empty, all the results are returned.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The information about a service.
	Services []*ListServicesResponseBodyServices `json:"services,omitempty" xml:"services,omitempty" type:"Repeated"`
}

func (s ListServicesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListServicesResponseBody) GoString() string {
	return s.String()
}

func (s *ListServicesResponseBody) SetNextToken(v string) *ListServicesResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListServicesResponseBody) SetServices(v []*ListServicesResponseBodyServices) *ListServicesResponseBody {
	s.Services = v
	return s
}

type ListServicesResponseBodyServices struct {
	// The time when the service was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the service.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Specifies whether to allow functions to access the Internet. Valid values:
	//
	// *   **true**: allows functions in the specified service to access the Internet.
	// *   **false**: does not allow functions to access the Internet.
	InternetAccess *bool `json:"internetAccess,omitempty" xml:"internetAccess,omitempty"`
	// The time when the service was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The log configuration, which specifies a Logstore to store function execution logs.
	LogConfig *LogConfig `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	// The configurations of the NAS file system. The configuration allows functions in the specified service in Function Compute to access the NAS file system.
	NasConfig *NASConfig `json:"nasConfig,omitempty" xml:"nasConfig,omitempty"`
	// The OSS mount configurations.
	OssMountConfig *OSSMountConfig `json:"ossMountConfig,omitempty" xml:"ossMountConfig,omitempty"`
	// The RAM role that is used to grant required permissions to Function Compute. The RAM role is used in the following scenarios:
	//
	// *   Sends function execution logs to your Logstore.
	// *   Generates a token for a function to access other cloud resources during function execution.
	Role *string `json:"role,omitempty" xml:"role,omitempty"`
	// The unique ID generated by the system for the service.
	ServiceId *string `json:"serviceId,omitempty" xml:"serviceId,omitempty"`
	// The name of the service.
	ServiceName *string `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
	// The configuration of Tracing Analysis. After you configure Tracing Analysis for a service in Function Compute, you can record the execution duration of a request, view the amount of cold start time for a function, and record the execution duration of a function. For more information, see [Overview](~~189804~~).
	TracingConfig *TracingConfig `json:"tracingConfig,omitempty" xml:"tracingConfig,omitempty"`
	// The VPC configuration. The configuration allows a function to access the specified VPC.
	VpcConfig *VPCConfig `json:"vpcConfig,omitempty" xml:"vpcConfig,omitempty"`
}

func (s ListServicesResponseBodyServices) String() string {
	return tea.Prettify(s)
}

func (s ListServicesResponseBodyServices) GoString() string {
	return s.String()
}

func (s *ListServicesResponseBodyServices) SetCreatedTime(v string) *ListServicesResponseBodyServices {
	s.CreatedTime = &v
	return s
}

func (s *ListServicesResponseBodyServices) SetDescription(v string) *ListServicesResponseBodyServices {
	s.Description = &v
	return s
}

func (s *ListServicesResponseBodyServices) SetInternetAccess(v bool) *ListServicesResponseBodyServices {
	s.InternetAccess = &v
	return s
}

func (s *ListServicesResponseBodyServices) SetLastModifiedTime(v string) *ListServicesResponseBodyServices {
	s.LastModifiedTime = &v
	return s
}

func (s *ListServicesResponseBodyServices) SetLogConfig(v *LogConfig) *ListServicesResponseBodyServices {
	s.LogConfig = v
	return s
}

func (s *ListServicesResponseBodyServices) SetNasConfig(v *NASConfig) *ListServicesResponseBodyServices {
	s.NasConfig = v
	return s
}

func (s *ListServicesResponseBodyServices) SetOssMountConfig(v *OSSMountConfig) *ListServicesResponseBodyServices {
	s.OssMountConfig = v
	return s
}

func (s *ListServicesResponseBodyServices) SetRole(v string) *ListServicesResponseBodyServices {
	s.Role = &v
	return s
}

func (s *ListServicesResponseBodyServices) SetServiceId(v string) *ListServicesResponseBodyServices {
	s.ServiceId = &v
	return s
}

func (s *ListServicesResponseBodyServices) SetServiceName(v string) *ListServicesResponseBodyServices {
	s.ServiceName = &v
	return s
}

func (s *ListServicesResponseBodyServices) SetTracingConfig(v *TracingConfig) *ListServicesResponseBodyServices {
	s.TracingConfig = v
	return s
}

func (s *ListServicesResponseBodyServices) SetVpcConfig(v *VPCConfig) *ListServicesResponseBodyServices {
	s.VpcConfig = v
	return s
}

type ListServicesResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListServicesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListServicesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListServicesResponse) GoString() string {
	return s.String()
}

func (s *ListServicesResponse) SetHeaders(v map[string]*string) *ListServicesResponse {
	s.Headers = v
	return s
}

func (s *ListServicesResponse) SetStatusCode(v int32) *ListServicesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListServicesResponse) SetBody(v *ListServicesResponseBody) *ListServicesResponse {
	s.Body = v
	return s
}

type ListStatefulAsyncInvocationFunctionsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListStatefulAsyncInvocationFunctionsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationFunctionsHeaders) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationFunctionsHeaders) SetCommonHeaders(v map[string]*string) *ListStatefulAsyncInvocationFunctionsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListStatefulAsyncInvocationFunctionsHeaders) SetXFcAccountId(v string) *ListStatefulAsyncInvocationFunctionsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListStatefulAsyncInvocationFunctionsHeaders) SetXFcDate(v string) *ListStatefulAsyncInvocationFunctionsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListStatefulAsyncInvocationFunctionsHeaders) SetXFcTraceId(v string) *ListStatefulAsyncInvocationFunctionsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListStatefulAsyncInvocationFunctionsRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The starting position of the query. If this parameter is left empty, the query starts from the beginning. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListStatefulAsyncInvocationFunctionsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationFunctionsRequest) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationFunctionsRequest) SetLimit(v int32) *ListStatefulAsyncInvocationFunctionsRequest {
	s.Limit = &v
	return s
}

func (s *ListStatefulAsyncInvocationFunctionsRequest) SetNextToken(v string) *ListStatefulAsyncInvocationFunctionsRequest {
	s.NextToken = &v
	return s
}

type ListStatefulAsyncInvocationFunctionsResponseBody struct {
	// The details of returned data.
	Data []*AsyncConfigMeta `json:"data,omitempty" xml:"data,omitempty" type:"Repeated"`
	// The token used to obtain more results. If this parameter is left empty, all the results are returned.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListStatefulAsyncInvocationFunctionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationFunctionsResponseBody) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationFunctionsResponseBody) SetData(v []*AsyncConfigMeta) *ListStatefulAsyncInvocationFunctionsResponseBody {
	s.Data = v
	return s
}

func (s *ListStatefulAsyncInvocationFunctionsResponseBody) SetNextToken(v string) *ListStatefulAsyncInvocationFunctionsResponseBody {
	s.NextToken = &v
	return s
}

type ListStatefulAsyncInvocationFunctionsResponse struct {
	Headers    map[string]*string                                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListStatefulAsyncInvocationFunctionsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListStatefulAsyncInvocationFunctionsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationFunctionsResponse) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationFunctionsResponse) SetHeaders(v map[string]*string) *ListStatefulAsyncInvocationFunctionsResponse {
	s.Headers = v
	return s
}

func (s *ListStatefulAsyncInvocationFunctionsResponse) SetStatusCode(v int32) *ListStatefulAsyncInvocationFunctionsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListStatefulAsyncInvocationFunctionsResponse) SetBody(v *ListStatefulAsyncInvocationFunctionsResponseBody) *ListStatefulAsyncInvocationFunctionsResponse {
	s.Body = v
	return s
}

type ListStatefulAsyncInvocationsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The CRC-64 value of the function code package. This value is used to check data integrity. The value is automatically calculated by the tool.
	XFcCodeChecksum *string `json:"X-Fc-Code-Checksum,omitempty" xml:"X-Fc-Code-Checksum,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The invocation method.
	//
	// - **Sync**: synchronous invocation
	// - **Async**: asynchronous invocation
	XFcInvocationType *string `json:"X-Fc-Invocation-Type,omitempty" xml:"X-Fc-Invocation-Type,omitempty"`
	// The method used to return logs. Valid values:
	//
	// - **Tail**: returns the last 4 KB of logs that are generated for the current request.
	// - **None**: does not return logs for the current request. This is the default value.
	XFcLogType *string `json:"X-Fc-Log-Type,omitempty" xml:"X-Fc-Log-Type,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListStatefulAsyncInvocationsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationsHeaders) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationsHeaders) SetCommonHeaders(v map[string]*string) *ListStatefulAsyncInvocationsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListStatefulAsyncInvocationsHeaders) SetXFcAccountId(v string) *ListStatefulAsyncInvocationsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListStatefulAsyncInvocationsHeaders) SetXFcCodeChecksum(v string) *ListStatefulAsyncInvocationsHeaders {
	s.XFcCodeChecksum = &v
	return s
}

func (s *ListStatefulAsyncInvocationsHeaders) SetXFcDate(v string) *ListStatefulAsyncInvocationsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListStatefulAsyncInvocationsHeaders) SetXFcInvocationType(v string) *ListStatefulAsyncInvocationsHeaders {
	s.XFcInvocationType = &v
	return s
}

func (s *ListStatefulAsyncInvocationsHeaders) SetXFcLogType(v string) *ListStatefulAsyncInvocationsHeaders {
	s.XFcLogType = &v
	return s
}

func (s *ListStatefulAsyncInvocationsHeaders) SetXFcTraceId(v string) *ListStatefulAsyncInvocationsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListStatefulAsyncInvocationsRequest struct {
	// - **true**: returns the invocationPayload parameter in the response.
	// - **false**: does not return the invocationPayload parameter in the response.
	//
	// > The `invocationPayload` parameter indicates the input parameters of an asynchronous task.
	IncludePayload *bool `json:"includePayload,omitempty" xml:"includePayload,omitempty"`
	// The name prefix of the asynchronous invocation. The names of returned resources must contain the prefix. For example, if invocationidPrefix is set to job, the names of returned resources must start with job.
	InvocationIdPrefix *string `json:"invocationIdPrefix,omitempty" xml:"invocationIdPrefix,omitempty"`
	// The maximum number of asynchronous invocations to return. Valid values: [1, 100]. Default value: 50.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token used to obtain more results. If the number of resources exceeds the limit, the nextToken parameter is returned. You can include the parameter in subsequent calls to obtain more results. You do not need to provide this parameter in the first call.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The version or alias of the service to which the asynchronous task belongs.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The order in which the returned asynchronous invocations are sorted. Valid values:
	//
	// - **asc**: in ascending order
	// - **desc**: in descending order
	SortOrderByTime *string `json:"sortOrderByTime,omitempty" xml:"sortOrderByTime,omitempty"`
	// The start time of the asynchronous task.
	StartedTimeBegin *int64 `json:"startedTimeBegin,omitempty" xml:"startedTimeBegin,omitempty"`
	// The end time of the asynchronous task.
	StartedTimeEnd *int64 `json:"startedTimeEnd,omitempty" xml:"startedTimeEnd,omitempty"`
	// The status of the asynchronous task.
	//
	// - **Enqueued**: The asynchronous invocation is enqueued and is waiting to be executed.
	// - **Succeeded**: The invocation is successful.
	// - **Failed**: The invocation fails.
	// - **Running**: The invocation is being executed.
	// - **Stopped**: The invocation is terminated.
	// - **Stopping**: The invocation is being terminated.
	// - **Invalid**: The invocation is invalid and not executed due to specific reasons. For example, the function is deleted.
	// - **Expired**: The maximum validity period of messages is specified for asynchronous invocation. The invocation is discarded and not executed because the specified maximum validity period of messages expires.
	// - **Retrying**: The asynchronous invocation is being retried due to an execution error.
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s ListStatefulAsyncInvocationsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationsRequest) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationsRequest) SetIncludePayload(v bool) *ListStatefulAsyncInvocationsRequest {
	s.IncludePayload = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetInvocationIdPrefix(v string) *ListStatefulAsyncInvocationsRequest {
	s.InvocationIdPrefix = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetLimit(v int32) *ListStatefulAsyncInvocationsRequest {
	s.Limit = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetNextToken(v string) *ListStatefulAsyncInvocationsRequest {
	s.NextToken = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetQualifier(v string) *ListStatefulAsyncInvocationsRequest {
	s.Qualifier = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetSortOrderByTime(v string) *ListStatefulAsyncInvocationsRequest {
	s.SortOrderByTime = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetStartedTimeBegin(v int64) *ListStatefulAsyncInvocationsRequest {
	s.StartedTimeBegin = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetStartedTimeEnd(v int64) *ListStatefulAsyncInvocationsRequest {
	s.StartedTimeEnd = &v
	return s
}

func (s *ListStatefulAsyncInvocationsRequest) SetStatus(v string) *ListStatefulAsyncInvocationsRequest {
	s.Status = &v
	return s
}

type ListStatefulAsyncInvocationsResponseBody struct {
	// The information about asynchronous tasks.
	Invocations []*StatefulAsyncInvocation `json:"invocations,omitempty" xml:"invocations,omitempty" type:"Repeated"`
	// The token used to obtain more results. If this parameter is left empty, all the results are returned.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListStatefulAsyncInvocationsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationsResponseBody) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationsResponseBody) SetInvocations(v []*StatefulAsyncInvocation) *ListStatefulAsyncInvocationsResponseBody {
	s.Invocations = v
	return s
}

func (s *ListStatefulAsyncInvocationsResponseBody) SetNextToken(v string) *ListStatefulAsyncInvocationsResponseBody {
	s.NextToken = &v
	return s
}

type ListStatefulAsyncInvocationsResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListStatefulAsyncInvocationsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListStatefulAsyncInvocationsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListStatefulAsyncInvocationsResponse) GoString() string {
	return s.String()
}

func (s *ListStatefulAsyncInvocationsResponse) SetHeaders(v map[string]*string) *ListStatefulAsyncInvocationsResponse {
	s.Headers = v
	return s
}

func (s *ListStatefulAsyncInvocationsResponse) SetStatusCode(v int32) *ListStatefulAsyncInvocationsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListStatefulAsyncInvocationsResponse) SetBody(v *ListStatefulAsyncInvocationsResponseBody) *ListStatefulAsyncInvocationsResponse {
	s.Body = v
	return s
}

type ListTaggedResourcesHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListTaggedResourcesHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListTaggedResourcesHeaders) GoString() string {
	return s.String()
}

func (s *ListTaggedResourcesHeaders) SetCommonHeaders(v map[string]*string) *ListTaggedResourcesHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListTaggedResourcesHeaders) SetXFcAccountId(v string) *ListTaggedResourcesHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListTaggedResourcesHeaders) SetXFcDate(v string) *ListTaggedResourcesHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListTaggedResourcesHeaders) SetXFcTraceId(v string) *ListTaggedResourcesHeaders {
	s.XFcTraceId = &v
	return s
}

type ListTaggedResourcesRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token used to obtain more results. You do not need to provide this parameter in the first call. The tokens for subsequent queries are obtained from the returned results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
}

func (s ListTaggedResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTaggedResourcesRequest) GoString() string {
	return s.String()
}

func (s *ListTaggedResourcesRequest) SetLimit(v int32) *ListTaggedResourcesRequest {
	s.Limit = &v
	return s
}

func (s *ListTaggedResourcesRequest) SetNextToken(v string) *ListTaggedResourcesRequest {
	s.NextToken = &v
	return s
}

type ListTaggedResourcesResponseBody struct {
	// The token used to obtain more results. You do not need to provide this parameter in the first call. The tokens for subsequent queries are obtained from the returned results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The information about tagged services.
	Resources []*Resource `json:"resources,omitempty" xml:"resources,omitempty" type:"Repeated"`
}

func (s ListTaggedResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListTaggedResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *ListTaggedResourcesResponseBody) SetNextToken(v string) *ListTaggedResourcesResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListTaggedResourcesResponseBody) SetResources(v []*Resource) *ListTaggedResourcesResponseBody {
	s.Resources = v
	return s
}

type ListTaggedResourcesResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListTaggedResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListTaggedResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListTaggedResourcesResponse) GoString() string {
	return s.String()
}

func (s *ListTaggedResourcesResponse) SetHeaders(v map[string]*string) *ListTaggedResourcesResponse {
	s.Headers = v
	return s
}

func (s *ListTaggedResourcesResponse) SetStatusCode(v int32) *ListTaggedResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListTaggedResourcesResponse) SetBody(v *ListTaggedResourcesResponseBody) *ListTaggedResourcesResponse {
	s.Body = v
	return s
}

type ListTriggersHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the request is initiated on the client. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListTriggersHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListTriggersHeaders) GoString() string {
	return s.String()
}

func (s *ListTriggersHeaders) SetCommonHeaders(v map[string]*string) *ListTriggersHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListTriggersHeaders) SetXFcAccountId(v string) *ListTriggersHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListTriggersHeaders) SetXFcDate(v string) *ListTriggersHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListTriggersHeaders) SetXFcTraceId(v string) *ListTriggersHeaders {
	s.XFcTraceId = &v
	return s
}

type ListTriggersRequest struct {
	// The maximum number of resources to return. Default value: 20. Maximum value: 100. The number of returned resources is less than or equal to the specified number.
	Limit *int32 `json:"limit,omitempty" xml:"limit,omitempty"`
	// The token required to obtain more results. You do not need to provide this parameter in the first call. The tokens for subsequent queries are obtained from the returned results.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The prefix that the names of returned resources must contain.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The returned resources are sorted in alphabetical order, and the resources that include and follow the resource specified by the startKey parameter are returned.
	StartKey *string `json:"startKey,omitempty" xml:"startKey,omitempty"`
}

func (s ListTriggersRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTriggersRequest) GoString() string {
	return s.String()
}

func (s *ListTriggersRequest) SetLimit(v int32) *ListTriggersRequest {
	s.Limit = &v
	return s
}

func (s *ListTriggersRequest) SetNextToken(v string) *ListTriggersRequest {
	s.NextToken = &v
	return s
}

func (s *ListTriggersRequest) SetPrefix(v string) *ListTriggersRequest {
	s.Prefix = &v
	return s
}

func (s *ListTriggersRequest) SetStartKey(v string) *ListTriggersRequest {
	s.StartKey = &v
	return s
}

type ListTriggersResponseBody struct {
	// The token used to obtain more results. If this parameter is left empty, all the results are returned.
	NextToken *string `json:"nextToken,omitempty" xml:"nextToken,omitempty"`
	// The information about triggers.
	Triggers []*ListTriggersResponseBodyTriggers `json:"triggers,omitempty" xml:"triggers,omitempty" type:"Repeated"`
}

func (s ListTriggersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListTriggersResponseBody) GoString() string {
	return s.String()
}

func (s *ListTriggersResponseBody) SetNextToken(v string) *ListTriggersResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListTriggersResponseBody) SetTriggers(v []*ListTriggersResponseBodyTriggers) *ListTriggersResponseBody {
	s.Triggers = v
	return s
}

type ListTriggersResponseBodyTriggers struct {
	// The time when the trigger was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the trigger.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The domain name used to invoke the function by using HTTP. You can add this domain name as the prefix to the endpoint of Function Compute. This way, you can invoke the function that corresponds to the trigger by using HTTP. For example, `{domainName}.cn-shanghai.fc.aliyuncs.com`.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The ARN of the RAM role that is used by the event source to invoke the function.
	InvocationRole *string `json:"invocationRole,omitempty" xml:"invocationRole,omitempty"`
	// The time when the trigger was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The ARN of the event source.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
	// The configurations of the trigger. The configurations vary based on the trigger type. For more information about the format, see the following topics:
	//
	// *   OSS trigger: [OSSTriggerConfig](~~struct:OSSTriggerConfig~~).
	// *   Log Service trigger: [LogTriggerConfig](~~struct:LogTriggerConfig~~).
	// *   Time trigger: [TimeTriggerConfig](~~struct:TimeTriggerConfig~~).
	// *   HTTP trigger: [HTTPTriggerConfig](~~struct:HTTPTriggerConfig~~).
	// *   Tablestore trigger: Specify the **SourceArn** parameter and leave this parameter empty.
	// *   Alibaba Cloud CDN event trigger: [CDNEventsTriggerConfig](~~struct:CDNEventsTriggerConfig~~).
	// *   MNS topic trigger: [MnsTopicTriggerConfig](~~struct:MnsTopicTriggerConfig~~).
	TriggerConfig *string `json:"triggerConfig,omitempty" xml:"triggerConfig,omitempty"`
	// The unique ID of the trigger.
	TriggerId *string `json:"triggerId,omitempty" xml:"triggerId,omitempty"`
	// The name of the trigger.
	TriggerName *string `json:"triggerName,omitempty" xml:"triggerName,omitempty"`
	// The trigger type, such as **oss**, **log**, **tablestore**, **timer**, **http**, **cdn\_events**, and **mns\_topic**.
	TriggerType *string `json:"triggerType,omitempty" xml:"triggerType,omitempty"`
	// The public domain address. You can access HTTP triggers over the Internet by using HTTP or HTTPS.
	UrlInternet *string `json:"urlInternet,omitempty" xml:"urlInternet,omitempty"`
	// The private endpoint. In a VPC, you can access HTTP triggers by using HTTP or HTTPS.
	UrlIntranet *string `json:"urlIntranet,omitempty" xml:"urlIntranet,omitempty"`
}

func (s ListTriggersResponseBodyTriggers) String() string {
	return tea.Prettify(s)
}

func (s ListTriggersResponseBodyTriggers) GoString() string {
	return s.String()
}

func (s *ListTriggersResponseBodyTriggers) SetCreatedTime(v string) *ListTriggersResponseBodyTriggers {
	s.CreatedTime = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetDescription(v string) *ListTriggersResponseBodyTriggers {
	s.Description = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetDomainName(v string) *ListTriggersResponseBodyTriggers {
	s.DomainName = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetInvocationRole(v string) *ListTriggersResponseBodyTriggers {
	s.InvocationRole = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetLastModifiedTime(v string) *ListTriggersResponseBodyTriggers {
	s.LastModifiedTime = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetQualifier(v string) *ListTriggersResponseBodyTriggers {
	s.Qualifier = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetSourceArn(v string) *ListTriggersResponseBodyTriggers {
	s.SourceArn = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetTriggerConfig(v string) *ListTriggersResponseBodyTriggers {
	s.TriggerConfig = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetTriggerId(v string) *ListTriggersResponseBodyTriggers {
	s.TriggerId = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetTriggerName(v string) *ListTriggersResponseBodyTriggers {
	s.TriggerName = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetTriggerType(v string) *ListTriggersResponseBodyTriggers {
	s.TriggerType = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetUrlInternet(v string) *ListTriggersResponseBodyTriggers {
	s.UrlInternet = &v
	return s
}

func (s *ListTriggersResponseBodyTriggers) SetUrlIntranet(v string) *ListTriggersResponseBodyTriggers {
	s.UrlIntranet = &v
	return s
}

type ListTriggersResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListTriggersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListTriggersResponse) String() string {
	return tea.Prettify(s)
}

func (s ListTriggersResponse) GoString() string {
	return s.String()
}

func (s *ListTriggersResponse) SetHeaders(v map[string]*string) *ListTriggersResponse {
	s.Headers = v
	return s
}

func (s *ListTriggersResponse) SetStatusCode(v int32) *ListTriggersResponse {
	s.StatusCode = &v
	return s
}

func (s *ListTriggersResponse) SetBody(v *ListTriggersResponseBody) *ListTriggersResponse {
	s.Body = v
	return s
}

type ListVpcBindingsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ListVpcBindingsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListVpcBindingsHeaders) GoString() string {
	return s.String()
}

func (s *ListVpcBindingsHeaders) SetCommonHeaders(v map[string]*string) *ListVpcBindingsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListVpcBindingsHeaders) SetXFcAccountId(v string) *ListVpcBindingsHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ListVpcBindingsHeaders) SetXFcDate(v string) *ListVpcBindingsHeaders {
	s.XFcDate = &v
	return s
}

func (s *ListVpcBindingsHeaders) SetXFcTraceId(v string) *ListVpcBindingsHeaders {
	s.XFcTraceId = &v
	return s
}

type ListVpcBindingsResponseBody struct {
	// The IDs of bound VPCs.
	VpcIds []*string `json:"vpcIds,omitempty" xml:"vpcIds,omitempty" type:"Repeated"`
}

func (s ListVpcBindingsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListVpcBindingsResponseBody) GoString() string {
	return s.String()
}

func (s *ListVpcBindingsResponseBody) SetVpcIds(v []*string) *ListVpcBindingsResponseBody {
	s.VpcIds = v
	return s
}

type ListVpcBindingsResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListVpcBindingsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListVpcBindingsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListVpcBindingsResponse) GoString() string {
	return s.String()
}

func (s *ListVpcBindingsResponse) SetHeaders(v map[string]*string) *ListVpcBindingsResponse {
	s.Headers = v
	return s
}

func (s *ListVpcBindingsResponse) SetStatusCode(v int32) *ListVpcBindingsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListVpcBindingsResponse) SetBody(v *ListVpcBindingsResponseBody) *ListVpcBindingsResponse {
	s.Body = v
	return s
}

type PublishServiceVersionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ETag value of the service. This value is used to ensure that the modified service is consistent with the service to be modified. The ETag value is returned in the responses of the [CreateService](~~175256~~), [UpdateService](~~188167~~), and [GetService](~~189225~~) operations.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The start time when the function is invoked. Specify the time in the yyyy-mm-ddhh:mm:ss format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s PublishServiceVersionHeaders) String() string {
	return tea.Prettify(s)
}

func (s PublishServiceVersionHeaders) GoString() string {
	return s.String()
}

func (s *PublishServiceVersionHeaders) SetCommonHeaders(v map[string]*string) *PublishServiceVersionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PublishServiceVersionHeaders) SetIfMatch(v string) *PublishServiceVersionHeaders {
	s.IfMatch = &v
	return s
}

func (s *PublishServiceVersionHeaders) SetXFcAccountId(v string) *PublishServiceVersionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *PublishServiceVersionHeaders) SetXFcDate(v string) *PublishServiceVersionHeaders {
	s.XFcDate = &v
	return s
}

func (s *PublishServiceVersionHeaders) SetXFcTraceId(v string) *PublishServiceVersionHeaders {
	s.XFcTraceId = &v
	return s
}

type PublishServiceVersionRequest struct {
	// The description of the service version.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
}

func (s PublishServiceVersionRequest) String() string {
	return tea.Prettify(s)
}

func (s PublishServiceVersionRequest) GoString() string {
	return s.String()
}

func (s *PublishServiceVersionRequest) SetDescription(v string) *PublishServiceVersionRequest {
	s.Description = &v
	return s
}

type PublishServiceVersionResponseBody struct {
	// The time when the service version was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the service version.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The time when the service version was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The version of the service.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s PublishServiceVersionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s PublishServiceVersionResponseBody) GoString() string {
	return s.String()
}

func (s *PublishServiceVersionResponseBody) SetCreatedTime(v string) *PublishServiceVersionResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *PublishServiceVersionResponseBody) SetDescription(v string) *PublishServiceVersionResponseBody {
	s.Description = &v
	return s
}

func (s *PublishServiceVersionResponseBody) SetLastModifiedTime(v string) *PublishServiceVersionResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *PublishServiceVersionResponseBody) SetVersionId(v string) *PublishServiceVersionResponseBody {
	s.VersionId = &v
	return s
}

type PublishServiceVersionResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *PublishServiceVersionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s PublishServiceVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s PublishServiceVersionResponse) GoString() string {
	return s.String()
}

func (s *PublishServiceVersionResponse) SetHeaders(v map[string]*string) *PublishServiceVersionResponse {
	s.Headers = v
	return s
}

func (s *PublishServiceVersionResponse) SetStatusCode(v int32) *PublishServiceVersionResponse {
	s.StatusCode = &v
	return s
}

func (s *PublishServiceVersionResponse) SetBody(v *PublishServiceVersionResponseBody) *PublishServiceVersionResponse {
	s.Body = v
	return s
}

type PutFunctionAsyncInvokeConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s PutFunctionAsyncInvokeConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionAsyncInvokeConfigHeaders) GoString() string {
	return s.String()
}

func (s *PutFunctionAsyncInvokeConfigHeaders) SetCommonHeaders(v map[string]*string) *PutFunctionAsyncInvokeConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutFunctionAsyncInvokeConfigHeaders) SetXFcAccountId(v string) *PutFunctionAsyncInvokeConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigHeaders) SetXFcDate(v string) *PutFunctionAsyncInvokeConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigHeaders) SetXFcTraceId(v string) *PutFunctionAsyncInvokeConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type PutFunctionAsyncInvokeConfigRequest struct {
	// The configuration structure of the destination for asynchronous invocation.
	DestinationConfig *DestinationConfig `json:"destinationConfig,omitempty" xml:"destinationConfig,omitempty"`
	// The maximum validity period of messages. Valid values: 1 to 2592000. Unit: seconds.
	MaxAsyncEventAgeInSeconds *int64 `json:"maxAsyncEventAgeInSeconds,omitempty" xml:"maxAsyncEventAgeInSeconds,omitempty"`
	// The maximum number of retries allowed after an asynchronous invocation fails. Default value: 3. Valid values: 0 to 8.
	MaxAsyncRetryAttempts *int64 `json:"maxAsyncRetryAttempts,omitempty" xml:"maxAsyncRetryAttempts,omitempty"`
	// Specifies whether to enable the asynchronous task feature.
	//
	// - **true**: enables the asynchronous task feature.
	// - **false**: does not enable the asynchronous task feature.
	StatefulInvocation *bool `json:"statefulInvocation,omitempty" xml:"statefulInvocation,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s PutFunctionAsyncInvokeConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionAsyncInvokeConfigRequest) GoString() string {
	return s.String()
}

func (s *PutFunctionAsyncInvokeConfigRequest) SetDestinationConfig(v *DestinationConfig) *PutFunctionAsyncInvokeConfigRequest {
	s.DestinationConfig = v
	return s
}

func (s *PutFunctionAsyncInvokeConfigRequest) SetMaxAsyncEventAgeInSeconds(v int64) *PutFunctionAsyncInvokeConfigRequest {
	s.MaxAsyncEventAgeInSeconds = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigRequest) SetMaxAsyncRetryAttempts(v int64) *PutFunctionAsyncInvokeConfigRequest {
	s.MaxAsyncRetryAttempts = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigRequest) SetStatefulInvocation(v bool) *PutFunctionAsyncInvokeConfigRequest {
	s.StatefulInvocation = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigRequest) SetQualifier(v string) *PutFunctionAsyncInvokeConfigRequest {
	s.Qualifier = &v
	return s
}

type PutFunctionAsyncInvokeConfigResponseBody struct {
	// The creation time.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The configuration structure of the destination for asynchronous invocation.
	DestinationConfig *DestinationConfig `json:"destinationConfig,omitempty" xml:"destinationConfig,omitempty"`
	// The name of the function.
	Function *string `json:"function,omitempty" xml:"function,omitempty"`
	// The time when the configuration was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The maximum validity period of messages.
	MaxAsyncEventAgeInSeconds *int64 `json:"maxAsyncEventAgeInSeconds,omitempty" xml:"maxAsyncEventAgeInSeconds,omitempty"`
	// The maximum number of retries allowed after an asynchronous invocation fails.
	MaxAsyncRetryAttempts *int64 `json:"maxAsyncRetryAttempts,omitempty" xml:"maxAsyncRetryAttempts,omitempty"`
	// The qualifier.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The name of the service.
	Service *string `json:"service,omitempty" xml:"service,omitempty"`
	// Specifies whether to enable the asynchronous task feature.
	//
	// - **true**: enables the asynchronous task feature.
	// - **false**: does not enable the asynchronous task feature.
	StatefulInvocation *bool `json:"statefulInvocation,omitempty" xml:"statefulInvocation,omitempty"`
}

func (s PutFunctionAsyncInvokeConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionAsyncInvokeConfigResponseBody) GoString() string {
	return s.String()
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetCreatedTime(v string) *PutFunctionAsyncInvokeConfigResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetDestinationConfig(v *DestinationConfig) *PutFunctionAsyncInvokeConfigResponseBody {
	s.DestinationConfig = v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetFunction(v string) *PutFunctionAsyncInvokeConfigResponseBody {
	s.Function = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetLastModifiedTime(v string) *PutFunctionAsyncInvokeConfigResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetMaxAsyncEventAgeInSeconds(v int64) *PutFunctionAsyncInvokeConfigResponseBody {
	s.MaxAsyncEventAgeInSeconds = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetMaxAsyncRetryAttempts(v int64) *PutFunctionAsyncInvokeConfigResponseBody {
	s.MaxAsyncRetryAttempts = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetQualifier(v string) *PutFunctionAsyncInvokeConfigResponseBody {
	s.Qualifier = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetService(v string) *PutFunctionAsyncInvokeConfigResponseBody {
	s.Service = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponseBody) SetStatefulInvocation(v bool) *PutFunctionAsyncInvokeConfigResponseBody {
	s.StatefulInvocation = &v
	return s
}

type PutFunctionAsyncInvokeConfigResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *PutFunctionAsyncInvokeConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s PutFunctionAsyncInvokeConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionAsyncInvokeConfigResponse) GoString() string {
	return s.String()
}

func (s *PutFunctionAsyncInvokeConfigResponse) SetHeaders(v map[string]*string) *PutFunctionAsyncInvokeConfigResponse {
	s.Headers = v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponse) SetStatusCode(v int32) *PutFunctionAsyncInvokeConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *PutFunctionAsyncInvokeConfigResponse) SetBody(v *PutFunctionAsyncInvokeConfigResponseBody) *PutFunctionAsyncInvokeConfigResponse {
	s.Body = v
	return s
}

type PutFunctionOnDemandConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// If the ETag specified in the request matches the ETag value of the OndemandConfig, FC returns 200 OK. If the ETag specified in the request does not match the ETag value of the object, FC returns 412 Precondition Failed.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The start time when the function is invoked. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API, which is also the unique ID of the request.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s PutFunctionOnDemandConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionOnDemandConfigHeaders) GoString() string {
	return s.String()
}

func (s *PutFunctionOnDemandConfigHeaders) SetCommonHeaders(v map[string]*string) *PutFunctionOnDemandConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutFunctionOnDemandConfigHeaders) SetIfMatch(v string) *PutFunctionOnDemandConfigHeaders {
	s.IfMatch = &v
	return s
}

func (s *PutFunctionOnDemandConfigHeaders) SetXFcAccountId(v string) *PutFunctionOnDemandConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *PutFunctionOnDemandConfigHeaders) SetXFcDate(v string) *PutFunctionOnDemandConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *PutFunctionOnDemandConfigHeaders) SetXFcTraceId(v string) *PutFunctionOnDemandConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type PutFunctionOnDemandConfigRequest struct {
	// The maximum number of on-demand instances. For more information, see [Instance scaling limits](~~185038~~).
	MaximumInstanceCount *int64 `json:"maximumInstanceCount,omitempty" xml:"maximumInstanceCount,omitempty"`
	// The alias of the service or LATEST.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s PutFunctionOnDemandConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionOnDemandConfigRequest) GoString() string {
	return s.String()
}

func (s *PutFunctionOnDemandConfigRequest) SetMaximumInstanceCount(v int64) *PutFunctionOnDemandConfigRequest {
	s.MaximumInstanceCount = &v
	return s
}

func (s *PutFunctionOnDemandConfigRequest) SetQualifier(v string) *PutFunctionOnDemandConfigRequest {
	s.Qualifier = &v
	return s
}

type PutFunctionOnDemandConfigResponseBody struct {
	// The maximum number of instances.
	MaximumInstanceCount *int64 `json:"maximumInstanceCount,omitempty" xml:"maximumInstanceCount,omitempty"`
	// The description of the resource.
	Resource *string `json:"resource,omitempty" xml:"resource,omitempty"`
}

func (s PutFunctionOnDemandConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionOnDemandConfigResponseBody) GoString() string {
	return s.String()
}

func (s *PutFunctionOnDemandConfigResponseBody) SetMaximumInstanceCount(v int64) *PutFunctionOnDemandConfigResponseBody {
	s.MaximumInstanceCount = &v
	return s
}

func (s *PutFunctionOnDemandConfigResponseBody) SetResource(v string) *PutFunctionOnDemandConfigResponseBody {
	s.Resource = &v
	return s
}

type PutFunctionOnDemandConfigResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *PutFunctionOnDemandConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s PutFunctionOnDemandConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s PutFunctionOnDemandConfigResponse) GoString() string {
	return s.String()
}

func (s *PutFunctionOnDemandConfigResponse) SetHeaders(v map[string]*string) *PutFunctionOnDemandConfigResponse {
	s.Headers = v
	return s
}

func (s *PutFunctionOnDemandConfigResponse) SetStatusCode(v int32) *PutFunctionOnDemandConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *PutFunctionOnDemandConfigResponse) SetBody(v *PutFunctionOnDemandConfigResponseBody) *PutFunctionOnDemandConfigResponse {
	s.Body = v
	return s
}

type PutLayerACLHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the operation is called. The format is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request for Function Compute API.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s PutLayerACLHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutLayerACLHeaders) GoString() string {
	return s.String()
}

func (s *PutLayerACLHeaders) SetCommonHeaders(v map[string]*string) *PutLayerACLHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutLayerACLHeaders) SetXFcAccountId(v string) *PutLayerACLHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *PutLayerACLHeaders) SetXFcDate(v string) *PutLayerACLHeaders {
	s.XFcDate = &v
	return s
}

func (s *PutLayerACLHeaders) SetXFcTraceId(v string) *PutLayerACLHeaders {
	s.XFcTraceId = &v
	return s
}

type PutLayerACLRequest struct {
	// Specifies whether the layer is public.
	//
	// *   **true**: Public.
	// *   **false**: Not public.
	Public *bool `json:"public,omitempty" xml:"public,omitempty"`
}

func (s PutLayerACLRequest) String() string {
	return tea.Prettify(s)
}

func (s PutLayerACLRequest) GoString() string {
	return s.String()
}

func (s *PutLayerACLRequest) SetPublic(v bool) *PutLayerACLRequest {
	s.Public = &v
	return s
}

type PutLayerACLResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s PutLayerACLResponse) String() string {
	return tea.Prettify(s)
}

func (s PutLayerACLResponse) GoString() string {
	return s.String()
}

func (s *PutLayerACLResponse) SetHeaders(v map[string]*string) *PutLayerACLResponse {
	s.Headers = v
	return s
}

func (s *PutLayerACLResponse) SetStatusCode(v int32) *PutLayerACLResponse {
	s.StatusCode = &v
	return s
}

type PutProvisionConfigHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s PutProvisionConfigHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutProvisionConfigHeaders) GoString() string {
	return s.String()
}

func (s *PutProvisionConfigHeaders) SetCommonHeaders(v map[string]*string) *PutProvisionConfigHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutProvisionConfigHeaders) SetXFcAccountId(v string) *PutProvisionConfigHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *PutProvisionConfigHeaders) SetXFcDate(v string) *PutProvisionConfigHeaders {
	s.XFcDate = &v
	return s
}

func (s *PutProvisionConfigHeaders) SetXFcTraceId(v string) *PutProvisionConfigHeaders {
	s.XFcTraceId = &v
	return s
}

type PutProvisionConfigRequest struct {
	// Specifies whether to always allocate CPU resources. Default value: true.
	AlwaysAllocateCPU *bool `json:"alwaysAllocateCPU,omitempty" xml:"alwaysAllocateCPU,omitempty"`
	// The configurations of scheduled auto scaling.
	ScheduledActions []*ScheduledActions `json:"scheduledActions,omitempty" xml:"scheduledActions,omitempty" type:"Repeated"`
	// The number of provisioned instances. Value range: [1,100000].
	Target *int64 `json:"target,omitempty" xml:"target,omitempty"`
	// The configurations of metric-based auto scaling.
	TargetTrackingPolicies []*TargetTrackingPolicies `json:"targetTrackingPolicies,omitempty" xml:"targetTrackingPolicies,omitempty" type:"Repeated"`
	// The name of the alias.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s PutProvisionConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s PutProvisionConfigRequest) GoString() string {
	return s.String()
}

func (s *PutProvisionConfigRequest) SetAlwaysAllocateCPU(v bool) *PutProvisionConfigRequest {
	s.AlwaysAllocateCPU = &v
	return s
}

func (s *PutProvisionConfigRequest) SetScheduledActions(v []*ScheduledActions) *PutProvisionConfigRequest {
	s.ScheduledActions = v
	return s
}

func (s *PutProvisionConfigRequest) SetTarget(v int64) *PutProvisionConfigRequest {
	s.Target = &v
	return s
}

func (s *PutProvisionConfigRequest) SetTargetTrackingPolicies(v []*TargetTrackingPolicies) *PutProvisionConfigRequest {
	s.TargetTrackingPolicies = v
	return s
}

func (s *PutProvisionConfigRequest) SetQualifier(v string) *PutProvisionConfigRequest {
	s.Qualifier = &v
	return s
}

type PutProvisionConfigResponseBody struct {
	// Specifies whether to always allocate CPU to a function instance.
	AlwaysAllocateCPU *bool `json:"alwaysAllocateCPU,omitempty" xml:"alwaysAllocateCPU,omitempty"`
	// The actual number of provisioned instances.
	Current *int64 `json:"current,omitempty" xml:"current,omitempty"`
	// The description of the resource.
	Resource *string `json:"resource,omitempty" xml:"resource,omitempty"`
	// The configurations of scheduled auto scaling.
	ScheduledActions []*ScheduledActions `json:"scheduledActions,omitempty" xml:"scheduledActions,omitempty" type:"Repeated"`
	// The expected number of provisioned instances.
	Target *int64 `json:"target,omitempty" xml:"target,omitempty"`
	// The configurations of metric-based auto scaling.
	TargetTrackingPolicies []*TargetTrackingPolicies `json:"targetTrackingPolicies,omitempty" xml:"targetTrackingPolicies,omitempty" type:"Repeated"`
}

func (s PutProvisionConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s PutProvisionConfigResponseBody) GoString() string {
	return s.String()
}

func (s *PutProvisionConfigResponseBody) SetAlwaysAllocateCPU(v bool) *PutProvisionConfigResponseBody {
	s.AlwaysAllocateCPU = &v
	return s
}

func (s *PutProvisionConfigResponseBody) SetCurrent(v int64) *PutProvisionConfigResponseBody {
	s.Current = &v
	return s
}

func (s *PutProvisionConfigResponseBody) SetResource(v string) *PutProvisionConfigResponseBody {
	s.Resource = &v
	return s
}

func (s *PutProvisionConfigResponseBody) SetScheduledActions(v []*ScheduledActions) *PutProvisionConfigResponseBody {
	s.ScheduledActions = v
	return s
}

func (s *PutProvisionConfigResponseBody) SetTarget(v int64) *PutProvisionConfigResponseBody {
	s.Target = &v
	return s
}

func (s *PutProvisionConfigResponseBody) SetTargetTrackingPolicies(v []*TargetTrackingPolicies) *PutProvisionConfigResponseBody {
	s.TargetTrackingPolicies = v
	return s
}

type PutProvisionConfigResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *PutProvisionConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s PutProvisionConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s PutProvisionConfigResponse) GoString() string {
	return s.String()
}

func (s *PutProvisionConfigResponse) SetHeaders(v map[string]*string) *PutProvisionConfigResponse {
	s.Headers = v
	return s
}

func (s *PutProvisionConfigResponse) SetStatusCode(v int32) *PutProvisionConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *PutProvisionConfigResponse) SetBody(v *PutProvisionConfigResponseBody) *PutProvisionConfigResponse {
	s.Body = v
	return s
}

type RegisterEventSourceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s RegisterEventSourceHeaders) String() string {
	return tea.Prettify(s)
}

func (s RegisterEventSourceHeaders) GoString() string {
	return s.String()
}

func (s *RegisterEventSourceHeaders) SetCommonHeaders(v map[string]*string) *RegisterEventSourceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *RegisterEventSourceHeaders) SetXFcAccountId(v string) *RegisterEventSourceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *RegisterEventSourceHeaders) SetXFcDate(v string) *RegisterEventSourceHeaders {
	s.XFcDate = &v
	return s
}

func (s *RegisterEventSourceHeaders) SetXFcTraceId(v string) *RegisterEventSourceHeaders {
	s.XFcTraceId = &v
	return s
}

type RegisterEventSourceRequest struct {
	// The Alibaba Cloud Resource Name (ARN) of the event source.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s RegisterEventSourceRequest) String() string {
	return tea.Prettify(s)
}

func (s RegisterEventSourceRequest) GoString() string {
	return s.String()
}

func (s *RegisterEventSourceRequest) SetSourceArn(v string) *RegisterEventSourceRequest {
	s.SourceArn = &v
	return s
}

func (s *RegisterEventSourceRequest) SetQualifier(v string) *RegisterEventSourceRequest {
	s.Qualifier = &v
	return s
}

type RegisterEventSourceResponseBody struct {
	// The time when the event source was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The ARN of the event source.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
}

func (s RegisterEventSourceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RegisterEventSourceResponseBody) GoString() string {
	return s.String()
}

func (s *RegisterEventSourceResponseBody) SetCreatedTime(v string) *RegisterEventSourceResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *RegisterEventSourceResponseBody) SetSourceArn(v string) *RegisterEventSourceResponseBody {
	s.SourceArn = &v
	return s
}

type RegisterEventSourceResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RegisterEventSourceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RegisterEventSourceResponse) String() string {
	return tea.Prettify(s)
}

func (s RegisterEventSourceResponse) GoString() string {
	return s.String()
}

func (s *RegisterEventSourceResponse) SetHeaders(v map[string]*string) *RegisterEventSourceResponse {
	s.Headers = v
	return s
}

func (s *RegisterEventSourceResponse) SetStatusCode(v int32) *RegisterEventSourceResponse {
	s.StatusCode = &v
	return s
}

func (s *RegisterEventSourceResponse) SetBody(v *RegisterEventSourceResponseBody) *RegisterEventSourceResponse {
	s.Body = v
	return s
}

type ReleaseGPUInstanceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the function is invoked. The format of the value is: EEE,d MMM yyyy HH:mm:ss GMT.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s ReleaseGPUInstanceHeaders) String() string {
	return tea.Prettify(s)
}

func (s ReleaseGPUInstanceHeaders) GoString() string {
	return s.String()
}

func (s *ReleaseGPUInstanceHeaders) SetCommonHeaders(v map[string]*string) *ReleaseGPUInstanceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ReleaseGPUInstanceHeaders) SetXFcAccountId(v string) *ReleaseGPUInstanceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *ReleaseGPUInstanceHeaders) SetXFcDate(v string) *ReleaseGPUInstanceHeaders {
	s.XFcDate = &v
	return s
}

func (s *ReleaseGPUInstanceHeaders) SetXFcTraceId(v string) *ReleaseGPUInstanceHeaders {
	s.XFcTraceId = &v
	return s
}

type ReleaseGPUInstanceResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s ReleaseGPUInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s ReleaseGPUInstanceResponse) GoString() string {
	return s.String()
}

func (s *ReleaseGPUInstanceResponse) SetHeaders(v map[string]*string) *ReleaseGPUInstanceResponse {
	s.Headers = v
	return s
}

func (s *ReleaseGPUInstanceResponse) SetStatusCode(v int32) *ReleaseGPUInstanceResponse {
	s.StatusCode = &v
	return s
}

type StopStatefulAsyncInvocationHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s StopStatefulAsyncInvocationHeaders) String() string {
	return tea.Prettify(s)
}

func (s StopStatefulAsyncInvocationHeaders) GoString() string {
	return s.String()
}

func (s *StopStatefulAsyncInvocationHeaders) SetCommonHeaders(v map[string]*string) *StopStatefulAsyncInvocationHeaders {
	s.CommonHeaders = v
	return s
}

func (s *StopStatefulAsyncInvocationHeaders) SetXFcAccountId(v string) *StopStatefulAsyncInvocationHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *StopStatefulAsyncInvocationHeaders) SetXFcDate(v string) *StopStatefulAsyncInvocationHeaders {
	s.XFcDate = &v
	return s
}

func (s *StopStatefulAsyncInvocationHeaders) SetXFcTraceId(v string) *StopStatefulAsyncInvocationHeaders {
	s.XFcTraceId = &v
	return s
}

type StopStatefulAsyncInvocationRequest struct {
	// The version or alias of the service to which the asynchronous task belongs.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
}

func (s StopStatefulAsyncInvocationRequest) String() string {
	return tea.Prettify(s)
}

func (s StopStatefulAsyncInvocationRequest) GoString() string {
	return s.String()
}

func (s *StopStatefulAsyncInvocationRequest) SetQualifier(v string) *StopStatefulAsyncInvocationRequest {
	s.Qualifier = &v
	return s
}

type StopStatefulAsyncInvocationResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s StopStatefulAsyncInvocationResponse) String() string {
	return tea.Prettify(s)
}

func (s StopStatefulAsyncInvocationResponse) GoString() string {
	return s.String()
}

func (s *StopStatefulAsyncInvocationResponse) SetHeaders(v map[string]*string) *StopStatefulAsyncInvocationResponse {
	s.Headers = v
	return s
}

func (s *StopStatefulAsyncInvocationResponse) SetStatusCode(v int32) *StopStatefulAsyncInvocationResponse {
	s.StatusCode = &v
	return s
}

type TagResourceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the function is invoked. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s TagResourceHeaders) String() string {
	return tea.Prettify(s)
}

func (s TagResourceHeaders) GoString() string {
	return s.String()
}

func (s *TagResourceHeaders) SetCommonHeaders(v map[string]*string) *TagResourceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *TagResourceHeaders) SetXFcAccountId(v string) *TagResourceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *TagResourceHeaders) SetXFcDate(v string) *TagResourceHeaders {
	s.XFcDate = &v
	return s
}

func (s *TagResourceHeaders) SetXFcTraceId(v string) *TagResourceHeaders {
	s.XFcTraceId = &v
	return s
}

type TagResourceRequest struct {
	// The ARN of the resource.
	//
	// > You can use the value of this parameter to query the information about the resource, such as the account, service, and region information of the resource. You can manage tags only for services for top level resources.
	ResourceArn *string `json:"resourceArn,omitempty" xml:"resourceArn,omitempty"`
	// The tag dictionary.
	Tags map[string]*string `json:"tags,omitempty" xml:"tags,omitempty"`
}

func (s TagResourceRequest) String() string {
	return tea.Prettify(s)
}

func (s TagResourceRequest) GoString() string {
	return s.String()
}

func (s *TagResourceRequest) SetResourceArn(v string) *TagResourceRequest {
	s.ResourceArn = &v
	return s
}

func (s *TagResourceRequest) SetTags(v map[string]*string) *TagResourceRequest {
	s.Tags = v
	return s
}

type TagResourceResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s TagResourceResponse) String() string {
	return tea.Prettify(s)
}

func (s TagResourceResponse) GoString() string {
	return s.String()
}

func (s *TagResourceResponse) SetHeaders(v map[string]*string) *TagResourceResponse {
	s.Headers = v
	return s
}

func (s *TagResourceResponse) SetStatusCode(v int32) *TagResourceResponse {
	s.StatusCode = &v
	return s
}

type UntagResourceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when Function Compute API is called. Specify the time in the **EEE,d MMM yyyy HH:mm:ss GMT** format.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s UntagResourceHeaders) String() string {
	return tea.Prettify(s)
}

func (s UntagResourceHeaders) GoString() string {
	return s.String()
}

func (s *UntagResourceHeaders) SetCommonHeaders(v map[string]*string) *UntagResourceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UntagResourceHeaders) SetXFcAccountId(v string) *UntagResourceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *UntagResourceHeaders) SetXFcDate(v string) *UntagResourceHeaders {
	s.XFcDate = &v
	return s
}

func (s *UntagResourceHeaders) SetXFcTraceId(v string) *UntagResourceHeaders {
	s.XFcTraceId = &v
	return s
}

type UntagResourceRequest struct {
	// Specifies whether to remove all tags. This parameter takes effect only when no tag key is specified. Valid values:
	//   - **true**: removes all tags.
	//   - **false**: does not remove all tags.
	All *bool `json:"all,omitempty" xml:"all,omitempty"`
	// The ARN of the resource.
	//
	// > You can use the value of this parameter to query the information about the resource, such as the account, service, and region information of the resource. You can manage tags only for services for top level resources.
	ResourceArn *string `json:"resourceArn,omitempty" xml:"resourceArn,omitempty"`
	// The keys of the tags that you want to remove.
	TagKeys []*string `json:"tagKeys,omitempty" xml:"tagKeys,omitempty" type:"Repeated"`
}

func (s UntagResourceRequest) String() string {
	return tea.Prettify(s)
}

func (s UntagResourceRequest) GoString() string {
	return s.String()
}

func (s *UntagResourceRequest) SetAll(v bool) *UntagResourceRequest {
	s.All = &v
	return s
}

func (s *UntagResourceRequest) SetResourceArn(v string) *UntagResourceRequest {
	s.ResourceArn = &v
	return s
}

func (s *UntagResourceRequest) SetTagKeys(v []*string) *UntagResourceRequest {
	s.TagKeys = v
	return s
}

type UntagResourceResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UntagResourceResponse) String() string {
	return tea.Prettify(s)
}

func (s UntagResourceResponse) GoString() string {
	return s.String()
}

func (s *UntagResourceResponse) SetHeaders(v map[string]*string) *UntagResourceResponse {
	s.Headers = v
	return s
}

func (s *UntagResourceResponse) SetStatusCode(v int32) *UntagResourceResponse {
	s.StatusCode = &v
	return s
}

type UpdateAliasHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// If the ETag specified in the request matches the ETag value of the object, the object and 200 OK are returned. Otherwise, 412 Precondition Failed is returned.
	//
	// The ETag value of an object is used to check data integrity of the object. This parameter is empty by default.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the invocation request of Function Compute.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s UpdateAliasHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateAliasHeaders) GoString() string {
	return s.String()
}

func (s *UpdateAliasHeaders) SetCommonHeaders(v map[string]*string) *UpdateAliasHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateAliasHeaders) SetIfMatch(v string) *UpdateAliasHeaders {
	s.IfMatch = &v
	return s
}

func (s *UpdateAliasHeaders) SetXFcAccountId(v string) *UpdateAliasHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *UpdateAliasHeaders) SetXFcDate(v string) *UpdateAliasHeaders {
	s.XFcDate = &v
	return s
}

func (s *UpdateAliasHeaders) SetXFcTraceId(v string) *UpdateAliasHeaders {
	s.XFcTraceId = &v
	return s
}

type UpdateAliasRequest struct {
	// The canary release version to which the alias points and the weight of the canary release version.
	//
	// *   The canary release version takes effect only when the function is invoked.
	// *   The value consists of a version number and a specific weight. For example, 2:0.05 indicates that when a function is invoked, Version 2 is the canary release version, 5% of the traffic is distributed to the canary release version, and 95% of the traffic is distributed to the major version.
	AdditionalVersionWeight map[string]*float32 `json:"additionalVersionWeight,omitempty" xml:"additionalVersionWeight,omitempty"`
	// The description of the alias.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The canary release mode. Valid values:
	//
	// *   **Random**: random canary release. This is the default value.
	// *   **Content**: rule-based canary release.
	ResolvePolicy *string `json:"resolvePolicy,omitempty" xml:"resolvePolicy,omitempty"`
	// The canary release rule. Traffic that meets the canary release rule is routed to the canary release instance.
	RoutePolicy *RoutePolicy `json:"routePolicy,omitempty" xml:"routePolicy,omitempty"`
	// The ID of the version to which the alias points.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s UpdateAliasRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateAliasRequest) GoString() string {
	return s.String()
}

func (s *UpdateAliasRequest) SetAdditionalVersionWeight(v map[string]*float32) *UpdateAliasRequest {
	s.AdditionalVersionWeight = v
	return s
}

func (s *UpdateAliasRequest) SetDescription(v string) *UpdateAliasRequest {
	s.Description = &v
	return s
}

func (s *UpdateAliasRequest) SetResolvePolicy(v string) *UpdateAliasRequest {
	s.ResolvePolicy = &v
	return s
}

func (s *UpdateAliasRequest) SetRoutePolicy(v *RoutePolicy) *UpdateAliasRequest {
	s.RoutePolicy = v
	return s
}

func (s *UpdateAliasRequest) SetVersionId(v string) *UpdateAliasRequest {
	s.VersionId = &v
	return s
}

type UpdateAliasResponseBody struct {
	// The canary release version to which the alias points and the weight of the canary release version.
	//
	// *   The canary release version takes effect only when the function is invoked.
	// *   The value consists of a version number and a specific weight. For example, 2:0.05 indicates that when a function is invoked, Version 2 is the canary release version, 5% of the traffic is distributed to the canary release version, and 95% of the traffic is distributed to the major version.
	AdditionalVersionWeight map[string]*float32 `json:"additionalVersionWeight,omitempty" xml:"additionalVersionWeight,omitempty"`
	// The name of the alias.
	AliasName *string `json:"aliasName,omitempty" xml:"aliasName,omitempty"`
	// The time when the alias was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the alias.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The time when the alias was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The ID of the version to which the alias points.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s UpdateAliasResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpdateAliasResponseBody) GoString() string {
	return s.String()
}

func (s *UpdateAliasResponseBody) SetAdditionalVersionWeight(v map[string]*float32) *UpdateAliasResponseBody {
	s.AdditionalVersionWeight = v
	return s
}

func (s *UpdateAliasResponseBody) SetAliasName(v string) *UpdateAliasResponseBody {
	s.AliasName = &v
	return s
}

func (s *UpdateAliasResponseBody) SetCreatedTime(v string) *UpdateAliasResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *UpdateAliasResponseBody) SetDescription(v string) *UpdateAliasResponseBody {
	s.Description = &v
	return s
}

func (s *UpdateAliasResponseBody) SetLastModifiedTime(v string) *UpdateAliasResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *UpdateAliasResponseBody) SetVersionId(v string) *UpdateAliasResponseBody {
	s.VersionId = &v
	return s
}

type UpdateAliasResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpdateAliasResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpdateAliasResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateAliasResponse) GoString() string {
	return s.String()
}

func (s *UpdateAliasResponse) SetHeaders(v map[string]*string) *UpdateAliasResponse {
	s.Headers = v
	return s
}

func (s *UpdateAliasResponse) SetStatusCode(v int32) *UpdateAliasResponse {
	s.StatusCode = &v
	return s
}

func (s *UpdateAliasResponse) SetBody(v *UpdateAliasResponseBody) *UpdateAliasResponse {
	s.Body = v
	return s
}

type UpdateCustomDomainHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the operation is called. The format is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s UpdateCustomDomainHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateCustomDomainHeaders) GoString() string {
	return s.String()
}

func (s *UpdateCustomDomainHeaders) SetCommonHeaders(v map[string]*string) *UpdateCustomDomainHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateCustomDomainHeaders) SetXFcAccountId(v string) *UpdateCustomDomainHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *UpdateCustomDomainHeaders) SetXFcDate(v string) *UpdateCustomDomainHeaders {
	s.XFcDate = &v
	return s
}

func (s *UpdateCustomDomainHeaders) SetXFcTraceId(v string) *UpdateCustomDomainHeaders {
	s.XFcTraceId = &v
	return s
}

type UpdateCustomDomainRequest struct {
	// The configurations of the HTTPS certificate.
	CertConfig *CertConfig `json:"certConfig,omitempty" xml:"certConfig,omitempty"`
	// The protocol types supported by the domain name. Valid values:
	//
	// *   **HTTP**: Only HTTP is supported.
	// *   **HTTPS**: Only HTTPS is supported.
	// *   **HTTP,HTTPS**: HTTP and HTTPS are supported.
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// The route table that maps the paths to functions when the functions are invoked by using the custom domain name.
	RouteConfig *RouteConfig `json:"routeConfig,omitempty" xml:"routeConfig,omitempty"`
	// The Transport Layer Security (TLS) configuration.
	TlsConfig *TLSConfig `json:"tlsConfig,omitempty" xml:"tlsConfig,omitempty"`
	// The Web Application Firewall (WAF) configuration.
	WafConfig *WAFConfig `json:"wafConfig,omitempty" xml:"wafConfig,omitempty"`
}

func (s UpdateCustomDomainRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateCustomDomainRequest) GoString() string {
	return s.String()
}

func (s *UpdateCustomDomainRequest) SetCertConfig(v *CertConfig) *UpdateCustomDomainRequest {
	s.CertConfig = v
	return s
}

func (s *UpdateCustomDomainRequest) SetProtocol(v string) *UpdateCustomDomainRequest {
	s.Protocol = &v
	return s
}

func (s *UpdateCustomDomainRequest) SetRouteConfig(v *RouteConfig) *UpdateCustomDomainRequest {
	s.RouteConfig = v
	return s
}

func (s *UpdateCustomDomainRequest) SetTlsConfig(v *TLSConfig) *UpdateCustomDomainRequest {
	s.TlsConfig = v
	return s
}

func (s *UpdateCustomDomainRequest) SetWafConfig(v *WAFConfig) *UpdateCustomDomainRequest {
	s.WafConfig = v
	return s
}

type UpdateCustomDomainResponseBody struct {
	// The ID of your Alibaba Cloud account.
	AccountId *string `json:"accountId,omitempty" xml:"accountId,omitempty"`
	// The version of the API.
	ApiVersion *string `json:"apiVersion,omitempty" xml:"apiVersion,omitempty"`
	// The configurations of the HTTPS certificate.
	CertConfig *CertConfig `json:"certConfig,omitempty" xml:"certConfig,omitempty"`
	// The time when the custom domain name was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The domain name.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The time when the domain name was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The protocol type that is supported by the custom domain name.
	//
	// *   **HTTP**: Only HTTP is supported.
	// *   **HTTPS**: Only HTTPS is supported.
	// *   **HTTP,HTTPS**: HTTP and HTTPS are supported.
	Protocol *string `json:"protocol,omitempty" xml:"protocol,omitempty"`
	// The route table that maps the paths to functions when the functions are invoked by using the custom domain name.
	RouteConfig *RouteConfig `json:"routeConfig,omitempty" xml:"routeConfig,omitempty"`
	// The Transport Layer Security (TLS) configuration.
	TlsConfig *TLSConfig `json:"tlsConfig,omitempty" xml:"tlsConfig,omitempty"`
	// The Web Application Firewall (WAF) configuration.
	WafConfig *WAFConfig `json:"wafConfig,omitempty" xml:"wafConfig,omitempty"`
}

func (s UpdateCustomDomainResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpdateCustomDomainResponseBody) GoString() string {
	return s.String()
}

func (s *UpdateCustomDomainResponseBody) SetAccountId(v string) *UpdateCustomDomainResponseBody {
	s.AccountId = &v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetApiVersion(v string) *UpdateCustomDomainResponseBody {
	s.ApiVersion = &v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetCertConfig(v *CertConfig) *UpdateCustomDomainResponseBody {
	s.CertConfig = v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetCreatedTime(v string) *UpdateCustomDomainResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetDomainName(v string) *UpdateCustomDomainResponseBody {
	s.DomainName = &v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetLastModifiedTime(v string) *UpdateCustomDomainResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetProtocol(v string) *UpdateCustomDomainResponseBody {
	s.Protocol = &v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetRouteConfig(v *RouteConfig) *UpdateCustomDomainResponseBody {
	s.RouteConfig = v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetTlsConfig(v *TLSConfig) *UpdateCustomDomainResponseBody {
	s.TlsConfig = v
	return s
}

func (s *UpdateCustomDomainResponseBody) SetWafConfig(v *WAFConfig) *UpdateCustomDomainResponseBody {
	s.WafConfig = v
	return s
}

type UpdateCustomDomainResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpdateCustomDomainResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpdateCustomDomainResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateCustomDomainResponse) GoString() string {
	return s.String()
}

func (s *UpdateCustomDomainResponse) SetHeaders(v map[string]*string) *UpdateCustomDomainResponse {
	s.Headers = v
	return s
}

func (s *UpdateCustomDomainResponse) SetStatusCode(v int32) *UpdateCustomDomainResponse {
	s.StatusCode = &v
	return s
}

func (s *UpdateCustomDomainResponse) SetBody(v *UpdateCustomDomainResponseBody) *UpdateCustomDomainResponse {
	s.Body = v
	return s
}

type UpdateFunctionHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The parameter that is used to ensure that the modified resource is consistent with the resource to be modified. The value of this parameter is returned in the responses of the [CreateFunction](~~415747~~), [GetFunction](~~415750~~), and [UpdateFunction](~~415749~~) operations.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The CRC-64 value of the function code package.
	XFcCodeChecksum *string `json:"X-Fc-Code-Checksum,omitempty" xml:"X-Fc-Code-Checksum,omitempty"`
	// The time on which the function is invoked. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The trace ID of the request. The value is the same as that of the requestId parameter in the response.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s UpdateFunctionHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateFunctionHeaders) GoString() string {
	return s.String()
}

func (s *UpdateFunctionHeaders) SetCommonHeaders(v map[string]*string) *UpdateFunctionHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateFunctionHeaders) SetIfMatch(v string) *UpdateFunctionHeaders {
	s.IfMatch = &v
	return s
}

func (s *UpdateFunctionHeaders) SetXFcAccountId(v string) *UpdateFunctionHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *UpdateFunctionHeaders) SetXFcCodeChecksum(v string) *UpdateFunctionHeaders {
	s.XFcCodeChecksum = &v
	return s
}

func (s *UpdateFunctionHeaders) SetXFcDate(v string) *UpdateFunctionHeaders {
	s.XFcDate = &v
	return s
}

func (s *UpdateFunctionHeaders) SetXFcTraceId(v string) *UpdateFunctionHeaders {
	s.XFcTraceId = &v
	return s
}

type UpdateFunctionRequest struct {
	// The number of requests that can be concurrently processed by a single instance.
	InstanceConcurrency *int32 `json:"InstanceConcurrency,omitempty" xml:"InstanceConcurrency,omitempty"`
	// The port on which the HTTP server listens for the custom runtime or custom container runtime.
	CaPort *int32 `json:"caPort,omitempty" xml:"caPort,omitempty"`
	// The packaged code of the function. **Function code packages** can be provided with the following two methods. You must use only one of the methods in a request.
	//
	// *   Specify the name of the Object Storage Service (OSS) bucket and object where the code package is stored. The names are specified in the **ossBucketName** and **ossObjectName** parameters.
	// *   Specify the Base64-encoded content of the ZIP file by using the **zipFile** parameter.
	Code *Code `json:"code,omitempty" xml:"code,omitempty"`
	// The number of vCPUs of the function. The value must be a multiple of 0.05.
	Cpu *float32 `json:"cpu,omitempty" xml:"cpu,omitempty"`
	// The configuration of the custom container. After you configure the custom container, Function Compute can execute the function in a container created from a custom image.
	CustomContainerConfig *CustomContainerConfig `json:"customContainerConfig,omitempty" xml:"customContainerConfig,omitempty"`
	// The custom DNS configurations of the function.
	CustomDNS *CustomDNS `json:"customDNS,omitempty" xml:"customDNS,omitempty"`
	// The custom health check configurations of the function. This parameter is applicable to only custom runtimes and custom containers.
	CustomHealthCheckConfig *CustomHealthCheckConfig `json:"customHealthCheckConfig,omitempty" xml:"customHealthCheckConfig,omitempty"`
	// The configurations of the custom runtime.
	CustomRuntimeConfig *CustomRuntimeConfig `json:"customRuntimeConfig,omitempty" xml:"customRuntimeConfig,omitempty"`
	// The description of the function.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The disk size of the function. Unit: MB. Valid values: 512 and 10240.
	DiskSize *int32 `json:"diskSize,omitempty" xml:"diskSize,omitempty"`
	// The environment variables that are configured for the function. You can obtain the values of the environment variables from the function. For more information, see [Environment variables](~~69777~~).
	EnvironmentVariables map[string]*string `json:"environmentVariables,omitempty" xml:"environmentVariables,omitempty"`
	// The GPU memory capacity for the function. Unit: MB. The value must be a multiple of 1,024.
	GpuMemorySize *int32 `json:"gpuMemorySize,omitempty" xml:"gpuMemorySize,omitempty"`
	// The handler of the function. The format varies based on the programming language. For more information, see [Function handlers](~~157704~~).
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	// The timeout period for the execution of the Initializer hook. Unit: seconds. Default value: 3. Minimum value: 1. When the period ends, the execution of the Initializer hook is terminated.
	InitializationTimeout *int32 `json:"initializationTimeout,omitempty" xml:"initializationTimeout,omitempty"`
	// The handler of the Initializer hook. The format is determined by the programming language. For more information, see [Function handlers](~~157704~~).
	Initializer *string `json:"initializer,omitempty" xml:"initializer,omitempty"`
	// The lifecycle configurations of the instance.
	InstanceLifecycleConfig *InstanceLifecycleConfig `json:"instanceLifecycleConfig,omitempty" xml:"instanceLifecycleConfig,omitempty"`
	// The soft concurrency of the instance. You can use this parameter to implement graceful scale-up of instances. If the number of concurrent requests on an instance is greater than the value of soft concurrency, an instance scale-up is triggered. For example, if your instance requires a long time to start, you can specify a suitable soft concurrency to start the instance in advance.
	//
	// The value must be less than or equal to that of the **instanceConcurrency** parameter.
	InstanceSoftConcurrency *int32 `json:"instanceSoftConcurrency,omitempty" xml:"instanceSoftConcurrency,omitempty"`
	// The instance type of the function. Valid values:
	//
	// *   **e1**: elastic instance
	// *   **c1**: performance instance
	// *   **fc.gpu.tesla.1**: GPU-accelerated instance (Tesla T4)
	// *   **fc.gpu.ampere.1**: GPU-accelerated instance (Ampere A10)
	// *   **g1**: same as **fc.gpu.tesla.1**
	InstanceType *string `json:"instanceType,omitempty" xml:"instanceType,omitempty"`
	// The information about layers.
	//
	// > Multiple layers are merged based on the order of array subscripts. The content of a layer with a smaller subscript overwrites the file that has the same name as a layer with a larger subscript.
	Layers []*string `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The memory size for the function. Unit: MB. The memory size must be a multiple of 64. The memory size varies based on the function instance type. For more information, see [Instance types](~~179379~~).
	MemorySize *int32 `json:"memorySize,omitempty" xml:"memorySize,omitempty"`
	// The runtime environment of the function. Valid values: **nodejs16**, **nodejs14**, **nodejs12**, **nodejs10**, **nodejs8**, **nodejs6**, **nodejs4.4**, **python3.9**, **python3**, **python2.7**, **java11**, **java8**, **go1**, **php7.2**, **dotnetcore3.1**, **dotnetcore2.1**, **custom** and **custom-container**. For more information, see [Supported function runtime environments](~~73338~~).
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The timeout period for the execution of the function. Unit: seconds. Default value: 3. Minimum value: 1. When the period ends, the execution of the function is terminated.
	Timeout *int32 `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s UpdateFunctionRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateFunctionRequest) GoString() string {
	return s.String()
}

func (s *UpdateFunctionRequest) SetInstanceConcurrency(v int32) *UpdateFunctionRequest {
	s.InstanceConcurrency = &v
	return s
}

func (s *UpdateFunctionRequest) SetCaPort(v int32) *UpdateFunctionRequest {
	s.CaPort = &v
	return s
}

func (s *UpdateFunctionRequest) SetCode(v *Code) *UpdateFunctionRequest {
	s.Code = v
	return s
}

func (s *UpdateFunctionRequest) SetCpu(v float32) *UpdateFunctionRequest {
	s.Cpu = &v
	return s
}

func (s *UpdateFunctionRequest) SetCustomContainerConfig(v *CustomContainerConfig) *UpdateFunctionRequest {
	s.CustomContainerConfig = v
	return s
}

func (s *UpdateFunctionRequest) SetCustomDNS(v *CustomDNS) *UpdateFunctionRequest {
	s.CustomDNS = v
	return s
}

func (s *UpdateFunctionRequest) SetCustomHealthCheckConfig(v *CustomHealthCheckConfig) *UpdateFunctionRequest {
	s.CustomHealthCheckConfig = v
	return s
}

func (s *UpdateFunctionRequest) SetCustomRuntimeConfig(v *CustomRuntimeConfig) *UpdateFunctionRequest {
	s.CustomRuntimeConfig = v
	return s
}

func (s *UpdateFunctionRequest) SetDescription(v string) *UpdateFunctionRequest {
	s.Description = &v
	return s
}

func (s *UpdateFunctionRequest) SetDiskSize(v int32) *UpdateFunctionRequest {
	s.DiskSize = &v
	return s
}

func (s *UpdateFunctionRequest) SetEnvironmentVariables(v map[string]*string) *UpdateFunctionRequest {
	s.EnvironmentVariables = v
	return s
}

func (s *UpdateFunctionRequest) SetGpuMemorySize(v int32) *UpdateFunctionRequest {
	s.GpuMemorySize = &v
	return s
}

func (s *UpdateFunctionRequest) SetHandler(v string) *UpdateFunctionRequest {
	s.Handler = &v
	return s
}

func (s *UpdateFunctionRequest) SetInitializationTimeout(v int32) *UpdateFunctionRequest {
	s.InitializationTimeout = &v
	return s
}

func (s *UpdateFunctionRequest) SetInitializer(v string) *UpdateFunctionRequest {
	s.Initializer = &v
	return s
}

func (s *UpdateFunctionRequest) SetInstanceLifecycleConfig(v *InstanceLifecycleConfig) *UpdateFunctionRequest {
	s.InstanceLifecycleConfig = v
	return s
}

func (s *UpdateFunctionRequest) SetInstanceSoftConcurrency(v int32) *UpdateFunctionRequest {
	s.InstanceSoftConcurrency = &v
	return s
}

func (s *UpdateFunctionRequest) SetInstanceType(v string) *UpdateFunctionRequest {
	s.InstanceType = &v
	return s
}

func (s *UpdateFunctionRequest) SetLayers(v []*string) *UpdateFunctionRequest {
	s.Layers = v
	return s
}

func (s *UpdateFunctionRequest) SetMemorySize(v int32) *UpdateFunctionRequest {
	s.MemorySize = &v
	return s
}

func (s *UpdateFunctionRequest) SetRuntime(v string) *UpdateFunctionRequest {
	s.Runtime = &v
	return s
}

func (s *UpdateFunctionRequest) SetTimeout(v int32) *UpdateFunctionRequest {
	s.Timeout = &v
	return s
}

type UpdateFunctionResponseBody struct {
	// The port on which the HTTP server listens for the custom runtime or custom container runtime.
	CaPort *int32 `json:"caPort,omitempty" xml:"caPort,omitempty"`
	// The CRC-64 value of the function code package.
	CodeChecksum *string `json:"codeChecksum,omitempty" xml:"codeChecksum,omitempty"`
	// The size of the function code package that is returned by the system. Unit: bytes.
	CodeSize *int64 `json:"codeSize,omitempty" xml:"codeSize,omitempty"`
	// The number of vCPUs of the function. The value must be a multiple of 0.05.
	Cpu *float32 `json:"cpu,omitempty" xml:"cpu,omitempty"`
	// The time when the function was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The configurations of the custom container runtime. After you configure the custom container runtime, Function Compute can execute the function in a container created from a custom image.
	CustomContainerConfig *CustomContainerConfig `json:"customContainerConfig,omitempty" xml:"customContainerConfig,omitempty"`
	// The custom DNS configurations of the function.
	CustomDNS *CustomDNS `json:"customDNS,omitempty" xml:"customDNS,omitempty"`
	// The custom health check configuration of the function. This parameter is applicable only to custom runtimes and custom containers.
	CustomHealthCheckConfig *CustomHealthCheckConfig `json:"customHealthCheckConfig,omitempty" xml:"customHealthCheckConfig,omitempty"`
	// The configurations of the custom runtime.
	CustomRuntimeConfig *CustomRuntimeConfig `json:"customRuntimeConfig,omitempty" xml:"customRuntimeConfig,omitempty"`
	// The description of the function.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The disk size of the function. Unit: MB. Valid values: 512 and 10240.
	DiskSize *int32 `json:"diskSize,omitempty" xml:"diskSize,omitempty"`
	// The environment variables that are configured for the function. You can obtain the values of the environment variables from the function. For more information, see [Environment variables](~~69777~~).
	EnvironmentVariables map[string]*string `json:"environmentVariables,omitempty" xml:"environmentVariables,omitempty"`
	// The unique ID that is generated by the system for the function.
	FunctionId *string `json:"functionId,omitempty" xml:"functionId,omitempty"`
	// The name of the function.
	FunctionName *string `json:"functionName,omitempty" xml:"functionName,omitempty"`
	// The GPU memory capacity for the function. Unit: MB. The value must be a multiple of 1,024.
	GpuMemorySize *int32 `json:"gpuMemorySize,omitempty" xml:"gpuMemorySize,omitempty"`
	// The handler of the function.
	Handler *string `json:"handler,omitempty" xml:"handler,omitempty"`
	// The timeout period for the execution of the Initializer hook. Unit: seconds. Default value: 3. Minimum value: 1. When the period ends, the execution of the Initializer hook is terminated.
	InitializationTimeout *int32 `json:"initializationTimeout,omitempty" xml:"initializationTimeout,omitempty"`
	// The handler of the Initializer hook. The format is determined by the programming language.
	Initializer *string `json:"initializer,omitempty" xml:"initializer,omitempty"`
	// The number of requests that can be concurrently processed by a single instance.
	InstanceConcurrency *int32 `json:"instanceConcurrency,omitempty" xml:"instanceConcurrency,omitempty"`
	// The lifecycle configurations of the instance.
	InstanceLifecycleConfig *InstanceLifecycleConfig `json:"instanceLifecycleConfig,omitempty" xml:"instanceLifecycleConfig,omitempty"`
	// The soft concurrency of the instance. You can use this parameter to implement graceful scale-up of instances. If the number of concurrent requests on an instance is greater than the value of soft concurrency, an instance scale-up is triggered. For example, if your instance requires a long time to start, you can specify a suitable soft concurrency to start the instance in advance.
	//
	// The value must be less than or equal to that of the **instanceConcurrency** parameter.
	InstanceSoftConcurrency *int32 `json:"instanceSoftConcurrency,omitempty" xml:"instanceSoftConcurrency,omitempty"`
	// The instance type of the function. Valid values:
	//
	// *   **e1**: elastic instance
	// *   **c1**: performance instance
	// *   **fc.gpu.tesla.1**: GPU-accelerated instance (Tesla T4)
	// *   **fc.gpu.ampere.1**: GPU-accelerated instance (Ampere A10)
	// *   **g1**: same as **fc.gpu.tesla.1**
	InstanceType *string `json:"instanceType,omitempty" xml:"instanceType,omitempty"`
	// The time when the function was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// An array that consists of the information of layers.
	//
	// > Multiple layers are merged based on the order of array subscripts. The content of a layer with a smaller subscript overwrites the file that has the same name as a layer with a larger subscript.
	Layers []*string `json:"layers,omitempty" xml:"layers,omitempty" type:"Repeated"`
	// The memory size that is configured for the function. Unit: MB.
	MemorySize *int32 `json:"memorySize,omitempty" xml:"memorySize,omitempty"`
	// The runtime environment of the function. Valid values: **nodejs16**, **nodejs14**, **nodejs12**, **nodejs10**, **nodejs8**, **nodejs6**, **nodejs4.4**, **python3.9**, **python3**, **python2.7**, **java11**, **java8**, **go1**, **php7.2**, **dotnetcore3.1**, **dotnetcore2.1**, **custom** and **custom-container**. For more information, see [Supported function runtime environments](~~73338~~).
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The timeout period for the execution. Unit: seconds.
	Timeout *int32 `json:"timeout,omitempty" xml:"timeout,omitempty"`
}

func (s UpdateFunctionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpdateFunctionResponseBody) GoString() string {
	return s.String()
}

func (s *UpdateFunctionResponseBody) SetCaPort(v int32) *UpdateFunctionResponseBody {
	s.CaPort = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetCodeChecksum(v string) *UpdateFunctionResponseBody {
	s.CodeChecksum = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetCodeSize(v int64) *UpdateFunctionResponseBody {
	s.CodeSize = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetCpu(v float32) *UpdateFunctionResponseBody {
	s.Cpu = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetCreatedTime(v string) *UpdateFunctionResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetCustomContainerConfig(v *CustomContainerConfig) *UpdateFunctionResponseBody {
	s.CustomContainerConfig = v
	return s
}

func (s *UpdateFunctionResponseBody) SetCustomDNS(v *CustomDNS) *UpdateFunctionResponseBody {
	s.CustomDNS = v
	return s
}

func (s *UpdateFunctionResponseBody) SetCustomHealthCheckConfig(v *CustomHealthCheckConfig) *UpdateFunctionResponseBody {
	s.CustomHealthCheckConfig = v
	return s
}

func (s *UpdateFunctionResponseBody) SetCustomRuntimeConfig(v *CustomRuntimeConfig) *UpdateFunctionResponseBody {
	s.CustomRuntimeConfig = v
	return s
}

func (s *UpdateFunctionResponseBody) SetDescription(v string) *UpdateFunctionResponseBody {
	s.Description = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetDiskSize(v int32) *UpdateFunctionResponseBody {
	s.DiskSize = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetEnvironmentVariables(v map[string]*string) *UpdateFunctionResponseBody {
	s.EnvironmentVariables = v
	return s
}

func (s *UpdateFunctionResponseBody) SetFunctionId(v string) *UpdateFunctionResponseBody {
	s.FunctionId = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetFunctionName(v string) *UpdateFunctionResponseBody {
	s.FunctionName = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetGpuMemorySize(v int32) *UpdateFunctionResponseBody {
	s.GpuMemorySize = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetHandler(v string) *UpdateFunctionResponseBody {
	s.Handler = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetInitializationTimeout(v int32) *UpdateFunctionResponseBody {
	s.InitializationTimeout = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetInitializer(v string) *UpdateFunctionResponseBody {
	s.Initializer = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetInstanceConcurrency(v int32) *UpdateFunctionResponseBody {
	s.InstanceConcurrency = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetInstanceLifecycleConfig(v *InstanceLifecycleConfig) *UpdateFunctionResponseBody {
	s.InstanceLifecycleConfig = v
	return s
}

func (s *UpdateFunctionResponseBody) SetInstanceSoftConcurrency(v int32) *UpdateFunctionResponseBody {
	s.InstanceSoftConcurrency = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetInstanceType(v string) *UpdateFunctionResponseBody {
	s.InstanceType = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetLastModifiedTime(v string) *UpdateFunctionResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetLayers(v []*string) *UpdateFunctionResponseBody {
	s.Layers = v
	return s
}

func (s *UpdateFunctionResponseBody) SetMemorySize(v int32) *UpdateFunctionResponseBody {
	s.MemorySize = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetRuntime(v string) *UpdateFunctionResponseBody {
	s.Runtime = &v
	return s
}

func (s *UpdateFunctionResponseBody) SetTimeout(v int32) *UpdateFunctionResponseBody {
	s.Timeout = &v
	return s
}

type UpdateFunctionResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpdateFunctionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpdateFunctionResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateFunctionResponse) GoString() string {
	return s.String()
}

func (s *UpdateFunctionResponse) SetHeaders(v map[string]*string) *UpdateFunctionResponse {
	s.Headers = v
	return s
}

func (s *UpdateFunctionResponse) SetStatusCode(v int32) *UpdateFunctionResponse {
	s.StatusCode = &v
	return s
}

func (s *UpdateFunctionResponse) SetBody(v *UpdateFunctionResponseBody) *UpdateFunctionResponse {
	s.Body = v
	return s
}

type UpdateServiceHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The value used to ensure that the modified service is consistent with the service to be modified. The value is obtained from the responses of the [CreateService](~~175256~~), [UpdateService](~~188167~~), and [GetService](~~189225~~) operations.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the Function Compute API is called. The format is **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s UpdateServiceHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateServiceHeaders) GoString() string {
	return s.String()
}

func (s *UpdateServiceHeaders) SetCommonHeaders(v map[string]*string) *UpdateServiceHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateServiceHeaders) SetIfMatch(v string) *UpdateServiceHeaders {
	s.IfMatch = &v
	return s
}

func (s *UpdateServiceHeaders) SetXFcAccountId(v string) *UpdateServiceHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *UpdateServiceHeaders) SetXFcDate(v string) *UpdateServiceHeaders {
	s.XFcDate = &v
	return s
}

func (s *UpdateServiceHeaders) SetXFcTraceId(v string) *UpdateServiceHeaders {
	s.XFcTraceId = &v
	return s
}

type UpdateServiceRequest struct {
	// The description of the service.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Specifies whether to allow functions to access the Internet. Valid values:
	//
	// *   **true**: allows functions in the specified service to access the Internet.
	// *   **false**: does not allow functions in the specified service to access the Internet.
	InternetAccess *bool `json:"internetAccess,omitempty" xml:"internetAccess,omitempty"`
	// The log configuration. Function Compute writes function execution logs to the specified Logstore.
	LogConfig *LogConfig `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	// The configurations of the NAS file system. The configurations allow functions to access the specified NAS resources.
	NasConfig *NASConfig `json:"nasConfig,omitempty" xml:"nasConfig,omitempty"`
	// The OSS mount configurations.
	OssMountConfig *OSSMountConfig `json:"ossMountConfig,omitempty" xml:"ossMountConfig,omitempty"`
	// The RAM role that is used to grant required permissions to Function Compute. The RAM role is used in the following scenarios:
	//
	// *   Sends function execution logs to your Logstore.
	// *   Generates a token for a function to access other cloud resources during function execution.
	Role *string `json:"role,omitempty" xml:"role,omitempty"`
	// The configurations of Tracing Analysis. After you configure Tracing Analysis for a service in Function Compute, you can record the execution duration of a request, view the amount of cold start time for a function, and record the execution duration of a function. For more information, see [Overview](~~189804~~).
	TracingConfig *TracingConfig `json:"tracingConfig,omitempty" xml:"tracingConfig,omitempty"`
	// The virtual private cloud (VPC) configuration, which allows functions in the specified service in Function Compute to access the specified VPC.
	VpcConfig *VPCConfig `json:"vpcConfig,omitempty" xml:"vpcConfig,omitempty"`
}

func (s UpdateServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateServiceRequest) GoString() string {
	return s.String()
}

func (s *UpdateServiceRequest) SetDescription(v string) *UpdateServiceRequest {
	s.Description = &v
	return s
}

func (s *UpdateServiceRequest) SetInternetAccess(v bool) *UpdateServiceRequest {
	s.InternetAccess = &v
	return s
}

func (s *UpdateServiceRequest) SetLogConfig(v *LogConfig) *UpdateServiceRequest {
	s.LogConfig = v
	return s
}

func (s *UpdateServiceRequest) SetNasConfig(v *NASConfig) *UpdateServiceRequest {
	s.NasConfig = v
	return s
}

func (s *UpdateServiceRequest) SetOssMountConfig(v *OSSMountConfig) *UpdateServiceRequest {
	s.OssMountConfig = v
	return s
}

func (s *UpdateServiceRequest) SetRole(v string) *UpdateServiceRequest {
	s.Role = &v
	return s
}

func (s *UpdateServiceRequest) SetTracingConfig(v *TracingConfig) *UpdateServiceRequest {
	s.TracingConfig = v
	return s
}

func (s *UpdateServiceRequest) SetVpcConfig(v *VPCConfig) *UpdateServiceRequest {
	s.VpcConfig = v
	return s
}

type UpdateServiceResponseBody struct {
	// The time when the service was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the service.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Specifies whether to allow functions to access the Internet. Valid values:
	//
	// *   **true**: allows functions in the specified service to access the Internet.
	// *   **false**: does not allow functions in the specified service to access the Internet.
	InternetAccess *bool `json:"internetAccess,omitempty" xml:"internetAccess,omitempty"`
	// The time when the service was last modified.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The log configuration, which specifies a Logstore to store function execution logs.
	LogConfig *LogConfig `json:"logConfig,omitempty" xml:"logConfig,omitempty"`
	// The configurations of the NAS file system. The configuration allows functions in the specified service in Function Compute to access the NAS file system.
	NasConfig *NASConfig `json:"nasConfig,omitempty" xml:"nasConfig,omitempty"`
	// The OSS mount configurations.
	OssMountConfig *OSSMountConfig `json:"ossMountConfig,omitempty" xml:"ossMountConfig,omitempty"`
	// The RAM role that is used to grant required permissions to Function Compute. The RAM role is used in the following scenarios:
	//
	// *   Sends function execution logs to your Logstore.
	// *   Generates a token for a function to access other cloud resources during function execution.
	Role *string `json:"role,omitempty" xml:"role,omitempty"`
	// The unique ID generated by the system for the service.
	ServiceId *string `json:"serviceId,omitempty" xml:"serviceId,omitempty"`
	// The name of the service.
	ServiceName *string `json:"serviceName,omitempty" xml:"serviceName,omitempty"`
	// The configurations of Tracing Analysis. After you configure Tracing Analysis for a service in Function Compute, you can record the execution duration of a request, view the amount of cold start time for a function, and record the execution duration of a function. For more information, see [Overview](~~189804~~).
	TracingConfig *TracingConfig `json:"tracingConfig,omitempty" xml:"tracingConfig,omitempty"`
	// The VPC configuration. The configuration allows a function to access the specified VPC.
	VpcConfig *VPCConfig `json:"vpcConfig,omitempty" xml:"vpcConfig,omitempty"`
}

func (s UpdateServiceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpdateServiceResponseBody) GoString() string {
	return s.String()
}

func (s *UpdateServiceResponseBody) SetCreatedTime(v string) *UpdateServiceResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *UpdateServiceResponseBody) SetDescription(v string) *UpdateServiceResponseBody {
	s.Description = &v
	return s
}

func (s *UpdateServiceResponseBody) SetInternetAccess(v bool) *UpdateServiceResponseBody {
	s.InternetAccess = &v
	return s
}

func (s *UpdateServiceResponseBody) SetLastModifiedTime(v string) *UpdateServiceResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *UpdateServiceResponseBody) SetLogConfig(v *LogConfig) *UpdateServiceResponseBody {
	s.LogConfig = v
	return s
}

func (s *UpdateServiceResponseBody) SetNasConfig(v *NASConfig) *UpdateServiceResponseBody {
	s.NasConfig = v
	return s
}

func (s *UpdateServiceResponseBody) SetOssMountConfig(v *OSSMountConfig) *UpdateServiceResponseBody {
	s.OssMountConfig = v
	return s
}

func (s *UpdateServiceResponseBody) SetRole(v string) *UpdateServiceResponseBody {
	s.Role = &v
	return s
}

func (s *UpdateServiceResponseBody) SetServiceId(v string) *UpdateServiceResponseBody {
	s.ServiceId = &v
	return s
}

func (s *UpdateServiceResponseBody) SetServiceName(v string) *UpdateServiceResponseBody {
	s.ServiceName = &v
	return s
}

func (s *UpdateServiceResponseBody) SetTracingConfig(v *TracingConfig) *UpdateServiceResponseBody {
	s.TracingConfig = v
	return s
}

func (s *UpdateServiceResponseBody) SetVpcConfig(v *VPCConfig) *UpdateServiceResponseBody {
	s.VpcConfig = v
	return s
}

type UpdateServiceResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpdateServiceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpdateServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateServiceResponse) GoString() string {
	return s.String()
}

func (s *UpdateServiceResponse) SetHeaders(v map[string]*string) *UpdateServiceResponse {
	s.Headers = v
	return s
}

func (s *UpdateServiceResponse) SetStatusCode(v int32) *UpdateServiceResponse {
	s.StatusCode = &v
	return s
}

func (s *UpdateServiceResponse) SetBody(v *UpdateServiceResponseBody) *UpdateServiceResponse {
	s.Body = v
	return s
}

type UpdateTriggerHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// This parameter is used to ensure that the modified resource is consistent with the resource to be modified. You can obtain the parameter value from the responses of [CreateTrigger](~~190054~~), [GetTrigger](~~190056~~), and [UpdateTrigger](~~190055~~) operations.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// The ID of your Alibaba Cloud account.
	XFcAccountId *string `json:"X-Fc-Account-Id,omitempty" xml:"X-Fc-Account-Id,omitempty"`
	// The time when the request is initiated on the client. The format of the value is: **EEE,d MMM yyyy HH:mm:ss GMT**.
	XFcDate *string `json:"X-Fc-Date,omitempty" xml:"X-Fc-Date,omitempty"`
	// The custom request ID.
	XFcTraceId *string `json:"X-Fc-Trace-Id,omitempty" xml:"X-Fc-Trace-Id,omitempty"`
}

func (s UpdateTriggerHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateTriggerHeaders) GoString() string {
	return s.String()
}

func (s *UpdateTriggerHeaders) SetCommonHeaders(v map[string]*string) *UpdateTriggerHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateTriggerHeaders) SetIfMatch(v string) *UpdateTriggerHeaders {
	s.IfMatch = &v
	return s
}

func (s *UpdateTriggerHeaders) SetXFcAccountId(v string) *UpdateTriggerHeaders {
	s.XFcAccountId = &v
	return s
}

func (s *UpdateTriggerHeaders) SetXFcDate(v string) *UpdateTriggerHeaders {
	s.XFcDate = &v
	return s
}

func (s *UpdateTriggerHeaders) SetXFcTraceId(v string) *UpdateTriggerHeaders {
	s.XFcTraceId = &v
	return s
}

type UpdateTriggerRequest struct {
	// The description of the trigger.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The role that is used by the event source such as OSS to invoke the function. For more information, see [Overview](~~53102~~).
	InvocationRole *string `json:"invocationRole,omitempty" xml:"invocationRole,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The configurations of the trigger. The configurations vary based on the trigger type. For more information about the format, see the following topics:
	//
	// *   OSS trigger: [OSSTriggerConfig](~~struct:OSSTriggerConfig~~).
	// *   Log Service trigger: [LogTriggerConfig](~~struct:LogTriggerConfig~~).
	// *   Time trigger: [TimeTriggerConfig](~~struct:TimeTriggerConfig~~).
	// *   HTTP trigger: [HTTPTriggerConfig](~~struct:HTTPTriggerConfig~~).
	// *   Tablestore trigger: Specify the **SourceArn** parameter and leave this parameter empty.
	// *   Alibaba Cloud CDN event trigger: [CDNEventsTriggerConfig](~~struct:CDNEventsTriggerConfig~~).
	// *   MNS topic trigger: [MnsTopicTriggerConfig](~~struct:MnsTopicTriggerConfig~~).
	TriggerConfig *string `json:"triggerConfig,omitempty" xml:"triggerConfig,omitempty"`
}

func (s UpdateTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateTriggerRequest) GoString() string {
	return s.String()
}

func (s *UpdateTriggerRequest) SetDescription(v string) *UpdateTriggerRequest {
	s.Description = &v
	return s
}

func (s *UpdateTriggerRequest) SetInvocationRole(v string) *UpdateTriggerRequest {
	s.InvocationRole = &v
	return s
}

func (s *UpdateTriggerRequest) SetQualifier(v string) *UpdateTriggerRequest {
	s.Qualifier = &v
	return s
}

func (s *UpdateTriggerRequest) SetTriggerConfig(v string) *UpdateTriggerRequest {
	s.TriggerConfig = &v
	return s
}

type UpdateTriggerResponseBody struct {
	// The time when the audio or video file was created.
	CreatedTime *string `json:"createdTime,omitempty" xml:"createdTime,omitempty"`
	// The description of the trigger.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The domain name used to invoke the function by using HTTP. You can add this domain name as the prefix to the endpoint of Function Compute. This way, you can invoke the function that corresponds to the trigger by using HTTP. For example, `{domainName}.cn-shanghai.fc.aliyuncs.com`.
	DomainName *string `json:"domainName,omitempty" xml:"domainName,omitempty"`
	// The ARN of the RAM role that is used by the event source to invoke the function.
	InvocationRole *string `json:"invocationRole,omitempty" xml:"invocationRole,omitempty"`
	// The last modification time.
	LastModifiedTime *string `json:"lastModifiedTime,omitempty" xml:"lastModifiedTime,omitempty"`
	// The version or alias of the service.
	Qualifier *string `json:"qualifier,omitempty" xml:"qualifier,omitempty"`
	// The ARN of the event source.
	SourceArn *string `json:"sourceArn,omitempty" xml:"sourceArn,omitempty"`
	// The configurations of the trigger. The configurations vary based on the trigger type.
	TriggerConfig *string `json:"triggerConfig,omitempty" xml:"triggerConfig,omitempty"`
	// The unique ID of the trigger.
	TriggerId *string `json:"triggerId,omitempty" xml:"triggerId,omitempty"`
	// The name of the trigger.
	TriggerName *string `json:"triggerName,omitempty" xml:"triggerName,omitempty"`
	// The trigger type, such as **oss**, **log**, **tablestore**, **timer**, **http**, **cdn_events**, and **mns_topic**.
	TriggerType *string `json:"triggerType,omitempty" xml:"triggerType,omitempty"`
	// The public domain address. You can access HTTP triggers over the Internet by using HTTP or HTTPS.
	UrlInternet *string `json:"urlInternet,omitempty" xml:"urlInternet,omitempty"`
	// The private endpoint. In a VPC, you can access HTTP triggers by using HTTP or HTTPS.
	UrlIntranet *string `json:"urlIntranet,omitempty" xml:"urlIntranet,omitempty"`
}

func (s UpdateTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpdateTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *UpdateTriggerResponseBody) SetCreatedTime(v string) *UpdateTriggerResponseBody {
	s.CreatedTime = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetDescription(v string) *UpdateTriggerResponseBody {
	s.Description = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetDomainName(v string) *UpdateTriggerResponseBody {
	s.DomainName = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetInvocationRole(v string) *UpdateTriggerResponseBody {
	s.InvocationRole = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetLastModifiedTime(v string) *UpdateTriggerResponseBody {
	s.LastModifiedTime = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetQualifier(v string) *UpdateTriggerResponseBody {
	s.Qualifier = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetSourceArn(v string) *UpdateTriggerResponseBody {
	s.SourceArn = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetTriggerConfig(v string) *UpdateTriggerResponseBody {
	s.TriggerConfig = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetTriggerId(v string) *UpdateTriggerResponseBody {
	s.TriggerId = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetTriggerName(v string) *UpdateTriggerResponseBody {
	s.TriggerName = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetTriggerType(v string) *UpdateTriggerResponseBody {
	s.TriggerType = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetUrlInternet(v string) *UpdateTriggerResponseBody {
	s.UrlInternet = &v
	return s
}

func (s *UpdateTriggerResponseBody) SetUrlIntranet(v string) *UpdateTriggerResponseBody {
	s.UrlIntranet = &v
	return s
}

type UpdateTriggerResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpdateTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpdateTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateTriggerResponse) GoString() string {
	return s.String()
}

func (s *UpdateTriggerResponse) SetHeaders(v map[string]*string) *UpdateTriggerResponse {
	s.Headers = v
	return s
}

func (s *UpdateTriggerResponse) SetStatusCode(v int32) *UpdateTriggerResponse {
	s.StatusCode = &v
	return s
}

func (s *UpdateTriggerResponse) SetBody(v *UpdateTriggerResponseBody) *UpdateTriggerResponse {
	s.Body = v
	return s
}

type Client struct {
	openapi.Client
}

func NewClient(config *openapi.Config) (*Client, error) {
	client := new(Client)
	err := client.Init(config)
	return client, err
}

func (client *Client) Init(config *openapi.Config) (_err error) {
	_err = client.Client.Init(config)
	if _err != nil {
		return _err
	}
	client.EndpointRule = tea.String("regional")
	client.EndpointMap = map[string]*string{
		"ap-northeast-1":      tea.String("account-id.ap-northeast-1.fc.aliyuncs.com"),
		"ap-south-1":          tea.String("account-id.ap-south-1.fc.aliyuncs.com"),
		"ap-southeast-1":      tea.String("account-id.ap-southeast-1.fc.aliyuncs.com"),
		"ap-southeast-2":      tea.String("account-id.ap-southeast-2.fc.aliyuncs.com"),
		"ap-southeast-3":      tea.String("account-id.ap-southeast-3.fc.aliyuncs.com"),
		"ap-southeast-5":      tea.String("account-id.ap-southeast-5.fc.aliyuncs.com"),
		"cn-beijing":          tea.String("account-id.cn-beijing.fc.aliyuncs.com"),
		"cn-chengdu":          tea.String("account-id.cn-chengdu.fc.aliyuncs.com"),
		"cn-hangzhou":         tea.String("account-id.cn-hangzhou.fc.aliyuncs.com"),
		"cn-hangzhou-finance": tea.String("account-id.cn-hangzhou-finance.fc.aliyuncs.com"),
		"cn-hongkong":         tea.String("account-id.cn-hongkong.fc.aliyuncs.com"),
		"cn-huhehaote":        tea.String("account-id.cn-huhehaote.fc.aliyuncs.com"),
		"cn-north-2-gov-1":    tea.String("account-id.cn-north-2-gov-1.fc.aliyuncs.com"),
		"cn-qingdao":          tea.String("account-id.cn-qingdao.fc.aliyuncs.com"),
		"cn-shanghai":         tea.String("account-id.cn-shanghai.fc.aliyuncs.com"),
		"cn-shenzhen":         tea.String("account-id.cn-shenzhen.fc.aliyuncs.com"),
		"cn-zhangjiakou":      tea.String("account-id.cn-zhangjiakou.fc.aliyuncs.com"),
		"eu-central-1":        tea.String("account-id.eu-central-1.fc.aliyuncs.com"),
		"eu-west-1":           tea.String("account-id.eu-west-1.fc.aliyuncs.com"),
		"us-east-1":           tea.String("account-id.us-east-1.fc.aliyuncs.com"),
		"us-west-1":           tea.String("account-id.us-west-1.fc.aliyuncs.com"),
	}
	_err = client.CheckConfig(config)
	if _err != nil {
		return _err
	}
	client.Endpoint, _err = client.GetEndpoint(tea.String("fc-open"), client.RegionId, client.EndpointRule, client.Network, client.Suffix, client.EndpointMap, client.Endpoint)
	if _err != nil {
		return _err
	}

	return nil
}

func (client *Client) GetEndpoint(productId *string, regionId *string, endpointRule *string, network *string, suffix *string, endpointMap map[string]*string, endpoint *string) (_result *string, _err error) {
	if !tea.BoolValue(util.Empty(endpoint)) {
		_result = endpoint
		return _result, _err
	}

	if !tea.BoolValue(util.IsUnset(endpointMap)) && !tea.BoolValue(util.Empty(endpointMap[tea.StringValue(regionId)])) {
		_result = endpointMap[tea.StringValue(regionId)]
		return _result, _err
	}

	_body, _err := endpointutil.GetEndpointRules(productId, regionId, endpointRule, network, suffix)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ClaimGPUInstanceWithOptions(request *ClaimGPUInstanceRequest, headers *ClaimGPUInstanceHeaders, runtime *util.RuntimeOptions) (_result *ClaimGPUInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DiskPerformanceLevel)) {
		body["diskPerformanceLevel"] = request.DiskPerformanceLevel
	}

	if !tea.BoolValue(util.IsUnset(request.DiskSizeGigabytes)) {
		body["diskSizeGigabytes"] = request.DiskSizeGigabytes
	}

	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["imageId"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceType)) {
		body["instanceType"] = request.InstanceType
	}

	if !tea.BoolValue(util.IsUnset(request.InternetBandwidthOut)) {
		body["internetBandwidthOut"] = request.InternetBandwidthOut
	}

	if !tea.BoolValue(util.IsUnset(request.Password)) {
		body["password"] = request.Password
	}

	if !tea.BoolValue(util.IsUnset(request.Role)) {
		body["role"] = request.Role
	}

	if !tea.BoolValue(util.IsUnset(request.SgId)) {
		body["sgId"] = request.SgId
	}

	if !tea.BoolValue(util.IsUnset(request.SourceCidrIp)) {
		body["sourceCidrIp"] = request.SourceCidrIp
	}

	if !tea.BoolValue(util.IsUnset(request.TcpPortRange)) {
		body["tcpPortRange"] = request.TcpPortRange
	}

	if !tea.BoolValue(util.IsUnset(request.UdpPortRange)) {
		body["udpPortRange"] = request.UdpPortRange
	}

	if !tea.BoolValue(util.IsUnset(request.VpcId)) {
		body["vpcId"] = request.VpcId
	}

	if !tea.BoolValue(util.IsUnset(request.VswId)) {
		body["vswId"] = request.VswId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ClaimGPUInstance"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/gpuInstances"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ClaimGPUInstanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ClaimGPUInstance(request *ClaimGPUInstanceRequest) (_result *ClaimGPUInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ClaimGPUInstanceHeaders{}
	_result = &ClaimGPUInstanceResponse{}
	_body, _err := client.ClaimGPUInstanceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateAliasWithOptions(serviceName *string, request *CreateAliasRequest, headers *CreateAliasHeaders, runtime *util.RuntimeOptions) (_result *CreateAliasResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AdditionalVersionWeight)) {
		body["additionalVersionWeight"] = request.AdditionalVersionWeight
	}

	if !tea.BoolValue(util.IsUnset(request.AliasName)) {
		body["aliasName"] = request.AliasName
	}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.ResolvePolicy)) {
		body["resolvePolicy"] = request.ResolvePolicy
	}

	if !tea.BoolValue(util.IsUnset(request.RoutePolicy)) {
		body["routePolicy"] = request.RoutePolicy
	}

	if !tea.BoolValue(util.IsUnset(request.VersionId)) {
		body["versionId"] = request.VersionId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateAlias"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/aliases"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateAliasResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateAlias(serviceName *string, request *CreateAliasRequest) (_result *CreateAliasResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateAliasHeaders{}
	_result = &CreateAliasResponse{}
	_body, _err := client.CreateAliasWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateCustomDomainWithOptions(request *CreateCustomDomainRequest, headers *CreateCustomDomainHeaders, runtime *util.RuntimeOptions) (_result *CreateCustomDomainResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CertConfig)) {
		body["certConfig"] = request.CertConfig
	}

	if !tea.BoolValue(util.IsUnset(request.DomainName)) {
		body["domainName"] = request.DomainName
	}

	if !tea.BoolValue(util.IsUnset(request.Protocol)) {
		body["protocol"] = request.Protocol
	}

	if !tea.BoolValue(util.IsUnset(request.RouteConfig)) {
		body["routeConfig"] = request.RouteConfig
	}

	if !tea.BoolValue(util.IsUnset(request.TlsConfig)) {
		body["tlsConfig"] = request.TlsConfig
	}

	if !tea.BoolValue(util.IsUnset(request.WafConfig)) {
		body["wafConfig"] = request.WafConfig
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateCustomDomain"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/custom-domains"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateCustomDomainResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateCustomDomain(request *CreateCustomDomainRequest) (_result *CreateCustomDomainResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateCustomDomainHeaders{}
	_result = &CreateCustomDomainResponse{}
	_body, _err := client.CreateCustomDomainWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateFunctionWithOptions(serviceName *string, request *CreateFunctionRequest, headers *CreateFunctionHeaders, runtime *util.RuntimeOptions) (_result *CreateFunctionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CaPort)) {
		body["caPort"] = request.CaPort
	}

	if !tea.BoolValue(util.IsUnset(request.Code)) {
		body["code"] = request.Code
	}

	if !tea.BoolValue(util.IsUnset(request.Cpu)) {
		body["cpu"] = request.Cpu
	}

	if !tea.BoolValue(util.IsUnset(request.CustomContainerConfig)) {
		body["customContainerConfig"] = request.CustomContainerConfig
	}

	if !tea.BoolValue(util.IsUnset(request.CustomDNS)) {
		body["customDNS"] = request.CustomDNS
	}

	if !tea.BoolValue(util.IsUnset(request.CustomHealthCheckConfig)) {
		body["customHealthCheckConfig"] = request.CustomHealthCheckConfig
	}

	if !tea.BoolValue(util.IsUnset(request.CustomRuntimeConfig)) {
		body["customRuntimeConfig"] = request.CustomRuntimeConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.DiskSize)) {
		body["diskSize"] = request.DiskSize
	}

	if !tea.BoolValue(util.IsUnset(request.EnvironmentVariables)) {
		body["environmentVariables"] = request.EnvironmentVariables
	}

	if !tea.BoolValue(util.IsUnset(request.FunctionName)) {
		body["functionName"] = request.FunctionName
	}

	if !tea.BoolValue(util.IsUnset(request.GpuMemorySize)) {
		body["gpuMemorySize"] = request.GpuMemorySize
	}

	if !tea.BoolValue(util.IsUnset(request.Handler)) {
		body["handler"] = request.Handler
	}

	if !tea.BoolValue(util.IsUnset(request.InitializationTimeout)) {
		body["initializationTimeout"] = request.InitializationTimeout
	}

	if !tea.BoolValue(util.IsUnset(request.Initializer)) {
		body["initializer"] = request.Initializer
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceConcurrency)) {
		body["instanceConcurrency"] = request.InstanceConcurrency
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceLifecycleConfig)) {
		body["instanceLifecycleConfig"] = request.InstanceLifecycleConfig
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceSoftConcurrency)) {
		body["instanceSoftConcurrency"] = request.InstanceSoftConcurrency
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceType)) {
		body["instanceType"] = request.InstanceType
	}

	if !tea.BoolValue(util.IsUnset(request.Layers)) {
		body["layers"] = request.Layers
	}

	if !tea.BoolValue(util.IsUnset(request.MemorySize)) {
		body["memorySize"] = request.MemorySize
	}

	if !tea.BoolValue(util.IsUnset(request.Runtime)) {
		body["runtime"] = request.Runtime
	}

	if !tea.BoolValue(util.IsUnset(request.Timeout)) {
		body["timeout"] = request.Timeout
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcCodeChecksum)) {
		realHeaders["X-Fc-Code-Checksum"] = util.ToJSONString(headers.XFcCodeChecksum)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateFunction"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateFunctionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateFunction(serviceName *string, request *CreateFunctionRequest) (_result *CreateFunctionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateFunctionHeaders{}
	_result = &CreateFunctionResponse{}
	_body, _err := client.CreateFunctionWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateLayerVersionWithOptions(layerName *string, request *CreateLayerVersionRequest, headers *CreateLayerVersionHeaders, runtime *util.RuntimeOptions) (_result *CreateLayerVersionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Code)) {
		body["Code"] = request.Code
	}

	if !tea.BoolValue(util.IsUnset(request.CompatibleRuntime)) {
		body["compatibleRuntime"] = request.CompatibleRuntime
	}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateLayerVersion"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/layers/" + tea.StringValue(openapiutil.GetEncodeParam(layerName)) + "/versions"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateLayerVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateLayerVersion(layerName *string, request *CreateLayerVersionRequest) (_result *CreateLayerVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateLayerVersionHeaders{}
	_result = &CreateLayerVersionResponse{}
	_body, _err := client.CreateLayerVersionWithOptions(layerName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateServiceWithOptions(request *CreateServiceRequest, headers *CreateServiceHeaders, runtime *util.RuntimeOptions) (_result *CreateServiceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.InternetAccess)) {
		body["internetAccess"] = request.InternetAccess
	}

	if !tea.BoolValue(util.IsUnset(request.LogConfig)) {
		body["logConfig"] = request.LogConfig
	}

	if !tea.BoolValue(util.IsUnset(request.NasConfig)) {
		body["nasConfig"] = request.NasConfig
	}

	if !tea.BoolValue(util.IsUnset(request.OssMountConfig)) {
		body["ossMountConfig"] = request.OssMountConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Role)) {
		body["role"] = request.Role
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceName)) {
		body["serviceName"] = request.ServiceName
	}

	if !tea.BoolValue(util.IsUnset(request.TracingConfig)) {
		body["tracingConfig"] = request.TracingConfig
	}

	if !tea.BoolValue(util.IsUnset(request.VpcConfig)) {
		body["vpcConfig"] = request.VpcConfig
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateService"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateServiceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateService(request *CreateServiceRequest) (_result *CreateServiceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateServiceHeaders{}
	_result = &CreateServiceResponse{}
	_body, _err := client.CreateServiceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateTriggerWithOptions(serviceName *string, functionName *string, request *CreateTriggerRequest, headers *CreateTriggerHeaders, runtime *util.RuntimeOptions) (_result *CreateTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.InvocationRole)) {
		body["invocationRole"] = request.InvocationRole
	}

	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		body["qualifier"] = request.Qualifier
	}

	if !tea.BoolValue(util.IsUnset(request.SourceArn)) {
		body["sourceArn"] = request.SourceArn
	}

	if !tea.BoolValue(util.IsUnset(request.TriggerConfig)) {
		body["triggerConfig"] = request.TriggerConfig
	}

	if !tea.BoolValue(util.IsUnset(request.TriggerName)) {
		body["triggerName"] = request.TriggerName
	}

	if !tea.BoolValue(util.IsUnset(request.TriggerType)) {
		body["triggerType"] = request.TriggerType
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateTrigger"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/triggers"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateTrigger(serviceName *string, functionName *string, request *CreateTriggerRequest) (_result *CreateTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateTriggerHeaders{}
	_result = &CreateTriggerResponse{}
	_body, _err := client.CreateTriggerWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateVpcBindingWithOptions(serviceName *string, request *CreateVpcBindingRequest, headers *CreateVpcBindingHeaders, runtime *util.RuntimeOptions) (_result *CreateVpcBindingResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.VpcId)) {
		body["vpcId"] = request.VpcId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateVpcBinding"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/bindings"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &CreateVpcBindingResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateVpcBinding(serviceName *string, request *CreateVpcBindingRequest) (_result *CreateVpcBindingResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &CreateVpcBindingHeaders{}
	_result = &CreateVpcBindingResponse{}
	_body, _err := client.CreateVpcBindingWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteAliasWithOptions(serviceName *string, aliasName *string, headers *DeleteAliasHeaders, runtime *util.RuntimeOptions) (_result *DeleteAliasResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteAlias"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/aliases/" + tea.StringValue(openapiutil.GetEncodeParam(aliasName))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteAliasResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteAlias(serviceName *string, aliasName *string) (_result *DeleteAliasResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteAliasHeaders{}
	_result = &DeleteAliasResponse{}
	_body, _err := client.DeleteAliasWithOptions(serviceName, aliasName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteCustomDomainWithOptions(domainName *string, headers *DeleteCustomDomainHeaders, runtime *util.RuntimeOptions) (_result *DeleteCustomDomainResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteCustomDomain"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/custom-domains/" + tea.StringValue(openapiutil.GetEncodeParam(domainName))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteCustomDomainResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteCustomDomain(domainName *string) (_result *DeleteCustomDomainResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteCustomDomainHeaders{}
	_result = &DeleteCustomDomainResponse{}
	_body, _err := client.DeleteCustomDomainWithOptions(domainName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteFunctionWithOptions(serviceName *string, functionName *string, headers *DeleteFunctionHeaders, runtime *util.RuntimeOptions) (_result *DeleteFunctionResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteFunction"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteFunctionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteFunction(serviceName *string, functionName *string) (_result *DeleteFunctionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteFunctionHeaders{}
	_result = &DeleteFunctionResponse{}
	_body, _err := client.DeleteFunctionWithOptions(serviceName, functionName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteFunctionAsyncInvokeConfigWithOptions(serviceName *string, functionName *string, request *DeleteFunctionAsyncInvokeConfigRequest, headers *DeleteFunctionAsyncInvokeConfigHeaders, runtime *util.RuntimeOptions) (_result *DeleteFunctionAsyncInvokeConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteFunctionAsyncInvokeConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/async-invoke-config"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteFunctionAsyncInvokeConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteFunctionAsyncInvokeConfig(serviceName *string, functionName *string, request *DeleteFunctionAsyncInvokeConfigRequest) (_result *DeleteFunctionAsyncInvokeConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteFunctionAsyncInvokeConfigHeaders{}
	_result = &DeleteFunctionAsyncInvokeConfigResponse{}
	_body, _err := client.DeleteFunctionAsyncInvokeConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteFunctionOnDemandConfigWithOptions(serviceName *string, functionName *string, request *DeleteFunctionOnDemandConfigRequest, headers *DeleteFunctionOnDemandConfigHeaders, runtime *util.RuntimeOptions) (_result *DeleteFunctionOnDemandConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteFunctionOnDemandConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/on-demand-config"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteFunctionOnDemandConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteFunctionOnDemandConfig(serviceName *string, functionName *string, request *DeleteFunctionOnDemandConfigRequest) (_result *DeleteFunctionOnDemandConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteFunctionOnDemandConfigHeaders{}
	_result = &DeleteFunctionOnDemandConfigResponse{}
	_body, _err := client.DeleteFunctionOnDemandConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteLayerVersionWithOptions(layerName *string, version *string, headers *DeleteLayerVersionHeaders, runtime *util.RuntimeOptions) (_result *DeleteLayerVersionResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteLayerVersion"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/layers/" + tea.StringValue(openapiutil.GetEncodeParam(layerName)) + "/versions/" + tea.StringValue(openapiutil.GetEncodeParam(version))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteLayerVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteLayerVersion(layerName *string, version *string) (_result *DeleteLayerVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteLayerVersionHeaders{}
	_result = &DeleteLayerVersionResponse{}
	_body, _err := client.DeleteLayerVersionWithOptions(layerName, version, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteProvisionConfigWithOptions(serviceName *string, functionName *string, request *DeleteProvisionConfigRequest, headers *DeleteProvisionConfigHeaders, runtime *util.RuntimeOptions) (_result *DeleteProvisionConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteProvisionConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/provision-config"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteProvisionConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteProvisionConfig(serviceName *string, functionName *string, request *DeleteProvisionConfigRequest) (_result *DeleteProvisionConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteProvisionConfigHeaders{}
	_result = &DeleteProvisionConfigResponse{}
	_body, _err := client.DeleteProvisionConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteServiceWithOptions(serviceName *string, headers *DeleteServiceHeaders, runtime *util.RuntimeOptions) (_result *DeleteServiceResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteService"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteServiceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteService(serviceName *string) (_result *DeleteServiceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteServiceHeaders{}
	_result = &DeleteServiceResponse{}
	_body, _err := client.DeleteServiceWithOptions(serviceName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteServiceVersionWithOptions(serviceName *string, versionId *string, headers *DeleteServiceVersionHeaders, runtime *util.RuntimeOptions) (_result *DeleteServiceVersionResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteServiceVersion"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/versions/" + tea.StringValue(openapiutil.GetEncodeParam(versionId))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteServiceVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteServiceVersion(serviceName *string, versionId *string) (_result *DeleteServiceVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteServiceVersionHeaders{}
	_result = &DeleteServiceVersionResponse{}
	_body, _err := client.DeleteServiceVersionWithOptions(serviceName, versionId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteTriggerWithOptions(serviceName *string, functionName *string, triggerName *string, headers *DeleteTriggerHeaders, runtime *util.RuntimeOptions) (_result *DeleteTriggerResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteTrigger"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/triggers/" + tea.StringValue(openapiutil.GetEncodeParam(triggerName))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteTrigger(serviceName *string, functionName *string, triggerName *string) (_result *DeleteTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteTriggerHeaders{}
	_result = &DeleteTriggerResponse{}
	_body, _err := client.DeleteTriggerWithOptions(serviceName, functionName, triggerName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteVpcBindingWithOptions(serviceName *string, vpcId *string, headers *DeleteVpcBindingHeaders, runtime *util.RuntimeOptions) (_result *DeleteVpcBindingResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteVpcBinding"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/bindings/" + tea.StringValue(openapiutil.GetEncodeParam(vpcId))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteVpcBindingResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteVpcBinding(serviceName *string, vpcId *string) (_result *DeleteVpcBindingResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeleteVpcBindingHeaders{}
	_result = &DeleteVpcBindingResponse{}
	_body, _err := client.DeleteVpcBindingWithOptions(serviceName, vpcId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeregisterEventSourceWithOptions(serviceName *string, functionName *string, sourceArn *string, request *DeregisterEventSourceRequest, headers *DeregisterEventSourceHeaders, runtime *util.RuntimeOptions) (_result *DeregisterEventSourceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeregisterEventSource"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/event-sources/" + tea.StringValue(openapiutil.GetEncodeParam(sourceArn))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeregisterEventSourceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeregisterEventSource(serviceName *string, functionName *string, sourceArn *string, request *DeregisterEventSourceRequest) (_result *DeregisterEventSourceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &DeregisterEventSourceHeaders{}
	_result = &DeregisterEventSourceResponse{}
	_body, _err := client.DeregisterEventSourceWithOptions(serviceName, functionName, sourceArn, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetAccountSettingsWithOptions(headers *GetAccountSettingsHeaders, runtime *util.RuntimeOptions) (_result *GetAccountSettingsResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("GetAccountSettings"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/account-settings"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetAccountSettingsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetAccountSettings() (_result *GetAccountSettingsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetAccountSettingsHeaders{}
	_result = &GetAccountSettingsResponse{}
	_body, _err := client.GetAccountSettingsWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetAliasWithOptions(serviceName *string, aliasName *string, headers *GetAliasHeaders, runtime *util.RuntimeOptions) (_result *GetAliasResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("GetAlias"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/aliases/" + tea.StringValue(openapiutil.GetEncodeParam(aliasName))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetAliasResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetAlias(serviceName *string, aliasName *string) (_result *GetAliasResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetAliasHeaders{}
	_result = &GetAliasResponse{}
	_body, _err := client.GetAliasWithOptions(serviceName, aliasName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetCustomDomainWithOptions(domainName *string, headers *GetCustomDomainHeaders, runtime *util.RuntimeOptions) (_result *GetCustomDomainResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("GetCustomDomain"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/custom-domains/" + tea.StringValue(openapiutil.GetEncodeParam(domainName))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetCustomDomainResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetCustomDomain(domainName *string) (_result *GetCustomDomainResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetCustomDomainHeaders{}
	_result = &GetCustomDomainResponse{}
	_body, _err := client.GetCustomDomainWithOptions(domainName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetFunctionWithOptions(serviceName *string, functionName *string, request *GetFunctionRequest, headers *GetFunctionHeaders, runtime *util.RuntimeOptions) (_result *GetFunctionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetFunction"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetFunctionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetFunction(serviceName *string, functionName *string, request *GetFunctionRequest) (_result *GetFunctionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetFunctionHeaders{}
	_result = &GetFunctionResponse{}
	_body, _err := client.GetFunctionWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * StatefulAsyncInvocation indicates whether the asynchronous task feature is enabled. If the value of StatefulAsyncInvocation is true, the asynchronous task feature is enabled. All asynchronous invocations change to asynchronous task mode.
 *
 * @param request GetFunctionAsyncInvokeConfigRequest
 * @param headers GetFunctionAsyncInvokeConfigHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return GetFunctionAsyncInvokeConfigResponse
 */
func (client *Client) GetFunctionAsyncInvokeConfigWithOptions(serviceName *string, functionName *string, request *GetFunctionAsyncInvokeConfigRequest, headers *GetFunctionAsyncInvokeConfigHeaders, runtime *util.RuntimeOptions) (_result *GetFunctionAsyncInvokeConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetFunctionAsyncInvokeConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/async-invoke-config"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetFunctionAsyncInvokeConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * StatefulAsyncInvocation indicates whether the asynchronous task feature is enabled. If the value of StatefulAsyncInvocation is true, the asynchronous task feature is enabled. All asynchronous invocations change to asynchronous task mode.
 *
 * @param request GetFunctionAsyncInvokeConfigRequest
 * @return GetFunctionAsyncInvokeConfigResponse
 */
func (client *Client) GetFunctionAsyncInvokeConfig(serviceName *string, functionName *string, request *GetFunctionAsyncInvokeConfigRequest) (_result *GetFunctionAsyncInvokeConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetFunctionAsyncInvokeConfigHeaders{}
	_result = &GetFunctionAsyncInvokeConfigResponse{}
	_body, _err := client.GetFunctionAsyncInvokeConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetFunctionCodeWithOptions(serviceName *string, functionName *string, request *GetFunctionCodeRequest, headers *GetFunctionCodeHeaders, runtime *util.RuntimeOptions) (_result *GetFunctionCodeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetFunctionCode"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/code"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetFunctionCodeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetFunctionCode(serviceName *string, functionName *string, request *GetFunctionCodeRequest) (_result *GetFunctionCodeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetFunctionCodeHeaders{}
	_result = &GetFunctionCodeResponse{}
	_body, _err := client.GetFunctionCodeWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetFunctionOnDemandConfigWithOptions(serviceName *string, functionName *string, request *GetFunctionOnDemandConfigRequest, headers *GetFunctionOnDemandConfigHeaders, runtime *util.RuntimeOptions) (_result *GetFunctionOnDemandConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetFunctionOnDemandConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/on-demand-config"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetFunctionOnDemandConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetFunctionOnDemandConfig(serviceName *string, functionName *string, request *GetFunctionOnDemandConfigRequest) (_result *GetFunctionOnDemandConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetFunctionOnDemandConfigHeaders{}
	_result = &GetFunctionOnDemandConfigResponse{}
	_body, _err := client.GetFunctionOnDemandConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetLayerVersionWithOptions(layerName *string, version *string, headers *GetLayerVersionHeaders, runtime *util.RuntimeOptions) (_result *GetLayerVersionResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("GetLayerVersion"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/layers/" + tea.StringValue(openapiutil.GetEncodeParam(layerName)) + "/versions/" + tea.StringValue(openapiutil.GetEncodeParam(version))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetLayerVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetLayerVersion(layerName *string, version *string) (_result *GetLayerVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetLayerVersionHeaders{}
	_result = &GetLayerVersionResponse{}
	_body, _err := client.GetLayerVersionWithOptions(layerName, version, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetProvisionConfigWithOptions(serviceName *string, functionName *string, request *GetProvisionConfigRequest, headers *GetProvisionConfigHeaders, runtime *util.RuntimeOptions) (_result *GetProvisionConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetProvisionConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/provision-config"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetProvisionConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetProvisionConfig(serviceName *string, functionName *string, request *GetProvisionConfigRequest) (_result *GetProvisionConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetProvisionConfigHeaders{}
	_result = &GetProvisionConfigResponse{}
	_body, _err := client.GetProvisionConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetResourceTagsWithOptions(request *GetResourceTagsRequest, headers *GetResourceTagsHeaders, runtime *util.RuntimeOptions) (_result *GetResourceTagsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ResourceArn)) {
		query["resourceArn"] = request.ResourceArn
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetResourceTags"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/tag"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetResourceTagsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetResourceTags(request *GetResourceTagsRequest) (_result *GetResourceTagsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetResourceTagsHeaders{}
	_result = &GetResourceTagsResponse{}
	_body, _err := client.GetResourceTagsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetServiceWithOptions(serviceName *string, request *GetServiceRequest, headers *GetServiceHeaders, runtime *util.RuntimeOptions) (_result *GetServiceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetService"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetServiceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetService(serviceName *string, request *GetServiceRequest) (_result *GetServiceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetServiceHeaders{}
	_result = &GetServiceResponse{}
	_body, _err := client.GetServiceWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request GetStatefulAsyncInvocationRequest
 * @param headers GetStatefulAsyncInvocationHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return GetStatefulAsyncInvocationResponse
 */
func (client *Client) GetStatefulAsyncInvocationWithOptions(serviceName *string, functionName *string, invocationId *string, request *GetStatefulAsyncInvocationRequest, headers *GetStatefulAsyncInvocationHeaders, runtime *util.RuntimeOptions) (_result *GetStatefulAsyncInvocationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcCodeChecksum)) {
		realHeaders["X-Fc-Code-Checksum"] = util.ToJSONString(headers.XFcCodeChecksum)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcInvocationType)) {
		realHeaders["X-Fc-Invocation-Type"] = util.ToJSONString(headers.XFcInvocationType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcLogType)) {
		realHeaders["X-Fc-Log-Type"] = util.ToJSONString(headers.XFcLogType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetStatefulAsyncInvocation"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/stateful-async-invocations/" + tea.StringValue(openapiutil.GetEncodeParam(invocationId))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetStatefulAsyncInvocationResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request GetStatefulAsyncInvocationRequest
 * @return GetStatefulAsyncInvocationResponse
 */
func (client *Client) GetStatefulAsyncInvocation(serviceName *string, functionName *string, invocationId *string, request *GetStatefulAsyncInvocationRequest) (_result *GetStatefulAsyncInvocationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetStatefulAsyncInvocationHeaders{}
	_result = &GetStatefulAsyncInvocationResponse{}
	_body, _err := client.GetStatefulAsyncInvocationWithOptions(serviceName, functionName, invocationId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetTriggerWithOptions(serviceName *string, functionName *string, triggerName *string, headers *GetTriggerHeaders, runtime *util.RuntimeOptions) (_result *GetTriggerResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("GetTrigger"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/triggers/" + tea.StringValue(openapiutil.GetEncodeParam(triggerName))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetTrigger(serviceName *string, functionName *string, triggerName *string) (_result *GetTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &GetTriggerHeaders{}
	_result = &GetTriggerResponse{}
	_body, _err := client.GetTriggerWithOptions(serviceName, functionName, triggerName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) InvokeFunctionWithOptions(serviceName *string, functionName *string, request *InvokeFunctionRequest, headers *InvokeFunctionHeaders, runtime *util.RuntimeOptions) (_result *InvokeFunctionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	body := tea.String("")
	if !tea.BoolValue(util.IsUnset(request.Body)) {
		body = util.ToString(request.Body)
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcInvocationType)) {
		realHeaders["X-Fc-Invocation-Type"] = util.ToJSONString(headers.XFcInvocationType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcLogType)) {
		realHeaders["X-Fc-Log-Type"] = util.ToJSONString(headers.XFcLogType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcStatefulAsyncInvocationId)) {
		realHeaders["X-Fc-Stateful-Async-Invocation-Id"] = util.ToJSONString(headers.XFcStatefulAsyncInvocationId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
		Body:    body,
	}
	params := &openapi.Params{
		Action:      tea.String("InvokeFunction"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/invocations"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("byte"),
	}
	_result = &InvokeFunctionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) InvokeFunction(serviceName *string, functionName *string, request *InvokeFunctionRequest) (_result *InvokeFunctionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &InvokeFunctionHeaders{}
	_result = &InvokeFunctionResponse{}
	_body, _err := client.InvokeFunctionWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListAliasesWithOptions(serviceName *string, request *ListAliasesRequest, headers *ListAliasesHeaders, runtime *util.RuntimeOptions) (_result *ListAliasesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Prefix)) {
		query["prefix"] = request.Prefix
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListAliases"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/aliases"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListAliasesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListAliases(serviceName *string, request *ListAliasesRequest) (_result *ListAliasesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListAliasesHeaders{}
	_result = &ListAliasesResponse{}
	_body, _err := client.ListAliasesWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListCustomDomainsWithOptions(request *ListCustomDomainsRequest, headers *ListCustomDomainsHeaders, runtime *util.RuntimeOptions) (_result *ListCustomDomainsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Prefix)) {
		query["prefix"] = request.Prefix
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListCustomDomains"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/custom-domains"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListCustomDomainsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListCustomDomains(request *ListCustomDomainsRequest) (_result *ListCustomDomainsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListCustomDomainsHeaders{}
	_result = &ListCustomDomainsResponse{}
	_body, _err := client.ListCustomDomainsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListEventSourcesWithOptions(serviceName *string, functionName *string, request *ListEventSourcesRequest, headers *ListEventSourcesHeaders, runtime *util.RuntimeOptions) (_result *ListEventSourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListEventSources"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/event-sources"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListEventSourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListEventSources(serviceName *string, functionName *string, request *ListEventSourcesRequest) (_result *ListEventSourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListEventSourcesHeaders{}
	_result = &ListEventSourcesResponse{}
	_body, _err := client.ListEventSourcesWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * StatefulAsyncInvocation indicates whether the asynchronous task feature is enabled. If StatefulAsyncInvocation is set to true, the asynchronous task is enabled. All asynchronous invocations to the function corresponding to this configuration change to asynchronous task mode.
 *
 * @param request ListFunctionAsyncInvokeConfigsRequest
 * @param headers ListFunctionAsyncInvokeConfigsHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return ListFunctionAsyncInvokeConfigsResponse
 */
func (client *Client) ListFunctionAsyncInvokeConfigsWithOptions(serviceName *string, functionName *string, request *ListFunctionAsyncInvokeConfigsRequest, headers *ListFunctionAsyncInvokeConfigsHeaders, runtime *util.RuntimeOptions) (_result *ListFunctionAsyncInvokeConfigsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcCodeChecksum)) {
		realHeaders["X-Fc-Code-Checksum"] = util.ToJSONString(headers.XFcCodeChecksum)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcInvocationType)) {
		realHeaders["X-Fc-Invocation-Type"] = util.ToJSONString(headers.XFcInvocationType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcLogType)) {
		realHeaders["X-Fc-Log-Type"] = util.ToJSONString(headers.XFcLogType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListFunctionAsyncInvokeConfigs"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/async-invoke-configs"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListFunctionAsyncInvokeConfigsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * StatefulAsyncInvocation indicates whether the asynchronous task feature is enabled. If StatefulAsyncInvocation is set to true, the asynchronous task is enabled. All asynchronous invocations to the function corresponding to this configuration change to asynchronous task mode.
 *
 * @param request ListFunctionAsyncInvokeConfigsRequest
 * @return ListFunctionAsyncInvokeConfigsResponse
 */
func (client *Client) ListFunctionAsyncInvokeConfigs(serviceName *string, functionName *string, request *ListFunctionAsyncInvokeConfigsRequest) (_result *ListFunctionAsyncInvokeConfigsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListFunctionAsyncInvokeConfigsHeaders{}
	_result = &ListFunctionAsyncInvokeConfigsResponse{}
	_body, _err := client.ListFunctionAsyncInvokeConfigsWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListFunctionsWithOptions(serviceName *string, request *ListFunctionsRequest, headers *ListFunctionsHeaders, runtime *util.RuntimeOptions) (_result *ListFunctionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Prefix)) {
		query["prefix"] = request.Prefix
	}

	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListFunctions"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListFunctionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListFunctions(serviceName *string, request *ListFunctionsRequest) (_result *ListFunctionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListFunctionsHeaders{}
	_result = &ListFunctionsResponse{}
	_body, _err := client.ListFunctionsWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * The ListInstances operation allows you to query the available instances of a function.
 * Available instances are instances that are processing requests or can be scheduled to process requests. Available instances queried by the ListInstances operation are the same as those that can be used when you call the InvokeFunction operation with the same values specified for the `serviceName`, `functionName`, and `qualifier` parameters.
 *
 * @param request ListInstancesRequest
 * @param headers ListInstancesHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return ListInstancesResponse
 */
func (client *Client) ListInstancesWithOptions(serviceName *string, functionName *string, request *ListInstancesRequest, headers *ListInstancesHeaders, runtime *util.RuntimeOptions) (_result *ListInstancesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.InstanceIds)) {
		query["instanceIds"] = request.InstanceIds
	}

	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListInstances"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/instances"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListInstancesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * The ListInstances operation allows you to query the available instances of a function.
 * Available instances are instances that are processing requests or can be scheduled to process requests. Available instances queried by the ListInstances operation are the same as those that can be used when you call the InvokeFunction operation with the same values specified for the `serviceName`, `functionName`, and `qualifier` parameters.
 *
 * @param request ListInstancesRequest
 * @return ListInstancesResponse
 */
func (client *Client) ListInstances(serviceName *string, functionName *string, request *ListInstancesRequest) (_result *ListInstancesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListInstancesHeaders{}
	_result = &ListInstancesResponse{}
	_body, _err := client.ListInstancesWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListLayerVersionsWithOptions(layerName *string, request *ListLayerVersionsRequest, headers *ListLayerVersionsHeaders, runtime *util.RuntimeOptions) (_result *ListLayerVersionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.StartVersion)) {
		query["startVersion"] = request.StartVersion
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListLayerVersions"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/layers/" + tea.StringValue(openapiutil.GetEncodeParam(layerName)) + "/versions"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListLayerVersionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListLayerVersions(layerName *string, request *ListLayerVersionsRequest) (_result *ListLayerVersionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListLayerVersionsHeaders{}
	_result = &ListLayerVersionsResponse{}
	_body, _err := client.ListLayerVersionsWithOptions(layerName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListLayersWithOptions(request *ListLayersRequest, headers *ListLayersHeaders, runtime *util.RuntimeOptions) (_result *ListLayersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Official)) {
		query["official"] = request.Official
	}

	if !tea.BoolValue(util.IsUnset(request.Prefix)) {
		query["prefix"] = request.Prefix
	}

	if !tea.BoolValue(util.IsUnset(request.Public)) {
		query["public"] = request.Public
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListLayers"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/layers"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListLayersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListLayers(request *ListLayersRequest) (_result *ListLayersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListLayersHeaders{}
	_result = &ListLayersResponse{}
	_body, _err := client.ListLayersWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListOnDemandConfigsWithOptions(request *ListOnDemandConfigsRequest, headers *ListOnDemandConfigsHeaders, runtime *util.RuntimeOptions) (_result *ListOnDemandConfigsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Prefix)) {
		query["prefix"] = request.Prefix
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListOnDemandConfigs"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/on-demand-configs"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListOnDemandConfigsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListOnDemandConfigs(request *ListOnDemandConfigsRequest) (_result *ListOnDemandConfigsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListOnDemandConfigsHeaders{}
	_result = &ListOnDemandConfigsResponse{}
	_body, _err := client.ListOnDemandConfigsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListProvisionConfigsWithOptions(request *ListProvisionConfigsRequest, headers *ListProvisionConfigsHeaders, runtime *util.RuntimeOptions) (_result *ListProvisionConfigsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceName)) {
		query["serviceName"] = request.ServiceName
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListProvisionConfigs"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/provision-configs"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListProvisionConfigsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListProvisionConfigs(request *ListProvisionConfigsRequest) (_result *ListProvisionConfigsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListProvisionConfigsHeaders{}
	_result = &ListProvisionConfigsResponse{}
	_body, _err := client.ListProvisionConfigsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListReservedCapacitiesWithOptions(request *ListReservedCapacitiesRequest, headers *ListReservedCapacitiesHeaders, runtime *util.RuntimeOptions) (_result *ListReservedCapacitiesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListReservedCapacities"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/reserved-capacities"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListReservedCapacitiesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListReservedCapacities(request *ListReservedCapacitiesRequest) (_result *ListReservedCapacitiesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListReservedCapacitiesHeaders{}
	_result = &ListReservedCapacitiesResponse{}
	_body, _err := client.ListReservedCapacitiesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListServiceVersionsWithOptions(serviceName *string, request *ListServiceVersionsRequest, headers *ListServiceVersionsHeaders, runtime *util.RuntimeOptions) (_result *ListServiceVersionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Direction)) {
		query["direction"] = request.Direction
	}

	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListServiceVersions"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/versions"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListServiceVersionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListServiceVersions(serviceName *string, request *ListServiceVersionsRequest) (_result *ListServiceVersionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListServiceVersionsHeaders{}
	_result = &ListServiceVersionsResponse{}
	_body, _err := client.ListServiceVersionsWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListServicesWithOptions(request *ListServicesRequest, headers *ListServicesHeaders, runtime *util.RuntimeOptions) (_result *ListServicesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Prefix)) {
		query["prefix"] = request.Prefix
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListServices"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListServicesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListServices(request *ListServicesRequest) (_result *ListServicesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListServicesHeaders{}
	_result = &ListServicesResponse{}
	_body, _err := client.ListServicesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request ListStatefulAsyncInvocationFunctionsRequest
 * @param headers ListStatefulAsyncInvocationFunctionsHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return ListStatefulAsyncInvocationFunctionsResponse
 */
func (client *Client) ListStatefulAsyncInvocationFunctionsWithOptions(request *ListStatefulAsyncInvocationFunctionsRequest, headers *ListStatefulAsyncInvocationFunctionsHeaders, runtime *util.RuntimeOptions) (_result *ListStatefulAsyncInvocationFunctionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListStatefulAsyncInvocationFunctions"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/stateful-async-invocation-functions"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListStatefulAsyncInvocationFunctionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request ListStatefulAsyncInvocationFunctionsRequest
 * @return ListStatefulAsyncInvocationFunctionsResponse
 */
func (client *Client) ListStatefulAsyncInvocationFunctions(request *ListStatefulAsyncInvocationFunctionsRequest) (_result *ListStatefulAsyncInvocationFunctionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListStatefulAsyncInvocationFunctionsHeaders{}
	_result = &ListStatefulAsyncInvocationFunctionsResponse{}
	_body, _err := client.ListStatefulAsyncInvocationFunctionsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request ListStatefulAsyncInvocationsRequest
 * @param headers ListStatefulAsyncInvocationsHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return ListStatefulAsyncInvocationsResponse
 */
func (client *Client) ListStatefulAsyncInvocationsWithOptions(serviceName *string, functionName *string, request *ListStatefulAsyncInvocationsRequest, headers *ListStatefulAsyncInvocationsHeaders, runtime *util.RuntimeOptions) (_result *ListStatefulAsyncInvocationsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.IncludePayload)) {
		query["includePayload"] = request.IncludePayload
	}

	if !tea.BoolValue(util.IsUnset(request.InvocationIdPrefix)) {
		query["invocationIdPrefix"] = request.InvocationIdPrefix
	}

	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	if !tea.BoolValue(util.IsUnset(request.SortOrderByTime)) {
		query["sortOrderByTime"] = request.SortOrderByTime
	}

	if !tea.BoolValue(util.IsUnset(request.StartedTimeBegin)) {
		query["startedTimeBegin"] = request.StartedTimeBegin
	}

	if !tea.BoolValue(util.IsUnset(request.StartedTimeEnd)) {
		query["startedTimeEnd"] = request.StartedTimeEnd
	}

	if !tea.BoolValue(util.IsUnset(request.Status)) {
		query["status"] = request.Status
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcCodeChecksum)) {
		realHeaders["X-Fc-Code-Checksum"] = util.ToJSONString(headers.XFcCodeChecksum)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcInvocationType)) {
		realHeaders["X-Fc-Invocation-Type"] = util.ToJSONString(headers.XFcInvocationType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcLogType)) {
		realHeaders["X-Fc-Log-Type"] = util.ToJSONString(headers.XFcLogType)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListStatefulAsyncInvocations"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/stateful-async-invocations"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListStatefulAsyncInvocationsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request ListStatefulAsyncInvocationsRequest
 * @return ListStatefulAsyncInvocationsResponse
 */
func (client *Client) ListStatefulAsyncInvocations(serviceName *string, functionName *string, request *ListStatefulAsyncInvocationsRequest) (_result *ListStatefulAsyncInvocationsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListStatefulAsyncInvocationsHeaders{}
	_result = &ListStatefulAsyncInvocationsResponse{}
	_body, _err := client.ListStatefulAsyncInvocationsWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListTaggedResourcesWithOptions(request *ListTaggedResourcesRequest, headers *ListTaggedResourcesHeaders, runtime *util.RuntimeOptions) (_result *ListTaggedResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListTaggedResources"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/tags"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListTaggedResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListTaggedResources(request *ListTaggedResourcesRequest) (_result *ListTaggedResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListTaggedResourcesHeaders{}
	_result = &ListTaggedResourcesResponse{}
	_body, _err := client.ListTaggedResourcesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListTriggersWithOptions(serviceName *string, functionName *string, request *ListTriggersRequest, headers *ListTriggersHeaders, runtime *util.RuntimeOptions) (_result *ListTriggersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Limit)) {
		query["limit"] = request.Limit
	}

	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["nextToken"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.Prefix)) {
		query["prefix"] = request.Prefix
	}

	if !tea.BoolValue(util.IsUnset(request.StartKey)) {
		query["startKey"] = request.StartKey
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListTriggers"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/triggers"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListTriggersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListTriggers(serviceName *string, functionName *string, request *ListTriggersRequest) (_result *ListTriggersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListTriggersHeaders{}
	_result = &ListTriggersResponse{}
	_body, _err := client.ListTriggersWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListVpcBindingsWithOptions(serviceName *string, headers *ListVpcBindingsHeaders, runtime *util.RuntimeOptions) (_result *ListVpcBindingsResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("ListVpcBindings"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/bindings"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListVpcBindingsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListVpcBindings(serviceName *string) (_result *ListVpcBindingsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ListVpcBindingsHeaders{}
	_result = &ListVpcBindingsResponse{}
	_body, _err := client.ListVpcBindingsWithOptions(serviceName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PublishServiceVersionWithOptions(serviceName *string, request *PublishServiceVersionRequest, headers *PublishServiceVersionHeaders, runtime *util.RuntimeOptions) (_result *PublishServiceVersionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("PublishServiceVersion"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/versions"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &PublishServiceVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PublishServiceVersion(serviceName *string, request *PublishServiceVersionRequest) (_result *PublishServiceVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &PublishServiceVersionHeaders{}
	_result = &PublishServiceVersionResponse{}
	_body, _err := client.PublishServiceVersionWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * StatefulAsyncInvocation specifies the configurations of the asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request PutFunctionAsyncInvokeConfigRequest
 * @param headers PutFunctionAsyncInvokeConfigHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return PutFunctionAsyncInvokeConfigResponse
 */
func (client *Client) PutFunctionAsyncInvokeConfigWithOptions(serviceName *string, functionName *string, request *PutFunctionAsyncInvokeConfigRequest, headers *PutFunctionAsyncInvokeConfigHeaders, runtime *util.RuntimeOptions) (_result *PutFunctionAsyncInvokeConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DestinationConfig)) {
		body["destinationConfig"] = request.DestinationConfig
	}

	if !tea.BoolValue(util.IsUnset(request.MaxAsyncEventAgeInSeconds)) {
		body["maxAsyncEventAgeInSeconds"] = request.MaxAsyncEventAgeInSeconds
	}

	if !tea.BoolValue(util.IsUnset(request.MaxAsyncRetryAttempts)) {
		body["maxAsyncRetryAttempts"] = request.MaxAsyncRetryAttempts
	}

	if !tea.BoolValue(util.IsUnset(request.StatefulInvocation)) {
		body["statefulInvocation"] = request.StatefulInvocation
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("PutFunctionAsyncInvokeConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/async-invoke-config"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &PutFunctionAsyncInvokeConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * StatefulAsyncInvocation specifies the configurations of the asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request PutFunctionAsyncInvokeConfigRequest
 * @return PutFunctionAsyncInvokeConfigResponse
 */
func (client *Client) PutFunctionAsyncInvokeConfig(serviceName *string, functionName *string, request *PutFunctionAsyncInvokeConfigRequest) (_result *PutFunctionAsyncInvokeConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &PutFunctionAsyncInvokeConfigHeaders{}
	_result = &PutFunctionAsyncInvokeConfigResponse{}
	_body, _err := client.PutFunctionAsyncInvokeConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PutFunctionOnDemandConfigWithOptions(serviceName *string, functionName *string, request *PutFunctionOnDemandConfigRequest, headers *PutFunctionOnDemandConfigHeaders, runtime *util.RuntimeOptions) (_result *PutFunctionOnDemandConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.MaximumInstanceCount)) {
		body["maximumInstanceCount"] = request.MaximumInstanceCount
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("PutFunctionOnDemandConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/on-demand-config"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &PutFunctionOnDemandConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PutFunctionOnDemandConfig(serviceName *string, functionName *string, request *PutFunctionOnDemandConfigRequest) (_result *PutFunctionOnDemandConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &PutFunctionOnDemandConfigHeaders{}
	_result = &PutFunctionOnDemandConfigResponse{}
	_body, _err := client.PutFunctionOnDemandConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PutLayerACLWithOptions(layerName *string, request *PutLayerACLRequest, headers *PutLayerACLHeaders, runtime *util.RuntimeOptions) (_result *PutLayerACLResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Public)) {
		query["public"] = request.Public
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("PutLayerACL"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/layers/" + tea.StringValue(openapiutil.GetEncodeParam(layerName)) + "/acl"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &PutLayerACLResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PutLayerACL(layerName *string, request *PutLayerACLRequest) (_result *PutLayerACLResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &PutLayerACLHeaders{}
	_result = &PutLayerACLResponse{}
	_body, _err := client.PutLayerACLWithOptions(layerName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PutProvisionConfigWithOptions(serviceName *string, functionName *string, request *PutProvisionConfigRequest, headers *PutProvisionConfigHeaders, runtime *util.RuntimeOptions) (_result *PutProvisionConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AlwaysAllocateCPU)) {
		body["alwaysAllocateCPU"] = request.AlwaysAllocateCPU
	}

	if !tea.BoolValue(util.IsUnset(request.ScheduledActions)) {
		body["scheduledActions"] = request.ScheduledActions
	}

	if !tea.BoolValue(util.IsUnset(request.Target)) {
		body["target"] = request.Target
	}

	if !tea.BoolValue(util.IsUnset(request.TargetTrackingPolicies)) {
		body["targetTrackingPolicies"] = request.TargetTrackingPolicies
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("PutProvisionConfig"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/provision-config"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &PutProvisionConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PutProvisionConfig(serviceName *string, functionName *string, request *PutProvisionConfigRequest) (_result *PutProvisionConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &PutProvisionConfigHeaders{}
	_result = &PutProvisionConfigResponse{}
	_body, _err := client.PutProvisionConfigWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RegisterEventSourceWithOptions(serviceName *string, functionName *string, request *RegisterEventSourceRequest, headers *RegisterEventSourceHeaders, runtime *util.RuntimeOptions) (_result *RegisterEventSourceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.SourceArn)) {
		body["sourceArn"] = request.SourceArn
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("RegisterEventSource"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/event-sources"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &RegisterEventSourceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RegisterEventSource(serviceName *string, functionName *string, request *RegisterEventSourceRequest) (_result *RegisterEventSourceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &RegisterEventSourceHeaders{}
	_result = &RegisterEventSourceResponse{}
	_body, _err := client.RegisterEventSourceWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ReleaseGPUInstanceWithOptions(instanceId *string, headers *ReleaseGPUInstanceHeaders, runtime *util.RuntimeOptions) (_result *ReleaseGPUInstanceResponse, _err error) {
	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
	}
	params := &openapi.Params{
		Action:      tea.String("ReleaseGPUInstance"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/gpuInstances/" + tea.StringValue(openapiutil.GetEncodeParam(instanceId))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &ReleaseGPUInstanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ReleaseGPUInstance(instanceId *string) (_result *ReleaseGPUInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &ReleaseGPUInstanceHeaders{}
	_result = &ReleaseGPUInstanceResponse{}
	_body, _err := client.ReleaseGPUInstanceWithOptions(instanceId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request StopStatefulAsyncInvocationRequest
 * @param headers StopStatefulAsyncInvocationHeaders
 * @param runtime runtime options for this request RuntimeOptions
 * @return StopStatefulAsyncInvocationResponse
 */
func (client *Client) StopStatefulAsyncInvocationWithOptions(serviceName *string, functionName *string, invocationId *string, request *StopStatefulAsyncInvocationRequest, headers *StopStatefulAsyncInvocationHeaders, runtime *util.RuntimeOptions) (_result *StopStatefulAsyncInvocationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		query["qualifier"] = request.Qualifier
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("StopStatefulAsyncInvocation"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/stateful-async-invocations/" + tea.StringValue(openapiutil.GetEncodeParam(invocationId))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &StopStatefulAsyncInvocationResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * StatefulAsyncInvocation: asynchronous task. Asynchronous tasks allow you to manage the states on the basis of common asynchronous invocations, which is more suitable for task scenarios.
 *
 * @param request StopStatefulAsyncInvocationRequest
 * @return StopStatefulAsyncInvocationResponse
 */
func (client *Client) StopStatefulAsyncInvocation(serviceName *string, functionName *string, invocationId *string, request *StopStatefulAsyncInvocationRequest) (_result *StopStatefulAsyncInvocationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &StopStatefulAsyncInvocationHeaders{}
	_result = &StopStatefulAsyncInvocationResponse{}
	_body, _err := client.StopStatefulAsyncInvocationWithOptions(serviceName, functionName, invocationId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) TagResourceWithOptions(request *TagResourceRequest, headers *TagResourceHeaders, runtime *util.RuntimeOptions) (_result *TagResourceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ResourceArn)) {
		body["resourceArn"] = request.ResourceArn
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("TagResource"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/tag"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &TagResourceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) TagResource(request *TagResourceRequest) (_result *TagResourceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &TagResourceHeaders{}
	_result = &TagResourceResponse{}
	_body, _err := client.TagResourceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UntagResourceWithOptions(request *UntagResourceRequest, headers *UntagResourceHeaders, runtime *util.RuntimeOptions) (_result *UntagResourceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.All)) {
		body["all"] = request.All
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceArn)) {
		body["resourceArn"] = request.ResourceArn
	}

	if !tea.BoolValue(util.IsUnset(request.TagKeys)) {
		body["tagKeys"] = request.TagKeys
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UntagResource"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/tag"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UntagResourceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UntagResource(request *UntagResourceRequest) (_result *UntagResourceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UntagResourceHeaders{}
	_result = &UntagResourceResponse{}
	_body, _err := client.UntagResourceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateAliasWithOptions(serviceName *string, aliasName *string, request *UpdateAliasRequest, headers *UpdateAliasHeaders, runtime *util.RuntimeOptions) (_result *UpdateAliasResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AdditionalVersionWeight)) {
		body["additionalVersionWeight"] = request.AdditionalVersionWeight
	}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.ResolvePolicy)) {
		body["resolvePolicy"] = request.ResolvePolicy
	}

	if !tea.BoolValue(util.IsUnset(request.RoutePolicy)) {
		body["routePolicy"] = request.RoutePolicy
	}

	if !tea.BoolValue(util.IsUnset(request.VersionId)) {
		body["versionId"] = request.VersionId
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateAlias"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/aliases/" + tea.StringValue(openapiutil.GetEncodeParam(aliasName))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &UpdateAliasResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateAlias(serviceName *string, aliasName *string, request *UpdateAliasRequest) (_result *UpdateAliasResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UpdateAliasHeaders{}
	_result = &UpdateAliasResponse{}
	_body, _err := client.UpdateAliasWithOptions(serviceName, aliasName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateCustomDomainWithOptions(domainName *string, request *UpdateCustomDomainRequest, headers *UpdateCustomDomainHeaders, runtime *util.RuntimeOptions) (_result *UpdateCustomDomainResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CertConfig)) {
		body["certConfig"] = request.CertConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Protocol)) {
		body["protocol"] = request.Protocol
	}

	if !tea.BoolValue(util.IsUnset(request.RouteConfig)) {
		body["routeConfig"] = request.RouteConfig
	}

	if !tea.BoolValue(util.IsUnset(request.TlsConfig)) {
		body["tlsConfig"] = request.TlsConfig
	}

	if !tea.BoolValue(util.IsUnset(request.WafConfig)) {
		body["wafConfig"] = request.WafConfig
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateCustomDomain"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/custom-domains/" + tea.StringValue(openapiutil.GetEncodeParam(domainName))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &UpdateCustomDomainResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateCustomDomain(domainName *string, request *UpdateCustomDomainRequest) (_result *UpdateCustomDomainResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UpdateCustomDomainHeaders{}
	_result = &UpdateCustomDomainResponse{}
	_body, _err := client.UpdateCustomDomainWithOptions(domainName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateFunctionWithOptions(serviceName *string, functionName *string, request *UpdateFunctionRequest, headers *UpdateFunctionHeaders, runtime *util.RuntimeOptions) (_result *UpdateFunctionResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.InstanceConcurrency)) {
		body["InstanceConcurrency"] = request.InstanceConcurrency
	}

	if !tea.BoolValue(util.IsUnset(request.CaPort)) {
		body["caPort"] = request.CaPort
	}

	if !tea.BoolValue(util.IsUnset(request.Code)) {
		body["code"] = request.Code
	}

	if !tea.BoolValue(util.IsUnset(request.Cpu)) {
		body["cpu"] = request.Cpu
	}

	if !tea.BoolValue(util.IsUnset(request.CustomContainerConfig)) {
		body["customContainerConfig"] = request.CustomContainerConfig
	}

	if !tea.BoolValue(util.IsUnset(request.CustomDNS)) {
		body["customDNS"] = request.CustomDNS
	}

	if !tea.BoolValue(util.IsUnset(request.CustomHealthCheckConfig)) {
		body["customHealthCheckConfig"] = request.CustomHealthCheckConfig
	}

	if !tea.BoolValue(util.IsUnset(request.CustomRuntimeConfig)) {
		body["customRuntimeConfig"] = request.CustomRuntimeConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.DiskSize)) {
		body["diskSize"] = request.DiskSize
	}

	if !tea.BoolValue(util.IsUnset(request.EnvironmentVariables)) {
		body["environmentVariables"] = request.EnvironmentVariables
	}

	if !tea.BoolValue(util.IsUnset(request.GpuMemorySize)) {
		body["gpuMemorySize"] = request.GpuMemorySize
	}

	if !tea.BoolValue(util.IsUnset(request.Handler)) {
		body["handler"] = request.Handler
	}

	if !tea.BoolValue(util.IsUnset(request.InitializationTimeout)) {
		body["initializationTimeout"] = request.InitializationTimeout
	}

	if !tea.BoolValue(util.IsUnset(request.Initializer)) {
		body["initializer"] = request.Initializer
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceLifecycleConfig)) {
		body["instanceLifecycleConfig"] = request.InstanceLifecycleConfig
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceSoftConcurrency)) {
		body["instanceSoftConcurrency"] = request.InstanceSoftConcurrency
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceType)) {
		body["instanceType"] = request.InstanceType
	}

	if !tea.BoolValue(util.IsUnset(request.Layers)) {
		body["layers"] = request.Layers
	}

	if !tea.BoolValue(util.IsUnset(request.MemorySize)) {
		body["memorySize"] = request.MemorySize
	}

	if !tea.BoolValue(util.IsUnset(request.Runtime)) {
		body["runtime"] = request.Runtime
	}

	if !tea.BoolValue(util.IsUnset(request.Timeout)) {
		body["timeout"] = request.Timeout
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcCodeChecksum)) {
		realHeaders["X-Fc-Code-Checksum"] = util.ToJSONString(headers.XFcCodeChecksum)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateFunction"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &UpdateFunctionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateFunction(serviceName *string, functionName *string, request *UpdateFunctionRequest) (_result *UpdateFunctionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UpdateFunctionHeaders{}
	_result = &UpdateFunctionResponse{}
	_body, _err := client.UpdateFunctionWithOptions(serviceName, functionName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateServiceWithOptions(serviceName *string, request *UpdateServiceRequest, headers *UpdateServiceHeaders, runtime *util.RuntimeOptions) (_result *UpdateServiceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.InternetAccess)) {
		body["internetAccess"] = request.InternetAccess
	}

	if !tea.BoolValue(util.IsUnset(request.LogConfig)) {
		body["logConfig"] = request.LogConfig
	}

	if !tea.BoolValue(util.IsUnset(request.NasConfig)) {
		body["nasConfig"] = request.NasConfig
	}

	if !tea.BoolValue(util.IsUnset(request.OssMountConfig)) {
		body["ossMountConfig"] = request.OssMountConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Role)) {
		body["role"] = request.Role
	}

	if !tea.BoolValue(util.IsUnset(request.TracingConfig)) {
		body["tracingConfig"] = request.TracingConfig
	}

	if !tea.BoolValue(util.IsUnset(request.VpcConfig)) {
		body["vpcConfig"] = request.VpcConfig
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateService"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &UpdateServiceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateService(serviceName *string, request *UpdateServiceRequest) (_result *UpdateServiceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UpdateServiceHeaders{}
	_result = &UpdateServiceResponse{}
	_body, _err := client.UpdateServiceWithOptions(serviceName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateTriggerWithOptions(serviceName *string, functionName *string, triggerName *string, request *UpdateTriggerRequest, headers *UpdateTriggerHeaders, runtime *util.RuntimeOptions) (_result *UpdateTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.InvocationRole)) {
		body["invocationRole"] = request.InvocationRole
	}

	if !tea.BoolValue(util.IsUnset(request.Qualifier)) {
		body["qualifier"] = request.Qualifier
	}

	if !tea.BoolValue(util.IsUnset(request.TriggerConfig)) {
		body["triggerConfig"] = request.TriggerConfig
	}

	realHeaders := make(map[string]*string)
	if !tea.BoolValue(util.IsUnset(headers.CommonHeaders)) {
		realHeaders = headers.CommonHeaders
	}

	if !tea.BoolValue(util.IsUnset(headers.IfMatch)) {
		realHeaders["If-Match"] = util.ToJSONString(headers.IfMatch)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcAccountId)) {
		realHeaders["X-Fc-Account-Id"] = util.ToJSONString(headers.XFcAccountId)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcDate)) {
		realHeaders["X-Fc-Date"] = util.ToJSONString(headers.XFcDate)
	}

	if !tea.BoolValue(util.IsUnset(headers.XFcTraceId)) {
		realHeaders["X-Fc-Trace-Id"] = util.ToJSONString(headers.XFcTraceId)
	}

	req := &openapi.OpenApiRequest{
		Headers: realHeaders,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateTrigger"),
		Version:     tea.String("2021-04-06"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/2021-04-06/services/" + tea.StringValue(openapiutil.GetEncodeParam(serviceName)) + "/functions/" + tea.StringValue(openapiutil.GetEncodeParam(functionName)) + "/triggers/" + tea.StringValue(openapiutil.GetEncodeParam(triggerName))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &UpdateTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateTrigger(serviceName *string, functionName *string, triggerName *string, request *UpdateTriggerRequest) (_result *UpdateTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := &UpdateTriggerHeaders{}
	_result = &UpdateTriggerResponse{}
	_body, _err := client.UpdateTriggerWithOptions(serviceName, functionName, triggerName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) InvokeHTTPTrigger(url *string, method *string, body []byte, headers *http.Header) (_result *http.Response, _err error) {
	cred := client.Credential
	utilClient, _err := fcutil.NewClient(cred)
	if _err != nil {
		return _result, _err
	}

	_body, _err := utilClient.InvokeHTTPTrigger(url, method, body, headers)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) InvokeAnonymousHTTPTrigger(url *string, method *string, body []byte, headers *http.Header) (_result *http.Response, _err error) {
	cred := client.Credential
	utilClient, _err := fcutil.NewClient(cred)
	if _err != nil {
		return _result, _err
	}

	_body, _err := utilClient.InvokeAnonymousHTTPTrigger(url, method, body, headers)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SendHTTPRequestWithAuthorization(req *http.Request) (_result *http.Response, _err error) {
	cred := client.Credential
	utilClient, _err := fcutil.NewClient(cred)
	if _err != nil {
		return _result, _err
	}

	_body, _err := utilClient.SendHTTPRequestWithAuthorization(req)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SendHTTPRequest(req *http.Request) (_result *http.Response, _err error) {
	cred := client.Credential
	utilClient, _err := fcutil.NewClient(cred)
	if _err != nil {
		return _result, _err
	}

	_body, _err := utilClient.SendHTTPRequest(req)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SignRequest(req *http.Request) (_result *http.Request, _err error) {
	cred := client.Credential
	utilClient, _err := fcutil.NewClient(cred)
	if _err != nil {
		return _result, _err
	}

	_body, _err := utilClient.SignRequest(req)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SignRequestWithContentMD5(req *http.Request, contentMD5 *string) (_result *http.Request, _err error) {
	cred := client.Credential
	utilClient, _err := fcutil.NewClient(cred)
	if _err != nil {
		return _result, _err
	}

	_body, _err := utilClient.SignRequestWithContentMD5(req, contentMD5)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) BuildHTTPRequest(url *string, method *string, body []byte, headers *http.Header) (_result *http.Request, _err error) {
	cred := client.Credential
	utilClient, _err := fcutil.NewClient(cred)
	if _err != nil {
		return _result, _err
	}

	_body, _err := utilClient.BuildHTTPRequest(url, method, body, headers)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}
