package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionWebLockConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionWebLockConfigCreate,
		Read:   resourceAlicloudThreatDetectionWebLockConfigRead,
		Delete: resourceAlicloudThreatDetectionWebLockConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"defence_mode": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"block", "audit"}, false),
				Type:         schema.TypeString,
			},
			"dir": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"exclusive_dir": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"exclusive_file": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"exclusive_file_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"inclusive_file_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"local_backup_dir": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"mode": {
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"whitelist", "blacklist"}, false),
				Type:         schema.TypeString,
			},
			"uuid": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudThreatDetectionWebLockConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	request["DefenceMode"] = d.Get("defence_mode")

	request["Dir"] = d.Get("dir")
	if v, ok := d.GetOk("exclusive_dir"); ok {
		request["ExclusiveDir"] = v
	}
	if v, ok := d.GetOk("exclusive_file"); ok {
		request["ExclusiveFile"] = v
	}
	if v, ok := d.GetOk("exclusive_file_type"); ok {
		request["ExclusiveFileType"] = v
	}
	if v, ok := d.GetOk("inclusive_file_type"); ok {
		request["InclusiveFileType"] = v
	}
	request["LocalBackupDir"] = d.Get("local_backup_dir")
	request["Mode"] = d.Get("mode")
	request["Uuid"] = d.Get("uuid")

	var response map[string]interface{}
	action := "ModifyWebLockStart"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_web_lock_config", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Uuid"]))

	return resourceAlicloudThreatDetectionWebLockConfigRead(d, meta)
}

func resourceAlicloudThreatDetectionWebLockConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	object, err := sasService.DescribeThreatDetectionWebLockConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_web_lock_config sasService.DescribeThreatDetectionWebLockConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("uuid", object["Uuid"])
	d.Set("defence_mode", object["DefenceMode"])
	d.Set("dir", object["Dir"])
	d.Set("exclusive_dir", object["ExclusiveDir"])
	d.Set("exclusive_file", object["ExclusiveFile"])
	d.Set("exclusive_file_type", object["ExclusiveFileType"])
	d.Set("inclusive_file_type", object["InclusiveFileType"])
	d.Set("local_backup_dir", object["LocalBackupDir"])
	d.Set("mode", object["Mode"])

	return nil
}

func resourceAlicloudThreatDetectionWebLockConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{
		"Uuid": d.Id(),
	}

	action := "ModifyWebLockUnbind"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
