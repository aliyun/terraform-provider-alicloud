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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAdbDbClusterLakeVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAdbDbClusterLakeVersionCreate,
		Read:   resourceAlicloudAdbDbClusterLakeVersionRead,
		Update: resourceAlicloudAdbDbClusterLakeVersionUpdate,
		Delete: resourceAlicloudAdbDbClusterLakeVersionDelete,
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
				ValidateFunc: validation.StringInSlice([]string{"5.0"}, false),
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
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
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
			"enable_default_resource_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"commodity_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired": {
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
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_id": {
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

func resourceAlicloudAdbDbClusterLakeVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}

	request["ComputeResource"] = d.Get("compute_resource")

	request["DBClusterVersion"] = d.Get("db_cluster_version")

	request["PayType"] = convertAdbDbClusterLakeVersionPaymentTypeRequest(d.Get("payment_type"))

	request["StorageResource"] = d.Get("storage_resource")

	request["VPCId"] = d.Get("vpc_id")

	request["VSwitchId"] = d.Get("vswitch_id")

	request["ZoneId"] = d.Get("zone_id")

	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}

	if v, ok := d.GetOkExists("enable_default_resource_group"); ok {
		request["EnableDefaultResourcePool"] = v
	}

	var response map[string]interface{}
	action := "CreateDBCluster"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
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

	return resourceAlicloudAdbDbClusterLakeVersionUpdate(d, meta)
}

func resourceAlicloudAdbDbClusterLakeVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	object, err := adbService.DescribeAdbDbClusterLakeVersion(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version adbService.DescribeAdbDbClusterLakeVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("commodity_code", object["CommodityCode"])
	d.Set("compute_resource", object["ComputeResource"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("create_time", object["CreationTime"])
	d.Set("db_cluster_version", object["DBVersion"])
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("expire_time", object["ExpireTime"])
	d.Set("expired", object["Expired"])
	d.Set("lock_mode", object["LockMode"])
	d.Set("lock_reason", object["LockReason"])
	d.Set("payment_type", convertAdbDbClusterLakeVersionPaymentTypeResponse(object["PayType"]))
	d.Set("port", object["Port"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["DBClusterStatus"])
	d.Set("storage_resource", object["StorageResource"])
	d.Set("vpc_id", object["VPCId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("db_cluster_description", object["DBClusterDescription"])

	SecurityIPsObject, err := adbService.DescribeClusterAccessWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ips", SecurityIPsObject["SecurityIPList"])

	return nil
}

func resourceAlicloudAdbDbClusterLakeVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
		"RegionId":    client.RegionId,
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
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
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

	if update {
		action := "ModifyClusterAccessWhiteList"

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-12-01"), StringPointer("AK"), nil, modifyClusterAccessWhiteListReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, modifyClusterAccessWhiteListReq)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
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
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-12-01"), StringPointer("AK"), nil, modifyDBClusterDescriptionReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, modifyDBClusterDescriptionReq)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudAdbDbClusterLakeVersionRead(d, meta)
}

func resourceAlicloudAdbDbClusterLakeVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DBClusterId": d.Id(),
		"RegionId":    client.RegionId,
	}

	action := "DeleteDBCluster"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-12-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
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
