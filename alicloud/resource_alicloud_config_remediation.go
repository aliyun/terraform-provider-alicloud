// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudConfigRemediation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigRemediationCreate,
		Read:   resourceAlicloudConfigRemediationRead,
		Update: resourceAlicloudConfigRemediationUpdate,
		Delete: resourceAlicloudConfigRemediationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config_rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"invoke_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NON_EXECUTION", "AUTO_EXECUTION", "MANUAL_EXECUTION", "NOT_CONFIG"}, false),
			},
			"params": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
			},
			"remediation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remediation_source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALIYUN", "CUSTOM"}, false),
			},
			"remediation_template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remediation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"OOS", "FC"}, false),
			},
		},
	}
}

func resourceAlicloudConfigRemediationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateRemediation"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("config_rule_id"); ok {
		request["ConfigRuleId"] = v
	}
	if v, ok := d.GetOk("remediation_type"); ok {
		request["RemediationType"] = v
	}
	if v, ok := d.GetOk("remediation_template_id"); ok {
		request["RemediationTemplateId"] = v
	}
	if v, ok := d.GetOk("invoke_type"); ok {
		request["InvokeType"] = v
	}
	if v, ok := d.GetOk("remediation_source_type"); ok {
		request["SourceType"] = v
	}
	if v, ok := d.GetOk("params"); ok {
		request["Params"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_remediation", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RemediationId"]))

	return resourceAlicloudConfigRemediationRead(d, meta)
}

func resourceAlicloudConfigRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := ConfigServiceV2{client}

	objectRaw, err := configServiceV2.DescribeConfigRemediation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_remediation DescribeConfigRemediation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_rule_id", objectRaw["ConfigRuleId"])
	d.Set("invoke_type", objectRaw["InvokeType"])
	d.Set("params", objectRaw["RemediationOriginParams"])
	d.Set("remediation_source_type", objectRaw["RemediationSourceType"])
	d.Set("remediation_template_id", objectRaw["RemediationTemplateId"])
	d.Set("remediation_type", objectRaw["RemediationType"])
	d.Set("remediation_id", objectRaw["RemediationId"])

	return nil
}

func resourceAlicloudConfigRemediationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	update := false
	update = false
	action := "UpdateRemediation"

	request = make(map[string]interface{})

	request["RemediationId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	if d.HasChange("params") {
		update = true
		if v, ok := d.GetOk("params"); ok {
			request["Params"] = v
		}
	}
	if d.HasChange("invoke_type") {
		update = true
		if v, ok := d.GetOk("invoke_type"); ok {
			request["InvokeType"] = v
		}
	}
	if d.HasChange("remediation_template_id") {
		update = true
		if v, ok := d.GetOk("remediation_template_id"); ok {
			request["RemediationTemplateId"] = v
		}
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err := client.RpcPost("Config", "2020-09-07", action, nil, request, true)
			request["ClientToken"] = buildClientToken(action)

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
		d.SetPartial("params")
		d.SetPartial("invoke_type")
		d.SetPartial("remediation_template_id")
	}

	return resourceAlicloudConfigRemediationRead(d, meta)
}

func resourceAlicloudConfigRemediationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteRemediations"
	var request map[string]interface{}
	request = make(map[string]interface{})

	request["RemediationIds"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err := client.RpcPost("Config", "2020-09-07", action, nil, request, true)

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
		if IsExpectedErrors(err, []string{"RemediationConfigNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
