// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudConfigAggregateRemediation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudConfigAggregateRemediationCreate,
		Read:   resourceAliCloudConfigAggregateRemediationRead,
		Update: resourceAliCloudConfigAggregateRemediationUpdate,
		Delete: resourceAliCloudConfigAggregateRemediationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"invoke_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remediation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remediation_origin_params": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remediation_source_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"remediation_template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remediation_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudConfigAggregateRemediationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAggregateRemediation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("aggregator_id"); ok {
		request["AggregatorId"] = v
	}

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("remediation_source_type"); ok {
		request["SourceType"] = v
	}
	request["InvokeType"] = d.Get("invoke_type")
	request["ConfigRuleId"] = d.Get("config_rule_id")
	request["RemediationTemplateId"] = d.Get("remediation_template_id")
	request["Params"] = d.Get("remediation_origin_params")
	request["RemediationType"] = d.Get("remediation_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_aggregate_remediation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["AggregatorId"], response["RemediationId"]))

	return resourceAliCloudConfigAggregateRemediationRead(d, meta)
}

func resourceAliCloudConfigAggregateRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := ConfigServiceV2{client}

	objectRaw, err := configServiceV2.DescribeConfigAggregateRemediation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_aggregate_remediation DescribeConfigAggregateRemediation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_rule_id", objectRaw["ConfigRuleId"])
	d.Set("invoke_type", objectRaw["InvokeType"])
	d.Set("remediation_origin_params", objectRaw["RemediationOriginParams"])
	d.Set("remediation_source_type", objectRaw["RemediationSourceType"])
	d.Set("remediation_template_id", objectRaw["RemediationTemplateId"])
	d.Set("remediation_type", objectRaw["RemediationType"])
	d.Set("aggregator_id", objectRaw["AggregatorId"])
	d.Set("remediation_id", objectRaw["RemediationId"])

	return nil
}

func resourceAliCloudConfigAggregateRemediationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateAggregateRemediation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RemediationId"] = parts[1]
	request["AggregatorId"] = parts[0]

	if d.HasChange("invoke_type") {
		update = true
	}
	request["InvokeType"] = d.Get("invoke_type")
	if d.HasChange("remediation_template_id") {
		update = true
	}
	request["RemediationTemplateId"] = d.Get("remediation_template_id")
	if d.HasChange("remediation_origin_params") {
		update = true
	}
	request["Params"] = d.Get("remediation_origin_params")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
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

	return resourceAliCloudConfigAggregateRemediationRead(d, meta)
}

func resourceAliCloudConfigAggregateRemediationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAggregateRemediations"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RemediationIds"] = parts[1]
	request["AggregatorId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, query, request, true)
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
