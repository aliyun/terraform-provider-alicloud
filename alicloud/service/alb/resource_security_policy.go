package alb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityPolicyCreate,
		Read:   resourceSecurityPolicyRead,
		Update: resourceSecurityPolicyUpdate,
		Delete: resourceSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(16 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ciphers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.TagsSchema(),
			"tls_versions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSecurityPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSecurityPolicy"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["SecurityPolicyName"] = d.Get("security_policy_name")
	request["TLSVersions"] = d.Get("tls_versions")
	request["Ciphers"] = d.Get("ciphers")
	request["ClientToken"] = helper.BuildClientToken("CreateSecurityPolicy")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alb_security_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SecurityPolicyId"]))

	return resourceSecurityPolicyRead(d, meta)
}
func resourceSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbSecurityPolicy(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_security_policy albService.DescribeAlbSecurityPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("ciphers", object["Ciphers"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("security_policy_name", object["SecurityPolicyName"])
	d.Set("status", object["SecurityPolicyStatus"])
	d.Set("tls_versions", object["TLSVersions"])

	tagResp, err := albService.ListTagResources(d.Id(), "securitypolicy")
	d.Set("tags", tags.TagsToMap(tagResp))

	return nil
}
func resourceSecurityPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := albService.SetResourceTags(d, "securitypolicy"); err != nil {
			return errs.WrapError(err)
		}
	}
	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["NewResourceGroupId"] = v
	}
	request["ResourceType"] = "securitypolicy"
	if update {
		action := "MoveResourceGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, false)
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
	update = false
	updateSecurityPolicyAttributeReq := map[string]interface{}{
		"SecurityPolicyId": d.Id(),
	}
	if d.HasChange("ciphers") {
		update = true
		updateSecurityPolicyAttributeReq[""] = d.Get("ciphers")
	}
	if d.HasChange("security_policy_name") {
		update = true
		updateSecurityPolicyAttributeReq["SecurityPolicyName"] = d.Get("security_policy_name")
	}
	if d.HasChange("tls_versions") {
		update = true
		updateSecurityPolicyAttributeReq[""] = d.Get("tls_versions")
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			updateSecurityPolicyAttributeReq["DryRun"] = v
		}
		if v, ok := d.GetOk("tls_versions"); ok {
			updateSecurityPolicyAttributeReq["TLSVersions"] = v
		}
		if v, ok := d.GetOk("ciphers"); ok {
			updateSecurityPolicyAttributeReq["Ciphers"] = v
		}
		action := "UpdateSecurityPolicyAttribute"
		request["ClientToken"] = helper.BuildClientToken("UpdateSecurityPolicyAttribute")
		wait := helper.IncrementalWait(3*time.Second, 30*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, updateSecurityPolicyAttributeReq, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectStatus.SecurityPolicy"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, updateSecurityPolicyAttributeReq)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbSecurityPolicyStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceSecurityPolicyRead(d, meta)
}
func resourceSecurityPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecurityPolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"SecurityPolicyId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = helper.BuildClientToken("DeleteSecurityPolicy")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound.SecurityPolicy"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
