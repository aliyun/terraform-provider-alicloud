// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudElasticsearchLogstash() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudElasticsearchLogstashCreate,
		Read:   resourceAliCloudElasticsearchLogstashRead,
		Update: resourceAliCloudElasticsearchLogstashUpdate,
		Delete: resourceAliCloudElasticsearchLogstashDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(181 * time.Minute),
			Update: schema.DefaultTimeout(181 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_config": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"vpc"}, false),
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
						"vs_area": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"node_amount": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_spec": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"cloud_ssd", "cloud_efficiency"}, false),
						},
						"spec": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"payment_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pricing_cycle": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: StringInSlice([]string{"Year", "Month"}, false),
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"auto_renew_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"7.4_with_X-Pack", "6.7_with_X-Pack"}, false),
			},
		},
	}
}

func resourceAliCloudElasticsearchLogstashCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/openapi/logstashes")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	query["clientToken"] = StringPointer(buildClientToken(action))

	paymentInfo := make(map[string]interface{})

	if v := d.Get("payment_info"); !IsNil(v) {
		duration1, _ := jsonpath.Get("$[0].duration", v)
		if duration1 != nil && duration1 != "" {
			paymentInfo["duration"] = duration1
		}
		autoRenewDuration1, _ := jsonpath.Get("$[0].auto_renew_duration", v)
		if autoRenewDuration1 != nil && autoRenewDuration1 != "" {
			paymentInfo["autoRenewDuration"] = autoRenewDuration1
		}
		pricingCycle1, _ := jsonpath.Get("$[0].pricing_cycle", v)
		if pricingCycle1 != nil && pricingCycle1 != "" {
			paymentInfo["pricingCycle"] = pricingCycle1
		}
		autoRenew, _ := jsonpath.Get("$[0].auto_renew", v)
		if autoRenew != nil && autoRenew != "" {
			paymentInfo["isAutoRenew"] = autoRenew
		}

		request["paymentInfo"] = paymentInfo
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["paymentType"] = convertElasticsearchLogstashpaymentTypeRequest(v.(string))
	}
	networkConfig := make(map[string]interface{})

	if v := d.Get("network_config"); v != nil {
		vSwitchId, _ := jsonpath.Get("$[0].vswitch_id", v)
		if vSwitchId != nil && vSwitchId != "" {
			networkConfig["vswitchId"] = vSwitchId
		}
		type1, _ := jsonpath.Get("$[0].type", v)
		if type1 != nil && type1 != "" {
			networkConfig["type"] = type1
		}
		vsArea1, _ := jsonpath.Get("$[0].vs_area", v)
		if vsArea1 != nil && vsArea1 != "" {
			networkConfig["vsArea"] = vsArea1
		}
		vpcId1, _ := jsonpath.Get("$[0].vpc_id", v)
		if vpcId1 != nil && vpcId1 != "" {
			networkConfig["vpcId"] = vpcId1
		}

		request["networkConfig"] = networkConfig
	}

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["version"] = d.Get("version")
	request["nodeAmount"] = d.Get("node_amount")
	nodeSpec := make(map[string]interface{})

	if v := d.Get("node_spec"); v != nil {
		diskType1, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType1 != nil && diskType1 != "" {
			nodeSpec["diskType"] = diskType1
		}
		disk1, _ := jsonpath.Get("$[0].disk", v)
		if disk1 != nil && disk1 != "" {
			nodeSpec["disk"] = disk1
		}
		spec1, _ := jsonpath.Get("$[0].spec", v)
		if spec1 != nil && spec1 != "" {
			nodeSpec["spec"] = spec1
		}

		request["nodeSpec"] = nodeSpec
	}

	body = request
	wait := incrementalWait(10*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, nil, body, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"TokenPreviousRequestProcessError", "ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_elasticsearch_logstash", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Result.instanceId", response)
	d.SetId(fmt.Sprint(id))

	elasticsearchServiceV2 := ElasticsearchServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, elasticsearchServiceV2.ElasticsearchLogstashStateRefreshFunc(d.Id(), "status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudElasticsearchLogstashUpdate(d, meta)
}

func resourceAliCloudElasticsearchLogstashRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchServiceV2 := ElasticsearchServiceV2{client}

	objectRaw, err := elasticsearchServiceV2.DescribeElasticsearchLogstash(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_elasticsearch_logstash DescribeElasticsearchLogstash Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createdAt"])
	d.Set("description", objectRaw["description"])
	d.Set("node_amount", objectRaw["nodeAmount"])
	d.Set("payment_type", objectRaw["paymentType"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("status", objectRaw["status"])
	d.Set("updated_at", objectRaw["updatedAt"])
	d.Set("version", objectRaw["version"])

	networkConfigMaps := make([]map[string]interface{}, 0)
	networkConfigMap := make(map[string]interface{})
	networkConfigRaw := make(map[string]interface{})
	if objectRaw["networkConfig"] != nil {
		networkConfigRaw = objectRaw["networkConfig"].(map[string]interface{})
	}
	if len(networkConfigRaw) > 0 {
		networkConfigMap["type"] = networkConfigRaw["type"]
		networkConfigMap["vswitch_id"] = networkConfigRaw["vswitchId"]
		networkConfigMap["vpc_id"] = networkConfigRaw["vpcId"]
		networkConfigMap["vs_area"] = networkConfigRaw["vsArea"]

		networkConfigMaps = append(networkConfigMaps, networkConfigMap)
	}
	if err := d.Set("network_config", networkConfigMaps); err != nil {
		return err
	}
	nodeSpecMaps := make([]map[string]interface{}, 0)
	nodeSpecMap := make(map[string]interface{})
	nodeSpecRaw := make(map[string]interface{})
	if objectRaw["nodeSpec"] != nil {
		nodeSpecRaw = objectRaw["nodeSpec"].(map[string]interface{})
	}
	if len(nodeSpecRaw) > 0 {
		nodeSpecMap["disk"] = nodeSpecRaw["disk"]
		nodeSpecMap["disk_type"] = nodeSpecRaw["diskType"]
		nodeSpecMap["spec"] = nodeSpecRaw["spec"]

		nodeSpecMaps = append(nodeSpecMaps, nodeSpecMap)
	}
	if err := d.Set("node_spec", nodeSpecMaps); err != nil {
		return err
	}
	tagsMaps := objectRaw["tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudElasticsearchLogstashUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	InstanceId := d.Id()
	action := fmt.Sprintf("/openapi/logstashes/%s", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	query["clientToken"] = StringPointer(buildClientToken(action))
	if !d.IsNewResource() && d.HasChange("node_spec") {
		update = true
	}
	nodeSpec := make(map[string]interface{})

	if v := d.Get("node_spec"); v != nil {
		diskType1, _ := jsonpath.Get("$[0].disk_type", v)
		if diskType1 != nil && diskType1 != "" {
			nodeSpec["diskType"] = diskType1
		}
		disk1, _ := jsonpath.Get("$[0].disk", v)
		if disk1 != nil && disk1 != "" {
			nodeSpec["disk"] = disk1
		}
		spec1, _ := jsonpath.Get("$[0].spec", v)
		if spec1 != nil && spec1 != "" {
			nodeSpec["spec"] = spec1
		}

		request["nodeSpec"] = nodeSpec
	}

	if !d.IsNewResource() && d.HasChange("node_amount") {
		update = true
	}
	request["nodeAmount"] = d.Get("node_amount")
	body = request
	if update {
		wait := incrementalWait(30*time.Second, 30*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("elasticsearch", "2017-06-13", action, query, nil, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"TokenPreviousRequestProcessError", "ConcurrencyUpdateInstanceConflict", "InstanceDuplicateScheduledTask", "ServiceUnavailable", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, elasticsearchServiceV2.ElasticsearchLogstashStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/logstashes/%s/description", InstanceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	query["clientToken"] = StringPointer(buildClientToken(action))
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	body = request
	if update {
		wait := incrementalWait(10*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, nil, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"TokenPreviousRequestProcessError", "ServiceUnavailable"}) || NeedRetry(err) {
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
	action = fmt.Sprintf("/openapi/tags")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("tags") {
		update = true
	}
	if v, ok := d.GetOk("tags"); ok || d.HasChange("tags") {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	request["ResourceType"] = "logstash"
	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "ResourceIds.0", d.Id())
	_ = json.Unmarshal([]byte(jsonString), &request)

	body = request
	if update {
		wait := incrementalWait(10*time.Second, 10*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, nil, body, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
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

	return resourceAliCloudElasticsearchLogstashRead(d, meta)
}

func resourceAliCloudElasticsearchLogstashDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	InstanceId := d.Id()
	action := fmt.Sprintf("/openapi/logstashes/%s", InstanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(30*time.Second, 30*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("elasticsearch", "2017-06-13", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"TokenPreviousRequestProcessError", "ServiceUnavailable", "InstanceActivating"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound", "ResourceNotfound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertElasticsearchLogstashpaymentTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Subscription":
		return "prepaid"
	case "PayAsYouGo":
		return "postpaid"
	}
	return source
}
