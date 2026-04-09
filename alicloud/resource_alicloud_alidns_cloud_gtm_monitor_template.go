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
)

func resourceAliCloudAlidnsCloudGtmMonitorTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlidnsCloudGtmMonitorTemplateCreate,
		Read:   resourceAliCloudAlidnsCloudGtmMonitorTemplateRead,
		Update: resourceAliCloudAlidnsCloudGtmMonitorTemplateUpdate,
		Delete: resourceAliCloudAlidnsCloudGtmMonitorTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"evaluation_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{0, 1, 2, 3}),
			},
			"extend_info": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"failure_rate": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntInSlice([]int{0, 20, 50, 80, 100}),
			},
			"interval": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"15", "60", "300", "900", "1800", "3600"}, false),
			},
			"ip_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"isp_city_nodes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"isp_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ping", "tcp", "http", "https"}, false),
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"2000", "3000", "5000", "10000"}, false),
			},
		},
	}
}

func resourceAliCloudAlidnsCloudGtmMonitorTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCloudGtmMonitorTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	request["Interval"] = d.Get("interval")
	request["Protocol"] = d.Get("protocol")
	request["IpVersion"] = d.Get("ip_version")
	if v, ok := d.GetOk("isp_city_nodes"); ok {
		ispCityNodesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["CityCode"] = dataLoopTmp["city_code"]
			dataLoopMap["IspCode"] = dataLoopTmp["isp_code"]
			ispCityNodesMapsArray = append(ispCityNodesMapsArray, dataLoopMap)
		}
		ispCityNodesMapsJson, err := json.Marshal(ispCityNodesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["IspCityNodes"] = string(ispCityNodesMapsJson)
	}

	request["EvaluationCount"] = d.Get("evaluation_count")
	request["Timeout"] = d.Get("timeout")
	if v, ok := d.GetOk("extend_info"); ok {
		request["ExtendInfo"] = v
	}
	request["FailureRate"] = d.Get("failure_rate")
	request["Name"] = d.Get("name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_cloud_gtm_monitor_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return resourceAliCloudAlidnsCloudGtmMonitorTemplateUpdate(d, meta)
}

func resourceAliCloudAlidnsCloudGtmMonitorTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsServiceV2 := AlidnsServiceV2{client}

	objectRaw, err := alidnsServiceV2.DescribeAlidnsCloudGtmMonitorTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_cloud_gtm_monitor_template DescribeAlidnsCloudGtmMonitorTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("evaluation_count", objectRaw["EvaluationCount"])
	d.Set("extend_info", objectRaw["ExtendInfo"])
	d.Set("failure_rate", objectRaw["FailureRate"])
	d.Set("interval", objectRaw["Interval"])
	d.Set("ip_version", objectRaw["IpVersion"])
	d.Set("name", objectRaw["Name"])
	d.Set("protocol", objectRaw["Protocol"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("timeout", objectRaw["Timeout"])

	ispCityNodeRaw, _ := jsonpath.Get("$.IspCityNodes.IspCityNode", objectRaw)
	ispCityNodesMaps := make([]map[string]interface{}, 0)
	if ispCityNodeRaw != nil {
		for _, ispCityNodeChildRaw := range convertToInterfaceArray(ispCityNodeRaw) {
			ispCityNodesMap := make(map[string]interface{})
			ispCityNodeChildRaw := ispCityNodeChildRaw.(map[string]interface{})
			ispCityNodesMap["city_code"] = ispCityNodeChildRaw["CityCode"]
			ispCityNodesMap["isp_code"] = ispCityNodeChildRaw["IspCode"]

			ispCityNodesMaps = append(ispCityNodesMaps, ispCityNodesMap)
		}
	}
	if err := d.Set("isp_city_nodes", ispCityNodesMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudAlidnsCloudGtmMonitorTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateCloudGtmMonitorTemplate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("interval") {
		update = true
	}
	request["Interval"] = d.Get("interval")
	if !d.IsNewResource() && d.HasChange("isp_city_nodes") {
		update = true
	}
	if v, ok := d.GetOk("isp_city_nodes"); ok || d.HasChange("isp_city_nodes") {
		ispCityNodesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["CityCode"] = dataLoopTmp["city_code"]
			dataLoopMap["IspCode"] = dataLoopTmp["isp_code"]
			ispCityNodesMapsArray = append(ispCityNodesMapsArray, dataLoopMap)
		}
		ispCityNodesMapsJson, err := json.Marshal(ispCityNodesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["IspCityNodes"] = string(ispCityNodesMapsJson)
	}

	if !d.IsNewResource() && d.HasChange("evaluation_count") {
		update = true
	}
	request["EvaluationCount"] = d.Get("evaluation_count")
	if !d.IsNewResource() && d.HasChange("timeout") {
		update = true
	}
	request["Timeout"] = d.Get("timeout")
	if !d.IsNewResource() && d.HasChange("extend_info") {
		update = true
	}
	request["ExtendInfo"] = d.Get("extend_info")

	if !d.IsNewResource() && d.HasChange("failure_rate") {
		update = true
	}
	request["FailureRate"] = d.Get("failure_rate")
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
	}
	request["Name"] = d.Get("name")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
	action = "UpdateCloudGtmMonitorTemplateRemark"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("remark") {
		update = true
		request["Remark"] = d.Get("remark")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudAlidnsCloudGtmMonitorTemplateRead(d, meta)
}

func resourceAliCloudAlidnsCloudGtmMonitorTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCloudGtmMonitorTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TemplateId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
