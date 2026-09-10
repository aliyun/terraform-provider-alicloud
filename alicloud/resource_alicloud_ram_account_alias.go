// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRamAccountAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamAccountAliasCreate,
		Read:   resourceAliCloudRamAccountAliasRead,
		Update: resourceAliCloudRamAccountAliasUpdate,
		Delete: resourceAliCloudRamAccountAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudRamAccountAliasCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "SetAccountAlias"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["AccountAlias"] = d.Get("account_alias")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_account_alias", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := client.AccountId()
	d.SetId(accountId)

	return resourceAliCloudRamAccountAliasRead(d, meta)
}

func resourceAliCloudRamAccountAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamAccountAlias(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_account_alias DescribeRamAccountAlias Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_alias", objectRaw["AccountAlias"])

	return nil
}

func resourceAliCloudRamAccountAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "SetAccountAlias"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if d.HasChange("account_alias") {
		update = true
	}
	request["AccountAlias"] = d.Get("account_alias")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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

	return resourceAliCloudRamAccountAliasRead(d, meta)
}

func resourceAliCloudRamAccountAliasDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Account Alias. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
