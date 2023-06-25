// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"cloud_ssd", "cloud_efficiency"}, false),
						},
						"spec": {
							Type:     schema.TypeString,
							Required: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"payment_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pricing_cycle": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"Year", "Month"}, false),
						},
						"auto_renew": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"auto_renew_duration": {
							Type:     schema.TypeInt,
							Optional: true,
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
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
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
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	query["clientToken"] = StringPointer(buildClientToken(action))

	if v, ok := d.GetOk("payment_type"); ok {
		request["paymentType"] = convertElasticsearchpaymentTypeRequest(v.(string))
	}
	request["version"] = d.Get("version")
	request["nodeAmount"] = d.Get("node_amount")
	objectDataLocalMap := make(map[string]interface{})
	if v, ok := d.GetOk("node_spec"); ok {
		nodeNative, _ := jsonpath.Get("$[0].disk", v)
		if nodeNative != "" {
			objectDataLocalMap["disk"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].disk_type", v)
		if nodeNative1 != "" {
			objectDataLocalMap["diskType"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].spec", v)
		if nodeNative2 != "" {
			objectDataLocalMap["spec"] = nodeNative2
		}
	}
	request["nodeSpec"] = objectDataLocalMap

	objectDataLocalMap1 := make(map[string]interface{})
	if v, ok := d.GetOk("network_config"); ok {
		nodeNative3, _ := jsonpath.Get("$[0].type", v)
		if nodeNative3 != "" {
			objectDataLocalMap1["type"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].vpc_id", v)
		if nodeNative4 != "" {
			objectDataLocalMap1["vpcId"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].vswitch_id", v)
		if nodeNative5 != "" {
			objectDataLocalMap1["vswitchId"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].vs_area", v)
		if nodeNative6 != "" {
			objectDataLocalMap1["vsArea"] = nodeNative6
		}
	}
	request["networkConfig"] = objectDataLocalMap1

	objectDataLocalMap2 := make(map[string]interface{})
	if v, ok := d.GetOk("payment_info"); ok {
		nodeNative7, _ := jsonpath.Get("$[0].duration", v)
		if nodeNative7 != "" {
			objectDataLocalMap2["duration"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].pricing_cycle", v)
		if nodeNative8 != "" {
			objectDataLocalMap2["pricingCycle"] = nodeNative8
		}
		nodeNative9, _ := jsonpath.Get("$[0].auto_renew_duration", v)
		if nodeNative9 != "" {
			objectDataLocalMap2["autoRenewDuration"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].auto_renew", v)
		if nodeNative10 != "" {
			objectDataLocalMap2["isAutoRenew"] = nodeNative10
		}
	}
	request["paymentInfo"] = objectDataLocalMap2

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	body["body"] = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})
		query["clientToken"] = StringPointer(buildClientToken(action))

		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable", "TokenPreviousRequestProcessError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_elasticsearch_logstash", action, AlibabaCloudSdkGoERROR)
	}

	id, err := jsonpath.Get("$.Result.instanceId", response)
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
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["status"])
	d.Set("version", objectRaw["version"])
	networkConfigMaps := make([]map[string]interface{}, 0)
	networkConfigMap := make(map[string]interface{})
	networkConfig1Raw := make(map[string]interface{})
	if objectRaw["networkConfig"] != nil {
		networkConfig1Raw = objectRaw["networkConfig"].(map[string]interface{})
	}
	if len(networkConfig1Raw) > 0 {
		networkConfigMap["type"] = networkConfig1Raw["type"]
		networkConfigMap["vswitch_id"] = networkConfig1Raw["vswitchId"]
		networkConfigMap["vpc_id"] = networkConfig1Raw["vpcId"]
		networkConfigMap["vs_area"] = networkConfig1Raw["vsArea"]
		networkConfigMaps = append(networkConfigMaps, networkConfigMap)
	}
	d.Set("network_config", networkConfigMaps)
	nodeSpecMaps := make([]map[string]interface{}, 0)
	nodeSpecMap := make(map[string]interface{})
	nodeSpec1Raw := make(map[string]interface{})
	if objectRaw["nodeSpec"] != nil {
		nodeSpec1Raw = objectRaw["nodeSpec"].(map[string]interface{})
	}
	if len(nodeSpec1Raw) > 0 {
		nodeSpecMap["disk"] = nodeSpec1Raw["disk"]
		nodeSpecMap["disk_type"] = nodeSpec1Raw["diskType"]
		nodeSpecMap["spec"] = nodeSpec1Raw["spec"]
		nodeSpecMaps = append(nodeSpecMaps, nodeSpecMap)
	}
	d.Set("node_spec", nodeSpecMaps)
	tagsMaps := objectRaw["Tags"]
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
	d.Partial(true)
	InstanceId := d.Id()
	action := fmt.Sprintf("/openapi/logstashes/%s", InstanceId)
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	query["clientToken"] = StringPointer(buildClientToken(action))
	if d.HasChange("node_spec") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v, ok := d.GetOk("node_spec"); ok {
		nodeNative, _ := jsonpath.Get("$[0].spec", v)
		if nodeNative != "" {
			objectDataLocalMap["spec"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].disk", v)
		if nodeNative1 != "" {
			objectDataLocalMap["disk"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].disk_type", v)
		if nodeNative2 != "" {
			objectDataLocalMap["diskType"] = nodeNative2
		}
	}
	request["nodeSpec"] = objectDataLocalMap

	if !d.IsNewResource() && d.HasChange("node_amount") {
		update = true
	}
	request["nodeAmount"] = d.Get("node_amount")
	body["body"] = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})
			query["clientToken"] = StringPointer(buildClientToken(action))

			if err != nil {
				if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "TokenPreviousRequestProcessError", "InstanceStatusNotSupportCurrentAction", "ServiceUnavailable", "InstanceDuplicateScheduledTask"}) || NeedRetry(err) {
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
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 60*time.Second, elasticsearchServiceV2.ElasticsearchLogstashStateRefreshFunc(d.Id(), "status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("node_amount")
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/logstashes/%s/description", InstanceId)
	conn, err = client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	query["clientToken"] = StringPointer(buildClientToken(action))
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["description"] = d.Get("description")
	}

	body["body"] = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})
			query["clientToken"] = StringPointer(buildClientToken(action))

			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "TokenPreviousRequestProcessError"}) || NeedRetry(err) {
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
		d.SetPartial("description")
	}
	update = false
	InstanceId = d.Id()
	action = fmt.Sprintf("/openapi/instances/%s/resourcegroup", InstanceId)
	conn, err = client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	query["clientToken"] = StringPointer(buildClientToken(action))
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		request["resourceGroupId"] = d.Get("resource_group_id")
	}

	body["body"] = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})
			query["clientToken"] = StringPointer(buildClientToken(action))

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
		d.SetPartial("resource_group_id")
	}

	update = false
	if d.HasChange("tags") {
		update = true
		elasticsearchServiceV2 := ElasticsearchServiceV2{client}
		if err := elasticsearchServiceV2.SetResourceTags(d, "LOGSTASH"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudElasticsearchLogstashRead(d, meta)
}

func resourceAliCloudElasticsearchLogstashDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	InstanceId := d.Id()
	action := fmt.Sprintf("/openapi/logstashes/%s", InstanceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["InstanceId"] = d.Id()

	body["body"] = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})

		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable", "InstanceActivating", "TokenPreviousRequestProcessError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound", "ResourceNotfound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertElasticsearchpaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "Subscription":
		return "prepaid"
	case "PayAsYouGo":
		return "postpaid"
	}
	return source
}
