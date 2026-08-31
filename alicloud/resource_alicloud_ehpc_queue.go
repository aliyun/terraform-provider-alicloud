package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudEhpcQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEhpcQueueCreate,
		Read:   resourceAliCloudEhpcQueueRead,
		Update: resourceAliCloudEhpcQueueUpdate,
		Delete: resourceAliCloudEhpcQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"compute_nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_renew_period": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"system_disk": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"level": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"enable_ht": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"period": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spot_price_limit": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"period_unit": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_scale_in": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_scale_out": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"hostname_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hostname_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"initial_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"inter_connect": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudEhpcQueueCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateQueue"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		request["ClusterId"] = v
	}

	queue := make(map[string]interface{})

	if v := d.Get("compute_nodes"); !IsNil(v) {
		localData, err := jsonpath.Get("$", v)
		if err != nil {
			localData = make([]interface{}, 0)
		}
		localMaps := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(localData) {
			dataLoopTmp := make(map[string]interface{})
			if dataLoop != nil {
				dataLoopTmp = dataLoop.(map[string]interface{})
			}
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["InstanceType"] = dataLoopTmp["instance_type"]
			dataLoopMap["PeriodUnit"] = dataLoopTmp["period_unit"]
			localData1 := make(map[string]interface{})
			level1, _ := jsonpath.Get("$[0].level", dataLoopTmp["system_disk"])
			if level1 != nil && level1 != "" {
				localData1["Level"] = level1
			}
			category1, _ := jsonpath.Get("$[0].category", dataLoopTmp["system_disk"])
			if category1 != nil && category1 != "" {
				localData1["Category"] = category1
			}
			size1, _ := jsonpath.Get("$[0].size", dataLoopTmp["system_disk"])
			if size1 != nil && size1 != "" {
				localData1["Size"] = size1
			}
			if len(localData1) > 0 {
				dataLoopMap["SystemDisk"] = localData1
			}
			dataLoopMap["Period"] = dataLoopTmp["period"]
			dataLoopMap["AutoRenewPeriod"] = dataLoopTmp["auto_renew_period"]
			dataLoopMap["SpotStrategy"] = dataLoopTmp["spot_strategy"]
			dataLoopMap["AutoRenew"] = dataLoopTmp["auto_renew"]
			dataLoopMap["EnableHT"] = dataLoopTmp["enable_ht"]
			dataLoopMap["Duration"] = dataLoopTmp["duration"]
			dataLoopMap["InstanceChargeType"] = dataLoopTmp["instance_charge_type"]
			dataLoopMap["SpotPriceLimit"] = dataLoopTmp["spot_price_limit"]
			dataLoopMap["ImageId"] = dataLoopTmp["image_id"]
			localMaps = append(localMaps, dataLoopMap)
		}
		queue["ComputeNodes"] = localMaps

	}

	if v, ok := d.GetOk("initial_count"); ok {
		queue["InitialCount"] = v
	}

	if v, ok := d.GetOk("hostname_suffix"); ok {
		queue["HostnameSuffix"] = v
	}

	if v, ok := d.GetOk("queue_name"); ok {
		queue["QueueName"] = v
	}

	if v, ok := d.GetOk("enable_scale_in"); ok {
		queue["EnableScaleIn"] = v
	}

	if v, ok := d.GetOk("min_count"); ok {
		queue["MinCount"] = v
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		vSwitchIds1, _ := jsonpath.Get("$", v)
		if vSwitchIds1 != nil && vSwitchIds1 != "" {
			queue["VSwitchIds"] = vSwitchIds1
		}
	}

	if v, ok := d.GetOk("hostname_prefix"); ok {
		queue["HostnamePrefix"] = v
	}

	if v, ok := d.GetOk("max_count"); ok {
		queue["MaxCount"] = v
	}

	if v, ok := d.GetOk("inter_connect"); ok {
		queue["InterConnect"] = v
	}

	if v, ok := d.GetOk("enable_scale_out"); ok {
		queue["EnableScaleOut"] = v
	}

	queueJson, err := json.Marshal(queue)
	if err != nil {
		return WrapError(err)
	}
	request["Queue"] = string(queueJson)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ehpc_queue", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ClusterId"], response["Name"]))

	return resourceAliCloudEhpcQueueRead(d, meta)
}

func resourceAliCloudEhpcQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ehpcServiceV2 := EhpcServiceV2{client}

	objectRaw, err := ehpcServiceV2.DescribeEhpcQueue(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ehpc_queue DescribeEhpcQueue Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enable_scale_in", objectRaw["EnableScaleIn"])
	d.Set("enable_scale_out", objectRaw["EnableScaleOut"])
	d.Set("hostname_prefix", objectRaw["HostnamePrefix"])
	d.Set("hostname_suffix", objectRaw["HostnameSuffix"])
	d.Set("initial_count", objectRaw["InitialCount"])
	d.Set("inter_connect", objectRaw["InterConnect"])
	d.Set("max_count", objectRaw["MaxCount"])
	d.Set("min_count", objectRaw["MinCount"])
	d.Set("queue_name", objectRaw["QueueName"])

	computeNodesRaw := objectRaw["ComputeNodes"]
	computeNodesMaps := make([]map[string]interface{}, 0)
	if computeNodesRaw != nil {
		for _, computeNodesChildRaw := range convertToInterfaceArray(computeNodesRaw) {
			computeNodesMap := make(map[string]interface{})
			computeNodesChildRaw := computeNodesChildRaw.(map[string]interface{})
			computeNodesMap["auto_renew"] = computeNodesChildRaw["AutoRenew"]
			computeNodesMap["auto_renew_period"] = computeNodesChildRaw["AutoRenewPeriod"]
			computeNodesMap["duration"] = computeNodesChildRaw["Duration"]
			computeNodesMap["enable_ht"] = computeNodesChildRaw["EnableHT"]
			computeNodesMap["image_id"] = computeNodesChildRaw["ImageId"]
			computeNodesMap["instance_charge_type"] = computeNodesChildRaw["InstanceChargeType"]
			computeNodesMap["instance_type"] = computeNodesChildRaw["InstanceType"]
			computeNodesMap["period"] = computeNodesChildRaw["Period"]
			computeNodesMap["period_unit"] = computeNodesChildRaw["PeriodUnit"]
			computeNodesMap["spot_price_limit"] = computeNodesChildRaw["SpotPriceLimit"]
			computeNodesMap["spot_strategy"] = computeNodesChildRaw["SpotStrategy"]

			systemDiskMaps := make([]map[string]interface{}, 0)
			systemDiskMap := make(map[string]interface{})
			systemDiskRaw := make(map[string]interface{})
			if computeNodesChildRaw["SystemDisk"] != nil {
				systemDiskRaw = computeNodesChildRaw["SystemDisk"].(map[string]interface{})
			}
			if len(systemDiskRaw) > 0 {
				systemDiskMap["category"] = systemDiskRaw["Category"]
				systemDiskMap["level"] = systemDiskRaw["Level"]
				systemDiskMap["size"] = systemDiskRaw["Size"]

				systemDiskMaps = append(systemDiskMaps, systemDiskMap)
			}
			computeNodesMap["system_disk"] = systemDiskMaps
			computeNodesMaps = append(computeNodesMaps, computeNodesMap)
		}
	}
	if err := d.Set("compute_nodes", computeNodesMaps); err != nil {
		return err
	}
	vSwitchIdsRaw := make([]interface{}, 0)
	if objectRaw["VSwitchIds"] != nil {
		vSwitchIdsRaw = convertToInterfaceArray(objectRaw["VSwitchIds"])
	}

	d.Set("vswitch_ids", vSwitchIdsRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("cluster_id", parts[0])

	return nil
}

func resourceAliCloudEhpcQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateQueue"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ClusterId"] = parts[0]

	queue := make(map[string]interface{})

	if d.HasChange("compute_nodes") {
		update = true
	}
	if v := d.Get("compute_nodes"); !IsNil(v) || d.HasChange("compute_nodes") {
		localData, err := jsonpath.Get("$", v)
		if err != nil {
			localData = make([]interface{}, 0)
		}
		localMaps := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(localData) {
			dataLoopTmp := make(map[string]interface{})
			if dataLoop != nil {
				dataLoopTmp = dataLoop.(map[string]interface{})
			}
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["InstanceType"] = dataLoopTmp["instance_type"]
			dataLoopMap["PeriodUnit"] = dataLoopTmp["period_unit"]
			if !IsNil(dataLoopTmp["system_disk"]) {
				localData1 := make(map[string]interface{})
				level1, _ := jsonpath.Get("$[0].level", dataLoopTmp["system_disk"])
				if level1 != nil && level1 != "" {
					localData1["Level"] = level1
				}
				category1, _ := jsonpath.Get("$[0].category", dataLoopTmp["system_disk"])
				if category1 != nil && category1 != "" {
					localData1["Category"] = category1
				}
				size1, _ := jsonpath.Get("$[0].size", dataLoopTmp["system_disk"])
				if size1 != nil && size1 != "" {
					localData1["Size"] = size1
				}
				if len(localData1) > 0 {
					dataLoopMap["SystemDisk"] = localData1
				}
			}
			dataLoopMap["Period"] = dataLoopTmp["period"]
			dataLoopMap["AutoRenewPeriod"] = dataLoopTmp["auto_renew_period"]
			dataLoopMap["SpotStrategy"] = dataLoopTmp["spot_strategy"]
			dataLoopMap["AutoRenew"] = dataLoopTmp["auto_renew"]
			dataLoopMap["EnableHT"] = dataLoopTmp["enable_ht"]
			dataLoopMap["Duration"] = dataLoopTmp["duration"]
			dataLoopMap["InstanceChargeType"] = dataLoopTmp["instance_charge_type"]
			dataLoopMap["SpotPriceLimit"] = dataLoopTmp["spot_price_limit"]
			dataLoopMap["ImageId"] = dataLoopTmp["image_id"]
			localMaps = append(localMaps, dataLoopMap)
		}
		queue["ComputeNodes"] = localMaps

	}

	if d.HasChange("hostname_suffix") {
		update = true
	}
	if v, ok := d.GetOk("hostname_suffix"); ok {
		queue["HostnameSuffix"] = v
	}

	if d.HasChange("queue_name") {
		update = true
	}
	if v, ok := d.GetOk("queue_name"); ok {
		queue["QueueName"] = v
	}

	if d.HasChange("enable_scale_in") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_scale_in"); ok {
		queue["EnableScaleIn"] = v
	}

	if d.HasChange("min_count") {
		update = true
	}
	if v, ok := d.GetOkExists("min_count"); ok {
		queue["MinCount"] = v
	}

	if d.HasChange("vswitch_ids") {
		update = true
	}
	vSwitchIds1, _ := jsonpath.Get("$", d.Get("vswitch_ids"))
	if vSwitchIds1 != nil && vSwitchIds1 != "" {
		queue["VSwitchIds"] = vSwitchIds1
	}

	if d.HasChange("hostname_prefix") {
		update = true
	}
	if v, ok := d.GetOk("hostname_prefix"); ok {
		queue["HostnamePrefix"] = v
	}

	if d.HasChange("max_count") {
		update = true
	}
	if v, ok := d.GetOkExists("max_count"); ok {
		queue["MaxCount"] = v
	}

	if d.HasChange("inter_connect") {
		update = true
	}
	if v, ok := d.GetOk("inter_connect"); ok {
		queue["InterConnect"] = v
	}

	if d.HasChange("enable_scale_out") {
		update = true
	}
	if v, ok := d.GetOkExists("enable_scale_out"); ok {
		queue["EnableScaleOut"] = v
	}

	queueJson, err := json.Marshal(queue)
	if err != nil {
		return WrapError(err)
	}
	request["Queue"] = string(queueJson)

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)
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

	return resourceAliCloudEhpcQueueRead(d, meta)
}

func resourceAliCloudEhpcQueueDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteQueues"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ClusterId"] = parts[0]

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "QueueNames.0", parts[1])
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("EHPC", "2024-07-30", action, query, request, true)
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

	return nil
}
