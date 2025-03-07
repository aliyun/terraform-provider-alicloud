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

func resourceAliCloudThreatDetectionAntiBruteForceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionAntiBruteForceRuleCreate,
		Read:   resourceAliCloudThreatDetectionAntiBruteForceRuleRead,
		Update: resourceAliCloudThreatDetectionAntiBruteForceRuleUpdate,
		Delete: resourceAliCloudThreatDetectionAntiBruteForceRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"anti_brute_force_rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_rule": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"fail_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"forbidden_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"span": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"uuid_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudThreatDetectionAntiBruteForceRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAntiBruteForceRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOkExists("default_rule"); ok {
		request["DefaultRule"] = v
	}
	request["FailCount"] = d.Get("fail_count")
	request["ForbiddenTime"] = d.Get("forbidden_time")
	request["Name"] = d.Get("anti_brute_force_rule_name")
	request["Span"] = d.Get("span")
	if v, ok := d.GetOk("uuid_list"); ok {
		uuidListMapsArray := v.(*schema.Set).List()
		request["UuidList"] = uuidListMapsArray
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_anti_brute_force_rule", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.CreateAntiBruteForceRule.RuleId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionAntiBruteForceRuleRead(d, meta)
}

func resourceAliCloudThreatDetectionAntiBruteForceRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionAntiBruteForceRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_anti_brute_force_rule DescribeThreatDetectionAntiBruteForceRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("anti_brute_force_rule_name", objectRaw["Name"])
	d.Set("default_rule", objectRaw["DefaultRule"])
	d.Set("fail_count", objectRaw["FailCount"])
	d.Set("forbidden_time", objectRaw["ForbiddenTime"])
	d.Set("span", objectRaw["Span"])

	uuidListRaw := make([]interface{}, 0)
	if objectRaw["UuidList"] != nil {
		uuidListRaw = objectRaw["UuidList"].([]interface{})
	}

	err = d.Set("uuid_list", uuidListRaw)
	if err != nil {
		return err
	}

	return nil
}

func resourceAliCloudThreatDetectionAntiBruteForceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyAntiBruteForceRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = d.Id()

	if d.HasChange("default_rule") {
		update = true
		request["DefaultRule"] = d.Get("default_rule")
	}

	if d.HasChange("fail_count") {
		update = true
	}
	request["FailCount"] = d.Get("fail_count")
	if d.HasChange("forbidden_time") {
		update = true
	}
	request["ForbiddenTime"] = d.Get("forbidden_time")
	if d.HasChange("anti_brute_force_rule_name") {
		update = true
	}
	request["Name"] = d.Get("anti_brute_force_rule_name")
	if d.HasChange("span") {
		update = true
	}
	request["Span"] = d.Get("span")
	if d.HasChange("uuid_list") {
		update = true
	}
	if v, ok := d.GetOk("uuid_list"); ok || d.HasChange("uuid_list") {
		uuidListMapsArray := v.(*schema.Set).List()
		request["UuidList"] = uuidListMapsArray
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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

	return resourceAliCloudThreatDetectionAntiBruteForceRuleRead(d, meta)
}

func resourceAliCloudThreatDetectionAntiBruteForceRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAntiBruteForceRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Ids.1"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

func convertThreatDetectionAntiBruteForceRuleRulesIdResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	}
	return source
}
