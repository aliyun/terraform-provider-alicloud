package cddc

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDedicatedHostAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedHostAccountCreate,
		Read:   resourceDedicatedHostAccountRead,
		Update: resourceDedicatedHostAccountUpdate,
		Delete: resourceDedicatedHostAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.225.1 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{0,14}[a-z0-9]$`), "The account name must be 2 to 16 characters in length, and can contain lower case letters, digits, and underscore(_)."),
			},
			"account_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w!@#$%^&*()_+=]{6,32}$`), "The account password must be 6 to 32 characters in length, and can contain letters, digits, and special characters（!@#$%^&*()_+-=)."),
			},
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Admin", "Normal"}, false),
			},
			"dedicated_host_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDedicatedHostAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDedicatedHostAccount"
	request := make(map[string]interface{})
	var err error
	request["AccountName"] = d.Get("account_name")
	request["AccountPassword"] = d.Get("account_password")
	if v, ok := d.GetOk("account_type"); ok {
		request["AccountType"] = v
	}
	request["DedicatedHostId"] = d.Get("dedicated_host_id")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken("CreateDedicatedHostAccount")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cddc_dedicated_host_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DedicatedHostId"], ":", request["AccountName"]))
	cddcService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(request["AccountName"])}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cddcService.DedicatedHostAccountStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDedicatedHostAccountRead(d, meta)
}
func resourceDedicatedHostAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := NewService(client)
	object, err := cddcService.DescribeDedicatedHostAccount(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cddc_dedicated_host_account cddcService.DescribeDedicatedHostAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("account_name", object["AccountName"])
	d.Set("account_type", object["AccountType"])
	d.Set("dedicated_host_id", object["DedicatedHostId"])
	return nil
}
func resourceDedicatedHostAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"AccountName":     parts[1],
		"DedicatedHostId": parts[0],
	}
	if d.HasChange("account_password") {
		request["AccountPassword"] = d.Get("account_password")
	}
	request["RegionId"] = client.RegionId
	action := "ModifyDedicatedHostAccount"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
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
	return resourceDedicatedHostAccountRead(d, meta)
}
func resourceDedicatedHostAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteDedicatedHostAccount"
	var response map[string]interface{}
	request := map[string]interface{}{
		"AccountName":     parts[1],
		"DedicatedHostId": parts[0],
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
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
	// wait for complete async operation
	time.Sleep(2 * time.Minute)
	return nil
}
