package slb

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

func ResourceTlsCipherPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTlsCipherPolicyCreate,
		Read:   resourceTlsCipherPolicyRead,
		Update: resourceTlsCipherPolicyUpdate,
		Delete: resourceTlsCipherPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ciphers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_cipher_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tls_versions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTlsCipherPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTLSCipherPolicy"
	request := make(map[string]interface{})
	var err error
	request["Ciphers"] = d.Get("ciphers")
	request["RegionId"] = client.RegionId
	request["Name"] = d.Get("tls_cipher_policy_name")
	request["TLSVersions"] = d.Get("tls_versions")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_slb_tls_cipher_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TLSCipherPolicyId"]))

	return resourceTlsCipherPolicyRead(d, meta)
}
func resourceTlsCipherPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := NewService(client)
	object, err := slbService.DescribeSlbTlsCipherPolicy(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_slb_tls_cipher_policy slbService.DescribeSlbTlsCipherPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("ciphers", object["Ciphers"])
	d.Set("status", object["Status"])
	d.Set("tls_cipher_policy_name", object["Name"])
	d.Set("tls_versions", object["TLSVersions"])
	return nil
}
func resourceTlsCipherPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"TLSCipherPolicyId": d.Id(),
	}
	if d.HasChange("ciphers") {
		update = true
	}
	request["Ciphers"] = d.Get("ciphers")
	request["RegionId"] = client.RegionId
	if d.HasChange("tls_cipher_policy_name") {
		update = true
	}
	request["Name"] = d.Get("tls_cipher_policy_name")
	if d.HasChange("tls_versions") {
		update = true
	}
	request["TLSVersions"] = d.Get("tls_versions")
	if update {
		action := "SetTLSCipherPolicyAttribute"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
	return resourceTlsCipherPolicyRead(d, meta)
}
func resourceTlsCipherPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTLSCipherPolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"TLSCipherPolicyId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Slb", "2014-05-15", action, nil, request, false)
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
	return nil
}
