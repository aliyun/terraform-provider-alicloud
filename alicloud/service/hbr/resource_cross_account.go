// Package hbr. This file is generated automatically. Please do not modify it manually, thank you!
package hbr

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCrossAccount() *schema.Resource {
	return &schema.Resource{
		Create: ResourceCrossAccountCreate,
		Read:   ResourceCrossAccountRead,
		Delete: ResourceCrossAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cross_account_role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cross_account_user_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceCrossAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddCrossAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("cross_account_user_id"); ok {
		request["CrossAccountUserId"] = v
	}
	if v, ok := d.GetOk("cross_account_role_name"); ok {
		request["CrossAccountRoleName"] = v
	}

	if v, ok := d.GetOk("alias"); ok {
		request["Alias"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_hbr_cross_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["CrossAccountUserId"], request["CrossAccountRoleName"]))

	return ResourceCrossAccountRead(d, meta)
}

func ResourceCrossAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrServiceV2 := HbrServiceV2{client}

	objectRaw, err := hbrServiceV2.DescribeHbrCrossAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_cross_account DescribeHbrCrossAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["Alias"] != nil {
		d.Set("alias", objectRaw["Alias"])
	}
	if objectRaw["CreatedTime"] != nil {
		d.Set("create_time", objectRaw["CreatedTime"])
	}
	if objectRaw["CrossAccountRoleName"] != nil {
		d.Set("cross_account_role_name", objectRaw["CrossAccountRoleName"])
	}
	if objectRaw["CrossAccountUserId"] != nil {
		d.Set("cross_account_user_id", objectRaw["CrossAccountUserId"])
	}

	return nil
}

func ResourceCrossAccountDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteCrossAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["CrossAccountUserId"] = parts[0]
	request["CrossAccountRoleName"] = parts[1]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("hbr", "2017-09-08", action, query, request, true)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
