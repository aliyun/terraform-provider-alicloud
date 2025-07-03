package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaCustomScenePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaCustomScenePolicyCreate,
		Read:   resourceAliCloudEsaCustomScenePolicyRead,
		Update: resourceAliCloudEsaCustomScenePolicyUpdate,
		Delete: resourceAliCloudEsaCustomScenePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_scene_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"site_ids": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudEsaCustomScenePolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCustomScenePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Template"] = d.Get("template")
	request["StartTime"] = d.Get("create_time")
	request["Name"] = d.Get("custom_scene_policy_name")
	request["SiteIds"] = d.Get("site_ids")
	request["EndTime"] = d.Get("end_time")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_custom_scene_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PolicyId"]))

	return resourceAliCloudEsaCustomScenePolicyUpdate(d, meta)
}

func resourceAliCloudEsaCustomScenePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaCustomScenePolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_custom_scene_policy DescribeEsaCustomScenePolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["StartTime"])
	d.Set("custom_scene_policy_name", objectRaw["Name"])
	d.Set("end_time", objectRaw["EndTime"])
	d.Set("status", objectRaw["Status"])
	d.Set("template", objectRaw["Template"])
	d.Set("site_ids", objectRaw["SiteIds"])

	return nil
}

func resourceAliCloudEsaCustomScenePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	if d.HasChange("status") {
		var err error
		esaServiceV2 := EsaServiceV2{client}
		object, err := esaServiceV2.DescribeEsaCustomScenePolicy(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Disabled" {
				action := "DisableCustomScenePolicy"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["PolicyId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
				esaServiceV2 := EsaServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Disabled"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, esaServiceV2.EsaCustomScenePolicyStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
			if target == "Running" {
				action := "EnableCustomScenePolicy"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["PolicyId"] = d.Id()
				request["RegionId"] = client.RegionId
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
				esaServiceV2 := EsaServiceV2{client}
				stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, esaServiceV2.EsaCustomScenePolicyStateRefreshFunc(d.Id(), "Status", []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}

			}
		}
	}

	var err error
	action := "UpdateCustomScenePolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PolicyId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("template") {
		update = true
	}
	request["Template"] = d.Get("template")
	if !d.IsNewResource() && d.HasChange("create_time") {
		update = true
	}
	request["StartTime"] = d.Get("create_time")
	if !d.IsNewResource() && d.HasChange("custom_scene_policy_name") {
		update = true
	}
	request["Name"] = d.Get("custom_scene_policy_name")
	if !d.IsNewResource() && d.HasChange("site_ids") {
		update = true
	}
	request["SiteIds"] = d.Get("site_ids")
	if !d.IsNewResource() && d.HasChange("end_time") {
		update = true
	}
	request["EndTime"] = d.Get("end_time")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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

	return resourceAliCloudEsaCustomScenePolicyRead(d, meta)
}

func resourceAliCloudEsaCustomScenePolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCustomScenePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PolicyId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)

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

	return nil
}
