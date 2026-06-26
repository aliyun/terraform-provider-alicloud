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
)

func resourceAliCloudApigHttpApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigHttpApiCreate,
		Read:   resourceAliCloudApigHttpApiRead,
		Update: resourceAliCloudApigHttpApiUpdate,
		Delete: resourceAliCloudApigHttpApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
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
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudApigHttpApiCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/http-apis")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	if v, ok := d.GetOk("protocols"); ok {
		protocolsMapsArray := convertToInterfaceArray(v)

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
		aiProtocolsMapsArray := convertToInterfaceArray(v)

		request["aiProtocols"] = aiProtocolsMapsArray
	}

	if v, ok := d.GetOk("deploy_configs"); ok {
		deployConfigsMapsArray := make([]interface{}, 0)
		for _, item := range v.([]interface{}) {
			deployConfigMap := make(map[string]interface{})
			if err := json.Unmarshal([]byte(item.(string)), &deployConfigMap); err != nil {
				return WrapError(err)
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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_http_api", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.httpApiId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigHttpApiUpdate(d, meta)
}

func resourceAliCloudApigHttpApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigHttpApi(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_http_api DescribeApigHttpApi Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["name"] != nil {
		d.Set("http_api_name", objectRaw["name"])
	}
	if objectRaw["type"] != nil {
		d.Set("type", objectRaw["type"])
	}
	if objectRaw["enableAuth"] != nil {
		d.Set("enable_auth", objectRaw["enableAuth"])
	}

	protocolsRaw := make([]interface{}, 0)
	if objectRaw["protocols"] != nil {
		protocolsRaw = convertToInterfaceArray(objectRaw["protocols"])
	}

	d.Set("protocols", protocolsRaw)

	return nil
}

func resourceAliCloudApigHttpApiUpdate(d *schema.ResourceData, meta interface{}) error {
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
		protocolsMapsArray := convertToInterfaceArray(v)

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
					return WrapError(err)
				}
				deployConfigsMapsArray = append(deployConfigsMapsArray, deployConfigMap)
			}
			request["deployConfigs"] = deployConfigsMapsArray
		}
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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
	update = false
	action = fmt.Sprintf("/move-resource-group")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	query["ResourceType"] = StringPointer("HttpApi")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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

	return resourceAliCloudApigHttpApiRead(d, meta)
}

func resourceAliCloudApigHttpApiDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	httpApiId := d.Id()
	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
		if IsExpectedErrors(err, []string{"Error.DatabaseError.RecordNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
