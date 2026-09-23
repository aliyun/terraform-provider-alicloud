// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudStarRocksInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudStarRocksInstanceCreate,
		Read:   resourceAliCloudStarRocksInstanceRead,
		Update: resourceAliCloudStarRocksInstanceUpdate,
		Delete: resourceAliCloudStarRocksInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"admin_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backend_node_groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resident_node_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"local_storage_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"cluster_zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"frontend_node_groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"storage_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resident_node_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"local_storage_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"observer_node_groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cu": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disk_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spec_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"storage_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resident_node_number": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"local_storage_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"oss_accessing_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pay_type": {
				Type:     schema.TypeString,
				Required: true,
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
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"run_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitches": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudStarRocksInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/webapi/cluster/createV1")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["clientToken"] = buildClientToken(action)

	request["RegionId"] = client.RegionId
	request["InstanceName"] = d.Get("instance_name")
	request["AdminPassword"] = d.Get("admin_password")
	request["Version"] = d.Get("version")
	request["RunMode"] = d.Get("run_mode")
	request["PackageType"] = d.Get("package_type")
	request["PayType"] = d.Get("pay_type")
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("pricing_cycle"); ok {
		request["PricingCycle"] = v
	}
	if v, ok := d.GetOk("oss_accessing_role_name"); ok {
		request["OssAccessingRoleName"] = v
	}
	if v, ok := d.GetOkExists("duration"); ok {
		request["Duration"] = v
	}
	if v, ok := d.GetOk("frontend_node_groups"); ok {
		frontendNodeGroupsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["cu"] = dataLoopTmp["cu"]
			dataLoopMap["storageSize"] = dataLoopTmp["storage_size"]
			dataLoopMap["storagePerformanceLevel"] = dataLoopTmp["storage_performance_level"]
			dataLoopMap["diskNumber"] = dataLoopTmp["disk_number"]
			dataLoopMap["residentNodeNumber"] = dataLoopTmp["resident_node_number"]
			dataLoopMap["specType"] = dataLoopTmp["spec_type"]
			dataLoopMap["localStorageInstanceType"] = dataLoopTmp["local_storage_instance_type"]
			dataLoopMap["zoneId"] = dataLoopTmp["zone_id"]
			frontendNodeGroupsMapsArray = append(frontendNodeGroupsMapsArray, dataLoopMap)
		}
		request["FrontendNodeGroups"] = frontendNodeGroupsMapsArray
	}

	if v, ok := d.GetOk("backend_node_groups"); ok {
		backendNodeGroupsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["cu"] = dataLoop1Tmp["cu"]
			dataLoop1Map["storageSize"] = dataLoop1Tmp["storage_size"]
			dataLoop1Map["storagePerformanceLevel"] = dataLoop1Tmp["storage_performance_level"]
			dataLoop1Map["diskNumber"] = dataLoop1Tmp["disk_number"]
			dataLoop1Map["residentNodeNumber"] = dataLoop1Tmp["resident_node_number"]
			dataLoop1Map["specType"] = dataLoop1Tmp["spec_type"]
			dataLoop1Map["localStorageInstanceType"] = dataLoop1Tmp["local_storage_instance_type"]
			dataLoop1Map["zoneId"] = dataLoop1Tmp["zone_id"]
			backendNodeGroupsMapsArray = append(backendNodeGroupsMapsArray, dataLoop1Map)
		}
		request["BackendNodeGroups"] = backendNodeGroupsMapsArray
	}

	if v, ok := d.GetOk("vswitches"); ok {
		vSwitchesMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range v.([]interface{}) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["VswId"] = dataLoop2Tmp["vswitch_id"]
			dataLoop2Map["ZoneId"] = dataLoop2Tmp["zone_id"]
			vSwitchesMapsArray = append(vSwitchesMapsArray, dataLoop2Map)
		}
		request["VSwitches"] = vSwitchesMapsArray
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	request["ZoneId"] = d.Get("cluster_zone_id")
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOkExists("encrypted"); ok {
		request["Encrypted"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}
	if v, ok := d.GetOk("promotion_option_no"); ok {
		request["PromotionOptionNo"] = v
	}
	if v, ok := d.GetOk("observer_node_groups"); ok {
		observerNodeGroupsMapsArray := make([]interface{}, 0)
		for _, dataLoop4 := range v.([]interface{}) {
			dataLoop4Tmp := dataLoop4.(map[string]interface{})
			dataLoop4Map := make(map[string]interface{})
			dataLoop4Map["cu"] = dataLoop4Tmp["cu"]
			dataLoop4Map["storageSize"] = dataLoop4Tmp["storage_size"]
			dataLoop4Map["storagePerformanceLevel"] = dataLoop4Tmp["storage_performance_level"]
			dataLoop4Map["diskNumber"] = dataLoop4Tmp["disk_number"]
			dataLoop4Map["residentNodeNumber"] = dataLoop4Tmp["resident_node_number"]
			dataLoop4Map["specType"] = dataLoop4Tmp["spec_type"]
			dataLoop4Map["localStorageInstanceType"] = dataLoop4Tmp["local_storage_instance_type"]
			dataLoop4Map["zoneId"] = dataLoop4Tmp["zone_id"]
			observerNodeGroupsMapsArray = append(observerNodeGroupsMapsArray, dataLoop4Map)
		}
		request["ObserverNodeGroups"] = observerNodeGroupsMapsArray
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_star_rocks_instance", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	starRocksServiceV2 := StarRocksServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, starRocksServiceV2.StarRocksInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudStarRocksInstanceUpdate(d, meta)
}

func resourceAliCloudStarRocksInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	starRocksServiceV2 := StarRocksServiceV2{client}

	objectRaw, err := starRocksServiceV2.DescribeStarRocksInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_star_rocks_instance DescribeStarRocksInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["BeginTime"])
	d.Set("encrypted", objectRaw["Encrypted"])
	d.Set("instance_name", objectRaw["InstanceName"])
	d.Set("kms_key_id", objectRaw["KmsKeyId"])
	d.Set("package_type", objectRaw["PackageType"])
	d.Set("pay_type", objectRaw["PayType"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("run_mode", objectRaw["RunMode"])
	d.Set("status", objectRaw["InstanceStatus"])
	d.Set("version", objectRaw["Version"])
	d.Set("vpc_id", objectRaw["VpcId"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))
	vSwitchesRaw := objectRaw["VSwitches"]
	vswitchesMaps := make([]map[string]interface{}, 0)
	if vSwitchesRaw != nil {
		for _, vSwitchesChildRaw := range vSwitchesRaw.([]interface{}) {
			vswitchesMap := make(map[string]interface{})
			vSwitchesChildRaw := vSwitchesChildRaw.(map[string]interface{})
			vswitchesMap["vswitch_id"] = vSwitchesChildRaw["VswId"]
			vswitchesMap["zone_id"] = vSwitchesChildRaw["ZoneId"]

			vswitchesMaps = append(vswitchesMaps, vswitchesMap)
		}
	}
	if err := d.Set("vswitches", vswitchesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudStarRocksInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := fmt.Sprintf("/webapi/cluster/update_name")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("instance_name") {
		update = true
	}
	if v, ok := d.GetOk("instance_name"); ok {
		query["ClusterName"] = StringPointer(v.(string))
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
	}
	update = false
	action = fmt.Sprintf("/webapi/resourceGroup/change")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["InstanceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["NewResourceGroupId"] = StringPointer(v.(string))
	}

	query["ResourceType"] = StringPointer("instance")

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
	}

	if d.HasChange("tags") {
		starRocksServiceV2 := StarRocksServiceV2{client}
		if err := starRocksServiceV2.SetResourceTags(d, "instance"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudStarRocksInstanceRead(d, meta)
}

func resourceAliCloudStarRocksInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/webapi/cluster/release")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = StringPointer(d.Id())
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
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 20*time.Second, starRocksServiceV2.StarRocksInstanceStateRefreshFunc(d.Id(), "InstanceStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
