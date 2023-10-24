// This file is auto-generated, don't edit it. Thanks.
/**
 *
 */
package client

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	endpointutil "github.com/alibabacloud-go/endpoint-util/service"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Addon struct {
	Config   *string `json:"config,omitempty" xml:"config,omitempty"`
	Disabled *bool   `json:"disabled,omitempty" xml:"disabled,omitempty"`
	Name     *string `json:"name,omitempty" xml:"name,omitempty"`
}

func (s Addon) String() string {
	return tea.Prettify(s)
}

func (s Addon) GoString() string {
	return s.String()
}

func (s *Addon) SetConfig(v string) *Addon {
	s.Config = &v
	return s
}

func (s *Addon) SetDisabled(v bool) *Addon {
	s.Disabled = &v
	return s
}

func (s *Addon) SetName(v string) *Addon {
	s.Name = &v
	return s
}

type DataDisk struct {
	AutoSnapshotPolicyId *string `json:"auto_snapshot_policy_id,omitempty" xml:"auto_snapshot_policy_id,omitempty"`
	BurstingEnabled      *bool   `json:"bursting_enabled,omitempty" xml:"bursting_enabled,omitempty"`
	Category             *string `json:"category,omitempty" xml:"category,omitempty"`
	Encrypted            *string `json:"encrypted,omitempty" xml:"encrypted,omitempty"`
	PerformanceLevel     *string `json:"performance_level,omitempty" xml:"performance_level,omitempty"`
	ProvisionedIops      *int64  `json:"provisioned_iops,omitempty" xml:"provisioned_iops,omitempty"`
	Size                 *int64  `json:"size,omitempty" xml:"size,omitempty"`
}

func (s DataDisk) String() string {
	return tea.Prettify(s)
}

func (s DataDisk) GoString() string {
	return s.String()
}

func (s *DataDisk) SetAutoSnapshotPolicyId(v string) *DataDisk {
	s.AutoSnapshotPolicyId = &v
	return s
}

func (s *DataDisk) SetBurstingEnabled(v bool) *DataDisk {
	s.BurstingEnabled = &v
	return s
}

func (s *DataDisk) SetCategory(v string) *DataDisk {
	s.Category = &v
	return s
}

func (s *DataDisk) SetEncrypted(v string) *DataDisk {
	s.Encrypted = &v
	return s
}

func (s *DataDisk) SetPerformanceLevel(v string) *DataDisk {
	s.PerformanceLevel = &v
	return s
}

func (s *DataDisk) SetProvisionedIops(v int64) *DataDisk {
	s.ProvisionedIops = &v
	return s
}

func (s *DataDisk) SetSize(v int64) *DataDisk {
	s.Size = &v
	return s
}

type MaintenanceWindow struct {
	Duration        *string `json:"duration,omitempty" xml:"duration,omitempty"`
	Enable          *bool   `json:"enable,omitempty" xml:"enable,omitempty"`
	MaintenanceTime *string `json:"maintenance_time,omitempty" xml:"maintenance_time,omitempty"`
	WeeklyPeriod    *string `json:"weekly_period,omitempty" xml:"weekly_period,omitempty"`
}

func (s MaintenanceWindow) String() string {
	return tea.Prettify(s)
}

func (s MaintenanceWindow) GoString() string {
	return s.String()
}

func (s *MaintenanceWindow) SetDuration(v string) *MaintenanceWindow {
	s.Duration = &v
	return s
}

func (s *MaintenanceWindow) SetEnable(v bool) *MaintenanceWindow {
	s.Enable = &v
	return s
}

func (s *MaintenanceWindow) SetMaintenanceTime(v string) *MaintenanceWindow {
	s.MaintenanceTime = &v
	return s
}

func (s *MaintenanceWindow) SetWeeklyPeriod(v string) *MaintenanceWindow {
	s.WeeklyPeriod = &v
	return s
}

type Nodepool struct {
	AutoScaling *NodepoolAutoScaling `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	// Deprecated
	Count *int64 `json:"count,omitempty" xml:"count,omitempty"`
	// Deprecated
	InterconnectConfig *NodepoolInterconnectConfig `json:"interconnect_config,omitempty" xml:"interconnect_config,omitempty" type:"Struct"`
	InterconnectMode   *string                     `json:"interconnect_mode,omitempty" xml:"interconnect_mode,omitempty"`
	KubernetesConfig   *NodepoolKubernetesConfig   `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	Management         *NodepoolManagement         `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
	MaxNodes           *int64                      `json:"max_nodes,omitempty" xml:"max_nodes,omitempty"`
	NodepoolInfo       *NodepoolNodepoolInfo       `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	ScalingGroup       *NodepoolScalingGroup       `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	TeeConfig          *NodepoolTeeConfig          `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
}

func (s Nodepool) String() string {
	return tea.Prettify(s)
}

func (s Nodepool) GoString() string {
	return s.String()
}

func (s *Nodepool) SetAutoScaling(v *NodepoolAutoScaling) *Nodepool {
	s.AutoScaling = v
	return s
}

func (s *Nodepool) SetCount(v int64) *Nodepool {
	s.Count = &v
	return s
}

func (s *Nodepool) SetInterconnectConfig(v *NodepoolInterconnectConfig) *Nodepool {
	s.InterconnectConfig = v
	return s
}

func (s *Nodepool) SetInterconnectMode(v string) *Nodepool {
	s.InterconnectMode = &v
	return s
}

func (s *Nodepool) SetKubernetesConfig(v *NodepoolKubernetesConfig) *Nodepool {
	s.KubernetesConfig = v
	return s
}

func (s *Nodepool) SetManagement(v *NodepoolManagement) *Nodepool {
	s.Management = v
	return s
}

func (s *Nodepool) SetMaxNodes(v int64) *Nodepool {
	s.MaxNodes = &v
	return s
}

func (s *Nodepool) SetNodepoolInfo(v *NodepoolNodepoolInfo) *Nodepool {
	s.NodepoolInfo = v
	return s
}

func (s *Nodepool) SetScalingGroup(v *NodepoolScalingGroup) *Nodepool {
	s.ScalingGroup = v
	return s
}

func (s *Nodepool) SetTeeConfig(v *NodepoolTeeConfig) *Nodepool {
	s.TeeConfig = v
	return s
}

type NodepoolAutoScaling struct {
	// Deprecated
	EipBandwidth *int64 `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	// Deprecated
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	Enable                *bool   `json:"enable,omitempty" xml:"enable,omitempty"`
	// Deprecated
	IsBondEip    *bool   `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	MaxInstances *int64  `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	MinInstances *int64  `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	Type         *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s NodepoolAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s NodepoolAutoScaling) GoString() string {
	return s.String()
}

func (s *NodepoolAutoScaling) SetEipBandwidth(v int64) *NodepoolAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *NodepoolAutoScaling) SetEipInternetChargeType(v string) *NodepoolAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *NodepoolAutoScaling) SetEnable(v bool) *NodepoolAutoScaling {
	s.Enable = &v
	return s
}

func (s *NodepoolAutoScaling) SetIsBondEip(v bool) *NodepoolAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *NodepoolAutoScaling) SetMaxInstances(v int64) *NodepoolAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *NodepoolAutoScaling) SetMinInstances(v int64) *NodepoolAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *NodepoolAutoScaling) SetType(v string) *NodepoolAutoScaling {
	s.Type = &v
	return s
}

type NodepoolInterconnectConfig struct {
	// Deprecated
	Bandwidth *int64 `json:"bandwidth,omitempty" xml:"bandwidth,omitempty"`
	// Deprecated
	CcnId *string `json:"ccn_id,omitempty" xml:"ccn_id,omitempty"`
	// Deprecated
	CcnRegionId *string `json:"ccn_region_id,omitempty" xml:"ccn_region_id,omitempty"`
	// Deprecated
	CenId *string `json:"cen_id,omitempty" xml:"cen_id,omitempty"`
	// Deprecated
	ImprovedPeriod *string `json:"improved_period,omitempty" xml:"improved_period,omitempty"`
}

func (s NodepoolInterconnectConfig) String() string {
	return tea.Prettify(s)
}

func (s NodepoolInterconnectConfig) GoString() string {
	return s.String()
}

func (s *NodepoolInterconnectConfig) SetBandwidth(v int64) *NodepoolInterconnectConfig {
	s.Bandwidth = &v
	return s
}

func (s *NodepoolInterconnectConfig) SetCcnId(v string) *NodepoolInterconnectConfig {
	s.CcnId = &v
	return s
}

func (s *NodepoolInterconnectConfig) SetCcnRegionId(v string) *NodepoolInterconnectConfig {
	s.CcnRegionId = &v
	return s
}

func (s *NodepoolInterconnectConfig) SetCenId(v string) *NodepoolInterconnectConfig {
	s.CenId = &v
	return s
}

func (s *NodepoolInterconnectConfig) SetImprovedPeriod(v string) *NodepoolInterconnectConfig {
	s.ImprovedPeriod = &v
	return s
}

type NodepoolKubernetesConfig struct {
	CmsEnabled     *bool    `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	CpuPolicy      *string  `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	Labels         []*Tag   `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	NodeNameMode   *string  `json:"node_name_mode,omitempty" xml:"node_name_mode,omitempty"`
	Runtime        *string  `json:"runtime,omitempty" xml:"runtime,omitempty"`
	RuntimeVersion *string  `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	Taints         []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	UserData       *string  `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s NodepoolKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s NodepoolKubernetesConfig) GoString() string {
	return s.String()
}

func (s *NodepoolKubernetesConfig) SetCmsEnabled(v bool) *NodepoolKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *NodepoolKubernetesConfig) SetCpuPolicy(v string) *NodepoolKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *NodepoolKubernetesConfig) SetLabels(v []*Tag) *NodepoolKubernetesConfig {
	s.Labels = v
	return s
}

func (s *NodepoolKubernetesConfig) SetNodeNameMode(v string) *NodepoolKubernetesConfig {
	s.NodeNameMode = &v
	return s
}

func (s *NodepoolKubernetesConfig) SetRuntime(v string) *NodepoolKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *NodepoolKubernetesConfig) SetRuntimeVersion(v string) *NodepoolKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *NodepoolKubernetesConfig) SetTaints(v []*Taint) *NodepoolKubernetesConfig {
	s.Taints = v
	return s
}

func (s *NodepoolKubernetesConfig) SetUserData(v string) *NodepoolKubernetesConfig {
	s.UserData = &v
	return s
}

type NodepoolManagement struct {
	AutoRepair    *bool                            `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	Enable        *bool                            `json:"enable,omitempty" xml:"enable,omitempty"`
	UpgradeConfig *NodepoolManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s NodepoolManagement) String() string {
	return tea.Prettify(s)
}

func (s NodepoolManagement) GoString() string {
	return s.String()
}

func (s *NodepoolManagement) SetAutoRepair(v bool) *NodepoolManagement {
	s.AutoRepair = &v
	return s
}

func (s *NodepoolManagement) SetEnable(v bool) *NodepoolManagement {
	s.Enable = &v
	return s
}

func (s *NodepoolManagement) SetUpgradeConfig(v *NodepoolManagementUpgradeConfig) *NodepoolManagement {
	s.UpgradeConfig = v
	return s
}

type NodepoolManagementUpgradeConfig struct {
	AutoUpgrade     *bool  `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	MaxUnavailable  *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
	Surge           *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
}

func (s NodepoolManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s NodepoolManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *NodepoolManagementUpgradeConfig) SetAutoUpgrade(v bool) *NodepoolManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *NodepoolManagementUpgradeConfig) SetMaxUnavailable(v int64) *NodepoolManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

func (s *NodepoolManagementUpgradeConfig) SetSurge(v int64) *NodepoolManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *NodepoolManagementUpgradeConfig) SetSurgePercentage(v int64) *NodepoolManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

type NodepoolNodepoolInfo struct {
	Name            *string `json:"name,omitempty" xml:"name,omitempty"`
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	Type            *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s NodepoolNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s NodepoolNodepoolInfo) GoString() string {
	return s.String()
}

func (s *NodepoolNodepoolInfo) SetName(v string) *NodepoolNodepoolInfo {
	s.Name = &v
	return s
}

func (s *NodepoolNodepoolInfo) SetResourceGroupId(v string) *NodepoolNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

func (s *NodepoolNodepoolInfo) SetType(v string) *NodepoolNodepoolInfo {
	s.Type = &v
	return s
}

type NodepoolScalingGroup struct {
	AutoRenew                           *bool                                   `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	AutoRenewPeriod                     *int64                                  `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	CompensateWithOnDemand              *bool                                   `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	DataDisks                           []*DataDisk                             `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	DeploymentsetId                     *string                                 `json:"deploymentset_id,omitempty" xml:"deploymentset_id,omitempty"`
	DesiredSize                         *int64                                  `json:"desired_size,omitempty" xml:"desired_size,omitempty"`
	ImageId                             *string                                 `json:"image_id,omitempty" xml:"image_id,omitempty"`
	ImageType                           *string                                 `json:"image_type,omitempty" xml:"image_type,omitempty"`
	InstanceChargeType                  *string                                 `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	InstanceTypes                       []*string                               `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	InternetChargeType                  *string                                 `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	InternetMaxBandwidthOut             *int64                                  `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
	KeyPair                             *string                                 `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	LoginPassword                       *string                                 `json:"login_password,omitempty" xml:"login_password,omitempty"`
	MultiAzPolicy                       *string                                 `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	OnDemandBaseCapacity                *int64                                  `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	OnDemandPercentageAboveBaseCapacity *int64                                  `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	Period                              *int64                                  `json:"period,omitempty" xml:"period,omitempty"`
	PeriodUnit                          *string                                 `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	Platform                            *string                                 `json:"platform,omitempty" xml:"platform,omitempty"`
	PrivatePoolOptions                  *NodepoolScalingGroupPrivatePoolOptions `json:"private_pool_options,omitempty" xml:"private_pool_options,omitempty" type:"Struct"`
	RdsInstances                        []*string                               `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	ScalingPolicy                       *string                                 `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	SecurityGroupId                     *string                                 `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	SecurityGroupIds                    []*string                               `json:"security_group_ids,omitempty" xml:"security_group_ids,omitempty" type:"Repeated"`
	SpotInstancePools                   *int64                                  `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	SpotInstanceRemedy                  *bool                                   `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	SpotPriceLimit                      []*NodepoolScalingGroupSpotPriceLimit   `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	SpotStrategy                        *string                                 `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	SystemDiskBurstingEnabled           *bool                                   `json:"system_disk_bursting_enabled,omitempty" xml:"system_disk_bursting_enabled,omitempty"`
	SystemDiskCategory                  *string                                 `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	SystemDiskPerformanceLevel          *string                                 `json:"system_disk_performance_level,omitempty" xml:"system_disk_performance_level,omitempty"`
	SystemDiskProvisionedIops           *int64                                  `json:"system_disk_provisioned_iops,omitempty" xml:"system_disk_provisioned_iops,omitempty"`
	SystemDiskSize                      *int64                                  `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	Tags                                []*NodepoolScalingGroupTags             `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	VswitchIds                          []*string                               `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
}

func (s NodepoolScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s NodepoolScalingGroup) GoString() string {
	return s.String()
}

func (s *NodepoolScalingGroup) SetAutoRenew(v bool) *NodepoolScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *NodepoolScalingGroup) SetAutoRenewPeriod(v int64) *NodepoolScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *NodepoolScalingGroup) SetCompensateWithOnDemand(v bool) *NodepoolScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *NodepoolScalingGroup) SetDataDisks(v []*DataDisk) *NodepoolScalingGroup {
	s.DataDisks = v
	return s
}

func (s *NodepoolScalingGroup) SetDeploymentsetId(v string) *NodepoolScalingGroup {
	s.DeploymentsetId = &v
	return s
}

func (s *NodepoolScalingGroup) SetDesiredSize(v int64) *NodepoolScalingGroup {
	s.DesiredSize = &v
	return s
}

func (s *NodepoolScalingGroup) SetImageId(v string) *NodepoolScalingGroup {
	s.ImageId = &v
	return s
}

func (s *NodepoolScalingGroup) SetImageType(v string) *NodepoolScalingGroup {
	s.ImageType = &v
	return s
}

func (s *NodepoolScalingGroup) SetInstanceChargeType(v string) *NodepoolScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *NodepoolScalingGroup) SetInstanceTypes(v []*string) *NodepoolScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *NodepoolScalingGroup) SetInternetChargeType(v string) *NodepoolScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *NodepoolScalingGroup) SetInternetMaxBandwidthOut(v int64) *NodepoolScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

func (s *NodepoolScalingGroup) SetKeyPair(v string) *NodepoolScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *NodepoolScalingGroup) SetLoginPassword(v string) *NodepoolScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *NodepoolScalingGroup) SetMultiAzPolicy(v string) *NodepoolScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *NodepoolScalingGroup) SetOnDemandBaseCapacity(v int64) *NodepoolScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *NodepoolScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *NodepoolScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *NodepoolScalingGroup) SetPeriod(v int64) *NodepoolScalingGroup {
	s.Period = &v
	return s
}

func (s *NodepoolScalingGroup) SetPeriodUnit(v string) *NodepoolScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *NodepoolScalingGroup) SetPlatform(v string) *NodepoolScalingGroup {
	s.Platform = &v
	return s
}

func (s *NodepoolScalingGroup) SetPrivatePoolOptions(v *NodepoolScalingGroupPrivatePoolOptions) *NodepoolScalingGroup {
	s.PrivatePoolOptions = v
	return s
}

func (s *NodepoolScalingGroup) SetRdsInstances(v []*string) *NodepoolScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *NodepoolScalingGroup) SetScalingPolicy(v string) *NodepoolScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *NodepoolScalingGroup) SetSecurityGroupId(v string) *NodepoolScalingGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *NodepoolScalingGroup) SetSecurityGroupIds(v []*string) *NodepoolScalingGroup {
	s.SecurityGroupIds = v
	return s
}

func (s *NodepoolScalingGroup) SetSpotInstancePools(v int64) *NodepoolScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *NodepoolScalingGroup) SetSpotInstanceRemedy(v bool) *NodepoolScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *NodepoolScalingGroup) SetSpotPriceLimit(v []*NodepoolScalingGroupSpotPriceLimit) *NodepoolScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *NodepoolScalingGroup) SetSpotStrategy(v string) *NodepoolScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *NodepoolScalingGroup) SetSystemDiskBurstingEnabled(v bool) *NodepoolScalingGroup {
	s.SystemDiskBurstingEnabled = &v
	return s
}

func (s *NodepoolScalingGroup) SetSystemDiskCategory(v string) *NodepoolScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *NodepoolScalingGroup) SetSystemDiskPerformanceLevel(v string) *NodepoolScalingGroup {
	s.SystemDiskPerformanceLevel = &v
	return s
}

func (s *NodepoolScalingGroup) SetSystemDiskProvisionedIops(v int64) *NodepoolScalingGroup {
	s.SystemDiskProvisionedIops = &v
	return s
}

func (s *NodepoolScalingGroup) SetSystemDiskSize(v int64) *NodepoolScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *NodepoolScalingGroup) SetTags(v []*NodepoolScalingGroupTags) *NodepoolScalingGroup {
	s.Tags = v
	return s
}

func (s *NodepoolScalingGroup) SetVswitchIds(v []*string) *NodepoolScalingGroup {
	s.VswitchIds = v
	return s
}

type NodepoolScalingGroupPrivatePoolOptions struct {
	Id            *string `json:"id,omitempty" xml:"id,omitempty"`
	MatchCriteria *string `json:"match_criteria,omitempty" xml:"match_criteria,omitempty"`
}

func (s NodepoolScalingGroupPrivatePoolOptions) String() string {
	return tea.Prettify(s)
}

func (s NodepoolScalingGroupPrivatePoolOptions) GoString() string {
	return s.String()
}

func (s *NodepoolScalingGroupPrivatePoolOptions) SetId(v string) *NodepoolScalingGroupPrivatePoolOptions {
	s.Id = &v
	return s
}

func (s *NodepoolScalingGroupPrivatePoolOptions) SetMatchCriteria(v string) *NodepoolScalingGroupPrivatePoolOptions {
	s.MatchCriteria = &v
	return s
}

type NodepoolScalingGroupSpotPriceLimit struct {
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	PriceLimit   *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s NodepoolScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s NodepoolScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *NodepoolScalingGroupSpotPriceLimit) SetInstanceType(v string) *NodepoolScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *NodepoolScalingGroupSpotPriceLimit) SetPriceLimit(v string) *NodepoolScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type NodepoolScalingGroupTags struct {
	Key   *string `json:"key,omitempty" xml:"key,omitempty"`
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s NodepoolScalingGroupTags) String() string {
	return tea.Prettify(s)
}

func (s NodepoolScalingGroupTags) GoString() string {
	return s.String()
}

func (s *NodepoolScalingGroupTags) SetKey(v string) *NodepoolScalingGroupTags {
	s.Key = &v
	return s
}

func (s *NodepoolScalingGroupTags) SetValue(v string) *NodepoolScalingGroupTags {
	s.Value = &v
	return s
}

type NodepoolTeeConfig struct {
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s NodepoolTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s NodepoolTeeConfig) GoString() string {
	return s.String()
}

func (s *NodepoolTeeConfig) SetTeeEnable(v bool) *NodepoolTeeConfig {
	s.TeeEnable = &v
	return s
}

type Runtime struct {
	Name    *string `json:"name,omitempty" xml:"name,omitempty"`
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s Runtime) String() string {
	return tea.Prettify(s)
}

func (s Runtime) GoString() string {
	return s.String()
}

func (s *Runtime) SetName(v string) *Runtime {
	s.Name = &v
	return s
}

func (s *Runtime) SetVersion(v string) *Runtime {
	s.Version = &v
	return s
}

type Tag struct {
	Key   *string `json:"key,omitempty" xml:"key,omitempty"`
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
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

type Taint struct {
	Effect *string `json:"effect,omitempty" xml:"effect,omitempty"`
	Key    *string `json:"key,omitempty" xml:"key,omitempty"`
	Value  *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s Taint) String() string {
	return tea.Prettify(s)
}

func (s Taint) GoString() string {
	return s.String()
}

func (s *Taint) SetEffect(v string) *Taint {
	s.Effect = &v
	return s
}

func (s *Taint) SetKey(v string) *Taint {
	s.Key = &v
	return s
}

func (s *Taint) SetValue(v string) *Taint {
	s.Value = &v
	return s
}

type StandardComponentsValue struct {
	Name        *string `json:"name,omitempty" xml:"name,omitempty"`
	Version     *string `json:"version,omitempty" xml:"version,omitempty"`
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	Required    *string `json:"required,omitempty" xml:"required,omitempty"`
	Disabled    *bool   `json:"disabled,omitempty" xml:"disabled,omitempty"`
}

func (s StandardComponentsValue) String() string {
	return tea.Prettify(s)
}

func (s StandardComponentsValue) GoString() string {
	return s.String()
}

func (s *StandardComponentsValue) SetName(v string) *StandardComponentsValue {
	s.Name = &v
	return s
}

func (s *StandardComponentsValue) SetVersion(v string) *StandardComponentsValue {
	s.Version = &v
	return s
}

func (s *StandardComponentsValue) SetDescription(v string) *StandardComponentsValue {
	s.Description = &v
	return s
}

func (s *StandardComponentsValue) SetRequired(v string) *StandardComponentsValue {
	s.Required = &v
	return s
}

func (s *StandardComponentsValue) SetDisabled(v bool) *StandardComponentsValue {
	s.Disabled = &v
	return s
}

type QuotasValue struct {
	Quota         *string `json:"quota,omitempty" xml:"quota,omitempty"`
	OperationCode *string `json:"operation_code,omitempty" xml:"operation_code,omitempty"`
	Adjustable    *bool   `json:"adjustable,omitempty" xml:"adjustable,omitempty"`
	Unit          *string `json:"unit,omitempty" xml:"unit,omitempty"`
}

func (s QuotasValue) String() string {
	return tea.Prettify(s)
}

func (s QuotasValue) GoString() string {
	return s.String()
}

func (s *QuotasValue) SetQuota(v string) *QuotasValue {
	s.Quota = &v
	return s
}

func (s *QuotasValue) SetOperationCode(v string) *QuotasValue {
	s.OperationCode = &v
	return s
}

func (s *QuotasValue) SetAdjustable(v bool) *QuotasValue {
	s.Adjustable = &v
	return s
}

func (s *QuotasValue) SetUnit(v string) *QuotasValue {
	s.Unit = &v
	return s
}

type AttachInstancesRequest struct {
	CpuPolicy        *string   `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	FormatDisk       *bool     `json:"format_disk,omitempty" xml:"format_disk,omitempty"`
	ImageId          *string   `json:"image_id,omitempty" xml:"image_id,omitempty"`
	Instances        []*string `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
	IsEdgeWorker     *bool     `json:"is_edge_worker,omitempty" xml:"is_edge_worker,omitempty"`
	KeepInstanceName *bool     `json:"keep_instance_name,omitempty" xml:"keep_instance_name,omitempty"`
	KeyPair          *string   `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	NodepoolId       *string   `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	Password         *string   `json:"password,omitempty" xml:"password,omitempty"`
	RdsInstances     []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	Runtime          *Runtime  `json:"runtime,omitempty" xml:"runtime,omitempty"`
	Tags             []*Tag    `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	UserData         *string   `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s AttachInstancesRequest) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesRequest) GoString() string {
	return s.String()
}

func (s *AttachInstancesRequest) SetCpuPolicy(v string) *AttachInstancesRequest {
	s.CpuPolicy = &v
	return s
}

func (s *AttachInstancesRequest) SetFormatDisk(v bool) *AttachInstancesRequest {
	s.FormatDisk = &v
	return s
}

func (s *AttachInstancesRequest) SetImageId(v string) *AttachInstancesRequest {
	s.ImageId = &v
	return s
}

func (s *AttachInstancesRequest) SetInstances(v []*string) *AttachInstancesRequest {
	s.Instances = v
	return s
}

func (s *AttachInstancesRequest) SetIsEdgeWorker(v bool) *AttachInstancesRequest {
	s.IsEdgeWorker = &v
	return s
}

func (s *AttachInstancesRequest) SetKeepInstanceName(v bool) *AttachInstancesRequest {
	s.KeepInstanceName = &v
	return s
}

func (s *AttachInstancesRequest) SetKeyPair(v string) *AttachInstancesRequest {
	s.KeyPair = &v
	return s
}

func (s *AttachInstancesRequest) SetNodepoolId(v string) *AttachInstancesRequest {
	s.NodepoolId = &v
	return s
}

func (s *AttachInstancesRequest) SetPassword(v string) *AttachInstancesRequest {
	s.Password = &v
	return s
}

func (s *AttachInstancesRequest) SetRdsInstances(v []*string) *AttachInstancesRequest {
	s.RdsInstances = v
	return s
}

func (s *AttachInstancesRequest) SetRuntime(v *Runtime) *AttachInstancesRequest {
	s.Runtime = v
	return s
}

func (s *AttachInstancesRequest) SetTags(v []*Tag) *AttachInstancesRequest {
	s.Tags = v
	return s
}

func (s *AttachInstancesRequest) SetUserData(v string) *AttachInstancesRequest {
	s.UserData = &v
	return s
}

type AttachInstancesResponseBody struct {
	// The details of the added nodes.
	List []*AttachInstancesResponseBodyList `json:"list,omitempty" xml:"list,omitempty" type:"Repeated"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s AttachInstancesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesResponseBody) GoString() string {
	return s.String()
}

func (s *AttachInstancesResponseBody) SetList(v []*AttachInstancesResponseBodyList) *AttachInstancesResponseBody {
	s.List = v
	return s
}

func (s *AttachInstancesResponseBody) SetTaskId(v string) *AttachInstancesResponseBody {
	s.TaskId = &v
	return s
}

type AttachInstancesResponseBodyList struct {
	// The code that indicates the task result.
	Code *string `json:"code,omitempty" xml:"code,omitempty"`
	// The ID of the instance.
	InstanceId *string `json:"instanceId,omitempty" xml:"instanceId,omitempty"`
	// Indicates whether the ECS instances are successfully added to the ACK cluster.
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
}

func (s AttachInstancesResponseBodyList) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesResponseBodyList) GoString() string {
	return s.String()
}

func (s *AttachInstancesResponseBodyList) SetCode(v string) *AttachInstancesResponseBodyList {
	s.Code = &v
	return s
}

func (s *AttachInstancesResponseBodyList) SetInstanceId(v string) *AttachInstancesResponseBodyList {
	s.InstanceId = &v
	return s
}

func (s *AttachInstancesResponseBodyList) SetMessage(v string) *AttachInstancesResponseBodyList {
	s.Message = &v
	return s
}

type AttachInstancesResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *AttachInstancesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s AttachInstancesResponse) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesResponse) GoString() string {
	return s.String()
}

func (s *AttachInstancesResponse) SetHeaders(v map[string]*string) *AttachInstancesResponse {
	s.Headers = v
	return s
}

func (s *AttachInstancesResponse) SetStatusCode(v int32) *AttachInstancesResponse {
	s.StatusCode = &v
	return s
}

func (s *AttachInstancesResponse) SetBody(v *AttachInstancesResponseBody) *AttachInstancesResponse {
	s.Body = v
	return s
}

type AttachInstancesToNodePoolRequest struct {
	FormatDisk       *bool     `json:"format_disk,omitempty" xml:"format_disk,omitempty"`
	Instances        []*string `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
	KeepInstanceName *bool     `json:"keep_instance_name,omitempty" xml:"keep_instance_name,omitempty"`
	Password         *string   `json:"password,omitempty" xml:"password,omitempty"`
}

func (s AttachInstancesToNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesToNodePoolRequest) GoString() string {
	return s.String()
}

func (s *AttachInstancesToNodePoolRequest) SetFormatDisk(v bool) *AttachInstancesToNodePoolRequest {
	s.FormatDisk = &v
	return s
}

func (s *AttachInstancesToNodePoolRequest) SetInstances(v []*string) *AttachInstancesToNodePoolRequest {
	s.Instances = v
	return s
}

func (s *AttachInstancesToNodePoolRequest) SetKeepInstanceName(v bool) *AttachInstancesToNodePoolRequest {
	s.KeepInstanceName = &v
	return s
}

func (s *AttachInstancesToNodePoolRequest) SetPassword(v string) *AttachInstancesToNodePoolRequest {
	s.Password = &v
	return s
}

type AttachInstancesToNodePoolResponseBody struct {
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	TaskId    *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s AttachInstancesToNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesToNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *AttachInstancesToNodePoolResponseBody) SetRequestId(v string) *AttachInstancesToNodePoolResponseBody {
	s.RequestId = &v
	return s
}

func (s *AttachInstancesToNodePoolResponseBody) SetTaskId(v string) *AttachInstancesToNodePoolResponseBody {
	s.TaskId = &v
	return s
}

type AttachInstancesToNodePoolResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *AttachInstancesToNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s AttachInstancesToNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s AttachInstancesToNodePoolResponse) GoString() string {
	return s.String()
}

func (s *AttachInstancesToNodePoolResponse) SetHeaders(v map[string]*string) *AttachInstancesToNodePoolResponse {
	s.Headers = v
	return s
}

func (s *AttachInstancesToNodePoolResponse) SetStatusCode(v int32) *AttachInstancesToNodePoolResponse {
	s.StatusCode = &v
	return s
}

func (s *AttachInstancesToNodePoolResponse) SetBody(v *AttachInstancesToNodePoolResponseBody) *AttachInstancesToNodePoolResponse {
	s.Body = v
	return s
}

type CancelClusterUpgradeResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s CancelClusterUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelClusterUpgradeResponse) GoString() string {
	return s.String()
}

func (s *CancelClusterUpgradeResponse) SetHeaders(v map[string]*string) *CancelClusterUpgradeResponse {
	s.Headers = v
	return s
}

func (s *CancelClusterUpgradeResponse) SetStatusCode(v int32) *CancelClusterUpgradeResponse {
	s.StatusCode = &v
	return s
}

type CancelComponentUpgradeResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s CancelComponentUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelComponentUpgradeResponse) GoString() string {
	return s.String()
}

func (s *CancelComponentUpgradeResponse) SetHeaders(v map[string]*string) *CancelComponentUpgradeResponse {
	s.Headers = v
	return s
}

func (s *CancelComponentUpgradeResponse) SetStatusCode(v int32) *CancelComponentUpgradeResponse {
	s.StatusCode = &v
	return s
}

type CancelTaskResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s CancelTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelTaskResponse) GoString() string {
	return s.String()
}

func (s *CancelTaskResponse) SetHeaders(v map[string]*string) *CancelTaskResponse {
	s.Headers = v
	return s
}

func (s *CancelTaskResponse) SetStatusCode(v int32) *CancelTaskResponse {
	s.StatusCode = &v
	return s
}

type CancelWorkflowRequest struct {
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
}

func (s CancelWorkflowRequest) String() string {
	return tea.Prettify(s)
}

func (s CancelWorkflowRequest) GoString() string {
	return s.String()
}

func (s *CancelWorkflowRequest) SetAction(v string) *CancelWorkflowRequest {
	s.Action = &v
	return s
}

type CancelWorkflowResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s CancelWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s CancelWorkflowResponse) GoString() string {
	return s.String()
}

func (s *CancelWorkflowResponse) SetHeaders(v map[string]*string) *CancelWorkflowResponse {
	s.Headers = v
	return s
}

func (s *CancelWorkflowResponse) SetStatusCode(v int32) *CancelWorkflowResponse {
	s.StatusCode = &v
	return s
}

type CheckControlPlaneLogEnableResponseBody struct {
	Aliuid     *string   `json:"aliuid,omitempty" xml:"aliuid,omitempty"`
	Components []*string `json:"components,omitempty" xml:"components,omitempty" type:"Repeated"`
	LogProject *string   `json:"log_project,omitempty" xml:"log_project,omitempty"`
	LogTtl     *string   `json:"log_ttl,omitempty" xml:"log_ttl,omitempty"`
}

func (s CheckControlPlaneLogEnableResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CheckControlPlaneLogEnableResponseBody) GoString() string {
	return s.String()
}

func (s *CheckControlPlaneLogEnableResponseBody) SetAliuid(v string) *CheckControlPlaneLogEnableResponseBody {
	s.Aliuid = &v
	return s
}

func (s *CheckControlPlaneLogEnableResponseBody) SetComponents(v []*string) *CheckControlPlaneLogEnableResponseBody {
	s.Components = v
	return s
}

func (s *CheckControlPlaneLogEnableResponseBody) SetLogProject(v string) *CheckControlPlaneLogEnableResponseBody {
	s.LogProject = &v
	return s
}

func (s *CheckControlPlaneLogEnableResponseBody) SetLogTtl(v string) *CheckControlPlaneLogEnableResponseBody {
	s.LogTtl = &v
	return s
}

type CheckControlPlaneLogEnableResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CheckControlPlaneLogEnableResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CheckControlPlaneLogEnableResponse) String() string {
	return tea.Prettify(s)
}

func (s CheckControlPlaneLogEnableResponse) GoString() string {
	return s.String()
}

func (s *CheckControlPlaneLogEnableResponse) SetHeaders(v map[string]*string) *CheckControlPlaneLogEnableResponse {
	s.Headers = v
	return s
}

func (s *CheckControlPlaneLogEnableResponse) SetStatusCode(v int32) *CheckControlPlaneLogEnableResponse {
	s.StatusCode = &v
	return s
}

func (s *CheckControlPlaneLogEnableResponse) SetBody(v *CheckControlPlaneLogEnableResponseBody) *CheckControlPlaneLogEnableResponse {
	s.Body = v
	return s
}

type CreateAutoscalingConfigRequest struct {
	CoolDownDuration           *string `json:"cool_down_duration,omitempty" xml:"cool_down_duration,omitempty"`
	DaemonsetEvictionForNodes  *bool   `json:"daemonset_eviction_for_nodes,omitempty" xml:"daemonset_eviction_for_nodes,omitempty"`
	Expander                   *string `json:"expander,omitempty" xml:"expander,omitempty"`
	GpuUtilizationThreshold    *string `json:"gpu_utilization_threshold,omitempty" xml:"gpu_utilization_threshold,omitempty"`
	MaxGracefulTerminationSec  *int32  `json:"max_graceful_termination_sec,omitempty" xml:"max_graceful_termination_sec,omitempty"`
	MinReplicaCount            *int32  `json:"min_replica_count,omitempty" xml:"min_replica_count,omitempty"`
	RecycleNodeDeletionEnabled *bool   `json:"recycle_node_deletion_enabled,omitempty" xml:"recycle_node_deletion_enabled,omitempty"`
	ScaleDownEnabled           *bool   `json:"scale_down_enabled,omitempty" xml:"scale_down_enabled,omitempty"`
	ScaleUpFromZero            *bool   `json:"scale_up_from_zero,omitempty" xml:"scale_up_from_zero,omitempty"`
	ScanInterval               *string `json:"scan_interval,omitempty" xml:"scan_interval,omitempty"`
	SkipNodesWithLocalStorage  *bool   `json:"skip_nodes_with_local_storage,omitempty" xml:"skip_nodes_with_local_storage,omitempty"`
	SkipNodesWithSystemPods    *bool   `json:"skip_nodes_with_system_pods,omitempty" xml:"skip_nodes_with_system_pods,omitempty"`
	UnneededDuration           *string `json:"unneeded_duration,omitempty" xml:"unneeded_duration,omitempty"`
	UtilizationThreshold       *string `json:"utilization_threshold,omitempty" xml:"utilization_threshold,omitempty"`
}

func (s CreateAutoscalingConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateAutoscalingConfigRequest) GoString() string {
	return s.String()
}

func (s *CreateAutoscalingConfigRequest) SetCoolDownDuration(v string) *CreateAutoscalingConfigRequest {
	s.CoolDownDuration = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetDaemonsetEvictionForNodes(v bool) *CreateAutoscalingConfigRequest {
	s.DaemonsetEvictionForNodes = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetExpander(v string) *CreateAutoscalingConfigRequest {
	s.Expander = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetGpuUtilizationThreshold(v string) *CreateAutoscalingConfigRequest {
	s.GpuUtilizationThreshold = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetMaxGracefulTerminationSec(v int32) *CreateAutoscalingConfigRequest {
	s.MaxGracefulTerminationSec = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetMinReplicaCount(v int32) *CreateAutoscalingConfigRequest {
	s.MinReplicaCount = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetRecycleNodeDeletionEnabled(v bool) *CreateAutoscalingConfigRequest {
	s.RecycleNodeDeletionEnabled = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetScaleDownEnabled(v bool) *CreateAutoscalingConfigRequest {
	s.ScaleDownEnabled = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetScaleUpFromZero(v bool) *CreateAutoscalingConfigRequest {
	s.ScaleUpFromZero = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetScanInterval(v string) *CreateAutoscalingConfigRequest {
	s.ScanInterval = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetSkipNodesWithLocalStorage(v bool) *CreateAutoscalingConfigRequest {
	s.SkipNodesWithLocalStorage = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetSkipNodesWithSystemPods(v bool) *CreateAutoscalingConfigRequest {
	s.SkipNodesWithSystemPods = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetUnneededDuration(v string) *CreateAutoscalingConfigRequest {
	s.UnneededDuration = &v
	return s
}

func (s *CreateAutoscalingConfigRequest) SetUtilizationThreshold(v string) *CreateAutoscalingConfigRequest {
	s.UtilizationThreshold = &v
	return s
}

type CreateAutoscalingConfigResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s CreateAutoscalingConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateAutoscalingConfigResponse) GoString() string {
	return s.String()
}

func (s *CreateAutoscalingConfigResponse) SetHeaders(v map[string]*string) *CreateAutoscalingConfigResponse {
	s.Headers = v
	return s
}

func (s *CreateAutoscalingConfigResponse) SetStatusCode(v int32) *CreateAutoscalingConfigResponse {
	s.StatusCode = &v
	return s
}

type CreateClusterRequest struct {
	Addons                           []*Addon    `json:"addons,omitempty" xml:"addons,omitempty" type:"Repeated"`
	ApiAudiences                     *string     `json:"api_audiences,omitempty" xml:"api_audiences,omitempty"`
	ChargeType                       *string     `json:"charge_type,omitempty" xml:"charge_type,omitempty"`
	CisEnabled                       *bool       `json:"cis_enabled,omitempty" xml:"cis_enabled,omitempty"`
	CloudMonitorFlags                *bool       `json:"cloud_monitor_flags,omitempty" xml:"cloud_monitor_flags,omitempty"`
	ClusterDomain                    *string     `json:"cluster_domain,omitempty" xml:"cluster_domain,omitempty"`
	ClusterSpec                      *string     `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	ClusterType                      *string     `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	ContainerCidr                    *string     `json:"container_cidr,omitempty" xml:"container_cidr,omitempty"`
	ControlplaneLogComponents        []*string   `json:"controlplane_log_components,omitempty" xml:"controlplane_log_components,omitempty" type:"Repeated"`
	ControlplaneLogProject           *string     `json:"controlplane_log_project,omitempty" xml:"controlplane_log_project,omitempty"`
	ControlplaneLogTtl               *string     `json:"controlplane_log_ttl,omitempty" xml:"controlplane_log_ttl,omitempty"`
	CpuPolicy                        *string     `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	CustomSan                        *string     `json:"custom_san,omitempty" xml:"custom_san,omitempty"`
	DeletionProtection               *bool       `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	DisableRollback                  *bool       `json:"disable_rollback,omitempty" xml:"disable_rollback,omitempty"`
	EnableRrsa                       *bool       `json:"enable_rrsa,omitempty" xml:"enable_rrsa,omitempty"`
	EncryptionProviderKey            *string     `json:"encryption_provider_key,omitempty" xml:"encryption_provider_key,omitempty"`
	EndpointPublicAccess             *bool       `json:"endpoint_public_access,omitempty" xml:"endpoint_public_access,omitempty"`
	FormatDisk                       *bool       `json:"format_disk,omitempty" xml:"format_disk,omitempty"`
	ImageId                          *string     `json:"image_id,omitempty" xml:"image_id,omitempty"`
	ImageType                        *string     `json:"image_type,omitempty" xml:"image_type,omitempty"`
	Instances                        []*string   `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
	IpStack                          *string     `json:"ip_stack,omitempty" xml:"ip_stack,omitempty"`
	IsEnterpriseSecurityGroup        *bool       `json:"is_enterprise_security_group,omitempty" xml:"is_enterprise_security_group,omitempty"`
	KeepInstanceName                 *bool       `json:"keep_instance_name,omitempty" xml:"keep_instance_name,omitempty"`
	KeyPair                          *string     `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	KubernetesVersion                *string     `json:"kubernetes_version,omitempty" xml:"kubernetes_version,omitempty"`
	LoadBalancerSpec                 *string     `json:"load_balancer_spec,omitempty" xml:"load_balancer_spec,omitempty"`
	LoggingType                      *string     `json:"logging_type,omitempty" xml:"logging_type,omitempty"`
	LoginPassword                    *string     `json:"login_password,omitempty" xml:"login_password,omitempty"`
	MasterAutoRenew                  *bool       `json:"master_auto_renew,omitempty" xml:"master_auto_renew,omitempty"`
	MasterAutoRenewPeriod            *int64      `json:"master_auto_renew_period,omitempty" xml:"master_auto_renew_period,omitempty"`
	MasterCount                      *int64      `json:"master_count,omitempty" xml:"master_count,omitempty"`
	MasterInstanceChargeType         *string     `json:"master_instance_charge_type,omitempty" xml:"master_instance_charge_type,omitempty"`
	MasterInstanceTypes              []*string   `json:"master_instance_types,omitempty" xml:"master_instance_types,omitempty" type:"Repeated"`
	MasterPeriod                     *int64      `json:"master_period,omitempty" xml:"master_period,omitempty"`
	MasterPeriodUnit                 *string     `json:"master_period_unit,omitempty" xml:"master_period_unit,omitempty"`
	MasterSystemDiskCategory         *string     `json:"master_system_disk_category,omitempty" xml:"master_system_disk_category,omitempty"`
	MasterSystemDiskPerformanceLevel *string     `json:"master_system_disk_performance_level,omitempty" xml:"master_system_disk_performance_level,omitempty"`
	MasterSystemDiskSize             *int64      `json:"master_system_disk_size,omitempty" xml:"master_system_disk_size,omitempty"`
	MasterSystemDiskSnapshotPolicyId *string     `json:"master_system_disk_snapshot_policy_id,omitempty" xml:"master_system_disk_snapshot_policy_id,omitempty"`
	MasterVswitchIds                 []*string   `json:"master_vswitch_ids,omitempty" xml:"master_vswitch_ids,omitempty" type:"Repeated"`
	Name                             *string     `json:"name,omitempty" xml:"name,omitempty"`
	NatGateway                       *bool       `json:"nat_gateway,omitempty" xml:"nat_gateway,omitempty"`
	NodeCidrMask                     *string     `json:"node_cidr_mask,omitempty" xml:"node_cidr_mask,omitempty"`
	NodeNameMode                     *string     `json:"node_name_mode,omitempty" xml:"node_name_mode,omitempty"`
	NodePortRange                    *string     `json:"node_port_range,omitempty" xml:"node_port_range,omitempty"`
	Nodepools                        []*Nodepool `json:"nodepools,omitempty" xml:"nodepools,omitempty" type:"Repeated"`
	// Deprecated
	NumOfNodes            *int64    `json:"num_of_nodes,omitempty" xml:"num_of_nodes,omitempty"`
	OsType                *string   `json:"os_type,omitempty" xml:"os_type,omitempty"`
	Period                *int64    `json:"period,omitempty" xml:"period,omitempty"`
	PeriodUnit            *string   `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	Platform              *string   `json:"platform,omitempty" xml:"platform,omitempty"`
	PodVswitchIds         []*string `json:"pod_vswitch_ids,omitempty" xml:"pod_vswitch_ids,omitempty" type:"Repeated"`
	Profile               *string   `json:"profile,omitempty" xml:"profile,omitempty"`
	ProxyMode             *string   `json:"proxy_mode,omitempty" xml:"proxy_mode,omitempty"`
	RdsInstances          []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	RegionId              *string   `json:"region_id,omitempty" xml:"region_id,omitempty"`
	ResourceGroupId       *string   `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	Runtime               *Runtime  `json:"runtime,omitempty" xml:"runtime,omitempty"`
	SecurityGroupId       *string   `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	ServiceAccountIssuer  *string   `json:"service_account_issuer,omitempty" xml:"service_account_issuer,omitempty"`
	ServiceCidr           *string   `json:"service_cidr,omitempty" xml:"service_cidr,omitempty"`
	ServiceDiscoveryTypes []*string `json:"service_discovery_types,omitempty" xml:"service_discovery_types,omitempty" type:"Repeated"`
	SnatEntry             *bool     `json:"snat_entry,omitempty" xml:"snat_entry,omitempty"`
	SocEnabled            *bool     `json:"soc_enabled,omitempty" xml:"soc_enabled,omitempty"`
	SshFlags              *bool     `json:"ssh_flags,omitempty" xml:"ssh_flags,omitempty"`
	Tags                  []*Tag    `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	Taints                []*Taint  `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	TimeoutMins           *int64    `json:"timeout_mins,omitempty" xml:"timeout_mins,omitempty"`
	Timezone              *string   `json:"timezone,omitempty" xml:"timezone,omitempty"`
	UserCa                *string   `json:"user_ca,omitempty" xml:"user_ca,omitempty"`
	UserData              *string   `json:"user_data,omitempty" xml:"user_data,omitempty"`
	Vpcid                 *string   `json:"vpcid,omitempty" xml:"vpcid,omitempty"`
	VswitchIds            []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	// Deprecated
	WorkerAutoRenew *bool `json:"worker_auto_renew,omitempty" xml:"worker_auto_renew,omitempty"`
	// Deprecated
	WorkerAutoRenewPeriod *int64 `json:"worker_auto_renew_period,omitempty" xml:"worker_auto_renew_period,omitempty"`
	// Deprecated
	WorkerDataDisks []*CreateClusterRequestWorkerDataDisks `json:"worker_data_disks,omitempty" xml:"worker_data_disks,omitempty" type:"Repeated"`
	// Deprecated
	WorkerInstanceChargeType *string `json:"worker_instance_charge_type,omitempty" xml:"worker_instance_charge_type,omitempty"`
	// Deprecated
	WorkerInstanceTypes []*string `json:"worker_instance_types,omitempty" xml:"worker_instance_types,omitempty" type:"Repeated"`
	// Deprecated
	WorkerPeriod *int64 `json:"worker_period,omitempty" xml:"worker_period,omitempty"`
	// Deprecated
	WorkerPeriodUnit *string `json:"worker_period_unit,omitempty" xml:"worker_period_unit,omitempty"`
	// Deprecated
	WorkerSystemDiskCategory *string `json:"worker_system_disk_category,omitempty" xml:"worker_system_disk_category,omitempty"`
	// Deprecated
	WorkerSystemDiskPerformanceLevel *string `json:"worker_system_disk_performance_level,omitempty" xml:"worker_system_disk_performance_level,omitempty"`
	// Deprecated
	WorkerSystemDiskSize *int64 `json:"worker_system_disk_size,omitempty" xml:"worker_system_disk_size,omitempty"`
	// Deprecated
	WorkerSystemDiskSnapshotPolicyId *string `json:"worker_system_disk_snapshot_policy_id,omitempty" xml:"worker_system_disk_snapshot_policy_id,omitempty"`
	// Deprecated
	WorkerVswitchIds []*string `json:"worker_vswitch_ids,omitempty" xml:"worker_vswitch_ids,omitempty" type:"Repeated"`
	ZoneId           *string   `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
}

func (s CreateClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterRequest) GoString() string {
	return s.String()
}

func (s *CreateClusterRequest) SetAddons(v []*Addon) *CreateClusterRequest {
	s.Addons = v
	return s
}

func (s *CreateClusterRequest) SetApiAudiences(v string) *CreateClusterRequest {
	s.ApiAudiences = &v
	return s
}

func (s *CreateClusterRequest) SetChargeType(v string) *CreateClusterRequest {
	s.ChargeType = &v
	return s
}

func (s *CreateClusterRequest) SetCisEnabled(v bool) *CreateClusterRequest {
	s.CisEnabled = &v
	return s
}

func (s *CreateClusterRequest) SetCloudMonitorFlags(v bool) *CreateClusterRequest {
	s.CloudMonitorFlags = &v
	return s
}

func (s *CreateClusterRequest) SetClusterDomain(v string) *CreateClusterRequest {
	s.ClusterDomain = &v
	return s
}

func (s *CreateClusterRequest) SetClusterSpec(v string) *CreateClusterRequest {
	s.ClusterSpec = &v
	return s
}

func (s *CreateClusterRequest) SetClusterType(v string) *CreateClusterRequest {
	s.ClusterType = &v
	return s
}

func (s *CreateClusterRequest) SetContainerCidr(v string) *CreateClusterRequest {
	s.ContainerCidr = &v
	return s
}

func (s *CreateClusterRequest) SetControlplaneLogComponents(v []*string) *CreateClusterRequest {
	s.ControlplaneLogComponents = v
	return s
}

func (s *CreateClusterRequest) SetControlplaneLogProject(v string) *CreateClusterRequest {
	s.ControlplaneLogProject = &v
	return s
}

func (s *CreateClusterRequest) SetControlplaneLogTtl(v string) *CreateClusterRequest {
	s.ControlplaneLogTtl = &v
	return s
}

func (s *CreateClusterRequest) SetCpuPolicy(v string) *CreateClusterRequest {
	s.CpuPolicy = &v
	return s
}

func (s *CreateClusterRequest) SetCustomSan(v string) *CreateClusterRequest {
	s.CustomSan = &v
	return s
}

func (s *CreateClusterRequest) SetDeletionProtection(v bool) *CreateClusterRequest {
	s.DeletionProtection = &v
	return s
}

func (s *CreateClusterRequest) SetDisableRollback(v bool) *CreateClusterRequest {
	s.DisableRollback = &v
	return s
}

func (s *CreateClusterRequest) SetEnableRrsa(v bool) *CreateClusterRequest {
	s.EnableRrsa = &v
	return s
}

func (s *CreateClusterRequest) SetEncryptionProviderKey(v string) *CreateClusterRequest {
	s.EncryptionProviderKey = &v
	return s
}

func (s *CreateClusterRequest) SetEndpointPublicAccess(v bool) *CreateClusterRequest {
	s.EndpointPublicAccess = &v
	return s
}

func (s *CreateClusterRequest) SetFormatDisk(v bool) *CreateClusterRequest {
	s.FormatDisk = &v
	return s
}

func (s *CreateClusterRequest) SetImageId(v string) *CreateClusterRequest {
	s.ImageId = &v
	return s
}

func (s *CreateClusterRequest) SetImageType(v string) *CreateClusterRequest {
	s.ImageType = &v
	return s
}

func (s *CreateClusterRequest) SetInstances(v []*string) *CreateClusterRequest {
	s.Instances = v
	return s
}

func (s *CreateClusterRequest) SetIpStack(v string) *CreateClusterRequest {
	s.IpStack = &v
	return s
}

func (s *CreateClusterRequest) SetIsEnterpriseSecurityGroup(v bool) *CreateClusterRequest {
	s.IsEnterpriseSecurityGroup = &v
	return s
}

func (s *CreateClusterRequest) SetKeepInstanceName(v bool) *CreateClusterRequest {
	s.KeepInstanceName = &v
	return s
}

func (s *CreateClusterRequest) SetKeyPair(v string) *CreateClusterRequest {
	s.KeyPair = &v
	return s
}

func (s *CreateClusterRequest) SetKubernetesVersion(v string) *CreateClusterRequest {
	s.KubernetesVersion = &v
	return s
}

func (s *CreateClusterRequest) SetLoadBalancerSpec(v string) *CreateClusterRequest {
	s.LoadBalancerSpec = &v
	return s
}

func (s *CreateClusterRequest) SetLoggingType(v string) *CreateClusterRequest {
	s.LoggingType = &v
	return s
}

func (s *CreateClusterRequest) SetLoginPassword(v string) *CreateClusterRequest {
	s.LoginPassword = &v
	return s
}

func (s *CreateClusterRequest) SetMasterAutoRenew(v bool) *CreateClusterRequest {
	s.MasterAutoRenew = &v
	return s
}

func (s *CreateClusterRequest) SetMasterAutoRenewPeriod(v int64) *CreateClusterRequest {
	s.MasterAutoRenewPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetMasterCount(v int64) *CreateClusterRequest {
	s.MasterCount = &v
	return s
}

func (s *CreateClusterRequest) SetMasterInstanceChargeType(v string) *CreateClusterRequest {
	s.MasterInstanceChargeType = &v
	return s
}

func (s *CreateClusterRequest) SetMasterInstanceTypes(v []*string) *CreateClusterRequest {
	s.MasterInstanceTypes = v
	return s
}

func (s *CreateClusterRequest) SetMasterPeriod(v int64) *CreateClusterRequest {
	s.MasterPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetMasterPeriodUnit(v string) *CreateClusterRequest {
	s.MasterPeriodUnit = &v
	return s
}

func (s *CreateClusterRequest) SetMasterSystemDiskCategory(v string) *CreateClusterRequest {
	s.MasterSystemDiskCategory = &v
	return s
}

func (s *CreateClusterRequest) SetMasterSystemDiskPerformanceLevel(v string) *CreateClusterRequest {
	s.MasterSystemDiskPerformanceLevel = &v
	return s
}

func (s *CreateClusterRequest) SetMasterSystemDiskSize(v int64) *CreateClusterRequest {
	s.MasterSystemDiskSize = &v
	return s
}

func (s *CreateClusterRequest) SetMasterSystemDiskSnapshotPolicyId(v string) *CreateClusterRequest {
	s.MasterSystemDiskSnapshotPolicyId = &v
	return s
}

func (s *CreateClusterRequest) SetMasterVswitchIds(v []*string) *CreateClusterRequest {
	s.MasterVswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetName(v string) *CreateClusterRequest {
	s.Name = &v
	return s
}

func (s *CreateClusterRequest) SetNatGateway(v bool) *CreateClusterRequest {
	s.NatGateway = &v
	return s
}

func (s *CreateClusterRequest) SetNodeCidrMask(v string) *CreateClusterRequest {
	s.NodeCidrMask = &v
	return s
}

func (s *CreateClusterRequest) SetNodeNameMode(v string) *CreateClusterRequest {
	s.NodeNameMode = &v
	return s
}

func (s *CreateClusterRequest) SetNodePortRange(v string) *CreateClusterRequest {
	s.NodePortRange = &v
	return s
}

func (s *CreateClusterRequest) SetNodepools(v []*Nodepool) *CreateClusterRequest {
	s.Nodepools = v
	return s
}

func (s *CreateClusterRequest) SetNumOfNodes(v int64) *CreateClusterRequest {
	s.NumOfNodes = &v
	return s
}

func (s *CreateClusterRequest) SetOsType(v string) *CreateClusterRequest {
	s.OsType = &v
	return s
}

func (s *CreateClusterRequest) SetPeriod(v int64) *CreateClusterRequest {
	s.Period = &v
	return s
}

func (s *CreateClusterRequest) SetPeriodUnit(v string) *CreateClusterRequest {
	s.PeriodUnit = &v
	return s
}

func (s *CreateClusterRequest) SetPlatform(v string) *CreateClusterRequest {
	s.Platform = &v
	return s
}

func (s *CreateClusterRequest) SetPodVswitchIds(v []*string) *CreateClusterRequest {
	s.PodVswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetProfile(v string) *CreateClusterRequest {
	s.Profile = &v
	return s
}

func (s *CreateClusterRequest) SetProxyMode(v string) *CreateClusterRequest {
	s.ProxyMode = &v
	return s
}

func (s *CreateClusterRequest) SetRdsInstances(v []*string) *CreateClusterRequest {
	s.RdsInstances = v
	return s
}

func (s *CreateClusterRequest) SetRegionId(v string) *CreateClusterRequest {
	s.RegionId = &v
	return s
}

func (s *CreateClusterRequest) SetResourceGroupId(v string) *CreateClusterRequest {
	s.ResourceGroupId = &v
	return s
}

func (s *CreateClusterRequest) SetRuntime(v *Runtime) *CreateClusterRequest {
	s.Runtime = v
	return s
}

func (s *CreateClusterRequest) SetSecurityGroupId(v string) *CreateClusterRequest {
	s.SecurityGroupId = &v
	return s
}

func (s *CreateClusterRequest) SetServiceAccountIssuer(v string) *CreateClusterRequest {
	s.ServiceAccountIssuer = &v
	return s
}

func (s *CreateClusterRequest) SetServiceCidr(v string) *CreateClusterRequest {
	s.ServiceCidr = &v
	return s
}

func (s *CreateClusterRequest) SetServiceDiscoveryTypes(v []*string) *CreateClusterRequest {
	s.ServiceDiscoveryTypes = v
	return s
}

func (s *CreateClusterRequest) SetSnatEntry(v bool) *CreateClusterRequest {
	s.SnatEntry = &v
	return s
}

func (s *CreateClusterRequest) SetSocEnabled(v bool) *CreateClusterRequest {
	s.SocEnabled = &v
	return s
}

func (s *CreateClusterRequest) SetSshFlags(v bool) *CreateClusterRequest {
	s.SshFlags = &v
	return s
}

func (s *CreateClusterRequest) SetTags(v []*Tag) *CreateClusterRequest {
	s.Tags = v
	return s
}

func (s *CreateClusterRequest) SetTaints(v []*Taint) *CreateClusterRequest {
	s.Taints = v
	return s
}

func (s *CreateClusterRequest) SetTimeoutMins(v int64) *CreateClusterRequest {
	s.TimeoutMins = &v
	return s
}

func (s *CreateClusterRequest) SetTimezone(v string) *CreateClusterRequest {
	s.Timezone = &v
	return s
}

func (s *CreateClusterRequest) SetUserCa(v string) *CreateClusterRequest {
	s.UserCa = &v
	return s
}

func (s *CreateClusterRequest) SetUserData(v string) *CreateClusterRequest {
	s.UserData = &v
	return s
}

func (s *CreateClusterRequest) SetVpcid(v string) *CreateClusterRequest {
	s.Vpcid = &v
	return s
}

func (s *CreateClusterRequest) SetVswitchIds(v []*string) *CreateClusterRequest {
	s.VswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetWorkerAutoRenew(v bool) *CreateClusterRequest {
	s.WorkerAutoRenew = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerAutoRenewPeriod(v int64) *CreateClusterRequest {
	s.WorkerAutoRenewPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerDataDisks(v []*CreateClusterRequestWorkerDataDisks) *CreateClusterRequest {
	s.WorkerDataDisks = v
	return s
}

func (s *CreateClusterRequest) SetWorkerInstanceChargeType(v string) *CreateClusterRequest {
	s.WorkerInstanceChargeType = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerInstanceTypes(v []*string) *CreateClusterRequest {
	s.WorkerInstanceTypes = v
	return s
}

func (s *CreateClusterRequest) SetWorkerPeriod(v int64) *CreateClusterRequest {
	s.WorkerPeriod = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerPeriodUnit(v string) *CreateClusterRequest {
	s.WorkerPeriodUnit = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerSystemDiskCategory(v string) *CreateClusterRequest {
	s.WorkerSystemDiskCategory = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerSystemDiskPerformanceLevel(v string) *CreateClusterRequest {
	s.WorkerSystemDiskPerformanceLevel = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerSystemDiskSize(v int64) *CreateClusterRequest {
	s.WorkerSystemDiskSize = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerSystemDiskSnapshotPolicyId(v string) *CreateClusterRequest {
	s.WorkerSystemDiskSnapshotPolicyId = &v
	return s
}

func (s *CreateClusterRequest) SetWorkerVswitchIds(v []*string) *CreateClusterRequest {
	s.WorkerVswitchIds = v
	return s
}

func (s *CreateClusterRequest) SetZoneId(v string) *CreateClusterRequest {
	s.ZoneId = &v
	return s
}

type CreateClusterRequestWorkerDataDisks struct {
	Category         *string `json:"category,omitempty" xml:"category,omitempty"`
	Encrypted        *string `json:"encrypted,omitempty" xml:"encrypted,omitempty"`
	PerformanceLevel *string `json:"performance_level,omitempty" xml:"performance_level,omitempty"`
	Size             *string `json:"size,omitempty" xml:"size,omitempty"`
}

func (s CreateClusterRequestWorkerDataDisks) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterRequestWorkerDataDisks) GoString() string {
	return s.String()
}

func (s *CreateClusterRequestWorkerDataDisks) SetCategory(v string) *CreateClusterRequestWorkerDataDisks {
	s.Category = &v
	return s
}

func (s *CreateClusterRequestWorkerDataDisks) SetEncrypted(v string) *CreateClusterRequestWorkerDataDisks {
	s.Encrypted = &v
	return s
}

func (s *CreateClusterRequestWorkerDataDisks) SetPerformanceLevel(v string) *CreateClusterRequestWorkerDataDisks {
	s.PerformanceLevel = &v
	return s
}

func (s *CreateClusterRequestWorkerDataDisks) SetSize(v string) *CreateClusterRequestWorkerDataDisks {
	s.Size = &v
	return s
}

type CreateClusterResponseBody struct {
	// The ID of the cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s CreateClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterResponseBody) GoString() string {
	return s.String()
}

func (s *CreateClusterResponseBody) SetClusterId(v string) *CreateClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *CreateClusterResponseBody) SetRequestId(v string) *CreateClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *CreateClusterResponseBody) SetTaskId(v string) *CreateClusterResponseBody {
	s.TaskId = &v
	return s
}

type CreateClusterResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterResponse) GoString() string {
	return s.String()
}

func (s *CreateClusterResponse) SetHeaders(v map[string]*string) *CreateClusterResponse {
	s.Headers = v
	return s
}

func (s *CreateClusterResponse) SetStatusCode(v int32) *CreateClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateClusterResponse) SetBody(v *CreateClusterResponseBody) *CreateClusterResponse {
	s.Body = v
	return s
}

type CreateClusterNodePoolRequest struct {
	AutoScaling        *CreateClusterNodePoolRequestAutoScaling        `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	Count              *int64                                          `json:"count,omitempty" xml:"count,omitempty"`
	InterconnectConfig *CreateClusterNodePoolRequestInterconnectConfig `json:"interconnect_config,omitempty" xml:"interconnect_config,omitempty" type:"Struct"`
	InterconnectMode   *string                                         `json:"interconnect_mode,omitempty" xml:"interconnect_mode,omitempty"`
	KubernetesConfig   *CreateClusterNodePoolRequestKubernetesConfig   `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	Management         *CreateClusterNodePoolRequestManagement         `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
	MaxNodes           *int64                                          `json:"max_nodes,omitempty" xml:"max_nodes,omitempty"`
	NodepoolInfo       *CreateClusterNodePoolRequestNodepoolInfo       `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	ScalingGroup       *CreateClusterNodePoolRequestScalingGroup       `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	TeeConfig          *CreateClusterNodePoolRequestTeeConfig          `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
}

func (s CreateClusterNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequest) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequest) SetAutoScaling(v *CreateClusterNodePoolRequestAutoScaling) *CreateClusterNodePoolRequest {
	s.AutoScaling = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetCount(v int64) *CreateClusterNodePoolRequest {
	s.Count = &v
	return s
}

func (s *CreateClusterNodePoolRequest) SetInterconnectConfig(v *CreateClusterNodePoolRequestInterconnectConfig) *CreateClusterNodePoolRequest {
	s.InterconnectConfig = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetInterconnectMode(v string) *CreateClusterNodePoolRequest {
	s.InterconnectMode = &v
	return s
}

func (s *CreateClusterNodePoolRequest) SetKubernetesConfig(v *CreateClusterNodePoolRequestKubernetesConfig) *CreateClusterNodePoolRequest {
	s.KubernetesConfig = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetManagement(v *CreateClusterNodePoolRequestManagement) *CreateClusterNodePoolRequest {
	s.Management = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetMaxNodes(v int64) *CreateClusterNodePoolRequest {
	s.MaxNodes = &v
	return s
}

func (s *CreateClusterNodePoolRequest) SetNodepoolInfo(v *CreateClusterNodePoolRequestNodepoolInfo) *CreateClusterNodePoolRequest {
	s.NodepoolInfo = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetScalingGroup(v *CreateClusterNodePoolRequestScalingGroup) *CreateClusterNodePoolRequest {
	s.ScalingGroup = v
	return s
}

func (s *CreateClusterNodePoolRequest) SetTeeConfig(v *CreateClusterNodePoolRequestTeeConfig) *CreateClusterNodePoolRequest {
	s.TeeConfig = v
	return s
}

type CreateClusterNodePoolRequestAutoScaling struct {
	EipBandwidth          *int64  `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	Enable                *bool   `json:"enable,omitempty" xml:"enable,omitempty"`
	IsBondEip             *bool   `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	MaxInstances          *int64  `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	MinInstances          *int64  `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	Type                  *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CreateClusterNodePoolRequestAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestAutoScaling) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetEipBandwidth(v int64) *CreateClusterNodePoolRequestAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetEipInternetChargeType(v string) *CreateClusterNodePoolRequestAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetEnable(v bool) *CreateClusterNodePoolRequestAutoScaling {
	s.Enable = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetIsBondEip(v bool) *CreateClusterNodePoolRequestAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetMaxInstances(v int64) *CreateClusterNodePoolRequestAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetMinInstances(v int64) *CreateClusterNodePoolRequestAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *CreateClusterNodePoolRequestAutoScaling) SetType(v string) *CreateClusterNodePoolRequestAutoScaling {
	s.Type = &v
	return s
}

type CreateClusterNodePoolRequestInterconnectConfig struct {
	Bandwidth      *int64  `json:"bandwidth,omitempty" xml:"bandwidth,omitempty"`
	CcnId          *string `json:"ccn_id,omitempty" xml:"ccn_id,omitempty"`
	CcnRegionId    *string `json:"ccn_region_id,omitempty" xml:"ccn_region_id,omitempty"`
	CenId          *string `json:"cen_id,omitempty" xml:"cen_id,omitempty"`
	ImprovedPeriod *string `json:"improved_period,omitempty" xml:"improved_period,omitempty"`
}

func (s CreateClusterNodePoolRequestInterconnectConfig) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestInterconnectConfig) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestInterconnectConfig) SetBandwidth(v int64) *CreateClusterNodePoolRequestInterconnectConfig {
	s.Bandwidth = &v
	return s
}

func (s *CreateClusterNodePoolRequestInterconnectConfig) SetCcnId(v string) *CreateClusterNodePoolRequestInterconnectConfig {
	s.CcnId = &v
	return s
}

func (s *CreateClusterNodePoolRequestInterconnectConfig) SetCcnRegionId(v string) *CreateClusterNodePoolRequestInterconnectConfig {
	s.CcnRegionId = &v
	return s
}

func (s *CreateClusterNodePoolRequestInterconnectConfig) SetCenId(v string) *CreateClusterNodePoolRequestInterconnectConfig {
	s.CenId = &v
	return s
}

func (s *CreateClusterNodePoolRequestInterconnectConfig) SetImprovedPeriod(v string) *CreateClusterNodePoolRequestInterconnectConfig {
	s.ImprovedPeriod = &v
	return s
}

type CreateClusterNodePoolRequestKubernetesConfig struct {
	CmsEnabled     *bool    `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	CpuPolicy      *string  `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	Labels         []*Tag   `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	NodeNameMode   *string  `json:"node_name_mode,omitempty" xml:"node_name_mode,omitempty"`
	Runtime        *string  `json:"runtime,omitempty" xml:"runtime,omitempty"`
	RuntimeVersion *string  `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	Taints         []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	UserData       *string  `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s CreateClusterNodePoolRequestKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestKubernetesConfig) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetCmsEnabled(v bool) *CreateClusterNodePoolRequestKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetCpuPolicy(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetLabels(v []*Tag) *CreateClusterNodePoolRequestKubernetesConfig {
	s.Labels = v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetNodeNameMode(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.NodeNameMode = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetRuntime(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetRuntimeVersion(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetTaints(v []*Taint) *CreateClusterNodePoolRequestKubernetesConfig {
	s.Taints = v
	return s
}

func (s *CreateClusterNodePoolRequestKubernetesConfig) SetUserData(v string) *CreateClusterNodePoolRequestKubernetesConfig {
	s.UserData = &v
	return s
}

type CreateClusterNodePoolRequestManagement struct {
	AutoRepair    *bool                                                `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	Enable        *bool                                                `json:"enable,omitempty" xml:"enable,omitempty"`
	UpgradeConfig *CreateClusterNodePoolRequestManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s CreateClusterNodePoolRequestManagement) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestManagement) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestManagement) SetAutoRepair(v bool) *CreateClusterNodePoolRequestManagement {
	s.AutoRepair = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagement) SetEnable(v bool) *CreateClusterNodePoolRequestManagement {
	s.Enable = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagement) SetUpgradeConfig(v *CreateClusterNodePoolRequestManagementUpgradeConfig) *CreateClusterNodePoolRequestManagement {
	s.UpgradeConfig = v
	return s
}

type CreateClusterNodePoolRequestManagementUpgradeConfig struct {
	AutoUpgrade     *bool  `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	MaxUnavailable  *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
	Surge           *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
}

func (s CreateClusterNodePoolRequestManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetAutoUpgrade(v bool) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetMaxUnavailable(v int64) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetSurge(v int64) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *CreateClusterNodePoolRequestManagementUpgradeConfig) SetSurgePercentage(v int64) *CreateClusterNodePoolRequestManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

type CreateClusterNodePoolRequestNodepoolInfo struct {
	Name            *string `json:"name,omitempty" xml:"name,omitempty"`
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	Type            *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CreateClusterNodePoolRequestNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestNodepoolInfo) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestNodepoolInfo) SetName(v string) *CreateClusterNodePoolRequestNodepoolInfo {
	s.Name = &v
	return s
}

func (s *CreateClusterNodePoolRequestNodepoolInfo) SetResourceGroupId(v string) *CreateClusterNodePoolRequestNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

func (s *CreateClusterNodePoolRequestNodepoolInfo) SetType(v string) *CreateClusterNodePoolRequestNodepoolInfo {
	s.Type = &v
	return s
}

type CreateClusterNodePoolRequestScalingGroup struct {
	AutoRenew                           *bool                                                       `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	AutoRenewPeriod                     *int64                                                      `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	CompensateWithOnDemand              *bool                                                       `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	DataDisks                           []*DataDisk                                                 `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	DeploymentsetId                     *string                                                     `json:"deploymentset_id,omitempty" xml:"deploymentset_id,omitempty"`
	DesiredSize                         *int64                                                      `json:"desired_size,omitempty" xml:"desired_size,omitempty"`
	ImageId                             *string                                                     `json:"image_id,omitempty" xml:"image_id,omitempty"`
	ImageType                           *string                                                     `json:"image_type,omitempty" xml:"image_type,omitempty"`
	InstanceChargeType                  *string                                                     `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	InstanceTypes                       []*string                                                   `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	InternetChargeType                  *string                                                     `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	InternetMaxBandwidthOut             *int64                                                      `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
	KeyPair                             *string                                                     `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	LoginPassword                       *string                                                     `json:"login_password,omitempty" xml:"login_password,omitempty"`
	MultiAzPolicy                       *string                                                     `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	OnDemandBaseCapacity                *int64                                                      `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	OnDemandPercentageAboveBaseCapacity *int64                                                      `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	Period                              *int64                                                      `json:"period,omitempty" xml:"period,omitempty"`
	PeriodUnit                          *string                                                     `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	Platform                            *string                                                     `json:"platform,omitempty" xml:"platform,omitempty"`
	PrivatePoolOptions                  *CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions `json:"private_pool_options,omitempty" xml:"private_pool_options,omitempty" type:"Struct"`
	RdsInstances                        []*string                                                   `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	ScalingPolicy                       *string                                                     `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	SecurityGroupId                     *string                                                     `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	SecurityGroupIds                    []*string                                                   `json:"security_group_ids,omitempty" xml:"security_group_ids,omitempty" type:"Repeated"`
	SpotInstancePools                   *int64                                                      `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	SpotInstanceRemedy                  *bool                                                       `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	SpotPriceLimit                      []*CreateClusterNodePoolRequestScalingGroupSpotPriceLimit   `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	SpotStrategy                        *string                                                     `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	SystemDiskBurstingEnabled           *bool                                                       `json:"system_disk_bursting_enabled,omitempty" xml:"system_disk_bursting_enabled,omitempty"`
	SystemDiskCategory                  *string                                                     `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	SystemDiskPerformanceLevel          *string                                                     `json:"system_disk_performance_level,omitempty" xml:"system_disk_performance_level,omitempty"`
	SystemDiskProvisionedIops           *int64                                                      `json:"system_disk_provisioned_iops,omitempty" xml:"system_disk_provisioned_iops,omitempty"`
	SystemDiskSize                      *int64                                                      `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	Tags                                []*CreateClusterNodePoolRequestScalingGroupTags             `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	VswitchIds                          []*string                                                   `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
}

func (s CreateClusterNodePoolRequestScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestScalingGroup) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetAutoRenew(v bool) *CreateClusterNodePoolRequestScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetAutoRenewPeriod(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetCompensateWithOnDemand(v bool) *CreateClusterNodePoolRequestScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetDataDisks(v []*DataDisk) *CreateClusterNodePoolRequestScalingGroup {
	s.DataDisks = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetDeploymentsetId(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.DeploymentsetId = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetDesiredSize(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.DesiredSize = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetImageId(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.ImageId = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetImageType(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.ImageType = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInstanceChargeType(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInstanceTypes(v []*string) *CreateClusterNodePoolRequestScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInternetChargeType(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetInternetMaxBandwidthOut(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetKeyPair(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetLoginPassword(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetMultiAzPolicy(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetOnDemandBaseCapacity(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetPeriod(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.Period = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetPeriodUnit(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetPlatform(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.Platform = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetPrivatePoolOptions(v *CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions) *CreateClusterNodePoolRequestScalingGroup {
	s.PrivatePoolOptions = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetRdsInstances(v []*string) *CreateClusterNodePoolRequestScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetScalingPolicy(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSecurityGroupId(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSecurityGroupIds(v []*string) *CreateClusterNodePoolRequestScalingGroup {
	s.SecurityGroupIds = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotInstancePools(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotInstanceRemedy(v bool) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotPriceLimit(v []*CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSpotStrategy(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSystemDiskBurstingEnabled(v bool) *CreateClusterNodePoolRequestScalingGroup {
	s.SystemDiskBurstingEnabled = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSystemDiskCategory(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSystemDiskPerformanceLevel(v string) *CreateClusterNodePoolRequestScalingGroup {
	s.SystemDiskPerformanceLevel = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSystemDiskProvisionedIops(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.SystemDiskProvisionedIops = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetSystemDiskSize(v int64) *CreateClusterNodePoolRequestScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetTags(v []*CreateClusterNodePoolRequestScalingGroupTags) *CreateClusterNodePoolRequestScalingGroup {
	s.Tags = v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroup) SetVswitchIds(v []*string) *CreateClusterNodePoolRequestScalingGroup {
	s.VswitchIds = v
	return s
}

type CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions struct {
	Id            *string `json:"id,omitempty" xml:"id,omitempty"`
	MatchCriteria *string `json:"match_criteria,omitempty" xml:"match_criteria,omitempty"`
}

func (s CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions) SetId(v string) *CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions {
	s.Id = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions) SetMatchCriteria(v string) *CreateClusterNodePoolRequestScalingGroupPrivatePoolOptions {
	s.MatchCriteria = &v
	return s
}

type CreateClusterNodePoolRequestScalingGroupSpotPriceLimit struct {
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	PriceLimit   *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) SetInstanceType(v string) *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit) SetPriceLimit(v string) *CreateClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type CreateClusterNodePoolRequestScalingGroupTags struct {
	Key   *string `json:"key,omitempty" xml:"key,omitempty"`
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s CreateClusterNodePoolRequestScalingGroupTags) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestScalingGroupTags) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestScalingGroupTags) SetKey(v string) *CreateClusterNodePoolRequestScalingGroupTags {
	s.Key = &v
	return s
}

func (s *CreateClusterNodePoolRequestScalingGroupTags) SetValue(v string) *CreateClusterNodePoolRequestScalingGroupTags {
	s.Value = &v
	return s
}

type CreateClusterNodePoolRequestTeeConfig struct {
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s CreateClusterNodePoolRequestTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolRequestTeeConfig) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolRequestTeeConfig) SetTeeEnable(v bool) *CreateClusterNodePoolRequestTeeConfig {
	s.TeeEnable = &v
	return s
}

type CreateClusterNodePoolResponseBody struct {
	// The ID of the node pool that is created.
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
}

func (s CreateClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolResponseBody) SetNodepoolId(v string) *CreateClusterNodePoolResponseBody {
	s.NodepoolId = &v
	return s
}

type CreateClusterNodePoolResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *CreateClusterNodePoolResponse) SetHeaders(v map[string]*string) *CreateClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *CreateClusterNodePoolResponse) SetStatusCode(v int32) *CreateClusterNodePoolResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateClusterNodePoolResponse) SetBody(v *CreateClusterNodePoolResponseBody) *CreateClusterNodePoolResponse {
	s.Body = v
	return s
}

type CreateEdgeMachineRequest struct {
	Hostname *string `json:"hostname,omitempty" xml:"hostname,omitempty"`
	Model    *string `json:"model,omitempty" xml:"model,omitempty"`
	Sn       *string `json:"sn,omitempty" xml:"sn,omitempty"`
}

func (s CreateEdgeMachineRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateEdgeMachineRequest) GoString() string {
	return s.String()
}

func (s *CreateEdgeMachineRequest) SetHostname(v string) *CreateEdgeMachineRequest {
	s.Hostname = &v
	return s
}

func (s *CreateEdgeMachineRequest) SetModel(v string) *CreateEdgeMachineRequest {
	s.Model = &v
	return s
}

func (s *CreateEdgeMachineRequest) SetSn(v string) *CreateEdgeMachineRequest {
	s.Sn = &v
	return s
}

type CreateEdgeMachineResponseBody struct {
	// The ID of the cloud-native box.
	EdgeMachineId *string `json:"edge_machine_id,omitempty" xml:"edge_machine_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
}

func (s CreateEdgeMachineResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateEdgeMachineResponseBody) GoString() string {
	return s.String()
}

func (s *CreateEdgeMachineResponseBody) SetEdgeMachineId(v string) *CreateEdgeMachineResponseBody {
	s.EdgeMachineId = &v
	return s
}

func (s *CreateEdgeMachineResponseBody) SetRequestId(v string) *CreateEdgeMachineResponseBody {
	s.RequestId = &v
	return s
}

type CreateEdgeMachineResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateEdgeMachineResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateEdgeMachineResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateEdgeMachineResponse) GoString() string {
	return s.String()
}

func (s *CreateEdgeMachineResponse) SetHeaders(v map[string]*string) *CreateEdgeMachineResponse {
	s.Headers = v
	return s
}

func (s *CreateEdgeMachineResponse) SetStatusCode(v int32) *CreateEdgeMachineResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateEdgeMachineResponse) SetBody(v *CreateEdgeMachineResponseBody) *CreateEdgeMachineResponse {
	s.Body = v
	return s
}

type CreateKubernetesTriggerRequest struct {
	Action    *string `json:"action,omitempty" xml:"action,omitempty"`
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	Type      *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CreateKubernetesTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateKubernetesTriggerRequest) GoString() string {
	return s.String()
}

func (s *CreateKubernetesTriggerRequest) SetAction(v string) *CreateKubernetesTriggerRequest {
	s.Action = &v
	return s
}

func (s *CreateKubernetesTriggerRequest) SetClusterId(v string) *CreateKubernetesTriggerRequest {
	s.ClusterId = &v
	return s
}

func (s *CreateKubernetesTriggerRequest) SetProjectId(v string) *CreateKubernetesTriggerRequest {
	s.ProjectId = &v
	return s
}

func (s *CreateKubernetesTriggerRequest) SetType(v string) *CreateKubernetesTriggerRequest {
	s.Type = &v
	return s
}

type CreateKubernetesTriggerResponseBody struct {
	// The action that the trigger performs. For example, a value of `redeploy` indicates that the trigger redeploys the application.
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
	// The ID of the ACK cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The ID of the trigger.
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// The name of the project.
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	// The type of trigger.
	//
	// Valid values:
	//
	// *   `deployment`: performs actions on Deployments.
	// *   `application`: performs actions on applications that are deployed in Application Center.
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CreateKubernetesTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateKubernetesTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *CreateKubernetesTriggerResponseBody) SetAction(v string) *CreateKubernetesTriggerResponseBody {
	s.Action = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetClusterId(v string) *CreateKubernetesTriggerResponseBody {
	s.ClusterId = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetId(v string) *CreateKubernetesTriggerResponseBody {
	s.Id = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetProjectId(v string) *CreateKubernetesTriggerResponseBody {
	s.ProjectId = &v
	return s
}

func (s *CreateKubernetesTriggerResponseBody) SetType(v string) *CreateKubernetesTriggerResponseBody {
	s.Type = &v
	return s
}

type CreateKubernetesTriggerResponse struct {
	Headers    map[string]*string                   `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                               `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateKubernetesTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateKubernetesTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateKubernetesTriggerResponse) GoString() string {
	return s.String()
}

func (s *CreateKubernetesTriggerResponse) SetHeaders(v map[string]*string) *CreateKubernetesTriggerResponse {
	s.Headers = v
	return s
}

func (s *CreateKubernetesTriggerResponse) SetStatusCode(v int32) *CreateKubernetesTriggerResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateKubernetesTriggerResponse) SetBody(v *CreateKubernetesTriggerResponseBody) *CreateKubernetesTriggerResponse {
	s.Body = v
	return s
}

type CreateTemplateRequest struct {
	Description  *string `json:"description,omitempty" xml:"description,omitempty"`
	Name         *string `json:"name,omitempty" xml:"name,omitempty"`
	Tags         *string `json:"tags,omitempty" xml:"tags,omitempty"`
	Template     *string `json:"template,omitempty" xml:"template,omitempty"`
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
}

func (s CreateTemplateRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateTemplateRequest) GoString() string {
	return s.String()
}

func (s *CreateTemplateRequest) SetDescription(v string) *CreateTemplateRequest {
	s.Description = &v
	return s
}

func (s *CreateTemplateRequest) SetName(v string) *CreateTemplateRequest {
	s.Name = &v
	return s
}

func (s *CreateTemplateRequest) SetTags(v string) *CreateTemplateRequest {
	s.Tags = &v
	return s
}

func (s *CreateTemplateRequest) SetTemplate(v string) *CreateTemplateRequest {
	s.Template = &v
	return s
}

func (s *CreateTemplateRequest) SetTemplateType(v string) *CreateTemplateRequest {
	s.TemplateType = &v
	return s
}

type CreateTemplateResponseBody struct {
	// The ID of the template.
	TemplateId *string `json:"template_id,omitempty" xml:"template_id,omitempty"`
}

func (s CreateTemplateResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateTemplateResponseBody) GoString() string {
	return s.String()
}

func (s *CreateTemplateResponseBody) SetTemplateId(v string) *CreateTemplateResponseBody {
	s.TemplateId = &v
	return s
}

type CreateTemplateResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *CreateTemplateResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s CreateTemplateResponse) String() string {
	return tea.Prettify(s)
}

func (s CreateTemplateResponse) GoString() string {
	return s.String()
}

func (s *CreateTemplateResponse) SetHeaders(v map[string]*string) *CreateTemplateResponse {
	s.Headers = v
	return s
}

func (s *CreateTemplateResponse) SetStatusCode(v int32) *CreateTemplateResponse {
	s.StatusCode = &v
	return s
}

func (s *CreateTemplateResponse) SetBody(v *CreateTemplateResponseBody) *CreateTemplateResponse {
	s.Body = v
	return s
}

type CreateTriggerRequest struct {
	Action    *string `json:"action,omitempty" xml:"action,omitempty"`
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	Type      *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CreateTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s CreateTriggerRequest) GoString() string {
	return s.String()
}

func (s *CreateTriggerRequest) SetAction(v string) *CreateTriggerRequest {
	s.Action = &v
	return s
}

func (s *CreateTriggerRequest) SetClusterId(v string) *CreateTriggerRequest {
	s.ClusterId = &v
	return s
}

func (s *CreateTriggerRequest) SetProjectId(v string) *CreateTriggerRequest {
	s.ProjectId = &v
	return s
}

func (s *CreateTriggerRequest) SetType(v string) *CreateTriggerRequest {
	s.Type = &v
	return s
}

type CreateTriggerResponseBody struct {
	// The action that the trigger performs. For example, a value of `redeploy` indicates that the trigger redeploys the application.
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
	// The ID of the cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The ID of the trigger.
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// The name of the project.
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	// The type of trigger. Default value: deployment.
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s CreateTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s CreateTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *CreateTriggerResponseBody) SetAction(v string) *CreateTriggerResponseBody {
	s.Action = &v
	return s
}

func (s *CreateTriggerResponseBody) SetClusterId(v string) *CreateTriggerResponseBody {
	s.ClusterId = &v
	return s
}

func (s *CreateTriggerResponseBody) SetId(v string) *CreateTriggerResponseBody {
	s.Id = &v
	return s
}

func (s *CreateTriggerResponseBody) SetProjectId(v string) *CreateTriggerResponseBody {
	s.ProjectId = &v
	return s
}

func (s *CreateTriggerResponseBody) SetType(v string) *CreateTriggerResponseBody {
	s.Type = &v
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

type DeleteAlertContactResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteAlertContactResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAlertContactResponse) GoString() string {
	return s.String()
}

func (s *DeleteAlertContactResponse) SetHeaders(v map[string]*string) *DeleteAlertContactResponse {
	s.Headers = v
	return s
}

func (s *DeleteAlertContactResponse) SetStatusCode(v int32) *DeleteAlertContactResponse {
	s.StatusCode = &v
	return s
}

type DeleteAlertContactGroupResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteAlertContactGroupResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteAlertContactGroupResponse) GoString() string {
	return s.String()
}

func (s *DeleteAlertContactGroupResponse) SetHeaders(v map[string]*string) *DeleteAlertContactGroupResponse {
	s.Headers = v
	return s
}

func (s *DeleteAlertContactGroupResponse) SetStatusCode(v int32) *DeleteAlertContactGroupResponse {
	s.StatusCode = &v
	return s
}

type DeleteClusterRequest struct {
	KeepSlb            *bool     `json:"keep_slb,omitempty" xml:"keep_slb,omitempty"`
	RetainAllResources *bool     `json:"retain_all_resources,omitempty" xml:"retain_all_resources,omitempty"`
	RetainResources    []*string `json:"retain_resources,omitempty" xml:"retain_resources,omitempty" type:"Repeated"`
}

func (s DeleteClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterRequest) GoString() string {
	return s.String()
}

func (s *DeleteClusterRequest) SetKeepSlb(v bool) *DeleteClusterRequest {
	s.KeepSlb = &v
	return s
}

func (s *DeleteClusterRequest) SetRetainAllResources(v bool) *DeleteClusterRequest {
	s.RetainAllResources = &v
	return s
}

func (s *DeleteClusterRequest) SetRetainResources(v []*string) *DeleteClusterRequest {
	s.RetainResources = v
	return s
}

type DeleteClusterShrinkRequest struct {
	KeepSlb               *bool   `json:"keep_slb,omitempty" xml:"keep_slb,omitempty"`
	RetainAllResources    *bool   `json:"retain_all_resources,omitempty" xml:"retain_all_resources,omitempty"`
	RetainResourcesShrink *string `json:"retain_resources,omitempty" xml:"retain_resources,omitempty"`
}

func (s DeleteClusterShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterShrinkRequest) GoString() string {
	return s.String()
}

func (s *DeleteClusterShrinkRequest) SetKeepSlb(v bool) *DeleteClusterShrinkRequest {
	s.KeepSlb = &v
	return s
}

func (s *DeleteClusterShrinkRequest) SetRetainAllResources(v bool) *DeleteClusterShrinkRequest {
	s.RetainAllResources = &v
	return s
}

func (s *DeleteClusterShrinkRequest) SetRetainResourcesShrink(v string) *DeleteClusterShrinkRequest {
	s.RetainResourcesShrink = &v
	return s
}

type DeleteClusterResponseBody struct {
	// 任务ID。
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s DeleteClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteClusterResponseBody) SetTaskId(v string) *DeleteClusterResponseBody {
	s.TaskId = &v
	return s
}

type DeleteClusterResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterResponse) GoString() string {
	return s.String()
}

func (s *DeleteClusterResponse) SetHeaders(v map[string]*string) *DeleteClusterResponse {
	s.Headers = v
	return s
}

func (s *DeleteClusterResponse) SetStatusCode(v int32) *DeleteClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteClusterResponse) SetBody(v *DeleteClusterResponseBody) *DeleteClusterResponse {
	s.Body = v
	return s
}

type DeleteClusterNodepoolRequest struct {
	// false
	Force *bool `json:"force,omitempty" xml:"force,omitempty"`
}

func (s DeleteClusterNodepoolRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodepoolRequest) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodepoolRequest) SetForce(v bool) *DeleteClusterNodepoolRequest {
	s.Force = &v
	return s
}

type DeleteClusterNodepoolResponseBody struct {
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
}

func (s DeleteClusterNodepoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodepoolResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodepoolResponseBody) SetRequestId(v string) *DeleteClusterNodepoolResponseBody {
	s.RequestId = &v
	return s
}

type DeleteClusterNodepoolResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteClusterNodepoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteClusterNodepoolResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodepoolResponse) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodepoolResponse) SetHeaders(v map[string]*string) *DeleteClusterNodepoolResponse {
	s.Headers = v
	return s
}

func (s *DeleteClusterNodepoolResponse) SetStatusCode(v int32) *DeleteClusterNodepoolResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteClusterNodepoolResponse) SetBody(v *DeleteClusterNodepoolResponseBody) *DeleteClusterNodepoolResponse {
	s.Body = v
	return s
}

type DeleteClusterNodesRequest struct {
	DrainNode   *bool     `json:"drain_node,omitempty" xml:"drain_node,omitempty"`
	Nodes       []*string `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
	ReleaseNode *bool     `json:"release_node,omitempty" xml:"release_node,omitempty"`
}

func (s DeleteClusterNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodesRequest) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodesRequest) SetDrainNode(v bool) *DeleteClusterNodesRequest {
	s.DrainNode = &v
	return s
}

func (s *DeleteClusterNodesRequest) SetNodes(v []*string) *DeleteClusterNodesRequest {
	s.Nodes = v
	return s
}

func (s *DeleteClusterNodesRequest) SetReleaseNode(v bool) *DeleteClusterNodesRequest {
	s.ReleaseNode = &v
	return s
}

type DeleteClusterNodesResponseBody struct {
	// The ID of the ACK cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s DeleteClusterNodesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodesResponseBody) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodesResponseBody) SetClusterId(v string) *DeleteClusterNodesResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DeleteClusterNodesResponseBody) SetRequestId(v string) *DeleteClusterNodesResponseBody {
	s.RequestId = &v
	return s
}

func (s *DeleteClusterNodesResponseBody) SetTaskId(v string) *DeleteClusterNodesResponseBody {
	s.TaskId = &v
	return s
}

type DeleteClusterNodesResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeleteClusterNodesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeleteClusterNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteClusterNodesResponse) GoString() string {
	return s.String()
}

func (s *DeleteClusterNodesResponse) SetHeaders(v map[string]*string) *DeleteClusterNodesResponse {
	s.Headers = v
	return s
}

func (s *DeleteClusterNodesResponse) SetStatusCode(v int32) *DeleteClusterNodesResponse {
	s.StatusCode = &v
	return s
}

func (s *DeleteClusterNodesResponse) SetBody(v *DeleteClusterNodesResponseBody) *DeleteClusterNodesResponse {
	s.Body = v
	return s
}

type DeleteEdgeMachineRequest struct {
	Force *string `json:"force,omitempty" xml:"force,omitempty"`
}

func (s DeleteEdgeMachineRequest) String() string {
	return tea.Prettify(s)
}

func (s DeleteEdgeMachineRequest) GoString() string {
	return s.String()
}

func (s *DeleteEdgeMachineRequest) SetForce(v string) *DeleteEdgeMachineRequest {
	s.Force = &v
	return s
}

type DeleteEdgeMachineResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteEdgeMachineResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteEdgeMachineResponse) GoString() string {
	return s.String()
}

func (s *DeleteEdgeMachineResponse) SetHeaders(v map[string]*string) *DeleteEdgeMachineResponse {
	s.Headers = v
	return s
}

func (s *DeleteEdgeMachineResponse) SetStatusCode(v int32) *DeleteEdgeMachineResponse {
	s.StatusCode = &v
	return s
}

type DeleteKubernetesTriggerResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteKubernetesTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteKubernetesTriggerResponse) GoString() string {
	return s.String()
}

func (s *DeleteKubernetesTriggerResponse) SetHeaders(v map[string]*string) *DeleteKubernetesTriggerResponse {
	s.Headers = v
	return s
}

func (s *DeleteKubernetesTriggerResponse) SetStatusCode(v int32) *DeleteKubernetesTriggerResponse {
	s.StatusCode = &v
	return s
}

type DeletePolicyInstanceRequest struct {
	InstanceName *string `json:"instance_name,omitempty" xml:"instance_name,omitempty"`
}

func (s DeletePolicyInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s DeletePolicyInstanceRequest) GoString() string {
	return s.String()
}

func (s *DeletePolicyInstanceRequest) SetInstanceName(v string) *DeletePolicyInstanceRequest {
	s.InstanceName = &v
	return s
}

type DeletePolicyInstanceResponseBody struct {
	// The policy instances that are deleted.
	Instances []*string `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
}

func (s DeletePolicyInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeletePolicyInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *DeletePolicyInstanceResponseBody) SetInstances(v []*string) *DeletePolicyInstanceResponseBody {
	s.Instances = v
	return s
}

type DeletePolicyInstanceResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeletePolicyInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeletePolicyInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s DeletePolicyInstanceResponse) GoString() string {
	return s.String()
}

func (s *DeletePolicyInstanceResponse) SetHeaders(v map[string]*string) *DeletePolicyInstanceResponse {
	s.Headers = v
	return s
}

func (s *DeletePolicyInstanceResponse) SetStatusCode(v int32) *DeletePolicyInstanceResponse {
	s.StatusCode = &v
	return s
}

func (s *DeletePolicyInstanceResponse) SetBody(v *DeletePolicyInstanceResponseBody) *DeletePolicyInstanceResponse {
	s.Body = v
	return s
}

type DeleteTemplateResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s DeleteTemplateResponse) String() string {
	return tea.Prettify(s)
}

func (s DeleteTemplateResponse) GoString() string {
	return s.String()
}

func (s *DeleteTemplateResponse) SetHeaders(v map[string]*string) *DeleteTemplateResponse {
	s.Headers = v
	return s
}

func (s *DeleteTemplateResponse) SetStatusCode(v int32) *DeleteTemplateResponse {
	s.StatusCode = &v
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

type DeployPolicyInstanceRequest struct {
	Action     *string                `json:"action,omitempty" xml:"action,omitempty"`
	Namespaces []*string              `json:"namespaces,omitempty" xml:"namespaces,omitempty" type:"Repeated"`
	Parameters map[string]interface{} `json:"parameters,omitempty" xml:"parameters,omitempty"`
}

func (s DeployPolicyInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s DeployPolicyInstanceRequest) GoString() string {
	return s.String()
}

func (s *DeployPolicyInstanceRequest) SetAction(v string) *DeployPolicyInstanceRequest {
	s.Action = &v
	return s
}

func (s *DeployPolicyInstanceRequest) SetNamespaces(v []*string) *DeployPolicyInstanceRequest {
	s.Namespaces = v
	return s
}

func (s *DeployPolicyInstanceRequest) SetParameters(v map[string]interface{}) *DeployPolicyInstanceRequest {
	s.Parameters = v
	return s
}

type DeployPolicyInstanceResponseBody struct {
	// The policy instances that are deployed.
	Instances []*string `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
}

func (s DeployPolicyInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DeployPolicyInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *DeployPolicyInstanceResponseBody) SetInstances(v []*string) *DeployPolicyInstanceResponseBody {
	s.Instances = v
	return s
}

type DeployPolicyInstanceResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DeployPolicyInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DeployPolicyInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s DeployPolicyInstanceResponse) GoString() string {
	return s.String()
}

func (s *DeployPolicyInstanceResponse) SetHeaders(v map[string]*string) *DeployPolicyInstanceResponse {
	s.Headers = v
	return s
}

func (s *DeployPolicyInstanceResponse) SetStatusCode(v int32) *DeployPolicyInstanceResponse {
	s.StatusCode = &v
	return s
}

func (s *DeployPolicyInstanceResponse) SetBody(v *DeployPolicyInstanceResponseBody) *DeployPolicyInstanceResponse {
	s.Body = v
	return s
}

type DescirbeWorkflowResponseBody struct {
	// The time when the workflow was created.
	CreateTime *string `json:"create_time,omitempty" xml:"create_time,omitempty"`
	// The duration of the workflow.
	Duration *string `json:"duration,omitempty" xml:"duration,omitempty"`
	// The time when the workflow ended.
	FinishTime *string `json:"finish_time,omitempty" xml:"finish_time,omitempty"`
	// The size of the input data.
	InputDataSize *string `json:"input_data_size,omitempty" xml:"input_data_size,omitempty"`
	// The name of the workflow.
	JobName *string `json:"job_name,omitempty" xml:"job_name,omitempty"`
	// The namespace to which the workflow belongs.
	JobNamespace *string `json:"job_namespace,omitempty" xml:"job_namespace,omitempty"`
	// The size of the output data.
	OutputDataSize *string `json:"output_data_size,omitempty" xml:"output_data_size,omitempty"`
	// The current state of the workflow.
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// The number of base pairs.
	TotalBases *string `json:"total_bases,omitempty" xml:"total_bases,omitempty"`
	// The number of reads.
	TotalReads *string `json:"total_reads,omitempty" xml:"total_reads,omitempty"`
	// The user input parameters.
	UserInputData *string `json:"user_input_data,omitempty" xml:"user_input_data,omitempty"`
}

func (s DescirbeWorkflowResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescirbeWorkflowResponseBody) GoString() string {
	return s.String()
}

func (s *DescirbeWorkflowResponseBody) SetCreateTime(v string) *DescirbeWorkflowResponseBody {
	s.CreateTime = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetDuration(v string) *DescirbeWorkflowResponseBody {
	s.Duration = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetFinishTime(v string) *DescirbeWorkflowResponseBody {
	s.FinishTime = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetInputDataSize(v string) *DescirbeWorkflowResponseBody {
	s.InputDataSize = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetJobName(v string) *DescirbeWorkflowResponseBody {
	s.JobName = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetJobNamespace(v string) *DescirbeWorkflowResponseBody {
	s.JobNamespace = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetOutputDataSize(v string) *DescirbeWorkflowResponseBody {
	s.OutputDataSize = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetStatus(v string) *DescirbeWorkflowResponseBody {
	s.Status = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetTotalBases(v string) *DescirbeWorkflowResponseBody {
	s.TotalBases = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetTotalReads(v string) *DescirbeWorkflowResponseBody {
	s.TotalReads = &v
	return s
}

func (s *DescirbeWorkflowResponseBody) SetUserInputData(v string) *DescirbeWorkflowResponseBody {
	s.UserInputData = &v
	return s
}

type DescirbeWorkflowResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescirbeWorkflowResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescirbeWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s DescirbeWorkflowResponse) GoString() string {
	return s.String()
}

func (s *DescirbeWorkflowResponse) SetHeaders(v map[string]*string) *DescirbeWorkflowResponse {
	s.Headers = v
	return s
}

func (s *DescirbeWorkflowResponse) SetStatusCode(v int32) *DescirbeWorkflowResponse {
	s.StatusCode = &v
	return s
}

func (s *DescirbeWorkflowResponse) SetBody(v *DescirbeWorkflowResponseBody) *DescirbeWorkflowResponse {
	s.Body = v
	return s
}

type DescribeAddonsRequest struct {
	ClusterProfile *string `json:"cluster_profile,omitempty" xml:"cluster_profile,omitempty"`
	ClusterSpec    *string `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	ClusterType    *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	ClusterVersion *string `json:"cluster_version,omitempty" xml:"cluster_version,omitempty"`
	Region         *string `json:"region,omitempty" xml:"region,omitempty"`
}

func (s DescribeAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsRequest) GoString() string {
	return s.String()
}

func (s *DescribeAddonsRequest) SetClusterProfile(v string) *DescribeAddonsRequest {
	s.ClusterProfile = &v
	return s
}

func (s *DescribeAddonsRequest) SetClusterSpec(v string) *DescribeAddonsRequest {
	s.ClusterSpec = &v
	return s
}

func (s *DescribeAddonsRequest) SetClusterType(v string) *DescribeAddonsRequest {
	s.ClusterType = &v
	return s
}

func (s *DescribeAddonsRequest) SetClusterVersion(v string) *DescribeAddonsRequest {
	s.ClusterVersion = &v
	return s
}

func (s *DescribeAddonsRequest) SetRegion(v string) *DescribeAddonsRequest {
	s.Region = &v
	return s
}

type DescribeAddonsResponseBody struct {
	// The details of the returned components.
	ComponentGroups    []*DescribeAddonsResponseBodyComponentGroups `json:"ComponentGroups,omitempty" xml:"ComponentGroups,omitempty" type:"Repeated"`
	StandardComponents map[string]*StandardComponentsValue          `json:"StandardComponents,omitempty" xml:"StandardComponents,omitempty"`
}

func (s DescribeAddonsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponseBody) SetComponentGroups(v []*DescribeAddonsResponseBodyComponentGroups) *DescribeAddonsResponseBody {
	s.ComponentGroups = v
	return s
}

func (s *DescribeAddonsResponseBody) SetStandardComponents(v map[string]*StandardComponentsValue) *DescribeAddonsResponseBody {
	s.StandardComponents = v
	return s
}

type DescribeAddonsResponseBodyComponentGroups struct {
	// The name of the component group.
	GroupName *string `json:"group_name,omitempty" xml:"group_name,omitempty"`
	// The names of the components in the component group.
	Items []*DescribeAddonsResponseBodyComponentGroupsItems `json:"items,omitempty" xml:"items,omitempty" type:"Repeated"`
}

func (s DescribeAddonsResponseBodyComponentGroups) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponseBodyComponentGroups) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponseBodyComponentGroups) SetGroupName(v string) *DescribeAddonsResponseBodyComponentGroups {
	s.GroupName = &v
	return s
}

func (s *DescribeAddonsResponseBodyComponentGroups) SetItems(v []*DescribeAddonsResponseBodyComponentGroupsItems) *DescribeAddonsResponseBodyComponentGroups {
	s.Items = v
	return s
}

type DescribeAddonsResponseBodyComponentGroupsItems struct {
	// The name of the component.
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
}

func (s DescribeAddonsResponseBodyComponentGroupsItems) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponseBodyComponentGroupsItems) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponseBodyComponentGroupsItems) SetName(v string) *DescribeAddonsResponseBodyComponentGroupsItems {
	s.Name = &v
	return s
}

type DescribeAddonsResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeAddonsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeAddonsResponse) GoString() string {
	return s.String()
}

func (s *DescribeAddonsResponse) SetHeaders(v map[string]*string) *DescribeAddonsResponse {
	s.Headers = v
	return s
}

func (s *DescribeAddonsResponse) SetStatusCode(v int32) *DescribeAddonsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeAddonsResponse) SetBody(v *DescribeAddonsResponseBody) *DescribeAddonsResponse {
	s.Body = v
	return s
}

type DescribeClusterAddonInstanceResponseBody struct {
	Config  *string `json:"config,omitempty" xml:"config,omitempty"`
	Name    *string `json:"name,omitempty" xml:"name,omitempty"`
	State   *string `json:"state,omitempty" xml:"state,omitempty"`
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s DescribeClusterAddonInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonInstanceResponseBody) SetConfig(v string) *DescribeClusterAddonInstanceResponseBody {
	s.Config = &v
	return s
}

func (s *DescribeClusterAddonInstanceResponseBody) SetName(v string) *DescribeClusterAddonInstanceResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeClusterAddonInstanceResponseBody) SetState(v string) *DescribeClusterAddonInstanceResponseBody {
	s.State = &v
	return s
}

func (s *DescribeClusterAddonInstanceResponseBody) SetVersion(v string) *DescribeClusterAddonInstanceResponseBody {
	s.Version = &v
	return s
}

type DescribeClusterAddonInstanceResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterAddonInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonInstanceResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonInstanceResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonInstanceResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonInstanceResponse) SetStatusCode(v int32) *DescribeClusterAddonInstanceResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterAddonInstanceResponse) SetBody(v *DescribeClusterAddonInstanceResponseBody) *DescribeClusterAddonInstanceResponse {
	s.Body = v
	return s
}

type DescribeClusterAddonMetadataResponseBody struct {
	// The schema of component parameters.
	ConfigSchema *string `json:"config_schema,omitempty" xml:"config_schema,omitempty"`
	// The name of the component.
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// The version of the component.
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s DescribeClusterAddonMetadataResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonMetadataResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonMetadataResponseBody) SetConfigSchema(v string) *DescribeClusterAddonMetadataResponseBody {
	s.ConfigSchema = &v
	return s
}

func (s *DescribeClusterAddonMetadataResponseBody) SetName(v string) *DescribeClusterAddonMetadataResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeClusterAddonMetadataResponseBody) SetVersion(v string) *DescribeClusterAddonMetadataResponseBody {
	s.Version = &v
	return s
}

type DescribeClusterAddonMetadataResponse struct {
	Headers    map[string]*string                        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterAddonMetadataResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonMetadataResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonMetadataResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonMetadataResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonMetadataResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonMetadataResponse) SetStatusCode(v int32) *DescribeClusterAddonMetadataResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterAddonMetadataResponse) SetBody(v *DescribeClusterAddonMetadataResponseBody) *DescribeClusterAddonMetadataResponse {
	s.Body = v
	return s
}

type DescribeClusterAddonUpgradeStatusResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       map[string]interface{} `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonUpgradeStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonUpgradeStatusResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonUpgradeStatusResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonUpgradeStatusResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonUpgradeStatusResponse) SetStatusCode(v int32) *DescribeClusterAddonUpgradeStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterAddonUpgradeStatusResponse) SetBody(v map[string]interface{}) *DescribeClusterAddonUpgradeStatusResponse {
	s.Body = v
	return s
}

type DescribeClusterAddonsUpgradeStatusRequest struct {
	ComponentIds []*string `json:"componentIds,omitempty" xml:"componentIds,omitempty" type:"Repeated"`
}

func (s DescribeClusterAddonsUpgradeStatusRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsUpgradeStatusRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsUpgradeStatusRequest) SetComponentIds(v []*string) *DescribeClusterAddonsUpgradeStatusRequest {
	s.ComponentIds = v
	return s
}

type DescribeClusterAddonsUpgradeStatusShrinkRequest struct {
	ComponentIdsShrink *string `json:"componentIds,omitempty" xml:"componentIds,omitempty"`
}

func (s DescribeClusterAddonsUpgradeStatusShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsUpgradeStatusShrinkRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsUpgradeStatusShrinkRequest) SetComponentIdsShrink(v string) *DescribeClusterAddonsUpgradeStatusShrinkRequest {
	s.ComponentIdsShrink = &v
	return s
}

type DescribeClusterAddonsUpgradeStatusResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       map[string]interface{} `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonsUpgradeStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsUpgradeStatusResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsUpgradeStatusResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonsUpgradeStatusResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonsUpgradeStatusResponse) SetStatusCode(v int32) *DescribeClusterAddonsUpgradeStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterAddonsUpgradeStatusResponse) SetBody(v map[string]interface{}) *DescribeClusterAddonsUpgradeStatusResponse {
	s.Body = v
	return s
}

type DescribeClusterAddonsVersionResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       map[string]interface{} `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAddonsVersionResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAddonsVersionResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAddonsVersionResponse) SetHeaders(v map[string]*string) *DescribeClusterAddonsVersionResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAddonsVersionResponse) SetStatusCode(v int32) *DescribeClusterAddonsVersionResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterAddonsVersionResponse) SetBody(v map[string]interface{}) *DescribeClusterAddonsVersionResponse {
	s.Body = v
	return s
}

type DescribeClusterAttachScriptsRequest struct {
	Arch             *string   `json:"arch,omitempty" xml:"arch,omitempty"`
	FormatDisk       *bool     `json:"format_disk,omitempty" xml:"format_disk,omitempty"`
	KeepInstanceName *bool     `json:"keep_instance_name,omitempty" xml:"keep_instance_name,omitempty"`
	NodepoolId       *string   `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	Options          *string   `json:"options,omitempty" xml:"options,omitempty"`
	RdsInstances     []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
}

func (s DescribeClusterAttachScriptsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAttachScriptsRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterAttachScriptsRequest) SetArch(v string) *DescribeClusterAttachScriptsRequest {
	s.Arch = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetFormatDisk(v bool) *DescribeClusterAttachScriptsRequest {
	s.FormatDisk = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetKeepInstanceName(v bool) *DescribeClusterAttachScriptsRequest {
	s.KeepInstanceName = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetNodepoolId(v string) *DescribeClusterAttachScriptsRequest {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetOptions(v string) *DescribeClusterAttachScriptsRequest {
	s.Options = &v
	return s
}

func (s *DescribeClusterAttachScriptsRequest) SetRdsInstances(v []*string) *DescribeClusterAttachScriptsRequest {
	s.RdsInstances = v
	return s
}

type DescribeClusterAttachScriptsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *string            `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterAttachScriptsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterAttachScriptsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterAttachScriptsResponse) SetHeaders(v map[string]*string) *DescribeClusterAttachScriptsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterAttachScriptsResponse) SetStatusCode(v int32) *DescribeClusterAttachScriptsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterAttachScriptsResponse) SetBody(v string) *DescribeClusterAttachScriptsResponse {
	s.Body = &v
	return s
}

type DescribeClusterDetailResponseBody struct {
	// The ID of the queried ACK cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The type of the managed Kubernetes cluster. This parameter is returned for a managed Kubernetes cluster. Valid values:
	//
	// *   `ack.pro.small`: professional managed Kubernetes cluster.
	// *   `ack.standard`: standard managed Kubernetes cluster.
	ClusterSpec *string `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	// The type of the cluster. Valid values:
	//
	// *   `Kubernetes`: dedicated Kubernetes cluster
	// *   `ManagedKubernetes`: managed Kubernetes cluster
	// *   `Ask`: ASK cluster
	// *   `ExternalKubernetes`: registered external Kubernetes cluster
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	// The time when the cluster was created.
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// The current Kubernetes version of the cluster. For more information about the Kubernetes versions supported by ACK, see [Release notes for Kubernetes versions](~~185269~~).
	CurrentVersion *string `json:"current_version,omitempty" xml:"current_version,omitempty"`
	// Indicates whether deletion protection is enabled. If deletion protection is enabled, the cluster cannot be deleted in the ACK console or by calling the API. Valid values:
	//
	// *   `true`: Deletion protection is enabled. You cannot delete the cluster in the ACK console or by calling the API.
	// *   `false`: Deletion protection is not enabled. You can delete the cluster in the ACK console or by calling the API.
	DeletionProtection *bool `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	// The Docker version that is used by the cluster.
	DockerVersion *string `json:"docker_version,omitempty" xml:"docker_version,omitempty"`
	// The ID of the Server Load Balancer (SLB) instance that is used for the Ingress of the cluster.
	ExternalLoadbalancerId *string `json:"external_loadbalancer_id,omitempty" xml:"external_loadbalancer_id,omitempty"`
	// The Kubernetes version that is initially used by the cluster.
	InitVersion *string `json:"init_version,omitempty" xml:"init_version,omitempty"`
	// The maintenance window of the cluster. This feature is available in only professional managed Kubernetes clusters.
	MaintenanceWindow *MaintenanceWindow `json:"maintenance_window,omitempty" xml:"maintenance_window,omitempty"`
	// The address of the cluster. It includes an internal endpoint and a public endpoint.
	MasterUrl *string `json:"master_url,omitempty" xml:"master_url,omitempty"`
	// The metadata of the cluster.
	MetaData *string `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	// The name of the cluster.
	//
	// The name must be 1 to 63 characters in length, and can contain digits, letters, and hyphens (-). It cannot start with a hyphen (-).
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// The network mode of the cluster. Valid values:
	//
	// *   `classic`: the classic network
	// *   `vpc`: virtual private cloud (VPC)
	// *   `overlay`: overlay network
	// *   `calico`: network powered by Calico
	//
	// Default value`: vpc`.
	NetworkMode *string `json:"network_mode,omitempty" xml:"network_mode,omitempty"`
	// The Kubernetes version to which the cluster can be upgraded.
	NextVersion *string            `json:"next_version,omitempty" xml:"next_version,omitempty"`
	Parameters  map[string]*string `json:"parameters,omitempty" xml:"parameters,omitempty"`
	// Indicates whether Alibaba Cloud DNS PrivateZone is enabled.
	//
	// *   `true`: indicates that Alibaba Cloud DNS PrivateZone is enabled.
	// *   `false`: indicates that Alibaba Cloud DNS PrivateZone is not enabled.
	PrivateZone *bool `json:"private_zone,omitempty" xml:"private_zone,omitempty"`
	// Indicates the scenario in which the cluster is used. Valid values:
	//
	// *   `Default`: indicates that the cluster is used in non-edge computing scenarios.
	// *   `Edge`: indicates that the ACK cluster is used in edge computing scenarios.
	Profile *string `json:"profile,omitempty" xml:"profile,omitempty"`
	// The ID of the region where the cluster is deployed.
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// The ID of the resource group to which the cluster belongs.
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// The ID of the security group to which the instances of the cluster belong.
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// The number of nodes in the cluster. Master nodes and worker nodes are included.
	Size *int64 `json:"size,omitempty" xml:"size,omitempty"`
	// The state of the cluster. Valid values:
	//
	// *   `initial`: The cluster is being created.
	// *   `failed`: The cluster failed to be created.
	// *   `running`: The cluster is running.
	// *   `updating`: The cluster is being upgraded.
	// *   `updating_failed`: The cluster failed to be upgraded.
	// *   `scaling`: The cluster is being scaled.
	// *   `waiting`: The registered cluster is waiting for connecting.
	// *   `disconnected`: The registeredcluster is disconnected.
	// *   `stopped`: The cluster is stopped.
	// *   `deleting`: The cluster is being deleted.
	// *   `deleted`: The cluster is deleted.
	// *   `delete_failed`: The cluster failed to be deleted.
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// The pod CIDR block. It must be a valid and private CIDR block, and must be one of the following CIDR blocks or their subnets:
	//
	// *   10.0.0.0/8
	// *   172.16-31.0.0/12-16
	// *   192.168.0.0/16
	//
	// The pod CIDR block cannot overlap with that of the VPC or those of the ACK clusters that are deployed in the VPC.
	//
	// For more information about the network segmentation of ACK clusters, see [Plan CIDR blocks for ACK clusters in a VPC](~~186964~~).
	SubnetCidr *string `json:"subnet_cidr,omitempty" xml:"subnet_cidr,omitempty"`
	// The labels of the cluster.
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// The time when the cluster was updated.
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
	// The ID of the VPC where the cluster is deployed. This parameter is required when you create an ACK cluster.
	VpcId *string `json:"vpc_id,omitempty" xml:"vpc_id,omitempty"`
	// The IDs of the vSwitches. You can select one to three vSwitches when you create an ACK cluster. vSwitches in different zones are recommended to ensure high availability.
	VswitchId *string `json:"vswitch_id,omitempty" xml:"vswitch_id,omitempty"`
	// The name of the worker RAM role. The RAM role is assigned to the worker nodes of the cluster and allows the worker nodes to manage Elastic Compute Service (ECS) instances.
	WorkerRamRoleName *string `json:"worker_ram_role_name,omitempty" xml:"worker_ram_role_name,omitempty"`
	// The ID of the zone where the cluster is deployed.
	ZoneId *string `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
}

func (s DescribeClusterDetailResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterDetailResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterDetailResponseBody) SetClusterId(v string) *DescribeClusterDetailResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetClusterSpec(v string) *DescribeClusterDetailResponseBody {
	s.ClusterSpec = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetClusterType(v string) *DescribeClusterDetailResponseBody {
	s.ClusterType = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetCreated(v string) *DescribeClusterDetailResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetCurrentVersion(v string) *DescribeClusterDetailResponseBody {
	s.CurrentVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetDeletionProtection(v bool) *DescribeClusterDetailResponseBody {
	s.DeletionProtection = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetDockerVersion(v string) *DescribeClusterDetailResponseBody {
	s.DockerVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetExternalLoadbalancerId(v string) *DescribeClusterDetailResponseBody {
	s.ExternalLoadbalancerId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetInitVersion(v string) *DescribeClusterDetailResponseBody {
	s.InitVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetMaintenanceWindow(v *MaintenanceWindow) *DescribeClusterDetailResponseBody {
	s.MaintenanceWindow = v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetMasterUrl(v string) *DescribeClusterDetailResponseBody {
	s.MasterUrl = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetMetaData(v string) *DescribeClusterDetailResponseBody {
	s.MetaData = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetName(v string) *DescribeClusterDetailResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetNetworkMode(v string) *DescribeClusterDetailResponseBody {
	s.NetworkMode = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetNextVersion(v string) *DescribeClusterDetailResponseBody {
	s.NextVersion = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetParameters(v map[string]*string) *DescribeClusterDetailResponseBody {
	s.Parameters = v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetPrivateZone(v bool) *DescribeClusterDetailResponseBody {
	s.PrivateZone = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetProfile(v string) *DescribeClusterDetailResponseBody {
	s.Profile = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetRegionId(v string) *DescribeClusterDetailResponseBody {
	s.RegionId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetResourceGroupId(v string) *DescribeClusterDetailResponseBody {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetSecurityGroupId(v string) *DescribeClusterDetailResponseBody {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetSize(v int64) *DescribeClusterDetailResponseBody {
	s.Size = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetState(v string) *DescribeClusterDetailResponseBody {
	s.State = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetSubnetCidr(v string) *DescribeClusterDetailResponseBody {
	s.SubnetCidr = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetTags(v []*Tag) *DescribeClusterDetailResponseBody {
	s.Tags = v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetUpdated(v string) *DescribeClusterDetailResponseBody {
	s.Updated = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetVpcId(v string) *DescribeClusterDetailResponseBody {
	s.VpcId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetVswitchId(v string) *DescribeClusterDetailResponseBody {
	s.VswitchId = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetWorkerRamRoleName(v string) *DescribeClusterDetailResponseBody {
	s.WorkerRamRoleName = &v
	return s
}

func (s *DescribeClusterDetailResponseBody) SetZoneId(v string) *DescribeClusterDetailResponseBody {
	s.ZoneId = &v
	return s
}

type DescribeClusterDetailResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterDetailResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterDetailResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterDetailResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterDetailResponse) SetHeaders(v map[string]*string) *DescribeClusterDetailResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterDetailResponse) SetStatusCode(v int32) *DescribeClusterDetailResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterDetailResponse) SetBody(v *DescribeClusterDetailResponseBody) *DescribeClusterDetailResponse {
	s.Body = v
	return s
}

type DescribeClusterEventsRequest struct {
	PageNumber *int64  `json:"page_number,omitempty" xml:"page_number,omitempty"`
	PageSize   *int64  `json:"page_size,omitempty" xml:"page_size,omitempty"`
	TaskId     *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s DescribeClusterEventsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterEventsRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterEventsRequest) SetPageNumber(v int64) *DescribeClusterEventsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeClusterEventsRequest) SetPageSize(v int64) *DescribeClusterEventsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeClusterEventsRequest) SetTaskId(v string) *DescribeClusterEventsRequest {
	s.TaskId = &v
	return s
}

type DescribeClusterEventsResponseBody struct {
	Events   []*DescribeClusterEventsResponseBodyEvents `json:"events,omitempty" xml:"events,omitempty" type:"Repeated"`
	PageInfo *DescribeClusterEventsResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
}

func (s DescribeClusterEventsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterEventsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterEventsResponseBody) SetEvents(v []*DescribeClusterEventsResponseBodyEvents) *DescribeClusterEventsResponseBody {
	s.Events = v
	return s
}

func (s *DescribeClusterEventsResponseBody) SetPageInfo(v *DescribeClusterEventsResponseBodyPageInfo) *DescribeClusterEventsResponseBody {
	s.PageInfo = v
	return s
}

type DescribeClusterEventsResponseBodyEvents struct {
	ClusterId *string                                      `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	Data      *DescribeClusterEventsResponseBodyEventsData `json:"data,omitempty" xml:"data,omitempty" type:"Struct"`
	EventId   *string                                      `json:"event_id,omitempty" xml:"event_id,omitempty"`
	Source    *string                                      `json:"source,omitempty" xml:"source,omitempty"`
	Subject   *string                                      `json:"subject,omitempty" xml:"subject,omitempty"`
	Time      *string                                      `json:"time,omitempty" xml:"time,omitempty"`
	Type      *string                                      `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeClusterEventsResponseBodyEvents) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterEventsResponseBodyEvents) GoString() string {
	return s.String()
}

func (s *DescribeClusterEventsResponseBodyEvents) SetClusterId(v string) *DescribeClusterEventsResponseBodyEvents {
	s.ClusterId = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyEvents) SetData(v *DescribeClusterEventsResponseBodyEventsData) *DescribeClusterEventsResponseBodyEvents {
	s.Data = v
	return s
}

func (s *DescribeClusterEventsResponseBodyEvents) SetEventId(v string) *DescribeClusterEventsResponseBodyEvents {
	s.EventId = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyEvents) SetSource(v string) *DescribeClusterEventsResponseBodyEvents {
	s.Source = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyEvents) SetSubject(v string) *DescribeClusterEventsResponseBodyEvents {
	s.Subject = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyEvents) SetTime(v string) *DescribeClusterEventsResponseBodyEvents {
	s.Time = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyEvents) SetType(v string) *DescribeClusterEventsResponseBodyEvents {
	s.Type = &v
	return s
}

type DescribeClusterEventsResponseBodyEventsData struct {
	Level   *string `json:"level,omitempty" xml:"level,omitempty"`
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
	Reason  *string `json:"reason,omitempty" xml:"reason,omitempty"`
}

func (s DescribeClusterEventsResponseBodyEventsData) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterEventsResponseBodyEventsData) GoString() string {
	return s.String()
}

func (s *DescribeClusterEventsResponseBodyEventsData) SetLevel(v string) *DescribeClusterEventsResponseBodyEventsData {
	s.Level = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyEventsData) SetMessage(v string) *DescribeClusterEventsResponseBodyEventsData {
	s.Message = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyEventsData) SetReason(v string) *DescribeClusterEventsResponseBodyEventsData {
	s.Reason = &v
	return s
}

type DescribeClusterEventsResponseBodyPageInfo struct {
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	PageSize   *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	TotalCount *int64 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeClusterEventsResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterEventsResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeClusterEventsResponseBodyPageInfo) SetPageNumber(v int64) *DescribeClusterEventsResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyPageInfo) SetPageSize(v int64) *DescribeClusterEventsResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeClusterEventsResponseBodyPageInfo) SetTotalCount(v int64) *DescribeClusterEventsResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeClusterEventsResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterEventsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterEventsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterEventsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterEventsResponse) SetHeaders(v map[string]*string) *DescribeClusterEventsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterEventsResponse) SetStatusCode(v int32) *DescribeClusterEventsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterEventsResponse) SetBody(v *DescribeClusterEventsResponseBody) *DescribeClusterEventsResponse {
	s.Body = v
	return s
}

type DescribeClusterLogsResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribeClusterLogsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeClusterLogsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterLogsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterLogsResponse) SetHeaders(v map[string]*string) *DescribeClusterLogsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterLogsResponse) SetStatusCode(v int32) *DescribeClusterLogsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterLogsResponse) SetBody(v []*DescribeClusterLogsResponseBody) *DescribeClusterLogsResponse {
	s.Body = v
	return s
}

type DescribeClusterLogsResponseBody struct {
	ID         *int64  `json:"ID,omitempty" xml:"ID,omitempty"`
	ClusterId  *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	ClusterLog *string `json:"cluster_log,omitempty" xml:"cluster_log,omitempty"`
	Created    *string `json:"created,omitempty" xml:"created,omitempty"`
	Updated    *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeClusterLogsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterLogsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterLogsResponseBody) SetID(v int64) *DescribeClusterLogsResponseBody {
	s.ID = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetClusterId(v string) *DescribeClusterLogsResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetClusterLog(v string) *DescribeClusterLogsResponseBody {
	s.ClusterLog = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetCreated(v string) *DescribeClusterLogsResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClusterLogsResponseBody) SetUpdated(v string) *DescribeClusterLogsResponseBody {
	s.Updated = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBody struct {
	// The auto scaling configurations of the queried node pool.
	AutoScaling        *DescribeClusterNodePoolDetailResponseBodyAutoScaling        `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	InterconnectConfig *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig `json:"interconnect_config,omitempty" xml:"interconnect_config,omitempty" type:"Struct"`
	// The network type of the edge node pool. Valid values: basic and enhanced. This parameter takes effect only for edge node pools.
	InterconnectMode *string `json:"interconnect_mode,omitempty" xml:"interconnect_mode,omitempty"`
	// The configurations of the cluster where the node pool is deployed.
	KubernetesConfig *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	// The configurations about the managed node pool feature.
	Management *DescribeClusterNodePoolDetailResponseBodyManagement `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
	// The maximum number of nodes that are supported by the edge node pool. The value of this parameter must be equal to or greater than 0. A value of 0 indicates that the number of nodes in the node pool is limited only by the quota of nodes in the cluster. In most cases, this parameter is set to a value larger than 0 for edge node pools. This parameter is set to 0 for node pools of the ess type or default edge node pools.
	MaxNodes *int64 `json:"max_nodes,omitempty" xml:"max_nodes,omitempty"`
	// The configurations of the node pool.
	NodepoolInfo *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	// The configurations of the scaling group.
	ScalingGroup *DescribeClusterNodePoolDetailResponseBodyScalingGroup `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	Status       *DescribeClusterNodePoolDetailResponseBodyStatus       `json:"status,omitempty" xml:"status,omitempty" type:"Struct"`
	// The configurations of confidential computing.
	TeeConfig *DescribeClusterNodePoolDetailResponseBodyTeeConfig `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolDetailResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetAutoScaling(v *DescribeClusterNodePoolDetailResponseBodyAutoScaling) *DescribeClusterNodePoolDetailResponseBody {
	s.AutoScaling = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetInterconnectConfig(v *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) *DescribeClusterNodePoolDetailResponseBody {
	s.InterconnectConfig = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetInterconnectMode(v string) *DescribeClusterNodePoolDetailResponseBody {
	s.InterconnectMode = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetKubernetesConfig(v *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) *DescribeClusterNodePoolDetailResponseBody {
	s.KubernetesConfig = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetManagement(v *DescribeClusterNodePoolDetailResponseBodyManagement) *DescribeClusterNodePoolDetailResponseBody {
	s.Management = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetMaxNodes(v int64) *DescribeClusterNodePoolDetailResponseBody {
	s.MaxNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetNodepoolInfo(v *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) *DescribeClusterNodePoolDetailResponseBody {
	s.NodepoolInfo = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetScalingGroup(v *DescribeClusterNodePoolDetailResponseBodyScalingGroup) *DescribeClusterNodePoolDetailResponseBody {
	s.ScalingGroup = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetStatus(v *DescribeClusterNodePoolDetailResponseBodyStatus) *DescribeClusterNodePoolDetailResponseBody {
	s.Status = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBody) SetTeeConfig(v *DescribeClusterNodePoolDetailResponseBodyTeeConfig) *DescribeClusterNodePoolDetailResponseBody {
	s.TeeConfig = v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyAutoScaling struct {
	// The peak bandwidth of the elastic IP address (EIP) that is associated with the node pool.
	EipBandwidth *int64 `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	// The billing method of the EIP. Valid values:
	//
	// *   `PayByBandwidth`: pay-by-bandwidth
	// *   `PayByTraffic`: pay-by-data-transfer
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	// Indicates whether auto scaling is enabled. Valid values:
	//
	// *   `true`: Auto scaling is enabled.
	// *   `false`: Auto scaling is disabled. If this parameter is set to false, other parameters in the `auto_scaling` section do not take effect.
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// Indicates whether an EIP is associated with the node pool. Valid values:
	//
	// *   `true`: An EIP is associated with the node pool.
	// *   `false`: No EIP is associated with the node pool.
	IsBondEip *bool `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	// The maximum number of Elastic Compute Service (ECS) instances supported by the node pool.
	MaxInstances *int64 `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	// The minimum number of ECS instances that must be kept in the node pool.
	MinInstances *int64 `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	// The instance types that can be used for the auto scaling of the node pool. Valid values:
	//
	// *   `cpu`: regular instance
	// *   `gpu`: GPU-accelerated instance
	// *   `gpushare`: shared GPU-accelerated instance
	// *   `spot`: preemptible instance
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyAutoScaling) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetEipBandwidth(v int64) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetEipInternetChargeType(v string) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetEnable(v bool) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetIsBondEip(v bool) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetMaxInstances(v int64) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetMinInstances(v int64) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyAutoScaling) SetType(v string) *DescribeClusterNodePoolDetailResponseBodyAutoScaling {
	s.Type = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyInterconnectConfig struct {
	Bandwidth      *int64  `json:"bandwidth,omitempty" xml:"bandwidth,omitempty"`
	CcnId          *string `json:"ccn_id,omitempty" xml:"ccn_id,omitempty"`
	CcnRegionId    *string `json:"ccn_region_id,omitempty" xml:"ccn_region_id,omitempty"`
	CenId          *string `json:"cen_id,omitempty" xml:"cen_id,omitempty"`
	ImprovedPeriod *string `json:"improved_period,omitempty" xml:"improved_period,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) SetBandwidth(v int64) *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig {
	s.Bandwidth = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) SetCcnId(v string) *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig {
	s.CcnId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) SetCcnRegionId(v string) *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig {
	s.CcnRegionId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) SetCenId(v string) *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig {
	s.CenId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig) SetImprovedPeriod(v string) *DescribeClusterNodePoolDetailResponseBodyInterconnectConfig {
	s.ImprovedPeriod = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyKubernetesConfig struct {
	// Indicates whether the CloudMonitor agent is installed on ECS nodes in the cluster. After the CloudMonitor agent is installed, you can view monitoring information about the ECS instances in the CloudMonitor console. Installation is recommended. Valid values:
	//
	// *   `true`: The CloudMonitor agent is installed on ECS nodes.
	// *   `false`: The CloudMonitor agent is not installed on ECS nodes.
	CmsEnabled *bool `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	// The CPU management policy of the nodes in the node pool. The following policies are supported if the Kubernetes version of the cluster is 1.12.6 or later.
	//
	// *   `static`: allows pods with specific resource characteristics on the node to be granted enhanced CPU affinity and exclusivity.
	// *   `none`: indicates that the default CPU affinity is used.
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// The labels of the nodes in the node pool. You can add labels to the nodes in the cluster. You must add labels based on the following rules:
	//
	// *   Each label is a case-sensitive key-value pair. You can add up to 20 labels.
	// *   A key must be unique and cannot exceed 64 characters in length. A value can be empty and cannot exceed 128 characters in length. Keys and values cannot start with `aliyun`, `acs:`, `https://`, or `http://`. For more information, see [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set).
	Labels []*Tag `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	// A custom node name consists of a prefix, an IP substring, and a suffix.
	//
	// *   The prefix and suffix can contain multiple parts that are separated by periods (.). Each part can contain lowercase letters, digits, and hyphens (-). A custom node name must start and end with a digit or lowercase letter.
	// *   The IP substring length specifies the number of digits to be truncated from the end of the node IP address. The IP substring length ranges from 5 to 12.
	//
	// For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, the IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test.
	NodeNameMode *string `json:"node_name_mode,omitempty" xml:"node_name_mode,omitempty"`
	// The name of the container runtime.
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The version of the container runtime.
	RuntimeVersion *string `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	// The taints of the nodes. Taints are added to nodes to prevent pods from being scheduled to inappropriate nodes. However, toleration rules allow pods to be scheduled to nodes with matching taints. For more information, see [taint-and-toleration](https://kubernetes.io/zh/docs/concepts/scheduling-eviction/taint-and-toleration/).
	Taints []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// The user-defined data of the node pool. For more information, see [Generate user-defined data](~~49121~~).
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetCmsEnabled(v bool) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetCpuPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetLabels(v []*Tag) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.Labels = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetNodeNameMode(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.NodeNameMode = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetRuntime(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetRuntimeVersion(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetTaints(v []*Taint) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.Taints = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig) SetUserData(v string) *DescribeClusterNodePoolDetailResponseBodyKubernetesConfig {
	s.UserData = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyManagement struct {
	// Indicates whether enable auto repair is enabled. This parameter takes effect only when `enable=true` is specified.
	//
	// *   `true`: Auto repair is enabled.
	// *   `false`: Auto repair is disabled.
	AutoRepair *bool `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	// Indicates whether to enable the managed node pool feature is enabled. Valid values:
	//
	// *   `true`: The managed node pool feature is enabled.
	// *   `false`: The managed node pool feature is disabled. Other parameters in this section take effect only when `enable=true` is specified.
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// The configurations of auto upgrade. The configurations take effect only when `enable=true` is specified.
	UpgradeConfig *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolDetailResponseBodyManagement) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyManagement) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagement) SetAutoRepair(v bool) *DescribeClusterNodePoolDetailResponseBodyManagement {
	s.AutoRepair = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagement) SetEnable(v bool) *DescribeClusterNodePoolDetailResponseBodyManagement {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagement) SetUpgradeConfig(v *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) *DescribeClusterNodePoolDetailResponseBodyManagement {
	s.UpgradeConfig = v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig struct {
	// Indicates whether auto upgrade is enabled. Valid values:
	//
	// *   `true`: Auto upgrade is enabled.
	// *   `false`: Auto upgrade is disabled.
	AutoUpgrade *bool `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	// The maximum number of nodes that can be in the Unavailable state. Valid values: 1 to 1000.
	//
	// Default value: 1
	MaxUnavailable *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
	// The number of nodes that are temporarily added to the node pool during an auto upgrade.
	Surge *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	// The percentage of temporary nodes to the nodes in the node pool. You must set this parameter or `surge`.
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetAutoUpgrade(v bool) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetMaxUnavailable(v int64) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetSurge(v int64) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig) SetSurgePercentage(v int64) *DescribeClusterNodePoolDetailResponseBodyManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyNodepoolInfo struct {
	// The time when the node pool was created.
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// Indicates whether the node pool is a default node pool. A Container Service for Kubernetes (ACK) cluster usually has only one default node pool. Valid values:
	//
	// `true`: The node pool is a default node pool.
	//
	// `false`: The node pool is not a default node pool.
	IsDefault *bool `json:"is_default,omitempty" xml:"is_default,omitempty"`
	// The name of the node pool.
	//
	// The name must be 1 to 63 characters in length, and can contain digits, letters, and hyphens (-). It cannot start with a hyphen (-).
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// The ID of the node pool.
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// The ID of the region where the node pool is deployed.
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// The ID of the resource group to which the node pool belongs.
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// The type of the node pool.
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// The time when the node pool was last updated.
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetCreated(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Created = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetIsDefault(v bool) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.IsDefault = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetName(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Name = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetNodepoolId(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetRegionId(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.RegionId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetResourceGroupId(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetType(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Type = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo) SetUpdated(v string) *DescribeClusterNodePoolDetailResponseBodyNodepoolInfo {
	s.Updated = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyScalingGroup struct {
	// Indicates whether auto-renewal is enabled for the nodes in the node pool. This parameter takes effect only when `instance_charge_type` is set to `PrePaid`. Valid values:
	//
	// *   `true`: Auto-renewal is enabled.
	// *   `false`: Auto-renewal is disabled.
	AutoRenew *bool `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	// The duration of the auto-renewal. This parameter takes effect and is required only when `instance_charge_type` is set to `PrePaid`.
	//
	// If you specify `PeriodUnit=Month`, the valid values are 1, 2, 3, 6, and 12.
	AutoRenewPeriod *int64 `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	// Indicates whether pay-as-you-go instances are automatically created to meet the required number of ECS instances if preemptible instances cannot be created due to reasons such as cost or insufficient inventory. This parameter takes effect when `multi_az_policy` is set to `COST_OPTIMIZED`. Valid values:
	//
	// *   `true`: Pay-as-you-go instances are automatically created to meet the required number of ECS instances if preemptible instances cannot be created.
	// *   `false`: Pay-as-you-go instances are not automatically created to meet the required number of ECS instances if preemptible instances cannot be created.
	CompensateWithOnDemand *bool `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	// The configurations of the data disks that are attached to the nodes in the node pool. The configurations include the disk type and disk size.
	DataDisks []*DataDisk `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	// The ID of the deployment set to which the ECS instances in the node pool belong.
	DeploymentsetId *string `json:"deploymentset_id,omitempty" xml:"deploymentset_id,omitempty"`
	// The expected number of nodes in the node pool.
	DesiredSize *int64 `json:"desired_size,omitempty" xml:"desired_size,omitempty"`
	// The ID of the custom image. You can call the `DescribeKubernetesVersionMetadata` operation to query the images supported by ACK.
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// The billing method of the nodes in the node pool. Valid values:
	//
	// *   `PrePaid`: subscription
	// *   `PostPaid`: pay-as-you-go
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// The instance types of the nodes in the node pool.
	InstanceTypes []*string `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	// The billing method of the public IP address of the node.
	InternetChargeType *string `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	// The maximum outbound bandwidth of the public IP address of the node. Unit: Mbit/s. Valid values: 1 to 100.
	InternetMaxBandwidthOut *int64 `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
	// The name of the key pair. You must set this parameter or the `login_password` parameter. You must set `key_pair` if the node pool is a managed node pool.
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// The password for SSH logon. You must set this parameter or the `key_pair` parameter. The password must be 8 to 30 characters in length, and must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters.
	//
	// For security purposes, the returned password is encrypted.
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// The ECS instance scaling policy for a multi-zone scaling group. Valid values:
	//
	// *   `PRIORITY`: the scaling group is scaled based on the VSwitchIds.N parameter. If an ECS instance cannot be created in the zone where the vSwitch that has the highest priority resides, Auto Scaling creates the ECS instance in the zone where the vSwitch that has the next highest priority resides.
	//
	// *   `COST_OPTIMIZED`: ECS instances are created based on the vCPU unit price in ascending order. Preemptible instances are preferably created when preemptible instance types are specified in the scaling configuration. You can set the `CompensateWithOnDemand` parameter to specify whether to automatically create pay-as-you-go instances when preemptible instances cannot be created due to insufficient resources.
	//
	//     **
	//
	//     **Note** `COST_OPTIMIZED` is valid only when multiple instance types are specified or at least one preemptible instance type is specified.
	//
	// *   `BALANCE`: ECS instances are evenly distributed across multiple zones specified by the scaling group. If ECS instances become imbalanced among multiple zones due to insufficient inventory, you can call the RebalanceInstances operation of Auto Scaling to balance the instance distribution among zones. For more information, see [RebalanceInstances](~~71516~~)
	//
	// Default value: `PRIORITY`
	MultiAzPolicy *string `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	// The minimum number of pay-as-you-go instances that must be kept in the scaling group. Valid values: 0 to 1000. If the number of pay-as-you-go instances is less than the value of this parameter, Auto Scaling preferably creates pay-as-you-go instances.
	OnDemandBaseCapacity *int64 `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	// The percentage of pay-as-you-go instances among the extra instances that exceed the number specified by `on_demand_base_capacity`. Valid values: 0 to 100.
	OnDemandPercentageAboveBaseCapacity *int64 `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	// The subscription duration of worker nodes. This parameter takes effect and is required only when `instance_charge_type` is set to `PrePaid`.
	//
	// If `PeriodUnit=Month` is specified, the valid values are 1, 2, 3, 6, 12, 24, 36, 48, and 60.
	Period *int64 `json:"period,omitempty" xml:"period,omitempty"`
	// The billing cycle of the nodes. This parameter is required if `instance_charge_type` is set to `PrePaid`.
	//
	// Valid value: `Month`
	PeriodUnit *string `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	// The release version of the operating system. Valid values:
	//
	// *   `CentOS`
	// *   `AliyunLinux`
	// *   `Windows`
	// *   `WindowsCore`
	Platform *string `json:"platform,omitempty" xml:"platform,omitempty"`
	// The configurations of the private node pool.
	PrivatePoolOptions *DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions `json:"private_pool_options,omitempty" xml:"private_pool_options,omitempty" type:"Struct"`
	// The name of the worker Resource Access Management (RAM) role. The RAM role is assigned to the worker nodes of the cluster to allow the worker nodes to manage ECS instances.
	RamPolicy *string `json:"ram_policy,omitempty" xml:"ram_policy,omitempty"`
	// The IDs of the ApsaraDB RDS instances.
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// The ID of the scaling group.
	ScalingGroupId *string `json:"scaling_group_id,omitempty" xml:"scaling_group_id,omitempty"`
	// The scaling mode of the scaling group. Valid values:
	//
	// *   `release`: the standard mode. ECS instances are created and released based on the resource usage.
	// *   `recycle`: the swift mode. ECS instances are created, stopped, or started during scaling events. This reduces the time required for the next scale-out event. When the instance is stopped, you are charged only for the storage service. This does not apply to ECS instances attached with local disks.
	ScalingPolicy *string `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	// The ID of the security group to which the node pool is added. If the node pool is added to multiple security groups, the first ID in the value of `security_group_ids` is returned.
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// The IDs of the security groups to which the node pool is added.
	SecurityGroupIds []*string `json:"security_group_ids,omitempty" xml:"security_group_ids,omitempty" type:"Repeated"`
	// The number of instance types that are available for creating preemptible instances. Auto Scaling creates preemptible instances of multiple instance types that are available at the lowest cost. Valid values: 1 to 10.
	SpotInstancePools *int64 `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	// Indicates whether preemptible instances are supplemented when the number of preemptible instances drops below the specified minimum number. If this parameter is set to true, when the scaling group receives a system message that a preemptible instance is to be reclaimed, the scaling group attempts to create a new instance to replace this instance. Valid values: Valid values:
	//
	// *   `true`: Supplementation of preemptible instances is enabled.
	// *   `false`: Supplementation of preemptible instances is disabled.
	SpotInstanceRemedy *bool `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	// The bid configurations of preemptible instances.
	SpotPriceLimit []*DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	// The bidding policy of preemptible instances. Valid values:
	//
	// *   NoSpot: a non-preemptible instance.
	// *   SpotWithPriceLimit: a preemptible instance that is configured with the highest bid price.
	// *   SpotAsPriceGo: a preemptible instance for which the system automatically bids based on the current market price.
	//
	// For more information, see [Preemptible instances](~~157759~~).
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	// The type of system disk. Valid values:
	//
	// *   `cloud_efficiency`: ultra disk
	// *   `cloud_ssd`: standard SSD
	SystemDiskCategory *string `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	// The performance level (PL) of the system disk that you want to use for the node. This parameter takes effect only for enhanced SSDs (ESSDs).
	SystemDiskPerformanceLevel *string `json:"system_disk_performance_level,omitempty" xml:"system_disk_performance_level,omitempty"`
	// The system disk size of a node. Unit: GiB.
	//
	// Valid values: 20 to 500
	SystemDiskSize *int64 `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	// The labels that you want to add to the ECS instances.
	//
	// A key must be unique and cannot exceed 128 characters in length. Neither keys nor values can start with aliyun or acs:. Neither keys nor values can contain https:// or http://.
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// The IDs of vSwitches.
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroup) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetAutoRenew(v bool) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetAutoRenewPeriod(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetCompensateWithOnDemand(v bool) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetDataDisks(v []*DataDisk) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.DataDisks = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetDeploymentsetId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.DeploymentsetId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetDesiredSize(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.DesiredSize = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetImageId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.ImageId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInstanceChargeType(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInstanceTypes(v []*string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInternetChargeType(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetInternetMaxBandwidthOut(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetKeyPair(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetLoginPassword(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetMultiAzPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetOnDemandBaseCapacity(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetPeriod(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.Period = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetPeriodUnit(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetPlatform(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.Platform = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetPrivatePoolOptions(v *DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.PrivatePoolOptions = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetRamPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.RamPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetRdsInstances(v []*string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetScalingGroupId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.ScalingGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetScalingPolicy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSecurityGroupId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSecurityGroupIds(v []*string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SecurityGroupIds = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotInstancePools(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotInstanceRemedy(v bool) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotPriceLimit(v []*DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSpotStrategy(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSystemDiskCategory(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSystemDiskPerformanceLevel(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SystemDiskPerformanceLevel = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetSystemDiskSize(v int64) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetTags(v []*Tag) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.Tags = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroup) SetVswitchIds(v []*string) *DescribeClusterNodePoolDetailResponseBodyScalingGroup {
	s.VswitchIds = v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions struct {
	// The ID of the private node pool.
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// The type of private node pool. This parameter specifies the type of the private pool that you want to use to create instances. A private pool is generated when an elasticity assurance or a capacity reservation takes effect. You can select a private pool to start instances. Valid values:
	//
	// *   `Open`: open private pool. The system selects an open private pool to start instances. If no matching open private pools are available, the resources in the public pool are used.
	// *   `Target`: specific private pool. The system uses the resources of the specified private pool to start instances. If the specified private pool is unavailable, instances cannot be started.
	// *   `None`: no private pool is used. The resources of private pools are not used to start instances.
	MatchCriteria *string `json:"match_criteria,omitempty" xml:"match_criteria,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions) SetId(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions {
	s.Id = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions) SetMatchCriteria(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroupPrivatePoolOptions {
	s.MatchCriteria = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit struct {
	// The instance type of preemptible instances.
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// The price limit of a preemptible instance.
	//
	// Unit: USD/hour.
	PriceLimit *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) SetInstanceType(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit) SetPriceLimit(v string) *DescribeClusterNodePoolDetailResponseBodyScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyStatus struct {
	FailedNodes   *int64  `json:"failed_nodes,omitempty" xml:"failed_nodes,omitempty"`
	HealthyNodes  *int64  `json:"healthy_nodes,omitempty" xml:"healthy_nodes,omitempty"`
	InitialNodes  *int64  `json:"initial_nodes,omitempty" xml:"initial_nodes,omitempty"`
	OfflineNodes  *int64  `json:"offline_nodes,omitempty" xml:"offline_nodes,omitempty"`
	RemovingNodes *int64  `json:"removing_nodes,omitempty" xml:"removing_nodes,omitempty"`
	ServingNodes  *int64  `json:"serving_nodes,omitempty" xml:"serving_nodes,omitempty"`
	State         *string `json:"state,omitempty" xml:"state,omitempty"`
	TotalNodes    *int64  `json:"total_nodes,omitempty" xml:"total_nodes,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyStatus) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyStatus) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetFailedNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.FailedNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetHealthyNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.HealthyNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetInitialNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.InitialNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetOfflineNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.OfflineNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetRemovingNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.RemovingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetServingNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.ServingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetState(v string) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.State = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponseBodyStatus) SetTotalNodes(v int64) *DescribeClusterNodePoolDetailResponseBodyStatus {
	s.TotalNodes = &v
	return s
}

type DescribeClusterNodePoolDetailResponseBodyTeeConfig struct {
	// Indicates whether confidential computing is enabled. Valid values:
	//
	// *   `true`: Confidential computing is enabled.
	// *   `false`: Confidential computing is disabled.
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s DescribeClusterNodePoolDetailResponseBodyTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponseBodyTeeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponseBodyTeeConfig) SetTeeEnable(v bool) *DescribeClusterNodePoolDetailResponseBodyTeeConfig {
	s.TeeEnable = &v
	return s
}

type DescribeClusterNodePoolDetailResponse struct {
	Headers    map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterNodePoolDetailResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterNodePoolDetailResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolDetailResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolDetailResponse) SetHeaders(v map[string]*string) *DescribeClusterNodePoolDetailResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterNodePoolDetailResponse) SetStatusCode(v int32) *DescribeClusterNodePoolDetailResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterNodePoolDetailResponse) SetBody(v *DescribeClusterNodePoolDetailResponseBody) *DescribeClusterNodePoolDetailResponse {
	s.Body = v
	return s
}

type DescribeClusterNodePoolsResponseBody struct {
	// The list of the returned node pools.
	Nodepools []*DescribeClusterNodePoolsResponseBodyNodepools `json:"nodepools,omitempty" xml:"nodepools,omitempty" type:"Repeated"`
}

func (s DescribeClusterNodePoolsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBody) SetNodepools(v []*DescribeClusterNodePoolsResponseBodyNodepools) *DescribeClusterNodePoolsResponseBody {
	s.Nodepools = v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepools struct {
	// The configurations of auto scaling.
	AutoScaling        *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling        `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	InterconnectConfig *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig `json:"interconnect_config,omitempty" xml:"interconnect_config,omitempty" type:"Struct"`
	InterconnectMode   *string                                                          `json:"interconnect_mode,omitempty" xml:"interconnect_mode,omitempty"`
	// The configurations of the cluster.
	KubernetesConfig *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	// The configurations of managed node pools. Managed node pools are available only in professional managed Kubernetes clusters.
	Management *DescribeClusterNodePoolsResponseBodyNodepoolsManagement `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
	MaxNodes   *int64                                                   `json:"max_nodes,omitempty" xml:"max_nodes,omitempty"`
	// The information about the node pool.
	NodepoolInfo *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	// The configurations of the scaling group.
	ScalingGroup *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	// The status details about the node pool.
	Status *DescribeClusterNodePoolsResponseBodyNodepoolsStatus `json:"status,omitempty" xml:"status,omitempty" type:"Struct"`
	// The configurations of confidential computing.
	TeeConfig *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepools) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepools) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetAutoScaling(v *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.AutoScaling = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetInterconnectConfig(v *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.InterconnectConfig = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetInterconnectMode(v string) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.InterconnectMode = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetKubernetesConfig(v *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.KubernetesConfig = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetManagement(v *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.Management = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetMaxNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.MaxNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetNodepoolInfo(v *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.NodepoolInfo = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetScalingGroup(v *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.ScalingGroup = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetStatus(v *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.Status = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepools) SetTeeConfig(v *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) *DescribeClusterNodePoolsResponseBodyNodepools {
	s.TeeConfig = v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling struct {
	// The peak bandwidth of the elastic IP address (EIP).
	EipBandwidth *int64 `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	// The billing method of the EIP. Valid values:
	//
	// *   `PayByBandwidth`: pay-by-bandwidth
	// *   `PayByTraffic`: pay-by-data-transfer
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	// Indicates whether auto scaling is enabled.
	//
	// *   `true`: Auto scaling is enabled for the node pool.
	// *   `false`: Auto scaling is disabled for the node pool. If you set this parameter to `false`, other parameters in the `auto_scaling` section does not take effect.
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// Indicates whether an EIP is associated with the node pool. Valid values:
	//
	// *   `true`: An EIP is associated with the node pool.
	// *   `false`: No EIP is associated with the node pool.
	IsBondEip *bool `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	// The maximum number of Elastic Compute Service (ECS) instances supported by the node pool.
	MaxInstances *int64 `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	// The minimum number of ECS instances.
	MinInstances *int64 `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	// The minimum number of ECS instances that must be kept in the node pool. Valid values:
	//
	// *   `cpu`: regular instance
	// *   `gpu`: GPU-accelerated instance
	// *   `gpushare`: shared GPU-accelerated instance
	// *   `spot`: preemptible instance
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetEipBandwidth(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetEipInternetChargeType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetEnable(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetIsBondEip(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetMaxInstances(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetMinInstances(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling) SetType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsAutoScaling {
	s.Type = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig struct {
	Bandwidth      *int64  `json:"bandwidth,omitempty" xml:"bandwidth,omitempty"`
	CcnId          *string `json:"ccn_id,omitempty" xml:"ccn_id,omitempty"`
	CcnRegionId    *string `json:"ccn_region_id,omitempty" xml:"ccn_region_id,omitempty"`
	CenId          *string `json:"cen_id,omitempty" xml:"cen_id,omitempty"`
	ImprovedPeriod *string `json:"improved_period,omitempty" xml:"improved_period,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) SetBandwidth(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig {
	s.Bandwidth = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) SetCcnId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig {
	s.CcnId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) SetCcnRegionId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig {
	s.CcnRegionId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) SetCenId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig {
	s.CenId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig) SetImprovedPeriod(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsInterconnectConfig {
	s.ImprovedPeriod = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig struct {
	// Indicates where the CloudMonitor agent is installed on ECS nodes of the cluster. After the CloudMonitor agent is installed, you can view monitoring information about the ECS instances in the CloudMonitor console. Installation is recommended. Valid values:
	//
	// *   `true` The CloudMonitor agent is installed on ECS nodes.
	// *   `false`: The CloudMonitor agent is not installed on ECS nodes.
	CmsEnabled *bool `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	// The CPU management policy. The following policies are supported if the Kubernetes version of the cluster is 1.12.6 or later.
	//
	// *   `static`: This policy allows pods with specific resource characteristics on the node to be granted with enhanced CPU affinity and exclusivity.
	// *   `none`: indicates that the default CPU affinity is used.
	CpuPolicy *string `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	// The labels of the nodes. You can add labels to the nodes in the cluster. You must add labels based on the following rules:
	//
	// *   Each label is a case-sensitive key-value pair. You can add up to 20 labels.
	// *   A key must be unique and cannot exceed 64 characters in length. A value can be empty and cannot exceed 128 characters in length. Keys and values cannot start with `aliyun`, `acs:`, `https://`, or `http://`. For more information, see [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set).
	Labels       []*Tag  `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	NodeNameMode *string `json:"node_name_mode,omitempty" xml:"node_name_mode,omitempty"`
	// The name of the container runtime.
	Runtime *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
	// The version of the container runtime.
	RuntimeVersion *string `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	// The taints that are added to nodes. Taints are added to nodes to prevent pods from being scheduled to inappropriate nodes. However, toleration rules allow pods to be scheduled to nodes with matching taints. For more information, see [taint-and-toleration](https://kubernetes.io/zh/docs/concepts/scheduling-eviction/taint-and-toleration/).
	Taints []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	// The user-defined data of the node pool. For more information, see [Generate user-defined data](~~49121~~).
	UserData *string `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetCmsEnabled(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetCpuPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetLabels(v []*Tag) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.Labels = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetNodeNameMode(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.NodeNameMode = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetRuntime(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetRuntimeVersion(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetTaints(v []*Taint) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.Taints = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig) SetUserData(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsKubernetesConfig {
	s.UserData = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsManagement struct {
	// Indicates whether auto repair is enabled. Auto repair is enabled only when `enable=true` is specified.
	//
	// *   `true`: Auto repair is enabled.
	// *   `false`: Auto repair is disabled.
	AutoRepair *bool `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	// Indicates whether managed node pools are enabled. Valid values:
	//
	// *   `true`: Managed node pools are enabled.
	// *   `false`: Managed node pools are disabled. Other parameters in this section take effect only when `enable=true` is specified.
	Enable *bool `json:"enable,omitempty" xml:"enable,omitempty"`
	// The configurations of auto upgrade. The configurations take effect only when `enable=true` is specified.
	UpgradeConfig *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagement) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagement) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) SetAutoRepair(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsManagement {
	s.AutoRepair = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) SetEnable(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsManagement {
	s.Enable = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagement) SetUpgradeConfig(v *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) *DescribeClusterNodePoolsResponseBodyNodepoolsManagement {
	s.UpgradeConfig = v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig struct {
	// Indicates whether auto upgrade is enabled. Valid values:
	//
	// *   `true`: Auto upgrade is enabled.
	// *   `true`: Auto upgrade is disabled.
	AutoUpgrade *bool `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	// The maximum number of nodes that can be in the unschedulable state. Valid values: 1 to 1000.
	//
	// Default value: 1
	MaxUnavailable *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
	// The number of nodes that are temporarily added to the node pool during an auto upgrade.
	Surge *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	// The percentage of temporary nodes to the nodes in the node pool. You must set this parameter or `surge`.
	//
	// The number of extra nodes = The percentage of extra nodes × The number of nodes in the node pool. For example, the percentage of extra nodes is set to 50% and the number of nodes in the node pool is six. The number of extra nodes will be three.
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetAutoUpgrade(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetMaxUnavailable(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetSurge(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig) SetSurgePercentage(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo struct {
	// The time when the node pool was created.
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// Indicates whether the node pool is a default node pool. An ACK cluster usually has only one default node pool. Valid values:
	//
	// *   `true`: The node pool is a default node pool.
	// *   `false`: The node pool is not a default node pool.
	IsDefault *bool `json:"is_default,omitempty" xml:"is_default,omitempty"`
	// The name of the node pool.
	//
	// The name must be 1 to 63 characters in length, and can contain digits, letters, and hyphens (-). It cannot start with a hyphen (-).
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// The ID of the node pool.
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// The ID of the region where the node pool is deployed.
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// The ID of the resource group to which the node pool belongs.
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// The type of the node pool. Valid values:
	//
	// *   `edge`: edge node pools.
	// *   `ess`: cloud node pools.
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
	// The time when the node pool was last updated.
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetCreated(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Created = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetIsDefault(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.IsDefault = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetName(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Name = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetNodepoolId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetRegionId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.RegionId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetResourceGroupId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Type = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo) SetUpdated(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsNodepoolInfo {
	s.Updated = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup struct {
	// Indicates whether auto-renewal is enabled for the nodes in the node pool. This parameter takes effect only when `instance_charge_type` is set to `PrePaid`. Valid values:
	//
	// *   `true`: Auto-renewal is enabled.
	// *   `false`: Auto-renewal is disabled.
	AutoRenew *bool `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	// The duration of the auto-renewal. This parameter takes effect and is required only when `instance_charge_type` is set to `PrePaid`.
	//
	// If `PeriodUnit=Month` is specified, the valid values are 1, 2, 3, 6, and 12.
	AutoRenewPeriod *int64 `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	// Indicates whether pay-as-you-go instances are automatically created to meet the required number of ECS instances when the preemptible instances cannot be created due to reasons such as the cost or inventory availability. This parameter takes effect when `multi_az_policy` is set to `COST_OPTIMIZED`. Valid values:
	//
	// *   `true`: Pay-as-you-go instances are automatically created to meet the required number of ECS instances if preemptible instances cannot be created.
	// *   `false`: Pay-as-you-go instances are not created to meet the required number of ECS instances if preemptible instances cannot be created.
	CompensateWithOnDemand *bool `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	// The configurations of the data disks attached to the nodes in the node pool. The configurations include the disk type and disk size.
	DataDisks       []*DataDisk `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	DeploymentsetId *string     `json:"deploymentset_id,omitempty" xml:"deploymentset_id,omitempty"`
	DesiredSize     *int64      `json:"desired_size,omitempty" xml:"desired_size,omitempty"`
	// The ID of the custom image. You can call `DescribeKubernetesVersionMetadata` to query the images supported by ACK.
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// The billing method of the nodes in the node pool. Valid values:
	//
	// *   `PrePaid`: subscription
	// *   `PostPaid`: pay-as-you-go
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// The instance types of the nodes in the node pool.
	InstanceTypes []*string `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	// The billing method of the public IP address of the node.
	InternetChargeType *string `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	// The maximum outbound bandwidth of the public IP address of the node. Unit: Mbit/s. Valid values: 1 to 100.
	InternetMaxBandwidthOut *int64 `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
	// The name of the key pair. You must set this parameter or the `login_password` parameter.
	//
	// You must set `key_pair` if the node pool is a managed node pool.
	KeyPair *string `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	// The password for SSH logon. You must set this parameter or the `key_pair` parameter. The password must be 8 to 30 characters in length, and must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters.
	//
	// For security purposes, the returned password is encrypted.
	LoginPassword *string `json:"login_password,omitempty" xml:"login_password,omitempty"`
	// The ECS instance scaling policy for a multi-zone scaling group. Valid values:
	//
	// *   `PRIORITY`: the scaling group is scaled based on the VSwitchIds.N parameter. When an ECS instance cannot be created in the zone where the vSwitch with the highest priority resides, the system uses the vSwitch with the next highest priority to create the ECS instance.
	//
	// *   `COST_OPTIMIZED`: ECS instances are created based on the vCPU unit price in ascending order. Preemptible instances are preferentially created when multiple instance types are specified in the scaling configurations. You can set the `CompensateWithOnDemand` parameter to specify whether to automatically create pay-as-you-go instances when preemptible instances cannot be created due to insufficient resources.
	//
	//     **
	//
	//     **Note** `COST_OPTIMIZED` is valid only when multiple instance types are specified or at least one preemptible instance type is specified.
	//
	// *   `BALANCE`: ECS instances are evenly distributed across multiple zones specified by the scaling group. If ECS instances become imbalanced among multiple zones due to insufficient inventory, you can call `RebalanceInstances` of Auto Scaling (ESS) to balance the instance distribution among zones. For more information, see [RebalanceInstances](~~71516~~).
	MultiAzPolicy *string `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	// The minimum number of pay-as-you-go instances that must be kept in the scaling group. Valid values: 0 to 1000. When the number of pay-as-you-go instances is lower than this value, pay-as-you-go instances are preferentially created to meet the required number.
	OnDemandBaseCapacity *int64 `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	// The percentage of pay-as-you-go instances among the extra instances that exceed the number specified by `on_demand_base_capacity`. Valid values: 0 to 100.
	OnDemandPercentageAboveBaseCapacity *int64 `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	// The subscription duration of worker nodes. This parameter takes effect and is required only when `instance_charge_type` is set to `PrePaid`.
	//
	// If `PeriodUnit=Month` is specified, the valid values are 1, 2, 3, 6, 12, 24, 36, 48, and 60.
	Period *int64 `json:"period,omitempty" xml:"period,omitempty"`
	// The billing cycle of the nodes. This parameter takes effect only when `instance_charge_type` is set to `PrePaid`.
	//
	// Valid value: `Month`
	PeriodUnit *string `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	// The release version of the operating system. Valid values:
	//
	// *   `CentOS`
	// *   `AliyunLinux`
	// *   `Windows`
	// *   `WindowsCore`
	Platform           *string                                                                      `json:"platform,omitempty" xml:"platform,omitempty"`
	PrivatePoolOptions *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions `json:"private_pool_options,omitempty" xml:"private_pool_options,omitempty" type:"Struct"`
	// The name of the worker Resource Access Management (RAM) role. The RAM role is assigned to the worker nodes of the cluster to allow the worker nodes to manage ECS instances.
	RamPolicy *string `json:"ram_policy,omitempty" xml:"ram_policy,omitempty"`
	// The IDs of the ApsaraDB RDS instances.
	RdsInstances []*string `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	// The ID of the scaling group.
	ScalingGroupId *string `json:"scaling_group_id,omitempty" xml:"scaling_group_id,omitempty"`
	// The scaling mode of the scaling group. Valid values:
	//
	// *   `release`: the standard mode. ECS instances are created and released based on the resource usage.
	// *   `recycle`: the swift mode. ECS instances are created, stopped, or started during scaling events. This reduces the time required for the next scale-out event. When the instance is stopped, you are charged only for the storage service. This does not apply to ECS instances attached with local disks.
	ScalingPolicy *string `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	// The ID of the security group to which the node pool is added. If the node pool is added to multiple security groups, the first ID in the value of `security_group_ids` is returned.
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// The IDs of the security groups to which the node pool is added.
	SecurityGroupIds []*string `json:"security_group_ids,omitempty" xml:"security_group_ids,omitempty" type:"Repeated"`
	// The number of available instance types. The scaling group creates preemptible instances of multiple instance types at the lowest cost. Valid values: 1 to 10.
	SpotInstancePools *int64 `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	// Indicates whether preemptible instances are supplemented when the number of preemptible instances drops below the specified minimum number. If this parameter is set to true, when the scaling group receives a system message that a preemptible instance is to be reclaimed, the scaling group attempts to create a new instance to replace this instance. Valid values:
	//
	// *   `true`: Supplement to preemptible instances is enabled.
	// *   `false`: Supplement to preemptible instances is disabled.
	SpotInstanceRemedy *bool `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	// The bid configurations of preemptible instances.
	SpotPriceLimit []*DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	// The bidding policy of preemptible instances. Valid values:
	//
	// *   NoSpot: non-preemptible instance.
	// *   SpotWithPriceLimit: specifies the highest bid for the preemptible instance.
	// *   SpotAsPriceGo: automatically submits bids based on the up-to-date market price.
	//
	// For more information, see [Preemptible instances](~~157759~~).
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	// The type of system disk. Valid values:
	//
	// *   `cloud_efficiency`: ultra disk
	// *   `cloud_ssd`: standard SSD
	SystemDiskCategory         *string `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	SystemDiskPerformanceLevel *string `json:"system_disk_performance_level,omitempty" xml:"system_disk_performance_level,omitempty"`
	// The system disk size of a worker node. Unit: GiB.
	//
	// Valid values: 20 to 500
	SystemDiskSize *int64 `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	// The labels that are added only to ECS instances.
	//
	// A key must be unique and cannot exceed 128 characters in length. Neither keys nor values can start with aliyun or acs:. Neither keys nor values can contain https:// or http://.
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// The IDs of vSwitches.
	VswitchIds []*string `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetAutoRenew(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetAutoRenewPeriod(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetCompensateWithOnDemand(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetDataDisks(v []*DataDisk) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.DataDisks = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetDeploymentsetId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.DeploymentsetId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetDesiredSize(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.DesiredSize = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetImageId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.ImageId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInstanceChargeType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInstanceTypes(v []*string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInternetChargeType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetInternetMaxBandwidthOut(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetKeyPair(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetLoginPassword(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetMultiAzPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetOnDemandBaseCapacity(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetPeriod(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.Period = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetPeriodUnit(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetPlatform(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.Platform = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetPrivatePoolOptions(v *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.PrivatePoolOptions = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetRamPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.RamPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetRdsInstances(v []*string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetScalingGroupId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.ScalingGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetScalingPolicy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSecurityGroupId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSecurityGroupIds(v []*string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SecurityGroupIds = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotInstancePools(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotInstanceRemedy(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotPriceLimit(v []*DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSpotStrategy(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSystemDiskCategory(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSystemDiskPerformanceLevel(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SystemDiskPerformanceLevel = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetSystemDiskSize(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetTags(v []*Tag) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.Tags = v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup) SetVswitchIds(v []*string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroup {
	s.VswitchIds = v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions struct {
	Id            *string `json:"id,omitempty" xml:"id,omitempty"`
	MatchCriteria *string `json:"match_criteria,omitempty" xml:"match_criteria,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions) SetId(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions {
	s.Id = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions) SetMatchCriteria(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupPrivatePoolOptions {
	s.MatchCriteria = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit struct {
	// The instance type for preemptible instances.
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// The price limit of a preemptible instance. Unit: USD/hour.
	PriceLimit *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) SetInstanceType(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit) SetPriceLimit(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsStatus struct {
	// The number of failed nodes.
	FailedNodes *int64 `json:"failed_nodes,omitempty" xml:"failed_nodes,omitempty"`
	// The number of healthy nodes.
	HealthyNodes *int64 `json:"healthy_nodes,omitempty" xml:"healthy_nodes,omitempty"`
	// The number of nodes that are being created.
	InitialNodes *int64 `json:"initial_nodes,omitempty" xml:"initial_nodes,omitempty"`
	// The number of offline nodes.
	OfflineNodes *int64 `json:"offline_nodes,omitempty" xml:"offline_nodes,omitempty"`
	// The number of nodes that are being removed.
	RemovingNodes *int64 `json:"removing_nodes,omitempty" xml:"removing_nodes,omitempty"`
	// The number of running nodes.
	ServingNodes *int64 `json:"serving_nodes,omitempty" xml:"serving_nodes,omitempty"`
	// The status of the node pool. Valid values:
	//
	// *   `active`: The node pool is active.
	// *   `scaling`: The node pool is being scaled.
	// *   `removing`: Nodes are being removed from the node pool.
	// *   `deleting`: The node pool is being deleted.
	// *   `updating`: The node pool is being updated.
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// The total number of nodes in the node pool.
	TotalNodes *int64 `json:"total_nodes,omitempty" xml:"total_nodes,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsStatus) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsStatus) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetFailedNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.FailedNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetHealthyNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.HealthyNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetInitialNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.InitialNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetOfflineNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.OfflineNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetRemovingNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.RemovingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetServingNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.ServingNodes = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetState(v string) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.State = &v
	return s
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsStatus) SetTotalNodes(v int64) *DescribeClusterNodePoolsResponseBodyNodepoolsStatus {
	s.TotalNodes = &v
	return s
}

type DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig struct {
	// Indicates whether confidential computing is enabled. Valid values:
	//
	// *   `true`: confidential computing is enabled.
	// *   `false`: confidential computing is disabled.
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig) SetTeeEnable(v bool) *DescribeClusterNodePoolsResponseBodyNodepoolsTeeConfig {
	s.TeeEnable = &v
	return s
}

type DescribeClusterNodePoolsResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterNodePoolsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterNodePoolsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodePoolsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodePoolsResponse) SetHeaders(v map[string]*string) *DescribeClusterNodePoolsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterNodePoolsResponse) SetStatusCode(v int32) *DescribeClusterNodePoolsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterNodePoolsResponse) SetBody(v *DescribeClusterNodePoolsResponseBody) *DescribeClusterNodePoolsResponse {
	s.Body = v
	return s
}

type DescribeClusterNodesRequest struct {
	InstanceIds *string `json:"instanceIds,omitempty" xml:"instanceIds,omitempty"`
	NodepoolId  *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	PageNumber  *string `json:"pageNumber,omitempty" xml:"pageNumber,omitempty"`
	PageSize    *string `json:"pageSize,omitempty" xml:"pageSize,omitempty"`
	State       *string `json:"state,omitempty" xml:"state,omitempty"`
}

func (s DescribeClusterNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesRequest) SetInstanceIds(v string) *DescribeClusterNodesRequest {
	s.InstanceIds = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetNodepoolId(v string) *DescribeClusterNodesRequest {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetPageNumber(v string) *DescribeClusterNodesRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetPageSize(v string) *DescribeClusterNodesRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeClusterNodesRequest) SetState(v string) *DescribeClusterNodesRequest {
	s.State = &v
	return s
}

type DescribeClusterNodesResponseBody struct {
	// The details of the nodes that are returned.
	Nodes []*DescribeClusterNodesResponseBodyNodes `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
	// The pagination details.
	Page *DescribeClusterNodesResponseBodyPage `json:"page,omitempty" xml:"page,omitempty" type:"Struct"`
}

func (s DescribeClusterNodesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponseBody) SetNodes(v []*DescribeClusterNodesResponseBodyNodes) *DescribeClusterNodesResponseBody {
	s.Nodes = v
	return s
}

func (s *DescribeClusterNodesResponseBody) SetPage(v *DescribeClusterNodesResponseBodyPage) *DescribeClusterNodesResponseBody {
	s.Page = v
	return s
}

type DescribeClusterNodesResponseBodyNodes struct {
	// The time when the node was created.
	CreationTime *string `json:"creation_time,omitempty" xml:"creation_time,omitempty"`
	// The error message that was generated when the node was created.
	ErrorMessage *string `json:"error_message,omitempty" xml:"error_message,omitempty"`
	// The expiration time of the node.
	ExpiredTime *string `json:"expired_time,omitempty" xml:"expired_time,omitempty"`
	// The name of the host.
	HostName *string `json:"host_name,omitempty" xml:"host_name,omitempty"`
	// The ID of the system image that is used by the node.
	ImageId *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	// The billing method of the instance on which the node is deployed. Valid values:
	//
	// *   `PrePaid`: the subscription billing method. If the value is PrePaid, make sure that you have a sufficient balance or credit in your account. Otherwise, an `InvalidPayMethod` error is returned.
	// *   `PostPaid`: the pay-as-you-go billing method.
	InstanceChargeType *string `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	// The ID of the instance on which the node is deployed.
	InstanceId *string `json:"instance_id,omitempty" xml:"instance_id,omitempty"`
	// The name of the instance on which the node is deployed.
	InstanceName *string `json:"instance_name,omitempty" xml:"instance_name,omitempty"`
	// The role of the node. Valid values:
	//
	// *   Master: master node
	// *   Worker: worker node
	InstanceRole *string `json:"instance_role,omitempty" xml:"instance_role,omitempty"`
	// The status of the node.
	InstanceStatus *string `json:"instance_status,omitempty" xml:"instance_status,omitempty"`
	// The instance type of the node.
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	// The Elastic Compute Service (ECS) instance family of the node.
	InstanceTypeFamily *string `json:"instance_type_family,omitempty" xml:"instance_type_family,omitempty"`
	// The IP address of the node.
	IpAddress []*string `json:"ip_address,omitempty" xml:"ip_address,omitempty" type:"Repeated"`
	// Indicates whether the instance on which the node is deployed is provided by Alibaba Cloud. Valid values:
	//
	// *   `true`: The instance is provided by Alibaba Cloud.
	// *   `false`: The instance is not provided by Alibaba Cloud.
	IsAliyunNode *bool `json:"is_aliyun_node,omitempty" xml:"is_aliyun_node,omitempty"`
	// The name of the node. This name is the identifier of the node in the cluster.
	NodeName *string `json:"node_name,omitempty" xml:"node_name,omitempty"`
	// Indicates whether the node is ready. Valid values:
	//
	// *   `Ready`: The node is ready.
	// *   `NotReady`: The node is not ready.
	// *   `Unknown`: The status of the node is unknown.
	// *   `Offline`: The node is offline.
	NodeStatus *string `json:"node_status,omitempty" xml:"node_status,omitempty"`
	// The ID of the node pool.
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// Indicates how the node is initialized. A node can be manually created or created by using Resource Orchestration Service (ROS).
	Source *string `json:"source,omitempty" xml:"source,omitempty"`
	// The type of the preemptible instance. Valid values:
	//
	// *   NoSpot: a non-preemptible instance.
	// *   SpotWithPriceLimit: a preemptible instance that is configured with the highest bid price.
	// *   SpotAsPriceGo: a preemptible instance for which the system automatically bids based on the current market price.
	SpotStrategy *string `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	// The status of the node. Valid values:
	//
	// *   `pending`: The node is being created.
	// *   `running`: The node is running.
	// *   `starting`: The node is being started.
	// *   `stopping`: The node is being stopped.
	// *   `stopped`: The node is stopped.
	State *string `json:"state,omitempty" xml:"state,omitempty"`
}

func (s DescribeClusterNodesResponseBodyNodes) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponseBodyNodes) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponseBodyNodes) SetCreationTime(v string) *DescribeClusterNodesResponseBodyNodes {
	s.CreationTime = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetErrorMessage(v string) *DescribeClusterNodesResponseBodyNodes {
	s.ErrorMessage = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetExpiredTime(v string) *DescribeClusterNodesResponseBodyNodes {
	s.ExpiredTime = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetHostName(v string) *DescribeClusterNodesResponseBodyNodes {
	s.HostName = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetImageId(v string) *DescribeClusterNodesResponseBodyNodes {
	s.ImageId = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceChargeType(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceChargeType = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceId(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceId = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceName(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceName = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceRole(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceRole = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceStatus(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceStatus = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceType(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceType = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetInstanceTypeFamily(v string) *DescribeClusterNodesResponseBodyNodes {
	s.InstanceTypeFamily = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetIpAddress(v []*string) *DescribeClusterNodesResponseBodyNodes {
	s.IpAddress = v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetIsAliyunNode(v bool) *DescribeClusterNodesResponseBodyNodes {
	s.IsAliyunNode = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetNodeName(v string) *DescribeClusterNodesResponseBodyNodes {
	s.NodeName = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetNodeStatus(v string) *DescribeClusterNodesResponseBodyNodes {
	s.NodeStatus = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetNodepoolId(v string) *DescribeClusterNodesResponseBodyNodes {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetSource(v string) *DescribeClusterNodesResponseBodyNodes {
	s.Source = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetSpotStrategy(v string) *DescribeClusterNodesResponseBodyNodes {
	s.SpotStrategy = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyNodes) SetState(v string) *DescribeClusterNodesResponseBodyNodes {
	s.State = &v
	return s
}

type DescribeClusterNodesResponseBodyPage struct {
	// The page number of the returned page.
	PageNumber *int32 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// The number of entries returned per page.
	PageSize *int32 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// The total number of entries returned.
	TotalCount *int32 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeClusterNodesResponseBodyPage) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponseBodyPage) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponseBodyPage) SetPageNumber(v int32) *DescribeClusterNodesResponseBodyPage {
	s.PageNumber = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyPage) SetPageSize(v int32) *DescribeClusterNodesResponseBodyPage {
	s.PageSize = &v
	return s
}

func (s *DescribeClusterNodesResponseBodyPage) SetTotalCount(v int32) *DescribeClusterNodesResponseBodyPage {
	s.TotalCount = &v
	return s
}

type DescribeClusterNodesResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterNodesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterNodesResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterNodesResponse) SetHeaders(v map[string]*string) *DescribeClusterNodesResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterNodesResponse) SetStatusCode(v int32) *DescribeClusterNodesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterNodesResponse) SetBody(v *DescribeClusterNodesResponseBody) *DescribeClusterNodesResponse {
	s.Body = v
	return s
}

type DescribeClusterResourcesResponse struct {
	Headers    map[string]*string                      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribeClusterResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeClusterResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterResourcesResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterResourcesResponse) SetHeaders(v map[string]*string) *DescribeClusterResourcesResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterResourcesResponse) SetStatusCode(v int32) *DescribeClusterResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterResourcesResponse) SetBody(v []*DescribeClusterResourcesResponseBody) *DescribeClusterResourcesResponse {
	s.Body = v
	return s
}

type DescribeClusterResourcesResponseBody struct {
	ClusterId    *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	Created      *string `json:"created,omitempty" xml:"created,omitempty"`
	InstanceId   *string `json:"instance_id,omitempty" xml:"instance_id,omitempty"`
	ResourceInfo *string `json:"resource_info,omitempty" xml:"resource_info,omitempty"`
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	State        *string `json:"state,omitempty" xml:"state,omitempty"`
	AutoCreate   *int64  `json:"auto_create,omitempty" xml:"auto_create,omitempty"`
}

func (s DescribeClusterResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterResourcesResponseBody) SetClusterId(v string) *DescribeClusterResourcesResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetCreated(v string) *DescribeClusterResourcesResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetInstanceId(v string) *DescribeClusterResourcesResponseBody {
	s.InstanceId = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetResourceInfo(v string) *DescribeClusterResourcesResponseBody {
	s.ResourceInfo = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetResourceType(v string) *DescribeClusterResourcesResponseBody {
	s.ResourceType = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetState(v string) *DescribeClusterResourcesResponseBody {
	s.State = &v
	return s
}

func (s *DescribeClusterResourcesResponseBody) SetAutoCreate(v int64) *DescribeClusterResourcesResponseBody {
	s.AutoCreate = &v
	return s
}

type DescribeClusterTasksResponseBody struct {
	PageInfo  *DescribeClusterTasksResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
	RequestId *string                                   `json:"requestId,omitempty" xml:"requestId,omitempty"`
	Tasks     []*DescribeClusterTasksResponseBodyTasks  `json:"tasks,omitempty" xml:"tasks,omitempty" type:"Repeated"`
}

func (s DescribeClusterTasksResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterTasksResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterTasksResponseBody) SetPageInfo(v *DescribeClusterTasksResponseBodyPageInfo) *DescribeClusterTasksResponseBody {
	s.PageInfo = v
	return s
}

func (s *DescribeClusterTasksResponseBody) SetRequestId(v string) *DescribeClusterTasksResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeClusterTasksResponseBody) SetTasks(v []*DescribeClusterTasksResponseBodyTasks) *DescribeClusterTasksResponseBody {
	s.Tasks = v
	return s
}

type DescribeClusterTasksResponseBodyPageInfo struct {
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	PageSize   *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	TotalCount *int64 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeClusterTasksResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterTasksResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeClusterTasksResponseBodyPageInfo) SetPageNumber(v int64) *DescribeClusterTasksResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeClusterTasksResponseBodyPageInfo) SetPageSize(v int64) *DescribeClusterTasksResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeClusterTasksResponseBodyPageInfo) SetTotalCount(v int64) *DescribeClusterTasksResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeClusterTasksResponseBodyTasks struct {
	Created  *string                                     `json:"created,omitempty" xml:"created,omitempty"`
	Error    *DescribeClusterTasksResponseBodyTasksError `json:"error,omitempty" xml:"error,omitempty" type:"Struct"`
	State    *string                                     `json:"state,omitempty" xml:"state,omitempty"`
	TaskId   *string                                     `json:"task_id,omitempty" xml:"task_id,omitempty"`
	TaskType *string                                     `json:"task_type,omitempty" xml:"task_type,omitempty"`
	Updated  *string                                     `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeClusterTasksResponseBodyTasks) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterTasksResponseBodyTasks) GoString() string {
	return s.String()
}

func (s *DescribeClusterTasksResponseBodyTasks) SetCreated(v string) *DescribeClusterTasksResponseBodyTasks {
	s.Created = &v
	return s
}

func (s *DescribeClusterTasksResponseBodyTasks) SetError(v *DescribeClusterTasksResponseBodyTasksError) *DescribeClusterTasksResponseBodyTasks {
	s.Error = v
	return s
}

func (s *DescribeClusterTasksResponseBodyTasks) SetState(v string) *DescribeClusterTasksResponseBodyTasks {
	s.State = &v
	return s
}

func (s *DescribeClusterTasksResponseBodyTasks) SetTaskId(v string) *DescribeClusterTasksResponseBodyTasks {
	s.TaskId = &v
	return s
}

func (s *DescribeClusterTasksResponseBodyTasks) SetTaskType(v string) *DescribeClusterTasksResponseBodyTasks {
	s.TaskType = &v
	return s
}

func (s *DescribeClusterTasksResponseBodyTasks) SetUpdated(v string) *DescribeClusterTasksResponseBodyTasks {
	s.Updated = &v
	return s
}

type DescribeClusterTasksResponseBodyTasksError struct {
	Code    *string `json:"code,omitempty" xml:"code,omitempty"`
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
}

func (s DescribeClusterTasksResponseBodyTasksError) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterTasksResponseBodyTasksError) GoString() string {
	return s.String()
}

func (s *DescribeClusterTasksResponseBodyTasksError) SetCode(v string) *DescribeClusterTasksResponseBodyTasksError {
	s.Code = &v
	return s
}

func (s *DescribeClusterTasksResponseBodyTasksError) SetMessage(v string) *DescribeClusterTasksResponseBodyTasksError {
	s.Message = &v
	return s
}

type DescribeClusterTasksResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterTasksResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterTasksResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterTasksResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterTasksResponse) SetHeaders(v map[string]*string) *DescribeClusterTasksResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterTasksResponse) SetStatusCode(v int32) *DescribeClusterTasksResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterTasksResponse) SetBody(v *DescribeClusterTasksResponseBody) *DescribeClusterTasksResponse {
	s.Body = v
	return s
}

type DescribeClusterUserKubeconfigRequest struct {
	PrivateIpAddress         *bool  `json:"PrivateIpAddress,omitempty" xml:"PrivateIpAddress,omitempty"`
	TemporaryDurationMinutes *int64 `json:"TemporaryDurationMinutes,omitempty" xml:"TemporaryDurationMinutes,omitempty"`
}

func (s DescribeClusterUserKubeconfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterUserKubeconfigRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterUserKubeconfigRequest) SetPrivateIpAddress(v bool) *DescribeClusterUserKubeconfigRequest {
	s.PrivateIpAddress = &v
	return s
}

func (s *DescribeClusterUserKubeconfigRequest) SetTemporaryDurationMinutes(v int64) *DescribeClusterUserKubeconfigRequest {
	s.TemporaryDurationMinutes = &v
	return s
}

type DescribeClusterUserKubeconfigResponseBody struct {
	// The content of the kubeconfig file. For more information about the content of the kubeconfig file, see [Configure cluster credentials](~~86494~~).
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
	// The expiration time of the kubeconfig file. The value is the UTC time displayed in RFC3339 format.
	Expiration *string `json:"expiration,omitempty" xml:"expiration,omitempty"`
}

func (s DescribeClusterUserKubeconfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterUserKubeconfigResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterUserKubeconfigResponseBody) SetConfig(v string) *DescribeClusterUserKubeconfigResponseBody {
	s.Config = &v
	return s
}

func (s *DescribeClusterUserKubeconfigResponseBody) SetExpiration(v string) *DescribeClusterUserKubeconfigResponseBody {
	s.Expiration = &v
	return s
}

type DescribeClusterUserKubeconfigResponse struct {
	Headers    map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterUserKubeconfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterUserKubeconfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterUserKubeconfigResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterUserKubeconfigResponse) SetHeaders(v map[string]*string) *DescribeClusterUserKubeconfigResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterUserKubeconfigResponse) SetStatusCode(v int32) *DescribeClusterUserKubeconfigResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterUserKubeconfigResponse) SetBody(v *DescribeClusterUserKubeconfigResponseBody) *DescribeClusterUserKubeconfigResponse {
	s.Body = v
	return s
}

type DescribeClusterV2UserKubeconfigRequest struct {
	PrivateIpAddress *bool `json:"PrivateIpAddress,omitempty" xml:"PrivateIpAddress,omitempty"`
}

func (s DescribeClusterV2UserKubeconfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterV2UserKubeconfigRequest) GoString() string {
	return s.String()
}

func (s *DescribeClusterV2UserKubeconfigRequest) SetPrivateIpAddress(v bool) *DescribeClusterV2UserKubeconfigRequest {
	s.PrivateIpAddress = &v
	return s
}

type DescribeClusterV2UserKubeconfigResponseBody struct {
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
}

func (s DescribeClusterV2UserKubeconfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterV2UserKubeconfigResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterV2UserKubeconfigResponseBody) SetConfig(v string) *DescribeClusterV2UserKubeconfigResponseBody {
	s.Config = &v
	return s
}

type DescribeClusterV2UserKubeconfigResponse struct {
	Headers    map[string]*string                           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterV2UserKubeconfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterV2UserKubeconfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterV2UserKubeconfigResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterV2UserKubeconfigResponse) SetHeaders(v map[string]*string) *DescribeClusterV2UserKubeconfigResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterV2UserKubeconfigResponse) SetStatusCode(v int32) *DescribeClusterV2UserKubeconfigResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterV2UserKubeconfigResponse) SetBody(v *DescribeClusterV2UserKubeconfigResponseBody) *DescribeClusterV2UserKubeconfigResponse {
	s.Body = v
	return s
}

type DescribeClusterVulsResponseBody struct {
	VulRecords []*DescribeClusterVulsResponseBodyVulRecords `json:"vul_records,omitempty" xml:"vul_records,omitempty" type:"Repeated"`
}

func (s DescribeClusterVulsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterVulsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClusterVulsResponseBody) SetVulRecords(v []*DescribeClusterVulsResponseBodyVulRecords) *DescribeClusterVulsResponseBody {
	s.VulRecords = v
	return s
}

type DescribeClusterVulsResponseBodyVulRecords struct {
	CveList      []*string `json:"cve_list,omitempty" xml:"cve_list,omitempty" type:"Repeated"`
	Necessity    *string   `json:"necessity,omitempty" xml:"necessity,omitempty"`
	NodeCount    *int32    `json:"node_count,omitempty" xml:"node_count,omitempty"`
	NodepoolId   *string   `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	NodepoolName *string   `json:"nodepool_name,omitempty" xml:"nodepool_name,omitempty"`
	VulAliasName *string   `json:"vul_alias_name,omitempty" xml:"vul_alias_name,omitempty"`
	VulName      *string   `json:"vul_name,omitempty" xml:"vul_name,omitempty"`
	VulType      *string   `json:"vul_type,omitempty" xml:"vul_type,omitempty"`
}

func (s DescribeClusterVulsResponseBodyVulRecords) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterVulsResponseBodyVulRecords) GoString() string {
	return s.String()
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetCveList(v []*string) *DescribeClusterVulsResponseBodyVulRecords {
	s.CveList = v
	return s
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetNecessity(v string) *DescribeClusterVulsResponseBodyVulRecords {
	s.Necessity = &v
	return s
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetNodeCount(v int32) *DescribeClusterVulsResponseBodyVulRecords {
	s.NodeCount = &v
	return s
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetNodepoolId(v string) *DescribeClusterVulsResponseBodyVulRecords {
	s.NodepoolId = &v
	return s
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetNodepoolName(v string) *DescribeClusterVulsResponseBodyVulRecords {
	s.NodepoolName = &v
	return s
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetVulAliasName(v string) *DescribeClusterVulsResponseBodyVulRecords {
	s.VulAliasName = &v
	return s
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetVulName(v string) *DescribeClusterVulsResponseBodyVulRecords {
	s.VulName = &v
	return s
}

func (s *DescribeClusterVulsResponseBodyVulRecords) SetVulType(v string) *DescribeClusterVulsResponseBodyVulRecords {
	s.VulType = &v
	return s
}

type DescribeClusterVulsResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClusterVulsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClusterVulsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClusterVulsResponse) GoString() string {
	return s.String()
}

func (s *DescribeClusterVulsResponse) SetHeaders(v map[string]*string) *DescribeClusterVulsResponse {
	s.Headers = v
	return s
}

func (s *DescribeClusterVulsResponse) SetStatusCode(v int32) *DescribeClusterVulsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClusterVulsResponse) SetBody(v *DescribeClusterVulsResponseBody) *DescribeClusterVulsResponse {
	s.Body = v
	return s
}

type DescribeClustersRequest struct {
	// The cluster type.
	ClusterType *string `json:"clusterType,omitempty" xml:"clusterType,omitempty"`
	// The cluster name based on which the system performs fuzzy searches among the clusters that belong to the current Alibaba Cloud account.
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
}

func (s DescribeClustersRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersRequest) GoString() string {
	return s.String()
}

func (s *DescribeClustersRequest) SetClusterType(v string) *DescribeClustersRequest {
	s.ClusterType = &v
	return s
}

func (s *DescribeClustersRequest) SetName(v string) *DescribeClustersRequest {
	s.Name = &v
	return s
}

type DescribeClustersResponse struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribeClustersResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeClustersResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersResponse) GoString() string {
	return s.String()
}

func (s *DescribeClustersResponse) SetHeaders(v map[string]*string) *DescribeClustersResponse {
	s.Headers = v
	return s
}

func (s *DescribeClustersResponse) SetStatusCode(v int32) *DescribeClustersResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeClustersResponse) SetBody(v []*DescribeClustersResponseBody) *DescribeClustersResponse {
	s.Body = v
	return s
}

type DescribeClustersResponseBody struct {
	ClusterId              *string                             `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	ClusterType            *string                             `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	Created                *string                             `json:"created,omitempty" xml:"created,omitempty"`
	CurrentVersion         *string                             `json:"current_version,omitempty" xml:"current_version,omitempty"`
	DataDiskCategory       *string                             `json:"data_disk_category,omitempty" xml:"data_disk_category,omitempty"`
	DataDiskSize           *int64                              `json:"data_disk_size,omitempty" xml:"data_disk_size,omitempty"`
	DeletionProtection     *bool                               `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	DockerVersion          *string                             `json:"docker_version,omitempty" xml:"docker_version,omitempty"`
	ExternalLoadbalancerId *string                             `json:"external_loadbalancer_id,omitempty" xml:"external_loadbalancer_id,omitempty"`
	InitVersion            *string                             `json:"init_version,omitempty" xml:"init_version,omitempty"`
	MasterUrl              *string                             `json:"master_url,omitempty" xml:"master_url,omitempty"`
	MetaData               *string                             `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	Name                   *string                             `json:"name,omitempty" xml:"name,omitempty"`
	NetworkMode            *string                             `json:"network_mode,omitempty" xml:"network_mode,omitempty"`
	PrivateZone            *bool                               `json:"private_zone,omitempty" xml:"private_zone,omitempty"`
	Profile                *string                             `json:"profile,omitempty" xml:"profile,omitempty"`
	RegionId               *string                             `json:"region_id,omitempty" xml:"region_id,omitempty"`
	ResourceGroupId        *string                             `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	SecurityGroupId        *string                             `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	Size                   *int64                              `json:"size,omitempty" xml:"size,omitempty"`
	State                  *string                             `json:"state,omitempty" xml:"state,omitempty"`
	SubnetCidr             *string                             `json:"subnet_cidr,omitempty" xml:"subnet_cidr,omitempty"`
	Tags                   []*DescribeClustersResponseBodyTags `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	Updated                *string                             `json:"updated,omitempty" xml:"updated,omitempty"`
	VpcId                  *string                             `json:"vpc_id,omitempty" xml:"vpc_id,omitempty"`
	VswitchCidr            *string                             `json:"vswitch_cidr,omitempty" xml:"vswitch_cidr,omitempty"`
	VswitchId              *string                             `json:"vswitch_id,omitempty" xml:"vswitch_id,omitempty"`
	WorkerRamRoleName      *string                             `json:"worker_ram_role_name,omitempty" xml:"worker_ram_role_name,omitempty"`
	ZoneId                 *string                             `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
}

func (s DescribeClustersResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClustersResponseBody) SetClusterId(v string) *DescribeClustersResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetClusterType(v string) *DescribeClustersResponseBody {
	s.ClusterType = &v
	return s
}

func (s *DescribeClustersResponseBody) SetCreated(v string) *DescribeClustersResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeClustersResponseBody) SetCurrentVersion(v string) *DescribeClustersResponseBody {
	s.CurrentVersion = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDataDiskCategory(v string) *DescribeClustersResponseBody {
	s.DataDiskCategory = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDataDiskSize(v int64) *DescribeClustersResponseBody {
	s.DataDiskSize = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDeletionProtection(v bool) *DescribeClustersResponseBody {
	s.DeletionProtection = &v
	return s
}

func (s *DescribeClustersResponseBody) SetDockerVersion(v string) *DescribeClustersResponseBody {
	s.DockerVersion = &v
	return s
}

func (s *DescribeClustersResponseBody) SetExternalLoadbalancerId(v string) *DescribeClustersResponseBody {
	s.ExternalLoadbalancerId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetInitVersion(v string) *DescribeClustersResponseBody {
	s.InitVersion = &v
	return s
}

func (s *DescribeClustersResponseBody) SetMasterUrl(v string) *DescribeClustersResponseBody {
	s.MasterUrl = &v
	return s
}

func (s *DescribeClustersResponseBody) SetMetaData(v string) *DescribeClustersResponseBody {
	s.MetaData = &v
	return s
}

func (s *DescribeClustersResponseBody) SetName(v string) *DescribeClustersResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeClustersResponseBody) SetNetworkMode(v string) *DescribeClustersResponseBody {
	s.NetworkMode = &v
	return s
}

func (s *DescribeClustersResponseBody) SetPrivateZone(v bool) *DescribeClustersResponseBody {
	s.PrivateZone = &v
	return s
}

func (s *DescribeClustersResponseBody) SetProfile(v string) *DescribeClustersResponseBody {
	s.Profile = &v
	return s
}

func (s *DescribeClustersResponseBody) SetRegionId(v string) *DescribeClustersResponseBody {
	s.RegionId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetResourceGroupId(v string) *DescribeClustersResponseBody {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetSecurityGroupId(v string) *DescribeClustersResponseBody {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetSize(v int64) *DescribeClustersResponseBody {
	s.Size = &v
	return s
}

func (s *DescribeClustersResponseBody) SetState(v string) *DescribeClustersResponseBody {
	s.State = &v
	return s
}

func (s *DescribeClustersResponseBody) SetSubnetCidr(v string) *DescribeClustersResponseBody {
	s.SubnetCidr = &v
	return s
}

func (s *DescribeClustersResponseBody) SetTags(v []*DescribeClustersResponseBodyTags) *DescribeClustersResponseBody {
	s.Tags = v
	return s
}

func (s *DescribeClustersResponseBody) SetUpdated(v string) *DescribeClustersResponseBody {
	s.Updated = &v
	return s
}

func (s *DescribeClustersResponseBody) SetVpcId(v string) *DescribeClustersResponseBody {
	s.VpcId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetVswitchCidr(v string) *DescribeClustersResponseBody {
	s.VswitchCidr = &v
	return s
}

func (s *DescribeClustersResponseBody) SetVswitchId(v string) *DescribeClustersResponseBody {
	s.VswitchId = &v
	return s
}

func (s *DescribeClustersResponseBody) SetWorkerRamRoleName(v string) *DescribeClustersResponseBody {
	s.WorkerRamRoleName = &v
	return s
}

func (s *DescribeClustersResponseBody) SetZoneId(v string) *DescribeClustersResponseBody {
	s.ZoneId = &v
	return s
}

type DescribeClustersResponseBodyTags struct {
	Key   *string `json:"key,omitempty" xml:"key,omitempty"`
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s DescribeClustersResponseBodyTags) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersResponseBodyTags) GoString() string {
	return s.String()
}

func (s *DescribeClustersResponseBodyTags) SetKey(v string) *DescribeClustersResponseBodyTags {
	s.Key = &v
	return s
}

func (s *DescribeClustersResponseBodyTags) SetValue(v string) *DescribeClustersResponseBodyTags {
	s.Value = &v
	return s
}

type DescribeClustersV1Request struct {
	ClusterSpec *string `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	Name        *string `json:"name,omitempty" xml:"name,omitempty"`
	PageNumber  *int64  `json:"page_number,omitempty" xml:"page_number,omitempty"`
	PageSize    *int64  `json:"page_size,omitempty" xml:"page_size,omitempty"`
	Profile     *string `json:"profile,omitempty" xml:"profile,omitempty"`
	RegionId    *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
}

func (s DescribeClustersV1Request) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1Request) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1Request) SetClusterSpec(v string) *DescribeClustersV1Request {
	s.ClusterSpec = &v
	return s
}

func (s *DescribeClustersV1Request) SetClusterType(v string) *DescribeClustersV1Request {
	s.ClusterType = &v
	return s
}

func (s *DescribeClustersV1Request) SetName(v string) *DescribeClustersV1Request {
	s.Name = &v
	return s
}

func (s *DescribeClustersV1Request) SetPageNumber(v int64) *DescribeClustersV1Request {
	s.PageNumber = &v
	return s
}

func (s *DescribeClustersV1Request) SetPageSize(v int64) *DescribeClustersV1Request {
	s.PageSize = &v
	return s
}

func (s *DescribeClustersV1Request) SetProfile(v string) *DescribeClustersV1Request {
	s.Profile = &v
	return s
}

func (s *DescribeClustersV1Request) SetRegionId(v string) *DescribeClustersV1Request {
	s.RegionId = &v
	return s
}

type DescribeClustersV1ResponseBody struct {
	// The list of the details of the queried cluster.
	Clusters []*DescribeClustersV1ResponseBodyClusters `json:"clusters,omitempty" xml:"clusters,omitempty" type:"Repeated"`
	// The pagination details.
	PageInfo *DescribeClustersV1ResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
}

func (s DescribeClustersV1ResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1ResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1ResponseBody) SetClusters(v []*DescribeClustersV1ResponseBodyClusters) *DescribeClustersV1ResponseBody {
	s.Clusters = v
	return s
}

func (s *DescribeClustersV1ResponseBody) SetPageInfo(v *DescribeClustersV1ResponseBodyPageInfo) *DescribeClustersV1ResponseBody {
	s.PageInfo = v
	return s
}

type DescribeClustersV1ResponseBodyClusters struct {
	// The ID of the queried cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The type of the managed Kubernetes cluster. This parameter is returned for a managed Kubernetes cluster. Valid values:
	//
	// *   `ack.pro.small`: professional managed Kubernetes cluster
	// *   `ack.standard`: standard managed Kubernetes cluster
	ClusterSpec *string `json:"cluster_spec,omitempty" xml:"cluster_spec,omitempty"`
	// The type of the cluster. Valid values:
	//
	// *   `Kubernetes`: dedicated Kubernetes cluster
	// *   `ManagedKubernetes`: managed Kubernetes cluster
	// *   `Ask`: ASK cluster
	// *   `ExternalKubernetes`: registered external cluster
	ClusterType *string `json:"cluster_type,omitempty" xml:"cluster_type,omitempty"`
	// The time when the cluster was created.
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// The Kubernetes version of the cluster.
	CurrentVersion *string `json:"current_version,omitempty" xml:"current_version,omitempty"`
	// Indicates whether deletion protection is enabled for the cluster. After deletion protection is enabled, the cluster cannot be deleted in the console or by calling API operations. Valid values:
	//
	// *   `true`: deletion protection is enabled for the cluster. The cluster cannot be deleted in the ACK console or by calling API operations.
	// *   `false`: deletion protection is disabled for the cluster. The cluster can be deleted in the ACK console or by calling API operations.
	DeletionProtection *bool `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	// The Docker version that is used by the cluster.
	DockerVersion *string `json:"docker_version,omitempty" xml:"docker_version,omitempty"`
	// The ID of the Server Load Balancer (SLB) instance that is used for the Ingress of the cluster.
	//
	// The default SLB specification is slb.s1.small, which belongs to the high-performance instance type.
	ExternalLoadbalancerId *string `json:"external_loadbalancer_id,omitempty" xml:"external_loadbalancer_id,omitempty"`
	// The Kubernetes version of the cluster. The Kubernetes versions provided by ACK are consistent with the open source Kubernetes versions. We recommend that you select the latest Kubernetes version. If you do not specify a Kubernetes version, the latest Kubernetes version is used by default.
	//
	// You can create clusters of the latest two Kubernetes versions in the ACK console. You can create ACK clusters of earlier Kubernetes versions by calling API operations. For more information about the Kubernetes versions supported by ACK, see [Release notes for Kubernetes versions](~~185269~~).
	InitVersion *string `json:"init_version,omitempty" xml:"init_version,omitempty"`
	// The maintenance window of the cluster. This feature is available only in professional managed Kubernetes clusters.
	MaintenanceWindow *MaintenanceWindow `json:"maintenance_window,omitempty" xml:"maintenance_window,omitempty"`
	// The address of the cluster API server. It includes an internal endpoint and a public endpoint.
	MasterUrl *string `json:"master_url,omitempty" xml:"master_url,omitempty"`
	// The metadata of the cluster.
	MetaData *string `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	// The name of the cluster.
	//
	// The name must be 1 to 63 characters in length, and can contain digits, letters, and hyphens (-). It cannot start with a hyphen (-).
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// The network mode of the cluster. Valid values:
	//
	// *   `classic`: classic network
	// *   `vpc`: virtual private cloud (VPC)
	// *   `overlay`: overlay network
	// *   `calico`: network powered by Calico
	NetworkMode *string `json:"network_mode,omitempty" xml:"network_mode,omitempty"`
	// The Kubernetes version to which the cluster can be upgraded.
	NextVersion *string `json:"next_version,omitempty" xml:"next_version,omitempty"`
	// Indicates whether Alibaba Cloud DNS PrivateZone is enabled. Valid values:
	//
	// *   `true`: Alibaba Cloud DNS PrivateZone is enabled.
	// *   `false`: Alibaba Cloud DNS PrivateZone is disabled.
	PrivateZone *bool `json:"private_zone,omitempty" xml:"private_zone,omitempty"`
	// The identifier of the cluster. Valid values:
	//
	// *   `Edge`: The cluster is a managed edge Kubernetes cluster.
	// *   `Default`: The cluster is not a managed edge Kubernetes cluster.
	Profile *string `json:"profile,omitempty" xml:"profile,omitempty"`
	// The ID of the region where the cluster is deployed.
	RegionId *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	// The ID of the resource group to which the cluster belongs.
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
	// The ID of the security group to which the instances of the cluster belong.
	SecurityGroupId *string `json:"security_group_id,omitempty" xml:"security_group_id,omitempty"`
	// The number of nodes in the cluster. Master nodes and worker nodes are included.
	Size *int64 `json:"size,omitempty" xml:"size,omitempty"`
	// The status of the cluster. Valid values:
	//
	// *   `initial`: The cluster is being created.
	// *   `failed`: The cluster failed to be created.
	// *   `running`: The cluster is running.
	// *   `updating`: The cluster is being upgraded.
	// *   `updating_failed`: The cluster failed to be upgraded.
	// *   `scaling`: The cluster is being scaled.
	// *   `stopped`: The cluster is stopped.
	// *   `deleting`: The cluster is being deleted.
	// *   `deleted`: The cluster is deleted.
	// *   `delete_failed`: The cluster failed to be deleted.
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// The pod CIDR block. It must be a valid and private CIDR block, and must be one of the following CIDR blocks or their subnets:
	//
	// *   10.0.0.0/8
	// *   172.16-31.0.0/12-16
	// *   192.168.0.0/16
	//
	// The CIDR block of pods cannot overlap with the CIDR block of the VPC in which the cluster is deployed and the CIDR blocks of existing clusters in the VPC. You cannot modify the pod CIDR block after the cluster is created.
	//
	// For more information about subnetting for ACK clusters, see [Plan CIDR blocks for ACK clusters in a VPC](~~86500~~).
	SubnetCidr *string `json:"subnet_cidr,omitempty" xml:"subnet_cidr,omitempty"`
	// The labels of the cluster.
	Tags []*Tag `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	// The time when the cluster was updated.
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
	// The ID of the VPC where the cluster is deployed. You must specify a VPC when you create a cluster.
	VpcId *string `json:"vpc_id,omitempty" xml:"vpc_id,omitempty"`
	// The IDs of the vSwitches. You can select one to three vSwitches when you create a cluster. We recommend that you select vSwitches in different zones to ensure high availability.
	VswitchId *string `json:"vswitch_id,omitempty" xml:"vswitch_id,omitempty"`
	// The name of the worker Resource Access Management (RAM) role. The RAM role is assigned to the worker nodes that are created on Elastic Compute Service (ECS) instances.
	WorkerRamRoleName *string `json:"worker_ram_role_name,omitempty" xml:"worker_ram_role_name,omitempty"`
	// The ID of the zone where the cluster is deployed.
	ZoneId *string `json:"zone_id,omitempty" xml:"zone_id,omitempty"`
}

func (s DescribeClustersV1ResponseBodyClusters) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1ResponseBodyClusters) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1ResponseBodyClusters) SetClusterId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ClusterId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetClusterSpec(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ClusterSpec = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetClusterType(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ClusterType = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetCreated(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Created = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetCurrentVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.CurrentVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetDeletionProtection(v bool) *DescribeClustersV1ResponseBodyClusters {
	s.DeletionProtection = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetDockerVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.DockerVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetExternalLoadbalancerId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ExternalLoadbalancerId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetInitVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.InitVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetMaintenanceWindow(v *MaintenanceWindow) *DescribeClustersV1ResponseBodyClusters {
	s.MaintenanceWindow = v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetMasterUrl(v string) *DescribeClustersV1ResponseBodyClusters {
	s.MasterUrl = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetMetaData(v string) *DescribeClustersV1ResponseBodyClusters {
	s.MetaData = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetName(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Name = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetNetworkMode(v string) *DescribeClustersV1ResponseBodyClusters {
	s.NetworkMode = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetNextVersion(v string) *DescribeClustersV1ResponseBodyClusters {
	s.NextVersion = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetPrivateZone(v bool) *DescribeClustersV1ResponseBodyClusters {
	s.PrivateZone = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetProfile(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Profile = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetRegionId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.RegionId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetResourceGroupId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ResourceGroupId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetSecurityGroupId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.SecurityGroupId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetSize(v int64) *DescribeClustersV1ResponseBodyClusters {
	s.Size = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetState(v string) *DescribeClustersV1ResponseBodyClusters {
	s.State = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetSubnetCidr(v string) *DescribeClustersV1ResponseBodyClusters {
	s.SubnetCidr = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetTags(v []*Tag) *DescribeClustersV1ResponseBodyClusters {
	s.Tags = v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetUpdated(v string) *DescribeClustersV1ResponseBodyClusters {
	s.Updated = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetVpcId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.VpcId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetVswitchId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.VswitchId = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetWorkerRamRoleName(v string) *DescribeClustersV1ResponseBodyClusters {
	s.WorkerRamRoleName = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyClusters) SetZoneId(v string) *DescribeClustersV1ResponseBodyClusters {
	s.ZoneId = &v
	return s
}

type DescribeClustersV1ResponseBodyPageInfo struct {
	// The number of the returned page.
	PageNumber *int32 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// The number of entries returned per page.
	PageSize *int32 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// The total number of entries returned.
	TotalCount *int32 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeClustersV1ResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1ResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1ResponseBodyPageInfo) SetPageNumber(v int32) *DescribeClustersV1ResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyPageInfo) SetPageSize(v int32) *DescribeClustersV1ResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeClustersV1ResponseBodyPageInfo) SetTotalCount(v int32) *DescribeClustersV1ResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeClustersV1Response struct {
	Headers    map[string]*string              `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                          `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeClustersV1ResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeClustersV1Response) String() string {
	return tea.Prettify(s)
}

func (s DescribeClustersV1Response) GoString() string {
	return s.String()
}

func (s *DescribeClustersV1Response) SetHeaders(v map[string]*string) *DescribeClustersV1Response {
	s.Headers = v
	return s
}

func (s *DescribeClustersV1Response) SetStatusCode(v int32) *DescribeClustersV1Response {
	s.StatusCode = &v
	return s
}

func (s *DescribeClustersV1Response) SetBody(v *DescribeClustersV1ResponseBody) *DescribeClustersV1Response {
	s.Body = v
	return s
}

type DescribeEdgeMachineActiveProcessResponseBody struct {
	// The list of details about the activation progress.
	Logs *string `json:"logs,omitempty" xml:"logs,omitempty"`
	// The activation progress.
	Progress *int64 `json:"progress,omitempty" xml:"progress,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The status of the cloud-native box.
	State *string `json:"state,omitempty" xml:"state,omitempty"`
	// The current step of the activation process.
	Step *string `json:"step,omitempty" xml:"step,omitempty"`
}

func (s DescribeEdgeMachineActiveProcessResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachineActiveProcessResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachineActiveProcessResponseBody) SetLogs(v string) *DescribeEdgeMachineActiveProcessResponseBody {
	s.Logs = &v
	return s
}

func (s *DescribeEdgeMachineActiveProcessResponseBody) SetProgress(v int64) *DescribeEdgeMachineActiveProcessResponseBody {
	s.Progress = &v
	return s
}

func (s *DescribeEdgeMachineActiveProcessResponseBody) SetRequestId(v string) *DescribeEdgeMachineActiveProcessResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeEdgeMachineActiveProcessResponseBody) SetState(v string) *DescribeEdgeMachineActiveProcessResponseBody {
	s.State = &v
	return s
}

func (s *DescribeEdgeMachineActiveProcessResponseBody) SetStep(v string) *DescribeEdgeMachineActiveProcessResponseBody {
	s.Step = &v
	return s
}

type DescribeEdgeMachineActiveProcessResponse struct {
	Headers    map[string]*string                            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeEdgeMachineActiveProcessResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeEdgeMachineActiveProcessResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachineActiveProcessResponse) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachineActiveProcessResponse) SetHeaders(v map[string]*string) *DescribeEdgeMachineActiveProcessResponse {
	s.Headers = v
	return s
}

func (s *DescribeEdgeMachineActiveProcessResponse) SetStatusCode(v int32) *DescribeEdgeMachineActiveProcessResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeEdgeMachineActiveProcessResponse) SetBody(v *DescribeEdgeMachineActiveProcessResponseBody) *DescribeEdgeMachineActiveProcessResponse {
	s.Body = v
	return s
}

type DescribeEdgeMachineModelsResponseBody struct {
	// The list of details about the models of cloud-native boxes.
	Models []*DescribeEdgeMachineModelsResponseBodyModels `json:"models,omitempty" xml:"models,omitempty" type:"Repeated"`
}

func (s DescribeEdgeMachineModelsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachineModelsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachineModelsResponseBody) SetModels(v []*DescribeEdgeMachineModelsResponseBodyModels) *DescribeEdgeMachineModelsResponseBody {
	s.Models = v
	return s
}

type DescribeEdgeMachineModelsResponseBodyModels struct {
	// The number of CPU cores.
	Cpu *int32 `json:"cpu,omitempty" xml:"cpu,omitempty"`
	// The CPU architecture.
	CpuArch *string `json:"cpu_arch,omitempty" xml:"cpu_arch,omitempty"`
	// The time when the cloud-native box was created.
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// The description.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Indicates whether the Docker runtime is managed.
	ManageRuntime *int32 `json:"manage_runtime,omitempty" xml:"manage_runtime,omitempty"`
	// The memory size. Unit: GB.
	Memory *int32 `json:"memory,omitempty" xml:"memory,omitempty"`
	// The model of the cloud-native box.
	Model *string `json:"model,omitempty" xml:"model,omitempty"`
	// The ID of the cloud-native box.
	ModelId *string `json:"model_id,omitempty" xml:"model_id,omitempty"`
}

func (s DescribeEdgeMachineModelsResponseBodyModels) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachineModelsResponseBodyModels) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetCpu(v int32) *DescribeEdgeMachineModelsResponseBodyModels {
	s.Cpu = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetCpuArch(v string) *DescribeEdgeMachineModelsResponseBodyModels {
	s.CpuArch = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetCreated(v string) *DescribeEdgeMachineModelsResponseBodyModels {
	s.Created = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetDescription(v string) *DescribeEdgeMachineModelsResponseBodyModels {
	s.Description = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetManageRuntime(v int32) *DescribeEdgeMachineModelsResponseBodyModels {
	s.ManageRuntime = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetMemory(v int32) *DescribeEdgeMachineModelsResponseBodyModels {
	s.Memory = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetModel(v string) *DescribeEdgeMachineModelsResponseBodyModels {
	s.Model = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponseBodyModels) SetModelId(v string) *DescribeEdgeMachineModelsResponseBodyModels {
	s.ModelId = &v
	return s
}

type DescribeEdgeMachineModelsResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeEdgeMachineModelsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeEdgeMachineModelsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachineModelsResponse) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachineModelsResponse) SetHeaders(v map[string]*string) *DescribeEdgeMachineModelsResponse {
	s.Headers = v
	return s
}

func (s *DescribeEdgeMachineModelsResponse) SetStatusCode(v int32) *DescribeEdgeMachineModelsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeEdgeMachineModelsResponse) SetBody(v *DescribeEdgeMachineModelsResponseBody) *DescribeEdgeMachineModelsResponse {
	s.Body = v
	return s
}

type DescribeEdgeMachineTunnelConfigDetailResponseBody struct {
	// The name of the cloud-native box.
	DeviceName *string `json:"device_name,omitempty" xml:"device_name,omitempty"`
	// The model of the cloud-native box.
	Model *string `json:"model,omitempty" xml:"model,omitempty"`
	// The product key.
	ProductKey *string `json:"product_key,omitempty" xml:"product_key,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The serial number of the cloud-native box.
	Sn *string `json:"sn,omitempty" xml:"sn,omitempty"`
	// The token.
	Token *string `json:"token,omitempty" xml:"token,omitempty"`
	// The backend endpoint of the tunnel.
	TunnelEndpoint *string `json:"tunnel_endpoint,omitempty" xml:"tunnel_endpoint,omitempty"`
}

func (s DescribeEdgeMachineTunnelConfigDetailResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachineTunnelConfigDetailResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponseBody) SetDeviceName(v string) *DescribeEdgeMachineTunnelConfigDetailResponseBody {
	s.DeviceName = &v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponseBody) SetModel(v string) *DescribeEdgeMachineTunnelConfigDetailResponseBody {
	s.Model = &v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponseBody) SetProductKey(v string) *DescribeEdgeMachineTunnelConfigDetailResponseBody {
	s.ProductKey = &v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponseBody) SetRequestId(v string) *DescribeEdgeMachineTunnelConfigDetailResponseBody {
	s.RequestId = &v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponseBody) SetSn(v string) *DescribeEdgeMachineTunnelConfigDetailResponseBody {
	s.Sn = &v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponseBody) SetToken(v string) *DescribeEdgeMachineTunnelConfigDetailResponseBody {
	s.Token = &v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponseBody) SetTunnelEndpoint(v string) *DescribeEdgeMachineTunnelConfigDetailResponseBody {
	s.TunnelEndpoint = &v
	return s
}

type DescribeEdgeMachineTunnelConfigDetailResponse struct {
	Headers    map[string]*string                                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeEdgeMachineTunnelConfigDetailResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeEdgeMachineTunnelConfigDetailResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachineTunnelConfigDetailResponse) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponse) SetHeaders(v map[string]*string) *DescribeEdgeMachineTunnelConfigDetailResponse {
	s.Headers = v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponse) SetStatusCode(v int32) *DescribeEdgeMachineTunnelConfigDetailResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeEdgeMachineTunnelConfigDetailResponse) SetBody(v *DescribeEdgeMachineTunnelConfigDetailResponseBody) *DescribeEdgeMachineTunnelConfigDetailResponse {
	s.Body = v
	return s
}

type DescribeEdgeMachinesRequest struct {
	Hostname    *string `json:"hostname,omitempty" xml:"hostname,omitempty"`
	LifeState   *string `json:"life_state,omitempty" xml:"life_state,omitempty"`
	Model       *string `json:"model,omitempty" xml:"model,omitempty"`
	OnlineState *string `json:"online_state,omitempty" xml:"online_state,omitempty"`
	PageNumber  *int64  `json:"page_number,omitempty" xml:"page_number,omitempty"`
	PageSize    *int64  `json:"page_size,omitempty" xml:"page_size,omitempty"`
}

func (s DescribeEdgeMachinesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachinesRequest) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachinesRequest) SetHostname(v string) *DescribeEdgeMachinesRequest {
	s.Hostname = &v
	return s
}

func (s *DescribeEdgeMachinesRequest) SetLifeState(v string) *DescribeEdgeMachinesRequest {
	s.LifeState = &v
	return s
}

func (s *DescribeEdgeMachinesRequest) SetModel(v string) *DescribeEdgeMachinesRequest {
	s.Model = &v
	return s
}

func (s *DescribeEdgeMachinesRequest) SetOnlineState(v string) *DescribeEdgeMachinesRequest {
	s.OnlineState = &v
	return s
}

func (s *DescribeEdgeMachinesRequest) SetPageNumber(v int64) *DescribeEdgeMachinesRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeEdgeMachinesRequest) SetPageSize(v int64) *DescribeEdgeMachinesRequest {
	s.PageSize = &v
	return s
}

type DescribeEdgeMachinesResponseBody struct {
	// The list of details about cloud-native boxes.
	EdgeMachines []*DescribeEdgeMachinesResponseBodyEdgeMachines `json:"edge_machines,omitempty" xml:"edge_machines,omitempty" type:"Repeated"`
	// The pagination details.
	PageInfo *DescribeEdgeMachinesResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
}

func (s DescribeEdgeMachinesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachinesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachinesResponseBody) SetEdgeMachines(v []*DescribeEdgeMachinesResponseBodyEdgeMachines) *DescribeEdgeMachinesResponseBody {
	s.EdgeMachines = v
	return s
}

func (s *DescribeEdgeMachinesResponseBody) SetPageInfo(v *DescribeEdgeMachinesResponseBodyPageInfo) *DescribeEdgeMachinesResponseBody {
	s.PageInfo = v
	return s
}

type DescribeEdgeMachinesResponseBodyEdgeMachines struct {
	// The time when the cloud-native box was activated.
	ActiveTime *string `json:"active_time,omitempty" xml:"active_time,omitempty"`
	// The time when the cloud-native box was created.
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// The ID of the cloud-native box.
	EdgeMachineId *string `json:"edge_machine_id,omitempty" xml:"edge_machine_id,omitempty"`
	// The `hostname` of the cloud-native box.
	Hostname *string `json:"hostname,omitempty" xml:"hostname,omitempty"`
	// The lifecycle status of the cloud-native box.
	LifeState *string `json:"life_state,omitempty" xml:"life_state,omitempty"`
	// The model of the cloud-native box.
	Model *string `json:"model,omitempty" xml:"model,omitempty"`
	// The name of the cloud-native box.
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// The online status of the cloud-native box.
	OnlineState *string `json:"online_state,omitempty" xml:"online_state,omitempty"`
	// The serial number of the cloud-native box.
	Sn *string `json:"sn,omitempty" xml:"sn,omitempty"`
	// The time when the cloud-native box was last updated.
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeEdgeMachinesResponseBodyEdgeMachines) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachinesResponseBodyEdgeMachines) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetActiveTime(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.ActiveTime = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetCreated(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.Created = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetEdgeMachineId(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.EdgeMachineId = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetHostname(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.Hostname = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetLifeState(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.LifeState = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetModel(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.Model = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetName(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.Name = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetOnlineState(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.OnlineState = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetSn(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.Sn = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyEdgeMachines) SetUpdated(v string) *DescribeEdgeMachinesResponseBodyEdgeMachines {
	s.Updated = &v
	return s
}

type DescribeEdgeMachinesResponseBodyPageInfo struct {
	// The page number of the returned page.
	PageNumber *int32 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// The number of entries returned per page.
	PageSize *int32 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// The total number of entries returned.
	TotalCount *int32 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeEdgeMachinesResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachinesResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachinesResponseBodyPageInfo) SetPageNumber(v int32) *DescribeEdgeMachinesResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyPageInfo) SetPageSize(v int32) *DescribeEdgeMachinesResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeEdgeMachinesResponseBodyPageInfo) SetTotalCount(v int32) *DescribeEdgeMachinesResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeEdgeMachinesResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeEdgeMachinesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeEdgeMachinesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeEdgeMachinesResponse) GoString() string {
	return s.String()
}

func (s *DescribeEdgeMachinesResponse) SetHeaders(v map[string]*string) *DescribeEdgeMachinesResponse {
	s.Headers = v
	return s
}

func (s *DescribeEdgeMachinesResponse) SetStatusCode(v int32) *DescribeEdgeMachinesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeEdgeMachinesResponse) SetBody(v *DescribeEdgeMachinesResponseBody) *DescribeEdgeMachinesResponse {
	s.Body = v
	return s
}

type DescribeEventsRequest struct {
	ClusterId  *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	PageNumber *int64  `json:"page_number,omitempty" xml:"page_number,omitempty"`
	PageSize   *int64  `json:"page_size,omitempty" xml:"page_size,omitempty"`
	Type       *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeEventsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsRequest) GoString() string {
	return s.String()
}

func (s *DescribeEventsRequest) SetClusterId(v string) *DescribeEventsRequest {
	s.ClusterId = &v
	return s
}

func (s *DescribeEventsRequest) SetPageNumber(v int64) *DescribeEventsRequest {
	s.PageNumber = &v
	return s
}

func (s *DescribeEventsRequest) SetPageSize(v int64) *DescribeEventsRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeEventsRequest) SetType(v string) *DescribeEventsRequest {
	s.Type = &v
	return s
}

type DescribeEventsResponseBody struct {
	// The details of the event.
	Events   []*DescribeEventsResponseBodyEvents `json:"events,omitempty" xml:"events,omitempty" type:"Repeated"`
	PageInfo *DescribeEventsResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
}

func (s DescribeEventsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBody) SetEvents(v []*DescribeEventsResponseBodyEvents) *DescribeEventsResponseBody {
	s.Events = v
	return s
}

func (s *DescribeEventsResponseBody) SetPageInfo(v *DescribeEventsResponseBodyPageInfo) *DescribeEventsResponseBody {
	s.PageInfo = v
	return s
}

type DescribeEventsResponseBodyEvents struct {
	// The ID of the cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The description of the event.
	Data *DescribeEventsResponseBodyEventsData `json:"data,omitempty" xml:"data,omitempty" type:"Struct"`
	// The ID of the event.
	EventId *string `json:"event_id,omitempty" xml:"event_id,omitempty"`
	// The source of the event.
	Source *string `json:"source,omitempty" xml:"source,omitempty"`
	// The subject of the event.
	Subject *string `json:"subject,omitempty" xml:"subject,omitempty"`
	// The time when the event started.
	Time *string `json:"time,omitempty" xml:"time,omitempty"`
	// The type of the event. Valid values:
	//
	// *   `cluster_create`: cluster creation.
	// *   `cluster_scaleout`: cluster scale-out.
	// *   `cluster_attach`: adding existing nodes.
	// *   `cluster_delete`: cluster deletion.
	// *   `cluster_upgrade`: cluster upgrades.
	// *   `cluster_migrate`: cluster migration.
	// *   `cluster_node_delete`: node removal.
	// *   `cluster_node_drain`: node draining.
	// *   `cluster_modify`: cluster modifications.
	// *   `cluster_configuration_modify`: modifications to cluster control configurations.
	// *   `cluster_addon_install`: component installation.
	// *   `cluster_addon_upgrade`: component upgrades.
	// *   `cluster_addon_uninstall`: component uninstallation.
	// *   `runtime_upgrade`: runtime upgrades.
	// *   `nodepool_upgrade`: node pool upgrades.
	// *   `nodepool_update`: node pool updates.
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeEventsResponseBodyEvents) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBodyEvents) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBodyEvents) SetClusterId(v string) *DescribeEventsResponseBodyEvents {
	s.ClusterId = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetData(v *DescribeEventsResponseBodyEventsData) *DescribeEventsResponseBodyEvents {
	s.Data = v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetEventId(v string) *DescribeEventsResponseBodyEvents {
	s.EventId = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetSource(v string) *DescribeEventsResponseBodyEvents {
	s.Source = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetSubject(v string) *DescribeEventsResponseBodyEvents {
	s.Subject = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetTime(v string) *DescribeEventsResponseBodyEvents {
	s.Time = &v
	return s
}

func (s *DescribeEventsResponseBodyEvents) SetType(v string) *DescribeEventsResponseBodyEvents {
	s.Type = &v
	return s
}

type DescribeEventsResponseBodyEventsData struct {
	// The level of the event.
	Level *string `json:"level,omitempty" xml:"level,omitempty"`
	// The details of the event.
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
	// The state of the event.
	Reason *string `json:"reason,omitempty" xml:"reason,omitempty"`
}

func (s DescribeEventsResponseBodyEventsData) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBodyEventsData) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBodyEventsData) SetLevel(v string) *DescribeEventsResponseBodyEventsData {
	s.Level = &v
	return s
}

func (s *DescribeEventsResponseBodyEventsData) SetMessage(v string) *DescribeEventsResponseBodyEventsData {
	s.Message = &v
	return s
}

func (s *DescribeEventsResponseBodyEventsData) SetReason(v string) *DescribeEventsResponseBodyEventsData {
	s.Reason = &v
	return s
}

type DescribeEventsResponseBodyPageInfo struct {
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	PageSize   *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	TotalCount *int64 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeEventsResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponseBodyPageInfo) SetPageNumber(v int64) *DescribeEventsResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeEventsResponseBodyPageInfo) SetPageSize(v int64) *DescribeEventsResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeEventsResponseBodyPageInfo) SetTotalCount(v int64) *DescribeEventsResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeEventsResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeEventsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeEventsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeEventsResponse) GoString() string {
	return s.String()
}

func (s *DescribeEventsResponse) SetHeaders(v map[string]*string) *DescribeEventsResponse {
	s.Headers = v
	return s
}

func (s *DescribeEventsResponse) SetStatusCode(v int32) *DescribeEventsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeEventsResponse) SetBody(v *DescribeEventsResponseBody) *DescribeEventsResponse {
	s.Body = v
	return s
}

type DescribeExternalAgentRequest struct {
	AgentMode        *string `json:"AgentMode,omitempty" xml:"AgentMode,omitempty"`
	PrivateIpAddress *string `json:"PrivateIpAddress,omitempty" xml:"PrivateIpAddress,omitempty"`
}

func (s DescribeExternalAgentRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeExternalAgentRequest) GoString() string {
	return s.String()
}

func (s *DescribeExternalAgentRequest) SetAgentMode(v string) *DescribeExternalAgentRequest {
	s.AgentMode = &v
	return s
}

func (s *DescribeExternalAgentRequest) SetPrivateIpAddress(v string) *DescribeExternalAgentRequest {
	s.PrivateIpAddress = &v
	return s
}

type DescribeExternalAgentResponseBody struct {
	// The agent configurations in YAML format.
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
}

func (s DescribeExternalAgentResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeExternalAgentResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeExternalAgentResponseBody) SetConfig(v string) *DescribeExternalAgentResponseBody {
	s.Config = &v
	return s
}

type DescribeExternalAgentResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeExternalAgentResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeExternalAgentResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeExternalAgentResponse) GoString() string {
	return s.String()
}

func (s *DescribeExternalAgentResponse) SetHeaders(v map[string]*string) *DescribeExternalAgentResponse {
	s.Headers = v
	return s
}

func (s *DescribeExternalAgentResponse) SetStatusCode(v int32) *DescribeExternalAgentResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeExternalAgentResponse) SetBody(v *DescribeExternalAgentResponseBody) *DescribeExternalAgentResponse {
	s.Body = v
	return s
}

type DescribeKubernetesVersionMetadataRequest struct {
	ClusterType       *string `json:"ClusterType,omitempty" xml:"ClusterType,omitempty"`
	KubernetesVersion *string `json:"KubernetesVersion,omitempty" xml:"KubernetesVersion,omitempty"`
	Mode              *string `json:"Mode,omitempty" xml:"Mode,omitempty"`
	Profile           *string `json:"Profile,omitempty" xml:"Profile,omitempty"`
	Region            *string `json:"Region,omitempty" xml:"Region,omitempty"`
	Runtime           *string `json:"runtime,omitempty" xml:"runtime,omitempty"`
}

func (s DescribeKubernetesVersionMetadataRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataRequest) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataRequest) SetClusterType(v string) *DescribeKubernetesVersionMetadataRequest {
	s.ClusterType = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetKubernetesVersion(v string) *DescribeKubernetesVersionMetadataRequest {
	s.KubernetesVersion = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetMode(v string) *DescribeKubernetesVersionMetadataRequest {
	s.Mode = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetProfile(v string) *DescribeKubernetesVersionMetadataRequest {
	s.Profile = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetRegion(v string) *DescribeKubernetesVersionMetadataRequest {
	s.Region = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataRequest) SetRuntime(v string) *DescribeKubernetesVersionMetadataRequest {
	s.Runtime = &v
	return s
}

type DescribeKubernetesVersionMetadataResponse struct {
	Headers    map[string]*string                               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribeKubernetesVersionMetadataResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeKubernetesVersionMetadataResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataResponse) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataResponse) SetHeaders(v map[string]*string) *DescribeKubernetesVersionMetadataResponse {
	s.Headers = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponse) SetStatusCode(v int32) *DescribeKubernetesVersionMetadataResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponse) SetBody(v []*DescribeKubernetesVersionMetadataResponseBody) *DescribeKubernetesVersionMetadataResponse {
	s.Body = v
	return s
}

type DescribeKubernetesVersionMetadataResponseBody struct {
	Capabilities   map[string]interface{}                                 `json:"capabilities,omitempty" xml:"capabilities,omitempty"`
	Images         []*DescribeKubernetesVersionMetadataResponseBodyImages `json:"images,omitempty" xml:"images,omitempty" type:"Repeated"`
	MetaData       map[string]interface{}                                 `json:"meta_data,omitempty" xml:"meta_data,omitempty"`
	Runtimes       []*Runtime                                             `json:"runtimes,omitempty" xml:"runtimes,omitempty" type:"Repeated"`
	Version        *string                                                `json:"version,omitempty" xml:"version,omitempty"`
	ReleaseDate    *string                                                `json:"release_date,omitempty" xml:"release_date,omitempty"`
	ExpirationDate *string                                                `json:"expiration_date,omitempty" xml:"expiration_date,omitempty"`
	Creatable      *bool                                                  `json:"creatable,omitempty" xml:"creatable,omitempty"`
}

func (s DescribeKubernetesVersionMetadataResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetCapabilities(v map[string]interface{}) *DescribeKubernetesVersionMetadataResponseBody {
	s.Capabilities = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetImages(v []*DescribeKubernetesVersionMetadataResponseBodyImages) *DescribeKubernetesVersionMetadataResponseBody {
	s.Images = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetMetaData(v map[string]interface{}) *DescribeKubernetesVersionMetadataResponseBody {
	s.MetaData = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetRuntimes(v []*Runtime) *DescribeKubernetesVersionMetadataResponseBody {
	s.Runtimes = v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetVersion(v string) *DescribeKubernetesVersionMetadataResponseBody {
	s.Version = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetReleaseDate(v string) *DescribeKubernetesVersionMetadataResponseBody {
	s.ReleaseDate = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetExpirationDate(v string) *DescribeKubernetesVersionMetadataResponseBody {
	s.ExpirationDate = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBody) SetCreatable(v bool) *DescribeKubernetesVersionMetadataResponseBody {
	s.Creatable = &v
	return s
}

type DescribeKubernetesVersionMetadataResponseBodyImages struct {
	ImageId       *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	ImageName     *string `json:"image_name,omitempty" xml:"image_name,omitempty"`
	Platform      *string `json:"platform,omitempty" xml:"platform,omitempty"`
	OsVersion     *string `json:"os_version,omitempty" xml:"os_version,omitempty"`
	ImageType     *string `json:"image_type,omitempty" xml:"image_type,omitempty"`
	OsType        *string `json:"os_type,omitempty" xml:"os_type,omitempty"`
	ImageCategory *string `json:"image_category,omitempty" xml:"image_category,omitempty"`
	Architecture  *string `json:"architecture,omitempty" xml:"architecture,omitempty"`
}

func (s DescribeKubernetesVersionMetadataResponseBodyImages) String() string {
	return tea.Prettify(s)
}

func (s DescribeKubernetesVersionMetadataResponseBodyImages) GoString() string {
	return s.String()
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageId(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageId = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageName(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageName = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetPlatform(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.Platform = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetOsVersion(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.OsVersion = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageType(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageType = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetOsType(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.OsType = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetImageCategory(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.ImageCategory = &v
	return s
}

func (s *DescribeKubernetesVersionMetadataResponseBodyImages) SetArchitecture(v string) *DescribeKubernetesVersionMetadataResponseBodyImages {
	s.Architecture = &v
	return s
}

type DescribeNodePoolVulsRequest struct {
	Necessity *string `json:"necessity,omitempty" xml:"necessity,omitempty"`
}

func (s DescribeNodePoolVulsRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeNodePoolVulsRequest) GoString() string {
	return s.String()
}

func (s *DescribeNodePoolVulsRequest) SetNecessity(v string) *DescribeNodePoolVulsRequest {
	s.Necessity = &v
	return s
}

type DescribeNodePoolVulsResponseBody struct {
	VulRecords              []*DescribeNodePoolVulsResponseBodyVulRecords `json:"vul_records,omitempty" xml:"vul_records,omitempty" type:"Repeated"`
	VulsFixServicePurchased *bool                                         `json:"vuls_fix_service_purchased,omitempty" xml:"vuls_fix_service_purchased,omitempty"`
}

func (s DescribeNodePoolVulsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeNodePoolVulsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeNodePoolVulsResponseBody) SetVulRecords(v []*DescribeNodePoolVulsResponseBodyVulRecords) *DescribeNodePoolVulsResponseBody {
	s.VulRecords = v
	return s
}

func (s *DescribeNodePoolVulsResponseBody) SetVulsFixServicePurchased(v bool) *DescribeNodePoolVulsResponseBody {
	s.VulsFixServicePurchased = &v
	return s
}

type DescribeNodePoolVulsResponseBodyVulRecords struct {
	InstanceId *string                                              `json:"instance_id,omitempty" xml:"instance_id,omitempty"`
	NodeName   *string                                              `json:"node_name,omitempty" xml:"node_name,omitempty"`
	VulList    []*DescribeNodePoolVulsResponseBodyVulRecordsVulList `json:"vul_list,omitempty" xml:"vul_list,omitempty" type:"Repeated"`
}

func (s DescribeNodePoolVulsResponseBodyVulRecords) String() string {
	return tea.Prettify(s)
}

func (s DescribeNodePoolVulsResponseBodyVulRecords) GoString() string {
	return s.String()
}

func (s *DescribeNodePoolVulsResponseBodyVulRecords) SetInstanceId(v string) *DescribeNodePoolVulsResponseBodyVulRecords {
	s.InstanceId = &v
	return s
}

func (s *DescribeNodePoolVulsResponseBodyVulRecords) SetNodeName(v string) *DescribeNodePoolVulsResponseBodyVulRecords {
	s.NodeName = &v
	return s
}

func (s *DescribeNodePoolVulsResponseBodyVulRecords) SetVulList(v []*DescribeNodePoolVulsResponseBodyVulRecordsVulList) *DescribeNodePoolVulsResponseBodyVulRecords {
	s.VulList = v
	return s
}

type DescribeNodePoolVulsResponseBodyVulRecordsVulList struct {
	AliasName *string   `json:"alias_name,omitempty" xml:"alias_name,omitempty"`
	CveList   []*string `json:"cve_list,omitempty" xml:"cve_list,omitempty" type:"Repeated"`
	Name      *string   `json:"name,omitempty" xml:"name,omitempty"`
	Necessity *string   `json:"necessity,omitempty" xml:"necessity,omitempty"`
}

func (s DescribeNodePoolVulsResponseBodyVulRecordsVulList) String() string {
	return tea.Prettify(s)
}

func (s DescribeNodePoolVulsResponseBodyVulRecordsVulList) GoString() string {
	return s.String()
}

func (s *DescribeNodePoolVulsResponseBodyVulRecordsVulList) SetAliasName(v string) *DescribeNodePoolVulsResponseBodyVulRecordsVulList {
	s.AliasName = &v
	return s
}

func (s *DescribeNodePoolVulsResponseBodyVulRecordsVulList) SetCveList(v []*string) *DescribeNodePoolVulsResponseBodyVulRecordsVulList {
	s.CveList = v
	return s
}

func (s *DescribeNodePoolVulsResponseBodyVulRecordsVulList) SetName(v string) *DescribeNodePoolVulsResponseBodyVulRecordsVulList {
	s.Name = &v
	return s
}

func (s *DescribeNodePoolVulsResponseBodyVulRecordsVulList) SetNecessity(v string) *DescribeNodePoolVulsResponseBodyVulRecordsVulList {
	s.Necessity = &v
	return s
}

type DescribeNodePoolVulsResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeNodePoolVulsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeNodePoolVulsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeNodePoolVulsResponse) GoString() string {
	return s.String()
}

func (s *DescribeNodePoolVulsResponse) SetHeaders(v map[string]*string) *DescribeNodePoolVulsResponse {
	s.Headers = v
	return s
}

func (s *DescribeNodePoolVulsResponse) SetStatusCode(v int32) *DescribeNodePoolVulsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeNodePoolVulsResponse) SetBody(v *DescribeNodePoolVulsResponseBody) *DescribeNodePoolVulsResponse {
	s.Body = v
	return s
}

type DescribePoliciesResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       map[string]interface{} `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribePoliciesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePoliciesResponse) GoString() string {
	return s.String()
}

func (s *DescribePoliciesResponse) SetHeaders(v map[string]*string) *DescribePoliciesResponse {
	s.Headers = v
	return s
}

func (s *DescribePoliciesResponse) SetStatusCode(v int32) *DescribePoliciesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePoliciesResponse) SetBody(v map[string]interface{}) *DescribePoliciesResponse {
	s.Body = v
	return s
}

type DescribePolicyDetailsResponseBody struct {
	// The action of the policy. Valid values:
	//
	// *   `enforce`: blocks deployments that match the policy.
	// *   `inform`: generates alerts for deployments that match the policy.
	Action *string `json:"action,omitempty" xml:"action,omitempty"`
	// The type of the policy.
	Category *string `json:"category,omitempty" xml:"category,omitempty"`
	// The description of the policy.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// Indicates whether the policy is deleted. Valid values:
	//
	// *   0: The policy is not deleted.
	// *   1: The policy is deleted.
	IsDeleted *int32 `json:"is_deleted,omitempty" xml:"is_deleted,omitempty"`
	// The name of the policy that is returned.
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// Indicates whether parameters are required. Valid values:
	//
	// *   0: Parameters are required.
	// *   1: Parameters are optional.
	NoConfig *int32 `json:"no_config,omitempty" xml:"no_config,omitempty"`
	// The severity level of the policy. Valid values:
	//
	// *   `high`
	// *   `medium`
	// *   `low`
	Severity *string `json:"severity,omitempty" xml:"severity,omitempty"`
	// The content of the policy.
	Template *string `json:"template,omitempty" xml:"template,omitempty"`
}

func (s DescribePolicyDetailsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyDetailsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribePolicyDetailsResponseBody) SetAction(v string) *DescribePolicyDetailsResponseBody {
	s.Action = &v
	return s
}

func (s *DescribePolicyDetailsResponseBody) SetCategory(v string) *DescribePolicyDetailsResponseBody {
	s.Category = &v
	return s
}

func (s *DescribePolicyDetailsResponseBody) SetDescription(v string) *DescribePolicyDetailsResponseBody {
	s.Description = &v
	return s
}

func (s *DescribePolicyDetailsResponseBody) SetIsDeleted(v int32) *DescribePolicyDetailsResponseBody {
	s.IsDeleted = &v
	return s
}

func (s *DescribePolicyDetailsResponseBody) SetName(v string) *DescribePolicyDetailsResponseBody {
	s.Name = &v
	return s
}

func (s *DescribePolicyDetailsResponseBody) SetNoConfig(v int32) *DescribePolicyDetailsResponseBody {
	s.NoConfig = &v
	return s
}

func (s *DescribePolicyDetailsResponseBody) SetSeverity(v string) *DescribePolicyDetailsResponseBody {
	s.Severity = &v
	return s
}

func (s *DescribePolicyDetailsResponseBody) SetTemplate(v string) *DescribePolicyDetailsResponseBody {
	s.Template = &v
	return s
}

type DescribePolicyDetailsResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribePolicyDetailsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribePolicyDetailsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyDetailsResponse) GoString() string {
	return s.String()
}

func (s *DescribePolicyDetailsResponse) SetHeaders(v map[string]*string) *DescribePolicyDetailsResponse {
	s.Headers = v
	return s
}

func (s *DescribePolicyDetailsResponse) SetStatusCode(v int32) *DescribePolicyDetailsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePolicyDetailsResponse) SetBody(v *DescribePolicyDetailsResponseBody) *DescribePolicyDetailsResponse {
	s.Body = v
	return s
}

type DescribePolicyGovernanceInClusterResponseBody struct {
	// The audit logs of policies in the cluster.
	AdmitLog *DescribePolicyGovernanceInClusterResponseBodyAdmitLog `json:"admit_log,omitempty" xml:"admit_log,omitempty" type:"Struct"`
	// Details about the policies of different severity levels that are enabled for the cluster.
	OnState         []*DescribePolicyGovernanceInClusterResponseBodyOnState       `json:"on_state,omitempty" xml:"on_state,omitempty" type:"Repeated"`
	TotalViolations *DescribePolicyGovernanceInClusterResponseBodyTotalViolations `json:"totalViolations,omitempty" xml:"totalViolations,omitempty" type:"Struct"`
	Violations      *DescribePolicyGovernanceInClusterResponseBodyViolations      `json:"violations,omitempty" xml:"violations,omitempty" type:"Struct"`
}

func (s DescribePolicyGovernanceInClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBody) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBody) SetAdmitLog(v *DescribePolicyGovernanceInClusterResponseBodyAdmitLog) *DescribePolicyGovernanceInClusterResponseBody {
	s.AdmitLog = v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBody) SetOnState(v []*DescribePolicyGovernanceInClusterResponseBodyOnState) *DescribePolicyGovernanceInClusterResponseBody {
	s.OnState = v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBody) SetTotalViolations(v *DescribePolicyGovernanceInClusterResponseBodyTotalViolations) *DescribePolicyGovernanceInClusterResponseBody {
	s.TotalViolations = v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBody) SetViolations(v *DescribePolicyGovernanceInClusterResponseBodyViolations) *DescribePolicyGovernanceInClusterResponseBody {
	s.Violations = v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyAdmitLog struct {
	// The number of audit log entries.
	Count *int64 `json:"count,omitempty" xml:"count,omitempty"`
	// The audit log content.
	Log *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog `json:"log,omitempty" xml:"log,omitempty" type:"Struct"`
	// The status of the query. Valid values:
	//
	// *   `Complete`: The query succeeded and the complete query result is returned.
	// *   `Incomplete`: The query succeeded but the query result is incomplete. To obtain the complete query result, you must repeat the request.
	Progress *string `json:"progress,omitempty" xml:"progress,omitempty"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyAdmitLog) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyAdmitLog) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLog) SetCount(v int64) *DescribePolicyGovernanceInClusterResponseBodyAdmitLog {
	s.Count = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLog) SetLog(v *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) *DescribePolicyGovernanceInClusterResponseBodyAdmitLog {
	s.Log = v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLog) SetProgress(v string) *DescribePolicyGovernanceInClusterResponseBodyAdmitLog {
	s.Progress = &v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog struct {
	// The ID of the cluster that you want to query.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The type of the policy.
	ConstraintKind *string `json:"constraint_kind,omitempty" xml:"constraint_kind,omitempty"`
	// The message that appears when an event is generated by a policy.
	Msg *string `json:"msg,omitempty" xml:"msg,omitempty"`
	// The type of the resource.
	ResourceKind *string `json:"resource_kind,omitempty" xml:"resource_kind,omitempty"`
	// The name of the resource.
	ResourceName *string `json:"resource_name,omitempty" xml:"resource_name,omitempty"`
	// The namespace to which the resource belongs.
	ResourceNamespace *string `json:"resource_namespace,omitempty" xml:"resource_namespace,omitempty"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) SetClusterId(v string) *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog {
	s.ClusterId = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) SetConstraintKind(v string) *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog {
	s.ConstraintKind = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) SetMsg(v string) *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog {
	s.Msg = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) SetResourceKind(v string) *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog {
	s.ResourceKind = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) SetResourceName(v string) *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog {
	s.ResourceName = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog) SetResourceNamespace(v string) *DescribePolicyGovernanceInClusterResponseBodyAdmitLogLog {
	s.ResourceNamespace = &v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyOnState struct {
	// The number of policies that are enabled.
	EnabledCount *int32 `json:"enabled_count,omitempty" xml:"enabled_count,omitempty"`
	// The severity level of the policy.
	Severity *string `json:"severity,omitempty" xml:"severity,omitempty"`
	// The total number of policies of the severity level.
	Total *int32 `json:"total,omitempty" xml:"total,omitempty"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyOnState) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyOnState) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyOnState) SetEnabledCount(v int32) *DescribePolicyGovernanceInClusterResponseBodyOnState {
	s.EnabledCount = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyOnState) SetSeverity(v string) *DescribePolicyGovernanceInClusterResponseBodyOnState {
	s.Severity = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyOnState) SetTotal(v int32) *DescribePolicyGovernanceInClusterResponseBodyOnState {
	s.Total = &v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyTotalViolations struct {
	Deny *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny `json:"deny,omitempty" xml:"deny,omitempty" type:"Struct"`
	Warn *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn `json:"warn,omitempty" xml:"warn,omitempty" type:"Struct"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyTotalViolations) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyTotalViolations) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyTotalViolations) SetDeny(v *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny) *DescribePolicyGovernanceInClusterResponseBodyTotalViolations {
	s.Deny = v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyTotalViolations) SetWarn(v *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn) *DescribePolicyGovernanceInClusterResponseBodyTotalViolations {
	s.Warn = v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny struct {
	Severity   *string `json:"severity,omitempty" xml:"severity,omitempty"`
	Violations *int64  `json:"violations,omitempty" xml:"violations,omitempty"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny) SetSeverity(v string) *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny {
	s.Severity = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny) SetViolations(v int64) *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsDeny {
	s.Violations = &v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn struct {
	Severity   *string `json:"severity,omitempty" xml:"severity,omitempty"`
	Violations *int64  `json:"violations,omitempty" xml:"violations,omitempty"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn) SetSeverity(v string) *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn {
	s.Severity = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn) SetViolations(v int64) *DescribePolicyGovernanceInClusterResponseBodyTotalViolationsWarn {
	s.Violations = &v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyViolations struct {
	Deny *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny `json:"deny,omitempty" xml:"deny,omitempty" type:"Struct"`
	Warn *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn `json:"warn,omitempty" xml:"warn,omitempty" type:"Struct"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyViolations) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyViolations) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolations) SetDeny(v *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny) *DescribePolicyGovernanceInClusterResponseBodyViolations {
	s.Deny = v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolations) SetWarn(v *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn) *DescribePolicyGovernanceInClusterResponseBodyViolations {
	s.Warn = v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyViolationsDeny struct {
	PolicyDescription *string `json:"policyDescription,omitempty" xml:"policyDescription,omitempty"`
	PolicyName        *string `json:"policyName,omitempty" xml:"policyName,omitempty"`
	Severity          *string `json:"severity,omitempty" xml:"severity,omitempty"`
	Violations        *int64  `json:"violations,omitempty" xml:"violations,omitempty"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyViolationsDeny) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyViolationsDeny) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny) SetPolicyDescription(v string) *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny {
	s.PolicyDescription = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny) SetPolicyName(v string) *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny {
	s.PolicyName = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny) SetSeverity(v string) *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny {
	s.Severity = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny) SetViolations(v int64) *DescribePolicyGovernanceInClusterResponseBodyViolationsDeny {
	s.Violations = &v
	return s
}

type DescribePolicyGovernanceInClusterResponseBodyViolationsWarn struct {
	PolicyDescription *string `json:"policyDescription,omitempty" xml:"policyDescription,omitempty"`
	PolicyName        *string `json:"policyName,omitempty" xml:"policyName,omitempty"`
	Severity          *string `json:"severity,omitempty" xml:"severity,omitempty"`
	Violations        *int64  `json:"violations,omitempty" xml:"violations,omitempty"`
}

func (s DescribePolicyGovernanceInClusterResponseBodyViolationsWarn) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponseBodyViolationsWarn) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn) SetPolicyDescription(v string) *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn {
	s.PolicyDescription = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn) SetPolicyName(v string) *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn {
	s.PolicyName = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn) SetSeverity(v string) *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn {
	s.Severity = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn) SetViolations(v int64) *DescribePolicyGovernanceInClusterResponseBodyViolationsWarn {
	s.Violations = &v
	return s
}

type DescribePolicyGovernanceInClusterResponse struct {
	Headers    map[string]*string                             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribePolicyGovernanceInClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribePolicyGovernanceInClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyGovernanceInClusterResponse) GoString() string {
	return s.String()
}

func (s *DescribePolicyGovernanceInClusterResponse) SetHeaders(v map[string]*string) *DescribePolicyGovernanceInClusterResponse {
	s.Headers = v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponse) SetStatusCode(v int32) *DescribePolicyGovernanceInClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePolicyGovernanceInClusterResponse) SetBody(v *DescribePolicyGovernanceInClusterResponseBody) *DescribePolicyGovernanceInClusterResponse {
	s.Body = v
	return s
}

type DescribePolicyInstancesRequest struct {
	InstanceName *string `json:"instance_name,omitempty" xml:"instance_name,omitempty"`
	PolicyName   *string `json:"policy_name,omitempty" xml:"policy_name,omitempty"`
}

func (s DescribePolicyInstancesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyInstancesRequest) GoString() string {
	return s.String()
}

func (s *DescribePolicyInstancesRequest) SetInstanceName(v string) *DescribePolicyInstancesRequest {
	s.InstanceName = &v
	return s
}

func (s *DescribePolicyInstancesRequest) SetPolicyName(v string) *DescribePolicyInstancesRequest {
	s.PolicyName = &v
	return s
}

type DescribePolicyInstancesResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribePolicyInstancesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribePolicyInstancesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyInstancesResponse) GoString() string {
	return s.String()
}

func (s *DescribePolicyInstancesResponse) SetHeaders(v map[string]*string) *DescribePolicyInstancesResponse {
	s.Headers = v
	return s
}

func (s *DescribePolicyInstancesResponse) SetStatusCode(v int32) *DescribePolicyInstancesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePolicyInstancesResponse) SetBody(v []*DescribePolicyInstancesResponseBody) *DescribePolicyInstancesResponse {
	s.Body = v
	return s
}

type DescribePolicyInstancesResponseBody struct {
	AliUid            *string `json:"ali_uid,omitempty" xml:"ali_uid,omitempty"`
	ClusterId         *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	InstanceName      *string `json:"instance_name,omitempty" xml:"instance_name,omitempty"`
	PolicyName        *string `json:"policy_name,omitempty" xml:"policy_name,omitempty"`
	PolicyCategory    *string `json:"policy_category,omitempty" xml:"policy_category,omitempty"`
	PolicyDescription *string `json:"policy_description,omitempty" xml:"policy_description,omitempty"`
	PolicyParameters  *string `json:"policy_parameters,omitempty" xml:"policy_parameters,omitempty"`
	PolicySeverity    *string `json:"policy_severity,omitempty" xml:"policy_severity,omitempty"`
	PolicyScope       *string `json:"policy_scope,omitempty" xml:"policy_scope,omitempty"`
	PolicyAction      *string `json:"policy_action,omitempty" xml:"policy_action,omitempty"`
}

func (s DescribePolicyInstancesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyInstancesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribePolicyInstancesResponseBody) SetAliUid(v string) *DescribePolicyInstancesResponseBody {
	s.AliUid = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetClusterId(v string) *DescribePolicyInstancesResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetInstanceName(v string) *DescribePolicyInstancesResponseBody {
	s.InstanceName = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetPolicyName(v string) *DescribePolicyInstancesResponseBody {
	s.PolicyName = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetPolicyCategory(v string) *DescribePolicyInstancesResponseBody {
	s.PolicyCategory = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetPolicyDescription(v string) *DescribePolicyInstancesResponseBody {
	s.PolicyDescription = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetPolicyParameters(v string) *DescribePolicyInstancesResponseBody {
	s.PolicyParameters = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetPolicySeverity(v string) *DescribePolicyInstancesResponseBody {
	s.PolicySeverity = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetPolicyScope(v string) *DescribePolicyInstancesResponseBody {
	s.PolicyScope = &v
	return s
}

func (s *DescribePolicyInstancesResponseBody) SetPolicyAction(v string) *DescribePolicyInstancesResponseBody {
	s.PolicyAction = &v
	return s
}

type DescribePolicyInstancesStatusResponseBody struct {
	// Information about the number of policy instances of each severity level.
	InstancesSeverityCount map[string]interface{} `json:"instances_severity_count,omitempty" xml:"instances_severity_count,omitempty"`
	// Details about policy instances of different types.
	PolicyInstances []*DescribePolicyInstancesStatusResponseBodyPolicyInstances `json:"policy_instances,omitempty" xml:"policy_instances,omitempty" type:"Repeated"`
}

func (s DescribePolicyInstancesStatusResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyInstancesStatusResponseBody) GoString() string {
	return s.String()
}

func (s *DescribePolicyInstancesStatusResponseBody) SetInstancesSeverityCount(v map[string]interface{}) *DescribePolicyInstancesStatusResponseBody {
	s.InstancesSeverityCount = v
	return s
}

func (s *DescribePolicyInstancesStatusResponseBody) SetPolicyInstances(v []*DescribePolicyInstancesStatusResponseBodyPolicyInstances) *DescribePolicyInstancesStatusResponseBody {
	s.PolicyInstances = v
	return s
}

type DescribePolicyInstancesStatusResponseBodyPolicyInstances struct {
	// The type of the policy. For more information about different types of policies and their descriptions, see [Predefined security policies of ACK](https://www.alibabacloud.com/help/doc-detail/359819.html).
	PolicyCategory *string `json:"policy_category,omitempty" xml:"policy_category,omitempty"`
	// The description of the policy.
	PolicyDescription *string `json:"policy_description,omitempty" xml:"policy_description,omitempty"`
	// The number of policy instances that are deployed. If this parameter is empty, it indicates that no policy instance is deployed from the policy.
	PolicyInstancesCount *int64 `json:"policy_instances_count,omitempty" xml:"policy_instances_count,omitempty"`
	// The name of the policy.
	PolicyName *string `json:"policy_name,omitempty" xml:"policy_name,omitempty"`
	// The severity level of the policy.
	PolicySeverity *string `json:"policy_severity,omitempty" xml:"policy_severity,omitempty"`
}

func (s DescribePolicyInstancesStatusResponseBodyPolicyInstances) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyInstancesStatusResponseBodyPolicyInstances) GoString() string {
	return s.String()
}

func (s *DescribePolicyInstancesStatusResponseBodyPolicyInstances) SetPolicyCategory(v string) *DescribePolicyInstancesStatusResponseBodyPolicyInstances {
	s.PolicyCategory = &v
	return s
}

func (s *DescribePolicyInstancesStatusResponseBodyPolicyInstances) SetPolicyDescription(v string) *DescribePolicyInstancesStatusResponseBodyPolicyInstances {
	s.PolicyDescription = &v
	return s
}

func (s *DescribePolicyInstancesStatusResponseBodyPolicyInstances) SetPolicyInstancesCount(v int64) *DescribePolicyInstancesStatusResponseBodyPolicyInstances {
	s.PolicyInstancesCount = &v
	return s
}

func (s *DescribePolicyInstancesStatusResponseBodyPolicyInstances) SetPolicyName(v string) *DescribePolicyInstancesStatusResponseBodyPolicyInstances {
	s.PolicyName = &v
	return s
}

func (s *DescribePolicyInstancesStatusResponseBodyPolicyInstances) SetPolicySeverity(v string) *DescribePolicyInstancesStatusResponseBodyPolicyInstances {
	s.PolicySeverity = &v
	return s
}

type DescribePolicyInstancesStatusResponse struct {
	Headers    map[string]*string                         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribePolicyInstancesStatusResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribePolicyInstancesStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribePolicyInstancesStatusResponse) GoString() string {
	return s.String()
}

func (s *DescribePolicyInstancesStatusResponse) SetHeaders(v map[string]*string) *DescribePolicyInstancesStatusResponse {
	s.Headers = v
	return s
}

func (s *DescribePolicyInstancesStatusResponse) SetStatusCode(v int32) *DescribePolicyInstancesStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribePolicyInstancesStatusResponse) SetBody(v *DescribePolicyInstancesStatusResponseBody) *DescribePolicyInstancesStatusResponse {
	s.Body = v
	return s
}

type DescribeSubaccountK8sClusterUserConfigRequest struct {
	PrivateIpAddress         *bool  `json:"PrivateIpAddress,omitempty" xml:"PrivateIpAddress,omitempty"`
	TemporaryDurationMinutes *int64 `json:"TemporaryDurationMinutes,omitempty" xml:"TemporaryDurationMinutes,omitempty"`
}

func (s DescribeSubaccountK8sClusterUserConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeSubaccountK8sClusterUserConfigRequest) GoString() string {
	return s.String()
}

func (s *DescribeSubaccountK8sClusterUserConfigRequest) SetPrivateIpAddress(v bool) *DescribeSubaccountK8sClusterUserConfigRequest {
	s.PrivateIpAddress = &v
	return s
}

func (s *DescribeSubaccountK8sClusterUserConfigRequest) SetTemporaryDurationMinutes(v int64) *DescribeSubaccountK8sClusterUserConfigRequest {
	s.TemporaryDurationMinutes = &v
	return s
}

type DescribeSubaccountK8sClusterUserConfigResponseBody struct {
	// The content of the KubeConfig file. For more information about the content of the KubeConfig file, see [Configure cluster credentials](~~86494~~).
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
	// The expiration time of the KubeConfig file. The value is the UTC time displayed in RFC3339 format.
	Expiration *string `json:"expiration,omitempty" xml:"expiration,omitempty"`
}

func (s DescribeSubaccountK8sClusterUserConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeSubaccountK8sClusterUserConfigResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeSubaccountK8sClusterUserConfigResponseBody) SetConfig(v string) *DescribeSubaccountK8sClusterUserConfigResponseBody {
	s.Config = &v
	return s
}

func (s *DescribeSubaccountK8sClusterUserConfigResponseBody) SetExpiration(v string) *DescribeSubaccountK8sClusterUserConfigResponseBody {
	s.Expiration = &v
	return s
}

type DescribeSubaccountK8sClusterUserConfigResponse struct {
	Headers    map[string]*string                                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeSubaccountK8sClusterUserConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeSubaccountK8sClusterUserConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeSubaccountK8sClusterUserConfigResponse) GoString() string {
	return s.String()
}

func (s *DescribeSubaccountK8sClusterUserConfigResponse) SetHeaders(v map[string]*string) *DescribeSubaccountK8sClusterUserConfigResponse {
	s.Headers = v
	return s
}

func (s *DescribeSubaccountK8sClusterUserConfigResponse) SetStatusCode(v int32) *DescribeSubaccountK8sClusterUserConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeSubaccountK8sClusterUserConfigResponse) SetBody(v *DescribeSubaccountK8sClusterUserConfigResponseBody) *DescribeSubaccountK8sClusterUserConfigResponse {
	s.Body = v
	return s
}

type DescribeTaskInfoResponseBody struct {
	// The ID of the ACK cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The time when the task was created.
	Created      *string                               `json:"created,omitempty" xml:"created,omitempty"`
	CurrentStage *string                               `json:"current_stage,omitempty" xml:"current_stage,omitempty"`
	Error        *DescribeTaskInfoResponseBodyError    `json:"error,omitempty" xml:"error,omitempty" type:"Struct"`
	Events       []*DescribeTaskInfoResponseBodyEvents `json:"events,omitempty" xml:"events,omitempty" type:"Repeated"`
	Parameters   map[string]interface{}                `json:"parameters,omitempty" xml:"parameters,omitempty"`
	Stages       []*DescribeTaskInfoResponseBodyStages `json:"stages,omitempty" xml:"stages,omitempty" type:"Repeated"`
	// The state of the task. Valid values:
	//
	// *   `running`: The task is running.
	// *   `fail`: The task failed.
	// *   `success`: The task is complete.
	State  *string                             `json:"state,omitempty" xml:"state,omitempty"`
	Target *DescribeTaskInfoResponseBodyTarget `json:"target,omitempty" xml:"target,omitempty" type:"Struct"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
	// The execution result of the task.
	TaskResult []*DescribeTaskInfoResponseBodyTaskResult `json:"task_result,omitempty" xml:"task_result,omitempty" type:"Repeated"`
	// The task type. A value of `cluster_scaleout` indicates a scale-out task.
	TaskType *string `json:"task_type,omitempty" xml:"task_type,omitempty"`
	// The time when the task was updated.
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeTaskInfoResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBody) SetClusterId(v string) *DescribeTaskInfoResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetCreated(v string) *DescribeTaskInfoResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetCurrentStage(v string) *DescribeTaskInfoResponseBody {
	s.CurrentStage = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetError(v *DescribeTaskInfoResponseBodyError) *DescribeTaskInfoResponseBody {
	s.Error = v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetEvents(v []*DescribeTaskInfoResponseBodyEvents) *DescribeTaskInfoResponseBody {
	s.Events = v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetParameters(v map[string]interface{}) *DescribeTaskInfoResponseBody {
	s.Parameters = v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetStages(v []*DescribeTaskInfoResponseBodyStages) *DescribeTaskInfoResponseBody {
	s.Stages = v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetState(v string) *DescribeTaskInfoResponseBody {
	s.State = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetTarget(v *DescribeTaskInfoResponseBodyTarget) *DescribeTaskInfoResponseBody {
	s.Target = v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetTaskId(v string) *DescribeTaskInfoResponseBody {
	s.TaskId = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetTaskResult(v []*DescribeTaskInfoResponseBodyTaskResult) *DescribeTaskInfoResponseBody {
	s.TaskResult = v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetTaskType(v string) *DescribeTaskInfoResponseBody {
	s.TaskType = &v
	return s
}

func (s *DescribeTaskInfoResponseBody) SetUpdated(v string) *DescribeTaskInfoResponseBody {
	s.Updated = &v
	return s
}

type DescribeTaskInfoResponseBodyError struct {
	Code    *string `json:"code,omitempty" xml:"code,omitempty"`
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
}

func (s DescribeTaskInfoResponseBodyError) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBodyError) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBodyError) SetCode(v string) *DescribeTaskInfoResponseBodyError {
	s.Code = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyError) SetMessage(v string) *DescribeTaskInfoResponseBodyError {
	s.Message = &v
	return s
}

type DescribeTaskInfoResponseBodyEvents struct {
	Action    *string `json:"action,omitempty" xml:"action,omitempty"`
	Level     *string `json:"level,omitempty" xml:"level,omitempty"`
	Message   *string `json:"message,omitempty" xml:"message,omitempty"`
	Reason    *string `json:"reason,omitempty" xml:"reason,omitempty"`
	Source    *string `json:"source,omitempty" xml:"source,omitempty"`
	Timestamp *string `json:"timestamp,omitempty" xml:"timestamp,omitempty"`
}

func (s DescribeTaskInfoResponseBodyEvents) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBodyEvents) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBodyEvents) SetAction(v string) *DescribeTaskInfoResponseBodyEvents {
	s.Action = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyEvents) SetLevel(v string) *DescribeTaskInfoResponseBodyEvents {
	s.Level = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyEvents) SetMessage(v string) *DescribeTaskInfoResponseBodyEvents {
	s.Message = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyEvents) SetReason(v string) *DescribeTaskInfoResponseBodyEvents {
	s.Reason = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyEvents) SetSource(v string) *DescribeTaskInfoResponseBodyEvents {
	s.Source = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyEvents) SetTimestamp(v string) *DescribeTaskInfoResponseBodyEvents {
	s.Timestamp = &v
	return s
}

type DescribeTaskInfoResponseBodyStages struct {
	EndTime   *string                `json:"end_time,omitempty" xml:"end_time,omitempty"`
	Message   *string                `json:"message,omitempty" xml:"message,omitempty"`
	Outputs   map[string]interface{} `json:"outputs,omitempty" xml:"outputs,omitempty"`
	StartTime *string                `json:"start_time,omitempty" xml:"start_time,omitempty"`
	State     *string                `json:"state,omitempty" xml:"state,omitempty"`
}

func (s DescribeTaskInfoResponseBodyStages) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBodyStages) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBodyStages) SetEndTime(v string) *DescribeTaskInfoResponseBodyStages {
	s.EndTime = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyStages) SetMessage(v string) *DescribeTaskInfoResponseBodyStages {
	s.Message = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyStages) SetOutputs(v map[string]interface{}) *DescribeTaskInfoResponseBodyStages {
	s.Outputs = v
	return s
}

func (s *DescribeTaskInfoResponseBodyStages) SetStartTime(v string) *DescribeTaskInfoResponseBodyStages {
	s.StartTime = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyStages) SetState(v string) *DescribeTaskInfoResponseBodyStages {
	s.State = &v
	return s
}

type DescribeTaskInfoResponseBodyTarget struct {
	Id   *string `json:"id,omitempty" xml:"id,omitempty"`
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s DescribeTaskInfoResponseBodyTarget) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBodyTarget) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBodyTarget) SetId(v string) *DescribeTaskInfoResponseBodyTarget {
	s.Id = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyTarget) SetType(v string) *DescribeTaskInfoResponseBodyTarget {
	s.Type = &v
	return s
}

type DescribeTaskInfoResponseBodyTaskResult struct {
	// The resources that are managed by the task. For a scale-out task, the value of this parameter the ID of the instance that is added by the task.
	Data *string `json:"data,omitempty" xml:"data,omitempty"`
	// The state of the scaling of the resource. Valid values:
	//
	// *   `success`: The scale-out task is successful.
	// *   `failed`: The scale-out task failed.
	// *   `initail`: The scale-out task is initializing.
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s DescribeTaskInfoResponseBodyTaskResult) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponseBodyTaskResult) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponseBodyTaskResult) SetData(v string) *DescribeTaskInfoResponseBodyTaskResult {
	s.Data = &v
	return s
}

func (s *DescribeTaskInfoResponseBodyTaskResult) SetStatus(v string) *DescribeTaskInfoResponseBodyTaskResult {
	s.Status = &v
	return s
}

type DescribeTaskInfoResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeTaskInfoResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeTaskInfoResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTaskInfoResponse) GoString() string {
	return s.String()
}

func (s *DescribeTaskInfoResponse) SetHeaders(v map[string]*string) *DescribeTaskInfoResponse {
	s.Headers = v
	return s
}

func (s *DescribeTaskInfoResponse) SetStatusCode(v int32) *DescribeTaskInfoResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeTaskInfoResponse) SetBody(v *DescribeTaskInfoResponseBody) *DescribeTaskInfoResponse {
	s.Body = v
	return s
}

type DescribeTemplateAttributeRequest struct {
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
}

func (s DescribeTemplateAttributeRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplateAttributeRequest) GoString() string {
	return s.String()
}

func (s *DescribeTemplateAttributeRequest) SetTemplateType(v string) *DescribeTemplateAttributeRequest {
	s.TemplateType = &v
	return s
}

type DescribeTemplateAttributeResponse struct {
	Headers    map[string]*string                       `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                   `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribeTemplateAttributeResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeTemplateAttributeResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplateAttributeResponse) GoString() string {
	return s.String()
}

func (s *DescribeTemplateAttributeResponse) SetHeaders(v map[string]*string) *DescribeTemplateAttributeResponse {
	s.Headers = v
	return s
}

func (s *DescribeTemplateAttributeResponse) SetStatusCode(v int32) *DescribeTemplateAttributeResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeTemplateAttributeResponse) SetBody(v []*DescribeTemplateAttributeResponseBody) *DescribeTemplateAttributeResponse {
	s.Body = v
	return s
}

type DescribeTemplateAttributeResponseBody struct {
	Id                 *string `json:"id,omitempty" xml:"id,omitempty"`
	Acl                *string `json:"acl,omitempty" xml:"acl,omitempty"`
	Name               *string `json:"name,omitempty" xml:"name,omitempty"`
	Template           *string `json:"template,omitempty" xml:"template,omitempty"`
	TemplateType       *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
	Description        *string `json:"description,omitempty" xml:"description,omitempty"`
	Tags               *string `json:"tags,omitempty" xml:"tags,omitempty"`
	TemplateWithHistId *string `json:"template_with_hist_id,omitempty" xml:"template_with_hist_id,omitempty"`
	Created            *string `json:"created,omitempty" xml:"created,omitempty"`
	Updated            *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeTemplateAttributeResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplateAttributeResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTemplateAttributeResponseBody) SetId(v string) *DescribeTemplateAttributeResponseBody {
	s.Id = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetAcl(v string) *DescribeTemplateAttributeResponseBody {
	s.Acl = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetName(v string) *DescribeTemplateAttributeResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTemplate(v string) *DescribeTemplateAttributeResponseBody {
	s.Template = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTemplateType(v string) *DescribeTemplateAttributeResponseBody {
	s.TemplateType = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetDescription(v string) *DescribeTemplateAttributeResponseBody {
	s.Description = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTags(v string) *DescribeTemplateAttributeResponseBody {
	s.Tags = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetTemplateWithHistId(v string) *DescribeTemplateAttributeResponseBody {
	s.TemplateWithHistId = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetCreated(v string) *DescribeTemplateAttributeResponseBody {
	s.Created = &v
	return s
}

func (s *DescribeTemplateAttributeResponseBody) SetUpdated(v string) *DescribeTemplateAttributeResponseBody {
	s.Updated = &v
	return s
}

type DescribeTemplatesRequest struct {
	PageNum      *int64  `json:"page_num,omitempty" xml:"page_num,omitempty"`
	PageSize     *int64  `json:"page_size,omitempty" xml:"page_size,omitempty"`
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
}

func (s DescribeTemplatesRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesRequest) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesRequest) SetPageNum(v int64) *DescribeTemplatesRequest {
	s.PageNum = &v
	return s
}

func (s *DescribeTemplatesRequest) SetPageSize(v int64) *DescribeTemplatesRequest {
	s.PageSize = &v
	return s
}

func (s *DescribeTemplatesRequest) SetTemplateType(v string) *DescribeTemplatesRequest {
	s.TemplateType = &v
	return s
}

type DescribeTemplatesResponseBody struct {
	// The pagination details.
	PageInfo *DescribeTemplatesResponseBodyPageInfo `json:"page_info,omitempty" xml:"page_info,omitempty" type:"Struct"`
	// The list of the templates returned .
	Templates []*DescribeTemplatesResponseBodyTemplates `json:"templates,omitempty" xml:"templates,omitempty" type:"Repeated"`
}

func (s DescribeTemplatesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponseBody) SetPageInfo(v *DescribeTemplatesResponseBodyPageInfo) *DescribeTemplatesResponseBody {
	s.PageInfo = v
	return s
}

func (s *DescribeTemplatesResponseBody) SetTemplates(v []*DescribeTemplatesResponseBodyTemplates) *DescribeTemplatesResponseBody {
	s.Templates = v
	return s
}

type DescribeTemplatesResponseBodyPageInfo struct {
	// The page number of the returned page.
	PageNumber *int64 `json:"page_number,omitempty" xml:"page_number,omitempty"`
	// The maximum number of entries returned per page.
	PageSize *int64 `json:"page_size,omitempty" xml:"page_size,omitempty"`
	// The total number of entries returned.
	TotalCount *int64 `json:"total_count,omitempty" xml:"total_count,omitempty"`
}

func (s DescribeTemplatesResponseBodyPageInfo) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponseBodyPageInfo) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponseBodyPageInfo) SetPageNumber(v int64) *DescribeTemplatesResponseBodyPageInfo {
	s.PageNumber = &v
	return s
}

func (s *DescribeTemplatesResponseBodyPageInfo) SetPageSize(v int64) *DescribeTemplatesResponseBodyPageInfo {
	s.PageSize = &v
	return s
}

func (s *DescribeTemplatesResponseBodyPageInfo) SetTotalCount(v int64) *DescribeTemplatesResponseBodyPageInfo {
	s.TotalCount = &v
	return s
}

type DescribeTemplatesResponseBodyTemplates struct {
	// The access control policy of the template. Valid values:
	//
	// *   `private`: The template is private.
	// *   `public`: The template is public.
	// *   `shared`: The template can be shared.
	//
	// Default value: `private`.
	Acl *string `json:"acl,omitempty" xml:"acl,omitempty"`
	// The time when the template was created.
	Created *string `json:"created,omitempty" xml:"created,omitempty"`
	// The description of the template.
	Description *string `json:"description,omitempty" xml:"description,omitempty"`
	// The ID of the template.
	Id *string `json:"id,omitempty" xml:"id,omitempty"`
	// The name of the template.
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
	// The tag of the template. By default, the value is the name of the template.
	Tags *string `json:"tags,omitempty" xml:"tags,omitempty"`
	// The template content in YAML format.
	Template *string `json:"template,omitempty" xml:"template,omitempty"`
	// The type of the template. The value can be a custom value.
	//
	// *   If the value is `kubernetes`, it indicates that the template is displayed on the Templates page in the ACK console.
	// *   If the value is `compose`, it indicates that the template is displayed on the Container Service - Swarm page in the console. However, Container Service for Swarm is deprecated.
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
	// The ID of the parent template. The value of `template_with_hist_id` is the same for each template version. This allows you to manage different template versions.
	TemplateWithHistId *string `json:"template_with_hist_id,omitempty" xml:"template_with_hist_id,omitempty"`
	// The time when the template was updated.
	Updated *string `json:"updated,omitempty" xml:"updated,omitempty"`
}

func (s DescribeTemplatesResponseBodyTemplates) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponseBodyTemplates) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponseBodyTemplates) SetAcl(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Acl = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetCreated(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Created = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetDescription(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Description = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetId(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Id = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetName(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Name = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTags(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Tags = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTemplate(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Template = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTemplateType(v string) *DescribeTemplatesResponseBodyTemplates {
	s.TemplateType = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetTemplateWithHistId(v string) *DescribeTemplatesResponseBodyTemplates {
	s.TemplateWithHistId = &v
	return s
}

func (s *DescribeTemplatesResponseBodyTemplates) SetUpdated(v string) *DescribeTemplatesResponseBodyTemplates {
	s.Updated = &v
	return s
}

type DescribeTemplatesResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeTemplatesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeTemplatesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTemplatesResponse) GoString() string {
	return s.String()
}

func (s *DescribeTemplatesResponse) SetHeaders(v map[string]*string) *DescribeTemplatesResponse {
	s.Headers = v
	return s
}

func (s *DescribeTemplatesResponse) SetStatusCode(v int32) *DescribeTemplatesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeTemplatesResponse) SetBody(v *DescribeTemplatesResponseBody) *DescribeTemplatesResponse {
	s.Body = v
	return s
}

type DescribeTriggerRequest struct {
	Name      *string `json:"Name,omitempty" xml:"Name,omitempty"`
	Namespace *string `json:"Namespace,omitempty" xml:"Namespace,omitempty"`
	Type      *string `json:"Type,omitempty" xml:"Type,omitempty"`
	Action    *string `json:"action,omitempty" xml:"action,omitempty"`
}

func (s DescribeTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s DescribeTriggerRequest) GoString() string {
	return s.String()
}

func (s *DescribeTriggerRequest) SetName(v string) *DescribeTriggerRequest {
	s.Name = &v
	return s
}

func (s *DescribeTriggerRequest) SetNamespace(v string) *DescribeTriggerRequest {
	s.Namespace = &v
	return s
}

func (s *DescribeTriggerRequest) SetType(v string) *DescribeTriggerRequest {
	s.Type = &v
	return s
}

func (s *DescribeTriggerRequest) SetAction(v string) *DescribeTriggerRequest {
	s.Action = &v
	return s
}

type DescribeTriggerResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribeTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeTriggerResponse) GoString() string {
	return s.String()
}

func (s *DescribeTriggerResponse) SetHeaders(v map[string]*string) *DescribeTriggerResponse {
	s.Headers = v
	return s
}

func (s *DescribeTriggerResponse) SetStatusCode(v int32) *DescribeTriggerResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeTriggerResponse) SetBody(v []*DescribeTriggerResponseBody) *DescribeTriggerResponse {
	s.Body = v
	return s
}

type DescribeTriggerResponseBody struct {
	Id        *string `json:"id,omitempty" xml:"id,omitempty"`
	Name      *string `json:"name,omitempty" xml:"name,omitempty"`
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	Type      *string `json:"type,omitempty" xml:"type,omitempty"`
	Action    *string `json:"action,omitempty" xml:"action,omitempty"`
	Token     *string `json:"token,omitempty" xml:"token,omitempty"`
}

func (s DescribeTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeTriggerResponseBody) SetId(v string) *DescribeTriggerResponseBody {
	s.Id = &v
	return s
}

func (s *DescribeTriggerResponseBody) SetName(v string) *DescribeTriggerResponseBody {
	s.Name = &v
	return s
}

func (s *DescribeTriggerResponseBody) SetClusterId(v string) *DescribeTriggerResponseBody {
	s.ClusterId = &v
	return s
}

func (s *DescribeTriggerResponseBody) SetProjectId(v string) *DescribeTriggerResponseBody {
	s.ProjectId = &v
	return s
}

func (s *DescribeTriggerResponseBody) SetType(v string) *DescribeTriggerResponseBody {
	s.Type = &v
	return s
}

func (s *DescribeTriggerResponseBody) SetAction(v string) *DescribeTriggerResponseBody {
	s.Action = &v
	return s
}

func (s *DescribeTriggerResponseBody) SetToken(v string) *DescribeTriggerResponseBody {
	s.Token = &v
	return s
}

type DescribeUserClusterNamespacesResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*string          `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeUserClusterNamespacesResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserClusterNamespacesResponse) GoString() string {
	return s.String()
}

func (s *DescribeUserClusterNamespacesResponse) SetHeaders(v map[string]*string) *DescribeUserClusterNamespacesResponse {
	s.Headers = v
	return s
}

func (s *DescribeUserClusterNamespacesResponse) SetStatusCode(v int32) *DescribeUserClusterNamespacesResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeUserClusterNamespacesResponse) SetBody(v []*string) *DescribeUserClusterNamespacesResponse {
	s.Body = v
	return s
}

type DescribeUserPermissionResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*DescribeUserPermissionResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s DescribeUserPermissionResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserPermissionResponse) GoString() string {
	return s.String()
}

func (s *DescribeUserPermissionResponse) SetHeaders(v map[string]*string) *DescribeUserPermissionResponse {
	s.Headers = v
	return s
}

func (s *DescribeUserPermissionResponse) SetStatusCode(v int32) *DescribeUserPermissionResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeUserPermissionResponse) SetBody(v []*DescribeUserPermissionResponseBody) *DescribeUserPermissionResponse {
	s.Body = v
	return s
}

type DescribeUserPermissionResponseBody struct {
	ResourceId   *string `json:"resource_id,omitempty" xml:"resource_id,omitempty"`
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	RoleName     *string `json:"role_name,omitempty" xml:"role_name,omitempty"`
	RoleType     *string `json:"role_type,omitempty" xml:"role_type,omitempty"`
	IsOwner      *int64  `json:"is_owner,omitempty" xml:"is_owner,omitempty"`
	IsRamRole    *int64  `json:"is_ram_role,omitempty" xml:"is_ram_role,omitempty"`
}

func (s DescribeUserPermissionResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserPermissionResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeUserPermissionResponseBody) SetResourceId(v string) *DescribeUserPermissionResponseBody {
	s.ResourceId = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetResourceType(v string) *DescribeUserPermissionResponseBody {
	s.ResourceType = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetRoleName(v string) *DescribeUserPermissionResponseBody {
	s.RoleName = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetRoleType(v string) *DescribeUserPermissionResponseBody {
	s.RoleType = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetIsOwner(v int64) *DescribeUserPermissionResponseBody {
	s.IsOwner = &v
	return s
}

func (s *DescribeUserPermissionResponseBody) SetIsRamRole(v int64) *DescribeUserPermissionResponseBody {
	s.IsRamRole = &v
	return s
}

type DescribeUserQuotaResponseBody struct {
	// The quota of Container Service for Kubernetes (ACK) managed clusters. Default value: 20. To increase the quota, [go to the Quota Center page to submit a ticket](https://quotas.console.aliyun.com/products/csk/quotas).
	AmkClusterQuota *int64 `json:"amk_cluster_quota,omitempty" xml:"amk_cluster_quota,omitempty"`
	// The quota of serverless Kubernetes (ASK) clusters. Default value: 20. To increase the quota, [go to the Quota Center page to submit a ticket](https://quotas.console.aliyun.com/products/csk/quotas).
	AskClusterQuota *int64 `json:"ask_cluster_quota,omitempty" xml:"ask_cluster_quota,omitempty"`
	// The quota of node pools in an ACK cluster. Default value: 20. To increase the quota, [go to the Quota Center page to submit a ticket](https://quotas.console.aliyun.com/products/csk/quotas).
	ClusterNodepoolQuota *int64 `json:"cluster_nodepool_quota,omitempty" xml:"cluster_nodepool_quota,omitempty"`
	// The quota of clusters within an Alibaba Cloud account. Default value: 50. To increase the quota, [go to the Quota Center page to submit a ticket](https://quotas.console.aliyun.com/products/csk/quotas).
	ClusterQuota *int64 `json:"cluster_quota,omitempty" xml:"cluster_quota,omitempty"`
	// The quota of enhanced edge node pools.
	EdgeImprovedNodepoolQuota *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota `json:"edge_improved_nodepool_quota,omitempty" xml:"edge_improved_nodepool_quota,omitempty" type:"Struct"`
	// The quota of nodes in an ACK cluster. Default value: 100. To increase the quota, [go to the Quota Center page to submit a ticket](https://quotas.console.aliyun.com/products/csk/quotas).
	NodeQuota *int64                  `json:"node_quota,omitempty" xml:"node_quota,omitempty"`
	Quotas    map[string]*QuotasValue `json:"quotas,omitempty" xml:"quotas,omitempty"`
}

func (s DescribeUserQuotaResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserQuotaResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeUserQuotaResponseBody) SetAmkClusterQuota(v int64) *DescribeUserQuotaResponseBody {
	s.AmkClusterQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetAskClusterQuota(v int64) *DescribeUserQuotaResponseBody {
	s.AskClusterQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetClusterNodepoolQuota(v int64) *DescribeUserQuotaResponseBody {
	s.ClusterNodepoolQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetClusterQuota(v int64) *DescribeUserQuotaResponseBody {
	s.ClusterQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetEdgeImprovedNodepoolQuota(v *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota) *DescribeUserQuotaResponseBody {
	s.EdgeImprovedNodepoolQuota = v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetNodeQuota(v int64) *DescribeUserQuotaResponseBody {
	s.NodeQuota = &v
	return s
}

func (s *DescribeUserQuotaResponseBody) SetQuotas(v map[string]*QuotasValue) *DescribeUserQuotaResponseBody {
	s.Quotas = v
	return s
}

type DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota struct {
	// The maximum bandwidth of each enhanced node pool. Unit: Mbit/s.
	Bandwidth *int32 `json:"bandwidth,omitempty" xml:"bandwidth,omitempty"`
	// The quota of enhanced edge node pools within an Alibaba Cloud account.
	Count *int32 `json:"count,omitempty" xml:"count,omitempty"`
	// The maximum subscription duration of an enhanced edge node pool. Unit: months.
	//
	// >  Enhanced node pools use the pay-as-you-go billing method. Therefore, this parameter is not required.
	Period *int32 `json:"period,omitempty" xml:"period,omitempty"`
}

func (s DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota) GoString() string {
	return s.String()
}

func (s *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota) SetBandwidth(v int32) *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota {
	s.Bandwidth = &v
	return s
}

func (s *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota) SetCount(v int32) *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota {
	s.Count = &v
	return s
}

func (s *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota) SetPeriod(v int32) *DescribeUserQuotaResponseBodyEdgeImprovedNodepoolQuota {
	s.Period = &v
	return s
}

type DescribeUserQuotaResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeUserQuotaResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeUserQuotaResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeUserQuotaResponse) GoString() string {
	return s.String()
}

func (s *DescribeUserQuotaResponse) SetHeaders(v map[string]*string) *DescribeUserQuotaResponse {
	s.Headers = v
	return s
}

func (s *DescribeUserQuotaResponse) SetStatusCode(v int32) *DescribeUserQuotaResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeUserQuotaResponse) SetBody(v *DescribeUserQuotaResponseBody) *DescribeUserQuotaResponse {
	s.Body = v
	return s
}

type DescribeWorkflowsResponseBody struct {
	// The list of the jobs.
	Jobs []*DescribeWorkflowsResponseBodyJobs `json:"jobs,omitempty" xml:"jobs,omitempty" type:"Repeated"`
}

func (s DescribeWorkflowsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s DescribeWorkflowsResponseBody) GoString() string {
	return s.String()
}

func (s *DescribeWorkflowsResponseBody) SetJobs(v []*DescribeWorkflowsResponseBodyJobs) *DescribeWorkflowsResponseBody {
	s.Jobs = v
	return s
}

type DescribeWorkflowsResponseBodyJobs struct {
	// The ID of the ACK cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The time when the workflow was created.
	CreateTime *string `json:"create_time,omitempty" xml:"create_time,omitempty"`
	// The name of the workflow.
	JobName *string `json:"job_name,omitempty" xml:"job_name,omitempty"`
}

func (s DescribeWorkflowsResponseBodyJobs) String() string {
	return tea.Prettify(s)
}

func (s DescribeWorkflowsResponseBodyJobs) GoString() string {
	return s.String()
}

func (s *DescribeWorkflowsResponseBodyJobs) SetClusterId(v string) *DescribeWorkflowsResponseBodyJobs {
	s.ClusterId = &v
	return s
}

func (s *DescribeWorkflowsResponseBodyJobs) SetCreateTime(v string) *DescribeWorkflowsResponseBodyJobs {
	s.CreateTime = &v
	return s
}

func (s *DescribeWorkflowsResponseBodyJobs) SetJobName(v string) *DescribeWorkflowsResponseBodyJobs {
	s.JobName = &v
	return s
}

type DescribeWorkflowsResponse struct {
	Headers    map[string]*string             `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                         `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *DescribeWorkflowsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s DescribeWorkflowsResponse) String() string {
	return tea.Prettify(s)
}

func (s DescribeWorkflowsResponse) GoString() string {
	return s.String()
}

func (s *DescribeWorkflowsResponse) SetHeaders(v map[string]*string) *DescribeWorkflowsResponse {
	s.Headers = v
	return s
}

func (s *DescribeWorkflowsResponse) SetStatusCode(v int32) *DescribeWorkflowsResponse {
	s.StatusCode = &v
	return s
}

func (s *DescribeWorkflowsResponse) SetBody(v *DescribeWorkflowsResponseBody) *DescribeWorkflowsResponse {
	s.Body = v
	return s
}

type EdgeClusterAddEdgeMachineRequest struct {
	Expired    *int64  `json:"expired,omitempty" xml:"expired,omitempty"`
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	Options    *string `json:"options,omitempty" xml:"options,omitempty"`
}

func (s EdgeClusterAddEdgeMachineRequest) String() string {
	return tea.Prettify(s)
}

func (s EdgeClusterAddEdgeMachineRequest) GoString() string {
	return s.String()
}

func (s *EdgeClusterAddEdgeMachineRequest) SetExpired(v int64) *EdgeClusterAddEdgeMachineRequest {
	s.Expired = &v
	return s
}

func (s *EdgeClusterAddEdgeMachineRequest) SetNodepoolId(v string) *EdgeClusterAddEdgeMachineRequest {
	s.NodepoolId = &v
	return s
}

func (s *EdgeClusterAddEdgeMachineRequest) SetOptions(v string) *EdgeClusterAddEdgeMachineRequest {
	s.Options = &v
	return s
}

type EdgeClusterAddEdgeMachineResponseBody struct {
	// The ID of the cloud-native box.
	EdgeMachineId *string `json:"edge_machine_id,omitempty" xml:"edge_machine_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
}

func (s EdgeClusterAddEdgeMachineResponseBody) String() string {
	return tea.Prettify(s)
}

func (s EdgeClusterAddEdgeMachineResponseBody) GoString() string {
	return s.String()
}

func (s *EdgeClusterAddEdgeMachineResponseBody) SetEdgeMachineId(v string) *EdgeClusterAddEdgeMachineResponseBody {
	s.EdgeMachineId = &v
	return s
}

func (s *EdgeClusterAddEdgeMachineResponseBody) SetRequestId(v string) *EdgeClusterAddEdgeMachineResponseBody {
	s.RequestId = &v
	return s
}

type EdgeClusterAddEdgeMachineResponse struct {
	Headers    map[string]*string                     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *EdgeClusterAddEdgeMachineResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s EdgeClusterAddEdgeMachineResponse) String() string {
	return tea.Prettify(s)
}

func (s EdgeClusterAddEdgeMachineResponse) GoString() string {
	return s.String()
}

func (s *EdgeClusterAddEdgeMachineResponse) SetHeaders(v map[string]*string) *EdgeClusterAddEdgeMachineResponse {
	s.Headers = v
	return s
}

func (s *EdgeClusterAddEdgeMachineResponse) SetStatusCode(v int32) *EdgeClusterAddEdgeMachineResponse {
	s.StatusCode = &v
	return s
}

func (s *EdgeClusterAddEdgeMachineResponse) SetBody(v *EdgeClusterAddEdgeMachineResponseBody) *EdgeClusterAddEdgeMachineResponse {
	s.Body = v
	return s
}

type FixNodePoolVulsRequest struct {
	Nodes         []*string                            `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
	RolloutPolicy *FixNodePoolVulsRequestRolloutPolicy `json:"rollout_policy,omitempty" xml:"rollout_policy,omitempty" type:"Struct"`
	Vuls          []*string                            `json:"vuls,omitempty" xml:"vuls,omitempty" type:"Repeated"`
}

func (s FixNodePoolVulsRequest) String() string {
	return tea.Prettify(s)
}

func (s FixNodePoolVulsRequest) GoString() string {
	return s.String()
}

func (s *FixNodePoolVulsRequest) SetNodes(v []*string) *FixNodePoolVulsRequest {
	s.Nodes = v
	return s
}

func (s *FixNodePoolVulsRequest) SetRolloutPolicy(v *FixNodePoolVulsRequestRolloutPolicy) *FixNodePoolVulsRequest {
	s.RolloutPolicy = v
	return s
}

func (s *FixNodePoolVulsRequest) SetVuls(v []*string) *FixNodePoolVulsRequest {
	s.Vuls = v
	return s
}

type FixNodePoolVulsRequestRolloutPolicy struct {
	MaxParallelism *int64 `json:"max_parallelism,omitempty" xml:"max_parallelism,omitempty"`
}

func (s FixNodePoolVulsRequestRolloutPolicy) String() string {
	return tea.Prettify(s)
}

func (s FixNodePoolVulsRequestRolloutPolicy) GoString() string {
	return s.String()
}

func (s *FixNodePoolVulsRequestRolloutPolicy) SetMaxParallelism(v int64) *FixNodePoolVulsRequestRolloutPolicy {
	s.MaxParallelism = &v
	return s
}

type FixNodePoolVulsResponseBody struct {
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s FixNodePoolVulsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s FixNodePoolVulsResponseBody) GoString() string {
	return s.String()
}

func (s *FixNodePoolVulsResponseBody) SetTaskId(v string) *FixNodePoolVulsResponseBody {
	s.TaskId = &v
	return s
}

type FixNodePoolVulsResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *FixNodePoolVulsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s FixNodePoolVulsResponse) String() string {
	return tea.Prettify(s)
}

func (s FixNodePoolVulsResponse) GoString() string {
	return s.String()
}

func (s *FixNodePoolVulsResponse) SetHeaders(v map[string]*string) *FixNodePoolVulsResponse {
	s.Headers = v
	return s
}

func (s *FixNodePoolVulsResponse) SetStatusCode(v int32) *FixNodePoolVulsResponse {
	s.StatusCode = &v
	return s
}

func (s *FixNodePoolVulsResponse) SetBody(v *FixNodePoolVulsResponseBody) *FixNodePoolVulsResponse {
	s.Body = v
	return s
}

type GetKubernetesTriggerRequest struct {
	Name      *string `json:"Name,omitempty" xml:"Name,omitempty"`
	Namespace *string `json:"Namespace,omitempty" xml:"Namespace,omitempty"`
	Type      *string `json:"Type,omitempty" xml:"Type,omitempty"`
	Action    *string `json:"action,omitempty" xml:"action,omitempty"`
}

func (s GetKubernetesTriggerRequest) String() string {
	return tea.Prettify(s)
}

func (s GetKubernetesTriggerRequest) GoString() string {
	return s.String()
}

func (s *GetKubernetesTriggerRequest) SetName(v string) *GetKubernetesTriggerRequest {
	s.Name = &v
	return s
}

func (s *GetKubernetesTriggerRequest) SetNamespace(v string) *GetKubernetesTriggerRequest {
	s.Namespace = &v
	return s
}

func (s *GetKubernetesTriggerRequest) SetType(v string) *GetKubernetesTriggerRequest {
	s.Type = &v
	return s
}

func (s *GetKubernetesTriggerRequest) SetAction(v string) *GetKubernetesTriggerRequest {
	s.Action = &v
	return s
}

type GetKubernetesTriggerResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       []*GetKubernetesTriggerResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true" type:"Repeated"`
}

func (s GetKubernetesTriggerResponse) String() string {
	return tea.Prettify(s)
}

func (s GetKubernetesTriggerResponse) GoString() string {
	return s.String()
}

func (s *GetKubernetesTriggerResponse) SetHeaders(v map[string]*string) *GetKubernetesTriggerResponse {
	s.Headers = v
	return s
}

func (s *GetKubernetesTriggerResponse) SetStatusCode(v int32) *GetKubernetesTriggerResponse {
	s.StatusCode = &v
	return s
}

func (s *GetKubernetesTriggerResponse) SetBody(v []*GetKubernetesTriggerResponseBody) *GetKubernetesTriggerResponse {
	s.Body = v
	return s
}

type GetKubernetesTriggerResponseBody struct {
	Id        *string `json:"id,omitempty" xml:"id,omitempty"`
	Name      *string `json:"name,omitempty" xml:"name,omitempty"`
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	ProjectId *string `json:"project_id,omitempty" xml:"project_id,omitempty"`
	Type      *string `json:"type,omitempty" xml:"type,omitempty"`
	Action    *string `json:"action,omitempty" xml:"action,omitempty"`
	// Token
	Token *string `json:"token,omitempty" xml:"token,omitempty"`
}

func (s GetKubernetesTriggerResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetKubernetesTriggerResponseBody) GoString() string {
	return s.String()
}

func (s *GetKubernetesTriggerResponseBody) SetId(v string) *GetKubernetesTriggerResponseBody {
	s.Id = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetName(v string) *GetKubernetesTriggerResponseBody {
	s.Name = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetClusterId(v string) *GetKubernetesTriggerResponseBody {
	s.ClusterId = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetProjectId(v string) *GetKubernetesTriggerResponseBody {
	s.ProjectId = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetType(v string) *GetKubernetesTriggerResponseBody {
	s.Type = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetAction(v string) *GetKubernetesTriggerResponseBody {
	s.Action = &v
	return s
}

func (s *GetKubernetesTriggerResponseBody) SetToken(v string) *GetKubernetesTriggerResponseBody {
	s.Token = &v
	return s
}

type GetUpgradeStatusResponseBody struct {
	// The error message returned during the update.
	ErrorMessage *string `json:"error_message,omitempty" xml:"error_message,omitempty"`
	// The ID of the precheck report.
	PrecheckReportId *string `json:"precheck_report_id,omitempty" xml:"precheck_report_id,omitempty"`
	// The status of the update. Valid values:
	//
	// *   `success`: The update is successful.
	// *   `fail`: The update failed.
	// *   `pause`: The update is paused.
	// *   `running`: The update is in progress.
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
	// The current phase of the update. Valid values:
	//
	// *   `not_start`: The update is not started.
	// *   `prechecking`: The precheck is in progress.
	// *   `upgrading`: The update is in progress.
	// *   `pause`: The update is paused.
	// *   `success`: The update is successful.
	UpgradeStep *string `json:"upgrade_step,omitempty" xml:"upgrade_step,omitempty"`
	// The details of the update task.
	UpgradeTask *GetUpgradeStatusResponseBodyUpgradeTask `json:"upgrade_task,omitempty" xml:"upgrade_task,omitempty" type:"Struct"`
}

func (s GetUpgradeStatusResponseBody) String() string {
	return tea.Prettify(s)
}

func (s GetUpgradeStatusResponseBody) GoString() string {
	return s.String()
}

func (s *GetUpgradeStatusResponseBody) SetErrorMessage(v string) *GetUpgradeStatusResponseBody {
	s.ErrorMessage = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetPrecheckReportId(v string) *GetUpgradeStatusResponseBody {
	s.PrecheckReportId = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetStatus(v string) *GetUpgradeStatusResponseBody {
	s.Status = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetUpgradeStep(v string) *GetUpgradeStatusResponseBody {
	s.UpgradeStep = &v
	return s
}

func (s *GetUpgradeStatusResponseBody) SetUpgradeTask(v *GetUpgradeStatusResponseBodyUpgradeTask) *GetUpgradeStatusResponseBody {
	s.UpgradeTask = v
	return s
}

type GetUpgradeStatusResponseBodyUpgradeTask struct {
	// The description of the update task.
	Message *string `json:"message,omitempty" xml:"message,omitempty"`
	// The status of the update task. Valid values:
	//
	// *   `running`: The update task is being executed.
	// *   `Success`: The update task is successfully executed.
	// *   `Failed`: The update task failed.
	Status *string `json:"status,omitempty" xml:"status,omitempty"`
}

func (s GetUpgradeStatusResponseBodyUpgradeTask) String() string {
	return tea.Prettify(s)
}

func (s GetUpgradeStatusResponseBodyUpgradeTask) GoString() string {
	return s.String()
}

func (s *GetUpgradeStatusResponseBodyUpgradeTask) SetMessage(v string) *GetUpgradeStatusResponseBodyUpgradeTask {
	s.Message = &v
	return s
}

func (s *GetUpgradeStatusResponseBodyUpgradeTask) SetStatus(v string) *GetUpgradeStatusResponseBodyUpgradeTask {
	s.Status = &v
	return s
}

type GetUpgradeStatusResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *GetUpgradeStatusResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s GetUpgradeStatusResponse) String() string {
	return tea.Prettify(s)
}

func (s GetUpgradeStatusResponse) GoString() string {
	return s.String()
}

func (s *GetUpgradeStatusResponse) SetHeaders(v map[string]*string) *GetUpgradeStatusResponse {
	s.Headers = v
	return s
}

func (s *GetUpgradeStatusResponse) SetStatusCode(v int32) *GetUpgradeStatusResponse {
	s.StatusCode = &v
	return s
}

func (s *GetUpgradeStatusResponse) SetBody(v *GetUpgradeStatusResponseBody) *GetUpgradeStatusResponse {
	s.Body = v
	return s
}

type GrantPermissionsRequest struct {
	Body []*GrantPermissionsRequestBody `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s GrantPermissionsRequest) String() string {
	return tea.Prettify(s)
}

func (s GrantPermissionsRequest) GoString() string {
	return s.String()
}

func (s *GrantPermissionsRequest) SetBody(v []*GrantPermissionsRequestBody) *GrantPermissionsRequest {
	s.Body = v
	return s
}

type GrantPermissionsRequestBody struct {
	Cluster   *string `json:"cluster,omitempty" xml:"cluster,omitempty"`
	IsCustom  *bool   `json:"is_custom,omitempty" xml:"is_custom,omitempty"`
	IsRamRole *bool   `json:"is_ram_role,omitempty" xml:"is_ram_role,omitempty"`
	Namespace *string `json:"namespace,omitempty" xml:"namespace,omitempty"`
	RoleName  *string `json:"role_name,omitempty" xml:"role_name,omitempty"`
	RoleType  *string `json:"role_type,omitempty" xml:"role_type,omitempty"`
}

func (s GrantPermissionsRequestBody) String() string {
	return tea.Prettify(s)
}

func (s GrantPermissionsRequestBody) GoString() string {
	return s.String()
}

func (s *GrantPermissionsRequestBody) SetCluster(v string) *GrantPermissionsRequestBody {
	s.Cluster = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetIsCustom(v bool) *GrantPermissionsRequestBody {
	s.IsCustom = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetIsRamRole(v bool) *GrantPermissionsRequestBody {
	s.IsRamRole = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetNamespace(v string) *GrantPermissionsRequestBody {
	s.Namespace = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetRoleName(v string) *GrantPermissionsRequestBody {
	s.RoleName = &v
	return s
}

func (s *GrantPermissionsRequestBody) SetRoleType(v string) *GrantPermissionsRequestBody {
	s.RoleType = &v
	return s
}

type GrantPermissionsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s GrantPermissionsResponse) String() string {
	return tea.Prettify(s)
}

func (s GrantPermissionsResponse) GoString() string {
	return s.String()
}

func (s *GrantPermissionsResponse) SetHeaders(v map[string]*string) *GrantPermissionsResponse {
	s.Headers = v
	return s
}

func (s *GrantPermissionsResponse) SetStatusCode(v int32) *GrantPermissionsResponse {
	s.StatusCode = &v
	return s
}

type InstallClusterAddonsRequest struct {
	Body []*InstallClusterAddonsRequestBody `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s InstallClusterAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s InstallClusterAddonsRequest) GoString() string {
	return s.String()
}

func (s *InstallClusterAddonsRequest) SetBody(v []*InstallClusterAddonsRequestBody) *InstallClusterAddonsRequest {
	s.Body = v
	return s
}

type InstallClusterAddonsRequestBody struct {
	Config  *string `json:"config,omitempty" xml:"config,omitempty"`
	Name    *string `json:"name,omitempty" xml:"name,omitempty"`
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s InstallClusterAddonsRequestBody) String() string {
	return tea.Prettify(s)
}

func (s InstallClusterAddonsRequestBody) GoString() string {
	return s.String()
}

func (s *InstallClusterAddonsRequestBody) SetConfig(v string) *InstallClusterAddonsRequestBody {
	s.Config = &v
	return s
}

func (s *InstallClusterAddonsRequestBody) SetName(v string) *InstallClusterAddonsRequestBody {
	s.Name = &v
	return s
}

func (s *InstallClusterAddonsRequestBody) SetVersion(v string) *InstallClusterAddonsRequestBody {
	s.Version = &v
	return s
}

type InstallClusterAddonsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s InstallClusterAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s InstallClusterAddonsResponse) GoString() string {
	return s.String()
}

func (s *InstallClusterAddonsResponse) SetHeaders(v map[string]*string) *InstallClusterAddonsResponse {
	s.Headers = v
	return s
}

func (s *InstallClusterAddonsResponse) SetStatusCode(v int32) *InstallClusterAddonsResponse {
	s.StatusCode = &v
	return s
}

type ListTagResourcesRequest struct {
	NextToken    *string   `json:"next_token,omitempty" xml:"next_token,omitempty"`
	RegionId     *string   `json:"region_id,omitempty" xml:"region_id,omitempty"`
	ResourceIds  []*string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty" type:"Repeated"`
	ResourceType *string   `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	Tags         []*Tag    `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
}

func (s ListTagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesRequest) GoString() string {
	return s.String()
}

func (s *ListTagResourcesRequest) SetNextToken(v string) *ListTagResourcesRequest {
	s.NextToken = &v
	return s
}

func (s *ListTagResourcesRequest) SetRegionId(v string) *ListTagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *ListTagResourcesRequest) SetResourceIds(v []*string) *ListTagResourcesRequest {
	s.ResourceIds = v
	return s
}

func (s *ListTagResourcesRequest) SetResourceType(v string) *ListTagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *ListTagResourcesRequest) SetTags(v []*Tag) *ListTagResourcesRequest {
	s.Tags = v
	return s
}

type ListTagResourcesShrinkRequest struct {
	NextToken         *string `json:"next_token,omitempty" xml:"next_token,omitempty"`
	RegionId          *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	ResourceIdsShrink *string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty"`
	ResourceType      *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	TagsShrink        *string `json:"tags,omitempty" xml:"tags,omitempty"`
}

func (s ListTagResourcesShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesShrinkRequest) GoString() string {
	return s.String()
}

func (s *ListTagResourcesShrinkRequest) SetNextToken(v string) *ListTagResourcesShrinkRequest {
	s.NextToken = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetRegionId(v string) *ListTagResourcesShrinkRequest {
	s.RegionId = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetResourceIdsShrink(v string) *ListTagResourcesShrinkRequest {
	s.ResourceIdsShrink = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetResourceType(v string) *ListTagResourcesShrinkRequest {
	s.ResourceType = &v
	return s
}

func (s *ListTagResourcesShrinkRequest) SetTagsShrink(v string) *ListTagResourcesShrinkRequest {
	s.TagsShrink = &v
	return s
}

type ListTagResourcesResponseBody struct {
	// The token that is used to start the next query.
	NextToken *string `json:"next_token,omitempty" xml:"next_token,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The details of the queried labels and resources.
	TagResources *ListTagResourcesResponseBodyTagResources `json:"tag_resources,omitempty" xml:"tag_resources,omitempty" type:"Struct"`
}

func (s ListTagResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBody) SetNextToken(v string) *ListTagResourcesResponseBody {
	s.NextToken = &v
	return s
}

func (s *ListTagResourcesResponseBody) SetRequestId(v string) *ListTagResourcesResponseBody {
	s.RequestId = &v
	return s
}

func (s *ListTagResourcesResponseBody) SetTagResources(v *ListTagResourcesResponseBodyTagResources) *ListTagResourcesResponseBody {
	s.TagResources = v
	return s
}

type ListTagResourcesResponseBodyTagResources struct {
	// The labels of the resource.
	TagResource []*ListTagResourcesResponseBodyTagResourcesTagResource `json:"tag_resource,omitempty" xml:"tag_resource,omitempty" type:"Repeated"`
}

func (s ListTagResourcesResponseBodyTagResources) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBodyTagResources) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBodyTagResources) SetTagResource(v []*ListTagResourcesResponseBodyTagResourcesTagResource) *ListTagResourcesResponseBodyTagResources {
	s.TagResource = v
	return s
}

type ListTagResourcesResponseBodyTagResourcesTagResource struct {
	// The ID of the resource.
	ResourceId *string `json:"resource_id,omitempty" xml:"resource_id,omitempty"`
	// The type of the resource. For more information, see [Labels](~~110425~~).
	ResourceType *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	// The key of the label.
	TagKey *string `json:"tag_key,omitempty" xml:"tag_key,omitempty"`
	// The value of the label.
	TagValue *string `json:"tag_value,omitempty" xml:"tag_value,omitempty"`
}

func (s ListTagResourcesResponseBodyTagResourcesTagResource) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponseBodyTagResourcesTagResource) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetResourceId(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.ResourceId = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetResourceType(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.ResourceType = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetTagKey(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.TagKey = &v
	return s
}

func (s *ListTagResourcesResponseBodyTagResourcesTagResource) SetTagValue(v string) *ListTagResourcesResponseBodyTagResourcesTagResource {
	s.TagValue = &v
	return s
}

type ListTagResourcesResponse struct {
	Headers    map[string]*string            `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                        `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ListTagResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ListTagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s ListTagResourcesResponse) GoString() string {
	return s.String()
}

func (s *ListTagResourcesResponse) SetHeaders(v map[string]*string) *ListTagResourcesResponse {
	s.Headers = v
	return s
}

func (s *ListTagResourcesResponse) SetStatusCode(v int32) *ListTagResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *ListTagResourcesResponse) SetBody(v *ListTagResourcesResponseBody) *ListTagResourcesResponse {
	s.Body = v
	return s
}

type MigrateClusterRequest struct {
	OssBucketEndpoint *string `json:"oss_bucket_endpoint,omitempty" xml:"oss_bucket_endpoint,omitempty"`
	OssBucketName     *string `json:"oss_bucket_name,omitempty" xml:"oss_bucket_name,omitempty"`
}

func (s MigrateClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s MigrateClusterRequest) GoString() string {
	return s.String()
}

func (s *MigrateClusterRequest) SetOssBucketEndpoint(v string) *MigrateClusterRequest {
	s.OssBucketEndpoint = &v
	return s
}

func (s *MigrateClusterRequest) SetOssBucketName(v string) *MigrateClusterRequest {
	s.OssBucketName = &v
	return s
}

type MigrateClusterResponseBody struct {
	// The ID of the cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s MigrateClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s MigrateClusterResponseBody) GoString() string {
	return s.String()
}

func (s *MigrateClusterResponseBody) SetClusterId(v string) *MigrateClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *MigrateClusterResponseBody) SetRequestId(v string) *MigrateClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *MigrateClusterResponseBody) SetTaskId(v string) *MigrateClusterResponseBody {
	s.TaskId = &v
	return s
}

type MigrateClusterResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *MigrateClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s MigrateClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s MigrateClusterResponse) GoString() string {
	return s.String()
}

func (s *MigrateClusterResponse) SetHeaders(v map[string]*string) *MigrateClusterResponse {
	s.Headers = v
	return s
}

func (s *MigrateClusterResponse) SetStatusCode(v int32) *MigrateClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *MigrateClusterResponse) SetBody(v *MigrateClusterResponseBody) *MigrateClusterResponse {
	s.Body = v
	return s
}

type ModifyClusterRequest struct {
	ApiServerEip               *bool              `json:"api_server_eip,omitempty" xml:"api_server_eip,omitempty"`
	ApiServerEipId             *string            `json:"api_server_eip_id,omitempty" xml:"api_server_eip_id,omitempty"`
	DeletionProtection         *bool              `json:"deletion_protection,omitempty" xml:"deletion_protection,omitempty"`
	EnableRrsa                 *bool              `json:"enable_rrsa,omitempty" xml:"enable_rrsa,omitempty"`
	IngressDomainRebinding     *string            `json:"ingress_domain_rebinding,omitempty" xml:"ingress_domain_rebinding,omitempty"`
	IngressLoadbalancerId      *string            `json:"ingress_loadbalancer_id,omitempty" xml:"ingress_loadbalancer_id,omitempty"`
	InstanceDeletionProtection *bool              `json:"instance_deletion_protection,omitempty" xml:"instance_deletion_protection,omitempty"`
	MaintenanceWindow          *MaintenanceWindow `json:"maintenance_window,omitempty" xml:"maintenance_window,omitempty"`
	ResourceGroupId            *string            `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
}

func (s ModifyClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterRequest) SetApiServerEip(v bool) *ModifyClusterRequest {
	s.ApiServerEip = &v
	return s
}

func (s *ModifyClusterRequest) SetApiServerEipId(v string) *ModifyClusterRequest {
	s.ApiServerEipId = &v
	return s
}

func (s *ModifyClusterRequest) SetDeletionProtection(v bool) *ModifyClusterRequest {
	s.DeletionProtection = &v
	return s
}

func (s *ModifyClusterRequest) SetEnableRrsa(v bool) *ModifyClusterRequest {
	s.EnableRrsa = &v
	return s
}

func (s *ModifyClusterRequest) SetIngressDomainRebinding(v string) *ModifyClusterRequest {
	s.IngressDomainRebinding = &v
	return s
}

func (s *ModifyClusterRequest) SetIngressLoadbalancerId(v string) *ModifyClusterRequest {
	s.IngressLoadbalancerId = &v
	return s
}

func (s *ModifyClusterRequest) SetInstanceDeletionProtection(v bool) *ModifyClusterRequest {
	s.InstanceDeletionProtection = &v
	return s
}

func (s *ModifyClusterRequest) SetMaintenanceWindow(v *MaintenanceWindow) *ModifyClusterRequest {
	s.MaintenanceWindow = v
	return s
}

func (s *ModifyClusterRequest) SetResourceGroupId(v string) *ModifyClusterRequest {
	s.ResourceGroupId = &v
	return s
}

type ModifyClusterResponseBody struct {
	// The ID of the cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ModifyClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyClusterResponseBody) SetClusterId(v string) *ModifyClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *ModifyClusterResponseBody) SetRequestId(v string) *ModifyClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *ModifyClusterResponseBody) SetTaskId(v string) *ModifyClusterResponseBody {
	s.TaskId = &v
	return s
}

type ModifyClusterResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterResponse) SetHeaders(v map[string]*string) *ModifyClusterResponse {
	s.Headers = v
	return s
}

func (s *ModifyClusterResponse) SetStatusCode(v int32) *ModifyClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyClusterResponse) SetBody(v *ModifyClusterResponseBody) *ModifyClusterResponse {
	s.Body = v
	return s
}

type ModifyClusterAddonRequest struct {
	Config *string `json:"config,omitempty" xml:"config,omitempty"`
}

func (s ModifyClusterAddonRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterAddonRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterAddonRequest) SetConfig(v string) *ModifyClusterAddonRequest {
	s.Config = &v
	return s
}

type ModifyClusterAddonResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s ModifyClusterAddonResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterAddonResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterAddonResponse) SetHeaders(v map[string]*string) *ModifyClusterAddonResponse {
	s.Headers = v
	return s
}

func (s *ModifyClusterAddonResponse) SetStatusCode(v int32) *ModifyClusterAddonResponse {
	s.StatusCode = &v
	return s
}

type ModifyClusterConfigurationRequest struct {
	CustomizeConfig []*ModifyClusterConfigurationRequestCustomizeConfig `json:"customize_config,omitempty" xml:"customize_config,omitempty" type:"Repeated"`
}

func (s ModifyClusterConfigurationRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationRequest) SetCustomizeConfig(v []*ModifyClusterConfigurationRequestCustomizeConfig) *ModifyClusterConfigurationRequest {
	s.CustomizeConfig = v
	return s
}

type ModifyClusterConfigurationRequestCustomizeConfig struct {
	Configs []*ModifyClusterConfigurationRequestCustomizeConfigConfigs `json:"configs,omitempty" xml:"configs,omitempty" type:"Repeated"`
	Name    *string                                                    `json:"name,omitempty" xml:"name,omitempty"`
}

func (s ModifyClusterConfigurationRequestCustomizeConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationRequestCustomizeConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationRequestCustomizeConfig) SetConfigs(v []*ModifyClusterConfigurationRequestCustomizeConfigConfigs) *ModifyClusterConfigurationRequestCustomizeConfig {
	s.Configs = v
	return s
}

func (s *ModifyClusterConfigurationRequestCustomizeConfig) SetName(v string) *ModifyClusterConfigurationRequestCustomizeConfig {
	s.Name = &v
	return s
}

type ModifyClusterConfigurationRequestCustomizeConfigConfigs struct {
	Key   *string `json:"key,omitempty" xml:"key,omitempty"`
	Value *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s ModifyClusterConfigurationRequestCustomizeConfigConfigs) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationRequestCustomizeConfigConfigs) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationRequestCustomizeConfigConfigs) SetKey(v string) *ModifyClusterConfigurationRequestCustomizeConfigConfigs {
	s.Key = &v
	return s
}

func (s *ModifyClusterConfigurationRequestCustomizeConfigConfigs) SetValue(v string) *ModifyClusterConfigurationRequestCustomizeConfigConfigs {
	s.Value = &v
	return s
}

type ModifyClusterConfigurationResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s ModifyClusterConfigurationResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterConfigurationResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterConfigurationResponse) SetHeaders(v map[string]*string) *ModifyClusterConfigurationResponse {
	s.Headers = v
	return s
}

func (s *ModifyClusterConfigurationResponse) SetStatusCode(v int32) *ModifyClusterConfigurationResponse {
	s.StatusCode = &v
	return s
}

type ModifyClusterNodePoolRequest struct {
	AutoScaling      *ModifyClusterNodePoolRequestAutoScaling      `json:"auto_scaling,omitempty" xml:"auto_scaling,omitempty" type:"Struct"`
	KubernetesConfig *ModifyClusterNodePoolRequestKubernetesConfig `json:"kubernetes_config,omitempty" xml:"kubernetes_config,omitempty" type:"Struct"`
	Management       *ModifyClusterNodePoolRequestManagement       `json:"management,omitempty" xml:"management,omitempty" type:"Struct"`
	NodepoolInfo     *ModifyClusterNodePoolRequestNodepoolInfo     `json:"nodepool_info,omitempty" xml:"nodepool_info,omitempty" type:"Struct"`
	ScalingGroup     *ModifyClusterNodePoolRequestScalingGroup     `json:"scaling_group,omitempty" xml:"scaling_group,omitempty" type:"Struct"`
	TeeConfig        *ModifyClusterNodePoolRequestTeeConfig        `json:"tee_config,omitempty" xml:"tee_config,omitempty" type:"Struct"`
	UpdateNodes      *bool                                         `json:"update_nodes,omitempty" xml:"update_nodes,omitempty"`
}

func (s ModifyClusterNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequest) SetAutoScaling(v *ModifyClusterNodePoolRequestAutoScaling) *ModifyClusterNodePoolRequest {
	s.AutoScaling = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetKubernetesConfig(v *ModifyClusterNodePoolRequestKubernetesConfig) *ModifyClusterNodePoolRequest {
	s.KubernetesConfig = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetManagement(v *ModifyClusterNodePoolRequestManagement) *ModifyClusterNodePoolRequest {
	s.Management = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetNodepoolInfo(v *ModifyClusterNodePoolRequestNodepoolInfo) *ModifyClusterNodePoolRequest {
	s.NodepoolInfo = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetScalingGroup(v *ModifyClusterNodePoolRequestScalingGroup) *ModifyClusterNodePoolRequest {
	s.ScalingGroup = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetTeeConfig(v *ModifyClusterNodePoolRequestTeeConfig) *ModifyClusterNodePoolRequest {
	s.TeeConfig = v
	return s
}

func (s *ModifyClusterNodePoolRequest) SetUpdateNodes(v bool) *ModifyClusterNodePoolRequest {
	s.UpdateNodes = &v
	return s
}

type ModifyClusterNodePoolRequestAutoScaling struct {
	EipBandwidth          *int64  `json:"eip_bandwidth,omitempty" xml:"eip_bandwidth,omitempty"`
	EipInternetChargeType *string `json:"eip_internet_charge_type,omitempty" xml:"eip_internet_charge_type,omitempty"`
	Enable                *bool   `json:"enable,omitempty" xml:"enable,omitempty"`
	IsBondEip             *bool   `json:"is_bond_eip,omitempty" xml:"is_bond_eip,omitempty"`
	MaxInstances          *int64  `json:"max_instances,omitempty" xml:"max_instances,omitempty"`
	MinInstances          *int64  `json:"min_instances,omitempty" xml:"min_instances,omitempty"`
	Type                  *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s ModifyClusterNodePoolRequestAutoScaling) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestAutoScaling) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetEipBandwidth(v int64) *ModifyClusterNodePoolRequestAutoScaling {
	s.EipBandwidth = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetEipInternetChargeType(v string) *ModifyClusterNodePoolRequestAutoScaling {
	s.EipInternetChargeType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetEnable(v bool) *ModifyClusterNodePoolRequestAutoScaling {
	s.Enable = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetIsBondEip(v bool) *ModifyClusterNodePoolRequestAutoScaling {
	s.IsBondEip = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetMaxInstances(v int64) *ModifyClusterNodePoolRequestAutoScaling {
	s.MaxInstances = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetMinInstances(v int64) *ModifyClusterNodePoolRequestAutoScaling {
	s.MinInstances = &v
	return s
}

func (s *ModifyClusterNodePoolRequestAutoScaling) SetType(v string) *ModifyClusterNodePoolRequestAutoScaling {
	s.Type = &v
	return s
}

type ModifyClusterNodePoolRequestKubernetesConfig struct {
	CmsEnabled     *bool    `json:"cms_enabled,omitempty" xml:"cms_enabled,omitempty"`
	CpuPolicy      *string  `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	Labels         []*Tag   `json:"labels,omitempty" xml:"labels,omitempty" type:"Repeated"`
	Runtime        *string  `json:"runtime,omitempty" xml:"runtime,omitempty"`
	RuntimeVersion *string  `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
	Taints         []*Taint `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	UserData       *string  `json:"user_data,omitempty" xml:"user_data,omitempty"`
}

func (s ModifyClusterNodePoolRequestKubernetesConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestKubernetesConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetCmsEnabled(v bool) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.CmsEnabled = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetCpuPolicy(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.CpuPolicy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetLabels(v []*Tag) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.Labels = v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetRuntime(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.Runtime = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetRuntimeVersion(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.RuntimeVersion = &v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetTaints(v []*Taint) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.Taints = v
	return s
}

func (s *ModifyClusterNodePoolRequestKubernetesConfig) SetUserData(v string) *ModifyClusterNodePoolRequestKubernetesConfig {
	s.UserData = &v
	return s
}

type ModifyClusterNodePoolRequestManagement struct {
	AutoRepair    *bool                                                `json:"auto_repair,omitempty" xml:"auto_repair,omitempty"`
	Enable        *bool                                                `json:"enable,omitempty" xml:"enable,omitempty"`
	UpgradeConfig *ModifyClusterNodePoolRequestManagementUpgradeConfig `json:"upgrade_config,omitempty" xml:"upgrade_config,omitempty" type:"Struct"`
}

func (s ModifyClusterNodePoolRequestManagement) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestManagement) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestManagement) SetAutoRepair(v bool) *ModifyClusterNodePoolRequestManagement {
	s.AutoRepair = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagement) SetEnable(v bool) *ModifyClusterNodePoolRequestManagement {
	s.Enable = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagement) SetUpgradeConfig(v *ModifyClusterNodePoolRequestManagementUpgradeConfig) *ModifyClusterNodePoolRequestManagement {
	s.UpgradeConfig = v
	return s
}

type ModifyClusterNodePoolRequestManagementUpgradeConfig struct {
	AutoUpgrade     *bool  `json:"auto_upgrade,omitempty" xml:"auto_upgrade,omitempty"`
	MaxUnavailable  *int64 `json:"max_unavailable,omitempty" xml:"max_unavailable,omitempty"`
	Surge           *int64 `json:"surge,omitempty" xml:"surge,omitempty"`
	SurgePercentage *int64 `json:"surge_percentage,omitempty" xml:"surge_percentage,omitempty"`
}

func (s ModifyClusterNodePoolRequestManagementUpgradeConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestManagementUpgradeConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetAutoUpgrade(v bool) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.AutoUpgrade = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetMaxUnavailable(v int64) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.MaxUnavailable = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetSurge(v int64) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.Surge = &v
	return s
}

func (s *ModifyClusterNodePoolRequestManagementUpgradeConfig) SetSurgePercentage(v int64) *ModifyClusterNodePoolRequestManagementUpgradeConfig {
	s.SurgePercentage = &v
	return s
}

type ModifyClusterNodePoolRequestNodepoolInfo struct {
	Name            *string `json:"name,omitempty" xml:"name,omitempty"`
	ResourceGroupId *string `json:"resource_group_id,omitempty" xml:"resource_group_id,omitempty"`
}

func (s ModifyClusterNodePoolRequestNodepoolInfo) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestNodepoolInfo) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestNodepoolInfo) SetName(v string) *ModifyClusterNodePoolRequestNodepoolInfo {
	s.Name = &v
	return s
}

func (s *ModifyClusterNodePoolRequestNodepoolInfo) SetResourceGroupId(v string) *ModifyClusterNodePoolRequestNodepoolInfo {
	s.ResourceGroupId = &v
	return s
}

type ModifyClusterNodePoolRequestScalingGroup struct {
	AutoRenew                           *bool                                                       `json:"auto_renew,omitempty" xml:"auto_renew,omitempty"`
	AutoRenewPeriod                     *int64                                                      `json:"auto_renew_period,omitempty" xml:"auto_renew_period,omitempty"`
	CompensateWithOnDemand              *bool                                                       `json:"compensate_with_on_demand,omitempty" xml:"compensate_with_on_demand,omitempty"`
	DataDisks                           []*DataDisk                                                 `json:"data_disks,omitempty" xml:"data_disks,omitempty" type:"Repeated"`
	DesiredSize                         *int64                                                      `json:"desired_size,omitempty" xml:"desired_size,omitempty"`
	ImageId                             *string                                                     `json:"image_id,omitempty" xml:"image_id,omitempty"`
	InstanceChargeType                  *string                                                     `json:"instance_charge_type,omitempty" xml:"instance_charge_type,omitempty"`
	InstanceTypes                       []*string                                                   `json:"instance_types,omitempty" xml:"instance_types,omitempty" type:"Repeated"`
	InternetChargeType                  *string                                                     `json:"internet_charge_type,omitempty" xml:"internet_charge_type,omitempty"`
	InternetMaxBandwidthOut             *int64                                                      `json:"internet_max_bandwidth_out,omitempty" xml:"internet_max_bandwidth_out,omitempty"`
	KeyPair                             *string                                                     `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	LoginPassword                       *string                                                     `json:"login_password,omitempty" xml:"login_password,omitempty"`
	MultiAzPolicy                       *string                                                     `json:"multi_az_policy,omitempty" xml:"multi_az_policy,omitempty"`
	OnDemandBaseCapacity                *int64                                                      `json:"on_demand_base_capacity,omitempty" xml:"on_demand_base_capacity,omitempty"`
	OnDemandPercentageAboveBaseCapacity *int64                                                      `json:"on_demand_percentage_above_base_capacity,omitempty" xml:"on_demand_percentage_above_base_capacity,omitempty"`
	Period                              *int64                                                      `json:"period,omitempty" xml:"period,omitempty"`
	PeriodUnit                          *string                                                     `json:"period_unit,omitempty" xml:"period_unit,omitempty"`
	Platform                            *string                                                     `json:"platform,omitempty" xml:"platform,omitempty"`
	PrivatePoolOptions                  *ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions `json:"private_pool_options,omitempty" xml:"private_pool_options,omitempty" type:"Struct"`
	RdsInstances                        []*string                                                   `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	ScalingPolicy                       *string                                                     `json:"scaling_policy,omitempty" xml:"scaling_policy,omitempty"`
	SpotInstancePools                   *int64                                                      `json:"spot_instance_pools,omitempty" xml:"spot_instance_pools,omitempty"`
	SpotInstanceRemedy                  *bool                                                       `json:"spot_instance_remedy,omitempty" xml:"spot_instance_remedy,omitempty"`
	SpotPriceLimit                      []*ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit   `json:"spot_price_limit,omitempty" xml:"spot_price_limit,omitempty" type:"Repeated"`
	SpotStrategy                        *string                                                     `json:"spot_strategy,omitempty" xml:"spot_strategy,omitempty"`
	SystemDiskCategory                  *string                                                     `json:"system_disk_category,omitempty" xml:"system_disk_category,omitempty"`
	SystemDiskPerformanceLevel          *string                                                     `json:"system_disk_performance_level,omitempty" xml:"system_disk_performance_level,omitempty"`
	SystemDiskSize                      *int64                                                      `json:"system_disk_size,omitempty" xml:"system_disk_size,omitempty"`
	Tags                                []*Tag                                                      `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	VswitchIds                          []*string                                                   `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
}

func (s ModifyClusterNodePoolRequestScalingGroup) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestScalingGroup) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetAutoRenew(v bool) *ModifyClusterNodePoolRequestScalingGroup {
	s.AutoRenew = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetAutoRenewPeriod(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.AutoRenewPeriod = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetCompensateWithOnDemand(v bool) *ModifyClusterNodePoolRequestScalingGroup {
	s.CompensateWithOnDemand = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetDataDisks(v []*DataDisk) *ModifyClusterNodePoolRequestScalingGroup {
	s.DataDisks = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetDesiredSize(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.DesiredSize = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetImageId(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.ImageId = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInstanceChargeType(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.InstanceChargeType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInstanceTypes(v []*string) *ModifyClusterNodePoolRequestScalingGroup {
	s.InstanceTypes = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInternetChargeType(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.InternetChargeType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetInternetMaxBandwidthOut(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.InternetMaxBandwidthOut = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetKeyPair(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.KeyPair = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetLoginPassword(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.LoginPassword = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetMultiAzPolicy(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.MultiAzPolicy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetOnDemandBaseCapacity(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.OnDemandBaseCapacity = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetOnDemandPercentageAboveBaseCapacity(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.OnDemandPercentageAboveBaseCapacity = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetPeriod(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.Period = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetPeriodUnit(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.PeriodUnit = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetPlatform(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.Platform = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetPrivatePoolOptions(v *ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions) *ModifyClusterNodePoolRequestScalingGroup {
	s.PrivatePoolOptions = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetRdsInstances(v []*string) *ModifyClusterNodePoolRequestScalingGroup {
	s.RdsInstances = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetScalingPolicy(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.ScalingPolicy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotInstancePools(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotInstancePools = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotInstanceRemedy(v bool) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotInstanceRemedy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotPriceLimit(v []*ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotPriceLimit = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSpotStrategy(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.SpotStrategy = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSystemDiskCategory(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.SystemDiskCategory = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSystemDiskPerformanceLevel(v string) *ModifyClusterNodePoolRequestScalingGroup {
	s.SystemDiskPerformanceLevel = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetSystemDiskSize(v int64) *ModifyClusterNodePoolRequestScalingGroup {
	s.SystemDiskSize = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetTags(v []*Tag) *ModifyClusterNodePoolRequestScalingGroup {
	s.Tags = v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroup) SetVswitchIds(v []*string) *ModifyClusterNodePoolRequestScalingGroup {
	s.VswitchIds = v
	return s
}

type ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions struct {
	Id            *string `json:"id,omitempty" xml:"id,omitempty"`
	MatchCriteria *string `json:"match_criteria,omitempty" xml:"match_criteria,omitempty"`
}

func (s ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions) SetId(v string) *ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions {
	s.Id = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions) SetMatchCriteria(v string) *ModifyClusterNodePoolRequestScalingGroupPrivatePoolOptions {
	s.MatchCriteria = &v
	return s
}

type ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit struct {
	InstanceType *string `json:"instance_type,omitempty" xml:"instance_type,omitempty"`
	PriceLimit   *string `json:"price_limit,omitempty" xml:"price_limit,omitempty"`
}

func (s ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) SetInstanceType(v string) *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.InstanceType = &v
	return s
}

func (s *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit) SetPriceLimit(v string) *ModifyClusterNodePoolRequestScalingGroupSpotPriceLimit {
	s.PriceLimit = &v
	return s
}

type ModifyClusterNodePoolRequestTeeConfig struct {
	TeeEnable *bool `json:"tee_enable,omitempty" xml:"tee_enable,omitempty"`
}

func (s ModifyClusterNodePoolRequestTeeConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolRequestTeeConfig) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolRequestTeeConfig) SetTeeEnable(v bool) *ModifyClusterNodePoolRequestTeeConfig {
	s.TeeEnable = &v
	return s
}

type ModifyClusterNodePoolResponseBody struct {
	// The ID of the node pool.
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ModifyClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolResponseBody) SetNodepoolId(v string) *ModifyClusterNodePoolResponseBody {
	s.NodepoolId = &v
	return s
}

func (s *ModifyClusterNodePoolResponseBody) SetTaskId(v string) *ModifyClusterNodePoolResponseBody {
	s.TaskId = &v
	return s
}

type ModifyClusterNodePoolResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterNodePoolResponse) SetHeaders(v map[string]*string) *ModifyClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *ModifyClusterNodePoolResponse) SetStatusCode(v int32) *ModifyClusterNodePoolResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyClusterNodePoolResponse) SetBody(v *ModifyClusterNodePoolResponseBody) *ModifyClusterNodePoolResponse {
	s.Body = v
	return s
}

type ModifyClusterTagsRequest struct {
	Body []*Tag `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s ModifyClusterTagsRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterTagsRequest) GoString() string {
	return s.String()
}

func (s *ModifyClusterTagsRequest) SetBody(v []*Tag) *ModifyClusterTagsRequest {
	s.Body = v
	return s
}

type ModifyClusterTagsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s ModifyClusterTagsResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyClusterTagsResponse) GoString() string {
	return s.String()
}

func (s *ModifyClusterTagsResponse) SetHeaders(v map[string]*string) *ModifyClusterTagsResponse {
	s.Headers = v
	return s
}

func (s *ModifyClusterTagsResponse) SetStatusCode(v int32) *ModifyClusterTagsResponse {
	s.StatusCode = &v
	return s
}

type ModifyNodePoolNodeConfigRequest struct {
	KubeletConfig *ModifyNodePoolNodeConfigRequestKubeletConfig `json:"kubelet_config,omitempty" xml:"kubelet_config,omitempty" type:"Struct"`
	RollingPolicy *ModifyNodePoolNodeConfigRequestRollingPolicy `json:"rolling_policy,omitempty" xml:"rolling_policy,omitempty" type:"Struct"`
}

func (s ModifyNodePoolNodeConfigRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyNodePoolNodeConfigRequest) GoString() string {
	return s.String()
}

func (s *ModifyNodePoolNodeConfigRequest) SetKubeletConfig(v *ModifyNodePoolNodeConfigRequestKubeletConfig) *ModifyNodePoolNodeConfigRequest {
	s.KubeletConfig = v
	return s
}

func (s *ModifyNodePoolNodeConfigRequest) SetRollingPolicy(v *ModifyNodePoolNodeConfigRequestRollingPolicy) *ModifyNodePoolNodeConfigRequest {
	s.RollingPolicy = v
	return s
}

type ModifyNodePoolNodeConfigRequestKubeletConfig struct {
	CpuManagerPolicy        *string                `json:"cpuManagerPolicy,omitempty" xml:"cpuManagerPolicy,omitempty"`
	EventBurst              *int64                 `json:"eventBurst,omitempty" xml:"eventBurst,omitempty"`
	EventRecordQPS          *int64                 `json:"eventRecordQPS,omitempty" xml:"eventRecordQPS,omitempty"`
	EvictionHard            map[string]interface{} `json:"evictionHard,omitempty" xml:"evictionHard,omitempty"`
	EvictionSoft            map[string]interface{} `json:"evictionSoft,omitempty" xml:"evictionSoft,omitempty"`
	EvictionSoftGracePeriod map[string]interface{} `json:"evictionSoftGracePeriod,omitempty" xml:"evictionSoftGracePeriod,omitempty"`
	KubeAPIBurst            *int64                 `json:"kubeAPIBurst,omitempty" xml:"kubeAPIBurst,omitempty"`
	KubeAPIQPS              *int64                 `json:"kubeAPIQPS,omitempty" xml:"kubeAPIQPS,omitempty"`
	KubeReserved            map[string]interface{} `json:"kubeReserved,omitempty" xml:"kubeReserved,omitempty"`
	RegistryBurst           *int64                 `json:"registryBurst,omitempty" xml:"registryBurst,omitempty"`
	RegistryPullQPS         *int64                 `json:"registryPullQPS,omitempty" xml:"registryPullQPS,omitempty"`
	SerializeImagePulls     *bool                  `json:"serializeImagePulls,omitempty" xml:"serializeImagePulls,omitempty"`
	SystemReserved          map[string]interface{} `json:"systemReserved,omitempty" xml:"systemReserved,omitempty"`
}

func (s ModifyNodePoolNodeConfigRequestKubeletConfig) String() string {
	return tea.Prettify(s)
}

func (s ModifyNodePoolNodeConfigRequestKubeletConfig) GoString() string {
	return s.String()
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetCpuManagerPolicy(v string) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.CpuManagerPolicy = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetEventBurst(v int64) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.EventBurst = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetEventRecordQPS(v int64) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.EventRecordQPS = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetEvictionHard(v map[string]interface{}) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.EvictionHard = v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetEvictionSoft(v map[string]interface{}) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.EvictionSoft = v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetEvictionSoftGracePeriod(v map[string]interface{}) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.EvictionSoftGracePeriod = v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetKubeAPIBurst(v int64) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.KubeAPIBurst = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetKubeAPIQPS(v int64) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.KubeAPIQPS = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetKubeReserved(v map[string]interface{}) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.KubeReserved = v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetRegistryBurst(v int64) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.RegistryBurst = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetRegistryPullQPS(v int64) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.RegistryPullQPS = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetSerializeImagePulls(v bool) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.SerializeImagePulls = &v
	return s
}

func (s *ModifyNodePoolNodeConfigRequestKubeletConfig) SetSystemReserved(v map[string]interface{}) *ModifyNodePoolNodeConfigRequestKubeletConfig {
	s.SystemReserved = v
	return s
}

type ModifyNodePoolNodeConfigRequestRollingPolicy struct {
	MaxParallelism *int64 `json:"max_parallelism,omitempty" xml:"max_parallelism,omitempty"`
}

func (s ModifyNodePoolNodeConfigRequestRollingPolicy) String() string {
	return tea.Prettify(s)
}

func (s ModifyNodePoolNodeConfigRequestRollingPolicy) GoString() string {
	return s.String()
}

func (s *ModifyNodePoolNodeConfigRequestRollingPolicy) SetMaxParallelism(v int64) *ModifyNodePoolNodeConfigRequestRollingPolicy {
	s.MaxParallelism = &v
	return s
}

type ModifyNodePoolNodeConfigResponseBody struct {
	NodepoolId *string `json:"nodepool_id,omitempty" xml:"nodepool_id,omitempty"`
	RequestId  *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	TaskId     *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ModifyNodePoolNodeConfigResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyNodePoolNodeConfigResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyNodePoolNodeConfigResponseBody) SetNodepoolId(v string) *ModifyNodePoolNodeConfigResponseBody {
	s.NodepoolId = &v
	return s
}

func (s *ModifyNodePoolNodeConfigResponseBody) SetRequestId(v string) *ModifyNodePoolNodeConfigResponseBody {
	s.RequestId = &v
	return s
}

func (s *ModifyNodePoolNodeConfigResponseBody) SetTaskId(v string) *ModifyNodePoolNodeConfigResponseBody {
	s.TaskId = &v
	return s
}

type ModifyNodePoolNodeConfigResponse struct {
	Headers    map[string]*string                    `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                                `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyNodePoolNodeConfigResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyNodePoolNodeConfigResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyNodePoolNodeConfigResponse) GoString() string {
	return s.String()
}

func (s *ModifyNodePoolNodeConfigResponse) SetHeaders(v map[string]*string) *ModifyNodePoolNodeConfigResponse {
	s.Headers = v
	return s
}

func (s *ModifyNodePoolNodeConfigResponse) SetStatusCode(v int32) *ModifyNodePoolNodeConfigResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyNodePoolNodeConfigResponse) SetBody(v *ModifyNodePoolNodeConfigResponseBody) *ModifyNodePoolNodeConfigResponse {
	s.Body = v
	return s
}

type ModifyPolicyInstanceRequest struct {
	Action       *string                `json:"action,omitempty" xml:"action,omitempty"`
	InstanceName *string                `json:"instance_name,omitempty" xml:"instance_name,omitempty"`
	Namespaces   []*string              `json:"namespaces,omitempty" xml:"namespaces,omitempty" type:"Repeated"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" xml:"parameters,omitempty"`
}

func (s ModifyPolicyInstanceRequest) String() string {
	return tea.Prettify(s)
}

func (s ModifyPolicyInstanceRequest) GoString() string {
	return s.String()
}

func (s *ModifyPolicyInstanceRequest) SetAction(v string) *ModifyPolicyInstanceRequest {
	s.Action = &v
	return s
}

func (s *ModifyPolicyInstanceRequest) SetInstanceName(v string) *ModifyPolicyInstanceRequest {
	s.InstanceName = &v
	return s
}

func (s *ModifyPolicyInstanceRequest) SetNamespaces(v []*string) *ModifyPolicyInstanceRequest {
	s.Namespaces = v
	return s
}

func (s *ModifyPolicyInstanceRequest) SetParameters(v map[string]interface{}) *ModifyPolicyInstanceRequest {
	s.Parameters = v
	return s
}

type ModifyPolicyInstanceResponseBody struct {
	// The policy instance that is updated.
	Instances []*string `json:"instances,omitempty" xml:"instances,omitempty" type:"Repeated"`
}

func (s ModifyPolicyInstanceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ModifyPolicyInstanceResponseBody) GoString() string {
	return s.String()
}

func (s *ModifyPolicyInstanceResponseBody) SetInstances(v []*string) *ModifyPolicyInstanceResponseBody {
	s.Instances = v
	return s
}

type ModifyPolicyInstanceResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ModifyPolicyInstanceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ModifyPolicyInstanceResponse) String() string {
	return tea.Prettify(s)
}

func (s ModifyPolicyInstanceResponse) GoString() string {
	return s.String()
}

func (s *ModifyPolicyInstanceResponse) SetHeaders(v map[string]*string) *ModifyPolicyInstanceResponse {
	s.Headers = v
	return s
}

func (s *ModifyPolicyInstanceResponse) SetStatusCode(v int32) *ModifyPolicyInstanceResponse {
	s.StatusCode = &v
	return s
}

func (s *ModifyPolicyInstanceResponse) SetBody(v *ModifyPolicyInstanceResponseBody) *ModifyPolicyInstanceResponse {
	s.Body = v
	return s
}

type OpenAckServiceRequest struct {
	Type *string `json:"type,omitempty" xml:"type,omitempty"`
}

func (s OpenAckServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s OpenAckServiceRequest) GoString() string {
	return s.String()
}

func (s *OpenAckServiceRequest) SetType(v string) *OpenAckServiceRequest {
	s.Type = &v
	return s
}

type OpenAckServiceResponseBody struct {
	// The ID of the order.
	OrderId *string `json:"order_id,omitempty" xml:"order_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
}

func (s OpenAckServiceResponseBody) String() string {
	return tea.Prettify(s)
}

func (s OpenAckServiceResponseBody) GoString() string {
	return s.String()
}

func (s *OpenAckServiceResponseBody) SetOrderId(v string) *OpenAckServiceResponseBody {
	s.OrderId = &v
	return s
}

func (s *OpenAckServiceResponseBody) SetRequestId(v string) *OpenAckServiceResponseBody {
	s.RequestId = &v
	return s
}

type OpenAckServiceResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *OpenAckServiceResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s OpenAckServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s OpenAckServiceResponse) GoString() string {
	return s.String()
}

func (s *OpenAckServiceResponse) SetHeaders(v map[string]*string) *OpenAckServiceResponse {
	s.Headers = v
	return s
}

func (s *OpenAckServiceResponse) SetStatusCode(v int32) *OpenAckServiceResponse {
	s.StatusCode = &v
	return s
}

func (s *OpenAckServiceResponse) SetBody(v *OpenAckServiceResponseBody) *OpenAckServiceResponse {
	s.Body = v
	return s
}

type PauseClusterUpgradeResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s PauseClusterUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s PauseClusterUpgradeResponse) GoString() string {
	return s.String()
}

func (s *PauseClusterUpgradeResponse) SetHeaders(v map[string]*string) *PauseClusterUpgradeResponse {
	s.Headers = v
	return s
}

func (s *PauseClusterUpgradeResponse) SetStatusCode(v int32) *PauseClusterUpgradeResponse {
	s.StatusCode = &v
	return s
}

type PauseComponentUpgradeResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s PauseComponentUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s PauseComponentUpgradeResponse) GoString() string {
	return s.String()
}

func (s *PauseComponentUpgradeResponse) SetHeaders(v map[string]*string) *PauseComponentUpgradeResponse {
	s.Headers = v
	return s
}

func (s *PauseComponentUpgradeResponse) SetStatusCode(v int32) *PauseComponentUpgradeResponse {
	s.StatusCode = &v
	return s
}

type PauseTaskResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s PauseTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s PauseTaskResponse) GoString() string {
	return s.String()
}

func (s *PauseTaskResponse) SetHeaders(v map[string]*string) *PauseTaskResponse {
	s.Headers = v
	return s
}

func (s *PauseTaskResponse) SetStatusCode(v int32) *PauseTaskResponse {
	s.StatusCode = &v
	return s
}

type RemoveClusterNodesRequest struct {
	DrainNode   *bool     `json:"drain_node,omitempty" xml:"drain_node,omitempty"`
	Nodes       []*string `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
	ReleaseNode *bool     `json:"release_node,omitempty" xml:"release_node,omitempty"`
}

func (s RemoveClusterNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s RemoveClusterNodesRequest) GoString() string {
	return s.String()
}

func (s *RemoveClusterNodesRequest) SetDrainNode(v bool) *RemoveClusterNodesRequest {
	s.DrainNode = &v
	return s
}

func (s *RemoveClusterNodesRequest) SetNodes(v []*string) *RemoveClusterNodesRequest {
	s.Nodes = v
	return s
}

func (s *RemoveClusterNodesRequest) SetReleaseNode(v bool) *RemoveClusterNodesRequest {
	s.ReleaseNode = &v
	return s
}

type RemoveClusterNodesResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s RemoveClusterNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s RemoveClusterNodesResponse) GoString() string {
	return s.String()
}

func (s *RemoveClusterNodesResponse) SetHeaders(v map[string]*string) *RemoveClusterNodesResponse {
	s.Headers = v
	return s
}

func (s *RemoveClusterNodesResponse) SetStatusCode(v int32) *RemoveClusterNodesResponse {
	s.StatusCode = &v
	return s
}

type RemoveNodePoolNodesRequest struct {
	// true
	DrainNode *bool `json:"drain_node,omitempty" xml:"drain_node,omitempty"`
	// i-bp1c70fqbv1nlu9xxxxx
	InstanceIds []*string `json:"instance_ids,omitempty" xml:"instance_ids,omitempty" type:"Repeated"`
	// cn-hangzhou.172.16.xxx.xxx
	Nodes []*string `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
	// true
	ReleaseNode *bool `json:"release_node,omitempty" xml:"release_node,omitempty"`
}

func (s RemoveNodePoolNodesRequest) String() string {
	return tea.Prettify(s)
}

func (s RemoveNodePoolNodesRequest) GoString() string {
	return s.String()
}

func (s *RemoveNodePoolNodesRequest) SetDrainNode(v bool) *RemoveNodePoolNodesRequest {
	s.DrainNode = &v
	return s
}

func (s *RemoveNodePoolNodesRequest) SetInstanceIds(v []*string) *RemoveNodePoolNodesRequest {
	s.InstanceIds = v
	return s
}

func (s *RemoveNodePoolNodesRequest) SetNodes(v []*string) *RemoveNodePoolNodesRequest {
	s.Nodes = v
	return s
}

func (s *RemoveNodePoolNodesRequest) SetReleaseNode(v bool) *RemoveNodePoolNodesRequest {
	s.ReleaseNode = &v
	return s
}

type RemoveNodePoolNodesShrinkRequest struct {
	// true
	DrainNode *bool `json:"drain_node,omitempty" xml:"drain_node,omitempty"`
	// i-bp1c70fqbv1nlu9xxxxx
	InstanceIdsShrink *string `json:"instance_ids,omitempty" xml:"instance_ids,omitempty"`
	// cn-hangzhou.172.16.xxx.xxx
	NodesShrink *string `json:"nodes,omitempty" xml:"nodes,omitempty"`
	// true
	ReleaseNode *bool `json:"release_node,omitempty" xml:"release_node,omitempty"`
}

func (s RemoveNodePoolNodesShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s RemoveNodePoolNodesShrinkRequest) GoString() string {
	return s.String()
}

func (s *RemoveNodePoolNodesShrinkRequest) SetDrainNode(v bool) *RemoveNodePoolNodesShrinkRequest {
	s.DrainNode = &v
	return s
}

func (s *RemoveNodePoolNodesShrinkRequest) SetInstanceIdsShrink(v string) *RemoveNodePoolNodesShrinkRequest {
	s.InstanceIdsShrink = &v
	return s
}

func (s *RemoveNodePoolNodesShrinkRequest) SetNodesShrink(v string) *RemoveNodePoolNodesShrinkRequest {
	s.NodesShrink = &v
	return s
}

func (s *RemoveNodePoolNodesShrinkRequest) SetReleaseNode(v bool) *RemoveNodePoolNodesShrinkRequest {
	s.ReleaseNode = &v
	return s
}

type RemoveNodePoolNodesResponseBody struct {
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s RemoveNodePoolNodesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RemoveNodePoolNodesResponseBody) GoString() string {
	return s.String()
}

func (s *RemoveNodePoolNodesResponseBody) SetRequestId(v string) *RemoveNodePoolNodesResponseBody {
	s.RequestId = &v
	return s
}

func (s *RemoveNodePoolNodesResponseBody) SetTaskId(v string) *RemoveNodePoolNodesResponseBody {
	s.TaskId = &v
	return s
}

type RemoveNodePoolNodesResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RemoveNodePoolNodesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RemoveNodePoolNodesResponse) String() string {
	return tea.Prettify(s)
}

func (s RemoveNodePoolNodesResponse) GoString() string {
	return s.String()
}

func (s *RemoveNodePoolNodesResponse) SetHeaders(v map[string]*string) *RemoveNodePoolNodesResponse {
	s.Headers = v
	return s
}

func (s *RemoveNodePoolNodesResponse) SetStatusCode(v int32) *RemoveNodePoolNodesResponse {
	s.StatusCode = &v
	return s
}

func (s *RemoveNodePoolNodesResponse) SetBody(v *RemoveNodePoolNodesResponseBody) *RemoveNodePoolNodesResponse {
	s.Body = v
	return s
}

type RemoveWorkflowResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s RemoveWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s RemoveWorkflowResponse) GoString() string {
	return s.String()
}

func (s *RemoveWorkflowResponse) SetHeaders(v map[string]*string) *RemoveWorkflowResponse {
	s.Headers = v
	return s
}

func (s *RemoveWorkflowResponse) SetStatusCode(v int32) *RemoveWorkflowResponse {
	s.StatusCode = &v
	return s
}

type RepairClusterNodePoolRequest struct {
	Nodes []*string `json:"nodes,omitempty" xml:"nodes,omitempty" type:"Repeated"`
}

func (s RepairClusterNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s RepairClusterNodePoolRequest) GoString() string {
	return s.String()
}

func (s *RepairClusterNodePoolRequest) SetNodes(v []*string) *RepairClusterNodePoolRequest {
	s.Nodes = v
	return s
}

type RepairClusterNodePoolResponseBody struct {
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s RepairClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s RepairClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *RepairClusterNodePoolResponseBody) SetRequestId(v string) *RepairClusterNodePoolResponseBody {
	s.RequestId = &v
	return s
}

func (s *RepairClusterNodePoolResponseBody) SetTaskId(v string) *RepairClusterNodePoolResponseBody {
	s.TaskId = &v
	return s
}

type RepairClusterNodePoolResponse struct {
	Headers    map[string]*string                 `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *RepairClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s RepairClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s RepairClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *RepairClusterNodePoolResponse) SetHeaders(v map[string]*string) *RepairClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *RepairClusterNodePoolResponse) SetStatusCode(v int32) *RepairClusterNodePoolResponse {
	s.StatusCode = &v
	return s
}

func (s *RepairClusterNodePoolResponse) SetBody(v *RepairClusterNodePoolResponseBody) *RepairClusterNodePoolResponse {
	s.Body = v
	return s
}

type ResumeComponentUpgradeResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s ResumeComponentUpgradeResponse) String() string {
	return tea.Prettify(s)
}

func (s ResumeComponentUpgradeResponse) GoString() string {
	return s.String()
}

func (s *ResumeComponentUpgradeResponse) SetHeaders(v map[string]*string) *ResumeComponentUpgradeResponse {
	s.Headers = v
	return s
}

func (s *ResumeComponentUpgradeResponse) SetStatusCode(v int32) *ResumeComponentUpgradeResponse {
	s.StatusCode = &v
	return s
}

type ResumeTaskResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s ResumeTaskResponse) String() string {
	return tea.Prettify(s)
}

func (s ResumeTaskResponse) GoString() string {
	return s.String()
}

func (s *ResumeTaskResponse) SetHeaders(v map[string]*string) *ResumeTaskResponse {
	s.Headers = v
	return s
}

func (s *ResumeTaskResponse) SetStatusCode(v int32) *ResumeTaskResponse {
	s.StatusCode = &v
	return s
}

type ResumeUpgradeClusterResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s ResumeUpgradeClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ResumeUpgradeClusterResponse) GoString() string {
	return s.String()
}

func (s *ResumeUpgradeClusterResponse) SetHeaders(v map[string]*string) *ResumeUpgradeClusterResponse {
	s.Headers = v
	return s
}

func (s *ResumeUpgradeClusterResponse) SetStatusCode(v int32) *ResumeUpgradeClusterResponse {
	s.StatusCode = &v
	return s
}

type ScaleClusterRequest struct {
	CloudMonitorFlags        *bool                                 `json:"cloud_monitor_flags,omitempty" xml:"cloud_monitor_flags,omitempty"`
	Count                    *int64                                `json:"count,omitempty" xml:"count,omitempty"`
	CpuPolicy                *string                               `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	DisableRollback          *bool                                 `json:"disable_rollback,omitempty" xml:"disable_rollback,omitempty"`
	KeyPair                  *string                               `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	LoginPassword            *string                               `json:"login_password,omitempty" xml:"login_password,omitempty"`
	Tags                     []*ScaleClusterRequestTags            `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	Taints                   []*ScaleClusterRequestTaints          `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	VswitchIds               []*string                             `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	WorkerAutoRenew          *bool                                 `json:"worker_auto_renew,omitempty" xml:"worker_auto_renew,omitempty"`
	WorkerAutoRenewPeriod    *int64                                `json:"worker_auto_renew_period,omitempty" xml:"worker_auto_renew_period,omitempty"`
	WorkerDataDisk           *bool                                 `json:"worker_data_disk,omitempty" xml:"worker_data_disk,omitempty"`
	WorkerDataDisks          []*ScaleClusterRequestWorkerDataDisks `json:"worker_data_disks,omitempty" xml:"worker_data_disks,omitempty" type:"Repeated"`
	WorkerInstanceChargeType *string                               `json:"worker_instance_charge_type,omitempty" xml:"worker_instance_charge_type,omitempty"`
	WorkerInstanceTypes      []*string                             `json:"worker_instance_types,omitempty" xml:"worker_instance_types,omitempty" type:"Repeated"`
	WorkerPeriod             *int64                                `json:"worker_period,omitempty" xml:"worker_period,omitempty"`
	WorkerPeriodUnit         *string                               `json:"worker_period_unit,omitempty" xml:"worker_period_unit,omitempty"`
	WorkerSystemDiskCategory *string                               `json:"worker_system_disk_category,omitempty" xml:"worker_system_disk_category,omitempty"`
	WorkerSystemDiskSize     *int64                                `json:"worker_system_disk_size,omitempty" xml:"worker_system_disk_size,omitempty"`
}

func (s ScaleClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequest) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequest) SetCloudMonitorFlags(v bool) *ScaleClusterRequest {
	s.CloudMonitorFlags = &v
	return s
}

func (s *ScaleClusterRequest) SetCount(v int64) *ScaleClusterRequest {
	s.Count = &v
	return s
}

func (s *ScaleClusterRequest) SetCpuPolicy(v string) *ScaleClusterRequest {
	s.CpuPolicy = &v
	return s
}

func (s *ScaleClusterRequest) SetDisableRollback(v bool) *ScaleClusterRequest {
	s.DisableRollback = &v
	return s
}

func (s *ScaleClusterRequest) SetKeyPair(v string) *ScaleClusterRequest {
	s.KeyPair = &v
	return s
}

func (s *ScaleClusterRequest) SetLoginPassword(v string) *ScaleClusterRequest {
	s.LoginPassword = &v
	return s
}

func (s *ScaleClusterRequest) SetTags(v []*ScaleClusterRequestTags) *ScaleClusterRequest {
	s.Tags = v
	return s
}

func (s *ScaleClusterRequest) SetTaints(v []*ScaleClusterRequestTaints) *ScaleClusterRequest {
	s.Taints = v
	return s
}

func (s *ScaleClusterRequest) SetVswitchIds(v []*string) *ScaleClusterRequest {
	s.VswitchIds = v
	return s
}

func (s *ScaleClusterRequest) SetWorkerAutoRenew(v bool) *ScaleClusterRequest {
	s.WorkerAutoRenew = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerAutoRenewPeriod(v int64) *ScaleClusterRequest {
	s.WorkerAutoRenewPeriod = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerDataDisk(v bool) *ScaleClusterRequest {
	s.WorkerDataDisk = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerDataDisks(v []*ScaleClusterRequestWorkerDataDisks) *ScaleClusterRequest {
	s.WorkerDataDisks = v
	return s
}

func (s *ScaleClusterRequest) SetWorkerInstanceChargeType(v string) *ScaleClusterRequest {
	s.WorkerInstanceChargeType = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerInstanceTypes(v []*string) *ScaleClusterRequest {
	s.WorkerInstanceTypes = v
	return s
}

func (s *ScaleClusterRequest) SetWorkerPeriod(v int64) *ScaleClusterRequest {
	s.WorkerPeriod = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerPeriodUnit(v string) *ScaleClusterRequest {
	s.WorkerPeriodUnit = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerSystemDiskCategory(v string) *ScaleClusterRequest {
	s.WorkerSystemDiskCategory = &v
	return s
}

func (s *ScaleClusterRequest) SetWorkerSystemDiskSize(v int64) *ScaleClusterRequest {
	s.WorkerSystemDiskSize = &v
	return s
}

type ScaleClusterRequestTags struct {
	Key *string `json:"key,omitempty" xml:"key,omitempty"`
}

func (s ScaleClusterRequestTags) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequestTags) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequestTags) SetKey(v string) *ScaleClusterRequestTags {
	s.Key = &v
	return s
}

type ScaleClusterRequestTaints struct {
	Effect *string `json:"effect,omitempty" xml:"effect,omitempty"`
	Key    *string `json:"key,omitempty" xml:"key,omitempty"`
	Value  *string `json:"value,omitempty" xml:"value,omitempty"`
}

func (s ScaleClusterRequestTaints) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequestTaints) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequestTaints) SetEffect(v string) *ScaleClusterRequestTaints {
	s.Effect = &v
	return s
}

func (s *ScaleClusterRequestTaints) SetKey(v string) *ScaleClusterRequestTaints {
	s.Key = &v
	return s
}

func (s *ScaleClusterRequestTaints) SetValue(v string) *ScaleClusterRequestTaints {
	s.Value = &v
	return s
}

type ScaleClusterRequestWorkerDataDisks struct {
	Category  *string `json:"category,omitempty" xml:"category,omitempty"`
	Encrypted *string `json:"encrypted,omitempty" xml:"encrypted,omitempty"`
	Size      *string `json:"size,omitempty" xml:"size,omitempty"`
}

func (s ScaleClusterRequestWorkerDataDisks) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterRequestWorkerDataDisks) GoString() string {
	return s.String()
}

func (s *ScaleClusterRequestWorkerDataDisks) SetCategory(v string) *ScaleClusterRequestWorkerDataDisks {
	s.Category = &v
	return s
}

func (s *ScaleClusterRequestWorkerDataDisks) SetEncrypted(v string) *ScaleClusterRequestWorkerDataDisks {
	s.Encrypted = &v
	return s
}

func (s *ScaleClusterRequestWorkerDataDisks) SetSize(v string) *ScaleClusterRequestWorkerDataDisks {
	s.Size = &v
	return s
}

type ScaleClusterResponseBody struct {
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	TaskId    *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ScaleClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterResponseBody) GoString() string {
	return s.String()
}

func (s *ScaleClusterResponseBody) SetClusterId(v string) *ScaleClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *ScaleClusterResponseBody) SetRequestId(v string) *ScaleClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *ScaleClusterResponseBody) SetTaskId(v string) *ScaleClusterResponseBody {
	s.TaskId = &v
	return s
}

type ScaleClusterResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ScaleClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ScaleClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterResponse) GoString() string {
	return s.String()
}

func (s *ScaleClusterResponse) SetHeaders(v map[string]*string) *ScaleClusterResponse {
	s.Headers = v
	return s
}

func (s *ScaleClusterResponse) SetStatusCode(v int32) *ScaleClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *ScaleClusterResponse) SetBody(v *ScaleClusterResponseBody) *ScaleClusterResponse {
	s.Body = v
	return s
}

type ScaleClusterNodePoolRequest struct {
	Count *int64 `json:"count,omitempty" xml:"count,omitempty"`
}

func (s ScaleClusterNodePoolRequest) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterNodePoolRequest) GoString() string {
	return s.String()
}

func (s *ScaleClusterNodePoolRequest) SetCount(v int64) *ScaleClusterNodePoolRequest {
	s.Count = &v
	return s
}

type ScaleClusterNodePoolResponseBody struct {
	// The ID of the scaling task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ScaleClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *ScaleClusterNodePoolResponseBody) SetTaskId(v string) *ScaleClusterNodePoolResponseBody {
	s.TaskId = &v
	return s
}

type ScaleClusterNodePoolResponse struct {
	Headers    map[string]*string                `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                            `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ScaleClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ScaleClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s ScaleClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *ScaleClusterNodePoolResponse) SetHeaders(v map[string]*string) *ScaleClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *ScaleClusterNodePoolResponse) SetStatusCode(v int32) *ScaleClusterNodePoolResponse {
	s.StatusCode = &v
	return s
}

func (s *ScaleClusterNodePoolResponse) SetBody(v *ScaleClusterNodePoolResponseBody) *ScaleClusterNodePoolResponse {
	s.Body = v
	return s
}

type ScaleOutClusterRequest struct {
	CloudMonitorFlags        *bool                                    `json:"cloud_monitor_flags,omitempty" xml:"cloud_monitor_flags,omitempty"`
	Count                    *int64                                   `json:"count,omitempty" xml:"count,omitempty"`
	CpuPolicy                *string                                  `json:"cpu_policy,omitempty" xml:"cpu_policy,omitempty"`
	ImageId                  *string                                  `json:"image_id,omitempty" xml:"image_id,omitempty"`
	KeyPair                  *string                                  `json:"key_pair,omitempty" xml:"key_pair,omitempty"`
	LoginPassword            *string                                  `json:"login_password,omitempty" xml:"login_password,omitempty"`
	RdsInstances             []*string                                `json:"rds_instances,omitempty" xml:"rds_instances,omitempty" type:"Repeated"`
	Runtime                  *Runtime                                 `json:"runtime,omitempty" xml:"runtime,omitempty"`
	Tags                     []*Tag                                   `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
	Taints                   []*Taint                                 `json:"taints,omitempty" xml:"taints,omitempty" type:"Repeated"`
	UserData                 *string                                  `json:"user_data,omitempty" xml:"user_data,omitempty"`
	VswitchIds               []*string                                `json:"vswitch_ids,omitempty" xml:"vswitch_ids,omitempty" type:"Repeated"`
	WorkerAutoRenew          *bool                                    `json:"worker_auto_renew,omitempty" xml:"worker_auto_renew,omitempty"`
	WorkerAutoRenewPeriod    *int64                                   `json:"worker_auto_renew_period,omitempty" xml:"worker_auto_renew_period,omitempty"`
	WorkerDataDisks          []*ScaleOutClusterRequestWorkerDataDisks `json:"worker_data_disks,omitempty" xml:"worker_data_disks,omitempty" type:"Repeated"`
	WorkerInstanceChargeType *string                                  `json:"worker_instance_charge_type,omitempty" xml:"worker_instance_charge_type,omitempty"`
	WorkerInstanceTypes      []*string                                `json:"worker_instance_types,omitempty" xml:"worker_instance_types,omitempty" type:"Repeated"`
	WorkerPeriod             *int64                                   `json:"worker_period,omitempty" xml:"worker_period,omitempty"`
	WorkerPeriodUnit         *string                                  `json:"worker_period_unit,omitempty" xml:"worker_period_unit,omitempty"`
	WorkerSystemDiskCategory *string                                  `json:"worker_system_disk_category,omitempty" xml:"worker_system_disk_category,omitempty"`
	WorkerSystemDiskSize     *int64                                   `json:"worker_system_disk_size,omitempty" xml:"worker_system_disk_size,omitempty"`
}

func (s ScaleOutClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterRequest) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterRequest) SetCloudMonitorFlags(v bool) *ScaleOutClusterRequest {
	s.CloudMonitorFlags = &v
	return s
}

func (s *ScaleOutClusterRequest) SetCount(v int64) *ScaleOutClusterRequest {
	s.Count = &v
	return s
}

func (s *ScaleOutClusterRequest) SetCpuPolicy(v string) *ScaleOutClusterRequest {
	s.CpuPolicy = &v
	return s
}

func (s *ScaleOutClusterRequest) SetImageId(v string) *ScaleOutClusterRequest {
	s.ImageId = &v
	return s
}

func (s *ScaleOutClusterRequest) SetKeyPair(v string) *ScaleOutClusterRequest {
	s.KeyPair = &v
	return s
}

func (s *ScaleOutClusterRequest) SetLoginPassword(v string) *ScaleOutClusterRequest {
	s.LoginPassword = &v
	return s
}

func (s *ScaleOutClusterRequest) SetRdsInstances(v []*string) *ScaleOutClusterRequest {
	s.RdsInstances = v
	return s
}

func (s *ScaleOutClusterRequest) SetRuntime(v *Runtime) *ScaleOutClusterRequest {
	s.Runtime = v
	return s
}

func (s *ScaleOutClusterRequest) SetTags(v []*Tag) *ScaleOutClusterRequest {
	s.Tags = v
	return s
}

func (s *ScaleOutClusterRequest) SetTaints(v []*Taint) *ScaleOutClusterRequest {
	s.Taints = v
	return s
}

func (s *ScaleOutClusterRequest) SetUserData(v string) *ScaleOutClusterRequest {
	s.UserData = &v
	return s
}

func (s *ScaleOutClusterRequest) SetVswitchIds(v []*string) *ScaleOutClusterRequest {
	s.VswitchIds = v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerAutoRenew(v bool) *ScaleOutClusterRequest {
	s.WorkerAutoRenew = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerAutoRenewPeriod(v int64) *ScaleOutClusterRequest {
	s.WorkerAutoRenewPeriod = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerDataDisks(v []*ScaleOutClusterRequestWorkerDataDisks) *ScaleOutClusterRequest {
	s.WorkerDataDisks = v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerInstanceChargeType(v string) *ScaleOutClusterRequest {
	s.WorkerInstanceChargeType = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerInstanceTypes(v []*string) *ScaleOutClusterRequest {
	s.WorkerInstanceTypes = v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerPeriod(v int64) *ScaleOutClusterRequest {
	s.WorkerPeriod = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerPeriodUnit(v string) *ScaleOutClusterRequest {
	s.WorkerPeriodUnit = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerSystemDiskCategory(v string) *ScaleOutClusterRequest {
	s.WorkerSystemDiskCategory = &v
	return s
}

func (s *ScaleOutClusterRequest) SetWorkerSystemDiskSize(v int64) *ScaleOutClusterRequest {
	s.WorkerSystemDiskSize = &v
	return s
}

type ScaleOutClusterRequestWorkerDataDisks struct {
	AutoSnapshotPolicyId *string `json:"auto_snapshot_policy_id,omitempty" xml:"auto_snapshot_policy_id,omitempty"`
	Category             *string `json:"category,omitempty" xml:"category,omitempty"`
	Encrypted            *string `json:"encrypted,omitempty" xml:"encrypted,omitempty"`
	Size                 *string `json:"size,omitempty" xml:"size,omitempty"`
}

func (s ScaleOutClusterRequestWorkerDataDisks) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterRequestWorkerDataDisks) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetAutoSnapshotPolicyId(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.AutoSnapshotPolicyId = &v
	return s
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetCategory(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.Category = &v
	return s
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetEncrypted(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.Encrypted = &v
	return s
}

func (s *ScaleOutClusterRequestWorkerDataDisks) SetSize(v string) *ScaleOutClusterRequestWorkerDataDisks {
	s.Size = &v
	return s
}

type ScaleOutClusterResponseBody struct {
	// The ID of the cluster.
	ClusterId *string `json:"cluster_id,omitempty" xml:"cluster_id,omitempty"`
	// The ID of the request.
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	// The ID of the task.
	TaskId *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ScaleOutClusterResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterResponseBody) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterResponseBody) SetClusterId(v string) *ScaleOutClusterResponseBody {
	s.ClusterId = &v
	return s
}

func (s *ScaleOutClusterResponseBody) SetRequestId(v string) *ScaleOutClusterResponseBody {
	s.RequestId = &v
	return s
}

func (s *ScaleOutClusterResponseBody) SetTaskId(v string) *ScaleOutClusterResponseBody {
	s.TaskId = &v
	return s
}

type ScaleOutClusterResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ScaleOutClusterResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ScaleOutClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s ScaleOutClusterResponse) GoString() string {
	return s.String()
}

func (s *ScaleOutClusterResponse) SetHeaders(v map[string]*string) *ScaleOutClusterResponse {
	s.Headers = v
	return s
}

func (s *ScaleOutClusterResponse) SetStatusCode(v int32) *ScaleOutClusterResponse {
	s.StatusCode = &v
	return s
}

func (s *ScaleOutClusterResponse) SetBody(v *ScaleOutClusterResponseBody) *ScaleOutClusterResponse {
	s.Body = v
	return s
}

type ScanClusterVulsResponseBody struct {
	RequestId *string `json:"request_id,omitempty" xml:"request_id,omitempty"`
	TaskId    *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s ScanClusterVulsResponseBody) String() string {
	return tea.Prettify(s)
}

func (s ScanClusterVulsResponseBody) GoString() string {
	return s.String()
}

func (s *ScanClusterVulsResponseBody) SetRequestId(v string) *ScanClusterVulsResponseBody {
	s.RequestId = &v
	return s
}

func (s *ScanClusterVulsResponseBody) SetTaskId(v string) *ScanClusterVulsResponseBody {
	s.TaskId = &v
	return s
}

type ScanClusterVulsResponse struct {
	Headers    map[string]*string           `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                       `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *ScanClusterVulsResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s ScanClusterVulsResponse) String() string {
	return tea.Prettify(s)
}

func (s ScanClusterVulsResponse) GoString() string {
	return s.String()
}

func (s *ScanClusterVulsResponse) SetHeaders(v map[string]*string) *ScanClusterVulsResponse {
	s.Headers = v
	return s
}

func (s *ScanClusterVulsResponse) SetStatusCode(v int32) *ScanClusterVulsResponse {
	s.StatusCode = &v
	return s
}

func (s *ScanClusterVulsResponse) SetBody(v *ScanClusterVulsResponseBody) *ScanClusterVulsResponse {
	s.Body = v
	return s
}

type StartAlertResponseBody struct {
	// The message returned.
	Msg *string `json:"msg,omitempty" xml:"msg,omitempty"`
	// The status.
	Status *bool `json:"status,omitempty" xml:"status,omitempty"`
}

func (s StartAlertResponseBody) String() string {
	return tea.Prettify(s)
}

func (s StartAlertResponseBody) GoString() string {
	return s.String()
}

func (s *StartAlertResponseBody) SetMsg(v string) *StartAlertResponseBody {
	s.Msg = &v
	return s
}

func (s *StartAlertResponseBody) SetStatus(v bool) *StartAlertResponseBody {
	s.Status = &v
	return s
}

type StartAlertResponse struct {
	Headers    map[string]*string      `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                  `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *StartAlertResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s StartAlertResponse) String() string {
	return tea.Prettify(s)
}

func (s StartAlertResponse) GoString() string {
	return s.String()
}

func (s *StartAlertResponse) SetHeaders(v map[string]*string) *StartAlertResponse {
	s.Headers = v
	return s
}

func (s *StartAlertResponse) SetStatusCode(v int32) *StartAlertResponse {
	s.StatusCode = &v
	return s
}

func (s *StartAlertResponse) SetBody(v *StartAlertResponseBody) *StartAlertResponse {
	s.Body = v
	return s
}

type StartWorkflowRequest struct {
	MappingBamOutFilename      *string `json:"mapping_bam_out_filename,omitempty" xml:"mapping_bam_out_filename,omitempty"`
	MappingBamOutPath          *string `json:"mapping_bam_out_path,omitempty" xml:"mapping_bam_out_path,omitempty"`
	MappingBucketName          *string `json:"mapping_bucket_name,omitempty" xml:"mapping_bucket_name,omitempty"`
	MappingFastqFirstFilename  *string `json:"mapping_fastq_first_filename,omitempty" xml:"mapping_fastq_first_filename,omitempty"`
	MappingFastqPath           *string `json:"mapping_fastq_path,omitempty" xml:"mapping_fastq_path,omitempty"`
	MappingFastqSecondFilename *string `json:"mapping_fastq_second_filename,omitempty" xml:"mapping_fastq_second_filename,omitempty"`
	MappingIsMarkDup           *string `json:"mapping_is_mark_dup,omitempty" xml:"mapping_is_mark_dup,omitempty"`
	MappingOssRegion           *string `json:"mapping_oss_region,omitempty" xml:"mapping_oss_region,omitempty"`
	MappingReferencePath       *string `json:"mapping_reference_path,omitempty" xml:"mapping_reference_path,omitempty"`
	Service                    *string `json:"service,omitempty" xml:"service,omitempty"`
	WgsBucketName              *string `json:"wgs_bucket_name,omitempty" xml:"wgs_bucket_name,omitempty"`
	WgsFastqFirstFilename      *string `json:"wgs_fastq_first_filename,omitempty" xml:"wgs_fastq_first_filename,omitempty"`
	WgsFastqPath               *string `json:"wgs_fastq_path,omitempty" xml:"wgs_fastq_path,omitempty"`
	WgsFastqSecondFilename     *string `json:"wgs_fastq_second_filename,omitempty" xml:"wgs_fastq_second_filename,omitempty"`
	WgsOssRegion               *string `json:"wgs_oss_region,omitempty" xml:"wgs_oss_region,omitempty"`
	WgsReferencePath           *string `json:"wgs_reference_path,omitempty" xml:"wgs_reference_path,omitempty"`
	WgsVcfOutFilename          *string `json:"wgs_vcf_out_filename,omitempty" xml:"wgs_vcf_out_filename,omitempty"`
	WgsVcfOutPath              *string `json:"wgs_vcf_out_path,omitempty" xml:"wgs_vcf_out_path,omitempty"`
	WorkflowType               *string `json:"workflow_type,omitempty" xml:"workflow_type,omitempty"`
}

func (s StartWorkflowRequest) String() string {
	return tea.Prettify(s)
}

func (s StartWorkflowRequest) GoString() string {
	return s.String()
}

func (s *StartWorkflowRequest) SetMappingBamOutFilename(v string) *StartWorkflowRequest {
	s.MappingBamOutFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingBamOutPath(v string) *StartWorkflowRequest {
	s.MappingBamOutPath = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingBucketName(v string) *StartWorkflowRequest {
	s.MappingBucketName = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingFastqFirstFilename(v string) *StartWorkflowRequest {
	s.MappingFastqFirstFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingFastqPath(v string) *StartWorkflowRequest {
	s.MappingFastqPath = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingFastqSecondFilename(v string) *StartWorkflowRequest {
	s.MappingFastqSecondFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingIsMarkDup(v string) *StartWorkflowRequest {
	s.MappingIsMarkDup = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingOssRegion(v string) *StartWorkflowRequest {
	s.MappingOssRegion = &v
	return s
}

func (s *StartWorkflowRequest) SetMappingReferencePath(v string) *StartWorkflowRequest {
	s.MappingReferencePath = &v
	return s
}

func (s *StartWorkflowRequest) SetService(v string) *StartWorkflowRequest {
	s.Service = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsBucketName(v string) *StartWorkflowRequest {
	s.WgsBucketName = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsFastqFirstFilename(v string) *StartWorkflowRequest {
	s.WgsFastqFirstFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsFastqPath(v string) *StartWorkflowRequest {
	s.WgsFastqPath = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsFastqSecondFilename(v string) *StartWorkflowRequest {
	s.WgsFastqSecondFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsOssRegion(v string) *StartWorkflowRequest {
	s.WgsOssRegion = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsReferencePath(v string) *StartWorkflowRequest {
	s.WgsReferencePath = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsVcfOutFilename(v string) *StartWorkflowRequest {
	s.WgsVcfOutFilename = &v
	return s
}

func (s *StartWorkflowRequest) SetWgsVcfOutPath(v string) *StartWorkflowRequest {
	s.WgsVcfOutPath = &v
	return s
}

func (s *StartWorkflowRequest) SetWorkflowType(v string) *StartWorkflowRequest {
	s.WorkflowType = &v
	return s
}

type StartWorkflowResponseBody struct {
	// The name of the workflow that is created.
	JobName *string `json:"JobName,omitempty" xml:"JobName,omitempty"`
}

func (s StartWorkflowResponseBody) String() string {
	return tea.Prettify(s)
}

func (s StartWorkflowResponseBody) GoString() string {
	return s.String()
}

func (s *StartWorkflowResponseBody) SetJobName(v string) *StartWorkflowResponseBody {
	s.JobName = &v
	return s
}

type StartWorkflowResponse struct {
	Headers    map[string]*string         `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                     `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *StartWorkflowResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s StartWorkflowResponse) String() string {
	return tea.Prettify(s)
}

func (s StartWorkflowResponse) GoString() string {
	return s.String()
}

func (s *StartWorkflowResponse) SetHeaders(v map[string]*string) *StartWorkflowResponse {
	s.Headers = v
	return s
}

func (s *StartWorkflowResponse) SetStatusCode(v int32) *StartWorkflowResponse {
	s.StatusCode = &v
	return s
}

func (s *StartWorkflowResponse) SetBody(v *StartWorkflowResponseBody) *StartWorkflowResponse {
	s.Body = v
	return s
}

type StopAlertResponseBody struct {
	// The error message returned if the call fails.
	Msg *string `json:"msg,omitempty" xml:"msg,omitempty"`
	// A value of True indicates that the call succeeds. A value of False indicates that the call failed.
	Status *bool `json:"status,omitempty" xml:"status,omitempty"`
}

func (s StopAlertResponseBody) String() string {
	return tea.Prettify(s)
}

func (s StopAlertResponseBody) GoString() string {
	return s.String()
}

func (s *StopAlertResponseBody) SetMsg(v string) *StopAlertResponseBody {
	s.Msg = &v
	return s
}

func (s *StopAlertResponseBody) SetStatus(v bool) *StopAlertResponseBody {
	s.Status = &v
	return s
}

type StopAlertResponse struct {
	Headers    map[string]*string     `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                 `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *StopAlertResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s StopAlertResponse) String() string {
	return tea.Prettify(s)
}

func (s StopAlertResponse) GoString() string {
	return s.String()
}

func (s *StopAlertResponse) SetHeaders(v map[string]*string) *StopAlertResponse {
	s.Headers = v
	return s
}

func (s *StopAlertResponse) SetStatusCode(v int32) *StopAlertResponse {
	s.StatusCode = &v
	return s
}

func (s *StopAlertResponse) SetBody(v *StopAlertResponseBody) *StopAlertResponse {
	s.Body = v
	return s
}

type SyncClusterNodePoolResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s SyncClusterNodePoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s SyncClusterNodePoolResponseBody) GoString() string {
	return s.String()
}

func (s *SyncClusterNodePoolResponseBody) SetRequestId(v string) *SyncClusterNodePoolResponseBody {
	s.RequestId = &v
	return s
}

type SyncClusterNodePoolResponse struct {
	Headers    map[string]*string               `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                           `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *SyncClusterNodePoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s SyncClusterNodePoolResponse) String() string {
	return tea.Prettify(s)
}

func (s SyncClusterNodePoolResponse) GoString() string {
	return s.String()
}

func (s *SyncClusterNodePoolResponse) SetHeaders(v map[string]*string) *SyncClusterNodePoolResponse {
	s.Headers = v
	return s
}

func (s *SyncClusterNodePoolResponse) SetStatusCode(v int32) *SyncClusterNodePoolResponse {
	s.StatusCode = &v
	return s
}

func (s *SyncClusterNodePoolResponse) SetBody(v *SyncClusterNodePoolResponseBody) *SyncClusterNodePoolResponse {
	s.Body = v
	return s
}

type TagResourcesRequest struct {
	RegionId     *string   `json:"region_id,omitempty" xml:"region_id,omitempty"`
	ResourceIds  []*string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty" type:"Repeated"`
	ResourceType *string   `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	Tags         []*Tag    `json:"tags,omitempty" xml:"tags,omitempty" type:"Repeated"`
}

func (s TagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesRequest) GoString() string {
	return s.String()
}

func (s *TagResourcesRequest) SetRegionId(v string) *TagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *TagResourcesRequest) SetResourceIds(v []*string) *TagResourcesRequest {
	s.ResourceIds = v
	return s
}

func (s *TagResourcesRequest) SetResourceType(v string) *TagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *TagResourcesRequest) SetTags(v []*Tag) *TagResourcesRequest {
	s.Tags = v
	return s
}

type TagResourcesResponseBody struct {
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s TagResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *TagResourcesResponseBody) SetRequestId(v string) *TagResourcesResponseBody {
	s.RequestId = &v
	return s
}

type TagResourcesResponse struct {
	Headers    map[string]*string        `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                    `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *TagResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s TagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s TagResourcesResponse) GoString() string {
	return s.String()
}

func (s *TagResourcesResponse) SetHeaders(v map[string]*string) *TagResourcesResponse {
	s.Headers = v
	return s
}

func (s *TagResourcesResponse) SetStatusCode(v int32) *TagResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *TagResourcesResponse) SetBody(v *TagResourcesResponseBody) *TagResourcesResponse {
	s.Body = v
	return s
}

type UnInstallClusterAddonsRequest struct {
	Addons []*UnInstallClusterAddonsRequestAddons `json:"addons,omitempty" xml:"addons,omitempty" type:"Repeated"`
}

func (s UnInstallClusterAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s UnInstallClusterAddonsRequest) GoString() string {
	return s.String()
}

func (s *UnInstallClusterAddonsRequest) SetAddons(v []*UnInstallClusterAddonsRequestAddons) *UnInstallClusterAddonsRequest {
	s.Addons = v
	return s
}

type UnInstallClusterAddonsRequestAddons struct {
	Name *string `json:"name,omitempty" xml:"name,omitempty"`
}

func (s UnInstallClusterAddonsRequestAddons) String() string {
	return tea.Prettify(s)
}

func (s UnInstallClusterAddonsRequestAddons) GoString() string {
	return s.String()
}

func (s *UnInstallClusterAddonsRequestAddons) SetName(v string) *UnInstallClusterAddonsRequestAddons {
	s.Name = &v
	return s
}

type UnInstallClusterAddonsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UnInstallClusterAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s UnInstallClusterAddonsResponse) GoString() string {
	return s.String()
}

func (s *UnInstallClusterAddonsResponse) SetHeaders(v map[string]*string) *UnInstallClusterAddonsResponse {
	s.Headers = v
	return s
}

func (s *UnInstallClusterAddonsResponse) SetStatusCode(v int32) *UnInstallClusterAddonsResponse {
	s.StatusCode = &v
	return s
}

type UntagResourcesRequest struct {
	All          *bool     `json:"all,omitempty" xml:"all,omitempty"`
	RegionId     *string   `json:"region_id,omitempty" xml:"region_id,omitempty"`
	ResourceIds  []*string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty" type:"Repeated"`
	ResourceType *string   `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	TagKeys      []*string `json:"tag_keys,omitempty" xml:"tag_keys,omitempty" type:"Repeated"`
}

func (s UntagResourcesRequest) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesRequest) GoString() string {
	return s.String()
}

func (s *UntagResourcesRequest) SetAll(v bool) *UntagResourcesRequest {
	s.All = &v
	return s
}

func (s *UntagResourcesRequest) SetRegionId(v string) *UntagResourcesRequest {
	s.RegionId = &v
	return s
}

func (s *UntagResourcesRequest) SetResourceIds(v []*string) *UntagResourcesRequest {
	s.ResourceIds = v
	return s
}

func (s *UntagResourcesRequest) SetResourceType(v string) *UntagResourcesRequest {
	s.ResourceType = &v
	return s
}

func (s *UntagResourcesRequest) SetTagKeys(v []*string) *UntagResourcesRequest {
	s.TagKeys = v
	return s
}

type UntagResourcesShrinkRequest struct {
	All               *bool   `json:"all,omitempty" xml:"all,omitempty"`
	RegionId          *string `json:"region_id,omitempty" xml:"region_id,omitempty"`
	ResourceIdsShrink *string `json:"resource_ids,omitempty" xml:"resource_ids,omitempty"`
	ResourceType      *string `json:"resource_type,omitempty" xml:"resource_type,omitempty"`
	TagKeysShrink     *string `json:"tag_keys,omitempty" xml:"tag_keys,omitempty"`
}

func (s UntagResourcesShrinkRequest) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesShrinkRequest) GoString() string {
	return s.String()
}

func (s *UntagResourcesShrinkRequest) SetAll(v bool) *UntagResourcesShrinkRequest {
	s.All = &v
	return s
}

func (s *UntagResourcesShrinkRequest) SetRegionId(v string) *UntagResourcesShrinkRequest {
	s.RegionId = &v
	return s
}

func (s *UntagResourcesShrinkRequest) SetResourceIdsShrink(v string) *UntagResourcesShrinkRequest {
	s.ResourceIdsShrink = &v
	return s
}

func (s *UntagResourcesShrinkRequest) SetResourceType(v string) *UntagResourcesShrinkRequest {
	s.ResourceType = &v
	return s
}

func (s *UntagResourcesShrinkRequest) SetTagKeysShrink(v string) *UntagResourcesShrinkRequest {
	s.TagKeysShrink = &v
	return s
}

type UntagResourcesResponseBody struct {
	// The ID of the request.
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
}

func (s UntagResourcesResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesResponseBody) GoString() string {
	return s.String()
}

func (s *UntagResourcesResponseBody) SetRequestId(v string) *UntagResourcesResponseBody {
	s.RequestId = &v
	return s
}

type UntagResourcesResponse struct {
	Headers    map[string]*string          `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                      `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UntagResourcesResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UntagResourcesResponse) String() string {
	return tea.Prettify(s)
}

func (s UntagResourcesResponse) GoString() string {
	return s.String()
}

func (s *UntagResourcesResponse) SetHeaders(v map[string]*string) *UntagResourcesResponse {
	s.Headers = v
	return s
}

func (s *UntagResourcesResponse) SetStatusCode(v int32) *UntagResourcesResponse {
	s.StatusCode = &v
	return s
}

func (s *UntagResourcesResponse) SetBody(v *UntagResourcesResponseBody) *UntagResourcesResponse {
	s.Body = v
	return s
}

type UpdateContactGroupForAlertResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UpdateContactGroupForAlertResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateContactGroupForAlertResponse) GoString() string {
	return s.String()
}

func (s *UpdateContactGroupForAlertResponse) SetHeaders(v map[string]*string) *UpdateContactGroupForAlertResponse {
	s.Headers = v
	return s
}

func (s *UpdateContactGroupForAlertResponse) SetStatusCode(v int32) *UpdateContactGroupForAlertResponse {
	s.StatusCode = &v
	return s
}

type UpdateControlPlaneLogRequest struct {
	Aliuid     *string   `json:"aliuid,omitempty" xml:"aliuid,omitempty"`
	Components []*string `json:"components,omitempty" xml:"components,omitempty" type:"Repeated"`
	LogProject *string   `json:"log_project,omitempty" xml:"log_project,omitempty"`
	LogTtl     *string   `json:"log_ttl,omitempty" xml:"log_ttl,omitempty"`
}

func (s UpdateControlPlaneLogRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateControlPlaneLogRequest) GoString() string {
	return s.String()
}

func (s *UpdateControlPlaneLogRequest) SetAliuid(v string) *UpdateControlPlaneLogRequest {
	s.Aliuid = &v
	return s
}

func (s *UpdateControlPlaneLogRequest) SetComponents(v []*string) *UpdateControlPlaneLogRequest {
	s.Components = v
	return s
}

func (s *UpdateControlPlaneLogRequest) SetLogProject(v string) *UpdateControlPlaneLogRequest {
	s.LogProject = &v
	return s
}

func (s *UpdateControlPlaneLogRequest) SetLogTtl(v string) *UpdateControlPlaneLogRequest {
	s.LogTtl = &v
	return s
}

type UpdateControlPlaneLogResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UpdateControlPlaneLogResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateControlPlaneLogResponse) GoString() string {
	return s.String()
}

func (s *UpdateControlPlaneLogResponse) SetHeaders(v map[string]*string) *UpdateControlPlaneLogResponse {
	s.Headers = v
	return s
}

func (s *UpdateControlPlaneLogResponse) SetStatusCode(v int32) *UpdateControlPlaneLogResponse {
	s.StatusCode = &v
	return s
}

type UpdateK8sClusterUserConfigExpireRequest struct {
	ExpireHour *int64  `json:"expire_hour,omitempty" xml:"expire_hour,omitempty"`
	User       *string `json:"user,omitempty" xml:"user,omitempty"`
}

func (s UpdateK8sClusterUserConfigExpireRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateK8sClusterUserConfigExpireRequest) GoString() string {
	return s.String()
}

func (s *UpdateK8sClusterUserConfigExpireRequest) SetExpireHour(v int64) *UpdateK8sClusterUserConfigExpireRequest {
	s.ExpireHour = &v
	return s
}

func (s *UpdateK8sClusterUserConfigExpireRequest) SetUser(v string) *UpdateK8sClusterUserConfigExpireRequest {
	s.User = &v
	return s
}

type UpdateK8sClusterUserConfigExpireResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UpdateK8sClusterUserConfigExpireResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateK8sClusterUserConfigExpireResponse) GoString() string {
	return s.String()
}

func (s *UpdateK8sClusterUserConfigExpireResponse) SetHeaders(v map[string]*string) *UpdateK8sClusterUserConfigExpireResponse {
	s.Headers = v
	return s
}

func (s *UpdateK8sClusterUserConfigExpireResponse) SetStatusCode(v int32) *UpdateK8sClusterUserConfigExpireResponse {
	s.StatusCode = &v
	return s
}

type UpdateTemplateRequest struct {
	Description  *string `json:"description,omitempty" xml:"description,omitempty"`
	Name         *string `json:"name,omitempty" xml:"name,omitempty"`
	Tags         *string `json:"tags,omitempty" xml:"tags,omitempty"`
	Template     *string `json:"template,omitempty" xml:"template,omitempty"`
	TemplateType *string `json:"template_type,omitempty" xml:"template_type,omitempty"`
}

func (s UpdateTemplateRequest) String() string {
	return tea.Prettify(s)
}

func (s UpdateTemplateRequest) GoString() string {
	return s.String()
}

func (s *UpdateTemplateRequest) SetDescription(v string) *UpdateTemplateRequest {
	s.Description = &v
	return s
}

func (s *UpdateTemplateRequest) SetName(v string) *UpdateTemplateRequest {
	s.Name = &v
	return s
}

func (s *UpdateTemplateRequest) SetTags(v string) *UpdateTemplateRequest {
	s.Tags = &v
	return s
}

func (s *UpdateTemplateRequest) SetTemplate(v string) *UpdateTemplateRequest {
	s.Template = &v
	return s
}

func (s *UpdateTemplateRequest) SetTemplateType(v string) *UpdateTemplateRequest {
	s.TemplateType = &v
	return s
}

type UpdateTemplateResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UpdateTemplateResponse) String() string {
	return tea.Prettify(s)
}

func (s UpdateTemplateResponse) GoString() string {
	return s.String()
}

func (s *UpdateTemplateResponse) SetHeaders(v map[string]*string) *UpdateTemplateResponse {
	s.Headers = v
	return s
}

func (s *UpdateTemplateResponse) SetStatusCode(v int32) *UpdateTemplateResponse {
	s.StatusCode = &v
	return s
}

type UpgradeClusterRequest struct {
	// Deprecated
	ComponentName *string `json:"component_name,omitempty" xml:"component_name,omitempty"`
	MasterOnly    *bool   `json:"master_only,omitempty" xml:"master_only,omitempty"`
	NextVersion   *string `json:"next_version,omitempty" xml:"next_version,omitempty"`
	// Deprecated
	Version *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s UpgradeClusterRequest) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterRequest) GoString() string {
	return s.String()
}

func (s *UpgradeClusterRequest) SetComponentName(v string) *UpgradeClusterRequest {
	s.ComponentName = &v
	return s
}

func (s *UpgradeClusterRequest) SetMasterOnly(v bool) *UpgradeClusterRequest {
	s.MasterOnly = &v
	return s
}

func (s *UpgradeClusterRequest) SetNextVersion(v string) *UpgradeClusterRequest {
	s.NextVersion = &v
	return s
}

func (s *UpgradeClusterRequest) SetVersion(v string) *UpgradeClusterRequest {
	s.Version = &v
	return s
}

type UpgradeClusterResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UpgradeClusterResponse) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterResponse) GoString() string {
	return s.String()
}

func (s *UpgradeClusterResponse) SetHeaders(v map[string]*string) *UpgradeClusterResponse {
	s.Headers = v
	return s
}

func (s *UpgradeClusterResponse) SetStatusCode(v int32) *UpgradeClusterResponse {
	s.StatusCode = &v
	return s
}

type UpgradeClusterAddonsRequest struct {
	Body []*UpgradeClusterAddonsRequestBody `json:"body,omitempty" xml:"body,omitempty" type:"Repeated"`
}

func (s UpgradeClusterAddonsRequest) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterAddonsRequest) GoString() string {
	return s.String()
}

func (s *UpgradeClusterAddonsRequest) SetBody(v []*UpgradeClusterAddonsRequestBody) *UpgradeClusterAddonsRequest {
	s.Body = v
	return s
}

type UpgradeClusterAddonsRequestBody struct {
	ComponentName *string `json:"component_name,omitempty" xml:"component_name,omitempty"`
	Config        *string `json:"config,omitempty" xml:"config,omitempty"`
	NextVersion   *string `json:"next_version,omitempty" xml:"next_version,omitempty"`
	Policy        *string `json:"policy,omitempty" xml:"policy,omitempty"`
	Version       *string `json:"version,omitempty" xml:"version,omitempty"`
}

func (s UpgradeClusterAddonsRequestBody) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterAddonsRequestBody) GoString() string {
	return s.String()
}

func (s *UpgradeClusterAddonsRequestBody) SetComponentName(v string) *UpgradeClusterAddonsRequestBody {
	s.ComponentName = &v
	return s
}

func (s *UpgradeClusterAddonsRequestBody) SetConfig(v string) *UpgradeClusterAddonsRequestBody {
	s.Config = &v
	return s
}

func (s *UpgradeClusterAddonsRequestBody) SetNextVersion(v string) *UpgradeClusterAddonsRequestBody {
	s.NextVersion = &v
	return s
}

func (s *UpgradeClusterAddonsRequestBody) SetPolicy(v string) *UpgradeClusterAddonsRequestBody {
	s.Policy = &v
	return s
}

func (s *UpgradeClusterAddonsRequestBody) SetVersion(v string) *UpgradeClusterAddonsRequestBody {
	s.Version = &v
	return s
}

type UpgradeClusterAddonsResponse struct {
	Headers    map[string]*string `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32             `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
}

func (s UpgradeClusterAddonsResponse) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterAddonsResponse) GoString() string {
	return s.String()
}

func (s *UpgradeClusterAddonsResponse) SetHeaders(v map[string]*string) *UpgradeClusterAddonsResponse {
	s.Headers = v
	return s
}

func (s *UpgradeClusterAddonsResponse) SetStatusCode(v int32) *UpgradeClusterAddonsResponse {
	s.StatusCode = &v
	return s
}

type UpgradeClusterNodepoolRequest struct {
	ImageId           *string `json:"image_id,omitempty" xml:"image_id,omitempty"`
	KubernetesVersion *string `json:"kubernetes_version,omitempty" xml:"kubernetes_version,omitempty"`
	RuntimeType       *string `json:"runtime_type,omitempty" xml:"runtime_type,omitempty"`
	RuntimeVersion    *string `json:"runtime_version,omitempty" xml:"runtime_version,omitempty"`
}

func (s UpgradeClusterNodepoolRequest) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterNodepoolRequest) GoString() string {
	return s.String()
}

func (s *UpgradeClusterNodepoolRequest) SetImageId(v string) *UpgradeClusterNodepoolRequest {
	s.ImageId = &v
	return s
}

func (s *UpgradeClusterNodepoolRequest) SetKubernetesVersion(v string) *UpgradeClusterNodepoolRequest {
	s.KubernetesVersion = &v
	return s
}

func (s *UpgradeClusterNodepoolRequest) SetRuntimeType(v string) *UpgradeClusterNodepoolRequest {
	s.RuntimeType = &v
	return s
}

func (s *UpgradeClusterNodepoolRequest) SetRuntimeVersion(v string) *UpgradeClusterNodepoolRequest {
	s.RuntimeVersion = &v
	return s
}

type UpgradeClusterNodepoolResponseBody struct {
	// The ID of the request.
	RequestId *string `json:"RequestId,omitempty" xml:"RequestId,omitempty"`
	TaskId    *string `json:"task_id,omitempty" xml:"task_id,omitempty"`
}

func (s UpgradeClusterNodepoolResponseBody) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterNodepoolResponseBody) GoString() string {
	return s.String()
}

func (s *UpgradeClusterNodepoolResponseBody) SetRequestId(v string) *UpgradeClusterNodepoolResponseBody {
	s.RequestId = &v
	return s
}

func (s *UpgradeClusterNodepoolResponseBody) SetTaskId(v string) *UpgradeClusterNodepoolResponseBody {
	s.TaskId = &v
	return s
}

type UpgradeClusterNodepoolResponse struct {
	Headers    map[string]*string                  `json:"headers,omitempty" xml:"headers,omitempty" require:"true"`
	StatusCode *int32                              `json:"statusCode,omitempty" xml:"statusCode,omitempty" require:"true"`
	Body       *UpgradeClusterNodepoolResponseBody `json:"body,omitempty" xml:"body,omitempty" require:"true"`
}

func (s UpgradeClusterNodepoolResponse) String() string {
	return tea.Prettify(s)
}

func (s UpgradeClusterNodepoolResponse) GoString() string {
	return s.String()
}

func (s *UpgradeClusterNodepoolResponse) SetHeaders(v map[string]*string) *UpgradeClusterNodepoolResponse {
	s.Headers = v
	return s
}

func (s *UpgradeClusterNodepoolResponse) SetStatusCode(v int32) *UpgradeClusterNodepoolResponse {
	s.StatusCode = &v
	return s
}

func (s *UpgradeClusterNodepoolResponse) SetBody(v *UpgradeClusterNodepoolResponseBody) *UpgradeClusterNodepoolResponse {
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
	client.SignatureAlgorithm = tea.String("v2")
	client.EndpointRule = tea.String("regional")
	client.EndpointMap = map[string]*string{
		"ap-northeast-2-pop":          tea.String("cs.aliyuncs.com"),
		"cn-beijing-finance-pop":      tea.String("cs.aliyuncs.com"),
		"cn-beijing-gov-1":            tea.String("cs.aliyuncs.com"),
		"cn-beijing-nu16-b01":         tea.String("cs.aliyuncs.com"),
		"cn-edge-1":                   tea.String("cs.aliyuncs.com"),
		"cn-fujian":                   tea.String("cs.aliyuncs.com"),
		"cn-haidian-cm12-c01":         tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-bj-b01":          tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-internal-prod-1": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-internal-test-1": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-internal-test-2": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-internal-test-3": tea.String("cs.aliyuncs.com"),
		"cn-hangzhou-test-306":        tea.String("cs.aliyuncs.com"),
		"cn-hongkong-finance-pop":     tea.String("cs.aliyuncs.com"),
		"cn-qingdao-nebula":           tea.String("cs.aliyuncs.com"),
		"cn-shanghai-et15-b01":        tea.String("cs.aliyuncs.com"),
		"cn-shanghai-et2-b01":         tea.String("cs.aliyuncs.com"),
		"cn-shanghai-inner":           tea.String("cs.aliyuncs.com"),
		"cn-shanghai-internal-test-1": tea.String("cs.aliyuncs.com"),
		"cn-shenzhen-inner":           tea.String("cs.aliyuncs.com"),
		"cn-shenzhen-st4-d01":         tea.String("cs.aliyuncs.com"),
		"cn-shenzhen-su18-b01":        tea.String("cs.aliyuncs.com"),
		"cn-wuhan":                    tea.String("cs.aliyuncs.com"),
		"cn-yushanfang":               tea.String("cs.aliyuncs.com"),
		"cn-zhangbei":                 tea.String("cs.aliyuncs.com"),
		"cn-zhangbei-na61-b01":        tea.String("cs.aliyuncs.com"),
		"cn-zhangjiakou-na62-a01":     tea.String("cs.aliyuncs.com"),
		"cn-zhengzhou-nebula-1":       tea.String("cs.aliyuncs.com"),
		"eu-west-1-oxs":               tea.String("cs.aliyuncs.com"),
		"rus-west-1-pop":              tea.String("cs.aliyuncs.com"),
	}
	_err = client.CheckConfig(config)
	if _err != nil {
		return _err
	}
	client.Endpoint, _err = client.GetEndpoint(tea.String("cs"), client.RegionId, client.EndpointRule, client.Network, client.Suffix, client.EndpointMap, client.Endpoint)
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

func (client *Client) AttachInstancesWithOptions(ClusterId *string, request *AttachInstancesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *AttachInstancesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.FormatDisk)) {
		body["format_disk"] = request.FormatDisk
	}

	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["image_id"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.Instances)) {
		body["instances"] = request.Instances
	}

	if !tea.BoolValue(util.IsUnset(request.IsEdgeWorker)) {
		body["is_edge_worker"] = request.IsEdgeWorker
	}

	if !tea.BoolValue(util.IsUnset(request.KeepInstanceName)) {
		body["keep_instance_name"] = request.KeepInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolId)) {
		body["nodepool_id"] = request.NodepoolId
	}

	if !tea.BoolValue(util.IsUnset(request.Password)) {
		body["password"] = request.Password
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	if !tea.BoolValue(util.IsUnset(request.Runtime)) {
		body["runtime"] = request.Runtime
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.UserData)) {
		body["user_data"] = request.UserData
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("AttachInstances"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/attach"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &AttachInstancesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) AttachInstances(ClusterId *string, request *AttachInstancesRequest) (_result *AttachInstancesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &AttachInstancesResponse{}
	_body, _err := client.AttachInstancesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) AttachInstancesToNodePoolWithOptions(ClusterId *string, NodepoolId *string, request *AttachInstancesToNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *AttachInstancesToNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.FormatDisk)) {
		body["format_disk"] = request.FormatDisk
	}

	if !tea.BoolValue(util.IsUnset(request.Instances)) {
		body["instances"] = request.Instances
	}

	if !tea.BoolValue(util.IsUnset(request.KeepInstanceName)) {
		body["keep_instance_name"] = request.KeepInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.Password)) {
		body["password"] = request.Password
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("AttachInstancesToNodePool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId)) + "/attach"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &AttachInstancesToNodePoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) AttachInstancesToNodePool(ClusterId *string, NodepoolId *string, request *AttachInstancesToNodePoolRequest) (_result *AttachInstancesToNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &AttachInstancesToNodePoolResponse{}
	_body, _err := client.AttachInstancesToNodePoolWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelClusterUpgradeWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CancelClusterUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("CancelClusterUpgrade"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/upgrade/cancel"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &CancelClusterUpgradeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelClusterUpgrade(ClusterId *string) (_result *CancelClusterUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CancelClusterUpgradeResponse{}
	_body, _err := client.CancelClusterUpgradeWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelComponentUpgradeWithOptions(clusterId *string, componentId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CancelComponentUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("CancelComponentUpgrade"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/components/" + tea.StringValue(openapiutil.GetEncodeParam(componentId)) + "/cancel"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &CancelComponentUpgradeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelComponentUpgrade(clusterId *string, componentId *string) (_result *CancelComponentUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CancelComponentUpgradeResponse{}
	_body, _err := client.CancelComponentUpgradeWithOptions(clusterId, componentId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelTaskWithOptions(taskId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CancelTaskResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("CancelTask"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/tasks/" + tea.StringValue(openapiutil.GetEncodeParam(taskId)) + "/cancel"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &CancelTaskResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelTask(taskId *string) (_result *CancelTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CancelTaskResponse{}
	_body, _err := client.CancelTaskWithOptions(taskId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CancelWorkflowWithOptions(workflowName *string, request *CancelWorkflowRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CancelWorkflowResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Action)) {
		body["action"] = request.Action
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CancelWorkflow"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/gs/workflow/" + tea.StringValue(openapiutil.GetEncodeParam(workflowName))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &CancelWorkflowResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CancelWorkflow(workflowName *string, request *CancelWorkflowRequest) (_result *CancelWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CancelWorkflowResponse{}
	_body, _err := client.CancelWorkflowWithOptions(workflowName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CheckControlPlaneLogEnableWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CheckControlPlaneLogEnableResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("CheckControlPlaneLogEnable"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/controlplanelog"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CheckControlPlaneLogEnableResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CheckControlPlaneLogEnable(ClusterId *string) (_result *CheckControlPlaneLogEnableResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CheckControlPlaneLogEnableResponse{}
	_body, _err := client.CheckControlPlaneLogEnableWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateAutoscalingConfigWithOptions(ClusterId *string, request *CreateAutoscalingConfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateAutoscalingConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CoolDownDuration)) {
		body["cool_down_duration"] = request.CoolDownDuration
	}

	if !tea.BoolValue(util.IsUnset(request.DaemonsetEvictionForNodes)) {
		body["daemonset_eviction_for_nodes"] = request.DaemonsetEvictionForNodes
	}

	if !tea.BoolValue(util.IsUnset(request.Expander)) {
		body["expander"] = request.Expander
	}

	if !tea.BoolValue(util.IsUnset(request.GpuUtilizationThreshold)) {
		body["gpu_utilization_threshold"] = request.GpuUtilizationThreshold
	}

	if !tea.BoolValue(util.IsUnset(request.MaxGracefulTerminationSec)) {
		body["max_graceful_termination_sec"] = request.MaxGracefulTerminationSec
	}

	if !tea.BoolValue(util.IsUnset(request.MinReplicaCount)) {
		body["min_replica_count"] = request.MinReplicaCount
	}

	if !tea.BoolValue(util.IsUnset(request.RecycleNodeDeletionEnabled)) {
		body["recycle_node_deletion_enabled"] = request.RecycleNodeDeletionEnabled
	}

	if !tea.BoolValue(util.IsUnset(request.ScaleDownEnabled)) {
		body["scale_down_enabled"] = request.ScaleDownEnabled
	}

	if !tea.BoolValue(util.IsUnset(request.ScaleUpFromZero)) {
		body["scale_up_from_zero"] = request.ScaleUpFromZero
	}

	if !tea.BoolValue(util.IsUnset(request.ScanInterval)) {
		body["scan_interval"] = request.ScanInterval
	}

	if !tea.BoolValue(util.IsUnset(request.SkipNodesWithLocalStorage)) {
		body["skip_nodes_with_local_storage"] = request.SkipNodesWithLocalStorage
	}

	if !tea.BoolValue(util.IsUnset(request.SkipNodesWithSystemPods)) {
		body["skip_nodes_with_system_pods"] = request.SkipNodesWithSystemPods
	}

	if !tea.BoolValue(util.IsUnset(request.UnneededDuration)) {
		body["unneeded_duration"] = request.UnneededDuration
	}

	if !tea.BoolValue(util.IsUnset(request.UtilizationThreshold)) {
		body["utilization_threshold"] = request.UtilizationThreshold
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateAutoscalingConfig"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/cluster/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/autoscale/config/"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &CreateAutoscalingConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateAutoscalingConfig(ClusterId *string, request *CreateAutoscalingConfigRequest) (_result *CreateAutoscalingConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateAutoscalingConfigResponse{}
	_body, _err := client.CreateAutoscalingConfigWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateClusterWithOptions(request *CreateClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Addons)) {
		body["addons"] = request.Addons
	}

	if !tea.BoolValue(util.IsUnset(request.ApiAudiences)) {
		body["api_audiences"] = request.ApiAudiences
	}

	if !tea.BoolValue(util.IsUnset(request.ChargeType)) {
		body["charge_type"] = request.ChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.CisEnabled)) {
		body["cis_enabled"] = request.CisEnabled
	}

	if !tea.BoolValue(util.IsUnset(request.CloudMonitorFlags)) {
		body["cloud_monitor_flags"] = request.CloudMonitorFlags
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterDomain)) {
		body["cluster_domain"] = request.ClusterDomain
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterSpec)) {
		body["cluster_spec"] = request.ClusterSpec
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		body["cluster_type"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.ContainerCidr)) {
		body["container_cidr"] = request.ContainerCidr
	}

	if !tea.BoolValue(util.IsUnset(request.ControlplaneLogComponents)) {
		body["controlplane_log_components"] = request.ControlplaneLogComponents
	}

	if !tea.BoolValue(util.IsUnset(request.ControlplaneLogProject)) {
		body["controlplane_log_project"] = request.ControlplaneLogProject
	}

	if !tea.BoolValue(util.IsUnset(request.ControlplaneLogTtl)) {
		body["controlplane_log_ttl"] = request.ControlplaneLogTtl
	}

	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.CustomSan)) {
		body["custom_san"] = request.CustomSan
	}

	if !tea.BoolValue(util.IsUnset(request.DeletionProtection)) {
		body["deletion_protection"] = request.DeletionProtection
	}

	if !tea.BoolValue(util.IsUnset(request.DisableRollback)) {
		body["disable_rollback"] = request.DisableRollback
	}

	if !tea.BoolValue(util.IsUnset(request.EnableRrsa)) {
		body["enable_rrsa"] = request.EnableRrsa
	}

	if !tea.BoolValue(util.IsUnset(request.EncryptionProviderKey)) {
		body["encryption_provider_key"] = request.EncryptionProviderKey
	}

	if !tea.BoolValue(util.IsUnset(request.EndpointPublicAccess)) {
		body["endpoint_public_access"] = request.EndpointPublicAccess
	}

	if !tea.BoolValue(util.IsUnset(request.FormatDisk)) {
		body["format_disk"] = request.FormatDisk
	}

	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["image_id"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.ImageType)) {
		body["image_type"] = request.ImageType
	}

	if !tea.BoolValue(util.IsUnset(request.Instances)) {
		body["instances"] = request.Instances
	}

	if !tea.BoolValue(util.IsUnset(request.IpStack)) {
		body["ip_stack"] = request.IpStack
	}

	if !tea.BoolValue(util.IsUnset(request.IsEnterpriseSecurityGroup)) {
		body["is_enterprise_security_group"] = request.IsEnterpriseSecurityGroup
	}

	if !tea.BoolValue(util.IsUnset(request.KeepInstanceName)) {
		body["keep_instance_name"] = request.KeepInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.KubernetesVersion)) {
		body["kubernetes_version"] = request.KubernetesVersion
	}

	if !tea.BoolValue(util.IsUnset(request.LoadBalancerSpec)) {
		body["load_balancer_spec"] = request.LoadBalancerSpec
	}

	if !tea.BoolValue(util.IsUnset(request.LoggingType)) {
		body["logging_type"] = request.LoggingType
	}

	if !tea.BoolValue(util.IsUnset(request.LoginPassword)) {
		body["login_password"] = request.LoginPassword
	}

	if !tea.BoolValue(util.IsUnset(request.MasterAutoRenew)) {
		body["master_auto_renew"] = request.MasterAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.MasterAutoRenewPeriod)) {
		body["master_auto_renew_period"] = request.MasterAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.MasterCount)) {
		body["master_count"] = request.MasterCount
	}

	if !tea.BoolValue(util.IsUnset(request.MasterInstanceChargeType)) {
		body["master_instance_charge_type"] = request.MasterInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.MasterInstanceTypes)) {
		body["master_instance_types"] = request.MasterInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.MasterPeriod)) {
		body["master_period"] = request.MasterPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.MasterPeriodUnit)) {
		body["master_period_unit"] = request.MasterPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.MasterSystemDiskCategory)) {
		body["master_system_disk_category"] = request.MasterSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.MasterSystemDiskPerformanceLevel)) {
		body["master_system_disk_performance_level"] = request.MasterSystemDiskPerformanceLevel
	}

	if !tea.BoolValue(util.IsUnset(request.MasterSystemDiskSize)) {
		body["master_system_disk_size"] = request.MasterSystemDiskSize
	}

	if !tea.BoolValue(util.IsUnset(request.MasterSystemDiskSnapshotPolicyId)) {
		body["master_system_disk_snapshot_policy_id"] = request.MasterSystemDiskSnapshotPolicyId
	}

	if !tea.BoolValue(util.IsUnset(request.MasterVswitchIds)) {
		body["master_vswitch_ids"] = request.MasterVswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.NatGateway)) {
		body["nat_gateway"] = request.NatGateway
	}

	if !tea.BoolValue(util.IsUnset(request.NodeCidrMask)) {
		body["node_cidr_mask"] = request.NodeCidrMask
	}

	if !tea.BoolValue(util.IsUnset(request.NodeNameMode)) {
		body["node_name_mode"] = request.NodeNameMode
	}

	if !tea.BoolValue(util.IsUnset(request.NodePortRange)) {
		body["node_port_range"] = request.NodePortRange
	}

	if !tea.BoolValue(util.IsUnset(request.Nodepools)) {
		body["nodepools"] = request.Nodepools
	}

	if !tea.BoolValue(util.IsUnset(request.NumOfNodes)) {
		body["num_of_nodes"] = request.NumOfNodes
	}

	if !tea.BoolValue(util.IsUnset(request.OsType)) {
		body["os_type"] = request.OsType
	}

	if !tea.BoolValue(util.IsUnset(request.Period)) {
		body["period"] = request.Period
	}

	if !tea.BoolValue(util.IsUnset(request.PeriodUnit)) {
		body["period_unit"] = request.PeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.Platform)) {
		body["platform"] = request.Platform
	}

	if !tea.BoolValue(util.IsUnset(request.PodVswitchIds)) {
		body["pod_vswitch_ids"] = request.PodVswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.Profile)) {
		body["profile"] = request.Profile
	}

	if !tea.BoolValue(util.IsUnset(request.ProxyMode)) {
		body["proxy_mode"] = request.ProxyMode
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		body["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		body["resource_group_id"] = request.ResourceGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.Runtime)) {
		body["runtime"] = request.Runtime
	}

	if !tea.BoolValue(util.IsUnset(request.SecurityGroupId)) {
		body["security_group_id"] = request.SecurityGroupId
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceAccountIssuer)) {
		body["service_account_issuer"] = request.ServiceAccountIssuer
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceCidr)) {
		body["service_cidr"] = request.ServiceCidr
	}

	if !tea.BoolValue(util.IsUnset(request.ServiceDiscoveryTypes)) {
		body["service_discovery_types"] = request.ServiceDiscoveryTypes
	}

	if !tea.BoolValue(util.IsUnset(request.SnatEntry)) {
		body["snat_entry"] = request.SnatEntry
	}

	if !tea.BoolValue(util.IsUnset(request.SocEnabled)) {
		body["soc_enabled"] = request.SocEnabled
	}

	if !tea.BoolValue(util.IsUnset(request.SshFlags)) {
		body["ssh_flags"] = request.SshFlags
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Taints)) {
		body["taints"] = request.Taints
	}

	if !tea.BoolValue(util.IsUnset(request.TimeoutMins)) {
		body["timeout_mins"] = request.TimeoutMins
	}

	if !tea.BoolValue(util.IsUnset(request.Timezone)) {
		body["timezone"] = request.Timezone
	}

	if !tea.BoolValue(util.IsUnset(request.UserCa)) {
		body["user_ca"] = request.UserCa
	}

	if !tea.BoolValue(util.IsUnset(request.UserData)) {
		body["user_data"] = request.UserData
	}

	if !tea.BoolValue(util.IsUnset(request.Vpcid)) {
		body["vpcid"] = request.Vpcid
	}

	if !tea.BoolValue(util.IsUnset(request.VswitchIds)) {
		body["vswitch_ids"] = request.VswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenew)) {
		body["worker_auto_renew"] = request.WorkerAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenewPeriod)) {
		body["worker_auto_renew_period"] = request.WorkerAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisks)) {
		body["worker_data_disks"] = request.WorkerDataDisks
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceChargeType)) {
		body["worker_instance_charge_type"] = request.WorkerInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceTypes)) {
		body["worker_instance_types"] = request.WorkerInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriod)) {
		body["worker_period"] = request.WorkerPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriodUnit)) {
		body["worker_period_unit"] = request.WorkerPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskCategory)) {
		body["worker_system_disk_category"] = request.WorkerSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskPerformanceLevel)) {
		body["worker_system_disk_performance_level"] = request.WorkerSystemDiskPerformanceLevel
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSize)) {
		body["worker_system_disk_size"] = request.WorkerSystemDiskSize
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSnapshotPolicyId)) {
		body["worker_system_disk_snapshot_policy_id"] = request.WorkerSystemDiskSnapshotPolicyId
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerVswitchIds)) {
		body["worker_vswitch_ids"] = request.WorkerVswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.ZoneId)) {
		body["zone_id"] = request.ZoneId
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateCluster(request *CreateClusterRequest) (_result *CreateClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateClusterResponse{}
	_body, _err := client.CreateClusterWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateClusterNodePoolWithOptions(ClusterId *string, request *CreateClusterNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateClusterNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AutoScaling)) {
		body["auto_scaling"] = request.AutoScaling
	}

	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	if !tea.BoolValue(util.IsUnset(request.InterconnectConfig)) {
		body["interconnect_config"] = request.InterconnectConfig
	}

	if !tea.BoolValue(util.IsUnset(request.InterconnectMode)) {
		body["interconnect_mode"] = request.InterconnectMode
	}

	if !tea.BoolValue(util.IsUnset(request.KubernetesConfig)) {
		body["kubernetes_config"] = request.KubernetesConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Management)) {
		body["management"] = request.Management
	}

	if !tea.BoolValue(util.IsUnset(request.MaxNodes)) {
		body["max_nodes"] = request.MaxNodes
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolInfo)) {
		body["nodepool_info"] = request.NodepoolInfo
	}

	if !tea.BoolValue(util.IsUnset(request.ScalingGroup)) {
		body["scaling_group"] = request.ScalingGroup
	}

	if !tea.BoolValue(util.IsUnset(request.TeeConfig)) {
		body["tee_config"] = request.TeeConfig
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateClusterNodePool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateClusterNodePoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateClusterNodePool(ClusterId *string, request *CreateClusterNodePoolRequest) (_result *CreateClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateClusterNodePoolResponse{}
	_body, _err := client.CreateClusterNodePoolWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateEdgeMachineWithOptions(request *CreateEdgeMachineRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateEdgeMachineResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Hostname)) {
		body["hostname"] = request.Hostname
	}

	if !tea.BoolValue(util.IsUnset(request.Model)) {
		body["model"] = request.Model
	}

	if !tea.BoolValue(util.IsUnset(request.Sn)) {
		body["sn"] = request.Sn
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateEdgeMachine"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/edge_machines"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateEdgeMachineResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateEdgeMachine(request *CreateEdgeMachineRequest) (_result *CreateEdgeMachineResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateEdgeMachineResponse{}
	_body, _err := client.CreateEdgeMachineWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateKubernetesTriggerWithOptions(request *CreateKubernetesTriggerRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateKubernetesTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Action)) {
		body["action"] = request.Action
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterId)) {
		body["cluster_id"] = request.ClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.ProjectId)) {
		body["project_id"] = request.ProjectId
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		body["type"] = request.Type
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateKubernetesTrigger"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/triggers"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateKubernetesTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateKubernetesTrigger(request *CreateKubernetesTriggerRequest) (_result *CreateKubernetesTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateKubernetesTriggerResponse{}
	_body, _err := client.CreateKubernetesTriggerWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateTemplateWithOptions(request *CreateTemplateRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateTemplateResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Template)) {
		body["template"] = request.Template
	}

	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		body["template_type"] = request.TemplateType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateTemplate"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/templates"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &CreateTemplateResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) CreateTemplate(request *CreateTemplateRequest) (_result *CreateTemplateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateTemplateResponse{}
	_body, _err := client.CreateTemplateWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) CreateTriggerWithOptions(clusterId *string, request *CreateTriggerRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *CreateTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Action)) {
		body["action"] = request.Action
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterId)) {
		body["cluster_id"] = request.ClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.ProjectId)) {
		body["project_id"] = request.ProjectId
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		body["type"] = request.Type
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("CreateTrigger"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/triggers"),
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

func (client *Client) CreateTrigger(clusterId *string, request *CreateTriggerRequest) (_result *CreateTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &CreateTriggerResponse{}
	_body, _err := client.CreateTriggerWithOptions(clusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteAlertContactWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteAlertContactResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteAlertContact"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/alert/contacts"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteAlertContactResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteAlertContact() (_result *DeleteAlertContactResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteAlertContactResponse{}
	_body, _err := client.DeleteAlertContactWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteAlertContactGroupWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteAlertContactGroupResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteAlertContactGroup"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/alert/contact_groups"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteAlertContactGroupResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteAlertContactGroup() (_result *DeleteAlertContactGroupResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteAlertContactGroupResponse{}
	_body, _err := client.DeleteAlertContactGroupWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteClusterWithOptions(ClusterId *string, tmpReq *DeleteClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteClusterResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &DeleteClusterShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.RetainResources)) {
		request.RetainResourcesShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.RetainResources, tea.String("retain_resources"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.KeepSlb)) {
		query["keep_slb"] = request.KeepSlb
	}

	if !tea.BoolValue(util.IsUnset(request.RetainAllResources)) {
		query["retain_all_resources"] = request.RetainAllResources
	}

	if !tea.BoolValue(util.IsUnset(request.RetainResourcesShrink)) {
		query["retain_resources"] = request.RetainResourcesShrink
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteCluster(ClusterId *string, request *DeleteClusterRequest) (_result *DeleteClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteClusterResponse{}
	_body, _err := client.DeleteClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteClusterNodepoolWithOptions(ClusterId *string, NodepoolId *string, request *DeleteClusterNodepoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteClusterNodepoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Force)) {
		query["force"] = request.Force
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteClusterNodepool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteClusterNodepoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteClusterNodepool(ClusterId *string, NodepoolId *string, request *DeleteClusterNodepoolRequest) (_result *DeleteClusterNodepoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteClusterNodepoolResponse{}
	_body, _err := client.DeleteClusterNodepoolWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * >
 * *   When you remove a node, the pods that run on the node are migrated to other nodes. This may cause service interruptions. We recommend that you remove nodes during off-peak hours.
 * *   Unknown errors may occur when you remove nodes. Before you remove nodes, back up the data on the nodes.
 * *   Nodes remain in the unschedulable state when they are being removed.
 * *   You can remove only worker nodes by calling this operation.
 *
 * @param request DeleteClusterNodesRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return DeleteClusterNodesResponse
 */
func (client *Client) DeleteClusterNodesWithOptions(ClusterId *string, request *DeleteClusterNodesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteClusterNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DrainNode)) {
		body["drain_node"] = request.DrainNode
	}

	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		body["nodes"] = request.Nodes
	}

	if !tea.BoolValue(util.IsUnset(request.ReleaseNode)) {
		body["release_node"] = request.ReleaseNode
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteClusterNodes"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodes"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DeleteClusterNodesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * >
 * *   When you remove a node, the pods that run on the node are migrated to other nodes. This may cause service interruptions. We recommend that you remove nodes during off-peak hours.
 * *   Unknown errors may occur when you remove nodes. Before you remove nodes, back up the data on the nodes.
 * *   Nodes remain in the unschedulable state when they are being removed.
 * *   You can remove only worker nodes by calling this operation.
 *
 * @param request DeleteClusterNodesRequest
 * @return DeleteClusterNodesResponse
 */
func (client *Client) DeleteClusterNodes(ClusterId *string, request *DeleteClusterNodesRequest) (_result *DeleteClusterNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteClusterNodesResponse{}
	_body, _err := client.DeleteClusterNodesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteEdgeMachineWithOptions(edgeMachineid *string, request *DeleteEdgeMachineRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteEdgeMachineResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Force)) {
		query["force"] = request.Force
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteEdgeMachine"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/edge_machines/%5Bedge_machineid%5D"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteEdgeMachineResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteEdgeMachine(edgeMachineid *string, request *DeleteEdgeMachineRequest) (_result *DeleteEdgeMachineResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteEdgeMachineResponse{}
	_body, _err := client.DeleteEdgeMachineWithOptions(edgeMachineid, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteKubernetesTriggerWithOptions(Id *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteKubernetesTriggerResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteKubernetesTrigger"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/triggers/revoke/" + tea.StringValue(openapiutil.GetEncodeParam(Id))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteKubernetesTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteKubernetesTrigger(Id *string) (_result *DeleteKubernetesTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteKubernetesTriggerResponse{}
	_body, _err := client.DeleteKubernetesTriggerWithOptions(Id, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeletePolicyInstanceWithOptions(clusterId *string, policyName *string, request *DeletePolicyInstanceRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeletePolicyInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.InstanceName)) {
		query["instance_name"] = request.InstanceName
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DeletePolicyInstance"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/policies/" + tea.StringValue(openapiutil.GetEncodeParam(policyName))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DeletePolicyInstanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeletePolicyInstance(clusterId *string, policyName *string, request *DeletePolicyInstanceRequest) (_result *DeletePolicyInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeletePolicyInstanceResponse{}
	_body, _err := client.DeletePolicyInstanceWithOptions(clusterId, policyName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteTemplateWithOptions(TemplateId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteTemplateResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteTemplate"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/templates/" + tea.StringValue(openapiutil.GetEncodeParam(TemplateId))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &DeleteTemplateResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeleteTemplate(TemplateId *string) (_result *DeleteTemplateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteTemplateResponse{}
	_body, _err := client.DeleteTemplateWithOptions(TemplateId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeleteTriggerWithOptions(clusterId *string, Id *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeleteTriggerResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DeleteTrigger"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/%5Bcluster_id%5D/triggers/%5BId%5D"),
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

func (client *Client) DeleteTrigger(clusterId *string, Id *string) (_result *DeleteTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeleteTriggerResponse{}
	_body, _err := client.DeleteTriggerWithOptions(clusterId, Id, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DeployPolicyInstanceWithOptions(clusterId *string, policyName *string, request *DeployPolicyInstanceRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DeployPolicyInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Action)) {
		body["action"] = request.Action
	}

	if !tea.BoolValue(util.IsUnset(request.Namespaces)) {
		body["namespaces"] = request.Namespaces
	}

	if !tea.BoolValue(util.IsUnset(request.Parameters)) {
		body["parameters"] = request.Parameters
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("DeployPolicyInstance"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/policies/" + tea.StringValue(openapiutil.GetEncodeParam(policyName))),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DeployPolicyInstanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DeployPolicyInstance(clusterId *string, policyName *string, request *DeployPolicyInstanceRequest) (_result *DeployPolicyInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DeployPolicyInstanceResponse{}
	_body, _err := client.DeployPolicyInstanceWithOptions(clusterId, policyName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescirbeWorkflowWithOptions(workflowName *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescirbeWorkflowResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescirbeWorkflow"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/gs/workflow/" + tea.StringValue(openapiutil.GetEncodeParam(workflowName))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescirbeWorkflowResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescirbeWorkflow(workflowName *string) (_result *DescirbeWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescirbeWorkflowResponse{}
	_body, _err := client.DescirbeWorkflowWithOptions(workflowName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeAddonsWithOptions(request *DescribeAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClusterProfile)) {
		query["cluster_profile"] = request.ClusterProfile
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterSpec)) {
		query["cluster_spec"] = request.ClusterSpec
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["cluster_type"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterVersion)) {
		query["cluster_version"] = request.ClusterVersion
	}

	if !tea.BoolValue(util.IsUnset(request.Region)) {
		query["region"] = request.Region
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeAddons"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/components/metadata"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeAddonsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeAddons(request *DescribeAddonsRequest) (_result *DescribeAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeAddonsResponse{}
	_body, _err := client.DescribeAddonsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAddonInstanceWithOptions(ClusterID *string, AddonName *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonInstanceResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterAddonInstance"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterID)) + "/components/" + tea.StringValue(openapiutil.GetEncodeParam(AddonName)) + "/instance"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterAddonInstanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAddonInstance(ClusterID *string, AddonName *string) (_result *DescribeClusterAddonInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonInstanceResponse{}
	_body, _err := client.DescribeClusterAddonInstanceWithOptions(ClusterID, AddonName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAddonMetadataWithOptions(clusterId *string, componentId *string, version *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonMetadataResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterAddonMetadata"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/components/" + tea.StringValue(openapiutil.GetEncodeParam(componentId)) + "/metadata"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterAddonMetadataResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAddonMetadata(clusterId *string, componentId *string, version *string) (_result *DescribeClusterAddonMetadataResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonMetadataResponse{}
	_body, _err := client.DescribeClusterAddonMetadataWithOptions(clusterId, componentId, version, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * @deprecated
 *
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return DescribeClusterAddonUpgradeStatusResponse
 */
// Deprecated
func (client *Client) DescribeClusterAddonUpgradeStatusWithOptions(ClusterId *string, ComponentId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonUpgradeStatusResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterAddonUpgradeStatus"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/components/" + tea.StringValue(openapiutil.GetEncodeParam(ComponentId)) + "/upgradestatus"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterAddonUpgradeStatusResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * @deprecated
 *
 * @return DescribeClusterAddonUpgradeStatusResponse
 */
// Deprecated
func (client *Client) DescribeClusterAddonUpgradeStatus(ClusterId *string, ComponentId *string) (_result *DescribeClusterAddonUpgradeStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonUpgradeStatusResponse{}
	_body, _err := client.DescribeClusterAddonUpgradeStatusWithOptions(ClusterId, ComponentId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAddonsUpgradeStatusWithOptions(ClusterId *string, tmpReq *DescribeClusterAddonsUpgradeStatusRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonsUpgradeStatusResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &DescribeClusterAddonsUpgradeStatusShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.ComponentIds)) {
		request.ComponentIdsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.ComponentIds, tea.String("componentIds"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ComponentIdsShrink)) {
		query["componentIds"] = request.ComponentIdsShrink
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterAddonsUpgradeStatus"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/components/upgradestatus"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterAddonsUpgradeStatusResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAddonsUpgradeStatus(ClusterId *string, request *DescribeClusterAddonsUpgradeStatusRequest) (_result *DescribeClusterAddonsUpgradeStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonsUpgradeStatusResponse{}
	_body, _err := client.DescribeClusterAddonsUpgradeStatusWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAddonsVersionWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAddonsVersionResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterAddonsVersion"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/components/version"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterAddonsVersionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAddonsVersion(ClusterId *string) (_result *DescribeClusterAddonsVersionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAddonsVersionResponse{}
	_body, _err := client.DescribeClusterAddonsVersionWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterAttachScriptsWithOptions(ClusterId *string, request *DescribeClusterAttachScriptsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterAttachScriptsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Arch)) {
		body["arch"] = request.Arch
	}

	if !tea.BoolValue(util.IsUnset(request.FormatDisk)) {
		body["format_disk"] = request.FormatDisk
	}

	if !tea.BoolValue(util.IsUnset(request.KeepInstanceName)) {
		body["keep_instance_name"] = request.KeepInstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolId)) {
		body["nodepool_id"] = request.NodepoolId
	}

	if !tea.BoolValue(util.IsUnset(request.Options)) {
		body["options"] = request.Options
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterAttachScripts"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/attachscript"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("string"),
	}
	_result = &DescribeClusterAttachScriptsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterAttachScripts(ClusterId *string, request *DescribeClusterAttachScriptsRequest) (_result *DescribeClusterAttachScriptsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterAttachScriptsResponse{}
	_body, _err := client.DescribeClusterAttachScriptsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterDetailWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterDetailResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterDetail"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterDetailResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterDetail(ClusterId *string) (_result *DescribeClusterDetailResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterDetailResponse{}
	_body, _err := client.DescribeClusterDetailWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterEventsWithOptions(ClusterId *string, request *DescribeClusterEventsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterEventsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["page_number"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.TaskId)) {
		query["task_id"] = request.TaskId
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterEvents"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/events"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterEventsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterEvents(ClusterId *string, request *DescribeClusterEventsRequest) (_result *DescribeClusterEventsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterEventsResponse{}
	_body, _err := client.DescribeClusterEventsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterLogsWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterLogsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterLogs"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/logs"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeClusterLogsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterLogs(ClusterId *string) (_result *DescribeClusterLogsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterLogsResponse{}
	_body, _err := client.DescribeClusterLogsWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterNodePoolDetailWithOptions(ClusterId *string, NodepoolId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterNodePoolDetailResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterNodePoolDetail"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterNodePoolDetailResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterNodePoolDetail(ClusterId *string, NodepoolId *string) (_result *DescribeClusterNodePoolDetailResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterNodePoolDetailResponse{}
	_body, _err := client.DescribeClusterNodePoolDetailWithOptions(ClusterId, NodepoolId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterNodePoolsWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterNodePoolsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterNodePools"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterNodePoolsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterNodePools(ClusterId *string) (_result *DescribeClusterNodePoolsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterNodePoolsResponse{}
	_body, _err := client.DescribeClusterNodePoolsWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterNodesWithOptions(ClusterId *string, request *DescribeClusterNodesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.InstanceIds)) {
		query["instanceIds"] = request.InstanceIds
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolId)) {
		query["nodepool_id"] = request.NodepoolId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["pageNumber"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["pageSize"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.State)) {
		query["state"] = request.State
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterNodes"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodes"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterNodesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterNodes(ClusterId *string, request *DescribeClusterNodesRequest) (_result *DescribeClusterNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterNodesResponse{}
	_body, _err := client.DescribeClusterNodesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterResourcesWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterResourcesResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterResources"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/resources"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeClusterResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterResources(ClusterId *string) (_result *DescribeClusterResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterResourcesResponse{}
	_body, _err := client.DescribeClusterResourcesWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterTasksWithOptions(clusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterTasksResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterTasks"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/tasks"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterTasksResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterTasks(clusterId *string) (_result *DescribeClusterTasksResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterTasksResponse{}
	_body, _err := client.DescribeClusterTasksWithOptions(clusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterUserKubeconfigWithOptions(ClusterId *string, request *DescribeClusterUserKubeconfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterUserKubeconfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PrivateIpAddress)) {
		query["PrivateIpAddress"] = request.PrivateIpAddress
	}

	if !tea.BoolValue(util.IsUnset(request.TemporaryDurationMinutes)) {
		query["TemporaryDurationMinutes"] = request.TemporaryDurationMinutes
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterUserKubeconfig"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/k8s/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/user_config"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterUserKubeconfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterUserKubeconfig(ClusterId *string, request *DescribeClusterUserKubeconfigRequest) (_result *DescribeClusterUserKubeconfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterUserKubeconfigResponse{}
	_body, _err := client.DescribeClusterUserKubeconfigWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request DescribeClusterV2UserKubeconfigRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return DescribeClusterV2UserKubeconfigResponse
 */
// Deprecated
func (client *Client) DescribeClusterV2UserKubeconfigWithOptions(ClusterId *string, request *DescribeClusterV2UserKubeconfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterV2UserKubeconfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PrivateIpAddress)) {
		query["PrivateIpAddress"] = request.PrivateIpAddress
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterV2UserKubeconfig"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/k8s/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/user_config"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterV2UserKubeconfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request DescribeClusterV2UserKubeconfigRequest
 * @return DescribeClusterV2UserKubeconfigResponse
 */
// Deprecated
func (client *Client) DescribeClusterV2UserKubeconfig(ClusterId *string, request *DescribeClusterV2UserKubeconfigRequest) (_result *DescribeClusterV2UserKubeconfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterV2UserKubeconfigResponse{}
	_body, _err := client.DescribeClusterV2UserKubeconfigWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClusterVulsWithOptions(clusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClusterVulsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusterVuls"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/vuls"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClusterVulsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClusterVuls(clusterId *string) (_result *DescribeClusterVulsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClusterVulsResponse{}
	_body, _err := client.DescribeClusterVulsWithOptions(clusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request DescribeClustersRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return DescribeClustersResponse
 */
// Deprecated
func (client *Client) DescribeClustersWithOptions(request *DescribeClustersRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClustersResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["clusterType"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["name"] = request.Name
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClusters"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeClustersResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request DescribeClustersRequest
 * @return DescribeClustersResponse
 */
// Deprecated
func (client *Client) DescribeClusters(request *DescribeClustersRequest) (_result *DescribeClustersResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClustersResponse{}
	_body, _err := client.DescribeClustersWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeClustersV1WithOptions(request *DescribeClustersV1Request, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeClustersV1Response, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClusterSpec)) {
		query["cluster_spec"] = request.ClusterSpec
	}

	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["cluster_type"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["page_number"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.Profile)) {
		query["profile"] = request.Profile
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["region_id"] = request.RegionId
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeClustersV1"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v1/clusters"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeClustersV1Response{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeClustersV1(request *DescribeClustersV1Request) (_result *DescribeClustersV1Response, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeClustersV1Response{}
	_body, _err := client.DescribeClustersV1WithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeEdgeMachineActiveProcessWithOptions(edgeMachineid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeEdgeMachineActiveProcessResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeEdgeMachineActiveProcess"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/edge_machines/%5Bedge_machineid%5D/activeprocess"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeEdgeMachineActiveProcessResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeEdgeMachineActiveProcess(edgeMachineid *string) (_result *DescribeEdgeMachineActiveProcessResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeEdgeMachineActiveProcessResponse{}
	_body, _err := client.DescribeEdgeMachineActiveProcessWithOptions(edgeMachineid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeEdgeMachineModelsWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeEdgeMachineModelsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeEdgeMachineModels"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/edge_machines/models"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeEdgeMachineModelsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeEdgeMachineModels() (_result *DescribeEdgeMachineModelsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeEdgeMachineModelsResponse{}
	_body, _err := client.DescribeEdgeMachineModelsWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeEdgeMachineTunnelConfigDetailWithOptions(edgeMachineid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeEdgeMachineTunnelConfigDetailResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeEdgeMachineTunnelConfigDetail"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/edge_machines/%5Bedge_machineid%5D/tunnelconfig"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeEdgeMachineTunnelConfigDetailResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeEdgeMachineTunnelConfigDetail(edgeMachineid *string) (_result *DescribeEdgeMachineTunnelConfigDetailResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeEdgeMachineTunnelConfigDetailResponse{}
	_body, _err := client.DescribeEdgeMachineTunnelConfigDetailWithOptions(edgeMachineid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeEdgeMachinesWithOptions(request *DescribeEdgeMachinesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeEdgeMachinesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Hostname)) {
		query["hostname"] = request.Hostname
	}

	if !tea.BoolValue(util.IsUnset(request.LifeState)) {
		query["life_state"] = request.LifeState
	}

	if !tea.BoolValue(util.IsUnset(request.Model)) {
		query["model"] = request.Model
	}

	if !tea.BoolValue(util.IsUnset(request.OnlineState)) {
		query["online_state"] = request.OnlineState
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["page_number"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeEdgeMachines"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/edge_machines"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeEdgeMachinesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeEdgeMachines(request *DescribeEdgeMachinesRequest) (_result *DescribeEdgeMachinesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeEdgeMachinesResponse{}
	_body, _err := client.DescribeEdgeMachinesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeEventsWithOptions(request *DescribeEventsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeEventsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClusterId)) {
		query["cluster_id"] = request.ClusterId
	}

	if !tea.BoolValue(util.IsUnset(request.PageNumber)) {
		query["page_number"] = request.PageNumber
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		query["type"] = request.Type
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeEvents"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/events"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeEventsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeEvents(request *DescribeEventsRequest) (_result *DescribeEventsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeEventsResponse{}
	_body, _err := client.DescribeEventsWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * For more information, see [Register an external Kubernetes cluster](~~121053~~).
 *
 * @param request DescribeExternalAgentRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return DescribeExternalAgentResponse
 */
func (client *Client) DescribeExternalAgentWithOptions(ClusterId *string, request *DescribeExternalAgentRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeExternalAgentResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AgentMode)) {
		query["AgentMode"] = request.AgentMode
	}

	if !tea.BoolValue(util.IsUnset(request.PrivateIpAddress)) {
		query["PrivateIpAddress"] = request.PrivateIpAddress
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeExternalAgent"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/k8s/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/external/agent/deployment"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeExternalAgentResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * For more information, see [Register an external Kubernetes cluster](~~121053~~).
 *
 * @param request DescribeExternalAgentRequest
 * @return DescribeExternalAgentResponse
 */
func (client *Client) DescribeExternalAgent(ClusterId *string, request *DescribeExternalAgentRequest) (_result *DescribeExternalAgentResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeExternalAgentResponse{}
	_body, _err := client.DescribeExternalAgentWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeKubernetesVersionMetadataWithOptions(request *DescribeKubernetesVersionMetadataRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeKubernetesVersionMetadataResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ClusterType)) {
		query["ClusterType"] = request.ClusterType
	}

	if !tea.BoolValue(util.IsUnset(request.KubernetesVersion)) {
		query["KubernetesVersion"] = request.KubernetesVersion
	}

	if !tea.BoolValue(util.IsUnset(request.Mode)) {
		query["Mode"] = request.Mode
	}

	if !tea.BoolValue(util.IsUnset(request.Profile)) {
		query["Profile"] = request.Profile
	}

	if !tea.BoolValue(util.IsUnset(request.Region)) {
		query["Region"] = request.Region
	}

	if !tea.BoolValue(util.IsUnset(request.Runtime)) {
		query["runtime"] = request.Runtime
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeKubernetesVersionMetadata"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v1/metadata/versions"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeKubernetesVersionMetadataResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeKubernetesVersionMetadata(request *DescribeKubernetesVersionMetadataRequest) (_result *DescribeKubernetesVersionMetadataResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeKubernetesVersionMetadataResponse{}
	_body, _err := client.DescribeKubernetesVersionMetadataWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeNodePoolVulsWithOptions(clusterId *string, nodepoolId *string, request *DescribeNodePoolVulsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeNodePoolVulsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Necessity)) {
		query["necessity"] = request.Necessity
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeNodePoolVuls"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(nodepoolId)) + "/vuls"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeNodePoolVulsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeNodePoolVuls(clusterId *string, nodepoolId *string, request *DescribeNodePoolVulsRequest) (_result *DescribeNodePoolVulsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeNodePoolVulsResponse{}
	_body, _err := client.DescribeNodePoolVulsWithOptions(clusterId, nodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePoliciesWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribePoliciesResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePolicies"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/policies"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribePoliciesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePolicies() (_result *DescribePoliciesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribePoliciesResponse{}
	_body, _err := client.DescribePoliciesWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePolicyDetailsWithOptions(policyName *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribePolicyDetailsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePolicyDetails"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/policies/" + tea.StringValue(openapiutil.GetEncodeParam(policyName))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribePolicyDetailsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePolicyDetails(policyName *string) (_result *DescribePolicyDetailsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribePolicyDetailsResponse{}
	_body, _err := client.DescribePolicyDetailsWithOptions(policyName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePolicyGovernanceInClusterWithOptions(clusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribePolicyGovernanceInClusterResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePolicyGovernanceInCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/policygovernance"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribePolicyGovernanceInClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePolicyGovernanceInCluster(clusterId *string) (_result *DescribePolicyGovernanceInClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribePolicyGovernanceInClusterResponse{}
	_body, _err := client.DescribePolicyGovernanceInClusterWithOptions(clusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePolicyInstancesWithOptions(clusterId *string, request *DescribePolicyInstancesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribePolicyInstancesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.InstanceName)) {
		query["instance_name"] = request.InstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.PolicyName)) {
		query["policy_name"] = request.PolicyName
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePolicyInstances"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/policies"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribePolicyInstancesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePolicyInstances(clusterId *string, request *DescribePolicyInstancesRequest) (_result *DescribePolicyInstancesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribePolicyInstancesResponse{}
	_body, _err := client.DescribePolicyInstancesWithOptions(clusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribePolicyInstancesStatusWithOptions(clusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribePolicyInstancesStatusResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribePolicyInstancesStatus"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/policies/status"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribePolicyInstancesStatusResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribePolicyInstancesStatus(clusterId *string) (_result *DescribePolicyInstancesStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribePolicyInstancesStatusResponse{}
	_body, _err := client.DescribePolicyInstancesStatusWithOptions(clusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * >  You can call this operation only with an Alibaba Cloud account.
 *
 * @param request DescribeSubaccountK8sClusterUserConfigRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return DescribeSubaccountK8sClusterUserConfigResponse
 */
func (client *Client) DescribeSubaccountK8sClusterUserConfigWithOptions(ClusterId *string, Uid *string, request *DescribeSubaccountK8sClusterUserConfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeSubaccountK8sClusterUserConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PrivateIpAddress)) {
		query["PrivateIpAddress"] = request.PrivateIpAddress
	}

	if !tea.BoolValue(util.IsUnset(request.TemporaryDurationMinutes)) {
		query["TemporaryDurationMinutes"] = request.TemporaryDurationMinutes
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeSubaccountK8sClusterUserConfig"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/k8s/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/users/" + tea.StringValue(openapiutil.GetEncodeParam(Uid)) + "/user_config"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeSubaccountK8sClusterUserConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * >  You can call this operation only with an Alibaba Cloud account.
 *
 * @param request DescribeSubaccountK8sClusterUserConfigRequest
 * @return DescribeSubaccountK8sClusterUserConfigResponse
 */
func (client *Client) DescribeSubaccountK8sClusterUserConfig(ClusterId *string, Uid *string, request *DescribeSubaccountK8sClusterUserConfigRequest) (_result *DescribeSubaccountK8sClusterUserConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeSubaccountK8sClusterUserConfigResponse{}
	_body, _err := client.DescribeSubaccountK8sClusterUserConfigWithOptions(ClusterId, Uid, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTaskInfoWithOptions(taskId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeTaskInfoResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeTaskInfo"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/tasks/" + tea.StringValue(openapiutil.GetEncodeParam(taskId))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeTaskInfoResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTaskInfo(taskId *string) (_result *DescribeTaskInfoResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeTaskInfoResponse{}
	_body, _err := client.DescribeTaskInfoWithOptions(taskId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTemplateAttributeWithOptions(TemplateId *string, request *DescribeTemplateAttributeRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeTemplateAttributeResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		query["template_type"] = request.TemplateType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeTemplateAttribute"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/templates/" + tea.StringValue(openapiutil.GetEncodeParam(TemplateId))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeTemplateAttributeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTemplateAttribute(TemplateId *string, request *DescribeTemplateAttributeRequest) (_result *DescribeTemplateAttributeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeTemplateAttributeResponse{}
	_body, _err := client.DescribeTemplateAttributeWithOptions(TemplateId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTemplatesWithOptions(request *DescribeTemplatesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeTemplatesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.PageNum)) {
		query["page_num"] = request.PageNum
	}

	if !tea.BoolValue(util.IsUnset(request.PageSize)) {
		query["page_size"] = request.PageSize
	}

	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		query["template_type"] = request.TemplateType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeTemplates"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/templates"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeTemplatesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTemplates(request *DescribeTemplatesRequest) (_result *DescribeTemplatesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeTemplatesResponse{}
	_body, _err := client.DescribeTemplatesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeTriggerWithOptions(clusterId *string, request *DescribeTriggerRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["Name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.Namespace)) {
		query["Namespace"] = request.Namespace
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		query["Type"] = request.Type
	}

	if !tea.BoolValue(util.IsUnset(request.Action)) {
		query["action"] = request.Action
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeTrigger"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/triggers"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeTrigger(clusterId *string, request *DescribeTriggerRequest) (_result *DescribeTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeTriggerResponse{}
	_body, _err := client.DescribeTriggerWithOptions(clusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeUserClusterNamespacesWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeUserClusterNamespacesResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeUserClusterNamespaces"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/k8s/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/namespaces"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeUserClusterNamespacesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeUserClusterNamespaces(ClusterId *string) (_result *DescribeUserClusterNamespacesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeUserClusterNamespacesResponse{}
	_body, _err := client.DescribeUserClusterNamespacesWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeUserPermissionWithOptions(uid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeUserPermissionResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeUserPermission"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/permissions/users/" + tea.StringValue(openapiutil.GetEncodeParam(uid))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &DescribeUserPermissionResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeUserPermission(uid *string) (_result *DescribeUserPermissionResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeUserPermissionResponse{}
	_body, _err := client.DescribeUserPermissionWithOptions(uid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeUserQuotaWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeUserQuotaResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeUserQuota"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/quota"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeUserQuotaResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeUserQuota() (_result *DescribeUserQuotaResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeUserQuotaResponse{}
	_body, _err := client.DescribeUserQuotaWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) DescribeWorkflowsWithOptions(headers map[string]*string, runtime *util.RuntimeOptions) (_result *DescribeWorkflowsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("DescribeWorkflows"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/gs/workflows"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &DescribeWorkflowsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) DescribeWorkflows() (_result *DescribeWorkflowsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &DescribeWorkflowsResponse{}
	_body, _err := client.DescribeWorkflowsWithOptions(headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) EdgeClusterAddEdgeMachineWithOptions(clusterid *string, edgeMachineid *string, request *EdgeClusterAddEdgeMachineRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *EdgeClusterAddEdgeMachineResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Expired)) {
		body["expired"] = request.Expired
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolId)) {
		body["nodepool_id"] = request.NodepoolId
	}

	if !tea.BoolValue(util.IsUnset(request.Options)) {
		body["options"] = request.Options
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("EdgeClusterAddEdgeMachine"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/%5Bclusterid%5D/attachedgemachine/%5Bedge_machineid%5D"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &EdgeClusterAddEdgeMachineResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) EdgeClusterAddEdgeMachine(clusterid *string, edgeMachineid *string, request *EdgeClusterAddEdgeMachineRequest) (_result *EdgeClusterAddEdgeMachineResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &EdgeClusterAddEdgeMachineResponse{}
	_body, _err := client.EdgeClusterAddEdgeMachineWithOptions(clusterid, edgeMachineid, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) FixNodePoolVulsWithOptions(clusterId *string, nodepoolId *string, request *FixNodePoolVulsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *FixNodePoolVulsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		body["nodes"] = request.Nodes
	}

	if !tea.BoolValue(util.IsUnset(request.RolloutPolicy)) {
		body["rollout_policy"] = request.RolloutPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.Vuls)) {
		body["vuls"] = request.Vuls
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("FixNodePoolVuls"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(nodepoolId)) + "/vuls/fix"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &FixNodePoolVulsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) FixNodePoolVuls(clusterId *string, nodepoolId *string, request *FixNodePoolVulsRequest) (_result *FixNodePoolVulsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &FixNodePoolVulsResponse{}
	_body, _err := client.FixNodePoolVulsWithOptions(clusterId, nodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetKubernetesTriggerWithOptions(ClusterId *string, request *GetKubernetesTriggerRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *GetKubernetesTriggerResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Name)) {
		query["Name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.Namespace)) {
		query["Namespace"] = request.Namespace
	}

	if !tea.BoolValue(util.IsUnset(request.Type)) {
		query["Type"] = request.Type
	}

	if !tea.BoolValue(util.IsUnset(request.Action)) {
		query["action"] = request.Action
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("GetKubernetesTrigger"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/triggers/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId))),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("array"),
	}
	_result = &GetKubernetesTriggerResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetKubernetesTrigger(ClusterId *string, request *GetKubernetesTriggerRequest) (_result *GetKubernetesTriggerResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &GetKubernetesTriggerResponse{}
	_body, _err := client.GetKubernetesTriggerWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) GetUpgradeStatusWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *GetUpgradeStatusResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("GetUpgradeStatus"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/upgrade/status"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &GetUpgradeStatusResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) GetUpgradeStatus(ClusterId *string) (_result *GetUpgradeStatusResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &GetUpgradeStatusResponse{}
	_body, _err := client.GetUpgradeStatusWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * **Precautions**:
 * *   Make sure that you have granted the specified RAM user at least read permissions on the specified cluster by attaching RAM policies. Otherwise, the `ErrorRamPolicyConfig` error will be returned.
 *     For more information about how to authorize a RAM user by attaching RAM policies, see [Create a custom RAM policy](~~86485~~).
 * *   If you call this operation as a RAM user, make sure that this RAM user has the permissions to grant other RAM users the permissions to manage ACK clusters. Otherwise, the `StatusForbidden` or `ForbiddenGrantPermissions` errors will be returned. For more information, see [Use a RAM user to grant RBAC permissions to other RAM users](~~119035~~).
 * *   This operation overwrites the permissions that have been granted to the specified RAM user. When you call this operation, make sure that the required permissions are included.
 *
 * @param request GrantPermissionsRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return GrantPermissionsResponse
 */
func (client *Client) GrantPermissionsWithOptions(uid *string, request *GrantPermissionsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *GrantPermissionsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	params := &openapi.Params{
		Action:      tea.String("GrantPermissions"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/permissions/users/" + tea.StringValue(openapiutil.GetEncodeParam(uid))),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &GrantPermissionsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * **Precautions**:
 * *   Make sure that you have granted the specified RAM user at least read permissions on the specified cluster by attaching RAM policies. Otherwise, the `ErrorRamPolicyConfig` error will be returned.
 *     For more information about how to authorize a RAM user by attaching RAM policies, see [Create a custom RAM policy](~~86485~~).
 * *   If you call this operation as a RAM user, make sure that this RAM user has the permissions to grant other RAM users the permissions to manage ACK clusters. Otherwise, the `StatusForbidden` or `ForbiddenGrantPermissions` errors will be returned. For more information, see [Use a RAM user to grant RBAC permissions to other RAM users](~~119035~~).
 * *   This operation overwrites the permissions that have been granted to the specified RAM user. When you call this operation, make sure that the required permissions are included.
 *
 * @param request GrantPermissionsRequest
 * @return GrantPermissionsResponse
 */
func (client *Client) GrantPermissions(uid *string, request *GrantPermissionsRequest) (_result *GrantPermissionsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &GrantPermissionsResponse{}
	_body, _err := client.GrantPermissionsWithOptions(uid, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) InstallClusterAddonsWithOptions(ClusterId *string, request *InstallClusterAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *InstallClusterAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	params := &openapi.Params{
		Action:      tea.String("InstallClusterAddons"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/components/install"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &InstallClusterAddonsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) InstallClusterAddons(ClusterId *string, request *InstallClusterAddonsRequest) (_result *InstallClusterAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &InstallClusterAddonsResponse{}
	_body, _err := client.InstallClusterAddonsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ListTagResourcesWithOptions(tmpReq *ListTagResourcesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ListTagResourcesResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &ListTagResourcesShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.ResourceIds)) {
		request.ResourceIdsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.ResourceIds, tea.String("resource_ids"), tea.String("json"))
	}

	if !tea.BoolValue(util.IsUnset(tmpReq.Tags)) {
		request.TagsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.Tags, tea.String("tags"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.NextToken)) {
		query["next_token"] = request.NextToken
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceIdsShrink)) {
		query["resource_ids"] = request.ResourceIdsShrink
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		query["resource_type"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.TagsShrink)) {
		query["tags"] = request.TagsShrink
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("ListTagResources"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/tags"),
		Method:      tea.String("GET"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ListTagResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ListTagResources(request *ListTagResourcesRequest) (_result *ListTagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ListTagResourcesResponse{}
	_body, _err := client.ListTagResourcesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) MigrateClusterWithOptions(clusterId *string, request *MigrateClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *MigrateClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.OssBucketEndpoint)) {
		body["oss_bucket_endpoint"] = request.OssBucketEndpoint
	}

	if !tea.BoolValue(util.IsUnset(request.OssBucketName)) {
		body["oss_bucket_name"] = request.OssBucketName
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("MigrateCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/migrate"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &MigrateClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) MigrateCluster(clusterId *string, request *MigrateClusterRequest) (_result *MigrateClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &MigrateClusterResponse{}
	_body, _err := client.MigrateClusterWithOptions(clusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterWithOptions(ClusterId *string, request *ModifyClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ApiServerEip)) {
		body["api_server_eip"] = request.ApiServerEip
	}

	if !tea.BoolValue(util.IsUnset(request.ApiServerEipId)) {
		body["api_server_eip_id"] = request.ApiServerEipId
	}

	if !tea.BoolValue(util.IsUnset(request.DeletionProtection)) {
		body["deletion_protection"] = request.DeletionProtection
	}

	if !tea.BoolValue(util.IsUnset(request.EnableRrsa)) {
		body["enable_rrsa"] = request.EnableRrsa
	}

	if !tea.BoolValue(util.IsUnset(request.IngressDomainRebinding)) {
		body["ingress_domain_rebinding"] = request.IngressDomainRebinding
	}

	if !tea.BoolValue(util.IsUnset(request.IngressLoadbalancerId)) {
		body["ingress_loadbalancer_id"] = request.IngressLoadbalancerId
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceDeletionProtection)) {
		body["instance_deletion_protection"] = request.InstanceDeletionProtection
	}

	if !tea.BoolValue(util.IsUnset(request.MaintenanceWindow)) {
		body["maintenance_window"] = request.MaintenanceWindow
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceGroupId)) {
		body["resource_group_id"] = request.ResourceGroupId
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyCluster(ClusterId *string, request *ModifyClusterRequest) (_result *ModifyClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterResponse{}
	_body, _err := client.ModifyClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * You can use this API operation to modify the components in a Container Service for Kubernetes (ACK) cluster or the control plane components in an ACK Pro cluster.
 * *   To query the customizable parameters of a component, call the `DescribeClusterAddonMetadata` API operation. For more information, see [Query the metadata of a specified component version](https://www.alibabacloud.com/help/zh/container-service-for-kubernetes/latest/query).
 * *   For more information about the customizable parameters of control plane components in ACK Pro clusters, see [Customize the parameters of control plane components in ACK Pro clusters](https://www.alibabacloud.com/help/zh/container-service-for-kubernetes/latest/customize-control-plane-parameters-for-a-professional-kubernetes-cluster).
 * After you call this operation, the component may be redeployed and restarted. We recommend that you evaluate the impact before you call this operation.
 *
 * @param request ModifyClusterAddonRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return ModifyClusterAddonResponse
 */
func (client *Client) ModifyClusterAddonWithOptions(clusterId *string, componentId *string, request *ModifyClusterAddonRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterAddonResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Config)) {
		body["config"] = request.Config
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyClusterAddon"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/components/" + tea.StringValue(openapiutil.GetEncodeParam(componentId)) + "/config"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &ModifyClusterAddonResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * You can use this API operation to modify the components in a Container Service for Kubernetes (ACK) cluster or the control plane components in an ACK Pro cluster.
 * *   To query the customizable parameters of a component, call the `DescribeClusterAddonMetadata` API operation. For more information, see [Query the metadata of a specified component version](https://www.alibabacloud.com/help/zh/container-service-for-kubernetes/latest/query).
 * *   For more information about the customizable parameters of control plane components in ACK Pro clusters, see [Customize the parameters of control plane components in ACK Pro clusters](https://www.alibabacloud.com/help/zh/container-service-for-kubernetes/latest/customize-control-plane-parameters-for-a-professional-kubernetes-cluster).
 * After you call this operation, the component may be redeployed and restarted. We recommend that you evaluate the impact before you call this operation.
 *
 * @param request ModifyClusterAddonRequest
 * @return ModifyClusterAddonResponse
 */
func (client *Client) ModifyClusterAddon(clusterId *string, componentId *string, request *ModifyClusterAddonRequest) (_result *ModifyClusterAddonResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterAddonResponse{}
	_body, _err := client.ModifyClusterAddonWithOptions(clusterId, componentId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterConfigurationWithOptions(ClusterId *string, request *ModifyClusterConfigurationRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterConfigurationResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CustomizeConfig)) {
		body["customize_config"] = request.CustomizeConfig
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyClusterConfiguration"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/configuration"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &ModifyClusterConfigurationResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyClusterConfiguration(ClusterId *string, request *ModifyClusterConfigurationRequest) (_result *ModifyClusterConfigurationResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterConfigurationResponse{}
	_body, _err := client.ModifyClusterConfigurationWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterNodePoolWithOptions(ClusterId *string, NodepoolId *string, request *ModifyClusterNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.AutoScaling)) {
		body["auto_scaling"] = request.AutoScaling
	}

	if !tea.BoolValue(util.IsUnset(request.KubernetesConfig)) {
		body["kubernetes_config"] = request.KubernetesConfig
	}

	if !tea.BoolValue(util.IsUnset(request.Management)) {
		body["management"] = request.Management
	}

	if !tea.BoolValue(util.IsUnset(request.NodepoolInfo)) {
		body["nodepool_info"] = request.NodepoolInfo
	}

	if !tea.BoolValue(util.IsUnset(request.ScalingGroup)) {
		body["scaling_group"] = request.ScalingGroup
	}

	if !tea.BoolValue(util.IsUnset(request.TeeConfig)) {
		body["tee_config"] = request.TeeConfig
	}

	if !tea.BoolValue(util.IsUnset(request.UpdateNodes)) {
		body["update_nodes"] = request.UpdateNodes
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyClusterNodePool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyClusterNodePoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyClusterNodePool(ClusterId *string, NodepoolId *string, request *ModifyClusterNodePoolRequest) (_result *ModifyClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterNodePoolResponse{}
	_body, _err := client.ModifyClusterNodePoolWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyClusterTagsWithOptions(ClusterId *string, request *ModifyClusterTagsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyClusterTagsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyClusterTags"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/tags"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &ModifyClusterTagsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyClusterTags(ClusterId *string, request *ModifyClusterTagsRequest) (_result *ModifyClusterTagsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyClusterTagsResponse{}
	_body, _err := client.ModifyClusterTagsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyNodePoolNodeConfigWithOptions(ClusterId *string, NodepoolId *string, request *ModifyNodePoolNodeConfigRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyNodePoolNodeConfigResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.KubeletConfig)) {
		body["kubelet_config"] = request.KubeletConfig
	}

	if !tea.BoolValue(util.IsUnset(request.RollingPolicy)) {
		body["rolling_policy"] = request.RollingPolicy
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyNodePoolNodeConfig"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId)) + "/node_config"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyNodePoolNodeConfigResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyNodePoolNodeConfig(ClusterId *string, NodepoolId *string, request *ModifyNodePoolNodeConfigRequest) (_result *ModifyNodePoolNodeConfigResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyNodePoolNodeConfigResponse{}
	_body, _err := client.ModifyNodePoolNodeConfigWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ModifyPolicyInstanceWithOptions(clusterId *string, policyName *string, request *ModifyPolicyInstanceRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ModifyPolicyInstanceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Action)) {
		body["action"] = request.Action
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceName)) {
		body["instance_name"] = request.InstanceName
	}

	if !tea.BoolValue(util.IsUnset(request.Namespaces)) {
		body["namespaces"] = request.Namespaces
	}

	if !tea.BoolValue(util.IsUnset(request.Parameters)) {
		body["parameters"] = request.Parameters
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ModifyPolicyInstance"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/policies/" + tea.StringValue(openapiutil.GetEncodeParam(policyName))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ModifyPolicyInstanceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ModifyPolicyInstance(clusterId *string, policyName *string, request *ModifyPolicyInstanceRequest) (_result *ModifyPolicyInstanceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ModifyPolicyInstanceResponse{}
	_body, _err := client.ModifyPolicyInstanceWithOptions(clusterId, policyName, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * You can activate ACK with Alibaba Cloud accounts or RAM users that have the authority of AdministratorAccess.
 *
 * @param request OpenAckServiceRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return OpenAckServiceResponse
 */
func (client *Client) OpenAckServiceWithOptions(request *OpenAckServiceRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *OpenAckServiceResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Type)) {
		query["type"] = request.Type
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("OpenAckService"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/service/open"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &OpenAckServiceResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * You can activate ACK with Alibaba Cloud accounts or RAM users that have the authority of AdministratorAccess.
 *
 * @param request OpenAckServiceRequest
 * @return OpenAckServiceResponse
 */
func (client *Client) OpenAckService(request *OpenAckServiceRequest) (_result *OpenAckServiceResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &OpenAckServiceResponse{}
	_body, _err := client.OpenAckServiceWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PauseClusterUpgradeWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *PauseClusterUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("PauseClusterUpgrade"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/upgrade/pause"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &PauseClusterUpgradeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PauseClusterUpgrade(ClusterId *string) (_result *PauseClusterUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &PauseClusterUpgradeResponse{}
	_body, _err := client.PauseClusterUpgradeWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PauseComponentUpgradeWithOptions(clusterid *string, componentid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *PauseComponentUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("PauseComponentUpgrade"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterid)) + "/components/" + tea.StringValue(openapiutil.GetEncodeParam(componentid)) + "/pause"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &PauseComponentUpgradeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PauseComponentUpgrade(clusterid *string, componentid *string) (_result *PauseComponentUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &PauseComponentUpgradeResponse{}
	_body, _err := client.PauseComponentUpgradeWithOptions(clusterid, componentid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) PauseTaskWithOptions(taskId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *PauseTaskResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("PauseTask"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/tasks/" + tea.StringValue(openapiutil.GetEncodeParam(taskId)) + "/pause"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &PauseTaskResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) PauseTask(taskId *string) (_result *PauseTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &PauseTaskResponse{}
	_body, _err := client.PauseTaskWithOptions(taskId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request RemoveClusterNodesRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return RemoveClusterNodesResponse
 */
// Deprecated
func (client *Client) RemoveClusterNodesWithOptions(ClusterId *string, request *RemoveClusterNodesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *RemoveClusterNodesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DrainNode)) {
		body["drain_node"] = request.DrainNode
	}

	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		body["nodes"] = request.Nodes
	}

	if !tea.BoolValue(util.IsUnset(request.ReleaseNode)) {
		body["release_node"] = request.ReleaseNode
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("RemoveClusterNodes"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodes/remove"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &RemoveClusterNodesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request RemoveClusterNodesRequest
 * @return RemoveClusterNodesResponse
 */
// Deprecated
func (client *Client) RemoveClusterNodes(ClusterId *string, request *RemoveClusterNodesRequest) (_result *RemoveClusterNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &RemoveClusterNodesResponse{}
	_body, _err := client.RemoveClusterNodesWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * >
 * *   When you remove a node, the pods that run on the node are migrated to other nodes. This may cause service interruptions. We recommend that you remove nodes during off-peak hours.
 * *   Unknown errors may occur when you remove nodes. Before you remove nodes, back up the data on the nodes.
 * *   Nodes remain in the Unschedulable state when they are being removed.
 * *   You can remove only worker nodes. You cannot remove control planes.
 *
 * @param tmpReq RemoveNodePoolNodesRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return RemoveNodePoolNodesResponse
 */
func (client *Client) RemoveNodePoolNodesWithOptions(ClusterId *string, NodepoolId *string, tmpReq *RemoveNodePoolNodesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *RemoveNodePoolNodesResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &RemoveNodePoolNodesShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.InstanceIds)) {
		request.InstanceIdsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.InstanceIds, tea.String("instance_ids"), tea.String("json"))
	}

	if !tea.BoolValue(util.IsUnset(tmpReq.Nodes)) {
		request.NodesShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.Nodes, tea.String("nodes"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.DrainNode)) {
		query["drain_node"] = request.DrainNode
	}

	if !tea.BoolValue(util.IsUnset(request.InstanceIdsShrink)) {
		query["instance_ids"] = request.InstanceIdsShrink
	}

	if !tea.BoolValue(util.IsUnset(request.NodesShrink)) {
		query["nodes"] = request.NodesShrink
	}

	if !tea.BoolValue(util.IsUnset(request.ReleaseNode)) {
		query["release_node"] = request.ReleaseNode
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("RemoveNodePoolNodes"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId)) + "/nodes"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &RemoveNodePoolNodesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * >
 * *   When you remove a node, the pods that run on the node are migrated to other nodes. This may cause service interruptions. We recommend that you remove nodes during off-peak hours.
 * *   Unknown errors may occur when you remove nodes. Before you remove nodes, back up the data on the nodes.
 * *   Nodes remain in the Unschedulable state when they are being removed.
 * *   You can remove only worker nodes. You cannot remove control planes.
 *
 * @param request RemoveNodePoolNodesRequest
 * @return RemoveNodePoolNodesResponse
 */
func (client *Client) RemoveNodePoolNodes(ClusterId *string, NodepoolId *string, request *RemoveNodePoolNodesRequest) (_result *RemoveNodePoolNodesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &RemoveNodePoolNodesResponse{}
	_body, _err := client.RemoveNodePoolNodesWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RemoveWorkflowWithOptions(workflowName *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *RemoveWorkflowResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("RemoveWorkflow"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/gs/workflow/" + tea.StringValue(openapiutil.GetEncodeParam(workflowName))),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &RemoveWorkflowResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RemoveWorkflow(workflowName *string) (_result *RemoveWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &RemoveWorkflowResponse{}
	_body, _err := client.RemoveWorkflowWithOptions(workflowName, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) RepairClusterNodePoolWithOptions(clusterId *string, nodepoolId *string, request *RepairClusterNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *RepairClusterNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Nodes)) {
		body["nodes"] = request.Nodes
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("RepairClusterNodePool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(nodepoolId)) + "/repair"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &RepairClusterNodePoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) RepairClusterNodePool(clusterId *string, nodepoolId *string, request *RepairClusterNodePoolRequest) (_result *RepairClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &RepairClusterNodePoolResponse{}
	_body, _err := client.RepairClusterNodePoolWithOptions(clusterId, nodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ResumeComponentUpgradeWithOptions(clusterid *string, componentid *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ResumeComponentUpgradeResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("ResumeComponentUpgrade"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterid)) + "/components/" + tea.StringValue(openapiutil.GetEncodeParam(componentid)) + "/resume"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &ResumeComponentUpgradeResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ResumeComponentUpgrade(clusterid *string, componentid *string) (_result *ResumeComponentUpgradeResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ResumeComponentUpgradeResponse{}
	_body, _err := client.ResumeComponentUpgradeWithOptions(clusterid, componentid, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ResumeTaskWithOptions(taskId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ResumeTaskResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("ResumeTask"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/tasks/" + tea.StringValue(openapiutil.GetEncodeParam(taskId)) + "/resume"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &ResumeTaskResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ResumeTask(taskId *string) (_result *ResumeTaskResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ResumeTaskResponse{}
	_body, _err := client.ResumeTaskWithOptions(taskId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ResumeUpgradeClusterWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ResumeUpgradeClusterResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("ResumeUpgradeCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/upgrade/resume"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &ResumeUpgradeClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ResumeUpgradeCluster(ClusterId *string) (_result *ResumeUpgradeClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ResumeUpgradeClusterResponse{}
	_body, _err := client.ResumeUpgradeClusterWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request ScaleClusterRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return ScaleClusterResponse
 */
// Deprecated
func (client *Client) ScaleClusterWithOptions(ClusterId *string, request *ScaleClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ScaleClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CloudMonitorFlags)) {
		body["cloud_monitor_flags"] = request.CloudMonitorFlags
	}

	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.DisableRollback)) {
		body["disable_rollback"] = request.DisableRollback
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.LoginPassword)) {
		body["login_password"] = request.LoginPassword
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Taints)) {
		body["taints"] = request.Taints
	}

	if !tea.BoolValue(util.IsUnset(request.VswitchIds)) {
		body["vswitch_ids"] = request.VswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenew)) {
		body["worker_auto_renew"] = request.WorkerAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenewPeriod)) {
		body["worker_auto_renew_period"] = request.WorkerAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisk)) {
		body["worker_data_disk"] = request.WorkerDataDisk
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisks)) {
		body["worker_data_disks"] = request.WorkerDataDisks
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceChargeType)) {
		body["worker_instance_charge_type"] = request.WorkerInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceTypes)) {
		body["worker_instance_types"] = request.WorkerInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriod)) {
		body["worker_period"] = request.WorkerPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriodUnit)) {
		body["worker_period_unit"] = request.WorkerPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskCategory)) {
		body["worker_system_disk_category"] = request.WorkerSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSize)) {
		body["worker_system_disk_size"] = request.WorkerSystemDiskSize
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ScaleCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ScaleClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * @deprecated
 *
 * @param request ScaleClusterRequest
 * @return ScaleClusterResponse
 */
// Deprecated
func (client *Client) ScaleCluster(ClusterId *string, request *ScaleClusterRequest) (_result *ScaleClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ScaleClusterResponse{}
	_body, _err := client.ScaleClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ScaleClusterNodePoolWithOptions(ClusterId *string, NodepoolId *string, request *ScaleClusterNodePoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ScaleClusterNodePoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ScaleClusterNodePool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId))),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ScaleClusterNodePoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ScaleClusterNodePool(ClusterId *string, NodepoolId *string, request *ScaleClusterNodePoolRequest) (_result *ScaleClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ScaleClusterNodePoolResponse{}
	_body, _err := client.ScaleClusterNodePoolWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * >  The ScaleOutCluster API operation is phased out. You must call the node pool-related API operations to manage nodes. If you want to add worker nodes to an ACK cluster, call the ScaleClusterNodePool API operation. For more information, see [ScaleClusterNodePool](~~184928~~).
 *
 * @param request ScaleOutClusterRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return ScaleOutClusterResponse
 */
func (client *Client) ScaleOutClusterWithOptions(ClusterId *string, request *ScaleOutClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ScaleOutClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.CloudMonitorFlags)) {
		body["cloud_monitor_flags"] = request.CloudMonitorFlags
	}

	if !tea.BoolValue(util.IsUnset(request.Count)) {
		body["count"] = request.Count
	}

	if !tea.BoolValue(util.IsUnset(request.CpuPolicy)) {
		body["cpu_policy"] = request.CpuPolicy
	}

	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["image_id"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.KeyPair)) {
		body["key_pair"] = request.KeyPair
	}

	if !tea.BoolValue(util.IsUnset(request.LoginPassword)) {
		body["login_password"] = request.LoginPassword
	}

	if !tea.BoolValue(util.IsUnset(request.RdsInstances)) {
		body["rds_instances"] = request.RdsInstances
	}

	if !tea.BoolValue(util.IsUnset(request.Runtime)) {
		body["runtime"] = request.Runtime
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Taints)) {
		body["taints"] = request.Taints
	}

	if !tea.BoolValue(util.IsUnset(request.UserData)) {
		body["user_data"] = request.UserData
	}

	if !tea.BoolValue(util.IsUnset(request.VswitchIds)) {
		body["vswitch_ids"] = request.VswitchIds
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenew)) {
		body["worker_auto_renew"] = request.WorkerAutoRenew
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerAutoRenewPeriod)) {
		body["worker_auto_renew_period"] = request.WorkerAutoRenewPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerDataDisks)) {
		body["worker_data_disks"] = request.WorkerDataDisks
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceChargeType)) {
		body["worker_instance_charge_type"] = request.WorkerInstanceChargeType
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerInstanceTypes)) {
		body["worker_instance_types"] = request.WorkerInstanceTypes
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriod)) {
		body["worker_period"] = request.WorkerPeriod
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerPeriodUnit)) {
		body["worker_period_unit"] = request.WorkerPeriodUnit
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskCategory)) {
		body["worker_system_disk_category"] = request.WorkerSystemDiskCategory
	}

	if !tea.BoolValue(util.IsUnset(request.WorkerSystemDiskSize)) {
		body["worker_system_disk_size"] = request.WorkerSystemDiskSize
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("ScaleOutCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId))),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ScaleOutClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * >  The ScaleOutCluster API operation is phased out. You must call the node pool-related API operations to manage nodes. If you want to add worker nodes to an ACK cluster, call the ScaleClusterNodePool API operation. For more information, see [ScaleClusterNodePool](~~184928~~).
 *
 * @param request ScaleOutClusterRequest
 * @return ScaleOutClusterResponse
 */
func (client *Client) ScaleOutCluster(ClusterId *string, request *ScaleOutClusterRequest) (_result *ScaleOutClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ScaleOutClusterResponse{}
	_body, _err := client.ScaleOutClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) ScanClusterVulsWithOptions(clusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *ScanClusterVulsResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("ScanClusterVuls"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(clusterId)) + "/vuls/scan"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &ScanClusterVulsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) ScanClusterVuls(clusterId *string) (_result *ScanClusterVulsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &ScanClusterVulsResponse{}
	_body, _err := client.ScanClusterVulsWithOptions(clusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) StartAlertWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *StartAlertResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("StartAlert"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/alert/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/alert_rule/start"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &StartAlertResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) StartAlert(ClusterId *string) (_result *StartAlertResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &StartAlertResponse{}
	_body, _err := client.StartAlertWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) StartWorkflowWithOptions(request *StartWorkflowRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *StartWorkflowResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.MappingBamOutFilename)) {
		body["mapping_bam_out_filename"] = request.MappingBamOutFilename
	}

	if !tea.BoolValue(util.IsUnset(request.MappingBamOutPath)) {
		body["mapping_bam_out_path"] = request.MappingBamOutPath
	}

	if !tea.BoolValue(util.IsUnset(request.MappingBucketName)) {
		body["mapping_bucket_name"] = request.MappingBucketName
	}

	if !tea.BoolValue(util.IsUnset(request.MappingFastqFirstFilename)) {
		body["mapping_fastq_first_filename"] = request.MappingFastqFirstFilename
	}

	if !tea.BoolValue(util.IsUnset(request.MappingFastqPath)) {
		body["mapping_fastq_path"] = request.MappingFastqPath
	}

	if !tea.BoolValue(util.IsUnset(request.MappingFastqSecondFilename)) {
		body["mapping_fastq_second_filename"] = request.MappingFastqSecondFilename
	}

	if !tea.BoolValue(util.IsUnset(request.MappingIsMarkDup)) {
		body["mapping_is_mark_dup"] = request.MappingIsMarkDup
	}

	if !tea.BoolValue(util.IsUnset(request.MappingOssRegion)) {
		body["mapping_oss_region"] = request.MappingOssRegion
	}

	if !tea.BoolValue(util.IsUnset(request.MappingReferencePath)) {
		body["mapping_reference_path"] = request.MappingReferencePath
	}

	if !tea.BoolValue(util.IsUnset(request.Service)) {
		body["service"] = request.Service
	}

	if !tea.BoolValue(util.IsUnset(request.WgsBucketName)) {
		body["wgs_bucket_name"] = request.WgsBucketName
	}

	if !tea.BoolValue(util.IsUnset(request.WgsFastqFirstFilename)) {
		body["wgs_fastq_first_filename"] = request.WgsFastqFirstFilename
	}

	if !tea.BoolValue(util.IsUnset(request.WgsFastqPath)) {
		body["wgs_fastq_path"] = request.WgsFastqPath
	}

	if !tea.BoolValue(util.IsUnset(request.WgsFastqSecondFilename)) {
		body["wgs_fastq_second_filename"] = request.WgsFastqSecondFilename
	}

	if !tea.BoolValue(util.IsUnset(request.WgsOssRegion)) {
		body["wgs_oss_region"] = request.WgsOssRegion
	}

	if !tea.BoolValue(util.IsUnset(request.WgsReferencePath)) {
		body["wgs_reference_path"] = request.WgsReferencePath
	}

	if !tea.BoolValue(util.IsUnset(request.WgsVcfOutFilename)) {
		body["wgs_vcf_out_filename"] = request.WgsVcfOutFilename
	}

	if !tea.BoolValue(util.IsUnset(request.WgsVcfOutPath)) {
		body["wgs_vcf_out_path"] = request.WgsVcfOutPath
	}

	if !tea.BoolValue(util.IsUnset(request.WorkflowType)) {
		body["workflow_type"] = request.WorkflowType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("StartWorkflow"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/gs/workflow"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &StartWorkflowResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) StartWorkflow(request *StartWorkflowRequest) (_result *StartWorkflowResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &StartWorkflowResponse{}
	_body, _err := client.StartWorkflowWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) StopAlertWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *StopAlertResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("StopAlert"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/alert/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/alert_rule/stop"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &StopAlertResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) StopAlert(ClusterId *string) (_result *StopAlertResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &StopAlertResponse{}
	_body, _err := client.StopAlertWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) SyncClusterNodePoolWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *SyncClusterNodePoolResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("SyncClusterNodePool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/sync_nodepools"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &SyncClusterNodePoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) SyncClusterNodePool(ClusterId *string) (_result *SyncClusterNodePoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &SyncClusterNodePoolResponse{}
	_body, _err := client.SyncClusterNodePoolWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) TagResourcesWithOptions(request *TagResourcesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *TagResourcesResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		body["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceIds)) {
		body["resource_ids"] = request.ResourceIds
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		body["resource_type"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("TagResources"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/tags"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &TagResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) TagResources(request *TagResourcesRequest) (_result *TagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &TagResourcesResponse{}
	_body, _err := client.TagResourcesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UnInstallClusterAddonsWithOptions(ClusterId *string, request *UnInstallClusterAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UnInstallClusterAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Addons),
	}
	params := &openapi.Params{
		Action:      tea.String("UnInstallClusterAddons"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/components/uninstall"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UnInstallClusterAddonsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UnInstallClusterAddons(ClusterId *string, request *UnInstallClusterAddonsRequest) (_result *UnInstallClusterAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UnInstallClusterAddonsResponse{}
	_body, _err := client.UnInstallClusterAddonsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UntagResourcesWithOptions(tmpReq *UntagResourcesRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UntagResourcesResponse, _err error) {
	_err = util.ValidateModel(tmpReq)
	if _err != nil {
		return _result, _err
	}
	request := &UntagResourcesShrinkRequest{}
	openapiutil.Convert(tmpReq, request)
	if !tea.BoolValue(util.IsUnset(tmpReq.ResourceIds)) {
		request.ResourceIdsShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.ResourceIds, tea.String("resource_ids"), tea.String("json"))
	}

	if !tea.BoolValue(util.IsUnset(tmpReq.TagKeys)) {
		request.TagKeysShrink = openapiutil.ArrayToStringWithSpecifiedStyle(tmpReq.TagKeys, tea.String("tag_keys"), tea.String("json"))
	}

	query := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.All)) {
		query["all"] = request.All
	}

	if !tea.BoolValue(util.IsUnset(request.RegionId)) {
		query["region_id"] = request.RegionId
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceIdsShrink)) {
		query["resource_ids"] = request.ResourceIdsShrink
	}

	if !tea.BoolValue(util.IsUnset(request.ResourceType)) {
		query["resource_type"] = request.ResourceType
	}

	if !tea.BoolValue(util.IsUnset(request.TagKeysShrink)) {
		query["tag_keys"] = request.TagKeysShrink
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Query:   openapiutil.Query(query),
	}
	params := &openapi.Params{
		Action:      tea.String("UntagResources"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/tags"),
		Method:      tea.String("DELETE"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &UntagResourcesResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UntagResources(request *UntagResourcesRequest) (_result *UntagResourcesResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UntagResourcesResponse{}
	_body, _err := client.UntagResourcesWithOptions(request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateContactGroupForAlertWithOptions(ClusterId *string, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpdateContactGroupForAlertResponse, _err error) {
	req := &openapi.OpenApiRequest{
		Headers: headers,
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateContactGroupForAlert"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/alert/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/alert_rule/contact_groups"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UpdateContactGroupForAlertResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateContactGroupForAlert(ClusterId *string) (_result *UpdateContactGroupForAlertResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpdateContactGroupForAlertResponse{}
	_body, _err := client.UpdateContactGroupForAlertWithOptions(ClusterId, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateControlPlaneLogWithOptions(ClusterId *string, request *UpdateControlPlaneLogRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpdateControlPlaneLogResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Aliuid)) {
		body["aliuid"] = request.Aliuid
	}

	if !tea.BoolValue(util.IsUnset(request.Components)) {
		body["components"] = request.Components
	}

	if !tea.BoolValue(util.IsUnset(request.LogProject)) {
		body["log_project"] = request.LogProject
	}

	if !tea.BoolValue(util.IsUnset(request.LogTtl)) {
		body["log_ttl"] = request.LogTtl
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateControlPlaneLog"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/controlplanelog"),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UpdateControlPlaneLogResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateControlPlaneLog(ClusterId *string, request *UpdateControlPlaneLogRequest) (_result *UpdateControlPlaneLogResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpdateControlPlaneLogResponse{}
	_body, _err := client.UpdateControlPlaneLogWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * >
 * *   You can call this operation only with an Alibaba Cloud account.
 * *   After you revoke the kubeconfig file of a cluster, the validity period of the kubeconfig file that you specified becomes invalid. You can call this API operation to specify the validity period again.
 *
 * @param request UpdateK8sClusterUserConfigExpireRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return UpdateK8sClusterUserConfigExpireResponse
 */
func (client *Client) UpdateK8sClusterUserConfigExpireWithOptions(ClusterId *string, request *UpdateK8sClusterUserConfigExpireRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpdateK8sClusterUserConfigExpireResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ExpireHour)) {
		body["expire_hour"] = request.ExpireHour
	}

	if !tea.BoolValue(util.IsUnset(request.User)) {
		body["user"] = request.User
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateK8sClusterUserConfigExpire"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/k8s/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/user_config/expire"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UpdateK8sClusterUserConfigExpireResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * >
 * *   You can call this operation only with an Alibaba Cloud account.
 * *   After you revoke the kubeconfig file of a cluster, the validity period of the kubeconfig file that you specified becomes invalid. You can call this API operation to specify the validity period again.
 *
 * @param request UpdateK8sClusterUserConfigExpireRequest
 * @return UpdateK8sClusterUserConfigExpireResponse
 */
func (client *Client) UpdateK8sClusterUserConfigExpire(ClusterId *string, request *UpdateK8sClusterUserConfigExpireRequest) (_result *UpdateK8sClusterUserConfigExpireResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpdateK8sClusterUserConfigExpireResponse{}
	_body, _err := client.UpdateK8sClusterUserConfigExpireWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpdateTemplateWithOptions(TemplateId *string, request *UpdateTemplateRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpdateTemplateResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.Description)) {
		body["description"] = request.Description
	}

	if !tea.BoolValue(util.IsUnset(request.Name)) {
		body["name"] = request.Name
	}

	if !tea.BoolValue(util.IsUnset(request.Tags)) {
		body["tags"] = request.Tags
	}

	if !tea.BoolValue(util.IsUnset(request.Template)) {
		body["template"] = request.Template
	}

	if !tea.BoolValue(util.IsUnset(request.TemplateType)) {
		body["template_type"] = request.TemplateType
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpdateTemplate"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/templates/" + tea.StringValue(openapiutil.GetEncodeParam(TemplateId))),
		Method:      tea.String("PUT"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UpdateTemplateResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpdateTemplate(TemplateId *string, request *UpdateTemplateRequest) (_result *UpdateTemplateResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpdateTemplateResponse{}
	_body, _err := client.UpdateTemplateWithOptions(TemplateId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpgradeClusterWithOptions(ClusterId *string, request *UpgradeClusterRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpgradeClusterResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ComponentName)) {
		body["component_name"] = request.ComponentName
	}

	if !tea.BoolValue(util.IsUnset(request.MasterOnly)) {
		body["master_only"] = request.MasterOnly
	}

	if !tea.BoolValue(util.IsUnset(request.NextVersion)) {
		body["next_version"] = request.NextVersion
	}

	if !tea.BoolValue(util.IsUnset(request.Version)) {
		body["version"] = request.Version
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpgradeCluster"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/api/v2/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/upgrade"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UpgradeClusterResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpgradeCluster(ClusterId *string, request *UpgradeClusterRequest) (_result *UpgradeClusterResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpgradeClusterResponse{}
	_body, _err := client.UpgradeClusterWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

func (client *Client) UpgradeClusterAddonsWithOptions(ClusterId *string, request *UpgradeClusterAddonsRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpgradeClusterAddonsResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    util.ToArray(request.Body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpgradeClusterAddons"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/components/upgrade"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("none"),
	}
	_result = &UpgradeClusterAddonsResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

func (client *Client) UpgradeClusterAddons(ClusterId *string, request *UpgradeClusterAddonsRequest) (_result *UpgradeClusterAddonsResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpgradeClusterAddonsResponse{}
	_body, _err := client.UpgradeClusterAddonsWithOptions(ClusterId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}

/**
 * You can call the UpgradeClusterNodepool operation to update the Kubernetes version, OS version, or container runtime version of the nodes in a node pool.
 *
 * @param request UpgradeClusterNodepoolRequest
 * @param headers map
 * @param runtime runtime options for this request RuntimeOptions
 * @return UpgradeClusterNodepoolResponse
 */
func (client *Client) UpgradeClusterNodepoolWithOptions(ClusterId *string, NodepoolId *string, request *UpgradeClusterNodepoolRequest, headers map[string]*string, runtime *util.RuntimeOptions) (_result *UpgradeClusterNodepoolResponse, _err error) {
	_err = util.ValidateModel(request)
	if _err != nil {
		return _result, _err
	}
	body := map[string]interface{}{}
	if !tea.BoolValue(util.IsUnset(request.ImageId)) {
		body["image_id"] = request.ImageId
	}

	if !tea.BoolValue(util.IsUnset(request.KubernetesVersion)) {
		body["kubernetes_version"] = request.KubernetesVersion
	}

	if !tea.BoolValue(util.IsUnset(request.RuntimeType)) {
		body["runtime_type"] = request.RuntimeType
	}

	if !tea.BoolValue(util.IsUnset(request.RuntimeVersion)) {
		body["runtime_version"] = request.RuntimeVersion
	}

	req := &openapi.OpenApiRequest{
		Headers: headers,
		Body:    openapiutil.ParseToMap(body),
	}
	params := &openapi.Params{
		Action:      tea.String("UpgradeClusterNodepool"),
		Version:     tea.String("2015-12-15"),
		Protocol:    tea.String("HTTPS"),
		Pathname:    tea.String("/clusters/" + tea.StringValue(openapiutil.GetEncodeParam(ClusterId)) + "/nodepools/" + tea.StringValue(openapiutil.GetEncodeParam(NodepoolId)) + "/upgrade"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("ROA"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = &UpgradeClusterNodepoolResponse{}
	_body, _err := client.CallApi(params, req, runtime)
	if _err != nil {
		return _result, _err
	}
	_err = tea.Convert(_body, &_result)
	return _result, _err
}

/**
 * You can call the UpgradeClusterNodepool operation to update the Kubernetes version, OS version, or container runtime version of the nodes in a node pool.
 *
 * @param request UpgradeClusterNodepoolRequest
 * @return UpgradeClusterNodepoolResponse
 */
func (client *Client) UpgradeClusterNodepool(ClusterId *string, NodepoolId *string, request *UpgradeClusterNodepoolRequest) (_result *UpgradeClusterNodepoolResponse, _err error) {
	runtime := &util.RuntimeOptions{}
	headers := make(map[string]*string)
	_result = &UpgradeClusterNodepoolResponse{}
	_body, _err := client.UpgradeClusterNodepoolWithOptions(ClusterId, NodepoolId, request, headers, runtime)
	if _err != nil {
		return _result, _err
	}
	_result = _body
	return _result, _err
}
