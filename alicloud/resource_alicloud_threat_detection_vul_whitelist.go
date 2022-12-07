package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceAlicloudThreatDetectionVulWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionVulWhitelistCreate,
		Read:   resourceAlicloudThreatDetectionVulWhitelistRead,
		Update: resourceAlicloudThreatDetectionVulWhitelistUpdate,
		Delete: resourceAlicloudThreatDetectionVulWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"whitelist": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudThreatDetectionVulWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyCreateVulWhitelist"
	request := make(map[string]interface{})
	conn, err := client.NewThreatdetectionClient()
	if err != nil {
		return WrapError(err)
	}

	request["Whitelist"] = d.Get("whitelist")

	if v, ok := d.GetOk("target_info"); ok {
		request["TargetInfo"] = v
	}

	if v, ok := d.GetOk("reason"); ok {
		request["Reason"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_vul_whitelist", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.VulWhitelistList", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_vul_whitelist")
	} else {
		d.SetId(fmt.Sprint(v.(map[string]interface{})["Id"]))
	}

	return resourceAlicloudThreatDetectionVulWhitelistRead(d, meta)
}

func resourceAlicloudThreatDetectionVulWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}

	object, err := threatDetectionService.DescribeThreatDetectionVulWhitelist(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("whitelist", object["Whitelist"])
	d.Set("target_info", object["Target"])
	d.Set("reason", object["Reason"])

	return nil
}

func resourceAlicloudThreatDetectionVulWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if d.HasChange("target_info") {
		update = true
	}
	if v, ok := d.GetOk("target_info"); ok {
		request["TargetInfo"] = v
	}

	if d.HasChange("reason") {
		update = true
	}
	if v, ok := d.GetOk("reason"); ok {
		request["Reason"] = v
	}

	if update {
		action := "ModifyVulWhitelistTarget"
		conn, err := client.NewThreatdetectionClient()
		if err != nil {
			return WrapError(err)
		}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

	return resourceAlicloudThreatDetectionVulWhitelistRead(d, meta)
}

func resourceAlicloudThreatDetectionVulWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteVulWhitelist"
	var response map[string]interface{}

	conn, err := client.NewThreatdetectionClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
