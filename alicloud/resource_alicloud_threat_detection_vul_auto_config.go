package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudThreatDetectionVulAutoConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionVulAutoConfigCreate,
		Read:   resourceAlicloudThreatDetectionVulAutoConfigRead,
		Update: resourceAlicloudThreatDetectionVulAutoConfigUpdate,
		Delete: resourceAlicloudThreatDetectionVulAutoConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"necessity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"all_uuid": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"need_snapshot": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"snapshot_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"target_end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudThreatDetectionVulAutoConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddOrUpdateAutoFixConfig"
	request := make(map[string]interface{})
	var err error

	request["Type"] = d.Get("type")
	request["StartTime"] = d.Get("start_time")
	request["AllUuid"] = d.Get("all_uuid")
	request["NeedSnapshot"] = d.Get("need_snapshot")
	request["Enable"] = d.Get("enable")

	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("necessity"); ok {
		request["Necessity"] = v
	}
	if v, ok := d.GetOk("target_start_time"); ok {
		request["TargetStartTime"] = v
	}
	if v, ok := d.GetOk("snapshot_name"); ok {
		request["SnapshotName"] = v
	}
	if v, ok := d.GetOk("snapshot_time"); ok {
		request["SnapshotTime"] = v
	}
	if v, ok := d.GetOk("target_end_time"); ok {
		request["TargetEndTime"] = v
	}
	if v, ok := d.GetOk("rules"); ok {
		request["Rules"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_vul_auto_config", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ConfigId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_vul_auto_config")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudThreatDetectionVulAutoConfigRead(d, meta)
}

func resourceAlicloudThreatDetectionVulAutoConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionService := ThreatDetectionService{client}

	object, err := threatDetectionService.DescribeThreatDetectionVulAutoConfig(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_vul_auto_config DescribeThreatDetectionVulAutoConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config_id", object["ConfigId"])
	d.Set("type", object["Type"])
	d.Set("start_time", object["StartTime"])
	d.Set("period_unit", object["PeriodUnit"])
	d.Set("necessity", object["Necessity"])
	d.Set("target_start_time", object["TargetStartTime"])
	d.Set("all_uuid", object["AllUuid"])
	d.Set("need_snapshot", object["NeedSnapshot"])
	d.Set("snapshot_name", object["SnapshotName"])
	d.Set("snapshot_time", object["SnapshotTime"])
	d.Set("target_end_time", object["TargetEndTime"])
	d.Set("enable", object["Enable"])
	d.Set("region_id", object["RegionId"])

	return nil
}

func resourceAlicloudThreatDetectionVulAutoConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error

	// If any field other than "enable" changed, perform a full upsert via AddOrUpdateAutoFixConfig.
	fullUpdate := d.HasChange("type") || d.HasChange("start_time") || d.HasChange("period_unit") ||
		d.HasChange("necessity") || d.HasChange("target_start_time") || d.HasChange("all_uuid") ||
		d.HasChange("need_snapshot") || d.HasChange("snapshot_name") || d.HasChange("snapshot_time") ||
		d.HasChange("target_end_time") || d.HasChange("rules")

	if fullUpdate {
		request := map[string]interface{}{
			"ConfigId":     d.Id(),
			"Type":         d.Get("type"),
			"StartTime":    d.Get("start_time"),
			"AllUuid":      d.Get("all_uuid"),
			"NeedSnapshot": d.Get("need_snapshot"),
			"Enable":       d.Get("enable"),
		}
		if v, ok := d.GetOk("period_unit"); ok {
			request["PeriodUnit"] = v
		}
		if v, ok := d.GetOk("necessity"); ok {
			request["Necessity"] = v
		}
		if v, ok := d.GetOk("target_start_time"); ok {
			request["TargetStartTime"] = v
		}
		if v, ok := d.GetOk("snapshot_name"); ok {
			request["SnapshotName"] = v
		}
		if v, ok := d.GetOk("snapshot_time"); ok {
			request["SnapshotTime"] = v
		}
		if v, ok := d.GetOk("target_end_time"); ok {
			request["TargetEndTime"] = v
		}
		if v, ok := d.GetOk("rules"); ok {
			request["Rules"] = v
		}
		action := "AddOrUpdateAutoFixConfig"
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
	} else if d.HasChange("enable") {
		// Only the enable flag changed; use the lightweight status toggle API.
		request := map[string]interface{}{
			"ConfigId": d.Id(),
			"Enable":   d.Get("enable"),
		}
		action := "UpdateAutoFixConfigStatus"
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
	}

	return resourceAlicloudThreatDetectionVulAutoConfigRead(d, meta)
}

func resourceAlicloudThreatDetectionVulAutoConfigDelete(d *schema.ResourceData, meta interface{}) error {
	// The ThreatDetection VulAutoConfig API does not provide a delete operation.
	// AddOrUpdateAutoFixConfig is an upsert (add or update) and there is no
	// DeleteAutoFixConfig API. Removing the resource from Terraform state only
	// detaches management; the cloud-side configuration is preserved.
	// Users can disable the rule by setting enable = 0 instead.
	return nil
}
