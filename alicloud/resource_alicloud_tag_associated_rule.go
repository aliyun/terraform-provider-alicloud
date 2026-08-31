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

func resourceAliCloudTagAssociatedRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudTagAssociatedRuleCreate,
		Read:   resourceAliCloudTagAssociatedRuleRead,
		Update: resourceAliCloudTagAssociatedRuleUpdate,
		Delete: resourceAliCloudTagAssociatedRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"associated_setting_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudTagAssociatedRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAssociatedResourceRules"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("tag_keys"); ok {
		tagKeys1, _ := jsonpath.Get("$", v)
		if tagKeys1 != nil && tagKeys1 != "" {
			objectDataLocalMap["TagKeys"] = tagKeys1
		}
	}

	if v, ok := d.GetOk("status"); ok {
		objectDataLocalMap["Status"] = v
	}

	if v, ok := d.GetOk("associated_setting_name"); ok {
		objectDataLocalMap["SettingName"] = v
	}

	CreateRulesListMap := make([]interface{}, 0)
	CreateRulesListMap = append(CreateRulesListMap, objectDataLocalMap)
	request["CreateRulesList"] = CreateRulesListMap

	jsonString := convertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "CreateRulesList.0.SettingName", d.Get("associated_setting_name"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_tag_associated_rule", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("CreateRulesList[0].SettingName", request)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudTagAssociatedRuleRead(d, meta)
}

func resourceAliCloudTagAssociatedRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagServiceV2 := TagServiceV2{client}

	objectRaw, err := tagServiceV2.DescribeTagAssociatedRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_tag_associated_rule DescribeTagAssociatedRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("status", objectRaw["Status"])
	d.Set("associated_setting_name", objectRaw["SettingName"])

	tagKeysRaw := make([]interface{}, 0)
	if objectRaw["TagKeys"] != nil {
		tagKeysRaw = objectRaw["TagKeys"].([]interface{})
	}

	d.Set("tag_keys", tagKeysRaw)

	return nil
}

func resourceAliCloudTagAssociatedRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAssociatedResourceRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SettingName"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("tag_keys") {
		update = true
	}
	if v, ok := d.GetOk("tag_keys"); ok {
		tagKeysMapsArray := v.([]interface{})
		request["TagKeys"] = tagKeysMapsArray
	}

	if d.HasChange("status") {
		update = true
	}
	request["Status"] = d.Get("status")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, query, request, true)
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

	return resourceAliCloudTagAssociatedRuleRead(d, meta)
}

func resourceAliCloudTagAssociatedRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAssociatedResourceRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SettingName"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidSettingName.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
