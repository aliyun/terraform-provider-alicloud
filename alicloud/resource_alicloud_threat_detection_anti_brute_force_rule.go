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

func resourceAlicloudThreatDetectionAntiBruteForceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionAntiBruteForceRuleCreate,
		Read:   resourceAlicloudThreatDetectionAntiBruteForceRuleRead,
		Update: resourceAlicloudThreatDetectionAntiBruteForceRuleUpdate,
		Delete: resourceAlicloudThreatDetectionAntiBruteForceRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"anti_brute_force_rule_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"anti_brute_force_rule_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"default_rule": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeBool,
			},
			"fail_count": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"forbidden_time": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"span": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"uuid_list": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudThreatDetectionAntiBruteForceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("anti_brute_force_rule_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("default_rule"); ok {
		request["DefaultRule"] = v
	}
	if v, ok := d.GetOk("fail_count"); ok {
		request["FailCount"] = v
	}
	if v, ok := d.GetOk("forbidden_time"); ok {
		request["ForbiddenTime"] = v
	}
	if v, ok := d.GetOk("span"); ok {
		request["Span"] = v
	}
	if v, ok := d.GetOk("uuid_list"); ok {
		request["UuidList"] = v.([]interface{})
	}

	var response map[string]interface{}
	action := "CreateAntiBruteForceRule"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_anti_brute_force_rule", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.CreateAntiBruteForceRule.RuleId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_anti_brute_force_rule")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudThreatDetectionAntiBruteForceRuleRead(d, meta)
}

func resourceAlicloudThreatDetectionAntiBruteForceRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	object, err := sasService.DescribeThreatDetectionAntiBruteForceRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_anti_brute_force_rule sasService.DescribeThreatDetectionAntiBruteForceRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("anti_brute_force_rule_name", object["Name"])
	d.Set("default_rule", object["DefaultRule"])
	d.Set("fail_count", object["FailCount"])
	d.Set("forbidden_time", object["ForbiddenTime"])
	d.Set("span", object["Span"])
	uuidList, _ := jsonpath.Get("$.UuidList", object)
	d.Set("uuid_list", uuidList)

	return nil
}

func resourceAlicloudThreatDetectionAntiBruteForceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if d.HasChange("anti_brute_force_rule_name") {
		update = true
		if v, ok := d.GetOk("anti_brute_force_rule_name"); ok {
			request["Name"] = v
		}
	}
	if d.HasChange("default_rule") {
		update = true
		if v, ok := d.GetOk("default_rule"); ok {
			request["DefaultRule"] = v
		}
	}
	if d.HasChange("fail_count") {
		update = true
		if v, ok := d.GetOk("fail_count"); ok {
			request["FailCount"] = v
		}
	}
	if d.HasChange("forbidden_time") {
		update = true
		if v, ok := d.GetOk("forbidden_time"); ok {
			request["ForbiddenTime"] = v
		}
	}
	if d.HasChange("span") {
		update = true
		if v, ok := d.GetOk("span"); ok {
			request["Span"] = v
		}
	}
	if d.HasChange("uuid_list") {
		update = true
		if v, ok := d.GetOk("uuid_list"); ok {
			request["UuidList"] = v.([]interface{})
		}
	}

	if update {
		action := "ModifyAntiBruteForceRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudThreatDetectionAntiBruteForceRuleRead(d, meta)
}

func resourceAlicloudThreatDetectionAntiBruteForceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{

		"Ids.1": d.Id(),
	}

	action := "DeleteAntiBruteForceRule"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
