// Package governance. This file is generated automatically. Please do not modify it manually, thank you!
package governance

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGovernanceBaseline() *schema.Resource {
	return &schema.Resource{
		Create: resourceGovernanceBaselineCreate,
		Read:   resourceGovernanceBaselineRead,
		Update: resourceGovernanceBaselineUpdate,
		Delete: resourceGovernanceBaselineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"baseline_items": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validators.StringInSlice([]string{"1.0"}, false),
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"baseline_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGovernanceBaselineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccountFactoryBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["RegionId"] = client.RegionId

	if v, ok := d.GetOk("baseline_items"); ok {
		baselineItemsMaps := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Config"] = dataLoopTmp["config"]
			dataLoopMap["Name"] = dataLoopTmp["name"]
			dataLoopMap["Version"] = dataLoopTmp["version"]
			baselineItemsMaps = append(baselineItemsMaps, dataLoopMap)
		}
		request["BaselineItems"] = baselineItemsMaps
	}

	if v, ok := d.GetOk("baseline_name"); ok {
		request["BaselineName"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_governance_baseline", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BaselineId"]))

	return resourceGovernanceBaselineRead(d, meta)
}

func resourceGovernanceBaselineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)

	objectRaw, err := service.DescribeGovernanceBaseline(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_governance_baseline DescribeGovernanceBaseline Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["BaselineName"] != nil {
		d.Set("baseline_name", objectRaw["BaselineName"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}

	baselineItems1Raw := objectRaw["BaselineItems"]
	baselineItemsMaps := make([]map[string]interface{}, 0)
	if baselineItems1Raw != nil {
		for _, baselineItemsChild1Raw := range baselineItems1Raw.([]interface{}) {
			baselineItemsMap := make(map[string]interface{})
			baselineItemsChild1Raw := baselineItemsChild1Raw.(map[string]interface{})
			baselineItemsMap["config"] = baselineItemsChild1Raw["Config"]
			baselineItemsMap["name"] = baselineItemsChild1Raw["Name"]
			baselineItemsMap["version"] = baselineItemsChild1Raw["Version"]

			baselineItemsMaps = append(baselineItemsMaps, baselineItemsMap)
		}
	}
	if objectRaw["BaselineItems"] != nil {
		if err := d.Set("baseline_items", baselineItemsMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceGovernanceBaselineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateAccountFactoryBaseline"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["BaselineId"] = d.Id()
	query["RegionId"] = client.RegionId
	if d.HasChange("baseline_items") {
		update = true
		if v, ok := d.GetOk("baseline_items"); ok {
			baselineItemsMaps := make([]interface{}, 0)
			for _, dataLoop := range v.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Config"] = dataLoopTmp["config"]
				dataLoopMap["Name"] = dataLoopTmp["name"]
				dataLoopMap["Version"] = dataLoopTmp["version"]
				baselineItemsMaps = append(baselineItemsMaps, dataLoopMap)
			}
			request["BaselineItems"] = baselineItemsMaps
		}
	}

	if d.HasChange("baseline_name") {
		update = true
		request["BaselineName"] = d.Get("baseline_name")
	}

	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)
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

	return resourceGovernanceBaselineRead(d, meta)
}

func resourceGovernanceBaselineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccountFactoryBaseline"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["BaselineId"] = d.Id()
	query["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("governance", "2021-01-20", action, query, request, false)

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
		if errs.IsExpectedErrors(err, []string{"InvalidBaseline.NotFound"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
