// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package tag

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

func ResourceAssociatedRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAssociatedRuleCreate,
		Read:   ResourceAssociatedRuleRead,
		Update: ResourceAssociatedRuleUpdate,
		Delete: ResourceAssociatedRuleDelete,
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

func ResourceAssociatedRuleCreate(d *schema.ResourceData, meta interface{}) error {

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

	jsonString := convert.ConvertObjectToJsonString(request)
	jsonString, _ = sjson.Set(jsonString, "CreateRulesList.0.SettingName", d.Get("associated_setting_name"))
	_ = json.Unmarshal([]byte(jsonString), &request)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_tag_associated_rule", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("CreateRulesList[0].SettingName", request)
	d.SetId(fmt.Sprint(id))

	return ResourceAssociatedRuleRead(d, meta)
}

func ResourceAssociatedRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagServiceV2 := NewTagServiceV2(client)

	objectRaw, err := tagServiceV2.DescribeTagAssociatedRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_tag_associated_rule DescribeTagAssociatedRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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

func ResourceAssociatedRuleUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return ResourceAssociatedRuleRead(d, meta)
}

func ResourceAssociatedRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAssociatedResourceRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SettingName"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidSettingName.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
