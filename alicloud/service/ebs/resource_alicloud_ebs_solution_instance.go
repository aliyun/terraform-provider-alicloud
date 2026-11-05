package ebs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSolutionInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceSolutionInstanceCreate,
		Read:   resourceSolutionInstanceRead,
		Update: resourceSolutionInstanceUpdate,
		Delete: resourceSolutionInstanceDelete,
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
			"tags": tags.TagsSchemaWithElements(),
		},
	}
}

func resourceSolutionInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateSolutionInstance"
	var response map[string]interface{}
	query := make(map[string]interface{})
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"ClientToken": helper.BuildClientToken(action),
		"SolutionId":  d.Get("solution_id"),
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("solution_instance_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("parameters"); ok {
		parametersMaps := make([]map[string]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			parametersMaps = append(parametersMaps, map[string]interface{}{
				"ParameterKey":   dataLoop1Tmp["parameter_key"],
				"ParameterValue": dataLoop1Tmp["parameter_value"],
			})
		}
		request["Parameters"], _ = convert.ConvertListMapToJsonString(parametersMaps)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var err error
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ebs_solution_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["InstanceId"]))

	ebsServiceV2 := NewEbsServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"CREATE_COMPLETE"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, ebsServiceV2.EbsSolutionInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	tags, _ := d.Get("tags").(map[string]interface{})
	if err := ebsServiceV2.TagResources(d, "solutioninstance", tags); err != nil {
		return errs.WrapError(err)
	}

	return resourceSolutionInstanceRead(d, meta)
}

func resourceSolutionInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsServiceV2 := NewEbsServiceV2(client)

	objectRaw, err := ebsServiceV2.DescribeEbsSolutionInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_solution_instance DescribeEbsSolutionInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateAt"])
	d.Set("description", objectRaw["Description"])
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("solution_id", objectRaw["SolutionId"])
	d.Set("solution_instance_name", objectRaw["Name"])
	d.Set("status", objectRaw["Status"])

	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceSolutionInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
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
	request["ClientToken"] = helper.BuildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("solution_instance_name") {
		update = true
		request["Name"] = d.Get("solution_instance_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			request["ClientToken"] = helper.BuildClientToken(action)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	action = "ChangeResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ResourceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	if _, ok := d.GetOk("resource_group_id"); ok && d.HasChange("resource_group_id") {
		update = true
		request["NewResourceGroupId"] = d.Get("resource_group_id")
	}

	request["ResourceType"] = "solutioninstance"
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
			request["ClientToken"] = helper.BuildClientToken(action)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		ebsServiceV2 := NewEbsServiceV2(client)
		if err := ebsServiceV2.SetResourceTags(d, "solutioninstance"); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Partial(false)
	return resourceSolutionInstanceRead(d, meta)
}

func resourceSolutionInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSolutionInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["SolutionInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = helper.BuildClientToken(action)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ebs", "2021-07-30", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	ebsServiceV2 := NewEbsServiceV2(client)
	stateConf := retry.StateChangeConf{
		Delay:                     5 * time.Second,
		Pending:                   []string{},
		Refresh:                   ebsServiceV2.EbsSolutionInstanceStateRefreshFunc(d.Id(), "Status", []string{}),
		Target:                    []string{"DELETE_COMPLETE"},
		Timeout:                   d.Timeout(schema.TimeoutDelete),
		ContinuousTargetOccurence: 1,
	}
	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
