package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRamAccessKeyPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamAccessKeyPolicyRead,

		Schema: map[string]*schema.Schema{
			"user_access_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_principal_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_key_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudRamAccessKeyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	userAccessKeyId := d.Get("user_access_key_id").(string)
	userPrincipalName := ""
	if v, ok := d.GetOk("user_principal_name"); ok {
		userPrincipalName = v.(string)
	}

	// Build the composite ID used by DescribeRamAccessKeyPolicy
	var id string
	if userPrincipalName != "" {
		id = fmt.Sprintf("%s:%s", userPrincipalName, userAccessKeyId)
	} else {
		id = userAccessKeyId
	}

	ramServiceV2 := RamServiceV2{client}

	action := "GetAccessKeyPolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"UserAccessKeyId": userAccessKeyId,
	}
	if userPrincipalName != "" {
		request["UserPrincipalName"] = userPrincipalName
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, nil, request, true)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_access_key_policy", action, AlibabaCloudSdkGoERROR)
	}

	// Use the service helper to validate the policy exists (not empty/reset)
	_, err = ramServiceV2.DescribeRamAccessKeyPolicy(id)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_access_key_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(id)
	d.Set("access_key_policy", response["AccessKeyPolicy"])
	d.Set("user_access_key_id", userAccessKeyId)
	if userPrincipalName != "" {
		d.Set("user_principal_name", userPrincipalName)
	}

	return nil
}
