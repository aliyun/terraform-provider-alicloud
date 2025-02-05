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

func resourceAlicloudThreatDetectionHoneyPot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudThreatDetectionHoneyPotCreate,
		Read:   resourceAlicloudThreatDetectionHoneyPotRead,
		Update: resourceAlicloudThreatDetectionHoneyPotUpdate,
		Delete: resourceAlicloudThreatDetectionHoneyPotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"honeypot_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"honeypot_image_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"honeypot_image_name": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"honeypot_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"node_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"preset_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"state": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudThreatDetectionHoneyPotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("honeypot_image_id"); ok {
		request["HoneypotImageId"] = v
	}
	if v, ok := d.GetOk("honeypot_image_name"); ok {
		request["HoneypotImageName"] = v
	}
	if v, ok := d.GetOk("honeypot_name"); ok {
		request["HoneypotName"] = v
	}
	if v, ok := d.GetOk("node_id"); ok {
		request["NodeId"] = v
	}

	var response map[string]interface{}
	action := "CreateHoneypot"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_honey_pot", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.Data.HoneypotId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_threat_detection_honey_pot")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	return resourceAlicloudThreatDetectionHoneyPotRead(d, meta)
}

func resourceAlicloudThreatDetectionHoneyPotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}

	object, err := sasService.DescribeThreatDetectionHoneyPot(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_honey_pot sasService.DescribeThreatDetectionHoneyPot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("honeypot_image_id", object["HoneypotImageId"])
	d.Set("honeypot_image_name", object["HoneypotImageName"])
	d.Set("honeypot_name", object["HoneypotName"])
	d.Set("node_id", object["NodeId"])
	d.Set("preset_id", object["PresetId"])
	state, _ := jsonpath.Get("$.State", object)
	d.Set("state", state)
	d.Set("status", state.([]interface{})[0])

	return nil
}

func resourceAlicloudThreatDetectionHoneyPotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{
		"HoneypotId": d.Id(),
	}

	if d.HasChange("honeypot_name") {
		update = true
		if v, ok := d.GetOk("honeypot_name"); ok {
			request["HoneypotName"] = v
		}
	}

	if update {
		action := "UpdateHoneypot"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudThreatDetectionHoneyPotRead(d, meta)
}

func resourceAlicloudThreatDetectionHoneyPotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sasService := SasService{client}
	var err error

	request := map[string]interface{}{
		"HoneypotId": d.Id(),
	}

	action := "DeleteHoneypot"
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, sasService.ThreatDetectionHoneyPotStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
