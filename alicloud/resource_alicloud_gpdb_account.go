// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGpdbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbAccountCreate,
		Read:   resourceAliCloudGpdbAccountRead,
		Update: resourceAliCloudGpdbAccountUpdate,
		Delete: resourceAliCloudGpdbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5A-Za-z0-9_]+$"), "The description of the account."),
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5A-Za-z0-9_]+$"), "The account name."),
			},
			"account_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudGpdbAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AccountName"] = d.Get("account_name")
	query["DBInstanceId"] = d.Get("db_instance_id")

	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	request["AccountPassword"] = d.Get("account_password")
	if v, ok := d.GetOk("account_type"); ok {
		request["AccountType"] = v
	}
	if v, ok := d.GetOk("database_name"); ok {
		request["DatabaseName"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["DBInstanceId"], query["AccountName"]))

	gpdbServiceV2 := GpdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gpdbServiceV2.GpdbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudGpdbAccountRead(d, meta)
}

func resourceAliCloudGpdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_account DescribeGpdbAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["AccountDescription"] != nil {
		d.Set("account_description", objectRaw["AccountDescription"])
	}
	if objectRaw["AccountType"] != nil {
		d.Set("account_type", objectRaw["AccountType"])
	}
	if objectRaw["AccountStatus"] != nil {
		d.Set("status", objectRaw["AccountStatus"])
	}
	if objectRaw["AccountName"] != nil {
		d.Set("account_name", objectRaw["AccountName"])
	}
	if objectRaw["DBInstanceId"] != nil {
		d.Set("db_instance_id", objectRaw["DBInstanceId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])
	d.Set("account_name", parts[1])

	return nil
}

func resourceAliCloudGpdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	parts := strings.Split(d.Id(), ":")
	action := "ModifyAccountDescription"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AccountName"] = parts[1]
	query["DBInstanceId"] = parts[0]

	if d.HasChange("account_description") {
		update = true
	}
	request["AccountDescription"] = d.Get("account_description")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ResetAccountPassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AccountName"] = parts[1]
	query["DBInstanceId"] = parts[0]

	if d.HasChange("account_password") {
		update = true
	}
	request["AccountPassword"] = d.Get("account_password")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
	}

	d.Partial(false)
	return resourceAliCloudGpdbAccountRead(d, meta)
}

func resourceAliCloudGpdbAccountDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("account_type"); ok {
		if v == "Super" {
			log.Printf("[WARN] Cannot destroy resource alicloud_gpdb_account which account_type valued Super. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["DBInstanceId"] = parts[0]
	query["AccountName"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)

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

	return nil
}
func convertGpdbAccountStatusResponse(source interface{}) interface{} {
	switch source {
	case "0":
		return "Creating"
	case "1":
		return "Active"
	case "3":
		return "Deleting"
	}
	return source
}
