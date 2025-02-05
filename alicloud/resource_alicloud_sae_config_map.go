package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSaeConfigMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeConfigMapCreate,
		Read:   resourceAlicloudSaeConfigMapRead,
		Update: resourceAlicloudSaeConfigMapUpdate,
		Delete: resourceAlicloudSaeConfigMapDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudSaeConfigMapCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/sam/configmap/configMap"
	request := make(map[string]*string)
	var err error
	request["Data"] = StringPointer(d.Get("data").(string))
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	request["Name"], request["NamespaceId"] = StringPointer(d.Get("name").(string)), StringPointer(d.Get("namespace_id").(string))
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_config_map", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["ConfigMapId"]))

	return resourceAlicloudSaeConfigMapRead(d, meta)
}
func resourceAlicloudSaeConfigMapRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeConfigMap(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_config_map saeService.DescribeSaeConfigMap Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	resp, err := json.Marshal(object["Data"])
	if err != nil {
		return WrapError(err)
	}
	d.Set("data", string(resp))
	d.Set("description", object["Description"])
	d.Set("name", object["Name"])
	d.Set("namespace_id", object["NamespaceId"])
	return nil
}
func resourceAlicloudSaeConfigMapUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]*string{
		"ConfigMapId": StringPointer(d.Id()),
	}
	if d.HasChange("data") {
		update = true
	}
	request["Data"] = StringPointer(d.Get("data").(string))
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = StringPointer(v.(string))
	}
	if update {
		action := "/pop/v1/sam/configmap/configMap"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
				return WrapErrorf(Error(GetNotFoundMessage("SAE:ConfigMap", d.Id())), NotFoundMsg, ProviderERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
	}
	return resourceAlicloudSaeConfigMapRead(d, meta)
}
func resourceAlicloudSaeConfigMapDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/configmap/configMap"
	var response map[string]interface{}
	var err error
	request := map[string]*string{
		"ConfigMapId": StringPointer(d.Id()),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return nil
}
