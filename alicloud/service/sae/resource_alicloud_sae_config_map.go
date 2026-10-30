package sae

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	request["Data"] = helper.StringPointer(d.Get("data").(string))
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = helper.StringPointer(v.(string))
	}
	request["Name"], request["NamespaceId"] = helper.StringPointer(d.Get("name").(string)), helper.StringPointer(d.Get("namespace_id").(string))
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_sae_config_map", "POST "+action, errs.AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["ConfigMapId"]))

	return resourceAlicloudSaeConfigMapRead(d, meta)
}
func resourceAlicloudSaeConfigMapRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := NewService(client)
	object, err := saeService.DescribeSaeConfigMap(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_config_map saeService.DescribeSaeConfigMap Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	resp, err := json.Marshal(object["Data"])
	if err != nil {
		return errs.WrapError(err)
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
		"ConfigMapId": helper.StringPointer(d.Id()),
	}
	if d.HasChange("data") {
		update = true
	}
	request["Data"] = helper.StringPointer(d.Get("data").(string))
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = helper.StringPointer(v.(string))
	}
	if update {
		action := "/pop/v1/sam/configmap/configMap"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("sae", "2019-05-06", action, request, nil, nil, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
				return errs.WrapErrorf(errs.NotFoundErr("SAE:ConfigMap", d.Id()), errs.NotFoundMsg, errs.ProviderERROR)
			}
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "PUT "+action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)
	}
	return resourceAlicloudSaeConfigMapRead(d, meta)
}
func resourceAlicloudSaeConfigMapDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/configmap/configMap"
	var response map[string]interface{}
	var err error
	request := map[string]*string{
		"ConfigMapId": helper.StringPointer(d.Id()),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("sae", "2019-05-06", action, request, nil, nil, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFound.ConfigMap"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DELETE "+action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	return nil
}
