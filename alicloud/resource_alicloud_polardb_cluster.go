package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBClusterCreate,
		Read:   resourceAlicloudPolarDBClusterRead,
		Update: resourceAlicloudPolarDBClusterUpdate,
		Delete: resourceAlicloudPolarDBClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"db_minor_version": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"8.0.1", "8.0.2"}, false),
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
			},
			"db_node_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"modify_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Upgrade", "Downgrade"}, false),
				Optional:     true,
				Default:      "Upgrade",
			},
			"db_node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(1, 16),
				Computed:     true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"standby_az": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: polardbStandbyAzDiffSuppressFunc,
			},
			"pay_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				Optional:     true,
				Default:      PostPaid,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  RenewNotRenewal,
				ValidateFunc: StringInSlice([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}, false),
				DiffSuppressFunc: polardbPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 6, 12, 24, 36}),
				DiffSuppressFunc: polardbPostPaidAndRenewDiffSuppressFunc,
			},
			"period": {
				Type:             schema.TypeInt,
				ValidateFunc:     IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				DiffSuppressFunc: polardbPostPaidDiffSuppressFunc,
			},
			"db_cluster_ip_array": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_cluster_ip_array_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "default",
						},
						"security_ips": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"modify_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"Cover", "Append", "Delete"}, false),
						},
					},
				},
			},
			"security_ips": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"maintain_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateMaintainTimeRange,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringLenBetween(2, 256),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"collector_status": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Enable", "Disabled"}, false),
				Optional:     true,
				Computed:     true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			"tde_status": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
				Optional:     true,
				Default:      "Disabled",
			},
			"encrypt_new_tables": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"ON", "OFF"}, false),
				Optional:     true,
			},
			"encryption_key": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: polardbTDEAndEnabledDiffSuppressFunc,
			},
			"role_arn": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: polardbTDEAndEnabledDiffSuppressFunc,
				Computed:         true,
			},
			"tde_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"deletion_lock": {
				Type:         schema.TypeInt,
				ValidateFunc: IntInSlice([]int{0, 1}),
				Optional:     true,
			},
			"backup_retention_policy_on_cluster_deletion": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ALL", "LATEST", "NONE"}, false),
			},
			"imci_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ON", "OFF"}, false),
			},
			"sub_category": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Exclusive", "General"}, false),
				Optional:     true,
				Computed:     true,
			},
			"creation_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "Basic", "ArchiveNormal", "NormalMultimaster", "SENormal"}, false),
			},
			"creation_option": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{"Normal", "CloneFromPolarDB", "CloneFromRDS", "MigrationFromRDS", "CreateGdnStandby", "RecoverFromRecyclebin", "UpgradeFromPolarDB"}, false),
				Optional:     true,
				Computed:     true,
			},
			"source_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gdn_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clone_data_point": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PSL5", "PSL4", "ESSDPL0", "ESSDPL1", "ESSDPL2", "ESSDPL3", "ESSDAUTOPL"}, false),
				Computed:     true,
			},
			"provisioned_iops": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: polardbStorageTypeDiffSuppressFunc,
			},
			"storage_pay_type": {
				Type:         schema.TypeString,
				ValidateFunc: StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
			},
			"storage_space": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(20, 100000),
				Computed:     true,
			},
			"hot_standby_cluster": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ON", "OFF", "EQUAL"}, false),
				Computed:     true,
				ForceNew:     true,
			},
			"strict_consistency": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ON", "OFF"}, false),
				Computed:     true,
				ForceNew:     true,
			},
			"serverless_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AgileServerless", "SteadyServerless"}, false),
			},
			"serverless_steady_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"ON", "OFF"}, false),
			},
			"scale_min": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(0, 31),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
			},
			"scale_max": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(0, 32),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
			},
			"scale_ro_num_min": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(0, 15),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
			},
			"scale_ro_num_max": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(0, 15),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
			},
			"allow_shut_down": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"true", "false"}, false),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
			},
			"seconds_until_auto_pause": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     IntBetween(300, 86400),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
			},
			"scale_ap_ro_num_min": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(0, 7),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
				Computed:         true,
			},
			"scale_ap_ro_num_max": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(0, 7),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
				Computed:         true,
			},
			"serverless_rule_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"normal", "flexible"}, false),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
				Computed:         true,
			},
			"serverless_rule_cpu_shrink_threshold": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(10, 100),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
				Computed:         true,
			},
			"serverless_rule_cpu_enlarge_threshold": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(40, 100),
				DiffSuppressFunc: polardbServrelessTypeDiffSuppressFunc,
				Computed:         true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"PROXY", "DB", "ALL"}, false),
			},
			"from_time_service": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"true", "false"}, false),
			},
			"planned_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"planned_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"EXCLUSIVE", "GENERAL"}, false),
				DiffSuppressFunc: polardbProxyTypeDiffSuppressFunc,
			},
			"proxy_class": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: StringInSlice([]string{"polar.maxscale.g2.medium.c", "polar.maxscale.g2.large.c",
					"polar.maxscale.g2.xlarge.c", "polar.maxscale.g2.2xlarge.c", "polar.maxscale.g2.3xlarge.c",
					"polar.maxscale.g2.4xlarge.c", "polar.maxscale.g2.8xlarge.c"}, false),
				DiffSuppressFunc: polardbProxyClassDiffSuppressFunc,
			},
			"loose_polar_log_bin": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"ON", "OFF"}, false),
				DiffSuppressFunc: polardbDiffSuppressFunc,
			},
			"db_node_num": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     IntBetween(1, 16),
				DiffSuppressFunc: polardbProxyTypeDiffSuppressFunc,
			},
			"parameter_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lower_case_table_names": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateFunc:     IntInSlice([]int{0, 1}),
				DiffSuppressFunc: polardbDiffSuppressFunc,
			},
			"default_time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: StringInSlice([]string{"-12:00", "-11:00", "-10:00", "-9:00", "-8:00", "-7:00",
					"-6:00", "-5:00", "-4:00", "-3:00", "-2:00", "-1:00",
					"+0:00", "+1:00", "+2:00", "+3:00", "+4:00", "+5:00",
					"+6:00", "+7:00", "+8:00", "+9:00", "+10:00", "+11:00",
					"+12:00", "+13:00", "SYSTEM"}, false),
				DiffSuppressFunc: polardbDiffSuppressFunc,
			},
			"loose_xengine": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     StringInSlice([]string{"ON", "OFF"}, false),
				DiffSuppressFunc: polardbXengineDiffSuppressFunc,
				Computed:         true,
			},
			"loose_xengine_use_memory_pct": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(10, 90),
				Computed:     true,
			},
			"hot_replica_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ON", "OFF"}, false),
				RequiredWith: []string{"db_node_id"},
			},
			"db_node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_db_revision_version_code": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: polardbAndCreationDiffSuppressFunc,
			},
			"db_revision_version_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"release_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"revision_version_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"revision_version_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"release_note": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"compress_storage": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     StringInSlice([]string{"ON", "OFF"}, false),
				DiffSuppressFunc: polardbCompressStorageDiffSuppressFunc,
			},
		},
	}
}

func resourceAlicloudPolarDBClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request, err := buildPolarDBCreateRequest(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	action := "CreateDBCluster"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_cluster", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["DBClusterId"]))

	// wait cluster status change from Creating to running
	stateConf := BuildStateConf([]string{"Creating", "WarmCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	if v, ok := d.GetOk("db_type"); ok && v.(string) == "MySQL" {
		categoryConf := BuildStateConf([]string{}, []string{"Normal", "Basic", "ArchiveNormal", "NormalMultimaster", "SENormal"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polarDBService.PolarDBClusterCategoryRefreshFunc(d.Id(), []string{}))
		if _, err := categoryConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	allConfig := make(map[string]string)
	if v, ok := d.GetOk("loose_polar_log_bin"); ok {
		allConfig["loose_polar_log_bin"] = fmt.Sprint(v)
	}
	if v, ok := d.GetOk("default_time_zone"); ok {
		allConfig["default_time_zone"] = fmt.Sprint(v)
	}
	if v, ok := d.GetOk("lower_case_table_names"); ok {
		allConfig["lower_case_table_names"] = fmt.Sprint(v)
	}
	// wait instance parameter expect after modifying
	if len(allConfig) > 0 {
		if err := polarDBService.WaitForPolarDBParameter(d.Id(), DefaultLongTimeout, allConfig); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudPolarDBClusterUpdate(d, meta)
}

func resourceAlicloudPolarDBClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	d.Partial(true)

	if !d.IsNewResource() && (d.HasChange("default_time_zone") || d.HasChange("lower_case_table_names") || d.HasChange("loose_polar_log_bin") || d.HasChange("loose_xengine") || d.HasChange("loose_xengine_use_memory_pct")) {
		if err := polarDBService.CreateClusterParamsModifyParameters(d); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("parameters") {
		if err := polarDBService.ModifyParameters(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("parameters")
	}

	if err := polarDBService.setClusterTags(d); err != nil {
		return WrapError(err)
	}

	var err error
	payType := d.Get("pay_type").(string)
	if !d.IsNewResource() && d.HasChange("pay_type") {
		action := "TransformDBClusterPayType"
		requestPayType := convertPolarDBPayTypeUpdateRequest(payType)
		request := map[string]interface{}{
			"RegionId":    client.RegionId,
			"DBClusterId": d.Id(),
			"PayType":     requestPayType,
		}
		if payType == string(PrePaid) {
			period := d.Get("period").(int)
			request["UsedTime"] = strconv.Itoa(period)
			request["Period"] = Month
			if period > 9 {
				request["UsedTime"] = strconv.Itoa(period / 12)
				request["Period"] = Year
			}
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		//wait asynchronously cluster payType
		if err := polarDBService.WaitForPolarDBPayType(d.Id(), requestPayType, DefaultTimeout); err != nil {
			return WrapError(err)
		}
		if payType == string(PrePaid) {
			d.SetPartial("period")
		}
		d.SetPartial("pay_type")
	}

	if (d.Get("pay_type").(string) == string(PrePaid)) &&
		(d.HasChange("renewal_status") || d.HasChange("auto_renew_period")) {
		status := d.Get("renewal_status").(string)
		request := polardb.CreateModifyAutoRenewAttributeRequest()
		request.DBClusterIds = d.Id()
		request.RenewalStatus = status
		if status == string(RenewAutoRenewal) {
			period := d.Get("auto_renew_period").(int)
			request.Duration = strconv.Itoa(period)
			request.PeriodUnit = string(Month)
			if period > 9 {
				request.Duration = strconv.Itoa(period / 12)
				request.PeriodUnit = string(Year)
			}
		}
		//wait asynchronously cluster payType
		if err := polarDBService.WaitForPolarDBPayType(d.Id(), "Prepaid", DefaultLongTimeout); err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPolarDBClient(func(polardbClient *polardb.Client) (interface{}, error) {
				return polardbClient.ModifyAutoRenewAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"-999"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("maintain_time") {
		request := polardb.CreateModifyDBClusterMaintainTimeRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.MaintainTime = d.Get("maintain_time").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterMaintainTime(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("maintain_time")
	}

	if !d.IsNewResource() && d.HasChanges("target_db_revision_version_code") {
		versionInfo, err := polarDBService.DescribeDBClusterVersion(d.Id())
		if err != nil {
			return WrapError(err)
		}
		var isLatestVersion = versionInfo["IsLatestVersion"].(string)
		action := "UpgradeDBClusterVersion"
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		if v, ok := d.GetOk("upgrade_type"); ok {
			request["UpgradeType"] = v
		}
		if v, ok := d.GetOk("from_time_service"); ok {
			fromTimeService, _ := strconv.ParseBool(v.(string))
			request["FromTimeService"] = fromTimeService
		}
		if v, ok := d.GetOk("planned_start_time"); ok {
			request["PlannedStartTime"] = v
		}
		if v, ok := d.GetOk("planned_end_time"); ok {
			request["PlannedEndTime"] = v
		}
		if v, ok := d.GetOk("target_db_revision_version_code"); ok && isLatestVersion == "false" {
			request["TargetDBRevisionVersionCode"] = v
		}
		wait := incrementalWait(3*time.Minute, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"TaskExists"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		fromTimeService := d.Get("from_time_service")
		TargetDBRevisionVersionCode := d.Get("target_db_revision_version_code")
		if strings.EqualFold(fromTimeService.(string), "true") || TargetDBRevisionVersionCode != "" {
			// maxscale upgrade is relatively slow, wait cluster maxscale proxy status from  MinorVersionUpgrading to running
			proxyStatusConf := BuildStateConf([]string{"MinorVersionUpgrading"}, []string{"Running"},
				d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterProxyStateRefreshFunc(d.Id(), []string{""}))
			if _, err := proxyStatusConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			// wait cluster status change from ConfigSwitching to running
			stateConf := BuildStateConf([]string{"MinorVersionUpgrading"}, []string{"Running"},
				d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}
		d.SetPartial("upgrade_type")
		d.SetPartial("from_time_service")
		d.SetPartial("planned_start_time")
		d.SetPartial("planned_end_time")
		d.SetPartial("target_db_revision_version_code")

	}

	if d.HasChange("db_cluster_ip_array") {

		if err := polarDBService.ModifyDBClusterAccessWhitelist(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("db_cluster_ip_array")
	}

	if !d.IsNewResource() && d.HasChange("security_ips") {
		ipList := expandStringList(d.Get("security_ips").(*schema.Set).List())

		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		// default disable connect from outside
		if ipstr == "" {
			ipstr = LOCAL_HOST_IP
		}

		if err := polarDBService.ModifyDBSecurityIps(d.Id(), ipstr); err != nil {
			return WrapError(err)
		}
		d.SetPartial("security_ips")
	}

	if v, ok := d.GetOk("creation_category"); !ok || v.(string) != "Basic" {
		if d.HasChange("db_node_count") {
			cluster, err := polarDBService.DescribePolarDBCluster(d.Id())
			if err != nil {
				return WrapError(err)
			}
			currentDbNodeCount := len(cluster.DBNodes.DBNode)
			expectDbNodeCount := d.Get("db_node_count").(int)
			if expectDbNodeCount > currentDbNodeCount {
				//create node
				expandDbNodes := &[]polardb.CreateDBNodesDBNode{
					{
						TargetClass: cluster.DBNodeClass,
					},
				}
				request := polardb.CreateCreateDBNodesRequest()
				request.RegionId = client.RegionId
				request.DBClusterId = d.Id()
				request.DBNode = expandDbNodes
				if v, ok := d.GetOk("imci_switch"); ok && v.(string) != "" {
					request.ImciSwitch = v.(string)
				}

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
						return polarDBClient.CreateDBNodes(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw, request.RpcRequest, request)
					return nil
				})

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				// wait cluster status change from DBNodeCreating to running
				stateConf := BuildStateConf([]string{"DBNodeCreating"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			} else if expectDbNodeCount < currentDbNodeCount {
				//delete node
				deleteDbNodeId := ""
				for _, dbNode := range cluster.DBNodes.DBNode {
					if dbNode.DBNodeRole == "Reader" {
						deleteDbNodeId = dbNode.DBNodeId
					}
				}
				request := polardb.CreateDeleteDBNodesRequest()
				request.RegionId = client.RegionId
				request.DBClusterId = d.Id()
				request.DBNodeId = &[]string{
					deleteDbNodeId,
				}

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
						return polarDBClient.DeleteDBNodes(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw, request.RpcRequest, request)
					return nil
				})

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}

				stateConf := BuildStateConf([]string{"DBNodeDeleting"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
				if _, err = stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}

	if d.HasChange("collector_status") {
		request := polardb.CreateModifyDBClusterAuditLogCollectorRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.CollectorStatus = d.Get("collector_status").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterAuditLogCollector(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("collector_status")
	}

	if v, ok := d.GetOk("db_type"); ok && v.(string) == "MySQL" {
		if d.HasChange("tde_status") {
			if v, ok := d.GetOk("tde_status"); ok && v.(string) != "Disabled" {
				action := "ModifyDBClusterTDE"
				request := map[string]interface{}{
					"DBClusterId": d.Id(),
					"TDEStatus":   convertPolarDBTdeStatusUpdateRequest(v.(string)),
				}
				if s, ok := d.GetOk("encrypt_new_tables"); ok && s.(string) != "" {
					request["EncryptNewTables"] = s.(string)
				}
				if v, ok := d.GetOk("encryption_key"); ok && v.(string) != "" {
					request["EncryptionKey"] = v.(string)
				}
				if v, ok := d.GetOk("role_arn"); ok && v.(string) != "" {
					request["RoleArn"] = v.(string)
				}
				//retry
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				//wait tde status 'Enabled'

				stateConf := BuildStateConf([]string{}, []string{"Enabled"}, d.Timeout(schema.TimeoutUpdate), 3*time.Minute, polarDBService.PolarDBClusterTDEStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
				d.SetPartial("tde_status")
				d.SetPartial("encrypt_new_tables")
				d.SetPartial("encryption_key")
				d.SetPartial("role_arn")
			}
		}
	}

	if d.HasChange("security_group_ids") {
		securityGroupsList := expandStringList(d.Get("security_group_ids").(*schema.Set).List())
		securityGroupsStr := strings.Join(securityGroupsList[:], COMMA_SEPARATED)

		request := polardb.CreateModifyDBClusterAccessWhitelistRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.WhiteListType = "SecurityGroup"
		request.SecurityGroupIds = securityGroupsStr
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterAccessWhitelist(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("security_group_ids")
	}

	if d.HasChange("deletion_lock") {
		if v, ok := d.GetOk("pay_type"); ok && v.(string) == string(PrePaid) {
			return nil
		}
		action := "ModifyDBClusterDeletion"
		protection := d.Get("deletion_lock").(int)
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
			"Protection":  protection == 1,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ProviderERROR)
		}
		d.SetPartial("deletion_lock")
	}

	if d.HasChange("serverless_steady_switch") {
		// Enable steady state
		if d.HasChanges("scale_min", "scale_max", "scale_ro_num_min", "scale_ro_num_max", "scale_ap_ro_num_min", "scale_ap_ro_num_max") {
			action := "EnableDBClusterServerless"
			request := map[string]interface{}{
				"DBClusterId": d.Id(),
			}
			scaleMin := d.Get("scale_min")
			if scaleMin != nil {
				request["ScaleMin"] = scaleMin
			}
			scaleMax := d.Get("scale_max")
			if scaleMax != nil {
				request["ScaleMax"] = scaleMax
			}
			scaleRoNumMin := d.Get("scale_ro_num_min")
			if scaleRoNumMin != nil {
				request["ScaleRoNumMin"] = scaleRoNumMin
			}
			scaleRoNumMax := d.Get("scale_ro_num_max")
			if scaleRoNumMax != nil {
				request["ScaleRoNumMax"] = scaleRoNumMax
			}
			clusterAttribute, err := polarDBService.DescribePolarDBClusterAttribute(d.Id())
			if err != nil {
				return WrapError(err)
			}
			imciParamterSwitch := false
			for _, nodes := range clusterAttribute.DBNodes {
				if nodes.ImciSwitch == "ON" {
					imciParamterSwitch = true
				}
			}
			if imciParamterSwitch {
				ScaleApRoNumMin := d.Get("scale_ap_ro_num_min")
				if ScaleApRoNumMin != nil {
					scaleApRoNumMin := ScaleApRoNumMin.(int)
					request["ScaleApRoNumMin"] = strconv.Itoa(scaleApRoNumMin)
				}
				ScaleApRoNumMax := d.Get("scale_ap_ro_num_max")
				if ScaleApRoNumMax != nil {
					scaleApRoNumMax := ScaleApRoNumMax.(int)
					request["ScaleApRoNumMax"] = strconv.Itoa(scaleApRoNumMax)
				}
			}

			//retry
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"TaskExists", "OperationDenied.DBClusterStatus"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			// wait cluster status change from ConfigSwitching to running
			stateConf := BuildStateConf([]string{"ConfigSwitching", "Maintaining"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 8*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("scale_min")
			d.SetPartial("scale_max")
			d.SetPartial("scale_ro_num_min")
			d.SetPartial("scale_ro_num_max")
			d.SetPartial("scale_ap_ro_num_min")
			d.SetPartial("scale_ap_ro_num_max")
		}
		// Turn off steady state
		if u, ok := d.GetOk("serverless_steady_switch"); ok {
			switchValue := u.(string)
			cluster, err := polarDBService.DescribePolarDBCluster(d.Id())
			if err != nil {
				return WrapError(err)
			}

			if switchValue == "OFF" && "SteadyServerless" == cluster.ServerlessType {
				action := "DisableDBClusterServerless"
				request := map[string]interface{}{
					"DBClusterId": d.Id(),
				}
				//retry
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				// wait cluster status change from ConfigSwitching to running
				stateConf := BuildStateConf([]string{"ConfigSwitching", "Maintaining"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 8*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}

	if d.HasChange("compress_storage") {
		action := "ModifyDBCluster"
		compressStorage := d.Get("compress_storage").(string)
		request := map[string]interface{}{
			"DBClusterId":     d.Id(),
			"CompressStorage": compressStorage,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ProviderERROR)
		}
		// wait cluster status change from StorageExpanding to running
		stateConf := BuildStateConf([]string{"ConfigSwitching"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 4*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("compress_storage")
	}

	if d.HasChanges("scale_min", "scale_max", "allow_shut_down", "scale_ro_num_min", "scale_ro_num_max", "seconds_until_auto_pause", "scale_ap_ro_num_min", "scale_ap_ro_num_max", "serverless_rule_mode", "serverless_rule_cpu_shrink_threshold", "serverless_rule_cpu_enlarge_threshold") {
		action := "ModifyDBClusterServerlessConf"
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
		}
		scaleMin := d.Get("scale_min")
		if scaleMin != nil {
			request["ScaleMin"] = scaleMin
		}
		scaleMax := d.Get("scale_max")
		if scaleMax != nil {
			request["ScaleMax"] = scaleMax
		}
		scaleRoNumMin := d.Get("scale_ro_num_min")
		if scaleRoNumMin != nil {
			request["ScaleRoNumMin"] = scaleRoNumMin
		}
		scaleRoNumMax := d.Get("scale_ro_num_max")
		if scaleRoNumMax != nil {
			request["ScaleRoNumMax"] = scaleRoNumMax
		}
		if v, ok := d.GetOk("allow_shut_down"); ok && v.(string) != "" {
			request["AllowShutDown"] = v.(string)
		}
		if v, ok := d.GetOk("seconds_until_auto_pause"); ok {
			secondsUntilAutoPause := v.(int)
			request["SecondsUntilAutoPause"] = strconv.Itoa(secondsUntilAutoPause)
		}
		clusterAttribute, err := polarDBService.DescribePolarDBClusterAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}
		imciParamterSwitch := false
		for _, nodes := range clusterAttribute.DBNodes {
			if nodes.ImciSwitch == "ON" {
				imciParamterSwitch = true
			}
		}
		if imciParamterSwitch {
			if v, ok := d.GetOk("scale_ap_ro_num_min"); ok {
				scaleApRoNumMin := v.(int)
				request["ScaleApRoNumMin"] = strconv.Itoa(scaleApRoNumMin)
			}
			if v, ok := d.GetOk("scale_ap_ro_num_max"); ok {
				scaleApRoNumMax := v.(int)
				request["ScaleApRoNumMax"] = strconv.Itoa(scaleApRoNumMax)
			}
		}
		if v, ok := d.GetOk("serverless_rule_mode"); ok && v.(string) != "" {
			request["ServerlessRuleMode"] = v.(string)
		}

		if v, ok := d.GetOk("serverless_rule_cpu_shrink_threshold"); ok {
			serverlessRuleCpuShrinkThreshold := v.(int)
			request["ServerlessRuleCpuShrinkThreshold"] = strconv.Itoa(serverlessRuleCpuShrinkThreshold)
		}
		if v, ok := d.GetOk("serverless_rule_cpu_enlarge_threshold"); ok {
			serverlessRuleCpuEnlargeThreshold := v.(int)
			request["ServerlessRuleCpuEnlargeThreshold"] = strconv.Itoa(serverlessRuleCpuEnlargeThreshold)
		}
		//retry
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait cluster status change from ConfigSwitching to running
		stateConf := BuildStateConf([]string{"ConfigSwitching", "Stopped", "STARTING"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("scale_min")
		d.SetPartial("scale_max")
		d.SetPartial("scale_ro_num_min")
		d.SetPartial("scale_ro_num_max")
		d.SetPartial("allow_shut_down")
		d.SetPartial("seconds_until_auto_pause")
		d.SetPartial("serverless_rule_mode")
		d.SetPartial("serverless_rule_cpu_shrink_threshold")
		d.SetPartial("serverless_rule_cpu_enlarge_threshold")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudPolarDBClusterRead(d, meta)
	}

	if v, ok := d.GetOk("creation_category"); !ok || v.(string) != "Basic" {
		if d.HasChange("db_node_class") {
			request := polardb.CreateModifyDBNodeClassRequest()
			request.RegionId = client.RegionId
			request.DBClusterId = d.Id()
			request.ModifyType = d.Get("modify_type").(string)
			request.DBNodeTargetClass = d.Get("db_node_class").(string)
			if v, ok := d.GetOk("sub_category"); ok && v.(string) != "" {
				request.SubCategory = convertPolarDBSubCategoryUpdateRequest(v.(string))
			}
			//wait asynchronously cluster nodes num the same
			if err := polarDBService.WaitForPolarDBNodeClass(d.Id(), DefaultLongTimeout); err != nil {
				return WrapError(err)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
				raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
					return polarDBClient.ModifyDBNodeClass(request)
				})
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				if err != nil {
					if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"}) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			// wait cluster status change from Creating to running
			stateConf := BuildStateConf([]string{"ClassChanging", "ClassChanged"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{"Deleting"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("db_node_class")
		}
	}

	if d.HasChange("description") {
		request := polardb.CreateModifyDBClusterDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = d.Id()
		request.DBClusterDescription = d.Get("description").(string)

		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBClusterDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("description")
	}

	if d.HasChange("storage_type") {
		action := "ModifyDBClusterStoragePerformance"
		storageType := d.Get("storage_type").(string)
		modifyType := d.Get("modify_type").(string)
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
			"StorageType": storageType,
			"ModifyType":  modifyType,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait cluster status change from ClassChanging/ConfigSwitching to running
		stateConf := BuildStateConf([]string{"ClassChanging", "ConfigSwitching"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 4*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("storage_type")
	}

	if d.HasChange("provisioned_iops") {
		action := "ModifyDBClusterStoragePerformance"
		storageType := d.Get("storage_type").(string)
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
			"StorageType": storageType,
		}
		if v, ok := d.GetOk("provisioned_iops"); ok && v.(string) != "" {
			request["ProvisionedIops"] = v
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait cluster status change from ClassChanging/ConfigSwitching to running
		stateConf := BuildStateConf([]string{"ClassChanging", "ConfigSwitching"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 4*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("provisioned_iops")
	}

	if d.HasChange("storage_space") {
		action := "ModifyDBClusterStorageSpace"
		storageSpace := d.Get("storage_space").(int)
		request := map[string]interface{}{
			"DBClusterId":  d.Id(),
			"StorageSpace": storageSpace,
		}
		if v, ok := d.GetOk("planned_start_time"); ok && v.(string) != "" {
			request["PlannedStartTime"] = v
		}
		if v, ok := d.GetOk("planned_end_time"); ok && v.(string) != "" {
			request["PlannedEndTime"] = v
		}
		if v, ok := d.GetOk("sub_category"); ok && v.(string) != "" {
			request["SubCategory"] = convertPolarDBSubCategoryUpdateRequest(v.(string))
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) {
				return nil
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ProviderERROR)
		}
		// wait cluster status change from StorageExpanding to running
		stateConf := BuildStateConf([]string{"StorageExpanding"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 4*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("storage_space")
	}

	if d.HasChange("hot_replica_mode") {
		dbNodeIdIndex := ""
		if v, ok := d.GetOk("db_node_id"); ok && v.(string) != "" {
			if len(v.(string)) > 2 {
				dbNodeIdIndex = v.(string)
			} else {
				clusterAttribute, err := polarDBService.DescribePolarDBClusterAttribute(d.Id())
				if err != nil {
					return WrapError(err)
				}
				index := formatInt(v)
				if len(clusterAttribute.DBNodes) <= index {
					return WrapError(Error("The specified db_node_id exceeded DBNodes range."))
				}
				dbNodeIdIndex = clusterAttribute.DBNodes[index].DBNodeId
			}
		}
		if v, ok := d.GetOk("db_type"); ok && v.(string) == "MySQL" {
			action := "ModifyDBNodeHotReplicaMode"
			hotReplicaMode := d.Get("hot_replica_mode").(string)
			request := map[string]interface{}{
				"DBClusterId":    d.Id(),
				"HotReplicaMode": hotReplicaMode,
				"DBNodeId":       dbNodeIdIndex,
			}
			//retry
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			//wait tde status 'Running'
			stateConf := BuildStateConf([]string{"RoleSwitching"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			d.SetPartial("hot_replica_mode")
			d.SetPartial("db_node_id")
		}
	}

	if d.HasChange("standby_az") {
		action := "ModifyDBClusterPrimaryZone"
		standbyAz := d.Get("standby_az").(string)
		request := map[string]interface{}{
			"DBClusterId": d.Id(),
			"ZoneId":      standbyAz,
			"ZoneType":    "Standby",
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// wait cluster status change from ClassChanging/ConfigSwitching to running
		stateConf := BuildStateConf([]string{"StandbyTransing"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 4*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("standby_az")
	}

	if d.HasChange("resource_group_id") {
		action := "ModifyDBClusterResourceGroup"
		newResourceGroup := d.Get("resource_group_id").(string)
		request := map[string]interface{}{
			"DBClusterId":        d.Id(),
			"NewResourceGroupId": newResourceGroup,
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				addDebug(action, response, request)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBClusterRead(d, meta)
}

func resourceAlicloudPolarDBClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	clusterAttribute, err := polarDBService.DescribePolarDBClusterAttribute(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	cluster, err := polarDBService.DescribePolarDBCluster(d.Id())
	if err != nil {
		return WrapError(err)
	}
	whiteList, err := polarDBService.DescribeDBClusterAccessWhitelist(d.Id())
	if err != nil {
		return WrapError(err)
	}
	defaultSecurityIps := make([]string, 0)
	dbClusterIPArrays := make([]map[string]interface{}, 0)
	inDBClusterIPArrays := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("db_cluster_ip_array"); ok {
		for _, e := range v.(*schema.Set).List() {
			inDBClusterIPArrays = append(inDBClusterIPArrays, e.(map[string]interface{}))
		}
	}
	for _, white := range whiteList.Items.DBClusterIPArray {
		if white.DBClusterIPArrayAttribute == "hidden" {
			continue
		}
		// Judge whether input parameters are passed into modify_mode, if there is a modify_mode parameter is based on
		// db_cluster_ip_array_name、security_ips determines whether the whitelist is the same and assigns a local modify_mode
		modifyMode := ""
		for _, temp := range inDBClusterIPArrays {
			if temp["db_cluster_ip_array_name"] == nil || temp["security_ips"] == nil {
				continue
			}
			if temp["db_cluster_ip_array_name"] == white.DBClusterIPArrayName &&
				arrValueEqual(convertPolarDBIpsSetListToString(temp["security_ips"].(*schema.Set)), convertPolarDBIpsSetToString(white.SecurityIps)) {
				if temp["modify_mode"] != nil {
					modifyMode = temp["modify_mode"].(string)
				}
			}
		}
		clusterIdItem := map[string]interface{}{
			"db_cluster_ip_array_name": white.DBClusterIPArrayName,
			"security_ips":             convertPolarDBIpsSetToString(white.SecurityIps),
		}
		if modifyMode != "" {
			clusterIdItem["modify_mode"] = modifyMode
		}
		dbClusterIPArrays = append(dbClusterIPArrays, clusterIdItem)
		if white.DBClusterIPArrayName == "default" {
			defaultSecurityIps = convertPolarDBIpsSetToString(white.SecurityIps)
		}
	}
	d.Set("db_cluster_ip_array", dbClusterIPArrays)
	d.Set("security_ips", defaultSecurityIps)
	//describe endpoints
	var connectionString, port string
	endpoints, err := polarDBService.DescribePolarDBInstanceNetInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	for _, endpoint := range endpoints {
		if endpoint.EndpointType == "Cluster" {
			for _, item := range endpoint.AddressItems {
				if item.NetType == "Private" {
					connectionString = item.ConnectionString
					port = item.Port
					break
				}
			}
		}
	}
	if connectionString == "" {
		// Compatible with the new logic of cloud products, if there is a cluster address, return the connection_string and port of the cluster address, and if not, return the primary address
		for _, endpoint := range endpoints {
			if endpoint.EndpointType == "Primary" {
				for _, item := range endpoint.AddressItems {
					if item.NetType == "Private" {
						connectionString = item.ConnectionString
						port = item.Port
						break
					}
				}
			}
		}
	}
	d.Set("connection_string", connectionString)
	d.Set("port", port)

	d.Set("vswitch_id", clusterAttribute.VSwitchId)
	d.Set("pay_type", getChargeType(clusterAttribute.PayType))
	d.Set("id", clusterAttribute.DBClusterId)
	d.Set("description", clusterAttribute.DBClusterDescription)
	d.Set("db_type", clusterAttribute.DBType)
	d.Set("db_version", clusterAttribute.DBVersion)
	d.Set("maintain_time", clusterAttribute.MaintainTime)
	// Only compare the main availability zone, and randomly allocate backup availability zones in the background
	if len(clusterAttribute.DBNodes) > 0 {
		d.Set("zone_id", clusterAttribute.DBNodes[0].ZoneId)
	}
	d.Set("db_node_class", cluster.DBNodeClass)
	// db_node_count normal nodes, excluding backend generated sensitive nodes
	dbNodeCount := len(clusterAttribute.DBNodes)
	if clusterAttribute.ServerlessType == "SteadyServerless" {
		for _, nodes := range clusterAttribute.DBNodes {
			if nodes.ServerlessType != "SteadyServerless" {
				dbNodeCount--
			}
		}
	}
	d.Set("db_node_count", dbNodeCount)
	d.Set("resource_group_id", clusterAttribute.ResourceGroupId)
	d.Set("deletion_lock", clusterAttribute.DeletionLock)
	d.Set("creation_category", clusterAttribute.Category)
	d.Set("vpc_id", clusterAttribute.VPCId)
	d.Set("status", clusterAttribute.DBClusterStatus)
	d.Set("create_time", clusterAttribute.CreationTime)

	tags, err := polarDBService.DescribeTags(d.Id(), "cluster")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", polarDBService.tagsToMap(tags))

	if clusterAttribute.PayType == string(Prepaid) {
		clusterAutoRenew, err := polarDBService.DescribePolarDBAutoRenewAttribute(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return WrapError(err)
		}
		renewPeriod := 1
		if clusterAutoRenew != nil {
			renewPeriod = clusterAutoRenew.Duration
		}
		if clusterAutoRenew != nil && clusterAutoRenew.PeriodUnit == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("auto_renew_period", renewPeriod)
		//period, err := computePeriodByUnit(clusterAttribute.CreationTime, clusterAttribute.ExpireTime, d.Get("period").(int), "Month")
		//if err != nil {
		//	return WrapError(err)
		//}
		//d.Set("period", period)
		d.Set("renewal_status", clusterAutoRenew.RenewalStatus)
	}

	if err = polarDBService.RefreshParameters(d); err != nil {
		return WrapError(err)
	}

	clusterCollectStatus, err := polarDBService.DescribeDBAuditLogCollectorStatus(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("collector_status", clusterCollectStatus)

	clusterTDEStatus, err := polarDBService.DescribeDBClusterTDE(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("tde_status", clusterTDEStatus["TDEStatus"])
	d.Set("encrypt_new_tables", clusterTDEStatus["EncryptNewTables"])
	d.Set("encryption_key", clusterTDEStatus["EncryptionKey"])
	d.Set("tde_region", clusterTDEStatus["TDERegion"])
	tdeRegion := ""
	if v, ok := clusterTDEStatus["TDERegion"]; ok {
		tdeRegion = fmt.Sprint(v)
	}
	// Check if the current TDE is enabled, and then call the interface if it is enabled
	if "Disabled" != clusterTDEStatus["TDEStatus"].(string) {
		roleArnObj, err := polarDBService.CheckKMSAuthorized(d.Id(), tdeRegion)
		if err != nil {
			return WrapError(err)
		}
		d.Set("role_arn", roleArnObj["RoleArn"])
	}
	securityGroups, err := polarDBService.DescribeDBSecurityGroups(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_group_ids", securityGroups)
	clusterInfo, err := polarDBService.DescribeDBClusterAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}
	if clusterInfo["StorageType"] != nil {
		d.Set("storage_type", convertPolarDBStorageTypeDescribeRequest(clusterInfo["StorageType"].(string)))
	}
	if clusterInfo["ProvisionedIops"] != nil {
		d.Set("provisioned_iops", clusterInfo["ProvisionedIops"].(json.Number).String())
	}
	if clusterInfo["StorageSpace"] != nil {
		resultStorageSpace, _ := clusterInfo["StorageSpace"].(json.Number).Int64()
		var storageSpace = resultStorageSpace / 1024 / 1024 / 1024
		d.Set("storage_space", storageSpace)
	}
	if clusterInfo["StoragePayType"] != nil {
		d.Set("storage_pay_type", getChargeType(clusterInfo["StoragePayType"].(string)))
	}
	if clusterInfo["ServerlessType"] != nil {
		d.Set("serverless_type", clusterInfo["ServerlessType"].(string))
		serverlessInfo, err := polarDBService.DescribeDBClusterServerlessConfig(d.Id())
		if err != nil {
			return WrapError(err)
		}
		d.Set("scale_min", formatInt(serverlessInfo["ScaleMin"]))
		d.Set("scale_max", formatInt(serverlessInfo["ScaleMax"]))
		d.Set("scale_ro_num_min", formatInt(serverlessInfo["ScaleRoNumMin"]))
		d.Set("scale_ro_num_max", formatInt(serverlessInfo["ScaleRoNumMax"]))
		d.Set("allow_shut_down", serverlessInfo["AllowShutDown"])
		d.Set("seconds_until_auto_pause", formatInt(serverlessInfo["SecondsUntilAutoPause"]))
		d.Set("scale_ap_ro_num_min", formatInt(serverlessInfo["ScaleApRoNumMin"]))
		d.Set("scale_ap_ro_num_max", formatInt(serverlessInfo["ScaleApRoNumMax"]))
		d.Set("serverless_rule_mode", serverlessInfo["ServerlessRuleMode"])
		d.Set("serverless_rule_cpu_shrink_threshold", formatInt(serverlessInfo["ServerlessRuleCpuShrinkThreshold"]))
		d.Set("serverless_rule_cpu_enlarge_threshold", formatInt(serverlessInfo["ServerlessRuleCpuEnlargeThreshold"]))
		serverlessSwitch := ""
		if v, ok := serverlessInfo["Switch"]; ok {
			serverlessSwitch = fmt.Sprint(v)
			d.Set("serverless_steady_switch", convertPolarDBServerlessSteadySwitchReadResponse(serverlessSwitch))
		}
	}
	if v, ok := d.GetOk("db_node_id"); ok && v.(string) != "" {
		dbNodeIdIndex := v.(string)
		if len(dbNodeIdIndex) > 2 {
			for _, nodes := range clusterAttribute.DBNodes {
				if nodes.DBNodeId == dbNodeIdIndex {
					d.Set("db_node_id", nodes.DBNodeId)
					d.Set("hot_replica_mode", nodes.HotReplicaMode)
				}
			}
		} else {
			d.Set("db_node_id", dbNodeIdIndex)
			d.Set("hot_replica_mode", clusterAttribute.DBNodes[formatInt(dbNodeIdIndex)].HotReplicaMode)
		}
	}
	creationCategory, categoryOk := d.GetOk("creation_category")
	DBRevisionVersionList := make([]map[string]interface{}, 0)
	if dbType, ok := d.GetOk("db_type"); ok && dbType.(string) == "MySQL" && (creationCategory == "Normal" || creationCategory == "NormalMultimaster" || !categoryOk) {
		availableVersion, errs := polarDBService.DescribeDBClusterAvailableVersion(d.Id())
		if errs != nil {
			return WrapError(errs)
		}
		for _, versionList := range availableVersion.DBRevisionVersionList {
			versionListItem := map[string]interface{}{
				"release_type":          versionList.ReleaseType,
				"revision_version_name": versionList.RevisionVersionName,
				"revision_version_code": versionList.RevisionVersionCode,
				"release_note":          versionList.ReleaseNote,
			}
			DBRevisionVersionList = append(DBRevisionVersionList, versionListItem)
		}
	}
	d.Set("db_revision_version_list", DBRevisionVersionList)
	d.Set("target_db_revision_version_code", d.Get("target_db_revision_version_code"))
	d.Set("compress_storage", clusterAttribute.CompressStorageMode)
	d.Set("hot_standby_cluster", convertPolarDBHotStandbyClusterStatusReadResponse(clusterAttribute.HotStandbyCluster))
	d.Set("strict_consistency", clusterAttribute.StrictConsistency)

	standbyAz, err := polarDBService.DescribeDBClusterStandbyAz(d, DefaultTimeoutMedium)
	if err != nil {
		return WrapError(err)
	}
	d.Set("standby_az", standbyAz)

	versionInfo, err := polarDBService.DescribeDBClusterVersion(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("db_minor_version", versionInfo["DBMinorVersion"])

	return nil
}

func resourceAlicloudPolarDBClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}

	cluster, err := polarDBService.DescribePolarDBClusterAttribute(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	// Pre paid cluster can not be release.
	if PayType(cluster.PayType) == Prepaid {
		return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
	}

	request := polardb.CreateDeleteDBClusterRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Id()
	if v, ok := d.GetOk("backup_retention_policy_on_cluster_deletion"); ok && v.(string) != "" {
		request.BackupRetentionPolicyOnClusterDeletion = v.(string)
	}
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.DeleteDBCluster(request)
		})

		if err != nil && !NotFoundError(err) {
			if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus", "OperationDenied.PolarDBClusterStatus", "OperationDenied.ReadPolarDBClusterStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Creating", "Running", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, polarDBService.PolarDBClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func buildPolarDBCreateRequest(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := map[string]interface{}{
		"RegionId":             client.RegionId,
		"DBType":               Trim(d.Get("db_type").(string)),
		"DBVersion":            Trim(d.Get("db_version").(string)),
		"DBNodeClass":          d.Get("db_node_class").(string),
		"DBClusterDescription": d.Get("description").(string),
		"ClientToken":          buildClientToken("CreateDBCluster"),
		"CreationCategory":     d.Get("creation_category").(string),
		"CloneDataPoint":       d.Get("clone_data_point").(string),
	}

	v, exist := d.GetOk("creation_option")
	db, ok := d.GetOk("db_type")
	dbv, dbvok := d.GetOk("db_version")

	if exist && v.(string) == "CloneFromPolarDB" {
		request["SourceResourceId"] = d.Get("source_resource_id").(string)
		request["CreationOption"] = d.Get("creation_option").(string)
	}

	if exist && v.(string) == "CloneFromRDS" {
		request["CloneDataPoint"] = "LATEST"
	}

	if exist && v.(string) == "CreateGdnStandby" {
		if ok && db.(string) == "MySQL" {
			if dbvok && dbv.(string) == "8.0" {
				request["CreationOption"] = d.Get("creation_option").(string)
				request["GDNId"] = d.Get("gdn_id").(string)
			}
		}
	}

	if exist && v.(string) == "CloneFromRDS" {
		if ok && db.(string) == "MySQL" {
			request["CreationOption"] = d.Get("creation_option").(string)
			request["SourceResourceId"] = d.Get("source_resource_id").(string)
		}
	}

	if exist && v.(string) == "MigrationFromRDS" {
		request["CreationOption"] = d.Get("creation_option").(string)
		request["SourceResourceId"] = d.Get("source_resource_id").(string)
	}

	if exist && v.(string) == "UpgradeFromPolarDB" {
		request["CreationOption"] = d.Get("creation_option").(string)
		request["SourceResourceId"] = d.Get("source_resource_id").(string)
	}

	if v, ok := d.GetOk("storage_type"); ok && v.(string) != "" {
		request["StorageType"] = d.Get("storage_type").(string)
	}
	if v, ok := d.GetOk("db_minor_version"); ok && v.(string) != "" {
		request["DBMinorVersion"] = d.Get("db_minor_version").(string)
	}
	if v, ok := d.GetOk("provisioned_iops"); ok && v.(string) != "" {
		request["ProvisionedIops"], _ = strconv.ParseInt(v.(string), 10, 64)
	}
	if v, ok := d.GetOk("storage_space"); ok && v.(int) != 0 {
		request["StorageSpace"] = d.Get("storage_space").(int)
	}
	if v, ok := d.GetOk("storage_pay_type"); ok && v.(string) != "" {
		if v.(string) == string(PrePaid) {
			request["StoragePayType"] = string(Prepaid)
		}
		if v.(string) == string(PostPaid) {
			request["StoragePayType"] = string(Postpaid)
		}
	}

	if v, ok := d.GetOk("hot_standby_cluster"); ok && v.(string) != "" {
		request["HotStandbyCluster"] = d.Get("hot_standby_cluster").(string)
	}

	if v, ok := d.GetOk("creation_category"); ok && v.(string) != "" {
		if v.(string) == "SENormal" {
			if w, ok := d.GetOk("hot_standby_cluster"); ok && w.(string) != "" {
				if w.(string) == "ON" {
					// 标准版：STANDBY=开启；OFF=关闭；集群版：ON=开启；OFF=关闭；
					request["HotStandbyCluster"] = "STANDBY"
				}
			}

		}
	}

	if v, ok := d.GetOk("strict_consistency"); ok && v.(string) != "" {
		request["StrictConsistency"] = d.Get("strict_consistency").(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request["ResourceGroupId"] = v.(string)
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request["ZoneId"] = Trim(zone.(string))
	}

	if standbyAz, ok := d.GetOk("standby_az"); ok && Trim(standbyAz.(string)) != "" {
		request["StandbyAZ"] = Trim(standbyAz.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v.(string)
	}

	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v.(string)
	}

	if request["VSwitchId"] != nil {
		request["ClusterNetworkType"] = strings.ToUpper(string(Vpc))
		if request["ZoneId"] == nil || request["VPCId"] == nil {
			// check vswitchId in zone
			vsw, err := vpcService.DescribeVSwitch(request["VSwitchId"].(string))
			if err != nil {
				return nil, WrapError(err)
			}

			if v, ok := request["ZoneId"].(string); !ok || v == "" {
				request["ZoneId"] = vsw.ZoneId
			} else if request["ZoneId"] != vsw.ZoneId {
				return nil, WrapError(Error("The specified vswitch %s isn't in the zone %s.", vsw.VSwitchId, request["ZoneId"]))
			}

			if v, ok := request["VPCId"].(string); !ok || v == "" {
				request["VPCId"] = vsw.VpcId
			}
		}
	}

	payType := Trim(d.Get("pay_type").(string))
	request["PayType"] = string(Postpaid)
	if payType == string(PrePaid) {
		request["PayType"] = string(Prepaid)
	}
	if PayType(request["PayType"].(string)) == Prepaid {
		period := d.Get("period").(int)
		request["UsedTime"] = strconv.Itoa(period)
		request["Period"] = string(Month)
		if period > 9 {
			request["UsedTime"] = strconv.Itoa(period / 12)
			request["Period"] = string(Year)
		}
		if d.Get("renewal_status").(string) != string(RenewNotRenewal) {
			request["AutoRenew"] = requests.Boolean(strconv.FormatBool(true))
		} else {
			request["AutoRenew"] = requests.Boolean(strconv.FormatBool(false))
		}
	}

	request["TDEStatus"] = requests.NewBoolean(convertPolarDBTdeStatusCreateRequest(d.Get("tde_status").(string)))

	if v, ok := d.GetOk("serverless_type"); ok && v.(string) == "AgileServerless" {
		request["ServerlessType"] = d.Get("serverless_type").(string)

		if v, ok := d.GetOk("scale_min"); ok {
			scaleMin := v.(int)
			request["ScaleMin"] = strconv.Itoa(scaleMin)
		}
		if v, ok := d.GetOk("scale_max"); ok {
			scaleMax := v.(int)
			request["ScaleMax"] = strconv.Itoa(scaleMax)
		}
		if v, ok := d.GetOk("allow_shut_down"); ok && v.(string) != "" {
			request["AllowShutDown"] = d.Get("allow_shut_down").(string)
		}
		scaleRoNumMin := d.Get("scale_ro_num_min")
		if scaleRoNumMin != nil {
			request["ScaleRoNumMin"] = scaleRoNumMin
		}
		scaleRoNumMax := d.Get("scale_ro_num_max")
		if scaleRoNumMax != nil {
			request["ScaleRoNumMax"] = scaleRoNumMax
		}

	}

	if v, ok := d.GetOk("proxy_type"); ok {
		request["ProxyType"] = v.(string)
	}

	if v, ok := d.GetOk("proxy_class"); ok {
		request["ProxyClass"] = v.(string)
	}

	if v, ok := d.GetOk("loose_polar_log_bin"); ok {
		request["LoosePolarLogBin"] = v.(string)
	}

	if v, ok := d.GetOk("db_node_num"); ok {
		request["DBNodeNum"] = v.(int)
	}

	if v, ok := d.GetOk("parameter_group_id"); ok {
		request["ParameterGroupId"] = v.(string)
	}
	LowerCaseTableNames := d.Get("lower_case_table_names")
	if LowerCaseTableNames != nil {
		request["LowerCaseTableNames"] = LowerCaseTableNames.(int)
	}

	if v, ok := d.GetOk("default_time_zone"); ok {
		request["DefaultTimeZone"] = v.(string)
	}

	if v, ok := d.GetOk("security_ips"); ok {
		ipList := expandStringList(v.(*schema.Set).List())
		ipstr := strings.Join(ipList[:], COMMA_SEPARATED)
		request["SecurityIPList"] = ipstr
	}

	if v, ok := d.GetOk("loose_xengine"); ok {
		request["LooseXEngine"] = v.(string)
		if l, ok := d.GetOk("loose_xengine_use_memory_pct"); ok && v.(string) == "ON" {
			looseXEngineUseMemoryPct := l.(int)
			request["LooseXEngineUseMemoryPct"] = strconv.Itoa(looseXEngineUseMemoryPct)
		}
	}

	return request, nil
}

func convertPolarDBTdeStatusCreateRequest(source string) bool {
	switch source {
	case "Enabled":
		return true
	}
	return false
}

func convertPolarDBTdeStatusUpdateRequest(source string) string {
	switch source {
	case "Enabled":
		return "Enable"
	}
	return "Disable"
}

func convertPolarDBPayTypeUpdateRequest(source string) string {
	switch source {
	case "PrePaid":
		return "Prepaid"
	}
	return "Postpaid"
}
func convertPolarDBSubCategoryUpdateRequest(source string) string {
	switch source {
	case "Exclusive":
		return "normal_exclusive"
	}
	return "normal_general"
}
func convertPolarDBStorageTypeDescribeRequest(source string) string {
	switch source {
	case "HighPerformance":
		return "PSL5"
	case "Standard":
		return "PSL4"
	case "essdpl0":
		return "ESSDPL0"
	case "essdpl1":
		return "ESSDPL1"
	case "essdpl2":
		return "ESSDPL2"
	case "essdpl3":
		return "ESSDPL3"
	case "essdautopl":
		return "ESSDAUTOPL"
	}
	return source
}

func convertPolarDBServerlessSteadySwitchReadResponse(source string) string {
	switch source {
	case "1":
		return "ON"
	}
	return "OFF"
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func validateMaintainTimeRange(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	re := regexp.MustCompile(`^([01][0-9]|2[0-3]):00Z-([01][0-9]|2[0-3]):00Z$`)
	if !re.MatchString(value) {
		errors = append(errors, fmt.Errorf("The maintain_time must be on the hour. Example: 16:00Z-17:00Z"))
	} else {
		start, _ := strconv.Atoi(re.FindStringSubmatch(value)[1])
		end, _ := strconv.Atoi(re.FindStringSubmatch(value)[2])
		if end-start != 1 {
			errors = append(errors, fmt.Errorf(
				"The maintain_time interval must be 1 hour."))
		}
	}
	return
}

func convertPolarDBHotStandbyClusterStatusReadResponse(source string) string {
	switch source {
	case "StandbyClusterON":
		return "ON"
	case "equal":
		return "EQUAL"
	}
	return "OFF"
}
