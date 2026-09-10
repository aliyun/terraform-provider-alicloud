// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEbsSolutionInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEbsSolutionInstanceCreate,
		Read:   resourceAliCloudEbsSolutionInstanceRead,
		Update: resourceAliCloudEbsSolutionInstanceUpdate,
		Delete: resourceAliCloudEbsSolutionInstanceDelete,
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
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"solution_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"solution_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEbsSolutionInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSolutionInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["SolutionId"] = d.Get("solution_id")
	if v, ok := d.GetOk("solution_instance_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request["Tags"] = tagsMap
	}

	if v, ok := d.GetOk("parameters"); ok {
		parametersMaps := make([]map[string]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["ParameterKey"] = dataLoop1Tmp["parameter_key"]
			dataLoop1Map["ParameterValue"] = dataLoop1Tmp["parameter_value"]
			parametersMaps = append(parametersMaps, dataLoop1Map)
		}
		request["Parameters"], _ = convertListMapToJsonString(parametersMaps)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ebs_solution_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	ebsServiceV2 := EbsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"CREATE_COMPLETE"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, ebsServiceV2.EbsSolutionInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEbsSolutionInstanceRead(d, meta)
}

func resourceAliCloudEbsSolutionInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := EbsServiceV2{client}

	objectRaw, err := ebsServiceV2.DescribeEbsSolutionInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_solution_instance DescribeEbsSolutionInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateAt"])
	d.Set("description", objectRaw["Description"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("solution_id", objectRaw["SolutionId"])
	d.Set("solution_instance_name", objectRaw["Name"])
	d.Set("status", objectRaw["Status"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudEbsSolutionInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	action := "UpdateSolutionInstanceAttribute"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["SolutionInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("solution_instance_name") {
		update = true
		request["Name"] = d.Get("solution_instance_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("description")
		d.SetPartial("solution_instance_name")
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "solutioninstance"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}

	if d.HasChange("tags") {
		ebsServiceV2 := EbsServiceV2{client}
		if err := ebsServiceV2.SetResourceTags(d, "solutioninstance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAliCloudEbsSolutionInstanceRead(d, meta)
}

func resourceAliCloudEbsSolutionInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSolutionInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["SolutionInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ebsServiceV2 := EbsServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 20*time.Second, ebsServiceV2.EbsSolutionInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
