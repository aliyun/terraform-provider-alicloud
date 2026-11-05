// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package config

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceRemediation() *schema.Resource {
	return &schema.Resource{
		Create: ResourceRemediationCreate,
		Read:   ResourceRemediationRead,
		Update: ResourceRemediationUpdate,
		Delete: ResourceRemediationDelete,
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
				ValidateFunc: validation.StringIsJSON,
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

func ResourceRemediationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateRemediation"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = helper.BuildClientToken(action)

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Config", "2020-09-07", action, nil, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_config_remediation", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["RemediationId"]))

	return ResourceRemediationRead(d, meta)
}

func ResourceRemediationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configServiceV2 := NewConfigServiceV2(client)

	objectRaw, err := configServiceV2.DescribeConfigRemediation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_remediation DescribeConfigRemediation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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

func ResourceRemediationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	update := false
	update = false
	action := "UpdateRemediation"

	request = make(map[string]interface{})

	request["RemediationId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)

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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err := client.RpcPost("Config", "2020-09-07", action, nil, request, true)
			request["ClientToken"] = helper.BuildClientToken(action)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return ResourceRemediationRead(d, meta)
}

func ResourceRemediationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteRemediations"
	var request map[string]interface{}
	request = make(map[string]interface{})

	request["RemediationIds"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err := client.RpcPost("Config", "2020-09-07", action, nil, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"RemediationConfigNotExist"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
