// This file is auto-generated, don't edit it. Thanks.
/**
 *
 */
package client

import (
	"io"

	"github.com/alibabacloud-go/tea/tea"
)

type AccessControlList struct {
	Grant *string `json:"Grant,omitempty" xml:"Grant,omitempty"`
}

func (s AccessControlList) String() string {
	return tea.Prettify(s)
}

func (s AccessControlList) GoString() string {
	return s.String()
}

func (s *AccessControlList) SetGrant(v string) *AccessControlList {
	s.Grant = &v
	return s
}

type AccessControlPolicy struct {
	AccessControlList *AccessControlList `json:"AccessControlList,omitempty" xml:"AccessControlList,omitempty"`
	Owner             *Owner             `json:"Owner,omitempty" xml:"Owner,omitempty"`
}

func (s AccessControlPolicy) String() string {
	return tea.Prettify(s)
}

func (s AccessControlPolicy) GoString() string {
	return s.String()
}

func (s *AccessControlPolicy) SetAccessControlList(v *AccessControlList) *AccessControlPolicy {
	s.AccessControlList = v
	return s
}

func (s *AccessControlPolicy) SetOwner(v *Owner) *AccessControlPolicy {
	s.Owner = v
	return s
}

type AccessMonitorConfiguration struct {
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s AccessMonitorConfiguration) String() string {
	return tea.Prettify(s)
}

func (s AccessMonitorConfiguration) GoString() string {
	return s.String()
}

func (s *AccessMonitorConfiguration) SetStatus(v string) *AccessMonitorConfiguration {
	s.Status = &v
	return s
}

type AccessPoint struct {
	AccessPointName  *string                      `json:"AccessPointName,omitempty" xml:"AccessPointName,omitempty"`
	Alias            *string                      `json:"Alias,omitempty" xml:"Alias,omitempty"`
	Bucket           *string                      `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	NetworkOrigin    *string                      `json:"NetworkOrigin,omitempty" xml:"NetworkOrigin,omitempty"`
	Status           *string                      `json:"Status,omitempty" xml:"Status,omitempty"`
	VpcConfiguration *AccessPointVpcConfiguration `json:"VpcConfiguration,omitempty" xml:"VpcConfiguration,omitempty"`
}

func (s AccessPoint) String() string {
	return tea.Prettify(s)
}

func (s AccessPoint) GoString() string {
	return s.String()
}

func (s *AccessPoint) SetAccessPointName(v string) *AccessPoint {
	s.AccessPointName = &v
	return s
}

func (s *AccessPoint) SetAlias(v string) *AccessPoint {
	s.Alias = &v
	return s
}

func (s *AccessPoint) SetBucket(v string) *AccessPoint {
	s.Bucket = &v
	return s
}

func (s *AccessPoint) SetNetworkOrigin(v string) *AccessPoint {
	s.NetworkOrigin = &v
	return s
}

func (s *AccessPoint) SetStatus(v string) *AccessPoint {
	s.Status = &v
	return s
}

func (s *AccessPoint) SetVpcConfiguration(v *AccessPointVpcConfiguration) *AccessPoint {
	s.VpcConfiguration = v
	return s
}

type AccessPointVpcConfiguration struct {
	VpcId *string `json:"VpcId,omitempty" xml:"VpcId,omitempty"`
}

func (s AccessPointVpcConfiguration) String() string {
	return tea.Prettify(s)
}

func (s AccessPointVpcConfiguration) GoString() string {
	return s.String()
}

func (s *AccessPointVpcConfiguration) SetVpcId(v string) *AccessPointVpcConfiguration {
	s.VpcId = &v
	return s
}

type ApplyServerSideEncryptionByDefault struct {
	KMSDataEncryption *string `json:"KMSDataEncryption,omitempty" xml:"KMSDataEncryption,omitempty"`
	KMSMasterKeyID    *string `json:"KMSMasterKeyID,omitempty" xml:"KMSMasterKeyID,omitempty"`
	SSEAlgorithm      *string `json:"SSEAlgorithm,omitempty" xml:"SSEAlgorithm,omitempty"`
}

func (s ApplyServerSideEncryptionByDefault) String() string {
	return tea.Prettify(s)
}

func (s ApplyServerSideEncryptionByDefault) GoString() string {
	return s.String()
}

func (s *ApplyServerSideEncryptionByDefault) SetKMSDataEncryption(v string) *ApplyServerSideEncryptionByDefault {
	s.KMSDataEncryption = &v
	return s
}

func (s *ApplyServerSideEncryptionByDefault) SetKMSMasterKeyID(v string) *ApplyServerSideEncryptionByDefault {
	s.KMSMasterKeyID = &v
	return s
}

func (s *ApplyServerSideEncryptionByDefault) SetSSEAlgorithm(v string) *ApplyServerSideEncryptionByDefault {
	s.SSEAlgorithm = &v
	return s
}

type ArchiveDirectReadConfiguration struct {
	Enabled *bool `json:"Enabled,omitempty" xml:"Enabled,omitempty"`
}

func (s ArchiveDirectReadConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ArchiveDirectReadConfiguration) GoString() string {
	return s.String()
}

func (s *ArchiveDirectReadConfiguration) SetEnabled(v bool) *ArchiveDirectReadConfiguration {
	s.Enabled = &v
	return s
}

type Bucket struct {
	CreationDate     *string `json:"CreationDate,omitempty" xml:"CreationDate,omitempty"`
	ExtranetEndpoint *string `json:"ExtranetEndpoint,omitempty" xml:"ExtranetEndpoint,omitempty"`
	IntranetEndpoint *string `json:"IntranetEndpoint,omitempty" xml:"IntranetEndpoint,omitempty"`
	Location         *string `json:"Location,omitempty" xml:"Location,omitempty"`
	Name             *string `json:"Name,omitempty" xml:"Name,omitempty"`
	Region           *string `json:"Region,omitempty" xml:"Region,omitempty"`
	StorageClass     *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
}

func (s Bucket) String() string {
	return tea.Prettify(s)
}

func (s Bucket) GoString() string {
	return s.String()
}

func (s *Bucket) SetCreationDate(v string) *Bucket {
	s.CreationDate = &v
	return s
}

func (s *Bucket) SetExtranetEndpoint(v string) *Bucket {
	s.ExtranetEndpoint = &v
	return s
}

func (s *Bucket) SetIntranetEndpoint(v string) *Bucket {
	s.IntranetEndpoint = &v
	return s
}

func (s *Bucket) SetLocation(v string) *Bucket {
	s.Location = &v
	return s
}

func (s *Bucket) SetName(v string) *Bucket {
	s.Name = &v
	return s
}

func (s *Bucket) SetRegion(v string) *Bucket {
	s.Region = &v
	return s
}

func (s *Bucket) SetStorageClass(v string) *Bucket {
	s.StorageClass = &v
	return s
}

type BucketAntiDDOSConfiguration struct {
	Cnames *BucketAntiDDOSConfigurationCnames `json:"Cnames,omitempty" xml:"Cnames,omitempty" type:"Struct"`
}

func (s BucketAntiDDOSConfiguration) String() string {
	return tea.Prettify(s)
}

func (s BucketAntiDDOSConfiguration) GoString() string {
	return s.String()
}

func (s *BucketAntiDDOSConfiguration) SetCnames(v *BucketAntiDDOSConfigurationCnames) *BucketAntiDDOSConfiguration {
	s.Cnames = v
	return s
}

type BucketAntiDDOSConfigurationCnames struct {
	Domain []*string `json:"Domain,omitempty" xml:"Domain,omitempty" type:"Repeated"`
}

func (s BucketAntiDDOSConfigurationCnames) String() string {
	return tea.Prettify(s)
}

func (s BucketAntiDDOSConfigurationCnames) GoString() string {
	return s.String()
}

func (s *BucketAntiDDOSConfigurationCnames) SetDomain(v []*string) *BucketAntiDDOSConfigurationCnames {
	s.Domain = v
	return s
}

type BucketAntiDDOSInfo struct {
	ActiveTime *int64                    `json:"ActiveTime,omitempty" xml:"ActiveTime,omitempty"`
	Bucket     *string                   `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	Cnames     *BucketAntiDDOSInfoCnames `json:"Cnames,omitempty" xml:"Cnames,omitempty" type:"Struct"`
	Ctime      *int64                    `json:"Ctime,omitempty" xml:"Ctime,omitempty"`
	InstanceId *string                   `json:"InstanceId,omitempty" xml:"InstanceId,omitempty"`
	Mtime      *int64                    `json:"Mtime,omitempty" xml:"Mtime,omitempty"`
	Owner      *string                   `json:"Owner,omitempty" xml:"Owner,omitempty"`
	Status     *string                   `json:"Status,omitempty" xml:"Status,omitempty"`
	Type       *string                   `json:"Type,omitempty" xml:"Type,omitempty"`
}

func (s BucketAntiDDOSInfo) String() string {
	return tea.Prettify(s)
}

func (s BucketAntiDDOSInfo) GoString() string {
	return s.String()
}

func (s *BucketAntiDDOSInfo) SetActiveTime(v int64) *BucketAntiDDOSInfo {
	s.ActiveTime = &v
	return s
}

func (s *BucketAntiDDOSInfo) SetBucket(v string) *BucketAntiDDOSInfo {
	s.Bucket = &v
	return s
}

func (s *BucketAntiDDOSInfo) SetCnames(v *BucketAntiDDOSInfoCnames) *BucketAntiDDOSInfo {
	s.Cnames = v
	return s
}

func (s *BucketAntiDDOSInfo) SetCtime(v int64) *BucketAntiDDOSInfo {
	s.Ctime = &v
	return s
}

func (s *BucketAntiDDOSInfo) SetInstanceId(v string) *BucketAntiDDOSInfo {
	s.InstanceId = &v
	return s
}

func (s *BucketAntiDDOSInfo) SetMtime(v int64) *BucketAntiDDOSInfo {
	s.Mtime = &v
	return s
}

func (s *BucketAntiDDOSInfo) SetOwner(v string) *BucketAntiDDOSInfo {
	s.Owner = &v
	return s
}

func (s *BucketAntiDDOSInfo) SetStatus(v string) *BucketAntiDDOSInfo {
	s.Status = &v
	return s
}

func (s *BucketAntiDDOSInfo) SetType(v string) *BucketAntiDDOSInfo {
	s.Type = &v
	return s
}

type BucketAntiDDOSInfoCnames struct {
	Domain []*string `json:"Domain,omitempty" xml:"Domain,omitempty" type:"Repeated"`
}

func (s BucketAntiDDOSInfoCnames) String() string {
	return tea.Prettify(s)
}

func (s BucketAntiDDOSInfoCnames) GoString() string {
	return s.String()
}

func (s *BucketAntiDDOSInfoCnames) SetDomain(v []*string) *BucketAntiDDOSInfoCnames {
	s.Domain = v
	return s
}

type BucketCnameConfiguration struct {
	Cname *BucketCnameConfigurationCname `json:"Cname,omitempty" xml:"Cname,omitempty" type:"Struct"`
}

func (s BucketCnameConfiguration) String() string {
	return tea.Prettify(s)
}

func (s BucketCnameConfiguration) GoString() string {
	return s.String()
}

func (s *BucketCnameConfiguration) SetCname(v *BucketCnameConfigurationCname) *BucketCnameConfiguration {
	s.Cname = v
	return s
}

type BucketCnameConfigurationCname struct {
	CertificateConfiguration *CertificateConfiguration `json:"CertificateConfiguration,omitempty" xml:"CertificateConfiguration,omitempty"`
	Domain                   *string                   `json:"Domain,omitempty" xml:"Domain,omitempty"`
}

func (s BucketCnameConfigurationCname) String() string {
	return tea.Prettify(s)
}

func (s BucketCnameConfigurationCname) GoString() string {
	return s.String()
}

func (s *BucketCnameConfigurationCname) SetCertificateConfiguration(v *CertificateConfiguration) *BucketCnameConfigurationCname {
	s.CertificateConfiguration = v
	return s
}

func (s *BucketCnameConfigurationCname) SetDomain(v string) *BucketCnameConfigurationCname {
	s.Domain = &v
	return s
}

type BucketDataRedundancyTransition struct {
	Bucket                 *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	CreateTime             *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	EndTime                *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	EstimatedRemainingTime *string `json:"EstimatedRemainingTime,omitempty" xml:"EstimatedRemainingTime,omitempty"`
	ProcessPercentage      *int32  `json:"ProcessPercentage,omitempty" xml:"ProcessPercentage,omitempty"`
	StartTime              *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
	Status                 *string `json:"Status,omitempty" xml:"Status,omitempty"`
	TaskId                 *string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`
}

func (s BucketDataRedundancyTransition) String() string {
	return tea.Prettify(s)
}

func (s BucketDataRedundancyTransition) GoString() string {
	return s.String()
}

func (s *BucketDataRedundancyTransition) SetBucket(v string) *BucketDataRedundancyTransition {
	s.Bucket = &v
	return s
}

func (s *BucketDataRedundancyTransition) SetCreateTime(v string) *BucketDataRedundancyTransition {
	s.CreateTime = &v
	return s
}

func (s *BucketDataRedundancyTransition) SetEndTime(v string) *BucketDataRedundancyTransition {
	s.EndTime = &v
	return s
}

func (s *BucketDataRedundancyTransition) SetEstimatedRemainingTime(v string) *BucketDataRedundancyTransition {
	s.EstimatedRemainingTime = &v
	return s
}

func (s *BucketDataRedundancyTransition) SetProcessPercentage(v int32) *BucketDataRedundancyTransition {
	s.ProcessPercentage = &v
	return s
}

func (s *BucketDataRedundancyTransition) SetStartTime(v string) *BucketDataRedundancyTransition {
	s.StartTime = &v
	return s
}

func (s *BucketDataRedundancyTransition) SetStatus(v string) *BucketDataRedundancyTransition {
	s.Status = &v
	return s
}

func (s *BucketDataRedundancyTransition) SetTaskId(v string) *BucketDataRedundancyTransition {
	s.TaskId = &v
	return s
}

type BucketInfo struct {
	Bucket *BucketInfoBucket `json:"Bucket,omitempty" xml:"Bucket,omitempty" type:"Struct"`
}

func (s BucketInfo) String() string {
	return tea.Prettify(s)
}

func (s BucketInfo) GoString() string {
	return s.String()
}

func (s *BucketInfo) SetBucket(v *BucketInfoBucket) *BucketInfo {
	s.Bucket = v
	return s
}

type BucketInfoBucket struct {
	AccessControlList        *AccessControlList                        `json:"AccessControlList,omitempty" xml:"AccessControlList,omitempty"`
	AccessMonitor            *string                                   `json:"AccessMonitor,omitempty" xml:"AccessMonitor,omitempty"`
	BucketPolicy             *BucketInfoBucketBucketPolicy             `json:"BucketPolicy,omitempty" xml:"BucketPolicy,omitempty" type:"Struct"`
	Comment                  *string                                   `json:"Comment,omitempty" xml:"Comment,omitempty"`
	CreationDate             *string                                   `json:"CreationDate,omitempty" xml:"CreationDate,omitempty"`
	CrossRegionReplication   *string                                   `json:"CrossRegionReplication,omitempty" xml:"CrossRegionReplication,omitempty"`
	DataRedundancyType       *string                                   `json:"DataRedundancyType,omitempty" xml:"DataRedundancyType,omitempty"`
	ExtranetEndpoint         *string                                   `json:"ExtranetEndpoint,omitempty" xml:"ExtranetEndpoint,omitempty"`
	IntranetEndpoint         *string                                   `json:"IntranetEndpoint,omitempty" xml:"IntranetEndpoint,omitempty"`
	Location                 *string                                   `json:"Location,omitempty" xml:"Location,omitempty"`
	Name                     *string                                   `json:"Name,omitempty" xml:"Name,omitempty"`
	Owner                    *Owner                                    `json:"Owner,omitempty" xml:"Owner,omitempty"`
	ResourceGroupId          *string                                   `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
	ServerSideEncryptionRule *BucketInfoBucketServerSideEncryptionRule `json:"ServerSideEncryptionRule,omitempty" xml:"ServerSideEncryptionRule,omitempty" type:"Struct"`
	StorageClass             *string                                   `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
	TransferAcceleration     *string                                   `json:"TransferAcceleration,omitempty" xml:"TransferAcceleration,omitempty"`
	Versioning               *string                                   `json:"Versioning,omitempty" xml:"Versioning,omitempty"`
}

func (s BucketInfoBucket) String() string {
	return tea.Prettify(s)
}

func (s BucketInfoBucket) GoString() string {
	return s.String()
}

func (s *BucketInfoBucket) SetAccessControlList(v *AccessControlList) *BucketInfoBucket {
	s.AccessControlList = v
	return s
}

func (s *BucketInfoBucket) SetAccessMonitor(v string) *BucketInfoBucket {
	s.AccessMonitor = &v
	return s
}

func (s *BucketInfoBucket) SetBucketPolicy(v *BucketInfoBucketBucketPolicy) *BucketInfoBucket {
	s.BucketPolicy = v
	return s
}

func (s *BucketInfoBucket) SetComment(v string) *BucketInfoBucket {
	s.Comment = &v
	return s
}

func (s *BucketInfoBucket) SetCreationDate(v string) *BucketInfoBucket {
	s.CreationDate = &v
	return s
}

func (s *BucketInfoBucket) SetCrossRegionReplication(v string) *BucketInfoBucket {
	s.CrossRegionReplication = &v
	return s
}

func (s *BucketInfoBucket) SetDataRedundancyType(v string) *BucketInfoBucket {
	s.DataRedundancyType = &v
	return s
}

func (s *BucketInfoBucket) SetExtranetEndpoint(v string) *BucketInfoBucket {
	s.ExtranetEndpoint = &v
	return s
}

func (s *BucketInfoBucket) SetIntranetEndpoint(v string) *BucketInfoBucket {
	s.IntranetEndpoint = &v
	return s
}

func (s *BucketInfoBucket) SetLocation(v string) *BucketInfoBucket {
	s.Location = &v
	return s
}

func (s *BucketInfoBucket) SetName(v string) *BucketInfoBucket {
	s.Name = &v
	return s
}

func (s *BucketInfoBucket) SetOwner(v *Owner) *BucketInfoBucket {
	s.Owner = v
	return s
}

func (s *BucketInfoBucket) SetResourceGroupId(v string) *BucketInfoBucket {
	s.ResourceGroupId = &v
	return s
}

func (s *BucketInfoBucket) SetServerSideEncryptionRule(v *BucketInfoBucketServerSideEncryptionRule) *BucketInfoBucket {
	s.ServerSideEncryptionRule = v
	return s
}

func (s *BucketInfoBucket) SetStorageClass(v string) *BucketInfoBucket {
	s.StorageClass = &v
	return s
}

func (s *BucketInfoBucket) SetTransferAcceleration(v string) *BucketInfoBucket {
	s.TransferAcceleration = &v
	return s
}

func (s *BucketInfoBucket) SetVersioning(v string) *BucketInfoBucket {
	s.Versioning = &v
	return s
}

type BucketInfoBucketBucketPolicy struct {
	LogBucket *string `json:"LogBucket,omitempty" xml:"LogBucket,omitempty"`
	LogPrefix *string `json:"LogPrefix,omitempty" xml:"LogPrefix,omitempty"`
}

func (s BucketInfoBucketBucketPolicy) String() string {
	return tea.Prettify(s)
}

func (s BucketInfoBucketBucketPolicy) GoString() string {
	return s.String()
}

func (s *BucketInfoBucketBucketPolicy) SetLogBucket(v string) *BucketInfoBucketBucketPolicy {
	s.LogBucket = &v
	return s
}

func (s *BucketInfoBucketBucketPolicy) SetLogPrefix(v string) *BucketInfoBucketBucketPolicy {
	s.LogPrefix = &v
	return s
}

type BucketInfoBucketServerSideEncryptionRule struct {
	KMSDataEncryption *string `json:"KMSDataEncryption,omitempty" xml:"KMSDataEncryption,omitempty"`
	KMSMasterKeyID    *string `json:"KMSMasterKeyID,omitempty" xml:"KMSMasterKeyID,omitempty"`
	SSEAlgorithm      *string `json:"SSEAlgorithm,omitempty" xml:"SSEAlgorithm,omitempty"`
}

func (s BucketInfoBucketServerSideEncryptionRule) String() string {
	return tea.Prettify(s)
}

func (s BucketInfoBucketServerSideEncryptionRule) GoString() string {
	return s.String()
}

func (s *BucketInfoBucketServerSideEncryptionRule) SetKMSDataEncryption(v string) *BucketInfoBucketServerSideEncryptionRule {
	s.KMSDataEncryption = &v
	return s
}

func (s *BucketInfoBucketServerSideEncryptionRule) SetKMSMasterKeyID(v string) *BucketInfoBucketServerSideEncryptionRule {
	s.KMSMasterKeyID = &v
	return s
}

func (s *BucketInfoBucketServerSideEncryptionRule) SetSSEAlgorithm(v string) *BucketInfoBucketServerSideEncryptionRule {
	s.SSEAlgorithm = &v
	return s
}

type BucketLoggingStatus struct {
	LoggingEnabled *LoggingEnabled `json:"LoggingEnabled,omitempty" xml:"LoggingEnabled,omitempty"`
}

func (s BucketLoggingStatus) String() string {
	return tea.Prettify(s)
}

func (s BucketLoggingStatus) GoString() string {
	return s.String()
}

func (s *BucketLoggingStatus) SetLoggingEnabled(v *LoggingEnabled) *BucketLoggingStatus {
	s.LoggingEnabled = v
	return s
}

type BucketResourceGroupConfiguration struct {
	ResourceGroupId *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
}

func (s BucketResourceGroupConfiguration) String() string {
	return tea.Prettify(s)
}

func (s BucketResourceGroupConfiguration) GoString() string {
	return s.String()
}

func (s *BucketResourceGroupConfiguration) SetResourceGroupId(v string) *BucketResourceGroupConfiguration {
	s.ResourceGroupId = &v
	return s
}

type BucketStat struct {
	ArchiveObjectCount          *int64 `json:"ArchiveObjectCount,omitempty" xml:"ArchiveObjectCount,omitempty"`
	ArchiveRealStorage          *int64 `json:"ArchiveRealStorage,omitempty" xml:"ArchiveRealStorage,omitempty"`
	ArchiveStorage              *int64 `json:"ArchiveStorage,omitempty" xml:"ArchiveStorage,omitempty"`
	ColdArchiveObjectCount      *int64 `json:"ColdArchiveObjectCount,omitempty" xml:"ColdArchiveObjectCount,omitempty"`
	ColdArchiveRealStorage      *int64 `json:"ColdArchiveRealStorage,omitempty" xml:"ColdArchiveRealStorage,omitempty"`
	ColdArchiveStorage          *int64 `json:"ColdArchiveStorage,omitempty" xml:"ColdArchiveStorage,omitempty"`
	DeepColdArchiveObjectCount  *int64 `json:"DeepColdArchiveObjectCount,omitempty" xml:"DeepColdArchiveObjectCount,omitempty"`
	DeepColdArchiveRealStorage  *int64 `json:"DeepColdArchiveRealStorage,omitempty" xml:"DeepColdArchiveRealStorage,omitempty"`
	DeepColdArchiveStorage      *int64 `json:"DeepColdArchiveStorage,omitempty" xml:"DeepColdArchiveStorage,omitempty"`
	DeleteMarkerCount           *int64 `json:"DeleteMarkerCount,omitempty" xml:"DeleteMarkerCount,omitempty"`
	InfrequentAccessObjectCount *int64 `json:"InfrequentAccessObjectCount,omitempty" xml:"InfrequentAccessObjectCount,omitempty"`
	InfrequentAccessRealStorage *int64 `json:"InfrequentAccessRealStorage,omitempty" xml:"InfrequentAccessRealStorage,omitempty"`
	InfrequentAccessStorage     *int64 `json:"InfrequentAccessStorage,omitempty" xml:"InfrequentAccessStorage,omitempty"`
	LastModifiedTime            *int64 `json:"LastModifiedTime,omitempty" xml:"LastModifiedTime,omitempty"`
	LiveChannelCount            *int64 `json:"LiveChannelCount,omitempty" xml:"LiveChannelCount,omitempty"`
	MultipartPartCount          *int64 `json:"MultipartPartCount,omitempty" xml:"MultipartPartCount,omitempty"`
	MultipartUploadCount        *int64 `json:"MultipartUploadCount,omitempty" xml:"MultipartUploadCount,omitempty"`
	ObjectCount                 *int64 `json:"ObjectCount,omitempty" xml:"ObjectCount,omitempty"`
	ReservedCapacityObjectCount *int64 `json:"ReservedCapacityObjectCount,omitempty" xml:"ReservedCapacityObjectCount,omitempty"`
	ReservedCapacityStorage     *int64 `json:"ReservedCapacityStorage,omitempty" xml:"ReservedCapacityStorage,omitempty"`
	StandardObjectCount         *int64 `json:"StandardObjectCount,omitempty" xml:"StandardObjectCount,omitempty"`
	StandardStorage             *int64 `json:"StandardStorage,omitempty" xml:"StandardStorage,omitempty"`
	Storage                     *int64 `json:"Storage,omitempty" xml:"Storage,omitempty"`
}

func (s BucketStat) String() string {
	return tea.Prettify(s)
}

func (s BucketStat) GoString() string {
	return s.String()
}

func (s *BucketStat) SetArchiveObjectCount(v int64) *BucketStat {
	s.ArchiveObjectCount = &v
	return s
}

func (s *BucketStat) SetArchiveRealStorage(v int64) *BucketStat {
	s.ArchiveRealStorage = &v
	return s
}

func (s *BucketStat) SetArchiveStorage(v int64) *BucketStat {
	s.ArchiveStorage = &v
	return s
}

func (s *BucketStat) SetColdArchiveObjectCount(v int64) *BucketStat {
	s.ColdArchiveObjectCount = &v
	return s
}

func (s *BucketStat) SetColdArchiveRealStorage(v int64) *BucketStat {
	s.ColdArchiveRealStorage = &v
	return s
}

func (s *BucketStat) SetColdArchiveStorage(v int64) *BucketStat {
	s.ColdArchiveStorage = &v
	return s
}

func (s *BucketStat) SetDeepColdArchiveObjectCount(v int64) *BucketStat {
	s.DeepColdArchiveObjectCount = &v
	return s
}

func (s *BucketStat) SetDeepColdArchiveRealStorage(v int64) *BucketStat {
	s.DeepColdArchiveRealStorage = &v
	return s
}

func (s *BucketStat) SetDeepColdArchiveStorage(v int64) *BucketStat {
	s.DeepColdArchiveStorage = &v
	return s
}

func (s *BucketStat) SetDeleteMarkerCount(v int64) *BucketStat {
	s.DeleteMarkerCount = &v
	return s
}

func (s *BucketStat) SetInfrequentAccessObjectCount(v int64) *BucketStat {
	s.InfrequentAccessObjectCount = &v
	return s
}

func (s *BucketStat) SetInfrequentAccessRealStorage(v int64) *BucketStat {
	s.InfrequentAccessRealStorage = &v
	return s
}

func (s *BucketStat) SetInfrequentAccessStorage(v int64) *BucketStat {
	s.InfrequentAccessStorage = &v
	return s
}

func (s *BucketStat) SetLastModifiedTime(v int64) *BucketStat {
	s.LastModifiedTime = &v
	return s
}

func (s *BucketStat) SetLiveChannelCount(v int64) *BucketStat {
	s.LiveChannelCount = &v
	return s
}

func (s *BucketStat) SetMultipartPartCount(v int64) *BucketStat {
	s.MultipartPartCount = &v
	return s
}

func (s *BucketStat) SetMultipartUploadCount(v int64) *BucketStat {
	s.MultipartUploadCount = &v
	return s
}

func (s *BucketStat) SetObjectCount(v int64) *BucketStat {
	s.ObjectCount = &v
	return s
}

func (s *BucketStat) SetReservedCapacityObjectCount(v int64) *BucketStat {
	s.ReservedCapacityObjectCount = &v
	return s
}

func (s *BucketStat) SetReservedCapacityStorage(v int64) *BucketStat {
	s.ReservedCapacityStorage = &v
	return s
}

func (s *BucketStat) SetStandardObjectCount(v int64) *BucketStat {
	s.StandardObjectCount = &v
	return s
}

func (s *BucketStat) SetStandardStorage(v int64) *BucketStat {
	s.StandardStorage = &v
	return s
}

func (s *BucketStat) SetStorage(v int64) *BucketStat {
	s.Storage = &v
	return s
}

type CORSConfiguration struct {
	CORSRule     []*CORSRule `json:"CORSRule,omitempty" xml:"CORSRule,omitempty" type:"Repeated"`
	ResponseVary *bool       `json:"ResponseVary,omitempty" xml:"ResponseVary,omitempty"`
}

func (s CORSConfiguration) String() string {
	return tea.Prettify(s)
}

func (s CORSConfiguration) GoString() string {
	return s.String()
}

func (s *CORSConfiguration) SetCORSRule(v []*CORSRule) *CORSConfiguration {
	s.CORSRule = v
	return s
}

func (s *CORSConfiguration) SetResponseVary(v bool) *CORSConfiguration {
	s.ResponseVary = &v
	return s
}

type CORSRule struct {
	AllowedHeader []*string `json:"AllowedHeader,omitempty" xml:"AllowedHeader,omitempty" type:"Repeated"`
	AllowedMethod []*string `json:"AllowedMethod,omitempty" xml:"AllowedMethod,omitempty" type:"Repeated"`
	AllowedOrigin []*string `json:"AllowedOrigin,omitempty" xml:"AllowedOrigin,omitempty" type:"Repeated"`
	ExposeHeader  []*string `json:"ExposeHeader,omitempty" xml:"ExposeHeader,omitempty" type:"Repeated"`
	MaxAgeSeconds *int64    `json:"MaxAgeSeconds,omitempty" xml:"MaxAgeSeconds,omitempty"`
}

func (s CORSRule) String() string {
	return tea.Prettify(s)
}

func (s CORSRule) GoString() string {
	return s.String()
}

func (s *CORSRule) SetAllowedHeader(v []*string) *CORSRule {
	s.AllowedHeader = v
	return s
}

func (s *CORSRule) SetAllowedMethod(v []*string) *CORSRule {
	s.AllowedMethod = v
	return s
}

func (s *CORSRule) SetAllowedOrigin(v []*string) *CORSRule {
	s.AllowedOrigin = v
	return s
}

func (s *CORSRule) SetExposeHeader(v []*string) *CORSRule {
	s.ExposeHeader = v
	return s
}

func (s *CORSRule) SetMaxAgeSeconds(v int64) *CORSRule {
	s.MaxAgeSeconds = &v
	return s
}

type CSVInput struct {
	AllowQuotedRecordDelimiter *bool   `json:"AllowQuotedRecordDelimiter,omitempty" xml:"AllowQuotedRecordDelimiter,omitempty"`
	CommentCharacter           *string `json:"CommentCharacter,omitempty" xml:"CommentCharacter,omitempty"`
	FieldDelimiter             *string `json:"FieldDelimiter,omitempty" xml:"FieldDelimiter,omitempty"`
	FileHeaderInfo             *string `json:"FileHeaderInfo,omitempty" xml:"FileHeaderInfo,omitempty"`
	QuoteCharacter             *string `json:"QuoteCharacter,omitempty" xml:"QuoteCharacter,omitempty"`
	Range                      *string `json:"Range,omitempty" xml:"Range,omitempty"`
	RecordDelimiter            *string `json:"RecordDelimiter,omitempty" xml:"RecordDelimiter,omitempty"`
}

func (s CSVInput) String() string {
	return tea.Prettify(s)
}

func (s CSVInput) GoString() string {
	return s.String()
}

func (s *CSVInput) SetAllowQuotedRecordDelimiter(v bool) *CSVInput {
	s.AllowQuotedRecordDelimiter = &v
	return s
}

func (s *CSVInput) SetCommentCharacter(v string) *CSVInput {
	s.CommentCharacter = &v
	return s
}

func (s *CSVInput) SetFieldDelimiter(v string) *CSVInput {
	s.FieldDelimiter = &v
	return s
}

func (s *CSVInput) SetFileHeaderInfo(v string) *CSVInput {
	s.FileHeaderInfo = &v
	return s
}

func (s *CSVInput) SetQuoteCharacter(v string) *CSVInput {
	s.QuoteCharacter = &v
	return s
}

func (s *CSVInput) SetRange(v string) *CSVInput {
	s.Range = &v
	return s
}

func (s *CSVInput) SetRecordDelimiter(v string) *CSVInput {
	s.RecordDelimiter = &v
	return s
}

type CSVOutput struct {
	FieldDelimiter  *string `json:"FieldDelimiter,omitempty" xml:"FieldDelimiter,omitempty"`
	RecordDelimiter *string `json:"RecordDelimiter,omitempty" xml:"RecordDelimiter,omitempty"`
}

func (s CSVOutput) String() string {
	return tea.Prettify(s)
}

func (s CSVOutput) GoString() string {
	return s.String()
}

func (s *CSVOutput) SetFieldDelimiter(v string) *CSVOutput {
	s.FieldDelimiter = &v
	return s
}

func (s *CSVOutput) SetRecordDelimiter(v string) *CSVOutput {
	s.RecordDelimiter = &v
	return s
}

type CallbackPolicy struct {
	PolicyItem []*CallbackPolicyPolicyItem `json:"PolicyItem,omitempty" xml:"PolicyItem,omitempty" type:"Repeated"`
}

func (s CallbackPolicy) String() string {
	return tea.Prettify(s)
}

func (s CallbackPolicy) GoString() string {
	return s.String()
}

func (s *CallbackPolicy) SetPolicyItem(v []*CallbackPolicyPolicyItem) *CallbackPolicy {
	s.PolicyItem = v
	return s
}

type CallbackPolicyPolicyItem struct {
	Callback    *string `json:"Callback,omitempty" xml:"Callback,omitempty"`
	CallbackVar *string `json:"CallbackVar,omitempty" xml:"CallbackVar,omitempty"`
	PolicyName  *string `json:"PolicyName,omitempty" xml:"PolicyName,omitempty"`
}

func (s CallbackPolicyPolicyItem) String() string {
	return tea.Prettify(s)
}

func (s CallbackPolicyPolicyItem) GoString() string {
	return s.String()
}

func (s *CallbackPolicyPolicyItem) SetCallback(v string) *CallbackPolicyPolicyItem {
	s.Callback = &v
	return s
}

func (s *CallbackPolicyPolicyItem) SetCallbackVar(v string) *CallbackPolicyPolicyItem {
	s.CallbackVar = &v
	return s
}

func (s *CallbackPolicyPolicyItem) SetPolicyName(v string) *CallbackPolicyPolicyItem {
	s.PolicyName = &v
	return s
}

type CertificateConfiguration struct {
	CertId            *string `json:"CertId,omitempty" xml:"CertId,omitempty"`
	Certificate       *string `json:"Certificate,omitempty" xml:"Certificate,omitempty"`
	DeleteCertificate *bool   `json:"DeleteCertificate,omitempty" xml:"DeleteCertificate,omitempty"`
	Force             *bool   `json:"Force,omitempty" xml:"Force,omitempty"`
	PreviousCertId    *string `json:"PreviousCertId,omitempty" xml:"PreviousCertId,omitempty"`
	PrivateKey        *string `json:"PrivateKey,omitempty" xml:"PrivateKey,omitempty"`
}

func (s CertificateConfiguration) String() string {
	return tea.Prettify(s)
}

func (s CertificateConfiguration) GoString() string {
	return s.String()
}

func (s *CertificateConfiguration) SetCertId(v string) *CertificateConfiguration {
	s.CertId = &v
	return s
}

func (s *CertificateConfiguration) SetCertificate(v string) *CertificateConfiguration {
	s.Certificate = &v
	return s
}

func (s *CertificateConfiguration) SetDeleteCertificate(v bool) *CertificateConfiguration {
	s.DeleteCertificate = &v
	return s
}

func (s *CertificateConfiguration) SetForce(v bool) *CertificateConfiguration {
	s.Force = &v
	return s
}

func (s *CertificateConfiguration) SetPreviousCertId(v string) *CertificateConfiguration {
	s.PreviousCertId = &v
	return s
}

func (s *CertificateConfiguration) SetPrivateKey(v string) *CertificateConfiguration {
	s.PrivateKey = &v
	return s
}

type CnameCertificate struct {
	CertId         *string `json:"CertId,omitempty" xml:"CertId,omitempty"`
	CreationDate   *string `json:"CreationDate,omitempty" xml:"CreationDate,omitempty"`
	Fingerprint    *string `json:"Fingerprint,omitempty" xml:"Fingerprint,omitempty"`
	Status         *string `json:"Status,omitempty" xml:"Status,omitempty"`
	Type           *string `json:"Type,omitempty" xml:"Type,omitempty"`
	ValidEndDate   *string `json:"ValidEndDate,omitempty" xml:"ValidEndDate,omitempty"`
	ValidStartDate *string `json:"ValidStartDate,omitempty" xml:"ValidStartDate,omitempty"`
}

func (s CnameCertificate) String() string {
	return tea.Prettify(s)
}

func (s CnameCertificate) GoString() string {
	return s.String()
}

func (s *CnameCertificate) SetCertId(v string) *CnameCertificate {
	s.CertId = &v
	return s
}

func (s *CnameCertificate) SetCreationDate(v string) *CnameCertificate {
	s.CreationDate = &v
	return s
}

func (s *CnameCertificate) SetFingerprint(v string) *CnameCertificate {
	s.Fingerprint = &v
	return s
}

func (s *CnameCertificate) SetStatus(v string) *CnameCertificate {
	s.Status = &v
	return s
}

func (s *CnameCertificate) SetType(v string) *CnameCertificate {
	s.Type = &v
	return s
}

func (s *CnameCertificate) SetValidEndDate(v string) *CnameCertificate {
	s.ValidEndDate = &v
	return s
}

func (s *CnameCertificate) SetValidStartDate(v string) *CnameCertificate {
	s.ValidStartDate = &v
	return s
}

type CnameInfo struct {
	Certificate  *CnameCertificate `json:"Certificate,omitempty" xml:"Certificate,omitempty"`
	Domain       *string           `json:"Domain,omitempty" xml:"Domain,omitempty"`
	LastModified *string           `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
	Status       *string           `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s CnameInfo) String() string {
	return tea.Prettify(s)
}

func (s CnameInfo) GoString() string {
	return s.String()
}

func (s *CnameInfo) SetCertificate(v *CnameCertificate) *CnameInfo {
	s.Certificate = v
	return s
}

func (s *CnameInfo) SetDomain(v string) *CnameInfo {
	s.Domain = &v
	return s
}

func (s *CnameInfo) SetLastModified(v string) *CnameInfo {
	s.LastModified = &v
	return s
}

func (s *CnameInfo) SetStatus(v string) *CnameInfo {
	s.Status = &v
	return s
}

type CnameSummary struct {
	Certificate  *CnameCertificate `json:"Certificate,omitempty" xml:"Certificate,omitempty"`
	Domain       *string           `json:"Domain,omitempty" xml:"Domain,omitempty"`
	LastModified *string           `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
	Status       *string           `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s CnameSummary) String() string {
	return tea.Prettify(s)
}

func (s CnameSummary) GoString() string {
	return s.String()
}

func (s *CnameSummary) SetCertificate(v *CnameCertificate) *CnameSummary {
	s.Certificate = v
	return s
}

func (s *CnameSummary) SetDomain(v string) *CnameSummary {
	s.Domain = &v
	return s
}

func (s *CnameSummary) SetLastModified(v string) *CnameSummary {
	s.LastModified = &v
	return s
}

func (s *CnameSummary) SetStatus(v string) *CnameSummary {
	s.Status = &v
	return s
}

type CnameToken struct {
	Bucket     *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	Cname      *string `json:"Cname,omitempty" xml:"Cname,omitempty"`
	ExpireTime *string `json:"ExpireTime,omitempty" xml:"ExpireTime,omitempty"`
	Token      *string `json:"Token,omitempty" xml:"Token,omitempty"`
}

func (s CnameToken) String() string {
	return tea.Prettify(s)
}

func (s CnameToken) GoString() string {
	return s.String()
}

func (s *CnameToken) SetBucket(v string) *CnameToken {
	s.Bucket = &v
	return s
}

func (s *CnameToken) SetCname(v string) *CnameToken {
	s.Cname = &v
	return s
}

func (s *CnameToken) SetExpireTime(v string) *CnameToken {
	s.ExpireTime = &v
	return s
}

func (s *CnameToken) SetToken(v string) *CnameToken {
	s.Token = &v
	return s
}

type CommonPrefix struct {
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
}

func (s CommonPrefix) String() string {
	return tea.Prettify(s)
}

func (s CommonPrefix) GoString() string {
	return s.String()
}

func (s *CommonPrefix) SetPrefix(v string) *CommonPrefix {
	s.Prefix = &v
	return s
}

type CompleteMultipartUpload struct {
	Part []*CompleteMultipartUploadPart `json:"Part,omitempty" xml:"Part,omitempty" type:"Repeated"`
}

func (s CompleteMultipartUpload) String() string {
	return tea.Prettify(s)
}

func (s CompleteMultipartUpload) GoString() string {
	return s.String()
}

func (s *CompleteMultipartUpload) SetPart(v []*CompleteMultipartUploadPart) *CompleteMultipartUpload {
	s.Part = v
	return s
}

type CompleteMultipartUploadPart struct {
	ETag       *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	PartNumber *int64  `json:"PartNumber,omitempty" xml:"PartNumber,omitempty"`
}

func (s CompleteMultipartUploadPart) String() string {
	return tea.Prettify(s)
}

func (s CompleteMultipartUploadPart) GoString() string {
	return s.String()
}

func (s *CompleteMultipartUploadPart) SetETag(v string) *CompleteMultipartUploadPart {
	s.ETag = &v
	return s
}

func (s *CompleteMultipartUploadPart) SetPartNumber(v int64) *CompleteMultipartUploadPart {
	s.PartNumber = &v
	return s
}

type CopyObjectResult struct {
	ETag         *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
}

func (s CopyObjectResult) String() string {
	return tea.Prettify(s)
}

func (s CopyObjectResult) GoString() string {
	return s.String()
}

func (s *CopyObjectResult) SetETag(v string) *CopyObjectResult {
	s.ETag = &v
	return s
}

func (s *CopyObjectResult) SetLastModified(v string) *CopyObjectResult {
	s.LastModified = &v
	return s
}

type CopyPartResult struct {
	ETag         *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
}

func (s CopyPartResult) String() string {
	return tea.Prettify(s)
}

func (s CopyPartResult) GoString() string {
	return s.String()
}

func (s *CopyPartResult) SetETag(v string) *CopyPartResult {
	s.ETag = &v
	return s
}

func (s *CopyPartResult) SetLastModified(v string) *CopyPartResult {
	s.LastModified = &v
	return s
}

type CreateAccessPointConfiguration struct {
	AccessPointName  *string                      `json:"AccessPointName,omitempty" xml:"AccessPointName,omitempty"`
	NetworkOrigin    *string                      `json:"NetworkOrigin,omitempty" xml:"NetworkOrigin,omitempty"`
	VpcConfiguration *AccessPointVpcConfiguration `json:"VpcConfiguration,omitempty" xml:"VpcConfiguration,omitempty"`
}

func (s CreateAccessPointConfiguration) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointConfiguration) GoString() string {
	return s.String()
}

func (s *CreateAccessPointConfiguration) SetAccessPointName(v string) *CreateAccessPointConfiguration {
	s.AccessPointName = &v
	return s
}

func (s *CreateAccessPointConfiguration) SetNetworkOrigin(v string) *CreateAccessPointConfiguration {
	s.NetworkOrigin = &v
	return s
}

func (s *CreateAccessPointConfiguration) SetVpcConfiguration(v *AccessPointVpcConfiguration) *CreateAccessPointConfiguration {
	s.VpcConfiguration = v
	return s
}

type CreateAccessPointResult struct {
	AccessPointArn *string `json:"AccessPointArn,omitempty" xml:"AccessPointArn,omitempty"`
	Alias          *string `json:"Alias,omitempty" xml:"Alias,omitempty"`
}

func (s CreateAccessPointResult) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointResult) GoString() string {
	return s.String()
}

func (s *CreateAccessPointResult) SetAccessPointArn(v string) *CreateAccessPointResult {
	s.AccessPointArn = &v
	return s
}

func (s *CreateAccessPointResult) SetAlias(v string) *CreateAccessPointResult {
	s.Alias = &v
	return s
}

type CreateBucketConfiguration struct {
	DataRedundancyType *string `json:"DataRedundancyType,omitempty" xml:"DataRedundancyType,omitempty"`
	StorageClass       *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
}

func (s CreateBucketConfiguration) String() string {
	return tea.Prettify(s)
}

func (s CreateBucketConfiguration) GoString() string {
	return s.String()
}

func (s *CreateBucketConfiguration) SetDataRedundancyType(v string) *CreateBucketConfiguration {
	s.DataRedundancyType = &v
	return s
}

func (s *CreateBucketConfiguration) SetStorageClass(v string) *CreateBucketConfiguration {
	s.StorageClass = &v
	return s
}

type Delete struct {
	Objects []*ObjectIdentifier `json:"Object,omitempty" xml:"Object,omitempty" type:"Repeated"`
	Quiet   *bool               `json:"Quiet,omitempty" xml:"Quiet,omitempty"`
}

func (s Delete) String() string {
	return tea.Prettify(s)
}

func (s Delete) GoString() string {
	return s.String()
}

func (s *Delete) SetObjects(v []*ObjectIdentifier) *Delete {
	s.Objects = v
	return s
}

func (s *Delete) SetQuiet(v bool) *Delete {
	s.Quiet = &v
	return s
}

type DeleteMarkerEntry struct {
	IsLatest     *bool   `json:"IsLatest,omitempty" xml:"IsLatest,omitempty"`
	Key          *string `json:"Key,omitempty" xml:"Key,omitempty"`
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
	Owner        *Owner  `json:"Owner,omitempty" xml:"Owner,omitempty"`
	VersionId    *string `json:"VersionId,omitempty" xml:"VersionId,omitempty"`
}

func (s DeleteMarkerEntry) String() string {
	return tea.Prettify(s)
}

func (s DeleteMarkerEntry) GoString() string {
	return s.String()
}

func (s *DeleteMarkerEntry) SetIsLatest(v bool) *DeleteMarkerEntry {
	s.IsLatest = &v
	return s
}

func (s *DeleteMarkerEntry) SetKey(v string) *DeleteMarkerEntry {
	s.Key = &v
	return s
}

func (s *DeleteMarkerEntry) SetLastModified(v string) *DeleteMarkerEntry {
	s.LastModified = &v
	return s
}

func (s *DeleteMarkerEntry) SetOwner(v *Owner) *DeleteMarkerEntry {
	s.Owner = v
	return s
}

func (s *DeleteMarkerEntry) SetVersionId(v string) *DeleteMarkerEntry {
	s.VersionId = &v
	return s
}

type DeletedObject struct {
	DeleteMarker          *bool   `json:"DeleteMarker,omitempty" xml:"DeleteMarker,omitempty"`
	DeleteMarkerVersionId *string `json:"DeleteMarkerVersionId,omitempty" xml:"DeleteMarkerVersionId,omitempty"`
	Key                   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	VersionId             *string `json:"VersionId,omitempty" xml:"VersionId,omitempty"`
}

func (s DeletedObject) String() string {
	return tea.Prettify(s)
}

func (s DeletedObject) GoString() string {
	return s.String()
}

func (s *DeletedObject) SetDeleteMarker(v bool) *DeletedObject {
	s.DeleteMarker = &v
	return s
}

func (s *DeletedObject) SetDeleteMarkerVersionId(v string) *DeletedObject {
	s.DeleteMarkerVersionId = &v
	return s
}

func (s *DeletedObject) SetKey(v string) *DeletedObject {
	s.Key = &v
	return s
}

func (s *DeletedObject) SetVersionId(v string) *DeletedObject {
	s.VersionId = &v
	return s
}

type Error struct {
	Code      *string `json:"Code,omitempty" xml:"Code,omitempty"`
	HostId    *string `json:"HostId,omitempty" xml:"HostId,omitempty"`
	Message   *string `json:"Message,omitempty" xml:"Message,omitempty"`
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s Error) String() string {
	return tea.Prettify(s)
}

func (s Error) GoString() string {
	return s.String()
}

func (s *Error) SetCode(v string) *Error {
	s.Code = &v
	return s
}

func (s *Error) SetHostId(v string) *Error {
	s.HostId = &v
	return s
}

func (s *Error) SetMessage(v string) *Error {
	s.Message = &v
	return s
}

func (s *Error) SetRequestId(v string) *Error {
	s.RequestId = &v
	return s
}

type ErrorDocument struct {
	HttpStatus *int64  `json:"HttpStatus,omitempty" xml:"HttpStatus,omitempty"`
	Key        *string `json:"Key,omitempty" xml:"Key,omitempty"`
}

func (s ErrorDocument) String() string {
	return tea.Prettify(s)
}

func (s ErrorDocument) GoString() string {
	return s.String()
}

func (s *ErrorDocument) SetHttpStatus(v int64) *ErrorDocument {
	s.HttpStatus = &v
	return s
}

func (s *ErrorDocument) SetKey(v string) *ErrorDocument {
	s.Key = &v
	return s
}

type ExtendWormConfiguration struct {
	RetentionPeriodInDays *int32 `json:"RetentionPeriodInDays,omitempty" xml:"RetentionPeriodInDays,omitempty"`
}

func (s ExtendWormConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ExtendWormConfiguration) GoString() string {
	return s.String()
}

func (s *ExtendWormConfiguration) SetRetentionPeriodInDays(v int32) *ExtendWormConfiguration {
	s.RetentionPeriodInDays = &v
	return s
}

type GetAccessPointResult struct {
	AccessPointArn                 *string                         `json:"AccessPointArn,omitempty" xml:"AccessPointArn,omitempty"`
	AccessPointName                *string                         `json:"AccessPointName,omitempty" xml:"AccessPointName,omitempty"`
	AccountId                      *string                         `json:"AccountId,omitempty" xml:"AccountId,omitempty"`
	Alias                          *string                         `json:"Alias,omitempty" xml:"Alias,omitempty"`
	Bucket                         *string                         `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	CreationDate                   *string                         `json:"CreationDate,omitempty" xml:"CreationDate,omitempty"`
	Endpoints                      *GetAccessPointResultEndpoints  `json:"Endpoints,omitempty" xml:"Endpoints,omitempty" type:"Struct"`
	NetworkOrigin                  *string                         `json:"NetworkOrigin,omitempty" xml:"NetworkOrigin,omitempty"`
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
	Status                         *string                         `json:"Status,omitempty" xml:"Status,omitempty"`
	VpcConfiguration               *AccessPointVpcConfiguration    `json:"VpcConfiguration,omitempty" xml:"VpcConfiguration,omitempty"`
}

func (s GetAccessPointResult) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointResult) GoString() string {
	return s.String()
}

func (s *GetAccessPointResult) SetAccessPointArn(v string) *GetAccessPointResult {
	s.AccessPointArn = &v
	return s
}

func (s *GetAccessPointResult) SetAccessPointName(v string) *GetAccessPointResult {
	s.AccessPointName = &v
	return s
}

func (s *GetAccessPointResult) SetAccountId(v string) *GetAccessPointResult {
	s.AccountId = &v
	return s
}

func (s *GetAccessPointResult) SetAlias(v string) *GetAccessPointResult {
	s.Alias = &v
	return s
}

func (s *GetAccessPointResult) SetBucket(v string) *GetAccessPointResult {
	s.Bucket = &v
	return s
}

func (s *GetAccessPointResult) SetCreationDate(v string) *GetAccessPointResult {
	s.CreationDate = &v
	return s
}

func (s *GetAccessPointResult) SetEndpoints(v *GetAccessPointResultEndpoints) *GetAccessPointResult {
	s.Endpoints = v
	return s
}

func (s *GetAccessPointResult) SetNetworkOrigin(v string) *GetAccessPointResult {
	s.NetworkOrigin = &v
	return s
}

func (s *GetAccessPointResult) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *GetAccessPointResult {
	s.PublicAccessBlockConfiguration = v
	return s
}

func (s *GetAccessPointResult) SetStatus(v string) *GetAccessPointResult {
	s.Status = &v
	return s
}

func (s *GetAccessPointResult) SetVpcConfiguration(v *AccessPointVpcConfiguration) *GetAccessPointResult {
	s.VpcConfiguration = v
	return s
}

type GetAccessPointResultEndpoints struct {
	InternalEndpoint *string `json:"InternalEndpoint,omitempty" xml:"InternalEndpoint,omitempty"`
	PublicEndpoint   *string `json:"PublicEndpoint,omitempty" xml:"PublicEndpoint,omitempty"`
}

func (s GetAccessPointResultEndpoints) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointResultEndpoints) GoString() string {
	return s.String()
}

func (s *GetAccessPointResultEndpoints) SetInternalEndpoint(v string) *GetAccessPointResultEndpoints {
	s.InternalEndpoint = &v
	return s
}

func (s *GetAccessPointResultEndpoints) SetPublicEndpoint(v string) *GetAccessPointResultEndpoints {
	s.PublicEndpoint = &v
	return s
}

type HttpsConfiguration struct {
	TLS *HttpsConfigurationTLS `json:"TLS,omitempty" xml:"TLS,omitempty" type:"Struct"`
}

func (s HttpsConfiguration) String() string {
	return tea.Prettify(s)
}

func (s HttpsConfiguration) GoString() string {
	return s.String()
}

func (s *HttpsConfiguration) SetTLS(v *HttpsConfigurationTLS) *HttpsConfiguration {
	s.TLS = v
	return s
}

type HttpsConfigurationTLS struct {
	Enable     *bool     `json:"Enable,omitempty" xml:"Enable,omitempty"`
	TLSVersion []*string `json:"TLSVersion,omitempty" xml:"TLSVersion,omitempty" type:"Repeated"`
}

func (s HttpsConfigurationTLS) String() string {
	return tea.Prettify(s)
}

func (s HttpsConfigurationTLS) GoString() string {
	return s.String()
}

func (s *HttpsConfigurationTLS) SetEnable(v bool) *HttpsConfigurationTLS {
	s.Enable = &v
	return s
}

func (s *HttpsConfigurationTLS) SetTLSVersion(v []*string) *HttpsConfigurationTLS {
	s.TLSVersion = v
	return s
}

type IndexDocument struct {
	Suffix        *string `json:"Suffix,omitempty" xml:"Suffix,omitempty"`
	SupportSubDir *bool   `json:"SupportSubDir,omitempty" xml:"SupportSubDir,omitempty"`
	Type          *int64  `json:"Type,omitempty" xml:"Type,omitempty"`
}

func (s IndexDocument) String() string {
	return tea.Prettify(s)
}

func (s IndexDocument) GoString() string {
	return s.String()
}

func (s *IndexDocument) SetSuffix(v string) *IndexDocument {
	s.Suffix = &v
	return s
}

func (s *IndexDocument) SetSupportSubDir(v bool) *IndexDocument {
	s.SupportSubDir = &v
	return s
}

func (s *IndexDocument) SetType(v int64) *IndexDocument {
	s.Type = &v
	return s
}

type InitiateWormConfiguration struct {
	RetentionPeriodInDays *int32 `json:"RetentionPeriodInDays,omitempty" xml:"RetentionPeriodInDays,omitempty"`
}

func (s InitiateWormConfiguration) String() string {
	return tea.Prettify(s)
}

func (s InitiateWormConfiguration) GoString() string {
	return s.String()
}

func (s *InitiateWormConfiguration) SetRetentionPeriodInDays(v int32) *InitiateWormConfiguration {
	s.RetentionPeriodInDays = &v
	return s
}

type InputSerialization struct {
	Csv             *CSVInput  `json:"CSV,omitempty" xml:"CSV,omitempty"`
	CompressionType *string    `json:"CompressionType,omitempty" xml:"CompressionType,omitempty"`
	Json            *JSONInput `json:"JSON,omitempty" xml:"JSON,omitempty"`
}

func (s InputSerialization) String() string {
	return tea.Prettify(s)
}

func (s InputSerialization) GoString() string {
	return s.String()
}

func (s *InputSerialization) SetCsv(v *CSVInput) *InputSerialization {
	s.Csv = v
	return s
}

func (s *InputSerialization) SetCompressionType(v string) *InputSerialization {
	s.CompressionType = &v
	return s
}

func (s *InputSerialization) SetJson(v *JSONInput) *InputSerialization {
	s.Json = v
	return s
}

type InventoryConfiguration struct {
	Destination            *InventoryDestination                 `json:"Destination,omitempty" xml:"Destination,omitempty"`
	Filter                 *InventoryFilter                      `json:"Filter,omitempty" xml:"Filter,omitempty"`
	Id                     *string                               `json:"Id,omitempty" xml:"Id,omitempty"`
	IncludedObjectVersions *string                               `json:"IncludedObjectVersions,omitempty" xml:"IncludedObjectVersions,omitempty"`
	IsEnabled              *bool                                 `json:"IsEnabled,omitempty" xml:"IsEnabled,omitempty"`
	OptionalFields         *InventoryConfigurationOptionalFields `json:"OptionalFields,omitempty" xml:"OptionalFields,omitempty" type:"Struct"`
	Schedule               *InventorySchedule                    `json:"Schedule,omitempty" xml:"Schedule,omitempty"`
}

func (s InventoryConfiguration) String() string {
	return tea.Prettify(s)
}

func (s InventoryConfiguration) GoString() string {
	return s.String()
}

func (s *InventoryConfiguration) SetDestination(v *InventoryDestination) *InventoryConfiguration {
	s.Destination = v
	return s
}

func (s *InventoryConfiguration) SetFilter(v *InventoryFilter) *InventoryConfiguration {
	s.Filter = v
	return s
}

func (s *InventoryConfiguration) SetId(v string) *InventoryConfiguration {
	s.Id = &v
	return s
}

func (s *InventoryConfiguration) SetIncludedObjectVersions(v string) *InventoryConfiguration {
	s.IncludedObjectVersions = &v
	return s
}

func (s *InventoryConfiguration) SetIsEnabled(v bool) *InventoryConfiguration {
	s.IsEnabled = &v
	return s
}

func (s *InventoryConfiguration) SetOptionalFields(v *InventoryConfigurationOptionalFields) *InventoryConfiguration {
	s.OptionalFields = v
	return s
}

func (s *InventoryConfiguration) SetSchedule(v *InventorySchedule) *InventoryConfiguration {
	s.Schedule = v
	return s
}

type InventoryConfigurationOptionalFields struct {
	Fields []*string `json:"Field,omitempty" xml:"Field,omitempty" type:"Repeated"`
}

func (s InventoryConfigurationOptionalFields) String() string {
	return tea.Prettify(s)
}

func (s InventoryConfigurationOptionalFields) GoString() string {
	return s.String()
}

func (s *InventoryConfigurationOptionalFields) SetFields(v []*string) *InventoryConfigurationOptionalFields {
	s.Fields = v
	return s
}

type InventoryDestination struct {
	OSSBucketDestination *InventoryOSSBucketDestination `json:"OSSBucketDestination,omitempty" xml:"OSSBucketDestination,omitempty"`
}

func (s InventoryDestination) String() string {
	return tea.Prettify(s)
}

func (s InventoryDestination) GoString() string {
	return s.String()
}

func (s *InventoryDestination) SetOSSBucketDestination(v *InventoryOSSBucketDestination) *InventoryDestination {
	s.OSSBucketDestination = v
	return s
}

type InventoryEncryption struct {
	SSEKMS *SSEKMS `json:"SSE-KMS,omitempty" xml:"SSE-KMS,omitempty"`
	SSEOSS *string `json:"SSE-OSS,omitempty" xml:"SSE-OSS,omitempty"`
}

func (s InventoryEncryption) String() string {
	return tea.Prettify(s)
}

func (s InventoryEncryption) GoString() string {
	return s.String()
}

func (s *InventoryEncryption) SetSSEKMS(v *SSEKMS) *InventoryEncryption {
	s.SSEKMS = v
	return s
}

func (s *InventoryEncryption) SetSSEOSS(v string) *InventoryEncryption {
	s.SSEOSS = &v
	return s
}

type InventoryFilter struct {
	LastModifyBeginTimeStamp *int64  `json:"LastModifyBeginTimeStamp,omitempty" xml:"LastModifyBeginTimeStamp,omitempty"`
	LastModifyEndTimeStamp   *int64  `json:"LastModifyEndTimeStamp,omitempty" xml:"LastModifyEndTimeStamp,omitempty"`
	LowerSizeBound           *int64  `json:"LowerSizeBound,omitempty" xml:"LowerSizeBound,omitempty"`
	Prefix                   *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	StorageClass             *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
	Tags                     *string `json:"Tags,omitempty" xml:"Tags,omitempty"`
	TagsCondition            *string `json:"TagsCondition,omitempty" xml:"TagsCondition,omitempty"`
	UpperSizeBound           *int64  `json:"UpperSizeBound,omitempty" xml:"UpperSizeBound,omitempty"`
}

func (s InventoryFilter) String() string {
	return tea.Prettify(s)
}

func (s InventoryFilter) GoString() string {
	return s.String()
}

func (s *InventoryFilter) SetLastModifyBeginTimeStamp(v int64) *InventoryFilter {
	s.LastModifyBeginTimeStamp = &v
	return s
}

func (s *InventoryFilter) SetLastModifyEndTimeStamp(v int64) *InventoryFilter {
	s.LastModifyEndTimeStamp = &v
	return s
}

func (s *InventoryFilter) SetLowerSizeBound(v int64) *InventoryFilter {
	s.LowerSizeBound = &v
	return s
}

func (s *InventoryFilter) SetPrefix(v string) *InventoryFilter {
	s.Prefix = &v
	return s
}

func (s *InventoryFilter) SetStorageClass(v string) *InventoryFilter {
	s.StorageClass = &v
	return s
}

func (s *InventoryFilter) SetTags(v string) *InventoryFilter {
	s.Tags = &v
	return s
}

func (s *InventoryFilter) SetTagsCondition(v string) *InventoryFilter {
	s.TagsCondition = &v
	return s
}

func (s *InventoryFilter) SetUpperSizeBound(v int64) *InventoryFilter {
	s.UpperSizeBound = &v
	return s
}

type InventoryOSSBucketDestination struct {
	AccountId  *string              `json:"AccountId,omitempty" xml:"AccountId,omitempty"`
	Bucket     *string              `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	Encryption *InventoryEncryption `json:"Encryption,omitempty" xml:"Encryption,omitempty"`
	Format     *string              `json:"Format,omitempty" xml:"Format,omitempty"`
	Prefix     *string              `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	RoleArn    *string              `json:"RoleArn,omitempty" xml:"RoleArn,omitempty"`
}

func (s InventoryOSSBucketDestination) String() string {
	return tea.Prettify(s)
}

func (s InventoryOSSBucketDestination) GoString() string {
	return s.String()
}

func (s *InventoryOSSBucketDestination) SetAccountId(v string) *InventoryOSSBucketDestination {
	s.AccountId = &v
	return s
}

func (s *InventoryOSSBucketDestination) SetBucket(v string) *InventoryOSSBucketDestination {
	s.Bucket = &v
	return s
}

func (s *InventoryOSSBucketDestination) SetEncryption(v *InventoryEncryption) *InventoryOSSBucketDestination {
	s.Encryption = v
	return s
}

func (s *InventoryOSSBucketDestination) SetFormat(v string) *InventoryOSSBucketDestination {
	s.Format = &v
	return s
}

func (s *InventoryOSSBucketDestination) SetPrefix(v string) *InventoryOSSBucketDestination {
	s.Prefix = &v
	return s
}

func (s *InventoryOSSBucketDestination) SetRoleArn(v string) *InventoryOSSBucketDestination {
	s.RoleArn = &v
	return s
}

type InventorySchedule struct {
	Frequency *string `json:"Frequency,omitempty" xml:"Frequency,omitempty"`
}

func (s InventorySchedule) String() string {
	return tea.Prettify(s)
}

func (s InventorySchedule) GoString() string {
	return s.String()
}

func (s *InventorySchedule) SetFrequency(v string) *InventorySchedule {
	s.Frequency = &v
	return s
}

type JSONInput struct {
	ParseJsonNumberAsString *bool   `json:"ParseJsonNumberAsString,omitempty" xml:"ParseJsonNumberAsString,omitempty"`
	Range                   *string `json:"Range,omitempty" xml:"Range,omitempty"`
	Type                    *string `json:"Type,omitempty" xml:"Type,omitempty"`
}

func (s JSONInput) String() string {
	return tea.Prettify(s)
}

func (s JSONInput) GoString() string {
	return s.String()
}

func (s *JSONInput) SetParseJsonNumberAsString(v bool) *JSONInput {
	s.ParseJsonNumberAsString = &v
	return s
}

func (s *JSONInput) SetRange(v string) *JSONInput {
	s.Range = &v
	return s
}

func (s *JSONInput) SetType(v string) *JSONInput {
	s.Type = &v
	return s
}

type JSONOutput struct {
	RecordDelimiter *string `json:"RecordDelimiter,omitempty" xml:"RecordDelimiter,omitempty"`
}

func (s JSONOutput) String() string {
	return tea.Prettify(s)
}

func (s JSONOutput) GoString() string {
	return s.String()
}

func (s *JSONOutput) SetRecordDelimiter(v string) *JSONOutput {
	s.RecordDelimiter = &v
	return s
}

type LifecycleConfiguration struct {
	Rule []*LifecycleRule `json:"Rule,omitempty" xml:"Rule,omitempty" type:"Repeated"`
}

func (s LifecycleConfiguration) String() string {
	return tea.Prettify(s)
}

func (s LifecycleConfiguration) GoString() string {
	return s.String()
}

func (s *LifecycleConfiguration) SetRule(v []*LifecycleRule) *LifecycleConfiguration {
	s.Rule = v
	return s
}

type LifecycleRule struct {
	LifecycleAbortMultipartUpload *LifecycleRuleLifecycleAbortMultipartUpload `json:"AbortMultipartUpload,omitempty" xml:"AbortMultipartUpload,omitempty" type:"Struct"`
	AtimeBase                     *int64                                      `json:"AtimeBase,omitempty" xml:"AtimeBase,omitempty"`
	LifecycleExpiration           *LifecycleRuleLifecycleExpiration           `json:"Expiration,omitempty" xml:"Expiration,omitempty" type:"Struct"`
	Filter                        *LifecycleRuleFilter                        `json:"Filter,omitempty" xml:"Filter,omitempty" type:"Struct"`
	ID                            *string                                     `json:"ID,omitempty" xml:"ID,omitempty"`
	NoncurrentVersionExpiration   *LifecycleRuleNoncurrentVersionExpiration   `json:"NoncurrentVersionExpiration,omitempty" xml:"NoncurrentVersionExpiration,omitempty" type:"Struct"`
	NoncurrentVersionTransition   []*LifecycleRuleNoncurrentVersionTransition `json:"NoncurrentVersionTransition,omitempty" xml:"NoncurrentVersionTransition,omitempty" type:"Repeated"`
	Prefix                        *string                                     `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	Status                        *string                                     `json:"Status,omitempty" xml:"Status,omitempty"`
	Tag                           []*Tag                                      `json:"Tag,omitempty" xml:"Tag,omitempty" type:"Repeated"`
	LifecycleTransition           []*LifecycleRuleLifecycleTransition         `json:"Transition,omitempty" xml:"Transition,omitempty" type:"Repeated"`
}

func (s LifecycleRule) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRule) GoString() string {
	return s.String()
}

func (s *LifecycleRule) SetLifecycleAbortMultipartUpload(v *LifecycleRuleLifecycleAbortMultipartUpload) *LifecycleRule {
	s.LifecycleAbortMultipartUpload = v
	return s
}

func (s *LifecycleRule) SetAtimeBase(v int64) *LifecycleRule {
	s.AtimeBase = &v
	return s
}

func (s *LifecycleRule) SetLifecycleExpiration(v *LifecycleRuleLifecycleExpiration) *LifecycleRule {
	s.LifecycleExpiration = v
	return s
}

func (s *LifecycleRule) SetFilter(v *LifecycleRuleFilter) *LifecycleRule {
	s.Filter = v
	return s
}

func (s *LifecycleRule) SetID(v string) *LifecycleRule {
	s.ID = &v
	return s
}

func (s *LifecycleRule) SetNoncurrentVersionExpiration(v *LifecycleRuleNoncurrentVersionExpiration) *LifecycleRule {
	s.NoncurrentVersionExpiration = v
	return s
}

func (s *LifecycleRule) SetNoncurrentVersionTransition(v []*LifecycleRuleNoncurrentVersionTransition) *LifecycleRule {
	s.NoncurrentVersionTransition = v
	return s
}

func (s *LifecycleRule) SetPrefix(v string) *LifecycleRule {
	s.Prefix = &v
	return s
}

func (s *LifecycleRule) SetStatus(v string) *LifecycleRule {
	s.Status = &v
	return s
}

func (s *LifecycleRule) SetTag(v []*Tag) *LifecycleRule {
	s.Tag = v
	return s
}

func (s *LifecycleRule) SetLifecycleTransition(v []*LifecycleRuleLifecycleTransition) *LifecycleRule {
	s.LifecycleTransition = v
	return s
}

type LifecycleRuleLifecycleAbortMultipartUpload struct {
	CreatedBeforeDate *string `json:"CreatedBeforeDate,omitempty" xml:"CreatedBeforeDate,omitempty"`
	Date              *string `json:"Date,omitempty" xml:"Date,omitempty"`
	Days              *int32  `json:"Days,omitempty" xml:"Days,omitempty"`
}

func (s LifecycleRuleLifecycleAbortMultipartUpload) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRuleLifecycleAbortMultipartUpload) GoString() string {
	return s.String()
}

func (s *LifecycleRuleLifecycleAbortMultipartUpload) SetCreatedBeforeDate(v string) *LifecycleRuleLifecycleAbortMultipartUpload {
	s.CreatedBeforeDate = &v
	return s
}

func (s *LifecycleRuleLifecycleAbortMultipartUpload) SetDate(v string) *LifecycleRuleLifecycleAbortMultipartUpload {
	s.Date = &v
	return s
}

func (s *LifecycleRuleLifecycleAbortMultipartUpload) SetDays(v int32) *LifecycleRuleLifecycleAbortMultipartUpload {
	s.Days = &v
	return s
}

type LifecycleRuleLifecycleExpiration struct {
	CreatedBeforeDate         *string `json:"CreatedBeforeDate,omitempty" xml:"CreatedBeforeDate,omitempty"`
	Date                      *string `json:"Date,omitempty" xml:"Date,omitempty"`
	Days                      *int32  `json:"Days,omitempty" xml:"Days,omitempty"`
	ExpiredObjectDeleteMarker *bool   `json:"ExpiredObjectDeleteMarker,omitempty" xml:"ExpiredObjectDeleteMarker,omitempty"`
}

func (s LifecycleRuleLifecycleExpiration) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRuleLifecycleExpiration) GoString() string {
	return s.String()
}

func (s *LifecycleRuleLifecycleExpiration) SetCreatedBeforeDate(v string) *LifecycleRuleLifecycleExpiration {
	s.CreatedBeforeDate = &v
	return s
}

func (s *LifecycleRuleLifecycleExpiration) SetDate(v string) *LifecycleRuleLifecycleExpiration {
	s.Date = &v
	return s
}

func (s *LifecycleRuleLifecycleExpiration) SetDays(v int32) *LifecycleRuleLifecycleExpiration {
	s.Days = &v
	return s
}

func (s *LifecycleRuleLifecycleExpiration) SetExpiredObjectDeleteMarker(v bool) *LifecycleRuleLifecycleExpiration {
	s.ExpiredObjectDeleteMarker = &v
	return s
}

type LifecycleRuleFilter struct {
	Not                   *LifecycleRuleFilterNot `json:"Not,omitempty" xml:"Not,omitempty" type:"Struct"`
	ObjectSizeGreaterThan *int64                  `json:"ObjectSizeGreaterThan,omitempty" xml:"ObjectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64                  `json:"ObjectSizeLessThan,omitempty" xml:"ObjectSizeLessThan,omitempty"`
}

func (s LifecycleRuleFilter) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRuleFilter) GoString() string {
	return s.String()
}

func (s *LifecycleRuleFilter) SetNot(v *LifecycleRuleFilterNot) *LifecycleRuleFilter {
	s.Not = v
	return s
}

func (s *LifecycleRuleFilter) SetObjectSizeGreaterThan(v int64) *LifecycleRuleFilter {
	s.ObjectSizeGreaterThan = &v
	return s
}

func (s *LifecycleRuleFilter) SetObjectSizeLessThan(v int64) *LifecycleRuleFilter {
	s.ObjectSizeLessThan = &v
	return s
}

type LifecycleRuleFilterNot struct {
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	Tag    *Tag    `json:"Tag,omitempty" xml:"Tag,omitempty"`
}

func (s LifecycleRuleFilterNot) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRuleFilterNot) GoString() string {
	return s.String()
}

func (s *LifecycleRuleFilterNot) SetPrefix(v string) *LifecycleRuleFilterNot {
	s.Prefix = &v
	return s
}

func (s *LifecycleRuleFilterNot) SetTag(v *Tag) *LifecycleRuleFilterNot {
	s.Tag = v
	return s
}

type LifecycleRuleNoncurrentVersionExpiration struct {
	NoncurrentDays *int32 `json:"NoncurrentDays,omitempty" xml:"NoncurrentDays,omitempty"`
}

func (s LifecycleRuleNoncurrentVersionExpiration) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRuleNoncurrentVersionExpiration) GoString() string {
	return s.String()
}

func (s *LifecycleRuleNoncurrentVersionExpiration) SetNoncurrentDays(v int32) *LifecycleRuleNoncurrentVersionExpiration {
	s.NoncurrentDays = &v
	return s
}

type LifecycleRuleNoncurrentVersionTransition struct {
	AllowSmallFile       *bool   `json:"AllowSmallFile,omitempty" xml:"AllowSmallFile,omitempty"`
	IsAccessTime         *bool   `json:"IsAccessTime,omitempty" xml:"IsAccessTime,omitempty"`
	NoncurrentDays       *int32  `json:"NoncurrentDays,omitempty" xml:"NoncurrentDays,omitempty"`
	ReturnToStdWhenVisit *bool   `json:"ReturnToStdWhenVisit,omitempty" xml:"ReturnToStdWhenVisit,omitempty"`
	StorageClass         *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
}

func (s LifecycleRuleNoncurrentVersionTransition) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRuleNoncurrentVersionTransition) GoString() string {
	return s.String()
}

func (s *LifecycleRuleNoncurrentVersionTransition) SetAllowSmallFile(v bool) *LifecycleRuleNoncurrentVersionTransition {
	s.AllowSmallFile = &v
	return s
}

func (s *LifecycleRuleNoncurrentVersionTransition) SetIsAccessTime(v bool) *LifecycleRuleNoncurrentVersionTransition {
	s.IsAccessTime = &v
	return s
}

func (s *LifecycleRuleNoncurrentVersionTransition) SetNoncurrentDays(v int32) *LifecycleRuleNoncurrentVersionTransition {
	s.NoncurrentDays = &v
	return s
}

func (s *LifecycleRuleNoncurrentVersionTransition) SetReturnToStdWhenVisit(v bool) *LifecycleRuleNoncurrentVersionTransition {
	s.ReturnToStdWhenVisit = &v
	return s
}

func (s *LifecycleRuleNoncurrentVersionTransition) SetStorageClass(v string) *LifecycleRuleNoncurrentVersionTransition {
	s.StorageClass = &v
	return s
}

type LifecycleRuleLifecycleTransition struct {
	AllowSmallFile       *bool   `json:"AllowSmallFile,omitempty" xml:"AllowSmallFile,omitempty"`
	CreatedBeforeDate    *string `json:"CreatedBeforeDate,omitempty" xml:"CreatedBeforeDate,omitempty"`
	Days                 *int32  `json:"Days,omitempty" xml:"Days,omitempty"`
	IsAccessTime         *bool   `json:"IsAccessTime,omitempty" xml:"IsAccessTime,omitempty"`
	ReturnToStdWhenVisit *bool   `json:"ReturnToStdWhenVisit,omitempty" xml:"ReturnToStdWhenVisit,omitempty"`
	StorageClass         *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
}

func (s LifecycleRuleLifecycleTransition) String() string {
	return tea.Prettify(s)
}

func (s LifecycleRuleLifecycleTransition) GoString() string {
	return s.String()
}

func (s *LifecycleRuleLifecycleTransition) SetAllowSmallFile(v bool) *LifecycleRuleLifecycleTransition {
	s.AllowSmallFile = &v
	return s
}

func (s *LifecycleRuleLifecycleTransition) SetCreatedBeforeDate(v string) *LifecycleRuleLifecycleTransition {
	s.CreatedBeforeDate = &v
	return s
}

func (s *LifecycleRuleLifecycleTransition) SetDays(v int32) *LifecycleRuleLifecycleTransition {
	s.Days = &v
	return s
}

func (s *LifecycleRuleLifecycleTransition) SetIsAccessTime(v bool) *LifecycleRuleLifecycleTransition {
	s.IsAccessTime = &v
	return s
}

func (s *LifecycleRuleLifecycleTransition) SetReturnToStdWhenVisit(v bool) *LifecycleRuleLifecycleTransition {
	s.ReturnToStdWhenVisit = &v
	return s
}

func (s *LifecycleRuleLifecycleTransition) SetStorageClass(v string) *LifecycleRuleLifecycleTransition {
	s.StorageClass = &v
	return s
}

type ListAccessPointsResult struct {
	AccessPoints          *ListAccessPointsResultAccessPoints `json:"AccessPoints,omitempty" xml:"AccessPoints,omitempty" type:"Struct"`
	AccountId             *string                             `json:"AccountId,omitempty" xml:"AccountId,omitempty"`
	IsTruncated           *string                             `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	MaxKeys               *int32                              `json:"MaxKeys,omitempty" xml:"MaxKeys,omitempty"`
	NextContinuationToken *string                             `json:"NextContinuationToken,omitempty" xml:"NextContinuationToken,omitempty"`
}

func (s ListAccessPointsResult) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsResult) GoString() string {
	return s.String()
}

func (s *ListAccessPointsResult) SetAccessPoints(v *ListAccessPointsResultAccessPoints) *ListAccessPointsResult {
	s.AccessPoints = v
	return s
}

func (s *ListAccessPointsResult) SetAccountId(v string) *ListAccessPointsResult {
	s.AccountId = &v
	return s
}

func (s *ListAccessPointsResult) SetIsTruncated(v string) *ListAccessPointsResult {
	s.IsTruncated = &v
	return s
}

func (s *ListAccessPointsResult) SetMaxKeys(v int32) *ListAccessPointsResult {
	s.MaxKeys = &v
	return s
}

func (s *ListAccessPointsResult) SetNextContinuationToken(v string) *ListAccessPointsResult {
	s.NextContinuationToken = &v
	return s
}

type ListAccessPointsResultAccessPoints struct {
	AccessPoint []*AccessPoint `json:"AccessPoint,omitempty" xml:"AccessPoint,omitempty" type:"Repeated"`
}

func (s ListAccessPointsResultAccessPoints) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsResultAccessPoints) GoString() string {
	return s.String()
}

func (s *ListAccessPointsResultAccessPoints) SetAccessPoint(v []*AccessPoint) *ListAccessPointsResultAccessPoints {
	s.AccessPoint = v
	return s
}

type LiveChannel struct {
	Description  *string                 `json:"Description,omitempty" xml:"Description,omitempty"`
	LastModified *string                 `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
	Name         *string                 `json:"Name,omitempty" xml:"Name,omitempty"`
	PlayUrls     *LiveChannelPlayUrls    `json:"PlayUrls,omitempty" xml:"PlayUrls,omitempty"`
	PublishUrls  *LiveChannelPublishUrls `json:"PublishUrls,omitempty" xml:"PublishUrls,omitempty"`
	Status       *string                 `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s LiveChannel) String() string {
	return tea.Prettify(s)
}

func (s LiveChannel) GoString() string {
	return s.String()
}

func (s *LiveChannel) SetDescription(v string) *LiveChannel {
	s.Description = &v
	return s
}

func (s *LiveChannel) SetLastModified(v string) *LiveChannel {
	s.LastModified = &v
	return s
}

func (s *LiveChannel) SetName(v string) *LiveChannel {
	s.Name = &v
	return s
}

func (s *LiveChannel) SetPlayUrls(v *LiveChannelPlayUrls) *LiveChannel {
	s.PlayUrls = v
	return s
}

func (s *LiveChannel) SetPublishUrls(v *LiveChannelPublishUrls) *LiveChannel {
	s.PublishUrls = v
	return s
}

func (s *LiveChannel) SetStatus(v string) *LiveChannel {
	s.Status = &v
	return s
}

type LiveChannelAudio struct {
	Bandwidth  *int64  `json:"Bandwidth,omitempty" xml:"Bandwidth,omitempty"`
	Codec      *string `json:"Codec,omitempty" xml:"Codec,omitempty"`
	SampleRate *int64  `json:"SampleRate,omitempty" xml:"SampleRate,omitempty"`
}

func (s LiveChannelAudio) String() string {
	return tea.Prettify(s)
}

func (s LiveChannelAudio) GoString() string {
	return s.String()
}

func (s *LiveChannelAudio) SetBandwidth(v int64) *LiveChannelAudio {
	s.Bandwidth = &v
	return s
}

func (s *LiveChannelAudio) SetCodec(v string) *LiveChannelAudio {
	s.Codec = &v
	return s
}

func (s *LiveChannelAudio) SetSampleRate(v int64) *LiveChannelAudio {
	s.SampleRate = &v
	return s
}

type LiveChannelConfiguration struct {
	Description *string              `json:"Description,omitempty" xml:"Description,omitempty"`
	Snapshot    *LiveChannelSnapshot `json:"Snapshot,omitempty" xml:"Snapshot,omitempty"`
	Status      *string              `json:"Status,omitempty" xml:"Status,omitempty"`
	Target      *LiveChannelTarget   `json:"Target,omitempty" xml:"Target,omitempty"`
}

func (s LiveChannelConfiguration) String() string {
	return tea.Prettify(s)
}

func (s LiveChannelConfiguration) GoString() string {
	return s.String()
}

func (s *LiveChannelConfiguration) SetDescription(v string) *LiveChannelConfiguration {
	s.Description = &v
	return s
}

func (s *LiveChannelConfiguration) SetSnapshot(v *LiveChannelSnapshot) *LiveChannelConfiguration {
	s.Snapshot = v
	return s
}

func (s *LiveChannelConfiguration) SetStatus(v string) *LiveChannelConfiguration {
	s.Status = &v
	return s
}

func (s *LiveChannelConfiguration) SetTarget(v *LiveChannelTarget) *LiveChannelConfiguration {
	s.Target = v
	return s
}

type LiveChannelPlayUrls struct {
	Url *string `json:"Url,omitempty" xml:"Url,omitempty"`
}

func (s LiveChannelPlayUrls) String() string {
	return tea.Prettify(s)
}

func (s LiveChannelPlayUrls) GoString() string {
	return s.String()
}

func (s *LiveChannelPlayUrls) SetUrl(v string) *LiveChannelPlayUrls {
	s.Url = &v
	return s
}

type LiveChannelPublishUrls struct {
	Url *string `json:"Url,omitempty" xml:"Url,omitempty"`
}

func (s LiveChannelPublishUrls) String() string {
	return tea.Prettify(s)
}

func (s LiveChannelPublishUrls) GoString() string {
	return s.String()
}

func (s *LiveChannelPublishUrls) SetUrl(v string) *LiveChannelPublishUrls {
	s.Url = &v
	return s
}

type LiveChannelSnapshot struct {
	DestBucket  *string `json:"DestBucket,omitempty" xml:"DestBucket,omitempty"`
	Interval    *int64  `json:"Interval,omitempty" xml:"Interval,omitempty"`
	NotifyTopic *string `json:"NotifyTopic,omitempty" xml:"NotifyTopic,omitempty"`
	RoleName    *string `json:"RoleName,omitempty" xml:"RoleName,omitempty"`
}

func (s LiveChannelSnapshot) String() string {
	return tea.Prettify(s)
}

func (s LiveChannelSnapshot) GoString() string {
	return s.String()
}

func (s *LiveChannelSnapshot) SetDestBucket(v string) *LiveChannelSnapshot {
	s.DestBucket = &v
	return s
}

func (s *LiveChannelSnapshot) SetInterval(v int64) *LiveChannelSnapshot {
	s.Interval = &v
	return s
}

func (s *LiveChannelSnapshot) SetNotifyTopic(v string) *LiveChannelSnapshot {
	s.NotifyTopic = &v
	return s
}

func (s *LiveChannelSnapshot) SetRoleName(v string) *LiveChannelSnapshot {
	s.RoleName = &v
	return s
}

type LiveChannelTarget struct {
	FragCount    *int64  `json:"FragCount,omitempty" xml:"FragCount,omitempty"`
	FragDuration *int64  `json:"FragDuration,omitempty" xml:"FragDuration,omitempty"`
	PlaylistName *string `json:"PlaylistName,omitempty" xml:"PlaylistName,omitempty"`
	Type         *string `json:"Type,omitempty" xml:"Type,omitempty"`
}

func (s LiveChannelTarget) String() string {
	return tea.Prettify(s)
}

func (s LiveChannelTarget) GoString() string {
	return s.String()
}

func (s *LiveChannelTarget) SetFragCount(v int64) *LiveChannelTarget {
	s.FragCount = &v
	return s
}

func (s *LiveChannelTarget) SetFragDuration(v int64) *LiveChannelTarget {
	s.FragDuration = &v
	return s
}

func (s *LiveChannelTarget) SetPlaylistName(v string) *LiveChannelTarget {
	s.PlaylistName = &v
	return s
}

func (s *LiveChannelTarget) SetType(v string) *LiveChannelTarget {
	s.Type = &v
	return s
}

type LiveChannelVideo struct {
	Bandwidth *int64  `json:"Bandwidth,omitempty" xml:"Bandwidth,omitempty"`
	Codec     *string `json:"Codec,omitempty" xml:"Codec,omitempty"`
	FrameRate *int64  `json:"FrameRate,omitempty" xml:"FrameRate,omitempty"`
	Height    *int64  `json:"Height,omitempty" xml:"Height,omitempty"`
	Width     *int64  `json:"Width,omitempty" xml:"Width,omitempty"`
}

func (s LiveChannelVideo) String() string {
	return tea.Prettify(s)
}

func (s LiveChannelVideo) GoString() string {
	return s.String()
}

func (s *LiveChannelVideo) SetBandwidth(v int64) *LiveChannelVideo {
	s.Bandwidth = &v
	return s
}

func (s *LiveChannelVideo) SetCodec(v string) *LiveChannelVideo {
	s.Codec = &v
	return s
}

func (s *LiveChannelVideo) SetFrameRate(v int64) *LiveChannelVideo {
	s.FrameRate = &v
	return s
}

func (s *LiveChannelVideo) SetHeight(v int64) *LiveChannelVideo {
	s.Height = &v
	return s
}

func (s *LiveChannelVideo) SetWidth(v int64) *LiveChannelVideo {
	s.Width = &v
	return s
}

type LiveRecord struct {
	EndTime    *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	RemoteAddr *string `json:"RemoteAddr,omitempty" xml:"RemoteAddr,omitempty"`
	StartTime  *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s LiveRecord) String() string {
	return tea.Prettify(s)
}

func (s LiveRecord) GoString() string {
	return s.String()
}

func (s *LiveRecord) SetEndTime(v string) *LiveRecord {
	s.EndTime = &v
	return s
}

func (s *LiveRecord) SetRemoteAddr(v string) *LiveRecord {
	s.RemoteAddr = &v
	return s
}

func (s *LiveRecord) SetStartTime(v string) *LiveRecord {
	s.StartTime = &v
	return s
}

type LocationTransferType struct {
	Location      *string                            `json:"Location,omitempty" xml:"Location,omitempty"`
	TransferTypes *LocationTransferTypeTransferTypes `json:"TransferTypes,omitempty" xml:"TransferTypes,omitempty" type:"Struct"`
}

func (s LocationTransferType) String() string {
	return tea.Prettify(s)
}

func (s LocationTransferType) GoString() string {
	return s.String()
}

func (s *LocationTransferType) SetLocation(v string) *LocationTransferType {
	s.Location = &v
	return s
}

func (s *LocationTransferType) SetTransferTypes(v *LocationTransferTypeTransferTypes) *LocationTransferType {
	s.TransferTypes = v
	return s
}

type LocationTransferTypeTransferTypes struct {
	Type []*string `json:"Type,omitempty" xml:"Type,omitempty" type:"Repeated"`
}

func (s LocationTransferTypeTransferTypes) String() string {
	return tea.Prettify(s)
}

func (s LocationTransferTypeTransferTypes) GoString() string {
	return s.String()
}

func (s *LocationTransferTypeTransferTypes) SetType(v []*string) *LocationTransferTypeTransferTypes {
	s.Type = v
	return s
}

type LoggingEnabled struct {
	TargetBucket *string `json:"TargetBucket,omitempty" xml:"TargetBucket,omitempty"`
	TargetPrefix *string `json:"TargetPrefix,omitempty" xml:"TargetPrefix,omitempty"`
}

func (s LoggingEnabled) String() string {
	return tea.Prettify(s)
}

func (s LoggingEnabled) GoString() string {
	return s.String()
}

func (s *LoggingEnabled) SetTargetBucket(v string) *LoggingEnabled {
	s.TargetBucket = &v
	return s
}

func (s *LoggingEnabled) SetTargetPrefix(v string) *LoggingEnabled {
	s.TargetPrefix = &v
	return s
}

type MetaQuery struct {
	Aggregations *MetaQueryAggregations `json:"Aggregations,omitempty" xml:"Aggregations,omitempty" type:"Struct"`
	MaxResults   *int64                 `json:"MaxResults,omitempty" xml:"MaxResults,omitempty"`
	NextToken    *string                `json:"NextToken,omitempty" xml:"NextToken,omitempty"`
	Order        *string                `json:"Order,omitempty" xml:"Order,omitempty"`
	Query        *string                `json:"Query,omitempty" xml:"Query,omitempty"`
	Sort         *string                `json:"Sort,omitempty" xml:"Sort,omitempty"`
}

func (s MetaQuery) String() string {
	return tea.Prettify(s)
}

func (s MetaQuery) GoString() string {
	return s.String()
}

func (s *MetaQuery) SetAggregations(v *MetaQueryAggregations) *MetaQuery {
	s.Aggregations = v
	return s
}

func (s *MetaQuery) SetMaxResults(v int64) *MetaQuery {
	s.MaxResults = &v
	return s
}

func (s *MetaQuery) SetNextToken(v string) *MetaQuery {
	s.NextToken = &v
	return s
}

func (s *MetaQuery) SetOrder(v string) *MetaQuery {
	s.Order = &v
	return s
}

func (s *MetaQuery) SetQuery(v string) *MetaQuery {
	s.Query = &v
	return s
}

func (s *MetaQuery) SetSort(v string) *MetaQuery {
	s.Sort = &v
	return s
}

type MetaQueryAggregations struct {
	Aggregation []*MetaQueryAggregation `json:"Aggregation,omitempty" xml:"Aggregation,omitempty" type:"Repeated"`
}

func (s MetaQueryAggregations) String() string {
	return tea.Prettify(s)
}

func (s MetaQueryAggregations) GoString() string {
	return s.String()
}

func (s *MetaQueryAggregations) SetAggregation(v []*MetaQueryAggregation) *MetaQueryAggregations {
	s.Aggregation = v
	return s
}

type MetaQueryAggregation struct {
	Field     *string `json:"Field,omitempty" xml:"Field,omitempty"`
	Operation *string `json:"Operation,omitempty" xml:"Operation,omitempty"`
}

func (s MetaQueryAggregation) String() string {
	return tea.Prettify(s)
}

func (s MetaQueryAggregation) GoString() string {
	return s.String()
}

func (s *MetaQueryAggregation) SetField(v string) *MetaQueryAggregation {
	s.Field = &v
	return s
}

func (s *MetaQueryAggregation) SetOperation(v string) *MetaQueryAggregation {
	s.Operation = &v
	return s
}

type MetaQueryFile struct {
	ETag                                  *string                   `json:"ETag,omitempty" xml:"ETag,omitempty"`
	FileModifiedTime                      *string                   `json:"FileModifiedTime,omitempty" xml:"FileModifiedTime,omitempty"`
	Filename                              *string                   `json:"Filename,omitempty" xml:"Filename,omitempty"`
	OSSCRC64                              *string                   `json:"OSSCRC64,omitempty" xml:"OSSCRC64,omitempty"`
	OSSObjectType                         *string                   `json:"OSSObjectType,omitempty" xml:"OSSObjectType,omitempty"`
	OSSStorageClass                       *string                   `json:"OSSStorageClass,omitempty" xml:"OSSStorageClass,omitempty"`
	OSSTagging                            *MetaQueryFileOSSTagging  `json:"OSSTagging,omitempty" xml:"OSSTagging,omitempty" type:"Struct"`
	OSSTaggingCount                       *int64                    `json:"OSSTaggingCount,omitempty" xml:"OSSTaggingCount,omitempty"`
	OSSUserMeta                           *MetaQueryFileOSSUserMeta `json:"OSSUserMeta,omitempty" xml:"OSSUserMeta,omitempty" type:"Struct"`
	ObjectACL                             *string                   `json:"ObjectACL,omitempty" xml:"ObjectACL,omitempty"`
	ServerSideEncryption                  *string                   `json:"ServerSideEncryption,omitempty" xml:"ServerSideEncryption,omitempty"`
	ServerSideEncryptionCustomerAlgorithm *string                   `json:"ServerSideEncryptionCustomerAlgorithm,omitempty" xml:"ServerSideEncryptionCustomerAlgorithm,omitempty"`
	Size                                  *int64                    `json:"Size,omitempty" xml:"Size,omitempty"`
}

func (s MetaQueryFile) String() string {
	return tea.Prettify(s)
}

func (s MetaQueryFile) GoString() string {
	return s.String()
}

func (s *MetaQueryFile) SetETag(v string) *MetaQueryFile {
	s.ETag = &v
	return s
}

func (s *MetaQueryFile) SetFileModifiedTime(v string) *MetaQueryFile {
	s.FileModifiedTime = &v
	return s
}

func (s *MetaQueryFile) SetFilename(v string) *MetaQueryFile {
	s.Filename = &v
	return s
}

func (s *MetaQueryFile) SetOSSCRC64(v string) *MetaQueryFile {
	s.OSSCRC64 = &v
	return s
}

func (s *MetaQueryFile) SetOSSObjectType(v string) *MetaQueryFile {
	s.OSSObjectType = &v
	return s
}

func (s *MetaQueryFile) SetOSSStorageClass(v string) *MetaQueryFile {
	s.OSSStorageClass = &v
	return s
}

func (s *MetaQueryFile) SetOSSTagging(v *MetaQueryFileOSSTagging) *MetaQueryFile {
	s.OSSTagging = v
	return s
}

func (s *MetaQueryFile) SetOSSTaggingCount(v int64) *MetaQueryFile {
	s.OSSTaggingCount = &v
	return s
}

func (s *MetaQueryFile) SetOSSUserMeta(v *MetaQueryFileOSSUserMeta) *MetaQueryFile {
	s.OSSUserMeta = v
	return s
}

func (s *MetaQueryFile) SetObjectACL(v string) *MetaQueryFile {
	s.ObjectACL = &v
	return s
}

func (s *MetaQueryFile) SetServerSideEncryption(v string) *MetaQueryFile {
	s.ServerSideEncryption = &v
	return s
}

func (s *MetaQueryFile) SetServerSideEncryptionCustomerAlgorithm(v string) *MetaQueryFile {
	s.ServerSideEncryptionCustomerAlgorithm = &v
	return s
}

func (s *MetaQueryFile) SetSize(v int64) *MetaQueryFile {
	s.Size = &v
	return s
}

type MetaQueryFileOSSTagging struct {
	Tagging []*MetaQueryTagging `json:"Tagging,omitempty" xml:"Tagging,omitempty" type:"Repeated"`
}

func (s MetaQueryFileOSSTagging) String() string {
	return tea.Prettify(s)
}

func (s MetaQueryFileOSSTagging) GoString() string {
	return s.String()
}

func (s *MetaQueryFileOSSTagging) SetTagging(v []*MetaQueryTagging) *MetaQueryFileOSSTagging {
	s.Tagging = v
	return s
}

type MetaQueryFileOSSUserMeta struct {
	UserMeta []*MetaQueryUserMeta `json:"UserMeta,omitempty" xml:"UserMeta,omitempty" type:"Repeated"`
}

func (s MetaQueryFileOSSUserMeta) String() string {
	return tea.Prettify(s)
}

func (s MetaQueryFileOSSUserMeta) GoString() string {
	return s.String()
}

func (s *MetaQueryFileOSSUserMeta) SetUserMeta(v []*MetaQueryUserMeta) *MetaQueryFileOSSUserMeta {
	s.UserMeta = v
	return s
}

type MetaQueryTagging struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s MetaQueryTagging) String() string {
	return tea.Prettify(s)
}

func (s MetaQueryTagging) GoString() string {
	return s.String()
}

func (s *MetaQueryTagging) SetKey(v string) *MetaQueryTagging {
	s.Key = &v
	return s
}

func (s *MetaQueryTagging) SetValue(v string) *MetaQueryTagging {
	s.Value = &v
	return s
}

type MetaQueryUserMeta struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s MetaQueryUserMeta) String() string {
	return tea.Prettify(s)
}

func (s MetaQueryUserMeta) GoString() string {
	return s.String()
}

func (s *MetaQueryUserMeta) SetKey(v string) *MetaQueryUserMeta {
	s.Key = &v
	return s
}

func (s *MetaQueryUserMeta) SetValue(v string) *MetaQueryUserMeta {
	s.Value = &v
	return s
}

type ObjectIdentifier struct {
	Key       *string `json:"Key,omitempty" xml:"Key,omitempty"`
	VersionId *string `json:"VersionId,omitempty" xml:"VersionId,omitempty"`
}

func (s ObjectIdentifier) String() string {
	return tea.Prettify(s)
}

func (s ObjectIdentifier) GoString() string {
	return s.String()
}

func (s *ObjectIdentifier) SetKey(v string) *ObjectIdentifier {
	s.Key = &v
	return s
}

func (s *ObjectIdentifier) SetVersionId(v string) *ObjectIdentifier {
	s.VersionId = &v
	return s
}

type ObjectProcessConfiguration struct {
	AllowedFeatures              *ObjectProcessConfigurationAllowedFeatures              `json:"AllowedFeatures,omitempty" xml:"AllowedFeatures,omitempty" type:"Struct"`
	TransformationConfigurations *ObjectProcessConfigurationTransformationConfigurations `json:"TransformationConfigurations,omitempty" xml:"TransformationConfigurations,omitempty" type:"Struct"`
}

func (s ObjectProcessConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfiguration) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfiguration) SetAllowedFeatures(v *ObjectProcessConfigurationAllowedFeatures) *ObjectProcessConfiguration {
	s.AllowedFeatures = v
	return s
}

func (s *ObjectProcessConfiguration) SetTransformationConfigurations(v *ObjectProcessConfigurationTransformationConfigurations) *ObjectProcessConfiguration {
	s.TransformationConfigurations = v
	return s
}

type ObjectProcessConfigurationAllowedFeatures struct {
	AllowedFeature []*string `json:"AllowedFeature,omitempty" xml:"AllowedFeature,omitempty" type:"Repeated"`
}

func (s ObjectProcessConfigurationAllowedFeatures) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationAllowedFeatures) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationAllowedFeatures) SetAllowedFeature(v []*string) *ObjectProcessConfigurationAllowedFeatures {
	s.AllowedFeature = v
	return s
}

type ObjectProcessConfigurationTransformationConfigurations struct {
	TransformationConfiguration []*ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration `json:"TransformationConfiguration,omitempty" xml:"TransformationConfiguration,omitempty" type:"Repeated"`
}

func (s ObjectProcessConfigurationTransformationConfigurations) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationTransformationConfigurations) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationTransformationConfigurations) SetTransformationConfiguration(v []*ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration) *ObjectProcessConfigurationTransformationConfigurations {
	s.TransformationConfiguration = v
	return s
}

type ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration struct {
	Actions               *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationActions               `json:"Actions,omitempty" xml:"Actions,omitempty" type:"Struct"`
	ContentTransformation *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation `json:"ContentTransformation,omitempty" xml:"ContentTransformation,omitempty" type:"Struct"`
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration) SetActions(v *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationActions) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration {
	s.Actions = v
	return s
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration) SetContentTransformation(v *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfiguration {
	s.ContentTransformation = v
	return s
}

type ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationActions struct {
	Action []*string `json:"Action,omitempty" xml:"Action,omitempty" type:"Repeated"`
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationActions) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationActions) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationActions) SetAction(v []*string) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationActions {
	s.Action = v
	return s
}

type ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation struct {
	AdditionalFeatures *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeatures `json:"AdditionalFeatures,omitempty" xml:"AdditionalFeatures,omitempty" type:"Struct"`
	FunctionCompute    *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute    `json:"FunctionCompute,omitempty" xml:"FunctionCompute,omitempty" type:"Struct"`
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation) SetAdditionalFeatures(v *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeatures) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation {
	s.AdditionalFeatures = v
	return s
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation) SetFunctionCompute(v *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformation {
	s.FunctionCompute = v
	return s
}

type ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeatures struct {
	CustomForwardHeaders *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeaturesCustomForwardHeaders `json:"CustomForwardHeaders,omitempty" xml:"CustomForwardHeaders,omitempty" type:"Struct"`
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeatures) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeatures) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeatures) SetCustomForwardHeaders(v *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeaturesCustomForwardHeaders) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeatures {
	s.CustomForwardHeaders = v
	return s
}

type ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeaturesCustomForwardHeaders struct {
	CustomForwardHeader []*string `json:"CustomForwardHeader,omitempty" xml:"CustomForwardHeader,omitempty" type:"Repeated"`
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeaturesCustomForwardHeaders) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeaturesCustomForwardHeaders) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeaturesCustomForwardHeaders) SetCustomForwardHeader(v []*string) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationAdditionalFeaturesCustomForwardHeaders {
	s.CustomForwardHeader = v
	return s
}

type ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute struct {
	FunctionArn           *string `json:"FunctionArn,omitempty" xml:"FunctionArn,omitempty"`
	FunctionAssumeRoleArn *string `json:"FunctionAssumeRoleArn,omitempty" xml:"FunctionAssumeRoleArn,omitempty"`
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute) String() string {
	return tea.Prettify(s)
}

func (s ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute) GoString() string {
	return s.String()
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute) SetFunctionArn(v string) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute {
	s.FunctionArn = &v
	return s
}

func (s *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute) SetFunctionAssumeRoleArn(v string) *ObjectProcessConfigurationTransformationConfigurationsTransformationConfigurationContentTransformationFunctionCompute {
	s.FunctionAssumeRoleArn = &v
	return s
}

type ObjectSummary struct {
	ETag         *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	Key          *string `json:"Key,omitempty" xml:"Key,omitempty"`
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
	Owner        *Owner  `json:"Owner,omitempty" xml:"Owner,omitempty"`
	ResoreInfo   *string `json:"ResoreInfo,omitempty" xml:"ResoreInfo,omitempty"`
	Size         *int64  `json:"Size,omitempty" xml:"Size,omitempty"`
	StorageClass *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
	Type         *string `json:"Type,omitempty" xml:"Type,omitempty"`
}

func (s ObjectSummary) String() string {
	return tea.Prettify(s)
}

func (s ObjectSummary) GoString() string {
	return s.String()
}

func (s *ObjectSummary) SetETag(v string) *ObjectSummary {
	s.ETag = &v
	return s
}

func (s *ObjectSummary) SetKey(v string) *ObjectSummary {
	s.Key = &v
	return s
}

func (s *ObjectSummary) SetLastModified(v string) *ObjectSummary {
	s.LastModified = &v
	return s
}

func (s *ObjectSummary) SetOwner(v *Owner) *ObjectSummary {
	s.Owner = v
	return s
}

func (s *ObjectSummary) SetResoreInfo(v string) *ObjectSummary {
	s.ResoreInfo = &v
	return s
}

func (s *ObjectSummary) SetSize(v int64) *ObjectSummary {
	s.Size = &v
	return s
}

func (s *ObjectSummary) SetStorageClass(v string) *ObjectSummary {
	s.StorageClass = &v
	return s
}

func (s *ObjectSummary) SetType(v string) *ObjectSummary {
	s.Type = &v
	return s
}

type ObjectVersion struct {
	ETag         *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	IsLatest     *bool   `json:"IsLatest,omitempty" xml:"IsLatest,omitempty"`
	Key          *string `json:"Key,omitempty" xml:"Key,omitempty"`
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
	Owner        *Owner  `json:"Owner,omitempty" xml:"Owner,omitempty"`
	Size         *int64  `json:"Size,omitempty" xml:"Size,omitempty"`
	StorageClass *string `json:"StorageClass,omitempty" xml:"StorageClass,omitempty"`
	VersionId    *string `json:"VersionId,omitempty" xml:"VersionId,omitempty"`
}

func (s ObjectVersion) String() string {
	return tea.Prettify(s)
}

func (s ObjectVersion) GoString() string {
	return s.String()
}

func (s *ObjectVersion) SetETag(v string) *ObjectVersion {
	s.ETag = &v
	return s
}

func (s *ObjectVersion) SetIsLatest(v bool) *ObjectVersion {
	s.IsLatest = &v
	return s
}

func (s *ObjectVersion) SetKey(v string) *ObjectVersion {
	s.Key = &v
	return s
}

func (s *ObjectVersion) SetLastModified(v string) *ObjectVersion {
	s.LastModified = &v
	return s
}

func (s *ObjectVersion) SetOwner(v *Owner) *ObjectVersion {
	s.Owner = v
	return s
}

func (s *ObjectVersion) SetSize(v int64) *ObjectVersion {
	s.Size = &v
	return s
}

func (s *ObjectVersion) SetStorageClass(v string) *ObjectVersion {
	s.StorageClass = &v
	return s
}

func (s *ObjectVersion) SetVersionId(v string) *ObjectVersion {
	s.VersionId = &v
	return s
}

type OutputSerialization struct {
	Csv              *CSVOutput  `json:"CSV,omitempty" xml:"CSV,omitempty"`
	EnablePayloadCrc *bool       `json:"EnablePayloadCrc,omitempty" xml:"EnablePayloadCrc,omitempty"`
	Json             *JSONOutput `json:"JSON,omitempty" xml:"JSON,omitempty"`
	KeepAllColumns   *bool       `json:"KeepAllColumns,omitempty" xml:"KeepAllColumns,omitempty"`
	OutputHeader     *bool       `json:"OutputHeader,omitempty" xml:"OutputHeader,omitempty"`
	OutputRawData    *bool       `json:"OutputRawData,omitempty" xml:"OutputRawData,omitempty"`
}

func (s OutputSerialization) String() string {
	return tea.Prettify(s)
}

func (s OutputSerialization) GoString() string {
	return s.String()
}

func (s *OutputSerialization) SetCsv(v *CSVOutput) *OutputSerialization {
	s.Csv = v
	return s
}

func (s *OutputSerialization) SetEnablePayloadCrc(v bool) *OutputSerialization {
	s.EnablePayloadCrc = &v
	return s
}

func (s *OutputSerialization) SetJson(v *JSONOutput) *OutputSerialization {
	s.Json = v
	return s
}

func (s *OutputSerialization) SetKeepAllColumns(v bool) *OutputSerialization {
	s.KeepAllColumns = &v
	return s
}

func (s *OutputSerialization) SetOutputHeader(v bool) *OutputSerialization {
	s.OutputHeader = &v
	return s
}

func (s *OutputSerialization) SetOutputRawData(v bool) *OutputSerialization {
	s.OutputRawData = &v
	return s
}

type Owner struct {
	DisplayName *string `json:"DisplayName,omitempty" xml:"DisplayName,omitempty"`
	ID          *string `json:"ID,omitempty" xml:"ID,omitempty"`
}

func (s Owner) String() string {
	return tea.Prettify(s)
}

func (s Owner) GoString() string {
	return s.String()
}

func (s *Owner) SetDisplayName(v string) *Owner {
	s.DisplayName = &v
	return s
}

func (s *Owner) SetID(v string) *Owner {
	s.ID = &v
	return s
}

type Part struct {
	ETag         *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
	PartNumber   *int64  `json:"PartNumber,omitempty" xml:"PartNumber,omitempty"`
	Size         *int64  `json:"Size,omitempty" xml:"Size,omitempty"`
}

func (s Part) String() string {
	return tea.Prettify(s)
}

func (s Part) GoString() string {
	return s.String()
}

func (s *Part) SetETag(v string) *Part {
	s.ETag = &v
	return s
}

func (s *Part) SetLastModified(v string) *Part {
	s.LastModified = &v
	return s
}

func (s *Part) SetPartNumber(v int64) *Part {
	s.PartNumber = &v
	return s
}

func (s *Part) SetSize(v int64) *Part {
	s.Size = &v
	return s
}

type PublicAccessBlockConfiguration struct {
	BlockPublicAccess *bool `json:"BlockPublicAccess,omitempty" xml:"BlockPublicAccess,omitempty"`
}

func (s PublicAccessBlockConfiguration) String() string {
	return tea.Prettify(s)
}

func (s PublicAccessBlockConfiguration) GoString() string {
	return s.String()
}

func (s *PublicAccessBlockConfiguration) SetBlockPublicAccess(v bool) *PublicAccessBlockConfiguration {
	s.BlockPublicAccess = &v
	return s
}

type PutReplicationRule struct {
	Action                      *string                             `json:"Action,omitempty" xml:"Action,omitempty"`
	Destination                 *ReplicationDestination             `json:"Destination,omitempty" xml:"Destination,omitempty"`
	EncryptionConfiguration     *ReplicationEncryptionConfiguration `json:"EncryptionConfiguration,omitempty" xml:"EncryptionConfiguration,omitempty"`
	HistoricalObjectReplication *string                             `json:"HistoricalObjectReplication,omitempty" xml:"HistoricalObjectReplication,omitempty"`
	ID                          *string                             `json:"ID,omitempty" xml:"ID,omitempty"`
	PrefixSet                   *ReplicationPrefixSet               `json:"PrefixSet,omitempty" xml:"PrefixSet,omitempty"`
	RTC                         *RTC                                `json:"RTC,omitempty" xml:"RTC,omitempty"`
	SourceSelectionCriteria     *ReplicationSourceSelectionCriteria `json:"SourceSelectionCriteria,omitempty" xml:"SourceSelectionCriteria,omitempty"`
	SyncRole                    *string                             `json:"SyncRole,omitempty" xml:"SyncRole,omitempty"`
}

func (s PutReplicationRule) String() string {
	return tea.Prettify(s)
}

func (s PutReplicationRule) GoString() string {
	return s.String()
}

func (s *PutReplicationRule) SetAction(v string) *PutReplicationRule {
	s.Action = &v
	return s
}

func (s *PutReplicationRule) SetDestination(v *ReplicationDestination) *PutReplicationRule {
	s.Destination = v
	return s
}

func (s *PutReplicationRule) SetEncryptionConfiguration(v *ReplicationEncryptionConfiguration) *PutReplicationRule {
	s.EncryptionConfiguration = v
	return s
}

func (s *PutReplicationRule) SetHistoricalObjectReplication(v string) *PutReplicationRule {
	s.HistoricalObjectReplication = &v
	return s
}

func (s *PutReplicationRule) SetID(v string) *PutReplicationRule {
	s.ID = &v
	return s
}

func (s *PutReplicationRule) SetPrefixSet(v *ReplicationPrefixSet) *PutReplicationRule {
	s.PrefixSet = v
	return s
}

func (s *PutReplicationRule) SetRTC(v *RTC) *PutReplicationRule {
	s.RTC = v
	return s
}

func (s *PutReplicationRule) SetSourceSelectionCriteria(v *ReplicationSourceSelectionCriteria) *PutReplicationRule {
	s.SourceSelectionCriteria = v
	return s
}

func (s *PutReplicationRule) SetSyncRole(v string) *PutReplicationRule {
	s.SyncRole = &v
	return s
}

type RTC struct {
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s RTC) String() string {
	return tea.Prettify(s)
}

func (s RTC) GoString() string {
	return s.String()
}

func (s *RTC) SetStatus(v string) *RTC {
	s.Status = &v
	return s
}

type RefererConfiguration struct {
	AllowEmptyReferer        *bool                                 `json:"AllowEmptyReferer,omitempty" xml:"AllowEmptyReferer,omitempty"`
	AllowTruncateQueryString *bool                                 `json:"AllowTruncateQueryString,omitempty" xml:"AllowTruncateQueryString,omitempty"`
	RefererBlacklist         *RefererConfigurationRefererBlacklist `json:"RefererBlacklist,omitempty" xml:"RefererBlacklist,omitempty" type:"Struct"`
	RefererList              *RefererConfigurationRefererList      `json:"RefererList,omitempty" xml:"RefererList,omitempty" type:"Struct"`
	TruncatePath             *bool                                 `json:"TruncatePath,omitempty" xml:"TruncatePath,omitempty"`
}

func (s RefererConfiguration) String() string {
	return tea.Prettify(s)
}

func (s RefererConfiguration) GoString() string {
	return s.String()
}

func (s *RefererConfiguration) SetAllowEmptyReferer(v bool) *RefererConfiguration {
	s.AllowEmptyReferer = &v
	return s
}

func (s *RefererConfiguration) SetAllowTruncateQueryString(v bool) *RefererConfiguration {
	s.AllowTruncateQueryString = &v
	return s
}

func (s *RefererConfiguration) SetRefererBlacklist(v *RefererConfigurationRefererBlacklist) *RefererConfiguration {
	s.RefererBlacklist = v
	return s
}

func (s *RefererConfiguration) SetRefererList(v *RefererConfigurationRefererList) *RefererConfiguration {
	s.RefererList = v
	return s
}

func (s *RefererConfiguration) SetTruncatePath(v bool) *RefererConfiguration {
	s.TruncatePath = &v
	return s
}

type RefererConfigurationRefererBlacklist struct {
	Referer []*string `json:"Referer,omitempty" xml:"Referer,omitempty" type:"Repeated"`
}

func (s RefererConfigurationRefererBlacklist) String() string {
	return tea.Prettify(s)
}

func (s RefererConfigurationRefererBlacklist) GoString() string {
	return s.String()
}

func (s *RefererConfigurationRefererBlacklist) SetReferer(v []*string) *RefererConfigurationRefererBlacklist {
	s.Referer = v
	return s
}

type RefererConfigurationRefererList struct {
	Referer []*string `json:"Referer,omitempty" xml:"Referer,omitempty" type:"Repeated"`
}

func (s RefererConfigurationRefererList) String() string {
	return tea.Prettify(s)
}

func (s RefererConfigurationRefererList) GoString() string {
	return s.String()
}

func (s *RefererConfigurationRefererList) SetReferer(v []*string) *RefererConfigurationRefererList {
	s.Referer = v
	return s
}

type RegionInfo struct {
	AccelerateEndpoint *string `json:"AccelerateEndpoint,omitempty" xml:"AccelerateEndpoint,omitempty"`
	InternalEndpoint   *string `json:"InternalEndpoint,omitempty" xml:"InternalEndpoint,omitempty"`
	InternetEndpoint   *string `json:"InternetEndpoint,omitempty" xml:"InternetEndpoint,omitempty"`
	Region             *string `json:"Region,omitempty" xml:"Region,omitempty"`
}

func (s RegionInfo) String() string {
	return tea.Prettify(s)
}

func (s RegionInfo) GoString() string {
	return s.String()
}

func (s *RegionInfo) SetAccelerateEndpoint(v string) *RegionInfo {
	s.AccelerateEndpoint = &v
	return s
}

func (s *RegionInfo) SetInternalEndpoint(v string) *RegionInfo {
	s.InternalEndpoint = &v
	return s
}

func (s *RegionInfo) SetInternetEndpoint(v string) *RegionInfo {
	s.InternetEndpoint = &v
	return s
}

func (s *RegionInfo) SetRegion(v string) *RegionInfo {
	s.Region = &v
	return s
}

type ReplicationConfiguration struct {
	Rule *PutReplicationRule `json:"Rule,omitempty" xml:"Rule,omitempty"`
}

func (s ReplicationConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ReplicationConfiguration) GoString() string {
	return s.String()
}

func (s *ReplicationConfiguration) SetRule(v *PutReplicationRule) *ReplicationConfiguration {
	s.Rule = v
	return s
}

type ReplicationDestination struct {
	Bucket       *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	Location     *string `json:"Location,omitempty" xml:"Location,omitempty"`
	TransferType *string `json:"TransferType,omitempty" xml:"TransferType,omitempty"`
}

func (s ReplicationDestination) String() string {
	return tea.Prettify(s)
}

func (s ReplicationDestination) GoString() string {
	return s.String()
}

func (s *ReplicationDestination) SetBucket(v string) *ReplicationDestination {
	s.Bucket = &v
	return s
}

func (s *ReplicationDestination) SetLocation(v string) *ReplicationDestination {
	s.Location = &v
	return s
}

func (s *ReplicationDestination) SetTransferType(v string) *ReplicationDestination {
	s.TransferType = &v
	return s
}

type ReplicationEncryptionConfiguration struct {
	ReplicaKmsKeyID *string `json:"ReplicaKmsKeyID,omitempty" xml:"ReplicaKmsKeyID,omitempty"`
}

func (s ReplicationEncryptionConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ReplicationEncryptionConfiguration) GoString() string {
	return s.String()
}

func (s *ReplicationEncryptionConfiguration) SetReplicaKmsKeyID(v string) *ReplicationEncryptionConfiguration {
	s.ReplicaKmsKeyID = &v
	return s
}

type ReplicationPrefixSet struct {
	Prefixs []*string `json:"Prefix,omitempty" xml:"Prefix,omitempty" type:"Repeated"`
}

func (s ReplicationPrefixSet) String() string {
	return tea.Prettify(s)
}

func (s ReplicationPrefixSet) GoString() string {
	return s.String()
}

func (s *ReplicationPrefixSet) SetPrefixs(v []*string) *ReplicationPrefixSet {
	s.Prefixs = v
	return s
}

type ReplicationProgressRule struct {
	Action                      *string                          `json:"Action,omitempty" xml:"Action,omitempty"`
	Destination                 *ReplicationDestination          `json:"Destination,omitempty" xml:"Destination,omitempty"`
	HistoricalObjectReplication *string                          `json:"HistoricalObjectReplication,omitempty" xml:"HistoricalObjectReplication,omitempty"`
	ID                          *string                          `json:"ID,omitempty" xml:"ID,omitempty"`
	PrefixSet                   *ReplicationPrefixSet            `json:"PrefixSet,omitempty" xml:"PrefixSet,omitempty"`
	Progress                    *ReplicationProgressRuleProgress `json:"Progress,omitempty" xml:"Progress,omitempty" type:"Struct"`
	Status                      *string                          `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s ReplicationProgressRule) String() string {
	return tea.Prettify(s)
}

func (s ReplicationProgressRule) GoString() string {
	return s.String()
}

func (s *ReplicationProgressRule) SetAction(v string) *ReplicationProgressRule {
	s.Action = &v
	return s
}

func (s *ReplicationProgressRule) SetDestination(v *ReplicationDestination) *ReplicationProgressRule {
	s.Destination = v
	return s
}

func (s *ReplicationProgressRule) SetHistoricalObjectReplication(v string) *ReplicationProgressRule {
	s.HistoricalObjectReplication = &v
	return s
}

func (s *ReplicationProgressRule) SetID(v string) *ReplicationProgressRule {
	s.ID = &v
	return s
}

func (s *ReplicationProgressRule) SetPrefixSet(v *ReplicationPrefixSet) *ReplicationProgressRule {
	s.PrefixSet = v
	return s
}

func (s *ReplicationProgressRule) SetProgress(v *ReplicationProgressRuleProgress) *ReplicationProgressRule {
	s.Progress = v
	return s
}

func (s *ReplicationProgressRule) SetStatus(v string) *ReplicationProgressRule {
	s.Status = &v
	return s
}

type ReplicationProgressRuleProgress struct {
	HistoricalObject *string `json:"HistoricalObject,omitempty" xml:"HistoricalObject,omitempty"`
	NewObject        *string `json:"NewObject,omitempty" xml:"NewObject,omitempty"`
}

func (s ReplicationProgressRuleProgress) String() string {
	return tea.Prettify(s)
}

func (s ReplicationProgressRuleProgress) GoString() string {
	return s.String()
}

func (s *ReplicationProgressRuleProgress) SetHistoricalObject(v string) *ReplicationProgressRuleProgress {
	s.HistoricalObject = &v
	return s
}

func (s *ReplicationProgressRuleProgress) SetNewObject(v string) *ReplicationProgressRuleProgress {
	s.NewObject = &v
	return s
}

type ReplicationRule struct {
	Action                      *string                             `json:"Action,omitempty" xml:"Action,omitempty"`
	Destination                 *ReplicationDestination             `json:"Destination,omitempty" xml:"Destination,omitempty"`
	EncryptionConfiguration     *ReplicationEncryptionConfiguration `json:"EncryptionConfiguration,omitempty" xml:"EncryptionConfiguration,omitempty"`
	HistoricalObjectReplication *string                             `json:"HistoricalObjectReplication,omitempty" xml:"HistoricalObjectReplication,omitempty"`
	ID                          *string                             `json:"ID,omitempty" xml:"ID,omitempty"`
	PrefixSet                   *ReplicationPrefixSet               `json:"PrefixSet,omitempty" xml:"PrefixSet,omitempty"`
	RTC                         *RTC                                `json:"RTC,omitempty" xml:"RTC,omitempty"`
	SourceSelectionCriteria     *ReplicationSourceSelectionCriteria `json:"SourceSelectionCriteria,omitempty" xml:"SourceSelectionCriteria,omitempty"`
	Status                      *string                             `json:"Status,omitempty" xml:"Status,omitempty"`
	SyncRole                    *string                             `json:"SyncRole,omitempty" xml:"SyncRole,omitempty"`
}

func (s ReplicationRule) String() string {
	return tea.Prettify(s)
}

func (s ReplicationRule) GoString() string {
	return s.String()
}

func (s *ReplicationRule) SetAction(v string) *ReplicationRule {
	s.Action = &v
	return s
}

func (s *ReplicationRule) SetDestination(v *ReplicationDestination) *ReplicationRule {
	s.Destination = v
	return s
}

func (s *ReplicationRule) SetEncryptionConfiguration(v *ReplicationEncryptionConfiguration) *ReplicationRule {
	s.EncryptionConfiguration = v
	return s
}

func (s *ReplicationRule) SetHistoricalObjectReplication(v string) *ReplicationRule {
	s.HistoricalObjectReplication = &v
	return s
}

func (s *ReplicationRule) SetID(v string) *ReplicationRule {
	s.ID = &v
	return s
}

func (s *ReplicationRule) SetPrefixSet(v *ReplicationPrefixSet) *ReplicationRule {
	s.PrefixSet = v
	return s
}

func (s *ReplicationRule) SetRTC(v *RTC) *ReplicationRule {
	s.RTC = v
	return s
}

func (s *ReplicationRule) SetSourceSelectionCriteria(v *ReplicationSourceSelectionCriteria) *ReplicationRule {
	s.SourceSelectionCriteria = v
	return s
}

func (s *ReplicationRule) SetStatus(v string) *ReplicationRule {
	s.Status = &v
	return s
}

func (s *ReplicationRule) SetSyncRole(v string) *ReplicationRule {
	s.SyncRole = &v
	return s
}

type ReplicationRuleProgress struct {
	Action    *string               `json:"Action,omitempty" xml:"Action,omitempty"`
	ID        *string               `json:"ID,omitempty" xml:"ID,omitempty"`
	PrefixSet *ReplicationPrefixSet `json:"PrefixSet,omitempty" xml:"PrefixSet,omitempty"`
}

func (s ReplicationRuleProgress) String() string {
	return tea.Prettify(s)
}

func (s ReplicationRuleProgress) GoString() string {
	return s.String()
}

func (s *ReplicationRuleProgress) SetAction(v string) *ReplicationRuleProgress {
	s.Action = &v
	return s
}

func (s *ReplicationRuleProgress) SetID(v string) *ReplicationRuleProgress {
	s.ID = &v
	return s
}

func (s *ReplicationRuleProgress) SetPrefixSet(v *ReplicationPrefixSet) *ReplicationRuleProgress {
	s.PrefixSet = v
	return s
}

type ReplicationRules struct {
	Ids []*string `json:"ID,omitempty" xml:"ID,omitempty" type:"Repeated"`
}

func (s ReplicationRules) String() string {
	return tea.Prettify(s)
}

func (s ReplicationRules) GoString() string {
	return s.String()
}

func (s *ReplicationRules) SetIds(v []*string) *ReplicationRules {
	s.Ids = v
	return s
}

type ReplicationSourceSelectionCriteria struct {
	SseKmsEncryptedObjects *ReplicationSourceSelectionCriteriaSseKmsEncryptedObjects `json:"SseKmsEncryptedObjects,omitempty" xml:"SseKmsEncryptedObjects,omitempty" type:"Struct"`
}

func (s ReplicationSourceSelectionCriteria) String() string {
	return tea.Prettify(s)
}

func (s ReplicationSourceSelectionCriteria) GoString() string {
	return s.String()
}

func (s *ReplicationSourceSelectionCriteria) SetSseKmsEncryptedObjects(v *ReplicationSourceSelectionCriteriaSseKmsEncryptedObjects) *ReplicationSourceSelectionCriteria {
	s.SseKmsEncryptedObjects = v
	return s
}

type ReplicationSourceSelectionCriteriaSseKmsEncryptedObjects struct {
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s ReplicationSourceSelectionCriteriaSseKmsEncryptedObjects) String() string {
	return tea.Prettify(s)
}

func (s ReplicationSourceSelectionCriteriaSseKmsEncryptedObjects) GoString() string {
	return s.String()
}

func (s *ReplicationSourceSelectionCriteriaSseKmsEncryptedObjects) SetStatus(v string) *ReplicationSourceSelectionCriteriaSseKmsEncryptedObjects {
	s.Status = &v
	return s
}

type RequestPaymentConfiguration struct {
	Payer *string `json:"Payer,omitempty" xml:"Payer,omitempty"`
}

func (s RequestPaymentConfiguration) String() string {
	return tea.Prettify(s)
}

func (s RequestPaymentConfiguration) GoString() string {
	return s.String()
}

func (s *RequestPaymentConfiguration) SetPayer(v string) *RequestPaymentConfiguration {
	s.Payer = &v
	return s
}

type ResponseHeaderConfiguration struct {
	Rule []*ResponseHeaderConfigurationRule `json:"Rule,omitempty" xml:"Rule,omitempty" type:"Repeated"`
}

func (s ResponseHeaderConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeaderConfiguration) GoString() string {
	return s.String()
}

func (s *ResponseHeaderConfiguration) SetRule(v []*ResponseHeaderConfigurationRule) *ResponseHeaderConfiguration {
	s.Rule = v
	return s
}

type ResponseHeaderConfigurationRule struct {
	Filters     *ResponseHeaderConfigurationRuleFilters     `json:"Filters,omitempty" xml:"Filters,omitempty" type:"Struct"`
	HideHeaders *ResponseHeaderConfigurationRuleHideHeaders `json:"HideHeaders,omitempty" xml:"HideHeaders,omitempty" type:"Struct"`
	Name        *string                                     `json:"Name,omitempty" xml:"Name,omitempty"`
}

func (s ResponseHeaderConfigurationRule) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeaderConfigurationRule) GoString() string {
	return s.String()
}

func (s *ResponseHeaderConfigurationRule) SetFilters(v *ResponseHeaderConfigurationRuleFilters) *ResponseHeaderConfigurationRule {
	s.Filters = v
	return s
}

func (s *ResponseHeaderConfigurationRule) SetHideHeaders(v *ResponseHeaderConfigurationRuleHideHeaders) *ResponseHeaderConfigurationRule {
	s.HideHeaders = v
	return s
}

func (s *ResponseHeaderConfigurationRule) SetName(v string) *ResponseHeaderConfigurationRule {
	s.Name = &v
	return s
}

type ResponseHeaderConfigurationRuleFilters struct {
	Operation []*string `json:"Operation,omitempty" xml:"Operation,omitempty" type:"Repeated"`
}

func (s ResponseHeaderConfigurationRuleFilters) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeaderConfigurationRuleFilters) GoString() string {
	return s.String()
}

func (s *ResponseHeaderConfigurationRuleFilters) SetOperation(v []*string) *ResponseHeaderConfigurationRuleFilters {
	s.Operation = v
	return s
}

type ResponseHeaderConfigurationRuleHideHeaders struct {
	Header []*string `json:"Header,omitempty" xml:"Header,omitempty" type:"Repeated"`
}

func (s ResponseHeaderConfigurationRuleHideHeaders) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeaderConfigurationRuleHideHeaders) GoString() string {
	return s.String()
}

func (s *ResponseHeaderConfigurationRuleHideHeaders) SetHeader(v []*string) *ResponseHeaderConfigurationRuleHideHeaders {
	s.Header = v
	return s
}

type RestoreRequest struct {
	Days          *int64                       `json:"Days,omitempty" xml:"Days,omitempty"`
	JobParameters *RestoreRequestJobParameters `json:"JobParameters,omitempty" xml:"JobParameters,omitempty" type:"Struct"`
}

func (s RestoreRequest) String() string {
	return tea.Prettify(s)
}

func (s RestoreRequest) GoString() string {
	return s.String()
}

func (s *RestoreRequest) SetDays(v int64) *RestoreRequest {
	s.Days = &v
	return s
}

func (s *RestoreRequest) SetJobParameters(v *RestoreRequestJobParameters) *RestoreRequest {
	s.JobParameters = v
	return s
}

type RestoreRequestJobParameters struct {
	Tier *string `json:"Tier,omitempty" xml:"Tier,omitempty"`
}

func (s RestoreRequestJobParameters) String() string {
	return tea.Prettify(s)
}

func (s RestoreRequestJobParameters) GoString() string {
	return s.String()
}

func (s *RestoreRequestJobParameters) SetTier(v string) *RestoreRequestJobParameters {
	s.Tier = &v
	return s
}

type RoutingRule struct {
	Condition  *RoutingRuleCondition `json:"Condition,omitempty" xml:"Condition,omitempty"`
	LuaConfig  *RoutingRuleLuaConfig `json:"LuaConfig,omitempty" xml:"LuaConfig,omitempty"`
	Redirect   *RoutingRuleRedirect  `json:"Redirect,omitempty" xml:"Redirect,omitempty"`
	RuleNumber *int64                `json:"RuleNumber,omitempty" xml:"RuleNumber,omitempty"`
}

func (s RoutingRule) String() string {
	return tea.Prettify(s)
}

func (s RoutingRule) GoString() string {
	return s.String()
}

func (s *RoutingRule) SetCondition(v *RoutingRuleCondition) *RoutingRule {
	s.Condition = v
	return s
}

func (s *RoutingRule) SetLuaConfig(v *RoutingRuleLuaConfig) *RoutingRule {
	s.LuaConfig = v
	return s
}

func (s *RoutingRule) SetRedirect(v *RoutingRuleRedirect) *RoutingRule {
	s.Redirect = v
	return s
}

func (s *RoutingRule) SetRuleNumber(v int64) *RoutingRule {
	s.RuleNumber = &v
	return s
}

type RoutingRuleCondition struct {
	HttpErrorCodeReturnedEquals *int64                               `json:"HttpErrorCodeReturnedEquals,omitempty" xml:"HttpErrorCodeReturnedEquals,omitempty"`
	IncludeHeader               []*RoutingRuleConditionIncludeHeader `json:"IncludeHeader,omitempty" xml:"IncludeHeader,omitempty" type:"Repeated"`
	KeyPrefixEquals             *string                              `json:"KeyPrefixEquals,omitempty" xml:"KeyPrefixEquals,omitempty"`
	KeySuffixEquals             *string                              `json:"KeySuffixEquals,omitempty" xml:"KeySuffixEquals,omitempty"`
}

func (s RoutingRuleCondition) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleCondition) GoString() string {
	return s.String()
}

func (s *RoutingRuleCondition) SetHttpErrorCodeReturnedEquals(v int64) *RoutingRuleCondition {
	s.HttpErrorCodeReturnedEquals = &v
	return s
}

func (s *RoutingRuleCondition) SetIncludeHeader(v []*RoutingRuleConditionIncludeHeader) *RoutingRuleCondition {
	s.IncludeHeader = v
	return s
}

func (s *RoutingRuleCondition) SetKeyPrefixEquals(v string) *RoutingRuleCondition {
	s.KeyPrefixEquals = &v
	return s
}

func (s *RoutingRuleCondition) SetKeySuffixEquals(v string) *RoutingRuleCondition {
	s.KeySuffixEquals = &v
	return s
}

type RoutingRuleConditionIncludeHeader struct {
	EndsWith   *string `json:"EndsWith,omitempty" xml:"EndsWith,omitempty"`
	Equals     *string `json:"Equals,omitempty" xml:"Equals,omitempty"`
	Key        *string `json:"Key,omitempty" xml:"Key,omitempty"`
	StartsWith *string `json:"StartsWith,omitempty" xml:"StartsWith,omitempty"`
}

func (s RoutingRuleConditionIncludeHeader) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleConditionIncludeHeader) GoString() string {
	return s.String()
}

func (s *RoutingRuleConditionIncludeHeader) SetEndsWith(v string) *RoutingRuleConditionIncludeHeader {
	s.EndsWith = &v
	return s
}

func (s *RoutingRuleConditionIncludeHeader) SetEquals(v string) *RoutingRuleConditionIncludeHeader {
	s.Equals = &v
	return s
}

func (s *RoutingRuleConditionIncludeHeader) SetKey(v string) *RoutingRuleConditionIncludeHeader {
	s.Key = &v
	return s
}

func (s *RoutingRuleConditionIncludeHeader) SetStartsWith(v string) *RoutingRuleConditionIncludeHeader {
	s.StartsWith = &v
	return s
}

type RoutingRuleLuaConfig struct {
	Script *string `json:"Script,omitempty" xml:"Script,omitempty"`
}

func (s RoutingRuleLuaConfig) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleLuaConfig) GoString() string {
	return s.String()
}

func (s *RoutingRuleLuaConfig) SetScript(v string) *RoutingRuleLuaConfig {
	s.Script = &v
	return s
}

type RoutingRuleRedirect struct {
	EnableReplacePrefix            *bool                                     `json:"EnableReplacePrefix,omitempty" xml:"EnableReplacePrefix,omitempty"`
	HostName                       *string                                   `json:"HostName,omitempty" xml:"HostName,omitempty"`
	HttpRedirectCode               *int64                                    `json:"HttpRedirectCode,omitempty" xml:"HttpRedirectCode,omitempty"`
	MirrorAllowGetImageInfo        *bool                                     `json:"MirrorAllowGetImageInfo,omitempty" xml:"MirrorAllowGetImageInfo,omitempty"`
	MirrorAllowHeadObject          *bool                                     `json:"MirrorAllowHeadObject,omitempty" xml:"MirrorAllowHeadObject,omitempty"`
	MirrorAllowVideoSnapshot       *bool                                     `json:"MirrorAllowVideoSnapshot,omitempty" xml:"MirrorAllowVideoSnapshot,omitempty"`
	MirrorAsyncStatus              *int64                                    `json:"MirrorAsyncStatus,omitempty" xml:"MirrorAsyncStatus,omitempty"`
	MirrorAuth                     *RoutingRuleRedirectMirrorAuth            `json:"MirrorAuth,omitempty" xml:"MirrorAuth,omitempty" type:"Struct"`
	MirrorCheckMd5                 *bool                                     `json:"MirrorCheckMd5,omitempty" xml:"MirrorCheckMd5,omitempty"`
	MirrorDstRegion                *string                                   `json:"MirrorDstRegion,omitempty" xml:"MirrorDstRegion,omitempty"`
	MirrorDstSlaveVpcId            *string                                   `json:"MirrorDstSlaveVpcId,omitempty" xml:"MirrorDstSlaveVpcId,omitempty"`
	MirrorDstVpcId                 *string                                   `json:"MirrorDstVpcId,omitempty" xml:"MirrorDstVpcId,omitempty"`
	MirrorFollowRedirect           *bool                                     `json:"MirrorFollowRedirect,omitempty" xml:"MirrorFollowRedirect,omitempty"`
	MirrorHeaders                  *RoutingRuleRedirectMirrorHeaders         `json:"MirrorHeaders,omitempty" xml:"MirrorHeaders,omitempty" type:"Struct"`
	MirrorIsExpressTunnel          *bool                                     `json:"MirrorIsExpressTunnel,omitempty" xml:"MirrorIsExpressTunnel,omitempty"`
	MirrorMultiAlternates          *RoutingRuleRedirectMirrorMultiAlternates `json:"MirrorMultiAlternates,omitempty" xml:"MirrorMultiAlternates,omitempty" type:"Struct"`
	MirrorPassOriginalSlashes      *bool                                     `json:"MirrorPassOriginalSlashes,omitempty" xml:"MirrorPassOriginalSlashes,omitempty"`
	MirrorPassQueryString          *bool                                     `json:"MirrorPassQueryString,omitempty" xml:"MirrorPassQueryString,omitempty"`
	MirrorProxyPass                *bool                                     `json:"MirrorProxyPass,omitempty" xml:"MirrorProxyPass,omitempty"`
	MirrorReturnHeaders            *RoutingRuleRedirectMirrorReturnHeaders   `json:"MirrorReturnHeaders,omitempty" xml:"MirrorReturnHeaders,omitempty" type:"Struct"`
	MirrorRole                     *string                                   `json:"MirrorRole,omitempty" xml:"MirrorRole,omitempty"`
	MirrorSNI                      *bool                                     `json:"MirrorSNI,omitempty" xml:"MirrorSNI,omitempty"`
	MirrorSaveOssMeta              *bool                                     `json:"MirrorSaveOssMeta,omitempty" xml:"MirrorSaveOssMeta,omitempty"`
	MirrorSwitchAllErrors          *bool                                     `json:"MirrorSwitchAllErrors,omitempty" xml:"MirrorSwitchAllErrors,omitempty"`
	MirrorTaggings                 *RoutingRuleRedirectMirrorTaggings        `json:"MirrorTaggings,omitempty" xml:"MirrorTaggings,omitempty" type:"Struct"`
	MirrorTunnelId                 *string                                   `json:"MirrorTunnelId,omitempty" xml:"MirrorTunnelId,omitempty"`
	MirrorURL                      *string                                   `json:"MirrorURL,omitempty" xml:"MirrorURL,omitempty"`
	MirrorURLProbe                 *string                                   `json:"MirrorURLProbe,omitempty" xml:"MirrorURLProbe,omitempty"`
	MirrorURLSlave                 *string                                   `json:"MirrorURLSlave,omitempty" xml:"MirrorURLSlave,omitempty"`
	MirrorUserLastModified         *bool                                     `json:"MirrorUserLastModified,omitempty" xml:"MirrorUserLastModified,omitempty"`
	MirrorUsingRole                *bool                                     `json:"MirrorUsingRole,omitempty" xml:"MirrorUsingRole,omitempty"`
	PassQueryString                *bool                                     `json:"PassQueryString,omitempty" xml:"PassQueryString,omitempty"`
	Protocol                       *string                                   `json:"Protocol,omitempty" xml:"Protocol,omitempty"`
	RedirectType                   *string                                   `json:"RedirectType,omitempty" xml:"RedirectType,omitempty"`
	ReplaceKeyPrefixWith           *string                                   `json:"ReplaceKeyPrefixWith,omitempty" xml:"ReplaceKeyPrefixWith,omitempty"`
	ReplaceKeyWith                 *string                                   `json:"ReplaceKeyWith,omitempty" xml:"ReplaceKeyWith,omitempty"`
	TransparentMirrorResponseCodes *string                                   `json:"TransparentMirrorResponseCodes,omitempty" xml:"TransparentMirrorResponseCodes,omitempty"`
}

func (s RoutingRuleRedirect) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirect) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirect) SetEnableReplacePrefix(v bool) *RoutingRuleRedirect {
	s.EnableReplacePrefix = &v
	return s
}

func (s *RoutingRuleRedirect) SetHostName(v string) *RoutingRuleRedirect {
	s.HostName = &v
	return s
}

func (s *RoutingRuleRedirect) SetHttpRedirectCode(v int64) *RoutingRuleRedirect {
	s.HttpRedirectCode = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorAllowGetImageInfo(v bool) *RoutingRuleRedirect {
	s.MirrorAllowGetImageInfo = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorAllowHeadObject(v bool) *RoutingRuleRedirect {
	s.MirrorAllowHeadObject = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorAllowVideoSnapshot(v bool) *RoutingRuleRedirect {
	s.MirrorAllowVideoSnapshot = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorAsyncStatus(v int64) *RoutingRuleRedirect {
	s.MirrorAsyncStatus = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorAuth(v *RoutingRuleRedirectMirrorAuth) *RoutingRuleRedirect {
	s.MirrorAuth = v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorCheckMd5(v bool) *RoutingRuleRedirect {
	s.MirrorCheckMd5 = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorDstRegion(v string) *RoutingRuleRedirect {
	s.MirrorDstRegion = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorDstSlaveVpcId(v string) *RoutingRuleRedirect {
	s.MirrorDstSlaveVpcId = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorDstVpcId(v string) *RoutingRuleRedirect {
	s.MirrorDstVpcId = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorFollowRedirect(v bool) *RoutingRuleRedirect {
	s.MirrorFollowRedirect = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorHeaders(v *RoutingRuleRedirectMirrorHeaders) *RoutingRuleRedirect {
	s.MirrorHeaders = v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorIsExpressTunnel(v bool) *RoutingRuleRedirect {
	s.MirrorIsExpressTunnel = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorMultiAlternates(v *RoutingRuleRedirectMirrorMultiAlternates) *RoutingRuleRedirect {
	s.MirrorMultiAlternates = v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorPassOriginalSlashes(v bool) *RoutingRuleRedirect {
	s.MirrorPassOriginalSlashes = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorPassQueryString(v bool) *RoutingRuleRedirect {
	s.MirrorPassQueryString = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorProxyPass(v bool) *RoutingRuleRedirect {
	s.MirrorProxyPass = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorReturnHeaders(v *RoutingRuleRedirectMirrorReturnHeaders) *RoutingRuleRedirect {
	s.MirrorReturnHeaders = v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorRole(v string) *RoutingRuleRedirect {
	s.MirrorRole = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorSNI(v bool) *RoutingRuleRedirect {
	s.MirrorSNI = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorSaveOssMeta(v bool) *RoutingRuleRedirect {
	s.MirrorSaveOssMeta = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorSwitchAllErrors(v bool) *RoutingRuleRedirect {
	s.MirrorSwitchAllErrors = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorTaggings(v *RoutingRuleRedirectMirrorTaggings) *RoutingRuleRedirect {
	s.MirrorTaggings = v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorTunnelId(v string) *RoutingRuleRedirect {
	s.MirrorTunnelId = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorURL(v string) *RoutingRuleRedirect {
	s.MirrorURL = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorURLProbe(v string) *RoutingRuleRedirect {
	s.MirrorURLProbe = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorURLSlave(v string) *RoutingRuleRedirect {
	s.MirrorURLSlave = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorUserLastModified(v bool) *RoutingRuleRedirect {
	s.MirrorUserLastModified = &v
	return s
}

func (s *RoutingRuleRedirect) SetMirrorUsingRole(v bool) *RoutingRuleRedirect {
	s.MirrorUsingRole = &v
	return s
}

func (s *RoutingRuleRedirect) SetPassQueryString(v bool) *RoutingRuleRedirect {
	s.PassQueryString = &v
	return s
}

func (s *RoutingRuleRedirect) SetProtocol(v string) *RoutingRuleRedirect {
	s.Protocol = &v
	return s
}

func (s *RoutingRuleRedirect) SetRedirectType(v string) *RoutingRuleRedirect {
	s.RedirectType = &v
	return s
}

func (s *RoutingRuleRedirect) SetReplaceKeyPrefixWith(v string) *RoutingRuleRedirect {
	s.ReplaceKeyPrefixWith = &v
	return s
}

func (s *RoutingRuleRedirect) SetReplaceKeyWith(v string) *RoutingRuleRedirect {
	s.ReplaceKeyWith = &v
	return s
}

func (s *RoutingRuleRedirect) SetTransparentMirrorResponseCodes(v string) *RoutingRuleRedirect {
	s.TransparentMirrorResponseCodes = &v
	return s
}

type RoutingRuleRedirectMirrorAuth struct {
	AccessKeyId     *string `json:"AccessKeyId,omitempty" xml:"AccessKeyId,omitempty"`
	AccessKeySecret *string `json:"AccessKeySecret,omitempty" xml:"AccessKeySecret,omitempty"`
	AuthType        *string `json:"AuthType,omitempty" xml:"AuthType,omitempty"`
	Region          *string `json:"Region,omitempty" xml:"Region,omitempty"`
}

func (s RoutingRuleRedirectMirrorAuth) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorAuth) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorAuth) SetAccessKeyId(v string) *RoutingRuleRedirectMirrorAuth {
	s.AccessKeyId = &v
	return s
}

func (s *RoutingRuleRedirectMirrorAuth) SetAccessKeySecret(v string) *RoutingRuleRedirectMirrorAuth {
	s.AccessKeySecret = &v
	return s
}

func (s *RoutingRuleRedirectMirrorAuth) SetAuthType(v string) *RoutingRuleRedirectMirrorAuth {
	s.AuthType = &v
	return s
}

func (s *RoutingRuleRedirectMirrorAuth) SetRegion(v string) *RoutingRuleRedirectMirrorAuth {
	s.Region = &v
	return s
}

type RoutingRuleRedirectMirrorHeaders struct {
	Pass    []*string                              `json:"Pass,omitempty" xml:"Pass,omitempty" type:"Repeated"`
	PassAll *bool                                  `json:"PassAll,omitempty" xml:"PassAll,omitempty"`
	Remove  []*string                              `json:"Remove,omitempty" xml:"Remove,omitempty" type:"Repeated"`
	Set     []*RoutingRuleRedirectMirrorHeadersSet `json:"Set,omitempty" xml:"Set,omitempty" type:"Repeated"`
}

func (s RoutingRuleRedirectMirrorHeaders) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorHeaders) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorHeaders) SetPass(v []*string) *RoutingRuleRedirectMirrorHeaders {
	s.Pass = v
	return s
}

func (s *RoutingRuleRedirectMirrorHeaders) SetPassAll(v bool) *RoutingRuleRedirectMirrorHeaders {
	s.PassAll = &v
	return s
}

func (s *RoutingRuleRedirectMirrorHeaders) SetRemove(v []*string) *RoutingRuleRedirectMirrorHeaders {
	s.Remove = v
	return s
}

func (s *RoutingRuleRedirectMirrorHeaders) SetSet(v []*RoutingRuleRedirectMirrorHeadersSet) *RoutingRuleRedirectMirrorHeaders {
	s.Set = v
	return s
}

type RoutingRuleRedirectMirrorHeadersSet struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s RoutingRuleRedirectMirrorHeadersSet) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorHeadersSet) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorHeadersSet) SetKey(v string) *RoutingRuleRedirectMirrorHeadersSet {
	s.Key = &v
	return s
}

func (s *RoutingRuleRedirectMirrorHeadersSet) SetValue(v string) *RoutingRuleRedirectMirrorHeadersSet {
	s.Value = &v
	return s
}

type RoutingRuleRedirectMirrorMultiAlternates struct {
	MirrorMultiAlternate []*RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate `json:"MirrorMultiAlternate,omitempty" xml:"MirrorMultiAlternate,omitempty" type:"Repeated"`
}

func (s RoutingRuleRedirectMirrorMultiAlternates) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorMultiAlternates) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorMultiAlternates) SetMirrorMultiAlternate(v []*RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate) *RoutingRuleRedirectMirrorMultiAlternates {
	s.MirrorMultiAlternate = v
	return s
}

type RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate struct {
	MirrorMultiAlternateDstRegion *string `json:"MirrorMultiAlternateDstRegion,omitempty" xml:"MirrorMultiAlternateDstRegion,omitempty"`
	MirrorMultiAlternateNumber    *int64  `json:"MirrorMultiAlternateNumber,omitempty" xml:"MirrorMultiAlternateNumber,omitempty"`
	MirrorMultiAlternateURL       *string `json:"MirrorMultiAlternateURL,omitempty" xml:"MirrorMultiAlternateURL,omitempty"`
	MirrorMultiAlternateVpcId     *string `json:"MirrorMultiAlternateVpcId,omitempty" xml:"MirrorMultiAlternateVpcId,omitempty"`
}

func (s RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate) SetMirrorMultiAlternateDstRegion(v string) *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate {
	s.MirrorMultiAlternateDstRegion = &v
	return s
}

func (s *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate) SetMirrorMultiAlternateNumber(v int64) *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate {
	s.MirrorMultiAlternateNumber = &v
	return s
}

func (s *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate) SetMirrorMultiAlternateURL(v string) *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate {
	s.MirrorMultiAlternateURL = &v
	return s
}

func (s *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate) SetMirrorMultiAlternateVpcId(v string) *RoutingRuleRedirectMirrorMultiAlternatesMirrorMultiAlternate {
	s.MirrorMultiAlternateVpcId = &v
	return s
}

type RoutingRuleRedirectMirrorReturnHeaders struct {
	ReturnHeader []*RoutingRuleRedirectMirrorReturnHeadersReturnHeader `json:"ReturnHeader,omitempty" xml:"ReturnHeader,omitempty" type:"Repeated"`
}

func (s RoutingRuleRedirectMirrorReturnHeaders) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorReturnHeaders) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorReturnHeaders) SetReturnHeader(v []*RoutingRuleRedirectMirrorReturnHeadersReturnHeader) *RoutingRuleRedirectMirrorReturnHeaders {
	s.ReturnHeader = v
	return s
}

type RoutingRuleRedirectMirrorReturnHeadersReturnHeader struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s RoutingRuleRedirectMirrorReturnHeadersReturnHeader) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorReturnHeadersReturnHeader) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorReturnHeadersReturnHeader) SetKey(v string) *RoutingRuleRedirectMirrorReturnHeadersReturnHeader {
	s.Key = &v
	return s
}

func (s *RoutingRuleRedirectMirrorReturnHeadersReturnHeader) SetValue(v string) *RoutingRuleRedirectMirrorReturnHeadersReturnHeader {
	s.Value = &v
	return s
}

type RoutingRuleRedirectMirrorTaggings struct {
	Taggings []*RoutingRuleRedirectMirrorTaggingsTaggings `json:"Taggings,omitempty" xml:"Taggings,omitempty" type:"Repeated"`
}

func (s RoutingRuleRedirectMirrorTaggings) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorTaggings) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorTaggings) SetTaggings(v []*RoutingRuleRedirectMirrorTaggingsTaggings) *RoutingRuleRedirectMirrorTaggings {
	s.Taggings = v
	return s
}

type RoutingRuleRedirectMirrorTaggingsTaggings struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s RoutingRuleRedirectMirrorTaggingsTaggings) String() string {
	return tea.Prettify(s)
}

func (s RoutingRuleRedirectMirrorTaggingsTaggings) GoString() string {
	return s.String()
}

func (s *RoutingRuleRedirectMirrorTaggingsTaggings) SetKey(v string) *RoutingRuleRedirectMirrorTaggingsTaggings {
	s.Key = &v
	return s
}

func (s *RoutingRuleRedirectMirrorTaggingsTaggings) SetValue(v string) *RoutingRuleRedirectMirrorTaggingsTaggings {
	s.Value = &v
	return s
}

type RtcConfiguration struct {
	ID  *string `json:"ID,omitempty" xml:"ID,omitempty"`
	RTC *RTC    `json:"RTC,omitempty" xml:"RTC,omitempty"`
}

func (s RtcConfiguration) String() string {
	return tea.Prettify(s)
}

func (s RtcConfiguration) GoString() string {
	return s.String()
}

func (s *RtcConfiguration) SetID(v string) *RtcConfiguration {
	s.ID = &v
	return s
}

func (s *RtcConfiguration) SetRTC(v *RTC) *RtcConfiguration {
	s.RTC = v
	return s
}

type SSEKMS struct {
	KeyId *string `json:"KeyId,omitempty" xml:"KeyId,omitempty"`
}

func (s SSEKMS) String() string {
	return tea.Prettify(s)
}

func (s SSEKMS) GoString() string {
	return s.String()
}

func (s *SSEKMS) SetKeyId(v string) *SSEKMS {
	s.KeyId = &v
	return s
}

type SelectMetaRequest struct {
	InputSerialization *InputSerialization `json:"InputSerialization,omitempty" xml:"InputSerialization,omitempty"`
	OverwriteIfExists  *bool               `json:"OverwriteIfExists,omitempty" xml:"OverwriteIfExists,omitempty"`
}

func (s SelectMetaRequest) String() string {
	return tea.Prettify(s)
}

func (s SelectMetaRequest) GoString() string {
	return s.String()
}

func (s *SelectMetaRequest) SetInputSerialization(v *InputSerialization) *SelectMetaRequest {
	s.InputSerialization = v
	return s
}

func (s *SelectMetaRequest) SetOverwriteIfExists(v bool) *SelectMetaRequest {
	s.OverwriteIfExists = &v
	return s
}

type SelectMetaStatus struct {
	ColsCount         *int64  `json:"ColsCount,omitempty" xml:"ColsCount,omitempty"`
	ErrorMessage      *string `json:"ErrorMessage,omitempty" xml:"ErrorMessage,omitempty"`
	Offset            *int64  `json:"Offset,omitempty" xml:"Offset,omitempty"`
	RowsCount         *int64  `json:"RowsCount,omitempty" xml:"RowsCount,omitempty"`
	SplitsCount       *int64  `json:"SplitsCount,omitempty" xml:"SplitsCount,omitempty"`
	Status            *int64  `json:"Status,omitempty" xml:"Status,omitempty"`
	TotalScannedBytes *int64  `json:"TotalScannedBytes,omitempty" xml:"TotalScannedBytes,omitempty"`
}

func (s SelectMetaStatus) String() string {
	return tea.Prettify(s)
}

func (s SelectMetaStatus) GoString() string {
	return s.String()
}

func (s *SelectMetaStatus) SetColsCount(v int64) *SelectMetaStatus {
	s.ColsCount = &v
	return s
}

func (s *SelectMetaStatus) SetErrorMessage(v string) *SelectMetaStatus {
	s.ErrorMessage = &v
	return s
}

func (s *SelectMetaStatus) SetOffset(v int64) *SelectMetaStatus {
	s.Offset = &v
	return s
}

func (s *SelectMetaStatus) SetRowsCount(v int64) *SelectMetaStatus {
	s.RowsCount = &v
	return s
}

func (s *SelectMetaStatus) SetSplitsCount(v int64) *SelectMetaStatus {
	s.SplitsCount = &v
	return s
}

func (s *SelectMetaStatus) SetStatus(v int64) *SelectMetaStatus {
	s.Status = &v
	return s
}

func (s *SelectMetaStatus) SetTotalScannedBytes(v int64) *SelectMetaStatus {
	s.TotalScannedBytes = &v
	return s
}

type SelectRequest struct {
	Expression          *string               `json:"Expression,omitempty" xml:"Expression,omitempty"`
	InputSerialization  *InputSerialization   `json:"InputSerialization,omitempty" xml:"InputSerialization,omitempty"`
	Options             *SelectRequestOptions `json:"Options,omitempty" xml:"Options,omitempty"`
	OutputSerialization *OutputSerialization  `json:"OutputSerialization,omitempty" xml:"OutputSerialization,omitempty"`
}

func (s SelectRequest) String() string {
	return tea.Prettify(s)
}

func (s SelectRequest) GoString() string {
	return s.String()
}

func (s *SelectRequest) SetExpression(v string) *SelectRequest {
	s.Expression = &v
	return s
}

func (s *SelectRequest) SetInputSerialization(v *InputSerialization) *SelectRequest {
	s.InputSerialization = v
	return s
}

func (s *SelectRequest) SetOptions(v *SelectRequestOptions) *SelectRequest {
	s.Options = v
	return s
}

func (s *SelectRequest) SetOutputSerialization(v *OutputSerialization) *SelectRequest {
	s.OutputSerialization = v
	return s
}

type SelectRequestOptions struct {
	MaxSkippedRecordsAllowed *int64 `json:"MaxSkippedRecordsAllowed,omitempty" xml:"MaxSkippedRecordsAllowed,omitempty"`
	SkipPartialDataRecord    *bool  `json:"SkipPartialDataRecord,omitempty" xml:"SkipPartialDataRecord,omitempty"`
}

func (s SelectRequestOptions) String() string {
	return tea.Prettify(s)
}

func (s SelectRequestOptions) GoString() string {
	return s.String()
}

func (s *SelectRequestOptions) SetMaxSkippedRecordsAllowed(v int64) *SelectRequestOptions {
	s.MaxSkippedRecordsAllowed = &v
	return s
}

func (s *SelectRequestOptions) SetSkipPartialDataRecord(v bool) *SelectRequestOptions {
	s.SkipPartialDataRecord = &v
	return s
}

type ServerSideEncryptionRule struct {
	ApplyServerSideEncryptionByDefault *ApplyServerSideEncryptionByDefault `json:"ApplyServerSideEncryptionByDefault,omitempty" xml:"ApplyServerSideEncryptionByDefault,omitempty"`
}

func (s ServerSideEncryptionRule) String() string {
	return tea.Prettify(s)
}

func (s ServerSideEncryptionRule) GoString() string {
	return s.String()
}

func (s *ServerSideEncryptionRule) SetApplyServerSideEncryptionByDefault(v *ApplyServerSideEncryptionByDefault) *ServerSideEncryptionRule {
	s.ApplyServerSideEncryptionByDefault = v
	return s
}

type Style struct {
	Content *string `json:"Content,omitempty" xml:"Content,omitempty"`
}

func (s Style) String() string {
	return tea.Prettify(s)
}

func (s Style) GoString() string {
	return s.String()
}

func (s *Style) SetContent(v string) *Style {
	s.Content = &v
	return s
}

type StyleInfo struct {
	Category       *string `json:"Category,omitempty" xml:"Category,omitempty"`
	Content        *string `json:"Content,omitempty" xml:"Content,omitempty"`
	CreateTime     *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	LastModifyTime *string `json:"LastModifyTime,omitempty" xml:"LastModifyTime,omitempty"`
	Name           *string `json:"Name,omitempty" xml:"Name,omitempty"`
}

func (s StyleInfo) String() string {
	return tea.Prettify(s)
}

func (s StyleInfo) GoString() string {
	return s.String()
}

func (s *StyleInfo) SetCategory(v string) *StyleInfo {
	s.Category = &v
	return s
}

func (s *StyleInfo) SetContent(v string) *StyleInfo {
	s.Content = &v
	return s
}

func (s *StyleInfo) SetCreateTime(v string) *StyleInfo {
	s.CreateTime = &v
	return s
}

func (s *StyleInfo) SetLastModifyTime(v string) *StyleInfo {
	s.LastModifyTime = &v
	return s
}

func (s *StyleInfo) SetName(v string) *StyleInfo {
	s.Name = &v
	return s
}

type Tag struct {
	Key   *string `json:"Key,omitempty" xml:"Key,omitempty"`
	Value *string `json:"Value,omitempty" xml:"Value,omitempty"`
}

func (s Tag) String() string {
	return tea.Prettify(s)
}

func (s Tag) GoString() string {
	return s.String()
}

func (s *Tag) SetKey(v string) *Tag {
	s.Key = &v
	return s
}

func (s *Tag) SetValue(v string) *Tag {
	s.Value = &v
	return s
}

type TagSet struct {
	Tags []*Tag `json:"Tag,omitempty" xml:"Tag,omitempty" type:"Repeated"`
}

func (s TagSet) String() string {
	return tea.Prettify(s)
}

func (s TagSet) GoString() string {
	return s.String()
}

func (s *TagSet) SetTags(v []*Tag) *TagSet {
	s.Tags = v
	return s
}

type Tagging struct {
	TagSet *TagSet `json:"TagSet,omitempty" xml:"TagSet,omitempty"`
}

func (s Tagging) String() string {
	return tea.Prettify(s)
}

func (s Tagging) GoString() string {
	return s.String()
}

func (s *Tagging) SetTagSet(v *TagSet) *Tagging {
	s.TagSet = v
	return s
}

type TransferAccelerationConfiguration struct {
	Enabled *bool `json:"Enabled,omitempty" xml:"Enabled,omitempty"`
}

func (s TransferAccelerationConfiguration) String() string {
	return tea.Prettify(s)
}

func (s TransferAccelerationConfiguration) GoString() string {
	return s.String()
}

func (s *TransferAccelerationConfiguration) SetEnabled(v bool) *TransferAccelerationConfiguration {
	s.Enabled = &v
	return s
}

type Upload struct {
	Initiated *string `json:"Initiated,omitempty" xml:"Initiated,omitempty"`
	Key       *string `json:"Key,omitempty" xml:"Key,omitempty"`
	UploadId  *string `json:"UploadId,omitempty" xml:"UploadId,omitempty"`
}

func (s Upload) String() string {
	return tea.Prettify(s)
}

func (s Upload) GoString() string {
	return s.String()
}

func (s *Upload) SetInitiated(v string) *Upload {
	s.Initiated = &v
	return s
}

func (s *Upload) SetKey(v string) *Upload {
	s.Key = &v
	return s
}

func (s *Upload) SetUploadId(v string) *Upload {
	s.UploadId = &v
	return s
}

type UserAntiDDOSInfo struct {
	ActiveTime *int64  `json:"ActiveTime,omitempty" xml:"ActiveTime,omitempty"`
	Ctime      *int64  `json:"Ctime,omitempty" xml:"Ctime,omitempty"`
	InstanceId *string `json:"InstanceId,omitempty" xml:"InstanceId,omitempty"`
	Mtime      *int64  `json:"Mtime,omitempty" xml:"Mtime,omitempty"`
	Owner      *string `json:"Owner,omitempty" xml:"Owner,omitempty"`
	Status     *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s UserAntiDDOSInfo) String() string {
	return tea.Prettify(s)
}

func (s UserAntiDDOSInfo) GoString() string {
	return s.String()
}

func (s *UserAntiDDOSInfo) SetActiveTime(v int64) *UserAntiDDOSInfo {
	s.ActiveTime = &v
	return s
}

func (s *UserAntiDDOSInfo) SetCtime(v int64) *UserAntiDDOSInfo {
	s.Ctime = &v
	return s
}

func (s *UserAntiDDOSInfo) SetInstanceId(v string) *UserAntiDDOSInfo {
	s.InstanceId = &v
	return s
}

func (s *UserAntiDDOSInfo) SetMtime(v int64) *UserAntiDDOSInfo {
	s.Mtime = &v
	return s
}

func (s *UserAntiDDOSInfo) SetOwner(v string) *UserAntiDDOSInfo {
	s.Owner = &v
	return s
}

func (s *UserAntiDDOSInfo) SetStatus(v string) *UserAntiDDOSInfo {
	s.Status = &v
	return s
}

type UserDefinedLogFieldsConfiguration struct {
	HeaderSet *UserDefinedLogFieldsConfigurationHeaderSet `json:"HeaderSet,omitempty" xml:"HeaderSet,omitempty" type:"Struct"`
	ParamSet  *UserDefinedLogFieldsConfigurationParamSet  `json:"ParamSet,omitempty" xml:"ParamSet,omitempty" type:"Struct"`
}

func (s UserDefinedLogFieldsConfiguration) String() string {
	return tea.Prettify(s)
}

func (s UserDefinedLogFieldsConfiguration) GoString() string {
	return s.String()
}

func (s *UserDefinedLogFieldsConfiguration) SetHeaderSet(v *UserDefinedLogFieldsConfigurationHeaderSet) *UserDefinedLogFieldsConfiguration {
	s.HeaderSet = v
	return s
}

func (s *UserDefinedLogFieldsConfiguration) SetParamSet(v *UserDefinedLogFieldsConfigurationParamSet) *UserDefinedLogFieldsConfiguration {
	s.ParamSet = v
	return s
}

type UserDefinedLogFieldsConfigurationHeaderSet struct {
	Header []*string `json:"header,omitempty" xml:"header,omitempty" type:"Repeated"`
}

func (s UserDefinedLogFieldsConfigurationHeaderSet) String() string {
	return tea.Prettify(s)
}

func (s UserDefinedLogFieldsConfigurationHeaderSet) GoString() string {
	return s.String()
}

func (s *UserDefinedLogFieldsConfigurationHeaderSet) SetHeader(v []*string) *UserDefinedLogFieldsConfigurationHeaderSet {
	s.Header = v
	return s
}

type UserDefinedLogFieldsConfigurationParamSet struct {
	Parameter []*string `json:"parameter,omitempty" xml:"parameter,omitempty" type:"Repeated"`
}

func (s UserDefinedLogFieldsConfigurationParamSet) String() string {
	return tea.Prettify(s)
}

func (s UserDefinedLogFieldsConfigurationParamSet) GoString() string {
	return s.String()
}

func (s *UserDefinedLogFieldsConfigurationParamSet) SetParameter(v []*string) *UserDefinedLogFieldsConfigurationParamSet {
	s.Parameter = v
	return s
}

type VersioningConfiguration struct {
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s VersioningConfiguration) String() string {
	return tea.Prettify(s)
}

func (s VersioningConfiguration) GoString() string {
	return s.String()
}

func (s *VersioningConfiguration) SetStatus(v string) *VersioningConfiguration {
	s.Status = &v
	return s
}

type WebsiteConfiguration struct {
	ErrorDocument *ErrorDocument                    `json:"ErrorDocument,omitempty" xml:"ErrorDocument,omitempty"`
	IndexDocument *IndexDocument                    `json:"IndexDocument,omitempty" xml:"IndexDocument,omitempty"`
	RoutingRules  *WebsiteConfigurationRoutingRules `json:"RoutingRules,omitempty" xml:"RoutingRules,omitempty" type:"Struct"`
}

func (s WebsiteConfiguration) String() string {
	return tea.Prettify(s)
}

func (s WebsiteConfiguration) GoString() string {
	return s.String()
}

func (s *WebsiteConfiguration) SetErrorDocument(v *ErrorDocument) *WebsiteConfiguration {
	s.ErrorDocument = v
	return s
}

func (s *WebsiteConfiguration) SetIndexDocument(v *IndexDocument) *WebsiteConfiguration {
	s.IndexDocument = v
	return s
}

func (s *WebsiteConfiguration) SetRoutingRules(v *WebsiteConfigurationRoutingRules) *WebsiteConfiguration {
	s.RoutingRules = v
	return s
}

type WebsiteConfigurationRoutingRules struct {
	RoutingRule []*RoutingRule `json:"RoutingRule,omitempty" xml:"RoutingRule,omitempty" type:"Repeated"`
}

func (s WebsiteConfigurationRoutingRules) String() string {
	return tea.Prettify(s)
}

func (s WebsiteConfigurationRoutingRules) GoString() string {
	return s.String()
}

func (s *WebsiteConfigurationRoutingRules) SetRoutingRule(v []*RoutingRule) *WebsiteConfigurationRoutingRules {
	s.RoutingRule = v
	return s
}

type AbortBucketWormResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s AbortBucketWormResponse) String() string {
	return tea.Prettify(s)
}

func (s AbortBucketWormResponse) GoString() string {
	return s.String()
}

func (s *AbortBucketWormResponse) SetHeaders(v map[string]*string) *AbortBucketWormResponse {
	s.Headers = v
	return s
}

func (s *AbortBucketWormResponse) SetStatusCode(v int32) *AbortBucketWormResponse {
	s.StatusCode = &v
	return s
}

type AbortMultipartUploadRequest struct {
	// The ID of the multipart upload task.
	UploadId *string `json:"uploadId,omitempty" xml:"uploadId,omitempty"`
}

func (s AbortMultipartUploadRequest) String() string {
	return tea.Prettify(s)
}

func (s AbortMultipartUploadRequest) GoString() string {
	return s.String()
}

func (s *AbortMultipartUploadRequest) SetUploadId(v string) *AbortMultipartUploadRequest {
	s.UploadId = &v
	return s
}

type AbortMultipartUploadResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s AbortMultipartUploadResponse) String() string {
	return tea.Prettify(s)
}

func (s AbortMultipartUploadResponse) GoString() string {
	return s.String()
}

func (s *AbortMultipartUploadResponse) SetHeaders(v map[string]*string) *AbortMultipartUploadResponse {
	s.Headers = v
	return s
}

func (s *AbortMultipartUploadResponse) SetStatusCode(v int32) *AbortMultipartUploadResponse {
	s.StatusCode = &v
	return s
}

type AppendObjectHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The web page caching behavior for the object. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	CacheControl *string `json:"Cache-Control,omitempty" xml:"Cache-Control,omitempty"`
	// The name of the object when the object is downloaded. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	ContentDisposition *string `json:"Content-Disposition,omitempty" xml:"Content-Disposition,omitempty"`
	// The encoding format of the object content. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	ContentEncoding *string `json:"Content-Encoding,omitempty" xml:"Content-Encoding,omitempty"`
	// The Content-MD5 header value is a string calculated by using the MD5 algorithm. The header is used to check whether the content of the received message is the same as that of the sent message.
	// To obtain the value of the Content-MD5 header, calculate a 128-bit number based on the message content except for the header, and then encode the number in Base64.
	// Default value: null.
	// Limits: none.
	ContentMD5 *string `json:"Content-MD5,omitempty" xml:"Content-MD5,omitempty"`
	// The expiration time. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	Expires *string `json:"Expires,omitempty" xml:"Expires,omitempty"`
	// You can add parameters whose names are prefixed with x-oss-meta-* when you call the AppendObject operation. These parameters cannot be included in the requests when you append objects to an existing object. Parameters whose names are prefixed with x-oss-meta-* are considered the metadata of the object.
	// You can configure multiple parameters whose name are prefixed with x-oss-meta- for an object. However, the total size of user metadata cannot exceed 8 KB.
	// The name of parameters whose name are prefixed with x-oss-meta- can contain hyphens (-), numbers, and letters. Uppercase letters are converted to lowercase letters. Other characters such as underscores (_) are not supported.
	MetaData map[string]*string `json:"x-oss-meta-*,omitempty" xml:"x-oss-meta-*,omitempty"`
	// The access control list (ACL) of the object. Default value: default.  Valid values:
	//
	// - default: The ACL of the object is the same as that of the bucket in which the object is stored.
	// - private: The ACL of the object is private. Only the owner of the object and authorized users can read and write this object.
	// - public-read: The ACL of the object is public-read. Only the owner of the object and authorized users can read and write this object. Other users can only read the object. Exercise caution when you set the object ACL to this value.
	// - public-read-write: The ACL of the object is public-read-write. All users can read and write this object. Exercise caution when you set the object ACL to this value.
	//
	// For more information about the ACL, see [ACL](~~100676~~).
	Acl *string `json:"x-oss-object-acl,omitempty" xml:"x-oss-object-acl,omitempty"`
	// The method used to encrypt objects on the specified OSS server.
	// Valid values:
	//
	// - AES256: Keys managed by OSS are used for encryption and decryption (SSE-OSS).
	// - KMS: Keys managed by Key Management Service (KMS) are used for encryption and decryption.
	// - SM4: The SM4 block cipher algorithm is used for encryption and decryption.
	ServerSideEncryption *string `json:"x-oss-server-side-encryption,omitempty" xml:"x-oss-server-side-encryption,omitempty"`
	// The storage class of the object that you want to upload. Valid values:
	//
	// - Standard
	// - IA
	// - Archive
	// If you specify the object storage class when you upload an object, the storage class of the uploaded object is the specified value regardless of the storage class of the bucket to which the object is uploaded. If you set x-oss-storage-class to Standard when you upload an object to an IA bucket, the object is stored as a Standard object.
	// For more information about storage classes, see the "Overview" topic in Developer Guide.
	//
	// ><notice> The value that you specify takes effect only when you call the AppendObject operation on an object for the first time.
	StorageClass *string `json:"x-oss-storage-class,omitempty" xml:"x-oss-storage-class,omitempty"`
}

func (s AppendObjectHeaders) String() string {
	return tea.Prettify(s)
}

func (s AppendObjectHeaders) GoString() string {
	return s.String()
}

func (s *AppendObjectHeaders) SetCommonHeaders(v map[string]*string) *AppendObjectHeaders {
	s.CommonHeaders = v
	return s
}

func (s *AppendObjectHeaders) SetCacheControl(v string) *AppendObjectHeaders {
	s.CacheControl = &v
	return s
}

func (s *AppendObjectHeaders) SetContentDisposition(v string) *AppendObjectHeaders {
	s.ContentDisposition = &v
	return s
}

func (s *AppendObjectHeaders) SetContentEncoding(v string) *AppendObjectHeaders {
	s.ContentEncoding = &v
	return s
}

func (s *AppendObjectHeaders) SetContentMD5(v string) *AppendObjectHeaders {
	s.ContentMD5 = &v
	return s
}

func (s *AppendObjectHeaders) SetExpires(v string) *AppendObjectHeaders {
	s.Expires = &v
	return s
}

func (s *AppendObjectHeaders) SetMetaData(v map[string]*string) *AppendObjectHeaders {
	s.MetaData = v
	return s
}

func (s *AppendObjectHeaders) SetAcl(v string) *AppendObjectHeaders {
	s.Acl = &v
	return s
}

func (s *AppendObjectHeaders) SetServerSideEncryption(v string) *AppendObjectHeaders {
	s.ServerSideEncryption = &v
	return s
}

func (s *AppendObjectHeaders) SetStorageClass(v string) *AppendObjectHeaders {
	s.StorageClass = &v
	return s
}

type AppendObjectRequest struct {
	// The request body.
	Body io.Reader `json:"body,omitempty" xml:"body,omitempty"`
	// The position from which the AppendObject operation starts.  Each time an AppendObject operation succeeds, the x-oss-next-append-position header is included in the response to specify the position from which the next AppendObject operation starts. The value of position in the first AppendObject operation performed on an object must be 0. The value of position in subsequent AppendObject operations performed on the object is the current length of the object. For example, if the value of position specified in the first AppendObject request is 0 and the value of content-length is 65536, the value of position in the second AppendObject request must be 65536.
	//
	// - If the value of position in the AppendObject request is 0 and the name of the object that you want to append is unique, you can set headers such as x-oss-server-side-encryption in an AppendObject request in the same way as you set in a PutObject request. If you add the x-oss-server-side-encryption header to an AppendObject request, the x-oss-server-side-encryption header is included in the response to the request. If you want to modify metadata, you can call the CopyObject operation.
	// - If you call an AppendObject operation to append a 0 KB object whose position value is valid to an Appendable object, the status of the Appendable object is not changed.
	Position *int64 `json:"position,omitempty" xml:"position,omitempty"`
}

func (s AppendObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s AppendObjectRequest) GoString() string {
	return s.String()
}

func (s *AppendObjectRequest) SetBody(v io.Reader) *AppendObjectRequest {
	s.Body = v
	return s
}

func (s *AppendObjectRequest) SetPosition(v int64) *AppendObjectRequest {
	s.Position = &v
	return s
}

type AppendObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s AppendObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s AppendObjectResponse) GoString() string {
	return s.String()
}

func (s *AppendObjectResponse) SetHeaders(v map[string]*string) *AppendObjectResponse {
	s.Headers = v
	return s
}

func (s *AppendObjectResponse) SetStatusCode(v int32) *AppendObjectResponse {
	s.StatusCode = &v
	return s
}

type CloseMetaQueryResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s CloseMetaQueryResponse) String() string {
	return tea.Prettify(s)
}

func (s CloseMetaQueryResponse) GoString() string {
	return s.String()
}

func (s *CloseMetaQueryResponse) SetHeaders(v map[string]*string) *CloseMetaQueryResponse {
	s.Headers = v
	return s
}

func (s *CloseMetaQueryResponse) SetStatusCode(v int32) *CloseMetaQueryResponse {
	s.StatusCode = &v
	return s
}

type CompleteBucketWormRequest struct {
	// The ID of the retention policy.
	WormId *string `json:"wormId,omitempty" xml:"wormId,omitempty"`
}

func (s CompleteBucketWormRequest) String() string {
	return tea.Prettify(s)
}

func (s CompleteBucketWormRequest) GoString() string {
	return s.String()
}

func (s *CompleteBucketWormRequest) SetWormId(v string) *CompleteBucketWormRequest {
	s.WormId = &v
	return s
}

type CompleteBucketWormResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s CompleteBucketWormResponse) String() string {
	return tea.Prettify(s)
}

func (s CompleteBucketWormResponse) GoString() string {
	return s.String()
}

func (s *CompleteBucketWormResponse) SetHeaders(v map[string]*string) *CompleteBucketWormResponse {
	s.Headers = v
	return s
}

func (s *CompleteBucketWormResponse) SetStatusCode(v int32) *CompleteBucketWormResponse {
	s.StatusCode = &v
	return s
}

type CompleteMultipartUploadHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// Specifies whether to list all parts that are uploaded by using the current upload ID.
	//
	// Valid value: yes.
	//
	// - If x-oss-complete-all is set to yes in the request, OSS lists all parts that are uploaded by using the current upload ID, sorts the parts by part number, and then performs the CompleteMultipartUpload operation. When OSS performs the CompleteMultipartUpload operation, OSS cannot detect the parts that are not uploaded or currently being uploaded. Before you call the CompleteMultipartUpload operation, make sure that all parts are uploaded.
	//
	// - If x-oss-complete-all is specified in the request, the request body cannot be specified. Otherwise, an error occurs.
	//
	// - If x-oss-complete-all is specified in the request, the format of the response remains unchanged.
	CompleteAll *string `json:"x-oss-complete-all,omitempty" xml:"x-oss-complete-all,omitempty"`
	// Specifieswhethertheobjectwith the sameobjectname is overwritten when you call the CompleteMultipartUpload operation.
	//
	// - If the value of x-oss-forbid-overwrite is not specified or set to false, the existing object can be overwritten by the object that has the same name.
	// - If the value of x-oss-forbid-overwrite is set to true, the existing object cannot be overwritten by the object that has the same name.
	//
	// - The x-oss-forbid-overwrite request header is invalid if versioning is enabled or suspended for the bucket. In this case, the existing object can be overwritten by the object that has the same name when you call the CompleteMultipartUpload operation.
	// - If you specify the x-oss-forbid-overwrite request header, the queries per second (QPS) performance of OSS may be degraded. If you want to configure the x-oss-forbid-overwrite header in a large number of requests (QPS > 1,000), submit a ticket.
	ForbidOverwrite *string `json:"x-oss-forbid-overwrite,omitempty" xml:"x-oss-forbid-overwrite,omitempty"`
}

func (s CompleteMultipartUploadHeaders) String() string {
	return tea.Prettify(s)
}

func (s CompleteMultipartUploadHeaders) GoString() string {
	return s.String()
}

func (s *CompleteMultipartUploadHeaders) SetCommonHeaders(v map[string]*string) *CompleteMultipartUploadHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CompleteMultipartUploadHeaders) SetCompleteAll(v string) *CompleteMultipartUploadHeaders {
	s.CompleteAll = &v
	return s
}

func (s *CompleteMultipartUploadHeaders) SetForbidOverwrite(v string) *CompleteMultipartUploadHeaders {
	s.ForbidOverwrite = &v
	return s
}

type CompleteMultipartUploadRequest struct {
	// The container that stores the content of the CompleteMultipartUpload request.
	CompleteMultipartUpload *CompleteMultipartUpload `json:"CompleteMultipartUpload,omitempty" xml:"CompleteMultipartUpload,omitempty"`
	// The encodingtype of the object name in the response. Only URL encoding is supported.
	// The object name can contain characters that are encoded in UTF-8. However, the XML 1.0 standard cannot be used to parse control characters, such as characters with an ASCII value from 0 to 10. You can configure this parameter to encode the object name in the response.
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
	// The identifier of the multipart upload task.
	UploadId *string `json:"uploadId,omitempty" xml:"uploadId,omitempty"`
}

func (s CompleteMultipartUploadRequest) String() string {
	return tea.Prettify(s)
}

func (s CompleteMultipartUploadRequest) GoString() string {
	return s.String()
}

func (s *CompleteMultipartUploadRequest) SetCompleteMultipartUpload(v *CompleteMultipartUpload) *CompleteMultipartUploadRequest {
	s.CompleteMultipartUpload = v
	return s
}

func (s *CompleteMultipartUploadRequest) SetEncodingType(v string) *CompleteMultipartUploadRequest {
	s.EncodingType = &v
	return s
}

func (s *CompleteMultipartUploadRequest) SetUploadId(v string) *CompleteMultipartUploadRequest {
	s.UploadId = &v
	return s
}

type CompleteMultipartUploadResponseBody struct {
	// The container that stores the results of the CompleteMultipartUpload request.
	CompleteMultipartUploadResult *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult `json:"CompleteMultipartUploadResult,omitempty" xml:"CompleteMultipartUploadResult,omitempty" type:"Struct"`
}

func (s CompleteMultipartUploadResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CompleteMultipartUploadResponseBody) GoString() string {
	return s.String()
}

func (s *CompleteMultipartUploadResponseBody) SetCompleteMultipartUploadResult(v *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) *CompleteMultipartUploadResponseBody {
	s.CompleteMultipartUploadResult = v
	return s
}

type CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult struct {
	// The name of the bucket that contains the object you want to restore.
	Bucket *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	// The ETag that is generated when an object is created. ETags are used to identify the content of objects.
	//
	// If an object is created by calling the CompleteMultipartUpload operation, the ETag value is not the MD5 hash of the object content but a unique value calculated based on a specific rule.
	//
	// > The ETag of an object can be used to check whether the object content is modified. However, we recommend that you use the MD5 hash of an object rather than the ETag value of the object to verify data integrity.
	ETag *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	// The encoding type of the object name in the response. If this parameter is specified in the request, the object name is encoded in the response.
	EncodingType *string `json:"EncodingType,omitempty" xml:"EncodingType,omitempty"`
	// The name of the uploaded object.
	Key *string `json:"Key,omitempty" xml:"Key,omitempty"`
	// The URL that is used to access the uploaded object.
	Location *string `json:"Location,omitempty" xml:"Location,omitempty"`
}

func (s CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) String() string {
	return tea.Prettify(s)
}

func (s CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) GoString() string {
	return s.String()
}

func (s *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) SetBucket(v string) *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult {
	s.Bucket = &v
	return s
}

func (s *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) SetETag(v string) *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult {
	s.ETag = &v
	return s
}

func (s *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) SetEncodingType(v string) *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult {
	s.EncodingType = &v
	return s
}

func (s *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) SetKey(v string) *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult {
	s.Key = &v
	return s
}

func (s *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult) SetLocation(v string) *CompleteMultipartUploadResponseBodyCompleteMultipartUploadResult {
	s.Location = &v
	return s
}

type CompleteMultipartUploadResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *CompleteMultipartUploadResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s CompleteMultipartUploadResponse) String() string {
	return tea.Prettify(s)
}

func (s CompleteMultipartUploadResponse) GoString() string {
	return s.String()
}

func (s *CompleteMultipartUploadResponse) SetHeaders(v map[string]*string) *CompleteMultipartUploadResponse {
	s.Headers = v
	return s
}

func (s *CompleteMultipartUploadResponse) SetStatusCode(v int32) *CompleteMultipartUploadResponse {
	s.StatusCode = &v
	return s
}

func (s *CompleteMultipartUploadResponse) SetBody(v *CompleteMultipartUploadResponseBody) *CompleteMultipartUploadResponse {
	s.Body = v
	return s
}

type CopyObjectHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The path of the source object. By default, this header is left empty.
	CopySource *string `json:"x-oss-copy-source,omitempty" xml:"x-oss-copy-source,omitempty"`
	// The object copy condition. If the ETag value of the source object is the same as the ETag value that you specify in the request, OSS copies the object and returns 200 OK. By default, this header is left empty.
	CopySourceIfMatch *string `json:"x-oss-copy-source-if-match,omitempty" xml:"x-oss-copy-source-if-match,omitempty"`
	// If the source object is modified after the time that you specify in the request, OSS copies the object. By default, this header is left empty.
	CopySourceIfModifiedSince *string `json:"x-oss-copy-source-if-modified-since,omitempty" xml:"x-oss-copy-source-if-modified-since,omitempty"`
	// The object copy condition. If the ETag value of the source object is different from the ETag value that you specify in the request, OSS copies the object and returns 200 OK. By default, this header is left empty.
	CopySourceIfNoneMatch *string `json:"x-oss-copy-source-if-none-match,omitempty" xml:"x-oss-copy-source-if-none-match,omitempty"`
	// The object copy condition. If the time that you specify in the request is the same as or later than the modification time of the object, OSS copies the object and returns 200 OK. By default, this header is left empty.
	CopySourceIfUnmodifiedSince *string `json:"x-oss-copy-source-if-unmodified-since,omitempty" xml:"x-oss-copy-source-if-unmodified-since,omitempty"`
	// Specifies whether the CopyObject operation overwrites objects with the same name. The **x-oss-forbid-overwrite** request header does not take effect when versioning is enabled or suspended for the destination bucket. In this case, the CopyObject operation overwrites the existing object that has the same name as the destination object.
	//
	// *   If you do not specify the **x-oss-forbid-overwrite** header or set the header to **false**, an existing object that has the same name as the object that you want to copy is overwritten.****
	// *   If you set the **x-oss-forbid-overwrite** header to **true**, an existing object that has the same name as the object that you want to copy is not overwritten.
	//
	// If you specify the **x-oss-forbid-overwrite** header, the queries per second (QPS) performance of OSS may be degraded. If you want to specify the **x-oss-forbid-overwrite** header in a large number of requests (QPS greater than 1,000), contact technical support. Default value: false.
	ForbidOverwrite *string `json:"x-oss-forbid-overwrite,omitempty" xml:"x-oss-forbid-overwrite,omitempty"`
	// You can add parameters that contain the x-oss-meta- prefix when you create an append object. You cannot include these parameters in the requests when you append objects to an existing append object. Parameters that contain the x-oss-meta-\* prefix are considered the metadata of the object. You can specify multiple parameters that contain the x-oss-meta- prefix for an object. The total size of the metadata cannot exceed 8 KB. The names of parameters that contain the x-oss-meta- prefix can contain hyphens (-), digits, and letters. Uppercase letters are converted into lowercase letters. Other characters such as underscores (\_) are not supported.
	MetaData map[string]*string `json:"x-oss-meta-*,omitempty" xml:"x-oss-meta-*,omitempty"`
	// The method that is used to configure the metadata of the destination object. Default value: COPY.
	//
	// *   **COPY**: The metadata of the source object is copied to the destination object. The **x-oss-server-side-encryption** attribute of the source object is not copied to the destination object. The **x-oss-server-side-encryption** header in the CopyObject request specifies the method that is used to encrypt the destination object.
	// *   **REPLACE**: The metadata that you specify in the request is used as the metadata of the destination object.
	//
	// >  If the path of the source object is the same as the path of the destination object and versioning is disabled for the bucket in which the source and destination objects are stored, the metadata that you specify in the CopyObject request is used as the metadata of the destination object regardless of the value of the x-oss-metadata-directive header.
	MetadataDirective *string `json:"x-oss-metadata-directive,omitempty" xml:"x-oss-metadata-directive,omitempty"`
	// The access control list (ACL) of the destination object when the object is created. Default value: default.
	//
	// Valid values:
	//
	// *   default: The ACL of the object is the same as the ACL of the bucket in which the object is stored.
	// *   private: The ACL of the object is private. Only the owner of the object and authorized users have read and write permissions on the object. Other users do not have permissions on the object.
	// *   public-read: The ACL of the object is public-read. Only the owner of the object and authorized users have read and write permissions on the object. Other users have only read permissions on the object. Exercise caution when you set the ACL of the bucket to this value.
	// *   public-read-write: The ACL of the object is public-read-write. All users have read and write permissions on the object. Exercise caution when you set the ACL of the bucket to this value.
	//
	// For more information about ACLs, see [Object ACL](~~100676~~).
	Acl *string `json:"x-oss-object-acl,omitempty" xml:"x-oss-object-acl,omitempty"`
	// The server side data encryption algorithm. Invalid value: SM4
	XOssServerSideDataEncryption *string `json:"x-oss-server-side-data-encryption,omitempty" xml:"x-oss-server-side-data-encryption,omitempty"`
	// The entropy coding-based encryption algorithm that OSS uses to encrypt an object when you create the object. The valid values of the header are **AES256** and **KMS**. You must activate Key Management Service (KMS) in the OSS console before you can use the KMS encryption algorithm. Otherwise, the KmsServiceNotEnabled error is returned.
	//
	// *   If you do not specify the **x-oss-server-side-encryption** header in the CopyObject request, the destination object is not encrypted on the server regardless of whether the source object is encrypted on the server.
	// *   If you specify the **x-oss-server-side-encryption** header in the CopyObject request, the destination object is encrypted on the server after the CopyObject operation is performed regardless of whether the source object is encrypted on the server. In addition, the response to a CopyObject request contains the **x-oss-server-side-encryption** header whose value is the encryption algorithm of the destination object. When the destination object is downloaded, the **x-oss-server-side-encryption** header is included in the response. The value of this header is the encryption algorithm of the destination object.
	ServerSideEncryption *string `json:"x-oss-server-side-encryption,omitempty" xml:"x-oss-server-side-encryption,omitempty"`
	// The ID of the customer master key (CMK) that is managed by KMS. This parameter is available only if you set **x-oss-server-side-encryption** to KMS.
	SseKeyId *string `json:"x-oss-server-side-encryption-key-id,omitempty" xml:"x-oss-server-side-encryption-key-id,omitempty"`
	// The storage class of the object that you want to upload. Default value: Standard. If you specify a storage class when you upload the object, the storage class applies regardless of the storage class of the bucket to which you upload the object. For example, if you set **x-oss-storage-class** to Standard when you upload an object to an IA bucket, the storage class of the uploaded object is Standard.
	//
	// Valid values:
	//
	// *   Standard
	// *   IA
	// *   Archive
	// *   ColdArchive
	//
	// For more information about storage classes, see [Overview](~~51374~~).
	StorageClass *string `json:"x-oss-storage-class,omitempty" xml:"x-oss-storage-class,omitempty"`
	// The tag of the destination object. You can add multiple tags to the destination object. Example: TagA=A\&TagB=B.
	//
	// >  The tag key and tag value must be URL-encoded. If a key-value pair does not contain an equal sign (=), the tag value is considered an empty string.
	Tagging *string `json:"x-oss-tagging,omitempty" xml:"x-oss-tagging,omitempty"`
	// The method that is used to add tags to the destination object. Default value: Copy. Valid values:
	//
	// *   **Copy**: The tags of the source object are copied to the destination object.
	// *   **Replace**: The tags that you specify in the request are added to the destination object.
	TaggingDirective *string `json:"x-oss-tagging-directive,omitempty" xml:"x-oss-tagging-directive,omitempty"`
}

func (s CopyObjectHeaders) String() string {
	return tea.Prettify(s)
}

func (s CopyObjectHeaders) GoString() string {
	return s.String()
}

func (s *CopyObjectHeaders) SetCommonHeaders(v map[string]*string) *CopyObjectHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CopyObjectHeaders) SetCopySource(v string) *CopyObjectHeaders {
	s.CopySource = &v
	return s
}

func (s *CopyObjectHeaders) SetCopySourceIfMatch(v string) *CopyObjectHeaders {
	s.CopySourceIfMatch = &v
	return s
}

func (s *CopyObjectHeaders) SetCopySourceIfModifiedSince(v string) *CopyObjectHeaders {
	s.CopySourceIfModifiedSince = &v
	return s
}

func (s *CopyObjectHeaders) SetCopySourceIfNoneMatch(v string) *CopyObjectHeaders {
	s.CopySourceIfNoneMatch = &v
	return s
}

func (s *CopyObjectHeaders) SetCopySourceIfUnmodifiedSince(v string) *CopyObjectHeaders {
	s.CopySourceIfUnmodifiedSince = &v
	return s
}

func (s *CopyObjectHeaders) SetForbidOverwrite(v string) *CopyObjectHeaders {
	s.ForbidOverwrite = &v
	return s
}

func (s *CopyObjectHeaders) SetMetaData(v map[string]*string) *CopyObjectHeaders {
	s.MetaData = v
	return s
}

func (s *CopyObjectHeaders) SetMetadataDirective(v string) *CopyObjectHeaders {
	s.MetadataDirective = &v
	return s
}

func (s *CopyObjectHeaders) SetAcl(v string) *CopyObjectHeaders {
	s.Acl = &v
	return s
}

func (s *CopyObjectHeaders) SetXOssServerSideDataEncryption(v string) *CopyObjectHeaders {
	s.XOssServerSideDataEncryption = &v
	return s
}

func (s *CopyObjectHeaders) SetServerSideEncryption(v string) *CopyObjectHeaders {
	s.ServerSideEncryption = &v
	return s
}

func (s *CopyObjectHeaders) SetSseKeyId(v string) *CopyObjectHeaders {
	s.SseKeyId = &v
	return s
}

func (s *CopyObjectHeaders) SetStorageClass(v string) *CopyObjectHeaders {
	s.StorageClass = &v
	return s
}

func (s *CopyObjectHeaders) SetTagging(v string) *CopyObjectHeaders {
	s.Tagging = &v
	return s
}

func (s *CopyObjectHeaders) SetTaggingDirective(v string) *CopyObjectHeaders {
	s.TaggingDirective = &v
	return s
}

type CopyObjectResponseBody struct {
	// The results of the CopyObject operation.
	CopyObjectResult *CopyObjectResponseBodyCopyObjectResult `json:"CopyObjectResult,omitempty" xml:"CopyObjectResult,omitempty" type:"Struct"`
}

func (s CopyObjectResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CopyObjectResponseBody) GoString() string {
	return s.String()
}

func (s *CopyObjectResponseBody) SetCopyObjectResult(v *CopyObjectResponseBodyCopyObjectResult) *CopyObjectResponseBody {
	s.CopyObjectResult = v
	return s
}

type CopyObjectResponseBodyCopyObjectResult struct {
	// The ETag value of the destination object.
	ETag *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	// The time when the destination object was last modified.
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
}

func (s CopyObjectResponseBodyCopyObjectResult) String() string {
	return tea.Prettify(s)
}

func (s CopyObjectResponseBodyCopyObjectResult) GoString() string {
	return s.String()
}

func (s *CopyObjectResponseBodyCopyObjectResult) SetETag(v string) *CopyObjectResponseBodyCopyObjectResult {
	s.ETag = &v
	return s
}

func (s *CopyObjectResponseBodyCopyObjectResult) SetLastModified(v string) *CopyObjectResponseBodyCopyObjectResult {
	s.LastModified = &v
	return s
}

type CopyObjectResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *CopyObjectResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s CopyObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s CopyObjectResponse) GoString() string {
	return s.String()
}

func (s *CopyObjectResponse) SetHeaders(v map[string]*string) *CopyObjectResponse {
	s.Headers = v
	return s
}

func (s *CopyObjectResponse) SetStatusCode(v int32) *CopyObjectResponse {
	s.StatusCode = &v
	return s
}

func (s *CopyObjectResponse) SetBody(v *CopyObjectResponseBody) *CopyObjectResponse {
	s.Body = v
	return s
}

type CreateAccessPointRequest struct {
	// The container that stores the information about the access point.
	CreateAccessPointConfiguration *CreateAccessPointConfiguration `json:"CreateAccessPointConfiguration,omitempty" xml:"CreateAccessPointConfiguration,omitempty"`
}

func (s CreateAccessPointRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointRequest) GoString() string {
	return s.String()
}

func (s *CreateAccessPointRequest) SetCreateAccessPointConfiguration(v *CreateAccessPointConfiguration) *CreateAccessPointRequest {
	s.CreateAccessPointConfiguration = v
	return s
}

type CreateAccessPointResponseBody struct {
	// The container that stores the information about the access point.
	CreateAccessPointResult *CreateAccessPointResult `json:"CreateAccessPointResult,omitempty" xml:"CreateAccessPointResult,omitempty"`
}

func (s CreateAccessPointResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointResponseBody) GoString() string {
	return s.String()
}

func (s *CreateAccessPointResponseBody) SetCreateAccessPointResult(v *CreateAccessPointResult) *CreateAccessPointResponseBody {
	s.CreateAccessPointResult = v
	return s
}

type CreateAccessPointResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *CreateAccessPointResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s CreateAccessPointResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointResponse) GoString() string {
	return s.String()
}

func (s *CreateAccessPointResponse) SetHeaders(v map[string]*string) *CreateAccessPointResponse {
	s.Headers = v
	return s
}

func (s *CreateAccessPointResponse) SetStatusCode(v int32) *CreateAccessPointResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateAccessPointResponse) SetBody(v *CreateAccessPointResponseBody) *CreateAccessPointResponse {
	s.Body = v
	return s
}

type CreateAccessPointForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s CreateAccessPointForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *CreateAccessPointForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *CreateAccessPointForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *CreateAccessPointForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *CreateAccessPointForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type CreateAccessPointForObjectProcessRequest struct {
	CreateAccessPointForObjectProcessConfiguration *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration `json:"CreateAccessPointForObjectProcessConfiguration,omitempty" xml:"CreateAccessPointForObjectProcessConfiguration,omitempty" type:"Struct"`
}

func (s CreateAccessPointForObjectProcessRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointForObjectProcessRequest) GoString() string {
	return s.String()
}

func (s *CreateAccessPointForObjectProcessRequest) SetCreateAccessPointForObjectProcessConfiguration(v *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration) *CreateAccessPointForObjectProcessRequest {
	s.CreateAccessPointForObjectProcessConfiguration = v
	return s
}

type CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration struct {
	AccessPointName                      *string                     `json:"AccessPointName,omitempty" xml:"AccessPointName,omitempty"`
	AllowAnonymousAccessForObjectProcess *string                     `json:"AllowAnonymousAccessForObjectProcess,omitempty" xml:"AllowAnonymousAccessForObjectProcess,omitempty"`
	ObjectProcessConfiguration           *ObjectProcessConfiguration `json:"ObjectProcessConfiguration,omitempty" xml:"ObjectProcessConfiguration,omitempty"`
}

func (s CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration) GoString() string {
	return s.String()
}

func (s *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration) SetAccessPointName(v string) *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration {
	s.AccessPointName = &v
	return s
}

func (s *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration) SetAllowAnonymousAccessForObjectProcess(v string) *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration {
	s.AllowAnonymousAccessForObjectProcess = &v
	return s
}

func (s *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration) SetObjectProcessConfiguration(v *ObjectProcessConfiguration) *CreateAccessPointForObjectProcessRequestCreateAccessPointForObjectProcessConfiguration {
	s.ObjectProcessConfiguration = v
	return s
}

type CreateAccessPointForObjectProcessResponseBody struct {
	CreateAccessPointForObjectProcessResult *CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult `json:"CreateAccessPointForObjectProcessResult,omitempty" xml:"CreateAccessPointForObjectProcessResult,omitempty" type:"Struct"`
}

func (s CreateAccessPointForObjectProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointForObjectProcessResponseBody) GoString() string {
	return s.String()
}

func (s *CreateAccessPointForObjectProcessResponseBody) SetCreateAccessPointForObjectProcessResult(v *CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult) *CreateAccessPointForObjectProcessResponseBody {
	s.CreateAccessPointForObjectProcessResult = v
	return s
}

type CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult struct {
	AccessPointForObjectProcessAlias *string `json:"AccessPointForObjectProcessAlias,omitempty" xml:"AccessPointForObjectProcessAlias,omitempty"`
	AccessPointForObjectProcessArn   *string `json:"AccessPointForObjectProcessArn,omitempty" xml:"AccessPointForObjectProcessArn,omitempty"`
}

func (s CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult) GoString() string {
	return s.String()
}

func (s *CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult) SetAccessPointForObjectProcessAlias(v string) *CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult {
	s.AccessPointForObjectProcessAlias = &v
	return s
}

func (s *CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult) SetAccessPointForObjectProcessArn(v string) *CreateAccessPointForObjectProcessResponseBodyCreateAccessPointForObjectProcessResult {
	s.AccessPointForObjectProcessArn = &v
	return s
}

type CreateAccessPointForObjectProcessResponse struct {
	Headers    map[string]*string                             `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                         `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *CreateAccessPointForObjectProcessResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s CreateAccessPointForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAccessPointForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *CreateAccessPointForObjectProcessResponse) SetHeaders(v map[string]*string) *CreateAccessPointForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *CreateAccessPointForObjectProcessResponse) SetStatusCode(v int32) *CreateAccessPointForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateAccessPointForObjectProcessResponse) SetBody(v *CreateAccessPointForObjectProcessResponseBody) *CreateAccessPointForObjectProcessResponse {
	s.Body = v
	return s
}

type CreateBucketDataRedundancyTransitionRequest struct {
	XOssTargetRedundancyType *string `json:"x-oss-target-redundancy-type,omitempty" xml:"x-oss-target-redundancy-type,omitempty"`
}

func (s CreateBucketDataRedundancyTransitionRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateBucketDataRedundancyTransitionRequest) GoString() string {
	return s.String()
}

func (s *CreateBucketDataRedundancyTransitionRequest) SetXOssTargetRedundancyType(v string) *CreateBucketDataRedundancyTransitionRequest {
	s.XOssTargetRedundancyType = &v
	return s
}

type CreateBucketDataRedundancyTransitionResponseBody struct {
	BucketDataRedundancyTransition *CreateBucketDataRedundancyTransitionResponseBodyBucketDataRedundancyTransition `json:"BucketDataRedundancyTransition,omitempty" xml:"BucketDataRedundancyTransition,omitempty" type:"Struct"`
}

func (s CreateBucketDataRedundancyTransitionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateBucketDataRedundancyTransitionResponseBody) GoString() string {
	return s.String()
}

func (s *CreateBucketDataRedundancyTransitionResponseBody) SetBucketDataRedundancyTransition(v *CreateBucketDataRedundancyTransitionResponseBodyBucketDataRedundancyTransition) *CreateBucketDataRedundancyTransitionResponseBody {
	s.BucketDataRedundancyTransition = v
	return s
}

type CreateBucketDataRedundancyTransitionResponseBodyBucketDataRedundancyTransition struct {
	TaskId *string `json:"TaskId,omitempty" xml:"TaskId,omitempty"`
}

func (s CreateBucketDataRedundancyTransitionResponseBodyBucketDataRedundancyTransition) String() string {
	return tea.Prettify(s)
}

func (s CreateBucketDataRedundancyTransitionResponseBodyBucketDataRedundancyTransition) GoString() string {
	return s.String()
}

func (s *CreateBucketDataRedundancyTransitionResponseBodyBucketDataRedundancyTransition) SetTaskId(v string) *CreateBucketDataRedundancyTransitionResponseBodyBucketDataRedundancyTransition {
	s.TaskId = &v
	return s
}

type CreateBucketDataRedundancyTransitionResponse struct {
	Headers    map[string]*string                                `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                            `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *CreateBucketDataRedundancyTransitionResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s CreateBucketDataRedundancyTransitionResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateBucketDataRedundancyTransitionResponse) GoString() string {
	return s.String()
}

func (s *CreateBucketDataRedundancyTransitionResponse) SetHeaders(v map[string]*string) *CreateBucketDataRedundancyTransitionResponse {
	s.Headers = v
	return s
}

func (s *CreateBucketDataRedundancyTransitionResponse) SetStatusCode(v int32) *CreateBucketDataRedundancyTransitionResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateBucketDataRedundancyTransitionResponse) SetBody(v *CreateBucketDataRedundancyTransitionResponseBody) *CreateBucketDataRedundancyTransitionResponse {
	s.Body = v
	return s
}

type CreateCnameTokenRequest struct {
	// The container that stores the CNAME record.
	BucketCnameConfiguration *CreateCnameTokenRequestBucketCnameConfiguration `json:"BucketCnameConfiguration,omitempty" xml:"BucketCnameConfiguration,omitempty" type:"Struct"`
}

func (s CreateCnameTokenRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateCnameTokenRequest) GoString() string {
	return s.String()
}

func (s *CreateCnameTokenRequest) SetBucketCnameConfiguration(v *CreateCnameTokenRequestBucketCnameConfiguration) *CreateCnameTokenRequest {
	s.BucketCnameConfiguration = v
	return s
}

type CreateCnameTokenRequestBucketCnameConfiguration struct {
	// The container that stores the CNAME information.
	Cname *CreateCnameTokenRequestBucketCnameConfigurationCname `json:"Cname,omitempty" xml:"Cname,omitempty" type:"Struct"`
}

func (s CreateCnameTokenRequestBucketCnameConfiguration) String() string {
	return tea.Prettify(s)
}

func (s CreateCnameTokenRequestBucketCnameConfiguration) GoString() string {
	return s.String()
}

func (s *CreateCnameTokenRequestBucketCnameConfiguration) SetCname(v *CreateCnameTokenRequestBucketCnameConfigurationCname) *CreateCnameTokenRequestBucketCnameConfiguration {
	s.Cname = v
	return s
}

type CreateCnameTokenRequestBucketCnameConfigurationCname struct {
	// The custom domain name.
	Domain *string `json:"Domain,omitempty" xml:"Domain,omitempty"`
}

func (s CreateCnameTokenRequestBucketCnameConfigurationCname) String() string {
	return tea.Prettify(s)
}

func (s CreateCnameTokenRequestBucketCnameConfigurationCname) GoString() string {
	return s.String()
}

func (s *CreateCnameTokenRequestBucketCnameConfigurationCname) SetDomain(v string) *CreateCnameTokenRequestBucketCnameConfigurationCname {
	s.Domain = &v
	return s
}

type CreateCnameTokenResponseBody struct {
	// The container in which the CNAME token is stored.
	CnameToken *CnameToken `json:"CnameToken,omitempty" xml:"CnameToken,omitempty"`
}

func (s CreateCnameTokenResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateCnameTokenResponseBody) GoString() string {
	return s.String()
}

func (s *CreateCnameTokenResponseBody) SetCnameToken(v *CnameToken) *CreateCnameTokenResponseBody {
	s.CnameToken = v
	return s
}

type CreateCnameTokenResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *CreateCnameTokenResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s CreateCnameTokenResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateCnameTokenResponse) GoString() string {
	return s.String()
}

func (s *CreateCnameTokenResponse) SetHeaders(v map[string]*string) *CreateCnameTokenResponse {
	s.Headers = v
	return s
}

func (s *CreateCnameTokenResponse) SetStatusCode(v int32) *CreateCnameTokenResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateCnameTokenResponse) SetBody(v *CreateCnameTokenResponseBody) *CreateCnameTokenResponse {
	s.Body = v
	return s
}

type CreateSelectObjectMetaRequest struct {
	// The container that stores SelectMetaRequest information.
	SelectMetaRequest *SelectMetaRequest `json:"SelectMetaRequest,omitempty" xml:"SelectMetaRequest,omitempty"`
}

func (s CreateSelectObjectMetaRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateSelectObjectMetaRequest) GoString() string {
	return s.String()
}

func (s *CreateSelectObjectMetaRequest) SetSelectMetaRequest(v *SelectMetaRequest) *CreateSelectObjectMetaRequest {
	s.SelectMetaRequest = v
	return s
}

type CreateSelectObjectMetaResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *SelectMetaStatus  `json:"body,omitempty" xml:"body,omitempty"`
}

func (s CreateSelectObjectMetaResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateSelectObjectMetaResponse) GoString() string {
	return s.String()
}

func (s *CreateSelectObjectMetaResponse) SetHeaders(v map[string]*string) *CreateSelectObjectMetaResponse {
	s.Headers = v
	return s
}

func (s *CreateSelectObjectMetaResponse) SetStatusCode(v int32) *CreateSelectObjectMetaResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateSelectObjectMetaResponse) SetBody(v *SelectMetaStatus) *CreateSelectObjectMetaResponse {
	s.Body = v
	return s
}

type DeleteAccessPointHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The name of the access point.
	XOssAccessPointName *string `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s DeleteAccessPointHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointHeaders) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointHeaders) SetCommonHeaders(v map[string]*string) *DeleteAccessPointHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteAccessPointHeaders) SetXOssAccessPointName(v string) *DeleteAccessPointHeaders {
	s.XOssAccessPointName = &v
	return s
}

type DeleteAccessPointResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteAccessPointResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointResponse) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointResponse) SetHeaders(v map[string]*string) *DeleteAccessPointResponse {
	s.Headers = v
	return s
}

func (s *DeleteAccessPointResponse) SetStatusCode(v int32) *DeleteAccessPointResponse {
	s.StatusCode = &v
	return s
}

type DeleteAccessPointForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s DeleteAccessPointForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *DeleteAccessPointForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteAccessPointForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *DeleteAccessPointForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type DeleteAccessPointForObjectProcessResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteAccessPointForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointForObjectProcessResponse) SetHeaders(v map[string]*string) *DeleteAccessPointForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *DeleteAccessPointForObjectProcessResponse) SetStatusCode(v int32) *DeleteAccessPointForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

type DeleteAccessPointPolicyHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The name of the access point.
	XOssAccessPointName *string `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s DeleteAccessPointPolicyHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointPolicyHeaders) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointPolicyHeaders) SetCommonHeaders(v map[string]*string) *DeleteAccessPointPolicyHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteAccessPointPolicyHeaders) SetXOssAccessPointName(v string) *DeleteAccessPointPolicyHeaders {
	s.XOssAccessPointName = &v
	return s
}

type DeleteAccessPointPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteAccessPointPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointPolicyResponse) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointPolicyResponse) SetHeaders(v map[string]*string) *DeleteAccessPointPolicyResponse {
	s.Headers = v
	return s
}

func (s *DeleteAccessPointPolicyResponse) SetStatusCode(v int32) *DeleteAccessPointPolicyResponse {
	s.StatusCode = &v
	return s
}

type DeleteAccessPointPolicyForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s DeleteAccessPointPolicyForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointPolicyForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointPolicyForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *DeleteAccessPointPolicyForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteAccessPointPolicyForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *DeleteAccessPointPolicyForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type DeleteAccessPointPolicyForObjectProcessResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteAccessPointPolicyForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointPolicyForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointPolicyForObjectProcessResponse) SetHeaders(v map[string]*string) *DeleteAccessPointPolicyForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *DeleteAccessPointPolicyForObjectProcessResponse) SetStatusCode(v int32) *DeleteAccessPointPolicyForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

type DeleteAccessPointPublicAccessBlockRequest struct {
	XOssAccessPointName *string `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s DeleteAccessPointPublicAccessBlockRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointPublicAccessBlockRequest) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointPublicAccessBlockRequest) SetXOssAccessPointName(v string) *DeleteAccessPointPublicAccessBlockRequest {
	s.XOssAccessPointName = &v
	return s
}

type DeleteAccessPointPublicAccessBlockResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteAccessPointPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAccessPointPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *DeleteAccessPointPublicAccessBlockResponse) SetHeaders(v map[string]*string) *DeleteAccessPointPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *DeleteAccessPointPublicAccessBlockResponse) SetStatusCode(v int32) *DeleteAccessPointPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketResponse) SetHeaders(v map[string]*string) *DeleteBucketResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketResponse) SetStatusCode(v int32) *DeleteBucketResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketCallbackPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketCallbackPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketCallbackPolicyResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketCallbackPolicyResponse) SetHeaders(v map[string]*string) *DeleteBucketCallbackPolicyResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketCallbackPolicyResponse) SetStatusCode(v int32) *DeleteBucketCallbackPolicyResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketCorsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketCorsResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketCorsResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketCorsResponse) SetHeaders(v map[string]*string) *DeleteBucketCorsResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketCorsResponse) SetStatusCode(v int32) *DeleteBucketCorsResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketDataRedundancyTransitionRequest struct {
	XOssRedundancyTransitionTaskid *string `json:"x-oss-redundancy-transition-taskid,omitempty" xml:"x-oss-redundancy-transition-taskid,omitempty"`
}

func (s DeleteBucketDataRedundancyTransitionRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketDataRedundancyTransitionRequest) GoString() string {
	return s.String()
}

func (s *DeleteBucketDataRedundancyTransitionRequest) SetXOssRedundancyTransitionTaskid(v string) *DeleteBucketDataRedundancyTransitionRequest {
	s.XOssRedundancyTransitionTaskid = &v
	return s
}

type DeleteBucketDataRedundancyTransitionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketDataRedundancyTransitionResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketDataRedundancyTransitionResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketDataRedundancyTransitionResponse) SetHeaders(v map[string]*string) *DeleteBucketDataRedundancyTransitionResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketDataRedundancyTransitionResponse) SetStatusCode(v int32) *DeleteBucketDataRedundancyTransitionResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketEncryptionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketEncryptionResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketEncryptionResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketEncryptionResponse) SetHeaders(v map[string]*string) *DeleteBucketEncryptionResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketEncryptionResponse) SetStatusCode(v int32) *DeleteBucketEncryptionResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketInventoryRequest struct {
	// The name of the inventory that you want to delete.
	InventoryId *string `json:"inventoryId,omitempty" xml:"inventoryId,omitempty"`
}

func (s DeleteBucketInventoryRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketInventoryRequest) GoString() string {
	return s.String()
}

func (s *DeleteBucketInventoryRequest) SetInventoryId(v string) *DeleteBucketInventoryRequest {
	s.InventoryId = &v
	return s
}

type DeleteBucketInventoryResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketInventoryResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketInventoryResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketInventoryResponse) SetHeaders(v map[string]*string) *DeleteBucketInventoryResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketInventoryResponse) SetStatusCode(v int32) *DeleteBucketInventoryResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketLifecycleResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketLifecycleResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketLifecycleResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketLifecycleResponse) SetHeaders(v map[string]*string) *DeleteBucketLifecycleResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketLifecycleResponse) SetStatusCode(v int32) *DeleteBucketLifecycleResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketLoggingResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketLoggingResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketLoggingResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketLoggingResponse) SetHeaders(v map[string]*string) *DeleteBucketLoggingResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketLoggingResponse) SetStatusCode(v int32) *DeleteBucketLoggingResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketPolicyResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketPolicyResponse) SetHeaders(v map[string]*string) *DeleteBucketPolicyResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketPolicyResponse) SetStatusCode(v int32) *DeleteBucketPolicyResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketPublicAccessBlockResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketPublicAccessBlockResponse) SetHeaders(v map[string]*string) *DeleteBucketPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketPublicAccessBlockResponse) SetStatusCode(v int32) *DeleteBucketPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketReplicationRequest struct {
	// The container that is used to store the data replication rule to delete.
	ReplicationRules *DeleteBucketReplicationRequestReplicationRules `json:"ReplicationRules,omitempty" xml:"ReplicationRules,omitempty" type:"Struct"`
}

func (s DeleteBucketReplicationRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketReplicationRequest) GoString() string {
	return s.String()
}

func (s *DeleteBucketReplicationRequest) SetReplicationRules(v *DeleteBucketReplicationRequestReplicationRules) *DeleteBucketReplicationRequest {
	s.ReplicationRules = v
	return s
}

type DeleteBucketReplicationRequestReplicationRules struct {
	ID *string `json:"ID,omitempty" xml:"ID,omitempty"`
}

func (s DeleteBucketReplicationRequestReplicationRules) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketReplicationRequestReplicationRules) GoString() string {
	return s.String()
}

func (s *DeleteBucketReplicationRequestReplicationRules) SetID(v string) *DeleteBucketReplicationRequestReplicationRules {
	s.ID = &v
	return s
}

type DeleteBucketReplicationResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketReplicationResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketReplicationResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketReplicationResponse) SetHeaders(v map[string]*string) *DeleteBucketReplicationResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketReplicationResponse) SetStatusCode(v int32) *DeleteBucketReplicationResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketResponseHeaderResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketResponseHeaderResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketResponseHeaderResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketResponseHeaderResponse) SetHeaders(v map[string]*string) *DeleteBucketResponseHeaderResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketResponseHeaderResponse) SetStatusCode(v int32) *DeleteBucketResponseHeaderResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketTagsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketTagsResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketTagsResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketTagsResponse) SetHeaders(v map[string]*string) *DeleteBucketTagsResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketTagsResponse) SetStatusCode(v int32) *DeleteBucketTagsResponse {
	s.StatusCode = &v
	return s
}

type DeleteBucketWebsiteResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteBucketWebsiteResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteBucketWebsiteResponse) GoString() string {
	return s.String()
}

func (s *DeleteBucketWebsiteResponse) SetHeaders(v map[string]*string) *DeleteBucketWebsiteResponse {
	s.Headers = v
	return s
}

func (s *DeleteBucketWebsiteResponse) SetStatusCode(v int32) *DeleteBucketWebsiteResponse {
	s.StatusCode = &v
	return s
}

type DeleteCnameRequest struct {
	// The container that stores the CNAME record.
	BucketCnameConfiguration *DeleteCnameRequestBucketCnameConfiguration `json:"BucketCnameConfiguration,omitempty" xml:"BucketCnameConfiguration,omitempty" type:"Struct"`
}

func (s DeleteCnameRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteCnameRequest) GoString() string {
	return s.String()
}

func (s *DeleteCnameRequest) SetBucketCnameConfiguration(v *DeleteCnameRequestBucketCnameConfiguration) *DeleteCnameRequest {
	s.BucketCnameConfiguration = v
	return s
}

type DeleteCnameRequestBucketCnameConfiguration struct {
	Cname *DeleteCnameRequestBucketCnameConfigurationCname `json:"Cname,omitempty" xml:"Cname,omitempty" type:"Struct"`
}

func (s DeleteCnameRequestBucketCnameConfiguration) String() string {
	return tea.Prettify(s)
}

func (s DeleteCnameRequestBucketCnameConfiguration) GoString() string {
	return s.String()
}

func (s *DeleteCnameRequestBucketCnameConfiguration) SetCname(v *DeleteCnameRequestBucketCnameConfigurationCname) *DeleteCnameRequestBucketCnameConfiguration {
	s.Cname = v
	return s
}

type DeleteCnameRequestBucketCnameConfigurationCname struct {
	Domain *string `json:"Domain,omitempty" xml:"Domain,omitempty"`
}

func (s DeleteCnameRequestBucketCnameConfigurationCname) String() string {
	return tea.Prettify(s)
}

func (s DeleteCnameRequestBucketCnameConfigurationCname) GoString() string {
	return s.String()
}

func (s *DeleteCnameRequestBucketCnameConfigurationCname) SetDomain(v string) *DeleteCnameRequestBucketCnameConfigurationCname {
	s.Domain = &v
	return s
}

type DeleteCnameResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteCnameResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteCnameResponse) GoString() string {
	return s.String()
}

func (s *DeleteCnameResponse) SetHeaders(v map[string]*string) *DeleteCnameResponse {
	s.Headers = v
	return s
}

func (s *DeleteCnameResponse) SetStatusCode(v int32) *DeleteCnameResponse {
	s.StatusCode = &v
	return s
}

type DeleteLiveChannelResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteLiveChannelResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteLiveChannelResponse) GoString() string {
	return s.String()
}

func (s *DeleteLiveChannelResponse) SetHeaders(v map[string]*string) *DeleteLiveChannelResponse {
	s.Headers = v
	return s
}

func (s *DeleteLiveChannelResponse) SetStatusCode(v int32) *DeleteLiveChannelResponse {
	s.StatusCode = &v
	return s
}

type DeleteMultipleObjectsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	ContentMd5    *string            `json:"content-md5,omitempty" xml:"content-md5,omitempty"`
}

func (s DeleteMultipleObjectsHeaders) String() string {
	return tea.Prettify(s)
}

func (s DeleteMultipleObjectsHeaders) GoString() string {
	return s.String()
}

func (s *DeleteMultipleObjectsHeaders) SetCommonHeaders(v map[string]*string) *DeleteMultipleObjectsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *DeleteMultipleObjectsHeaders) SetContentMd5(v string) *DeleteMultipleObjectsHeaders {
	s.ContentMd5 = &v
	return s
}

type DeleteMultipleObjectsRequest struct {
	Delete *Delete `json:"Delete,omitempty" xml:"Delete,omitempty"`
	// The encoding type of the object name in the response. The value of the Key parameter is UTF-8 encoded. If the Key parameter includes control characters that are not supported by the XML 1.0 standard, you can specify this header to encode the value of the Key parameter in the response.
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
}

func (s DeleteMultipleObjectsRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteMultipleObjectsRequest) GoString() string {
	return s.String()
}

func (s *DeleteMultipleObjectsRequest) SetDelete(v *Delete) *DeleteMultipleObjectsRequest {
	s.Delete = v
	return s
}

func (s *DeleteMultipleObjectsRequest) SetEncodingType(v string) *DeleteMultipleObjectsRequest {
	s.EncodingType = &v
	return s
}

type DeleteMultipleObjectsResponseBody struct {
	DeleteResult *DeleteMultipleObjectsResponseBodyDeleteResult `json:"DeleteResult,omitempty" xml:"DeleteResult,omitempty" type:"Struct"`
}

func (s DeleteMultipleObjectsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteMultipleObjectsResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteMultipleObjectsResponseBody) SetDeleteResult(v *DeleteMultipleObjectsResponseBodyDeleteResult) *DeleteMultipleObjectsResponseBody {
	s.DeleteResult = v
	return s
}

type DeleteMultipleObjectsResponseBodyDeleteResult struct {
	Deleted      []*DeletedObject `json:"Deleted,omitempty" xml:"Deleted,omitempty" type:"Repeated"`
	EncodingType *string          `json:"EncodingType,omitempty" xml:"EncodingType,omitempty"`
}

func (s DeleteMultipleObjectsResponseBodyDeleteResult) String() string {
	return tea.Prettify(s)
}

func (s DeleteMultipleObjectsResponseBodyDeleteResult) GoString() string {
	return s.String()
}

func (s *DeleteMultipleObjectsResponseBodyDeleteResult) SetDeleted(v []*DeletedObject) *DeleteMultipleObjectsResponseBodyDeleteResult {
	s.Deleted = v
	return s
}

func (s *DeleteMultipleObjectsResponseBodyDeleteResult) SetEncodingType(v string) *DeleteMultipleObjectsResponseBodyDeleteResult {
	s.EncodingType = &v
	return s
}

type DeleteMultipleObjectsResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *DeleteMultipleObjectsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s DeleteMultipleObjectsResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteMultipleObjectsResponse) GoString() string {
	return s.String()
}

func (s *DeleteMultipleObjectsResponse) SetHeaders(v map[string]*string) *DeleteMultipleObjectsResponse {
	s.Headers = v
	return s
}

func (s *DeleteMultipleObjectsResponse) SetStatusCode(v int32) *DeleteMultipleObjectsResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteMultipleObjectsResponse) SetBody(v *DeleteMultipleObjectsResponseBody) *DeleteMultipleObjectsResponse {
	s.Body = v
	return s
}

type DeleteObjectRequest struct {
	// The version ID of the object.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s DeleteObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteObjectRequest) GoString() string {
	return s.String()
}

func (s *DeleteObjectRequest) SetVersionId(v string) *DeleteObjectRequest {
	s.VersionId = &v
	return s
}

type DeleteObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteObjectResponse) GoString() string {
	return s.String()
}

func (s *DeleteObjectResponse) SetHeaders(v map[string]*string) *DeleteObjectResponse {
	s.Headers = v
	return s
}

func (s *DeleteObjectResponse) SetStatusCode(v int32) *DeleteObjectResponse {
	s.StatusCode = &v
	return s
}

type DeleteObjectTaggingRequest struct {
	// The version ID of the object that you want to delete.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s DeleteObjectTaggingRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteObjectTaggingRequest) GoString() string {
	return s.String()
}

func (s *DeleteObjectTaggingRequest) SetVersionId(v string) *DeleteObjectTaggingRequest {
	s.VersionId = &v
	return s
}

type DeleteObjectTaggingResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteObjectTaggingResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteObjectTaggingResponse) GoString() string {
	return s.String()
}

func (s *DeleteObjectTaggingResponse) SetHeaders(v map[string]*string) *DeleteObjectTaggingResponse {
	s.Headers = v
	return s
}

func (s *DeleteObjectTaggingResponse) SetStatusCode(v int32) *DeleteObjectTaggingResponse {
	s.StatusCode = &v
	return s
}

type DeletePublicAccessBlockResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeletePublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s DeletePublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *DeletePublicAccessBlockResponse) SetHeaders(v map[string]*string) *DeletePublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *DeletePublicAccessBlockResponse) SetStatusCode(v int32) *DeletePublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

type DeleteStyleRequest struct {
	// The name of the image style.
	StyleName *string `json:"styleName,omitempty" xml:"styleName,omitempty"`
}

func (s DeleteStyleRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteStyleRequest) GoString() string {
	return s.String()
}

func (s *DeleteStyleRequest) SetStyleName(v string) *DeleteStyleRequest {
	s.StyleName = &v
	return s
}

type DeleteStyleResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteStyleResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteStyleResponse) GoString() string {
	return s.String()
}

func (s *DeleteStyleResponse) SetHeaders(v map[string]*string) *DeleteStyleResponse {
	s.Headers = v
	return s
}

func (s *DeleteStyleResponse) SetStatusCode(v int32) *DeleteStyleResponse {
	s.StatusCode = &v
	return s
}

type DeleteUserDefinedLogFieldsConfigResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteUserDefinedLogFieldsConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteUserDefinedLogFieldsConfigResponse) GoString() string {
	return s.String()
}

func (s *DeleteUserDefinedLogFieldsConfigResponse) SetHeaders(v map[string]*string) *DeleteUserDefinedLogFieldsConfigResponse {
	s.Headers = v
	return s
}

func (s *DeleteUserDefinedLogFieldsConfigResponse) SetStatusCode(v int32) *DeleteUserDefinedLogFieldsConfigResponse {
	s.StatusCode = &v
	return s
}

type DescribeRegionsRequest struct {
	// The region ID of the request.
	Regions *string `json:"regions,omitempty" xml:"regions,omitempty"`
}

func (s DescribeRegionsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsRequest) GoString() string {
	return s.String()
}

func (s *DescribeRegionsRequest) SetRegions(v string) *DescribeRegionsRequest {
	s.Regions = &v
	return s
}

type DescribeRegionsResponseBody struct {
	// The information about the regions.
	RegionInfoList *DescribeRegionsResponseBodyRegionInfoList `json:"RegionInfoList,omitempty" xml:"RegionInfoList,omitempty" type:"Struct"`
}

func (s DescribeRegionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponseBody) SetRegionInfoList(v *DescribeRegionsResponseBodyRegionInfoList) *DescribeRegionsResponseBody {
	s.RegionInfoList = v
	return s
}

type DescribeRegionsResponseBodyRegionInfoList struct {
	// The information about the regions.
	RegionInfos []*RegionInfo `json:"RegionInfo,omitempty" xml:"RegionInfo,omitempty" type:"Repeated"`
}

func (s DescribeRegionsResponseBodyRegionInfoList) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponseBodyRegionInfoList) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponseBodyRegionInfoList) SetRegionInfos(v []*RegionInfo) *DescribeRegionsResponseBodyRegionInfoList {
	s.RegionInfos = v
	return s
}

type DescribeRegionsResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *DescribeRegionsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s DescribeRegionsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeRegionsResponse) GoString() string {
	return s.String()
}

func (s *DescribeRegionsResponse) SetHeaders(v map[string]*string) *DescribeRegionsResponse {
	s.Headers = v
	return s
}

func (s *DescribeRegionsResponse) SetStatusCode(v int32) *DescribeRegionsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeRegionsResponse) SetBody(v *DescribeRegionsResponseBody) *DescribeRegionsResponse {
	s.Body = v
	return s
}

type DoMetaQueryRequest struct {
	// The containerfor query conditions.
	MetaQuery *MetaQuery `json:"MetaQuery,omitempty" xml:"MetaQuery,omitempty"`
}

func (s DoMetaQueryRequest) String() string {
	return tea.Prettify(s)
}

func (s DoMetaQueryRequest) GoString() string {
	return s.String()
}

func (s *DoMetaQueryRequest) SetMetaQuery(v *MetaQuery) *DoMetaQueryRequest {
	s.MetaQuery = v
	return s
}

type DoMetaQueryResponseBody struct {
	// The container for the query result.
	MetaQuery *DoMetaQueryResponseBodyMetaQuery `json:"MetaQuery,omitempty" xml:"MetaQuery,omitempty" type:"Struct"`
}

func (s DoMetaQueryResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DoMetaQueryResponseBody) GoString() string {
	return s.String()
}

func (s *DoMetaQueryResponseBody) SetMetaQuery(v *DoMetaQueryResponseBodyMetaQuery) *DoMetaQueryResponseBody {
	s.MetaQuery = v
	return s
}

type DoMetaQueryResponseBodyMetaQuery struct {
	// The container for the information about objects.
	Files *DoMetaQueryResponseBodyMetaQueryFiles `json:"Files,omitempty" xml:"Files,omitempty" type:"Struct"`
	// The token that is used for the next query when the total number of objects exceeds the value of MaxResults.
	// The value of NextToken is used to return the unreturned results in the next query.
	//
	// This parameter has a value only when not all objects are returned.
	NextToken *string `json:"NextToken,omitempty" xml:"NextToken,omitempty"`
}

func (s DoMetaQueryResponseBodyMetaQuery) String() string {
	return tea.Prettify(s)
}

func (s DoMetaQueryResponseBodyMetaQuery) GoString() string {
	return s.String()
}

func (s *DoMetaQueryResponseBodyMetaQuery) SetFiles(v *DoMetaQueryResponseBodyMetaQueryFiles) *DoMetaQueryResponseBodyMetaQuery {
	s.Files = v
	return s
}

func (s *DoMetaQueryResponseBodyMetaQuery) SetNextToken(v string) *DoMetaQueryResponseBodyMetaQuery {
	s.NextToken = &v
	return s
}

type DoMetaQueryResponseBodyMetaQueryFiles struct {
	// The list of file information.
	File []*MetaQueryFile `json:"File,omitempty" xml:"File,omitempty" type:"Repeated"`
}

func (s DoMetaQueryResponseBodyMetaQueryFiles) String() string {
	return tea.Prettify(s)
}

func (s DoMetaQueryResponseBodyMetaQueryFiles) GoString() string {
	return s.String()
}

func (s *DoMetaQueryResponseBodyMetaQueryFiles) SetFile(v []*MetaQueryFile) *DoMetaQueryResponseBodyMetaQueryFiles {
	s.File = v
	return s
}

type DoMetaQueryResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *DoMetaQueryResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s DoMetaQueryResponse) String() string {
	return tea.Prettify(s)
}

func (s DoMetaQueryResponse) GoString() string {
	return s.String()
}

func (s *DoMetaQueryResponse) SetHeaders(v map[string]*string) *DoMetaQueryResponse {
	s.Headers = v
	return s
}

func (s *DoMetaQueryResponse) SetStatusCode(v int32) *DoMetaQueryResponse {
	s.StatusCode = &v
	return s
}

func (s *DoMetaQueryResponse) SetBody(v *DoMetaQueryResponseBody) *DoMetaQueryResponse {
	s.Body = v
	return s
}

type ExtendBucketWormRequest struct {
	// The parameters for WORM extension.
	ExtendWormConfiguration *ExtendWormConfiguration `json:"ExtendWormConfiguration,omitempty" xml:"ExtendWormConfiguration,omitempty"`
	// The ID of the retention policy.
	//
	// >  If the ID of the retention policy that specifies the number of days for which objects can be retained does not exist, the HTTP status code 404 is returned.
	WormId *string `json:"wormId,omitempty" xml:"wormId,omitempty"`
}

func (s ExtendBucketWormRequest) String() string {
	return tea.Prettify(s)
}

func (s ExtendBucketWormRequest) GoString() string {
	return s.String()
}

func (s *ExtendBucketWormRequest) SetExtendWormConfiguration(v *ExtendWormConfiguration) *ExtendBucketWormRequest {
	s.ExtendWormConfiguration = v
	return s
}

func (s *ExtendBucketWormRequest) SetWormId(v string) *ExtendBucketWormRequest {
	s.WormId = &v
	return s
}

type ExtendBucketWormResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s ExtendBucketWormResponse) String() string {
	return tea.Prettify(s)
}

func (s ExtendBucketWormResponse) GoString() string {
	return s.String()
}

func (s *ExtendBucketWormResponse) SetHeaders(v map[string]*string) *ExtendBucketWormResponse {
	s.Headers = v
	return s
}

func (s *ExtendBucketWormResponse) SetStatusCode(v int32) *ExtendBucketWormResponse {
	s.StatusCode = &v
	return s
}

type GetAccessPointHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The name of the access point.
	XOssAccessPointName *string `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s GetAccessPointHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointHeaders) GoString() string {
	return s.String()
}

func (s *GetAccessPointHeaders) SetCommonHeaders(v map[string]*string) *GetAccessPointHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAccessPointHeaders) SetXOssAccessPointName(v string) *GetAccessPointHeaders {
	s.XOssAccessPointName = &v
	return s
}

type GetAccessPointResponseBody struct {
	// The container that stores the information about the access point.
	GetAccessPointResult *GetAccessPointResult `json:"GetAccessPointResult,omitempty" xml:"GetAccessPointResult,omitempty"`
}

func (s GetAccessPointResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointResponseBody) GoString() string {
	return s.String()
}

func (s *GetAccessPointResponseBody) SetGetAccessPointResult(v *GetAccessPointResult) *GetAccessPointResponseBody {
	s.GetAccessPointResult = v
	return s
}

type GetAccessPointResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetAccessPointResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAccessPointResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointResponse) GoString() string {
	return s.String()
}

func (s *GetAccessPointResponse) SetHeaders(v map[string]*string) *GetAccessPointResponse {
	s.Headers = v
	return s
}

func (s *GetAccessPointResponse) SetStatusCode(v int32) *GetAccessPointResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAccessPointResponse) SetBody(v *GetAccessPointResponseBody) *GetAccessPointResponse {
	s.Body = v
	return s
}

type GetAccessPointConfigForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s GetAccessPointConfigForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointConfigForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *GetAccessPointConfigForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *GetAccessPointConfigForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAccessPointConfigForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *GetAccessPointConfigForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type GetAccessPointConfigForObjectProcessResponseBody struct {
	GetAccessPointConfigForObjectProcessResult *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult `json:"GetAccessPointConfigForObjectProcessResult,omitempty" xml:"GetAccessPointConfigForObjectProcessResult,omitempty" type:"Struct"`
}

func (s GetAccessPointConfigForObjectProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointConfigForObjectProcessResponseBody) GoString() string {
	return s.String()
}

func (s *GetAccessPointConfigForObjectProcessResponseBody) SetGetAccessPointConfigForObjectProcessResult(v *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult) *GetAccessPointConfigForObjectProcessResponseBody {
	s.GetAccessPointConfigForObjectProcessResult = v
	return s
}

type GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult struct {
	AllowAnonymousAccessForObjectProcess *string                         `json:"AllowAnonymousAccessForObjectProcess,omitempty" xml:"AllowAnonymousAccessForObjectProcess,omitempty"`
	ObjectProcessConfiguration           *ObjectProcessConfiguration     `json:"ObjectProcessConfiguration,omitempty" xml:"ObjectProcessConfiguration,omitempty"`
	PublicAccessBlockConfiguration       *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
}

func (s GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult) GoString() string {
	return s.String()
}

func (s *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult) SetAllowAnonymousAccessForObjectProcess(v string) *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult {
	s.AllowAnonymousAccessForObjectProcess = &v
	return s
}

func (s *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult) SetObjectProcessConfiguration(v *ObjectProcessConfiguration) *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult {
	s.ObjectProcessConfiguration = v
	return s
}

func (s *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *GetAccessPointConfigForObjectProcessResponseBodyGetAccessPointConfigForObjectProcessResult {
	s.PublicAccessBlockConfiguration = v
	return s
}

type GetAccessPointConfigForObjectProcessResponse struct {
	Headers    map[string]*string                                `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                            `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetAccessPointConfigForObjectProcessResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAccessPointConfigForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointConfigForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *GetAccessPointConfigForObjectProcessResponse) SetHeaders(v map[string]*string) *GetAccessPointConfigForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *GetAccessPointConfigForObjectProcessResponse) SetStatusCode(v int32) *GetAccessPointConfigForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAccessPointConfigForObjectProcessResponse) SetBody(v *GetAccessPointConfigForObjectProcessResponseBody) *GetAccessPointConfigForObjectProcessResponse {
	s.Body = v
	return s
}

type GetAccessPointForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s GetAccessPointForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *GetAccessPointForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *GetAccessPointForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAccessPointForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *GetAccessPointForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type GetAccessPointForObjectProcessResponseBody struct {
	GetAccessPointForObjectProcessResult *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult `json:"GetAccessPointForObjectProcessResult,omitempty" xml:"GetAccessPointForObjectProcessResult,omitempty" type:"Struct"`
}

func (s GetAccessPointForObjectProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointForObjectProcessResponseBody) GoString() string {
	return s.String()
}

func (s *GetAccessPointForObjectProcessResponseBody) SetGetAccessPointForObjectProcessResult(v *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) *GetAccessPointForObjectProcessResponseBody {
	s.GetAccessPointForObjectProcessResult = v
	return s
}

type GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult struct {
	AccessPointForObjectProcessAlias     *string                                                                                  `json:"AccessPointForObjectProcessAlias,omitempty" xml:"AccessPointForObjectProcessAlias,omitempty"`
	AccessPointForObjectProcessArn       *string                                                                                  `json:"AccessPointForObjectProcessArn,omitempty" xml:"AccessPointForObjectProcessArn,omitempty"`
	AccessPointName                      *string                                                                                  `json:"AccessPointName,omitempty" xml:"AccessPointName,omitempty"`
	AccessPointNameForObjectProcess      *string                                                                                  `json:"AccessPointNameForObjectProcess,omitempty" xml:"AccessPointNameForObjectProcess,omitempty"`
	AccountId                            *string                                                                                  `json:"AccountId,omitempty" xml:"AccountId,omitempty"`
	AllowAnonymousAccessForObjectProcess *string                                                                                  `json:"AllowAnonymousAccessForObjectProcess,omitempty" xml:"AllowAnonymousAccessForObjectProcess,omitempty"`
	CreationDate                         *string                                                                                  `json:"CreationDate,omitempty" xml:"CreationDate,omitempty"`
	Endpoints                            *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints `json:"Endpoints,omitempty" xml:"Endpoints,omitempty" type:"Struct"`
	PublicAccessBlockConfiguration       *PublicAccessBlockConfiguration                                                          `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
	Status                               *string                                                                                  `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) GoString() string {
	return s.String()
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetAccessPointForObjectProcessAlias(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.AccessPointForObjectProcessAlias = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetAccessPointForObjectProcessArn(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.AccessPointForObjectProcessArn = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetAccessPointName(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.AccessPointName = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetAccessPointNameForObjectProcess(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.AccessPointNameForObjectProcess = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetAccountId(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.AccountId = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetAllowAnonymousAccessForObjectProcess(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.AllowAnonymousAccessForObjectProcess = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetCreationDate(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.CreationDate = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetEndpoints(v *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.Endpoints = v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.PublicAccessBlockConfiguration = v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult) SetStatus(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResult {
	s.Status = &v
	return s
}

type GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints struct {
	InternalEndpoint *string `json:"InternalEndpoint,omitempty" xml:"InternalEndpoint,omitempty"`
	PublicEndpoint   *string `json:"PublicEndpoint,omitempty" xml:"PublicEndpoint,omitempty"`
}

func (s GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints) GoString() string {
	return s.String()
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints) SetInternalEndpoint(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints {
	s.InternalEndpoint = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints) SetPublicEndpoint(v string) *GetAccessPointForObjectProcessResponseBodyGetAccessPointForObjectProcessResultEndpoints {
	s.PublicEndpoint = &v
	return s
}

type GetAccessPointForObjectProcessResponse struct {
	Headers    map[string]*string                          `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                      `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetAccessPointForObjectProcessResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAccessPointForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *GetAccessPointForObjectProcessResponse) SetHeaders(v map[string]*string) *GetAccessPointForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *GetAccessPointForObjectProcessResponse) SetStatusCode(v int32) *GetAccessPointForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAccessPointForObjectProcessResponse) SetBody(v *GetAccessPointForObjectProcessResponseBody) *GetAccessPointForObjectProcessResponse {
	s.Body = v
	return s
}

type GetAccessPointPolicyHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The name of the access point.
	XOssAccessPointName *string `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s GetAccessPointPolicyHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointPolicyHeaders) GoString() string {
	return s.String()
}

func (s *GetAccessPointPolicyHeaders) SetCommonHeaders(v map[string]*string) *GetAccessPointPolicyHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAccessPointPolicyHeaders) SetXOssAccessPointName(v string) *GetAccessPointPolicyHeaders {
	s.XOssAccessPointName = &v
	return s
}

type GetAccessPointPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *string            `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAccessPointPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointPolicyResponse) GoString() string {
	return s.String()
}

func (s *GetAccessPointPolicyResponse) SetHeaders(v map[string]*string) *GetAccessPointPolicyResponse {
	s.Headers = v
	return s
}

func (s *GetAccessPointPolicyResponse) SetStatusCode(v int32) *GetAccessPointPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAccessPointPolicyResponse) SetBody(v string) *GetAccessPointPolicyResponse {
	s.Body = &v
	return s
}

type GetAccessPointPolicyForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s GetAccessPointPolicyForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointPolicyForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *GetAccessPointPolicyForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *GetAccessPointPolicyForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetAccessPointPolicyForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *GetAccessPointPolicyForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type GetAccessPointPolicyForObjectProcessResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *string            `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAccessPointPolicyForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointPolicyForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *GetAccessPointPolicyForObjectProcessResponse) SetHeaders(v map[string]*string) *GetAccessPointPolicyForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *GetAccessPointPolicyForObjectProcessResponse) SetStatusCode(v int32) *GetAccessPointPolicyForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAccessPointPolicyForObjectProcessResponse) SetBody(v string) *GetAccessPointPolicyForObjectProcessResponse {
	s.Body = &v
	return s
}

type GetAccessPointPublicAccessBlockRequest struct {
	XOssAccessPointName *string `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s GetAccessPointPublicAccessBlockRequest) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointPublicAccessBlockRequest) GoString() string {
	return s.String()
}

func (s *GetAccessPointPublicAccessBlockRequest) SetXOssAccessPointName(v string) *GetAccessPointPublicAccessBlockRequest {
	s.XOssAccessPointName = &v
	return s
}

type GetAccessPointPublicAccessBlockResponseBody struct {
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
}

func (s GetAccessPointPublicAccessBlockResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointPublicAccessBlockResponseBody) GoString() string {
	return s.String()
}

func (s *GetAccessPointPublicAccessBlockResponseBody) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *GetAccessPointPublicAccessBlockResponseBody {
	s.PublicAccessBlockConfiguration = v
	return s
}

type GetAccessPointPublicAccessBlockResponse struct {
	Headers    map[string]*string                           `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                       `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetAccessPointPublicAccessBlockResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAccessPointPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAccessPointPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *GetAccessPointPublicAccessBlockResponse) SetHeaders(v map[string]*string) *GetAccessPointPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *GetAccessPointPublicAccessBlockResponse) SetStatusCode(v int32) *GetAccessPointPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAccessPointPublicAccessBlockResponse) SetBody(v *GetAccessPointPublicAccessBlockResponseBody) *GetAccessPointPublicAccessBlockResponse {
	s.Body = v
	return s
}

type GetBucketAccessMonitorResponseBody struct {
	// The container that stores access monitor configuration.
	AccessMonitorConfiguration *GetBucketAccessMonitorResponseBodyAccessMonitorConfiguration `json:"AccessMonitorConfiguration,omitempty" xml:"AccessMonitorConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketAccessMonitorResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketAccessMonitorResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketAccessMonitorResponseBody) SetAccessMonitorConfiguration(v *GetBucketAccessMonitorResponseBodyAccessMonitorConfiguration) *GetBucketAccessMonitorResponseBody {
	s.AccessMonitorConfiguration = v
	return s
}

type GetBucketAccessMonitorResponseBodyAccessMonitorConfiguration struct {
	// The access tracking status of the bucket. Valid values:
	//
	// - Enabled: Access tracking is enabled.
	//
	// - Disabled: Access tracking is disabled.
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s GetBucketAccessMonitorResponseBodyAccessMonitorConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketAccessMonitorResponseBodyAccessMonitorConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketAccessMonitorResponseBodyAccessMonitorConfiguration) SetStatus(v string) *GetBucketAccessMonitorResponseBodyAccessMonitorConfiguration {
	s.Status = &v
	return s
}

type GetBucketAccessMonitorResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketAccessMonitorResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketAccessMonitorResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketAccessMonitorResponse) GoString() string {
	return s.String()
}

func (s *GetBucketAccessMonitorResponse) SetHeaders(v map[string]*string) *GetBucketAccessMonitorResponse {
	s.Headers = v
	return s
}

func (s *GetBucketAccessMonitorResponse) SetStatusCode(v int32) *GetBucketAccessMonitorResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketAccessMonitorResponse) SetBody(v *GetBucketAccessMonitorResponseBody) *GetBucketAccessMonitorResponse {
	s.Body = v
	return s
}

type GetBucketAclResponseBody struct {
	// The container that stores the ACL information.
	AccessControlPolicy *GetBucketAclResponseBodyAccessControlPolicy `json:"AccessControlPolicy,omitempty" xml:"AccessControlPolicy,omitempty" type:"Struct"`
}

func (s GetBucketAclResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketAclResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketAclResponseBody) SetAccessControlPolicy(v *GetBucketAclResponseBodyAccessControlPolicy) *GetBucketAclResponseBody {
	s.AccessControlPolicy = v
	return s
}

type GetBucketAclResponseBodyAccessControlPolicy struct {
	// The class of the container that stores the ACL information.
	AccessControlList *GetBucketAclResponseBodyAccessControlPolicyAccessControlList `json:"AccessControlList,omitempty" xml:"AccessControlList,omitempty" type:"Struct"`
	// The container that stores the information about the bucket owner.
	Owner *Owner `json:"Owner,omitempty" xml:"Owner,omitempty"`
}

func (s GetBucketAclResponseBodyAccessControlPolicy) String() string {
	return tea.Prettify(s)
}

func (s GetBucketAclResponseBodyAccessControlPolicy) GoString() string {
	return s.String()
}

func (s *GetBucketAclResponseBodyAccessControlPolicy) SetAccessControlList(v *GetBucketAclResponseBodyAccessControlPolicyAccessControlList) *GetBucketAclResponseBodyAccessControlPolicy {
	s.AccessControlList = v
	return s
}

func (s *GetBucketAclResponseBodyAccessControlPolicy) SetOwner(v *Owner) *GetBucketAclResponseBodyAccessControlPolicy {
	s.Owner = v
	return s
}

type GetBucketAclResponseBodyAccessControlPolicyAccessControlList struct {
	// The ACL of the bucket.
	Grant *string `json:"Grant,omitempty" xml:"Grant,omitempty"`
}

func (s GetBucketAclResponseBodyAccessControlPolicyAccessControlList) String() string {
	return tea.Prettify(s)
}

func (s GetBucketAclResponseBodyAccessControlPolicyAccessControlList) GoString() string {
	return s.String()
}

func (s *GetBucketAclResponseBodyAccessControlPolicyAccessControlList) SetGrant(v string) *GetBucketAclResponseBodyAccessControlPolicyAccessControlList {
	s.Grant = &v
	return s
}

type GetBucketAclResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketAclResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketAclResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketAclResponse) GoString() string {
	return s.String()
}

func (s *GetBucketAclResponse) SetHeaders(v map[string]*string) *GetBucketAclResponse {
	s.Headers = v
	return s
}

func (s *GetBucketAclResponse) SetStatusCode(v int32) *GetBucketAclResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketAclResponse) SetBody(v *GetBucketAclResponseBody) *GetBucketAclResponse {
	s.Body = v
	return s
}

type GetBucketArchiveDirectReadResponseBody struct {
	ArchiveDirectReadConfiguration *ArchiveDirectReadConfiguration `json:"ArchiveDirectReadConfiguration,omitempty" xml:"ArchiveDirectReadConfiguration,omitempty"`
}

func (s GetBucketArchiveDirectReadResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketArchiveDirectReadResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketArchiveDirectReadResponseBody) SetArchiveDirectReadConfiguration(v *ArchiveDirectReadConfiguration) *GetBucketArchiveDirectReadResponseBody {
	s.ArchiveDirectReadConfiguration = v
	return s
}

type GetBucketArchiveDirectReadResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketArchiveDirectReadResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketArchiveDirectReadResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketArchiveDirectReadResponse) GoString() string {
	return s.String()
}

func (s *GetBucketArchiveDirectReadResponse) SetHeaders(v map[string]*string) *GetBucketArchiveDirectReadResponse {
	s.Headers = v
	return s
}

func (s *GetBucketArchiveDirectReadResponse) SetStatusCode(v int32) *GetBucketArchiveDirectReadResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketArchiveDirectReadResponse) SetBody(v *GetBucketArchiveDirectReadResponseBody) *GetBucketArchiveDirectReadResponse {
	s.Body = v
	return s
}

type GetBucketCallbackPolicyResponseBody struct {
	BucketCallbackPolicy *CallbackPolicy `json:"BucketCallbackPolicy,omitempty" xml:"BucketCallbackPolicy,omitempty"`
}

func (s GetBucketCallbackPolicyResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketCallbackPolicyResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketCallbackPolicyResponseBody) SetBucketCallbackPolicy(v *CallbackPolicy) *GetBucketCallbackPolicyResponseBody {
	s.BucketCallbackPolicy = v
	return s
}

type GetBucketCallbackPolicyResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketCallbackPolicyResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketCallbackPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketCallbackPolicyResponse) GoString() string {
	return s.String()
}

func (s *GetBucketCallbackPolicyResponse) SetHeaders(v map[string]*string) *GetBucketCallbackPolicyResponse {
	s.Headers = v
	return s
}

func (s *GetBucketCallbackPolicyResponse) SetStatusCode(v int32) *GetBucketCallbackPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketCallbackPolicyResponse) SetBody(v *GetBucketCallbackPolicyResponseBody) *GetBucketCallbackPolicyResponse {
	s.Body = v
	return s
}

type GetBucketCorsResponseBody struct {
	// The container that stores CORS configuration.
	CORSConfiguration *GetBucketCorsResponseBodyCORSConfiguration `json:"CORSConfiguration,omitempty" xml:"CORSConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketCorsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketCorsResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketCorsResponseBody) SetCORSConfiguration(v *GetBucketCorsResponseBodyCORSConfiguration) *GetBucketCorsResponseBody {
	s.CORSConfiguration = v
	return s
}

type GetBucketCorsResponseBodyCORSConfiguration struct {
	// The container that stores CORS rules. Up to 10 rules can be configured for a bucket.
	CORSRule []*CORSRule `json:"CORSRule,omitempty" xml:"CORSRule,omitempty" type:"Repeated"`
	// Indicates whether the Vary: Origin header was returned. Default value: false.
	// - true: The Vary: Origin header is returned regardless whether the request is a cross-origin request or whether the cross-origin request succeeds.
	// - false: The Vary: Origin header is not returned.
	ResponseVary *bool `json:"ResponseVary,omitempty" xml:"ResponseVary,omitempty"`
}

func (s GetBucketCorsResponseBodyCORSConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketCorsResponseBodyCORSConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketCorsResponseBodyCORSConfiguration) SetCORSRule(v []*CORSRule) *GetBucketCorsResponseBodyCORSConfiguration {
	s.CORSRule = v
	return s
}

func (s *GetBucketCorsResponseBodyCORSConfiguration) SetResponseVary(v bool) *GetBucketCorsResponseBodyCORSConfiguration {
	s.ResponseVary = &v
	return s
}

type GetBucketCorsResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketCorsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketCorsResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketCorsResponse) GoString() string {
	return s.String()
}

func (s *GetBucketCorsResponse) SetHeaders(v map[string]*string) *GetBucketCorsResponse {
	s.Headers = v
	return s
}

func (s *GetBucketCorsResponse) SetStatusCode(v int32) *GetBucketCorsResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketCorsResponse) SetBody(v *GetBucketCorsResponseBody) *GetBucketCorsResponse {
	s.Body = v
	return s
}

type GetBucketDataRedundancyTransitionRequest struct {
	XOssRedundancyTransitionTaskid *string `json:"x-oss-redundancy-transition-taskid,omitempty" xml:"x-oss-redundancy-transition-taskid,omitempty"`
}

func (s GetBucketDataRedundancyTransitionRequest) String() string {
	return tea.Prettify(s)
}

func (s GetBucketDataRedundancyTransitionRequest) GoString() string {
	return s.String()
}

func (s *GetBucketDataRedundancyTransitionRequest) SetXOssRedundancyTransitionTaskid(v string) *GetBucketDataRedundancyTransitionRequest {
	s.XOssRedundancyTransitionTaskid = &v
	return s
}

type GetBucketDataRedundancyTransitionResponseBody struct {
	BucketDataRedundancyTransition *BucketDataRedundancyTransition `json:"BucketDataRedundancyTransition,omitempty" xml:"BucketDataRedundancyTransition,omitempty"`
}

func (s GetBucketDataRedundancyTransitionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketDataRedundancyTransitionResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketDataRedundancyTransitionResponseBody) SetBucketDataRedundancyTransition(v *BucketDataRedundancyTransition) *GetBucketDataRedundancyTransitionResponseBody {
	s.BucketDataRedundancyTransition = v
	return s
}

type GetBucketDataRedundancyTransitionResponse struct {
	Headers    map[string]*string                             `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                         `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketDataRedundancyTransitionResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketDataRedundancyTransitionResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketDataRedundancyTransitionResponse) GoString() string {
	return s.String()
}

func (s *GetBucketDataRedundancyTransitionResponse) SetHeaders(v map[string]*string) *GetBucketDataRedundancyTransitionResponse {
	s.Headers = v
	return s
}

func (s *GetBucketDataRedundancyTransitionResponse) SetStatusCode(v int32) *GetBucketDataRedundancyTransitionResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketDataRedundancyTransitionResponse) SetBody(v *GetBucketDataRedundancyTransitionResponseBody) *GetBucketDataRedundancyTransitionResponse {
	s.Body = v
	return s
}

type GetBucketEncryptionResponseBody struct {
	// The container that stores server-side encryption rules.
	ServerSideEncryptionRule *GetBucketEncryptionResponseBodyServerSideEncryptionRule `json:"ServerSideEncryptionRule,omitempty" xml:"ServerSideEncryptionRule,omitempty" type:"Struct"`
}

func (s GetBucketEncryptionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketEncryptionResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketEncryptionResponseBody) SetServerSideEncryptionRule(v *GetBucketEncryptionResponseBodyServerSideEncryptionRule) *GetBucketEncryptionResponseBody {
	s.ServerSideEncryptionRule = v
	return s
}

type GetBucketEncryptionResponseBodyServerSideEncryptionRule struct {
	// The container that stores the default server-side encryption method.
	ApplyServerSideEncryptionByDefault *ApplyServerSideEncryptionByDefault `json:"ApplyServerSideEncryptionByDefault,omitempty" xml:"ApplyServerSideEncryptionByDefault,omitempty"`
}

func (s GetBucketEncryptionResponseBodyServerSideEncryptionRule) String() string {
	return tea.Prettify(s)
}

func (s GetBucketEncryptionResponseBodyServerSideEncryptionRule) GoString() string {
	return s.String()
}

func (s *GetBucketEncryptionResponseBodyServerSideEncryptionRule) SetApplyServerSideEncryptionByDefault(v *ApplyServerSideEncryptionByDefault) *GetBucketEncryptionResponseBodyServerSideEncryptionRule {
	s.ApplyServerSideEncryptionByDefault = v
	return s
}

type GetBucketEncryptionResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketEncryptionResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketEncryptionResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketEncryptionResponse) GoString() string {
	return s.String()
}

func (s *GetBucketEncryptionResponse) SetHeaders(v map[string]*string) *GetBucketEncryptionResponse {
	s.Headers = v
	return s
}

func (s *GetBucketEncryptionResponse) SetStatusCode(v int32) *GetBucketEncryptionResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketEncryptionResponse) SetBody(v *GetBucketEncryptionResponseBody) *GetBucketEncryptionResponse {
	s.Body = v
	return s
}

type GetBucketHttpsConfigResponseBody struct {
	// The container that stores HTTPS configurations.
	HttpsConfiguration *HttpsConfiguration `json:"HttpsConfiguration,omitempty" xml:"HttpsConfiguration,omitempty"`
}

func (s GetBucketHttpsConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketHttpsConfigResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketHttpsConfigResponseBody) SetHttpsConfiguration(v *HttpsConfiguration) *GetBucketHttpsConfigResponseBody {
	s.HttpsConfiguration = v
	return s
}

type GetBucketHttpsConfigResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketHttpsConfigResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketHttpsConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketHttpsConfigResponse) GoString() string {
	return s.String()
}

func (s *GetBucketHttpsConfigResponse) SetHeaders(v map[string]*string) *GetBucketHttpsConfigResponse {
	s.Headers = v
	return s
}

func (s *GetBucketHttpsConfigResponse) SetStatusCode(v int32) *GetBucketHttpsConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketHttpsConfigResponse) SetBody(v *GetBucketHttpsConfigResponseBody) *GetBucketHttpsConfigResponse {
	s.Body = v
	return s
}

type GetBucketInfoResponseBody struct {
	// The container that stores the information about the bucket.
	BucketInfo *BucketInfo `json:"BucketInfo,omitempty" xml:"BucketInfo,omitempty"`
}

func (s GetBucketInfoResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketInfoResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketInfoResponseBody) SetBucketInfo(v *BucketInfo) *GetBucketInfoResponseBody {
	s.BucketInfo = v
	return s
}

type GetBucketInfoResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketInfoResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketInfoResponse) GoString() string {
	return s.String()
}

func (s *GetBucketInfoResponse) SetHeaders(v map[string]*string) *GetBucketInfoResponse {
	s.Headers = v
	return s
}

func (s *GetBucketInfoResponse) SetStatusCode(v int32) *GetBucketInfoResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketInfoResponse) SetBody(v *GetBucketInfoResponseBody) *GetBucketInfoResponse {
	s.Body = v
	return s
}

type GetBucketInventoryRequest struct {
	// The name of the inventory to be queried.
	InventoryId *string `json:"inventoryId,omitempty" xml:"inventoryId,omitempty"`
}

func (s GetBucketInventoryRequest) String() string {
	return tea.Prettify(s)
}

func (s GetBucketInventoryRequest) GoString() string {
	return s.String()
}

func (s *GetBucketInventoryRequest) SetInventoryId(v string) *GetBucketInventoryRequest {
	s.InventoryId = &v
	return s
}

type GetBucketInventoryResponseBody struct {
	// The inventory task configured for a bucket.
	InventoryConfiguration *InventoryConfiguration `json:"InventoryConfiguration,omitempty" xml:"InventoryConfiguration,omitempty"`
}

func (s GetBucketInventoryResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketInventoryResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketInventoryResponseBody) SetInventoryConfiguration(v *InventoryConfiguration) *GetBucketInventoryResponseBody {
	s.InventoryConfiguration = v
	return s
}

type GetBucketInventoryResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketInventoryResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketInventoryResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketInventoryResponse) GoString() string {
	return s.String()
}

func (s *GetBucketInventoryResponse) SetHeaders(v map[string]*string) *GetBucketInventoryResponse {
	s.Headers = v
	return s
}

func (s *GetBucketInventoryResponse) SetStatusCode(v int32) *GetBucketInventoryResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketInventoryResponse) SetBody(v *GetBucketInventoryResponseBody) *GetBucketInventoryResponse {
	s.Body = v
	return s
}

type GetBucketLifecycleResponseBody struct {
	// The container that stores the lifecycle rules configured for the bucket.
	LifecycleConfiguration *GetBucketLifecycleResponseBodyLifecycleConfiguration `json:"LifecycleConfiguration,omitempty" xml:"LifecycleConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketLifecycleResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLifecycleResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketLifecycleResponseBody) SetLifecycleConfiguration(v *GetBucketLifecycleResponseBodyLifecycleConfiguration) *GetBucketLifecycleResponseBody {
	s.LifecycleConfiguration = v
	return s
}

type GetBucketLifecycleResponseBodyLifecycleConfiguration struct {
	// The container that stores the lifecycle rules.
	Rule []*LifecycleRule `json:"Rule,omitempty" xml:"Rule,omitempty" type:"Repeated"`
}

func (s GetBucketLifecycleResponseBodyLifecycleConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLifecycleResponseBodyLifecycleConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketLifecycleResponseBodyLifecycleConfiguration) SetRule(v []*LifecycleRule) *GetBucketLifecycleResponseBodyLifecycleConfiguration {
	s.Rule = v
	return s
}

type GetBucketLifecycleResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketLifecycleResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketLifecycleResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLifecycleResponse) GoString() string {
	return s.String()
}

func (s *GetBucketLifecycleResponse) SetHeaders(v map[string]*string) *GetBucketLifecycleResponse {
	s.Headers = v
	return s
}

func (s *GetBucketLifecycleResponse) SetStatusCode(v int32) *GetBucketLifecycleResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketLifecycleResponse) SetBody(v *GetBucketLifecycleResponseBody) *GetBucketLifecycleResponse {
	s.Body = v
	return s
}

type GetBucketLocationResponseBody struct {
	// The region in which the bucket resides.\
	// Examples: oss-cn-hangzhou, oss-cn-shanghai, oss-cn-qingdao, oss-cn-beijing, oss-cn-zhangjiakou, oss-cn-hongkong, oss-cn-shenzhen, oss-us-west-1, oss-us-east-1, and oss-ap-southeast-1.
	//
	// \
	// For more information about the regions in which buckets reside, see [Regions and endpoints](~~31837~~).
	LocationConstraint *string `json:"LocationConstraint,omitempty" xml:"LocationConstraint,omitempty"`
}

func (s GetBucketLocationResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLocationResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketLocationResponseBody) SetLocationConstraint(v string) *GetBucketLocationResponseBody {
	s.LocationConstraint = &v
	return s
}

type GetBucketLocationResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketLocationResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketLocationResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLocationResponse) GoString() string {
	return s.String()
}

func (s *GetBucketLocationResponse) SetHeaders(v map[string]*string) *GetBucketLocationResponse {
	s.Headers = v
	return s
}

func (s *GetBucketLocationResponse) SetStatusCode(v int32) *GetBucketLocationResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketLocationResponse) SetBody(v *GetBucketLocationResponseBody) *GetBucketLocationResponse {
	s.Body = v
	return s
}

type GetBucketLoggingResponseBody struct {
	// Indicates the container used to store access logging configuration of a bucket.
	BucketLoggingStatus *GetBucketLoggingResponseBodyBucketLoggingStatus `json:"BucketLoggingStatus,omitempty" xml:"BucketLoggingStatus,omitempty" type:"Struct"`
}

func (s GetBucketLoggingResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLoggingResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketLoggingResponseBody) SetBucketLoggingStatus(v *GetBucketLoggingResponseBodyBucketLoggingStatus) *GetBucketLoggingResponseBody {
	s.BucketLoggingStatus = v
	return s
}

type GetBucketLoggingResponseBodyBucketLoggingStatus struct {
	// Indicates the container used to store access logging information. This element is returned if it is enabled and is not returned if it is disabled.
	LoggingEnabled *LoggingEnabled `json:"LoggingEnabled,omitempty" xml:"LoggingEnabled,omitempty"`
}

func (s GetBucketLoggingResponseBodyBucketLoggingStatus) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLoggingResponseBodyBucketLoggingStatus) GoString() string {
	return s.String()
}

func (s *GetBucketLoggingResponseBodyBucketLoggingStatus) SetLoggingEnabled(v *LoggingEnabled) *GetBucketLoggingResponseBodyBucketLoggingStatus {
	s.LoggingEnabled = v
	return s
}

type GetBucketLoggingResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketLoggingResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketLoggingResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketLoggingResponse) GoString() string {
	return s.String()
}

func (s *GetBucketLoggingResponse) SetHeaders(v map[string]*string) *GetBucketLoggingResponse {
	s.Headers = v
	return s
}

func (s *GetBucketLoggingResponse) SetStatusCode(v int32) *GetBucketLoggingResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketLoggingResponse) SetBody(v *GetBucketLoggingResponseBody) *GetBucketLoggingResponse {
	s.Body = v
	return s
}

type GetBucketPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *string            `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketPolicyResponse) GoString() string {
	return s.String()
}

func (s *GetBucketPolicyResponse) SetHeaders(v map[string]*string) *GetBucketPolicyResponse {
	s.Headers = v
	return s
}

func (s *GetBucketPolicyResponse) SetStatusCode(v int32) *GetBucketPolicyResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketPolicyResponse) SetBody(v string) *GetBucketPolicyResponse {
	s.Body = &v
	return s
}

type GetBucketPolicyStatusResponseBody struct {
	PolicyStatus *GetBucketPolicyStatusResponseBodyPolicyStatus `json:"PolicyStatus,omitempty" xml:"PolicyStatus,omitempty" type:"Struct"`
}

func (s GetBucketPolicyStatusResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketPolicyStatusResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketPolicyStatusResponseBody) SetPolicyStatus(v *GetBucketPolicyStatusResponseBodyPolicyStatus) *GetBucketPolicyStatusResponseBody {
	s.PolicyStatus = v
	return s
}

type GetBucketPolicyStatusResponseBodyPolicyStatus struct {
	IsPublic *bool `json:"IsPublic,omitempty" xml:"IsPublic,omitempty"`
}

func (s GetBucketPolicyStatusResponseBodyPolicyStatus) String() string {
	return tea.Prettify(s)
}

func (s GetBucketPolicyStatusResponseBodyPolicyStatus) GoString() string {
	return s.String()
}

func (s *GetBucketPolicyStatusResponseBodyPolicyStatus) SetIsPublic(v bool) *GetBucketPolicyStatusResponseBodyPolicyStatus {
	s.IsPublic = &v
	return s
}

type GetBucketPolicyStatusResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketPolicyStatusResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketPolicyStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketPolicyStatusResponse) GoString() string {
	return s.String()
}

func (s *GetBucketPolicyStatusResponse) SetHeaders(v map[string]*string) *GetBucketPolicyStatusResponse {
	s.Headers = v
	return s
}

func (s *GetBucketPolicyStatusResponse) SetStatusCode(v int32) *GetBucketPolicyStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketPolicyStatusResponse) SetBody(v *GetBucketPolicyStatusResponseBody) *GetBucketPolicyStatusResponse {
	s.Body = v
	return s
}

type GetBucketPublicAccessBlockResponseBody struct {
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
}

func (s GetBucketPublicAccessBlockResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketPublicAccessBlockResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketPublicAccessBlockResponseBody) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *GetBucketPublicAccessBlockResponseBody {
	s.PublicAccessBlockConfiguration = v
	return s
}

type GetBucketPublicAccessBlockResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketPublicAccessBlockResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *GetBucketPublicAccessBlockResponse) SetHeaders(v map[string]*string) *GetBucketPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *GetBucketPublicAccessBlockResponse) SetStatusCode(v int32) *GetBucketPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketPublicAccessBlockResponse) SetBody(v *GetBucketPublicAccessBlockResponseBody) *GetBucketPublicAccessBlockResponse {
	s.Body = v
	return s
}

type GetBucketRefererResponseBody struct {
	// The container that stores the hotlink protection configurations.
	RefererConfiguration *RefererConfiguration `json:"RefererConfiguration,omitempty" xml:"RefererConfiguration,omitempty"`
}

func (s GetBucketRefererResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketRefererResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketRefererResponseBody) SetRefererConfiguration(v *RefererConfiguration) *GetBucketRefererResponseBody {
	s.RefererConfiguration = v
	return s
}

type GetBucketRefererResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketRefererResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketRefererResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketRefererResponse) GoString() string {
	return s.String()
}

func (s *GetBucketRefererResponse) SetHeaders(v map[string]*string) *GetBucketRefererResponse {
	s.Headers = v
	return s
}

func (s *GetBucketRefererResponse) SetStatusCode(v int32) *GetBucketRefererResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketRefererResponse) SetBody(v *GetBucketRefererResponseBody) *GetBucketRefererResponse {
	s.Body = v
	return s
}

type GetBucketReplicationResponseBody struct {
	// The container that stores data replication configurations.
	ReplicationConfiguration *GetBucketReplicationResponseBodyReplicationConfiguration `json:"ReplicationConfiguration,omitempty" xml:"ReplicationConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketReplicationResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationResponseBody) SetReplicationConfiguration(v *GetBucketReplicationResponseBodyReplicationConfiguration) *GetBucketReplicationResponseBody {
	s.ReplicationConfiguration = v
	return s
}

type GetBucketReplicationResponseBodyReplicationConfiguration struct {
	// The container that stores the data replication rules.
	Rule []*ReplicationRule `json:"Rule,omitempty" xml:"Rule,omitempty" type:"Repeated"`
}

func (s GetBucketReplicationResponseBodyReplicationConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationResponseBodyReplicationConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationResponseBodyReplicationConfiguration) SetRule(v []*ReplicationRule) *GetBucketReplicationResponseBodyReplicationConfiguration {
	s.Rule = v
	return s
}

type GetBucketReplicationResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketReplicationResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketReplicationResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationResponse) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationResponse) SetHeaders(v map[string]*string) *GetBucketReplicationResponse {
	s.Headers = v
	return s
}

func (s *GetBucketReplicationResponse) SetStatusCode(v int32) *GetBucketReplicationResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketReplicationResponse) SetBody(v *GetBucketReplicationResponseBody) *GetBucketReplicationResponse {
	s.Body = v
	return s
}

type GetBucketReplicationLocationResponseBody struct {
	// The container that stores the region in which the destination bucket can be located.
	ReplicationLocation *GetBucketReplicationLocationResponseBodyReplicationLocation `json:"ReplicationLocation,omitempty" xml:"ReplicationLocation,omitempty" type:"Struct"`
}

func (s GetBucketReplicationLocationResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationLocationResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationLocationResponseBody) SetReplicationLocation(v *GetBucketReplicationLocationResponseBodyReplicationLocation) *GetBucketReplicationLocationResponseBody {
	s.ReplicationLocation = v
	return s
}

type GetBucketReplicationLocationResponseBodyReplicationLocation struct {
	// The regions in which the destination bucket can be located.
	Location []*string `json:"Location,omitempty" xml:"Location,omitempty" type:"Repeated"`
	// The container that stores regions in which the RTC can be enabled.
	LocationRTCConstraint *GetBucketReplicationLocationResponseBodyReplicationLocationLocationRTCConstraint `json:"LocationRTCConstraint,omitempty" xml:"LocationRTCConstraint,omitempty" type:"Struct"`
	// The container that stores regions in which the destination bucket can be located with TransferType specified.
	LocationTransferTypeConstraint *GetBucketReplicationLocationResponseBodyReplicationLocationLocationTransferTypeConstraint `json:"LocationTransferTypeConstraint,omitempty" xml:"LocationTransferTypeConstraint,omitempty" type:"Struct"`
}

func (s GetBucketReplicationLocationResponseBodyReplicationLocation) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationLocationResponseBodyReplicationLocation) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationLocationResponseBodyReplicationLocation) SetLocation(v []*string) *GetBucketReplicationLocationResponseBodyReplicationLocation {
	s.Location = v
	return s
}

func (s *GetBucketReplicationLocationResponseBodyReplicationLocation) SetLocationRTCConstraint(v *GetBucketReplicationLocationResponseBodyReplicationLocationLocationRTCConstraint) *GetBucketReplicationLocationResponseBodyReplicationLocation {
	s.LocationRTCConstraint = v
	return s
}

func (s *GetBucketReplicationLocationResponseBodyReplicationLocation) SetLocationTransferTypeConstraint(v *GetBucketReplicationLocationResponseBodyReplicationLocationLocationTransferTypeConstraint) *GetBucketReplicationLocationResponseBodyReplicationLocation {
	s.LocationTransferTypeConstraint = v
	return s
}

type GetBucketReplicationLocationResponseBodyReplicationLocationLocationRTCConstraint struct {
	// The regions where RTC is supported.
	Location []*string `json:"Location,omitempty" xml:"Location,omitempty" type:"Repeated"`
}

func (s GetBucketReplicationLocationResponseBodyReplicationLocationLocationRTCConstraint) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationLocationResponseBodyReplicationLocationLocationRTCConstraint) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationLocationResponseBodyReplicationLocationLocationRTCConstraint) SetLocation(v []*string) *GetBucketReplicationLocationResponseBodyReplicationLocationLocationRTCConstraint {
	s.Location = v
	return s
}

type GetBucketReplicationLocationResponseBodyReplicationLocationLocationTransferTypeConstraint struct {
	// The container that stores regions in which the destination bucket can be located with the TransferType information.
	LocationTransferType []*LocationTransferType `json:"LocationTransferType,omitempty" xml:"LocationTransferType,omitempty" type:"Repeated"`
}

func (s GetBucketReplicationLocationResponseBodyReplicationLocationLocationTransferTypeConstraint) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationLocationResponseBodyReplicationLocationLocationTransferTypeConstraint) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationLocationResponseBodyReplicationLocationLocationTransferTypeConstraint) SetLocationTransferType(v []*LocationTransferType) *GetBucketReplicationLocationResponseBodyReplicationLocationLocationTransferTypeConstraint {
	s.LocationTransferType = v
	return s
}

type GetBucketReplicationLocationResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketReplicationLocationResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketReplicationLocationResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationLocationResponse) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationLocationResponse) SetHeaders(v map[string]*string) *GetBucketReplicationLocationResponse {
	s.Headers = v
	return s
}

func (s *GetBucketReplicationLocationResponse) SetStatusCode(v int32) *GetBucketReplicationLocationResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketReplicationLocationResponse) SetBody(v *GetBucketReplicationLocationResponseBody) *GetBucketReplicationLocationResponse {
	s.Body = v
	return s
}

type GetBucketReplicationProgressRequest struct {
	// The ID of the data replication rule. You can call the GetBucketReplication operation to query the ID.
	RuleId *string `json:"rule-id,omitempty" xml:"rule-id,omitempty"`
}

func (s GetBucketReplicationProgressRequest) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationProgressRequest) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationProgressRequest) SetRuleId(v string) *GetBucketReplicationProgressRequest {
	s.RuleId = &v
	return s
}

type GetBucketReplicationProgressResponseBody struct {
	// The container that is used to store the progress of data replication tasks.
	ReplicationProgress *GetBucketReplicationProgressResponseBodyReplicationProgress `json:"ReplicationProgress,omitempty" xml:"ReplicationProgress,omitempty" type:"Struct"`
}

func (s GetBucketReplicationProgressResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationProgressResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationProgressResponseBody) SetReplicationProgress(v *GetBucketReplicationProgressResponseBodyReplicationProgress) *GetBucketReplicationProgressResponseBody {
	s.ReplicationProgress = v
	return s
}

type GetBucketReplicationProgressResponseBodyReplicationProgress struct {
	// The container that stores the progress of the data replication task corresponding to each data replication rule.
	Rule []*ReplicationProgressRule `json:"Rule,omitempty" xml:"Rule,omitempty" type:"Repeated"`
}

func (s GetBucketReplicationProgressResponseBodyReplicationProgress) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationProgressResponseBodyReplicationProgress) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationProgressResponseBodyReplicationProgress) SetRule(v []*ReplicationProgressRule) *GetBucketReplicationProgressResponseBodyReplicationProgress {
	s.Rule = v
	return s
}

type GetBucketReplicationProgressResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketReplicationProgressResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketReplicationProgressResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketReplicationProgressResponse) GoString() string {
	return s.String()
}

func (s *GetBucketReplicationProgressResponse) SetHeaders(v map[string]*string) *GetBucketReplicationProgressResponse {
	s.Headers = v
	return s
}

func (s *GetBucketReplicationProgressResponse) SetStatusCode(v int32) *GetBucketReplicationProgressResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketReplicationProgressResponse) SetBody(v *GetBucketReplicationProgressResponseBody) *GetBucketReplicationProgressResponse {
	s.Body = v
	return s
}

type GetBucketRequestPaymentResponseBody struct {
	// Indicates the container for the payer.
	RequestPaymentConfiguration *GetBucketRequestPaymentResponseBodyRequestPaymentConfiguration `json:"RequestPaymentConfiguration,omitempty" xml:"RequestPaymentConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketRequestPaymentResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketRequestPaymentResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketRequestPaymentResponseBody) SetRequestPaymentConfiguration(v *GetBucketRequestPaymentResponseBodyRequestPaymentConfiguration) *GetBucketRequestPaymentResponseBody {
	s.RequestPaymentConfiguration = v
	return s
}

type GetBucketRequestPaymentResponseBodyRequestPaymentConfiguration struct {
	// Indicates who pays the download and request fees.
	Payer *string `json:"Payer,omitempty" xml:"Payer,omitempty"`
}

func (s GetBucketRequestPaymentResponseBodyRequestPaymentConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketRequestPaymentResponseBodyRequestPaymentConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketRequestPaymentResponseBodyRequestPaymentConfiguration) SetPayer(v string) *GetBucketRequestPaymentResponseBodyRequestPaymentConfiguration {
	s.Payer = &v
	return s
}

type GetBucketRequestPaymentResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketRequestPaymentResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketRequestPaymentResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketRequestPaymentResponse) GoString() string {
	return s.String()
}

func (s *GetBucketRequestPaymentResponse) SetHeaders(v map[string]*string) *GetBucketRequestPaymentResponse {
	s.Headers = v
	return s
}

func (s *GetBucketRequestPaymentResponse) SetStatusCode(v int32) *GetBucketRequestPaymentResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketRequestPaymentResponse) SetBody(v *GetBucketRequestPaymentResponseBody) *GetBucketRequestPaymentResponse {
	s.Body = v
	return s
}

type GetBucketResourceGroupResponseBody struct {
	// The container that stores the ID of the resource group.
	BucketResourceGroupConfiguration *GetBucketResourceGroupResponseBodyBucketResourceGroupConfiguration `json:"BucketResourceGroupConfiguration,omitempty" xml:"BucketResourceGroupConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketResourceGroupResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketResourceGroupResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketResourceGroupResponseBody) SetBucketResourceGroupConfiguration(v *GetBucketResourceGroupResponseBodyBucketResourceGroupConfiguration) *GetBucketResourceGroupResponseBody {
	s.BucketResourceGroupConfiguration = v
	return s
}

type GetBucketResourceGroupResponseBodyBucketResourceGroupConfiguration struct {
	// The ID of the resource group to which the bucket belongs.
	//
	// If this element is not specified, the bucket is moved to the default resource group.
	ResourceGroupId *string `json:"ResourceGroupId,omitempty" xml:"ResourceGroupId,omitempty"`
}

func (s GetBucketResourceGroupResponseBodyBucketResourceGroupConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketResourceGroupResponseBodyBucketResourceGroupConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketResourceGroupResponseBodyBucketResourceGroupConfiguration) SetResourceGroupId(v string) *GetBucketResourceGroupResponseBodyBucketResourceGroupConfiguration {
	s.ResourceGroupId = &v
	return s
}

type GetBucketResourceGroupResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketResourceGroupResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketResourceGroupResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketResourceGroupResponse) GoString() string {
	return s.String()
}

func (s *GetBucketResourceGroupResponse) SetHeaders(v map[string]*string) *GetBucketResourceGroupResponse {
	s.Headers = v
	return s
}

func (s *GetBucketResourceGroupResponse) SetStatusCode(v int32) *GetBucketResourceGroupResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketResourceGroupResponse) SetBody(v *GetBucketResourceGroupResponseBody) *GetBucketResourceGroupResponse {
	s.Body = v
	return s
}

type GetBucketResponseHeaderResponseBody struct {
	ResponseHeaderConfiguration *ResponseHeaderConfiguration `json:"ResponseHeaderConfiguration,omitempty" xml:"ResponseHeaderConfiguration,omitempty"`
}

func (s GetBucketResponseHeaderResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketResponseHeaderResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketResponseHeaderResponseBody) SetResponseHeaderConfiguration(v *ResponseHeaderConfiguration) *GetBucketResponseHeaderResponseBody {
	s.ResponseHeaderConfiguration = v
	return s
}

type GetBucketResponseHeaderResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketResponseHeaderResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketResponseHeaderResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketResponseHeaderResponse) GoString() string {
	return s.String()
}

func (s *GetBucketResponseHeaderResponse) SetHeaders(v map[string]*string) *GetBucketResponseHeaderResponse {
	s.Headers = v
	return s
}

func (s *GetBucketResponseHeaderResponse) SetStatusCode(v int32) *GetBucketResponseHeaderResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketResponseHeaderResponse) SetBody(v *GetBucketResponseHeaderResponseBody) *GetBucketResponseHeaderResponse {
	s.Body = v
	return s
}

type GetBucketStatResponseBody struct {
	// The container that stores all information returned for the GetBucketStat request.
	BucketStat *BucketStat `json:"BucketStat,omitempty" xml:"BucketStat,omitempty"`
}

func (s GetBucketStatResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketStatResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketStatResponseBody) SetBucketStat(v *BucketStat) *GetBucketStatResponseBody {
	s.BucketStat = v
	return s
}

type GetBucketStatResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketStatResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketStatResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketStatResponse) GoString() string {
	return s.String()
}

func (s *GetBucketStatResponse) SetHeaders(v map[string]*string) *GetBucketStatResponse {
	s.Headers = v
	return s
}

func (s *GetBucketStatResponse) SetStatusCode(v int32) *GetBucketStatResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketStatResponse) SetBody(v *GetBucketStatResponseBody) *GetBucketStatResponse {
	s.Body = v
	return s
}

type GetBucketTagsResponseBody struct {
	// The container that stores the returned tags of the bucket.
	// > If no tags are configured for the bucket, an XML message body is returned in which the Tagging element is empty.
	Tagging *GetBucketTagsResponseBodyTagging `json:"Tagging,omitempty" xml:"Tagging,omitempty" type:"Struct"`
}

func (s GetBucketTagsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketTagsResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketTagsResponseBody) SetTagging(v *GetBucketTagsResponseBodyTagging) *GetBucketTagsResponseBody {
	s.Tagging = v
	return s
}

type GetBucketTagsResponseBodyTagging struct {
	// The container that stores the returned tags of the bucket.
	TagSet *TagSet `json:"TagSet,omitempty" xml:"TagSet,omitempty"`
}

func (s GetBucketTagsResponseBodyTagging) String() string {
	return tea.Prettify(s)
}

func (s GetBucketTagsResponseBodyTagging) GoString() string {
	return s.String()
}

func (s *GetBucketTagsResponseBodyTagging) SetTagSet(v *TagSet) *GetBucketTagsResponseBodyTagging {
	s.TagSet = v
	return s
}

type GetBucketTagsResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketTagsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketTagsResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketTagsResponse) GoString() string {
	return s.String()
}

func (s *GetBucketTagsResponse) SetHeaders(v map[string]*string) *GetBucketTagsResponse {
	s.Headers = v
	return s
}

func (s *GetBucketTagsResponse) SetStatusCode(v int32) *GetBucketTagsResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketTagsResponse) SetBody(v *GetBucketTagsResponseBody) *GetBucketTagsResponse {
	s.Body = v
	return s
}

type GetBucketTransferAccelerationResponseBody struct {
	// The container that stores the transfer acceleration configurations.
	TransferAccelerationConfiguration *GetBucketTransferAccelerationResponseBodyTransferAccelerationConfiguration `json:"TransferAccelerationConfiguration,omitempty" xml:"TransferAccelerationConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketTransferAccelerationResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketTransferAccelerationResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketTransferAccelerationResponseBody) SetTransferAccelerationConfiguration(v *GetBucketTransferAccelerationResponseBodyTransferAccelerationConfiguration) *GetBucketTransferAccelerationResponseBody {
	s.TransferAccelerationConfiguration = v
	return s
}

type GetBucketTransferAccelerationResponseBodyTransferAccelerationConfiguration struct {
	// Whether the transfer acceleration is enabled for this bucket.
	Enabled *bool `json:"Enabled,omitempty" xml:"Enabled,omitempty"`
}

func (s GetBucketTransferAccelerationResponseBodyTransferAccelerationConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketTransferAccelerationResponseBodyTransferAccelerationConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketTransferAccelerationResponseBodyTransferAccelerationConfiguration) SetEnabled(v bool) *GetBucketTransferAccelerationResponseBodyTransferAccelerationConfiguration {
	s.Enabled = &v
	return s
}

type GetBucketTransferAccelerationResponse struct {
	Headers    map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketTransferAccelerationResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketTransferAccelerationResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketTransferAccelerationResponse) GoString() string {
	return s.String()
}

func (s *GetBucketTransferAccelerationResponse) SetHeaders(v map[string]*string) *GetBucketTransferAccelerationResponse {
	s.Headers = v
	return s
}

func (s *GetBucketTransferAccelerationResponse) SetStatusCode(v int32) *GetBucketTransferAccelerationResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketTransferAccelerationResponse) SetBody(v *GetBucketTransferAccelerationResponseBody) *GetBucketTransferAccelerationResponse {
	s.Body = v
	return s
}

type GetBucketVersioningResponseBody struct {
	// The container that stores the versioning state of the bucket.
	VersioningConfiguration *GetBucketVersioningResponseBodyVersioningConfiguration `json:"VersioningConfiguration,omitempty" xml:"VersioningConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketVersioningResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketVersioningResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketVersioningResponseBody) SetVersioningConfiguration(v *GetBucketVersioningResponseBodyVersioningConfiguration) *GetBucketVersioningResponseBody {
	s.VersioningConfiguration = v
	return s
}

type GetBucketVersioningResponseBodyVersioningConfiguration struct {
	// The versioning state of the bucket.
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s GetBucketVersioningResponseBodyVersioningConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketVersioningResponseBodyVersioningConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketVersioningResponseBodyVersioningConfiguration) SetStatus(v string) *GetBucketVersioningResponseBodyVersioningConfiguration {
	s.Status = &v
	return s
}

type GetBucketVersioningResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketVersioningResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketVersioningResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketVersioningResponse) GoString() string {
	return s.String()
}

func (s *GetBucketVersioningResponse) SetHeaders(v map[string]*string) *GetBucketVersioningResponse {
	s.Headers = v
	return s
}

func (s *GetBucketVersioningResponse) SetStatusCode(v int32) *GetBucketVersioningResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketVersioningResponse) SetBody(v *GetBucketVersioningResponseBody) *GetBucketVersioningResponse {
	s.Body = v
	return s
}

type GetBucketWebsiteResponseBody struct {
	// The containers of the website configuration.
	WebsiteConfiguration *WebsiteConfiguration `json:"WebsiteConfiguration,omitempty" xml:"WebsiteConfiguration,omitempty"`
}

func (s GetBucketWebsiteResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketWebsiteResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketWebsiteResponseBody) SetWebsiteConfiguration(v *WebsiteConfiguration) *GetBucketWebsiteResponseBody {
	s.WebsiteConfiguration = v
	return s
}

type GetBucketWebsiteResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketWebsiteResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketWebsiteResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketWebsiteResponse) GoString() string {
	return s.String()
}

func (s *GetBucketWebsiteResponse) SetHeaders(v map[string]*string) *GetBucketWebsiteResponse {
	s.Headers = v
	return s
}

func (s *GetBucketWebsiteResponse) SetStatusCode(v int32) *GetBucketWebsiteResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketWebsiteResponse) SetBody(v *GetBucketWebsiteResponseBody) *GetBucketWebsiteResponse {
	s.Body = v
	return s
}

type GetBucketWormResponseBody struct {
	// The container that stores the information about retention policies of the bucket.
	WormConfiguration *GetBucketWormResponseBodyWormConfiguration `json:"WormConfiguration,omitempty" xml:"WormConfiguration,omitempty" type:"Struct"`
}

func (s GetBucketWormResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetBucketWormResponseBody) GoString() string {
	return s.String()
}

func (s *GetBucketWormResponseBody) SetWormConfiguration(v *GetBucketWormResponseBodyWormConfiguration) *GetBucketWormResponseBody {
	s.WormConfiguration = v
	return s
}

type GetBucketWormResponseBodyWormConfiguration struct {
	// The time at which the retention policy was created.
	CreationDate *string `json:"CreationDate,omitempty" xml:"CreationDate,omitempty"`
	// The time at which the retention policy will be expired.
	ExpirationDate *string `json:"ExpirationDate,omitempty" xml:"ExpirationDate,omitempty"`
	// The number of days for which objects can be retained.
	RetentionPeriodInDays *int32 `json:"RetentionPeriodInDays,omitempty" xml:"RetentionPeriodInDays,omitempty"`
	// The status of the retention policy. Valid values:
	//
	// - InProgress: indicates that the retention policy is in the InProgress state. By default, a retention policy is in the InProgress state after it is created. The policy remains in this state for 24 hours.
	// - Locked: indicates that the retention policy is in the Locked state.
	State *string `json:"State,omitempty" xml:"State,omitempty"`
	// The ID of the retention policy.
	//
	// >Note If the specified retention policy ID that is used to query the retention policy configurations of the bucket does not exist, OSS returns the 404 error code.
	WormId *string `json:"WormId,omitempty" xml:"WormId,omitempty"`
}

func (s GetBucketWormResponseBodyWormConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetBucketWormResponseBodyWormConfiguration) GoString() string {
	return s.String()
}

func (s *GetBucketWormResponseBodyWormConfiguration) SetCreationDate(v string) *GetBucketWormResponseBodyWormConfiguration {
	s.CreationDate = &v
	return s
}

func (s *GetBucketWormResponseBodyWormConfiguration) SetExpirationDate(v string) *GetBucketWormResponseBodyWormConfiguration {
	s.ExpirationDate = &v
	return s
}

func (s *GetBucketWormResponseBodyWormConfiguration) SetRetentionPeriodInDays(v int32) *GetBucketWormResponseBodyWormConfiguration {
	s.RetentionPeriodInDays = &v
	return s
}

func (s *GetBucketWormResponseBodyWormConfiguration) SetState(v string) *GetBucketWormResponseBodyWormConfiguration {
	s.State = &v
	return s
}

func (s *GetBucketWormResponseBodyWormConfiguration) SetWormId(v string) *GetBucketWormResponseBodyWormConfiguration {
	s.WormId = &v
	return s
}

type GetBucketWormResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetBucketWormResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetBucketWormResponse) String() string {
	return tea.Prettify(s)
}

func (s GetBucketWormResponse) GoString() string {
	return s.String()
}

func (s *GetBucketWormResponse) SetHeaders(v map[string]*string) *GetBucketWormResponse {
	s.Headers = v
	return s
}

func (s *GetBucketWormResponse) SetStatusCode(v int32) *GetBucketWormResponse {
	s.StatusCode = &v
	return s
}

func (s *GetBucketWormResponse) SetBody(v *GetBucketWormResponseBody) *GetBucketWormResponse {
	s.Body = v
	return s
}

type GetCnameTokenRequest struct {
	// The name of the CNAME record that is mapped to the bucket.
	Cname *string `json:"cname,omitempty" xml:"cname,omitempty"`
}

func (s GetCnameTokenRequest) String() string {
	return tea.Prettify(s)
}

func (s GetCnameTokenRequest) GoString() string {
	return s.String()
}

func (s *GetCnameTokenRequest) SetCname(v string) *GetCnameTokenRequest {
	s.Cname = &v
	return s
}

type GetCnameTokenResponseBody struct {
	// The container in which the CNAME token is stored.
	CnameToken *CnameToken `json:"CnameToken,omitempty" xml:"CnameToken,omitempty"`
}

func (s GetCnameTokenResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetCnameTokenResponseBody) GoString() string {
	return s.String()
}

func (s *GetCnameTokenResponseBody) SetCnameToken(v *CnameToken) *GetCnameTokenResponseBody {
	s.CnameToken = v
	return s
}

type GetCnameTokenResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetCnameTokenResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetCnameTokenResponse) String() string {
	return tea.Prettify(s)
}

func (s GetCnameTokenResponse) GoString() string {
	return s.String()
}

func (s *GetCnameTokenResponse) SetHeaders(v map[string]*string) *GetCnameTokenResponse {
	s.Headers = v
	return s
}

func (s *GetCnameTokenResponse) SetStatusCode(v int32) *GetCnameTokenResponse {
	s.StatusCode = &v
	return s
}

func (s *GetCnameTokenResponse) SetBody(v *GetCnameTokenResponseBody) *GetCnameTokenResponse {
	s.Body = v
	return s
}

type GetLiveChannelHistoryResponseBody struct {
	// The container that stores the returned results of the GetLiveChannelHistory request.
	LiveChannelHistory *GetLiveChannelHistoryResponseBodyLiveChannelHistory `json:"LiveChannelHistory,omitempty" xml:"LiveChannelHistory,omitempty" type:"Struct"`
}

func (s GetLiveChannelHistoryResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelHistoryResponseBody) GoString() string {
	return s.String()
}

func (s *GetLiveChannelHistoryResponseBody) SetLiveChannelHistory(v *GetLiveChannelHistoryResponseBodyLiveChannelHistory) *GetLiveChannelHistoryResponseBody {
	s.LiveChannelHistory = v
	return s
}

type GetLiveChannelHistoryResponseBodyLiveChannelHistory struct {
	// The container that stores a list of stream pushing records.
	LiveRecord []*LiveRecord `json:"LiveRecord,omitempty" xml:"LiveRecord,omitempty" type:"Repeated"`
}

func (s GetLiveChannelHistoryResponseBodyLiveChannelHistory) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelHistoryResponseBodyLiveChannelHistory) GoString() string {
	return s.String()
}

func (s *GetLiveChannelHistoryResponseBodyLiveChannelHistory) SetLiveRecord(v []*LiveRecord) *GetLiveChannelHistoryResponseBodyLiveChannelHistory {
	s.LiveRecord = v
	return s
}

type GetLiveChannelHistoryResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetLiveChannelHistoryResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetLiveChannelHistoryResponse) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelHistoryResponse) GoString() string {
	return s.String()
}

func (s *GetLiveChannelHistoryResponse) SetHeaders(v map[string]*string) *GetLiveChannelHistoryResponse {
	s.Headers = v
	return s
}

func (s *GetLiveChannelHistoryResponse) SetStatusCode(v int32) *GetLiveChannelHistoryResponse {
	s.StatusCode = &v
	return s
}

func (s *GetLiveChannelHistoryResponse) SetBody(v *GetLiveChannelHistoryResponseBody) *GetLiveChannelHistoryResponse {
	s.Body = v
	return s
}

type GetLiveChannelInfoResponseBody struct {
	// The container that stores the returned results of the GetLiveChannelInfo request.
	LiveChannelConfiguration *GetLiveChannelInfoResponseBodyLiveChannelConfiguration `json:"LiveChannelConfiguration,omitempty" xml:"LiveChannelConfiguration,omitempty" type:"Struct"`
}

func (s GetLiveChannelInfoResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelInfoResponseBody) GoString() string {
	return s.String()
}

func (s *GetLiveChannelInfoResponseBody) SetLiveChannelConfiguration(v *GetLiveChannelInfoResponseBodyLiveChannelConfiguration) *GetLiveChannelInfoResponseBody {
	s.LiveChannelConfiguration = v
	return s
}

type GetLiveChannelInfoResponseBodyLiveChannelConfiguration struct {
	// The description of the LiveChannel.
	Description *string `json:"Description,omitempty" xml:"Description,omitempty"`
	// The status of the LiveChannel.
	//
	// Valid values:
	// - enabled: indicates that the LiveChannel is enabled.
	// - disabled: indicates that the LiveChannel is disabled.
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
	// The container that stores the configurations used by the LiveChannel to store uploaded data.
	// > FragDuration, FragCount, and PlaylistName are returned only when the value of Type is HLS.
	Target *LiveChannelTarget `json:"Target,omitempty" xml:"Target,omitempty"`
}

func (s GetLiveChannelInfoResponseBodyLiveChannelConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelInfoResponseBodyLiveChannelConfiguration) GoString() string {
	return s.String()
}

func (s *GetLiveChannelInfoResponseBodyLiveChannelConfiguration) SetDescription(v string) *GetLiveChannelInfoResponseBodyLiveChannelConfiguration {
	s.Description = &v
	return s
}

func (s *GetLiveChannelInfoResponseBodyLiveChannelConfiguration) SetStatus(v string) *GetLiveChannelInfoResponseBodyLiveChannelConfiguration {
	s.Status = &v
	return s
}

func (s *GetLiveChannelInfoResponseBodyLiveChannelConfiguration) SetTarget(v *LiveChannelTarget) *GetLiveChannelInfoResponseBodyLiveChannelConfiguration {
	s.Target = v
	return s
}

type GetLiveChannelInfoResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetLiveChannelInfoResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetLiveChannelInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelInfoResponse) GoString() string {
	return s.String()
}

func (s *GetLiveChannelInfoResponse) SetHeaders(v map[string]*string) *GetLiveChannelInfoResponse {
	s.Headers = v
	return s
}

func (s *GetLiveChannelInfoResponse) SetStatusCode(v int32) *GetLiveChannelInfoResponse {
	s.StatusCode = &v
	return s
}

func (s *GetLiveChannelInfoResponse) SetBody(v *GetLiveChannelInfoResponseBody) *GetLiveChannelInfoResponse {
	s.Body = v
	return s
}

type GetLiveChannelStatResponseBody struct {
	// The container that stores the returned results of the GetLiveChannelStat request.
	LiveChannelStat *GetLiveChannelStatResponseBodyLiveChannelStat `json:"LiveChannelStat,omitempty" xml:"LiveChannelStat,omitempty" type:"Struct"`
}

func (s GetLiveChannelStatResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelStatResponseBody) GoString() string {
	return s.String()
}

func (s *GetLiveChannelStatResponseBody) SetLiveChannelStat(v *GetLiveChannelStatResponseBodyLiveChannelStat) *GetLiveChannelStatResponseBody {
	s.LiveChannelStat = v
	return s
}

type GetLiveChannelStatResponseBodyLiveChannelStat struct {
	// The container that stores audio stream information if Status is set to Live.
	// >Video and audio containers can be returned only if Status is set to Live. However, these two containers may not necessarily be returned if Status is set to Live. For example, if the client has connected to the LiveChannel but no audio or video stream is sent, these two containers are not returned.
	Audio *LiveChannelAudio `json:"Audio,omitempty" xml:"Audio,omitempty"`
	// If Status is set to Live, this element indicates the time when the current client starts to ingest streams. The value of the element is in the ISO 8601 format.
	ConnectedTime *string `json:"ConnectedTime,omitempty" xml:"ConnectedTime,omitempty"`
	// If Status is set to Live, this element indicates the IP address of the current client that ingests streams.
	RemoteAddr *string `json:"RemoteAddr,omitempty" xml:"RemoteAddr,omitempty"`
	// The current stream ingestion status of the LiveChannel. Valid value: Disabled、Live、Idle。
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
	// The container that stores video stream information if Status is set to Live.
	//
	// >Video and audio containers can be returned only if Status is set to Live. However, these two containers may not necessarily be returned if Status is set to Live. For example, if the client has connected to the LiveChannel but no audio or video stream is sent, these two containers are not returned.
	Video *LiveChannelVideo `json:"Video,omitempty" xml:"Video,omitempty"`
}

func (s GetLiveChannelStatResponseBodyLiveChannelStat) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelStatResponseBodyLiveChannelStat) GoString() string {
	return s.String()
}

func (s *GetLiveChannelStatResponseBodyLiveChannelStat) SetAudio(v *LiveChannelAudio) *GetLiveChannelStatResponseBodyLiveChannelStat {
	s.Audio = v
	return s
}

func (s *GetLiveChannelStatResponseBodyLiveChannelStat) SetConnectedTime(v string) *GetLiveChannelStatResponseBodyLiveChannelStat {
	s.ConnectedTime = &v
	return s
}

func (s *GetLiveChannelStatResponseBodyLiveChannelStat) SetRemoteAddr(v string) *GetLiveChannelStatResponseBodyLiveChannelStat {
	s.RemoteAddr = &v
	return s
}

func (s *GetLiveChannelStatResponseBodyLiveChannelStat) SetStatus(v string) *GetLiveChannelStatResponseBodyLiveChannelStat {
	s.Status = &v
	return s
}

func (s *GetLiveChannelStatResponseBodyLiveChannelStat) SetVideo(v *LiveChannelVideo) *GetLiveChannelStatResponseBodyLiveChannelStat {
	s.Video = v
	return s
}

type GetLiveChannelStatResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetLiveChannelStatResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetLiveChannelStatResponse) String() string {
	return tea.Prettify(s)
}

func (s GetLiveChannelStatResponse) GoString() string {
	return s.String()
}

func (s *GetLiveChannelStatResponse) SetHeaders(v map[string]*string) *GetLiveChannelStatResponse {
	s.Headers = v
	return s
}

func (s *GetLiveChannelStatResponse) SetStatusCode(v int32) *GetLiveChannelStatResponse {
	s.StatusCode = &v
	return s
}

func (s *GetLiveChannelStatResponse) SetBody(v *GetLiveChannelStatResponseBody) *GetLiveChannelStatResponse {
	s.Body = v
	return s
}

type GetMetaQueryStatusResponseBody struct {
	// The container that stores the metadata information.
	MetaQueryStatus *GetMetaQueryStatusResponseBodyMetaQueryStatus `json:"MetaQueryStatus,omitempty" xml:"MetaQueryStatus,omitempty" type:"Struct"`
}

func (s GetMetaQueryStatusResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetMetaQueryStatusResponseBody) GoString() string {
	return s.String()
}

func (s *GetMetaQueryStatusResponseBody) SetMetaQueryStatus(v *GetMetaQueryStatusResponseBodyMetaQueryStatus) *GetMetaQueryStatusResponseBody {
	s.MetaQueryStatus = v
	return s
}

type GetMetaQueryStatusResponseBodyMetaQueryStatus struct {
	// The time when the metadata index library was created. The value follows the RFC 3339 standard in the YYYY-MM-DDTHH:mm:ss+TIMEZONE format. YYYY-MM-DD indicates the year, month, and day. T indicates the beginning of the time element. HH:mm:ss indicates the hour, minute, and second. TIMEZONE indicates the time zone.
	CreateTime *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	// The scan type. Valid values:
	// - FullScanning: Full scanning is in progress.
	// - IncrementalScanning: Incremental scanning is in progress.
	Phase *string `json:"Phase,omitempty" xml:"Phase,omitempty"`
	// The status of the metadata index library. Valid values:
	// - Ready: The metadata index library is being prepared after it is created.
	// In this case, the metadata index library cannot be used to query data.
	//
	// - Stop: The metadata index library is paused.
	// - Running: The metadata index library is running.
	// - Retrying: The metadata index library failed to be created and is being created again.
	// - Failed: The metadata index library failed to be created.
	// - Deleted: The metadata index library is deleted.
	State *string `json:"State,omitempty" xml:"State,omitempty"`
	// The time when the metadata index library was updated. The value follows the RFC 3339 standard in the YYYY-MM-DDTHH:mm:ss+TIMEZONE format. YYYY-MM-DD indicates the year, month, and day. T indicates the beginning of the time element. HH:mm:ss indicates the hour, minute, and second. TIMEZONE indicates the time zone.
	UpdateTime *string `json:"UpdateTime,omitempty" xml:"UpdateTime,omitempty"`
}

func (s GetMetaQueryStatusResponseBodyMetaQueryStatus) String() string {
	return tea.Prettify(s)
}

func (s GetMetaQueryStatusResponseBodyMetaQueryStatus) GoString() string {
	return s.String()
}

func (s *GetMetaQueryStatusResponseBodyMetaQueryStatus) SetCreateTime(v string) *GetMetaQueryStatusResponseBodyMetaQueryStatus {
	s.CreateTime = &v
	return s
}

func (s *GetMetaQueryStatusResponseBodyMetaQueryStatus) SetPhase(v string) *GetMetaQueryStatusResponseBodyMetaQueryStatus {
	s.Phase = &v
	return s
}

func (s *GetMetaQueryStatusResponseBodyMetaQueryStatus) SetState(v string) *GetMetaQueryStatusResponseBodyMetaQueryStatus {
	s.State = &v
	return s
}

func (s *GetMetaQueryStatusResponseBodyMetaQueryStatus) SetUpdateTime(v string) *GetMetaQueryStatusResponseBodyMetaQueryStatus {
	s.UpdateTime = &v
	return s
}

type GetMetaQueryStatusResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetMetaQueryStatusResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetMetaQueryStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s GetMetaQueryStatusResponse) GoString() string {
	return s.String()
}

func (s *GetMetaQueryStatusResponse) SetHeaders(v map[string]*string) *GetMetaQueryStatusResponse {
	s.Headers = v
	return s
}

func (s *GetMetaQueryStatusResponse) SetStatusCode(v int32) *GetMetaQueryStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *GetMetaQueryStatusResponse) SetBody(v *GetMetaQueryStatusResponseBody) *GetMetaQueryStatusResponse {
	s.Body = v
	return s
}

type GetObjectHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The encoding type at the client side.
	// If you want an object to be returned in the GZIP format, you must include the Accept-Encoding:gzip header in your request. OSS determines whether to return the object compressed in the GZip format based on the Content-Type header and whether the size of the object is larger than or equal to 1 KB.
	//
	// > If an object is compressed in the GZip format, the response OSS returns does not include the ETag value of the object.
	// >   - OSS supports the following Content-Type values to compress the object in the GZip format: text/cache-manifest, text/xml, text/plain, text/css, application/javascript, application/x-javascript, application/rss+xml, application/json, and text/json.
	//
	// Default value: null
	AcceptEncoding *string `json:"Accept-Encoding,omitempty" xml:"Accept-Encoding,omitempty"`
	// If the ETag specified in the request matches the ETag value of the object, OSS transmits the object and returns 200 OK. If the ETag specified in the request does not match the ETag value of the object, OSS returns 412 Precondition Failed.
	// The ETag value of an object is used to check whether the content of the object has changed. You can check data integrity by using the ETag value.
	// Default value: null
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// If the time specified in this header is earlier than the object modified time or is invalid, OSS returns the object and 200 OK. If the time specified in this header is later than or the same as the object modified time, OSS returns 304 Not Modified.
	// The time must be in GMT. Example: `Fri, 13 Nov 2015 14:47:53 GMT`.
	// Default value: null
	IfModifiedSince *string `json:"If-Modified-Since,omitempty" xml:"If-Modified-Since,omitempty"`
	// If the ETag specified in the request does not match the ETag value of the object, OSS transmits the object and returns 200 OK. If the ETag specified in the request matches the ETag value of the object, OSS returns 304 Not Modified.
	// You can specify both the **If-Match** and **If-None-Match** headers in a request.
	// Default value: null
	IfNoneMatch *string `json:"If-None-Match,omitempty" xml:"If-None-Match,omitempty"`
	// If the time specified in this header is the same as or later than the object modified time, OSS returns the object and 200 OK. If the time specified in this header is earlier than the object modified time, OSS returns 412 Precondition Failed.
	//
	// The time must be in GMT. Example: `Fri, 13 Nov 2015 14:47:53 GMT`.
	// You can specify both the **If-Modified-Since** and **If-Unmodified-Since** headers in a request.
	// Default value: null
	IfUnmodifiedSince *string `json:"If-Unmodified-Since,omitempty" xml:"If-Unmodified-Since,omitempty"`
	// The range of data of the object to be returned.
	//   - If the value of Range is valid, OSS returns the response that includes the total size of the object and the range of data returned. For example, Content-Range: bytes 0~9/44 indicates that the total size of the object is 44 bytes, and the range of data returned is the first 10 bytes.
	//   - However, if the value of Range is invalid, the entire object is returned, and the response returned by OSS excludes Content-Range.
	//
	// Default value: null
	Range *string `json:"Range,omitempty" xml:"Range,omitempty"`
}

func (s GetObjectHeaders) String() string {
	return tea.Prettify(s)
}

func (s GetObjectHeaders) GoString() string {
	return s.String()
}

func (s *GetObjectHeaders) SetCommonHeaders(v map[string]*string) *GetObjectHeaders {
	s.CommonHeaders = v
	return s
}

func (s *GetObjectHeaders) SetAcceptEncoding(v string) *GetObjectHeaders {
	s.AcceptEncoding = &v
	return s
}

func (s *GetObjectHeaders) SetIfMatch(v string) *GetObjectHeaders {
	s.IfMatch = &v
	return s
}

func (s *GetObjectHeaders) SetIfModifiedSince(v string) *GetObjectHeaders {
	s.IfModifiedSince = &v
	return s
}

func (s *GetObjectHeaders) SetIfNoneMatch(v string) *GetObjectHeaders {
	s.IfNoneMatch = &v
	return s
}

func (s *GetObjectHeaders) SetIfUnmodifiedSince(v string) *GetObjectHeaders {
	s.IfUnmodifiedSince = &v
	return s
}

func (s *GetObjectHeaders) SetRange(v string) *GetObjectHeaders {
	s.Range = &v
	return s
}

type GetObjectRequest struct {
	// The cache-control header in the response that OSS returns.
	ResponseCacheControl *string `json:"response-cache-control,omitempty" xml:"response-cache-control,omitempty"`
	// The content-disposition header in the response that OSS returns.
	ResponseContentDisposition *string `json:"response-content-disposition,omitempty" xml:"response-content-disposition,omitempty"`
	// The content-encoding header in the response that OSS returns.
	ResponseContentEncoding *string `json:"response-content-encoding,omitempty" xml:"response-content-encoding,omitempty"`
	// The content-language header in the response that OSS returns.
	ResponseContentLanguage *string `json:"response-content-language,omitempty" xml:"response-content-language,omitempty"`
	// The content-type header in the response that OSS returns.
	ResponseContentType *string `json:"response-content-type,omitempty" xml:"response-content-type,omitempty"`
	// The expires header in the response that OSS returns.
	ResponseExpires *string `json:"response-expires,omitempty" xml:"response-expires,omitempty"`
	// The version ID of the object that you want to query.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s GetObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s GetObjectRequest) GoString() string {
	return s.String()
}

func (s *GetObjectRequest) SetResponseCacheControl(v string) *GetObjectRequest {
	s.ResponseCacheControl = &v
	return s
}

func (s *GetObjectRequest) SetResponseContentDisposition(v string) *GetObjectRequest {
	s.ResponseContentDisposition = &v
	return s
}

func (s *GetObjectRequest) SetResponseContentEncoding(v string) *GetObjectRequest {
	s.ResponseContentEncoding = &v
	return s
}

func (s *GetObjectRequest) SetResponseContentLanguage(v string) *GetObjectRequest {
	s.ResponseContentLanguage = &v
	return s
}

func (s *GetObjectRequest) SetResponseContentType(v string) *GetObjectRequest {
	s.ResponseContentType = &v
	return s
}

func (s *GetObjectRequest) SetResponseExpires(v string) *GetObjectRequest {
	s.ResponseExpires = &v
	return s
}

func (s *GetObjectRequest) SetVersionId(v string) *GetObjectRequest {
	s.VersionId = &v
	return s
}

type GetObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       io.Reader          `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s GetObjectResponse) GoString() string {
	return s.String()
}

func (s *GetObjectResponse) SetHeaders(v map[string]*string) *GetObjectResponse {
	s.Headers = v
	return s
}

func (s *GetObjectResponse) SetStatusCode(v int32) *GetObjectResponse {
	s.StatusCode = &v
	return s
}

func (s *GetObjectResponse) SetBody(v io.Reader) *GetObjectResponse {
	s.Body = v
	return s
}

type GetObjectAclRequest struct {
	// The verison id of the target object.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s GetObjectAclRequest) String() string {
	return tea.Prettify(s)
}

func (s GetObjectAclRequest) GoString() string {
	return s.String()
}

func (s *GetObjectAclRequest) SetVersionId(v string) *GetObjectAclRequest {
	s.VersionId = &v
	return s
}

type GetObjectAclResponseBody struct {
	// The container that stores the results of the GetObjectACL request.
	AccessControlPolicy *GetObjectAclResponseBodyAccessControlPolicy `json:"AccessControlPolicy,omitempty" xml:"AccessControlPolicy,omitempty" type:"Struct"`
}

func (s GetObjectAclResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetObjectAclResponseBody) GoString() string {
	return s.String()
}

func (s *GetObjectAclResponseBody) SetAccessControlPolicy(v *GetObjectAclResponseBodyAccessControlPolicy) *GetObjectAclResponseBody {
	s.AccessControlPolicy = v
	return s
}

type GetObjectAclResponseBodyAccessControlPolicy struct {
	// The container that stores the ACL information.
	AccessControlList *GetObjectAclResponseBodyAccessControlPolicyAccessControlList `json:"AccessControlList,omitempty" xml:"AccessControlList,omitempty" type:"Struct"`
	// The container that stores the information about the bucket owner.
	Owner *Owner `json:"Owner,omitempty" xml:"Owner,omitempty"`
}

func (s GetObjectAclResponseBodyAccessControlPolicy) String() string {
	return tea.Prettify(s)
}

func (s GetObjectAclResponseBodyAccessControlPolicy) GoString() string {
	return s.String()
}

func (s *GetObjectAclResponseBodyAccessControlPolicy) SetAccessControlList(v *GetObjectAclResponseBodyAccessControlPolicyAccessControlList) *GetObjectAclResponseBodyAccessControlPolicy {
	s.AccessControlList = v
	return s
}

func (s *GetObjectAclResponseBodyAccessControlPolicy) SetOwner(v *Owner) *GetObjectAclResponseBodyAccessControlPolicy {
	s.Owner = v
	return s
}

type GetObjectAclResponseBodyAccessControlPolicyAccessControlList struct {
	// The ACL of the object. Default value: default.
	ACL *string `json:"Grant,omitempty" xml:"Grant,omitempty"`
}

func (s GetObjectAclResponseBodyAccessControlPolicyAccessControlList) String() string {
	return tea.Prettify(s)
}

func (s GetObjectAclResponseBodyAccessControlPolicyAccessControlList) GoString() string {
	return s.String()
}

func (s *GetObjectAclResponseBodyAccessControlPolicyAccessControlList) SetACL(v string) *GetObjectAclResponseBodyAccessControlPolicyAccessControlList {
	s.ACL = &v
	return s
}

type GetObjectAclResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetObjectAclResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetObjectAclResponse) String() string {
	return tea.Prettify(s)
}

func (s GetObjectAclResponse) GoString() string {
	return s.String()
}

func (s *GetObjectAclResponse) SetHeaders(v map[string]*string) *GetObjectAclResponse {
	s.Headers = v
	return s
}

func (s *GetObjectAclResponse) SetStatusCode(v int32) *GetObjectAclResponse {
	s.StatusCode = &v
	return s
}

func (s *GetObjectAclResponse) SetBody(v *GetObjectAclResponseBody) *GetObjectAclResponse {
	s.Body = v
	return s
}

type GetObjectMetaRequest struct {
	// The versionID of the object.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s GetObjectMetaRequest) String() string {
	return tea.Prettify(s)
}

func (s GetObjectMetaRequest) GoString() string {
	return s.String()
}

func (s *GetObjectMetaRequest) SetVersionId(v string) *GetObjectMetaRequest {
	s.VersionId = &v
	return s
}

type GetObjectMetaResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s GetObjectMetaResponse) String() string {
	return tea.Prettify(s)
}

func (s GetObjectMetaResponse) GoString() string {
	return s.String()
}

func (s *GetObjectMetaResponse) SetHeaders(v map[string]*string) *GetObjectMetaResponse {
	s.Headers = v
	return s
}

func (s *GetObjectMetaResponse) SetStatusCode(v int32) *GetObjectMetaResponse {
	s.StatusCode = &v
	return s
}

type GetObjectTaggingRequest struct {
	// The versionID of the object that you want to query.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s GetObjectTaggingRequest) String() string {
	return tea.Prettify(s)
}

func (s GetObjectTaggingRequest) GoString() string {
	return s.String()
}

func (s *GetObjectTaggingRequest) SetVersionId(v string) *GetObjectTaggingRequest {
	s.VersionId = &v
	return s
}

type GetObjectTaggingResponseBody struct {
	// The container that stores the returned tag of the bucket.
	Tagging *GetObjectTaggingResponseBodyTagging `json:"Tagging,omitempty" xml:"Tagging,omitempty" type:"Struct"`
}

func (s GetObjectTaggingResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetObjectTaggingResponseBody) GoString() string {
	return s.String()
}

func (s *GetObjectTaggingResponseBody) SetTagging(v *GetObjectTaggingResponseBodyTagging) *GetObjectTaggingResponseBody {
	s.Tagging = v
	return s
}

type GetObjectTaggingResponseBodyTagging struct {
	// The tag set of the target object.
	TagSet *TagSet `json:"TagSet,omitempty" xml:"TagSet,omitempty"`
}

func (s GetObjectTaggingResponseBodyTagging) String() string {
	return tea.Prettify(s)
}

func (s GetObjectTaggingResponseBodyTagging) GoString() string {
	return s.String()
}

func (s *GetObjectTaggingResponseBodyTagging) SetTagSet(v *TagSet) *GetObjectTaggingResponseBodyTagging {
	s.TagSet = v
	return s
}

type GetObjectTaggingResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetObjectTaggingResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetObjectTaggingResponse) String() string {
	return tea.Prettify(s)
}

func (s GetObjectTaggingResponse) GoString() string {
	return s.String()
}

func (s *GetObjectTaggingResponse) SetHeaders(v map[string]*string) *GetObjectTaggingResponse {
	s.Headers = v
	return s
}

func (s *GetObjectTaggingResponse) SetStatusCode(v int32) *GetObjectTaggingResponse {
	s.StatusCode = &v
	return s
}

func (s *GetObjectTaggingResponse) SetBody(v *GetObjectTaggingResponseBody) *GetObjectTaggingResponse {
	s.Body = v
	return s
}

type GetPublicAccessBlockResponseBody struct {
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
}

func (s GetPublicAccessBlockResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetPublicAccessBlockResponseBody) GoString() string {
	return s.String()
}

func (s *GetPublicAccessBlockResponseBody) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *GetPublicAccessBlockResponseBody {
	s.PublicAccessBlockConfiguration = v
	return s
}

type GetPublicAccessBlockResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetPublicAccessBlockResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s GetPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *GetPublicAccessBlockResponse) SetHeaders(v map[string]*string) *GetPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *GetPublicAccessBlockResponse) SetStatusCode(v int32) *GetPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

func (s *GetPublicAccessBlockResponse) SetBody(v *GetPublicAccessBlockResponseBody) *GetPublicAccessBlockResponse {
	s.Body = v
	return s
}

type GetStyleRequest struct {
	// The name of the image style.
	StyleName *string `json:"styleName,omitempty" xml:"styleName,omitempty"`
}

func (s GetStyleRequest) String() string {
	return tea.Prettify(s)
}

func (s GetStyleRequest) GoString() string {
	return s.String()
}

func (s *GetStyleRequest) SetStyleName(v string) *GetStyleRequest {
	s.StyleName = &v
	return s
}

type GetStyleResponseBody struct {
	// The container that stores the information about the image style.
	Style *StyleInfo `json:"Style,omitempty" xml:"Style,omitempty"`
}

func (s GetStyleResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetStyleResponseBody) GoString() string {
	return s.String()
}

func (s *GetStyleResponseBody) SetStyle(v *StyleInfo) *GetStyleResponseBody {
	s.Style = v
	return s
}

type GetStyleResponse struct {
	Headers    map[string]*string    `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetStyleResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetStyleResponse) String() string {
	return tea.Prettify(s)
}

func (s GetStyleResponse) GoString() string {
	return s.String()
}

func (s *GetStyleResponse) SetHeaders(v map[string]*string) *GetStyleResponse {
	s.Headers = v
	return s
}

func (s *GetStyleResponse) SetStatusCode(v int32) *GetStyleResponse {
	s.StatusCode = &v
	return s
}

func (s *GetStyleResponse) SetBody(v *GetStyleResponseBody) *GetStyleResponse {
	s.Body = v
	return s
}

type GetSymlinkRequest struct {
	// The version of the object to which the symbolic link points.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s GetSymlinkRequest) String() string {
	return tea.Prettify(s)
}

func (s GetSymlinkRequest) GoString() string {
	return s.String()
}

func (s *GetSymlinkRequest) SetVersionId(v string) *GetSymlinkRequest {
	s.VersionId = &v
	return s
}

type GetSymlinkResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s GetSymlinkResponse) String() string {
	return tea.Prettify(s)
}

func (s GetSymlinkResponse) GoString() string {
	return s.String()
}

func (s *GetSymlinkResponse) SetHeaders(v map[string]*string) *GetSymlinkResponse {
	s.Headers = v
	return s
}

func (s *GetSymlinkResponse) SetStatusCode(v int32) *GetSymlinkResponse {
	s.StatusCode = &v
	return s
}

type GetUserAntiDDosInfoResponseBody struct {
	// The container that stores the list of Anti-DDoS instances.
	AntiDDOSListConfiguration *GetUserAntiDDosInfoResponseBodyAntiDDOSListConfiguration `json:"AntiDDOSListConfiguration,omitempty" xml:"AntiDDOSListConfiguration,omitempty" type:"Struct"`
}

func (s GetUserAntiDDosInfoResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetUserAntiDDosInfoResponseBody) GoString() string {
	return s.String()
}

func (s *GetUserAntiDDosInfoResponseBody) SetAntiDDOSListConfiguration(v *GetUserAntiDDosInfoResponseBodyAntiDDOSListConfiguration) *GetUserAntiDDosInfoResponseBody {
	s.AntiDDOSListConfiguration = v
	return s
}

type GetUserAntiDDosInfoResponseBodyAntiDDOSListConfiguration struct {
	// The container that stores information about the Anti-DDoS instance.
	AntiDDOSConfiguration []*UserAntiDDOSInfo `json:"AntiDDOSConfiguration,omitempty" xml:"AntiDDOSConfiguration,omitempty" type:"Repeated"`
}

func (s GetUserAntiDDosInfoResponseBodyAntiDDOSListConfiguration) String() string {
	return tea.Prettify(s)
}

func (s GetUserAntiDDosInfoResponseBodyAntiDDOSListConfiguration) GoString() string {
	return s.String()
}

func (s *GetUserAntiDDosInfoResponseBodyAntiDDOSListConfiguration) SetAntiDDOSConfiguration(v []*UserAntiDDOSInfo) *GetUserAntiDDosInfoResponseBodyAntiDDOSListConfiguration {
	s.AntiDDOSConfiguration = v
	return s
}

type GetUserAntiDDosInfoResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetUserAntiDDosInfoResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetUserAntiDDosInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s GetUserAntiDDosInfoResponse) GoString() string {
	return s.String()
}

func (s *GetUserAntiDDosInfoResponse) SetHeaders(v map[string]*string) *GetUserAntiDDosInfoResponse {
	s.Headers = v
	return s
}

func (s *GetUserAntiDDosInfoResponse) SetStatusCode(v int32) *GetUserAntiDDosInfoResponse {
	s.StatusCode = &v
	return s
}

func (s *GetUserAntiDDosInfoResponse) SetBody(v *GetUserAntiDDosInfoResponseBody) *GetUserAntiDDosInfoResponse {
	s.Body = v
	return s
}

type GetUserDefinedLogFieldsConfigResponseBody struct {
	UserDefinedLogFieldsConfiguration *UserDefinedLogFieldsConfiguration `json:"UserDefinedLogFieldsConfiguration,omitempty" xml:"UserDefinedLogFieldsConfiguration,omitempty"`
}

func (s GetUserDefinedLogFieldsConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetUserDefinedLogFieldsConfigResponseBody) GoString() string {
	return s.String()
}

func (s *GetUserDefinedLogFieldsConfigResponseBody) SetUserDefinedLogFieldsConfiguration(v *UserDefinedLogFieldsConfiguration) *GetUserDefinedLogFieldsConfigResponseBody {
	s.UserDefinedLogFieldsConfiguration = v
	return s
}

type GetUserDefinedLogFieldsConfigResponse struct {
	Headers    map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetUserDefinedLogFieldsConfigResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetUserDefinedLogFieldsConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s GetUserDefinedLogFieldsConfigResponse) GoString() string {
	return s.String()
}

func (s *GetUserDefinedLogFieldsConfigResponse) SetHeaders(v map[string]*string) *GetUserDefinedLogFieldsConfigResponse {
	s.Headers = v
	return s
}

func (s *GetUserDefinedLogFieldsConfigResponse) SetStatusCode(v int32) *GetUserDefinedLogFieldsConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *GetUserDefinedLogFieldsConfigResponse) SetBody(v *GetUserDefinedLogFieldsConfigResponseBody) *GetUserDefinedLogFieldsConfigResponse {
	s.Body = v
	return s
}

type GetVodPlaylistRequest struct {
	// The end time of the time range during which the TS files that you want to query are generated in the Unix timestamp format.
	// > The value of EndTime must be greater than the value of StartTime. The duration between EndTime and StartTime must be less than one day.
	EndTime *string `json:"endTime,omitempty" xml:"endTime,omitempty"`
	// The start time of the time range during which the TS files that you want to query are generated in the Unix timestamp format.
	StartTime *string `json:"startTime,omitempty" xml:"startTime,omitempty"`
}

func (s GetVodPlaylistRequest) String() string {
	return tea.Prettify(s)
}

func (s GetVodPlaylistRequest) GoString() string {
	return s.String()
}

func (s *GetVodPlaylistRequest) SetEndTime(v string) *GetVodPlaylistRequest {
	s.EndTime = &v
	return s
}

func (s *GetVodPlaylistRequest) SetStartTime(v string) *GetVodPlaylistRequest {
	s.StartTime = &v
	return s
}

type GetVodPlaylistResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       io.Reader          `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetVodPlaylistResponse) String() string {
	return tea.Prettify(s)
}

func (s GetVodPlaylistResponse) GoString() string {
	return s.String()
}

func (s *GetVodPlaylistResponse) SetHeaders(v map[string]*string) *GetVodPlaylistResponse {
	s.Headers = v
	return s
}

func (s *GetVodPlaylistResponse) SetStatusCode(v int32) *GetVodPlaylistResponse {
	s.StatusCode = &v
	return s
}

func (s *GetVodPlaylistResponse) SetBody(v io.Reader) *GetVodPlaylistResponse {
	s.Body = v
	return s
}

type HeadObjectHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// If the ETag value that is specified in the request matches the ETag value of the object, OSS returns 200 OK and the metadata of the object. Otherwise, OSS returns 412 precondition failed.
	// Default value: null.
	IfMatch *string `json:"If-Match,omitempty" xml:"If-Match,omitempty"`
	// If the time that is specified in the request is earlier than the time when the object is modified, OSS returns 200 OK and the metadata of the object. Otherwise, OSS returns 304 not modified.
	// Default value: null.
	IfModifiedSince *string `json:"If-Modified-Since,omitempty" xml:"If-Modified-Since,omitempty"`
	// If the ETag value that is specified in the request does not match the ETag value of the object, OSS returns 200 OK and the metadata of the object. Otherwise, OSS returns 304 Not Modified.
	// Default value: null.
	IfNoneMatch *string `json:"If-None-Match,omitempty" xml:"If-None-Match,omitempty"`
	// If the time that is specified in the request is later than or the same as the time when the object is modified, OSS returns 200 OK and the metadata of the object. Otherwise, OSS returns 412 precondition failed.
	// Default value: null.
	IfUnmodifiedSince *string `json:"If-Unmodified-Since,omitempty" xml:"If-Unmodified-Since,omitempty"`
}

func (s HeadObjectHeaders) String() string {
	return tea.Prettify(s)
}

func (s HeadObjectHeaders) GoString() string {
	return s.String()
}

func (s *HeadObjectHeaders) SetCommonHeaders(v map[string]*string) *HeadObjectHeaders {
	s.CommonHeaders = v
	return s
}

func (s *HeadObjectHeaders) SetIfMatch(v string) *HeadObjectHeaders {
	s.IfMatch = &v
	return s
}

func (s *HeadObjectHeaders) SetIfModifiedSince(v string) *HeadObjectHeaders {
	s.IfModifiedSince = &v
	return s
}

func (s *HeadObjectHeaders) SetIfNoneMatch(v string) *HeadObjectHeaders {
	s.IfNoneMatch = &v
	return s
}

func (s *HeadObjectHeaders) SetIfUnmodifiedSince(v string) *HeadObjectHeaders {
	s.IfUnmodifiedSince = &v
	return s
}

type HeadObjectRequest struct {
	// The version ID of the object for which you want to query metadata.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s HeadObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s HeadObjectRequest) GoString() string {
	return s.String()
}

func (s *HeadObjectRequest) SetVersionId(v string) *HeadObjectRequest {
	s.VersionId = &v
	return s
}

type HeadObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s HeadObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s HeadObjectResponse) GoString() string {
	return s.String()
}

func (s *HeadObjectResponse) SetHeaders(v map[string]*string) *HeadObjectResponse {
	s.Headers = v
	return s
}

func (s *HeadObjectResponse) SetStatusCode(v int32) *HeadObjectResponse {
	s.StatusCode = &v
	return s
}

type InitBucketAntiDDosInfoHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of the Anti-DDoS instance.
	DefenderInstance *string `json:"x-oss-defender-instance,omitempty" xml:"x-oss-defender-instance,omitempty"`
	// The type of the Anti-DDoS instance. Set the value to AntiDDosPremimum.
	DefenderType *string `json:"x-oss-defender-type,omitempty" xml:"x-oss-defender-type,omitempty"`
}

func (s InitBucketAntiDDosInfoHeaders) String() string {
	return tea.Prettify(s)
}

func (s InitBucketAntiDDosInfoHeaders) GoString() string {
	return s.String()
}

func (s *InitBucketAntiDDosInfoHeaders) SetCommonHeaders(v map[string]*string) *InitBucketAntiDDosInfoHeaders {
	s.CommonHeaders = v
	return s
}

func (s *InitBucketAntiDDosInfoHeaders) SetDefenderInstance(v string) *InitBucketAntiDDosInfoHeaders {
	s.DefenderInstance = &v
	return s
}

func (s *InitBucketAntiDDosInfoHeaders) SetDefenderType(v string) *InitBucketAntiDDosInfoHeaders {
	s.DefenderType = &v
	return s
}

type InitBucketAntiDDosInfoRequest struct {
	// The container that stores the configurations of Anti-DDoS instances.
	AntiDDOSConfiguration *BucketAntiDDOSConfiguration `json:"AntiDDOSConfiguration,omitempty" xml:"AntiDDOSConfiguration,omitempty"`
}

func (s InitBucketAntiDDosInfoRequest) String() string {
	return tea.Prettify(s)
}

func (s InitBucketAntiDDosInfoRequest) GoString() string {
	return s.String()
}

func (s *InitBucketAntiDDosInfoRequest) SetAntiDDOSConfiguration(v *BucketAntiDDOSConfiguration) *InitBucketAntiDDosInfoRequest {
	s.AntiDDOSConfiguration = v
	return s
}

type InitBucketAntiDDosInfoResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s InitBucketAntiDDosInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s InitBucketAntiDDosInfoResponse) GoString() string {
	return s.String()
}

func (s *InitBucketAntiDDosInfoResponse) SetHeaders(v map[string]*string) *InitBucketAntiDDosInfoResponse {
	s.Headers = v
	return s
}

func (s *InitBucketAntiDDosInfoResponse) SetStatusCode(v int32) *InitBucketAntiDDosInfoResponse {
	s.StatusCode = &v
	return s
}

type InitUserAntiDDosInfoResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s InitUserAntiDDosInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s InitUserAntiDDosInfoResponse) GoString() string {
	return s.String()
}

func (s *InitUserAntiDDosInfoResponse) SetHeaders(v map[string]*string) *InitUserAntiDDosInfoResponse {
	s.Headers = v
	return s
}

func (s *InitUserAntiDDosInfoResponse) SetStatusCode(v int32) *InitUserAntiDDosInfoResponse {
	s.StatusCode = &v
	return s
}

type InitiateBucketWormRequest struct {
	// The parameters for WORM initialization.
	InitiateWormConfiguration *InitiateWormConfiguration `json:"InitiateWormConfiguration,omitempty" xml:"InitiateWormConfiguration,omitempty"`
}

func (s InitiateBucketWormRequest) String() string {
	return tea.Prettify(s)
}

func (s InitiateBucketWormRequest) GoString() string {
	return s.String()
}

func (s *InitiateBucketWormRequest) SetInitiateWormConfiguration(v *InitiateWormConfiguration) *InitiateBucketWormRequest {
	s.InitiateWormConfiguration = v
	return s
}

type InitiateBucketWormResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s InitiateBucketWormResponse) String() string {
	return tea.Prettify(s)
}

func (s InitiateBucketWormResponse) GoString() string {
	return s.String()
}

func (s *InitiateBucketWormResponse) SetHeaders(v map[string]*string) *InitiateBucketWormResponse {
	s.Headers = v
	return s
}

func (s *InitiateBucketWormResponse) SetStatusCode(v int32) *InitiateBucketWormResponse {
	s.StatusCode = &v
	return s
}

type InitiateMultipartUploadHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The caching behavior of the web page when the object is downloaded. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	CacheControl *string `json:"Cache-Control,omitempty" xml:"Cache-Control,omitempty"`
	// The name of the object when the object is downloaded. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	ContentDisposition *string `json:"Content-Disposition,omitempty" xml:"Content-Disposition,omitempty"`
	// The content encoding format of the object when the object is downloaded. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	ContentEncoding *string `json:"Content-Encoding,omitempty" xml:"Content-Encoding,omitempty"`
	// The expiration time of the request. Unit: milliseconds. For more information, see **[RFC 2616](https://www.ietf.org/rfc/rfc2616.txt)**.
	// Default value: null.
	Expires *string `json:"Expires,omitempty" xml:"Expires,omitempty"`
	// Specifies whether the InitiateMultipartUpload operation overwrites the existing object that has the same name as the object that you want to upload. When versioning is enabled or suspended for the bucket to which you want to upload the object, the **x-oss-forbid-overwrite** header does not take effect. In this case, the InitiateMultipartUpload operation overwrites the existing object that has the same name as the object that you want to upload.
	//   - If you do not specify the **x-oss-forbid-overwrite** header or set the **x-oss-forbid-overwrite** header to **false**, the object that is uploaded by calling the PutObject operation overwrites the existing object that has the same name.
	//   - If the value of **x-oss-forbid-overwrite** is set to **true**, existing objects cannot be overwritten by objects that have the same names.
	//
	// If you specify the **x-oss-forbid-overwrite** request header, the queries per second (QPS) performance of OSS is degraded. If you want to use the **x-oss-forbid-overwrite** request header to perform a large number of operations (QPS greater than 1,000), contact technical support
	ForbidOverwrite *string `json:"x-oss-forbid-overwrite,omitempty" xml:"x-oss-forbid-overwrite,omitempty"`
	// The algorithm that is used to encrypt the object that you want to upload. If this header is not specified, the object is encrypted by using AES-256. This header is valid only when **x-oss-server-side-encryption** is set to KMS.
	// Valid value: SM4.
	SseDataEncryption *string `json:"x-oss-server-side-data-encryption,omitempty" xml:"x-oss-server-side-data-encryption,omitempty"`
	// The server-side encryption method that is used to encrypt each part of the object that you want to upload.
	// Valid values: **AES256**, **KMS**, and **SM4**.
	// > You must activate Key Management Service (KMS) before you set this header to KMS.
	//
	//
	// If you specify this header in the request, this header is included in the response. OSS uses the method specified by this header to encrypt each uploaded part. When you download the object, the x-oss-server-side-encryption header is included in the response and the header value is set to the algorithm that is used to encrypt the object.
	ServerSideEncryption *string `json:"x-oss-server-side-encryption,omitempty" xml:"x-oss-server-side-encryption,omitempty"`
	// The ID of the CMK that is managed by KMS.
	// This header is valid only when **x-oss-server-side-encryption** is set to KMS.
	SseKeyId *string `json:"x-oss-server-side-encryption-key-id,omitempty" xml:"x-oss-server-side-encryption-key-id,omitempty"`
	// The storage class of the bucket. Default value: Standard.  Valid values:
	//
	// - Standard
	// - IA
	// - Archive
	// - ColdArchive
	StorageClass *string `json:"x-oss-storage-class,omitempty" xml:"x-oss-storage-class,omitempty"`
	// The tag of the object. You can configure multiple tags for the object. Example: TagA=A&amp;TagB=B.
	// > The key and value of a tag must be URL-encoded. If a tag does not contain an equal sign (=), the value of the tag is considered an empty string.
	Tagging *string `json:"x-oss-tagging,omitempty" xml:"x-oss-tagging,omitempty"`
}

func (s InitiateMultipartUploadHeaders) String() string {
	return tea.Prettify(s)
}

func (s InitiateMultipartUploadHeaders) GoString() string {
	return s.String()
}

func (s *InitiateMultipartUploadHeaders) SetCommonHeaders(v map[string]*string) *InitiateMultipartUploadHeaders {
	s.CommonHeaders = v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetCacheControl(v string) *InitiateMultipartUploadHeaders {
	s.CacheControl = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetContentDisposition(v string) *InitiateMultipartUploadHeaders {
	s.ContentDisposition = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetContentEncoding(v string) *InitiateMultipartUploadHeaders {
	s.ContentEncoding = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetExpires(v string) *InitiateMultipartUploadHeaders {
	s.Expires = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetForbidOverwrite(v string) *InitiateMultipartUploadHeaders {
	s.ForbidOverwrite = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetSseDataEncryption(v string) *InitiateMultipartUploadHeaders {
	s.SseDataEncryption = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetServerSideEncryption(v string) *InitiateMultipartUploadHeaders {
	s.ServerSideEncryption = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetSseKeyId(v string) *InitiateMultipartUploadHeaders {
	s.SseKeyId = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetStorageClass(v string) *InitiateMultipartUploadHeaders {
	s.StorageClass = &v
	return s
}

func (s *InitiateMultipartUploadHeaders) SetTagging(v string) *InitiateMultipartUploadHeaders {
	s.Tagging = &v
	return s
}

type InitiateMultipartUploadRequest struct {
	// The method used to encode the object name in the response. Only URL encoding is supported. The object name can contain characters encoded in UTF-8. However, the XML 1.0 standard cannot be used to parse specific control characters, such as characters whose ASCII values range from 0 to 10. You can configure the encoding-type parameter to encode object names that include characters that cannot be parsed by XML 1.0 in the response.
	// <br>Default value: null
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
}

func (s InitiateMultipartUploadRequest) String() string {
	return tea.Prettify(s)
}

func (s InitiateMultipartUploadRequest) GoString() string {
	return s.String()
}

func (s *InitiateMultipartUploadRequest) SetEncodingType(v string) *InitiateMultipartUploadRequest {
	s.EncodingType = &v
	return s
}

type InitiateMultipartUploadResponseBody struct {
	// The container that stores the results of the InitiateMultipartUpload request.
	InitiateMultipartUploadResult *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult `json:"InitiateMultipartUploadResult,omitempty" xml:"InitiateMultipartUploadResult,omitempty" type:"Struct"`
}

func (s InitiateMultipartUploadResponseBody) String() string {
	return tea.Prettify(s)
}

func (s InitiateMultipartUploadResponseBody) GoString() string {
	return s.String()
}

func (s *InitiateMultipartUploadResponseBody) SetInitiateMultipartUploadResult(v *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult) *InitiateMultipartUploadResponseBody {
	s.InitiateMultipartUploadResult = v
	return s
}

type InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult struct {
	// The name of the bucket to which the object is uploaded by the multipart upload task.
	Bucket *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	// The encoding type of the object name in the response. If the encoding-type parameter is specified in the request, the object name in the response is encoded.
	EncodingType *string `json:"EncodingType,omitempty" xml:"EncodingType,omitempty"`
	// The name of the object that is uploaded by the multipart upload task.
	Key *string `json:"Key,omitempty" xml:"Key,omitempty"`
	// The upload ID that uniquely identifies the multipart upload task. The upload ID is used to call UploadPart and CompleteMultipartUpload later.
	UploadId *string `json:"UploadId,omitempty" xml:"UploadId,omitempty"`
}

func (s InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult) String() string {
	return tea.Prettify(s)
}

func (s InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult) GoString() string {
	return s.String()
}

func (s *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult) SetBucket(v string) *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult {
	s.Bucket = &v
	return s
}

func (s *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult) SetEncodingType(v string) *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult {
	s.EncodingType = &v
	return s
}

func (s *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult) SetKey(v string) *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult {
	s.Key = &v
	return s
}

func (s *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult) SetUploadId(v string) *InitiateMultipartUploadResponseBodyInitiateMultipartUploadResult {
	s.UploadId = &v
	return s
}

type InitiateMultipartUploadResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *InitiateMultipartUploadResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s InitiateMultipartUploadResponse) String() string {
	return tea.Prettify(s)
}

func (s InitiateMultipartUploadResponse) GoString() string {
	return s.String()
}

func (s *InitiateMultipartUploadResponse) SetHeaders(v map[string]*string) *InitiateMultipartUploadResponse {
	s.Headers = v
	return s
}

func (s *InitiateMultipartUploadResponse) SetStatusCode(v int32) *InitiateMultipartUploadResponse {
	s.StatusCode = &v
	return s
}

func (s *InitiateMultipartUploadResponse) SetBody(v *InitiateMultipartUploadResponseBody) *InitiateMultipartUploadResponse {
	s.Body = v
	return s
}

type ListAccessPointsRequest struct {
	// The token from which the listing operation starts. You must specify the value of NextContinuationToken that is obtained from the previous query as the value of continuation-token.
	ContinuationToken *string `json:"continuation-token,omitempty" xml:"continuation-token,omitempty"`
	// The maximum number of access points that can be returned. Valid values:
	//
	// *   For user-level access points: (0,1000].
	// *   For bucket-level access points: (0,100].
	MaxKeys *int64 `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
}

func (s ListAccessPointsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsRequest) GoString() string {
	return s.String()
}

func (s *ListAccessPointsRequest) SetContinuationToken(v string) *ListAccessPointsRequest {
	s.ContinuationToken = &v
	return s
}

func (s *ListAccessPointsRequest) SetMaxKeys(v int64) *ListAccessPointsRequest {
	s.MaxKeys = &v
	return s
}

type ListAccessPointsResponseBody struct {
	// The container that stores the information about access points.
	ListAccessPointsResult *ListAccessPointsResult `json:"ListAccessPointsResult,omitempty" xml:"ListAccessPointsResult,omitempty"`
}

func (s ListAccessPointsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsResponseBody) GoString() string {
	return s.String()
}

func (s *ListAccessPointsResponseBody) SetListAccessPointsResult(v *ListAccessPointsResult) *ListAccessPointsResponseBody {
	s.ListAccessPointsResult = v
	return s
}

type ListAccessPointsResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListAccessPointsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListAccessPointsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsResponse) GoString() string {
	return s.String()
}

func (s *ListAccessPointsResponse) SetHeaders(v map[string]*string) *ListAccessPointsResponse {
	s.Headers = v
	return s
}

func (s *ListAccessPointsResponse) SetStatusCode(v int32) *ListAccessPointsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListAccessPointsResponse) SetBody(v *ListAccessPointsResponseBody) *ListAccessPointsResponse {
	s.Body = v
	return s
}

type ListAccessPointsForObjectProcessRequest struct {
	ContinuationToken *string `json:"continuation-token,omitempty" xml:"continuation-token,omitempty"`
	MaxKeys           *int64  `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
}

func (s ListAccessPointsForObjectProcessRequest) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsForObjectProcessRequest) GoString() string {
	return s.String()
}

func (s *ListAccessPointsForObjectProcessRequest) SetContinuationToken(v string) *ListAccessPointsForObjectProcessRequest {
	s.ContinuationToken = &v
	return s
}

func (s *ListAccessPointsForObjectProcessRequest) SetMaxKeys(v int64) *ListAccessPointsForObjectProcessRequest {
	s.MaxKeys = &v
	return s
}

type ListAccessPointsForObjectProcessResponseBody struct {
	ListAccessPointsForObjectProcessResult *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult `json:"ListAccessPointsForObjectProcessResult,omitempty" xml:"ListAccessPointsForObjectProcessResult,omitempty" type:"Struct"`
}

func (s ListAccessPointsForObjectProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsForObjectProcessResponseBody) GoString() string {
	return s.String()
}

func (s *ListAccessPointsForObjectProcessResponseBody) SetListAccessPointsForObjectProcessResult(v *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult) *ListAccessPointsForObjectProcessResponseBody {
	s.ListAccessPointsForObjectProcessResult = v
	return s
}

type ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult struct {
	AccessPointsForObjectProcess *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcess `json:"AccessPointsForObjectProcess,omitempty" xml:"AccessPointsForObjectProcess,omitempty" type:"Struct"`
	AccountId                    *string                                                                                                         `json:"AccountId,omitempty" xml:"AccountId,omitempty"`
	IsTruncated                  *bool                                                                                                           `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	NextContinuationToken        *string                                                                                                         `json:"NextContinuationToken,omitempty" xml:"NextContinuationToken,omitempty"`
}

func (s ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult) GoString() string {
	return s.String()
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult) SetAccessPointsForObjectProcess(v *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcess) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult {
	s.AccessPointsForObjectProcess = v
	return s
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult) SetAccountId(v string) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult {
	s.AccountId = &v
	return s
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult) SetIsTruncated(v bool) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult {
	s.IsTruncated = &v
	return s
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult) SetNextContinuationToken(v string) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResult {
	s.NextContinuationToken = &v
	return s
}

type ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcess struct {
	AccessPointForObjectProcess []*ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess `json:"AccessPointForObjectProcess,omitempty" xml:"AccessPointForObjectProcess,omitempty" type:"Repeated"`
}

func (s ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcess) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcess) GoString() string {
	return s.String()
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcess) SetAccessPointForObjectProcess(v []*ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcess {
	s.AccessPointForObjectProcess = v
	return s
}

type ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess struct {
	AccessPointForObjectProcessAlias     *string `json:"AccessPointForObjectProcessAlias,omitempty" xml:"AccessPointForObjectProcessAlias,omitempty"`
	AccessPointName                      *string `json:"AccessPointName,omitempty" xml:"AccessPointName,omitempty"`
	AccessPointNameForObjectProcess      *string `json:"AccessPointNameForObjectProcess,omitempty" xml:"AccessPointNameForObjectProcess,omitempty"`
	AllowAnonymousAccessForObjectProcess *string `json:"AllowAnonymousAccessForObjectProcess,omitempty" xml:"AllowAnonymousAccessForObjectProcess,omitempty"`
	Status                               *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) GoString() string {
	return s.String()
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) SetAccessPointForObjectProcessAlias(v string) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess {
	s.AccessPointForObjectProcessAlias = &v
	return s
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) SetAccessPointName(v string) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess {
	s.AccessPointName = &v
	return s
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) SetAccessPointNameForObjectProcess(v string) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess {
	s.AccessPointNameForObjectProcess = &v
	return s
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) SetAllowAnonymousAccessForObjectProcess(v string) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess {
	s.AllowAnonymousAccessForObjectProcess = &v
	return s
}

func (s *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess) SetStatus(v string) *ListAccessPointsForObjectProcessResponseBodyListAccessPointsForObjectProcessResultAccessPointsForObjectProcessAccessPointForObjectProcess {
	s.Status = &v
	return s
}

type ListAccessPointsForObjectProcessResponse struct {
	Headers    map[string]*string                            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListAccessPointsForObjectProcessResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListAccessPointsForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s ListAccessPointsForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *ListAccessPointsForObjectProcessResponse) SetHeaders(v map[string]*string) *ListAccessPointsForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *ListAccessPointsForObjectProcessResponse) SetStatusCode(v int32) *ListAccessPointsForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

func (s *ListAccessPointsForObjectProcessResponse) SetBody(v *ListAccessPointsForObjectProcessResponseBody) *ListAccessPointsForObjectProcessResponse {
	s.Body = v
	return s
}

type ListBucketAntiDDosInfoRequest struct {
	// The name of the Anti-DDoS instance from which the list starts. The Anti-DDoS instances whose names are alphabetically after the value of marker are returned.
	//
	// >  You can set marker to an empty string in the first request. If IsTruncated is returned in the response and the value of IsTruncated is true, you must use the value of Marker in the response as the value of marker in the next request.
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
	// The maximum number of Anti-DDoS instances that can be returned.
	//
	// Valid values: 1 to 100.
	//
	// Default value: 100.
	MaxKeys *string `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
}

func (s ListBucketAntiDDosInfoRequest) String() string {
	return tea.Prettify(s)
}

func (s ListBucketAntiDDosInfoRequest) GoString() string {
	return s.String()
}

func (s *ListBucketAntiDDosInfoRequest) SetMarker(v string) *ListBucketAntiDDosInfoRequest {
	s.Marker = &v
	return s
}

func (s *ListBucketAntiDDosInfoRequest) SetMaxKeys(v string) *ListBucketAntiDDosInfoRequest {
	s.MaxKeys = &v
	return s
}

type ListBucketAntiDDosInfoResponseBody struct {
	// The container that stores the protection list of an Anti-DDoS instance of a bucket.
	AntiDDOSListConfiguration *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration `json:"AntiDDOSListConfiguration,omitempty" xml:"AntiDDOSListConfiguration,omitempty" type:"Struct"`
}

func (s ListBucketAntiDDosInfoResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListBucketAntiDDosInfoResponseBody) GoString() string {
	return s.String()
}

func (s *ListBucketAntiDDosInfoResponseBody) SetAntiDDOSListConfiguration(v *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration) *ListBucketAntiDDosInfoResponseBody {
	s.AntiDDOSListConfiguration = v
	return s
}

type ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration struct {
	// The container that stores information about the Anti-DDoS instance.
	AntiDDOSConfiguration []*BucketAntiDDOSInfo `json:"AntiDDOSConfiguration,omitempty" xml:"AntiDDOSConfiguration,omitempty" type:"Repeated"`
	// Indicates whether all Anti-DDoS instances are returned.
	//
	// - true: All Anti-DDoS instances are returned.
	//
	// - false: Not all Anti-DDoS instances are returned.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// The Anti-DDoS instances whose names are alphabetically after the specified marker.
	Marker *string `json:"Marker,omitempty" xml:"Marker,omitempty"`
}

func (s ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration) String() string {
	return tea.Prettify(s)
}

func (s ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration) GoString() string {
	return s.String()
}

func (s *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration) SetAntiDDOSConfiguration(v []*BucketAntiDDOSInfo) *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration {
	s.AntiDDOSConfiguration = v
	return s
}

func (s *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration) SetIsTruncated(v bool) *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration {
	s.IsTruncated = &v
	return s
}

func (s *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration) SetMarker(v string) *ListBucketAntiDDosInfoResponseBodyAntiDDOSListConfiguration {
	s.Marker = &v
	return s
}

type ListBucketAntiDDosInfoResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListBucketAntiDDosInfoResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListBucketAntiDDosInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s ListBucketAntiDDosInfoResponse) GoString() string {
	return s.String()
}

func (s *ListBucketAntiDDosInfoResponse) SetHeaders(v map[string]*string) *ListBucketAntiDDosInfoResponse {
	s.Headers = v
	return s
}

func (s *ListBucketAntiDDosInfoResponse) SetStatusCode(v int32) *ListBucketAntiDDosInfoResponse {
	s.StatusCode = &v
	return s
}

func (s *ListBucketAntiDDosInfoResponse) SetBody(v *ListBucketAntiDDosInfoResponseBody) *ListBucketAntiDDosInfoResponse {
	s.Body = v
	return s
}

type ListBucketDataRedundancyTransitionResponseBody struct {
	ListBucketDataRedundancyTransition *ListBucketDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition `json:"ListBucketDataRedundancyTransition,omitempty" xml:"ListBucketDataRedundancyTransition,omitempty" type:"Struct"`
}

func (s ListBucketDataRedundancyTransitionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListBucketDataRedundancyTransitionResponseBody) GoString() string {
	return s.String()
}

func (s *ListBucketDataRedundancyTransitionResponseBody) SetListBucketDataRedundancyTransition(v *ListBucketDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) *ListBucketDataRedundancyTransitionResponseBody {
	s.ListBucketDataRedundancyTransition = v
	return s
}

type ListBucketDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition struct {
	BucketDataRedundancyTransition *BucketDataRedundancyTransition `json:"BucketDataRedundancyTransition,omitempty" xml:"BucketDataRedundancyTransition,omitempty"`
}

func (s ListBucketDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) String() string {
	return tea.Prettify(s)
}

func (s ListBucketDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) GoString() string {
	return s.String()
}

func (s *ListBucketDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) SetBucketDataRedundancyTransition(v *BucketDataRedundancyTransition) *ListBucketDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition {
	s.BucketDataRedundancyTransition = v
	return s
}

type ListBucketDataRedundancyTransitionResponse struct {
	Headers    map[string]*string                              `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                          `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListBucketDataRedundancyTransitionResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListBucketDataRedundancyTransitionResponse) String() string {
	return tea.Prettify(s)
}

func (s ListBucketDataRedundancyTransitionResponse) GoString() string {
	return s.String()
}

func (s *ListBucketDataRedundancyTransitionResponse) SetHeaders(v map[string]*string) *ListBucketDataRedundancyTransitionResponse {
	s.Headers = v
	return s
}

func (s *ListBucketDataRedundancyTransitionResponse) SetStatusCode(v int32) *ListBucketDataRedundancyTransitionResponse {
	s.StatusCode = &v
	return s
}

func (s *ListBucketDataRedundancyTransitionResponse) SetBody(v *ListBucketDataRedundancyTransitionResponseBody) *ListBucketDataRedundancyTransitionResponse {
	s.Body = v
	return s
}

type ListBucketInventoryRequest struct {
	// Specify the start position of the list operation. You can obtain this token from the NextContinuationToken field of last ListBucketInventory\"s result.
	ContinuationToken *string `json:"continuation-token,omitempty" xml:"continuation-token,omitempty"`
}

func (s ListBucketInventoryRequest) String() string {
	return tea.Prettify(s)
}

func (s ListBucketInventoryRequest) GoString() string {
	return s.String()
}

func (s *ListBucketInventoryRequest) SetContinuationToken(v string) *ListBucketInventoryRequest {
	s.ContinuationToken = &v
	return s
}

type ListBucketInventoryResponseBody struct {
	// The container that stores inventory configuration list.
	ListInventoryConfigurationsResult *ListBucketInventoryResponseBodyListInventoryConfigurationsResult `json:"ListInventoryConfigurationsResult,omitempty" xml:"ListInventoryConfigurationsResult,omitempty" type:"Struct"`
}

func (s ListBucketInventoryResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListBucketInventoryResponseBody) GoString() string {
	return s.String()
}

func (s *ListBucketInventoryResponseBody) SetListInventoryConfigurationsResult(v *ListBucketInventoryResponseBodyListInventoryConfigurationsResult) *ListBucketInventoryResponseBody {
	s.ListInventoryConfigurationsResult = v
	return s
}

type ListBucketInventoryResponseBodyListInventoryConfigurationsResult struct {
	// The container that stores inventory configurations.
	InventoryConfigurations []*InventoryConfiguration `json:"InventoryConfiguration,omitempty" xml:"InventoryConfiguration,omitempty" type:"Repeated"`
	// Specifies whether to list all inventory tasks configured for the bucket.
	// Valid values: true and false
	// - The value of false indicates that all inventory tasks configured for the bucket are listed.
	// - The value of true indicates that not all inventory tasks configured for the bucket are listed. To list the next page of inventory configurations, set the continuation-token parameter in the next request to the value of the NextContinuationToken header in the response to the current request.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// If the value of IsTruncated in the response is true and value of this header is not null, set the continuation-token parameter in the next request to the value of this header.
	NextContinuationToken *string `json:"NextContinuationToken,omitempty" xml:"NextContinuationToken,omitempty"`
}

func (s ListBucketInventoryResponseBodyListInventoryConfigurationsResult) String() string {
	return tea.Prettify(s)
}

func (s ListBucketInventoryResponseBodyListInventoryConfigurationsResult) GoString() string {
	return s.String()
}

func (s *ListBucketInventoryResponseBodyListInventoryConfigurationsResult) SetInventoryConfigurations(v []*InventoryConfiguration) *ListBucketInventoryResponseBodyListInventoryConfigurationsResult {
	s.InventoryConfigurations = v
	return s
}

func (s *ListBucketInventoryResponseBodyListInventoryConfigurationsResult) SetIsTruncated(v bool) *ListBucketInventoryResponseBodyListInventoryConfigurationsResult {
	s.IsTruncated = &v
	return s
}

func (s *ListBucketInventoryResponseBodyListInventoryConfigurationsResult) SetNextContinuationToken(v string) *ListBucketInventoryResponseBodyListInventoryConfigurationsResult {
	s.NextContinuationToken = &v
	return s
}

type ListBucketInventoryResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListBucketInventoryResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListBucketInventoryResponse) String() string {
	return tea.Prettify(s)
}

func (s ListBucketInventoryResponse) GoString() string {
	return s.String()
}

func (s *ListBucketInventoryResponse) SetHeaders(v map[string]*string) *ListBucketInventoryResponse {
	s.Headers = v
	return s
}

func (s *ListBucketInventoryResponse) SetStatusCode(v int32) *ListBucketInventoryResponse {
	s.StatusCode = &v
	return s
}

func (s *ListBucketInventoryResponse) SetBody(v *ListBucketInventoryResponseBody) *ListBucketInventoryResponse {
	s.Body = v
	return s
}

type ListBucketsHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ID of the resource group to which the bucket belongs.
	XOssResourceGroupId *string `json:"x-oss-resource-group-id,omitempty" xml:"x-oss-resource-group-id,omitempty"`
}

func (s ListBucketsHeaders) String() string {
	return tea.Prettify(s)
}

func (s ListBucketsHeaders) GoString() string {
	return s.String()
}

func (s *ListBucketsHeaders) SetCommonHeaders(v map[string]*string) *ListBucketsHeaders {
	s.CommonHeaders = v
	return s
}

func (s *ListBucketsHeaders) SetXOssResourceGroupId(v string) *ListBucketsHeaders {
	s.XOssResourceGroupId = &v
	return s
}

type ListBucketsRequest struct {
	// The name of the bucket from which the buckets start to return. The buckets whose names are alphabetically after the value of marker are returned. If this parameter is not specified, all results are returned. By default, this parameter is left empty.
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
	// The maximum number of buckets that can be returned. Valid values: 1 to 1000. Default value: 100
	MaxKeys *int64 `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
	// The prefix that the names of returned buckets must contain. If this parameter is not specified, prefixes are not used to filter returned buckets. By default, this parameter is left empty.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
}

func (s ListBucketsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListBucketsRequest) GoString() string {
	return s.String()
}

func (s *ListBucketsRequest) SetMarker(v string) *ListBucketsRequest {
	s.Marker = &v
	return s
}

func (s *ListBucketsRequest) SetMaxKeys(v int64) *ListBucketsRequest {
	s.MaxKeys = &v
	return s
}

func (s *ListBucketsRequest) SetPrefix(v string) *ListBucketsRequest {
	s.Prefix = &v
	return s
}

type ListBucketsResponseBody struct {
	// The container that stores the result of ListBuckets(GetService) request.
	ListAllMyBucketsResult *ListBucketsResponseBodyListAllMyBucketsResult `json:"ListAllMyBucketsResult,omitempty" xml:"ListAllMyBucketsResult,omitempty" type:"Struct"`
}

func (s ListBucketsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListBucketsResponseBody) GoString() string {
	return s.String()
}

func (s *ListBucketsResponseBody) SetListAllMyBucketsResult(v *ListBucketsResponseBodyListAllMyBucketsResult) *ListBucketsResponseBody {
	s.ListAllMyBucketsResult = v
	return s
}

type ListBucketsResponseBodyListAllMyBucketsResult struct {
	// The container that stores the information about multiple buckets.
	Buckets *ListBucketsResponseBodyListAllMyBucketsResultBuckets `json:"Buckets,omitempty" xml:"Buckets,omitempty" type:"Struct"`
	// Indicates whether all results are returned. Valid values:
	// - true: All results are not returned in the response.
	// - false: All results are returned in the response.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// The name of the bucket from which the buckets are returned.
	Marker *string `json:"Marker,omitempty" xml:"Marker,omitempty"`
	// The maximum number of buckets that can be returned.
	MaxKeys *int64 `json:"MaxKeys,omitempty" xml:"MaxKeys,omitempty"`
	// The marker for the next ListBuckets (GetService) request. You can use the value of this parameter as the value of marker in the next ListBuckets (GetService) request to retrieve the unreturned results.
	NextMarker *string `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	// The container that stores the information about the bucket owner.
	Owner *Owner `json:"Owner,omitempty" xml:"Owner,omitempty"`
	// The prefix contained in the names of returned buckets.
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
}

func (s ListBucketsResponseBodyListAllMyBucketsResult) String() string {
	return tea.Prettify(s)
}

func (s ListBucketsResponseBodyListAllMyBucketsResult) GoString() string {
	return s.String()
}

func (s *ListBucketsResponseBodyListAllMyBucketsResult) SetBuckets(v *ListBucketsResponseBodyListAllMyBucketsResultBuckets) *ListBucketsResponseBodyListAllMyBucketsResult {
	s.Buckets = v
	return s
}

func (s *ListBucketsResponseBodyListAllMyBucketsResult) SetIsTruncated(v bool) *ListBucketsResponseBodyListAllMyBucketsResult {
	s.IsTruncated = &v
	return s
}

func (s *ListBucketsResponseBodyListAllMyBucketsResult) SetMarker(v string) *ListBucketsResponseBodyListAllMyBucketsResult {
	s.Marker = &v
	return s
}

func (s *ListBucketsResponseBodyListAllMyBucketsResult) SetMaxKeys(v int64) *ListBucketsResponseBodyListAllMyBucketsResult {
	s.MaxKeys = &v
	return s
}

func (s *ListBucketsResponseBodyListAllMyBucketsResult) SetNextMarker(v string) *ListBucketsResponseBodyListAllMyBucketsResult {
	s.NextMarker = &v
	return s
}

func (s *ListBucketsResponseBodyListAllMyBucketsResult) SetOwner(v *Owner) *ListBucketsResponseBodyListAllMyBucketsResult {
	s.Owner = v
	return s
}

func (s *ListBucketsResponseBodyListAllMyBucketsResult) SetPrefix(v string) *ListBucketsResponseBodyListAllMyBucketsResult {
	s.Prefix = &v
	return s
}

type ListBucketsResponseBodyListAllMyBucketsResultBuckets struct {
	// The container that stores the information list of multiple buckets.
	Bucket []*Bucket `json:"Bucket,omitempty" xml:"Bucket,omitempty" type:"Repeated"`
}

func (s ListBucketsResponseBodyListAllMyBucketsResultBuckets) String() string {
	return tea.Prettify(s)
}

func (s ListBucketsResponseBodyListAllMyBucketsResultBuckets) GoString() string {
	return s.String()
}

func (s *ListBucketsResponseBodyListAllMyBucketsResultBuckets) SetBucket(v []*Bucket) *ListBucketsResponseBodyListAllMyBucketsResultBuckets {
	s.Bucket = v
	return s
}

type ListBucketsResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListBucketsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListBucketsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListBucketsResponse) GoString() string {
	return s.String()
}

func (s *ListBucketsResponse) SetHeaders(v map[string]*string) *ListBucketsResponse {
	s.Headers = v
	return s
}

func (s *ListBucketsResponse) SetStatusCode(v int32) *ListBucketsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListBucketsResponse) SetBody(v *ListBucketsResponseBody) *ListBucketsResponse {
	s.Body = v
	return s
}

type ListCnameResponseBody struct {
	// The container that is used to query information about all CNAME records.
	ListCnameResult *ListCnameResponseBodyListCnameResult `json:"ListCnameResult,omitempty" xml:"ListCnameResult,omitempty" type:"Struct"`
}

func (s ListCnameResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListCnameResponseBody) GoString() string {
	return s.String()
}

func (s *ListCnameResponseBody) SetListCnameResult(v *ListCnameResponseBodyListCnameResult) *ListCnameResponseBody {
	s.ListCnameResult = v
	return s
}

type ListCnameResponseBodyListCnameResult struct {
	// The name of the bucket to which the CNAME records you want to query are mapped.
	Bucket *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	// The container that is used to store the information about all CNAME records.
	Cname []*CnameInfo `json:"Cname,omitempty" xml:"Cname,omitempty" type:"Repeated"`
	// The name of the bucket owner.
	Owner *string `json:"Owner,omitempty" xml:"Owner,omitempty"`
}

func (s ListCnameResponseBodyListCnameResult) String() string {
	return tea.Prettify(s)
}

func (s ListCnameResponseBodyListCnameResult) GoString() string {
	return s.String()
}

func (s *ListCnameResponseBodyListCnameResult) SetBucket(v string) *ListCnameResponseBodyListCnameResult {
	s.Bucket = &v
	return s
}

func (s *ListCnameResponseBodyListCnameResult) SetCname(v []*CnameInfo) *ListCnameResponseBodyListCnameResult {
	s.Cname = v
	return s
}

func (s *ListCnameResponseBodyListCnameResult) SetOwner(v string) *ListCnameResponseBodyListCnameResult {
	s.Owner = &v
	return s
}

type ListCnameResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListCnameResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListCnameResponse) String() string {
	return tea.Prettify(s)
}

func (s ListCnameResponse) GoString() string {
	return s.String()
}

func (s *ListCnameResponse) SetHeaders(v map[string]*string) *ListCnameResponse {
	s.Headers = v
	return s
}

func (s *ListCnameResponse) SetStatusCode(v int32) *ListCnameResponse {
	s.StatusCode = &v
	return s
}

func (s *ListCnameResponse) SetBody(v *ListCnameResponseBody) *ListCnameResponse {
	s.Body = v
	return s
}

type ListLiveChannelRequest struct {
	// The name of the LiveChannel from which the list operation starts. LiveChannels whose names are alphabetically after the value of the marker parameter are returned.
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
	// The maximum number of LiveChannels that can be returned for the current request. The value of max-keys cannot exceed 1000.
	// Default value: 100.
	MaxKeys *int64 `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
	// The prefix that the names of the LiveChannels that you want to return must contain. If you specify a prefix in the request, the specified prefix is included in the response.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
}

func (s ListLiveChannelRequest) String() string {
	return tea.Prettify(s)
}

func (s ListLiveChannelRequest) GoString() string {
	return s.String()
}

func (s *ListLiveChannelRequest) SetMarker(v string) *ListLiveChannelRequest {
	s.Marker = &v
	return s
}

func (s *ListLiveChannelRequest) SetMaxKeys(v int64) *ListLiveChannelRequest {
	s.MaxKeys = &v
	return s
}

func (s *ListLiveChannelRequest) SetPrefix(v string) *ListLiveChannelRequest {
	s.Prefix = &v
	return s
}

type ListLiveChannelResponseBody struct {
	// The container that stores the results of the ListLiveChannel request.
	ListLiveChannelResult *ListLiveChannelResponseBodyListLiveChannelResult `json:"ListLiveChannelResult,omitempty" xml:"ListLiveChannelResult,omitempty" type:"Struct"`
}

func (s ListLiveChannelResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListLiveChannelResponseBody) GoString() string {
	return s.String()
}

func (s *ListLiveChannelResponseBody) SetListLiveChannelResult(v *ListLiveChannelResponseBodyListLiveChannelResult) *ListLiveChannelResponseBody {
	s.ListLiveChannelResult = v
	return s
}

type ListLiveChannelResponseBodyListLiveChannelResult struct {
	// Indicates whether all results are returned.
	// - true: All results are returned.
	// - false: Not all results are returned.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// The container that stores the information about each returned LiveChannel.
	LiveChannels []*LiveChannel `json:"LiveChannel,omitempty" xml:"LiveChannel,omitempty" type:"Repeated"`
	// The name of the LiveChannel after which the ListLiveChannel operation starts.
	Marker *string `json:"Marker,omitempty" xml:"Marker,omitempty"`
	// The maximum number of returned LiveChannels in the response.
	MaxKeys *int64 `json:"MaxKeys,omitempty" xml:"MaxKeys,omitempty"`
	// If not all results are returned, the NextMarker parameter is included in the response to indicate the Marker value of the next request.
	NextMarker *string `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	// The prefix that the names of the returned LiveChannels contain.
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
}

func (s ListLiveChannelResponseBodyListLiveChannelResult) String() string {
	return tea.Prettify(s)
}

func (s ListLiveChannelResponseBodyListLiveChannelResult) GoString() string {
	return s.String()
}

func (s *ListLiveChannelResponseBodyListLiveChannelResult) SetIsTruncated(v bool) *ListLiveChannelResponseBodyListLiveChannelResult {
	s.IsTruncated = &v
	return s
}

func (s *ListLiveChannelResponseBodyListLiveChannelResult) SetLiveChannels(v []*LiveChannel) *ListLiveChannelResponseBodyListLiveChannelResult {
	s.LiveChannels = v
	return s
}

func (s *ListLiveChannelResponseBodyListLiveChannelResult) SetMarker(v string) *ListLiveChannelResponseBodyListLiveChannelResult {
	s.Marker = &v
	return s
}

func (s *ListLiveChannelResponseBodyListLiveChannelResult) SetMaxKeys(v int64) *ListLiveChannelResponseBodyListLiveChannelResult {
	s.MaxKeys = &v
	return s
}

func (s *ListLiveChannelResponseBodyListLiveChannelResult) SetNextMarker(v string) *ListLiveChannelResponseBodyListLiveChannelResult {
	s.NextMarker = &v
	return s
}

func (s *ListLiveChannelResponseBodyListLiveChannelResult) SetPrefix(v string) *ListLiveChannelResponseBodyListLiveChannelResult {
	s.Prefix = &v
	return s
}

type ListLiveChannelResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListLiveChannelResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListLiveChannelResponse) String() string {
	return tea.Prettify(s)
}

func (s ListLiveChannelResponse) GoString() string {
	return s.String()
}

func (s *ListLiveChannelResponse) SetHeaders(v map[string]*string) *ListLiveChannelResponse {
	s.Headers = v
	return s
}

func (s *ListLiveChannelResponse) SetStatusCode(v int32) *ListLiveChannelResponse {
	s.StatusCode = &v
	return s
}

func (s *ListLiveChannelResponse) SetBody(v *ListLiveChannelResponseBody) *ListLiveChannelResponse {
	s.Body = v
	return s
}

type ListMultipartUploadsRequest struct {
	// The character used to group objects by name. Objects whose names contain the same string that ranges from the specified prefix to the delimiter that appears for the first time are grouped as a CommonPrefixes element.
	Delimiter *string `json:"delimiter,omitempty" xml:"delimiter,omitempty"`
	// The encoding type of the object name in the response. Values of Delimiter, KeyMarker, Prefix, NextKeyMarker, and Key can be encoded in UTF-8. However, the XML 1.0 standard cannot be used to parse control characters such as characters with an American Standard Code for Information Interchange (ASCII) value from 0 to 10. You can set the encoding-type parameter to encode values of Delimiter, KeyMarker, Prefix, NextKeyMarker, and Key in the response.
	//
	// Default value: null
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
	// This parameter is used together with the upload-id-marker parameter to specify the position from which the next list begins.
	//
	// - If the upload-id-marker parameter is not set, Object Storage Service (OSS) returns all multipart upload tasks in which object names are alphabetically after the key-marker value.
	//
	// - If the upload-id-marker parameter is set, the response includes the following tasks:
	//
	//   - Multipart upload tasks in which object names are alphabetically after the key-marker value in alphabetical order
	//
	//   - Multipart upload tasks in which object names are the same as the key-marker parameter value but whose upload IDs are greater than the upload-id-marker parameter value
	KeyMarker *string `json:"key-marker,omitempty" xml:"key-marker,omitempty"`
	// The maximumnumber of multipart upload tasks that can be returned for the current request. Default value: 1000. Maximum value: 1000.
	MaxUploads *int64 `json:"max-uploads,omitempty" xml:"max-uploads,omitempty"`
	// The prefix that the returned object names must contain. If you specify a prefix in the request, the specified prefix is included in the response.
	//
	// >You can use prefixes to group and manage objects in buckets in the same way you manage a folder in a file system.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The upload ID of the multipart upload task after which the list begins. This parameter is used together with the key-marker parameter.
	//
	// - If the key-marker parameter is not set, OSS ignores the upload-id-marker parameter.
	//
	// - If the key-marker parameter is configured, the query result includes:
	//
	//   - Multipart upload tasks in which object names are alphabetically after the key-marker value in alphabetical order
	//
	//   - Multipart upload tasks in which object names are the same as the key-marker parameter value but whose upload IDs are greater than the upload-id-marker parameter value
	UploadIdMarker *string `json:"upload-id-marker,omitempty" xml:"upload-id-marker,omitempty"`
}

func (s ListMultipartUploadsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListMultipartUploadsRequest) GoString() string {
	return s.String()
}

func (s *ListMultipartUploadsRequest) SetDelimiter(v string) *ListMultipartUploadsRequest {
	s.Delimiter = &v
	return s
}

func (s *ListMultipartUploadsRequest) SetEncodingType(v string) *ListMultipartUploadsRequest {
	s.EncodingType = &v
	return s
}

func (s *ListMultipartUploadsRequest) SetKeyMarker(v string) *ListMultipartUploadsRequest {
	s.KeyMarker = &v
	return s
}

func (s *ListMultipartUploadsRequest) SetMaxUploads(v int64) *ListMultipartUploadsRequest {
	s.MaxUploads = &v
	return s
}

func (s *ListMultipartUploadsRequest) SetPrefix(v string) *ListMultipartUploadsRequest {
	s.Prefix = &v
	return s
}

func (s *ListMultipartUploadsRequest) SetUploadIdMarker(v string) *ListMultipartUploadsRequest {
	s.UploadIdMarker = &v
	return s
}

type ListMultipartUploadsResponseBody struct {
	// The container that stores the response to the ListMultipartUpload request.
	ListMultipartUploadsResult *ListMultipartUploadsResponseBodyListMultipartUploadsResult `json:"ListMultipartUploadsResult,omitempty" xml:"ListMultipartUploadsResult,omitempty" type:"Struct"`
}

func (s ListMultipartUploadsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListMultipartUploadsResponseBody) GoString() string {
	return s.String()
}

func (s *ListMultipartUploadsResponseBody) SetListMultipartUploadsResult(v *ListMultipartUploadsResponseBodyListMultipartUploadsResult) *ListMultipartUploadsResponseBody {
	s.ListMultipartUploadsResult = v
	return s
}

type ListMultipartUploadsResponseBodyListMultipartUploadsResult struct {
	// The name of the bucket.
	Bucket *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	// If the delimiter parameter is specified in the request, the response contains the CommonPrefixes parameter. The objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in the CommonPrefixes parameter.
	CommonPrefixes []*CommonPrefix `json:"CommonPrefixes,omitempty" xml:"CommonPrefixes,omitempty" type:"Repeated"`
	// The character used to group objects by name. If you specify the Delimiter parameter in the request, the response contains the CommonPrefixes element. Objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in
	Delimiter *string `json:"Delimiter,omitempty" xml:"Delimiter,omitempty"`
	// The method used to encode the object name in the response. If encoding-type is specified in the request, values of those elements including Delimiter, KeyMarker, Prefix, NextKeyMarker, and Key are encoded in the returned result.
	EncodingType *string `json:"EncodingType,omitempty" xml:"EncodingType,omitempty"`
	// Indicates whether the list of multipart upload tasks returned in the response is truncated. Default value: false. Valid values:
	//
	// - true: Only part of the results are returned this time.
	//
	// - false: All results are returned.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// The name of the object that corresponds to the multipart upload task after which the list begins.
	KeyMarker *string `json:"KeyMarker,omitempty" xml:"KeyMarker,omitempty"`
	// The maximum number of multipart upload tasks returned by OSS.
	MaxUploads *int64 `json:"MaxUploads,omitempty" xml:"MaxUploads,omitempty"`
	// The object name marker in the response for the next request to return the remaining results.
	NextKeyMarker *string `json:"NextKeyMarker,omitempty" xml:"NextKeyMarker,omitempty"`
	// The NextUploadMarker value that is used for the UploadMarker value in the next request if the response does not contain all required results.
	NextUploadIdMarker *string `json:"NextUploadIdMarker,omitempty" xml:"NextUploadIdMarker,omitempty"`
	// The prefix that the returned object names must contain. If you specify a prefix in the request, the specified prefix is included in the response.
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	// The ID list of the multipart upload tasks.
	Uploads []*Upload `json:"Upload,omitempty" xml:"Upload,omitempty" type:"Repeated"`
	// The upload ID of the multipart upload task after which the list begins.
	UploadIdMarker *string `json:"UploadIdMarker,omitempty" xml:"UploadIdMarker,omitempty"`
}

func (s ListMultipartUploadsResponseBodyListMultipartUploadsResult) String() string {
	return tea.Prettify(s)
}

func (s ListMultipartUploadsResponseBodyListMultipartUploadsResult) GoString() string {
	return s.String()
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetBucket(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.Bucket = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetCommonPrefixes(v []*CommonPrefix) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.CommonPrefixes = v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetDelimiter(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.Delimiter = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetEncodingType(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.EncodingType = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetIsTruncated(v bool) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.IsTruncated = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetKeyMarker(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.KeyMarker = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetMaxUploads(v int64) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.MaxUploads = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetNextKeyMarker(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.NextKeyMarker = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetNextUploadIdMarker(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.NextUploadIdMarker = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetPrefix(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.Prefix = &v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetUploads(v []*Upload) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.Uploads = v
	return s
}

func (s *ListMultipartUploadsResponseBodyListMultipartUploadsResult) SetUploadIdMarker(v string) *ListMultipartUploadsResponseBodyListMultipartUploadsResult {
	s.UploadIdMarker = &v
	return s
}

type ListMultipartUploadsResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListMultipartUploadsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListMultipartUploadsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListMultipartUploadsResponse) GoString() string {
	return s.String()
}

func (s *ListMultipartUploadsResponse) SetHeaders(v map[string]*string) *ListMultipartUploadsResponse {
	s.Headers = v
	return s
}

func (s *ListMultipartUploadsResponse) SetStatusCode(v int32) *ListMultipartUploadsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListMultipartUploadsResponse) SetBody(v *ListMultipartUploadsResponseBody) *ListMultipartUploadsResponse {
	s.Body = v
	return s
}

type ListObjectVersionsRequest struct {
	// The character that is used to group objects by name. If you specify prefix and delimiter in the request, the response contains CommonPrefixes. The objects whose name contains the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in CommonPrefixes. If you specify prefix and set delimiter to a forward slash (/), only the objects in the directory are listed. The subdirectories in the directory are returned in CommonPrefixes. Objects and subdirectories in the subdirectories are not listed.
	//
	// By default, this parameter is left empty.
	Delimiter *string `json:"delimiter,omitempty" xml:"delimiter,omitempty"`
	// The encoding type of the content in the response. By default, this parameter is left empty. Set the value to URL.
	//
	// >  The values of Delimiter, Marker, Prefix, NextMarker, and Key are UTF-8 encoded. If the value of Delimiter, Marker, Prefix, NextMarker, or Key contains a control character that is not supported by Extensible Markup Language (XML) 1.0, you can specify encoding-type to encode the value in the response.
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
	// The name of the object after which the GetBucketVersions (ListObjectVersions) operation begins. If this parameter is specified, objects whose name is alphabetically after the value of key-marker are returned. Use key-marker and version-id-marker in combination. The value of key-marker must be less than 1,024 bytes in length.
	//
	// By default, this parameter is left empty.
	//
	// >  You must also specify key-marker if you specify version-id-marker.
	KeyMarker *string `json:"key-marker,omitempty" xml:"key-marker,omitempty"`
	// The maximum number of objects to be returned. If the number of returned objects exceeds the value of max-keys, the response contains `NextKeyMarker` and `NextVersionIdMarker`. Specify the values of `NextKeyMarker` and `NextVersionIdMarker` as the markers for the next request. Valid values: 1 to 999. Default value: 100.
	MaxKeys *int64 `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
	// The prefix that the names of returned objects must contain.
	//
	// *   The value of prefix must be less than 1,024 bytes in length.
	// *   If you specify prefix, the names of the returned objects contain the prefix.
	//
	// If you set prefix to a directory name, the objects whose name starts with the prefix are listed. The returned objects consist of all objects and subdirectories in the directory.
	//
	// By default, this parameter is left empty.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The version ID of the object specified in key-marker after which the GetBucketVersions (ListObjectVersions) operation begins. The versions are returned from the latest version to the earliest version. If version-id-marker is not specified, the GetBucketVersions (ListObjectVersions) operation starts from the latest version of the object whose name is alphabetically after the value of key-marker by default.
	//
	// By default, this parameter is left empty.
	//
	// Valid values: version IDs.
	VersionIdMarker *string `json:"version-id-marker,omitempty" xml:"version-id-marker,omitempty"`
}

func (s ListObjectVersionsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListObjectVersionsRequest) GoString() string {
	return s.String()
}

func (s *ListObjectVersionsRequest) SetDelimiter(v string) *ListObjectVersionsRequest {
	s.Delimiter = &v
	return s
}

func (s *ListObjectVersionsRequest) SetEncodingType(v string) *ListObjectVersionsRequest {
	s.EncodingType = &v
	return s
}

func (s *ListObjectVersionsRequest) SetKeyMarker(v string) *ListObjectVersionsRequest {
	s.KeyMarker = &v
	return s
}

func (s *ListObjectVersionsRequest) SetMaxKeys(v int64) *ListObjectVersionsRequest {
	s.MaxKeys = &v
	return s
}

func (s *ListObjectVersionsRequest) SetPrefix(v string) *ListObjectVersionsRequest {
	s.Prefix = &v
	return s
}

func (s *ListObjectVersionsRequest) SetVersionIdMarker(v string) *ListObjectVersionsRequest {
	s.VersionIdMarker = &v
	return s
}

type ListObjectVersionsResponseBody struct {
	// The container that stores the results of the ListObjectVersions (GetBucketVersions) request.
	ListVersionsResult *ListObjectVersionsResponseBodyListVersionsResult `json:"ListVersionsResult,omitempty" xml:"ListVersionsResult,omitempty" type:"Struct"`
}

func (s ListObjectVersionsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListObjectVersionsResponseBody) GoString() string {
	return s.String()
}

func (s *ListObjectVersionsResponseBody) SetListVersionsResult(v *ListObjectVersionsResponseBodyListVersionsResult) *ListObjectVersionsResponseBody {
	s.ListVersionsResult = v
	return s
}

type ListObjectVersionsResponseBodyListVersionsResult struct {
	// Objects whose names contain the same string that ranges from the prefix to the next occurrence of the delimiter are grouped as a single result element
	CommonPrefixes []*CommonPrefix `json:"CommonPrefixes,omitempty" xml:"CommonPrefixes,omitempty" type:"Repeated"`
	// The container that stores delete markers
	DeleteMarkers []*DeleteMarkerEntry `json:"DeleteMarker,omitempty" xml:"DeleteMarker,omitempty" type:"Repeated"`
	// The character that is used to group objects by name. The objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result parameter in CommonPrefixes.
	Delimiter *string `json:"Delimiter,omitempty" xml:"Delimiter,omitempty"`
	// The encoding type of the content in the response. If you specify encoding-type in the request, the values of Delimiter, Marker, Prefix, NextMarker, and Key are encoded in the response.
	EncodingType *string `json:"EncodingType,omitempty" xml:"EncodingType,omitempty"`
	// Indicates whether the returned results are truncated.
	//
	// - true: indicates that not all results are returned for the request.
	// - false: indicates that all results are returned for the request.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// Indicates the object from which the ListObjectVersions (GetBucketVersions) operation starts.
	KeyMarker *string `json:"KeyMarker,omitempty" xml:"KeyMarker,omitempty"`
	// The maximum number of objects that can be returned in the response.
	MaxKeys *int64 `json:"MaxKeys,omitempty" xml:"MaxKeys,omitempty"`
	// The bucket name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// If not all results are returned for the request, the NextKeyMarker parameter is included in the response to indicate the key-marker value of the next ListObjectVersions (GetBucketVersions) request.
	NextKeyMarker *string `json:"NextKeyMarker,omitempty" xml:"NextKeyMarker,omitempty"`
	// If not all results are returned for the request, the NextVersionIdMarker parameter is included in the response to indicate the version-id-marker value of the next ListObjectVersions (GetBucketVersions) request.
	NextVersionIdMarker *string `json:"NextVersionIdMarker,omitempty" xml:"NextVersionIdMarker,omitempty"`
	// The prefix contained in the names of the returned objects.
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	// The container that stores the versions of objects except for delete markers
	Versions []*ObjectVersion `json:"Version,omitempty" xml:"Version,omitempty" type:"Repeated"`
	// The version from which the ListObjectVersions (GetBucketVersions) operation starts. This parameter is used together with KeyMarker.
	VersionIdMarker *string `json:"VersionIdMarker,omitempty" xml:"VersionIdMarker,omitempty"`
}

func (s ListObjectVersionsResponseBodyListVersionsResult) String() string {
	return tea.Prettify(s)
}

func (s ListObjectVersionsResponseBodyListVersionsResult) GoString() string {
	return s.String()
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetCommonPrefixes(v []*CommonPrefix) *ListObjectVersionsResponseBodyListVersionsResult {
	s.CommonPrefixes = v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetDeleteMarkers(v []*DeleteMarkerEntry) *ListObjectVersionsResponseBodyListVersionsResult {
	s.DeleteMarkers = v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetDelimiter(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.Delimiter = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetEncodingType(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.EncodingType = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetIsTruncated(v bool) *ListObjectVersionsResponseBodyListVersionsResult {
	s.IsTruncated = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetKeyMarker(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.KeyMarker = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetMaxKeys(v int64) *ListObjectVersionsResponseBodyListVersionsResult {
	s.MaxKeys = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetName(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.Name = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetNextKeyMarker(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.NextKeyMarker = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetNextVersionIdMarker(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.NextVersionIdMarker = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetPrefix(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.Prefix = &v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetVersions(v []*ObjectVersion) *ListObjectVersionsResponseBodyListVersionsResult {
	s.Versions = v
	return s
}

func (s *ListObjectVersionsResponseBodyListVersionsResult) SetVersionIdMarker(v string) *ListObjectVersionsResponseBodyListVersionsResult {
	s.VersionIdMarker = &v
	return s
}

type ListObjectVersionsResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListObjectVersionsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListObjectVersionsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListObjectVersionsResponse) GoString() string {
	return s.String()
}

func (s *ListObjectVersionsResponse) SetHeaders(v map[string]*string) *ListObjectVersionsResponse {
	s.Headers = v
	return s
}

func (s *ListObjectVersionsResponse) SetStatusCode(v int32) *ListObjectVersionsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListObjectVersionsResponse) SetBody(v *ListObjectVersionsResponseBody) *ListObjectVersionsResponse {
	s.Body = v
	return s
}

type ListObjectsRequest struct {
	// The character that is used to group objects by name. If you specify delimiter in the request, the response contains CommonPrefixes. The objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in CommonPrefixes.
	Delimiter *string `json:"delimiter,omitempty" xml:"delimiter,omitempty"`
	// The encoding format of the content in the response.
	//
	// >  The value of Delimiter, Marker, Prefix, NextMarker, and Key are UTF-8 encoded. If the values of Delimiter, Marker, Prefix, NextMarker, and Key contain a control character that is not supported by Extensible Markup Language (XML) 1.0, you can specify encoding-type to encode the value in the response.
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
	// The name of the object after which the GetBucket (ListObjects) operation begins. If this parameter is specified, objects whose names are alphabetically after the value of marker are returned.\
	// The objects are returned by page based on marker. The value of marker can be up to 1,024 bytes.\
	// If the value of marker does not exist in the list when you perform a conditional query, the GetBucket (ListObjects) operation starts from the object whose name is alphabetically after the value of marker.
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
	// The maximum number of objects that can be returned. If the number of objects to be returned exceeds the value of max-keys specified in the request, NextMarker is included in the returned response. The value of NextMarker is used as the value of marker for the next request.\
	// Valid values: 1 to 999.\
	// Default value: 100.
	MaxKeys *int64 `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
	// The prefix that must be contained in names of the returned objects.
	//
	// *   The value of prefix can be up to 1,024 bytes in length.
	// *   If you specify prefix, the names of the returned objects contain the prefix.
	//
	// If you set prefix to a directory name, the object whose names start with this prefix are listed. The objects consist of all recursive objects and subdirectories in this directory.\
	// If you set prefix to a directory name and set delimiter to a forward slash (/), only the objects in the directory are listed. The subdirectories in the directory are listed in CommonPrefixes. Recursive objects and subdirectories in the subdirectories are not listed.\
	// For example, a bucket contains the following three objects: fun/test.jpg, fun/movie/001.avi, and fun/movie/007.avi. If prefix is set to fun/, the three objects are returned. If prefix is set to fun/ and delimiter is set to a forward slash (/), fun/test.jpg and fun/movie/ are returned.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
}

func (s ListObjectsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsRequest) GoString() string {
	return s.String()
}

func (s *ListObjectsRequest) SetDelimiter(v string) *ListObjectsRequest {
	s.Delimiter = &v
	return s
}

func (s *ListObjectsRequest) SetEncodingType(v string) *ListObjectsRequest {
	s.EncodingType = &v
	return s
}

func (s *ListObjectsRequest) SetMarker(v string) *ListObjectsRequest {
	s.Marker = &v
	return s
}

func (s *ListObjectsRequest) SetMaxKeys(v int64) *ListObjectsRequest {
	s.MaxKeys = &v
	return s
}

func (s *ListObjectsRequest) SetPrefix(v string) *ListObjectsRequest {
	s.Prefix = &v
	return s
}

type ListObjectsResponseBody struct {
	// The container that stores the information about the returned objects.
	ListBucketResult *ListObjectsResponseBodyListBucketResult `json:"ListBucketResult,omitempty" xml:"ListBucketResult,omitempty" type:"Struct"`
}

func (s ListObjectsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsResponseBody) GoString() string {
	return s.String()
}

func (s *ListObjectsResponseBody) SetListBucketResult(v *ListObjectsResponseBodyListBucketResult) *ListObjectsResponseBody {
	s.ListBucketResult = v
	return s
}

type ListObjectsResponseBodyListBucketResult struct {
	// If delimiter is specified in the request, the response contains CommonPrefixes. The objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in CommonPrefixes.
	CommonPrefixes []*CommonPrefix `json:"CommonPrefixes,omitempty" xml:"CommonPrefixes,omitempty" type:"Repeated"`
	// The container that stores the metadata of the returned objects.
	Contents []*ObjectSummary `json:"Contents,omitempty" xml:"Contents,omitempty" type:"Repeated"`
	// The character that is used to group objects by name. The objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in CommonPrefixes.
	Delimiter *string `json:"Delimiter,omitempty" xml:"Delimiter,omitempty"`
	// The encoding type of the content in the response. If you specify encoding-type in the request, the values of Delimiter, Marker, Prefix, NextMarker, and Key are encoded in the response.
	EncodingType *string `json:"EncodingType,omitempty" xml:"EncodingType,omitempty"`
	// Indicates whether the returned list in the result is truncated. Valid values:
	// - true
	// - false
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// The name of the object after which the GetBucket (ListObjects) operation begins.
	Marker *string `json:"Marker,omitempty" xml:"Marker,omitempty"`
	// The maximum number of returned objects in the response.
	MaxKeys *int32 `json:"MaxKeys,omitempty" xml:"MaxKeys,omitempty"`
	// The name of the bucket.
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// If not all results are returned, NextMarker is included in the response to indicate the value of marker in the next request.
	NextMarker *string `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	// The prefix in the names of the returned objects.
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
}

func (s ListObjectsResponseBodyListBucketResult) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsResponseBodyListBucketResult) GoString() string {
	return s.String()
}

func (s *ListObjectsResponseBodyListBucketResult) SetCommonPrefixes(v []*CommonPrefix) *ListObjectsResponseBodyListBucketResult {
	s.CommonPrefixes = v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetContents(v []*ObjectSummary) *ListObjectsResponseBodyListBucketResult {
	s.Contents = v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetDelimiter(v string) *ListObjectsResponseBodyListBucketResult {
	s.Delimiter = &v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetEncodingType(v string) *ListObjectsResponseBodyListBucketResult {
	s.EncodingType = &v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetIsTruncated(v bool) *ListObjectsResponseBodyListBucketResult {
	s.IsTruncated = &v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetMarker(v string) *ListObjectsResponseBodyListBucketResult {
	s.Marker = &v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetMaxKeys(v int32) *ListObjectsResponseBodyListBucketResult {
	s.MaxKeys = &v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetName(v string) *ListObjectsResponseBodyListBucketResult {
	s.Name = &v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetNextMarker(v string) *ListObjectsResponseBodyListBucketResult {
	s.NextMarker = &v
	return s
}

func (s *ListObjectsResponseBodyListBucketResult) SetPrefix(v string) *ListObjectsResponseBodyListBucketResult {
	s.Prefix = &v
	return s
}

type ListObjectsResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListObjectsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListObjectsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsResponse) GoString() string {
	return s.String()
}

func (s *ListObjectsResponse) SetHeaders(v map[string]*string) *ListObjectsResponse {
	s.Headers = v
	return s
}

func (s *ListObjectsResponse) SetStatusCode(v int32) *ListObjectsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListObjectsResponse) SetBody(v *ListObjectsResponseBody) *ListObjectsResponse {
	s.Body = v
	return s
}

type ListObjectsV2Request struct {
	// The token from which the list operation starts. You can obtain the token from NextContinuationToken in the response of the ListObjectsV2 request.
	ContinuationToken *string `json:"continuation-token,omitempty" xml:"continuation-token,omitempty"`
	// The character that is used to group objects by name. If you specify delimiter in the request, the response contains CommonPrefixes. The objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in CommonPrefixes.
	Delimiter *string `json:"delimiter,omitempty" xml:"delimiter,omitempty"`
	// The encoding format of the returned objects in the response.
	//
	// >  The values of Delimiter, StartAfter, Prefix, NextContinuationToken, and Key are UTF-8 encoded. If the value of Delimiter, StartAfter, Prefix, NextContinuationToken, or Key contains a control character that is not supported by Extensible Markup Language (XML) 1.0, you can specify encoding-type to encode the value in the response.
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
	// Specifies whether to include the information about the bucket owner in the response. Valid values:
	//
	// *   true
	// *   false
	FetchOwner *bool `json:"fetch-owner,omitempty" xml:"fetch-owner,omitempty"`
	// The maximum number of objects to be returned.\
	// Valid values: 1 to 999.\
	// Default value: 100.
	//
	// >  If the number of returned objects exceeds the value of max-keys, the response contains NextContinuationToken.Use the value of NextContinuationToken as the value of continuation-token in the next request.
	MaxKeys *int64 `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
	// The prefix that must be contained in names of the returned objects.\
	//
	//
	// *   The value of prefix can be up to 1,024 bytes in length.
	// *   If you specify prefix, the names of the returned objects contain the prefix.
	//
	// If you set prefix to a directory name, the objects whose names start with this prefix are listed. The objects consist of all objects and subdirectories in this directory.\
	// If you set prefix to a directory name and set delimiter to a forward slash (/), only the objects in the directory are listed. The subdirectories in the directory are returned in CommonPrefixes. Objects and subdirectories in the subdirectories are not listed.\
	// For example, a bucket contains the following three objects: fun/test.jpg, fun/movie/001.avi, and fun/movie/007.avi. If prefix is set to fun/, the three objects are returned. If prefix is set to fun/ and delimiter is set to a forward slash (/), fun/test.jpg and fun/movie/ are returned.
	Prefix *string `json:"prefix,omitempty" xml:"prefix,omitempty"`
	// The name of the object after which the list operation begins. If this parameter is specified, objects whose names are alphabetically after the value of start-after are returned.\
	// The objects are returned by page based on start-after. The value of start-after can be up to 1,024 bytes in length.\
	// If the value of start-after does not exist when you perform a conditional query, the list starts from the object whose name is alphabetically after the value of start-after.
	StartAfter *string `json:"start-after,omitempty" xml:"start-after,omitempty"`
}

func (s ListObjectsV2Request) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsV2Request) GoString() string {
	return s.String()
}

func (s *ListObjectsV2Request) SetContinuationToken(v string) *ListObjectsV2Request {
	s.ContinuationToken = &v
	return s
}

func (s *ListObjectsV2Request) SetDelimiter(v string) *ListObjectsV2Request {
	s.Delimiter = &v
	return s
}

func (s *ListObjectsV2Request) SetEncodingType(v string) *ListObjectsV2Request {
	s.EncodingType = &v
	return s
}

func (s *ListObjectsV2Request) SetFetchOwner(v bool) *ListObjectsV2Request {
	s.FetchOwner = &v
	return s
}

func (s *ListObjectsV2Request) SetMaxKeys(v int64) *ListObjectsV2Request {
	s.MaxKeys = &v
	return s
}

func (s *ListObjectsV2Request) SetPrefix(v string) *ListObjectsV2Request {
	s.Prefix = &v
	return s
}

func (s *ListObjectsV2Request) SetStartAfter(v string) *ListObjectsV2Request {
	s.StartAfter = &v
	return s
}

type ListObjectsV2ResponseBody struct {
	// The container that stores the metadata of returned objects.
	ListBucketResult *ListObjectsV2ResponseBodyListBucketResult `json:"ListBucketResult,omitempty" xml:"ListBucketResult,omitempty" type:"Struct"`
}

func (s ListObjectsV2ResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsV2ResponseBody) GoString() string {
	return s.String()
}

func (s *ListObjectsV2ResponseBody) SetListBucketResult(v *ListObjectsV2ResponseBodyListBucketResult) *ListObjectsV2ResponseBody {
	s.ListBucketResult = v
	return s
}

type ListObjectsV2ResponseBodyListBucketResult struct {
	// Objects whose names contain the same string that ranges from the prefix to the next occurrence of the delimiter are grouped as a single result element
	CommonPrefixes []*CommonPrefix `json:"CommonPrefixes,omitempty" xml:"CommonPrefixes,omitempty" type:"Repeated"`
	// The container that stores the metadata of the returned objects.
	Contents []*ObjectSummary `json:"Contents,omitempty" xml:"Contents,omitempty" type:"Repeated"`
	// If continuation-token is specified in the request, the response contains ContinuationToken.
	ContinuationToken *string `json:"ContinuationToken,omitempty" xml:"ContinuationToken,omitempty"`
	// The character that is used to group objects by name. The objects whose names contain the same string from the prefix to the next occurrence of the delimiter are grouped as a single result element in CommonPrefixes.
	Delimiter *string `json:"Delimiter,omitempty" xml:"Delimiter,omitempty"`
	// The encoding type of the content in the response. If you specify encoding-type in the request, the values of Delimiter, StartAfter, Prefix, NextContinuationToken, and Key are encoded in the response.
	EncodingType *string `json:"EncodingType,omitempty" xml:"EncodingType,omitempty"`
	// Indicates whether the returned results are truncated. Valid values:
	//
	// - true
	// - false
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// The number of objects returned for this request. If delimiter is specified in the request, the value of KeyCount is the sum of the values of Key and CommonPrefixes.
	KeyCount *int32 `json:"KeyCount,omitempty" xml:"KeyCount,omitempty"`
	// The maximum number of returned objects in the response.
	MaxKeys *int32 `json:"MaxKeys,omitempty" xml:"MaxKeys,omitempty"`
	// The name of the bucket.
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// The token from which the next list operation starts. Use the value of NextContinuationToken as the value of continuation-token in the next request.
	NextContinuationToken *string `json:"NextContinuationToken,omitempty" xml:"NextContinuationToken,omitempty"`
	// The prefix in the names of the returned objects.
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	// If start-after is specified in the request, the response contains StartAfter.
	StartAfter *string `json:"StartAfter,omitempty" xml:"StartAfter,omitempty"`
}

func (s ListObjectsV2ResponseBodyListBucketResult) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsV2ResponseBodyListBucketResult) GoString() string {
	return s.String()
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetCommonPrefixes(v []*CommonPrefix) *ListObjectsV2ResponseBodyListBucketResult {
	s.CommonPrefixes = v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetContents(v []*ObjectSummary) *ListObjectsV2ResponseBodyListBucketResult {
	s.Contents = v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetContinuationToken(v string) *ListObjectsV2ResponseBodyListBucketResult {
	s.ContinuationToken = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetDelimiter(v string) *ListObjectsV2ResponseBodyListBucketResult {
	s.Delimiter = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetEncodingType(v string) *ListObjectsV2ResponseBodyListBucketResult {
	s.EncodingType = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetIsTruncated(v bool) *ListObjectsV2ResponseBodyListBucketResult {
	s.IsTruncated = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetKeyCount(v int32) *ListObjectsV2ResponseBodyListBucketResult {
	s.KeyCount = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetMaxKeys(v int32) *ListObjectsV2ResponseBodyListBucketResult {
	s.MaxKeys = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetName(v string) *ListObjectsV2ResponseBodyListBucketResult {
	s.Name = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetNextContinuationToken(v string) *ListObjectsV2ResponseBodyListBucketResult {
	s.NextContinuationToken = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetPrefix(v string) *ListObjectsV2ResponseBodyListBucketResult {
	s.Prefix = &v
	return s
}

func (s *ListObjectsV2ResponseBodyListBucketResult) SetStartAfter(v string) *ListObjectsV2ResponseBodyListBucketResult {
	s.StartAfter = &v
	return s
}

type ListObjectsV2Response struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListObjectsV2ResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListObjectsV2Response) String() string {
	return tea.Prettify(s)
}

func (s ListObjectsV2Response) GoString() string {
	return s.String()
}

func (s *ListObjectsV2Response) SetHeaders(v map[string]*string) *ListObjectsV2Response {
	s.Headers = v
	return s
}

func (s *ListObjectsV2Response) SetStatusCode(v int32) *ListObjectsV2Response {
	s.StatusCode = &v
	return s
}

func (s *ListObjectsV2Response) SetBody(v *ListObjectsV2ResponseBody) *ListObjectsV2Response {
	s.Body = v
	return s
}

type ListPartsRequest struct {
	// The maximum number of parts that can be returned by OSS.
	//
	// Default value: 1000.
	//
	// Maximum value: 1000.
	EncodingType *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
	// The maximum number of parts that can be returned by OSS.
	//
	// Default value: 1000.
	//
	// Maximum value: 1000.
	MaxParts *int64 `json:"max-parts,omitempty" xml:"max-parts,omitempty"`
	// The position from which the list starts. All parts whose part numbers are greater than the value of this parameter are listed.
	//
	// By default, this parameter is left empty.
	PartNumberMarker *int64 `json:"part-number-marker,omitempty" xml:"part-number-marker,omitempty"`
	// The ID of the multipart upload task.
	//
	// By default, this parameter is left empty.
	UploadId *string `json:"uploadId,omitempty" xml:"uploadId,omitempty"`
}

func (s ListPartsRequest) String() string {
	return tea.Prettify(s)
}

func (s ListPartsRequest) GoString() string {
	return s.String()
}

func (s *ListPartsRequest) SetEncodingType(v string) *ListPartsRequest {
	s.EncodingType = &v
	return s
}

func (s *ListPartsRequest) SetMaxParts(v int64) *ListPartsRequest {
	s.MaxParts = &v
	return s
}

func (s *ListPartsRequest) SetPartNumberMarker(v int64) *ListPartsRequest {
	s.PartNumberMarker = &v
	return s
}

func (s *ListPartsRequest) SetUploadId(v string) *ListPartsRequest {
	s.UploadId = &v
	return s
}

type ListPartsShrinkRequest struct {
	// The maximum number of parts that can be returned by OSS.
	//
	// Default value: 1000.
	//
	// Maximum value: 1000.
	EncodingTypeShrink *string `json:"encoding-type,omitempty" xml:"encoding-type,omitempty"`
	// The maximum number of parts that can be returned by OSS.
	//
	// Default value: 1000.
	//
	// Maximum value: 1000.
	MaxParts *int64 `json:"max-parts,omitempty" xml:"max-parts,omitempty"`
	// The position from which the list starts. All parts whose part numbers are greater than the value of this parameter are listed.
	//
	// By default, this parameter is left empty.
	PartNumberMarker *int64 `json:"part-number-marker,omitempty" xml:"part-number-marker,omitempty"`
	// The ID of the multipart upload task.
	//
	// By default, this parameter is left empty.
	UploadId *string `json:"uploadId,omitempty" xml:"uploadId,omitempty"`
}

func (s ListPartsShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s ListPartsShrinkRequest) GoString() string {
	return s.String()
}

func (s *ListPartsShrinkRequest) SetEncodingTypeShrink(v string) *ListPartsShrinkRequest {
	s.EncodingTypeShrink = &v
	return s
}

func (s *ListPartsShrinkRequest) SetMaxParts(v int64) *ListPartsShrinkRequest {
	s.MaxParts = &v
	return s
}

func (s *ListPartsShrinkRequest) SetPartNumberMarker(v int64) *ListPartsShrinkRequest {
	s.PartNumberMarker = &v
	return s
}

func (s *ListPartsShrinkRequest) SetUploadId(v string) *ListPartsShrinkRequest {
	s.UploadId = &v
	return s
}

type ListPartsResponseBody struct {
	// The container that stores the response of the ListParts request.
	ListPartResult *ListPartsResponseBodyListPartResult `json:"ListPartResult,omitempty" xml:"ListPartResult,omitempty" type:"Struct"`
}

func (s ListPartsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListPartsResponseBody) GoString() string {
	return s.String()
}

func (s *ListPartsResponseBody) SetListPartResult(v *ListPartsResponseBodyListPartResult) *ListPartsResponseBody {
	s.ListPartResult = v
	return s
}

type ListPartsResponseBodyListPartResult struct {
	// The name of the bucket.
	Bucket *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	// Indicates whether the list of parts returned in the response has been truncated. A value of true indicates that the response does not contain all required results. A value of false indicates that the response contains all required results.
	//
	// Valid values: true and false.
	IsTruncated *bool `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	// The name of the object.
	Key *string `json:"Key,omitempty" xml:"Key,omitempty"`
	// The maximum number of parts in the response.
	MaxParts *int64 `json:"MaxParts,omitempty" xml:"MaxParts,omitempty"`
	// The NextPartNumberMarker value that is used for the PartNumberMarker value in a subsequent request when the response does not contain all required results.
	NextPartNumberMarker *int64 `json:"NextPartNumberMarker,omitempty" xml:"NextPartNumberMarker,omitempty"`
	// The list of all parts.
	Part []*Part `json:"Part,omitempty" xml:"Part,omitempty" type:"Repeated"`
	// The position from which the list starts. All parts whose part numbers are greater than the value of this parameter are listed.
	PartNumberMarker *int64 `json:"PartNumberMarker,omitempty" xml:"PartNumberMarker,omitempty"`
	// The ID of the upload task.
	UploadId *string `json:"UploadId,omitempty" xml:"UploadId,omitempty"`
}

func (s ListPartsResponseBodyListPartResult) String() string {
	return tea.Prettify(s)
}

func (s ListPartsResponseBodyListPartResult) GoString() string {
	return s.String()
}

func (s *ListPartsResponseBodyListPartResult) SetBucket(v string) *ListPartsResponseBodyListPartResult {
	s.Bucket = &v
	return s
}

func (s *ListPartsResponseBodyListPartResult) SetIsTruncated(v bool) *ListPartsResponseBodyListPartResult {
	s.IsTruncated = &v
	return s
}

func (s *ListPartsResponseBodyListPartResult) SetKey(v string) *ListPartsResponseBodyListPartResult {
	s.Key = &v
	return s
}

func (s *ListPartsResponseBodyListPartResult) SetMaxParts(v int64) *ListPartsResponseBodyListPartResult {
	s.MaxParts = &v
	return s
}

func (s *ListPartsResponseBodyListPartResult) SetNextPartNumberMarker(v int64) *ListPartsResponseBodyListPartResult {
	s.NextPartNumberMarker = &v
	return s
}

func (s *ListPartsResponseBodyListPartResult) SetPart(v []*Part) *ListPartsResponseBodyListPartResult {
	s.Part = v
	return s
}

func (s *ListPartsResponseBodyListPartResult) SetPartNumberMarker(v int64) *ListPartsResponseBodyListPartResult {
	s.PartNumberMarker = &v
	return s
}

func (s *ListPartsResponseBodyListPartResult) SetUploadId(v string) *ListPartsResponseBodyListPartResult {
	s.UploadId = &v
	return s
}

type ListPartsResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListPartsResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListPartsResponse) String() string {
	return tea.Prettify(s)
}

func (s ListPartsResponse) GoString() string {
	return s.String()
}

func (s *ListPartsResponse) SetHeaders(v map[string]*string) *ListPartsResponse {
	s.Headers = v
	return s
}

func (s *ListPartsResponse) SetStatusCode(v int32) *ListPartsResponse {
	s.StatusCode = &v
	return s
}

func (s *ListPartsResponse) SetBody(v *ListPartsResponseBody) *ListPartsResponse {
	s.Body = v
	return s
}

type ListStyleResponseBody struct {
	// The container that was used to query the information about image styles.
	StyleList *ListStyleResponseBodyStyleList `json:"StyleList,omitempty" xml:"StyleList,omitempty" type:"Struct"`
}

func (s ListStyleResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListStyleResponseBody) GoString() string {
	return s.String()
}

func (s *ListStyleResponseBody) SetStyleList(v *ListStyleResponseBodyStyleList) *ListStyleResponseBody {
	s.StyleList = v
	return s
}

type ListStyleResponseBodyStyleList struct {
	// The list of styles.
	Style []*StyleInfo `json:"Style,omitempty" xml:"Style,omitempty" type:"Repeated"`
}

func (s ListStyleResponseBodyStyleList) String() string {
	return tea.Prettify(s)
}

func (s ListStyleResponseBodyStyleList) GoString() string {
	return s.String()
}

func (s *ListStyleResponseBodyStyleList) SetStyle(v []*StyleInfo) *ListStyleResponseBodyStyleList {
	s.Style = v
	return s
}

type ListStyleResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListStyleResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListStyleResponse) String() string {
	return tea.Prettify(s)
}

func (s ListStyleResponse) GoString() string {
	return s.String()
}

func (s *ListStyleResponse) SetHeaders(v map[string]*string) *ListStyleResponse {
	s.Headers = v
	return s
}

func (s *ListStyleResponse) SetStatusCode(v int32) *ListStyleResponse {
	s.StatusCode = &v
	return s
}

func (s *ListStyleResponse) SetBody(v *ListStyleResponseBody) *ListStyleResponse {
	s.Body = v
	return s
}

type ListUserDataRedundancyTransitionRequest struct {
	ContinuationToken *string `json:"continuation-token,omitempty" xml:"continuation-token,omitempty"`
	MaxKeys           *int32  `json:"max-keys,omitempty" xml:"max-keys,omitempty"`
}

func (s ListUserDataRedundancyTransitionRequest) String() string {
	return tea.Prettify(s)
}

func (s ListUserDataRedundancyTransitionRequest) GoString() string {
	return s.String()
}

func (s *ListUserDataRedundancyTransitionRequest) SetContinuationToken(v string) *ListUserDataRedundancyTransitionRequest {
	s.ContinuationToken = &v
	return s
}

func (s *ListUserDataRedundancyTransitionRequest) SetMaxKeys(v int32) *ListUserDataRedundancyTransitionRequest {
	s.MaxKeys = &v
	return s
}

type ListUserDataRedundancyTransitionResponseBody struct {
	ListBucketDataRedundancyTransition *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition `json:"ListBucketDataRedundancyTransition,omitempty" xml:"ListBucketDataRedundancyTransition,omitempty" type:"Struct"`
}

func (s ListUserDataRedundancyTransitionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListUserDataRedundancyTransitionResponseBody) GoString() string {
	return s.String()
}

func (s *ListUserDataRedundancyTransitionResponseBody) SetListBucketDataRedundancyTransition(v *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) *ListUserDataRedundancyTransitionResponseBody {
	s.ListBucketDataRedundancyTransition = v
	return s
}

type ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition struct {
	BucketDataRedundancyTransition []*BucketDataRedundancyTransition `json:"BucketDataRedundancyTransition,omitempty" xml:"BucketDataRedundancyTransition,omitempty" type:"Repeated"`
	IsTruncated                    *bool                             `json:"IsTruncated,omitempty" xml:"IsTruncated,omitempty"`
	NextContinuationToken          *string                           `json:"NextContinuationToken,omitempty" xml:"NextContinuationToken,omitempty"`
}

func (s ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) String() string {
	return tea.Prettify(s)
}

func (s ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) GoString() string {
	return s.String()
}

func (s *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) SetBucketDataRedundancyTransition(v []*BucketDataRedundancyTransition) *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition {
	s.BucketDataRedundancyTransition = v
	return s
}

func (s *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) SetIsTruncated(v bool) *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition {
	s.IsTruncated = &v
	return s
}

func (s *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition) SetNextContinuationToken(v string) *ListUserDataRedundancyTransitionResponseBodyListBucketDataRedundancyTransition {
	s.NextContinuationToken = &v
	return s
}

type ListUserDataRedundancyTransitionResponse struct {
	Headers    map[string]*string                            `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                                        `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListUserDataRedundancyTransitionResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListUserDataRedundancyTransitionResponse) String() string {
	return tea.Prettify(s)
}

func (s ListUserDataRedundancyTransitionResponse) GoString() string {
	return s.String()
}

func (s *ListUserDataRedundancyTransitionResponse) SetHeaders(v map[string]*string) *ListUserDataRedundancyTransitionResponse {
	s.Headers = v
	return s
}

func (s *ListUserDataRedundancyTransitionResponse) SetStatusCode(v int32) *ListUserDataRedundancyTransitionResponse {
	s.StatusCode = &v
	return s
}

func (s *ListUserDataRedundancyTransitionResponse) SetBody(v *ListUserDataRedundancyTransitionResponseBody) *ListUserDataRedundancyTransitionResponse {
	s.Body = v
	return s
}

type OpenMetaQueryResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s OpenMetaQueryResponse) String() string {
	return tea.Prettify(s)
}

func (s OpenMetaQueryResponse) GoString() string {
	return s.String()
}

func (s *OpenMetaQueryResponse) SetHeaders(v map[string]*string) *OpenMetaQueryResponse {
	s.Headers = v
	return s
}

func (s *OpenMetaQueryResponse) SetStatusCode(v int32) *OpenMetaQueryResponse {
	s.StatusCode = &v
	return s
}

type OptionObjectHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The custom headers to be sent in the actual cross-origin request. You can configure multiple custom headers in a cross-origin request. Custom headers are separated by commas (,). By default, this header is left empty.
	AccessControlRequestHeaders *string `json:"Access-Control-Request-Headers,omitempty" xml:"Access-Control-Request-Headers,omitempty"`
	// The method to be used in the actual cross-origin request. You can specify only one Access-Control-Request-Method header in a cross-origin request. By default, this header is left empty.
	AccessControlRequestMethod *string `json:"Access-Control-Request-Method,omitempty" xml:"Access-Control-Request-Method,omitempty"`
	// The origin of the request. It is used to identify a cross-origin request. You can specify only one Origin header in a cross-origin request. By default, this header is left empty.
	Origin *string `json:"Origin,omitempty" xml:"Origin,omitempty"`
}

func (s OptionObjectHeaders) String() string {
	return tea.Prettify(s)
}

func (s OptionObjectHeaders) GoString() string {
	return s.String()
}

func (s *OptionObjectHeaders) SetCommonHeaders(v map[string]*string) *OptionObjectHeaders {
	s.CommonHeaders = v
	return s
}

func (s *OptionObjectHeaders) SetAccessControlRequestHeaders(v string) *OptionObjectHeaders {
	s.AccessControlRequestHeaders = &v
	return s
}

func (s *OptionObjectHeaders) SetAccessControlRequestMethod(v string) *OptionObjectHeaders {
	s.AccessControlRequestMethod = &v
	return s
}

func (s *OptionObjectHeaders) SetOrigin(v string) *OptionObjectHeaders {
	s.Origin = &v
	return s
}

type OptionObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s OptionObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s OptionObjectResponse) GoString() string {
	return s.String()
}

func (s *OptionObjectResponse) SetHeaders(v map[string]*string) *OptionObjectResponse {
	s.Headers = v
	return s
}

func (s *OptionObjectResponse) SetStatusCode(v int32) *OptionObjectResponse {
	s.StatusCode = &v
	return s
}

type PostObjectRequest struct {
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
}

func (s PostObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s PostObjectRequest) GoString() string {
	return s.String()
}

func (s *PostObjectRequest) SetKey(v string) *PostObjectRequest {
	s.Key = &v
	return s
}

type PostObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PostObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s PostObjectResponse) GoString() string {
	return s.String()
}

func (s *PostObjectResponse) SetHeaders(v map[string]*string) *PostObjectResponse {
	s.Headers = v
	return s
}

func (s *PostObjectResponse) SetStatusCode(v int32) *PostObjectResponse {
	s.StatusCode = &v
	return s
}

type PostVodPlaylistRequest struct {
	// The end time of the time range during which the TS files that you want to query are generated,
	// which is a Unix timestamp.
	// > The value of EndTime must be later than the value of StartTime. The duration between EndTime and StartTime must be shorter than one day.
	EndTime *string `json:"endTime,omitempty" xml:"endTime,omitempty"`
	// The start time of the time range during which the TS files that you want to query are generated, which is a Unix timestamp.
	StartTime *string `json:"startTime,omitempty" xml:"startTime,omitempty"`
}

func (s PostVodPlaylistRequest) String() string {
	return tea.Prettify(s)
}

func (s PostVodPlaylistRequest) GoString() string {
	return s.String()
}

func (s *PostVodPlaylistRequest) SetEndTime(v string) *PostVodPlaylistRequest {
	s.EndTime = &v
	return s
}

func (s *PostVodPlaylistRequest) SetStartTime(v string) *PostVodPlaylistRequest {
	s.StartTime = &v
	return s
}

type PostVodPlaylistResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PostVodPlaylistResponse) String() string {
	return tea.Prettify(s)
}

func (s PostVodPlaylistResponse) GoString() string {
	return s.String()
}

func (s *PostVodPlaylistResponse) SetHeaders(v map[string]*string) *PostVodPlaylistResponse {
	s.Headers = v
	return s
}

func (s *PostVodPlaylistResponse) SetStatusCode(v int32) *PostVodPlaylistResponse {
	s.StatusCode = &v
	return s
}

type PutAccessPointConfigForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s PutAccessPointConfigForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointConfigForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *PutAccessPointConfigForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *PutAccessPointConfigForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutAccessPointConfigForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *PutAccessPointConfigForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type PutAccessPointConfigForObjectProcessRequest struct {
	PutAccessPointConfigForObjectProcessConfiguration *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration `json:"PutAccessPointConfigForObjectProcessConfiguration,omitempty" xml:"PutAccessPointConfigForObjectProcessConfiguration,omitempty" type:"Struct"`
}

func (s PutAccessPointConfigForObjectProcessRequest) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointConfigForObjectProcessRequest) GoString() string {
	return s.String()
}

func (s *PutAccessPointConfigForObjectProcessRequest) SetPutAccessPointConfigForObjectProcessConfiguration(v *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration) *PutAccessPointConfigForObjectProcessRequest {
	s.PutAccessPointConfigForObjectProcessConfiguration = v
	return s
}

type PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration struct {
	AllowAnonymousAccessForObjectProcess *string                         `json:"AllowAnonymousAccessForObjectProcess,omitempty" xml:"AllowAnonymousAccessForObjectProcess,omitempty"`
	ObjectProcessConfiguration           *ObjectProcessConfiguration     `json:"ObjectProcessConfiguration,omitempty" xml:"ObjectProcessConfiguration,omitempty"`
	PublicAccessBlockConfiguration       *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
}

func (s PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration) GoString() string {
	return s.String()
}

func (s *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration) SetAllowAnonymousAccessForObjectProcess(v string) *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration {
	s.AllowAnonymousAccessForObjectProcess = &v
	return s
}

func (s *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration) SetObjectProcessConfiguration(v *ObjectProcessConfiguration) *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration {
	s.ObjectProcessConfiguration = v
	return s
}

func (s *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *PutAccessPointConfigForObjectProcessRequestPutAccessPointConfigForObjectProcessConfiguration {
	s.PublicAccessBlockConfiguration = v
	return s
}

type PutAccessPointConfigForObjectProcessResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutAccessPointConfigForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointConfigForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *PutAccessPointConfigForObjectProcessResponse) SetHeaders(v map[string]*string) *PutAccessPointConfigForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *PutAccessPointConfigForObjectProcessResponse) SetStatusCode(v int32) *PutAccessPointConfigForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

type PutAccessPointPolicyHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The name of the access point.
	XOssAccessPointName *string `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s PutAccessPointPolicyHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPolicyHeaders) GoString() string {
	return s.String()
}

func (s *PutAccessPointPolicyHeaders) SetCommonHeaders(v map[string]*string) *PutAccessPointPolicyHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutAccessPointPolicyHeaders) SetXOssAccessPointName(v string) *PutAccessPointPolicyHeaders {
	s.XOssAccessPointName = &v
	return s
}

type PutAccessPointPolicyRequest struct {
	// The configurations of the access point policy.
	Body *string `json:"body,omitempty" xml:"body,omitempty"`
}

func (s PutAccessPointPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPolicyRequest) GoString() string {
	return s.String()
}

func (s *PutAccessPointPolicyRequest) SetBody(v string) *PutAccessPointPolicyRequest {
	s.Body = &v
	return s
}

type PutAccessPointPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutAccessPointPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPolicyResponse) GoString() string {
	return s.String()
}

func (s *PutAccessPointPolicyResponse) SetHeaders(v map[string]*string) *PutAccessPointPolicyResponse {
	s.Headers = v
	return s
}

func (s *PutAccessPointPolicyResponse) SetStatusCode(v int32) *PutAccessPointPolicyResponse {
	s.StatusCode = &v
	return s
}

type PutAccessPointPolicyForObjectProcessHeaders struct {
	CommonHeaders                       map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	XOssAccessPointForObjectProcessName *string            `json:"x-oss-access-point-for-object-process-name,omitempty" xml:"x-oss-access-point-for-object-process-name,omitempty"`
}

func (s PutAccessPointPolicyForObjectProcessHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPolicyForObjectProcessHeaders) GoString() string {
	return s.String()
}

func (s *PutAccessPointPolicyForObjectProcessHeaders) SetCommonHeaders(v map[string]*string) *PutAccessPointPolicyForObjectProcessHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutAccessPointPolicyForObjectProcessHeaders) SetXOssAccessPointForObjectProcessName(v string) *PutAccessPointPolicyForObjectProcessHeaders {
	s.XOssAccessPointForObjectProcessName = &v
	return s
}

type PutAccessPointPolicyForObjectProcessRequest struct {
	Body *string `json:"body,omitempty" xml:"body,omitempty"`
}

func (s PutAccessPointPolicyForObjectProcessRequest) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPolicyForObjectProcessRequest) GoString() string {
	return s.String()
}

func (s *PutAccessPointPolicyForObjectProcessRequest) SetBody(v string) *PutAccessPointPolicyForObjectProcessRequest {
	s.Body = &v
	return s
}

type PutAccessPointPolicyForObjectProcessResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutAccessPointPolicyForObjectProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPolicyForObjectProcessResponse) GoString() string {
	return s.String()
}

func (s *PutAccessPointPolicyForObjectProcessResponse) SetHeaders(v map[string]*string) *PutAccessPointPolicyForObjectProcessResponse {
	s.Headers = v
	return s
}

func (s *PutAccessPointPolicyForObjectProcessResponse) SetStatusCode(v int32) *PutAccessPointPolicyForObjectProcessResponse {
	s.StatusCode = &v
	return s
}

type PutAccessPointPublicAccessBlockRequest struct {
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
	XOssAccessPointName            *string                         `json:"x-oss-access-point-name,omitempty" xml:"x-oss-access-point-name,omitempty"`
}

func (s PutAccessPointPublicAccessBlockRequest) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPublicAccessBlockRequest) GoString() string {
	return s.String()
}

func (s *PutAccessPointPublicAccessBlockRequest) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *PutAccessPointPublicAccessBlockRequest {
	s.PublicAccessBlockConfiguration = v
	return s
}

func (s *PutAccessPointPublicAccessBlockRequest) SetXOssAccessPointName(v string) *PutAccessPointPublicAccessBlockRequest {
	s.XOssAccessPointName = &v
	return s
}

type PutAccessPointPublicAccessBlockResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutAccessPointPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s PutAccessPointPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *PutAccessPointPublicAccessBlockResponse) SetHeaders(v map[string]*string) *PutAccessPointPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *PutAccessPointPublicAccessBlockResponse) SetStatusCode(v int32) *PutAccessPointPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

type PutBucketHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The access control list (ACL) of the bucket to be created. Valid values:
	//
	// *   public-read-write
	// *   public-read
	// *   private (default)
	//
	// For more information, see [Bucket ACL](~~31843~~).
	Acl *string `json:"x-oss-acl,omitempty" xml:"x-oss-acl,omitempty"`
	// The ID of the resource group.
	//
	// *   If you include the header in the request and specify the ID of the resource group, the bucket that you create belongs to the resource group. If the specified resource group ID is rg-default-id, the bucket that you create belongs to the default resource group.
	// *   If you do not include the header in the request, the bucket that you create belongs to the default resource group.
	//
	// You can obtain the ID of a resource group in the Resource Management console or by calling the ListResourceGroups operation. For more information, see [View basic information of a resource group](~~151181~~) and [ListResourceGroups](~~158855~~).
	//
	// >  You cannot configure a resource group for an Anywhere Bucket.
	XOssResourceGroupId *string `json:"x-oss-resource-group-id,omitempty" xml:"x-oss-resource-group-id,omitempty"`
}

func (s PutBucketHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutBucketHeaders) GoString() string {
	return s.String()
}

func (s *PutBucketHeaders) SetCommonHeaders(v map[string]*string) *PutBucketHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutBucketHeaders) SetAcl(v string) *PutBucketHeaders {
	s.Acl = &v
	return s
}

func (s *PutBucketHeaders) SetXOssResourceGroupId(v string) *PutBucketHeaders {
	s.XOssResourceGroupId = &v
	return s
}

type PutBucketRequest struct {
	// The container that stores the information about the bucket to be created.
	CreateBucketConfiguration *CreateBucketConfiguration `json:"CreateBucketConfiguration,omitempty" xml:"CreateBucketConfiguration,omitempty"`
}

func (s PutBucketRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRequest) GoString() string {
	return s.String()
}

func (s *PutBucketRequest) SetCreateBucketConfiguration(v *CreateBucketConfiguration) *PutBucketRequest {
	s.CreateBucketConfiguration = v
	return s
}

type PutBucketResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketResponse) GoString() string {
	return s.String()
}

func (s *PutBucketResponse) SetHeaders(v map[string]*string) *PutBucketResponse {
	s.Headers = v
	return s
}

func (s *PutBucketResponse) SetStatusCode(v int32) *PutBucketResponse {
	s.StatusCode = &v
	return s
}

type PutBucketAccessMonitorRequest struct {
	// The access tracking configurations of the bucket.
	AccessMonitorConfiguration *AccessMonitorConfiguration `json:"AccessMonitorConfiguration,omitempty" xml:"AccessMonitorConfiguration,omitempty"`
}

func (s PutBucketAccessMonitorRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketAccessMonitorRequest) GoString() string {
	return s.String()
}

func (s *PutBucketAccessMonitorRequest) SetAccessMonitorConfiguration(v *AccessMonitorConfiguration) *PutBucketAccessMonitorRequest {
	s.AccessMonitorConfiguration = v
	return s
}

type PutBucketAccessMonitorResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketAccessMonitorResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketAccessMonitorResponse) GoString() string {
	return s.String()
}

func (s *PutBucketAccessMonitorResponse) SetHeaders(v map[string]*string) *PutBucketAccessMonitorResponse {
	s.Headers = v
	return s
}

func (s *PutBucketAccessMonitorResponse) SetStatusCode(v int32) *PutBucketAccessMonitorResponse {
	s.StatusCode = &v
	return s
}

type PutBucketAclHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The ACL that you want to configure or modify for the bucket. The x-oss-acl header is included in PutBucketAcl requests to configure or modify the ACL of the bucket. If this header is not included, the ACL configurations do not take effect.\
	// Valid values:
	//
	// *   public-read-write: All users can read and write objects in the bucket. Exercise caution when you set the value to public-read-write.
	// *   public-read: Only the owner and authorized users of the bucket can read and write objects in the bucket. Other users can only read objects in the bucket. Exercise caution when you set the value to public-read.
	// *   private: Only the owner and authorized users of this bucket can read and write objects in the bucket. Other users cannot access objects in the bucket.
	Acl *string `json:"x-oss-acl,omitempty" xml:"x-oss-acl,omitempty"`
}

func (s PutBucketAclHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutBucketAclHeaders) GoString() string {
	return s.String()
}

func (s *PutBucketAclHeaders) SetCommonHeaders(v map[string]*string) *PutBucketAclHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutBucketAclHeaders) SetAcl(v string) *PutBucketAclHeaders {
	s.Acl = &v
	return s
}

type PutBucketAclResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketAclResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketAclResponse) GoString() string {
	return s.String()
}

func (s *PutBucketAclResponse) SetHeaders(v map[string]*string) *PutBucketAclResponse {
	s.Headers = v
	return s
}

func (s *PutBucketAclResponse) SetStatusCode(v int32) *PutBucketAclResponse {
	s.StatusCode = &v
	return s
}

type PutBucketArchiveDirectReadRequest struct {
	ArchiveDirectReadConfiguration *ArchiveDirectReadConfiguration `json:"ArchiveDirectReadConfiguration,omitempty" xml:"ArchiveDirectReadConfiguration,omitempty"`
}

func (s PutBucketArchiveDirectReadRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketArchiveDirectReadRequest) GoString() string {
	return s.String()
}

func (s *PutBucketArchiveDirectReadRequest) SetArchiveDirectReadConfiguration(v *ArchiveDirectReadConfiguration) *PutBucketArchiveDirectReadRequest {
	s.ArchiveDirectReadConfiguration = v
	return s
}

type PutBucketArchiveDirectReadResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketArchiveDirectReadResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketArchiveDirectReadResponse) GoString() string {
	return s.String()
}

func (s *PutBucketArchiveDirectReadResponse) SetHeaders(v map[string]*string) *PutBucketArchiveDirectReadResponse {
	s.Headers = v
	return s
}

func (s *PutBucketArchiveDirectReadResponse) SetStatusCode(v int32) *PutBucketArchiveDirectReadResponse {
	s.StatusCode = &v
	return s
}

type PutBucketCallbackPolicyRequest struct {
	BucketCallbackPolicy *CallbackPolicy `json:"BucketCallbackPolicy,omitempty" xml:"BucketCallbackPolicy,omitempty"`
}

func (s PutBucketCallbackPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketCallbackPolicyRequest) GoString() string {
	return s.String()
}

func (s *PutBucketCallbackPolicyRequest) SetBucketCallbackPolicy(v *CallbackPolicy) *PutBucketCallbackPolicyRequest {
	s.BucketCallbackPolicy = v
	return s
}

type PutBucketCallbackPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketCallbackPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketCallbackPolicyResponse) GoString() string {
	return s.String()
}

func (s *PutBucketCallbackPolicyResponse) SetHeaders(v map[string]*string) *PutBucketCallbackPolicyResponse {
	s.Headers = v
	return s
}

func (s *PutBucketCallbackPolicyResponse) SetStatusCode(v int32) *PutBucketCallbackPolicyResponse {
	s.StatusCode = &v
	return s
}

type PutBucketCorsRequest struct {
	// The container that stores CORS rules.
	//
	// You can configure up to 10 CORS rules for a bucket. The XML message body in a request can be up to 16 KB in size.
	CORSConfiguration *CORSConfiguration `json:"CORSConfiguration,omitempty" xml:"CORSConfiguration,omitempty"`
}

func (s PutBucketCorsRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketCorsRequest) GoString() string {
	return s.String()
}

func (s *PutBucketCorsRequest) SetCORSConfiguration(v *CORSConfiguration) *PutBucketCorsRequest {
	s.CORSConfiguration = v
	return s
}

type PutBucketCorsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketCorsResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketCorsResponse) GoString() string {
	return s.String()
}

func (s *PutBucketCorsResponse) SetHeaders(v map[string]*string) *PutBucketCorsResponse {
	s.Headers = v
	return s
}

func (s *PutBucketCorsResponse) SetStatusCode(v int32) *PutBucketCorsResponse {
	s.StatusCode = &v
	return s
}

type PutBucketEncryptionRequest struct {
	// The container that stores server-side encryption rules.
	ServerSideEncryptionRule *ServerSideEncryptionRule `json:"ServerSideEncryptionRule,omitempty" xml:"ServerSideEncryptionRule,omitempty"`
}

func (s PutBucketEncryptionRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketEncryptionRequest) GoString() string {
	return s.String()
}

func (s *PutBucketEncryptionRequest) SetServerSideEncryptionRule(v *ServerSideEncryptionRule) *PutBucketEncryptionRequest {
	s.ServerSideEncryptionRule = v
	return s
}

type PutBucketEncryptionResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketEncryptionResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketEncryptionResponse) GoString() string {
	return s.String()
}

func (s *PutBucketEncryptionResponse) SetHeaders(v map[string]*string) *PutBucketEncryptionResponse {
	s.Headers = v
	return s
}

func (s *PutBucketEncryptionResponse) SetStatusCode(v int32) *PutBucketEncryptionResponse {
	s.StatusCode = &v
	return s
}

type PutBucketHttpsConfigRequest struct {
	// The container that stores HTTPS configurations.
	HttpsConfiguration *HttpsConfiguration `json:"HttpsConfiguration,omitempty" xml:"HttpsConfiguration,omitempty"`
}

func (s PutBucketHttpsConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketHttpsConfigRequest) GoString() string {
	return s.String()
}

func (s *PutBucketHttpsConfigRequest) SetHttpsConfiguration(v *HttpsConfiguration) *PutBucketHttpsConfigRequest {
	s.HttpsConfiguration = v
	return s
}

type PutBucketHttpsConfigResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketHttpsConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketHttpsConfigResponse) GoString() string {
	return s.String()
}

func (s *PutBucketHttpsConfigResponse) SetHeaders(v map[string]*string) *PutBucketHttpsConfigResponse {
	s.Headers = v
	return s
}

func (s *PutBucketHttpsConfigResponse) SetStatusCode(v int32) *PutBucketHttpsConfigResponse {
	s.StatusCode = &v
	return s
}

type PutBucketInventoryRequest struct {
	// The container that stores the Inventory configuration.
	InventoryConfiguration *InventoryConfiguration `json:"InventoryConfiguration,omitempty" xml:"InventoryConfiguration,omitempty"`
	// The name of the inventory.
	InventoryId *string `json:"inventoryId,omitempty" xml:"inventoryId,omitempty"`
}

func (s PutBucketInventoryRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketInventoryRequest) GoString() string {
	return s.String()
}

func (s *PutBucketInventoryRequest) SetInventoryConfiguration(v *InventoryConfiguration) *PutBucketInventoryRequest {
	s.InventoryConfiguration = v
	return s
}

func (s *PutBucketInventoryRequest) SetInventoryId(v string) *PutBucketInventoryRequest {
	s.InventoryId = &v
	return s
}

type PutBucketInventoryResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketInventoryResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketInventoryResponse) GoString() string {
	return s.String()
}

func (s *PutBucketInventoryResponse) SetHeaders(v map[string]*string) *PutBucketInventoryResponse {
	s.Headers = v
	return s
}

func (s *PutBucketInventoryResponse) SetStatusCode(v int32) *PutBucketInventoryResponse {
	s.StatusCode = &v
	return s
}

type PutBucketLifecycleHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// Specifies whether to allow overlapped prefixes. Valid values:
	//
	// true: Overlapped prefixes are allowed.
	//
	// false: Overlapped prefixes are not allowed.
	XOssAllowSameActionOverlap *string `json:"x-oss-allow-same-action-overlap,omitempty" xml:"x-oss-allow-same-action-overlap,omitempty"`
}

func (s PutBucketLifecycleHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutBucketLifecycleHeaders) GoString() string {
	return s.String()
}

func (s *PutBucketLifecycleHeaders) SetCommonHeaders(v map[string]*string) *PutBucketLifecycleHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutBucketLifecycleHeaders) SetXOssAllowSameActionOverlap(v string) *PutBucketLifecycleHeaders {
	s.XOssAllowSameActionOverlap = &v
	return s
}

type PutBucketLifecycleRequest struct {
	// The container that stores the lifecycle configuration.
	LifecycleConfiguration *LifecycleConfiguration `json:"LifecycleConfiguration,omitempty" xml:"LifecycleConfiguration,omitempty"`
}

func (s PutBucketLifecycleRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketLifecycleRequest) GoString() string {
	return s.String()
}

func (s *PutBucketLifecycleRequest) SetLifecycleConfiguration(v *LifecycleConfiguration) *PutBucketLifecycleRequest {
	s.LifecycleConfiguration = v
	return s
}

type PutBucketLifecycleResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketLifecycleResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketLifecycleResponse) GoString() string {
	return s.String()
}

func (s *PutBucketLifecycleResponse) SetHeaders(v map[string]*string) *PutBucketLifecycleResponse {
	s.Headers = v
	return s
}

func (s *PutBucketLifecycleResponse) SetStatusCode(v int32) *PutBucketLifecycleResponse {
	s.StatusCode = &v
	return s
}

type PutBucketLoggingRequest struct {
	// The container that stores the logging status information.
	BucketLoggingStatus *BucketLoggingStatus `json:"BucketLoggingStatus,omitempty" xml:"BucketLoggingStatus,omitempty"`
}

func (s PutBucketLoggingRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketLoggingRequest) GoString() string {
	return s.String()
}

func (s *PutBucketLoggingRequest) SetBucketLoggingStatus(v *BucketLoggingStatus) *PutBucketLoggingRequest {
	s.BucketLoggingStatus = v
	return s
}

type PutBucketLoggingResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketLoggingResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketLoggingResponse) GoString() string {
	return s.String()
}

func (s *PutBucketLoggingResponse) SetHeaders(v map[string]*string) *PutBucketLoggingResponse {
	s.Headers = v
	return s
}

func (s *PutBucketLoggingResponse) SetStatusCode(v int32) *PutBucketLoggingResponse {
	s.StatusCode = &v
	return s
}

type PutBucketPolicyRequest struct {
	// The request parameters.
	Policy *string `json:"body,omitempty" xml:"body,omitempty"`
}

func (s PutBucketPolicyRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketPolicyRequest) GoString() string {
	return s.String()
}

func (s *PutBucketPolicyRequest) SetPolicy(v string) *PutBucketPolicyRequest {
	s.Policy = &v
	return s
}

type PutBucketPolicyResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketPolicyResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketPolicyResponse) GoString() string {
	return s.String()
}

func (s *PutBucketPolicyResponse) SetHeaders(v map[string]*string) *PutBucketPolicyResponse {
	s.Headers = v
	return s
}

func (s *PutBucketPolicyResponse) SetStatusCode(v int32) *PutBucketPolicyResponse {
	s.StatusCode = &v
	return s
}

type PutBucketPublicAccessBlockRequest struct {
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
}

func (s PutBucketPublicAccessBlockRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketPublicAccessBlockRequest) GoString() string {
	return s.String()
}

func (s *PutBucketPublicAccessBlockRequest) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *PutBucketPublicAccessBlockRequest {
	s.PublicAccessBlockConfiguration = v
	return s
}

type PutBucketPublicAccessBlockResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *PutBucketPublicAccessBlockResponse) SetHeaders(v map[string]*string) *PutBucketPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *PutBucketPublicAccessBlockResponse) SetStatusCode(v int32) *PutBucketPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

type PutBucketRedundancyTypeRequest struct {
	DataRedundancyTypeConfiguration *PutBucketRedundancyTypeRequestDataRedundancyTypeConfiguration `json:"DataRedundancyTypeConfiguration,omitempty" xml:"DataRedundancyTypeConfiguration,omitempty" type:"Struct"`
}

func (s PutBucketRedundancyTypeRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRedundancyTypeRequest) GoString() string {
	return s.String()
}

func (s *PutBucketRedundancyTypeRequest) SetDataRedundancyTypeConfiguration(v *PutBucketRedundancyTypeRequestDataRedundancyTypeConfiguration) *PutBucketRedundancyTypeRequest {
	s.DataRedundancyTypeConfiguration = v
	return s
}

type PutBucketRedundancyTypeRequestDataRedundancyTypeConfiguration struct {
	Type *string `json:"Type,omitempty" xml:"Type,omitempty"`
}

func (s PutBucketRedundancyTypeRequestDataRedundancyTypeConfiguration) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRedundancyTypeRequestDataRedundancyTypeConfiguration) GoString() string {
	return s.String()
}

func (s *PutBucketRedundancyTypeRequestDataRedundancyTypeConfiguration) SetType(v string) *PutBucketRedundancyTypeRequestDataRedundancyTypeConfiguration {
	s.Type = &v
	return s
}

type PutBucketRedundancyTypeResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketRedundancyTypeResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRedundancyTypeResponse) GoString() string {
	return s.String()
}

func (s *PutBucketRedundancyTypeResponse) SetHeaders(v map[string]*string) *PutBucketRedundancyTypeResponse {
	s.Headers = v
	return s
}

func (s *PutBucketRedundancyTypeResponse) SetStatusCode(v int32) *PutBucketRedundancyTypeResponse {
	s.StatusCode = &v
	return s
}

type PutBucketRefererRequest struct {
	// The container that stores the hotlink protection configurations.
	RefererConfiguration *RefererConfiguration `json:"RefererConfiguration,omitempty" xml:"RefererConfiguration,omitempty"`
}

func (s PutBucketRefererRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRefererRequest) GoString() string {
	return s.String()
}

func (s *PutBucketRefererRequest) SetRefererConfiguration(v *RefererConfiguration) *PutBucketRefererRequest {
	s.RefererConfiguration = v
	return s
}

type PutBucketRefererResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketRefererResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRefererResponse) GoString() string {
	return s.String()
}

func (s *PutBucketRefererResponse) SetHeaders(v map[string]*string) *PutBucketRefererResponse {
	s.Headers = v
	return s
}

func (s *PutBucketRefererResponse) SetStatusCode(v int32) *PutBucketRefererResponse {
	s.StatusCode = &v
	return s
}

type PutBucketReplicationRequest struct {
	// The container that stores data replication configurations.
	ReplicationConfiguration *ReplicationConfiguration `json:"ReplicationConfiguration,omitempty" xml:"ReplicationConfiguration,omitempty"`
}

func (s PutBucketReplicationRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketReplicationRequest) GoString() string {
	return s.String()
}

func (s *PutBucketReplicationRequest) SetReplicationConfiguration(v *ReplicationConfiguration) *PutBucketReplicationRequest {
	s.ReplicationConfiguration = v
	return s
}

type PutBucketReplicationResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketReplicationResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketReplicationResponse) GoString() string {
	return s.String()
}

func (s *PutBucketReplicationResponse) SetHeaders(v map[string]*string) *PutBucketReplicationResponse {
	s.Headers = v
	return s
}

func (s *PutBucketReplicationResponse) SetStatusCode(v int32) *PutBucketReplicationResponse {
	s.StatusCode = &v
	return s
}

type PutBucketRequestPaymentRequest struct {
	// The container that stores pay-by-requester configurations.
	RequestPaymentConfiguration *RequestPaymentConfiguration `json:"RequestPaymentConfiguration,omitempty" xml:"RequestPaymentConfiguration,omitempty"`
}

func (s PutBucketRequestPaymentRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRequestPaymentRequest) GoString() string {
	return s.String()
}

func (s *PutBucketRequestPaymentRequest) SetRequestPaymentConfiguration(v *RequestPaymentConfiguration) *PutBucketRequestPaymentRequest {
	s.RequestPaymentConfiguration = v
	return s
}

type PutBucketRequestPaymentResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketRequestPaymentResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRequestPaymentResponse) GoString() string {
	return s.String()
}

func (s *PutBucketRequestPaymentResponse) SetHeaders(v map[string]*string) *PutBucketRequestPaymentResponse {
	s.Headers = v
	return s
}

func (s *PutBucketRequestPaymentResponse) SetStatusCode(v int32) *PutBucketRequestPaymentResponse {
	s.StatusCode = &v
	return s
}

type PutBucketResourceGroupRequest struct {
	// The container that contains the ID of the resource group.
	BucketResourceGroupConfiguration *BucketResourceGroupConfiguration `json:"BucketResourceGroupConfiguration,omitempty" xml:"BucketResourceGroupConfiguration,omitempty"`
}

func (s PutBucketResourceGroupRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketResourceGroupRequest) GoString() string {
	return s.String()
}

func (s *PutBucketResourceGroupRequest) SetBucketResourceGroupConfiguration(v *BucketResourceGroupConfiguration) *PutBucketResourceGroupRequest {
	s.BucketResourceGroupConfiguration = v
	return s
}

type PutBucketResourceGroupResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketResourceGroupResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketResourceGroupResponse) GoString() string {
	return s.String()
}

func (s *PutBucketResourceGroupResponse) SetHeaders(v map[string]*string) *PutBucketResourceGroupResponse {
	s.Headers = v
	return s
}

func (s *PutBucketResourceGroupResponse) SetStatusCode(v int32) *PutBucketResourceGroupResponse {
	s.StatusCode = &v
	return s
}

type PutBucketResponseHeaderRequest struct {
	ResponseHeaderConfiguration *ResponseHeaderConfiguration `json:"ResponseHeaderConfiguration,omitempty" xml:"ResponseHeaderConfiguration,omitempty"`
}

func (s PutBucketResponseHeaderRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketResponseHeaderRequest) GoString() string {
	return s.String()
}

func (s *PutBucketResponseHeaderRequest) SetResponseHeaderConfiguration(v *ResponseHeaderConfiguration) *PutBucketResponseHeaderRequest {
	s.ResponseHeaderConfiguration = v
	return s
}

type PutBucketResponseHeaderResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketResponseHeaderResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketResponseHeaderResponse) GoString() string {
	return s.String()
}

func (s *PutBucketResponseHeaderResponse) SetHeaders(v map[string]*string) *PutBucketResponseHeaderResponse {
	s.Headers = v
	return s
}

func (s *PutBucketResponseHeaderResponse) SetStatusCode(v int32) *PutBucketResponseHeaderResponse {
	s.StatusCode = &v
	return s
}

type PutBucketRtcRequest struct {
	// The container that stores the RTC configurations.
	ReplicationRule *RtcConfiguration `json:"ReplicationRule,omitempty" xml:"ReplicationRule,omitempty"`
}

func (s PutBucketRtcRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRtcRequest) GoString() string {
	return s.String()
}

func (s *PutBucketRtcRequest) SetReplicationRule(v *RtcConfiguration) *PutBucketRtcRequest {
	s.ReplicationRule = v
	return s
}

type PutBucketRtcResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketRtcResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketRtcResponse) GoString() string {
	return s.String()
}

func (s *PutBucketRtcResponse) SetHeaders(v map[string]*string) *PutBucketRtcResponse {
	s.Headers = v
	return s
}

func (s *PutBucketRtcResponse) SetStatusCode(v int32) *PutBucketRtcResponse {
	s.StatusCode = &v
	return s
}

type PutBucketTagsRequest struct {
	// The container used to store TagSet.
	Tagging *Tagging `json:"Tagging,omitempty" xml:"Tagging,omitempty"`
}

func (s PutBucketTagsRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketTagsRequest) GoString() string {
	return s.String()
}

func (s *PutBucketTagsRequest) SetTagging(v *Tagging) *PutBucketTagsRequest {
	s.Tagging = v
	return s
}

type PutBucketTagsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketTagsResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketTagsResponse) GoString() string {
	return s.String()
}

func (s *PutBucketTagsResponse) SetHeaders(v map[string]*string) *PutBucketTagsResponse {
	s.Headers = v
	return s
}

func (s *PutBucketTagsResponse) SetStatusCode(v int32) *PutBucketTagsResponse {
	s.StatusCode = &v
	return s
}

type PutBucketTransferAccelerationRequest struct {
	// The container that stores the transfer acceleration configurations.
	TransferAccelerationConfiguration *TransferAccelerationConfiguration `json:"TransferAccelerationConfiguration,omitempty" xml:"TransferAccelerationConfiguration,omitempty"`
}

func (s PutBucketTransferAccelerationRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketTransferAccelerationRequest) GoString() string {
	return s.String()
}

func (s *PutBucketTransferAccelerationRequest) SetTransferAccelerationConfiguration(v *TransferAccelerationConfiguration) *PutBucketTransferAccelerationRequest {
	s.TransferAccelerationConfiguration = v
	return s
}

type PutBucketTransferAccelerationResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketTransferAccelerationResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketTransferAccelerationResponse) GoString() string {
	return s.String()
}

func (s *PutBucketTransferAccelerationResponse) SetHeaders(v map[string]*string) *PutBucketTransferAccelerationResponse {
	s.Headers = v
	return s
}

func (s *PutBucketTransferAccelerationResponse) SetStatusCode(v int32) *PutBucketTransferAccelerationResponse {
	s.StatusCode = &v
	return s
}

type PutBucketVersioningRequest struct {
	// The container that stores the versioning state of the bucket.
	VersioningConfiguration *VersioningConfiguration `json:"VersioningConfiguration,omitempty" xml:"VersioningConfiguration,omitempty"`
}

func (s PutBucketVersioningRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketVersioningRequest) GoString() string {
	return s.String()
}

func (s *PutBucketVersioningRequest) SetVersioningConfiguration(v *VersioningConfiguration) *PutBucketVersioningRequest {
	s.VersioningConfiguration = v
	return s
}

type PutBucketVersioningResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketVersioningResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketVersioningResponse) GoString() string {
	return s.String()
}

func (s *PutBucketVersioningResponse) SetHeaders(v map[string]*string) *PutBucketVersioningResponse {
	s.Headers = v
	return s
}

func (s *PutBucketVersioningResponse) SetStatusCode(v int32) *PutBucketVersioningResponse {
	s.StatusCode = &v
	return s
}

type PutBucketWebsiteRequest struct {
	// The container that stores the website configuration.
	WebsiteConfiguration *WebsiteConfiguration `json:"WebsiteConfiguration,omitempty" xml:"WebsiteConfiguration,omitempty"`
}

func (s PutBucketWebsiteRequest) String() string {
	return tea.Prettify(s)
}

func (s PutBucketWebsiteRequest) GoString() string {
	return s.String()
}

func (s *PutBucketWebsiteRequest) SetWebsiteConfiguration(v *WebsiteConfiguration) *PutBucketWebsiteRequest {
	s.WebsiteConfiguration = v
	return s
}

type PutBucketWebsiteResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutBucketWebsiteResponse) String() string {
	return tea.Prettify(s)
}

func (s PutBucketWebsiteResponse) GoString() string {
	return s.String()
}

func (s *PutBucketWebsiteResponse) SetHeaders(v map[string]*string) *PutBucketWebsiteResponse {
	s.Headers = v
	return s
}

func (s *PutBucketWebsiteResponse) SetStatusCode(v int32) *PutBucketWebsiteResponse {
	s.StatusCode = &v
	return s
}

type PutCnameRequest struct {
	// The container that stores the CNAME record.
	BucketCnameConfiguration *BucketCnameConfiguration `json:"BucketCnameConfiguration,omitempty" xml:"BucketCnameConfiguration,omitempty"`
}

func (s PutCnameRequest) String() string {
	return tea.Prettify(s)
}

func (s PutCnameRequest) GoString() string {
	return s.String()
}

func (s *PutCnameRequest) SetBucketCnameConfiguration(v *BucketCnameConfiguration) *PutCnameRequest {
	s.BucketCnameConfiguration = v
	return s
}

type PutCnameResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutCnameResponse) String() string {
	return tea.Prettify(s)
}

func (s PutCnameResponse) GoString() string {
	return s.String()
}

func (s *PutCnameResponse) SetHeaders(v map[string]*string) *PutCnameResponse {
	s.Headers = v
	return s
}

func (s *PutCnameResponse) SetStatusCode(v int32) *PutCnameResponse {
	s.StatusCode = &v
	return s
}

type PutLiveChannelRequest struct {
	// The container that stores the configurations of the LiveChannel.
	LiveChannelConfiguration *LiveChannelConfiguration `json:"LiveChannelConfiguration,omitempty" xml:"LiveChannelConfiguration,omitempty"`
}

func (s PutLiveChannelRequest) String() string {
	return tea.Prettify(s)
}

func (s PutLiveChannelRequest) GoString() string {
	return s.String()
}

func (s *PutLiveChannelRequest) SetLiveChannelConfiguration(v *LiveChannelConfiguration) *PutLiveChannelRequest {
	s.LiveChannelConfiguration = v
	return s
}

type PutLiveChannelResponseBody struct {
	// The container that stores the result of the CreateLiveChannel request.
	CreateLiveChannelResult *PutLiveChannelResponseBodyCreateLiveChannelResult `json:"CreateLiveChannelResult,omitempty" xml:"CreateLiveChannelResult,omitempty" type:"Struct"`
}

func (s PutLiveChannelResponseBody) String() string {
	return tea.Prettify(s)
}

func (s PutLiveChannelResponseBody) GoString() string {
	return s.String()
}

func (s *PutLiveChannelResponseBody) SetCreateLiveChannelResult(v *PutLiveChannelResponseBodyCreateLiveChannelResult) *PutLiveChannelResponseBody {
	s.CreateLiveChannelResult = v
	return s
}

type PutLiveChannelResponseBodyCreateLiveChannelResult struct {
	// 保存播放地址的容器。
	PlayUrls *LiveChannelPlayUrls `json:"PlayUrls,omitempty" xml:"PlayUrls,omitempty"`
	// 保存推流地址的容器。
	PublishUrls *LiveChannelPublishUrls `json:"PublishUrls,omitempty" xml:"PublishUrls,omitempty"`
}

func (s PutLiveChannelResponseBodyCreateLiveChannelResult) String() string {
	return tea.Prettify(s)
}

func (s PutLiveChannelResponseBodyCreateLiveChannelResult) GoString() string {
	return s.String()
}

func (s *PutLiveChannelResponseBodyCreateLiveChannelResult) SetPlayUrls(v *LiveChannelPlayUrls) *PutLiveChannelResponseBodyCreateLiveChannelResult {
	s.PlayUrls = v
	return s
}

func (s *PutLiveChannelResponseBodyCreateLiveChannelResult) SetPublishUrls(v *LiveChannelPublishUrls) *PutLiveChannelResponseBodyCreateLiveChannelResult {
	s.PublishUrls = v
	return s
}

type PutLiveChannelResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *PutLiveChannelResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s PutLiveChannelResponse) String() string {
	return tea.Prettify(s)
}

func (s PutLiveChannelResponse) GoString() string {
	return s.String()
}

func (s *PutLiveChannelResponse) SetHeaders(v map[string]*string) *PutLiveChannelResponse {
	s.Headers = v
	return s
}

func (s *PutLiveChannelResponse) SetStatusCode(v int32) *PutLiveChannelResponse {
	s.StatusCode = &v
	return s
}

func (s *PutLiveChannelResponse) SetBody(v *PutLiveChannelResponseBody) *PutLiveChannelResponse {
	s.Body = v
	return s
}

type PutLiveChannelStatusRequest struct {
	// The status of the LiveChannel.
	// Valid values:
	// - enabled: enables the LiveChannel.
	// - disabled: disables the LiveChannel.
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s PutLiveChannelStatusRequest) String() string {
	return tea.Prettify(s)
}

func (s PutLiveChannelStatusRequest) GoString() string {
	return s.String()
}

func (s *PutLiveChannelStatusRequest) SetStatus(v string) *PutLiveChannelStatusRequest {
	s.Status = &v
	return s
}

type PutLiveChannelStatusResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutLiveChannelStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s PutLiveChannelStatusResponse) GoString() string {
	return s.String()
}

func (s *PutLiveChannelStatusResponse) SetHeaders(v map[string]*string) *PutLiveChannelStatusResponse {
	s.Headers = v
	return s
}

func (s *PutLiveChannelStatusResponse) SetStatusCode(v int32) *PutLiveChannelStatusResponse {
	s.StatusCode = &v
	return s
}

type PutObjectHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// Specifies whether the object that is uploaded by calling the PutObject operation overwrites the existing object that has the same name.  When versioning is enabled or suspended for the bucket to which you want to upload the object, the **x-oss-forbid-overwrite** header does not take effect. In this case, the object that is uploaded by calling the PutObject operation overwrites the existing object that has the same name.
	//   - If you do not specify the **x-oss-forbid-overwrite** header or set the **x-oss-forbid-overwrite** header to **false**, the object that is uploaded by calling the PutObject operation overwrites the existing object that has the same name.
	//   - If the value of **x-oss-forbid-overwrite** is set to **true**, existing objects cannot be overwritten by objects that have the same names.
	//
	// If you specify the **x-oss-forbid-overwrite** request header, the queries per second (QPS) performance of OSS is degraded. If you want to use the **x-oss-forbid-overwrite** request header to perform a large number of operations (QPS greater than 1,000), contact technical support.
	// Default value: **false**.
	ForbidOverwrite *bool `json:"x-oss-forbid-overwrite,omitempty" xml:"x-oss-forbid-overwrite,omitempty"`
	// If the PutObject request contains a parameter prefixed with **x-oss-meta-***, the parameter is considered as the user metadata. Example: `x-oss-meta-location`. You can specify multiple similar headers for an object. However, the user metadata of an object cannot exceed 8 KB in size.
	//
	// The names of user metadata headers can contain letters, digits, and hyphens (-). Uppercase letters are converted to lowercase letters. Other characters such as underscores (_) are not supported.
	MetaData map[string]*string `json:"x-oss-meta-*,omitempty" xml:"x-oss-meta-*,omitempty"`
	// The access control list (ACL) of the object. Default value: default.
	// Valid values:
	//
	// - default: The ACL of the object is the same as that of the bucket in which the object is stored.
	// - private: The ACL of the object is private. Only the owner of the object and authorized users can read and write this object.
	// - public-read: The ACL of the object is public-read. Only the owner of the object and authorized users can read and write this object. Other users can only read the object. Exercise caution when you set the object ACL to this value.
	// - public-read-write: The ACL of the object is public-read-write. All users can read and write this object. Exercise caution when you set the object ACL to this value.
	//
	// For more information about the ACL, see **[ACL](~~100676~~)**.
	Acl *string `json:"x-oss-object-acl,omitempty" xml:"x-oss-object-acl,omitempty"`
	// The encryption method on the server side when an object is created.
	//
	// Valid values: **AES256**, **KMS**, and **SM4**.
	//
	// If you specify the header, the header is returned in the response. OSS uses the method that is specified by this header to encrypt the uploaded object. When you download the encrypted object, the **x-oss-server-side-encryption** header is included in the response and the header value is set to the algorithm that is used to encrypt the object.
	SseDataEncryption *string `json:"x-oss-server-side-data-encryption,omitempty" xml:"x-oss-server-side-data-encryption,omitempty"`
	// The method that is used to encrypt the object on the OSS server when the object is created.
	//
	// Valid values: **AES256**, **KMS**, and **SM4****.
	//
	// If you specify the header, the header is returned in the response. OSS uses the method that is specified by this header to encrypt the uploaded object. When you download the encrypted object, the **x-oss-server-side-encryption** header is included in the response and the header value is set to the algorithm that is used to encrypt the object.
	ServerSideEncryption *string `json:"x-oss-server-side-encryption,omitempty" xml:"x-oss-server-side-encryption,omitempty"`
	// The ID of the customer master key (CMK) managed by Key Management Service (KMS).
	// This header is valid only when the **x-oss-server-side-encryption** header is set to KMS.
	SseKeyId *string `json:"x-oss-server-side-encryption-key-id,omitempty" xml:"x-oss-server-side-encryption-key-id,omitempty"`
	// The storage class of the bucket. Default value: Standard.  Valid values:
	//
	// - Standard
	// - IA
	// - Archive
	// - ColdArchive
	StorageClass *string `json:"x-oss-storage-class,omitempty" xml:"x-oss-storage-class,omitempty"`
	// The tag of the object. You can configure multiple tags for the object. Example: TagA=A&TagB=B.
	// > The key and value of a tag must be URL-encoded. If a tag does not contain an equal sign (=), the value of the tag is considered an empty string.
	Tagging *string `json:"x-oss-tagging,omitempty" xml:"x-oss-tagging,omitempty"`
}

func (s PutObjectHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutObjectHeaders) GoString() string {
	return s.String()
}

func (s *PutObjectHeaders) SetCommonHeaders(v map[string]*string) *PutObjectHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutObjectHeaders) SetForbidOverwrite(v bool) *PutObjectHeaders {
	s.ForbidOverwrite = &v
	return s
}

func (s *PutObjectHeaders) SetMetaData(v map[string]*string) *PutObjectHeaders {
	s.MetaData = v
	return s
}

func (s *PutObjectHeaders) SetAcl(v string) *PutObjectHeaders {
	s.Acl = &v
	return s
}

func (s *PutObjectHeaders) SetSseDataEncryption(v string) *PutObjectHeaders {
	s.SseDataEncryption = &v
	return s
}

func (s *PutObjectHeaders) SetServerSideEncryption(v string) *PutObjectHeaders {
	s.ServerSideEncryption = &v
	return s
}

func (s *PutObjectHeaders) SetSseKeyId(v string) *PutObjectHeaders {
	s.SseKeyId = &v
	return s
}

func (s *PutObjectHeaders) SetStorageClass(v string) *PutObjectHeaders {
	s.StorageClass = &v
	return s
}

func (s *PutObjectHeaders) SetTagging(v string) *PutObjectHeaders {
	s.Tagging = &v
	return s
}

type PutObjectRequest struct {
	// The body of the request.
	Body io.Reader `json:"body,omitempty" xml:"body,omitempty"`
}

func (s PutObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s PutObjectRequest) GoString() string {
	return s.String()
}

func (s *PutObjectRequest) SetBody(v io.Reader) *PutObjectRequest {
	s.Body = v
	return s
}

type PutObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s PutObjectResponse) GoString() string {
	return s.String()
}

func (s *PutObjectResponse) SetHeaders(v map[string]*string) *PutObjectResponse {
	s.Headers = v
	return s
}

func (s *PutObjectResponse) SetStatusCode(v int32) *PutObjectResponse {
	s.StatusCode = &v
	return s
}

type PutObjectAclHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The access control list (ACL) of the object.
	Acl *string `json:"x-oss-object-acl,omitempty" xml:"x-oss-object-acl,omitempty"`
}

func (s PutObjectAclHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutObjectAclHeaders) GoString() string {
	return s.String()
}

func (s *PutObjectAclHeaders) SetCommonHeaders(v map[string]*string) *PutObjectAclHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutObjectAclHeaders) SetAcl(v string) *PutObjectAclHeaders {
	s.Acl = &v
	return s
}

type PutObjectAclRequest struct {
	// The version id of the object.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s PutObjectAclRequest) String() string {
	return tea.Prettify(s)
}

func (s PutObjectAclRequest) GoString() string {
	return s.String()
}

func (s *PutObjectAclRequest) SetVersionId(v string) *PutObjectAclRequest {
	s.VersionId = &v
	return s
}

type PutObjectAclResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutObjectAclResponse) String() string {
	return tea.Prettify(s)
}

func (s PutObjectAclResponse) GoString() string {
	return s.String()
}

func (s *PutObjectAclResponse) SetHeaders(v map[string]*string) *PutObjectAclResponse {
	s.Headers = v
	return s
}

func (s *PutObjectAclResponse) SetStatusCode(v int32) *PutObjectAclResponse {
	s.StatusCode = &v
	return s
}

type PutObjectTaggingRequest struct {
	// The container of the tag set.
	Tagging *Tagging `json:"Tagging,omitempty" xml:"Tagging,omitempty"`
	// The version id of the target object.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s PutObjectTaggingRequest) String() string {
	return tea.Prettify(s)
}

func (s PutObjectTaggingRequest) GoString() string {
	return s.String()
}

func (s *PutObjectTaggingRequest) SetTagging(v *Tagging) *PutObjectTaggingRequest {
	s.Tagging = v
	return s
}

func (s *PutObjectTaggingRequest) SetVersionId(v string) *PutObjectTaggingRequest {
	s.VersionId = &v
	return s
}

type PutObjectTaggingResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutObjectTaggingResponse) String() string {
	return tea.Prettify(s)
}

func (s PutObjectTaggingResponse) GoString() string {
	return s.String()
}

func (s *PutObjectTaggingResponse) SetHeaders(v map[string]*string) *PutObjectTaggingResponse {
	s.Headers = v
	return s
}

func (s *PutObjectTaggingResponse) SetStatusCode(v int32) *PutObjectTaggingResponse {
	s.StatusCode = &v
	return s
}

type PutPublicAccessBlockRequest struct {
	PublicAccessBlockConfiguration *PublicAccessBlockConfiguration `json:"PublicAccessBlockConfiguration,omitempty" xml:"PublicAccessBlockConfiguration,omitempty"`
}

func (s PutPublicAccessBlockRequest) String() string {
	return tea.Prettify(s)
}

func (s PutPublicAccessBlockRequest) GoString() string {
	return s.String()
}

func (s *PutPublicAccessBlockRequest) SetPublicAccessBlockConfiguration(v *PublicAccessBlockConfiguration) *PutPublicAccessBlockRequest {
	s.PublicAccessBlockConfiguration = v
	return s
}

type PutPublicAccessBlockResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutPublicAccessBlockResponse) String() string {
	return tea.Prettify(s)
}

func (s PutPublicAccessBlockResponse) GoString() string {
	return s.String()
}

func (s *PutPublicAccessBlockResponse) SetHeaders(v map[string]*string) *PutPublicAccessBlockResponse {
	s.Headers = v
	return s
}

func (s *PutPublicAccessBlockResponse) SetStatusCode(v int32) *PutPublicAccessBlockResponse {
	s.StatusCode = &v
	return s
}

type PutStyleRequest struct {
	// The style content.
	Style    *Style  `json:"Style,omitempty" xml:"Style,omitempty"`
	Category *string `json:"category,omitempty" xml:"category,omitempty"`
	// The name of the image style.
	StyleName *string `json:"styleName,omitempty" xml:"styleName,omitempty"`
}

func (s PutStyleRequest) String() string {
	return tea.Prettify(s)
}

func (s PutStyleRequest) GoString() string {
	return s.String()
}

func (s *PutStyleRequest) SetStyle(v *Style) *PutStyleRequest {
	s.Style = v
	return s
}

func (s *PutStyleRequest) SetCategory(v string) *PutStyleRequest {
	s.Category = &v
	return s
}

func (s *PutStyleRequest) SetStyleName(v string) *PutStyleRequest {
	s.StyleName = &v
	return s
}

type PutStyleResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutStyleResponse) String() string {
	return tea.Prettify(s)
}

func (s PutStyleResponse) GoString() string {
	return s.String()
}

func (s *PutStyleResponse) SetHeaders(v map[string]*string) *PutStyleResponse {
	s.Headers = v
	return s
}

func (s *PutStyleResponse) SetStatusCode(v int32) *PutStyleResponse {
	s.StatusCode = &v
	return s
}

type PutSymlinkHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// Specifies whether the PutSymlink operation overwrites the object that has the same name as that of the symbolic link you want to create.
	//   - If the value of **x-oss-forbid-overwrite** is not specified or set to **false**, existing objects can be overwritten by objects that have the same names.
	//   - If the value of **x-oss-forbid-overwrite** is set to **true**, existing objects cannot be overwritten by objects that have the same names.
	//
	// If you specify the **x-oss-forbid-overwrite** request header, the queries per second (QPS) performance of OSS is degraded. If you want to use the **x-oss-forbid-overwrite** request header to perform a large number of operations (QPS greater than 1,000), contact technical support.
	// > The **x-oss-forbid-overwrite** request header is invalid when versioning is enabled or suspended for the destination bucket. In this case, the object with the same name can be overwritten.
	ForbidOverwrite *string `json:"x-oss-forbid-overwrite,omitempty" xml:"x-oss-forbid-overwrite,omitempty"`
	// The access control list (ACL) of the object. Default value: default.
	// Valid values:
	//
	// - default: The ACL of the object is the same as that of the bucket in which the object is stored.
	// - private: The ACL of the object is private. Only the owner of the object and authorized users can read and write this object.
	// - public-read: The ACL of the object is public-read. Only the owner of the object and authorized users can read and write this object. Other users can only read the object. Exercise caution when you set the object ACL to this value.
	// - public-read-write: The ACL of the object is public-read-write. All users can read and write this object. Exercise caution when you set the object ACL to this value.
	//
	// For more information about the ACL, see **[ACL](~~100676~~)**.
	Acl *string `json:"x-oss-object-acl,omitempty" xml:"x-oss-object-acl,omitempty"`
	// The storage class of the bucket. Default value: Standard.  Valid values:
	//
	// - Standard
	// - IA
	// - Archive
	// - ColdArchive
	StorageClass *string `json:"x-oss-storage-class,omitempty" xml:"x-oss-storage-class,omitempty"`
	// The target object to which the symbolic link points.
	// The naming conventions for target objects are the same as those for objects.
	//   - Similar to ObjectName, TargetObjectName must be URL-encoded.
	//   - The target object to which a symbolic link points cannot be a symbolic link.
	SymlinkTargetKey *string `json:"x-oss-symlink-target,omitempty" xml:"x-oss-symlink-target,omitempty"`
}

func (s PutSymlinkHeaders) String() string {
	return tea.Prettify(s)
}

func (s PutSymlinkHeaders) GoString() string {
	return s.String()
}

func (s *PutSymlinkHeaders) SetCommonHeaders(v map[string]*string) *PutSymlinkHeaders {
	s.CommonHeaders = v
	return s
}

func (s *PutSymlinkHeaders) SetForbidOverwrite(v string) *PutSymlinkHeaders {
	s.ForbidOverwrite = &v
	return s
}

func (s *PutSymlinkHeaders) SetAcl(v string) *PutSymlinkHeaders {
	s.Acl = &v
	return s
}

func (s *PutSymlinkHeaders) SetStorageClass(v string) *PutSymlinkHeaders {
	s.StorageClass = &v
	return s
}

func (s *PutSymlinkHeaders) SetSymlinkTargetKey(v string) *PutSymlinkHeaders {
	s.SymlinkTargetKey = &v
	return s
}

type PutSymlinkResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutSymlinkResponse) String() string {
	return tea.Prettify(s)
}

func (s PutSymlinkResponse) GoString() string {
	return s.String()
}

func (s *PutSymlinkResponse) SetHeaders(v map[string]*string) *PutSymlinkResponse {
	s.Headers = v
	return s
}

func (s *PutSymlinkResponse) SetStatusCode(v int32) *PutSymlinkResponse {
	s.StatusCode = &v
	return s
}

type PutUserDefinedLogFieldsConfigRequest struct {
	UserDefinedLogFieldsConfiguration *UserDefinedLogFieldsConfiguration `json:"UserDefinedLogFieldsConfiguration,omitempty" xml:"UserDefinedLogFieldsConfiguration,omitempty"`
}

func (s PutUserDefinedLogFieldsConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s PutUserDefinedLogFieldsConfigRequest) GoString() string {
	return s.String()
}

func (s *PutUserDefinedLogFieldsConfigRequest) SetUserDefinedLogFieldsConfiguration(v *UserDefinedLogFieldsConfiguration) *PutUserDefinedLogFieldsConfigRequest {
	s.UserDefinedLogFieldsConfiguration = v
	return s
}

type PutUserDefinedLogFieldsConfigResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s PutUserDefinedLogFieldsConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s PutUserDefinedLogFieldsConfigResponse) GoString() string {
	return s.String()
}

func (s *PutUserDefinedLogFieldsConfigResponse) SetHeaders(v map[string]*string) *PutUserDefinedLogFieldsConfigResponse {
	s.Headers = v
	return s
}

func (s *PutUserDefinedLogFieldsConfigResponse) SetStatusCode(v int32) *PutUserDefinedLogFieldsConfigResponse {
	s.StatusCode = &v
	return s
}

type RestoreObjectRequest struct {
	// The container that stores information about the RestoreObject request.
	RestoreRequest *RestoreRequest `json:"RestoreRequest,omitempty" xml:"RestoreRequest,omitempty"`
	// The version number of the object that you want to restore.
	VersionId *string `json:"versionId,omitempty" xml:"versionId,omitempty"`
}

func (s RestoreObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s RestoreObjectRequest) GoString() string {
	return s.String()
}

func (s *RestoreObjectRequest) SetRestoreRequest(v *RestoreRequest) *RestoreObjectRequest {
	s.RestoreRequest = v
	return s
}

func (s *RestoreObjectRequest) SetVersionId(v string) *RestoreObjectRequest {
	s.VersionId = &v
	return s
}

type RestoreObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s RestoreObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s RestoreObjectResponse) GoString() string {
	return s.String()
}

func (s *RestoreObjectResponse) SetHeaders(v map[string]*string) *RestoreObjectResponse {
	s.Headers = v
	return s
}

func (s *RestoreObjectResponse) SetStatusCode(v int32) *RestoreObjectResponse {
	s.StatusCode = &v
	return s
}

type SelectObjectRequest struct {
	// The container that stores the SelectObject request.
	SelectRequest *SelectRequest `json:"SelectRequest,omitempty" xml:"SelectRequest,omitempty"`
}

func (s SelectObjectRequest) String() string {
	return tea.Prettify(s)
}

func (s SelectObjectRequest) GoString() string {
	return s.String()
}

func (s *SelectObjectRequest) SetSelectRequest(v *SelectRequest) *SelectObjectRequest {
	s.SelectRequest = v
	return s
}

type SelectObjectResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       io.Reader          `json:"body,omitempty" xml:"body,omitempty"`
}

func (s SelectObjectResponse) String() string {
	return tea.Prettify(s)
}

func (s SelectObjectResponse) GoString() string {
	return s.String()
}

func (s *SelectObjectResponse) SetHeaders(v map[string]*string) *SelectObjectResponse {
	s.Headers = v
	return s
}

func (s *SelectObjectResponse) SetStatusCode(v int32) *SelectObjectResponse {
	s.StatusCode = &v
	return s
}

func (s *SelectObjectResponse) SetBody(v io.Reader) *SelectObjectResponse {
	s.Body = v
	return s
}

type UpdateBucketAntiDDosInfoHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The Anti-DDoS instance ID.
	DefenderInstance *string `json:"x-oss-defender-instance,omitempty" xml:"x-oss-defender-instance,omitempty"`
	// The new status of the Anti-DDoS instance. Valid values:
	//
	// *   Init: You must specify the custom domain name that you want to protect.
	// *   Defending: You can select whether to specify the custom domain name that you want to protect.
	// *   HaltDefending: You do not need to specify the custom domain name that you want to protect.
	DefenderStatus *string `json:"x-oss-defender-status,omitempty" xml:"x-oss-defender-status,omitempty"`
}

func (s UpdateBucketAntiDDosInfoHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateBucketAntiDDosInfoHeaders) GoString() string {
	return s.String()
}

func (s *UpdateBucketAntiDDosInfoHeaders) SetCommonHeaders(v map[string]*string) *UpdateBucketAntiDDosInfoHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateBucketAntiDDosInfoHeaders) SetDefenderInstance(v string) *UpdateBucketAntiDDosInfoHeaders {
	s.DefenderInstance = &v
	return s
}

func (s *UpdateBucketAntiDDosInfoHeaders) SetDefenderStatus(v string) *UpdateBucketAntiDDosInfoHeaders {
	s.DefenderStatus = &v
	return s
}

type UpdateBucketAntiDDosInfoRequest struct {
	// The container that stores the configurations of Anti-DDoS instances.
	AntiDDOSConfiguration *BucketAntiDDOSConfiguration `json:"AntiDDOSConfiguration,omitempty" xml:"AntiDDOSConfiguration,omitempty"`
}

func (s UpdateBucketAntiDDosInfoRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateBucketAntiDDosInfoRequest) GoString() string {
	return s.String()
}

func (s *UpdateBucketAntiDDosInfoRequest) SetAntiDDOSConfiguration(v *BucketAntiDDOSConfiguration) *UpdateBucketAntiDDosInfoRequest {
	s.AntiDDOSConfiguration = v
	return s
}

type UpdateBucketAntiDDosInfoResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s UpdateBucketAntiDDosInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateBucketAntiDDosInfoResponse) GoString() string {
	return s.String()
}

func (s *UpdateBucketAntiDDosInfoResponse) SetHeaders(v map[string]*string) *UpdateBucketAntiDDosInfoResponse {
	s.Headers = v
	return s
}

func (s *UpdateBucketAntiDDosInfoResponse) SetStatusCode(v int32) *UpdateBucketAntiDDosInfoResponse {
	s.StatusCode = &v
	return s
}

type UpdateUserAntiDDosInfoHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The Anti-DDoS instance ID.
	DefenderInstance *string `json:"x-oss-defender-instance,omitempty" xml:"x-oss-defender-instance,omitempty"`
	// The new status of the Anti-DDoS instance. Set the value to HaltDefending, which indicates that the Anti-DDos protection is disabled for a bucket.
	DefenderStatus *string `json:"x-oss-defender-status,omitempty" xml:"x-oss-defender-status,omitempty"`
}

func (s UpdateUserAntiDDosInfoHeaders) String() string {
	return tea.Prettify(s)
}

func (s UpdateUserAntiDDosInfoHeaders) GoString() string {
	return s.String()
}

func (s *UpdateUserAntiDDosInfoHeaders) SetCommonHeaders(v map[string]*string) *UpdateUserAntiDDosInfoHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UpdateUserAntiDDosInfoHeaders) SetDefenderInstance(v string) *UpdateUserAntiDDosInfoHeaders {
	s.DefenderInstance = &v
	return s
}

func (s *UpdateUserAntiDDosInfoHeaders) SetDefenderStatus(v string) *UpdateUserAntiDDosInfoHeaders {
	s.DefenderStatus = &v
	return s
}

type UpdateUserAntiDDosInfoResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s UpdateUserAntiDDosInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateUserAntiDDosInfoResponse) GoString() string {
	return s.String()
}

func (s *UpdateUserAntiDDosInfoResponse) SetHeaders(v map[string]*string) *UpdateUserAntiDDosInfoResponse {
	s.Headers = v
	return s
}

func (s *UpdateUserAntiDDosInfoResponse) SetStatusCode(v int32) *UpdateUserAntiDDosInfoResponse {
	s.StatusCode = &v
	return s
}

type UploadPartRequest struct {
	// The request body.
	Body io.Reader `json:"body,omitempty" xml:"body,omitempty"`
	// The number that identifies a part.
	//
	// Valid values: 1 to 10000.
	//
	// The size of a part ranges from 100 KB to 5 GB.
	// > In multipart upload, each part except the last part must be larger than or equal to 100 KB in size. When you call the UploadPart operation, the size of each part is not verified because not all parts have been uploaded and OSS does not know which part is the last part. The size of each part is verified only when you call CompleteMultipartUpload.
	PartNumber *int64 `json:"partNumber,omitempty" xml:"partNumber,omitempty"`
	// The ID that identifies the object to which the part that you want to upload belongs.
	UploadId *string `json:"uploadId,omitempty" xml:"uploadId,omitempty"`
}

func (s UploadPartRequest) String() string {
	return tea.Prettify(s)
}

func (s UploadPartRequest) GoString() string {
	return s.String()
}

func (s *UploadPartRequest) SetBody(v io.Reader) *UploadPartRequest {
	s.Body = v
	return s
}

func (s *UploadPartRequest) SetPartNumber(v int64) *UploadPartRequest {
	s.PartNumber = &v
	return s
}

func (s *UploadPartRequest) SetUploadId(v string) *UploadPartRequest {
	s.UploadId = &v
	return s
}

type UploadPartResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s UploadPartResponse) String() string {
	return tea.Prettify(s)
}

func (s UploadPartResponse) GoString() string {
	return s.String()
}

func (s *UploadPartResponse) SetHeaders(v map[string]*string) *UploadPartResponse {
	s.Headers = v
	return s
}

func (s *UploadPartResponse) SetStatusCode(v int32) *UploadPartResponse {
	s.StatusCode = &v
	return s
}

type UploadPartCopyHeaders struct {
	CommonHeaders map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	// The address to access the source object. You must have permissions to read the source object.
	// <br>Default value: null
	CopySource *string `json:"x-oss-copy-source,omitempty" xml:"x-oss-copy-source,omitempty"`
	// The copy operation condition. If the ETag value of the source object is the same as the ETag value provided by the user, OSS copies data. Otherwise, OSS returns 412 Precondition Failed.
	// <br>Default value: null
	CopySourceIfMatch *string `json:"x-oss-copy-source-if-match,omitempty" xml:"x-oss-copy-source-if-match,omitempty"`
	// The object transfer condition. If the specified time is earlier than the actual modified time of the object, the system transfers the object normally and returns 200 OK. Otherwise, OSS returns 304 Not Modified.
	// <br>Default value: null
	// <br>Time format: ddd, dd MMM yyyy HH:mm:ss GMT. Example: Fri, 13 Nov 2015 14:47:53 GMT.
	CopySourceIfModifiedSince *string `json:"x-oss-copy-source-if-modified-since,omitempty" xml:"x-oss-copy-source-if-modified-since,omitempty"`
	// The object transfer condition. If the input ETag value does not match the ETag value of the object, the system transfers the object normally and returns 200 OK. Otherwise, OSS returns 304 Not Modified.
	// <br>Default value: null
	CopySourceIfNoneMatch *string `json:"x-oss-copy-source-if-none-match,omitempty" xml:"x-oss-copy-source-if-none-match,omitempty"`
	// The object transfer condition. If the specified time is the same as or later than the actual modified time of the object, OSS transfers the object normally and returns 200 OK. Otherwise, OSS returns 412 Precondition Failed.
	// <br>Default value: null
	CopySourceIfUnmodifiedSince *string `json:"x-oss-copy-source-if-unmodified-since,omitempty" xml:"x-oss-copy-source-if-unmodified-since,omitempty"`
	// The range of bytes to copy data from the source object. For example, if you specify bytes to 0 to 9, the system transfers byte 0 to byte 9, a total of 10 bytes.
	// <br>Default value: null
	//
	// - If the x-oss-copy-source-range request header is not specified, the entire source object is copied.
	// - If the x-oss-copy-source-range request header is specified, the response contains the length of the entire object and the range of bytes to be copied for this operation. For example, Content-Range: bytes 0~9/44 indicates that the length of the entire object is 44 bytes. The range of bytes to be copied is byte 0 to byte 9.
	// - If the specified range does not conform to the range conventions, OSS copies the entire object and does not include Content-Range in the response.
	CopySourceRange *string `json:"x-oss-copy-source-range,omitempty" xml:"x-oss-copy-source-range,omitempty"`
}

func (s UploadPartCopyHeaders) String() string {
	return tea.Prettify(s)
}

func (s UploadPartCopyHeaders) GoString() string {
	return s.String()
}

func (s *UploadPartCopyHeaders) SetCommonHeaders(v map[string]*string) *UploadPartCopyHeaders {
	s.CommonHeaders = v
	return s
}

func (s *UploadPartCopyHeaders) SetCopySource(v string) *UploadPartCopyHeaders {
	s.CopySource = &v
	return s
}

func (s *UploadPartCopyHeaders) SetCopySourceIfMatch(v string) *UploadPartCopyHeaders {
	s.CopySourceIfMatch = &v
	return s
}

func (s *UploadPartCopyHeaders) SetCopySourceIfModifiedSince(v string) *UploadPartCopyHeaders {
	s.CopySourceIfModifiedSince = &v
	return s
}

func (s *UploadPartCopyHeaders) SetCopySourceIfNoneMatch(v string) *UploadPartCopyHeaders {
	s.CopySourceIfNoneMatch = &v
	return s
}

func (s *UploadPartCopyHeaders) SetCopySourceIfUnmodifiedSince(v string) *UploadPartCopyHeaders {
	s.CopySourceIfUnmodifiedSince = &v
	return s
}

func (s *UploadPartCopyHeaders) SetCopySourceRange(v string) *UploadPartCopyHeaders {
	s.CopySourceRange = &v
	return s
}

type UploadPartCopyRequest struct {
	// The number of parts.
	PartNumber *int64 `json:"partNumber,omitempty" xml:"partNumber,omitempty"`
	// The ID that identifies the object to which the parts to upload belong.
	UploadId *string `json:"uploadId,omitempty" xml:"uploadId,omitempty"`
}

func (s UploadPartCopyRequest) String() string {
	return tea.Prettify(s)
}

func (s UploadPartCopyRequest) GoString() string {
	return s.String()
}

func (s *UploadPartCopyRequest) SetPartNumber(v int64) *UploadPartCopyRequest {
	s.PartNumber = &v
	return s
}

func (s *UploadPartCopyRequest) SetUploadId(v string) *UploadPartCopyRequest {
	s.UploadId = &v
	return s
}

type UploadPartCopyResponseBody struct {
	// The container that stores the copy result.
	CopyPartResult *UploadPartCopyResponseBodyCopyPartResult `json:"CopyPartResult,omitempty" xml:"CopyPartResult,omitempty" type:"Struct"`
}

func (s UploadPartCopyResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UploadPartCopyResponseBody) GoString() string {
	return s.String()
}

func (s *UploadPartCopyResponseBody) SetCopyPartResult(v *UploadPartCopyResponseBodyCopyPartResult) *UploadPartCopyResponseBody {
	s.CopyPartResult = v
	return s
}

type UploadPartCopyResponseBodyCopyPartResult struct {
	// The ETag of the copied part.
	ETag *string `json:"ETag,omitempty" xml:"ETag,omitempty"`
	// The last modified time of copy source.
	LastModified *string `json:"LastModified,omitempty" xml:"LastModified,omitempty"`
}

func (s UploadPartCopyResponseBodyCopyPartResult) String() string {
	return tea.Prettify(s)
}

func (s UploadPartCopyResponseBodyCopyPartResult) GoString() string {
	return s.String()
}

func (s *UploadPartCopyResponseBodyCopyPartResult) SetETag(v string) *UploadPartCopyResponseBodyCopyPartResult {
	s.ETag = &v
	return s
}

func (s *UploadPartCopyResponseBodyCopyPartResult) SetLastModified(v string) *UploadPartCopyResponseBodyCopyPartResult {
	s.LastModified = &v
	return s
}

type UploadPartCopyResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *UploadPartCopyResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s UploadPartCopyResponse) String() string {
	return tea.Prettify(s)
}

func (s UploadPartCopyResponse) GoString() string {
	return s.String()
}

func (s *UploadPartCopyResponse) SetHeaders(v map[string]*string) *UploadPartCopyResponse {
	s.Headers = v
	return s
}

func (s *UploadPartCopyResponse) SetStatusCode(v int32) *UploadPartCopyResponse {
	s.StatusCode = &v
	return s
}

func (s *UploadPartCopyResponse) SetBody(v *UploadPartCopyResponseBody) *UploadPartCopyResponse {
	s.Body = v
	return s
}

type WriteGetObjectResponseHeaders struct {
	CommonHeaders                   map[string]*string `json:"commonHeaders,omitempty" xml:"commonHeaders,omitempty"`
	ContentLength                   *int64             `json:"Content-Length,omitempty" xml:"Content-Length,omitempty"`
	XOssFwdHeaderAcceptRanges       *string            `json:"x-oss-fwd-header-Accept-Ranges,omitempty" xml:"x-oss-fwd-header-Accept-Ranges,omitempty"`
	XOssFwdHeaderCacheControl       *string            `json:"x-oss-fwd-header-Cache-Control,omitempty" xml:"x-oss-fwd-header-Cache-Control,omitempty"`
	XOssFwdHeaderContentDisposition *string            `json:"x-oss-fwd-header-Content-Disposition,omitempty" xml:"x-oss-fwd-header-Content-Disposition,omitempty"`
	XOssFwdHeaderContentEncoding    *string            `json:"x-oss-fwd-header-Content-Encoding,omitempty" xml:"x-oss-fwd-header-Content-Encoding,omitempty"`
	XOssFwdHeaderContentLanguage    *string            `json:"x-oss-fwd-header-Content-Language,omitempty" xml:"x-oss-fwd-header-Content-Language,omitempty"`
	XOssFwdHeaderContentRange       *string            `json:"x-oss-fwd-header-Content-Range,omitempty" xml:"x-oss-fwd-header-Content-Range,omitempty"`
	XOssFwdHeaderContentType        *string            `json:"x-oss-fwd-header-Content-Type,omitempty" xml:"x-oss-fwd-header-Content-Type,omitempty"`
	XOssFwdHeaderETag               *string            `json:"x-oss-fwd-header-ETag,omitempty" xml:"x-oss-fwd-header-ETag,omitempty"`
	XOssFwdHeaderExpires            *string            `json:"x-oss-fwd-header-Expires,omitempty" xml:"x-oss-fwd-header-Expires,omitempty"`
	XOssFwdHeaderLastModified       *string            `json:"x-oss-fwd-header-Last-Modified,omitempty" xml:"x-oss-fwd-header-Last-Modified,omitempty"`
	XOssFwdStatus                   *string            `json:"x-oss-fwd-status,omitempty" xml:"x-oss-fwd-status,omitempty"`
	XOssRequestRoute                *string            `json:"x-oss-request-route,omitempty" xml:"x-oss-request-route,omitempty"`
	XOssRequestToken                *string            `json:"x-oss-request-token,omitempty" xml:"x-oss-request-token,omitempty"`
}

func (s WriteGetObjectResponseHeaders) String() string {
	return tea.Prettify(s)
}

func (s WriteGetObjectResponseHeaders) GoString() string {
	return s.String()
}

func (s *WriteGetObjectResponseHeaders) SetCommonHeaders(v map[string]*string) *WriteGetObjectResponseHeaders {
	s.CommonHeaders = v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetContentLength(v int64) *WriteGetObjectResponseHeaders {
	s.ContentLength = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderAcceptRanges(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderAcceptRanges = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderCacheControl(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderCacheControl = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderContentDisposition(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderContentDisposition = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderContentEncoding(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderContentEncoding = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderContentLanguage(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderContentLanguage = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderContentRange(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderContentRange = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderContentType(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderContentType = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderETag(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderETag = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderExpires(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderExpires = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdHeaderLastModified(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdHeaderLastModified = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssFwdStatus(v string) *WriteGetObjectResponseHeaders {
	s.XOssFwdStatus = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssRequestRoute(v string) *WriteGetObjectResponseHeaders {
	s.XOssRequestRoute = &v
	return s
}

func (s *WriteGetObjectResponseHeaders) SetXOssRequestToken(v string) *WriteGetObjectResponseHeaders {
	s.XOssRequestToken = &v
	return s
}

type WriteGetObjectResponseResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s WriteGetObjectResponseResponse) String() string {
	return tea.Prettify(s)
}

func (s WriteGetObjectResponseResponse) GoString() string {
	return s.String()
}

func (s *WriteGetObjectResponseResponse) SetHeaders(v map[string]*string) *WriteGetObjectResponseResponse {
	s.Headers = v
	return s
}

func (s *WriteGetObjectResponseResponse) SetStatusCode(v int32) *WriteGetObjectResponseResponse {
	s.StatusCode = &v
	return s
}

// for hcs-mgw
type AddressDetail struct {
	// This parameter is required.
	//
	// example:
	//
	// test_access_id
	AccessId *string `json:"AccessId,omitempty" xml:"AccessId,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_secret_key
	AccessSecret *string `json:"AccessSecret,omitempty" xml:"AccessSecret,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// ossinv
	AddressType *string `json:"AddressType,omitempty" xml:"AddressType,omitempty"`
	AgentList   *string `json:"AgentList,omitempty" xml:"AgentList,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_bucket
	Bucket *string `json:"Bucket,omitempty" xml:"Bucket,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_domain
	Domain *string `json:"Domain,omitempty" xml:"Domain,omitempty"`
	// example:
	//
	// test_inv_access_id
	InvAccessId *string `json:"InvAccessId,omitempty" xml:"InvAccessId,omitempty"`
	// example:
	//
	// test_inv_secret_key
	InvAccessSecret *string `json:"InvAccessSecret,omitempty" xml:"InvAccessSecret,omitempty"`
	// example:
	//
	// test_inv_bucket
	InvBucket *string `json:"InvBucket,omitempty" xml:"InvBucket,omitempty"`
	// example:
	//
	// test_inv_domain
	InvDomain *string `json:"InvDomain,omitempty" xml:"InvDomain,omitempty"`
	// example:
	//
	// oss
	InvLocation *string `json:"InvLocation,omitempty" xml:"InvLocation,omitempty"`
	// example:
	//
	// manifest.json
	InvPath *string `json:"InvPath,omitempty" xml:"InvPath,omitempty"`
	// example:
	//
	// test_inv_region_id
	InvRegionId *string `json:"InvRegionId,omitempty" xml:"InvRegionId,omitempty"`
	// example:
	//
	// test_inv_role
	InvRole *string `json:"InvRole,omitempty" xml:"InvRole,omitempty"`
	// example:
	//
	// test_prefix
	Prefix *string `json:"Prefix,omitempty" xml:"Prefix,omitempty"`
	// example:
	//
	// test_region_id
	RegionId *string `json:"RegionId,omitempty" xml:"RegionId,omitempty"`
	// example:
	//
	// test_role
	Role *string `json:"Role,omitempty" xml:"Role,omitempty"`
}

func (s AddressDetail) String() string {
	return tea.Prettify(s)
}

func (s AddressDetail) GoString() string {
	return s.String()
}

func (s *AddressDetail) SetAccessId(v string) *AddressDetail {
	s.AccessId = &v
	return s
}

func (s *AddressDetail) SetAccessSecret(v string) *AddressDetail {
	s.AccessSecret = &v
	return s
}

func (s *AddressDetail) SetAddressType(v string) *AddressDetail {
	s.AddressType = &v
	return s
}

func (s *AddressDetail) SetAgentList(v string) *AddressDetail {
	s.AgentList = &v
	return s
}

func (s *AddressDetail) SetBucket(v string) *AddressDetail {
	s.Bucket = &v
	return s
}

func (s *AddressDetail) SetDomain(v string) *AddressDetail {
	s.Domain = &v
	return s
}

func (s *AddressDetail) SetInvAccessId(v string) *AddressDetail {
	s.InvAccessId = &v
	return s
}

func (s *AddressDetail) SetInvAccessSecret(v string) *AddressDetail {
	s.InvAccessSecret = &v
	return s
}

func (s *AddressDetail) SetInvBucket(v string) *AddressDetail {
	s.InvBucket = &v
	return s
}

func (s *AddressDetail) SetInvDomain(v string) *AddressDetail {
	s.InvDomain = &v
	return s
}

func (s *AddressDetail) SetInvLocation(v string) *AddressDetail {
	s.InvLocation = &v
	return s
}

func (s *AddressDetail) SetInvPath(v string) *AddressDetail {
	s.InvPath = &v
	return s
}

func (s *AddressDetail) SetInvRegionId(v string) *AddressDetail {
	s.InvRegionId = &v
	return s
}

func (s *AddressDetail) SetInvRole(v string) *AddressDetail {
	s.InvRole = &v
	return s
}

func (s *AddressDetail) SetPrefix(v string) *AddressDetail {
	s.Prefix = &v
	return s
}

func (s *AddressDetail) SetRegionId(v string) *AddressDetail {
	s.RegionId = &v
	return s
}

func (s *AddressDetail) SetRole(v string) *AddressDetail {
	s.Role = &v
	return s
}

type Audit struct {
	// example:
	//
	// off
	LogMode *string `json:"LogMode,omitempty" xml:"LogMode,omitempty"`
}

func (s Audit) String() string {
	return tea.Prettify(s)
}

func (s Audit) GoString() string {
	return s.String()
}

func (s *Audit) SetLogMode(v string) *Audit {
	s.LogMode = &v
	return s
}

type CreateAddressInfo struct {
	// This parameter is required.
	AddressDetail *AddressDetail `json:"AddressDetail,omitempty" xml:"AddressDetail,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// example:
	//
	// K1:V1,K2:V2
	Tags *string `json:"Tags,omitempty" xml:"Tags,omitempty"`
}

func (s CreateAddressInfo) String() string {
	return tea.Prettify(s)
}

func (s CreateAddressInfo) GoString() string {
	return s.String()
}

func (s *CreateAddressInfo) SetAddressDetail(v *AddressDetail) *CreateAddressInfo {
	s.AddressDetail = v
	return s
}

func (s *CreateAddressInfo) SetName(v string) *CreateAddressInfo {
	s.Name = &v
	return s
}

func (s *CreateAddressInfo) SetTags(v string) *CreateAddressInfo {
	s.Tags = &v
	return s
}

type CreateAgentInfo struct {
	// This parameter is required.
	//
	// example:
	//
	// vpc
	AgentEndpoint *string `json:"AgentEndpoint,omitempty" xml:"AgentEndpoint,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// default
	DeployMethod *string `json:"DeployMethod,omitempty" xml:"DeployMethod,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// example:
	//
	// K1:V1,K2:V2
	Tags *string `json:"Tags,omitempty" xml:"Tags,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_tunnel_id
	TunnelId *string `json:"TunnelId,omitempty" xml:"TunnelId,omitempty"`
}

func (s CreateAgentInfo) String() string {
	return tea.Prettify(s)
}

func (s CreateAgentInfo) GoString() string {
	return s.String()
}

func (s *CreateAgentInfo) SetAgentEndpoint(v string) *CreateAgentInfo {
	s.AgentEndpoint = &v
	return s
}

func (s *CreateAgentInfo) SetDeployMethod(v string) *CreateAgentInfo {
	s.DeployMethod = &v
	return s
}

func (s *CreateAgentInfo) SetName(v string) *CreateAgentInfo {
	s.Name = &v
	return s
}

func (s *CreateAgentInfo) SetTags(v string) *CreateAgentInfo {
	s.Tags = &v
	return s
}

func (s *CreateAgentInfo) SetTunnelId(v string) *CreateAgentInfo {
	s.TunnelId = &v
	return s
}

type CreateJobInfo struct {
	Audit *Audit `json:"Audit,omitempty" xml:"Audit,omitempty"`
	// example:
	//
	// false
	ConvertSymlinkTarget *bool `json:"ConvertSymlinkTarget,omitempty" xml:"ConvertSymlinkTarget,omitempty"`
	CreateReport         *bool `json:"CreateReport,omitempty" xml:"CreateReport,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_dest_address
	DestAddress           *string     `json:"DestAddress,omitempty" xml:"DestAddress,omitempty"`
	EnableMultiVersioning *bool       `json:"EnableMultiVersioning,omitempty" xml:"EnableMultiVersioning,omitempty"`
	FilterRule            *FilterRule `json:"FilterRule,omitempty" xml:"FilterRule,omitempty"`
	ImportQos             *ImportQos  `json:"ImportQos,omitempty" xml:"ImportQos,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// always
	OverwriteMode *string       `json:"OverwriteMode,omitempty" xml:"OverwriteMode,omitempty"`
	ParentVersion *string       `json:"ParentVersion,omitempty" xml:"ParentVersion,omitempty"`
	ScheduleRule  *ScheduleRule `json:"ScheduleRule,omitempty" xml:"ScheduleRule,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_src_address
	SrcAddress *string `json:"SrcAddress,omitempty" xml:"SrcAddress,omitempty"`
	// example:
	//
	// K1:V1,K2:V2
	Tags *string `json:"Tags,omitempty" xml:"Tags,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// all
	TransferMode *string `json:"TransferMode,omitempty" xml:"TransferMode,omitempty"`
}

func (s CreateJobInfo) String() string {
	return tea.Prettify(s)
}

func (s CreateJobInfo) GoString() string {
	return s.String()
}

func (s *CreateJobInfo) SetAudit(v *Audit) *CreateJobInfo {
	s.Audit = v
	return s
}

func (s *CreateJobInfo) SetConvertSymlinkTarget(v bool) *CreateJobInfo {
	s.ConvertSymlinkTarget = &v
	return s
}

func (s *CreateJobInfo) SetCreateReport(v bool) *CreateJobInfo {
	s.CreateReport = &v
	return s
}

func (s *CreateJobInfo) SetDestAddress(v string) *CreateJobInfo {
	s.DestAddress = &v
	return s
}

func (s *CreateJobInfo) SetEnableMultiVersioning(v bool) *CreateJobInfo {
	s.EnableMultiVersioning = &v
	return s
}

func (s *CreateJobInfo) SetFilterRule(v *FilterRule) *CreateJobInfo {
	s.FilterRule = v
	return s
}

func (s *CreateJobInfo) SetImportQos(v *ImportQos) *CreateJobInfo {
	s.ImportQos = v
	return s
}

func (s *CreateJobInfo) SetName(v string) *CreateJobInfo {
	s.Name = &v
	return s
}

func (s *CreateJobInfo) SetOverwriteMode(v string) *CreateJobInfo {
	s.OverwriteMode = &v
	return s
}

func (s *CreateJobInfo) SetParentVersion(v string) *CreateJobInfo {
	s.ParentVersion = &v
	return s
}

func (s *CreateJobInfo) SetScheduleRule(v *ScheduleRule) *CreateJobInfo {
	s.ScheduleRule = v
	return s
}

func (s *CreateJobInfo) SetSrcAddress(v string) *CreateJobInfo {
	s.SrcAddress = &v
	return s
}

func (s *CreateJobInfo) SetTags(v string) *CreateJobInfo {
	s.Tags = &v
	return s
}

func (s *CreateJobInfo) SetTransferMode(v string) *CreateJobInfo {
	s.TransferMode = &v
	return s
}

type CreateReportInfo struct {
	// example:
	//
	// test_job_name
	JobName *string `json:"JobName,omitempty" xml:"JobName,omitempty"`
	// example:
	//
	// 1
	RuntimeId *int32 `json:"RuntimeId,omitempty" xml:"RuntimeId,omitempty"`
	// example:
	//
	// test_job_id
	Version *string `json:"Version,omitempty" xml:"Version,omitempty"`
}

func (s CreateReportInfo) String() string {
	return tea.Prettify(s)
}

func (s CreateReportInfo) GoString() string {
	return s.String()
}

func (s *CreateReportInfo) SetJobName(v string) *CreateReportInfo {
	s.JobName = &v
	return s
}

func (s *CreateReportInfo) SetRuntimeId(v int32) *CreateReportInfo {
	s.RuntimeId = &v
	return s
}

func (s *CreateReportInfo) SetVersion(v string) *CreateReportInfo {
	s.Version = &v
	return s
}

type CreateTunnelInfo struct {
	// example:
	//
	// K1:V1,K2:V2
	Tags      *string    `json:"Tags,omitempty" xml:"Tags,omitempty"`
	TunnelQos *TunnelQos `json:"TunnelQos,omitempty" xml:"TunnelQos,omitempty"`
}

func (s CreateTunnelInfo) String() string {
	return tea.Prettify(s)
}

func (s CreateTunnelInfo) GoString() string {
	return s.String()
}

func (s *CreateTunnelInfo) SetTags(v string) *CreateTunnelInfo {
	s.Tags = &v
	return s
}

func (s *CreateTunnelInfo) SetTunnelQos(v *TunnelQos) *CreateTunnelInfo {
	s.TunnelQos = v
	return s
}

type FileTypeFilters struct {
	// example:
	//
	// fasle
	ExcludeDir *bool `json:"ExcludeDir,omitempty" xml:"ExcludeDir,omitempty"`
	// example:
	//
	// fasle
	ExcludeSymlink *bool `json:"ExcludeSymlink,omitempty" xml:"ExcludeSymlink,omitempty"`
}

func (s FileTypeFilters) String() string {
	return tea.Prettify(s)
}

func (s FileTypeFilters) GoString() string {
	return s.String()
}

func (s *FileTypeFilters) SetExcludeDir(v bool) *FileTypeFilters {
	s.ExcludeDir = &v
	return s
}

func (s *FileTypeFilters) SetExcludeSymlink(v bool) *FileTypeFilters {
	s.ExcludeSymlink = &v
	return s
}

type FilterRule struct {
	FileTypeFilters     *FileTypeFilters     `json:"FileTypeFilters,omitempty" xml:"FileTypeFilters,omitempty"`
	KeyFilters          *KeyFilters          `json:"KeyFilters,omitempty" xml:"KeyFilters,omitempty"`
	LastModifiedFilters *LastModifiedFilters `json:"LastModifiedFilters,omitempty" xml:"LastModifiedFilters,omitempty"`
}

func (s FilterRule) String() string {
	return tea.Prettify(s)
}

func (s FilterRule) GoString() string {
	return s.String()
}

func (s *FilterRule) SetFileTypeFilters(v *FileTypeFilters) *FilterRule {
	s.FileTypeFilters = v
	return s
}

func (s *FilterRule) SetKeyFilters(v *KeyFilters) *FilterRule {
	s.KeyFilters = v
	return s
}

func (s *FilterRule) SetLastModifiedFilters(v *LastModifiedFilters) *FilterRule {
	s.LastModifiedFilters = v
	return s
}

type GetAddressResp struct {
	AddressDetail *AddressDetail `json:"AddressDetail,omitempty" xml:"AddressDetail,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	CreateTime *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	ModifyTime *string `json:"ModifyTime,omitempty" xml:"ModifyTime,omitempty"`
	// example:
	//
	// test_name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// example:
	//
	// test_owner
	Owner *string `json:"Owner,omitempty" xml:"Owner,omitempty"`
	// example:
	//
	// avaliable
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
	// example:
	//
	// K1:V1,K2:V2
	Tags         *string     `json:"Tags,omitempty" xml:"Tags,omitempty"`
	VerifyResult *VerifyResp `json:"VerifyResult,omitempty" xml:"VerifyResult,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	VerifyTime *string `json:"VerifyTime,omitempty" xml:"VerifyTime,omitempty"`
	// example:
	//
	// test_id
	Version *string `json:"Version,omitempty" xml:"Version,omitempty"`
}

func (s GetAddressResp) String() string {
	return tea.Prettify(s)
}

func (s GetAddressResp) GoString() string {
	return s.String()
}

func (s *GetAddressResp) SetAddressDetail(v *AddressDetail) *GetAddressResp {
	s.AddressDetail = v
	return s
}

func (s *GetAddressResp) SetCreateTime(v string) *GetAddressResp {
	s.CreateTime = &v
	return s
}

func (s *GetAddressResp) SetModifyTime(v string) *GetAddressResp {
	s.ModifyTime = &v
	return s
}

func (s *GetAddressResp) SetName(v string) *GetAddressResp {
	s.Name = &v
	return s
}

func (s *GetAddressResp) SetOwner(v string) *GetAddressResp {
	s.Owner = &v
	return s
}

func (s *GetAddressResp) SetStatus(v string) *GetAddressResp {
	s.Status = &v
	return s
}

func (s *GetAddressResp) SetTags(v string) *GetAddressResp {
	s.Tags = &v
	return s
}

func (s *GetAddressResp) SetVerifyResult(v *VerifyResp) *GetAddressResp {
	s.VerifyResult = v
	return s
}

func (s *GetAddressResp) SetVerifyTime(v string) *GetAddressResp {
	s.VerifyTime = &v
	return s
}

func (s *GetAddressResp) SetVersion(v string) *GetAddressResp {
	s.Version = &v
	return s
}

type GetAgentResp struct {
	ActivationKey *string `json:"ActivationKey,omitempty" xml:"ActivationKey,omitempty"`
	// example:
	//
	// vpc
	AgentEndpoint *string `json:"AgentEndpoint,omitempty" xml:"AgentEndpoint,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	CreateTime *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	// example:
	//
	// default
	DeployMethod *string `json:"DeployMethod,omitempty" xml:"DeployMethod,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	ModifyTime *string `json:"ModifyTime,omitempty" xml:"ModifyTime,omitempty"`
	// example:
	//
	// test_name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// example:
	//
	// test_owner
	Owner *string `json:"Owner,omitempty" xml:"Owner,omitempty"`
	// example:
	//
	// K1:V1,K2:V2
	Tags *string `json:"Tags,omitempty" xml:"Tags,omitempty"`
	// example:
	//
	// test_tunnel_id
	TunnelId *string `json:"TunnelId,omitempty" xml:"TunnelId,omitempty"`
	// example:
	//
	// test_agent_id
	Version *string `json:"Version,omitempty" xml:"Version,omitempty"`
}

func (s GetAgentResp) String() string {
	return tea.Prettify(s)
}

func (s GetAgentResp) GoString() string {
	return s.String()
}

func (s *GetAgentResp) SetActivationKey(v string) *GetAgentResp {
	s.ActivationKey = &v
	return s
}

func (s *GetAgentResp) SetAgentEndpoint(v string) *GetAgentResp {
	s.AgentEndpoint = &v
	return s
}

func (s *GetAgentResp) SetCreateTime(v string) *GetAgentResp {
	s.CreateTime = &v
	return s
}

func (s *GetAgentResp) SetDeployMethod(v string) *GetAgentResp {
	s.DeployMethod = &v
	return s
}

func (s *GetAgentResp) SetModifyTime(v string) *GetAgentResp {
	s.ModifyTime = &v
	return s
}

func (s *GetAgentResp) SetName(v string) *GetAgentResp {
	s.Name = &v
	return s
}

func (s *GetAgentResp) SetOwner(v string) *GetAgentResp {
	s.Owner = &v
	return s
}

func (s *GetAgentResp) SetTags(v string) *GetAgentResp {
	s.Tags = &v
	return s
}

func (s *GetAgentResp) SetTunnelId(v string) *GetAgentResp {
	s.TunnelId = &v
	return s
}

func (s *GetAgentResp) SetVersion(v string) *GetAgentResp {
	s.Version = &v
	return s
}

type GetAgentStatusResp struct {
	// example:
	//
	// OK
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s GetAgentStatusResp) String() string {
	return tea.Prettify(s)
}

func (s GetAgentStatusResp) GoString() string {
	return s.String()
}

func (s *GetAgentStatusResp) SetStatus(v string) *GetAgentStatusResp {
	s.Status = &v
	return s
}

type GetJobResp struct {
	Audit *Audit `json:"Audit,omitempty" xml:"Audit,omitempty"`
	// example:
	//
	// false
	ConvertSymlinkTarget *bool `json:"ConvertSymlinkTarget,omitempty" xml:"ConvertSymlinkTarget,omitempty"`
	// example:
	//
	// false
	CreateReport *bool `json:"CreateReport,omitempty" xml:"CreateReport,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	CreateTime *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	// example:
	//
	// test_dest_address
	DestAddress *string `json:"DestAddress,omitempty" xml:"DestAddress,omitempty"`
	// example:
	//
	// false
	EnableMultiVersioning *bool       `json:"EnableMultiVersioning,omitempty" xml:"EnableMultiVersioning,omitempty"`
	FilterRule            *FilterRule `json:"FilterRule,omitempty" xml:"FilterRule,omitempty"`
	ImportQos             *ImportQos  `json:"ImportQos,omitempty" xml:"ImportQos,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	ModifyTime *string `json:"ModifyTime,omitempty" xml:"ModifyTime,omitempty"`
	// example:
	//
	// test_name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// example:
	//
	// always
	OverwriteMode *string       `json:"OverwriteMode,omitempty" xml:"OverwriteMode,omitempty"`
	Owner         *string       `json:"Owner,omitempty" xml:"Owner,omitempty"`
	ParentName    *string       `json:"ParentName,omitempty" xml:"ParentName,omitempty"`
	ParentVersion *string       `json:"ParentVersion,omitempty" xml:"ParentVersion,omitempty"`
	ScheduleRule  *ScheduleRule `json:"ScheduleRule,omitempty" xml:"ScheduleRule,omitempty"`
	// example:
	//
	// test_src_address
	SrcAddress *string `json:"SrcAddress,omitempty" xml:"SrcAddress,omitempty"`
	// example:
	//
	// IMPORT_JOB_DOING
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
	// example:
	//
	// K1:V1,K2:V2
	Tags *string `json:"Tags,omitempty" xml:"Tags,omitempty"`
	// example:
	//
	// all
	TransferMode *string `json:"TransferMode,omitempty" xml:"TransferMode,omitempty"`
	// example:
	//
	// test_id
	Version *string `json:"Version,omitempty" xml:"Version,omitempty"`
}

func (s GetJobResp) String() string {
	return tea.Prettify(s)
}

func (s GetJobResp) GoString() string {
	return s.String()
}

func (s *GetJobResp) SetAudit(v *Audit) *GetJobResp {
	s.Audit = v
	return s
}

func (s *GetJobResp) SetConvertSymlinkTarget(v bool) *GetJobResp {
	s.ConvertSymlinkTarget = &v
	return s
}

func (s *GetJobResp) SetCreateReport(v bool) *GetJobResp {
	s.CreateReport = &v
	return s
}

func (s *GetJobResp) SetCreateTime(v string) *GetJobResp {
	s.CreateTime = &v
	return s
}

func (s *GetJobResp) SetDestAddress(v string) *GetJobResp {
	s.DestAddress = &v
	return s
}

func (s *GetJobResp) SetEnableMultiVersioning(v bool) *GetJobResp {
	s.EnableMultiVersioning = &v
	return s
}

func (s *GetJobResp) SetFilterRule(v *FilterRule) *GetJobResp {
	s.FilterRule = v
	return s
}

func (s *GetJobResp) SetImportQos(v *ImportQos) *GetJobResp {
	s.ImportQos = v
	return s
}

func (s *GetJobResp) SetModifyTime(v string) *GetJobResp {
	s.ModifyTime = &v
	return s
}

func (s *GetJobResp) SetName(v string) *GetJobResp {
	s.Name = &v
	return s
}

func (s *GetJobResp) SetOverwriteMode(v string) *GetJobResp {
	s.OverwriteMode = &v
	return s
}

func (s *GetJobResp) SetOwner(v string) *GetJobResp {
	s.Owner = &v
	return s
}

func (s *GetJobResp) SetParentName(v string) *GetJobResp {
	s.ParentName = &v
	return s
}

func (s *GetJobResp) SetParentVersion(v string) *GetJobResp {
	s.ParentVersion = &v
	return s
}

func (s *GetJobResp) SetScheduleRule(v *ScheduleRule) *GetJobResp {
	s.ScheduleRule = v
	return s
}

func (s *GetJobResp) SetSrcAddress(v string) *GetJobResp {
	s.SrcAddress = &v
	return s
}

func (s *GetJobResp) SetStatus(v string) *GetJobResp {
	s.Status = &v
	return s
}

func (s *GetJobResp) SetTags(v string) *GetJobResp {
	s.Tags = &v
	return s
}

func (s *GetJobResp) SetTransferMode(v string) *GetJobResp {
	s.TransferMode = &v
	return s
}

func (s *GetJobResp) SetVersion(v string) *GetJobResp {
	s.Version = &v
	return s
}

type GetJobResultResp struct {
	// example:
	//
	// ossinv
	AddressType *string `json:"AddressType,omitempty" xml:"AddressType,omitempty"`
	// example:
	//
	// 800
	CopiedObjectCount *int64 `json:"CopiedObjectCount,omitempty" xml:"CopiedObjectCount,omitempty"`
	// example:
	//
	// 800
	CopiedObjectSize *int64 `json:"CopiedObjectSize,omitempty" xml:"CopiedObjectSize,omitempty"`
	// example:
	//
	// 200
	FailedObjectCount *int64 `json:"FailedObjectCount,omitempty" xml:"FailedObjectCount,omitempty"`
	// example:
	//
	// test_access_id
	InvAccessId *string `json:"InvAccessId,omitempty" xml:"InvAccessId,omitempty"`
	// example:
	//
	// test_secret_key
	InvAccessSecret *string `json:"InvAccessSecret,omitempty" xml:"InvAccessSecret,omitempty"`
	// example:
	//
	// test_sys_bucket
	InvBucket *string `json:"InvBucket,omitempty" xml:"InvBucket,omitempty"`
	// example:
	//
	// test_domain
	InvDomain *string `json:"InvDomain,omitempty" xml:"InvDomain,omitempty"`
	// example:
	//
	// oss
	InvLocation *string `json:"InvLocation,omitempty" xml:"InvLocation,omitempty"`
	// example:
	//
	// mainfest.json
	InvPath *string `json:"InvPath,omitempty" xml:"InvPath,omitempty"`
	// example:
	//
	// test_region_id
	InvRegionId *string `json:"InvRegionId,omitempty" xml:"InvRegionId,omitempty"`
	// example:
	//
	// Ready
	ReadyRetry *string `json:"ReadyRetry,omitempty" xml:"ReadyRetry,omitempty"`
	// example:
	//
	// 1000
	TotalObjectCount *int64 `json:"TotalObjectCount,omitempty" xml:"TotalObjectCount,omitempty"`
	// example:
	//
	// 1000
	TotalObjectSize *int64 `json:"TotalObjectSize,omitempty" xml:"TotalObjectSize,omitempty"`
	// example:
	//
	// test_job_id
	Version *string `json:"Version,omitempty" xml:"Version,omitempty"`
}

func (s GetJobResultResp) String() string {
	return tea.Prettify(s)
}

func (s GetJobResultResp) GoString() string {
	return s.String()
}

func (s *GetJobResultResp) SetAddressType(v string) *GetJobResultResp {
	s.AddressType = &v
	return s
}

func (s *GetJobResultResp) SetCopiedObjectCount(v int64) *GetJobResultResp {
	s.CopiedObjectCount = &v
	return s
}

func (s *GetJobResultResp) SetCopiedObjectSize(v int64) *GetJobResultResp {
	s.CopiedObjectSize = &v
	return s
}

func (s *GetJobResultResp) SetFailedObjectCount(v int64) *GetJobResultResp {
	s.FailedObjectCount = &v
	return s
}

func (s *GetJobResultResp) SetInvAccessId(v string) *GetJobResultResp {
	s.InvAccessId = &v
	return s
}

func (s *GetJobResultResp) SetInvAccessSecret(v string) *GetJobResultResp {
	s.InvAccessSecret = &v
	return s
}

func (s *GetJobResultResp) SetInvBucket(v string) *GetJobResultResp {
	s.InvBucket = &v
	return s
}

func (s *GetJobResultResp) SetInvDomain(v string) *GetJobResultResp {
	s.InvDomain = &v
	return s
}

func (s *GetJobResultResp) SetInvLocation(v string) *GetJobResultResp {
	s.InvLocation = &v
	return s
}

func (s *GetJobResultResp) SetInvPath(v string) *GetJobResultResp {
	s.InvPath = &v
	return s
}

func (s *GetJobResultResp) SetInvRegionId(v string) *GetJobResultResp {
	s.InvRegionId = &v
	return s
}

func (s *GetJobResultResp) SetReadyRetry(v string) *GetJobResultResp {
	s.ReadyRetry = &v
	return s
}

func (s *GetJobResultResp) SetTotalObjectCount(v int64) *GetJobResultResp {
	s.TotalObjectCount = &v
	return s
}

func (s *GetJobResultResp) SetTotalObjectSize(v int64) *GetJobResultResp {
	s.TotalObjectSize = &v
	return s
}

func (s *GetJobResultResp) SetVersion(v string) *GetJobResultResp {
	s.Version = &v
	return s
}

type GetReportResp struct {
	// example:
	//
	// 800
	CopiedCount  *int64  `json:"CopiedCount,omitempty" xml:"CopiedCount,omitempty"`
	ErrorMessage *string `json:"ErrorMessage,omitempty" xml:"ErrorMessage,omitempty"`
	// example:
	//
	// 100
	FailedCount *int64 `json:"FailedCount,omitempty" xml:"FailedCount,omitempty"`
	// example:
	//
	// test_failed_prefix/
	FailedListPrefix *string `json:"FailedListPrefix,omitempty" xml:"FailedListPrefix,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	JobCreateTime *string `json:"JobCreateTime,omitempty" xml:"JobCreateTime,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	JobEndTime *string `json:"JobEndTime,omitempty" xml:"JobEndTime,omitempty"`
	// example:
	//
	// 1000
	JobExecuteTime *string `json:"JobExecuteTime,omitempty" xml:"JobExecuteTime,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	ReportCreateTime *string `json:"ReportCreateTime,omitempty" xml:"ReportCreateTime,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	ReportEndTime *string `json:"ReportEndTime,omitempty" xml:"ReportEndTime,omitempty"`
	// example:
	//
	// 100
	SkippedCount *int64 `json:"SkippedCount,omitempty" xml:"SkippedCount,omitempty"`
	// example:
	//
	// test_skipped_prefix/
	SkippedListPrefix *string `json:"SkippedListPrefix,omitempty" xml:"SkippedListPrefix,omitempty"`
	// example:
	//
	// Running
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
	// example:
	//
	// 1000
	TotalCount *int64 `json:"TotalCount,omitempty" xml:"TotalCount,omitempty"`
	// example:
	//
	// test_total_prefix/
	TotalListPrefix *string `json:"TotalListPrefix,omitempty" xml:"TotalListPrefix,omitempty"`
}

func (s GetReportResp) String() string {
	return tea.Prettify(s)
}

func (s GetReportResp) GoString() string {
	return s.String()
}

func (s *GetReportResp) SetCopiedCount(v int64) *GetReportResp {
	s.CopiedCount = &v
	return s
}

func (s *GetReportResp) SetErrorMessage(v string) *GetReportResp {
	s.ErrorMessage = &v
	return s
}

func (s *GetReportResp) SetFailedCount(v int64) *GetReportResp {
	s.FailedCount = &v
	return s
}

func (s *GetReportResp) SetFailedListPrefix(v string) *GetReportResp {
	s.FailedListPrefix = &v
	return s
}

func (s *GetReportResp) SetJobCreateTime(v string) *GetReportResp {
	s.JobCreateTime = &v
	return s
}

func (s *GetReportResp) SetJobEndTime(v string) *GetReportResp {
	s.JobEndTime = &v
	return s
}

func (s *GetReportResp) SetJobExecuteTime(v string) *GetReportResp {
	s.JobExecuteTime = &v
	return s
}

func (s *GetReportResp) SetReportCreateTime(v string) *GetReportResp {
	s.ReportCreateTime = &v
	return s
}

func (s *GetReportResp) SetReportEndTime(v string) *GetReportResp {
	s.ReportEndTime = &v
	return s
}

func (s *GetReportResp) SetSkippedCount(v int64) *GetReportResp {
	s.SkippedCount = &v
	return s
}

func (s *GetReportResp) SetSkippedListPrefix(v string) *GetReportResp {
	s.SkippedListPrefix = &v
	return s
}

func (s *GetReportResp) SetStatus(v string) *GetReportResp {
	s.Status = &v
	return s
}

func (s *GetReportResp) SetTotalCount(v int64) *GetReportResp {
	s.TotalCount = &v
	return s
}

func (s *GetReportResp) SetTotalListPrefix(v string) *GetReportResp {
	s.TotalListPrefix = &v
	return s
}

type GetTunnelResp struct {
	// example:
	//
	// 2024-05-01 12:00:00
	CreateTime *string `json:"CreateTime,omitempty" xml:"CreateTime,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	ModifyTime *string `json:"ModifyTime,omitempty" xml:"ModifyTime,omitempty"`
	// example:
	//
	// test_owner
	Owner *string `json:"Owner,omitempty" xml:"Owner,omitempty"`
	// example:
	//
	// K1:V1,K2:V2
	Tags *string `json:"Tags,omitempty" xml:"Tags,omitempty"`
	// example:
	//
	// test_tunnel_id
	TunnelId  *string    `json:"TunnelId,omitempty" xml:"TunnelId,omitempty"`
	TunnelQos *TunnelQos `json:"TunnelQos,omitempty" xml:"TunnelQos,omitempty"`
}

func (s GetTunnelResp) String() string {
	return tea.Prettify(s)
}

func (s GetTunnelResp) GoString() string {
	return s.String()
}

func (s *GetTunnelResp) SetCreateTime(v string) *GetTunnelResp {
	s.CreateTime = &v
	return s
}

func (s *GetTunnelResp) SetModifyTime(v string) *GetTunnelResp {
	s.ModifyTime = &v
	return s
}

func (s *GetTunnelResp) SetOwner(v string) *GetTunnelResp {
	s.Owner = &v
	return s
}

func (s *GetTunnelResp) SetTags(v string) *GetTunnelResp {
	s.Tags = &v
	return s
}

func (s *GetTunnelResp) SetTunnelId(v string) *GetTunnelResp {
	s.TunnelId = &v
	return s
}

func (s *GetTunnelResp) SetTunnelQos(v *TunnelQos) *GetTunnelResp {
	s.TunnelQos = v
	return s
}

type ImportQos struct {
	// example:
	//
	// 1073741824
	MaxBandWidth *int64 `json:"MaxBandWidth,omitempty" xml:"MaxBandWidth,omitempty"`
	// example:
	//
	// 1000
	MaxImportTaskQps *int64 `json:"MaxImportTaskQps,omitempty" xml:"MaxImportTaskQps,omitempty"`
}

func (s ImportQos) String() string {
	return tea.Prettify(s)
}

func (s ImportQos) GoString() string {
	return s.String()
}

func (s *ImportQos) SetMaxBandWidth(v int64) *ImportQos {
	s.MaxBandWidth = &v
	return s
}

func (s *ImportQos) SetMaxImportTaskQps(v int64) *ImportQos {
	s.MaxImportTaskQps = &v
	return s
}

type JobHistory struct {
	// example:
	//
	// 2
	CommitId *string `json:"CommitId,omitempty" xml:"CommitId,omitempty"`
	// example:
	//
	// 900
	CopiedCount *int64 `json:"CopiedCount,omitempty" xml:"CopiedCount,omitempty"`
	// example:
	//
	// 1000
	CopiedSize *int64 `json:"CopiedSize,omitempty" xml:"CopiedSize,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	EndTime *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	// example:
	//
	// 100
	FailedCount *int64 `json:"FailedCount,omitempty" xml:"FailedCount,omitempty"`
	// example:
	//
	// test_id
	JobVersion *string `json:"JobVersion,omitempty" xml:"JobVersion,omitempty"`
	// example:
	//
	// Listing
	ListStatus *string `json:"ListStatus,omitempty" xml:"ListStatus,omitempty"`
	Message    *string `json:"Message,omitempty" xml:"Message,omitempty"`
	// example:
	//
	// test_name
	Name *string `json:"Name,omitempty" xml:"Name,omitempty"`
	// example:
	//
	// user
	Operator *string `json:"Operator,omitempty" xml:"Operator,omitempty"`
	// example:
	//
	// 1
	RuntimeId *string `json:"RuntimeId,omitempty" xml:"RuntimeId,omitempty"`
	// example:
	//
	// Normal
	RuntimeState *string `json:"RuntimeState,omitempty" xml:"RuntimeState,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	StartTime *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
	// example:
	//
	// IMPORT_JOB_DOING
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
	// example:
	//
	// 1000
	TotalCount *int64 `json:"TotalCount,omitempty" xml:"TotalCount,omitempty"`
	// example:
	//
	// 1000
	TotalSize *int64 `json:"TotalSize,omitempty" xml:"TotalSize,omitempty"`
}

func (s JobHistory) String() string {
	return tea.Prettify(s)
}

func (s JobHistory) GoString() string {
	return s.String()
}

func (s *JobHistory) SetCommitId(v string) *JobHistory {
	s.CommitId = &v
	return s
}

func (s *JobHistory) SetCopiedCount(v int64) *JobHistory {
	s.CopiedCount = &v
	return s
}

func (s *JobHistory) SetCopiedSize(v int64) *JobHistory {
	s.CopiedSize = &v
	return s
}

func (s *JobHistory) SetEndTime(v string) *JobHistory {
	s.EndTime = &v
	return s
}

func (s *JobHistory) SetFailedCount(v int64) *JobHistory {
	s.FailedCount = &v
	return s
}

func (s *JobHistory) SetJobVersion(v string) *JobHistory {
	s.JobVersion = &v
	return s
}

func (s *JobHistory) SetListStatus(v string) *JobHistory {
	s.ListStatus = &v
	return s
}

func (s *JobHistory) SetMessage(v string) *JobHistory {
	s.Message = &v
	return s
}

func (s *JobHistory) SetName(v string) *JobHistory {
	s.Name = &v
	return s
}

func (s *JobHistory) SetOperator(v string) *JobHistory {
	s.Operator = &v
	return s
}

func (s *JobHistory) SetRuntimeId(v string) *JobHistory {
	s.RuntimeId = &v
	return s
}

func (s *JobHistory) SetRuntimeState(v string) *JobHistory {
	s.RuntimeState = &v
	return s
}

func (s *JobHistory) SetStartTime(v string) *JobHistory {
	s.StartTime = &v
	return s
}

func (s *JobHistory) SetStatus(v string) *JobHistory {
	s.Status = &v
	return s
}

func (s *JobHistory) SetTotalCount(v int64) *JobHistory {
	s.TotalCount = &v
	return s
}

func (s *JobHistory) SetTotalSize(v int64) *JobHistory {
	s.TotalSize = &v
	return s
}

type KeyFilterItem struct {
	Regex []*string `json:"Regex,omitempty" xml:"Regex,omitempty" type:"Repeated"`
}

func (s KeyFilterItem) String() string {
	return tea.Prettify(s)
}

func (s KeyFilterItem) GoString() string {
	return s.String()
}

func (s *KeyFilterItem) SetRegex(v []*string) *KeyFilterItem {
	s.Regex = v
	return s
}

type KeyFilters struct {
	Excludes *KeyFilterItem `json:"Excludes,omitempty" xml:"Excludes,omitempty"`
	Includes *KeyFilterItem `json:"Includes,omitempty" xml:"Includes,omitempty"`
}

func (s KeyFilters) String() string {
	return tea.Prettify(s)
}

func (s KeyFilters) GoString() string {
	return s.String()
}

func (s *KeyFilters) SetExcludes(v *KeyFilterItem) *KeyFilters {
	s.Excludes = v
	return s
}

func (s *KeyFilters) SetIncludes(v *KeyFilterItem) *KeyFilters {
	s.Includes = v
	return s
}

type LastModifiedFilters struct {
	Excludes *LastModifyFilterItem `json:"Excludes,omitempty" xml:"Excludes,omitempty"`
	Includes *LastModifyFilterItem `json:"Includes,omitempty" xml:"Includes,omitempty"`
}

func (s LastModifiedFilters) String() string {
	return tea.Prettify(s)
}

func (s LastModifiedFilters) GoString() string {
	return s.String()
}

func (s *LastModifiedFilters) SetExcludes(v *LastModifyFilterItem) *LastModifiedFilters {
	s.Excludes = v
	return s
}

func (s *LastModifiedFilters) SetIncludes(v *LastModifyFilterItem) *LastModifiedFilters {
	s.Includes = v
	return s
}

type LastModifyFilterItem struct {
	TimeFilter []*TimeFilter `json:"TimeFilter,omitempty" xml:"TimeFilter,omitempty" type:"Repeated"`
}

func (s LastModifyFilterItem) String() string {
	return tea.Prettify(s)
}

func (s LastModifyFilterItem) GoString() string {
	return s.String()
}

func (s *LastModifyFilterItem) SetTimeFilter(v []*TimeFilter) *LastModifyFilterItem {
	s.TimeFilter = v
	return s
}

type ListAddressResp struct {
	ImportAddress []*GetAddressResp `json:"ImportAddress,omitempty" xml:"ImportAddress,omitempty" type:"Repeated"`
	// example:
	//
	// test_marker
	NextMarker *string `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	// example:
	//
	// true
	Truncated *bool `json:"Truncated,omitempty" xml:"Truncated,omitempty"`
}

func (s ListAddressResp) String() string {
	return tea.Prettify(s)
}

func (s ListAddressResp) GoString() string {
	return s.String()
}

func (s *ListAddressResp) SetImportAddress(v []*GetAddressResp) *ListAddressResp {
	s.ImportAddress = v
	return s
}

func (s *ListAddressResp) SetNextMarker(v string) *ListAddressResp {
	s.NextMarker = &v
	return s
}

func (s *ListAddressResp) SetTruncated(v bool) *ListAddressResp {
	s.Truncated = &v
	return s
}

type ListAgentResp struct {
	ImportAgent []*GetAgentResp `json:"ImportAgent,omitempty" xml:"ImportAgent,omitempty" type:"Repeated"`
	// example:
	//
	// test_next_marker
	NextMarker *string `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	// example:
	//
	// true
	Truncated *bool `json:"Truncated,omitempty" xml:"Truncated,omitempty"`
}

func (s ListAgentResp) String() string {
	return tea.Prettify(s)
}

func (s ListAgentResp) GoString() string {
	return s.String()
}

func (s *ListAgentResp) SetImportAgent(v []*GetAgentResp) *ListAgentResp {
	s.ImportAgent = v
	return s
}

func (s *ListAgentResp) SetNextMarker(v string) *ListAgentResp {
	s.NextMarker = &v
	return s
}

func (s *ListAgentResp) SetTruncated(v bool) *ListAgentResp {
	s.Truncated = &v
	return s
}

type ListJobHistoryResp struct {
	JobHistory []*JobHistory `json:"JobHistory,omitempty" xml:"JobHistory,omitempty" type:"Repeated"`
	// example:
	//
	// test_next_marker
	NextMarker *string `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	// example:
	//
	// true
	Truncated *string `json:"Truncated,omitempty" xml:"Truncated,omitempty"`
}

func (s ListJobHistoryResp) String() string {
	return tea.Prettify(s)
}

func (s ListJobHistoryResp) GoString() string {
	return s.String()
}

func (s *ListJobHistoryResp) SetJobHistory(v []*JobHistory) *ListJobHistoryResp {
	s.JobHistory = v
	return s
}

func (s *ListJobHistoryResp) SetNextMarker(v string) *ListJobHistoryResp {
	s.NextMarker = &v
	return s
}

func (s *ListJobHistoryResp) SetTruncated(v string) *ListJobHistoryResp {
	s.Truncated = &v
	return s
}

type ListJobInfo struct {
	ImportJob []*CreateJobInfo `json:"ImportJob,omitempty" xml:"ImportJob,omitempty" type:"Repeated"`
	// example:
	//
	// test_next_marker
	NextMarker *string `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	// example:
	//
	// true
	Truncated *bool `json:"Truncated,omitempty" xml:"Truncated,omitempty"`
}

func (s ListJobInfo) String() string {
	return tea.Prettify(s)
}

func (s ListJobInfo) GoString() string {
	return s.String()
}

func (s *ListJobInfo) SetImportJob(v []*CreateJobInfo) *ListJobInfo {
	s.ImportJob = v
	return s
}

func (s *ListJobInfo) SetNextMarker(v string) *ListJobInfo {
	s.NextMarker = &v
	return s
}

func (s *ListJobInfo) SetTruncated(v bool) *ListJobInfo {
	s.Truncated = &v
	return s
}

type ListJobResp struct {
	ImportJob  []*GetJobResp `json:"ImportJob,omitempty" xml:"ImportJob,omitempty" type:"Repeated"`
	NextMarker *string       `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	Truncated  *bool         `json:"Truncated,omitempty" xml:"Truncated,omitempty"`
}

func (s ListJobResp) String() string {
	return tea.Prettify(s)
}

func (s ListJobResp) GoString() string {
	return s.String()
}

func (s *ListJobResp) SetImportJob(v []*GetJobResp) *ListJobResp {
	s.ImportJob = v
	return s
}

func (s *ListJobResp) SetNextMarker(v string) *ListJobResp {
	s.NextMarker = &v
	return s
}

func (s *ListJobResp) SetTruncated(v bool) *ListJobResp {
	s.Truncated = &v
	return s
}

type ListTunnelResp struct {
	ImportTunnel []*GetTunnelResp `json:"ImportTunnel,omitempty" xml:"ImportTunnel,omitempty" type:"Repeated"`
	NextMarker   *string          `json:"NextMarker,omitempty" xml:"NextMarker,omitempty"`
	Truncated    *bool            `json:"Truncated,omitempty" xml:"Truncated,omitempty"`
}

func (s ListTunnelResp) String() string {
	return tea.Prettify(s)
}

func (s ListTunnelResp) GoString() string {
	return s.String()
}

func (s *ListTunnelResp) SetImportTunnel(v []*GetTunnelResp) *ListTunnelResp {
	s.ImportTunnel = v
	return s
}

func (s *ListTunnelResp) SetNextMarker(v string) *ListTunnelResp {
	s.NextMarker = &v
	return s
}

func (s *ListTunnelResp) SetTruncated(v bool) *ListTunnelResp {
	s.Truncated = &v
	return s
}

type ScheduleRule struct {
	MaxScheduleCount      *int64  `json:"MaxScheduleCount,omitempty" xml:"MaxScheduleCount,omitempty"`
	StartCronExpression   *string `json:"StartCronExpression,omitempty" xml:"StartCronExpression,omitempty"`
	SuspendCronExpression *string `json:"SuspendCronExpression,omitempty" xml:"SuspendCronExpression,omitempty"`
}

func (s ScheduleRule) String() string {
	return tea.Prettify(s)
}

func (s ScheduleRule) GoString() string {
	return s.String()
}

func (s *ScheduleRule) SetMaxScheduleCount(v int64) *ScheduleRule {
	s.MaxScheduleCount = &v
	return s
}

func (s *ScheduleRule) SetStartCronExpression(v string) *ScheduleRule {
	s.StartCronExpression = &v
	return s
}

func (s *ScheduleRule) SetSuspendCronExpression(v string) *ScheduleRule {
	s.SuspendCronExpression = &v
	return s
}

type TimeFilter struct {
	EndTime   *string `json:"EndTime,omitempty" xml:"EndTime,omitempty"`
	StartTime *string `json:"StartTime,omitempty" xml:"StartTime,omitempty"`
}

func (s TimeFilter) String() string {
	return tea.Prettify(s)
}

func (s TimeFilter) GoString() string {
	return s.String()
}

func (s *TimeFilter) SetEndTime(v string) *TimeFilter {
	s.EndTime = &v
	return s
}

func (s *TimeFilter) SetStartTime(v string) *TimeFilter {
	s.StartTime = &v
	return s
}

type TunnelQos struct {
	// example:
	//
	// 1073741824
	MaxBandwidth *int64 `json:"MaxBandwidth,omitempty" xml:"MaxBandwidth,omitempty"`
	// example:
	//
	// 100
	MaxQps *int32 `json:"MaxQps,omitempty" xml:"MaxQps,omitempty"`
}

func (s TunnelQos) String() string {
	return tea.Prettify(s)
}

func (s TunnelQos) GoString() string {
	return s.String()
}

func (s *TunnelQos) SetMaxBandwidth(v int64) *TunnelQos {
	s.MaxBandwidth = &v
	return s
}

func (s *TunnelQos) SetMaxQps(v int32) *TunnelQos {
	s.MaxQps = &v
	return s
}

type UpdateAddressInfo struct {
	AgentList *string `json:"AgentList,omitempty" xml:"AgentList,omitempty"`
}

func (s UpdateAddressInfo) String() string {
	return tea.Prettify(s)
}

func (s UpdateAddressInfo) GoString() string {
	return s.String()
}

func (s *UpdateAddressInfo) SetAgentList(v string) *UpdateAddressInfo {
	s.AgentList = &v
	return s
}

type UpdateJobInfo struct {
	ImportQos *ImportQos `json:"ImportQos,omitempty" xml:"ImportQos,omitempty"`
	// example:
	//
	// IMPORT_JOB_LAUNCHING
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
}

func (s UpdateJobInfo) String() string {
	return tea.Prettify(s)
}

func (s UpdateJobInfo) GoString() string {
	return s.String()
}

func (s *UpdateJobInfo) SetImportQos(v *ImportQos) *UpdateJobInfo {
	s.ImportQos = v
	return s
}

func (s *UpdateJobInfo) SetStatus(v string) *UpdateJobInfo {
	s.Status = &v
	return s
}

type UpdateTunnelInfo struct {
	Tags      *string    `json:"Tags,omitempty" xml:"Tags,omitempty"`
	TunnelQos *TunnelQos `json:"TunnelQos,omitempty" xml:"TunnelQos,omitempty"`
}

func (s UpdateTunnelInfo) String() string {
	return tea.Prettify(s)
}

func (s UpdateTunnelInfo) GoString() string {
	return s.String()
}

func (s *UpdateTunnelInfo) SetTags(v string) *UpdateTunnelInfo {
	s.Tags = &v
	return s
}

func (s *UpdateTunnelInfo) SetTunnelQos(v *TunnelQos) *UpdateTunnelInfo {
	s.TunnelQos = v
	return s
}

type VerifyAddressResp struct {
	ErrorCode    *string `json:"ErrorCode,omitempty" xml:"ErrorCode,omitempty"`
	ErrorMessage *string `json:"ErrorMessage,omitempty" xml:"ErrorMessage,omitempty"`
	// example:
	//
	// avaliable
	Status *string `json:"Status,omitempty" xml:"Status,omitempty"`
	// example:
	//
	// 2024-05-01 12:00:00
	VerifyTime *string `json:"VerifyTime,omitempty" xml:"VerifyTime,omitempty"`
}

func (s VerifyAddressResp) String() string {
	return tea.Prettify(s)
}

func (s VerifyAddressResp) GoString() string {
	return s.String()
}

func (s *VerifyAddressResp) SetErrorCode(v string) *VerifyAddressResp {
	s.ErrorCode = &v
	return s
}

func (s *VerifyAddressResp) SetErrorMessage(v string) *VerifyAddressResp {
	s.ErrorMessage = &v
	return s
}

func (s *VerifyAddressResp) SetStatus(v string) *VerifyAddressResp {
	s.Status = &v
	return s
}

func (s *VerifyAddressResp) SetVerifyTime(v string) *VerifyAddressResp {
	s.VerifyTime = &v
	return s
}

type VerifyResp struct {
	ErrorCode *string `json:"ErrorCode,omitempty" xml:"ErrorCode,omitempty"`
	ErrorMsg  *string `json:"ErrorMsg,omitempty" xml:"ErrorMsg,omitempty"`
	// example:
	//
	// 200
	HttpCode *string `json:"HttpCode,omitempty" xml:"HttpCode,omitempty"`
}

func (s VerifyResp) String() string {
	return tea.Prettify(s)
}

func (s VerifyResp) GoString() string {
	return s.String()
}

func (s *VerifyResp) SetErrorCode(v string) *VerifyResp {
	s.ErrorCode = &v
	return s
}

func (s *VerifyResp) SetErrorMsg(v string) *VerifyResp {
	s.ErrorMsg = &v
	return s
}

func (s *VerifyResp) SetHttpCode(v string) *VerifyResp {
	s.HttpCode = &v
	return s
}

type CreateAddressRequest struct {
	ImportAddress *CreateAddressInfo `json:"ImportAddress,omitempty" xml:"ImportAddress,omitempty"`
}

func (s CreateAddressRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAddressRequest) GoString() string {
	return s.String()
}

func (s *CreateAddressRequest) SetImportAddress(v *CreateAddressInfo) *CreateAddressRequest {
	s.ImportAddress = v
	return s
}

type CreateAddressResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s CreateAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAddressResponse) GoString() string {
	return s.String()
}

func (s *CreateAddressResponse) SetHeaders(v map[string]*string) *CreateAddressResponse {
	s.Headers = v
	return s
}

func (s *CreateAddressResponse) SetStatusCode(v int32) *CreateAddressResponse {
	s.StatusCode = &v
	return s
}

type CreateAgentRequest struct {
	ImportAgent *CreateAgentInfo `json:"ImportAgent,omitempty" xml:"ImportAgent,omitempty"`
}

func (s CreateAgentRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAgentRequest) GoString() string {
	return s.String()
}

func (s *CreateAgentRequest) SetImportAgent(v *CreateAgentInfo) *CreateAgentRequest {
	s.ImportAgent = v
	return s
}

type CreateAgentResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s CreateAgentResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAgentResponse) GoString() string {
	return s.String()
}

func (s *CreateAgentResponse) SetHeaders(v map[string]*string) *CreateAgentResponse {
	s.Headers = v
	return s
}

func (s *CreateAgentResponse) SetStatusCode(v int32) *CreateAgentResponse {
	s.StatusCode = &v
	return s
}

type CreateJobRequest struct {
	// This parameter is required.
	ImportJob *CreateJobInfo `json:"ImportJob,omitempty" xml:"ImportJob,omitempty"`
}

func (s CreateJobRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateJobRequest) GoString() string {
	return s.String()
}

func (s *CreateJobRequest) SetImportJob(v *CreateJobInfo) *CreateJobRequest {
	s.ImportJob = v
	return s
}

type CreateJobResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s CreateJobResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateJobResponse) GoString() string {
	return s.String()
}

func (s *CreateJobResponse) SetHeaders(v map[string]*string) *CreateJobResponse {
	s.Headers = v
	return s
}

func (s *CreateJobResponse) SetStatusCode(v int32) *CreateJobResponse {
	s.StatusCode = &v
	return s
}

type CreateReportRequest struct {
	CreateReport *CreateReportInfo `json:"CreateReport,omitempty" xml:"CreateReport,omitempty"`
}

func (s CreateReportRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateReportRequest) GoString() string {
	return s.String()
}

func (s *CreateReportRequest) SetCreateReport(v *CreateReportInfo) *CreateReportRequest {
	s.CreateReport = v
	return s
}

type CreateReportResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s CreateReportResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateReportResponse) GoString() string {
	return s.String()
}

func (s *CreateReportResponse) SetHeaders(v map[string]*string) *CreateReportResponse {
	s.Headers = v
	return s
}

func (s *CreateReportResponse) SetStatusCode(v int32) *CreateReportResponse {
	s.StatusCode = &v
	return s
}

type CreateTunnelRequest struct {
	ImportTunnel *CreateTunnelInfo `json:"ImportTunnel,omitempty" xml:"ImportTunnel,omitempty"`
}

func (s CreateTunnelRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateTunnelRequest) GoString() string {
	return s.String()
}

func (s *CreateTunnelRequest) SetImportTunnel(v *CreateTunnelInfo) *CreateTunnelRequest {
	s.ImportTunnel = v
	return s
}

type CreateTunnelResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s CreateTunnelResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateTunnelResponse) GoString() string {
	return s.String()
}

func (s *CreateTunnelResponse) SetHeaders(v map[string]*string) *CreateTunnelResponse {
	s.Headers = v
	return s
}

func (s *CreateTunnelResponse) SetStatusCode(v int32) *CreateTunnelResponse {
	s.StatusCode = &v
	return s
}

type DeleteAddressResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAddressResponse) GoString() string {
	return s.String()
}

func (s *DeleteAddressResponse) SetHeaders(v map[string]*string) *DeleteAddressResponse {
	s.Headers = v
	return s
}

func (s *DeleteAddressResponse) SetStatusCode(v int32) *DeleteAddressResponse {
	s.StatusCode = &v
	return s
}

type DeleteAgentResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteAgentResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAgentResponse) GoString() string {
	return s.String()
}

func (s *DeleteAgentResponse) SetHeaders(v map[string]*string) *DeleteAgentResponse {
	s.Headers = v
	return s
}

func (s *DeleteAgentResponse) SetStatusCode(v int32) *DeleteAgentResponse {
	s.StatusCode = &v
	return s
}

type DeleteJobRequest struct {
	// example:
	//
	// true
	ForceDelete *bool `json:"forceDelete,omitempty" xml:"forceDelete,omitempty"`
}

func (s DeleteJobRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteJobRequest) GoString() string {
	return s.String()
}

func (s *DeleteJobRequest) SetForceDelete(v bool) *DeleteJobRequest {
	s.ForceDelete = &v
	return s
}

type DeleteJobResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteJobResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteJobResponse) GoString() string {
	return s.String()
}

func (s *DeleteJobResponse) SetHeaders(v map[string]*string) *DeleteJobResponse {
	s.Headers = v
	return s
}

func (s *DeleteJobResponse) SetStatusCode(v int32) *DeleteJobResponse {
	s.StatusCode = &v
	return s
}

type DeleteTunnelResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s DeleteTunnelResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteTunnelResponse) GoString() string {
	return s.String()
}

func (s *DeleteTunnelResponse) SetHeaders(v map[string]*string) *DeleteTunnelResponse {
	s.Headers = v
	return s
}

func (s *DeleteTunnelResponse) SetStatusCode(v int32) *DeleteTunnelResponse {
	s.StatusCode = &v
	return s
}

type GetAddressResponseBody struct {
	// 222
	ImportAddress *GetAddressResp `json:"ImportAddress,omitempty" xml:"ImportAddress,omitempty"`
}

func (s GetAddressResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAddressResponseBody) GoString() string {
	return s.String()
}

func (s *GetAddressResponseBody) SetImportAddress(v *GetAddressResp) *GetAddressResponseBody {
	s.ImportAddress = v
	return s
}

type GetAddressResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetAddressResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAddressResponse) GoString() string {
	return s.String()
}

func (s *GetAddressResponse) SetHeaders(v map[string]*string) *GetAddressResponse {
	s.Headers = v
	return s
}

func (s *GetAddressResponse) SetStatusCode(v int32) *GetAddressResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAddressResponse) SetBody(v *GetAddressResponseBody) *GetAddressResponse {
	s.Body = v
	return s
}

type GetAgentResponseBody struct {
	// 2
	ImportAgent *GetAgentResp `json:"ImportAgent,omitempty" xml:"ImportAgent,omitempty"`
}

func (s GetAgentResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAgentResponseBody) GoString() string {
	return s.String()
}

func (s *GetAgentResponseBody) SetImportAgent(v *GetAgentResp) *GetAgentResponseBody {
	s.ImportAgent = v
	return s
}

type GetAgentResponse struct {
	Headers    map[string]*string    `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetAgentResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAgentResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAgentResponse) GoString() string {
	return s.String()
}

func (s *GetAgentResponse) SetHeaders(v map[string]*string) *GetAgentResponse {
	s.Headers = v
	return s
}

func (s *GetAgentResponse) SetStatusCode(v int32) *GetAgentResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAgentResponse) SetBody(v *GetAgentResponseBody) *GetAgentResponse {
	s.Body = v
	return s
}

type GetAgentStatusResponseBody struct {
	// 2
	ImportAgentStatus *GetAgentStatusResp `json:"ImportAgentStatus,omitempty" xml:"ImportAgentStatus,omitempty"`
}

func (s GetAgentStatusResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetAgentStatusResponseBody) GoString() string {
	return s.String()
}

func (s *GetAgentStatusResponseBody) SetImportAgentStatus(v *GetAgentStatusResp) *GetAgentStatusResponseBody {
	s.ImportAgentStatus = v
	return s
}

type GetAgentStatusResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetAgentStatusResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetAgentStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s GetAgentStatusResponse) GoString() string {
	return s.String()
}

func (s *GetAgentStatusResponse) SetHeaders(v map[string]*string) *GetAgentStatusResponse {
	s.Headers = v
	return s
}

func (s *GetAgentStatusResponse) SetStatusCode(v int32) *GetAgentStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *GetAgentStatusResponse) SetBody(v *GetAgentStatusResponseBody) *GetAgentStatusResponse {
	s.Body = v
	return s
}

type GetJobRequest struct {
	// example:
	//
	// false
	ByVersion *string `json:"byVersion,omitempty" xml:"byVersion,omitempty"`
}

func (s GetJobRequest) String() string {
	return tea.Prettify(s)
}

func (s GetJobRequest) GoString() string {
	return s.String()
}

func (s *GetJobRequest) SetByVersion(v string) *GetJobRequest {
	s.ByVersion = &v
	return s
}

type GetJobResponseBody struct {
	ImportJob *GetJobResp `json:"ImportJob,omitempty" xml:"ImportJob,omitempty"`
}

func (s GetJobResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetJobResponseBody) GoString() string {
	return s.String()
}

func (s *GetJobResponseBody) SetImportJob(v *GetJobResp) *GetJobResponseBody {
	s.ImportJob = v
	return s
}

type GetJobResponse struct {
	Headers    map[string]*string  `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32              `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetJobResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetJobResponse) String() string {
	return tea.Prettify(s)
}

func (s GetJobResponse) GoString() string {
	return s.String()
}

func (s *GetJobResponse) SetHeaders(v map[string]*string) *GetJobResponse {
	s.Headers = v
	return s
}

func (s *GetJobResponse) SetStatusCode(v int32) *GetJobResponse {
	s.StatusCode = &v
	return s
}

func (s *GetJobResponse) SetBody(v *GetJobResponseBody) *GetJobResponse {
	s.Body = v
	return s
}

type GetJobResultRequest struct {
	// This parameter is required.
	//
	// example:
	//
	// 2
	RuntimeId *int32 `json:"runtimeId,omitempty" xml:"runtimeId,omitempty"`
}

func (s GetJobResultRequest) String() string {
	return tea.Prettify(s)
}

func (s GetJobResultRequest) GoString() string {
	return s.String()
}

func (s *GetJobResultRequest) SetRuntimeId(v int32) *GetJobResultRequest {
	s.RuntimeId = &v
	return s
}

type GetJobResultResponseBody struct {
	// 1
	ImportJobResult *GetJobResultResp `json:"ImportJobResult,omitempty" xml:"ImportJobResult,omitempty"`
}

func (s GetJobResultResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetJobResultResponseBody) GoString() string {
	return s.String()
}

func (s *GetJobResultResponseBody) SetImportJobResult(v *GetJobResultResp) *GetJobResultResponseBody {
	s.ImportJobResult = v
	return s
}

type GetJobResultResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetJobResultResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetJobResultResponse) String() string {
	return tea.Prettify(s)
}

func (s GetJobResultResponse) GoString() string {
	return s.String()
}

func (s *GetJobResultResponse) SetHeaders(v map[string]*string) *GetJobResultResponse {
	s.Headers = v
	return s
}

func (s *GetJobResultResponse) SetStatusCode(v int32) *GetJobResultResponse {
	s.StatusCode = &v
	return s
}

func (s *GetJobResultResponse) SetBody(v *GetJobResultResponseBody) *GetJobResultResponse {
	s.Body = v
	return s
}

type GetReportRequest struct {
	// example:
	//
	// 1
	RuntimeId *int32 `json:"runtimeId,omitempty" xml:"runtimeId,omitempty"`
	// This parameter is required.
	//
	// example:
	//
	// test_job_id
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s GetReportRequest) String() string {
	return tea.Prettify(s)
}

func (s GetReportRequest) GoString() string {
	return s.String()
}

func (s *GetReportRequest) SetRuntimeId(v int32) *GetReportRequest {
	s.RuntimeId = &v
	return s
}

func (s *GetReportRequest) SetVersion(v string) *GetReportRequest {
	s.Version = &v
	return s
}

type GetReportResponseBody struct {
	GetReportResponse *GetReportResp `json:"GetReportResponse,omitempty" xml:"GetReportResponse,omitempty"`
}

func (s GetReportResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetReportResponseBody) GoString() string {
	return s.String()
}

func (s *GetReportResponseBody) SetGetReportResponse(v *GetReportResp) *GetReportResponseBody {
	s.GetReportResponse = v
	return s
}

type GetReportResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetReportResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetReportResponse) String() string {
	return tea.Prettify(s)
}

func (s GetReportResponse) GoString() string {
	return s.String()
}

func (s *GetReportResponse) SetHeaders(v map[string]*string) *GetReportResponse {
	s.Headers = v
	return s
}

func (s *GetReportResponse) SetStatusCode(v int32) *GetReportResponse {
	s.StatusCode = &v
	return s
}

func (s *GetReportResponse) SetBody(v *GetReportResponseBody) *GetReportResponse {
	s.Body = v
	return s
}

type GetTunnelResponseBody struct {
	ImportTunnel *GetTunnelResp `json:"ImportTunnel,omitempty" xml:"ImportTunnel,omitempty"`
}

func (s GetTunnelResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetTunnelResponseBody) GoString() string {
	return s.String()
}

func (s *GetTunnelResponseBody) SetImportTunnel(v *GetTunnelResp) *GetTunnelResponseBody {
	s.ImportTunnel = v
	return s
}

type GetTunnelResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *GetTunnelResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s GetTunnelResponse) String() string {
	return tea.Prettify(s)
}

func (s GetTunnelResponse) GoString() string {
	return s.String()
}

func (s *GetTunnelResponse) SetHeaders(v map[string]*string) *GetTunnelResponse {
	s.Headers = v
	return s
}

func (s *GetTunnelResponse) SetStatusCode(v int32) *GetTunnelResponse {
	s.StatusCode = &v
	return s
}

func (s *GetTunnelResponse) SetBody(v *GetTunnelResponseBody) *GetTunnelResponse {
	s.Body = v
	return s
}

type ListAddressRequest struct {
	// example:
	//
	// 100
	Count *int32 `json:"count,omitempty" xml:"count,omitempty"`
	// example:
	//
	// test_marker
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
}

func (s ListAddressRequest) String() string {
	return tea.Prettify(s)
}

func (s ListAddressRequest) GoString() string {
	return s.String()
}

func (s *ListAddressRequest) SetCount(v int32) *ListAddressRequest {
	s.Count = &v
	return s
}

func (s *ListAddressRequest) SetMarker(v string) *ListAddressRequest {
	s.Marker = &v
	return s
}

type ListAddressResponseBody struct {
	ImportAddressList *ListAddressResp `json:"ImportAddressList,omitempty" xml:"ImportAddressList,omitempty"`
}

func (s ListAddressResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListAddressResponseBody) GoString() string {
	return s.String()
}

func (s *ListAddressResponseBody) SetImportAddressList(v *ListAddressResp) *ListAddressResponseBody {
	s.ImportAddressList = v
	return s
}

type ListAddressResponse struct {
	Headers    map[string]*string       `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                   `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListAddressResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s ListAddressResponse) GoString() string {
	return s.String()
}

func (s *ListAddressResponse) SetHeaders(v map[string]*string) *ListAddressResponse {
	s.Headers = v
	return s
}

func (s *ListAddressResponse) SetStatusCode(v int32) *ListAddressResponse {
	s.StatusCode = &v
	return s
}

func (s *ListAddressResponse) SetBody(v *ListAddressResponseBody) *ListAddressResponse {
	s.Body = v
	return s
}

type ListAgentRequest struct {
	// example:
	//
	// 100
	Count *int32 `json:"count,omitempty" xml:"count,omitempty"`
	// example:
	//
	// test_agent
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
}

func (s ListAgentRequest) String() string {
	return tea.Prettify(s)
}

func (s ListAgentRequest) GoString() string {
	return s.String()
}

func (s *ListAgentRequest) SetCount(v int32) *ListAgentRequest {
	s.Count = &v
	return s
}

func (s *ListAgentRequest) SetMarker(v string) *ListAgentRequest {
	s.Marker = &v
	return s
}

type ListAgentResponseBody struct {
	ImportAgentList *ListAgentResp `json:"ImportAgentList,omitempty" xml:"ImportAgentList,omitempty"`
}

func (s ListAgentResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListAgentResponseBody) GoString() string {
	return s.String()
}

func (s *ListAgentResponseBody) SetImportAgentList(v *ListAgentResp) *ListAgentResponseBody {
	s.ImportAgentList = v
	return s
}

type ListAgentResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListAgentResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListAgentResponse) String() string {
	return tea.Prettify(s)
}

func (s ListAgentResponse) GoString() string {
	return s.String()
}

func (s *ListAgentResponse) SetHeaders(v map[string]*string) *ListAgentResponse {
	s.Headers = v
	return s
}

func (s *ListAgentResponse) SetStatusCode(v int32) *ListAgentResponse {
	s.StatusCode = &v
	return s
}

func (s *ListAgentResponse) SetBody(v *ListAgentResponseBody) *ListAgentResponse {
	s.Body = v
	return s
}

type ListJobRequest struct {
	// example:
	//
	// true
	All *bool `json:"all,omitempty" xml:"all,omitempty"`
	// example:
	//
	// 1000
	Count *int32 `json:"count,omitempty" xml:"count,omitempty"`
	// example:
	//
	// test_marker
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
	// example:
	//
	// test_parent_job_name
	ParentName *string `json:"parentName,omitempty" xml:"parentName,omitempty"`
}

func (s ListJobRequest) String() string {
	return tea.Prettify(s)
}

func (s ListJobRequest) GoString() string {
	return s.String()
}

func (s *ListJobRequest) SetAll(v bool) *ListJobRequest {
	s.All = &v
	return s
}

func (s *ListJobRequest) SetCount(v int32) *ListJobRequest {
	s.Count = &v
	return s
}

func (s *ListJobRequest) SetMarker(v string) *ListJobRequest {
	s.Marker = &v
	return s
}

func (s *ListJobRequest) SetParentName(v string) *ListJobRequest {
	s.ParentName = &v
	return s
}

type ListJobResponseBody struct {
	ImportJobList *ListJobResp `json:"ImportJobList,omitempty" xml:"ImportJobList,omitempty"`
}

func (s ListJobResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListJobResponseBody) GoString() string {
	return s.String()
}

func (s *ListJobResponseBody) SetImportJobList(v *ListJobResp) *ListJobResponseBody {
	s.ImportJobList = v
	return s
}

type ListJobResponse struct {
	Headers    map[string]*string   `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32               `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListJobResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListJobResponse) String() string {
	return tea.Prettify(s)
}

func (s ListJobResponse) GoString() string {
	return s.String()
}

func (s *ListJobResponse) SetHeaders(v map[string]*string) *ListJobResponse {
	s.Headers = v
	return s
}

func (s *ListJobResponse) SetStatusCode(v int32) *ListJobResponse {
	s.StatusCode = &v
	return s
}

func (s *ListJobResponse) SetBody(v *ListJobResponseBody) *ListJobResponse {
	s.Body = v
	return s
}

type ListJobHistoryRequest struct {
	// example:
	//
	// 100
	Count  *int32  `json:"count,omitempty" xml:"count,omitempty"`
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
	// example:
	//
	// 1
	RuntimeId *int32 `json:"runtimeId,omitempty" xml:"runtimeId,omitempty"`
}

func (s ListJobHistoryRequest) String() string {
	return tea.Prettify(s)
}

func (s ListJobHistoryRequest) GoString() string {
	return s.String()
}

func (s *ListJobHistoryRequest) SetCount(v int32) *ListJobHistoryRequest {
	s.Count = &v
	return s
}

func (s *ListJobHistoryRequest) SetMarker(v string) *ListJobHistoryRequest {
	s.Marker = &v
	return s
}

func (s *ListJobHistoryRequest) SetRuntimeId(v int32) *ListJobHistoryRequest {
	s.RuntimeId = &v
	return s
}

type ListJobHistoryResponseBody struct {
	JobHistoryList *ListJobHistoryResp `json:"JobHistoryList,omitempty" xml:"JobHistoryList,omitempty"`
}

func (s ListJobHistoryResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListJobHistoryResponseBody) GoString() string {
	return s.String()
}

func (s *ListJobHistoryResponseBody) SetJobHistoryList(v *ListJobHistoryResp) *ListJobHistoryResponseBody {
	s.JobHistoryList = v
	return s
}

type ListJobHistoryResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListJobHistoryResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListJobHistoryResponse) String() string {
	return tea.Prettify(s)
}

func (s ListJobHistoryResponse) GoString() string {
	return s.String()
}

func (s *ListJobHistoryResponse) SetHeaders(v map[string]*string) *ListJobHistoryResponse {
	s.Headers = v
	return s
}

func (s *ListJobHistoryResponse) SetStatusCode(v int32) *ListJobHistoryResponse {
	s.StatusCode = &v
	return s
}

func (s *ListJobHistoryResponse) SetBody(v *ListJobHistoryResponseBody) *ListJobHistoryResponse {
	s.Body = v
	return s
}

type ListTunnelRequest struct {
	// example:
	//
	// 2
	Count *int32 `json:"count,omitempty" xml:"count,omitempty"`
	// example:
	//
	// 1
	Marker *string `json:"marker,omitempty" xml:"marker,omitempty"`
}

func (s ListTunnelRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTunnelRequest) GoString() string {
	return s.String()
}

func (s *ListTunnelRequest) SetCount(v int32) *ListTunnelRequest {
	s.Count = &v
	return s
}

func (s *ListTunnelRequest) SetMarker(v string) *ListTunnelRequest {
	s.Marker = &v
	return s
}

type ListTunnelResponseBody struct {
	// 2
	ImportTunnelList *ListTunnelResp `json:"ImportTunnelList,omitempty" xml:"ImportTunnelList,omitempty"`
}

func (s ListTunnelResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListTunnelResponseBody) GoString() string {
	return s.String()
}

func (s *ListTunnelResponseBody) SetImportTunnelList(v *ListTunnelResp) *ListTunnelResponseBody {
	s.ImportTunnelList = v
	return s
}

type ListTunnelResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *ListTunnelResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s ListTunnelResponse) String() string {
	return tea.Prettify(s)
}

func (s ListTunnelResponse) GoString() string {
	return s.String()
}

func (s *ListTunnelResponse) SetHeaders(v map[string]*string) *ListTunnelResponse {
	s.Headers = v
	return s
}

func (s *ListTunnelResponse) SetStatusCode(v int32) *ListTunnelResponse {
	s.StatusCode = &v
	return s
}

func (s *ListTunnelResponse) SetBody(v *ListTunnelResponseBody) *ListTunnelResponse {
	s.Body = v
	return s
}

type UpdateAddressRequest struct {
	ImportAddress *UpdateAddressInfo `json:"ImportAddress,omitempty" xml:"ImportAddress,omitempty"`
}

func (s UpdateAddressRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateAddressRequest) GoString() string {
	return s.String()
}

func (s *UpdateAddressRequest) SetImportAddress(v *UpdateAddressInfo) *UpdateAddressRequest {
	s.ImportAddress = v
	return s
}

type UpdateAddressResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s UpdateAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateAddressResponse) GoString() string {
	return s.String()
}

func (s *UpdateAddressResponse) SetHeaders(v map[string]*string) *UpdateAddressResponse {
	s.Headers = v
	return s
}

func (s *UpdateAddressResponse) SetStatusCode(v int32) *UpdateAddressResponse {
	s.StatusCode = &v
	return s
}

type UpdateJobRequest struct {
	ImportJob *UpdateJobInfo `json:"ImportJob,omitempty" xml:"ImportJob,omitempty"`
}

func (s UpdateJobRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateJobRequest) GoString() string {
	return s.String()
}

func (s *UpdateJobRequest) SetImportJob(v *UpdateJobInfo) *UpdateJobRequest {
	s.ImportJob = v
	return s
}

type UpdateJobResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s UpdateJobResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateJobResponse) GoString() string {
	return s.String()
}

func (s *UpdateJobResponse) SetHeaders(v map[string]*string) *UpdateJobResponse {
	s.Headers = v
	return s
}

func (s *UpdateJobResponse) SetStatusCode(v int32) *UpdateJobResponse {
	s.StatusCode = &v
	return s
}

type UpdateTunnelRequest struct {
	ImportTunnel *UpdateTunnelInfo `json:"ImportTunnel,omitempty" xml:"ImportTunnel,omitempty"`
}

func (s UpdateTunnelRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateTunnelRequest) GoString() string {
	return s.String()
}

func (s *UpdateTunnelRequest) SetImportTunnel(v *UpdateTunnelInfo) *UpdateTunnelRequest {
	s.ImportTunnel = v
	return s
}

type UpdateTunnelResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
}

func (s UpdateTunnelResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateTunnelResponse) GoString() string {
	return s.String()
}

func (s *UpdateTunnelResponse) SetHeaders(v map[string]*string) *UpdateTunnelResponse {
	s.Headers = v
	return s
}

func (s *UpdateTunnelResponse) SetStatusCode(v int32) *UpdateTunnelResponse {
	s.StatusCode = &v
	return s
}

type VerifyAddressResponseBody struct {
	// 1
	VerifyAddressResponse *VerifyAddressResp `json:"VerifyAddressResponse,omitempty" xml:"VerifyAddressResponse,omitempty"`
}

func (s VerifyAddressResponseBody) String() string {
	return tea.Prettify(s)
}

func (s VerifyAddressResponseBody) GoString() string {
	return s.String()
}

func (s *VerifyAddressResponseBody) SetVerifyAddressResponse(v *VerifyAddressResp) *VerifyAddressResponseBody {
	s.VerifyAddressResponse = v
	return s
}

type VerifyAddressResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty"`
	Body       *VerifyAddressResponseBody `json:"body,omitempty" xml:"body,omitempty"`
}

func (s VerifyAddressResponse) String() string {
	return tea.Prettify(s)
}

func (s VerifyAddressResponse) GoString() string {
	return s.String()
}

func (s *VerifyAddressResponse) SetHeaders(v map[string]*string) *VerifyAddressResponse {
	s.Headers = v
	return s
}

func (s *VerifyAddressResponse) SetStatusCode(v int32) *VerifyAddressResponse {
	s.StatusCode = &v
	return s
}

func (s *VerifyAddressResponse) SetBody(v *VerifyAddressResponseBody) *VerifyAddressResponse {
	s.Body = v
	return s
}
