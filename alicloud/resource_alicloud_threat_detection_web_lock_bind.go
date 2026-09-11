package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudThreatDetectionWebLockBind() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionWebLockBindCreate,
		Read:   resourceAlicloudThreatDetectionWebLockBindRead,
		Update: resourceAlicloudThreatDetectionWebLockBindUpdate,
		Delete: resourceAlicloudThreatDetectionWebLockBindDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"exclusive_file": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"lang": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
				Type:         schema.TypeString,
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
			"status": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
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

func resourceAlicloudThreatDetectionWebLockBindCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})

	request["DefenceMode"] = d.Get("defence_mode")
	request["Dir"] = d.Get("dir")
	// ExclusiveDir and ExclusiveFile are rejected with IllegalParam unless sent
	// as JSON array strings (verified against the live API).
	if v, ok := d.GetOk("exclusive_dir"); ok {
		request["ExclusiveDir"] = webLockBindToJsonArray(v.([]interface{}))
	}
	if v, ok := d.GetOk("exclusive_file"); ok {
		request["ExclusiveFile"] = webLockBindToJsonArray(v.([]interface{}))
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
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"BindDataExist"}) {
				return nil
			}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_web_lock_bind", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Uuid"]))

	// status and lang are update-only fields (ModifyWebLockStatus); apply them
	// after the bind is started so the desired state is reached in one apply.
	if v, ok := d.GetOk("status"); ok {
		if err := webLockBindUpdateStatus(client, d.Id(), fmt.Sprint(v), d.Get("lang")); err != nil {
			return err
		}
	} else if v, ok := d.GetOk("lang"); ok {
		if err := webLockBindUpdateStatus(client, d.Id(), "", fmt.Sprint(v)); err != nil {
			return err
		}
	}

	return resourceAlicloudThreatDetectionWebLockBindRead(d, meta)
}

func resourceAlicloudThreatDetectionWebLockBindRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	object, err := sasService.DescribeThreatDetectionWebLockBind(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_web_lock_bind sasService.DescribeThreatDetectionWebLockBind Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("uuid", object["Uuid"])
	d.Set("status", object["Status"])

	// DescribeWebLockBindList only returns the bind status; the protection
	// config (dir, mode, backup dir, file types, ...) is served by
	// DescribeWebLockConfigList.
	config, err := sasService.DescribeThreatDetectionWebLockConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_web_lock_bind sasService.DescribeThreatDetectionWebLockConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("defence_mode", config["DefenceMode"])
	d.Set("dir", config["Dir"])
	d.Set("exclusive_dir", webLockBindFromJsonArray(config["ExclusiveDir"]))
	d.Set("exclusive_file", webLockBindFromJsonArray(config["ExclusiveFile"]))
	d.Set("exclusive_file_type", config["ExclusiveFileType"])
	d.Set("inclusive_file_type", config["InclusiveFileType"])
	d.Set("local_backup_dir", config["LocalBackupDir"])
	d.Set("mode", config["Mode"])

	return nil
}

func resourceAlicloudThreatDetectionWebLockBindUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// ModifyWebLockStatus only updates status and lang; all other fields are
	// ForceNew and trigger recreation.
	if d.HasChange("status") || d.HasChange("lang") {
		if err := webLockBindUpdateStatus(client, d.Id(), d.Get("status"), d.Get("lang")); err != nil {
			return err
		}
	}

	return resourceAlicloudThreatDetectionWebLockBindRead(d, meta)
}

func resourceAlicloudThreatDetectionWebLockBindDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	if _, err := sasService.DescribeThreatDetectionWebLockBind(d.Id()); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	// ModifyWebLockUnbind rejects the request with "StillOn" while the
	// protection is still switched on; turn it off first and wait for the
	// bind entry to report status=off.
	if err := webLockBindUpdateStatus(client, d.Id(), "off", nil); err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"Asset not bind.", "AssetNotBind"}) {
			return nil
		}
		return WrapError(err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"on"},
		Target:     []string{"off"},
		Refresh:    webLockBindStatusRefreshFunc(sasService, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	request := map[string]interface{}{
		"Uuid": d.Id(),
	}

	action := "ModifyWebLockUnbind"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"Asset not bind.", "AssetNotBind"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	// Unbind is asynchronous: the bind entry keeps showing up in
	// DescribeWebLockBindList for a while. Wait until it is really gone so an
	// immediate re-create of the same server does not hit a stale entry.
	removeConf := &resource.StateChangeConf{
		Pending:    []string{"on", "off"},
		Target:     []string{"removed"},
		Refresh:    webLockBindRemovedRefreshFunc(sasService, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := removeConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

// webLockBindRemovedRefreshFunc reports "removed" once the bind entry no
// longer appears in DescribeWebLockBindList.
func webLockBindRemovedRefreshFunc(sasService SasService, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := sasService.DescribeThreatDetectionWebLockBind(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "removed", nil
			}
			return nil, "", err
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

// webLockBindStatusRefreshFunc polls DescribeWebLockBindList for the bind
// entry status; a missing entry is reported as "off" so deletion can proceed.
func webLockBindStatusRefreshFunc(sasService SasService, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := sasService.DescribeThreatDetectionWebLockBind(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "off", nil
			}
			return nil, "", err
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

// webLockBindToJsonArray marshals a schema list into the JSON array string
// form that ExclusiveDir/ExclusiveFile require.
func webLockBindToJsonArray(list []interface{}) string {
	items := make([]string, 0, len(list))
	for _, item := range list {
		items = append(items, fmt.Sprint(item))
	}
	out, _ := json.Marshal(items)
	return string(out)
}

// webLockBindFromJsonArray decodes the JSON array string returned for
// ExclusiveDir/ExclusiveFile back into a list; a non-empty plain string is
// returned as a single-element list for robustness.
func webLockBindFromJsonArray(v interface{}) []interface{} {
	s := fmt.Sprint(v)
	if s == "" || s == "<nil>" {
		return nil
	}
	var items []interface{}
	if err := json.Unmarshal([]byte(s), &items); err == nil {
		return items
	}
	return []interface{}{s}
}

// webLockBindUpdateStatus calls ModifyWebLockStatus to update the protection
// status and/or language of a web lock bind entry.
func webLockBindUpdateStatus(client *connectivity.AliyunClient, id, status, lang interface{}) error {
	request := map[string]interface{}{
		"Uuid": id,
	}
	if status != nil && fmt.Sprint(status) != "" {
		request["Status"] = status
	}
	if lang != nil && fmt.Sprint(lang) != "" {
		request["Lang"] = lang
	}

	action := "ModifyWebLockStatus"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
