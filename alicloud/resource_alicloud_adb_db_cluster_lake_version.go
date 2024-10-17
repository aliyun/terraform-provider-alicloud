package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"compute_resource": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_resource": {
				Type:     schema.TypeString,
				Required: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo"}, false),
			},
			"security_ips": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Type:     schema.TypeString,
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
	request["ComputeResource"] = d.Get("compute_resource")
	request["StorageResource"] = d.Get("storage_resource")
	request["PayType"] = convertAdbDbClusterLakeVersionPaymentTypeRequest(d.Get("payment_type"))

	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
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
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, false)
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

	SecurityIPsObject, err := adbService.DescribeClusterAccessWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ips", SecurityIPsObject["SecurityIPList"])

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
	request["ComputeResource"] = d.Get("compute_resource")

	if !d.IsNewResource() && d.HasChange("storage_resource") {
		update = true
	}
	request["StorageResource"] = d.Get("storage_resource")

	if update {
		action := "ModifyDBCluster"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("compute_resource")
		d.SetPartial("storage_resource")
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

	if update {
		action := "ModifyClusterAccessWhiteList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyClusterAccessWhiteListReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
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

		d.SetPartial("security_ips")
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
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterDescriptionReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
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

		d.SetPartial("db_cluster_description")
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
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterResourceGroupReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
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

		d.SetPartial("resource_group_id")
	}

	d.Partial(false)

	return resourceAliCloudAdbDbClusterLakeVersionRead(d, meta)
}

func resourceAliCloudAdbDbClusterLakeVersionDelete(d *schema.ResourceData, meta interface{}) error {
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
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, false)
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

	adbService := AdbService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertAdbDbClusterLakeVersionPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	}

	return source
}

func convertAdbDbClusterLakeVersionPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	}

	return source
}
