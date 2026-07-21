// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/sjson"
)

func resourceAliCloudThreatDetectionAssetBind() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionAssetBindCreate,
		Read:   resourceAliCloudThreatDetectionAssetBindRead,
		Update: resourceAliCloudThreatDetectionAssetBindUpdate,
		Delete: resourceAliCloudThreatDetectionAssetBindDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auth_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudThreatDetectionAssetBindCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "UpdatePostPaidBindRel"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOkExists("auth_version"); ok {
		objectDataLocalMap["Version"] = v
	}

	if v, ok := d.GetOkExists("uuid"); ok {
		objectDataLocalMap["UuidList"] = v
	}

	BindActionMap := make([]interface{}, 0)
	BindActionMap = append(BindActionMap, objectDataLocalMap)
	request["BindAction"] = BindActionMap

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "BindAction.0.UuidList.0", d.Get("uuid"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_asset_bind", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("BindAction[0].UuidList[0]", request)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionAssetBindRead(d, meta)
}

func resourceAliCloudThreatDetectionAssetBindRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionAssetBind(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_asset_bind DescribeThreatDetectionAssetBind Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auth_version", objectRaw["AuthVersion"])
	d.Set("uuid", objectRaw["Uuid"])

	d.Set("uuid", d.Id())

	return nil
}

func resourceAliCloudThreatDetectionAssetBindUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdatePostPaidBindRel"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("auth_version") {
		update = true
		objectDataLocalMap["Version"] = d.Get("auth_version")
	}

	if d.HasChange("uuid") {
		update = true
		objectDataLocalMap["UuidList"] = d.Get("uuid")
	}

	BindActionMap := make([]interface{}, 0)
	BindActionMap = append(BindActionMap, objectDataLocalMap)
	request["BindAction"] = BindActionMap

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "BindAction.0.UuidList.0", d.Id())
	_ = json.Unmarshal([]byte(jsonString), &request)

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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

	return resourceAliCloudThreatDetectionAssetBindRead(d, meta)
}

func resourceAliCloudThreatDetectionAssetBindDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Asset Bind. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
