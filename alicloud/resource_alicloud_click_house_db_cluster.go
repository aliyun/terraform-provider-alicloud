// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudClickHouseDBCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseDBClusterCreate,
		Read:   resourceAliCloudClickHouseDBClusterRead,
		Update: resourceAliCloudClickHouseDBClusterUpdate,
		Delete: resourceAliCloudClickHouseDBClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Basic", "HighAvailability"}, true),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_cluster_ip_array_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_cluster_description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"db_cluster_network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_cluster_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"19.15.2.2", "20.3.10.75", "20.8.7.15", "21.8.10.19", "22.8.5.29"}, false),
			},
			"db_node_group_count": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringLenBetween(1, 48),
			},
			"db_node_storage": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(100, 32000),
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encryption_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modify_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Cover", "Append", "Delete"}, true),
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, true),
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Year", "Month"}, true),
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"renewal_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"restart_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_ips": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Creating", "Deleting", "Restarting", "Preparing", "Running"}, false),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"CloudESSD_PL0", "CloudESSD_PL1", "CloudESSD_PL2", "CloudESSD_PL3", "CloudEfficiency", "cloud_essd", "cloud_efficiency", "cloud_essd_pl2", "cloud_essd_pl3"}, true),
			},
			"used_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id_bak": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id_bak_second": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudClickHouseDBClusterCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("used_time"); ok {
		request["UsedTime"] = v
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}
	request["VSwitchId"] = d.Get("vswitch_id")
	if v, ok := d.GetOk("encryption_key"); ok {
		request["EncryptionKey"] = v
	}
	if v, ok := d.GetOk("encryption_type"); ok {
		request["EncryptionType"] = v
	}
	request["DBClusterCategory"] = d.Get("category")
	request["PayType"] = convertClickHouseDBClusterPayTypeRequest(d.Get("payment_type").(string))
	request["DbNodeStorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("db_cluster_name"); ok {
		request["DBClusterDescription"] = v
	}
	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["DBClusterClass"] = d.Get("db_cluster_class")
	request["DBNodeGroupCount"] = d.Get("db_node_group_count")
	request["DBNodeStorage"] = d.Get("db_node_storage")
	request["DBClusterNetworkType"] = d.Get("db_cluster_network_type")
	request["DBClusterVersion"] = d.Get("db_cluster_version")
	if v, ok := d.GetOk("vswitch_id_bak"); ok {
		request["VSwitchBak"] = v
	}
	if v, ok := d.GetOk("vswitch_id_bak_second"); ok {
		request["VSwitchBak2"] = v
	}
	if (request["ZoneId"] == nil || request["VpcId"] == nil) && request["VSwitchId"] != nil {
		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitchWithTeadsl(request["VSwitchId"].(string))
		if err != nil {
			return WrapError(err)
		}
		if v, ok := request["VPCId"].(string); !ok || v == "" {
			request["VPCId"] = vsw["VpcId"]
		}
		if v, ok := request["ZoneId"].(string); !ok || v == "" {
			request["ZoneId"] = vsw["ZoneId"]
		}
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_db_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DBClusterId"]))

	clickHouseServiceV2 := ClickHouseServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 50*time.Second, clickHouseServiceV2.ClickHouseDBClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudClickHouseDBClusterUpdate(d, meta)
}

func resourceAliCloudClickHouseDBClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseDBCluster(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_db_cluster DescribeClickHouseDBCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("category", objectRaw["Category"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("db_cluster_name", objectRaw["DBClusterDescription"])
	d.Set("db_cluster_description", objectRaw["DBClusterDescription"])
	d.Set("db_cluster_network_type", objectRaw["DBClusterNetworkType"])
	d.Set("db_node_storage", objectRaw["DBNodeStorage"])
	d.Set("encryption_key", objectRaw["EncryptionKey"])
	d.Set("encryption_type", objectRaw["EncryptionType"])
	d.Set("maintain_time", objectRaw["MaintainTime"])
	d.Set("payment_type", convertClickHouseDBClusterDBClusterPayTypeResponse(objectRaw["PayType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["DBClusterStatus"])
	d.Set("storage_type", convertClickHouseDbClusterStorageTypeResponse(objectRaw["StorageType"].(string)))
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("connection_string", objectRaw["ConnectionString"])
	d.Set("port", objectRaw["Port"])

	if convertClickHouseDBClusterPayTypeRequest(d.Get("payment_type").(string)) == "Prepaid" {
		objectRaw, err = clickHouseServiceV2.DescribeDescribeAutoRenewAttribute(d.Id())
		if err != nil {
			return WrapError(err)
		}
	}

	d.Set("duration", objectRaw["Duration"])
	d.Set("period_unit", objectRaw["PeriodUnit"])
	d.Set("renewal_status", objectRaw["RenewalStatus"])

	objectRaw, err = clickHouseServiceV2.DescribeDescribeDBClusterAccessWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("db_cluster_ip_array_name", objectRaw["DBClusterIPArrayName"])
	d.Set("security_ips", objectRaw["SecurityIPList"])

	return nil
}

func resourceAliCloudClickHouseDBClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "ModifyDBClusterDescription"
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBClusterId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("db_cluster_name") {
		update = true
		request["DBClusterDescription"] = d.Get("db_cluster_name")
	}
	if !d.IsNewResource() && d.HasChange("db_cluster_description") {
		update = true
		request["DBClusterDescription"] = d.Get("db_cluster_description")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)

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
		clickHouseServiceV2 := ClickHouseServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 50*time.Second, clickHouseServiceV2.ClickHouseDBClusterStateRefreshFunc(d.Id(), "DBClusterStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("db_cluster_name")
		d.SetPartial("db_cluster_description")
	}
	update = false
	action = "ModifyDBClusterMaintainTime"
	conn, err = client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBClusterId"] = d.Id()
	if d.HasChange("maintain_time") {
		update = true
		request["MaintainTime"] = d.Get("maintain_time")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)

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
		d.SetPartial("maintain_time")
	}
	update = false
	action = "ModifyDBClusterAccessWhiteList"
	conn, err = client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBClusterId"] = d.Id()
	if v, ok := d.GetOk("modify_mode"); ok {
		request["ModifyMode"] = v
	}
	if d.HasChange("security_ips") {
		update = true
		request["SecurityIps"] = d.Get("security_ips")
	}

	if d.HasChange("db_cluster_ip_array_name") {
		update = true
		request["DBClusterIPArrayName"] = d.Get("db_cluster_ip_array_name")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)

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
		d.SetPartial("security_ips")
		d.SetPartial("db_cluster_ip_array_name")
	}
	update = false
	action = "ModifyDBCluster"
	conn, err = client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBClusterId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["DBNodeGroupCount"] = d.Get("db_node_group_count")
	if !d.IsNewResource() && d.HasChange("db_node_storage") {
		update = true
	}
	request["DBNodeStorage"] = d.Get("db_node_storage")
	request["DBClusterClass"] = d.Get("db_cluster_class")
	if !d.IsNewResource() && d.HasChange("storage_type") {
		update = true
	}
	request["DbNodeStorageType"] = d.Get("storage_type")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) || NeedRetry(err) {
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
		d.SetPartial("storage_type")
	}
	update = false
	action = "ChangeResourceGroup"
	conn, err = client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if _, ok := d.GetOk("resource_group_id"); ok && !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}
	request["ResourceId"] = d.Id()

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if IsExpectedErrors(err, []string{"SystemError"}) || NeedRetry(err) {
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
		d.SetPartial("resource_group_id")
		d.SetPartial("db_cluster_name")
		d.SetPartial("db_cluster_description")
	}
	update = false
	action = "ModifyAutoRenewAttribute"
	conn, err = client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBClusterIds"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("renewal_status") {
		update = true
		request["RenewalStatus"] = d.Get("renewal_status")
	}

	if d.HasChange("duration") {
		update = true
		request["Duration"] = d.Get("duration")
	}

	if d.HasChange("period_unit") {
		update = true
		request["PeriodUnit"] = d.Get("period_unit")
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)

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
		d.SetPartial("renewal_status")
		d.SetPartial("duration")
		d.SetPartial("period_unit")
	}

	if d.HasChange("status") {
		clickhouseService := ClickhouseService{client}
		object, err := clickhouseService.DescribeClickHouseDbCluster(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DBClusterStatus"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"DBClusterId": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "RestartInstance"
				conn, err := client.NewClickhouseClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectDBInstanceState"}) || NeedRetry(err) {
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
			}
			d.SetPartial("status")
		}
		stateConf := BuildStateConf([]string{"RESTARTING"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, clickhouseService.ClickHouseDbClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceAliCloudClickHouseDBClusterRead(d, meta)
}

func resourceAliCloudClickHouseDBClusterDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBCluster"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["DBClusterId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), query, request, &runtime)

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

	return nil
}

func convertClickHouseDBClusterDBClusterPayTypeResponse(source interface{}) interface{} {
	switch source {
	case "Prepaid":
		return "Subscription"
	case "Postpaid":
		return "PayAsYouGo"
	}
	return source
}
func convertClickHouseDBClusterPayTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "Prepaid"
	case "PayAsYouGo":
		return "Postpaid"
	}
	return source
}

func convertClickHouseDbClusterPaymentTypeRequest(source string) string {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}

func convertClickHouseDbClusterPaymentTypeResponse(source string) string {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}

func convertClickHouseDbClusterStorageTypeResponse(source string) string {
	switch source {
	case "CloudESSD":
		return "cloud_essd"
	case "CloudEfficiency":
		return "cloud_efficiency"
	case "CloudESSD_PL2":
		return "cloud_essd_pl2"
	case "CloudESSD_PL3":
		return "cloud_essd_pl3"

	}
	return source
}
