// Package apig. This file is generated automatically. Please do not modify it manually, thank you!
package apig

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
)

func ResourceHttpApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceHttpApiCreate,
		Read:   resourceHttpApiRead,
		Update: resourceHttpApiUpdate,
		Delete: resourceHttpApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ai_protocols": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deploy_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"http_api_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"model_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocols": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceHttpApiCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/http-apis")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = helper.StringPointer(client.RegionId)

	if v, ok := d.GetOk("protocols"); ok {
		protocolsMapsArray := convert.ConvertToInterfaceArray(v)

		request["protocols"] = protocolsMapsArray
	}

	if v, ok := d.GetOkExists("enable_auth"); ok {
		request["enableAuth"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("ai_protocols"); ok {
		aiProtocolsMapsArray := convert.ConvertToInterfaceArray(v)

		request["aiProtocols"] = aiProtocolsMapsArray
	}

	if v, ok := d.GetOk("deploy_configs"); ok {
		deployConfigsMapsArray := make([]interface{}, 0)
		for _, item := range v.([]interface{}) {
			deployConfigMap := make(map[string]interface{})
			if err := json.Unmarshal([]byte(item.(string)), &deployConfigMap); err != nil {
				return errs.WrapError(err)
			}
			deployConfigsMapsArray = append(deployConfigsMapsArray, deployConfigMap)
		}
		request["deployConfigs"] = deployConfigsMapsArray
	}

	if v, ok := d.GetOk("type"); ok {
		request["type"] = v
	}
	request["name"] = d.Get("http_api_name")
	if v, ok := d.GetOk("base_path"); ok {
		request["basePath"] = v
	}
	if v, ok := d.GetOk("model_category"); ok {
		request["modelCategory"] = v
	}
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_apig_http_api", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.httpApiId", response)
	d.SetId(fmt.Sprint(id))

	return resourceHttpApiRead(d, meta)
}

func resourceHttpApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := NewService(client)

	objectRaw, err := apigServiceV2.DescribeApigHttpApi(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_http_api DescribeApigHttpApi Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["name"] != nil {
		d.Set("http_api_name", objectRaw["name"])
	}
	if objectRaw["resourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["resourceGroupId"])
	}
	if objectRaw["type"] != nil {
		d.Set("type", objectRaw["type"])
	}
	if objectRaw["enableAuth"] != nil {
		d.Set("enable_auth", objectRaw["enableAuth"])
	}

	protocolsRaw := make([]interface{}, 0)
	if objectRaw["protocols"] != nil {
		protocolsRaw = convert.ConvertToInterfaceArray(objectRaw["protocols"])
	}

	d.Set("protocols", protocolsRaw)

	return nil
}

func resourceHttpApiUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	httpApiId := d.Id()
	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("protocols") {
		update = true
	}
	if v, ok := d.GetOk("protocols"); ok || d.HasChange("protocols") {
		protocolsMapsArray := convert.ConvertToInterfaceArray(v)

		request["protocols"] = protocolsMapsArray
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if !d.IsNewResource() && d.HasChange("deploy_configs") {
		update = true
		if v, ok := d.GetOk("deploy_configs"); ok || d.HasChange("deploy_configs") {
			deployConfigsMapsArray := make([]interface{}, 0)
			for _, item := range v.([]interface{}) {
				deployConfigMap := make(map[string]interface{})
				if err := json.Unmarshal([]byte(item.(string)), &deployConfigMap); err != nil {
					return errs.WrapError(err)
				}
				deployConfigsMapsArray = append(deployConfigsMapsArray, deployConfigMap)
			}
			request["deployConfigs"] = deployConfigsMapsArray
		}
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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
	update = false
	action = fmt.Sprintf("/move-resource-group")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = helper.StringPointer(d.Id())
	query["RegionId"] = helper.StringPointer(client.RegionId)
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["ResourceGroupId"] = helper.StringPointer(v.(string))
	}

	query["ResourceType"] = helper.StringPointer("HttpApi")
	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		apigServiceV2 := NewService(client)
		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 35*time.Second, apigServiceV2.ApigHttpApiStateRefreshFunc(d.Id(), "resourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceHttpApiRead(d, meta)
}

func resourceHttpApiDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	httpApiId := d.Id()
	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
		if errs.IsExpectedErrors(err, []string{"Error.DatabaseError.RecordNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
