// Package polardb. This file is generated automatically. Please do not modify it manually, thank you!
// Nested attributes require Terraform Plugin Protocol 6.
package polardb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const backupTypeName = "alicloud_polardb_backup"

var (
	_ resource.Resource                = &backupResource{}
	_ resource.ResourceWithImportState = &backupResource{}
)

func NewBackupResource() resource.Resource {
	return &backupResource{}
}

type backupResource struct {
	fwadapt.ResourceBase
}

type backupResourceModel struct {
	Id                                           types.String                     `tfsdk:"id"`
	FwTimeouts                                   types.Object                     `tfsdk:"timeouts"`
	AdvancedDataPolicies                         types.List                       `tfsdk:"advanced_data_policies"`
	BackupEndTime                                types.String                     `tfsdk:"backup_end_time"`
	BackupFrequency                              types.String                     `tfsdk:"backup_frequency"`
	BackupId                                     types.String                     `tfsdk:"backup_id"`
	BackupJobId                                  types.String                     `tfsdk:"backup_job_id"`
	BackupMethod                                 types.String                     `tfsdk:"backup_method"`
	BackupMode                                   types.String                     `tfsdk:"backup_mode"`
	BackupPolicyLevel                            types.String                     `tfsdk:"backup_policy_level"`
	BackupProgressStatus                         types.String                     `tfsdk:"backup_progress_status"`
	BackupRetentionPolicyOnClusterDeletion       types.String                     `tfsdk:"backup_retention_policy_on_cluster_deletion"`
	BackupSetSize                                types.String                     `tfsdk:"backup_set_size"`
	BackupStartTime                              types.String                     `tfsdk:"backup_start_time"`
	BackupType                                   types.String                     `tfsdk:"backup_type"`
	BackupsLevel                                 types.String                     `tfsdk:"backups_level"`
	Comment                                      types.String                     `tfsdk:"comment"`
	ConsistentTime                               types.String                     `tfsdk:"consistent_time"`
	CreateTime                                   types.String                     `tfsdk:"create_time"`
	DBClusterId                                  types.String                     `tfsdk:"db_cluster_id"`
	DataLevel1BackupFrequency                    types.String                     `tfsdk:"data_level1_backup_frequency"`
	DataLevel1BackupPeriod                       types.String                     `tfsdk:"data_level1_backup_period"`
	DataLevel1BackupRetentionPeriod              types.String                     `tfsdk:"data_level1_backup_retention_period"`
	DataLevel1BackupTime                         types.String                     `tfsdk:"data_level1_backup_time"`
	DataLevel2BackupAnotherRegionRegion          types.String                     `tfsdk:"data_level2_backup_another_region_region"`
	DataLevel2BackupAnotherRegionRetentionPeriod types.String                     `tfsdk:"data_level2_backup_another_region_retention_period"`
	DataLevel2BackupPeriod                       types.String                     `tfsdk:"data_level2_backup_period"`
	DataLevel2BackupRetentionPeriod              types.String                     `tfsdk:"data_level2_backup_retention_period"`
	IsAvail                                      types.String                     `tfsdk:"is_avail"`
	JobMode                                      types.String                     `tfsdk:"job_mode"`
	PreferredBackupPeriod                        preferredBackupPeriodStringValue `tfsdk:"preferred_backup_period"`
	PreferredBackupTime                          types.String                     `tfsdk:"preferred_backup_time"`
	PreferredNextBackupTime                      types.String                     `tfsdk:"preferred_next_backup_time"`
	Process                                      types.String                     `tfsdk:"process"`
	Status                                       types.String                     `tfsdk:"status"`
	StoreStatus                                  types.String                     `tfsdk:"store_status"`
	TaskAction                                   types.String                     `tfsdk:"task_action"`
}

func (r *backupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = backupTypeName
}

func (r *backupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": fwTimeoutBlock(),
			"advanced_data_policies": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"action_type": schema.StringAttribute{
							Optional:    true,
							Description: "Action type",
						},
						"only_preserve_one_each_hour": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Only preserve one each hour",
						},
						"storage_class": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Storage class",
						},
						"retention_value": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Retention value",
						},
						"bak_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Backup type",
						},
						"filter_key": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Filter key",
						},
						"filter_value": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Filter value",
						},
						"dest_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Destination type",
						},
						"filter_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Filter type",
						},
						"src_region": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Source region",
						},
						"auto_created": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Auto created",
						},
						"only_preserve_one_each_day": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Only preserve one each day",
						},
						"src_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Source type",
						},
						"dump_action": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Dump action",
						},
						"dest_region": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Destination region",
						},
						"policy_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Policy ID",
						},
						"retention_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Retention type",
						},
					},
				},
				Description: "Advanced data backup policies.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_end_time": schema.StringAttribute{
				Computed:    true,
				Description: "The end time of the query. The end time must be later than the start time. The format is `YYYY-MM-DDThh:mmZ` (UTC time).",
			},
			"backup_frequency": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Backup frequency. Valid values are as follows:  \n\n- **Normal** (default): Standard backup performed once per day at a scheduled time.  \n- **2/24H**: Enhanced backup performed every 2 hours.  \n- **3/24H**: Enhanced backup performed every 3 hours.  \n- **4/24H**: Enhanced backup performed every 4 hours.  \n\n> * After enhanced backup is enabled, all backups completed within any 24-hour period are retained. For backups older than 24 hours, the system retains only the first backup completed after 00:00 UTC each day and deletes all others.  \n> * After enhanced backup is enabled, the **PreferredBackupPeriod** parameter defaults to all days of the week (Monday through Sunday).  ",
			},
			"backup_id": schema.StringAttribute{
				Computed:    true,
				Description: "The backup ID.",
			},
			"backup_job_id": schema.StringAttribute{
				Computed:    true,
				Description: "The backup job ID.",
			},
			"backup_method": schema.StringAttribute{
				Computed:    true,
				Description: "The data backup method. Only snapshot backup is supported. The value is fixed to **Snapshot**.  ",
			},
			"backup_mode": schema.StringAttribute{
				Computed:    true,
				Description: "Backup mode. Valid values are as follows:  \n* **Automated**: System-initiated automatic backup  \n* **Manual**: Manual backup",
			},
			"backup_policy_level": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The backup policy level.",
			},
			"backup_progress_status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the backup task. Valid values:\n* **NoStart**: Not started\n* **Preparing**: Preparing\n* **Waiting**: Waiting\n* **Uploading**: Uploading\n* **Checking**: Checking\n* **Finished**: Finished",
			},
			"backup_retention_policy_on_cluster_deletion": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to retain backups when a cluster is deleted:\n\n- **ALL**: Retains all backups permanently.\n\n- **LATEST**: Retains only the latest backup permanently.\n\n- **NONE**: Does not retain any backup sets.\n\n> Default value: NONE.",
			},
			"backup_set_size": schema.StringAttribute{
				Computed:    true,
				Description: "The size of the backup file, in bytes.\n\n> Because snapshots share data blocks, the amount of space freed after deleting a target snapshot backup is significantly less than the size of the backup file. For more information, see [Backup FAQ](https://help.aliyun.com/document_detail/164881.html).",
			},
			"backup_start_time": schema.StringAttribute{
				Computed:    true,
				Description: "The start time of the query. The format is `YYYY-MM-DDThh:mmZ` (UTC time).",
			},
			"backup_type": schema.StringAttribute{
				Computed:    true,
				Description: "The backup type. Only full backup is supported. The value is fixed to **FullBackup**.  ",
			},
			"backups_level": schema.StringAttribute{
				Computed:    true,
				Description: "Backup level. Valid values are as follows:  \n* **Level-1**: Level-1 backup  \n* **Level-2**: Level-2 backup",
			},
			"comment": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "The backup comment.",
			},
			"consistent_time": schema.StringAttribute{
				Computed:    true,
				Description: "The point in time of the consistent snapshot, formatted as a Unix timestamp in seconds.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Start time of the backup job in UTC. The format is `YYYY-MM-DDThh:mm:ssZ`.  ",
			},
			"db_cluster_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The cluster ID.\n\n> You can call the [DescribeDBClusters](https://help.aliyun.com/document_detail/98094.html) operation to view information about all clusters in the target region, including their cluster IDs.",
			},
			"data_level1_backup_frequency": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Backup frequency. Valid values are as follows:\n\n- **Normal** (default): Standard backup performed once per day at a scheduled time.\n- **2/24H**: High-frequency backup performed every 2 hours.\n- **3/24H**: High-frequency backup performed every 3 hours.\n- **4/24H**: High-frequency backup performed every 4 hours.\n\n> * PolarDB for PostgreSQL (compatible with Oracle) and PolarDB for PostgreSQL do not support this parameter.\n> * If cross-region backup is not supported in the region where your PolarDB for MySQL cluster is deployed, this parameter is not supported. For regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html).\n> * After advanced backup is enabled, we recommend that you do not use this parameter. Instead, use the AdvancedDataPolicies parameter.",
			},
			"data_level1_backup_period": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The primary backup cycle. Valid values are as follows:  \n* **Monday**: Monday  \n* **Tuesday**: Tuesday  \n* **Wednesday**: Wednesday  \n* **Thursday**: Thursday  \n* **Friday**: Friday  \n* **Saturday**: Saturday  \n* **Sunday**: Sunday  \n\n> * * You must select at least two days. Separate multiple values with commas (,).  \n>* * This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL.  \n>* * If the region where your PolarDB for MySQL cluster resides does not support cross-region backup, this parameter is unavailable. For regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html).  \n>* * After advanced backup is enabled, we recommend that you do not use this parameter. Use the **AdvancedDataPolicies** parameter instead.  ",
			},
			"data_level1_backup_retention_period": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Retention period for level-1 backups, in days. Valid values range from 3 to 14.  \n\n> * After advanced backup is enabled, this parameter no longer takes effect. Use the **AdvancedDataPolicies** parameter instead.  ",
			},
			"data_level1_backup_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The time window for automatic backups, in the format `hh:mmZ-hh:mmZ` (UTC). The specified time must be on the hour and span exactly one hour, such as `14:00Z-15:00Z`.\n\n> * This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL.\n> * If cross-region backup is not supported in the region where your PolarDB for MySQL cluster resides, this parameter is not supported. For regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html).",
			},
			"data_level2_backup_another_region_region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region for cross-region Level-2 backup. For regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html).\n\n> After advanced backup is enabled, we recommend that you do not use this parameter. Use the AdvancedDataPolicies parameter instead.",
			},
			"data_level2_backup_another_region_retention_period": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Retention period for cross-region Level-2 backup. Valid values are as follows:\n\n- **0**: Disables the Level-2 backup feature.\n\n- **30–7300**: Retention period for Level-2 backup, in days.\n\n- **-1**: Retains Level-2 backups permanently.\n\n>\n>- When you create a cluster, the default value is **0**, which disables cross-region Level-2 backup.\n>- After advanced backup is enabled, we recommend that you do not use this parameter. Use the AdvancedDataPolicies parameter instead.",
			},
			"data_level2_backup_period": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The secondary backup cycle. Valid values are as follows:  \n* **Monday**: Monday  \n* **Tuesday**: Tuesday  \n* **Wednesday**: Wednesday  \n* **Thursday**: Thursday  \n* **Friday**: Friday  \n* **Saturday**: Saturday  \n* **Sunday**: Sunday  \n\n> * * You must select at least two days. Separate multiple values with commas (,).  \n>* * This parameter is not supported by PolarDB for PostgreSQL (compatible with Oracle) or PolarDB for PostgreSQL.  \n>* * If the region where your PolarDB for MySQL cluster resides does not support cross-region backup, this parameter is unavailable. For regions that support cross-region backup, see [Overview](https://help.aliyun.com/document_detail/72672.html).  \n>* * After advanced backup is enabled, we recommend that you do not use this parameter. Use the **AdvancedDataPolicies** parameter instead.  ",
			},
			"data_level2_backup_retention_period": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Retention period for secondary backups. Valid values are as follows:\n * 0: Disables the secondary backup feature.\n * 30–7300: Retention period for secondary backups, in days.\n * -1: Retains secondary backups indefinitely.\n\n >\n>- * When you create a cluster, the default value is **0**, which disables the secondary backup feature.\n>- * After advanced backup is enabled, we recommend that you do not use this parameter. Instead, use the AdvancedDataPolicies parameter.",
			},
			"is_avail": schema.StringAttribute{
				Computed:    true,
				Description: "Indicates whether the backup set is available. Valid values are as follows:  \n* **0**: Unavailable  \n* **1**: Available.  ",
			},
			"job_mode": schema.StringAttribute{
				Computed:    true,
				Description: "The backup mode. Valid values are as follows:  \n* **Automated**: System-initiated automatic backup  \n* **Manual**: Manual backup.  ",
			},
			"preferred_backup_period": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				CustomType:  preferredBackupPeriodStringType{},
				Description: "Data backup cycle. Valid values are as follows:\n\n- Monday: Monday\n- Tuesday: Tuesday\n- Wednesday: Wednesday\n- Thursday: Thursday\n- Friday: Friday\n- Saturday: Saturday\n- Sunday: Sunday\n> The order of weekdays is not significant; PolarDB normalizes it, and the resource treats the value as an unordered weekday set. After advanced backup is enabled, we recommend that you do not use this parameter. Instead, use the AdvancedDataPolicies parameter.",
			},
			"preferred_backup_time": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Time window for performing automatic backups, in UTC. The format is `hh:mmZ-hh:mmZ`. The specified time must be on the hour and span exactly one hour, such as `14:00Z-15:00Z`.  ",
			},
			"preferred_next_backup_time": schema.StringAttribute{
				Computed:    true,
				Description: "The scheduled time for the next backup, in the format `YYYY-MM-DDThh:mmZ` (UTC).",
			},
			"process": schema.StringAttribute{
				Computed:    true,
				Description: "The progress percentage of the backup task.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The backup status. Valid values are as follows:  \n* **Success**: Backup completed  \n* **Failed**: Backup failed  ",
			},
			"store_status": schema.StringAttribute{
				Computed:    true,
				Description: "Storage status.",
			},
			"task_action": schema.StringAttribute{
				Computed:    true,
				Description: "Backup task type. Valid values are as follows:\n\n* **TempBackupTask**: Temporary backup task  \n* **NormalBackupTask**: Normal backup task",
			},
		},
	}
}

func (r *backupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan backupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"ClientToken": buildClientToken("CreateBackup"),
	}
	// gap#9: RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Client().RegionId
	createValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AdvancedDataPolicies.IsNull() && !plan.AdvancedDataPolicies.IsUnknown() {
		createValues["AdvancedDataPolicies"] = fwValueToInterface(plan.AdvancedDataPolicies, map[string]fwWireField{"action_type": {Name: "ActionType", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_hour": {Name: "OnlyPreserveOneEachHour", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "storage_class": {Name: "StorageClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_value": {Name: "RetentionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "bak_type": {Name: "BakType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_key": {Name: "FilterKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_value": {Name: "FilterValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_type": {Name: "DestType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_type": {Name: "FilterType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_region": {Name: "SrcRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "auto_created": {Name: "AutoCreated", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_day": {Name: "OnlyPreserveOneEachDay", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_type": {Name: "SrcType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dump_action": {Name: "DumpAction", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_region": {Name: "DestRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "policy_id": {Name: "PolicyId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_type": {Name: "RetentionType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.BackupEndTime.IsNull() && !plan.BackupEndTime.IsUnknown() {
		createValues["BackupEndTime"] = fwValueToInterface(plan.BackupEndTime, map[string]fwWireField{}, false)
	}
	if !plan.BackupFrequency.IsNull() && !plan.BackupFrequency.IsUnknown() {
		createValues["BackupFrequency"] = fwValueToInterface(plan.BackupFrequency, map[string]fwWireField{}, false)
	}
	if !plan.BackupId.IsNull() && !plan.BackupId.IsUnknown() {
		createValues["BackupId"] = fwValueToInterface(plan.BackupId, map[string]fwWireField{}, false)
	}
	if !plan.BackupJobId.IsNull() && !plan.BackupJobId.IsUnknown() {
		createValues["BackupJobId"] = fwValueToInterface(plan.BackupJobId, map[string]fwWireField{}, false)
	}
	if !plan.BackupMethod.IsNull() && !plan.BackupMethod.IsUnknown() {
		createValues["BackupMethod"] = fwValueToInterface(plan.BackupMethod, map[string]fwWireField{}, false)
	}
	if !plan.BackupMode.IsNull() && !plan.BackupMode.IsUnknown() {
		createValues["BackupMode"] = fwValueToInterface(plan.BackupMode, map[string]fwWireField{}, false)
	}
	if !plan.BackupPolicyLevel.IsNull() && !plan.BackupPolicyLevel.IsUnknown() {
		createValues["BackupPolicyLevel"] = fwValueToInterface(plan.BackupPolicyLevel, map[string]fwWireField{}, false)
	}
	if !plan.BackupProgressStatus.IsNull() && !plan.BackupProgressStatus.IsUnknown() {
		createValues["BackupProgressStatus"] = fwValueToInterface(plan.BackupProgressStatus, map[string]fwWireField{}, false)
	}
	if !plan.BackupRetentionPolicyOnClusterDeletion.IsNull() && !plan.BackupRetentionPolicyOnClusterDeletion.IsUnknown() {
		createValues["BackupRetentionPolicyOnClusterDeletion"] = fwValueToInterface(plan.BackupRetentionPolicyOnClusterDeletion, map[string]fwWireField{}, false)
	}
	if !plan.BackupSetSize.IsNull() && !plan.BackupSetSize.IsUnknown() {
		createValues["BackupSetSize"] = fwValueToInterface(plan.BackupSetSize, map[string]fwWireField{}, false)
	}
	if !plan.BackupStartTime.IsNull() && !plan.BackupStartTime.IsUnknown() {
		createValues["BackupStartTime"] = fwValueToInterface(plan.BackupStartTime, map[string]fwWireField{}, false)
	}
	if !plan.BackupType.IsNull() && !plan.BackupType.IsUnknown() {
		createValues["BackupType"] = fwValueToInterface(plan.BackupType, map[string]fwWireField{}, false)
	}
	if !plan.BackupsLevel.IsNull() && !plan.BackupsLevel.IsUnknown() {
		createValues["BackupsLevel"] = fwValueToInterface(plan.BackupsLevel, map[string]fwWireField{}, false)
	}
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		createValues["Comment"] = fwValueToInterface(plan.Comment, map[string]fwWireField{}, false)
	}
	if !plan.ConsistentTime.IsNull() && !plan.ConsistentTime.IsUnknown() {
		createValues["ConsistentTime"] = fwValueToInterface(plan.ConsistentTime, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		createValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.DBClusterId.IsNull() && !plan.DBClusterId.IsUnknown() {
		createValues["DBClusterId"] = fwValueToInterface(plan.DBClusterId, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupFrequency.IsNull() && !plan.DataLevel1BackupFrequency.IsUnknown() {
		createValues["DataLevel1BackupFrequency"] = fwValueToInterface(plan.DataLevel1BackupFrequency, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupPeriod.IsNull() && !plan.DataLevel1BackupPeriod.IsUnknown() {
		createValues["DataLevel1BackupPeriod"] = fwValueToInterface(plan.DataLevel1BackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupRetentionPeriod.IsNull() && !plan.DataLevel1BackupRetentionPeriod.IsUnknown() {
		createValues["DataLevel1BackupRetentionPeriod"] = fwValueToInterface(plan.DataLevel1BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupTime.IsNull() && !plan.DataLevel1BackupTime.IsUnknown() {
		createValues["DataLevel1BackupTime"] = fwValueToInterface(plan.DataLevel1BackupTime, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupAnotherRegionRegion.IsNull() && !plan.DataLevel2BackupAnotherRegionRegion.IsUnknown() {
		createValues["DataLevel2BackupAnotherRegionRegion"] = fwValueToInterface(plan.DataLevel2BackupAnotherRegionRegion, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupAnotherRegionRetentionPeriod.IsNull() && !plan.DataLevel2BackupAnotherRegionRetentionPeriod.IsUnknown() {
		createValues["DataLevel2BackupAnotherRegionRetentionPeriod"] = fwValueToInterface(plan.DataLevel2BackupAnotherRegionRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupPeriod.IsNull() && !plan.DataLevel2BackupPeriod.IsUnknown() {
		createValues["DataLevel2BackupPeriod"] = fwValueToInterface(plan.DataLevel2BackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupRetentionPeriod.IsNull() && !plan.DataLevel2BackupRetentionPeriod.IsUnknown() {
		createValues["DataLevel2BackupRetentionPeriod"] = fwValueToInterface(plan.DataLevel2BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.IsAvail.IsNull() && !plan.IsAvail.IsUnknown() {
		createValues["IsAvail"] = fwValueToInterface(plan.IsAvail, map[string]fwWireField{}, false)
	}
	if !plan.JobMode.IsNull() && !plan.JobMode.IsUnknown() {
		createValues["JobMode"] = fwValueToInterface(plan.JobMode, map[string]fwWireField{}, false)
	}
	if !plan.PreferredBackupPeriod.IsNull() && !plan.PreferredBackupPeriod.IsUnknown() {
		createValues["PreferredBackupPeriod"] = fwValueToInterface(plan.PreferredBackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.PreferredBackupTime.IsNull() && !plan.PreferredBackupTime.IsUnknown() {
		createValues["PreferredBackupTime"] = fwValueToInterface(plan.PreferredBackupTime, map[string]fwWireField{}, false)
	}
	if !plan.PreferredNextBackupTime.IsNull() && !plan.PreferredNextBackupTime.IsUnknown() {
		createValues["PreferredNextBackupTime"] = fwValueToInterface(plan.PreferredNextBackupTime, map[string]fwWireField{}, false)
	}
	if !plan.Process.IsNull() && !plan.Process.IsUnknown() {
		createValues["Process"] = fwValueToInterface(plan.Process, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		createValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	if !plan.StoreStatus.IsNull() && !plan.StoreStatus.IsUnknown() {
		createValues["StoreStatus"] = fwValueToInterface(plan.StoreStatus, map[string]fwWireField{}, false)
	}
	if !plan.TaskAction.IsNull() && !plan.TaskAction.IsUnknown() {
		createValues["TaskAction"] = fwValueToInterface(plan.TaskAction, map[string]fwWireField{}, false)
	}
	if err := fwApplyMappings(request, createValues, []fwMapping{
		{Source: "$.Comment", Target: "Comment", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	response, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "CreateBackup", request, operationTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateBackup", backupTypeName),
			errs.WrapErrorf(err, "Creating %s: calling CreateBackup", backupTypeName).Error(),
		)
		return
	}

	identityValues0 := fwWalkMappedPath(request, fwMappingPath("DBClusterId"), nil)
	if len(identityValues0) == 0 || identityValues0[0].Value == nil || fmt.Sprint(identityValues0[0].Value) == "" {
		resp.Diagnostics.AddError("Creating alicloud_polardb_backup: missing identity", "The create operation did not return or accept every identity component")
		return
	}
	dbClusterId := fmt.Sprint(identityValues0[0].Value)
	// F1: CreateBackup returns BackupJobId (the background task id). The backup SET BackupId
	// (keyed by DescribeBackups/DeleteBackup) materializes asynchronously and differs from
	// BackupJobId; it is resolved after the task finishes and paired with DBClusterId to form
	// the composite identity "DBClusterId:BackupId" so Read/Delete target the exact owned
	// backup instead of "last wins" over every backup set in the cluster.
	backupJobId := ""
	if values := fwWalkMappedPath(response, fwMappingPath("BackupJobId"), nil); len(values) > 0 && values[0].Value != nil {
		backupJobId = fmt.Sprint(values[0].Value)
	}
	if backupJobId == "" || backupJobId == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateBackup", backupTypeName),
			fmt.Sprintf("CreateBackup did not return a BackupJobId: %v", response),
		)
		return
	}
	// F2: CreateBackup only submits the background task and returns BackupJobId; the backup
	// SET BackupId (keyed by DescribeBackups/DeleteBackup) materializes asynchronously and is
	// NOT the BackupJobId. Wait for the task to reach Finished, then resolve the real backup-set
	// BackupId before composing the composite identity. Otherwise Read/Delete would query
	// DescribeBackups with the BackupJobId (a task id) and never match the owned backup set.
	if waitErr := backupWaitForBackupTask(r.Client(), dbClusterId, backupJobId, operationTimeout); waitErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for backup task", backupTypeName),
			errs.WrapErrorf(waitErr, "Creating %s: waiting for backup task %s on cluster %s", backupTypeName, backupJobId, dbClusterId).Error(),
		)
		return
	}
	comment := ""
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		comment = plan.Comment.ValueString()
	}
	// The comment is create-only and not reliably read back, so the post-create
	// state preserves the planned value. Pin it to the value actually sent here:
	// an omitted (unknown) comment would otherwise stay unknown after apply and
	// Terraform rejects the result object.
	plan.Comment = types.StringValue(comment)
	backupId, resolveErr := backupResolveBackupId(r.Client(), dbClusterId, backupJobId, comment)
	if resolveErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: resolving backup set id", backupTypeName),
			errs.WrapErrorf(resolveErr, "Creating %s: resolving backup set id for task %s on cluster %s", backupTypeName, backupJobId, dbClusterId).Error(),
		)
		return
	}
	if backupId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: resolving backup set id", backupTypeName),
			fmt.Sprintf("CreateBackup task %s on cluster %s finished but no matching backup set was found", backupJobId, dbClusterId),
		)
		return
	}
	id := backupComposeID(dbClusterId, backupId)

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.BackupId = types.StringValue(backupId)
	trackedPlan.BackupJobId = types.StringValue(backupJobId)
	trackedPlan.AdvancedDataPolicies = fwKnownValue(trackedPlan.AdvancedDataPolicies).(types.List)
	trackedPlan.BackupEndTime = fwKnownValue(trackedPlan.BackupEndTime).(types.String)
	trackedPlan.BackupFrequency = fwKnownValue(trackedPlan.BackupFrequency).(types.String)
	trackedPlan.BackupId = fwKnownValue(trackedPlan.BackupId).(types.String)
	trackedPlan.BackupJobId = fwKnownValue(trackedPlan.BackupJobId).(types.String)
	trackedPlan.BackupMethod = fwKnownValue(trackedPlan.BackupMethod).(types.String)
	trackedPlan.BackupMode = fwKnownValue(trackedPlan.BackupMode).(types.String)
	trackedPlan.BackupPolicyLevel = fwKnownValue(trackedPlan.BackupPolicyLevel).(types.String)
	trackedPlan.BackupProgressStatus = fwKnownValue(trackedPlan.BackupProgressStatus).(types.String)
	trackedPlan.BackupRetentionPolicyOnClusterDeletion = fwKnownValue(trackedPlan.BackupRetentionPolicyOnClusterDeletion).(types.String)
	trackedPlan.BackupSetSize = fwKnownValue(trackedPlan.BackupSetSize).(types.String)
	trackedPlan.BackupStartTime = fwKnownValue(trackedPlan.BackupStartTime).(types.String)
	trackedPlan.BackupType = fwKnownValue(trackedPlan.BackupType).(types.String)
	trackedPlan.BackupsLevel = fwKnownValue(trackedPlan.BackupsLevel).(types.String)
	trackedPlan.Comment = fwKnownValue(trackedPlan.Comment).(types.String)
	trackedPlan.ConsistentTime = fwKnownValue(trackedPlan.ConsistentTime).(types.String)
	trackedPlan.CreateTime = fwKnownValue(trackedPlan.CreateTime).(types.String)
	trackedPlan.DBClusterId = fwKnownValue(trackedPlan.DBClusterId).(types.String)
	trackedPlan.DataLevel1BackupFrequency = fwKnownValue(trackedPlan.DataLevel1BackupFrequency).(types.String)
	trackedPlan.DataLevel1BackupPeriod = fwKnownValue(trackedPlan.DataLevel1BackupPeriod).(types.String)
	trackedPlan.DataLevel1BackupRetentionPeriod = fwKnownValue(trackedPlan.DataLevel1BackupRetentionPeriod).(types.String)
	trackedPlan.DataLevel1BackupTime = fwKnownValue(trackedPlan.DataLevel1BackupTime).(types.String)
	trackedPlan.DataLevel2BackupAnotherRegionRegion = fwKnownValue(trackedPlan.DataLevel2BackupAnotherRegionRegion).(types.String)
	trackedPlan.DataLevel2BackupAnotherRegionRetentionPeriod = fwKnownValue(trackedPlan.DataLevel2BackupAnotherRegionRetentionPeriod).(types.String)
	trackedPlan.DataLevel2BackupPeriod = fwKnownValue(trackedPlan.DataLevel2BackupPeriod).(types.String)
	trackedPlan.DataLevel2BackupRetentionPeriod = fwKnownValue(trackedPlan.DataLevel2BackupRetentionPeriod).(types.String)
	trackedPlan.IsAvail = fwKnownValue(trackedPlan.IsAvail).(types.String)
	trackedPlan.JobMode = fwKnownValue(trackedPlan.JobMode).(types.String)
	trackedPlan.PreferredBackupPeriod = fwKnownValue(trackedPlan.PreferredBackupPeriod).(preferredBackupPeriodStringValue)
	trackedPlan.PreferredBackupTime = fwKnownValue(trackedPlan.PreferredBackupTime).(types.String)
	trackedPlan.PreferredNextBackupTime = fwKnownValue(trackedPlan.PreferredNextBackupTime).(types.String)
	trackedPlan.Process = fwKnownValue(trackedPlan.Process).(types.String)
	trackedPlan.Status = fwKnownValue(trackedPlan.Status).(types.String)
	trackedPlan.StoreStatus = fwKnownValue(trackedPlan.StoreStatus).(types.String)
	trackedPlan.TaskAction = fwKnownValue(trackedPlan.TaskAction).(types.String)
	resp.Diagnostics.Append(resp.State.Set(ctx, trackedPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := backupGet(r.Client(), id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", backupTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", backupTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: after create", backupTypeName),
			fmt.Sprintf("The resource %s is not visible via DescribeBackupPolicy right after creation", id),
		)
		return
	}
	// Apply configurable properties represented only by secondary update APIs.
	state := backupState(object, &plan)
	state.Id = types.StringValue(id)
	plan.Id = types.StringValue(id)
	planValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AdvancedDataPolicies.IsNull() && !plan.AdvancedDataPolicies.IsUnknown() {
		planValues["AdvancedDataPolicies"] = fwValueToInterface(plan.AdvancedDataPolicies, map[string]fwWireField{"action_type": {Name: "ActionType", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_hour": {Name: "OnlyPreserveOneEachHour", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "storage_class": {Name: "StorageClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_value": {Name: "RetentionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "bak_type": {Name: "BakType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_key": {Name: "FilterKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_value": {Name: "FilterValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_type": {Name: "DestType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_type": {Name: "FilterType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_region": {Name: "SrcRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "auto_created": {Name: "AutoCreated", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_day": {Name: "OnlyPreserveOneEachDay", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_type": {Name: "SrcType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dump_action": {Name: "DumpAction", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_region": {Name: "DestRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "policy_id": {Name: "PolicyId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_type": {Name: "RetentionType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.BackupEndTime.IsNull() && !plan.BackupEndTime.IsUnknown() {
		planValues["BackupEndTime"] = fwValueToInterface(plan.BackupEndTime, map[string]fwWireField{}, false)
	}
	if !plan.BackupFrequency.IsNull() && !plan.BackupFrequency.IsUnknown() {
		planValues["BackupFrequency"] = fwValueToInterface(plan.BackupFrequency, map[string]fwWireField{}, false)
	}
	if !plan.BackupId.IsNull() && !plan.BackupId.IsUnknown() {
		planValues["BackupId"] = fwValueToInterface(plan.BackupId, map[string]fwWireField{}, false)
	}
	if !plan.BackupJobId.IsNull() && !plan.BackupJobId.IsUnknown() {
		planValues["BackupJobId"] = fwValueToInterface(plan.BackupJobId, map[string]fwWireField{}, false)
	}
	if !plan.BackupMethod.IsNull() && !plan.BackupMethod.IsUnknown() {
		planValues["BackupMethod"] = fwValueToInterface(plan.BackupMethod, map[string]fwWireField{}, false)
	}
	if !plan.BackupMode.IsNull() && !plan.BackupMode.IsUnknown() {
		planValues["BackupMode"] = fwValueToInterface(plan.BackupMode, map[string]fwWireField{}, false)
	}
	if !plan.BackupPolicyLevel.IsNull() && !plan.BackupPolicyLevel.IsUnknown() {
		planValues["BackupPolicyLevel"] = fwValueToInterface(plan.BackupPolicyLevel, map[string]fwWireField{}, false)
	}
	if !plan.BackupProgressStatus.IsNull() && !plan.BackupProgressStatus.IsUnknown() {
		planValues["BackupProgressStatus"] = fwValueToInterface(plan.BackupProgressStatus, map[string]fwWireField{}, false)
	}
	if !plan.BackupRetentionPolicyOnClusterDeletion.IsNull() && !plan.BackupRetentionPolicyOnClusterDeletion.IsUnknown() {
		planValues["BackupRetentionPolicyOnClusterDeletion"] = fwValueToInterface(plan.BackupRetentionPolicyOnClusterDeletion, map[string]fwWireField{}, false)
	}
	if !plan.BackupSetSize.IsNull() && !plan.BackupSetSize.IsUnknown() {
		planValues["BackupSetSize"] = fwValueToInterface(plan.BackupSetSize, map[string]fwWireField{}, false)
	}
	if !plan.BackupStartTime.IsNull() && !plan.BackupStartTime.IsUnknown() {
		planValues["BackupStartTime"] = fwValueToInterface(plan.BackupStartTime, map[string]fwWireField{}, false)
	}
	if !plan.BackupType.IsNull() && !plan.BackupType.IsUnknown() {
		planValues["BackupType"] = fwValueToInterface(plan.BackupType, map[string]fwWireField{}, false)
	}
	if !plan.BackupsLevel.IsNull() && !plan.BackupsLevel.IsUnknown() {
		planValues["BackupsLevel"] = fwValueToInterface(plan.BackupsLevel, map[string]fwWireField{}, false)
	}
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		planValues["Comment"] = fwValueToInterface(plan.Comment, map[string]fwWireField{}, false)
	}
	if !plan.ConsistentTime.IsNull() && !plan.ConsistentTime.IsUnknown() {
		planValues["ConsistentTime"] = fwValueToInterface(plan.ConsistentTime, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		planValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.DBClusterId.IsNull() && !plan.DBClusterId.IsUnknown() {
		planValues["DBClusterId"] = fwValueToInterface(plan.DBClusterId, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupFrequency.IsNull() && !plan.DataLevel1BackupFrequency.IsUnknown() {
		planValues["DataLevel1BackupFrequency"] = fwValueToInterface(plan.DataLevel1BackupFrequency, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupPeriod.IsNull() && !plan.DataLevel1BackupPeriod.IsUnknown() {
		planValues["DataLevel1BackupPeriod"] = fwValueToInterface(plan.DataLevel1BackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupRetentionPeriod.IsNull() && !plan.DataLevel1BackupRetentionPeriod.IsUnknown() {
		planValues["DataLevel1BackupRetentionPeriod"] = fwValueToInterface(plan.DataLevel1BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupTime.IsNull() && !plan.DataLevel1BackupTime.IsUnknown() {
		planValues["DataLevel1BackupTime"] = fwValueToInterface(plan.DataLevel1BackupTime, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupAnotherRegionRegion.IsNull() && !plan.DataLevel2BackupAnotherRegionRegion.IsUnknown() {
		planValues["DataLevel2BackupAnotherRegionRegion"] = fwValueToInterface(plan.DataLevel2BackupAnotherRegionRegion, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupAnotherRegionRetentionPeriod.IsNull() && !plan.DataLevel2BackupAnotherRegionRetentionPeriod.IsUnknown() {
		planValues["DataLevel2BackupAnotherRegionRetentionPeriod"] = fwValueToInterface(plan.DataLevel2BackupAnotherRegionRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupPeriod.IsNull() && !plan.DataLevel2BackupPeriod.IsUnknown() {
		planValues["DataLevel2BackupPeriod"] = fwValueToInterface(plan.DataLevel2BackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupRetentionPeriod.IsNull() && !plan.DataLevel2BackupRetentionPeriod.IsUnknown() {
		planValues["DataLevel2BackupRetentionPeriod"] = fwValueToInterface(plan.DataLevel2BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.IsAvail.IsNull() && !plan.IsAvail.IsUnknown() {
		planValues["IsAvail"] = fwValueToInterface(plan.IsAvail, map[string]fwWireField{}, false)
	}
	if !plan.JobMode.IsNull() && !plan.JobMode.IsUnknown() {
		planValues["JobMode"] = fwValueToInterface(plan.JobMode, map[string]fwWireField{}, false)
	}
	if !plan.PreferredBackupPeriod.IsNull() && !plan.PreferredBackupPeriod.IsUnknown() {
		planValues["PreferredBackupPeriod"] = fwValueToInterface(plan.PreferredBackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.PreferredBackupTime.IsNull() && !plan.PreferredBackupTime.IsUnknown() {
		planValues["PreferredBackupTime"] = fwValueToInterface(plan.PreferredBackupTime, map[string]fwWireField{}, false)
	}
	if !plan.PreferredNextBackupTime.IsNull() && !plan.PreferredNextBackupTime.IsUnknown() {
		planValues["PreferredNextBackupTime"] = fwValueToInterface(plan.PreferredNextBackupTime, map[string]fwWireField{}, false)
	}
	if !plan.Process.IsNull() && !plan.Process.IsUnknown() {
		planValues["Process"] = fwValueToInterface(plan.Process, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		planValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	if !plan.StoreStatus.IsNull() && !plan.StoreStatus.IsUnknown() {
		planValues["StoreStatus"] = fwValueToInterface(plan.StoreStatus, map[string]fwWireField{}, false)
	}
	if !plan.TaskAction.IsNull() && !plan.TaskAction.IsUnknown() {
		planValues["TaskAction"] = fwValueToInterface(plan.TaskAction, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.AdvancedDataPolicies.IsNull() && !state.AdvancedDataPolicies.IsUnknown() {
		stateValues["AdvancedDataPolicies"] = fwValueToInterface(state.AdvancedDataPolicies, map[string]fwWireField{"action_type": {Name: "ActionType", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_hour": {Name: "OnlyPreserveOneEachHour", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "storage_class": {Name: "StorageClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_value": {Name: "RetentionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "bak_type": {Name: "BakType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_key": {Name: "FilterKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_value": {Name: "FilterValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_type": {Name: "DestType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_type": {Name: "FilterType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_region": {Name: "SrcRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "auto_created": {Name: "AutoCreated", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_day": {Name: "OnlyPreserveOneEachDay", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_type": {Name: "SrcType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dump_action": {Name: "DumpAction", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_region": {Name: "DestRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "policy_id": {Name: "PolicyId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_type": {Name: "RetentionType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.BackupEndTime.IsNull() && !state.BackupEndTime.IsUnknown() {
		stateValues["BackupEndTime"] = fwValueToInterface(state.BackupEndTime, map[string]fwWireField{}, false)
	}
	if !state.BackupFrequency.IsNull() && !state.BackupFrequency.IsUnknown() {
		stateValues["BackupFrequency"] = fwValueToInterface(state.BackupFrequency, map[string]fwWireField{}, false)
	}
	if !state.BackupId.IsNull() && !state.BackupId.IsUnknown() {
		stateValues["BackupId"] = fwValueToInterface(state.BackupId, map[string]fwWireField{}, false)
	}
	if !state.BackupJobId.IsNull() && !state.BackupJobId.IsUnknown() {
		stateValues["BackupJobId"] = fwValueToInterface(state.BackupJobId, map[string]fwWireField{}, false)
	}
	if !state.BackupMethod.IsNull() && !state.BackupMethod.IsUnknown() {
		stateValues["BackupMethod"] = fwValueToInterface(state.BackupMethod, map[string]fwWireField{}, false)
	}
	if !state.BackupMode.IsNull() && !state.BackupMode.IsUnknown() {
		stateValues["BackupMode"] = fwValueToInterface(state.BackupMode, map[string]fwWireField{}, false)
	}
	if !state.BackupPolicyLevel.IsNull() && !state.BackupPolicyLevel.IsUnknown() {
		stateValues["BackupPolicyLevel"] = fwValueToInterface(state.BackupPolicyLevel, map[string]fwWireField{}, false)
	}
	if !state.BackupProgressStatus.IsNull() && !state.BackupProgressStatus.IsUnknown() {
		stateValues["BackupProgressStatus"] = fwValueToInterface(state.BackupProgressStatus, map[string]fwWireField{}, false)
	}
	if !state.BackupRetentionPolicyOnClusterDeletion.IsNull() && !state.BackupRetentionPolicyOnClusterDeletion.IsUnknown() {
		stateValues["BackupRetentionPolicyOnClusterDeletion"] = fwValueToInterface(state.BackupRetentionPolicyOnClusterDeletion, map[string]fwWireField{}, false)
	}
	if !state.BackupSetSize.IsNull() && !state.BackupSetSize.IsUnknown() {
		stateValues["BackupSetSize"] = fwValueToInterface(state.BackupSetSize, map[string]fwWireField{}, false)
	}
	if !state.BackupStartTime.IsNull() && !state.BackupStartTime.IsUnknown() {
		stateValues["BackupStartTime"] = fwValueToInterface(state.BackupStartTime, map[string]fwWireField{}, false)
	}
	if !state.BackupType.IsNull() && !state.BackupType.IsUnknown() {
		stateValues["BackupType"] = fwValueToInterface(state.BackupType, map[string]fwWireField{}, false)
	}
	if !state.BackupsLevel.IsNull() && !state.BackupsLevel.IsUnknown() {
		stateValues["BackupsLevel"] = fwValueToInterface(state.BackupsLevel, map[string]fwWireField{}, false)
	}
	if !state.Comment.IsNull() && !state.Comment.IsUnknown() {
		stateValues["Comment"] = fwValueToInterface(state.Comment, map[string]fwWireField{}, false)
	}
	if !state.ConsistentTime.IsNull() && !state.ConsistentTime.IsUnknown() {
		stateValues["ConsistentTime"] = fwValueToInterface(state.ConsistentTime, map[string]fwWireField{}, false)
	}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		stateValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.DBClusterId.IsNull() && !state.DBClusterId.IsUnknown() {
		stateValues["DBClusterId"] = fwValueToInterface(state.DBClusterId, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupFrequency.IsNull() && !state.DataLevel1BackupFrequency.IsUnknown() {
		stateValues["DataLevel1BackupFrequency"] = fwValueToInterface(state.DataLevel1BackupFrequency, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupPeriod.IsNull() && !state.DataLevel1BackupPeriod.IsUnknown() {
		stateValues["DataLevel1BackupPeriod"] = fwValueToInterface(state.DataLevel1BackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupRetentionPeriod.IsNull() && !state.DataLevel1BackupRetentionPeriod.IsUnknown() {
		stateValues["DataLevel1BackupRetentionPeriod"] = fwValueToInterface(state.DataLevel1BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupTime.IsNull() && !state.DataLevel1BackupTime.IsUnknown() {
		stateValues["DataLevel1BackupTime"] = fwValueToInterface(state.DataLevel1BackupTime, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupAnotherRegionRegion.IsNull() && !state.DataLevel2BackupAnotherRegionRegion.IsUnknown() {
		stateValues["DataLevel2BackupAnotherRegionRegion"] = fwValueToInterface(state.DataLevel2BackupAnotherRegionRegion, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupAnotherRegionRetentionPeriod.IsNull() && !state.DataLevel2BackupAnotherRegionRetentionPeriod.IsUnknown() {
		stateValues["DataLevel2BackupAnotherRegionRetentionPeriod"] = fwValueToInterface(state.DataLevel2BackupAnotherRegionRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupPeriod.IsNull() && !state.DataLevel2BackupPeriod.IsUnknown() {
		stateValues["DataLevel2BackupPeriod"] = fwValueToInterface(state.DataLevel2BackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupRetentionPeriod.IsNull() && !state.DataLevel2BackupRetentionPeriod.IsUnknown() {
		stateValues["DataLevel2BackupRetentionPeriod"] = fwValueToInterface(state.DataLevel2BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.IsAvail.IsNull() && !state.IsAvail.IsUnknown() {
		stateValues["IsAvail"] = fwValueToInterface(state.IsAvail, map[string]fwWireField{}, false)
	}
	if !state.JobMode.IsNull() && !state.JobMode.IsUnknown() {
		stateValues["JobMode"] = fwValueToInterface(state.JobMode, map[string]fwWireField{}, false)
	}
	if !state.PreferredBackupPeriod.IsNull() && !state.PreferredBackupPeriod.IsUnknown() {
		stateValues["PreferredBackupPeriod"] = fwValueToInterface(state.PreferredBackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.PreferredBackupTime.IsNull() && !state.PreferredBackupTime.IsUnknown() {
		stateValues["PreferredBackupTime"] = fwValueToInterface(state.PreferredBackupTime, map[string]fwWireField{}, false)
	}
	if !state.PreferredNextBackupTime.IsNull() && !state.PreferredNextBackupTime.IsUnknown() {
		stateValues["PreferredNextBackupTime"] = fwValueToInterface(state.PreferredNextBackupTime, map[string]fwWireField{}, false)
	}
	if !state.Process.IsNull() && !state.Process.IsUnknown() {
		stateValues["Process"] = fwValueToInterface(state.Process, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		stateValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	if !state.StoreStatus.IsNull() && !state.StoreStatus.IsUnknown() {
		stateValues["StoreStatus"] = fwValueToInterface(state.StoreStatus, map[string]fwWireField{}, false)
	}
	if !state.TaskAction.IsNull() && !state.TaskAction.IsUnknown() {
		stateValues["TaskAction"] = fwValueToInterface(state.TaskAction, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.BackupFrequency, state.BackupFrequency) {
		updated0 = true
	}
	if fwPlannedChanged(plan.BackupPolicyLevel, state.BackupPolicyLevel) {
		updated0 = true
	}
	if fwPlannedChanged(plan.BackupRetentionPolicyOnClusterDeletion, state.BackupRetentionPolicyOnClusterDeletion) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupFrequency, state.DataLevel1BackupFrequency) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupPeriod, state.DataLevel1BackupPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupRetentionPeriod, state.DataLevel1BackupRetentionPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupTime, state.DataLevel1BackupTime) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupAnotherRegionRegion, state.DataLevel2BackupAnotherRegionRegion) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupAnotherRegionRetentionPeriod, state.DataLevel2BackupAnotherRegionRetentionPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupPeriod, state.DataLevel2BackupPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupRetentionPeriod, state.DataLevel2BackupRetentionPeriod) {
		updated0 = true
	}
	// F7: preferred_backup_period is an unordered weekday set; an API-reordered value
	// (e.g. Wednesday,Monday -> Monday,Wednesday) must not trigger a policy update, while a
	// genuine add/remove of weekdays still must.
	if !backupWeekdaySetEqual(plan.PreferredBackupPeriod.ValueString(), state.PreferredBackupPeriod.ValueString()) {
		updated0 = true
	}
	if fwPlannedChanged(plan.PreferredBackupTime, state.PreferredBackupTime) {
		updated0 = true
	}
	// F4: detect advanced policy changes so ModifyBackupPolicy is called even when only
	// AdvancedDataPolicies differ and all scalar policy fields are unchanged.
	if updated0 || fwPlannedChanged(plan.AdvancedDataPolicies, state.AdvancedDataPolicies) {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"DBClusterId": dbClusterId,
			"ClientToken": buildClientToken("ModifyBackupPolicy"),
		}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, operationValues0, []fwMapping{
			{Source: "$.AdvancedDataPolicies[*].ActionType", Target: "AdvancedDataPolicies[*].ActionType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].AutoCreated", Target: "AdvancedDataPolicies[*].AutoCreated", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].BakType", Target: "AdvancedDataPolicies[*].BakType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].DestRegion", Target: "AdvancedDataPolicies[*].DestRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].DestType", Target: "AdvancedDataPolicies[*].DestType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].DumpAction", Target: "AdvancedDataPolicies[*].DumpAction", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].FilterKey", Target: "AdvancedDataPolicies[*].FilterKey", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].FilterType", Target: "AdvancedDataPolicies[*].FilterType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].FilterValue", Target: "AdvancedDataPolicies[*].FilterValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachDay", Target: "AdvancedDataPolicies[*].OnlyPreserveOneEachDay", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachHour", Target: "AdvancedDataPolicies[*].OnlyPreserveOneEachHour", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].PolicyId", Target: "AdvancedDataPolicies[*].PolicyId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].RetentionType", Target: "AdvancedDataPolicies[*].RetentionType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].RetentionValue", Target: "AdvancedDataPolicies[*].RetentionValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].SrcRegion", Target: "AdvancedDataPolicies[*].SrcRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].SrcType", Target: "AdvancedDataPolicies[*].SrcType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].StorageClass", Target: "AdvancedDataPolicies[*].StorageClass", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupFrequency", Target: "BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupPolicyLevel", Target: "BackupPolicyLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupRetentionPolicyOnClusterDeletion", Target: "BackupRetentionPolicyOnClusterDeletion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupFrequency", Target: "DataLevel1BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupPeriod", Target: "DataLevel1BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupRetentionPeriod", Target: "DataLevel1BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupTime", Target: "DataLevel1BackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupAnotherRegionRegion", Target: "DataLevel2BackupAnotherRegionRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupAnotherRegionRetentionPeriod", Target: "DataLevel2BackupAnotherRegionRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupPeriod", Target: "DataLevel2BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupRetentionPeriod", Target: "DataLevel2BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PreferredBackupPeriod", Target: "PreferredBackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PreferredBackupTime", Target: "PreferredBackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building ModifyBackupPolicy request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "ModifyBackupPolicy", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_polardb_backup: ModifyBackupPolicy", err.Error())
			return
		}
	}
	object, err = backupGet(r.Client(), id)
	if err != nil || object == nil {
		resp.Diagnostics.AddError("Reading after post-create updates", fmt.Sprint(err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, backupState(object, &plan))...)
}

func (r *backupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state backupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := backupGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeBackupPolicy", backupTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeBackupPolicy", backupTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, backupState(object, &state))...)
}

func (r *backupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State = req.State
	var plan, state backupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "update", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	// comment is create-only (CreateBackup) and cannot be changed here, so an omitted
	// (unknown) value must be pinned to the prior state instead of being written back
	// unknown after apply.
	if plan.Comment.IsUnknown() {
		plan.Comment = state.Comment
	}
	// An advanced policy entry is validated as a complete chain and is identified by its
	// computed policy_id. The framework plans those computed sub-attributes as unknown, so
	// fill them from the prior state before the change-filtered request is built; otherwise
	// only the differing leaves are sent (e.g. {"ActionType":"UPDATE","RetentionValue":"14"})
	// and DBS rejects the request with DBS.NotExisted "the policy which you want to update is
	// not exist" (request id 01A166E4-9CB3-5BB4-9D05-4A6B60818C6E).
	plan.AdvancedDataPolicies = backupAlignAdvancedPolicies(plan.AdvancedDataPolicies, state.AdvancedDataPolicies)
	advancedPoliciesChanged := fwPlannedChanged(plan.AdvancedDataPolicies, state.AdvancedDataPolicies)

	planValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !plan.AdvancedDataPolicies.IsNull() && !plan.AdvancedDataPolicies.IsUnknown() {
		planValues["AdvancedDataPolicies"] = fwValueToInterface(plan.AdvancedDataPolicies, map[string]fwWireField{"action_type": {Name: "ActionType", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_hour": {Name: "OnlyPreserveOneEachHour", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "storage_class": {Name: "StorageClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_value": {Name: "RetentionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "bak_type": {Name: "BakType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_key": {Name: "FilterKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_value": {Name: "FilterValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_type": {Name: "DestType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_type": {Name: "FilterType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_region": {Name: "SrcRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "auto_created": {Name: "AutoCreated", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_day": {Name: "OnlyPreserveOneEachDay", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_type": {Name: "SrcType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dump_action": {Name: "DumpAction", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_region": {Name: "DestRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "policy_id": {Name: "PolicyId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_type": {Name: "RetentionType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !plan.BackupEndTime.IsNull() && !plan.BackupEndTime.IsUnknown() {
		planValues["BackupEndTime"] = fwValueToInterface(plan.BackupEndTime, map[string]fwWireField{}, false)
	}
	if !plan.BackupFrequency.IsNull() && !plan.BackupFrequency.IsUnknown() {
		planValues["BackupFrequency"] = fwValueToInterface(plan.BackupFrequency, map[string]fwWireField{}, false)
	}
	if !plan.BackupId.IsNull() && !plan.BackupId.IsUnknown() {
		planValues["BackupId"] = fwValueToInterface(plan.BackupId, map[string]fwWireField{}, false)
	}
	if !plan.BackupJobId.IsNull() && !plan.BackupJobId.IsUnknown() {
		planValues["BackupJobId"] = fwValueToInterface(plan.BackupJobId, map[string]fwWireField{}, false)
	}
	if !plan.BackupMethod.IsNull() && !plan.BackupMethod.IsUnknown() {
		planValues["BackupMethod"] = fwValueToInterface(plan.BackupMethod, map[string]fwWireField{}, false)
	}
	if !plan.BackupMode.IsNull() && !plan.BackupMode.IsUnknown() {
		planValues["BackupMode"] = fwValueToInterface(plan.BackupMode, map[string]fwWireField{}, false)
	}
	if !plan.BackupPolicyLevel.IsNull() && !plan.BackupPolicyLevel.IsUnknown() {
		planValues["BackupPolicyLevel"] = fwValueToInterface(plan.BackupPolicyLevel, map[string]fwWireField{}, false)
	}
	if !plan.BackupProgressStatus.IsNull() && !plan.BackupProgressStatus.IsUnknown() {
		planValues["BackupProgressStatus"] = fwValueToInterface(plan.BackupProgressStatus, map[string]fwWireField{}, false)
	}
	if !plan.BackupRetentionPolicyOnClusterDeletion.IsNull() && !plan.BackupRetentionPolicyOnClusterDeletion.IsUnknown() {
		planValues["BackupRetentionPolicyOnClusterDeletion"] = fwValueToInterface(plan.BackupRetentionPolicyOnClusterDeletion, map[string]fwWireField{}, false)
	}
	if !plan.BackupSetSize.IsNull() && !plan.BackupSetSize.IsUnknown() {
		planValues["BackupSetSize"] = fwValueToInterface(plan.BackupSetSize, map[string]fwWireField{}, false)
	}
	if !plan.BackupStartTime.IsNull() && !plan.BackupStartTime.IsUnknown() {
		planValues["BackupStartTime"] = fwValueToInterface(plan.BackupStartTime, map[string]fwWireField{}, false)
	}
	if !plan.BackupType.IsNull() && !plan.BackupType.IsUnknown() {
		planValues["BackupType"] = fwValueToInterface(plan.BackupType, map[string]fwWireField{}, false)
	}
	if !plan.BackupsLevel.IsNull() && !plan.BackupsLevel.IsUnknown() {
		planValues["BackupsLevel"] = fwValueToInterface(plan.BackupsLevel, map[string]fwWireField{}, false)
	}
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		planValues["Comment"] = fwValueToInterface(plan.Comment, map[string]fwWireField{}, false)
	}
	if !plan.ConsistentTime.IsNull() && !plan.ConsistentTime.IsUnknown() {
		planValues["ConsistentTime"] = fwValueToInterface(plan.ConsistentTime, map[string]fwWireField{}, false)
	}
	if !plan.CreateTime.IsNull() && !plan.CreateTime.IsUnknown() {
		planValues["CreateTime"] = fwValueToInterface(plan.CreateTime, map[string]fwWireField{}, false)
	}
	if !plan.DBClusterId.IsNull() && !plan.DBClusterId.IsUnknown() {
		planValues["DBClusterId"] = fwValueToInterface(plan.DBClusterId, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupFrequency.IsNull() && !plan.DataLevel1BackupFrequency.IsUnknown() {
		planValues["DataLevel1BackupFrequency"] = fwValueToInterface(plan.DataLevel1BackupFrequency, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupPeriod.IsNull() && !plan.DataLevel1BackupPeriod.IsUnknown() {
		planValues["DataLevel1BackupPeriod"] = fwValueToInterface(plan.DataLevel1BackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupRetentionPeriod.IsNull() && !plan.DataLevel1BackupRetentionPeriod.IsUnknown() {
		planValues["DataLevel1BackupRetentionPeriod"] = fwValueToInterface(plan.DataLevel1BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel1BackupTime.IsNull() && !plan.DataLevel1BackupTime.IsUnknown() {
		planValues["DataLevel1BackupTime"] = fwValueToInterface(plan.DataLevel1BackupTime, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupAnotherRegionRegion.IsNull() && !plan.DataLevel2BackupAnotherRegionRegion.IsUnknown() {
		planValues["DataLevel2BackupAnotherRegionRegion"] = fwValueToInterface(plan.DataLevel2BackupAnotherRegionRegion, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupAnotherRegionRetentionPeriod.IsNull() && !plan.DataLevel2BackupAnotherRegionRetentionPeriod.IsUnknown() {
		planValues["DataLevel2BackupAnotherRegionRetentionPeriod"] = fwValueToInterface(plan.DataLevel2BackupAnotherRegionRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupPeriod.IsNull() && !plan.DataLevel2BackupPeriod.IsUnknown() {
		planValues["DataLevel2BackupPeriod"] = fwValueToInterface(plan.DataLevel2BackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.DataLevel2BackupRetentionPeriod.IsNull() && !plan.DataLevel2BackupRetentionPeriod.IsUnknown() {
		planValues["DataLevel2BackupRetentionPeriod"] = fwValueToInterface(plan.DataLevel2BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !plan.IsAvail.IsNull() && !plan.IsAvail.IsUnknown() {
		planValues["IsAvail"] = fwValueToInterface(plan.IsAvail, map[string]fwWireField{}, false)
	}
	if !plan.JobMode.IsNull() && !plan.JobMode.IsUnknown() {
		planValues["JobMode"] = fwValueToInterface(plan.JobMode, map[string]fwWireField{}, false)
	}
	if !plan.PreferredBackupPeriod.IsNull() && !plan.PreferredBackupPeriod.IsUnknown() {
		planValues["PreferredBackupPeriod"] = fwValueToInterface(plan.PreferredBackupPeriod, map[string]fwWireField{}, false)
	}
	if !plan.PreferredBackupTime.IsNull() && !plan.PreferredBackupTime.IsUnknown() {
		planValues["PreferredBackupTime"] = fwValueToInterface(plan.PreferredBackupTime, map[string]fwWireField{}, false)
	}
	if !plan.PreferredNextBackupTime.IsNull() && !plan.PreferredNextBackupTime.IsUnknown() {
		planValues["PreferredNextBackupTime"] = fwValueToInterface(plan.PreferredNextBackupTime, map[string]fwWireField{}, false)
	}
	if !plan.Process.IsNull() && !plan.Process.IsUnknown() {
		planValues["Process"] = fwValueToInterface(plan.Process, map[string]fwWireField{}, false)
	}
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		planValues["Status"] = fwValueToInterface(plan.Status, map[string]fwWireField{}, false)
	}
	if !plan.StoreStatus.IsNull() && !plan.StoreStatus.IsUnknown() {
		planValues["StoreStatus"] = fwValueToInterface(plan.StoreStatus, map[string]fwWireField{}, false)
	}
	if !plan.TaskAction.IsNull() && !plan.TaskAction.IsUnknown() {
		planValues["TaskAction"] = fwValueToInterface(plan.TaskAction, map[string]fwWireField{}, false)
	}
	stateValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.AdvancedDataPolicies.IsNull() && !state.AdvancedDataPolicies.IsUnknown() {
		stateValues["AdvancedDataPolicies"] = fwValueToInterface(state.AdvancedDataPolicies, map[string]fwWireField{"action_type": {Name: "ActionType", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_hour": {Name: "OnlyPreserveOneEachHour", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "storage_class": {Name: "StorageClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_value": {Name: "RetentionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "bak_type": {Name: "BakType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_key": {Name: "FilterKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_value": {Name: "FilterValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_type": {Name: "DestType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_type": {Name: "FilterType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_region": {Name: "SrcRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "auto_created": {Name: "AutoCreated", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_day": {Name: "OnlyPreserveOneEachDay", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_type": {Name: "SrcType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dump_action": {Name: "DumpAction", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_region": {Name: "DestRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "policy_id": {Name: "PolicyId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_type": {Name: "RetentionType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.BackupEndTime.IsNull() && !state.BackupEndTime.IsUnknown() {
		stateValues["BackupEndTime"] = fwValueToInterface(state.BackupEndTime, map[string]fwWireField{}, false)
	}
	if !state.BackupFrequency.IsNull() && !state.BackupFrequency.IsUnknown() {
		stateValues["BackupFrequency"] = fwValueToInterface(state.BackupFrequency, map[string]fwWireField{}, false)
	}
	if !state.BackupId.IsNull() && !state.BackupId.IsUnknown() {
		stateValues["BackupId"] = fwValueToInterface(state.BackupId, map[string]fwWireField{}, false)
	}
	if !state.BackupJobId.IsNull() && !state.BackupJobId.IsUnknown() {
		stateValues["BackupJobId"] = fwValueToInterface(state.BackupJobId, map[string]fwWireField{}, false)
	}
	if !state.BackupMethod.IsNull() && !state.BackupMethod.IsUnknown() {
		stateValues["BackupMethod"] = fwValueToInterface(state.BackupMethod, map[string]fwWireField{}, false)
	}
	if !state.BackupMode.IsNull() && !state.BackupMode.IsUnknown() {
		stateValues["BackupMode"] = fwValueToInterface(state.BackupMode, map[string]fwWireField{}, false)
	}
	if !state.BackupPolicyLevel.IsNull() && !state.BackupPolicyLevel.IsUnknown() {
		stateValues["BackupPolicyLevel"] = fwValueToInterface(state.BackupPolicyLevel, map[string]fwWireField{}, false)
	}
	if !state.BackupProgressStatus.IsNull() && !state.BackupProgressStatus.IsUnknown() {
		stateValues["BackupProgressStatus"] = fwValueToInterface(state.BackupProgressStatus, map[string]fwWireField{}, false)
	}
	if !state.BackupRetentionPolicyOnClusterDeletion.IsNull() && !state.BackupRetentionPolicyOnClusterDeletion.IsUnknown() {
		stateValues["BackupRetentionPolicyOnClusterDeletion"] = fwValueToInterface(state.BackupRetentionPolicyOnClusterDeletion, map[string]fwWireField{}, false)
	}
	if !state.BackupSetSize.IsNull() && !state.BackupSetSize.IsUnknown() {
		stateValues["BackupSetSize"] = fwValueToInterface(state.BackupSetSize, map[string]fwWireField{}, false)
	}
	if !state.BackupStartTime.IsNull() && !state.BackupStartTime.IsUnknown() {
		stateValues["BackupStartTime"] = fwValueToInterface(state.BackupStartTime, map[string]fwWireField{}, false)
	}
	if !state.BackupType.IsNull() && !state.BackupType.IsUnknown() {
		stateValues["BackupType"] = fwValueToInterface(state.BackupType, map[string]fwWireField{}, false)
	}
	if !state.BackupsLevel.IsNull() && !state.BackupsLevel.IsUnknown() {
		stateValues["BackupsLevel"] = fwValueToInterface(state.BackupsLevel, map[string]fwWireField{}, false)
	}
	if !state.Comment.IsNull() && !state.Comment.IsUnknown() {
		stateValues["Comment"] = fwValueToInterface(state.Comment, map[string]fwWireField{}, false)
	}
	if !state.ConsistentTime.IsNull() && !state.ConsistentTime.IsUnknown() {
		stateValues["ConsistentTime"] = fwValueToInterface(state.ConsistentTime, map[string]fwWireField{}, false)
	}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		stateValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.DBClusterId.IsNull() && !state.DBClusterId.IsUnknown() {
		stateValues["DBClusterId"] = fwValueToInterface(state.DBClusterId, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupFrequency.IsNull() && !state.DataLevel1BackupFrequency.IsUnknown() {
		stateValues["DataLevel1BackupFrequency"] = fwValueToInterface(state.DataLevel1BackupFrequency, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupPeriod.IsNull() && !state.DataLevel1BackupPeriod.IsUnknown() {
		stateValues["DataLevel1BackupPeriod"] = fwValueToInterface(state.DataLevel1BackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupRetentionPeriod.IsNull() && !state.DataLevel1BackupRetentionPeriod.IsUnknown() {
		stateValues["DataLevel1BackupRetentionPeriod"] = fwValueToInterface(state.DataLevel1BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupTime.IsNull() && !state.DataLevel1BackupTime.IsUnknown() {
		stateValues["DataLevel1BackupTime"] = fwValueToInterface(state.DataLevel1BackupTime, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupAnotherRegionRegion.IsNull() && !state.DataLevel2BackupAnotherRegionRegion.IsUnknown() {
		stateValues["DataLevel2BackupAnotherRegionRegion"] = fwValueToInterface(state.DataLevel2BackupAnotherRegionRegion, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupAnotherRegionRetentionPeriod.IsNull() && !state.DataLevel2BackupAnotherRegionRetentionPeriod.IsUnknown() {
		stateValues["DataLevel2BackupAnotherRegionRetentionPeriod"] = fwValueToInterface(state.DataLevel2BackupAnotherRegionRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupPeriod.IsNull() && !state.DataLevel2BackupPeriod.IsUnknown() {
		stateValues["DataLevel2BackupPeriod"] = fwValueToInterface(state.DataLevel2BackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupRetentionPeriod.IsNull() && !state.DataLevel2BackupRetentionPeriod.IsUnknown() {
		stateValues["DataLevel2BackupRetentionPeriod"] = fwValueToInterface(state.DataLevel2BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.IsAvail.IsNull() && !state.IsAvail.IsUnknown() {
		stateValues["IsAvail"] = fwValueToInterface(state.IsAvail, map[string]fwWireField{}, false)
	}
	if !state.JobMode.IsNull() && !state.JobMode.IsUnknown() {
		stateValues["JobMode"] = fwValueToInterface(state.JobMode, map[string]fwWireField{}, false)
	}
	if !state.PreferredBackupPeriod.IsNull() && !state.PreferredBackupPeriod.IsUnknown() {
		stateValues["PreferredBackupPeriod"] = fwValueToInterface(state.PreferredBackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.PreferredBackupTime.IsNull() && !state.PreferredBackupTime.IsUnknown() {
		stateValues["PreferredBackupTime"] = fwValueToInterface(state.PreferredBackupTime, map[string]fwWireField{}, false)
	}
	if !state.PreferredNextBackupTime.IsNull() && !state.PreferredNextBackupTime.IsUnknown() {
		stateValues["PreferredNextBackupTime"] = fwValueToInterface(state.PreferredNextBackupTime, map[string]fwWireField{}, false)
	}
	if !state.Process.IsNull() && !state.Process.IsUnknown() {
		stateValues["Process"] = fwValueToInterface(state.Process, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		stateValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	if !state.StoreStatus.IsNull() && !state.StoreStatus.IsUnknown() {
		stateValues["StoreStatus"] = fwValueToInterface(state.StoreStatus, map[string]fwWireField{}, false)
	}
	if !state.TaskAction.IsNull() && !state.TaskAction.IsUnknown() {
		stateValues["TaskAction"] = fwValueToInterface(state.TaskAction, map[string]fwWireField{}, false)
	}
	updated0 := false
	if fwPlannedChanged(plan.BackupFrequency, state.BackupFrequency) {
		updated0 = true
	}
	if fwPlannedChanged(plan.BackupPolicyLevel, state.BackupPolicyLevel) {
		updated0 = true
	}
	if fwPlannedChanged(plan.BackupRetentionPolicyOnClusterDeletion, state.BackupRetentionPolicyOnClusterDeletion) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupFrequency, state.DataLevel1BackupFrequency) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupPeriod, state.DataLevel1BackupPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupRetentionPeriod, state.DataLevel1BackupRetentionPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel1BackupTime, state.DataLevel1BackupTime) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupAnotherRegionRegion, state.DataLevel2BackupAnotherRegionRegion) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupAnotherRegionRetentionPeriod, state.DataLevel2BackupAnotherRegionRetentionPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupPeriod, state.DataLevel2BackupPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.DataLevel2BackupRetentionPeriod, state.DataLevel2BackupRetentionPeriod) {
		updated0 = true
	}
	// F7: preferred_backup_period is an unordered weekday set; an API-reordered value
	// (e.g. Wednesday,Monday -> Monday,Wednesday) must not trigger a policy update, while a
	// genuine add/remove of weekdays still must.
	if !backupWeekdaySetEqual(plan.PreferredBackupPeriod.ValueString(), state.PreferredBackupPeriod.ValueString()) {
		updated0 = true
	}
	if fwPlannedChanged(plan.PreferredBackupTime, state.PreferredBackupTime) {
		updated0 = true
	}
	// F4: detect advanced policy changes so ModifyBackupPolicy is called even when only
	// AdvancedDataPolicies differ and all scalar policy fields are unchanged. A changed
	// entry must be sent whole: DBS cannot identify a policy from a delta (see the
	// DBS.NotExisted note above), so keep every advanced-policy field in that request.
	if updated0 || advancedPoliciesChanged {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"DBClusterId": state.DBClusterId.ValueString(),
			"ClientToken": buildClientToken("ModifyBackupPolicy"),
		}
		request0["RegionId"] = r.Client().RegionId
		requiredSources := []string{"$.DBClusterId"}
		if advancedPoliciesChanged {
			requiredSources = append(requiredSources, "$.AdvancedDataPolicies")
		}
		if err := fwApplyMappings(request0, operationValues0, fwChangedUpdateMappings([]fwMapping{
			{Source: "$.AdvancedDataPolicies[*].ActionType", Target: "AdvancedDataPolicies[*].ActionType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].AutoCreated", Target: "AdvancedDataPolicies[*].AutoCreated", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].BakType", Target: "AdvancedDataPolicies[*].BakType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].DestRegion", Target: "AdvancedDataPolicies[*].DestRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].DestType", Target: "AdvancedDataPolicies[*].DestType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].DumpAction", Target: "AdvancedDataPolicies[*].DumpAction", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].FilterKey", Target: "AdvancedDataPolicies[*].FilterKey", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].FilterType", Target: "AdvancedDataPolicies[*].FilterType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].FilterValue", Target: "AdvancedDataPolicies[*].FilterValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachDay", Target: "AdvancedDataPolicies[*].OnlyPreserveOneEachDay", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachHour", Target: "AdvancedDataPolicies[*].OnlyPreserveOneEachHour", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].PolicyId", Target: "AdvancedDataPolicies[*].PolicyId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].RetentionType", Target: "AdvancedDataPolicies[*].RetentionType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].RetentionValue", Target: "AdvancedDataPolicies[*].RetentionValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].SrcRegion", Target: "AdvancedDataPolicies[*].SrcRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].SrcType", Target: "AdvancedDataPolicies[*].SrcType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.AdvancedDataPolicies[*].StorageClass", Target: "AdvancedDataPolicies[*].StorageClass", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupFrequency", Target: "BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupPolicyLevel", Target: "BackupPolicyLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupRetentionPolicyOnClusterDeletion", Target: "BackupRetentionPolicyOnClusterDeletion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupFrequency", Target: "DataLevel1BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupPeriod", Target: "DataLevel1BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupRetentionPeriod", Target: "DataLevel1BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel1BackupTime", Target: "DataLevel1BackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupAnotherRegionRegion", Target: "DataLevel2BackupAnotherRegionRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupAnotherRegionRetentionPeriod", Target: "DataLevel2BackupAnotherRegionRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupPeriod", Target: "DataLevel2BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DataLevel2BackupRetentionPeriod", Target: "DataLevel2BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PreferredBackupPeriod", Target: "PreferredBackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.PreferredBackupTime", Target: "PreferredBackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, operationValues0, stateValues, requiredSources), false); err != nil {
			resp.Diagnostics.AddError("Building ModifyBackupPolicy request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "ModifyBackupPolicy", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_polardb_backup: ModifyBackupPolicy", err.Error())
			return
		}
	}
	// RestoreTable is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	object, err := backupGet(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeBackupPolicy", backupTypeName),
			errs.WrapErrorf(err, "Reading %s: calling DescribeBackupPolicy", backupTypeName).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, backupState(object, &plan))...)
}

func (r *backupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 5*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteValues := map[string]interface{}{"RegionId": r.Client().RegionId}
	if !state.AdvancedDataPolicies.IsNull() && !state.AdvancedDataPolicies.IsUnknown() {
		deleteValues["AdvancedDataPolicies"] = fwValueToInterface(state.AdvancedDataPolicies, map[string]fwWireField{"action_type": {Name: "ActionType", Single: false, Read: false, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_hour": {Name: "OnlyPreserveOneEachHour", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "storage_class": {Name: "StorageClass", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_value": {Name: "RetentionValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "bak_type": {Name: "BakType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_key": {Name: "FilterKey", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_value": {Name: "FilterValue", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_type": {Name: "DestType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "filter_type": {Name: "FilterType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_region": {Name: "SrcRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "auto_created": {Name: "AutoCreated", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "only_preserve_one_each_day": {Name: "OnlyPreserveOneEachDay", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "src_type": {Name: "SrcType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dump_action": {Name: "DumpAction", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "dest_region": {Name: "DestRegion", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "policy_id": {Name: "PolicyId", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}, "retention_type": {Name: "RetentionType", Single: false, Read: true, Input: true, Fields: map[string]fwWireField{}}}, false)
	}
	if !state.BackupEndTime.IsNull() && !state.BackupEndTime.IsUnknown() {
		deleteValues["BackupEndTime"] = fwValueToInterface(state.BackupEndTime, map[string]fwWireField{}, false)
	}
	if !state.BackupFrequency.IsNull() && !state.BackupFrequency.IsUnknown() {
		deleteValues["BackupFrequency"] = fwValueToInterface(state.BackupFrequency, map[string]fwWireField{}, false)
	}
	if !state.BackupId.IsNull() && !state.BackupId.IsUnknown() {
		deleteValues["BackupId"] = fwValueToInterface(state.BackupId, map[string]fwWireField{}, false)
	}
	if !state.BackupJobId.IsNull() && !state.BackupJobId.IsUnknown() {
		deleteValues["BackupJobId"] = fwValueToInterface(state.BackupJobId, map[string]fwWireField{}, false)
	}
	if !state.BackupMethod.IsNull() && !state.BackupMethod.IsUnknown() {
		deleteValues["BackupMethod"] = fwValueToInterface(state.BackupMethod, map[string]fwWireField{}, false)
	}
	if !state.BackupMode.IsNull() && !state.BackupMode.IsUnknown() {
		deleteValues["BackupMode"] = fwValueToInterface(state.BackupMode, map[string]fwWireField{}, false)
	}
	if !state.BackupPolicyLevel.IsNull() && !state.BackupPolicyLevel.IsUnknown() {
		deleteValues["BackupPolicyLevel"] = fwValueToInterface(state.BackupPolicyLevel, map[string]fwWireField{}, false)
	}
	if !state.BackupProgressStatus.IsNull() && !state.BackupProgressStatus.IsUnknown() {
		deleteValues["BackupProgressStatus"] = fwValueToInterface(state.BackupProgressStatus, map[string]fwWireField{}, false)
	}
	if !state.BackupRetentionPolicyOnClusterDeletion.IsNull() && !state.BackupRetentionPolicyOnClusterDeletion.IsUnknown() {
		deleteValues["BackupRetentionPolicyOnClusterDeletion"] = fwValueToInterface(state.BackupRetentionPolicyOnClusterDeletion, map[string]fwWireField{}, false)
	}
	if !state.BackupSetSize.IsNull() && !state.BackupSetSize.IsUnknown() {
		deleteValues["BackupSetSize"] = fwValueToInterface(state.BackupSetSize, map[string]fwWireField{}, false)
	}
	if !state.BackupStartTime.IsNull() && !state.BackupStartTime.IsUnknown() {
		deleteValues["BackupStartTime"] = fwValueToInterface(state.BackupStartTime, map[string]fwWireField{}, false)
	}
	if !state.BackupType.IsNull() && !state.BackupType.IsUnknown() {
		deleteValues["BackupType"] = fwValueToInterface(state.BackupType, map[string]fwWireField{}, false)
	}
	if !state.BackupsLevel.IsNull() && !state.BackupsLevel.IsUnknown() {
		deleteValues["BackupsLevel"] = fwValueToInterface(state.BackupsLevel, map[string]fwWireField{}, false)
	}
	if !state.Comment.IsNull() && !state.Comment.IsUnknown() {
		deleteValues["Comment"] = fwValueToInterface(state.Comment, map[string]fwWireField{}, false)
	}
	if !state.ConsistentTime.IsNull() && !state.ConsistentTime.IsUnknown() {
		deleteValues["ConsistentTime"] = fwValueToInterface(state.ConsistentTime, map[string]fwWireField{}, false)
	}
	if !state.CreateTime.IsNull() && !state.CreateTime.IsUnknown() {
		deleteValues["CreateTime"] = fwValueToInterface(state.CreateTime, map[string]fwWireField{}, false)
	}
	if !state.DBClusterId.IsNull() && !state.DBClusterId.IsUnknown() {
		deleteValues["DBClusterId"] = fwValueToInterface(state.DBClusterId, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupFrequency.IsNull() && !state.DataLevel1BackupFrequency.IsUnknown() {
		deleteValues["DataLevel1BackupFrequency"] = fwValueToInterface(state.DataLevel1BackupFrequency, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupPeriod.IsNull() && !state.DataLevel1BackupPeriod.IsUnknown() {
		deleteValues["DataLevel1BackupPeriod"] = fwValueToInterface(state.DataLevel1BackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupRetentionPeriod.IsNull() && !state.DataLevel1BackupRetentionPeriod.IsUnknown() {
		deleteValues["DataLevel1BackupRetentionPeriod"] = fwValueToInterface(state.DataLevel1BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel1BackupTime.IsNull() && !state.DataLevel1BackupTime.IsUnknown() {
		deleteValues["DataLevel1BackupTime"] = fwValueToInterface(state.DataLevel1BackupTime, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupAnotherRegionRegion.IsNull() && !state.DataLevel2BackupAnotherRegionRegion.IsUnknown() {
		deleteValues["DataLevel2BackupAnotherRegionRegion"] = fwValueToInterface(state.DataLevel2BackupAnotherRegionRegion, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupAnotherRegionRetentionPeriod.IsNull() && !state.DataLevel2BackupAnotherRegionRetentionPeriod.IsUnknown() {
		deleteValues["DataLevel2BackupAnotherRegionRetentionPeriod"] = fwValueToInterface(state.DataLevel2BackupAnotherRegionRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupPeriod.IsNull() && !state.DataLevel2BackupPeriod.IsUnknown() {
		deleteValues["DataLevel2BackupPeriod"] = fwValueToInterface(state.DataLevel2BackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.DataLevel2BackupRetentionPeriod.IsNull() && !state.DataLevel2BackupRetentionPeriod.IsUnknown() {
		deleteValues["DataLevel2BackupRetentionPeriod"] = fwValueToInterface(state.DataLevel2BackupRetentionPeriod, map[string]fwWireField{}, false)
	}
	if !state.IsAvail.IsNull() && !state.IsAvail.IsUnknown() {
		deleteValues["IsAvail"] = fwValueToInterface(state.IsAvail, map[string]fwWireField{}, false)
	}
	if !state.JobMode.IsNull() && !state.JobMode.IsUnknown() {
		deleteValues["JobMode"] = fwValueToInterface(state.JobMode, map[string]fwWireField{}, false)
	}
	if !state.PreferredBackupPeriod.IsNull() && !state.PreferredBackupPeriod.IsUnknown() {
		deleteValues["PreferredBackupPeriod"] = fwValueToInterface(state.PreferredBackupPeriod, map[string]fwWireField{}, false)
	}
	if !state.PreferredBackupTime.IsNull() && !state.PreferredBackupTime.IsUnknown() {
		deleteValues["PreferredBackupTime"] = fwValueToInterface(state.PreferredBackupTime, map[string]fwWireField{}, false)
	}
	if !state.PreferredNextBackupTime.IsNull() && !state.PreferredNextBackupTime.IsUnknown() {
		deleteValues["PreferredNextBackupTime"] = fwValueToInterface(state.PreferredNextBackupTime, map[string]fwWireField{}, false)
	}
	if !state.Process.IsNull() && !state.Process.IsUnknown() {
		deleteValues["Process"] = fwValueToInterface(state.Process, map[string]fwWireField{}, false)
	}
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		deleteValues["Status"] = fwValueToInterface(state.Status, map[string]fwWireField{}, false)
	}
	if !state.StoreStatus.IsNull() && !state.StoreStatus.IsUnknown() {
		deleteValues["StoreStatus"] = fwValueToInterface(state.StoreStatus, map[string]fwWireField{}, false)
	}
	if !state.TaskAction.IsNull() && !state.TaskAction.IsUnknown() {
		deleteValues["TaskAction"] = fwValueToInterface(state.TaskAction, map[string]fwWireField{}, false)
	}
	// F1: the resource id is the composite "DBClusterId:BackupId"; resolve both parts so
	// DeleteBackup targets the exact owned backup. The prior line passed the full composite
	// id as DBClusterId, which would delete the wrong/foreign backup or fail identity checks.
	deleteClusterId, deleteBackupId := backupParseID(state.Id.ValueString())
	deleteValues["DBClusterId"] = deleteClusterId
	if deleteBackupId != "" {
		deleteValues["BackupId"] = deleteBackupId
	}
	selected := false
	deleteEligible0, deleteConditionErr0 := fwDeleteEligible(deleteValues, []fwDeleteCondition{})
	if deleteConditionErr0 != nil {
		resp.Diagnostics.AddError("Selecting deletion operation", deleteConditionErr0.Error())
		return
	}
	selected = selected || deleteEligible0
	if !selected {
		resp.Diagnostics.AddError("Selecting deletion operation", "No deletion operation matches the resource state; the resource remains tracked")
		return
	}
	if deleteEligible0 {
		request0 := map[string]interface{}{"ClientToken": buildClientToken("DeleteBackup")}
		request0["RegionId"] = r.Client().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.BackupId", Target: "BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteBackup request", err.Error())
			return
		}
		if _, err := retryRpcPost(r.Client(), "polardb", "2017-08-01", "DeleteBackup", request0, operationTimeout); err != nil && !isExpectedErrors(err, backupDeleteNotFoundCodes) {
			resp.Diagnostics.AddError("Deleting alicloud_polardb_backup: DeleteBackup", err.Error())
			return
		}
	}
}

func (r *backupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// backupState maps the Get response object onto the resource model. Fields present in the
// response are read via their responseKey (e.g. status -> ListenerStatus); request-only fields
// (responseKey empty, e.g. dry_run) are preserved from the prior plan/state instead.
func backupState(object map[string]interface{}, prior *backupResourceModel) backupResourceModel {
	// Get has already merged and normalized every response mapping.
	str := func(key string) types.String {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return types.StringNull()
			}
			cur, ok = m[seg]
			if !ok {
				return types.StringNull()
			}
		}
		// gap#7c: present-but-empty strings must return ValueString(""), not Null.
		// v1 d.Set(field, "") treats empty as set (not absence); TestCheckResourceAttrSet
		// expects a non-null value for API-returned "" (e.g. DhcpOptionsSetId when VPC has
		// no DHCP set). The earlier && v != "" collapsed present-empty to Null, failing Set
		// checks. Absent fields still return Null via the !ok branches above.
		if v, ok := cur.(string); ok {
			return types.StringValue(v)
		}
		return types.StringNull()
	}
	// lookup traverses a dotted responseKey (e.g. "RegionStatuses.RegionStatus") segment by segment
	// from the Get response object, returning the raw value (array/map/scalar) or nil. gap#5: array
	// fields' responseKey is the array's true path from the read-tree [*] child, which may be a
	// two-level unwrap (object["RegionStatuses"]["RegionStatus"]) that a literal object[key] lookup
	// cannot reach. Single-segment keys (e.g. "IpWhiteList") traverse one step == object[key].
	lookup := func(key string) interface{} {
		cur := interface{}(object)
		for _, seg := range strings.Split(key, ".") {
			m, ok := cur.(map[string]interface{})
			if !ok {
				return nil
			}
			cur, ok = m[seg]
			if !ok {
				return nil
			}
		}
		return cur
	}
	// F1: the resource id is the composite "DBClusterId:BackupId". The stable BackupId is
	// injected by backupGet and materialized by DescribeBackups; fall back to BackupJobId
	// (the create task id, equal to the backup set id) when BackupId is absent. A bare
	// DBClusterId-only id would let Read/Delete lose the owned backup identity ("last wins").
	backupStateId := str("DBClusterId").ValueString()
	backupSetId := str("BackupId").ValueString()
	if backupSetId == "" {
		backupSetId = str("BackupJobId").ValueString()
	}
	if backupSetId != "" {
		backupStateId = backupComposeID(backupStateId, backupSetId)
	}
	m := backupResourceModel{
		FwTimeouts:                             types.ObjectNull(fwTimeoutTypes()),
		Id:                                     types.StringValue(backupStateId),
		AdvancedDataPolicies:                   backupAdvancedDataPoliciesToTypeList(lookup("AdvancedDataPolicies")),
		BackupEndTime:                          str("BackupEndTime"),
		BackupFrequency:                        str("BackupFrequency"),
		BackupId:                               str("BackupId"),
		BackupJobId:                            str("BackupJobId"),
		BackupMethod:                           str("BackupMethod"),
		BackupMode:                             str("BackupMode"),
		BackupPolicyLevel:                      str("BackupPolicyLevel"),
		BackupProgressStatus:                   str("BackupProgressStatus"),
		BackupRetentionPolicyOnClusterDeletion: str("BackupRetentionPolicyOnClusterDeletion"),
		BackupSetSize:                          str("BackupSetSize"),
		BackupStartTime:                        str("BackupStartTime"),
		BackupType:                             str("BackupType"),
		BackupsLevel:                           str("BackupsLevel"),
		Comment:                                str("Comment"),
		ConsistentTime:                         str("ConsistentTime"),
		CreateTime:                             str("CreateTime"),
		DBClusterId:                            str("DBClusterId"),
		DataLevel1BackupFrequency:              str("DataLevel1BackupFrequency"),
		DataLevel1BackupPeriod:                 str("DataLevel1BackupPeriod"),
		DataLevel1BackupRetentionPeriod:        str("DataLevel1BackupRetentionPeriod"),
		DataLevel1BackupTime:                   str("DataLevel1BackupTime"),
		DataLevel2BackupAnotherRegionRegion:    str("DataLevel2BackupAnotherRegionRegion"),
		DataLevel2BackupAnotherRegionRetentionPeriod: str("DataLevel2BackupAnotherRegionRetentionPeriod"),
		DataLevel2BackupPeriod:                       str("DataLevel2BackupPeriod"),
		DataLevel2BackupRetentionPeriod:              str("DataLevel2BackupRetentionPeriod"),
		IsAvail:                                      str("IsAvail"),
		JobMode:                                      str("JobMode"),
		PreferredBackupPeriod:                        preferredBackupPeriodStringValue{str("PreferredBackupPeriod")},
		PreferredBackupTime:                          str("PreferredBackupTime"),
		PreferredNextBackupTime:                      str("PreferredNextBackupTime"),
		Process:                                      str("Process"),
		Status:                                       str("Status"),
		StoreStatus:                                  str("StoreStatus"),
		TaskAction:                                   str("TaskAction"),
	}
	if prior != nil {
		m.FwTimeouts = prior.FwTimeouts
		// advanced_data_policies is a configuration-owned nested block, and Plugin
		// Framework nested blocks cannot be Computed. DescribeBackupPolicy returns the
		// system-generated policies next to the configured ones, so taking the read
		// result would make the applied state diverge from the planned block count
		// ("block count changed from 0 to 1"). Keep the configured entries instead and
		// only fill their computed attributes (policy_id, ...) from the matching
		// read-back entry, so the state stays known and later plans are stable.
		m.AdvancedDataPolicies = backupAlignAdvancedPolicies(prior.AdvancedDataPolicies, m.AdvancedDataPolicies)
		// Comment is a create-only field (set via CreateBackup) that is not
		// reliably read back from DescribeBackups right after creation: the
		// backup-set list may not include the just-created backup yet, or its
		// entries may carry empty comments, so the response mapping yields null.
		// Preserve the prior/plan value when the read returns null so apply
		// stays consistent; this mirrors the test's import-ignore for comment.
		if m.Comment.IsNull() {
			m.Comment = prior.Comment
		}
		// Under BackupPolicyLevel=Advanced the server ignores PreferredBackupPeriod and
		// reports all seven days (ModifyBackupPolicy doc: 开启高级备份后，该参数不再生效).
		// Mirroring that read-back would contradict the configured value in apply and
		// leave a permanent diff, so keep the plan/state value while advanced is enabled.
		if !prior.PreferredBackupPeriod.IsNull() && !prior.PreferredBackupPeriod.IsUnknown() &&
			!m.BackupPolicyLevel.IsNull() && !m.BackupPolicyLevel.IsUnknown() &&
			m.BackupPolicyLevel.ValueString() == "Advanced" {
			m.PreferredBackupPeriod = prior.PreferredBackupPeriod
		}
	}
	return m
}
