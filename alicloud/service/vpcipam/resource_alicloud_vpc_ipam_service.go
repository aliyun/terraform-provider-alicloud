// Package vpcipam. This file is generated automatically. Please do not modify it manually, thank you!
package vpcipam

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func ResourceVpcIpamService() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudVpcIpamServiceCreate,
		Read:   resourceAliCloudVpcIpamServiceRead,
		Delete: resourceAliCloudVpcIpamServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudVpcIpamServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "OpenVpcIpamService"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("VpcIpam", "2023-02-28", action, query, request, true)
		if errs.IsExpectedErrors(err, []string{"CURRENT_USER_QUANTITY_EXCEED", "OperationFailed.ExceedPurchaseLimit"}) {
			return nil
		}
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_vpc_ipam_service", action, errs.AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudVpcIpamServiceRead(d, meta)
}

func resourceAliCloudVpcIpamServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcIpamServiceV2 := NewService(client)

	objectRaw, err := vpcIpamServiceV2.DescribeVpcIpamService(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_ipam_service DescribeVpcIpamService Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["Enabled"] != nil {
		d.Set("enabled", objectRaw["Enabled"])
	}

	return nil
}

func resourceAliCloudVpcIpamServiceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Service. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
