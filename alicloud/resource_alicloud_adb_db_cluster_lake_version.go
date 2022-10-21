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
			"commodity_code": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"compute_resource": {
				Required: true,
				Type:     schema.TypeString,
			},
			"connection_string": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"db_cluster_version": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"5.0"}, false),
				Type:         schema.TypeString,
			},
			"enable_default_resource_group": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"engine": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"engine_version": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"expire_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"expired": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"lock_mode": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"lock_reason": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"payment_type": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
				Type:         schema.TypeString,
			},
			"port": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"storage_resource": {
				Required: true,
				Type:     schema.TypeString,
			},
			"vpc_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vswitch_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"zone_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
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
	if v, ok := d.GetOkExists("enable_default_resource_group"); ok {
		request["EnableDefaultResourcePool"] = v
	}
	request["PayType"] = convertAdbDbClusterLakeVersionPaymentTypeRequest(d.Get("payment_type"))

	request["StorageResource"] = d.Get("storage_resource")

	request["VPCId"] = d.Get("vpc_id")

	request["VSwitchId"] = d.Get("vswitch_id")

	request["ZoneId"] = d.Get("zone_id")

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
	return resourceAlicloudAdbDbClusterLakeVersionRead(d, meta)
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

	request["ComputeResource"] = d.Get("compute_resource")
	if !d.IsNewResource() && d.HasChange("compute_resource") {
		update = true
	}
	request["StorageResource"] = d.Get("storage_resource")
	if !d.IsNewResource() && d.HasChange("storage_resource") {
		update = true
	}

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
