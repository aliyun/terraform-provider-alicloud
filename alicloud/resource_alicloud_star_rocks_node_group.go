// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudStarRocksNodeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudStarRocksNodeGroupCreate,
		Read:   resourceAliCloudStarRocksNodeGroupRead,
		Update: resourceAliCloudStarRocksNodeGroupUpdate,
		Delete: resourceAliCloudStarRocksNodeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"fast_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"local_storage_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pricing_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"promotion_option_no": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resident_node_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"spec_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudStarRocksNodeGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/webapi/lifecycle/createNodeGroup")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["instanceId"] = v
	}
	query["RegionId"] = StringPointer(client.RegionId)

	if v, ok := d.GetOk("pay_type"); ok {
		request["payType"] = v
	}
	if v, ok := d.GetOkExists("duration"); ok {
		request["duration"] = v
	}
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["pricingCycle"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["autoRenew"] = v
	}
	if v, ok := d.GetOk("node_group_name"); ok {
		request["nodeGroupName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("spec_type"); ok {
		request["specType"] = v
	}
	if v, ok := d.GetOkExists("cu"); ok {
		request["cu"] = v
	}
	if v, ok := d.GetOkExists("storage_size"); ok {
		request["storageSize"] = v
	}
	if v, ok := d.GetOk("storage_performance_level"); ok {
		request["storagePerformanceLevel"] = v
	}
	if v, ok := d.GetOkExists("disk_number"); ok {
		request["diskNumber"] = v
	}
	if v, ok := d.GetOkExists("resident_node_number"); ok {
		request["residentNodeNumber"] = v
	}
	if v, ok := d.GetOk("local_storage_instance_type"); ok {
		request["localStorageInstanceType"] = v
	}
	if v, ok := d.GetOk("promotion_option_no"); ok {
		request["promotionOptionNo"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_star_rocks_node_group", action, AlibabaCloudSdkGoERROR)
	}

	DataNodeGroupIdVar, _ := jsonpath.Get("$.Data.NodeGroupId", response)
	d.SetId(fmt.Sprintf("%v:%v", request["instanceId"], DataNodeGroupIdVar))

	starRocksServiceV2 := StarRocksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudStarRocksNodeGroupRead(d, meta)
}

func resourceAliCloudStarRocksNodeGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	starRocksServiceV2 := StarRocksServiceV2{client}

	objectRaw, err := starRocksServiceV2.DescribeStarRocksNodeGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_star_rocks_node_group DescribeStarRocksNodeGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["BeginTime"])
	d.Set("cu", objectRaw["Cu"])
	d.Set("description", objectRaw["Description"])
	d.Set("disk_number", objectRaw["DiskNumber"])
	d.Set("local_storage_instance_type", objectRaw["LocalStorageInstanceType"])
	d.Set("node_group_name", objectRaw["NodeGroupName"])
	d.Set("pay_type", objectRaw["PayType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resident_node_number", objectRaw["ResidentNodeNumber"])
	d.Set("spec_type", objectRaw["SpecType"])
	d.Set("status", objectRaw["Status"])
	d.Set("storage_performance_level", objectRaw["StoragePerformanceLevel"])
	d.Set("storage_size", objectRaw["StorageSize"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("node_group_id", objectRaw["NodeGroupId"])

	return nil
}

func resourceAliCloudStarRocksNodeGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/webapi/resourceChange/modifyCu")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(parts[0])
	query["NodeGroupId"] = StringPointer(parts[1])
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("cu") {
		update = true
	}
	if v, ok := d.GetOk("cu"); ok {
		query["Target"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("fast_mode"); ok {
		query["FastMode"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		query["PromotionOptionNo"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		starRocksServiceV2 := StarRocksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = fmt.Sprintf("/webapi/resourceChange/modifyDiskNumber")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(parts[0])
	query["NodeGroupId"] = StringPointer(parts[1])
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("disk_number") {
		update = true
	}
	if v, ok := d.GetOk("disk_number"); ok {
		query["Target"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("fast_mode"); ok {
		query["FastMode"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		query["PromotionOptionNo"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		starRocksServiceV2 := StarRocksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = fmt.Sprintf("/webapi/resourceChange/modifyDiskSize")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(parts[0])
	query["NodeGroupId"] = StringPointer(parts[1])
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("storage_size") {
		update = true
	}
	if v, ok := d.GetOk("storage_size"); ok {
		query["Target"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		query["PromotionOptionNo"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		starRocksServiceV2 := StarRocksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = fmt.Sprintf("/webapi/resourceChange/modifyDiskPerformanceLevel")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(parts[0])
	query["NodeGroupId"] = StringPointer(parts[1])
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("storage_performance_level") {
		update = true
	}
	if v, ok := d.GetOk("storage_performance_level"); ok {
		query["Target"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		query["PromotionOptionNo"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		starRocksServiceV2 := StarRocksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = fmt.Sprintf("/webapi/resourceChange/modifyNodeNumber")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(parts[0])
	query["NodeGroupId"] = StringPointer(parts[1])
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("resident_node_number") {
		update = true
	}
	if v, ok := d.GetOk("resident_node_number"); ok {
		query["Target"] = StringPointer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		query["PromotionOptionNo"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		starRocksServiceV2 := StarRocksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = fmt.Sprintf("/webapi/resourceChange/modifySpecType")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(parts[0])
	query["NodeGroupId"] = StringPointer(parts[1])
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("spec_type") {
		update = true
	}
	if v, ok := d.GetOk("spec_type"); ok {
		query["TargetSpecType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("fast_mode"); ok {
		query["FastMode"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("promotion_option_no"); ok {
		query["PromotionOptionNo"] = StringPointer(v.(string))
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		starRocksServiceV2 := StarRocksServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAliCloudStarRocksNodeGroupRead(d, meta)
}

func resourceAliCloudStarRocksNodeGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/webapi/lifecycle/releaseNodeGroup")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = StringPointer(parts[0])
	query["NodeGroupId"] = StringPointer(parts[1])
	query["RegionId"] = StringPointer(client.RegionId)

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaPost("starrocks", "2022-10-19", action, query, nil, body, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	starRocksServiceV2 := StarRocksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 20*time.Second, starRocksServiceV2.StarRocksNodeGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
