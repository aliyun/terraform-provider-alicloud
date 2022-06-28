package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudBastionhostHostShareKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostHostShareKeyCreate,
		Read:   resourceAlicloudBastionhostHostShareKeyRead,
		Update: resourceAlicloudBastionhostHostShareKeyUpdate,
		Delete: resourceAlicloudBastionhostHostShareKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"host_share_key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pass_phrase": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"host_share_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key_finger_print": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudBastionhostHostShareKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateHostShareKey"
	request := make(map[string]interface{})
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request["HostShareKeyName"] = d.Get("host_share_key_name")
	request["InstanceId"] = d.Get("instance_id")
	if v, ok := d.GetOk("pass_phrase"); ok {
		request["PassPhrase"] = v
	}
	request["PrivateKey"] = d.Get("private_key")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_bastionhost_host_share_key", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", response["HostShareKeyId"]))

	return resourceAlicloudBastionhostHostShareKeyUpdate(d, meta)
}
func resourceAlicloudBastionhostHostShareKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := YundunBastionhostService{client}
	object, err := yundunBastionhostService.DescribeBastionhostHostShareKey(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host_share_key yundunBastionhostService.DescribeBastionhostHostShareKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("host_share_key_id", parts[1])
	d.Set("instance_id", parts[0])
	d.Set("host_share_key_name", object["HostShareKeyName"])
	d.Set("private_key_finger_print", object["PrivateKeyFingerPrint"])
	return nil
}
func resourceAlicloudBastionhostHostShareKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
	}
	if !d.IsNewResource() && d.HasChange("host_share_key_name") {
		update = true
		request["HostShareKeyName"] = d.Get("host_share_key_name")
	}
	if !d.IsNewResource() && d.HasChange("private_key") {
		update = true
		request["PrivateKey"] = d.Get("private_key")
	}
	request["RegionId"] = client.RegionId
	if update {
		if v, ok := d.GetOk("pass_phrase"); ok {
			request["PassPhrase"] = v
		}
		action := "ModifyHostShareKey"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudBastionhostHostShareKeyRead(d, meta)
}
func resourceAlicloudBastionhostHostShareKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteHostShareKey"
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"HostShareKeyId": parts[1],
		"InstanceId":     parts[0],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
