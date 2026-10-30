// Package polardb. Framework-native alicloud_polardb_backup resource.
// Nested attributes require Terraform Plugin Protocol 6.
package polardb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
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
				Description: "Advanced data backup policies. Only applied when backup_policy_level is Advanced.",
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
				Description: "The end time of the backup set. The format is `YYYY-MM-DDThh:mmZ` (UTC time).",
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
				Description: "The start time of the backup set. The format is `YYYY-MM-DDThh:mmZ` (UTC time).",
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
				Description: "Storage status. Only backup sets of a released cluster report this field; it is empty for backup sets of a running cluster.",
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

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(plan.FwTimeouts, "create", 20*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The comment is create-only (CreateBackup). An omitted (unknown) comment is pinned to ""
	// so it does not stay unknown after apply, while an explicitly configured value -- including
	// "" -- is preserved verbatim: turning a configured "" into null would differ from the
	// configuration forever and, together with RequiresReplaceIfConfigured, recreate the backup
	// on every plan. An empty comment is simply not sent.
	comment := ""
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		comment = plan.Comment.ValueString()
	}
	if plan.Comment.IsUnknown() {
		plan.Comment = types.StringValue("")
	}

	request := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("CreateBackup"),
	}
	// RPC APIs require RegionId as a request param (separate from endpoint resolution);
	// the SDK RpcPost resolves the endpoint by region but does NOT inject RegionId into query/body.
	// v1 sets request["RegionId"] = client.RegionId authoritatively (provider region wins).
	request["RegionId"] = r.Meta().RegionId
	createValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
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
	createMappings := []fwMapping{
		{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}
	if comment != "" {
		createMappings = append([]fwMapping{{Source: "$.Comment", Target: "Comment", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}}}, createMappings...)
	}
	if err := fwApplyMappings(request, createValues, createMappings, false); err != nil {
		resp.Diagnostics.AddError("Building create request", err.Error())
		return
	}

	// CreateBackup returns the backup task id; the backup set it produces is resolved from
	// that task's exact start time, never from "the newest set in a window".
	response, err := retryRpcPost(r.Meta(), "polardb", "2017-08-01", "CreateBackup", request, operationTimeout)
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
	// CreateBackup answers with the backup task id, which is stored as backup_job_id. The
	// backup set the task produces is listed under its own BackupId, so the set is bound by
	// the exact task start time (see backupResolveCreatedBackupSet); binding an arbitrary
	// backup set (a repeated Comment or the largest BackupId) could record a foreign backup.
	jobId := ""
	if values := fwWalkMappedPath(response, fwMappingPath("BackupJobId"), nil); len(values) > 0 && values[0].Value != nil {
		jobId = fmt.Sprint(values[0].Value)
	}
	if jobId == "" || jobId == "<nil>" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateBackup", backupTypeName),
			fmt.Sprintf("CreateBackup did not return a BackupJobId: %v", response),
		)
		return
	}
	// Wait for the CreateBackup task itself first: CreateBackup answers asynchronously and
	// the backup set only materializes once the task reaches Finished. Waiting on the set
	// alone would consume the whole create budget while the task is still running.
	if waitErr := backupWaitForBackupTask(r.Meta(), dbClusterId, jobId, operationTimeout); waitErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for the backup task", backupTypeName),
			errs.WrapErrorf(waitErr, "Creating %s: waiting for backup task %s on cluster %s", backupTypeName, jobId, dbClusterId).Error(),
		)
		return
	}
	backupId, backupSet, resolveErr := backupResolveCreatedBackupSet(r.Meta(), dbClusterId, jobId, backupCreateResolveTimeout)
	if resolveErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: locating the created backup set", backupTypeName),
			errs.WrapErrorf(resolveErr, "Creating %s: locating the backup set produced by task %s on cluster %s", backupTypeName, jobId, dbClusterId).Error(),
		)
		return
	}
	id := backupComposeID(dbClusterId, backupId)
	if status := backupStringField(backupSet, "BackupStatus"); status == "Failed" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: backup set failed", backupTypeName),
			fmt.Sprintf("backup set %s on cluster %s finished with status Failed", backupId, dbClusterId),
		)
		return
	}

	// Persist the allocated identity before any asynchronous or configuration stage.
	trackedPlan := plan
	trackedPlan.Id = types.StringValue(id)
	trackedPlan.BackupId = types.StringValue(backupId)
	trackedPlan.BackupJobId = types.StringValue(jobId)
	// The post-create reads rebuild the model from the API response, which never carries the
	// create task id; keep it on the planned value so the state keeps backup_job_id.
	plan.BackupJobId = types.StringValue(jobId)
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

	object, err := backupGetPolicyWithRetry(r.Meta(), id, 3*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after create", backupTypeName),
			errs.WrapErrorf(err, "Creating %s: reading after create (%s)", backupTypeName, id).Error(),
		)
		return
	}
	if object != nil {
		object = backupMergeResolvedBackupSet(object, backupSet, backupId)
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
	// AdvancedDataPolicies is only valid for BackupPolicyLevel=Advanced (ModifyBackupPolicy
	// contract), so the block must not be sent under Normal: doing so makes the server reject
	// the whole request, including the valid Normal policy fields in it.
	advancedEnabled := backupAdvancedLevelEnabled(backupEffectivePolicyLevel(plan.BackupPolicyLevel, object["BackupPolicyLevel"]).ValueString())
	planValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
	if advancedEnabled && !plan.AdvancedDataPolicies.IsNull() && !plan.AdvancedDataPolicies.IsUnknown() {
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
	stateValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
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
	// genuine add/remove of weekdays still must. An omitted (unknown) value is not a change.
	if backupWeekdaySetChanged(plan.PreferredBackupPeriod, state.PreferredBackupPeriod) {
		updated0 = true
	}
	if fwPlannedChanged(plan.PreferredBackupTime, state.PreferredBackupTime) {
		updated0 = true
	}
	// F4: detect advanced policy changes so ModifyBackupPolicy is called even when only
	// AdvancedDataPolicies differ and all scalar policy fields are unchanged.
	if updated0 || (advancedEnabled && backupAdvancedPoliciesChanged(plan.AdvancedDataPolicies, state.AdvancedDataPolicies)) {
		operationValues0 := planValues
		request0 := map[string]interface{}{
			"DBClusterId": dbClusterId,
			"ClientToken": helper.BuildClientToken("ModifyBackupPolicy"),
		}
		request0["RegionId"] = r.Meta().RegionId
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
		if _, err := retryRpcPost(r.Meta(), "polardb", "2017-08-01", "ModifyBackupPolicy", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_polardb_backup: ModifyBackupPolicy", err.Error())
			return
		}
	}
	object, err = backupGetPolicyWithRetry(r.Meta(), id, 3*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after post-create updates", backupTypeName),
			errs.WrapErrorf(err, "Creating %s: reading %s after the post-create policy update", backupTypeName, id).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: reading after post-create updates", backupTypeName),
			fmt.Sprintf("backup %s on cluster %s is still not visible through DescribeBackupPolicy three minutes after the post-create policy update", id, dbClusterId),
		)
		return
	}
	object = backupMergeResolvedBackupSet(object, backupSet, backupId)
	resp.Diagnostics.Append(resp.State.Set(ctx, backupState(object, &plan))...)
}

func (r *backupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state backupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A backup set can lag behind the create task in the list APIs for a couple of minutes, and
	// treating that lag as "deleted" drops a live resource from state (observed after a refresh
	// on a just-created backup). backupReadExisting retries that lag, resolves a backup set that
	// outlived its released cluster (P1), and only acts on an absence it can support: repeated
	// complete empty answers, weighted against the cluster's documented retention contract.
	object, err := backupReadExisting(r.Meta(), state.Id.ValueString(), state.BackupStartTime.ValueString(), state.BackupRetentionPolicyOnClusterDeletion.ValueString(), backupGoneGrace)
	if err != nil {
		if errors.Is(err, errBackupReadInconclusive) || errors.Is(err, errBackupAbsenceUnproven) {
			// The backup set may still exist (a source could not answer, or the retention
			// policy keeps it although the released-cluster list does not report it yet), so
			// the current state is kept instead of removing or failing the resource.
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: state kept", backupTypeName),
				fmt.Sprintf("the backup set of %s may still exist; the current state is kept unchanged: %v", state.Id.ValueString(), err),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", backupTypeName),
			errs.WrapErrorf(err, "Reading %s (%s)", backupTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
	if object == nil {
		// Never drop the resource silently: the backup set is no longer reported by any of the
		// sources after backupAbsenceConfirmations complete empty answers, and the retention
		// policy does not require it to exist, so it is removed from state with a diagnostic.
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Removing %s from state", backupTypeName),
			fmt.Sprintf("backup %s is no longer reported by the cluster backup policy, the cluster backup list, or the released-cluster backup list, and its retention policy does not keep it; it is removed from the Terraform state", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}
	// The create task id is not part of the backup-set response, so the task-owned fields are
	// refreshed (and the id kept) from the tracked backup_job_id. A missing task record is not
	// an absence: the task list only covers running/recent tasks.
	if jobId := state.BackupJobId.ValueString(); jobId != "" {
		if task, taskErr := backupFindBackupTask(r.Meta(), state.DBClusterId.ValueString(), jobId); taskErr == nil && task != nil {
			for _, field := range backupTaskResponseFields {
				if value, ok := task[field.Source]; ok && value != nil {
					object[field.Target] = value
				}
			}
		}
	}
	next := backupState(object, &state)
	if detached, _ := object[backupDetachedReadKey].(bool); detached {
		// The cluster is released: its backup policy no longer exists, so the policy fields are
		// kept from state instead of being reported as removed.
		backupPreserveBackupPolicyFields(&next, &state)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, next)...)
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
	// not exist".
	plan.AdvancedDataPolicies = backupAlignAdvancedPolicies(plan.AdvancedDataPolicies, state.AdvancedDataPolicies)
	// AdvancedDataPolicies is only valid for BackupPolicyLevel=Advanced, so a block under
	// Normal is ignored by the cloud and must not be sent either (see Create).
	advancedEnabled := backupAdvancedLevelEnabled(backupEffectivePolicyLevel(plan.BackupPolicyLevel, state.BackupPolicyLevel).ValueString())
	// Removing the block or one of its entries is not a protocol-level delete; reject it
	// instead of reporting success while the cloud policy survives.
	if advancedEnabled && backupAdvancedPoliciesRemoved(plan.AdvancedDataPolicies, state.AdvancedDataPolicies) {
		resp.Diagnostics.AddError(
			"Updating alicloud_polardb_backup: advanced_data_policies",
			"advanced_data_policies entries are not removed by dropping the block or an entry; keep the entry and set action_type = \"DELETE\" to remove the cloud policy.",
		)
		return
	}
	advancedPoliciesChanged := advancedEnabled && backupAdvancedPoliciesChanged(plan.AdvancedDataPolicies, state.AdvancedDataPolicies)

	planValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
	if advancedEnabled && !plan.AdvancedDataPolicies.IsNull() && !plan.AdvancedDataPolicies.IsUnknown() {
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
	stateValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
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
	// genuine add/remove of weekdays still must. An omitted (unknown) value is not a change.
	if backupWeekdaySetChanged(plan.PreferredBackupPeriod, state.PreferredBackupPeriod) {
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
			"ClientToken": helper.BuildClientToken("ModifyBackupPolicy"),
		}
		request0["RegionId"] = r.Meta().RegionId
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
		if _, err := retryRpcPost(r.Meta(), "polardb", "2017-08-01", "ModifyBackupPolicy", request0, operationTimeout); err != nil {
			resp.Diagnostics.AddError("Updating alicloud_polardb_backup: ModifyBackupPolicy", err.Error())
			return
		}
	}
	// RestoreTable is marked as a separate association resource in metadata.
	// Shared input fields do not make it part of this resource's update contract.
	// A backup set can lag behind the create task in the list APIs for a couple of minutes, and
	// treating that lag as "deleted" drops a live resource from state (observed after a refresh
	// on a just-created backup). backupReadExisting retries that lag and also resolves a backup
	// set that outlived its released cluster (P1).
	object, err := backupReadExisting(r.Meta(), state.Id.ValueString(), state.BackupStartTime.ValueString(), state.BackupRetentionPolicyOnClusterDeletion.ValueString(), backupGoneGrace)
	if err != nil {
		if errors.Is(err, errBackupReadInconclusive) || errors.Is(err, errBackupAbsenceUnproven) {
			// The backup set may still exist: keep the current state instead of removing or
			// failing the resource (see Read).
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Reading %s: state kept", backupTypeName),
				fmt.Sprintf("the backup set of %s may still exist; the current state is kept unchanged: %v", state.Id.ValueString(), err),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s", backupTypeName),
			errs.WrapErrorf(err, "Reading %s (%s)", backupTypeName, state.Id.ValueString()).Error(),
		)
		return
	}
	if object == nil {
		// Never drop the resource silently: the backup set is no longer reported by any of the
		// three sources (cluster policy, cluster backup list, released-cluster backup list) after
		// the lag tolerance, so it is removed from state with an explicit diagnostic.
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Removing %s from state", backupTypeName),
			fmt.Sprintf("backup %s is no longer reported by the cluster backup policy, the cluster backup list, or the released-cluster backup list; it is removed from the Terraform state", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}
	next := backupState(object, &plan)
	if detached, _ := object[backupDetachedReadKey].(bool); detached {
		backupPreserveBackupPolicyFields(&next, &state)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, next)...)
}

func (r *backupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	operationTimeout, timeoutDiagnostics := fwTimeoutDuration(state.FwTimeouts, "delete", 10*time.Minute)
	resp.Diagnostics.Append(timeoutDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteValues := map[string]interface{}{"RegionId": r.Meta().RegionId}
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
	// DeleteBackup targets the exact owned backup. A missing backup-set segment cannot address
	// a backup through DeleteBackup (BackupId is required), so it is rejected instead of
	// requesting a deletion that the API would refuse.
	deleteClusterId, deleteBackupId, idErr := backupValidateID(state.Id.ValueString())
	if idErr != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s", backupTypeName),
			idErr.Error(),
		)
		return
	}
	deleteValues["DBClusterId"] = deleteClusterId
	deleteValues["BackupId"] = deleteBackupId
	{
		request0 := map[string]interface{}{"ClientToken": helper.BuildClientToken("DeleteBackup")}
		request0["RegionId"] = r.Meta().RegionId
		if err := fwApplyMappings(request0, deleteValues, []fwMapping{
			{Source: "$.BackupId", Target: "BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			resp.Diagnostics.AddError("Building DeleteBackup request", err.Error())
			return
		}
		_, deleteErr := retryRpcPost(r.Meta(), "polardb", "2017-08-01", "DeleteBackup", request0, operationTimeout)
		if deleteErr == nil {
			return
		}
		// Classify the DeleteBackup result before anything else: only a documented
		// NotFound means the set is already gone or belongs to a released cluster. A permission,
		// conflict, throttling or validation error means the backup still exists, so it must
		// propagate — the retention lookup below must never turn it into a successful destroy.
		if !errs.IsExpectedErrors(deleteErr, backupDeleteNotFoundCodes) {
			resp.Diagnostics.AddError("Deleting alicloud_polardb_backup: DeleteBackup", deleteErr.Error())
			return
		}
		// P1: a backup set of a released cluster cannot be deleted through DeleteBackup (the API
		// requires a running cluster). When the detached list still reports this exact backup,
		// the retention policy keeps it: stop tracking it without claiming it was deleted.
		retained, retainedErr := backupFindRetainedBackupSet(r.Meta(), deleteClusterId, deleteBackupId, state.BackupStartTime.ValueString())
		if retainedErr != nil {
			// The retained-backup list could not answer, so whether the set survives its released
			// cluster is unknown: fail instead of reporting a destroy that was not confirmed.
			resp.Diagnostics.AddError(
				"Deleting alicloud_polardb_backup: DeleteBackup",
				errs.WrapErrorf(retainedErr, "backup set %s of cluster %s was reported absent by DeleteBackup, but the released-cluster backup list could not confirm whether it is retained", deleteBackupId, deleteClusterId).Error(),
			)
			return
		}
		if retained != nil {
			resp.Diagnostics.AddWarning(
				fmt.Sprintf("Deleting %s: backup set retained", backupTypeName),
				fmt.Sprintf("backup set %s of cluster %s belongs to a released cluster and is retained by its backup retention policy; DeleteBackup requires a running cluster, so the backup set stays in the cloud and is only removed from the Terraform state", deleteBackupId, deleteClusterId),
			)
			return
		}
	}
}

func (r *backupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Reject an id without both identity segments: importing "<db_cluster_id>" alone would
	// record a resource that has no backup-set identity, which Read would always report as
	// present and Delete could never target.
	if _, _, err := backupValidateID(req.ID); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Importing %s", backupTypeName), err.Error())
		return
	}
	// Reject an id whose backup set does not exist: passing the id through alone would import
	// an empty state and report success for a backup that is not there.
	object, err := backupGetWithRetry(r.Meta(), req.ID, 10*time.Second)
	if err != nil {
		if errors.Is(err, errBackupResourceGone) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Importing %s", backupTypeName),
				fmt.Sprintf("backup %s does not exist", req.ID),
			)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s", backupTypeName),
			errs.WrapErrorf(err, "Importing %s (%s)", backupTypeName, req.ID).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Importing %s", backupTypeName),
			fmt.Sprintf("backup %s does not exist", req.ID),
		)
		return
	}
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
		// Present-but-empty strings must return ValueString(""), not Null.
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
	// from the Get response object, returning the raw value (array/map/scalar) or nil. Array
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
	// injected by backupGet and materialized by DescribeBackups; keep the tracked value when
	// a detached read reports none. A bare DBClusterId-only id would let Read/Delete lose the
	// owned backup identity ("last wins").
	backupStateId := str("DBClusterId").ValueString()
	backupSetId := str("BackupId").ValueString()
	if backupSetId == "" && prior != nil {
		backupSetId = prior.BackupId.ValueString()
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
		// backup_job_id is the CreateBackup task id: CreateBackup returns it, DescribeBackups
		// never reports it (verified against the live API), so a read keeps the tracked value
		// instead of drifting to null. An imported backup simply has none.
		if m.BackupJobId.IsNull() && !prior.BackupJobId.IsNull() && !prior.BackupJobId.IsUnknown() {
			m.BackupJobId = prior.BackupJobId
		}
		// advanced_data_policies is a configuration-owned nested block, and Plugin
		// Framework nested blocks cannot be Computed. DescribeBackupPolicy returns the
		// system-generated policies next to the configured ones, so taking the read
		// result would make the applied state diverge from the planned block count
		// ("block count changed from 0 to 1"). Keep the configured entries instead and
		// only fill their computed attributes (policy_id, ...) from the matching
		// read-back entry, so the state stays known and later plans are stable.
		// An import has no configuration to align with: there the user-defined policies the
		// API reports are the state, with the system-generated entries filtered out.
		if prior.AdvancedDataPolicies.IsNull() || prior.AdvancedDataPolicies.IsUnknown() {
			m.AdvancedDataPolicies = backupImportedAdvancedPolicies(m.AdvancedDataPolicies)
		} else {
			m.AdvancedDataPolicies = backupAlignAdvancedPolicies(prior.AdvancedDataPolicies, m.AdvancedDataPolicies)
		}
		// Comment is a create-only field (set via CreateBackup) that is not
		// reliably read back from DescribeBackups right after creation: the
		// backup-set list may not include the just-created backup yet, or its
		// entries may carry empty comments, so the response mapping yields null.
		// Preserve the prior/plan value when the read returns null; without a prior
		// value the comment is the empty string, matching both an omitted and an
		// explicitly empty comment instead of oscillating between "" and null.
		if m.Comment.IsNull() {
			if prior.Comment.IsNull() || prior.Comment.IsUnknown() {
				m.Comment = types.StringValue("")
			} else {
				m.Comment = prior.Comment
			}
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

// backupAdvancedLevelEnabled reports whether AdvancedDataPolicies may be sent for a policy
// level. ModifyBackupPolicy documents the parameter as supported only by clusters whose
// BackupPolicyLevel is Advanced; sending it under Normal makes the server reject the whole
// request, including the valid Normal policy fields carried by the same call.
func backupAdvancedLevelEnabled(level string) bool {
	return strings.EqualFold(level, "Advanced")
}

// backupEffectivePolicyLevel returns the policy level that will apply after this operation:
// the configured value, or the value already in state (or read from the cluster) when the
// attribute is omitted and therefore unknown in the plan.
func backupEffectivePolicyLevel(level types.String, fallback interface{}) types.String {
	if !level.IsNull() && !level.IsUnknown() {
		return level
	}
	switch value := fallback.(type) {
	case types.String:
		return value
	case string:
		return types.StringValue(value)
	case nil:
		return types.StringNull()
	}
	return types.StringNull()
}

// backupPreserveBackupPolicyFields keeps the cluster backup-policy fields of a tracked backup
// whose cluster has been released. DescribeDetachedBackups reports the retained backup set
// only, and the cluster policy it used to reference no longer exists, so reporting those
// fields as empty would fabricate a change that cannot be applied.
func backupPreserveBackupPolicyFields(m *backupResourceModel, prior *backupResourceModel) {
	m.AdvancedDataPolicies = prior.AdvancedDataPolicies
	if m.Comment.IsNull() {
		// The released-cluster backup list does not report the backup comment; the configured
		// value must survive a detached read instead of drifting to null.
		m.Comment = prior.Comment
	}
	m.BackupFrequency = prior.BackupFrequency
	m.BackupPolicyLevel = prior.BackupPolicyLevel
	m.BackupRetentionPolicyOnClusterDeletion = prior.BackupRetentionPolicyOnClusterDeletion
	m.DataLevel1BackupFrequency = prior.DataLevel1BackupFrequency
	m.DataLevel1BackupPeriod = prior.DataLevel1BackupPeriod
	m.DataLevel1BackupRetentionPeriod = prior.DataLevel1BackupRetentionPeriod
	m.DataLevel1BackupTime = prior.DataLevel1BackupTime
	m.DataLevel2BackupAnotherRegionRegion = prior.DataLevel2BackupAnotherRegionRegion
	m.DataLevel2BackupAnotherRegionRetentionPeriod = prior.DataLevel2BackupAnotherRegionRetentionPeriod
	m.DataLevel2BackupPeriod = prior.DataLevel2BackupPeriod
	m.DataLevel2BackupRetentionPeriod = prior.DataLevel2BackupRetentionPeriod
	m.PreferredBackupPeriod = prior.PreferredBackupPeriod
	m.PreferredBackupTime = prior.PreferredBackupTime
	m.PreferredNextBackupTime = prior.PreferredNextBackupTime
}
