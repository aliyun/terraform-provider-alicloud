package privatezone

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
)

func ResourceRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleAttachmentCreate,
		Read:   resourceRuleAttachmentRead,
		Update: resourceRuleAttachmentUpdate,
		Delete: resourceRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpcs": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceRuleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "BindResolverRuleVpc"
	request := make(map[string]interface{})
	var err error
	request["Lang"] = "en"
	request["RuleId"] = d.Get("rule_id")

	for k, vpcConfig := range d.Get("vpcs").(*schema.Set).List() {
		vpcConfigArg := vpcConfig.(map[string]interface{})
		request[fmt.Sprintf("Vpc.%d.VpcId", k+1)] = vpcConfigArg["vpc_id"]
		request[fmt.Sprintf("Vpc.%d.RegionId", k+1)] = vpcConfigArg["region_id"]
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {

		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_pvtz_rule_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["RuleId"]))

	return resourceRuleAttachmentRead(d, meta)
}
func resourceRuleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := Service{client}
	object, err := pvtzService.DescribeRuleAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_rule_attachment pvtzService.DescribeRuleAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	vpcsSli := make([]map[string]interface{}, 0)
	if vpcs, ok := object["BindVpcs"].([]interface{}); ok {
		for _, vpcConfigArgs := range vpcs {
			vpcConfigArg := vpcConfigArgs.(map[string]interface{})
			vpcsMap := make(map[string]interface{})
			vpcsMap["region_id"] = vpcConfigArg["RegionId"]
			vpcsMap["vpc_id"] = vpcConfigArg["VpcId"]
			vpcsSli = append(vpcsSli, vpcsMap)
		}
	}
	d.Set("vpcs", vpcsSli)
	d.Set("rule_id", d.Id())
	return nil
}
func resourceRuleAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("vpcs") {
		update = true
		for k, vpcConfig := range d.Get("vpcs").(*schema.Set).List() {
			vpcConfigArg := vpcConfig.(map[string]interface{})
			request[fmt.Sprintf("Vpc.%d.VpcId", k+1)] = vpcConfigArg["vpc_id"]
			request[fmt.Sprintf("Vpc.%d.RegionId", k+1)] = vpcConfigArg["region_id"]
		}
	}

	if update {
		request["Lang"] = "en"
		action := "BindResolverRuleVpc"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceRuleAttachmentRead(d, meta)
}
func resourceRuleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "BindResolverRuleVpc"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"RuleId": d.Id(),
	}

	request["Lang"] = "en"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
