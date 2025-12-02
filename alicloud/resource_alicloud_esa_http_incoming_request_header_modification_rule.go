// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaHttpIncomingRequestHeaderModificationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleCreate,
		Read:   resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleRead,
		Update: resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleUpdate,
		Delete: resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_header_modification": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operation": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"add", "del", "modify"}, false),
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, false),
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateHttpIncomingRequestHeaderModificationRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("request_header_modification"); ok {
		requestHeaderModificationMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Operation"] = dataLoopTmp["operation"]
			dataLoopMap["Type"] = dataLoopTmp["type"]
			dataLoopMap["Value"] = dataLoopTmp["value"]
			dataLoopMap["Name"] = dataLoopTmp["name"]
			requestHeaderModificationMapsArray = append(requestHeaderModificationMapsArray, dataLoopMap)
		}
		requestHeaderModificationMapsJson, err := json.Marshal(requestHeaderModificationMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RequestHeaderModification"] = string(requestHeaderModificationMapsJson)
	}

	if v, ok := d.GetOkExists("sequence"); ok {
		request["Sequence"] = v
	}
	if v, ok := d.GetOkExists("site_version"); ok {
		request["SiteVersion"] = v
	}
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_http_incoming_request_header_modification_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleRead(d, meta)
}

func resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaHttpIncomingRequestHeaderModificationRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_http_incoming_request_header_modification_rule DescribeEsaHttpIncomingRequestHeaderModificationRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("sequence", objectRaw["Sequence"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("config_id", objectRaw["ConfigId"])

	requestHeaderModificationRaw := objectRaw["RequestHeaderModification"]
	requestHeaderModificationMaps := make([]map[string]interface{}, 0)
	if requestHeaderModificationRaw != nil {
		for _, requestHeaderModificationChildRaw := range convertToInterfaceArray(requestHeaderModificationRaw) {
			requestHeaderModificationMap := make(map[string]interface{})
			requestHeaderModificationChildRaw := requestHeaderModificationChildRaw.(map[string]interface{})
			requestHeaderModificationMap["name"] = requestHeaderModificationChildRaw["Name"]
			requestHeaderModificationMap["operation"] = requestHeaderModificationChildRaw["Operation"]
			requestHeaderModificationMap["type"] = requestHeaderModificationChildRaw["Type"]
			requestHeaderModificationMap["value"] = requestHeaderModificationChildRaw["Value"]

			requestHeaderModificationMaps = append(requestHeaderModificationMaps, requestHeaderModificationMap)
		}
	}
	if err := d.Set("request_header_modification", requestHeaderModificationMaps); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprint(parts[0]))

	return nil
}

func resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateHttpIncomingRequestHeaderModificationRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if d.HasChange("request_header_modification") {
		update = true
	}
	if v, ok := d.GetOk("request_header_modification"); ok || d.HasChange("request_header_modification") {
		requestHeaderModificationMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Operation"] = dataLoopTmp["operation"]
			dataLoopMap["Value"] = dataLoopTmp["value"]
			dataLoopMap["Name"] = dataLoopTmp["name"]
			dataLoopMap["Type"] = dataLoopTmp["type"]
			requestHeaderModificationMapsArray = append(requestHeaderModificationMapsArray, dataLoopMap)
		}
		requestHeaderModificationMapsJson, err := json.Marshal(requestHeaderModificationMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RequestHeaderModification"] = string(requestHeaderModificationMapsJson)
	}

	if d.HasChange("sequence") {
		update = true
		request["Sequence"] = d.Get("sequence")
	}

	if d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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

	return resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleRead(d, meta)
}

func resourceAliCloudEsaHttpIncomingRequestHeaderModificationRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteHttpIncomingRequestHeaderModificationRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
