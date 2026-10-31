package alicloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudAdbDbClusterLakeVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAdbDbClusterLakeVersionCreate,
		Read:   resourceAliCloudAdbDbClusterLakeVersionRead,
		Update: resourceAliCloudAdbDbClusterLakeVersionUpdate,
		Delete: resourceAliCloudAdbDbClusterLakeVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customizeDiffAdbDbClusterLakeVersionAutoRenewal,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(72 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"5.0"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"secondary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secondary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_form": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"product_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"compute_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reserved_node_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reserved_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"disk_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"kms_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_ips": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"db_cluster_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			},
			"enable_default_resource_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_db_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"backup", "timepoint"}, false),
			},
			"restore_to_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commodity_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"VPC"}, false),
			},
			"audit_log_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"used_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"essd_cache_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"auto_renewal_period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Year", "Month"}, false),
			},
			"switch_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
			},
			"db_cluster_ip_array_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "Normal", "NotRenewal"}, false),
			},
			"maintain_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAdbDbClusterLakeVersionMaintainTime,
			},
			"auto_renewal_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 11),
			},
			"enable_compaction_service": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"modify_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Cover", "Append", "Delete"}, false),
			},
			"db_cluster_ip_array_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_essd_cache": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_lake_cache": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"XIHE", "SPARK"}, false),
			},
		},
	}
}

func resourceAliCloudAdbDbClusterLakeVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	var response map[string]interface{}
	var err error
	action := "CreateDBCluster"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["DBClusterVersion"] = d.Get("db_cluster_version")
	request["VPCId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["ZoneId"] = d.Get("zone_id")
	request["PayType"] = convertAdbDbClusterLakeVersionPaymentTypeRequest(d.Get("payment_type"))

	if v, ok := d.GetOk("secondary_vswitch_id"); ok {
		request["SecondaryVSwitchId"] = v
	}

	if v, ok := d.GetOk("secondary_zone_id"); ok {
		request["SecondaryZoneId"] = v
	}

	if v, ok := d.GetOk("product_form"); ok {
		request["ProductForm"] = v
	}

	if v, ok := d.GetOk("product_version"); ok {
		request["ProductVersion"] = v
	}

	if v, ok := d.GetOk("compute_resource"); ok {
		request["ComputeResource"] = v
	}

	if v, ok := d.GetOk("storage_resource"); ok {
		request["StorageResource"] = v
	}

	if v, ok := d.GetOk("reserved_node_size"); ok {
		request["ReservedNodeSize"] = v
	}

	if v, ok := d.GetOkExists("reserved_node_count"); ok {
		request["ReservedNodeCount"] = v
	}

	if v, ok := d.GetOkExists("disk_encryption"); ok {
		request["DiskEncryption"] = v
	}

	if v, ok := d.GetOk("kms_id"); ok {
		request["KmsId"] = v
	}

	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("db_cluster_network_type"); ok {
		request["DBClusterNetworkType"] = v
	}

	if request["PayType"] == string(Prepaid) {
		if v, ok := d.GetOkExists("period"); ok {
			usedTime := v.(int)
			request["UsedTime"] = strconv.Itoa(usedTime)
			request["Period"] = string(Month)
			if usedTime > 9 {
				request["UsedTime"] = strconv.Itoa(usedTime / 12)
				request["Period"] = string(Year)
			}
		}
	}

	if v, ok := d.GetOkExists("enable_default_resource_group"); ok {
		request["EnableDefaultResourcePool"] = v
	}

	if v, ok := d.GetOk("source_db_cluster_id"); ok {
		request["SourceDbClusterId"] = v
	}

	if v, ok := d.GetOk("backup_set_id"); ok {
		request["BackupSetId"] = v
	}

	if v, ok := d.GetOk("restore_type"); ok {
		request["RestoreType"] = v
	}

	if v, ok := d.GetOk("restore_to_time"); ok {
		request["RestoreToTime"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_db_cluster_lake_version", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.DBClusterId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_adb_db_cluster_lake_version")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAdbDbClusterLakeVersionUpdate(d, meta)
}

func resourceAliCloudAdbDbClusterLakeVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	object, err := adbService.DescribeAdbDbClusterLakeVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version adbService.DescribeAdbDbClusterLakeVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_version", object["DBVersion"])
	d.Set("vpc_id", object["VPCId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("compute_resource", object["ComputeResource"])
	d.Set("storage_resource", object["StorageResource"])
	d.Set("payment_type", convertAdbDbClusterLakeVersionPaymentTypeResponse(object["PayType"]))
	d.Set("secondary_vswitch_id", object["SecondaryVSwitchId"])
	d.Set("secondary_zone_id", object["SecondaryZoneId"])
	d.Set("product_form", object["ProductForm"])
	d.Set("product_version", object["ProductVersion"])
	d.Set("reserved_node_size", object["ReservedNodeSize"])
	d.Set("reserved_node_count", object["ReservedNodeCount"])
	d.Set("disk_encryption", object["DiskEncryption"])
	d.Set("kms_id", object["KmsId"])
	d.Set("db_cluster_description", object["DBClusterDescription"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("commodity_code", object["CommodityCode"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("port", object["Port"])
	d.Set("lock_mode", object["LockMode"])
	d.Set("lock_reason", object["LockReason"])
	d.Set("expired", object["Expired"])
	d.Set("create_time", object["CreationTime"])
	d.Set("expire_time", object["ExpireTime"])
	d.Set("status", object["DBClusterStatus"])
	d.Set("db_cluster_network_type", object["DBClusterNetworkType"])
	d.Set("maintain_time", object["MaintainTime"])

	// used_time is derived from period and only meaningful for a subscription cluster.
	if fmt.Sprint(d.Get("payment_type")) == "Subscription" {
		if v, ok := d.GetOk("period"); ok {
			d.Set("used_time", v)
		} else {
			d.Set("used_time", 0)
		}
	} else {
		d.Set("used_time", 0)
	}

	if essdCacheObject, err := adbService.DescribeAdbDbClusterLakeVersionEssdCacheConfig(d.Id()); err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
		log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version DescribeEssdCacheConfig not found: %s", err)
	} else {
		d.Set("enable_essd_cache", essdCacheObject["EnableEssdCache"])
		d.Set("essd_cache_size", essdCacheObject["EssdCacheSize"])
	}

	if lakeCacheObject, err := adbService.DescribeAdbDbClusterLakeVersionLakeCacheSize(d.Id()); err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
		log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version DescribeLakeCacheSize not found: %s", err)
	} else {
		d.Set("enable_lake_cache", lakeCacheObject["EnableLakeCache"])
		d.Set("capacity", lakeCacheObject["Capacity"])
	}

	if compactionObject, err := adbService.DescribeAdbDbClusterLakeVersionCompactionServiceSwitch(d.Id()); err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
		log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version DescribeCompactionServiceSwitch not found: %s", err)
	} else {
		d.Set("enable_compaction_service", compactionObject["EnableCompactionService"])
	}

	if autoRenewObject, err := adbService.DescribeAdbDbClusterLakeVersionAutoRenewalAttribute(d.Id()); err != nil {
		// PayAsYouGo clusters are not listed by this API, which is the only
		// NotFound case allowed to skip a partial read.
		if !NotFoundError(err) {
			return WrapError(err)
		}
		log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version DescribeAutoRenewalAttribute not found: %s", err)
	} else {
		d.Set("auto_renewal_status", autoRenewObject["AutoRenewalStatus"])
		d.Set("auto_renewal_period", autoRenewObject["AutoRenewalPeriod"])
		d.Set("auto_renewal_period_unit", autoRenewObject["AutoRenewalPeriodUnit"])
	}

	if auditLogObject, err := adbService.DescribeAdbDbClusterLakeVersionAuditLogConfig(d.Id()); err != nil {
		if !NotFoundError(err) {
			return WrapError(err)
		}
		log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version DescribeAuditLogConfig not found: %s", err)
	} else {
		d.Set("audit_log_status", auditLogObject["AuditLogStatus"])
	}

	ipArrayName := "default"
	if v, ok := d.GetOk("db_cluster_ip_array_name"); ok {
		ipArrayName = v.(string)
	}
	if ipArrayObject, err := adbService.DescribeClusterAccessWhiteList(d.Id(), ipArrayName); err != nil {
		// Groups tagged as hidden are not returned by the list API, so the
		// configured group cannot be resolved back here. Keep the state values
		// already known instead of overwriting them with another group.
		if !NotFoundError(err) {
			return WrapError(err)
		}
		log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version whitelist group %s not found: %s", ipArrayName, err)
	} else {
		d.Set("db_cluster_ip_array_name", ipArrayObject["DBClusterIPArrayName"])
		d.Set("security_ips", ipArrayObject["SecurityIPList"])
	}

	sslObject, err := adbService.DescribeAdbDbClusterLakeVersionDBClusterSSL(d.Id())
	if err != nil {
		log.Printf("[WARN] Resource alicloud_adb_db_cluster_lake_version DescribeAdbDbClusterLakeVersionDBClusterSSL Failed!!! %s", err)
	} else {
		d.Set("enable_ssl", sslObject["SSLEnabled"])
	}

	return nil
}

func resourceAliCloudAdbDbClusterLakeVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("compute_resource") {
		update = true
	}
	if v, ok := d.GetOk("compute_resource"); ok {
		request["ComputeResource"] = v
	}

	if !d.IsNewResource() && d.HasChange("storage_resource") {
		update = true
	}
	if v, ok := d.GetOk("storage_resource"); ok {
		request["StorageResource"] = v
	}

	if !d.IsNewResource() && d.HasChange("reserved_node_size") {
		update = true

		if v, ok := d.GetOk("reserved_node_size"); ok {
			request["ReservedNodeSize"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("reserved_node_count") {
		update = true

		if v, ok := d.GetOkExists("reserved_node_count"); ok {
			request["ReservedNodeCount"] = v
		}
	}

	if update {
		action := "ModifyDBCluster"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyClusterAccessWhiteListReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if d.HasChange("security_ips") {
		update = true
	}
	if v, ok := d.GetOk("security_ips"); ok {
		modifyClusterAccessWhiteListReq["SecurityIps"] = v
	}

	// The group selector and the modification mode are not values of their own:
	// they decide which whitelist group receives SecurityIps. They must be sent
	// whenever they are configured, otherwise the API falls back to the default
	// group and the configured group is never updated.
	if v, ok := d.GetOk("db_cluster_ip_array_name"); ok {
		modifyClusterAccessWhiteListReq["DBClusterIPArrayName"] = v
		if d.HasChange("db_cluster_ip_array_name") {
			update = true
		}
	}

	if v, ok := d.GetOk("db_cluster_ip_array_attribute"); ok {
		modifyClusterAccessWhiteListReq["DBClusterIPArrayAttribute"] = v
		if d.HasChange("db_cluster_ip_array_attribute") {
			update = true
		}
	}

	if v, ok := d.GetOk("modify_mode"); ok {
		modifyClusterAccessWhiteListReq["ModifyMode"] = v
		if d.HasChange("modify_mode") {
			update = true
		}
	}

	if update {
		// SecurityIps is required by ModifyClusterAccessWhiteList.
		if _, ok := modifyClusterAccessWhiteListReq["SecurityIps"]; !ok {
			return WrapErrorf(Error("security_ips must be set when the IP whitelist group is changed"), DefaultErrorMsg, d.Id(), "ModifyClusterAccessWhiteList", AlibabaCloudSdkGoERROR)
		}
	}

	if update {
		action := "ModifyClusterAccessWhiteList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyClusterAccessWhiteListReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyClusterAccessWhiteListReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterSSLReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if d.HasChange("enable_ssl") {
		update = true

		if v, ok := d.GetOkExists("enable_ssl"); ok {
			modifyDBClusterSSLReq["EnableSSL"] = v
		}
	}

	if update {
		action := "ModifyDBClusterSSL"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterSSLReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterSSLReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		sslEnabledStateConf := BuildStateConf([]string{}, []string{fmt.Sprint(modifyDBClusterSSLReq["EnableSSL"])}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionDBClusterSSLStateRefreshFunc(d, []string{}))
		if _, err := sslEnabledStateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterDescriptionReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("db_cluster_description") {
		update = true
	}
	if v, ok := d.GetOk("db_cluster_description"); ok {
		modifyDBClusterDescriptionReq["DBClusterDescription"] = v
	}

	if update {
		action := "ModifyDBClusterDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterDescriptionReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterDescriptionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterResourceGroupReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		modifyDBClusterResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ModifyDBClusterResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterResourceGroupReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterPayTypeReq := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DbClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	modifyDBClusterPayTypeReq["PayType"] = convertAdbDbClusterLakeVersionPaymentTypeRequest(d.Get("payment_type"))

	if modifyDBClusterPayTypeReq["PayType"] == string(Prepaid) {
		if v, ok := d.GetOkExists("period"); ok {
			usedTime := v.(int)
			modifyDBClusterPayTypeReq["UsedTime"] = strconv.Itoa(usedTime)
			modifyDBClusterPayTypeReq["Period"] = string(Month)
			if usedTime > 9 {
				modifyDBClusterPayTypeReq["UsedTime"] = strconv.Itoa(usedTime / 12)
				modifyDBClusterPayTypeReq["Period"] = string(Year)
			}
		}
	}

	if update {
		action := "ModifyDBClusterPayType"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterPayTypeReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterPayTypeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyAuditLogConfigReq := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DBClusterId": d.Id(),
	}

	if d.HasChange("audit_log_status") || d.HasChange("engine_type") {
		// AuditLogStatus is a required parameter of ModifyAuditLogConfig and
		// EngineType defaults to XIHE when omitted, so the current status must
		// always be carried to avoid resetting it while only changing the engine.
		status := fmt.Sprint(d.Get("audit_log_status"))
		if status == "" {
			return WrapErrorf(Error("audit_log_status must be set when engine_type is changed"), DefaultErrorMsg, d.Id(), "ModifyAuditLogConfig", AlibabaCloudSdkGoERROR)
		}
		update = true
		modifyAuditLogConfigReq["AuditLogStatus"] = status
		if v, ok := d.GetOk("engine_type"); ok {
			modifyAuditLogConfigReq["EngineType"] = v
		}
	}

	if update {
		action := "ModifyAuditLogConfig"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyAuditLogConfigReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyAuditLogConfigReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if succeeded, present := adbAuditLogUpdateSucceeded(response); present && !succeeded {
			return WrapErrorf(Error("ModifyAuditLogConfig reported UpdateSucceed=false"), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	modifyEssdCacheConfigReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	// EnableEssdCache is a required parameter of ModifyEssdCacheConfig, so the
	// current switch state is always carried when the request is needed.
	if v, ok := adbBoolConfigValue(d.GetRawConfig(), "enable_essd_cache"); ok {
		modifyEssdCacheConfigReq["EnableEssdCache"] = v
		if d.HasChange("enable_essd_cache") {
			update = true
		}
	} else {
		// A size change still needs the required switch, so carry the value
		// already known from the configuration or the last read.
		modifyEssdCacheConfigReq["EnableEssdCache"] = d.Get("enable_essd_cache")
	}

	if v, ok := d.GetOk("essd_cache_size"); ok {
		modifyEssdCacheConfigReq["EssdCacheSize"] = v
		if d.HasChange("essd_cache_size") {
			update = true
		}
	}

	if update {
		action := "ModifyEssdCacheConfig"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyEssdCacheConfigReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyEssdCacheConfigReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	modifyLakeCacheSizeReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	// EnableLakeCache is a required parameter of ModifyLakeCacheSize, so the
	// current switch state is always carried when the request is needed.
	if v, ok := adbBoolConfigValue(d.GetRawConfig(), "enable_lake_cache"); ok {
		modifyLakeCacheSizeReq["EnableLakeCache"] = v
		if d.HasChange("enable_lake_cache") {
			update = true
		}
	} else {
		modifyLakeCacheSizeReq["EnableLakeCache"] = d.Get("enable_lake_cache")
	}

	if v, ok := d.GetOk("capacity"); ok {
		modifyLakeCacheSizeReq["Capacity"] = v
		if d.HasChange("capacity") {
			update = true
		}
	}

	if update {
		action := "ModifyLakeCacheSize"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyLakeCacheSizeReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyLakeCacheSizeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if ok, present := adbLakeCacheResponseCodeOK(response); present && !ok {
			return WrapErrorf(Error(fmt.Sprintf("ModifyLakeCacheSize reported Code=%v", response["Code"])), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	update = false
	modifyCompactionServiceSwitchReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if v, ok := adbBoolConfigValue(d.GetRawConfig(), "enable_compaction_service"); ok {
		modifyCompactionServiceSwitchReq["EnableCompactionService"] = v
		if d.HasChange("enable_compaction_service") {
			update = true
		}
	}

	if update {
		action := "ModifyCompactionServiceSwitch"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyCompactionServiceSwitchReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyCompactionServiceSwitchReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	update = false
	modifyMaintainTimeReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if v, ok := d.GetOk("maintain_time"); ok && d.HasChange("maintain_time") {
		update = true
		modifyMaintainTimeReq["MaintainTime"] = v
	}

	if update {
		action := "ModifyDBClusterMaintainTime"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyMaintainTimeReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyMaintainTimeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	update = false
	modifyAutoRenewalAttributeReq := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DBClusterId": d.Id(),
	}

	if fmt.Sprint(d.Get("payment_type")) == "Subscription" {
		if v, ok := d.GetOk("auto_renewal_status"); ok && d.HasChange("auto_renewal_status") {
			update = true
			modifyAutoRenewalAttributeReq["AutoRenewalStatus"] = v
		}

		if v, ok := d.GetOk("auto_renewal_period"); ok && d.HasChange("auto_renewal_period") {
			update = true
			modifyAutoRenewalAttributeReq["AutoRenewalPeriod"] = strconv.Itoa(v.(int))
		}

		if v, ok := d.GetOk("auto_renewal_period_unit"); ok && d.HasChange("auto_renewal_period_unit") {
			update = true
			modifyAutoRenewalAttributeReq["AutoRenewalPeriodUnit"] = v
		}
	}

	if update {
		action := "ModifyAutoRenewalAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyAutoRenewalAttributeReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyAutoRenewalAttributeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	update = false
	upgradeKernelVersionReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && (d.HasChange("db_version") || d.HasChange("switch_mode")) {
		v, ok := d.GetOk("db_version")
		if !ok {
			return WrapErrorf(Error("db_version must be set when the kernel version is changed"), DefaultErrorMsg, d.Id(), "UpgradeKernelVersion", AlibabaCloudSdkGoERROR)
		}
		update = true
		upgradeKernelVersionReq["DBVersion"] = v
		if switchMode, ok := d.GetOk("switch_mode"); ok {
			upgradeKernelVersionReq["SwitchMode"] = switchMode
		}
	}

	if update {
		action := "UpgradeKernelVersion"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, upgradeKernelVersionReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, upgradeKernelVersionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)

	return resourceAliCloudAdbDbClusterLakeVersionRead(d, meta)
}

func resourceAliCloudAdbDbClusterLakeVersionDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok && fmt.Sprint(v) == "Subscription" {
		return WrapError(Error("'Subscription' instance cannot be destroyed: please convert it to 'PayAsYouGo' (set payment_type = \"PayAsYouGo\" and re-apply) or unsubscribe it in the console. To remove it from the state only, run 'terraform state rm'."))
	}

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBCluster"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
		"RegionId":    client.RegionId,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ACS.Cluster.InvalidState"}) {
				// The cluster can still be finishing a previous operation (for example
				// right after a failed apply), which is transient for deletion.
				wait()
				return retry.RetryableError(err)
			}
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	adbService := AdbService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertAdbDbClusterLakeVersionPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}

	return source
}

func convertAdbDbClusterLakeVersionPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}

	return source
}

func adbBoolConfigValue(rawConfig cty.Value, key string) (bool, bool) {
	if rawConfig == cty.NilVal || !rawConfig.IsKnown() || rawConfig.IsNull() {
		return false, false
	}
	if !rawConfig.Type().IsObjectType() || !rawConfig.Type().HasAttribute(key) {
		return false, false
	}
	value := rawConfig.GetAttr(key)
	if value.IsNull() || !value.IsKnown() || value.Type() != cty.Bool {
		return false, false
	}
	return value.True(), true
}

func adbAuditLogUpdateSucceeded(response map[string]interface{}) (bool, bool) {
	if response == nil {
		return false, false
	}
	value, ok := response["UpdateSucceed"]
	if !ok || value == nil {
		return false, false
	}
	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		return strings.EqualFold(typed, "true"), true
	}
	return strings.EqualFold(fmt.Sprint(value), "true"), true
}

func validateAdbDbClusterLakeVersionMaintainTime(i interface{}, k string) ([]string, []error) {
	value, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}
	matched := regexp.MustCompile(`^(\d{2}):00Z-(\d{2}):00Z$`).FindStringSubmatch(value)
	if matched == nil {
		return nil, []error{fmt.Errorf("%s must be an hourly window in the format HH:00Z-HH:00Z, for example 04:00Z-05:00Z", k)}
	}
	start, _ := strconv.Atoi(matched[1])
	end, _ := strconv.Atoi(matched[2])
	if start > 23 || end > 23 || end != (start+1)%24 {
		return nil, []error{fmt.Errorf("%s must start on the hour and span exactly one hour, for example 04:00Z-05:00Z", k)}
	}
	return nil, nil
}

func customizeDiffAdbDbClusterLakeVersionAutoRenewal(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	// Only validate a value the configuration actually sets; values restored from
	// state (import/refresh) must not block planning or destroying.
	if !diff.HasChange("auto_renewal_period") && !diff.HasChange("auto_renewal_period_unit") {
		return nil
	}
	period, ok := diff.Get("auto_renewal_period").(int)
	if !ok || period == 0 {
		return nil
	}
	if fmt.Sprint(diff.Get("auto_renewal_period_unit")) != "Year" {
		return nil
	}
	switch period {
	case 1, 2, 3, 5:
		return nil
	default:
		return fmt.Errorf("auto_renewal_period must be 1, 2, 3 or 5 when auto_renewal_period_unit is Year")
	}
}
