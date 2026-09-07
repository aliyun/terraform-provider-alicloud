package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudThreatDetectionBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionBackupPolicyCreate,
		Read:   resourceAlicloudThreatDetectionBackupPolicyRead,
		Update: resourceAlicloudThreatDetectionBackupPolicyUpdate,
		Delete: resourceAlicloudThreatDetectionBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"backup_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"1.0.0", "2.0.0"}, false),
			},
			"uuid_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policy_region_id": {
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

func resourceAlicloudThreatDetectionBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}
	var response map[string]interface{}
	action := "CreateBackupPolicy"
	request := make(map[string]interface{})
	var err error

	request["Name"] = d.Get("backup_policy_name")
	request["Policy"] = d.Get("policy")
	request["PolicyVersion"] = d.Get("policy_version")
	request["UuidList"] = d.Get("uuid_list")

	if v, ok := d.GetOk("policy_region_id"); ok {
		request["PolicyRegionId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_backup_policy", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.BackupPolicy", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_backup_policy")
	} else {
		d.SetId(fmt.Sprint(v.(map[string]interface{})["Id"]))
	}

	stateConf := BuildStateConf([]string{}, []string{"enabled"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, threatDetectionService.ThreatDetectionBackupPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudThreatDetectionBackupPolicyRead(d, meta)
}

func resourceAlicloudThreatDetectionBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}

	object, err := threatDetectionService.DescribeThreatDetectionBackupPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("backup_policy_name", object["Name"])
	d.Set("policy", object["Policy"])
	d.Set("policy_version", object["PolicyVersion"])
	d.Set("policy_region_id", object["RegionId"])
	d.Set("status", object["Status"])

	if uuidLists, ok := object["UuidList"]; ok {
		d.Set("uuid_list", uuidLists.([]interface{}))
	}

	return nil
}

func resourceAlicloudThreatDetectionBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	if d.HasChange("backup_policy_name") {
		update = true
	}
	request["Name"] = d.Get("backup_policy_name")

	if d.HasChange("policy") {
		update = true
	}
	request["Policy"] = d.Get("policy")

	if d.HasChange("uuid_list") {
		update = true
	}
	request["UuidList"] = d.Get("uuid_list")

	if d.HasChange("policy_region_id") {
		update = true
	}
	if v, ok := d.GetOk("policy_region_id"); ok {
		request["PolicyRegionId"] = v
	}

	if update {
		action := "ModifyBackupPolicy"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
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

		stateConf := BuildStateConf([]string{}, []string{"enabled"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, threatDetectionService.ThreatDetectionBackupPolicyStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudThreatDetectionBackupPolicyRead(d, meta)
}

func resourceAlicloudThreatDetectionBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}
	action := "DeleteBackupPolicy"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"Id": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, nil, request, false)
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, threatDetectionService.ThreatDetectionBackupPolicyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
