package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudHbaseMultiZoneCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudHbaseMultiZoneClusterCreate,
		Read:   resourceAlicloudHbaseMultiZoneClusterRead,
		Update: resourceAlicloudHbaseMultiZoneClusterUpdate,
		Delete: resourceAlicloudHbaseMultiZoneClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arch_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2.0"}, false),
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"2.0"}, false),
			},
			"arbiter_vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"arbiter_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"core_disk_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(400, 64000),
			},
			"core_disk_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"core_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"core_node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(2, 20),
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"immediate_delete_flag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_disk_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(400, 64000),
			},
			"log_disk_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(4, 400),
			},
			"master_instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"multi_zone_combination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_unit": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
				Optional:     true,
			},
			"primary_core_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"primary_vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"primary_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_ip_list": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standby_core_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"standby_vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"standby_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudHbaseMultiZoneClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMultiZoneCluster"
	request := make(map[string]interface{})
	conn, err := client.NewHbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request["ArbiterVSwitchId"] = d.Get("arbiter_vswitch_id")
	request["ArbiterZoneId"] = d.Get("arbiter_zone_id")
	request["ArchVersion"] = d.Get("arch_version")
	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}
	request["ClusterName"] = d.Get("cluster_name")
	request["CoreDiskSize"] = d.Get("core_disk_size")
	request["CoreDiskType"] = d.Get("core_disk_type")
	request["CoreInstanceType"] = d.Get("core_instance_type")
	request["CoreNodeCount"] = d.Get("core_node_count")
	request["Engine"] = d.Get("engine")
	request["EngineVersion"] = d.Get("engine_version")
	request["LogDiskSize"] = d.Get("log_disk_size")
	request["LogDiskType"] = d.Get("log_disk_type")
	request["LogInstanceType"] = d.Get("log_instance_type")
	request["LogNodeCount"] = d.Get("log_node_count")
	request["MasterInstanceType"] = d.Get("master_instance_type")
	request["MultiZoneCombination"] = d.Get("multi_zone_combination")
	request["PayType"] = convertHBaseSyncPaymentTypeRequest(d.Get("payment_type").(string))
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	request["PrimaryVSwitchId"] = d.Get("primary_vswitch_id")
	request["PrimaryZoneId"] = d.Get("primary_zone_id")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("security_ip_list"); ok {
		request["SecurityIPList"] = v
	}
	request["StandbyVSwitchId"] = d.Get("standby_vswitch_id")
	request["StandbyZoneId"] = d.Get("standby_zone_id")
	request["VpcId"] = d.Get("vpc_id")
	request["ClientToken"] = buildClientToken("CreateMultiZoneCluster")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbase_multi_zone_cluster", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ClusterId"]))
	hbaseService := HBaseService{client}
	stateConf := BuildStateConf([]string{"CREATING", "LAUNCHING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbaseService.HbaseMultiZoneClusterStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudHbaseMultiZoneClusterUpdate(d, meta)
}
func resourceAlicloudHbaseMultiZoneClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}
	object, err := hbaseService.DescribeHbaseMultiZoneCluster(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbase_multi_zone_cluster hbaseService.DescribeHbaseMultiZoneCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("arbiter_zone_id", object["ArbiterZoneId"])
	d.Set("arbiter_vswitch_id", object["ArbiterVSwitchIds"])
	d.Set("standby_vswitch_id", object["StandbyVSwitchIds"])
	d.Set("primary_vswitch_id", object["PrimaryVSwitchIds"])
	d.Set("cluster_name", object["ClusterName"])
	if v, ok := object["CoreDiskSize"]; ok && fmt.Sprint(v) != "0" {
		d.Set("core_disk_size", formatInt(v))
	}
	d.Set("core_disk_type", object["CoreDiskType"])
	d.Set("core_instance_type", object["CoreInstanceType"])
	if v, ok := object["CoreNodeCount"]; ok && fmt.Sprint(v) != "0" {
		d.Set("core_node_count", formatInt(v))
	}
	d.Set("engine", object["Engine"])
	if v, ok := object["LogDiskSize"]; ok && fmt.Sprint(v) != "0" {
		d.Set("log_disk_size", formatInt(v))
	}
	d.Set("log_disk_type", object["LogDiskType"])
	d.Set("log_instance_type", object["LogInstanceType"])
	if v, ok := object["LogNodeCount"]; ok && fmt.Sprint(v) != "0" {
		d.Set("log_node_count", formatInt(v))
	}
	d.Set("master_instance_type", object["MasterInstanceType"])
	d.Set("multi_zone_combination", object["MultiZoneCombination"])
	d.Set("payment_type", convertHBaseSyncPaymentTypeResponse(object["PayType"].(string)))
	d.Set("primary_zone_id", object["PrimaryZoneId"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("standby_zone_id", object["StandbyZoneId"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("vpc_id", object["VpcId"])
	d.Set("status", object["Status"])
	d.Set("engine_version", object["MajorVersion"])
	d.Set("arch_version", object["MajorVersion"])
	return nil
}
func resourceAlicloudHbaseMultiZoneClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbaseService := HBaseService{client}
	conn, err := client.NewHbaseClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := hbaseService.SetResourceTags(d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"ClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("core_disk_size") {
		update = true
	}
	request["CoreDiskSize"] = d.Get("core_disk_size")
	if !d.IsNewResource() && d.HasChange("log_disk_size") {
		update = true
		request["LogDiskSize"] = d.Get("log_disk_size")
	}
	if update {
		action := "ResizeMultiZoneClusterDiskSize"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("core_disk_size")
		d.SetPartial("log_disk_size")
	}
	update = false
	modifyMultiZoneClusterNodeTypeReq := map[string]interface{}{
		"ClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("core_instance_type") {
		update = true
	}
	modifyMultiZoneClusterNodeTypeReq["CoreInstanceType"] = d.Get("core_instance_type")
	if !d.IsNewResource() && d.HasChange("log_instance_type") {
		update = true
	}
	modifyMultiZoneClusterNodeTypeReq["LogInstanceType"] = d.Get("log_instance_type")
	if !d.IsNewResource() && d.HasChange("master_instance_type") {
		update = true
	}
	modifyMultiZoneClusterNodeTypeReq["MasterInstanceType"] = d.Get("master_instance_type")
	if update {
		action := "ModifyMultiZoneClusterNodeType"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, modifyMultiZoneClusterNodeTypeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyMultiZoneClusterNodeTypeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		hbaseService := HBaseService{client}
		stateConf := BuildStateConf([]string{"INSTANCE_LEVEL_MODIFY", "LAUNCHING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbaseService.HbaseMultiZoneClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("core_instance_type")
		d.SetPartial("log_instance_type")
		d.SetPartial("master_instance_type")
	}
	update = false
	resizeMultiZoneClusterNodeCountReq := map[string]interface{}{
		"ClusterId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("core_node_count") {
		update = true
	}
	resizeMultiZoneClusterNodeCountReq["CoreNodeCount"] = d.Get("core_node_count")
	if !d.IsNewResource() && d.HasChange("log_node_count") {
		update = true
	}
	resizeMultiZoneClusterNodeCountReq["LogNodeCount"] = d.Get("log_node_count")
	if update {
		if v, ok := d.GetOk("arbiter_vswitch_id"); ok {
			resizeMultiZoneClusterNodeCountReq["ArbiterVSwitchId"] = v
		}
		if v, ok := d.GetOk("primary_core_node_count"); ok {
			resizeMultiZoneClusterNodeCountReq["PrimaryCoreNodeCount"] = v
		}
		if v, ok := d.GetOk("primary_vswitch_id"); ok {
			resizeMultiZoneClusterNodeCountReq["PrimaryVSwitchId"] = v
		}
		if v, ok := d.GetOk("standby_core_node_count"); ok {
			resizeMultiZoneClusterNodeCountReq["StandbyCoreNodeCount"] = v
		}
		if v, ok := d.GetOk("standby_vswitch_id"); ok {
			resizeMultiZoneClusterNodeCountReq["StandbyVSwitchId"] = v
		}
		action := "ResizeMultiZoneClusterNodeCount"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, resizeMultiZoneClusterNodeCountReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, resizeMultiZoneClusterNodeCountReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		hbaseService := HBaseService{client}
		stateConf := BuildStateConf([]string{"INSTANCE_LEVEL_MODIFY", "LAUNCHING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, hbaseService.HbaseMultiZoneClusterStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("core_node_count")
		d.SetPartial("log_node_count")
	}
	d.Partial(false)
	return resourceAlicloudHbaseMultiZoneClusterRead(d, meta)
}
func resourceAlicloudHbaseMultiZoneClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMultiZoneCluster"
	var response map[string]interface{}
	conn, err := client.NewHbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ClusterId": d.Id(),
	}

	if v, ok := d.GetOkExists("immediate_delete_flag"); ok {
		request["ImmediateDeleteFlag"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertHBaseSyncPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}
	return source
}
func convertHBaseSyncPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}
	return source
}
