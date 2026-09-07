package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDirectMailMailAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDirectMailMailAddressCreate,
		Read:   resourceAlicloudDirectMailMailAddressRead,
		Update: resourceAlicloudDirectMailMailAddressUpdate,
		Delete: resourceAlicloudDirectMailMailAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"reply_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sendtype": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"batch", "trigger"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDirectMailMailAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMailAddress"
	request := make(map[string]interface{})
	var err error
	request["AccountName"] = d.Get("account_name")
	if v, ok := d.GetOk("reply_address"); ok {
		request["ReplyAddress"] = v
	}
	request["Sendtype"] = d.Get("sendtype")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_direct_mail_mail_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["MailAddressId"]))

	return resourceAlicloudDirectMailMailAddressUpdate(d, meta)
}
func resourceAlicloudDirectMailMailAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmService := DmService{client}
	object, err := dmService.DescribeDirectMailMailAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_direct_mail_mail_address dmService.DescribeDirectMailMailAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("account_name", object["AccountName"])
	d.Set("reply_address", object["ReplyAddress"])
	d.Set("sendtype", object["Sendtype"])
	d.Set("status", object["AccountStatus"])
	return nil
}
func resourceAlicloudDirectMailMailAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"MailAddressId": d.Id(),
	}
	if d.HasChange("password") {
		update = true
		if v, ok := d.GetOk("password"); ok {
			request["Password"] = v
		}
	}
	if !d.IsNewResource() && d.HasChange("reply_address") {
		update = true
		if v, ok := d.GetOk("reply_address"); ok {
			request["ReplyAddress"] = v
		}
	}
	if update {
		action := "ModifyMailAddress"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, false)
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
	return resourceAlicloudDirectMailMailAddressRead(d, meta)
}
func resourceAlicloudDirectMailMailAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMailAddress"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"MailAddressId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"InvalidMailAddressId.Malformed"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
